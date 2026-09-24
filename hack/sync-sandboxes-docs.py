#!/usr/bin/env python3
"""Import the pinned public cookbook export. See hack/sandboxes/README.md."""

import io
import json
from pathlib import Path
import subprocess
import tarfile


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

    reference = files["cookbook/outputs/api-reference/index.md"].decode()
    reference = reference.replace(
        "title: Cloud Sandboxes API reference",
        "title: Docker Sandboxes API reference\nlayout: api-reference-generated",
    ).replace(
        "linkTitle: API reference", "linkTitle: Sandboxes"
    ).replace(
        "    group: APIs and SDKs",
        "    badge:\n      color: violet\n      text: Experimental",
    ).replace(
        "](../_index.md)", "](/manuals/ai/sandboxes-api/_index.md)"
    ).replace(
        "](../connect-to-cloud-with-a-bearer-token.md)",
        "](/manuals/ai/sandboxes-api/authentication.md)",
    )
    front, body = reference.split("---\n", 2)[1:]
    notice = (
        "> [!NOTE]\n"
        "> The Docker Sandboxes API and SDK are experimental. Features, interfaces,\n"
        "> and behavior may change.\n\n"
    )
    outputs[REFERENCE / "index.md"] = (
        "---\n" + front + "---\n\n" + notice + body.lstrip()
    ).encode()
    # Preserve the published spec byte for byte, including streaming contracts.
    outputs[REFERENCE / "api.yaml"] = files["cookbook/outputs/api-reference/api.yaml"]

    for destination, content in outputs.items():
        destination.parent.mkdir(parents=True, exist_ok=True)
        destination.write_bytes(content)
    print(f"Imported {len(source['guides'])} recipes and the public API reference")


if __name__ == "__main__":
    main()
