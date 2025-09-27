---
title: glibc and musl support in Docker Hardened Images
linktitle: glibc and musl
description: Compare glibc and musl variants of DHIs to choose the right base image for your application’s compatibility, size, and performance needs.
keywords: glibc vs musl, alpine musl image, debian glibc container, docker hardened images compatibility, c library in containers
---

Docker Hardened Images (DHI) are built to prioritize security without
sacrificing compatibility with the broader open source and enterprise software
ecosystem. A key aspect of this compatibility is support for common Linux
standard libraries: `glibc` and `musl`.

## What are glibc and musl?

When you run Linux-based containers, the image's C library plays a key role in
how applications interact with the operating system. Most modern Linux
distributions rely on one of the following standard C libraries:

- `glibc` (GNU C Library): The standard C library on mainstream distributions
  like Debian, Ubuntu, and Red Hat Enterprise Linux. It is widely supported and
  typically considered the most compatible option across languages, frameworks,
  and enterprise software.

- `musl`: A lightweight alternative to `glibc`, commonly used in minimal
  distributions like Alpine Linux. While it offers smaller image sizes and
  performance benefits, `musl` is not always fully compatible with software that
  expects `glibc`.

## DHI compatibility

DHI images are available in both `glibc`-based (e.g., Debian) and `musl`-based
(e.g., Alpine) variants. For enterprise applications and language runtimes where
compatibility is critical, we recommend using DHI images based on glibc.

## What to choose, glibc or musl?

Docker Hardened Images are available in both glibc-based (Debian) and musl-based
(Alpine) variants, allowing you to choose the best fit for your workload.

Choose Debian-based (`glibc`) images if:

- You need broad compatibility with enterprise workloads, language runtimes, or
  proprietary software.
- You're using ecosystems like .NET, Java, or Python with native extensions that
  depend on `glibc`.
- You want to minimize the risk of runtime errors due to library
  incompatibilities.

Choose Alpine-based (`musl`) images if:

- You want a minimal footprint with smaller image sizes and reduced surface
  area.
- You're building a custom or tightly controlled application stack where
  dependencies are known and tested.
- You prioritize startup speed and lean deployments over maximum compatibility.

If you're unsure, start with a Debian-based image to ensure compatibility, and
evaluate Alpine once you're confident in your application's dependencies.