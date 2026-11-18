#!/usr/bin/env python3
"""Build the Docker Sandboxes capability table from a pinned public catalog."""

import argparse
import base64
import hashlib
import json
import re
import subprocess
import sys
from pathlib import Path


ROOT = Path(__file__).resolve().parents[2]
DIR = Path(__file__).resolve().parent
SOURCE = DIR / "source.json"
SCHEMA = DIR / "kit.schema.json"
SUPPORT = DIR / "support.json"
OUTPUT = ROOT / "data/kit_capabilities.json"
REPOSITORY = "docker/sandbox-kit-spec"
TYPE = re.compile(r"^com\.docker\.sandbox/[a-z0-9][a-z0-9-]*@[1-9][0-9]*$")
COMMIT = re.compile(r"^[0-9a-f]{40}$")
VERSION = re.compile(r"^v[0-9]+\.[0-9]+\.[0-9]+$")
STATUSES = {"supported", "limited", "declaration-only", "unsupported", "unverified"}
DOC_PREFIX = "docs/spec/capabilities/"


class CatalogError(ValueError):
    pass


def read_json(path):
    try:
        return json.loads(path.read_text(encoding="utf-8"))
    except (OSError, json.JSONDecodeError) as error:
        raise CatalogError(f"{path}: {error}") from error


def json_bytes(value):
    return (json.dumps(value, indent=2, ensure_ascii=False) + "\n").encode()


def inventory(schema_bytes, documentation):
    try:
        schema = json.loads(schema_bytes)
        branches = schema["definitions"]["capability"]["allOf"]
    except (ValueError, KeyError, TypeError) as error:
        raise CatalogError(f"invalid capability schema: {error}") from error

    found = {}
    for branch in branches:
        try:
            identifier = branch["if"]["properties"]["type"]["const"]
            reference = branch["then"]["properties"]["config"]["$ref"]
        except (KeyError, TypeError) as error:
            raise CatalogError(f"invalid capability branch: {error}") from error
        if not isinstance(identifier, str) or not TYPE.fullmatch(identifier):
            raise CatalogError(f"invalid capability ID: {identifier!r}")
        if identifier in found:
            raise CatalogError(f"duplicate capability ID: {identifier}")
        expected = f"capabilities/{identifier}.schema.json"
        if reference != expected:
            raise CatalogError(f"{identifier}: expected schema reference {expected}, got {reference}")
        found[identifier] = f"{DOC_PREFIX}{identifier}.md"

    if not found:
        raise CatalogError("schema has no capability types")
    if not isinstance(documentation, list) or documentation != sorted(set(documentation)):
        raise CatalogError("documentation paths must be a sorted list without duplicates")
    expected_docs = set(found.values())
    actual_docs = set(documentation)
    if expected_docs != actual_docs:
        raise CatalogError(
            "spec documentation mismatch: missing " + str(sorted(expected_docs - actual_docs))
            + "; extra " + str(sorted(actual_docs - expected_docs))
        )
    return found


def validate_support(support, identifiers):
    if not isinstance(support, dict):
        raise CatalogError("support.json must be an object")
    missing = sorted(identifiers - support.keys())
    extra = sorted(support.keys() - identifiers)
    if missing or extra:
        raise CatalogError(f"support.json classifications: missing {missing}; extra {extra}")
    for identifier, entry in support.items():
        if not isinstance(entry, dict) or set(entry) != {"local", "cloud", "restrictions"}:
            raise CatalogError(f"{identifier}: expected local, cloud, and restrictions")
        restrictions = entry["restrictions"]
        if not isinstance(restrictions, str) or not restrictions.strip():
            raise CatalogError(f"{identifier}: restrictions must be nonempty text")
        if any(character in restrictions for character in "|\r\n"):
            raise CatalogError(f"{identifier}: restrictions must fit in one Markdown table cell")
        for environment in ("local", "cloud"):
            value = entry[environment]
            if not isinstance(value, dict) or value.get("status") not in STATUSES:
                raise CatalogError(f"{identifier} {environment}: invalid status")
            status = value["status"]
            expected_keys = {"status", "since"} if status in {"supported", "limited", "declaration-only"} else {"status"}
            if set(value) != expected_keys:
                raise CatalogError(f"{identifier} {environment}: expected {sorted(expected_keys)}")
            if "since" in value and (not isinstance(value["since"], str) or not VERSION.fullmatch(value["since"])):
                raise CatalogError(f"{identifier} {environment}: invalid sbx version")


