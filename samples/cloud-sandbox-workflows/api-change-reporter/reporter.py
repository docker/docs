"""Report removed inline OpenAPI paths and operations, not full compatibility."""

import argparse
from html import escape
from pathlib import Path

import yaml

METHODS = {"get", "put", "post", "delete", "options", "head", "patch", "trace"}


def load_spec(path):
    with Path(path).open(encoding="utf-8") as source:
        spec = yaml.safe_load(source)
    if not isinstance(spec, dict) or not isinstance(spec.get("paths"), dict):
        raise ValueError("Expected an OpenAPI document with a paths object")
    for path, item in spec["paths"].items():
        if not isinstance(path, str) or not path.startswith("/"):
            raise ValueError("Expected paths beginning with /")
        if not isinstance(item, dict) or "$ref" in item:
            raise ValueError("Only inline Path Item objects are supported")
    return spec


def compare(before, after):
    old_paths, next_paths = before["paths"], after["paths"]
    removed_paths = sorted(set(old_paths) - set(next_paths))
    removed_operations = {}
    for path, item in sorted(old_paths.items()):
        missing = (set(item) & METHODS) - (set(next_paths.get(path, {})) & METHODS)
        if missing:
            removed_operations[path] = sorted(method.upper() for method in missing)
    return {"removed_paths": removed_paths, "removed_operations": removed_operations}


def render_report(result):
    operations = result["removed_operations"]
    count = sum(len(methods) for methods in operations.values())
    rows = (
        "".join(
            f"<tr><td><code>{escape(path)}</code></td>"
            f"<td>{escape(', '.join(methods))}</td></tr>"
            for path, methods in operations.items()
        )
        or '<tr><td colspan="2">No removed operations</td></tr>'
    )
    paths = ", ".join(escape(path) for path in result["removed_paths"]) or "None"
    return f"""<!doctype html>
<html lang="en"><meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>API-change report</title>
<style>
body {{ font: 18px/1.6 system-ui, sans-serif; color: #142b45;
        background: #eef4fb; margin: 0; padding: 48px 24px; }}
main {{ max-width: 780px; margin: auto; padding: 32px;
        background: white; border-radius: 12px; border: 1px solid #cad8e8; }}
h1 {{ margin-top: 0; }} table {{ width: 100%; border-collapse: collapse; }}
th, td {{ text-align: left; padding: 12px; border-bottom: 1px solid #cad8e8; }}
code {{ color: #0758a6; }} .scope {{ color: #475569; }}
</style><main>
<h1>API-change report</h1>
<p>Removed paths: {len(result['removed_paths'])} · Removed operations: {count}</p>
<table><caption>Removed operations grouped by path</caption>
<thead><tr><th scope="col">Path</th><th scope="col">Methods</th></tr></thead>
<tbody>{rows}</tbody></table>
<p>Removed paths: <code>{paths}</code></p>
<p class="scope">Checks inline paths and HTTP operations only.
This report does not establish complete API compatibility.</p>
</main></html>
"""


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("before", type=Path)
    parser.add_argument("after", type=Path)
    parser.add_argument("--html", type=Path, required=True)
    args = parser.parse_args()
    try:
        result = compare(load_spec(args.before), load_spec(args.after))
        args.html.parent.mkdir(parents=True, exist_ok=True)
        args.html.write_text(render_report(result), encoding="utf-8")
    except (OSError, ValueError, yaml.YAMLError) as error:
        parser.exit(2, f"error: {error}\n")
    print(f"Removed paths: {len(result['removed_paths'])}")
    print(f"Removed operations: {sum(map(len, result['removed_operations'].values()))}")
    for path, methods in result["removed_operations"].items():
        print(f"  {path}: {', '.join(methods)}")


if __name__ == "__main__":
    main()
