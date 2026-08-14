---
description: Learn about the latest features of Docker Build Cloud
keywords: docker build cloud, release notes, changelog, features, changes, delta, new, releases
title: Docker Build Cloud release notes
linkTitle: Release notes
tags: [Release notes]
---

This page contains information about the new features, improvements, known
issues, and bug fixes in Docker Build Cloud releases. 

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