---
title: Base image hardening
linktitle: Hardening
description: Learn how Docker Hardened Images are designed for security, with minimal components, nonroot execution, and secure-by-default configurations.
keywords: hardened base image, minimal container image, non-root containers, secure container configuration, remove package manager
---

## What is base image hardening?

Base image hardening is the process of securing the foundational layers of a
container image by minimizing what they include and configuring them with
security-first defaults. A hardened base image removes unnecessary components,
like shells, compilers, and package managers, which limits the available attack
surface, making it more difficult for an attacker to gain control or escalate
privileges inside the container.

Hardening also involves applying best practices like running as a non-root user,
reducing writable surfaces, and ensuring consistency through immutability. While
[Docker Official
Images](../../docker-hub/image-library/trusted-content.md#docker-official-images)
and [Docker Verified Publisher
Images](../../docker-hub/image-library/trusted-content.md#verified-publisher-images)
follow best practices for security, they may not be as hardened as Docker
Hardened Images, as they are designed to support a broader range of use cases.

## Why is it important?

Most containers inherit their security posture from the base image they use. If
the base image includes unnecessary tools or runs with elevated privileges,
every container built on top of it is exposed to those risks.

Hardening the base image:

- Reduces the attack surface by removing tools and libraries that could be exploited
- Enforces least privilege by dropping root access and restricting what the container can do
- Improves reliability and consistency by avoiding runtime changes and drift
- Aligns with secure software supply chain practices and helps meet compliance standards

Using hardened base images is a critical first step in securing the software you
build and run in containers.

## What's removed and why

Hardened images typically exclude common components that are risky or unnecessary in secure production environments:

| Removed component                                | Reason                                                                           |
|--------------------------------------------------|----------------------------------------------------------------------------------|
| Shells (e.g., `sh`, `bash`)                      | Prevents users or attackers from executing arbitrary commands inside containers  |
| Package managers (e.g., `apt`, `apk`)            | Disables the ability to install software post-build, reducing drift and exposure |
| Compilers and interpreters                       | Avoids introducing tools that could be used to run or inject malicious code      |
| Debugging tools (e.g., `strace`, `curl`, `wget`) | Reduces risk of exploitation or information leakage                              |
| Unused libraries or locales                      | Shrinks image size and minimizes attack vectors                                  |

## How Docker Hardened Images apply base image hardening

Docker Hardened Images (DHIs) apply base image hardening principles by design.
Each image is constructed to include only what is necessary for its specific
purpose, whether that’s building applications (with `-dev` or `-sdk` tags) or
running them in production.

### Docker Hardened Image traits

Docker Hardened Images are built to be:

- Minimal: Only essential libraries and binaries are included
- Immutable: Images are fixed at build time—no runtime installations
- Non-root by default: Containers run as an unprivileged user unless configured otherwise
- Purpose-scoped: Different tags are available for development (`-dev`), SDK-based builds (`-sdk`), and production runtime

These characteristics help enforce consistent, secure behavior across development, testing, and production environments.

### Docker Hardened Image compatibility considerations

Because Docker Hardened Images strip out many common tools, they may not work out of the box for all use cases. You may need to:

- Use multi-stage builds to compile code or install dependencies in a `-dev` image and copy the output into a hardened runtime image
- Replace shell scripts with equivalent entrypoint binaries or explicitly include a shell if needed
- Use [Docker Debug](../../../reference/cli/docker/debug.md) to temporarily inspect or troubleshoot containers without altering the base image

These trade-offs are intentional and help support best practices for building secure, reproducible, and production-ready containers.