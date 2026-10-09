"""Verify the companion examples locally; writes only to output/."""

import difflib
import hashlib
import json
from pathlib import Path
import shutil
import subprocess
import sys
import tempfile
from zipfile import ZipFile

ROOT = Path(__file__).resolve().parent
OUTPUT = ROOT / "output"


def run(command, cwd, log):
    result = subprocess.run(command, cwd=cwd, capture_output=True, text=True)
    (OUTPUT / log).write_text(result.stdout + result.stderr, encoding="utf-8")
    return result


def require_success(result):
    if result.returncode:
        raise RuntimeError(result.stdout + result.stderr)


def baseline(project, label):
    require_success(
        run(
            ["npm", "ci", "--ignore-scripts", "--no-audit", "--no-fund"],
            project,
            f"{label}-install.txt",
        )
    )
    result = run(["npm", "test"], project, f"{label}-tests.txt")
    if result.returncode:
        (OUTPUT / f"{label}-report.md").write_text(
            "# Investigation blocked\n\n"
            "The Express 4.21.2 baseline failed. No upgrade was attempted.\n\n"
            f"See `{label}-tests.txt`. Resolve the baseline failure before "
            "investigating Express 5.1.0.\n",
            encoding="utf-8",
        )
        return False
    return True


def verify_express(work):
    project = work / "express"
    shutil.copytree(
        ROOT / "express-upgrade",
        project,
        ignore=shutil.ignore_patterns("node_modules", "output"),
    )
    assert baseline(project, "express-baseline")
    original = (project / "app.js").read_text()
    baseline_files = {
        name: (project / name).read_text()
        for name in ["app.js", "package.json", "package-lock.json"]
    }
    require_success(
        run(
            [
                "npm",
                "install",
                "--save-exact",
                "express@5.1.0",
                "--ignore-scripts",
                "--no-audit",
                "--no-fund",
            ],
            project,
            "express-upgrade-install.txt",
        )
    )
    failed = run(["npm", "test"], project, "express-upgrade-tests.txt")
    assert failed.returncode != 0 and "Missing parameter name" in failed.stdout
    (project / "app.js").write_text(
        original.replace('app.get("*",', 'app.get("/*splat",')
    )
    partial = run(["npm", "test"], project, "express-partial-tests.txt")
    assert (
        partial.returncode != 0
        and "# pass 4" in partial.stdout
        and "# fail 4" in partial.stdout
    )
    routing_fix = original.replace('app.get("*",', 'app.get("/{*splat}",')
    (project / "app.js").write_text(routing_fix)
    routing = run(["npm", "test"], project, "express-routing-tests.txt")
    assert (
        routing.returncode != 0
        and "# pass 5" in routing.stdout
        and "# fail 3" in routing.stdout
    )
    parser_fix = routing_fix.replace(
        "const app = express();",
        'const app = express();\napp.set("query parser", "extended");',
    )
    (project / "app.js").write_text(parser_fix)
    parser = run(["npm", "test"], project, "express-parser-tests.txt")
    assert (
        parser.returncode != 0
        and "# pass 7" in parser.stdout
        and "# fail 1" in parser.stdout
    )
    fixed = parser_fix.replace(
        "const { filters } = req.body;", "const { filters } = req.body ?? {};"
    )
    (project / "app.js").write_text(fixed)
    require_success(run(["npm", "test"], project, "express-fixed-tests.txt"))
    require_success(
        run(
            ["npm", "ls", "express", "router", "path-to-regexp"],
            project,
            "express-versions.txt",
        )
    )
    shutil.copyfile(
        project / "package-lock.json", OUTPUT / "express-target-package-lock.json"
    )
    app_patch = "".join(
        difflib.unified_diff(
            original.splitlines(True),
            fixed.splitlines(True),
            fromfile="a/app.js",
            tofile="b/app.js",
        )
    )
    patch = "".join(
        "".join(
            difflib.unified_diff(
                content.splitlines(True),
                (project / name).read_text().splitlines(True),
                fromfile=f"a/{name}",
                tofile=f"b/{name}",
            )
        )
        for name, content in baseline_files.items()
    )
    (OUTPUT / "express-fix.patch").write_text(patch)
    (OUTPUT / "express-report.md").write_text(
        "# Plan the Express 5.1.0 migration\n\n"
        "Recommendation: preserve the API contract before deploying. The "
        "proposed compatibility fixes pass eight tests, but production traffic "
        "and security review remain outside this fixture's coverage.\n\n"
        "Local reproduction only; this report has not been posted to GitHub.\n\n"
        f"Baseline lockfile SHA-256: `{hashlib.sha256((ROOT / 'express-upgrade/package-lock.json').read_bytes()).hexdigest()}`.\n\n"
        f"Baseline app SHA-256: `{hashlib.sha256(original.encode()).hexdigest()}`.\n\n"
        "| Stage | Result |\n| --- | --- |\n"
        "| Express 4.21.2 | Eight tests pass |\n"
        "| Express 5.1.0, unchanged route | Application fails to load: Missing parameter name |\n"
        "| Express 5.1.0, /*splat | Four pass; root, nested filters, and empty search fail |\n"
        "| Express 5.1.0, /{*splat} | Five pass; nested filters and empty search fail |\n"
        "| Extended query parser restored | Seven pass; empty search fails |\n"
        "| Empty request body handled | Eight tests pass |\n\n"
        "Compatibility changes: name the fallback wildcard and include `/`; "
        "preserve nested query filters by selecting the extended parser; "
        "handle an unparsed request body as an empty search. Without these fixes, "
        "clients cannot open the app, filtered searches return unwanted tickets, "
        "and empty searches return HTTP 500.\n\n"
        f"```diff\n{app_patch}```\n\n"
        "Follow-up: verify real client requests and filter validation/limits, "
        "review dependency advisories, and test rollout and rollback. "
        "The full proposed patch, target lockfile, dependency versions, and "
        "individual test logs are saved alongside this report.\n",
        encoding="utf-8",
    )
    blocked = work / "blocked"
    shutil.copytree(
        ROOT / "express-upgrade",
        blocked,
        ignore=shutil.ignore_patterns("node_modules", "output"),
    )
    tests = blocked / "test/app.test.js"
    tests.write_text(
        tests.read_text().replace("response.status, 200", "response.status, 418", 1)
    )
    lock_before = (blocked / "package-lock.json").read_bytes()
    assert not baseline(blocked, "express-blocked-baseline")
    assert (blocked / "package-lock.json").read_bytes() == lock_before
    assert (
        json.loads((blocked / "package.json").read_text())["dependencies"]["express"]
        == "4.21.2"
    )


