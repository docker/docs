---
title: Docker Build
weight: 20
description: Get an overview of Docker Build to package and bundle your code and ship it anywhere
keywords: build, buildx, buildkit
params:
  sidebar:
    group: Open source
grid:
- title: Packaging your software
  description: 'Build and package your application to run it anywhere: locally or
    in the cloud.'
  icon: inventory_2
  link: /build/concepts/overview/
- title: Multi-stage builds
  description: Keep your images small and secure with minimal dependencies.
  icon: stairs
  link: /build/building/multi-stage/
- title: Multi-platform images
  description: Build, push, pull, and run images seamlessly on different computer
    architectures.
  icon: content_copy
  link: /build/building/multi-platform/
- title: BuildKit
  description: Explore BuildKit, the open source build engine.
  icon: construction
  link: /build/buildkit/
- title: Build drivers
  description: Configure where and how you run your builds.
  icon: engineering
  link: /build/builders/drivers/
- title: Exporters
  description: Export any artifact you like, not just Docker images.
  icon: output
  link: /build/exporters/
- title: Build caching
  description: Avoid unnecessary repetitions of costly operations, such as package
    installs.
  icon: cycle
  link: /build/cache/
- title: Bake
  description: Orchestrate your builds with Bake.
  icon: cake
  link: /build/bake/
aliases:
- /buildx/working-with-buildx/
- /develop/develop-images/build_enhancements/
---

Docker Build is one of Docker Engine's most used features. Every time you
create an image, you use Docker Build. Build is a key part of your software
development lifecycle that lets you package and bundle your code and ship it
anywhere.

Docker Build is more than a command for building images. It's a complete
ecosystem of tools and features that supports common workflow tasks and
advanced scenarios.

{{< grid >}}
