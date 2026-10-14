---
description: Learn about the latest features of Docker Build Cloud
keywords: docker build cloud, release notes, changelog, features, changes, delta, new, releases
title: Docker Build Cloud release notes
linkTitle: Release notes
tags: [Release notes]
---

This page contains information about the new features, improvements, known
issues, and bug fixes in Docker Build Cloud releases. 

## 2026-09-14

### Enhancements

- Cloud builders now run BuildKit v0.33.0, upgraded from v0.32.2. Notable
  changes:

  - The maximum size of an exported attestation, such as an SBOM, is raised
    from 40 MiB to 80 MiB. Builds whose SBOM exceeded the previous limit failed
    with an error ending in `sbom.spdx.json exceeds 41943040 bytes`.
  - The built-in Dockerfile frontend is updated to
    [v1.27.0](https://github.com/moby/buildkit/releases/tag/dockerfile%2F1.27.0).
    This only affects builds that don't pin a frontend with a `# syntax=`
    directive.
  - Credentials embedded in HTTP source URLs and Git bundle checkouts are now
    redacted from more places in build progress output and error messages.
    Builds that pass credentials with
    [build secrets](/manuals/build/building/secrets.md) weren't affected.
  - Git advice messages, such as detached `HEAD` guidance, are no longer shown
    in build output for Git contexts and `ADD` Git sources. To show them again,
    set the build argument `BUILDKIT_GIT_ADVICE=1`.

  For the full list of changes, see the
  [BuildKit release notes](https://github.com/moby/buildkit/releases/tag/v0.33.0).

- Docker Build Cloud is verified against Buildx v0.37.0, up from v0.36.1. Buildx
  is distributed with Docker Desktop and Docker Engine rather than with Docker
  Build Cloud, so your client version depends on how you installed Docker. See
  the [Buildx release notes](https://github.com/docker/buildx/releases) for
  client-side changes.

### Bug fixes

- Fixed an error containing `no active session for` that could occur when
  running concurrent builds that use a remote cache.
- Fixed a possible `failed to apply diffs: snapshot does not exist` error.
- Fixed a cache link that the remote cache exporter could silently drop.

## 2026-08-18

### Enhancements

- Cloud builders now run BuildKit v0.32.2, upgraded from v0.32.1. This release
  reverts a v0.32.1 workaround that pushed each attestation manifest only after
  the image manifest it refers to. BuildKit again follows the OCI distribution
  specification, which lets a referrer be pushed before its subject. If your
  registry rejects a manifest whose subject doesn't exist yet, build with
  `--output type=image,oci-artifact=false`. For details, see the
  [BuildKit release notes](https://github.com/moby/buildkit/releases/tag/v0.32.2).

- Docker Build Cloud is verified against Buildx v0.36.1, up from v0.36.0. See
  the [Buildx release notes](https://github.com/docker/buildx/releases) for
  client-side changes.

## 2026-08-03

### Enhancements

- Cloud builders now run BuildKit v0.32.1, upgraded from v0.20.0. This spans
  twelve minor BuildKit releases of new build features, performance
  improvements, and bug fixes. For the full list of changes, see the
  [BuildKit release notes](https://github.com/moby/buildkit/releases).

  This range includes the following changes to default behavior. Each of these
  is an [exporter attribute](/manuals/build/exporters/image-registry.md) that
  you pass with `--output`, or an
  [attestation attribute](/manuals/build/metadata/attestations/_index.md) that
  you pass with `--attest`:

  - Image results now use OCI media types by default. If your registry doesn't
    support OCI media types, build with
    `--output type=image,oci-mediatypes=false`. Treat this as a temporary
    measure until your registry supports OCI media types.
  - Attestations now use OCI artifact descriptors by default. If your registry
    doesn't support OCI artifacts, build with
    `--output type=image,oci-artifact=false`. Note that the attribute is
    singular.
  - [Provenance attestations](/manuals/build/metadata/attestations/slsa-provenance.md)
    now default to SLSA v1 instead of v0.2, which changes both the predicate
    type and the schema. If you verify provenance in your CI pipeline, confirm
    that your verifier supports the SLSA v1 predicate before your next build.
    A verifier that doesn't recognize the new predicate type may report that no
    attestation was found rather than failing outright, which can silently
    weaken a supply chain check. To keep the previous format while you update
    your tooling, build with `--attest type=provenance,version=v0.2`.

- Docker Build Cloud is verified against Buildx v0.36.0, up from v0.21.0. Buildx
  is the client you build with, and it's distributed with Docker Desktop and
  Docker Engine rather than with Docker Build Cloud, so your client version
  depends on how you installed Docker. For the client-side features available
  when building with a cloud builder, see the
  [Buildx release notes](https://github.com/docker/buildx/releases).

## 2025-12-19

### Bug fixes

- Fixed an error when building
  [Docker Hardened Images](/manuals/dhi/_index.md) with a cloud builder.
- Fixed builds failing with gRPC message size errors when transferring large
  amounts of data.
- Fixed a 500 error when pulling images whose attestations or signatures are
  resolved through the OCI referrers fallback, which registries use when they
  don't implement the referrers API directly.

## 2025-06-04

### Enhancements

- Build timeouts are now determined by your subscription plan rather than by a
  single fixed limit applied to every build.

## 2025-04-29

### Enhancements

- Improved build error messages. Failures caused by your Dockerfile or build
  context are now reported as build errors rather than internal errors. This
  includes denied base image pulls, invalid stage names, and invalid `chmod`
  and `mkdir` targets.

## 2025-04-09

### Enhancements

- Added a consolidated cloud usage report covering all builds in your
  organization.

## 2025-03-05

### Enhancements

- Build usage reports now support custom date ranges.

## 2025-02-24

### New

Added a new **Build settings** page where you can configure disk allocation, private resource access, and firewall settings for your cloud builders in your organization. These configurations help optimize storage, enable access to private registries, and secure outbound network traffic.