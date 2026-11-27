#!/usr/bin/env python3
"""Sync a verified published sandbox-kit-spec catalog.

Run with GH_TOKEN and optional KIT_SPEC_VERSION. Requires gh and permission
for GITHUB_TOKEN to create pull requests in the docs repository.
"""

import hashlib
import json
import os
from pathlib import Path
import re
import subprocess
import tempfile


REPOSITORY = "docker/sandbox-kit-spec"
CATALOG = "data/sandbox_kit_capabilities.json"
SUPPORT = "data/sandbox_kit_support.json"
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


def validate_data(catalog, support, version):
    digest = re.compile(r"[0-9a-f]{64}")
    ident = re.compile(r"com\.docker\.sandbox/[a-z0-9][a-z0-9-]*@[1-9][0-9]*")
    if type(catalog.get("format_version")) is not int or catalog["format_version"] != 1 or catalog.get("schema_version") != "3":
        raise ValueError("Unsupported catalog format or schema")
    for key in ("shared_spec_sha256", "conformance_sha256"):
        if not isinstance(catalog.get(key), str) or not digest.fullmatch(catalog[key]):
            raise ValueError(f"Invalid {key}")
    entries = catalog.get("capabilities")
    if not isinstance(entries, list) or not entries:
        raise ValueError("Empty capability catalog")
    ids = []
    base = f"https://github.com/{REPOSITORY}/blob/{version}/"
    for entry in entries:
        if not isinstance(entry, dict):
            raise ValueError("Invalid capability entry")
        name = entry.get("id")
        if not isinstance(name, str) or not ident.fullmatch(name) or (ids and name <= ids[-1]):
            raise ValueError("Invalid, duplicate or unsorted capability ID")
        ids.append(name)
        if not isinstance(entry.get("description"), str) or not entry["description"].strip():
            raise ValueError("Missing capability description")
        if entry.get("spec_url") != base + "docs/spec/capabilities/" + name + ".md" or entry.get("schema_url") != base + "schema/capabilities/" + name + ".schema.json":
            raise ValueError("Noncanonical capability URL")
        for key in ("schema_sha256", "schema_contract_sha256", "spec_sha256"):
            if not isinstance(entry.get(key), str) or not digest.fullmatch(entry[key]):
                raise ValueError(f"Invalid {key}")
    if not isinstance(support, dict) or type(support.get("format_version")) is not int or support["format_version"] != 1 or not isinstance(support.get("capabilities"), list):
        raise ValueError("Invalid support data")
    seen = set()
    for entry in support["capabilities"]:
        if not isinstance(entry, dict) or not isinstance(entry.get("id"), str) or not ident.fullmatch(entry["id"]) or entry["id"] in seen:
            raise ValueError("Invalid or duplicate support ID")
        seen.add(entry["id"])
        if set(entry) != {"id", "local", "cloud", "restrictions"}:
            raise ValueError("Invalid support fields")
        status = r"(?:Not verified|Unsupported|Declaration only|(?:Limited: )?v[0-9]+\.[0-9]+\.[0-9]+)"
        for key in ("local", "cloud"):
            if not isinstance(entry[key], str) or not re.fullmatch(status, entry[key]):
                raise ValueError("Invalid support status")
        if not isinstance(entry["restrictions"], str) or not entry["restrictions"].strip() or "\n" in entry["restrictions"] or "\r" in entry["restrictions"]:
            raise ValueError("Invalid support restrictions")
    if seen - set(ids):
        raise ValueError("Catalog omits known support entries")


def write_catalog(path, data):
    if path.exists() and path.read_bytes() == data:
        return
    with tempfile.NamedTemporaryFile(dir=path.parent, prefix=".kit-catalog-", delete=False) as output:
        name = output.name
        try:
            output.write(data)
            output.flush()
            os.fsync(output.fileno())
        except BaseException:
            os.unlink(name)
            raise
    try:
        os.replace(name, path)
    finally:
        if os.path.exists(name):
            os.unlink(name)


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
        data = catalog.read_bytes()
        parsed = json.loads(data, object_pairs_hook=unique_keys, parse_constant=invalid_constant)
        support = json.loads((ROOT / SUPPORT).read_bytes(), object_pairs_hook=unique_keys, parse_constant=invalid_constant)
        validate_data(parsed, support, version)
        write_catalog(ROOT / CATALOG, data)
    report(
        f"Synced specification contracts from [{version}](https://github.com/{REPOSITORY}/releases/tag/{version}). "
        "Review the matching Docker Sandboxes release before marking the PR ready.",
        available="true", version=version,
    )


if __name__ == "__main__":
    main()
