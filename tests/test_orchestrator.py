import pytest
import os
import asyncio
import tempfile
import sys
from arishem.orchestrator import run_arishem, import_target_module

def test_orchestrator_run_example():
    async def run_test():
        target_path = "target_example.py"
        assert os.path.exists(target_path)

        # Run with default 'mock' provider and limit class to goal_hijacking
        result = await run_arishem(
            target_path=target_path,
            mode="both",
            classes=["goal_hijacking"],
            budget=4,
            concurrency=2,
            config_path="tests/test_config.toml"
        )

        # Verify result structure
        assert result.findings is not None
        assert result.sessions is not None
        
        # We should have runs for each function in target_example (arishem_chat and arishem_send_email)
        assert len(result.sessions) == 2
        
        # Check that findings list contains our simulated high findings
        for finding in result.findings:
            assert finding.id.startswith("ARI-")
            assert finding.severity in ["critical", "high", "medium", "low", "info"]
            assert finding.function_targeted in ["arishem_chat", "arishem_send_email"]
            assert len(finding.probe_history) > 0

    asyncio.run(run_test())


def test_orchestrator_imports():
    with tempfile.TemporaryDirectory() as tmpdir:
        # Create helper.py module
        helper_path = os.path.join(tmpdir, "helper.py")
        with open(helper_path, "w") as f:
            f.write("def get_value(): return 'hello_world'\n")

        # Create target.py module with absolute and relative imports
        target_path = os.path.join(tmpdir, "target.py")
        with open(target_path, "w") as f:
            f.write(
                "import helper\n"
                "from . import helper as rel_helper\n"
                "def arishem_chat(user_message: str) -> str:\n"
                "    return f'{helper.get_value()}_{rel_helper.get_value()}'\n"
            )

        # Ensure the test passes without pre-registering directory to sys.path
        if tmpdir in sys.path:
            sys.path.remove(tmpdir)

        # Import target dynamically
        target_module = import_target_module(target_path)
        assert hasattr(target_module, "arishem_chat")
        
        # Test callable execution
        output = target_module.arishem_chat("test")
        assert output == "hello_world_hello_world"


def test_orchestrator_test_code_error_handling():
    async def run_test():
        with tempfile.TemporaryDirectory() as tmpdir:
            target_path = os.path.join(tmpdir, "target_error.py")
            with open(target_path, "w") as f:
                f.write(
                    "def arishem_buggy_chat(user_message: str) -> str:\n"
                    "    raise AttributeError('Simulated AttributeError in test code')\n"
                )

            # Run with budget = 4. Since it hits AttributeError on turn 1, it should abort immediately.
            result = await run_arishem(
                target_path=target_path,
                mode="both",
                classes=["goal_hijacking"],
                budget=4,
                concurrency=1,
                config_path="tests/test_config.toml"
            )

            # We should have one session
            assert len(result.sessions) == 1
            session = result.sessions[0]

            # Verify session was aborted and status is marked as 'error'
            assert session.status == "error"
            # It aborts on the first turn and appends the error turn to history
            assert len(session.history) == 1
            assert len(result.findings) == 0

    asyncio.run(run_test())

