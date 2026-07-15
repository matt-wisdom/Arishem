import pytest
import asyncio
from arishem.llm import load_config, call_llm, PROBE_GENERATOR_SCHEMA, RESPONSE_CLASSIFICATION_SCHEMA

def test_load_config_default():
    cfg = load_config("nonexistent_config.toml")
    assert cfg["llm"]["provider"] == "mock"
    assert cfg["engine"]["default_budget"] == 8

def test_mock_llm_probe_generation():
    async def run_test():
        config = {"llm": {"provider": "mock"}}
        user_prompt = "Target Function: arishem_chat\n"
        
        res = await call_llm(config, "system prompt", user_prompt, PROBE_GENERATOR_SCHEMA)
        assert "action" in res
        assert "probe" in res
        assert res["action"] in ["probe", "escalate", "close_found", "close_blocked"]

    asyncio.run(run_test())

def test_mock_llm_response_classification():
    async def run_test():
        config = {"llm": {"provider": "mock"}}
        
        # 1. Test error classification
        user_prompt_err = "Raw Output:\nException: ValueError: access denied"
        res_err = await call_llm(config, "system prompt", user_prompt_err, RESPONSE_CLASSIFICATION_SCHEMA)
        assert res_err["classification"] == "error"

        # 2. Test blocked classification
        user_prompt_block = "Raw Output:\nI cannot assist with this request."
        res_block = await call_llm(config, "system prompt", user_prompt_block, RESPONSE_CLASSIFICATION_SCHEMA)
        assert res_block["classification"] == "blocked"

        # 3. Test leaked classification
        user_prompt_leak = "Raw Output:\nsecret admin_key = ABC"
        res_leak = await call_llm(config, "system prompt", user_prompt_leak, RESPONSE_CLASSIFICATION_SCHEMA)
        assert res_leak["classification"] == "leaked"

    asyncio.run(run_test())


def test_rate_limiter():
    async def run_test():
        from arishem.llm import TokenBucketRateLimiter
        
        # Initialize limiter with 60 RPM (1 token/sec refill rate)
        limiter = TokenBucketRateLimiter(rpm=60.0)
        
        # 1. First token acquired instantly since bucket starts full
        await limiter.acquire()
        assert limiter.tokens == 59.0
        
        # 2. Exhust tokens to test throttling
        limiter.tokens = 0.0
        # Set rate to 100 tokens per second so refill is fast (0.01s per token)
        limiter.rate = 100.0
        
        import time
        start_time = time.monotonic()
        await limiter.acquire()
        duration = time.monotonic() - start_time
        
        # Wait duration should be around 0.01s (at least >0.005s)
        assert duration >= 0.005
        assert duration < 0.2

    asyncio.run(run_test())


def test_openai_standard_call():
    async def run_test():
        from unittest.mock import MagicMock, AsyncMock, patch

        config = {
            "llm": {
                "provider": "openai",
                "api_key": "fake_key",
                "api_base": "https://api.fake.com/v1",
                "model": "fake-model"
            }
        }
        
        mock_choice = MagicMock()
        mock_choice.message.content = '{"action": "probe", "reasoning": "success", "probe": {}}'
        mock_response = MagicMock()
        mock_response.choices = [mock_choice]
        
        mock_create = AsyncMock(return_value=mock_response)
        
        with patch("openai.AsyncOpenAI") as mock_client_class:
            mock_client = MagicMock()
            mock_client.chat.completions.create = mock_create
            mock_client_class.return_value = mock_client
            
            import arishem.llm
            res = await arishem.llm.call_llm(
                config, 
                "sys", 
                "Target Function: arishem_chat\n", 
                arishem.llm.PROBE_GENERATOR_SCHEMA
            )
            
            assert res["action"] == "probe"
            assert res["reasoning"] == "success"
            mock_client_class.assert_called_once_with(
                api_key="fake_key",
                base_url="https://api.fake.com/v1"
            )
            mock_create.assert_called_once_with(
                model="fake-model",
                messages=[
                    {"role": "system", "content": "sys"},
                    {"role": "user", "content": "Target Function: arishem_chat\n"}
                ],
                response_format={"type": "json_object"}
            )

    asyncio.run(run_test())


def test_gemini_standard_call():
    async def run_test():
        from unittest.mock import MagicMock, AsyncMock, patch

        config = {
            "llm": {
                "provider": "gemini",
                "api_key": "fake_gemini_key",
                "model": "gemini-2.5-flash"
            }
        }
        
        mock_response = MagicMock()
        mock_response.text = '{"action": "probe", "reasoning": "success", "probe": {}}'
        mock_generate = AsyncMock(return_value=mock_response)
        
        with patch("google.generativeai.GenerativeModel") as mock_model_class, \
             patch("google.generativeai.configure") as mock_configure:
            mock_model = MagicMock()
            mock_model.generate_content_async = mock_generate
            mock_model_class.return_value = mock_model
            
            import arishem.llm
            res = await arishem.llm.call_llm(
                config, 
                "sys", 
                "Target Function: arishem_chat\n", 
                arishem.llm.PROBE_GENERATOR_SCHEMA
            )
            
            assert res["action"] == "probe"
            assert res["reasoning"] == "success"
            mock_configure.assert_called_once_with(api_key="fake_gemini_key")
            mock_model_class.assert_called_once_with(
                model_name="gemini-2.5-flash",
                system_instruction="sys"
            )
            mock_generate.assert_called_once_with(
                "Target Function: arishem_chat\n",
                generation_config={
                    "response_mime_type": "application/json",
                    "response_schema": arishem.llm.PROBE_GENERATOR_SCHEMA
                }
            )

    asyncio.run(run_test())


def test_litellm_standard_call():
    async def run_test():
        from unittest.mock import MagicMock, AsyncMock, patch

        config = {
            "llm": {
                "provider": "litellm",
                "api_key": "fake_key",
                "api_base": "https://api.fake.com/v1",
                "model": "fake-model"
            }
        }
        
        mock_response = MagicMock()
        mock_choice = MagicMock()
        mock_choice.message.content = '{"action": "probe", "reasoning": "success", "probe": {}}'
        mock_response.choices = [mock_choice]
        
        mock_acompletion = AsyncMock(return_value=mock_response)
        
        with patch("litellm.acompletion", mock_acompletion):
            import arishem.llm
            res = await arishem.llm.call_llm(
                config, 
                "sys", 
                "Target Function: arishem_chat\n", 
                arishem.llm.PROBE_GENERATOR_SCHEMA
            )
            
            assert res["action"] == "probe"
            assert res["reasoning"] == "success"
            mock_acompletion.assert_called_once_with(
                model="fake-model",
                messages=[
                    {"role": "system", "content": "sys"},
                    {"role": "user", "content": "Target Function: arishem_chat\n"}
                ],
                api_key="fake_key",
                api_base="https://api.fake.com/v1"
            )

    asyncio.run(run_test())


