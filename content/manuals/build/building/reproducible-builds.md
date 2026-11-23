---
title: Reproducible builds
description: Pin build inputs and timestamps to make Docker builds reproducible
keywords: build, reproducibility, source policy, SOURCE_DATE_EPOCH, timestamps
weight: 65
---

Reproducing an image requires consistent inputs and build output. Pin your
[base images to digests](best-practices.md#pin-base-image-versions), use fixed
dependency versions, and make sure your build commands produce deterministic
output. BuildKit can also pin remote downloads and normalize timestamps.

## Pin an HTTP download to a registry blob

A source policy can replace an HTTP download with a single blob stored in a
container registry. Use a `docker-image+blob://` identifier with the blob's
digest. BuildKit downloads that blob as a file, without extracting it as an
image layer.

For example, after uploading an archive to a registry, use this rule to keep
an existing Dockerfile's HTTP `ADD` instruction while fetching the archived
copy. Save the policy as `policy.json`:

```json
{
  "rules": [
    {
      "action": "CONVERT",
      "selector": {
        "identifier": "https://example.com/source.tar.gz"
      },
      "updates": {
        "identifier": "docker-image+blob://registry.example.com/archive/source@sha256:<BLOB_DIGEST>",
        "attrs": {
          "http.filename": "source.tar.gz"
        }
      }
    }
  ]
}
```

Replace `<BLOB_DIGEST>` with the digest of the uploaded file's blob, not an
image manifest digest. The blob must already exist in the registry repository.

Apply the policy to your build:

```console
$ EXPERIMENTAL_BUILDKIT_SOURCE_POLICY=./policy.json docker buildx build .
```
