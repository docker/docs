---
title: Minimal or distroless images
linktitle: Distroless images
description: Learn how Docker Hardened Images use distroless variants to minimize attack surface and remove unnecessary components.
keywords: distroless container image, minimal docker image, secure base image, no shell container, reduced attack surface
---


Minimal images, sometimes called distroless images, are container images
stripped of unnecessary components such as package managers, shells, or even the
underlying operating system distribution. Docker Hardened Images (DHI) embrace
this minimal approach to reduce vulnerabilities and enforce secure software
delivery. [Docker Official
Images](../../docker-hub/image-library/trusted-content.md#docker-official-images)
and [Docker Verified Publisher
Images](../../docker-hub/image-library/trusted-content.md#verified-publisher-images)
follow similar best practices for minimalism and security but may not be as
stripped down to ensure compatibility with a wider range of use cases.

## What are minimal or distroless images?

Traditional container images include a full OS, often more than what is needed
to run an application. In contrast, minimal or distroless images include only:

- The application binary
- Its runtime dependencies (e.g., libc, Java, Python)
- Any explicitly required configuration or metadata

They typically exclude:

- OS tools (e.g., `ls`, `ps`, `cat`)
- Shells (e.g., `sh`, `bash`)
- Package managers (e.g., `apt`, `apk`)
- Debugging utilities (e.g., `curl`, `wget`, `strace`)

Docker Hardened Images are based on this model, ensuring a smaller and more
secure runtime surface.

## What you gain

| Benefit                | Description                                                                   |
|------------------------|-------------------------------------------------------------------------------|
| Smaller attack surface | Fewer components mean fewer vulnerabilities and less exposure to CVEs         |
| Faster startup         | Smaller image sizes result in faster pull and start times                     |
| Improved security      | Lack of shell and package manager limits what attackers can do if compromised |
| Better compliance      | Easier to audit and verify, especially with SBOMs and attestations            |

## Addressing common tradeoffs

Minimal and distroless images offer strong security benefits, but they can
change how you work with containers. Docker Hardened Images are designed to
maintain productivity while enhancing security.

| Concern           | How Docker Hardened Images help                                                                                                                                                                                         |
|-------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Debuggability | Hardened images exclude shells and CLI tools by default. Use [Docker Debug](../../../reference/cli/docker/debug.md) to temporarily attach a debug sidecar for troubleshooting without modifying the original container. |
| Familiarity   | DHI supports multiple base images, including Alpine and Debian variants, so you can choose a familiar environment while still benefiting from hardening practices.                                                        |
| Flexibility   | Runtime immutability helps secure your containers. Use multi-stage builds and CI/CD to control changes, and optionally use dev-focused base images during development.                                                  |

By balancing minimalism with practical tooling, Docker Hardened Images support
modern development workflows without compromising on security or reliability.

## Best practices for using minimal images

- Use multi-stage builds to separate build-time and runtime environments
- Validate image behavior using CI pipelines, not interactive inspection
- Include runtime-specific dependencies explicitly in your Dockerfile
- Use Docker Scout to continuously monitor for CVEs, even in minimal images

By adopting minimal or distroless images through Docker Hardened Images, you
gain a more secure, predictable, and production-ready container environment
that's designed for automation, clarity, and reduced risk.

