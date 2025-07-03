---
linktitle: Image types
title: Available types of Docker Hardened Images
description: Learn about the different image types, distributions, and variants offered in the Docker Hardened Images catalog.
keywords: docker hardened images, distroless containers, distroless images, docker distroless, alpine base image, debian base image, development containers, runtime containers, secure base image, multi-stage builds
weight: 20
---

Docker Hardened Images (DHI) is a comprehensive catalog of
security-hardened container images built to meet diverse
development and production needs.

## Framework and application images

DHI includes a selection of popular frameworks and application images, each
hardened and maintained to ensure security and compliance. These images
integrate seamlessly into existing workflows, allowing developers to focus on
building applications without compromising on security.

For example, you might find repositories like the following in the DHI catalog:

- `node`: framework for Node.js applications
- `python`: framework for Python applications
- `nginx`: web server image

## Compatibility options

Docker Hardened Images are available in different base image options, giving you
flexibility to choose the best match for your environment and workload
requirements:

- Debian-based images: A good fit if you're already working in glibc-based
  environments. Debian is widely used and offers strong compatibility across
  many language ecosystems and enterprise systems.

- Alpine-based images: A smaller and more lightweight option using musl libc.
  These images tend to be small and are therefore faster to pull and have a
  reduced footprint.

Each image maintains a minimal and secure runtime layer by removing
non-essential components like shells, package managers, and debugging tools.
This helps reduce the attack surface while retaining compatibility with common
runtime environments.

Example tags include:

- `3.9.23-alpine3.21`: Alpine-based image for Python 3.9.23
- `3.9.23-debian12`: Debian-based image for Python 3.9.23

If you're not sure which to choose, start with the base you're already familiar
with. Debian tends to offer the broadest compatibility.

## Development and runtime variants

To accommodate different stages of the application lifecycle, DHI offers all
language framework images and select application images in two variants:

- Development (dev) images: Equipped with necessary development tools and
libraries, these images facilitate the building and testing of applications in a
secure environment. They include a shell, package manager, a root user, and
other tools needed for development.

- Runtime images: Stripped of development tools, these images contain only the
essential components needed to run applications, ensuring a minimal attack
surface in production.

This separation supports multi-stage builds, enabling developers to compile code
in a secure build environment and deploy it using a lean runtime image.

For example, you might find tags like the following in a DHI repository:

- `3.9.23-debian12`: runtime image for Python 3.9.23
- `3.9.23-debian12-dev`: development image for Python 3.9.23

## FIPS variants

Some Docker Hardened Images include a `-fips` variant. These variants use
cryptographic modules that have been validated under [FIPS
140](../core-concepts/fips.md), a U.S. government standard for secure
cryptographic operations.

FIPS variants are designed to help organizations meet regulatory and compliance
requirements related to cryptographic use in sensitive or regulated
environments.

You can recognize FIPS variants by their tag that includes `-fips`.

For example:
- `3.13-fips`: FIPS variant of the Python 3.13 image
- `3.9.23-debian12-fips`: FIPS variant of the Debian-based Python 3.9.23 image

FIPS variants can be used in the same way as any other Docker Hardened Image and
are ideal for teams operating in regulated industries or under compliance
frameworks that require cryptographic validation.
