---
title: Example policies
linkTitle: Examples
description: Browse the example library of build policies
weight: 40
---

This page provides complete, working policy examples you can copy and adapt.
The examples are organized into two sections: getting started policies for
quick adoption, and production templates for comprehensive security.

If you're new to policies, start with the tutorials:
[Introduction](./intro.md), [Image validation](./validate-images.md), and [Git
validation](./validate-git.md). Those pages teach individual techniques. This
page shows complete policies combining those techniques.

## How to use these examples

1. **Copy** the policy code into a `Dockerfile.rego` file next to your
   Dockerfile
2. **Customize** any TODO comments with your specific values
3. **Test** by running `docker build .` and verifying the policy works as
   expected
4. **Refine** based on your team's needs

## Getting started

These policies work immediately with minimal or no customization. Use them to
adopt policies quickly and demonstrate value to your team.

### Development-friendly baseline

A permissive policy that allows typical development workflows while blocking
obvious security issues.

```rego
package docker

default allow := false

allow if input.local
allow if input.git

# Allow common public registries
allow if {
  input.image.host == "docker.io"  # Docker Hub
}

allow if {
  input.image.host == "ghcr.io" # GitHub Container Registry
}

# Require HTTPS for all downloads
allow if {
  input.http.schema == "https"
}

decision := {"allow": allow}
```

This policy allows local and Git contexts, images from Docker Hub and GitHub
Container Registry, and `ADD` downloads over HTTPS. It blocks HTTP downloads
and non-standard registries.

When to use: Starting point for teams new to policies. Provides basic security
without disrupting development workflows.

### Registry allowlist

Control which registries your builds can pull images from.

```rego
package docker

default allow := false

allow if input.local

# TODO: Add your internal registry hostname
allowed_registries := ["docker.io", "ghcr.io", "registry.company.com"]

allow if {
  input.image.host in allowed_registries
}

decision := {"allow": allow}
```

This policy restricts image pulls to approved registries. Customize and add
your internal registry to the list.

When to use: Enforce corporate policies about approved image sources. Prevents
developers from using arbitrary public registries.

### Pin base images to digests

Require digest references for reproducible builds.

```rego
package docker

default allow := false

allow if input.local

# Require digest references for all images
allow if {
  input.image.isCanonical
}

decision := {"allow": allow}
```

This policy requires images use digest references like
`alpine@sha256:abc123...` instead of tags like `alpine:3.19`. Digests are
immutable - the same digest always resolves to the same image content.

When to use: Ensure build reproducibility. Prevents builds from breaking when
upstream tags are updated. Required for compliance in some environments.

### Control external dependencies

Pin specific versions of dependencies downloaded during builds.

```rego
package docker

default allow := false

allow if input.local

# Allow any image (add restrictions as needed)
allow if input.image

# TODO: Add your allowed Git repositories and tags
allowed_repos := {
  "https://github.com/moby/buildkit.git": ["v0.26.1", "v0.27.0"],
}
# Only allow Git input from allowed_repos
allow if {
  some repo, versions in allowed_repos
  input.git.fullURL == repo
  input.git.tagName in versions
}

# TODO: Add your allowed downloads
allow if {
  input.http.url == "https://example.com/app-v1.0.tar.gz"
}

decision := {"allow": allow}
```

This policy creates allowlists for external dependencies. Add your Git
repositories with approved version tags, and URLs.

When to use: Control which external dependencies can be used in builds.
Prevents builds from pulling arbitrary versions or unverified downloads.

## Production templates

These templates demonstrate comprehensive security patterns. They require
customization but show best practices for production environments.

### Image attestation and provenance

Require images have provenance attestations from trusted builders.

```rego
package docker

default allow := false

allow if input.local

# TODO: Add your repository names
allowed_repos := ["myorg/backend", "myorg/frontend", "myorg/worker"]

# Production images need full attestations
allow if {
  some repo in allowed_repos
  input.image.repo == repo
  input.image.hasProvenance
  trusted_github_builder_ref(repo, "refs/heads/main")
}

# Allow official base images with digests
allow if {
  input.image.repo == "alpine"
  input.image.isCanonical
}

decision := {"allow": allow}
```

This template validates that your application images have provenance
attestations, and were built by GitHub Actions from your main branch. Base
images must use digests.

Customize:

- Replace `allowed_repos` with your image names
- Update the GitHub repository in `trusted_github_builder_ref()`
- Add rules for other base images you use

When to use: Enforce supply chain security for production deployments. Ensures
images are built by trusted CI/CD pipelines with auditable provenance.

### Signed Git releases

Enforce signed tags from trusted maintainers for Git dependencies.

