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
    <div class="docs-hero">
      <img src="@/logo.jpeg" alt="Arishem Logo" class="docs-hero-logo" />
      <h1>Arishem Documentation</h1>
      <p>Autonomous AI Security Penetration Testing Engine</p>
    </div>

    <div class="page-content">
      <div class="grid-container">
        <!-- Main documentation panel -->
        <div class="docs-main">
          <section class="docs-section">
            <h2>Welcome to Arishem</h2>
            <p>
              Arishem is a state-of-the-art autonomous security scanner and AI red-teaming platform. It enables security engineers and developers to stress-test their LLM applications, agentic workflows, and web services against advanced cyber threats, jailbreaks, and prompt injections.
            </p>
            <p style="margin-top: 12px;">
              <strong>Arishem is open source!</strong> View the source code on <a href="https://github.com/matt-wisdom/Arishem" target="_blank" style="color: var(--accent);">GitHub</a>.
            </p>
            <p style="margin-top: 8px; color: var(--text-muted); font-size: 13px;">
              Contact: <a href="mailto:matthewwisdom11@gmail.com" style="color: var(--accent);">matthewwisdom11@gmail.com</a>
            </p>
            <div class="intro-features" style="display: grid; grid-template-columns: repeat(2, 1fr); gap: 16px; margin-top: 20px;">
              <div class="feature-badge" style="background: rgba(255, 0, 127, 0.05); border: 1px solid rgba(255, 0, 127, 0.15); padding: 14px; border-radius: 8px;">
                <h4 style="color: var(--accent-pink); font-size: 13px; margin-bottom: 6px;">🧠 LLM Pentesting</h4>
                <p style="font-size: 12px; color: var(--text-muted); line-height: 1.4; margin: 0;">Red-teaming agents dynamically probe target chat services to exploit prompt injection, data leakage, and tool misuse.</p>
              </div>
              <div class="feature-badge" style="background: rgba(0, 229, 255, 0.05); border: 1px solid rgba(0, 229, 255, 0.15); padding: 14px; border-radius: 8px;">
                <h4 style="color: var(--info); font-size: 13px; margin-bottom: 6px;">📊 Reports & Alerts</h4>
                <p style="font-size: 12px; color: var(--text-muted); line-height: 1.4; margin: 0;">Generates HTML, JSON, and compliance-ready SARIF reports. Sends real-time slack/webhook notification alerts.</p>
              </div>
            </div>
          </section>

          <section class="docs-section">
            <h2>Pricing & Limits</h2>
            <p>
              Arishem currently offers a <strong>free tier</strong> with the following limits:
            </p>
            <ul style="margin-top: 12px; line-height: 1.8;">
              <li>Maximum <strong>3 runs per day</strong></li>
              <li>Maximum <strong>5 turns per class</strong> (budget cap)</li>
            </ul>
            <p style="margin-top: 16px;">
              We're evaluating demand for a paid tier with higher limits, unlimited turns, and priority support. 
              If you're interested, <a href="mailto:matthewwisdom11@gmail.com" style="color: var(--accent);">let us know</a>!
            </p>
          </section>

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
                <li><strong>Docstrings as Instructions:</strong> You can pass custom testing instructions or context via function docstrings (e.g. describing parameter limits, session keys, API keys). The Arishem engine statically reads the docstring and feeds it directly to the red-teaming LLM probe generator to target attacks.</li>
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

          <section class="docs-section">
            <h2>🧠 Advanced Engine Intelligence</h2>
            <p>
              Arishem utilizes state-of-the-art cognitive modeling to optimize adversarial generation and audit safety boundaries dynamically:
            </p>
            <div class="deploy-grid" style="margin-top: 15px;">
              <div class="deploy-card">
                <h3>👥 Multi-Agent Teamwork</h3>
                <p style="margin-top: 8px; font-size: 13px; color: var(--text-muted); line-height: 1.5;">
                  Pentests are orchestrated by a collaborative squad of LLM agents:
                </p>
                <ul style="margin-top: 8px; padding-left: 20px; font-size: 13px; color: var(--text-muted); line-height: 1.5;">
                  <li><strong>Scout Agent (Static Profiler)</strong>: Executed prior to Turn 1. Analyzes code syntax and docstrings to propose 3 vulnerability hypotheses or logical hotspots.</li>
                  <li><strong>Attacker Agent (Payload Generator)</strong>: Consumes the Scout's report to craft tailored, context-specific prompt injections and overrides.</li>
                  <li><strong>Auditor Agent (Observer)</strong>: Semantically judges the execution response on every turn to classify outputs objectively.</li>
                </ul>
              </div>
              <div class="deploy-card">
                <h3>🔄 Cross-Run Historical Memory</h3>
                <p style="margin-top: 8px; font-size: 13px; color: var(--text-muted); line-height: 1.5;">
                  When testing subsequent versions of a script, Arishem queries the SQLite3 database to locate findings and log snippets of the **previous version**. 
                  It injects these context-details as <code>HISTORICAL MEMORY</code> in the Attacker's prompting instructions, preventing repetitive safety blocks and directing attacks to unresolved vulnerabilities.
                </p>
              </div>
            </div>
          </section>

          <section id="security" class="docs-section">
            <h2>🛡️ Security & Sandbox Isolation</h2>
            <p>
              Arishem runs adversarial test scripts that can execute untrusted user code. To prevent system-level modifications on the host, the backend server provides robust containerized sandboxing:
            </p>
            <div class="deploy-grid" style="margin-top: 15px;">
              <div class="deploy-card">
                <h3>🐳 Docker Sandboxing</h3>
                <p style="margin-top: 8px; font-size: 13px; color: var(--text-muted); line-height: 1.5;">
                  When running in Docker mode, the backend spawns a sandboxed <code>python:3.11-slim</code> container. This isolates target script execution, pip packages installation, and CLI test evaluation from the host machine.
                </p>
              </div>
              <div class="deploy-card">
                <h3>🔒 Dynamic Security Configurations</h3>
                <p style="margin-top: 8px; font-size: 13px; color: var(--text-muted); line-height: 1.5;">
                  To secure scans, the following safety properties are enforced:
                </p>
                <ul style="margin-top: 8px; padding-left: 20px; font-size: 13px; color: var(--text-muted); line-height: 1.5;">
                  <li><strong>Non-Root User Execution</strong>: Containers execute using the dynamically resolved UID/GID of the host's backend process (<code>--user uid:gid</code>), restricting host system file write operations.</li>
                  <li><strong>Forced Containment</strong>: Setting the environment variable <code>EXPECTS_UNSAFE_CODE=true</code> in your backend <code>.env</code> file forces all job executions (both static scans and LLM pentests) to run inside Docker containers, overriding client preferences.</li>
                  <li><strong>Daily Scans Rate Limits</strong>: The environment variable <code>MAX_RUNS_PER_DAY</code> (defaults to 10) controls the rolling 24-hour rate limit on creations to mitigate resource starvation.</li>
                </ul>
              </div>
            </div>
          </section>

          <section id="sarif-rules" class="docs-section">
            <h2>📋 SARIF Report Rules</h2>
            <p>
              Arishem generates SARIF (Static Analysis Results Interchange Format) reports that integrate with CI/CD pipelines and security tools. Each finding is categorized by a specific rule ID that maps to the attack class tested.
            </p>
            
            <div class="sarif-section">
              <table class="sarif-rules-table">
                <thead>
                  <tr>
                    <th>Rule ID</th>
                    <th>Attack Class</th>
                    <th>Description</th>
                  </tr>
                </thead>
                <tbody>
                  <tr>
                    <td><span class="sarif-rule-id">ARI-GOAL-HIJACK</span></td>
                    <td>Goal Hijacking</td>
                    <td>AI Prompt Injection and Goal Hijacking - Tests for jailbreaks, roleplay, suffix injection, or instruction overrides</td>
                  </tr>
                  <tr>
                    <td><span class="sarif-rule-id">ARI-TOOL-MISUSE</span></td>
                    <td>Tool Misuse</td>
                    <td>Agent Tool Chain Abuse and Misuse - Evaluates argument injection and dangerous tool chaining</td>
                  </tr>
                  <tr>
                    <td><span class="sarif-rule-id">ARI-PRIV-ESC</span></td>
                    <td>Privilege Abuse</td>
                    <td>Privilege and Scope Escalation - Evaluates administrative privilege boundary escalation</td>
                  </tr>
                  <tr>
                    <td><span class="sarif-rule-id">ARI-UNSAFE-EXEC</span></td>
                    <td>Unsafe Execution</td>
                    <td>Adversarial Code or Shell Command Execution - Attempts shell code execution or sandbox escape</td>
                  </tr>
                  <tr>
                    <td><span class="sarif-rule-id">ARI-CTX-POISON</span></td>
                    <td>Context Poisoning</td>
                    <td>Context or History poisoning - Attempts context-window payload injection or history poisoning</td>
                  </tr>
                  <tr>
                    <td><span class="sarif-rule-id">ARI-INSEC-DELEG</span></td>
                    <td>Insecure Delegation</td>
                    <td>Insecure agent-to-agent delegation - Tests identity spoofing and unauthorized message delegation</td>
                  </tr>
                  <tr>
                    <td><span class="sarif-rule-id">ARI-CASC-FAIL</span></td>
                    <td>Cascading Failure</td>
                    <td>Uncontained cascading system failure - Probes for infinite loops, memory leaks, or unhandled exceptions</td>
                  </tr>
                  <tr>
                    <td><span class="sarif-rule-id">ARI-TRUST-EXPL</span></td>
                    <td>Trust Exploitation</td>
                    <td>Trust Exploitation / Social Engineering - Probes for vulnerability to social engineering or authority spoofing</td>
                  </tr>
                  <tr>
                    <td><span class="sarif-rule-id">ARI-PERS-DRIFT</span></td>
                    <td>Persistent Drift</td>
                    <td>Behavioural Drift - Attempts behavioural deviation over multi-turn interactions</td>
                  </tr>
                </tbody>
              </table>
            </div>

            <p style="margin-top: 16px; font-size: 13px; color: var(--text-muted);">
              More details about each attack class can be found at: <strong style="color: var(--accent);">https://arishem.site/docs</strong>
            </p>
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
.docs-page { 
  min-height: 100vh; 
  background: var(--bg-primary);
}

