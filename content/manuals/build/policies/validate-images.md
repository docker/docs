---
title: Validating image inputs
linkTitle: Image validation
description: Write policies to validate container images used in your builds
keywords: build policies, image validation, docker images, provenance, attestations, signatures
weight: 30
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

## Compare semantic versions

Restrict images to specific version ranges using Rego's `semver` functions:

```rego
package docker

default allow := false

allow if input.local

# Allow Go 1.21 or newer
allow if {
  input.image.repo == "golang"
  semver.is_valid(input.image.tag)
  semver.compare(input.image.tag, "1.21.0") >= 0
}

decision := {"allow": allow}
```

The `semver.compare(a, b)` function compares semantic versions and returns:

- `-1` if version `a` is less than `b`
- `0` if versions are equal
- `1` if version `a` is greater than `b`

Use `semver.is_valid()` to check if a tag is a valid semantic version before
comparing.

Restrict to specific version ranges:

```rego
allow if {
  input.image.repo == "node"
  version := input.image.tag
  semver.is_valid(version)
  semver.compare(version, "20.0.0") >= 0  # 20.0.0 or newer
  semver.compare(version, "21.0.0") < 0   # older than 21.0.0
}
```

This allows only Node.js 20.x versions. The pattern works for any image using
semantic versioning.

These `semver` functions are standard Rego built-ins documented in the [OPA
policy
reference](https://www.openpolicyagent.org/docs/latest/policy-reference/#semver).

## Require digest references

Tags like `alpine:3.22` can change - someone could push a new image with the
same tag. Digests like `alpine@sha256:abc123...` are immutable.

### Requiring users to provide digests

You can require that users always specify digests in their Dockerfiles:

```rego
package docker

default allow := false

allow if input.local

allow if {
  input.image.isCanonical
}

decision := {"allow": allow}
```

The `isCanonical` field is `true` when the user's reference includes a digest.
This policy would allow:

```dockerfile
FROM alpine@sha256:4b7ce07002c69e8f3d704a9c5d6fd3053be500b7f1c69fc0d80990c2ad8dd412
```

But reject tag-only references like `FROM alpine:3.22`.

### Pinning to specific digests

Alternatively (or additionally), you can validate that an image's actual digest
matches a specific value, regardless of how the user wrote the reference:

```rego
allow if {
  input.image.repo == "alpine"
  input.image.checksum == "sha256:4b7ce07002c69e8f3d704a9c5d6fd3053be500b7f1c69fc0d80990c2ad8dd412"
}

decision := {"allow": allow}
```

This checks the actual content digest of the pulled image. It would allow both:

```dockerfile
FROM alpine:3.22
FROM alpine@sha256:4b7ce...
```

As long as the resolved image has the specified digest. This is useful for
pinning critical base images to known-good versions.

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

The `host` field contains the registry hostname. Docker Hub images use
`"docker.io"` as the host value. Test with:

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
policies validate what Buildx resolves, not what you specify.

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

Modern images include [attestations](/build/metadata/attestations/):
machine-readable metadata about how the image was built.
[Provenance](/build/metadata/attestations/slsa-provenance/) attestations
describe the build process, and [SBOMs](/build/metadata/attestations/sbom/)
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

## Verify GitHub Actions signatures

For images built with GitHub Actions, verify they came from trusted workflows by
inspecting signature metadata:

```rego
allow if {
  input.image.repo == "myapp"
  input.image.hasProvenance
  some sig in input.image.signatures
  valid_github_signature(sig)
}

# Helper to validate GitHub Actions signature
valid_github_signature(sig) if {
  sig.signer.certificateIssuer == "CN=sigstore-intermediate,O=sigstore.dev"
  sig.signer.issuer == "https://token.actions.githubusercontent.com"
  startswith(sig.signer.buildSignerURI, "https://github.com/myorg/")
  sig.signer.runnerEnvironment == "github-hosted"
}

decision := {"allow": allow}
```

This pattern works with any GitHub Actions workflow using Sigstore keyless
signing. The signature metadata provides cryptographic proof of the build's
origin. For complete signature verification examples, see [Example
policies](./examples.md).

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
