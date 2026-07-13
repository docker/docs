---
title: Image digests
description: Learn how Docker Hardened Images help secure every stage of your software supply chain with signed metadata, provenance, and minimal attack surface.
keywords: docker image digest, pull image by digest, immutable container image, secure container reference, multi-platform manifest
---

## What are Docker image digests?

A Docker image digest is a unique, cryptographic identifier (SHA-256 hash)
representing the content of a Docker image. Unlike tags, which can be reused or
changed, a digest is immutable and ensures that the exact same image is pulled
every time. This guarantees consistency across different environments and
deployments.

For example, the digest for the `nginx:latest` image might look like:

```text
sha256:94a00394bc5a8ef503fb59db0a7d0ae9e1110866e8aee8ba40cd864cea69ea1a
```

This digest uniquely identifies the specific version of the `nginx:latest` image,
ensuring that any changes to the image content result in a different digest.

## Why are image digests important?

Using image digests instead of tags offers several advantages:

- Immutability: Once an image is built and its digest is generated, the content
  tied to that digest cannot change. This means that if you pull an image using
  its digest, you can be confident that you are retrieving exactly the same
  image that was originally built.

- Security: Digests help prevent supply chain attacks by ensuring that the image
  content has not been tampered with. Even a small change in the image content
  will result in a completely different digest.

- Consistency: Using digests ensures that the same image is used across
  different environments, reducing the risk of discrepancies between
  development, staging, and production environments.

## Docker Hardened Image digests

By using image digests to reference DHIs, you can ensure that your applications are
always using the exact same secure image version, enhancing security and
compliance

## View an image digest

### Use the Docker CLI

To view the image digest of a Docker image, you can use the following command. Replace
`<image-name>:<tag>` with the image name and tag.

```console
$ docker buildx imagetools inspect <image-name>:<tag>
```

### Use the Docker Hub UI

1. Go to [Docker Hub](https://hub.docker.com/) and sign in.
2. Navigate to your organization's namespace and open the mirrored DHI repository.
3. Select the **Tags** tab to view image variants.
4. Each tag in the list includes a **Digest** field showing the image's SHA-256 value.

## Pull an image by digest

Pulling an image by digest ensures that you are pulling the exact image version
identified by the specified digest.

To pull a Docker image using its digest, use the following command. Replace
`<image-name>` with the image name and `<digest>` with the image digest.

```console
$ docker pull <image-name>@sha256:<digest>
```

For example, to pull a `docs/dhi-python:3.13` image using its digest of
`94a00394bc5a8ef503fb59db0a7d0ae9e1110866e8aee8ba40cd864cea69ea1a`, you would
run:

```console
$ docker pull docs/dhi-python@sha256:94a00394bc5a8ef503fb59db0a7d0ae9e1110866e8aee8ba40cd864cea69ea1a
```

## Multi-platform images and manifests

Docker Hardened Images are published as multi-platform images, which means
a single image tag (like `docs/dhi-python:3.13`) can support multiple operating
systems and CPU architectures, such as `linux/amd64`, `linux/arm64`, and more.

Instead of pointing to a single image, a multi-platform tag points to a manifest
list (also called an index), which is a higher-level object that references
multiple image digests, one for each supported platform.

When you inspect a multi-platform image using `docker buildx imagetools inspect`, you'll see something like this:

```text
Name:      docs/dhi-python:3.13
MediaType: application/vnd.docker.distribution.manifest.list.v2+json
Digest:    sha256:6e05...d231

Manifests:
  Name:        docs/dhi-python:3.13@sha256:94a0...ea1a
  Platform:    linux/amd64
  ...

  Name:        docs/dhi-python:3.13@sha256:7f1d...bc43
  Platform:    linux/arm64
  ...
```

- The manifest list digest (`sha256:6e05...d231`) identifies the overall
  multi-platform image.
- Each platform-specific image has its own digest (e.g., `sha256:94a0...ea1a`
  for `linux/amd64`).

### Why this matters

- Reproducibility: If you're building or running containers on different
  architectures, using a tag alone will resolve to the appropriate image digest
  for your platform.
- Verification: You can pull and verify a specific image digest for your
  platform to ensure you're using the exact image version, not just the manifest
  list.
- Policy enforcement: When enforcing digest-based policies with Docker Scout,
  each platform variant is evaluated individually using its digest.