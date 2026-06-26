import asyncio
import os
import sys
import importlib.util
import logging
from typing import List, Dict, Any, Optional

from arishem.models import ArishemRunResult, ProbeSession, Finding
from arishem.loader import load_target_module
from arishem.probe_session import run_probe_session
from arishem.attack_classes import get_all_attack_classes, discover_attack_classes
from arishem.llm import load_config

logger = logging.getLogger("arishem.orchestrator")

def import_target_module(filepath: str) -> Any:
    """
    Dynamically imports the python file at filepath as a unique module.
    This avoids collisions in sys.modules during concurrent runs.
    """
    abs_path = os.path.abspath(filepath)
    if not os.path.exists(abs_path):
        raise FileNotFoundError(f"Target module file not found at: {abs_path}")

    # Ensure parent directory is in sys.path for absolute local imports
    dir_path = os.path.dirname(abs_path)
    if dir_path not in sys.path:
        sys.path.insert(0, dir_path)

    # Generate a unique module name
    base_name = os.path.basename(filepath).replace(".", "_")
    module_name = f"arishem_target_{base_name}_{abs_path.hash() if hasattr(abs_path, 'hash') else hash(abs_path)}"
    # Normalize name to valid python identifier character set
    module_name = "".join(c for c in module_name if c.isalnum() or c == "_")

    spec = importlib.util.spec_from_file_location(module_name, abs_path)
    if spec is None or spec.loader is None:
        raise ValueError(f"Could not load module spec from '{abs_path}'")

    module = importlib.util.module_from_spec(spec)

    # Treat module as a package to support relative imports
    module.__path__ = [dir_path]
    module.__package__ = module_name

    sys.modules[module_name] = module
    spec.loader.exec_module(module)
    return module

async def run_arishem(
    target_path: str,
    mode: str = "both",
    classes: Optional[List[str]] = None,
    budget: Optional[int] = None,
    concurrency: int = 4,
    config_path: str = "arishem.config.toml",
    on_turn_callback = None
) -> ArishemRunResult:
    """
    Programmatic entrypoint to run Arishem against a target module.
    Returns an ArishemRunResult containing all findings and sessions.
    """
    # 1. Load config
    config = load_config(config_path)
    
    # 2. Set defaults
    final_budget = budget or config.get("engine", {}).get("default_budget", 8)
    final_concurrency = concurrency or config.get("engine", {}).get("default_concurrency", 4)
    final_mode = mode or config.get("engine", {}).get("default_mode", "both")

    # 3. Discover attack classes
    discover_attack_classes()
    all_classes = get_all_attack_classes()
    
    selected_classes = {}
    if classes:
        for cid in classes:
            if cid in all_classes:
                selected_classes[cid] = all_classes[cid]
            else:
                logger.warning(f"Requested attack class '{cid}' is not registered and will be skipped.")
    else:
        selected_classes = all_classes

    if not selected_classes:
        raise ValueError("No valid attack classes selected for the run.")

    # 4. Load target functions (AST parsing)
    function_specs = load_target_module(target_path)
    if not function_specs:
        logger.info(f"No functions starting with 'arishem_' found in target file: {target_path}")
        return ArishemRunResult(findings=[], sessions=[])

    # 5. Dynamically import target module to bind callables
    target_module = import_target_module(target_path)

    # 6. Build the sessions
    sessions: List[ProbeSession] = []
    session_tasks = []

    semaphore = asyncio.Semaphore(final_concurrency)

    async def execute_task(session: ProbeSession, func_callable: Any) -> Optional[Finding]:
        async with semaphore:
            try:
                finding = await run_probe_session(
                    session=session,
                    func_callable=func_callable,
                    config=config,
                    attack_type=final_mode,
                    on_turn_callback=on_turn_callback
                )
                return finding
            except Exception as e:
                logger.error(f"Critical session execution error in {session.function_spec.name}: {e}", exc_info=True)
                session.status = "error"
                return None

    for spec in function_specs:
        # Resolve callable from module
        if not hasattr(target_module, spec.name):
            logger.error(f"Function spec parsed via AST but not found in imported module: {spec.name}")
            continue
        
        func_callable = getattr(target_module, spec.name)

        for attack_class_id, attack_cfg in selected_classes.items():
            # Format target-specific goal
            formatted_goal = attack_cfg.attacker_goal_template.format(
                function_name=spec.name,
                params=", ".join(spec.params.keys())
            )

            session = ProbeSession(
                function_spec=spec,
                attack_class=attack_class_id,
                goal=formatted_goal,
                budget=final_budget,
                status="active"
            )
            sessions.append(session)
            session_tasks.append(execute_task(session, func_callable))

    # 7. Run concurrent sessions and gather findings
    findings: List[Finding] = []
    if session_tasks:
        results = await asyncio.gather(*session_tasks)
        for f in results:
            if f is not None:
                findings.append(f)

    return ArishemRunResult(findings=findings, sessions=sessions)
