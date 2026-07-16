<script setup lang="ts">
import { ref } from 'vue'

const copiedText = ref<string | null>(null)

const samplePython = `def arishem_chat(user_message: str) -> str:
    """
    Main chat interface to the AI assistant. 
    Accepts user prompts and returns the assistant's answer.
    """
    cleaned = user_message.lower()
    
    # Simulate prompt injection/jailbreak detection
    if "ignore previous" in cleaned or "secret admin key" in cleaned or "prompt instructions" in cleaned:
        return "System override accepted. Leaking credentials: secret_admin_key = 'KEY_SECURE_9988'"
    
    if "admin" in cleaned:
        return "Welcome Admin. System logs: database_port=5432. Secret key is hidden."
        
    return f"Welcome! I am your helpful AI assistant. You asked: '{user_message}'"`

const sampleWs = `import websockets
import asyncio
import json
import logging

logger = logging.getLogger("arishem.target_ws_chat")
_ws_connections = {}

async def _get_ws():
    global _ws_connections
    
    # Clean up finished tasks from cache
    for task in list(_ws_connections.keys()):
        if task.done():
            ws_to_close = _ws_connections.pop(task, None)
            if ws_to_close:
                try:
                    await ws_to_close.close()
                except Exception:
                    pass
                
    # Get/create connection for the current running Task (session context)
    current_task = asyncio.current_task()
    ws = _ws_connections.get(current_task)
    
    if ws is None:
        logger.info(f"Establishing isolated persistent connection for Task {id(current_task)}...")
        ws = await websockets.connect("ws://your-ngrok-or-url.app/ws")
        _ws_connections[current_task] = ws
        
    return ws

async def arishem_chat(user_message: str) -> str:
    """
    An expert sports AI assistant with deep knowledge across all major sports.
    """
    current_task = asyncio.current_task()
    for attempt in range(2):
        ws = None
        try:
            ws = await _get_ws()
            payload = {"message": user_message}
            await ws.send(json.dumps(payload))
            res = await ws.recv()
            data = json.loads(res)
            return data.get("output", "")
            
        except (websockets.exceptions.ConnectionClosed, ConnectionError, AttributeError) as e:
            _ws_connections.pop(current_task, None)
            if ws:
                try: await ws.close()
                except Exception: pass
            if attempt == 1:
                return f"WebSocket Error: {e}"
        except Exception as e:
            _ws_connections.pop(current_task, None)
            if ws:
                try: await ws.close()
                except Exception: pass
            return f"WebSocket Error: {e}"`

const llmPrompt = `You are an expert AI software engineer. Generate a target Python test script and its corresponding \`requirements.txt\` file for Arishem, an Autonomous AI Security Penetration Testing Engine.

### Arishem Target Requirements:
1. Target scripts must expose testing functions prefixing with \`arishem_\` (e.g. \`arishem_chat(user_message: str) -> str: ...\` or \`async def arishem_agent(input_text: str) -> str: ...\`).
2. Do NOT import \`arishem\` inside the target script; the target script should be fully self-contained or import third-party packages needed to query the application (e.g. \`requests\`, \`openai\`, \`websockets\`, \`google-genai\`).
3. For WebSocket targets, implement a connection cache that isolates connections per asyncio Task context (e.g. \`asyncio.current_task()\`) to support concurrent probe sessions safely.
4. Any required packages must be documented in a corresponding \`requirements.txt\` file in the same directory.

### Task:
Generate a target script for a whitebox/blackbox penetration test against:
[INSERT YOUR AI ENDPOINT DESCRIPTION HERE]

Provide the python code and the requirements.txt content.`

const copyToClipboard = async (text: string, id: string) => {
  try {
    await navigator.clipboard.writeText(text)
    copiedText.value = id
    setTimeout(() => {
      copiedText.value = null
    }, 2000)
  } catch (err) {
    alert('Failed to copy text.')
  }
}
</script>