def make_data(schema_bytes, source, support):
    if not isinstance(source, dict) or set(source) != {"repository", "commit", "schema_sha256", "documentation"}:
        raise CatalogError("source.json has unexpected fields")
    if source["repository"] != REPOSITORY or not isinstance(source["commit"], str) or not COMMIT.fullmatch(source["commit"]):
        raise CatalogError("source.json must pin docker/sandbox-kit-spec at a 40-character commit SHA")
    checksum = hashlib.sha256(schema_bytes).hexdigest()
    if source["schema_sha256"] != checksum:
        raise CatalogError(f"schema checksum mismatch: expected {source['schema_sha256']}, got {checksum}")
    docs = inventory(schema_bytes, source["documentation"])
    validate_support(support, set(docs))
    return {
        "source": {"repository": REPOSITORY, "commit": source["commit"]},
        "capabilities": [
            {
                "id": identifier,
                "url": f"https://github.com/{REPOSITORY}/blob/{source['commit']}/{docs[identifier]}",
                **support[identifier],
            }
            for identifier in sorted(docs)
        ],
    }


def github_api(path):
    try:
        return json.loads(subprocess.check_output(["gh", "api", f"repos/{REPOSITORY}/{path}"], text=True))
    except (OSError, subprocess.CalledProcessError, json.JSONDecodeError) as error:
        raise CatalogError(f"GitHub API request failed for {path}: {error}") from error


def sync(commit, support):
    if not COMMIT.fullmatch(commit):
        raise CatalogError("--ref must be an immutable 40-character lowercase commit SHA")
    content = github_api(f"contents/schema/kit.schema.json?ref={commit}")
    if content.get("encoding") != "base64":
        raise CatalogError("GitHub API did not return base64 schema content")
    try:
        schema_bytes = base64.b64decode(content["content"].replace("\n", ""), validate=True)
    except (KeyError, ValueError) as error:
        raise CatalogError(f"invalid GitHub schema content: {error}") from error
    tree = github_api(f"git/trees/{commit}?recursive=1")
    if tree.get("truncated"):
        raise CatalogError("GitHub tree response was truncated")
    documentation = sorted(
        node["path"] for node in tree["tree"]
        if node["type"] == "blob"
        and node["path"].startswith(DOC_PREFIX + "com.docker.sandbox/")
        and node["path"].endswith(".md")
    )
    source = {
        "repository": REPOSITORY,
        "commit": commit,
        "schema_sha256": hashlib.sha256(schema_bytes).hexdigest(),
        "documentation": documentation,
    }
    data = make_data(schema_bytes, source, support)
    # Validate every file and every classification before changing the checkout.
    SCHEMA.write_bytes(schema_bytes)
    SOURCE.write_bytes(json_bytes(source))
    OUTPUT.write_bytes(json_bytes(data))


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    commands = parser.add_subparsers(dest="command", required=True)
    commands.add_parser("check")
    commands.add_parser("generate")
    sync_parser = commands.add_parser("sync")
    sync_parser.add_argument("--ref", required=True, help="full public spec commit SHA")
    args = parser.parse_args()
    try:
        support = read_json(SUPPORT)
        if args.command == "sync":
            sync(args.ref, support)
        else:
            source = read_json(SOURCE)
            try:
                schema_bytes = SCHEMA.read_bytes()
            except OSError as error:
                raise CatalogError(f"{SCHEMA}: {error}") from error
            expected = json_bytes(make_data(schema_bytes, source, support))
            if args.command == "generate":
                OUTPUT.write_bytes(expected)
            elif not OUTPUT.exists() or OUTPUT.read_bytes() != expected:
                raise CatalogError(f"{OUTPUT}: missing or stale; run python3 hack/kit-capabilities/sync.py generate")
    except CatalogError as error:
        parser.exit(1, f"{error}\n")


if __name__ == "__main__":
    main()
