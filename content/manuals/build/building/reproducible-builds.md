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

## Set build timestamps

The `SOURCE_DATE_EPOCH` build argument sets timestamps in image metadata and
in files exported with the `local` and `tar` exporters. Pass a Unix timestamp
in seconds:

```console
$ docker buildx build --build-arg SOURCE_DATE_EPOCH=1704067200 .
```

To also normalize file timestamps inside image layers, set the
[`rewrite-timestamp` exporter option](../exporters/image-registry.md):

```console
$ docker buildx build --build-arg SOURCE_DATE_EPOCH=1704067200 \
    --output type=image,name=<registry>/<image>,push=true,rewrite-timestamp=true .
```

### Define a default in the Dockerfile

You can define `SOURCE_DATE_EPOCH` as a global argument, before the first
`FROM` instruction:

```dockerfile
# syntax=docker/dockerfile:1
ARG SOURCE_DATE_EPOCH=1704067200
FROM alpine
```

An explicit caller-provided value takes precedence over the Dockerfile default.

### Derive the timestamp from the build context

Set `SOURCE_DATE_EPOCH=context` to derive a timestamp from the main build
context:

- Git contexts use the commit timestamp.
- HTTP contexts use the `Last-Modified` header. For an archive without that
  header, BuildKit uses the newest archive entry's modification time.
- Local contexts leave `SOURCE_DATE_EPOCH` unset, including checked-out Git
  repositories passed as local directories.

Pass `--build-arg SOURCE_DATE_EPOCH=context` or define
`ARG SOURCE_DATE_EPOCH=context` before the first `FROM` instruction. For local
contexts, pass a numeric timestamp instead.

For CI examples, see
[Reproducible builds with GitHub Actions](../ci/github-actions/reproducible-builds.md).
