#!/usr/bin/env python3
"""Import the pinned public cookbook export. See hack/sandboxes/README.md."""

import io
import json
from pathlib import Path
import subprocess
import tarfile
import tempfile


ROOT = Path(__file__).resolve().parent.parent
MANIFEST = ROOT / "hack/sandboxes/source.json"
GUIDES = ROOT / "content/manuals/ai/sandboxes-api/cookbook"
REFERENCE = ROOT / "content/reference/api/sandboxes"


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
        page = files[f"cookbook/outputs/guides/{name}"].decode()
        page = page.replace(
            "](install-the-sdks.md)", "](../sdks.md)"
        )
        outputs[GUIDES / name] = page.encode()

    specification = files["cookbook/outputs/api-reference/api.yaml"]
    # Temporary export correction, generated from the owning upstream source fix.
    # Apply before writing any outputs, and fail if the pinned export has drifted.
    if source.get("apiPatch"):
        patch = MANIFEST.parent / source["apiPatch"]
        with tempfile.TemporaryDirectory() as temporary:
            target = Path(temporary) / "api.yaml"
            target.write_bytes(specification)
            subprocess.run(
                ["git", "apply", "--whitespace=nowarn", str(patch)],
                cwd=temporary, check=True,
            )
            specification = target.read_bytes()
    outputs[REFERENCE / "api.yaml"] = specification

    for destination, content in outputs.items():
        destination.parent.mkdir(parents=True, exist_ok=True)
        destination.write_bytes(content)
    print(f"Imported {len(source['guides'])} recipes and the public OpenAPI specification")


if __name__ == "__main__":
    main()
