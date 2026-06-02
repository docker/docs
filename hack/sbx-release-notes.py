#!/usr/bin/env -S uv run --quiet --script
# /// script
# requires-python = ">=3.11"
# dependencies = ["jinja2"]
# ///
"""
Fetch recent stable releases from a GitHub releases page and splice them into
a docs markdown file between the BEGIN/END markers.

Usage (from repo root):

    ./hack/sbx-release-notes.py
    ./hack/sbx-release-notes.py --preset dhi
    GITHUB_TOKEN=$(gh auth token) ./hack/sbx-release-notes.py
    ./hack/sbx-release-notes.py --minor-releases 3
"""

from __future__ import annotations

import argparse
import json
import os
import re
import shutil
import subprocess
import sys
import urllib.request
from collections import defaultdict
from pathlib import Path

from jinja2 import Template

PRESETS: dict[str, dict] = {
    "sbx": {
        "repo": "docker/sbx-releases",
        "file": Path("content/manuals/ai/sandboxes/release-notes.md"),
    },
    "dhi": {
        "repo": "docker-hardened-images/dhictl",
        "file": Path("content/manuals/dhi/release-notes/cli.md"),
    },
}

DEFAULT_PRESET = "sbx"
DEFAULT_MINOR_RELEASES = 3

BEGIN = "<!-- BEGIN GENERATED RELEASES -->"
END = "<!-- END GENERATED RELEASES -->"

SEMVER = re.compile(r"^v(\d+)\.(\d+)\.(\d+)$")

TEMPLATE = Template(
    """\
{% for r in releases -%}
## {{ r.version }}

{{ '{{<' }} release-date date="{{ r.date }}" {{ '>}}' }}

[GitHub release]({{ r.url }})

{{ r.body }}

{% endfor -%}
"""
)


def fetch(repo: str) -> list[dict]:
    url = f"https://api.github.com/repos/{repo}/releases?per_page=100"
    req = urllib.request.Request(
        url,
        headers={
            "Accept": "application/vnd.github+json",
            "X-GitHub-Api-Version": "2022-11-28",
        },
    )
    if token := os.environ.get("GITHUB_TOKEN"):
        req.add_header("Authorization", f"Bearer {token}")
    with urllib.request.urlopen(req) as resp:
        return json.load(resp)


def parse_stable(raw: list[dict]) -> list[dict]:
    out = []
    for r in raw:
        if r.get("prerelease") or r.get("draft"):
            continue
        m = SEMVER.match(r["tag_name"])
        if not m:
            continue
        body = (r.get("body") or "").strip()
        if not body:
            continue
        major, minor, patch = (int(x) for x in m.groups())
        out.append(
            {
                "major": major,
                "minor": minor,
                "patch": patch,
                "version": f"{major}.{minor}.{patch}",
                "date": r["published_at"][:10],
                "url": r["html_url"],
                "body": normalize_body(shift_headings(body)),
            }
        )
    out.sort(key=lambda r: (r["major"], r["minor"], r["patch"]), reverse=True)
    return out


def pick_minor_releases(releases: list[dict], n: int) -> list[dict]:
    by_minor: dict[tuple[int, int], list[dict]] = defaultdict(list)
    for r in releases:
        by_minor[(r["major"], r["minor"])].append(r)
    latest_keys = sorted(by_minor.keys(), reverse=True)[:n]
    keep = set(latest_keys)
    return [r for r in releases if (r["major"], r["minor"]) in keep]


def shift_headings(body: str) -> str:
    """Demote ATX headings by one level so body H2s become H3 under the
    version's H2. Skips fenced code blocks."""
    lines = body.splitlines()
    in_fence = False
    for i, line in enumerate(lines):
        stripped = line.lstrip(" \t")
        if stripped.startswith(("```", "~~~")):
            in_fence = not in_fence
            continue
        if in_fence:
            continue
        if stripped.startswith("#"):
            indent = line[: len(line) - len(stripped)]
            lines[i] = f"{indent}#{stripped}"
    return "\n".join(lines)


def normalize_body(body: str) -> str:
    """Fix markdownlint issues in release body content:
    - Ensure a blank line follows each heading (MD022).
    - Add 'console' language tag to fenced code blocks that have none (MD040).
    Safe to run on content that already complies — no double blank lines are added."""
    lines = body.splitlines()
    result: list[str] = []
    in_fence = False

    for i, line in enumerate(lines):
        stripped = line.lstrip(" \t")

        if stripped.startswith(("```", "~~~")):
            if not in_fence:
                # Opening fence: add language tag if missing
                fence_marker = "```" if stripped.startswith("```") else "~~~"
                lang = stripped[len(fence_marker):].strip()
                if not lang:
                    indent = line[: len(line) - len(stripped)]
                    line = f"{indent}{fence_marker}console"
            in_fence = not in_fence
            result.append(line)
            continue

        result.append(line)

        # Outside fences: insert blank line after a heading if the next line is non-empty
        if not in_fence and stripped.startswith("#"):
            next_line = lines[i + 1] if i + 1 < len(lines) else ""
            if next_line.strip():
                result.append("")

    return "\n".join(result)


def splice(path: Path, generated: str) -> None:
    src = path.read_text()
    try:
        before, rest = src.split(BEGIN, 1)
        _, after = rest.split(END, 1)
    except ValueError:
        sys.exit(f"markers {BEGIN!r} / {END!r} not found in {path}")
    path.write_text(f"{before}{BEGIN}\n\n{generated}{END}{after}")


def main() -> None:
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument(
        "--preset",
        choices=list(PRESETS),
        default=DEFAULT_PRESET,
        help="Named preset that sets --repo and --file defaults (default: %(default)s)",
    )
    p.add_argument("--repo", default=None, help="GitHub repo (owner/name), overrides preset")
    p.add_argument("--file", type=Path, default=None, help="Target markdown file, overrides preset")
    p.add_argument("--minor-releases", type=int, default=DEFAULT_MINOR_RELEASES)
    args = p.parse_args()

    preset = PRESETS[args.preset]
    repo = args.repo or preset["repo"]
    file = args.file or preset["file"]

    releases = pick_minor_releases(parse_stable(fetch(repo)), args.minor_releases)
    if not releases:
        sys.exit("no stable releases found")

    generated = TEMPLATE.render(releases=releases)
    splice(file, generated)
    if shutil.which("npx"):
        subprocess.run(["npx", "--no-install", "prettier", "--write", str(file)], check=False)
    print(f"Wrote {len(releases)} releases (latest {args.minor_releases} minor releases) to {file}")


if __name__ == "__main__":
    main()
