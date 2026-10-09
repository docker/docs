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
        and "# pass 2" in partial.stdout
        and "# fail 1" in partial.stdout
    )
    fixed = original.replace('app.get("*",', 'app.get("/{*splat}",')
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
    patch = "".join(
        difflib.unified_diff(
            original.splitlines(True),
            fixed.splitlines(True),
            fromfile="a/app.js",
            tofile="b/app.js",
        )
    )
    (OUTPUT / "express-fix.patch").write_text(patch)
    (OUTPUT / "express-report.md").write_text(
        "# Express upgrade investigation\n\n"
        "Local fixture run; this report has not been posted to GitHub.\n\n"
        f"Baseline lockfile SHA-256: `{hashlib.sha256((ROOT / 'express-upgrade/package-lock.json').read_bytes()).hexdigest()}`.\n\n"
        "| Stage | Result |\n| --- | --- |\n"
        "| Express 4.21.2 | Three tests pass |\n"
        "| Express 5.1.0, unchanged route | Application fails to load: Missing parameter name |\n"
        "| Express 5.1.0, /*splat | Two pass; root route returns 404 |\n"
        "| Express 5.1.0, /{*splat} | Three tests pass |\n\n"
        "Express 5 requires a named wildcard. Braces include the root path. "
        "The API route remains registered before the fallback.\n\n"
        f"```diff\n{patch}```\n\n"
        "The dependency and lockfile also change. The target lockfile and "
        "individual test logs are saved alongside this report.\n\n"
        "Unresolved work: this fixture covers routing only. Review other "
        "migration changes and dependency advisories for a production application.\n",
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
        + "Local execution only. Console, sbx, GitHub delivery, and email are unverified.\n"
    )
    print(f"All local checks passed. Evidence and bundle: {OUTPUT}")


if __name__ == "__main__":
    main()
