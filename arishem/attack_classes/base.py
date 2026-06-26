import os
import importlib
import logging
from dataclasses import dataclass
from typing import Dict, List, Any

logger = logging.getLogger("arishem.attack_classes")

@dataclass
class AttackClassConfig:
    id: str
    name: str
    attacker_goal_template: str      # Template with placeholders like {function_name} and {params}
    attacker_system_prompt: str      # Persona and guidance
    example_payloads: List[Dict[str, Any]]  # Seeds for turn 1

# Central registry for attack classes
_ATTACK_CLASSES_REGISTRY: Dict[str, AttackClassConfig] = {}

def register_attack_class(config: AttackClassConfig):
    """Manually register an attack class config."""
    _ATTACK_CLASSES_REGISTRY[config.id] = config
    logger.debug(f"Registered attack class: {config.id}")

def get_attack_class(class_id: str) -> AttackClassConfig:
    """Retrieves an attack class by its unique id."""
    if not _ATTACK_CLASSES_REGISTRY:
        discover_attack_classes()
    if class_id not in _ATTACK_CLASSES_REGISTRY:
        raise ValueError(f"Attack class '{class_id}' is not registered.")
    return _ATTACK_CLASSES_REGISTRY[class_id]

def get_all_attack_classes() -> Dict[str, AttackClassConfig]:
    """Returns all registered attack classes, running discovery if necessary."""
    if not _ATTACK_CLASSES_REGISTRY:
        discover_attack_classes()
    return _ATTACK_CLASSES_REGISTRY

def discover_attack_classes():
    """
    Scans the directory of this file for Python modules and loads their 'config' attributes.
    """
    current_dir = os.path.dirname(os.path.abspath(__file__))
    
    for filename in os.listdir(current_dir):
        if filename.endswith(".py") and filename not in ("__init__.py", "base.py"):
            module_name = filename[:-3]
            try:
                # Dynamically load the module relative to this package
                module = importlib.import_module(f"arishem.attack_classes.{module_name}")
                if hasattr(module, "config"):
                    cfg = getattr(module, "config")
                    if isinstance(cfg, AttackClassConfig):
                        register_attack_class(cfg)
                    else:
                        logger.warning(f"Module {module_name} exports 'config' but it is not an AttackClassConfig instance.")
                else:
                    logger.warning(f"Module {module_name} does not export 'config'.")
            except Exception as e:
                logger.error(f"Failed to dynamically load attack class module '{module_name}': {e}")
