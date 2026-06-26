import os
import subprocess
import tempfile
import json
import pytest

def test_cli_list_classes():
    res = subprocess.run(
        ["python3", "-m", "arishem.cli", "list-classes"],
        capture_output=True,
        text=True
    )
    assert res.returncode == 0
    assert "Registered Arishem Attack Classes" in res.stdout
    assert "goal_hijacking" in res.stdout

def test_cli_run_and_report_flow():
    with tempfile.TemporaryDirectory() as tmpdir:
        results_json = os.path.join(tmpdir, "res.json")
        html_report = os.path.join(tmpdir, "rep.html")
        sarif_report = os.path.join(tmpdir, "rep.sarif")

        # 1. Test running pen test saving results to JSON
        res_run = subprocess.run([
            "python3", "-m", "arishem.cli", "run", "target_example.py",
            "--output", results_json,
            "--classes", "goal_hijacking,tool_misuse"
        ], capture_output=True, text=True)

        assert res_run.returncode == 0
        assert os.path.exists(results_json)

        # Verify output json parses
        with open(results_json, "r") as f:
            data = json.load(f)
            assert "findings" in data
            assert "sessions" in data

        # 2. Test generating HTML report via report subcommand
        res_html = subprocess.run([
            "python3", "-m", "arishem.cli", "report", results_json,
            "--format", "html", "--output", html_report
        ], capture_output=True, text=True)

        assert res_html.returncode == 0
        assert os.path.exists(html_report)
        with open(html_report, "r") as f:
            html_content = f.read()
            assert "<!DOCTYPE html>" in html_content
            assert "Arishem Security Assessment Report" in html_content

        # 3. Test generating SARIF report via report subcommand
        res_sarif = subprocess.run([
            "python3", "-m", "arishem.cli", "report", results_json,
            "--format", "sarif", "--output", sarif_report,
            "--target-src", "target_example.py"
        ], capture_output=True, text=True)

        assert res_sarif.returncode == 0
        assert os.path.exists(sarif_report)
        with open(sarif_report, "r") as f:
            sarif_data = json.load(f)
            assert sarif_data["version"] == "2.1.0"
            assert len(sarif_data["runs"]) > 0
