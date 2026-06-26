import json
import os
import logging
from typing import Dict, Any
from arishem.models import ArishemRunResult

logger = logging.getLogger("arishem.store")

def save_run_result(result: ArishemRunResult, filepath: str):
    """
    Saves the entire run result (findings and session traces) into a JSON file.
    Creates parent directories if necessary.
    """
    directory = os.path.dirname(filepath)
    if directory and not os.path.exists(directory):
        os.makedirs(directory, exist_ok=True)

    data = result.to_dict()
    with open(filepath, "w", encoding="utf-8") as f:
        json.dump(data, f, indent=2, ensure_ascii=False)
    logger.info(f"Successfully saved test results to: {filepath}")

def load_run_result(filepath: str) -> ArishemRunResult:
    """
    Loads and parses an ArishemRunResult from a JSON file.
    """
    if not os.path.exists(filepath):
        raise FileNotFoundError(f"Run result file not found at: {filepath}")

    with open(filepath, "r", encoding="utf-8") as f:
        data = json.load(f)
    
    return ArishemRunResult.from_dict(data)
