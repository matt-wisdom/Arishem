import argparse
import asyncio
import logging
import sys
import os
from typing import List, Optional

from rich.console import Console

from arishem.orchestrator import run_arishem
from arishem.attack_classes import discover_attack_classes, get_all_attack_classes
from arishem.store import save_run_result, load_run_result
from arishem.reporters.console import ConsoleReporter
from arishem.reporters.html import generate_html_report
from arishem.reporters.sarif import generate_sarif_report

console = Console()

def load_env():
    import os
    for path in [".env", "../.env", "../../.env"]:
        if os.path.exists(path):
            with open(path, "r") as f:
                for line in f:
                    line = line.strip()
                    if not line or line.startswith("#"):
                        continue
                    parts = line.split("=", 1)
                    if len(parts) == 2:
                        key, val = parts[0].strip(), parts[1].strip()
                        val = val.strip("'\"")
                        if key not in os.environ:
                            os.environ[key] = val
            break

def main():
    logging.basicConfig(
        level=logging.INFO,
        format="%(asctime)s | %(levelname)s | %(name)s | %(message)s",
        stream=sys.stdout
    )
    load_env()
    parser = argparse.ArgumentParser(
        description="Arishem — Autonomous AI Security Testing Engine",
        prog="arishem"
    )
    subparsers = parser.add_subparsers(dest="command", required=True)

    # 1. run subcommand
    run_parser = subparsers.add_parser("run", help="Run active probing penetration test against a target module")
    run_parser.add_argument("target", help="Path to Python file exposing arishem_ functions")
    run_parser.add_argument("--mode", choices=["whitebox", "blackbox", "both"], default="both",
                            help="Testing mode (default: both)")
    run_parser.add_argument("--classes", help="Comma-separated list of attack class ids to run")
    run_parser.add_argument("--budget", type=int, help="Max turns per probe session (default: 8)")
    run_parser.add_argument("--concurrency", type=int, help="Max concurrent probe sessions (default: 4)")
    run_parser.add_argument("--output", help="Path to save report (e.g. report.html or results.json)")
    run_parser.add_argument("--config", default="arishem.config.toml", help="Path to arishem.config.toml")
    run_parser.add_argument("--history-file", help="Path to a file containing historical knowledge from previous runs")

    # 2. list-classes subcommand
    subparsers.add_parser("list-classes", help="List all registered attack classes")

    # 3. report subcommand
    report_parser = subparsers.add_parser("report", help="Generate a report from a saved results JSON file")
    report_parser.add_argument("results_json", help="Path to saved JSON results file")
    report_parser.add_argument("--format", choices=["html", "sarif", "console"], default="console",
                               help="Format to output (default: console)")
    report_parser.add_argument("--output", help="Path to write output file (required for html and sarif)")
    report_parser.add_argument("--target-src", help="Path to target python file (useful for SARIF rendering, optional)")

    args = parser.parse_args()

    if args.command == "run":
        asyncio.run(handle_run(args))
    elif args.command == "list-classes":
        handle_list_classes()
    elif args.command == "report":
        handle_report(args)

