#!/usr/bin/env python3
"""Sync a verified published sandbox-kit-spec catalog using its released updater.

Run with GH_TOKEN and optional KIT_SPEC_VERSION. Requires gh, Go, and permission
for GITHUB_TOKEN to create pull requests in the docs repository.
"""

import hashlib
import json
import os
from pathlib import Path, PurePosixPath
import re
import subprocess
import tarfile
import tempfile


REPOSITORY = "docker/sandbox-kit-spec"
MODULE = f"github.com/{REPOSITORY}/v3"
PAGE = "content/manuals/ai/sandboxes/customize/author/capabilities.md"
ROOT = Path(__file__).resolve().parent.parent
NUMBER = r"(?:0|[1-9][0-9]*)"
TAG = re.compile(rf"v3\.{NUMBER}\.{NUMBER}(?:-([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?")


def run(*args, **kwargs):
    return subprocess.run(args, check=True, **kwargs)


def gh_json(*args):
    return json.loads(run("gh", *args, capture_output=True, text=True).stdout)


def validate_tag(tag):
    match = TAG.fullmatch(tag)
    if not match or any(
        part.isdigit() and len(part) > 1 and part.startswith("0")
        for part in (match.group(1) or "").split(".")
    ):
        raise ValueError(f"Not a canonical v3 module release tag: {tag!r}")
    return tag


def select_release(version):
    if not version:
        # Include milestone/prerelease tags; the latest-release API excludes them.
        releases = gh_json(
            "release", "list", "--repo", REPOSITORY, "--exclude-drafts",
            "--limit", "1", "--json", "tagName",
        )
        if not releases:
            raise ValueError("No published sandbox-kit-spec releases")
        version = releases[0]["tagName"]
    validate_tag(version)
    release = gh_json(
        "release", "view", version, "--repo", REPOSITORY,
        "--json", "tagName,isDraft,publishedAt,assets",
    )
    if release["tagName"] != version or release["isDraft"] or not release["publishedAt"]:
        raise ValueError("The requested tag must have a published release")
    return release


def invalid_constant(value):
    raise ValueError(f"Invalid JSON constant: {value}")


def unique_keys(pairs):
    result = {}
    for key, value in pairs:
        if key in result:
            raise ValueError(f"Duplicate JSON key: {key}")
        result[key] = value
    return result


def verify_catalog(directory, version, revision):
    catalog_path = directory / "capabilities.json"
    matches = []
    for line in (directory / "checksums.txt").read_text().splitlines():
        fields = line.split()
        if len(fields) >= 2 and fields[1].lstrip("*") == "capabilities.json":
            if len(fields) != 2:
                raise ValueError("Malformed capability catalog checksum")
            matches.append(fields[0])
    if len(matches) != 1 or not re.fullmatch(r"[0-9a-f]{64}", matches[0]):
        raise ValueError("Expected exactly one SHA-256 checksum for capabilities.json")
    if hashlib.sha256(catalog_path.read_bytes()).hexdigest() != matches[0]:
        raise ValueError("Capability catalog checksum mismatch")
    catalog = json.loads(
        catalog_path.read_text(), object_pairs_hook=unique_keys,
        parse_constant=invalid_constant,
    )
    if not isinstance(catalog, dict) or catalog.get("version") != version or catalog.get("revision") != revision:
        raise ValueError("Capability catalog version/revision does not match the release commit")
    return catalog_path


def extract_source(archive_path, destination):
    """Accept one GitHub archive root containing only directories and files."""
    with tarfile.open(archive_path, "r:gz") as archive:
        members = archive.getmembers()
        paths = set()
        roots = set()
        for member in members:
            path = PurePosixPath(member.name)
            if (
                not path.parts or path.is_absolute() or ".." in path.parts
                or "\\" in member.name or path in paths
                or not (member.isdir() or member.isfile())
            ):
                raise ValueError(f"Unsafe source archive member: {member.name!r}")
            paths.add(path)
            roots.add(path.parts[0])
        if len(roots) != 1:
            raise ValueError("Expected one source archive root")
        root = PurePosixPath(roots.pop())
        for required in ("go.mod", "cmd/capability-catalog/main.go"):
            if root / required not in paths:
                raise ValueError(f"Source archive is missing {required}")
        destination.mkdir()
        archive.extractall(destination, filter="data")
    source = destination / root
    if not re.search(
        rf"^module[ \t]+{re.escape(MODULE)}[ \t]*$",
        (source / "go.mod").read_text(), re.MULTILINE,
    ):
        raise ValueError("Source archive has an unexpected Go module")
    return source


def report(message, **outputs):
    print(message)
    if path := os.environ.get("GITHUB_STEP_SUMMARY"):
        with open(path, "a") as summary:
            summary.write(message + "\n")
    if path := os.environ.get("GITHUB_OUTPUT"):
        with open(path, "a") as output:
            for key, value in outputs.items():
                output.write(f"{key}={value}\n")


def main():
    release = select_release(os.environ.get("KIT_SPEC_VERSION", ""))
    version = release["tagName"]
    assets = [asset["name"] for asset in release["assets"]]
    if "capabilities.json" not in assets:
        report(f"Skipped {version}: this release has no capability catalog.", available="false")
        return
    if assets.count("capabilities.json") != 1 or assets.count("checksums.txt") != 1:
        raise ValueError("Release must contain exactly one catalog and checksum asset")
    revision = gh_json("api", f"repos/{REPOSITORY}/commits/{version}")["sha"]
    if not re.fullmatch(r"[0-9a-f]{40}", revision):
        raise ValueError("Release commit must have a full SHA-1 revision")
    with tempfile.TemporaryDirectory(prefix="kit-capabilities-") as temporary:
        directory = Path(temporary)
        run(
            "gh", "release", "download", version, "--repo", REPOSITORY,
            "--pattern", "capabilities.json", "--pattern", "checksums.txt",
            "--dir", str(directory),
        )
        catalog = verify_catalog(directory, version, revision)
        # Version-suffixed go run rejects the released module's replace directives.
        # Fetch the exact verified commit and run it as the main module instead.
        archive = directory / "source.tar.gz"
        with archive.open("wb") as output:
            run("gh", "api", f"repos/{REPOSITORY}/tarball/{revision}", stdout=output)
        source = extract_source(archive, directory / "source")
        env = dict(os.environ)
        for name in ("GH_TOKEN", "GITHUB_TOKEN"):
            env.pop(name, None)
        env["GOWORK"] = "off"
        # Preserve the released module, including its dependency replacements.
        # Source, assets, and Go module caches stay outside the docs tree.
        run(
            "go", "run", "./cmd/capability-catalog", "docs",
            "--catalog", str(catalog), "--page", str(ROOT / PAGE),
            cwd=source, env=env,
        )
    report(
        f"Synced specification contracts from [{version}](https://github.com/{REPOSITORY}/releases/tag/{version}). "
        "Review the matching Docker Sandboxes release before marking the PR ready.",
        available="true", version=version,
    )


if __name__ == "__main__":
    main()