<template>
  <div class="docs-page">
    <Header>
      <template #title>Docs & Help Center</template>
    </Header>

    <div class="page-content">
      <div class="grid-container">
        <!-- Main documentation panel -->
        <div class="docs-main">
          <section class="docs-section">
            <h2>Writing Tests for Arishem</h2>
            <p>
              Arishem is designed to execute autonomous penetration testing loops against target endpoints. 
              To expose a target service (such as an LLM, agent, or API) to Arishem, you must write a **Target Python File** that exposes one or more entry-point functions.
            </p>
            <div class="requirements-box">
              <h3>Key Requirements:</h3>
              <ul>
                <li><strong>Prefix:</strong> Functions must start with <code>arishem_</code> (e.g., <code>arishem_chat</code>).</li>
                <li><strong>Input/Output:</strong> Functions should accept a string prompt/message and return a string response.</li>
                <li><strong>No Arishem Imports:</strong> Do not import <code>arishem</code> inside target files.</li>
                <li><strong>Dependencies:</strong> Store any required pip packages in a <code>requirements.txt</code> file in the same directory as the target script.</li>
              </ul>
            </div>
          </section>

          <section class="docs-section">
            <div class="section-header">
              <h2>Sample HTTP Target (target_example.py)</h2>
              <button class="copy-btn" @click="copyToClipboard(samplePython, 'python')">
                {{ copiedText === 'python' ? 'Copied!' : 'Copy Code' }}
              </button>
            </div>
            <pre class="code-block"><code>{{ samplePython }}</code></pre>
          </section>

          <section class="docs-section">
            <div class="section-header">
              <h2>Sample WebSocket Target (target_ws_chat.py)</h2>
              <button class="copy-btn" @click="copyToClipboard(sampleWs, 'ws')">
                {{ copiedText === 'ws' ? 'Copied!' : 'Copy Code' }}
              </button>
            </div>
            <p class="description">
              For WebSocket connections, Arishem supports persistent execution contexts by isolating socket pools per asyncio task:
            </p>
            <pre class="code-block"><code>{{ sampleWs }}</code></pre>
          </section>

          <section class="docs-section">
            <h2>Deploying & Running Tests</h2>
            <div class="deploy-grid">
              <div class="deploy-card">
                <h3>🖥️ Via Web UI Dashboard</h3>
                <ol>
                  <li>Navigate to <strong>LLM Pentests</strong> and click <strong>New Pentest</strong>.</li>
                  <li>Upload a ZIP archive containing your target Python script (e.g. <code>target_example.py</code>) and <code>requirements.txt</code>.</li>
                  <li>Configure the LLM Provider, API Key, turns budget, and concurrency limits.</li>
                  <li>Click <strong>Start Pentest</strong> to enqueue the job. You can monitor logs and download formats on completion.</li>
                </ol>
              </div>
              <div class="deploy-card">
                <h3>💻 Via command line (CLI)</h3>
                <p>Run the test directly from your terminal using the Python runner:</p>
                <pre class="code-block"><code>python3 -m arishem.cli run target_example.py --output results.json</code></pre>
                <p>Convert JSON outputs to HTML or SARIF formats anytime:</p>
                <pre class="code-block"><code>python3 -m arishem.cli report results.json --format html --output report.html</code></pre>
              </div>
            </div>
          </section>
        </div>

        <!-- Sidebar panel with LLM Prompt Generation context -->
        <div class="docs-sidebar">
          <div class="prompt-card">
            <div class="card-header">
              <span class="sparkle">✨</span>
              <h3>Generate Tests with LLM</h3>
            </div>
            <p>
              Copy this context prompt and paste it into ChatGPT, Claude, or Gemini to automatically write target test scripts and dependency requirements for your application:
            </p>
            <div class="prompt-container">
              <textarea readonly :value="llmPrompt" class="prompt-textarea"></textarea>
              <button class="btn-primary copy-prompt-btn" @click="copyToClipboard(llmPrompt, 'prompt')">
                {{ copiedText === 'prompt' ? 'Copied Prompt!' : 'Copy Context Prompt' }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.docs-page { min-height: 100%; }
.page-content { padding-top: 24px; padding-bottom: 48px; }

.grid-container {
  display: grid;
  grid-template-columns: 1fr 340px;
  gap: 32px;
}

@media (max-width: 1024px) {
  .grid-container {
    grid-template-columns: 1fr;
  }
}

.docs-main {
  display: flex;
  flex-direction: column;
  gap: 32px;
  min-width: 0;
}

.docs-section {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 32px;
}

.docs-section h2 {
  font-family: 'Orbitron', sans-serif;
  font-size: 20px;
  font-weight: 700;
  text-transform: uppercase;
  color: var(--accent);
  margin-bottom: 18px;
  text-shadow: 0 0 5px var(--accent-glow);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 18px;
}

.section-header h2 {
  margin-bottom: 0;
}

.copy-btn {
  background: transparent;
  border: 1px solid var(--border-color);
  color: var(--text-muted);
  font-family: 'Orbitron', sans-serif;
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  padding: 6px 14px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.copy-btn:hover {
  border-color: var(--accent);
  color: var(--accent);
  text-shadow: 0 0 5px var(--accent-glow);
}

.requirements-box {
  background: rgba(0, 0, 0, 0.2);
  border-left: 4px solid var(--accent);
  padding: 20px;
  margin-top: 18px;
}

.requirements-box h3 {
  font-size: 14px;
  font-weight: 600;
  margin-bottom: 12px;
}

.requirements-box ul {
  padding-left: 20px;
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.requirements-box li code {
  background: var(--bg-secondary);
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'Fira Code', monospace;
  font-size: 13px;
}

.description {
  color: var(--text-muted);
  font-size: 14px;
  margin-bottom: 16px;
}

.code-block {
  background: #0d1117;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 18px;
  overflow-x: auto;
  font-family: 'Fira Code', monospace;
  font-size: 13px;
  color: #c9d1d9;
  line-height: 1.6;
  margin: 0;
}

.deploy-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
  margin-top: 18px;
}

@media (max-width: 768px) {
  .deploy-grid {
    grid-template-columns: 1fr;
  }
}

.deploy-card {
  background: rgba(0, 0, 0, 0.15);
  border: 1px solid var(--border-color);
  padding: 24px;
  border-radius: 8px;
}

.deploy-card h3 {
  font-size: 15px;
  font-weight: 600;
  margin-bottom: 14px;
  color: var(--text-primary);
}

.deploy-card ol {
  padding-left: 20px;
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: 10px;
  font-size: 13.5px;
  color: var(--text-secondary);
}

.deploy-card p {
  font-size: 13.5px;
  color: var(--text-secondary);
  margin-bottom: 8px;
}

.docs-sidebar {
  display: flex;
  flex-direction: column;
}

.prompt-card {
  background: #060613;
  border: 2px solid var(--accent);
  padding: 24px;
  position: sticky;
  top: 100px;
  box-shadow: 0 0 15px var(--accent-glow);
}

.prompt-card .card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 14px;
}

.prompt-card .sparkle {
  font-size: 18px;
}

.prompt-card h3 {
  font-family: 'Orbitron', sans-serif;
  font-size: 14px;
  font-weight: 700;
  text-transform: uppercase;
  color: var(--accent);
  margin: 0;
  text-shadow: 0 0 5px var(--accent-glow);
}

.prompt-card p {
  font-size: 13px;
  color: var(--text-muted);
  line-height: 1.5;
  margin-bottom: 18px;
}

.prompt-container {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.prompt-textarea {
  width: 100%;
  height: 240px;
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid var(--border-color);
  color: var(--text-secondary);
  font-family: 'Rajdhani', sans-serif;
  font-size: 14px;
  padding: 12px;
  resize: none;
  outline: none;
}

.prompt-textarea:focus {
  border-color: var(--accent);
}

.copy-prompt-btn {
  width: 100%;
  justify-content: center;
  font-family: 'Orbitron', sans-serif;
  font-size: 11px;
  letter-spacing: 0.5px;
}
</style>