async def handle_run(args):
    # Setup callback for live updating reporter
    reporter = ConsoleReporter()

    async def on_turn(session, turn):
        reporter.update_session(session, turn)

    # Parse classes list
    classes_list = None
    if args.classes:
        classes_list = [c.strip() for c in args.classes.split(",") if c.strip()]

    target_path = args.target
    if not os.path.exists(target_path):
        console.print(f"[bold red]Error: Target file '{target_path}' does not exist.[/bold red]")
        sys.exit(1)

    console.print(f"[bold cyan]Scanning target file:[/bold cyan] {target_path}...")
    
    # Run orchestrator
    try:
        # Pre-build run result structure to register sessions in ConsoleReporter
        # We do a quick loader pass to show initial statuses
        from arishem.loader import load_target_module
        from arishem.models import ProbeSession
        from arishem.attack_classes import get_all_attack_classes
        
        discover_attack_classes()
        all_classes = get_all_attack_classes()
        
        # Determine active classes
        active_classes = {}
        if classes_list:
            for cid in classes_list:
                if cid in all_classes:
                    active_classes[cid] = all_classes[cid]
        else:
            active_classes = all_classes

        specs = load_target_module(target_path)
        dummy_sessions = []
        for spec in specs:
            for acid in active_classes:
                dummy_sessions.append(ProbeSession(
                    function_spec=spec,
                    attack_class=acid,
                    goal="",
                    budget=args.budget or 8
                ))

        if not dummy_sessions:
            console.print("[yellow]No 'arishem_' functions found to probe.[/yellow]")
            sys.exit(0)

        # Start live reporting table
        reporter.start(dummy_sessions)

        # Parse history file if provided
        history_summary = ""
        if args.history_file and os.path.exists(args.history_file):
            try:
                with open(args.history_file, "r", encoding="utf-8") as f:
                    history_summary = f.read()
            except Exception as e:
                console.print(f"[yellow]Warning: Failed to read history file: {e}[/yellow]")

        # Run actual penetration test
        result = await run_arishem(
            target_path=target_path,
            mode=args.mode,
            classes=classes_list,
            budget=args.budget,
            concurrency=args.concurrency,
            config_path=args.config,
            on_turn_callback=on_turn,
            history_summary=history_summary
        )

        reporter.finish(result)

        # If output is specified, save appropriately
        if args.output:
            out_path = args.output
            if out_path.endswith(".html") or out_path.endswith(".htm"):
                generate_html_report(result, out_path)
                console.print(f"[green]HTML report saved to {out_path}[/green]")
            elif out_path.endswith(".sarif"):
                generate_sarif_report(result, out_path, target_path)
                console.print(f"[green]SARIF report saved to {out_path}[/green]")
            else:
                # Default to saving raw JSON
                save_run_result(result, out_path)
                console.print(f"[green]JSON results saved to {out_path}[/green]")
        else:
            # Auto-save results.json to current directory for convenience
            save_run_result(result, "results.json")
            console.print("[dim]Raw results autosaved to results.json[/dim]")

    except Exception as e:
        console.print(f"[bold red]Execution failed: {e}[/bold red]", style="red")
        import traceback
        logger = logging.getLogger("arishem")
        logger.debug(traceback.format_exc())
        sys.exit(1)

def handle_list_classes():
    discover_attack_classes()
    classes = get_all_attack_classes()
    
    console.print("\n[bold cyan]=== Registered Arishem Attack Classes ===[/bold cyan]\n")
    for cid, cfg in classes.items():
        console.print(f"[bold green]• {cfg.name}[/bold green] (ID: [yellow]{cid}[/yellow])")
        console.print(f"  [dim]Goal Template: {cfg.attacker_goal_template}[/dim]\n")

def handle_report(args):
    results_path = args.results_json
    if not os.path.exists(results_path):
        console.print(f"[bold red]Error: Results file '{results_path}' not found.[/bold red]")
        sys.exit(1)

    try:
        result = load_run_result(results_path)
    except Exception as e:
        console.print(f"[bold red]Error loading results file: {e}[/bold red]")
        sys.exit(1)

    fmt = args.format
    if fmt == "console":
        # Print findings via ConsoleReporter summary logic
        reporter = ConsoleReporter()
        reporter.finish(result)
    elif fmt == "html":
        if not args.output:
            console.print("[bold red]Error: --output file must be specified for HTML reports.[/bold red]")
            sys.exit(1)
        generate_html_report(result, args.output)
        console.print(f"[green]HTML report generated at {args.output}[/green]")
    elif fmt == "sarif":
        if not args.output:
            console.print("[bold red]Error: --output file must be specified for SARIF reports.[/bold red]")
            sys.exit(1)
        target = args.target_src or "target_module.py"
        generate_sarif_report(result, args.output, target)
        console.print(f"[green]SARIF report generated at {args.output}[/green]")

if __name__ == "__main__":
    main()
