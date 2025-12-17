---
title: Understanding roles and responsibilities for Docker Hardened Images
linkTitle: Responsibility overview
description: Understand the division of responsibilities between Docker, upstream projects, and you when using Docker Hardened Images.
keywords: software supply chain security, signed sbom, vex document, container provenance, image attestation
weight: 46
---

Docker Hardened Images (DHIs) are curated and maintained by Docker, and built
using upstream open source components. To deliver security, reliability, and
compliance, responsibilities are shared among three groups:

- Upstream maintainers: the developers and communities responsible for the
  open source software included in each image.
- Docker: the provider of hardened, signed, and maintained container images.
- You (the customer): the consumer who runs and, optionally, customizes DHIs
  in your environment.

This topic outlines who handles what, so you can use DHIs effectively and
securely.

## Releases

- Upstream: Publishes and maintains official releases of the software
  components included in DHIs. This includes versioning, changelogs, and
  deprecation notices.
- Docker: Builds, hardens, and signs Docker Hardened Images based on
  upstream versions. Docker maintains these images in line with upstream release
  timelines and internal policies.
- You: Ensure you're staying on supported versions of DHIs and upstream
  projects. Using outdated or unsupported components can introduce security
  risk.

## Patching

- Upstream: Maintains and updates the source code for each component,
  including fixing vulnerabilities in libraries and dependencies.
- Docker: Rebuilds and re-releases images with upstream patches applied. Docker
  monitors for vulnerabilities and publishes updates to affected images. Only
  DHI Enterprise includes SLAs. DHI Free offers a secure baseline but no
  guaranteed remediation timelines.
- You: Apply DHI updates in your environments and patch any software or
  dependencies you install on top of the base image.

## Testing

- Upstream: Defines the behavior and functionality of the original software,
  and is responsible for validating core features.
- Docker: Validates that DHIs start, run, and behave consistently with
  upstream expectations. Docker also runs security scans and includes a [testing
  attestation](../core-concepts/attestations.md) with each image.
- You: Test your application on top of DHIs and validate that any changes or
  customizations function as expected in your environment.

## Security and compliance

- Docker: Publishes signed SBOMs, VEX documents, provenance data, and CVE
  scan results with each image to support compliance and supply chain security.
  - For free DHI users: All security metadata and transparency features are
    included at no cost.
  - For DHI Enterprise users: Additional compliance variants (like FIPS and
    STIG) and customization capabilities are available, with automatic rebuilds
    when base images are patched.
- You: Integrate DHIs into your security and compliance workflows, including
  vulnerability management and auditing.

## Support

- Docker:
  - For free DHI users: Community support and public documentation are available.
  - For DHI Enterprise users: Access to Docker's enterprise support team for
    mission-critical applications.
- You: Monitor Docker's release notes, security advisories, and documentation
  for updates and best practices.

## Summary

Docker Hardened Images give you a secure foundation, complete with signed
metadata and upstream transparency. Your role is to make informed use of these
images, apply updates promptly, and validate that your configurations and
applications meet your internal requirements.