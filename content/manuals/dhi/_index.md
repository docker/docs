---
title: Docker Hardened Images
description: Secure, minimal, and production-ready base images
weight: 8
aliases:
  - /dhi/features/
  - /dhi/features/secure/
  - /dhi/features/integration/
  - /dhi/features/support/
  - /dhi/features/patching/
  - /dhi/features/flexible/
  - /dhi/features/helm/
params:
  sidebar:
    group: Supply chain security
  grid_sections:
    - title: Quickstart
      description: Follow a step-by-step guide to explore and run a Docker Hardened Image.
      icon: rocket-launch
      link: /dhi/get-started/
    - title: Explore
      description: Learn what Docker Hardened Images are, how they're built, and what sets them apart from typical base images.
      icon: information-circle
      link: /dhi/explore/
    - title: How-tos
      description: Step-by-step guides for using, verifying, scanning, and migrating to Docker Hardened Images.
      icon: play
      link: /dhi/how-to/
    - title: Security concepts
      description: Understand the secure supply chain principles that make Docker Hardened Images production-ready.
      icon: clipboard-document-check
      link: /dhi/explore/security-concepts/
    - title: Tools
      description: Browse and manage Docker Hardened Images using Docker Hub, the CLI, MCP server, or Terraform.
      icon: wrench-screwdriver
      link: /dhi/tools/
    - title: Additional resources
      description: Guides, Docker Hub catalog, GitHub repositories, and more.
      icon: link
      link: /dhi/resources/
    - title: Release notes
      description: New features, improvements, and changes in Docker Hardened Images.
      icon: newspaper
      link: /dhi/release-notes/platform/
---

Docker Hardened Images (DHI) provide minimal, secure, and production-ready
container images, Helm charts, and system packages maintained by Docker.
Designed to reduce vulnerabilities and simplify compliance, DHI integrates
easily into your existing Docker-based workflows with little to no retooling
required.

DHI is available in the following three subscriptions.

| Feature | Community | Select | Enterprise |
|---|---|---|---|
| Hardened, minimal images | ✅ | ✅ | ✅ |
| Near-zero CVEs | ✅ | ✅ | ✅ |
| Verifiable SBOMs & SLSA Build L3 provenance | ✅ | ✅ | ✅ |
| Full, unsuppressed CVE visibility | ✅ | ✅ | ✅ |
| Drop-in adoption, no workflow changes | ✅ | ✅ | ✅ |
| Full catalog of open source images under Apache 2.0 | ✅ | ✅ | ✅ |
| Built with Docker Hardened System Packages | ✅ | ✅ | ✅ |
| Upstream cadence for Docker-released patches | ✅ | ✅ | ✅ |
| FIPS/STIG variants | ❌ | ✅ | ✅ |
| Critical CVE fixes < 7 days with SLA-backed continuous patching | ❌ | ✅ | ✅ |
| Customizations | ❌ | ✅ Up to 5 | ✅ Unlimited |
| Access to Hardened System Packages repository | ❌ | ❌ | ✅ |
| Full catalog access available | ❌ | ❌ | ✅ |
| Extended Lifecycle Support add-on available | ❌ | ❌ | ✅ +5 years of hardened updates |

For pricing and more details, see the [Docker Hardened Images subscription
comparison](https://www.docker.com/products/hardened-images/#compare).

## Community features

DHI's core features are free to use, share, and build on under Apache 2.0.

- [Near-zero CVEs](/dhi/explore/security-concepts/cves/): continuously scanned and patched
  to maintain minimal known vulnerabilities
- [Distroless variants](/dhi/explore/security-concepts/distroless/): remove unnecessary
  components, reducing attack surface by up to 95%
- Non-root execution: containers run as non-root by default
- [Hardened system packages](/dhi/how-to/hardened-packages/): system packages
  built from source, cryptographically signed, and verified by Docker
- [SLSA Build Level 3 provenance](/dhi/explore/security-concepts/slsa/), [signed
  SBOMs](/dhi/explore/security-concepts/sbom/), [VEX statements](/dhi/explore/security-concepts/vex/),
  and [cryptographic signatures](/dhi/explore/security-concepts/signatures/) on every image
- Built on Alpine and Debian with [glibc and musl
  variants](/dhi/explore/security-concepts/glibc-musl/); dev and runtime image variants
  available
- Works with existing Docker workflows, CI/CD pipelines, and tools with no
  retooling required
- [Helm charts](/dhi/how-to/helm/): Docker-provided charts built from upstream
  sources, tested for compatibility with DHI, and available as OCI artifacts in
  the DHI catalog; include SLSA Level 3 provenance, SBOMs, and cryptographic
  signing

## Select and Enterprise features

For organizations with strict security or compliance requirements:

- [7-day SLA](https://docs.docker.com/go/dhi-sla/) for critical and high
  severity CVE remediation
- [FIPS-enabled](/dhi/explore/security-concepts/fips/) and
  [STIG-ready](/dhi/explore/security-concepts/stig/) compliance variants
- [Customization](/dhi/how-to/customize/): add packages, tools, certificates,
  and configurations (up to 5 with Select, unlimited with Enterprise)
- [Enterprise package repository](/dhi/how-to/hardened-packages/) access and
  full catalog access (Enterprise)
- Extended Lifecycle Support: post-EOL security patches, updated SBOMs,
  provenance, and signing (Enterprise add-on)

## Get started

Explore the sections below to get started with Docker Hardened Images, integrate
them into your workflow, and learn what makes them secure and enterprise-ready.

{{< grid items="grid_sections" >}}
