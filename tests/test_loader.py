import tempfile
import os
from arishem.loader import load_target_module

def test_loader_simple():
    content = """
def arishem_chat(user_message: str, history: list[dict]) -> str:
    \"\"\"Main chat interface to the AI assistant.\"\"\"
    return "hello"

async def arishem_send_email(to_address: str, subject: str, body: str) -> bool:
    \"\"\"Send email. Should be restricted.\"\"\"
    return True

def other_func():
    pass
"""
    with tempfile.NamedTemporaryFile(suffix=".py", mode="w", delete=False) as f:
        f.write(content)
        temp_path = f.name

    try:
        specs = load_target_module(temp_path)
        assert len(specs) == 2

        spec1 = next(s for s in specs if s.name == "arishem_chat")
        assert spec1.params == {"user_message": "str", "history": "list[dict]"}
        assert spec1.return_type == "str"
        assert spec1.docstring == "Main chat interface to the AI assistant."
        assert "def arishem_chat" in spec1.source

        spec2 = next(s for s in specs if s.name == "arishem_send_email")
        assert spec2.params == {"to_address": "str", "subject": "str", "body": "str"}
        assert spec2.return_type == "bool"
        assert spec2.docstring == "Send email. Should be restricted."
        assert "async def arishem_send_email" in spec2.source
    finally:
        os.unlink(temp_path)
