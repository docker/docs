---
title: Docker Hardened Images
description: Secure, minimal, and production-ready base images
weight: 8
params:
  sidebar:
    group: Products
    badge:
      color: green
      text: New
  grid_sections:
    - title: Quickstart
      description: Follow a step-by-step guide to explore and run a Docker Hardened Image.
      icon: rocket_launch
      link: /dhi/get-started/
    - title: Explore
      description: Learn what Docker Hardened Images are, how they're built, and what sets them apart from typical base images.
      icon: info
      link: /dhi/explore/
    - title: Features
      description: Discover the security, compliance, and enterprise-readiness features built into Docker Hardened Images.
      icon: lock
      link: /dhi/features/
    - title: How-tos
      description: Step-by-step guides for using, verifying, scanning, and migrating to Docker Hardened Images.
      icon: play_arrow
      link: /dhi/how-to/
    - title: Core concepts
      description: Understand the secure supply chain principles that make Docker Hardened Images production-ready.
      icon: fact_check
      link: /dhi/core-concepts/
    - title: Troubleshoot
      description: Resolve common issues with building, running, or debugging Docker Hardened Images.
      icon: help_center
      link: /dhi/troubleshoot/
    - title: Additional resources
      description: Links to blog posts, Docker Hub catalog, GitHub repositories, and more.
      icon: link
      link: /dhi/resources/
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
| Customizations | ❌ | Up to 5 | Unlimited |
| Access to Hardened System Packages repository | ❌ | ❌ | ✅ |
| Full catalog access available | ❌ | ❌ | ✅ |
| Extended Lifecycle Support add-on available | ❌ | ❌ | ✅<br><br>Includes:<br>✅ +5 years of hardened updates<br>✅ Maintains security updates after upstream EOL<br>✅ SBOMs & provenance<br>✅ Protects long-lived workloads |

For pricing and more details, see the [Docker Hardened Images subscription
comparison](https://www.docker.com/products/hardened-images/#compare).

Explore the sections below to get started with Docker Hardened Images, integrate
them into your workflow, and learn what makes them secure and enterprise-ready.

{{< grid
  items="grid_sections"
>}}