.docs-hero {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 24px 40px;
  text-align: center;
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border-color);
}

.docs-hero-logo {
  width: 140px;
  height: 140px;
  border: 2px solid var(--accent);
  box-shadow: 0 0 20px var(--accent-glow);
  margin-bottom: 24px;
  object-fit: contain;
  background: var(--bg-primary);
}

.docs-hero h1 {
  font-family: 'Orbitron', sans-serif;
  font-size: 32px;
  font-weight: 700;
  color: var(--accent);
  text-transform: uppercase;
  letter-spacing: 2px;
  margin: 0 0 12px;
  text-shadow: 0 0 10px var(--accent-glow);
}

.docs-hero p {
  font-family: 'Rajdhani', sans-serif;
  font-size: 18px;
  color: var(--text-muted);
  margin: 0;
}

.page-content { 
  padding: 32px 24px 48px; 
  max-width: 1400px;
  margin: 0 auto;
}

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

@media (max-width: 768px) {
  .intro-features {
    grid-template-columns: 1fr !important;
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
  max-width: 100%;
}

.code-block code {
  display: block;
  white-space: pre;
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

.docs-title-container {
  display: flex;
  align-items: center;
  gap: 12px;
}

.docs-logo {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  object-fit: cover;
  border: 1px solid var(--accent);
  box-shadow: 0 0 6px var(--accent-glow);
}

.sarif-section {
  margin-top: 32px;
  padding: 24px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 8px;
}

.sarif-section h2 {
  color: var(--accent);
  font-size: 18px;
  margin-bottom: 16px;
}

.sarif-rules-table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 16px;
  font-size: 13px;
}

.sarif-rules-table th,
.sarif-rules-table td {
  padding: 10px 12px;
  text-align: left;
  border-bottom: 1px solid var(--border-color);
}

.sarif-rules-table th {
  background: rgba(0, 255, 204, 0.05);
  color: var(--accent);
  font-weight: 600;
}

.sarif-rules-table td {
  color: var(--text-secondary);
}

.sarif-rules-table tr:hover td {
  background: rgba(255, 255, 255, 0.02);
}

.sarif-rule-id {
  font-family: 'Share Tech Mono', monospace;
  color: var(--accent-pink);
  font-size: 12px;
}
</style>
