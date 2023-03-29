---
title: Overview of Docker Build
description: Introduction and overview of Docker Build
keywords: build, buildx, buildkit
aliases:
  - /build/buildx/
  - /buildx/working-with-buildx/
  - /develop/develop-images/build_enhancements/
grid:
  - title: "Packaging your software"
    description:
      "Build and package your application to run it anywhere: locally or in the
      cloud."
    icon: "inventory_2"
    link: "/build/building/packaging"
  - title: "Multi-stage builds"
    description: "Keep your images small and secure with minimal dependencies."
    icon: "stairs"
    link: "/build/building/multi-stage"
  - title: "Multi-platform images"
    description:
      "Build, push, pull, and run images seamlessly on different computer
      architectures."
    icon: "home_storage"
    link: "/build/building/multi-platform/"
  - title: "Build drivers"
    description: "Configure where and how you run your builds."
    icon: "engineering"
    link: "/build/drivers/"
  - title: "Exporters"
    description: "Export any artifact you like, not just Docker images."
    icon: "output"
    link: "/build/exporters"
  - title: "Build caching"
    description:
      "Avoid unnecessary repetitions of costly operations, such as package
      installs."
    icon: "cycle"
    link: "/build/cache"
  - title: "Bake"
    description: "Orchestrate your builds with Bake."
    icon: "cake"
    link: "/build/bake"
  - title: "Continuous integration"
    description:
      "Learn how to use Docker in your continuous integration pipelines."
    icon: "all_inclusive"
    link: "/build/ci"
  - title: "Dockerfile frontend"
    description:
      "Learn about the Dockerfile frontend for BuildKit."
    icon: "all_inclusive"
    link: "/build/buildkit/dockerfile-frontend/"
  - title: "Configure BuildKit"
    description:
      "Take a deep dive into the internals of BuildKit to get the most out of your builds."
    icon: "all_inclusive"
    link: "/build/buildkit/configure/"
---

Docker Build is one of Docker Engine's most used features. Whenever you are
creating an image you are using Docker Build. Build is a key part of your
software development life cycle allowing you to package and bundle your code and
ship it anywhere.

The Docker Engine uses a client-server architecture and is composed of multiple components
and tools. The most common method of executing a build is by issuing a
[`docker build` command](../engine/reference/commandline/build.md). The CLI
sends the request to Docker Engine which, in turn, executes your build.

There are now two components in Engine that can be used to build an image.
Starting with the [18.09 release](../engine/release-notes/18.09.md#18090),
Engine is shipped with Moby [BuildKit](buildkit/index.md), the new component for
executing your builds by default.

The new client [Docker Buildx](https://github.com/docker/buildx){:target="blank" rel="noopener" class=""}
is a CLI plugin that extends the `docker` command with the full support of the
features provided by [BuildKit](buildkit/index.md) builder toolkit.
[`docker buildx build` command](../engine/reference/commandline/buildx_build.md)
provides the same user experience as `docker build` with many new features like
creating scoped [builder instances](drivers/index.md), building against
multiple nodes concurrently, outputs configuration, inline
[build caching](cache/index.md), and specifying target platform. In
addition, Buildx also supports new features that aren't yet available for
regular `docker build` like building manifest lists, distributed caching, and
exporting build results to OCI image tarballs.

Docker Build is more than a simple build command, and it's not only about
packaging your code. It's a whole ecosystem of tools and features that support
not only common workflow tasks but also provides support for more complex and
advanced scenarios.

{{< grid >}}
