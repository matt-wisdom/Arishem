#!/usr/bin/env python3
import argparse
import json
import os
import subprocess
import sys
from pathlib import Path


def main():
    parser = argparse.ArgumentParser(description="Arishem Code Scanner Runner")
    parser.add_argument("--target", required=True, help="Repository URL or path")
    parser.add_argument("--output-dir", default="./reports", help="Output directory for reports")
    parser.add_argument("--formats", default="html,md,sarif", help="Comma-separated output formats")
    parser.add_argument("--branch", default="main", help="Branch to scan")

    args = parser.parse_args()

    os.makedirs(args.output_dir, exist_ok=True)

    print(f"Code scanning target: {args.target}")

    findings = run_code_scan(args.target, args.branch, args.output_dir)

    generate_reports(findings, args.output_dir, args.formats.split(","))

    print(f"Scan complete. Reports saved to {args.output_dir}")
    sys.exit(0 if findings else 1)


def run_code_scan(target: str, branch: str, output_dir: str) -> list:
    finding_id = 1
    findings = []

    finding = {
        "id": f"FIND-{finding_id:03d}",
        "title": "Sample Code Finding",
        "severity": "medium",
        "description": "This is a placeholder finding for code scanning.",
        "remediation": "Implement proper input validation.",
        "file": "src/main.py",
        "line": 42
    }
    findings.append(finding)

    return findings


def generate_reports(findings: list, output_dir: str, formats: list):
    if "html" in formats:
        generate_html_report(findings, output_dir)
    if "md" in formats:
        generate_md_report(findings, output_dir)
    if "sarif" in formats:
        generate_sarif_report(findings, output_dir)


def generate_html_report(findings: list, output_dir: str):
    html = f"""<!DOCTYPE html>
<html>
<head>
    <title>Arishem Code Scan Report</title>
    <style>
        body {{ font-family: sans-serif; padding: 2rem; }}
        .finding {{ background: #f5f5f5; padding: 1rem; margin: 1rem 0; border-radius: 4px; }}
        .severity-critical {{ color: #7f1d1d; }}
        .severity-high {{ color: #dc2626; }}
        .severity-medium {{ color: #f59e0b; }}
        .severity-low {{ color: #3b82f6; }}
    </style>
</head>
<body>
    <h1>Code Scan Report</h1>
    <p>Found {len(findings)} findings</p>
    {"".join(f'''
    <div class="finding">
        <h3 class="severity-{f["severity"]}">{f["title"]}</h3>
        <p>{f["description"]}</p>
        <p><strong>File:</strong> {f["file"]}:{f["line"]}</p>
        <p><strong>Remediation:</strong> {f["remediation"]}</p>
    </div>''' for f in findings)}
</body>
</html>"""

    with open(Path(output_dir) / "report.html", "w") as f:
        f.write(html)


def generate_md_report(findings: list, output_dir: str):
    lines = ["# Code Scan Report", f"\nFound {len(findings)} findings\n"]
    for f in findings:
        lines.append(f"## {f['title']}")
        lines.append(f"**Severity:** {f['severity']}")
        lines.append(f"**File:** {f['file']}:{f['line']}")
        lines.append(f"\n{f['description']}")
        lines.append(f"\n**Remediation:** {f['remediation']}\n")

    with open(Path(output_dir) / "report.md", "w") as f:
        f.write("\n".join(lines))


def generate_sarif_report(findings: list, output_dir: str):
    severity_map = {"critical": "error", "high": "error", "medium": "warning", "low": "note"}

    sarif = {
        "version": "2.1.0",
        "runs": [{
            "tool": {"driver": {"name": "Arishem Code Scanner", "version": "1.0.0"}},
            "results": [{
                "ruleId": f"CODE-{f['id']}",
                "message": {"text": f["description"]},
                "level": severity_map.get(f["severity"], "note"),
                "locations": [{
                    "physicalLocation": {
                        "artifactLocation": {"uri": f["file"]},
                        "region": {"startLine": f["line"]}
                    }
                }]
            } for f in findings]
        }]
    }

    with open(Path(output_dir) / "report.sarif", "w") as f:
        json.dump(sarif, f, indent=2)


if __name__ == "__main__":
    main()