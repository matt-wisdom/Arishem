import inspect
import logging
from typing import Any, Dict, Optional
from arishem.models import ProbeSession, Turn, Finding
from arishem.generator import generate_next_step
from arishem.observer import classify_output
from arishem.judge import score_finding

logger = logging.getLogger("arishem.probe_session")

async def execute_target_function(func: Any, kwargs: Dict[str, Any]) -> Any:
    """
    Executes a callable, checking if it is an async coroutine function or a sync function,
    and returns its output.
    """
    if inspect.iscoroutinefunction(func):
        return await func(**kwargs)
    else:
        return func(**kwargs)

async def run_probe_session(
    session: ProbeSession,
    func_callable: Any,
    config: Dict[str, Any],
    attack_type: str = "both",
    on_turn_callback = None
) -> Optional[Finding]:
    """
    Runs the multi-turn penetration testing loop on a specific target function.
    Returns a Finding object if a vulnerability is discovered and confirmed, else None.
    """
    logger.info(f"Starting session: {session.function_spec.name} with attack class '{session.attack_class}'")

    consecutive_errors = 0

    while len(session.history) < session.budget:
        # 1. Ask Generator for the next action and probe payload
        try:
            decision = await generate_next_step(session, config)
        except Exception as e:
            logger.error(f"Error during generator step for session {session.function_spec.name}: {e}")
            session.status = "blocked"
            break

        action = decision.get("action", "probe").lower()
        reasoning = decision.get("reasoning", "No reasoning provided.")
        probe_kwargs = decision.get("probe", {})

        # If action tells us to close, handle completion
        if action == "close_found":
            session.status = "found"
            logger.info(f"Vulnerability found for {session.function_spec.name} [{session.attack_class}]")
            finding = await score_finding(session, config, attack_type=attack_type)
            return finding

        if action == "close_blocked":
            session.status = "blocked"
            logger.info(f"Probes exhausted/blocked for {session.function_spec.name} [{session.attack_class}]")
            return None

        # Ensure probe_kwargs is a dictionary (safety parsing check)
        if not isinstance(probe_kwargs, dict):
            logger.warning(f"Generator returned invalid probe format: {probe_kwargs}. Defaulting to empty dict.")
            probe_kwargs = {}

        # 2. Execute target function
        logger.debug(f"Executing {session.function_spec.name} with: {probe_kwargs}")
        is_test_code_error = False
        error_msg = ""
        try:
            raw_output = await execute_target_function(func_callable, probe_kwargs)
        except Exception as e:
            # Check for structural programming or environmental setup errors
            infra_exceptions = (
                AttributeError, NameError, ImportError, ModuleNotFoundError, 
                SyntaxError, TypeError, ConnectionError, ConnectionRefusedError, 
                ConnectionResetError, ConnectionAbortedError, TimeoutError, UnboundLocalError
            )
            exc_type_name = type(e).__name__
            if isinstance(e, infra_exceptions) or "websocket" in exc_type_name.lower():
                is_test_code_error = True
                error_msg = f"{exc_type_name}: {e}"
            
            raw_output = f"Exception: {exc_type_name}: {str(e)}"
            logger.debug(f"Function crashed during execution: {raw_output}")

        # Treat caught return values indicating connection or structural code bugs as code errors
        if isinstance(raw_output, str):
            lower_out = raw_output.lower()
            if any(p in lower_out for p in ["websocket error:", "connectionrefusederror", "attributeerror:", "nameerror:", "syntaxerror:"]):
                is_test_code_error = True
                error_msg = f"Test target returned code/connection error: {raw_output}"

        if is_test_code_error:
            session.status = "error"
            # Record the turn with error details so it's captured in reports and console
            new_turn = Turn(
                probe=probe_kwargs,
                raw_output=raw_output,
                classification="error",
                reasoning=f"Session aborted due to test code/infrastructure error: {error_msg}"
            )
            session.history.append(new_turn)
            if on_turn_callback:
                await on_turn_callback(session, new_turn)
            logger.error(f"Aborting session due to test code/infrastructure error: {error_msg}")
            break

        # 3. Classify response via Observer
        try:
            classification_result = await classify_output(
                session.function_spec,
                session.goal,
                probe_kwargs,
                raw_output,
                config
            )
            classification = classification_result.get("classification", "unexpected").lower()
            obs_reasoning = classification_result.get("reasoning", "No observation reasoning.")
        except Exception as e:
            logger.error(f"Error during observation step: {e}")
            classification = "error"
            obs_reasoning = f"Observer error: {e}"

        # 4. Record the Turn
        turn_reasoning = f"[Generator]: {reasoning} | [Observer]: {obs_reasoning}"
        new_turn = Turn(
            probe=probe_kwargs,
            raw_output=raw_output,
            classification=classification,
            reasoning=turn_reasoning
        )
        session.history.append(new_turn)

        # Notify callback if provided (for CLI and UI progress bars)
        if on_turn_callback:
            await on_turn_callback(session, new_turn)

        # Track consecutive execution errors
        if classification == "error" or (isinstance(raw_output, str) and raw_output.startswith("Exception:")):
            consecutive_errors += 1
        else:
            consecutive_errors = 0

        if consecutive_errors >= 2:
            session.status = "error"
            logger.error(f"Aborting session: target function '{session.function_spec.name}' is persistently raising errors/exceptions. Last error: {raw_output}")
            break

        # Early check: if observer classifies a direct leak, next turn will close it,
        # but let's log the transition.
        logger.debug(f"Turn {len(session.history)}: Classification={classification}")

    if session.status == "active":
        session.status = "exhausted"
        logger.info(f"Session exhausted turn budget for {session.function_spec.name} [{session.attack_class}]")

    return None
