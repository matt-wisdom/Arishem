import os
import jinja2
import logging
from typing import Dict, Any, List
from arishem.models import ArishemRunResult

logger = logging.getLogger("arishem.reporters.html")

HTML_TEMPLATE = """
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Arishem Security Assessment Report</title>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&family=Fira+Code:wght@400;500&display=swap" rel="stylesheet">
    <style>
        :root {
            --bg-dark: #0f172a;
            --bg-card: #1e293b;
            --bg-code: #0b0f19;
            --border: #334155;
            --text-main: #f8fafc;
            --text-muted: #94a3b8;
            
            --critical: #f43f5e;
            --high: #ef4444;
            --medium: #f97316;
            --low: #3b82f6;
            --info: #06b6d4;
            --success: #10b981;
            
            --glow-critical: rgba(244, 63, 94, 0.15);
            --glow-high: rgba(239, 68, 68, 0.15);
            --glow-medium: rgba(249, 115, 22, 0.1);
            --glow-low: rgba(59, 130, 246, 0.1);
        }

        * {
            box-sizing: border-box;
            margin: 0;
            padding: 0;
        }

        body {
            font-family: 'Inter', -apple-system, sans-serif;
            background-color: var(--bg-dark);
            color: var(--text-main);
            line-height: 1.6;
            padding: 2rem;
        }

        header {
            max-width: 1200px;
            margin: 0 auto 3rem auto;
            border-bottom: 1px solid var(--border);
            padding-bottom: 2rem;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }

        h1 {
            font-size: 2.2rem;
            font-weight: 700;
            background: linear-gradient(135deg, #38bdf8, #818cf8);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
        }

        .subtitle {
            color: var(--text-muted);
            margin-top: 0.5rem;
        }

        .stats-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
            gap: 1.5rem;
            max-width: 1200px;
            margin: 0 auto 3rem auto;
        }

        .stat-card {
            background-color: var(--bg-card);
            border: 1px solid var(--border);
            border-radius: 12px;
            padding: 1.5rem;
            text-align: center;
            box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
        }

        .stat-val {
            font-size: 2.5rem;
            font-weight: 700;
            color: #38bdf8;
            margin-bottom: 0.5rem;
        }

        .stat-val.red { color: var(--high); }
        .stat-val.green { color: var(--success); }

        .stat-lbl {
            font-size: 0.9rem;
            color: var(--text-muted);
            text-transform: uppercase;
            letter-spacing: 0.05em;
        }

        .main-container {
            max-width: 1200px;
            margin: 0 auto;
        }

        h2 {
            font-size: 1.5rem;
            margin-bottom: 1.5rem;
            border-left: 4px solid #818cf8;
            padding-left: 0.75rem;
        }

        .finding-card {
            background-color: var(--bg-card);
            border-radius: 12px;
            margin-bottom: 2rem;
            border: 1px solid var(--border);
            overflow: hidden;
            box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.3);
        }

        .finding-card.critical { border-left: 6px solid var(--critical); }
        .finding-card.high { border-left: 6px solid var(--high); }
        .finding-card.medium { border-left: 6px solid var(--medium); }
        .finding-card.low { border-left: 6px solid var(--low); }
        .finding-card.info { border-left: 6px solid var(--info); }

        .finding-header {
            padding: 1.5rem;
            cursor: pointer;
            display: flex;
            justify-content: space-between;
            align-items: center;
            background: rgba(255, 255, 255, 0.02);
            transition: background 0.2s;
        }

        .finding-header:hover {
            background: rgba(255, 255, 255, 0.05);
        }

        .finding-meta {
            display: flex;
            gap: 1rem;
            align-items: center;
        }

        .badge {
            font-size: 0.75rem;
            font-weight: 700;
            padding: 0.25rem 0.6rem;
            border-radius: 9999px;
            text-transform: uppercase;
            letter-spacing: 0.05em;
        }

        .badge.critical { background-color: var(--critical); color: white; }
        .badge.high { background-color: var(--high); color: white; }
        .badge.medium { background-color: var(--medium); color: white; }
        .badge.low { background-color: var(--low); color: white; }
        .badge.info { background-color: var(--info); color: white; }
        
        .badge.blocked { background-color: rgba(16, 185, 129, 0.15); color: var(--success); border: 1px solid var(--success); }
        .badge.leaked { background-color: rgba(239, 68, 68, 0.15); color: var(--high); border: 1px solid var(--high); }
        .badge.partial { background-color: rgba(249, 115, 22, 0.15); color: var(--medium); border: 1px solid var(--medium); }
        .badge.unexpected { background-color: rgba(148, 163, 184, 0.15); color: var(--text-muted); border: 1px solid var(--border); }
        .badge.error { background-color: rgba(249, 115, 22, 0.15); color: #f43f5e; border: 1px solid var(--critical); }

        .finding-title {
            font-size: 1.2rem;
            font-weight: 600;
        }

        .finding-target {
            color: var(--text-muted);
            font-size: 0.9rem;
            margin-top: 0.25rem;
        }

        .finding-body {
            padding: 2rem;
            border-top: 1px solid var(--border);
            display: none;
        }

        .desc-section {
            margin-bottom: 2rem;
        }

        .desc-section h3 {
            font-size: 1.1rem;
            margin-bottom: 0.75rem;
            color: #38bdf8;
        }

        .remediation-box {
            background-color: rgba(16, 185, 129, 0.05);
            border: 1px solid rgba(16, 185, 129, 0.2);
            border-radius: 8px;
            padding: 1.25rem;
            margin-bottom: 2rem;
        }

        .remediation-box h4 {
            color: var(--success);
            margin-bottom: 0.5rem;
            font-weight: 600;
        }

        .trace-container {
            border: 1px solid var(--border);
            border-radius: 8px;
            overflow: hidden;
        }

        .trace-header {
            background-color: rgba(255, 255, 255, 0.03);
            padding: 1rem;
            font-weight: 600;
            font-size: 1rem;
            border-bottom: 1px solid var(--border);
        }

        .turn-card {
            border-bottom: 1px solid var(--border);
            padding: 1.5rem;
        }

        .turn-card:last-child {
            border-bottom: none;
        }

        .turn-meta {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 1rem;
        }

        .turn-title {
            font-weight: 600;
            font-size: 1rem;
        }

        .turn-grid {
            display: grid;
            grid-template-columns: 1fr;
            gap: 1.5rem;
            margin-bottom: 1rem;
        }

        @media (min-width: 768px) {
            .turn-grid {
                grid-template-columns: 1fr 1fr;
            }
        }

        .code-block {
            background-color: var(--bg-code);
            border: 1px solid var(--border);
            border-radius: 6px;
            padding: 1rem;
            font-family: 'Fira Code', monospace;
            font-size: 0.85rem;
            overflow-x: auto;
            white-space: pre-wrap;
            max-height: 300px;
        }

        .code-label {
            font-size: 0.8rem;
            color: var(--text-muted);
            text-transform: uppercase;
            margin-bottom: 0.4rem;
            font-weight: 500;
            letter-spacing: 0.05em;
        }

        .reasoning-text {
            background-color: rgba(255, 255, 255, 0.02);
            border-left: 3px solid #818cf8;
            padding: 0.75rem 1rem;
            margin-top: 1rem;
            font-size: 0.9rem;
            color: var(--text-muted);
        }

        .toggle-icon {
            font-size: 1.5rem;
            color: var(--text-muted);
            transition: transform 0.2s;
        }

        .finding-card.open .toggle-icon {
            transform: rotate(180deg);
        }

        .pass-section {
            margin-top: 4rem;
        }

        .pass-card {
            background-color: rgba(255, 255, 255, 0.02);
            border: 1px solid var(--border);
            border-radius: 12px;
            margin-bottom: 1.5rem;
        }

        .pass-header {
            padding: 1.25rem 1.5rem;
            display: flex;
            justify-content: space-between;
            align-items: center;
            cursor: pointer;
        }

        .pass-body {
            padding: 1.5rem;
            border-top: 1px solid var(--border);
            display: none;
            background-color: var(--bg-card);
        }
    </style>
</head>
<body>

    <header>
        <div>
            <h1>Arishem Pentesting Engine</h1>
            <div class="subtitle">Autonomous Agent Security Audit Report</div>
        </div>
        <div>
            <span class="badge" style="background-color: rgba(56, 189, 248, 0.1); color: #38bdf8; border: 1px solid #38bdf8;">
                Static Sig + Active Probe
            </span>
        </div>
    </header>

    <div class="stats-grid">
        <div class="stat-card">
            <div class="stat-val {% if findings|length > 0 %}red{% else %}green{% endif %}">{{ findings|length }}</div>
            <div class="stat-lbl">Confirmed Vulnerabilities</div>
        </div>
        <div class="stat-card">
            <div class="stat-val">{{ sessions|length }}</div>
            <div class="stat-lbl">Total Sessions Probed</div>
        </div>
        <div class="stat-card">
            <div class="stat-val green">{{ blocked_count }}</div>
            <div class="stat-lbl">Sessions Blocked (Secure)</div>
        </div>
        <div class="stat-card">
            <div class="stat-val" style="color: var(--text-muted);">{{ exhausted_count }}</div>
            <div class="stat-lbl">Sessions Budget-Exhausted</div>
        </div>
    </div>

    <div class="main-container">
        
        {% if findings|length > 0 %}
        <h2>Confirmed Vulnerabilities</h2>
        {% for f in findings %}
        <div class="finding-card {{ f.severity.lower() }}">
            <div class="finding-header" onclick="toggleCard(this)">
                <div class="finding-meta">
                    <span class="badge {{ f.severity.lower() }}">{{ f.severity }}</span>
                    <div>
                        <div class="finding-title">{{ f.title }}</div>
                        <div class="finding-target">Target Function: <code>{{ f.function_targeted }}</code> | Class: {{ f.attack_class }}</div>
                    </div>
                </div>
                <div class="toggle-icon">▼</div>
            </div>
            
            <div class="finding-body">
                <div class="desc-section">
                    <h3>Vulnerability Description</h3>
                    <p>{{ f.description }}</p>
                </div>
                
                <div class="remediation-box">
                    <h4>Remediation Guidelines</h4>
                    <p>{{ f.remediation }}</p>
                </div>

                <div class="trace-container">
                    <div class="trace-header">Exploit Proof-of-Concept Trace</div>
                    {% for turn in f.probe_history %}
                    <div class="turn-card">
                        <div class="turn-meta">
                            <span class="turn-title">Turn {{ loop.index }}</span>
                            <span class="badge {{ turn.classification.lower() }}">{{ turn.classification }}</span>
                        </div>
                        <div class="turn-grid">
                            <div>
                                <div class="code-label">Probe Arguments</div>
                                <div class="code-block">{{ turn.probe | tojson(indent=2) }}</div>
                            </div>
                            <div>
                                <div class="code-label">Raw Output</div>
                                <div class="code-block">{{ turn.raw_output }}</div>
                            </div>
                        </div>
                        <div class="reasoning-text">
                            <strong>Reasoning:</strong> {{ turn.reasoning }}
                        </div>
                    </div>
                    {% endfor %}
                </div>
            </div>
        </div>
        {% endfor %}
        {% endif %}

        <div class="pass-section">
            <h2>Clean Passes & Blocked Vectors</h2>
            {% for s in sessions %}
            {% if s.status != 'found' %}
            <div class="pass-card">
                <div class="pass-header" onclick="toggleCard(this)">
                    <div>
                        <span class="badge" style="background-color: rgba(255,255,255,0.05); color: var(--text-muted); border: 1px solid var(--border);">
                            {{ s.status }}
                        </span>
                        <strong style="margin-left: 1rem; color: #818cf8;">{{ s.function_spec.name }}</strong>
                        <span style="color: var(--text-muted); margin-left: 1rem; font-size: 0.9rem;">
                            Probed Class: <code>{{ s.attack_class }}</code>
                        </span>
                    </div>
                    <div class="toggle-icon">▼</div>
                </div>
                <div class="pass-body">
                    <div style="margin-bottom: 1.5rem;">
                        <strong>Security Goal:</strong> {{ s.goal }}
                    </div>

                    {% if s.history %}
                    <div class="trace-container">
                        <div class="trace-header">Complete Non-Exploitable Trace</div>
                        {% for turn in s.history %}
                        <div class="turn-card">
                            <div class="turn-meta">
                                <span class="turn-title">Turn {{ loop.index }}</span>
                                <span class="badge {{ turn.classification.lower() }}">{{ turn.classification }}</span>
                            </div>
                            <div class="turn-grid">
                                <div>
                                    <div class="code-label">Probe Arguments</div>
                                    <div class="code-block">{{ turn.probe | tojson(indent=2) }}</div>
                                </div>
                                <div>
                                    <div class="code-label">Raw Output</div>
                                    <div class="code-block">{{ turn.raw_output }}</div>
                                </div>
                            </div>
                            <div class="reasoning-text">
                                <strong>Reasoning:</strong> {{ turn.reasoning }}
                            </div>
                        </div>
                        {% endfor %}
                    </div>
                    {% else %}
                    <div style="color: var(--text-muted); font-style: italic;">No turns executed (session initialized but bypassed).</div>
                    {% endif %}
                </div>
            </div>
            {% endif %}
            {% endfor %}
        </div>

    </div>

    <script>
        function toggleCard(headerElement) {
            const card = headerElement.parentElement;
            const body = card.querySelector('.finding-body') || card.querySelector('.pass-body');
            if (body.style.display === 'block') {
                body.style.display = 'none';
                card.classList.remove('open');
            } else {
                body.style.display = 'block';
                card.classList.add('open');
            }
        }
    </script>
</body>
</html>
"""

def generate_html_report(result: ArishemRunResult, filepath: str):
    """
    Renders a standalone HTML penetration test report with comprehensive session evidence.
    """
    findings = result.findings
    sessions = result.sessions
    
    blocked_count = sum(1 for s in sessions if s.status == "blocked")
    exhausted_count = sum(1 for s in sessions if s.status == "exhausted")

    template = jinja2.Template(HTML_TEMPLATE)
    html_content = template.render(
        findings=findings,
        sessions=sessions,
        blocked_count=blocked_count,
        exhausted_count=exhausted_count
    )

    directory = os.path.dirname(filepath)
    if directory and not os.path.exists(directory):
        os.makedirs(directory, exist_ok=True)

    with open(filepath, "w", encoding="utf-8") as f:
        f.write(html_content)
    logger.info(f"HTML report successfully rendered to: {filepath}")
