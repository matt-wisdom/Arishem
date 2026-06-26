import ast
from typing import List, Dict
from arishem.models import FunctionSpec

def load_target_module(filepath: str) -> List[FunctionSpec]:
    """
    Statically parses a python file using AST, extracting FunctionSpec definitions
    for all functions prefixed with 'arishem_'.
    Does not import or execute the module.
    """
    with open(filepath, "r", encoding="utf-8") as f:
        source_code = f.read()

    try:
        tree = ast.parse(source_code, filename=filepath)
    except SyntaxError as e:
        raise ValueError(f"Failed to parse target file '{filepath}' due to a syntax error: {e}")

    specs = []

    for node in tree.body:
        if isinstance(node, (ast.FunctionDef, ast.AsyncFunctionDef)):
            if node.name.startswith("arishem_"):
                # Extract parameters and their type annotations
                params = {}
                for arg in node.args.args:
                    param_name = arg.arg
                    if arg.annotation:
                        param_type = ast.unparse(arg.annotation).strip()
                    else:
                        param_type = "Any"
                    params[param_name] = param_type

                # Extract return type annotation
                if node.returns:
                    return_type = ast.unparse(node.returns).strip()
                else:
                    return_type = "None"

                # Extract docstring
                docstring = ast.get_docstring(node) or ""

                # Extract source code segment of the function
                source = ast.get_source_segment(source_code, node) or ast.unparse(node)

                spec = FunctionSpec(
                    name=node.name,
                    params=params,
                    return_type=return_type,
                    docstring=docstring,
                    source=source,
                )
                specs.append(spec)

    return specs