def verify_reporter(work):
    project = ROOT / "api-change-reporter"
    environment = work / "venv"
    subprocess.run([sys.executable, "-m", "venv", str(environment)], check=True)
    python = str(environment / "bin/python")
    require_success(
        run(
            [python, "-m", "pip", "install", "-r", "requirements.txt"],
            project,
            "reporter-install.txt",
        )
    )
    require_success(
        run(
            [python, "-m", "unittest", "discover", "-s", "tests", "-v"],
            project,
            "reporter-tests.txt",
        )
    )
    require_success(
        run(
            [
                python,
                "reporter.py",
                "examples/before.yaml",
                "examples/after.json",
                "--html",
                str(OUTPUT / "index.html"),
            ],
            project,
            "reporter-result.txt",
        )
    )
    archive = OUTPUT / "api-change-reporter.zip"
    with ZipFile(archive, "w") as bundle:
        for name in [
            "reporter.py",
            "requirements.txt",
            "README.md",
            "tests/test_reporter.py",
            "examples/before.yaml",
            "examples/after.json",
        ]:
            bundle.write(project / name, f"api-change-reporter/{name}")
        bundle.write(OUTPUT / "index.html", "api-change-reporter/report/index.html")
    extracted = work / "extracted"
    with ZipFile(archive) as bundle:
        assert bundle.testzip() is None
        bundle.extractall(extracted)
    require_success(
        run(
            [python, "-m", "unittest", "discover", "-s", "tests", "-v"],
            extracted / "api-change-reporter",
            "reporter-bundle-tests.txt",
        )
    )


def main():
    OUTPUT.mkdir(exist_ok=True)
    with tempfile.TemporaryDirectory(prefix="cloud-sandbox-examples-") as directory:
        work = Path(directory)
        verify_reporter(work)
        verify_express(work)
    (OUTPUT / "environment.txt").write_text(
        f"Python {sys.version.split()[0]}\n"
        + subprocess.check_output(["node", "--version"], text=True)
        + "Local execution only. Console, sbx, and GitHub delivery are unverified.\n"
    )
    print(f"All local checks passed. Evidence and bundle: {OUTPUT}")


if __name__ == "__main__":
    main()
