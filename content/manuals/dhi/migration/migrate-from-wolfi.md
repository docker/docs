---
title: Migrate from Wolfi
description: Step-by-step guide to migrate from Wolfi distribution images to Docker Hardened Images
weight: 30
keywords: wolfi, chainguard, migration, dhi
---

This guide helps you migrate from Wolfi-based images to Docker Hardened
Images (DHI). Generally, the migration process is straightforward since Wolfi is
Alpine-like and DHI provides an Alpine-based hardened image.

Like other hardened images, DHI provides comprehensive
[attestations](/dhi/core-concepts/attestations/) including SBOMs and provenance,
allowing you to [verify](/manuals/dhi/how-to/verify.md) image signatures and
[scan](/manuals/dhi/how-to/scan.md) for vulnerabilities to ensure the security
and integrity of your images.

## Migration steps

The following example demonstrates how to migrate a Dockerfile from a
Wolfi-based image to an Alpine-based Docker Hardened Image.

### Step 1: Update the base image in your Dockerfile

Update the base image in your application's Dockerfile to a hardened image. This
is typically going to be an image tagged as `dev` or `sdk` because it has the tools
needed to install packages and dependencies.

The following example diff snippet from a Dockerfile shows the old base image
replaced by the new hardened image.

> [!NOTE]
>
> You must authenticate to `dhi.io` before you can pull Docker Hardened Images.
> Use your Docker ID credentials (the same username and password you use for
> Docker Hub). If you don't have a Docker account, [create
> one](../../accounts/create-account.md) for free.
>
> Run `docker login dhi.io` to authenticate.

```diff
- ## Original base image
- FROM cgr.dev/chainguard/go:latest-dev

+ ## Updated to use hardened base image
+ FROM dhi.io/golang:1.25-alpine3.22-dev
```

Note that DHI does not have a `latest` tag in order to promote best practices
around image versioning. Ensure that you specify the appropriate version tag for your image. To find the right tag, explore the available tags in the [DHI Catalog](https://hub.docker.com/hardened-images/catalog/).

### Step 2: Update the runtime image in your Dockerfile

> [!NOTE]
>
> Multi-stage builds are recommended to keep your final image minimal and
> secure. Single-stage builds are supported, but they include the full `dev` image
> and therefore result in a larger image with a broader attack surface.

To ensure that your final image is as minimal as possible, you should use a
[multi-stage build](/manuals/build/building/multi-stage.md). All stages in your
Dockerfile should use a hardened image. While intermediary stages will typically
use images tagged as `dev` or `sdk`, your final runtime stage should use a runtime image.

Utilize the build stage to compile your application and copy the resulting
artifacts to the final runtime stage. This ensures that your final image is
minimal and secure.

The following example shows a multi-stage Dockerfile with a build stage and runtime stage:

```dockerfile
# Build stage
FROM dhi.io/golang:1.25-alpine3.22-dev AS builder
WORKDIR /app
COPY . .
RUN go build -o myapp

# Runtime stage
FROM dhi.io/golang:1.25-alpine3.22
WORKDIR /app
COPY --from=builder /app/myapp .
ENTRYPOINT ["/app/myapp"]
```

After updating your Dockerfile, build and test your application. If you encounter
issues, see the [Troubleshoot](/manuals/dhi/troubleshoot.md) guide for common
problems and solutions.

## Language-specific examples

See the examples section for language-specific migration examples:

- [Go](examples/go.md)
- [Python](examples/python.md)
- [Node.js](examples/node.md)
