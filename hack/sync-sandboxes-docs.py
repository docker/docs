#!/usr/bin/env python3
"""Import the pinned public cookbook export. See hack/sandboxes/README.md."""

import io
import json
from pathlib import Path
import subprocess
import tarfile


ROOT = Path(__file__).resolve().parent.parent
MANIFEST = ROOT / "hack/sandboxes/source.json"
GUIDES = ROOT / "content/manuals/ai/sandboxes-api/v1/cookbook"
REFERENCE = ROOT / "content/reference/api/sandboxes"


def prepare_guide(content, name):
    """Preserve published recipe URLs and point installation links at v1."""
    text = content.decode()
    alias = f"/ai/sandboxes-api/cookbook/{Path(name).stem}/"
    if "\naliases:\n" in text:
        text = text.replace("\naliases:\n", f"\naliases:\n  - {alias}\n", 1)
    else:
        text = text.replace("\n---\n", f"\naliases:\n  - {alias}\n---\n", 1)
    text = text.replace(
        "https://docs.docker.com/ai/sandboxes-api/install/",
        "/manuals/ai/sandboxes-api/v1/install.md",
    )
    return text.encode()


def main():
    source = json.loads(MANIFEST.read_text())
    archive = subprocess.check_output([
        "gh", "api",
        f"repos/{source['repository']}/tarball/{source['generatedCommit']}",
    ])
    with tarfile.open(fileobj=io.BytesIO(archive), mode="r:gz") as tar:
        files = {}
        for member in tar.getmembers():
            if member.isfile():
                path = member.name.split("/", 1)[1]
                files[path] = tar.extractfile(member).read()

    if files["SOURCE"].decode().strip() != source["sourceCommit"]:
        raise SystemExit("The generated export does not match the recorded source commit")

    # Prepare every output before changing the working tree. Fail if an upstream
    # page disappears instead of silently dropping a published recipe.
    outputs = {}
    for name in source["guides"]:
        if Path(name).name != name or not name.endswith(".md"):
            raise SystemExit(f"Invalid guide filename: {name}")
        outputs[GUIDES / name] = prepare_guide(
            files[f"cookbook/outputs/guides/{name}"], name
        )

    outputs[REFERENCE / "api.yaml"] = files["cookbook/outputs/api-reference/api.yaml"]

    for destination, content in outputs.items():
        destination.parent.mkdir(parents=True, exist_ok=True)
        destination.write_bytes(content)
    print(f"Imported {len(source['guides'])} recipes and the public OpenAPI specification")


if __name__ == "__main__":
    main()
