import logging
from typing import Dict, Any, List
from rich.console import Console
from rich.table import Table
from rich.live import Live
from rich.panel import Panel
from rich.text import Text

from arishem.models import ProbeSession, Turn, Finding, ArishemRunResult

console = Console()

class ConsoleReporter:
    def __init__(self):
        self.live: Optional[Live] = None
        self.sessions_status: Dict[str, Dict[str, Any]] = {}

    def start(self, sessions: List[ProbeSession]):
        """Initializes the live table monitor in the terminal."""
        for s in sessions:
            key = f"{s.function_spec.name}:{s.attack_class}"
            self.sessions_status[key] = {
                "function": s.function_spec.name,
                "class": s.attack_class,
                "turns": 0,
                "budget": s.budget,
                "status": "PENDING"
            }
        
        self.live = Live(self._generate_table(), refresh_per_second=4, console=console)
        self.live.start()

    def update_session(self, session: ProbeSession, turn: Turn):
        """Updates the status of a specific session on turn completion."""
        key = f"{session.function_spec.name}:{session.attack_class}"
        if key in self.sessions_status:
            self.sessions_status[key]["turns"] = len(session.history)
            self.sessions_status[key]["status"] = session.status.upper()
            if self.live:
                self.live.update(self._generate_table())

    def finish(self, result: ArishemRunResult):
        """Stops the live monitor and prints a summary of findings."""
        if self.live:
            # Sync final states
            for s in result.sessions:
                key = f"{s.function_spec.name}:{s.attack_class}"
                if key in self.sessions_status:
                    self.sessions_status[key]["turns"] = len(s.history)
                    self.sessions_status[key]["status"] = s.status.upper()
            self.live.update(self._generate_table())
            self.live.stop()

        console.print("\n[bold cyan]=== Penetration Test Execution Complete ===[/bold cyan]\n")

        # Check for sessions with error status and display details
        error_sessions = [s for s in result.sessions if s.status == "error"]
        if error_sessions:
            console.print("[bold red]✗ The following sessions aborted due to execution/infrastructure errors:[/bold red]\n")
            for s in error_sessions:
                console.print(f"[red]• Function: [bold]{s.function_spec.name}[/bold] | Attack Class: {s.attack_class}[/red]")
                if s.history:
                    last_turn = s.history[-1]
                    console.print(f"  [yellow]Error Message:[/yellow] {last_turn.raw_output}")
                else:
                    console.print("  [yellow]Error Message:[/yellow] Session aborted before executing any turns.")
                console.print()

        if not result.findings:
            console.print("[bold green]✔ All clean! No vulnerabilities confirmed in target functions.[/bold green]")
            return

        console.print(f"[bold red]✗ Confirmed {len(result.findings)} vulnerabilities:[/bold red]\n")

        for f in result.findings:
            severity_style = "bold red"
            if f.severity.lower() == "critical":
                severity_style = "bold white on red"
            elif f.severity.lower() == "high":
                severity_style = "bold red"
            elif f.severity.lower() == "medium":
                severity_style = "bold yellow"
            elif f.severity.lower() == "low":
                severity_style = "bold blue"

            panel_content = Text()
            panel_content.append(f"Function: {f.function_targeted}\n", style="bold")
            panel_content.append(f"Attack Class: {f.attack_class}\n\n", style="italic")
            panel_content.append(f"Description:\n{f.description}\n\n", style="default")
            panel_content.append(f"Remediation:\n", style="bold green")
            panel_content.append(f"{f.remediation}", style="green")

            console.print(Panel(
                panel_content,
                title=f"[{f.severity.upper()}] {f.title}",
                title_align="left",
                border_style=severity_style,
                padding=(1, 2)
            ))
            console.print()

    def _generate_table(self) -> Table:
        table = Table(title="Arishem Active Probing Monitor", expand=True)
        table.add_column("Target Function", style="cyan", no_wrap=True)
        table.add_column("Attack Class", style="magenta")
        table.add_column("Turns", justify="right", style="green")
        table.add_column("Status", justify="center")

        for _, status in self.sessions_status.items():
            status_str = status["status"]
            if status_str in ("PENDING", "ACTIVE"):
                status_color = "[yellow]ACTIVE[/yellow]"
            elif status_str == "FOUND":
                status_color = "[bold red]VULN FOUND[/bold red]"
            elif status_str == "BLOCKED":
                status_color = "[green]BLOCKED[/green]"
            elif status_str == "EXHAUSTED":
                status_color = "[blue]EXHAUSTED[/blue]"
            else:
                status_color = f"[white]{status_str}[/white]"

            table.add_row(
                status["function"],
                status["class"],
                f"{status['turns']}/{status['budget']}",
                status_color
            )
        return table
