---
title: Docker Build
weight: 20
description: Get an overview of Docker Build to package and bundle your code and ship it anywhere
keywords: build, buildx, buildkit
params:
  sidebar:
    group: Application development
grid:
- title: Packaging your software
  description: 'Build and package your application to run it anywhere: locally or
    in the cloud.'
  icon: archive
  link: /build/concepts/overview/
- title: Multi-stage builds
  description: Keep your images small and secure with minimal dependencies.
  icon: trending-up
  link: /build/building/multi-stage/
- title: Multi-platform images
  description: Build, push, pull, and run images seamlessly on different computer
    architectures.
  icon: copy
  link: /build/building/multi-platform/
- title: BuildKit
  description: Explore BuildKit, the open source build engine.
  icon: wrench
  link: /build/buildkit/
- title: Build drivers
  description: Configure where and how you run your builds.
  icon: wrench
  link: /build/builders/drivers/
- title: Docker Build Cloud
  description: Run builds on managed cloud builders with a shared cache.
  icon: /icons/logo-build-cloud.svg
  link: /build-cloud/
- title: Exporters
  description: Export any artifact you like, not just Docker images.
  icon: upload
  link: /build/exporters/
- title: Build caching
  description: Avoid unnecessary repetitions of costly operations, such as package
    installs.
  icon: refresh-cw
  link: /build/cache/
- title: Bake
  description: Orchestrate your builds with Bake.
  icon: cake
  link: /build/bake/
- title: CLI reference
  description: Find Buildx commands and options for building images and managing builders.
  icon: terminal
  link: /reference/cli/docker/buildx/
aliases:
- /buildx/working-with-buildx/
- /develop/develop-images/build_enhancements/
---

Docker Build is one of Docker Engine's most used features. Whenever you are
creating an image you are using Docker Build. Build is a key part of your
software development life cycle allowing you to package and bundle your code and
ship it anywhere.

Docker Build is more than a command for building images, and it's not only about
packaging your code. It's a whole ecosystem of tools and features that support
not only common workflow tasks but also provides support for more complex and
advanced scenarios.

## Explore the docs

{{< grid >}}
