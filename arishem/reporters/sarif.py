import json
import os
import logging
from typing import Dict, Any
from arishem.models import ArishemRunResult

logger = logging.getLogger("arishem.reporters.sarif")

def generate_sarif_report(result: ArishemRunResult, filepath: str, target_filepath: str):
    """
    Generates a standardized SARIF (v2.1.0) JSON report for GitHub security / CI integration.
    """
    # Build relative path to display in editors
    rel_path = os.path.relpath(target_filepath) if os.path.isabs(target_filepath) else target_filepath

    sarif_data = {
        "$schema": "https://docs.oasis-open.org/sarif/sarif/v2.1.0/cos02/schemas/sarif-schema-2.1.0.json",
        "version": "2.1.0",
        "runs": [
            {
                "tool": {
                    "driver": {
                        "name": "Arishem AI Pentesting Engine",
                        "semanticVersion": "0.1.0",
                        "rules": [
                            {
                                "id": "ARI-GOAL-HIJACK",
                                "name": "goal_hijacking",
                                "shortDescription": {"text": "AI Prompt Injection and Goal Hijacking"},
                                "helpUri": "https://github.com/google-deepmind/arishem"
                            },
                            {
                                "id": "ARI-TOOL-MISUSE",
                                "name": "tool_misuse",
                                "shortDescription": {"text": "Agent Tool Chain Abuse and Misuse"},
                                "helpUri": "https://github.com/google-deepmind/arishem"
                            },
                            {
                                "id": "ARI-PRIV-ESC",
                                "name": "privilege_abuse",
                                "shortDescription": {"text": "Privilege and Scope Escalation"},
                                "helpUri": "https://github.com/google-deepmind/arishem"
                            },
                            {
                                "id": "ARI-UNSAFE-EXEC",
                                "name": "unsafe_execution",
                                "shortDescription": {"text": "Adversarial Code or Shell Command Execution"},
                                "helpUri": "https://github.com/google-deepmind/arishem"
                            },
                            {
                                "id": "ARI-CTX-POISON",
                                "name": "context_poisoning",
                                "shortDescription": {"text": "Context or History poisoning"},
                                "helpUri": "https://github.com/google-deepmind/arishem"
                            },
                            {
                                "id": "ARI-INSEC-DELEG",
                                "name": "insecure_delegation",
                                "shortDescription": {"text": "Insecure agent-to-agent delegation"},
                                "helpUri": "https://github.com/google-deepmind/arishem"
                            },
                            {
                                "id": "ARI-CASC-FAIL",
                                "name": "cascading_failure",
                                "shortDescription": {"text": "Uncontained cascading system failure"},
                                "helpUri": "https://github.com/google-deepmind/arishem"
                            },
                            {
                                "id": "ARI-TRUST-EXPL",
                                "name": "trust_exploitation",
                                "shortDescription": {"text": "Human Trust Exploitation via Urgency"},
                                "helpUri": "https://github.com/google-deepmind/arishem"
                            },
                            {
                                "id": "ARI-PERS-DRIFT",
                                "name": "persistent_drift",
                                "shortDescription": {"text": "Persistent agent misalignment and drift"},
                                "helpUri": "https://github.com/google-deepmind/arishem"
                            }
                        ]
                    }
                },
                "results": []
            }
        ]
    }

    # Map finding severities to SARIF level values
    severity_map = {
        "critical": "error",
        "high": "error",
        "medium": "warning",
        "low": "note",
        "info": "note"
    }

    rules_by_name = {rule["name"]: rule["id"] for rule in sarif_data["runs"][0]["tool"]["driver"]["rules"]}

    results_list = []
    for f in result.findings:
        rule_id = rules_by_name.get(f.attack_class, "ARI-GOAL-HIJACK")
        level = severity_map.get(f.severity.lower(), "warning")

        sarif_finding = {
            "ruleId": rule_id,
            "message": {
                "text": f"{f.title}\n\nDescription: {f.description}\n\nRemediation: {f.remediation}"
            },
            "level": level,
            "locations": [
                {
                    "physicalLocation": {
                        "artifactLocation": {
                            "uri": rel_path
                        },
                        "region": {
                            "startLine": 1,
                            "startColumn": 1
                        }
                    }
                }
            ],
            "properties": {
                "targetedFunction": f.function_targeted,
                "attackType": f.attack_type,
                "findingId": f.id
            }
        }
        results_list.append(sarif_finding)

    sarif_data["runs"][0]["results"] = results_list

    directory = os.path.dirname(filepath)
    if directory and not os.path.exists(directory):
        os.makedirs(directory, exist_ok=True)

    with open(filepath, "w", encoding="utf-8") as f:
        json.dump(sarif_data, f, indent=2, ensure_ascii=False)
    logger.info(f"SARIF report successfully saved to: {filepath}")
