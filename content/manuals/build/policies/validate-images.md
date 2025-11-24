---
title: Validating image inputs
linkTitle: Image validation
description: Write policies to validate container images used in your builds
weight: 20
---

Container images are the most common build inputs. Every `FROM` instruction
pulls an image, and `COPY --from` references pull additional images. Validating
these images protects your build supply chain from compromised registries,
unexpected updates, and unauthorized base images.

This guide teaches you to write policies that validate image inputs,
progressing from basic allowlisting to advanced attestation checks.

## Prerequisites

You should understand the policy basics from the [Introduction](./intro.md):
creating policy files, basic Rego syntax, and how policies evaluate during
builds.

## What are image inputs?

Image inputs come from two Dockerfile instructions:

```dockerfile
# FROM instructions
FROM alpine:3.22
FROM golang:1.25-alpine AS builder

# COPY --from references
COPY --from=builder /app /app
COPY --from=nginx:latest /etc/nginx/nginx.conf /nginx.conf
```

Each of these references triggers a policy evaluation. Your policy can inspect
image metadata, verify attestations, and enforce constraints before the build
proceeds.

## Allowlist specific repositories

The simplest image policy restricts which repositories can be used. This
prevents developers from using arbitrary images that haven't been vetted.

Create a policy that only allows Alpine:

```rego {title="Dockerfile.rego"}
package docker

default allow := false

allow if input.local

allow if {
  input.image.repo == "alpine"
}

decision := {"allow": allow}
```

This policy:

- Denies all inputs by default
- Allows local build context
- Allows any image from the `alpine` repository (any tag or digest)

Test it with a Dockerfile:

```dockerfile {title="Dockerfile"}
FROM alpine
RUN echo "hello"
```

```console
$ docker build .
```

The build succeeds. Try changing to `FROM ubuntu`:

```console
$ docker build .
```

The build fails because `ubuntu` doesn't match the allowed repository.

## Require digest references

Tags like `alpine:3.22` can change - someone could push a new image with the
same tag. Digests like `alpine@sha256:abc123...` are immutable. Requiring
digests ensures builds are reproducible.

Update your policy to require digests:

```rego
package docker

default allow := false

allow if input.local

allow if {
  input.image.isCanonical
}

decision := {"allow": allow}
```

The `isCanonical` field is `true` when the image reference includes a digest.
Now update your Dockerfile:

```dockerfile
FROM alpine@sha256:4b7ce07002c69e8f3d704a9c5d6fd3053be500b7f1c69fc0d80990c2ad8dd412
RUN echo "hello"
```

The build succeeds with a digest reference. Without the digest, it fails.

You can also validate specific digests:

```rego
allow if {
  input.image.repo == "alpine"
  input.image.checksum == "sha256:4b7ce07002c69e8f3d704a9c5d6fd3053be500b7f1c69fc0d80990c2ad8dd412"
}

decision := {"allow": allow}
```

This pins the exact image content, useful for critical base images.

## Restrict registries

Control which registries your builds can pull from. This helps enforce
corporate policies or restrict to trusted sources.

```rego
package docker

default allow := false

allow if input.local

# Allow Docker Hub images
allow if {
  input.image.host == "docker.io"  # Docker Hub
  input.image.repo == "alpine"
}

# Allow images from internal registry
allow if {
  input.image.host == "registry.company.com"
}

decision := {"allow": allow}
```

The `host` field contains the registry hostname. Docker Hub images have an
empty `host` value. Test with:

```dockerfile
FROM alpine                                    # Allowed (Docker Hub)
FROM registry.company.com/myapp:latest         # Allowed (company registry)
FROM ghcr.io/someorg/image:latest              # Denied (wrong registry)
```

Use `fullRepo` when you need the complete path including registry:

```rego
allow if {
  input.image.fullRepo == "docker.io/library/alpine"
}
```

## Validate platform constraints

Multi-architecture images support different operating systems and CPU
architectures. You can restrict builds to specific platforms:

```rego
package docker

default allow := false

allow if input.local

allow if {
  input.image.os == "linux"
  input.image.arch in ["amd64", "arm64"]
}

decision := {"allow": allow}
```

This policy:

- Defines supported architectures in a list
- Checks `input.image.os` matches Linux
- Verifies `input.image.arch` is in the supported list

The `os` and `arch` fields come from the image manifest, reflecting the actual
image platform. This works with Docker's automatic platform selection -
policies validate what buildx resolves, not what you specify.

## Inspect image metadata

Images contain metadata like environment variables, labels, and working
directories. You can validate these to ensure images meet requirements.

Check for specific environment variables:

```rego
package docker

default allow := false

allow if input.local

allow if {
  input.image.repo == "golang"
  input.image.workingDir == "/go"
  some ver in input.image.env
  startswith(ver, "GOLANG_VERSION=")
  some toolchain in input.image.env
  toolchain == "GOTOOLCHAIN=local"
}

decision := {"allow": allow}
```

This policy validates the official Go image by checking:

- The working directory is `/go`
- The environment has `GOLANG_VERSION` set
- The environment includes `GOTOOLCHAIN=local`

The `input.image.env` field is an array of strings in `KEY=VALUE` format.
Use Rego's `some` iteration to search the array.

Check image labels:

```rego
allow if {
  input.image.labels["org.opencontainers.image.vendor"] == "Example Corp"
  input.image.labels["org.opencontainers.image.version"] != ""
}
```

The `labels` field is a map, so you access values with bracket notation.

## Require attestations and provenance

Modern images include attestations: machine-readable metadata about how the
image was built. Provenance attestations describe the build process, and SBOMs
list the software inside.

Require provenance:

```rego
package docker

default allow := false

allow if input.local

allow if {
  input.image.hasProvenance
}

decision := {"allow": allow}
```

The `hasProvenance` field is `true` when the image has provenance or SBOM
[attestations](../metadata/attestations/_index.md).

For images built with GitHub Actions, verify they came from trusted workflows:

```rego
allow if {
  input.image.repo == "tonistiigi/xx"
  trusted_github_builder("tonistiigi/xx")
}

decision := {"allow": allow}
```

The `trusted_github_builder()` function verifies the image was by Docker's
trusted GitHub Actions builder, using the [reusable GitHub Actions
workflow](https://github.com/docker/github-builder-experimental). For more
verification functions, including Sigstore signature verification, refer to
[Built-in functions](./built-ins.md).

## Combine multiple checks

Real policies often combine several checks. Multiple conditions in one `allow`
rule means AND - all must be true:

```rego
package docker

default allow := false

allow if input.local

# Production images need everything
allow if {
  input.image.repo == "alpine"
  input.image.isCanonical
  input.image.hasProvenance
}

decision := {"allow": allow}
```

Multiple `allow` rules means OR - any rule can match:

```rego
package docker

default allow := false

allow if input.local

# Allow Alpine with strict checks
allow if {
  input.image.repo == "alpine"
  input.image.isCanonical
}

# Allow Go with different checks
allow if {
  input.image.repo == "golang"
  input.image.workingDir == "/go"
}

decision := {"allow": allow}
```

Use this pattern to apply different requirements to different base images.

## Next steps

You now understand how to validate container images in build policies. To
continue learning:

- Learn [Git repository validation](./validate-git.md) for source code inputs
- Browse [Example policies](./examples.md) for complete policy patterns
- Read [Built-in functions](./built-ins.md) for signature verification and
  attestation checking
- Check the [Input reference](./inputs.md) for all available image fields
