---
title: Validating build inputs with policies
linkTitle: Validating builds
description: Secure your Docker builds by validating images, Git repositories, and dependencies with build policies
keywords: build policies, opa, rego, docker security, supply chain, attestations
weight: 70
params:
  sidebar:
    badge:
      color: blue
      text: Experimental
---

Building with Docker often involves downloading remote resources. These
external dependencies, such as Docker images, Git repositories, remote files,
and other artifacts, are called build inputs.

For example:

- Pulling images from a registry
- Cloning a source code repository
- Fetching files from a server over HTTPS

When consuming build inputs, it's a good idea to verify the contents are what
you expect them to be. One way to do this is to use the `--checksum` option for
the `ADD` Dockerfile instruction. This lets you verify the SHA256 checksum of a
remote resource when pulling it into a build:

```dockerfile
ADD --checksum=sha256:c0ff3312345… https://example.com/archive.tar.gz /
```

If the remote `archive.tar.gz` file does not match the checksum that the
Dockerfile expects, the build fails.

Checksums verify that content matches what you expect, but only for the `ADD`
instruction. They don't tell you anything about where the content came from or
how it was produced. You can't use checksums to enforce constraints like
"images must be signed" or "dependencies must come from approved sources."

Build policies solve this problem. They let you define rules that validate all
your build inputs, enforcing requirements like provenance attestations,
approved registries, and signed Git tags across your entire build process.

## Prerequisites

Build policies is currently an experimental feature. To try it out, you'll
need:

- Buildx 0.31.0 or later - Check your version: `docker buildx version`
- BuildKit 0.27.0 or later - Verify with: `docker buildx inspect --bootstrap`

If you're using Docker Desktop, ensure you're on a version that includes these
updates.

## Build policies

Buildx version 0.31.0 added support for build policies. Build policies are
rules for securing your Docker build supply chain, and help protect against
upstream compromises, malicious dependencies, and unauthorized modifications to
your build inputs.

Build policies let you enforce extended verifications on inputs used to build
your projects, such as:

- Docker images must use digest references (not tags alone)
- Images must have provenance attestations and cosign signatures
- Git tags are signed by maintainers with a PGP public key
- All remote artifacts must use HTTPS and include a checksum for verification

Build policies are defined in a declarative policy language, called Rego,
created for the [Open Policy Agent (OPA)](https://www.openpolicyagent.org/).
The following example shows a minimal build policy in Rego.

```rego {title="Dockerfile.rego"}
package docker

default allow := false

# Allow any local inputs for this build
# For example: a local build context, or a local Dockerfile
allow if input.local

# Allow images, but only if they have provenance attestations
allow if {
    input.image.hasProvenance
}

decision := {"allow": allow}
```

If the Dockerfile associated with this policy references an image with no
provenance attestation in a `FROM` instruction, the policy would be violated
and the build would fail.

## How policies work

When you run `docker buildx build`, Buildx:

1. Resolves all build inputs (images, Git repos, HTTP downloads)
2. Looks for a policy file matching your Dockerfile name (e.g.,
   `Dockerfile.rego`)
3. Evaluates each input against the policy before the build starts
4. Allows the build to proceed only if all inputs pass the policy

Policies are written in Rego (Open Policy Agent's policy language). You don't
need to be a Rego expert - the [Introduction](./intro.md) tutorial teaches you
everything needed.

Policy files live alongside your Dockerfile:

```text
project/
├── Dockerfile
├── Dockerfile.rego
└── src/
```

No additional configuration is needed - Buildx automatically finds and loads
the policy when you build.

## Use cases

Build policies help you enforce security and compliance requirements on your
Docker builds. Common scenarios where policies provide value:

### Enforce base image standards

Require all production Dockerfiles to use specific, approved base images with
digest references. Prevent developers from using arbitrary images that haven't
been vetted by your security team.

### Validate third-party dependencies

When your build downloads files, libraries, or tools from the internet, verify
they come from trusted sources and match expected checksums or signatures. This
protects against supply chain attacks where an upstream dependency is
compromised.

### Ensure signed releases

Require that all dependencies have valid signatures from trusted parties.

- Check GPG signatures for Git repositories you clone in your builds
- Verify provenance attestation signatures with Sigstore

### Meet compliance requirements

Some regulatory frameworks require evidence that you validate your build
inputs. Build policies give you an auditable, declarative way to demonstrate
you're checking dependencies against security standards.

### Separate development and production rules

Apply stricter validation for production builds while allowing more flexibility
during development. The same policy file can contain conditional rules based on
build context or target.

## Get started

Ready to start writing policies? The [Introduction](./intro.md) tutorial walks
you through creating your first policy and teaches the Rego basics you need.

For practical usage guidance, see [Using build policies](./usage.md).

For practical examples you can copy and adapt, see the [Example
policies](./examples.md) library.
