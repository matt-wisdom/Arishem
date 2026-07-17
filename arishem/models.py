from dataclasses import dataclass, field
from typing import Any, Dict, List, Optional
import json

@dataclass
class FunctionSpec:
    name: str                        # e.g., "arishem_chat"
    params: dict[str, str]           # Parameter name mapped to string representation of its type, e.g., {"user_message": "str"}
    return_type: str                 # String representation of return type, e.g., "str"
    docstring: str
    source: str                      # raw source for white-box analysis

    def to_dict(self) -> Dict[str, Any]:
        return {
            "name": self.name,
            "params": self.params,
            "return_type": self.return_type,
            "docstring": self.docstring,
            "source": self.source,
        }

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "FunctionSpec":
        return cls(
            name=data["name"],
            params=data["params"],
            return_type=data["return_type"],
            docstring=data.get("docstring", ""),
            source=data.get("source", ""),
        )

@dataclass
class Turn:
    probe: dict                      # kwargs passed to the function
    raw_output: Any                  # whatever the function returned (e.g., return value or error message)
    classification: str              # blocked | partial | leaked | unexpected | error
    reasoning: str                   # Arishem's observation note for this turn

    def to_dict(self) -> Dict[str, Any]:
        return {
            "probe": self.probe,
            "raw_output": self.raw_output,
            "classification": self.classification,
            "reasoning": self.reasoning,
        }

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "Turn":
        return cls(
            probe=data["probe"],
            raw_output=data["raw_output"],
            classification=data["classification"],
            reasoning=data["reasoning"],
        )

@dataclass
class ProbeSession:
    function_spec: FunctionSpec
    attack_class: str                # e.g. "goal_hijacking"
    goal: str                        # what Arishem is trying to achieve
    history: list[Turn] = field(default_factory=list)
    budget: int = 8                  # max turns
    status: str = "active"           # active | found | exhausted | blocked
    history_summary: str = ""
    scout_analysis: str = ""

    def to_dict(self) -> Dict[str, Any]:
        return {
            "function_spec": self.function_spec.to_dict(),
            "attack_class": self.attack_class,
            "goal": self.goal,
            "history": [t.to_dict() for t in self.history],
            "budget": self.budget,
            "status": self.status,
            "history_summary": self.history_summary,
            "scout_analysis": self.scout_analysis,
        }

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "ProbeSession":
        return cls(
            function_spec=FunctionSpec.from_dict(data["function_spec"]),
            attack_class=data["attack_class"],
            goal=data["goal"],
            history=[Turn.from_dict(t) for t in data.get("history", [])],
            budget=data.get("budget", 8),
            status=data.get("status", "active"),
            history_summary=data.get("history_summary", ""),
            scout_analysis=data.get("scout_analysis", ""),
        )

@dataclass
class Finding:
    id: str
    attack_class: str
    severity: str                    # critical | high | medium | low | info
    title: str
    attack_type: str                 # white-box | black-box
    description: str
    probe_history: list[Turn]        # full multi-turn trace as evidence
    function_targeted: str
    remediation: str

    def to_dict(self) -> Dict[str, Any]:
        return {
            "id": self.id,
            "attack_class": self.attack_class,
            "severity": self.severity,
            "title": self.title,
            "attack_type": self.attack_type,
            "description": self.description,
            "probe_history": [t.to_dict() for t in self.probe_history],
            "function_targeted": self.function_targeted,
            "remediation": self.remediation,
        }

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "Finding":
        return cls(
            id=data["id"],
            attack_class=data["attack_class"],
            severity=data["severity"],
            title=data["title"],
            attack_type=data["attack_type"],
            description=data["description"],
            probe_history=[Turn.from_dict(t) for t in data["probe_history"]],
            function_targeted=data["function_targeted"],
            remediation=data["remediation"],
        )

@dataclass
class ArishemRunResult:
    findings: List[Finding]
    sessions: List[ProbeSession]

    def to_dict(self) -> Dict[str, Any]:
        return {
            "findings": [f.to_dict() for f in self.findings],
            "sessions": [s.to_dict() for s in self.sessions],
        }

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> "ArishemRunResult":
        return cls(
            findings=[Finding.from_dict(f) for f in data.get("findings", [])],
            sessions=[ProbeSession.from_dict(s) for s in data.get("sessions", [])]
        )