```rego
package docker

default allow := false

allow if input.local

allow if input.image

# TODO: Replace with your repository URL
is_buildkit if {
    input.git.fullURL == "https://github.com/moby/buildkit.git"
}

is_version_tag if {
    is_buildkit
    regex.match(`^v[0-9]+\.[0-9]+\.[0-9]+$`, input.git.tagName)
}

# TODO: Add maintainer GitHub usernames
maintainers := ["tonistiigi", "crazy-max", "jsternberg"]

# Version tags must be signed by maintainers
allow if {
    is_version_tag
    some maintainer in maintainers
    git_signed_github_user(input.git, maintainer)
}

# Allow unsigned refs for development
allow if {
    is_buildkit
    not is_version_tag
}

decision := {"allow": allow}
```

This template requires production release tags to be signed by trusted
maintainers. Development branches and commits can be unsigned.

Customize:

- Replace the repository URL in `is_buildkit`
- Update the `maintainers` list with GitHub usernames
- Adjust the version tag regex pattern if needed

When to use: Validate that production dependencies come from signed releases.
Protects against compromised releases or unauthorized updates.

### Multi-registry policy

Apply different validation rules for internal vs external registries.

```rego
package docker

default allow := false

allow if input.local

# TODO: Replace with your internal registry hostname
internal_registry := "registry.company.com"

# Internal registry: basic validation
allow if {
  input.image.host == internal_registry
}

# External registries: strict validation
allow if {
  input.image.host != internal_registry
  input.image.host != ""
  input.image.isCanonical
  input.image.hasProvenance
}

# Docker Hub: allowlist specific images
allow if {
  input.image.host == ""
  # TODO: Add your approved base images
  input.image.repo in ["alpine", "golang", "node"]
  input.image.isCanonical
}

decision := {"allow": allow}
```

This template defines a trust boundary between internal and external image
sources. Internal images require minimal validation, while external images need
digests and provenance.

Customize:

- Set your internal registry hostname
- Add your approved Docker Hub base images
- Adjust validation requirements based on your security policies

When to use: Organizations with internal registries that need different rules
for internal vs external sources. Balances security with practical workflow
needs.

### Multi-environment policy

Apply different rules based on the build target or stage. For example, 

```rego
package docker

default allow := false

allow if input.local

# TODO: Define your environment detection logic
is_production if {
  input.env.target == "production"
}

is_development if {
  input.env.target == "development"
}

# Production: strict rules - only digest images with provenance
allow if {
  is_production
  input.image.isCanonical
  input.image.hasProvenance
}

# Development: permissive rules - any image
allow if {
  is_development
  input.image
}

# Staging inherits production rules (default target detection)
allow if {
  not is_production
  not is_development
  input.image.isCanonical
}

decision := {"allow": allow}
```

This template uses build targets to apply different validation levels.
Production requires attestations and digests, development is permissive, and
staging uses moderate rules.

Customize:

- Update environment detection logic (target names, build args, etc.)
- Adjust validation requirements for each environment
- Add more environments as needed

When to use: Teams with separate build configurations for different deployment
stages. Allows flexibility in development while enforcing strict rules for
production.

### Complete dependency pinning

Pin all external dependencies to specific versions across all input types.

```rego
package docker

default allow := false

allow if input.local

# TODO: Add your pinned images with exact digests
allowed_images := {
  "alpine": "sha256:4b7ce07002c69e8f3d704a9c5d6fd3053be500b7f1c69fc0d80990c2ad8dd412",
  "golang": "sha256:abc123...",
}

allow if {
  input.image
  some repo, digest in allowed_images
  input.image.repo == repo
  input.image.checksum == digest
}

# TODO: Add your pinned Git dependencies
allowed_git := {
  "https://github.com/moby/buildkit.git": {
    "tag": "v0.26.1",
    "commit": "abc123...",
  },
}

allow if {
  some url, version in allowed_git
  input.git.fullURL == url
  input.git.tagName == version.tag
  input.git.commitChecksum == version.commit
}

# TODO: Add your pinned HTTP downloads
allowed_downloads := {
  "https://releases.example.com/app-v1.0.tar.gz": "sha256:def456...",
}

allow if {
  some url, checksum in allowed_downloads
  input.http.url == url
  input.http.checksum == checksum
}

decision := {"allow": allow}
```

This template pins every external dependency to exact versions with cryptographic
verification. Images use digests, Git repos use commit SHAs, and downloads use
checksums.

Customize:

- Add all your dependencies with exact versions/checksums
- Maintain this file when updating dependencies
- Consider automating updates through CI/CD

When to use: Maximum reproducibility and security. Ensures builds always use
exact versions of all dependencies. Required for high-security or regulated
environments.

## Next steps

- Review [Built-in functions](./built-ins.md) for signature verification and
  attestation checking
- Check the [Input reference](./inputs.md) for all available fields you can
  validate
- Read the tutorials for detailed explanations:
  [Introduction](./intro.md), [Image validation](./validate-images.md), [Git
  validation](./validate-git.md)
