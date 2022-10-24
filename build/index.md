---
title: Overview of Docker Build
description: Introduction and overview of Docker Build
keywords: build, buildx, buildkit
redirect_from:
- /build/buildx/
- /buildx/working-with-buildx/
- /develop/develop-images/build_enhancements/
---

## Overview

Docker Build is one of Docker Engine's most used features. Whenever you are
creating an image you are using Docker Build. Build is a key part of your
software development life cycle allowing you to package and bundle your code
and ship it anywhere.

Engine uses a client-server architecture and is composed of multiple components
and tools. The most common method of executing a build is by issuing a
[`docker build` command](../engine/reference/commandline/build.md). The CLI
sends the request to Docker Engine which, in turn, executes your build.

There are now two components in Engine that can be used to build an image.
Starting with the [18.09 release](../engine/release-notes/18.09.md#18090),
Engine is shipped with Moby [BuildKit](buildkit/index.md), the new component for
executing your builds by default.

The new client [Docker Buildx](https://github.com/docker/buildx){:target="_blank" rel="noopener" class="_"},
is a CLI plugin that extends the docker command with the full support of the
features provided by [BuildKit](buildkit/index.md) builder toolkit. [`docker buildx build` command](../engine/reference/commandline/buildx_build.md)
provides the same user experience as `docker build` with many new features like
creating scoped [builder instances](building/drivers/index.md), building
against multiple nodes concurrently, outputs configuration, inline [build caching](building/cache/index.md),
and specifying target platform. In addition, Buildx also supports new features
that are not yet available for regular `docker build` like building manifest
lists, distributed caching, and exporting build results to OCI image tarballs.

Docker Build is way more than a simple build command and is not only about
packaging your code, it's a whole ecosystem of tools and features that support
not only common workflow tasks but also provides support for more complex and
advanced scenarios.

## Building your images

### Packaging your software

Build and package your application to run it anywhere: locally using Docker
Desktop, or in the cloud using Docker Engine and Kubernetes:

[Packaging your software](building/packaging.md){: .button .outline-btn }

### Choosing a build driver

Run Buildx with different configurations depending on the scenario you are
working on, regardless of whether you are using your local machine or a remote
cluster, all from the comfort of your local working environment:

[Choosing a build driver](building/drivers/index.md){: .button .outline-btn }

### Optimizing builds with cache

Improve build performance by using a persistent shared build cache to avoid
repeating costly operations such as package installs, file downloads, or code
build steps:

[Optimizing builds with cache](./building/cache/index.md){: .button .outline-btn }

### Multi-stage builds

Use the multi-stage feature to selectively copy artifacts from one stage to
another, leaving behind everything you don't want in the final image, so you
keep your images small and secure with minimal dependencies:

[Multi-stage builds](building/multi-stage.md){: .button .outline-btn }

### Multi-platform images

Using the standard Docker tooling and processes, you can start to build, push,
pull, and run images seamlessly on different computer architectures:

[Multi-platform images](building/multi-platform.md){: .button .outline-btn }

## Continuous integration

### GitHub Actions

Automate your image builds to run in GitHub actions using the official docker
build actions:

* [GitHub Action to build and push Docker images with Buildx](https://github.com/docker/build-push-action).
* [GitHub Action to extract metadata from Git reference and GitHub events](https://github.com/docker/metadata-action/).

## Customizing your builds

### Select your build output format

Choose from a variety of available output formats, to export any artifact you
like from BuildKit, not just docker images. See [Set the export action for the build result](../engine/reference/commandline/buildx_build.md#output).

### Managing build secrets

Securely access protected repositories and resources at build time without
leaking data into the final build or the cache.

### High-level builds with Bake

Bake provides support for high-level build concepts using a file definition
that goes beyond invoking a single build command. Bake allows all the services
to be built concurrently as part of a single request:

[High-level builds with Bake](customize/bake/index.md){: .button .outline-btn }

## BuildKit

### Custom Dockerfile syntax

Use experimental versions of the Dockerfile frontend, or even just bring your
own to BuildKit using the power of custom frontends.

[Custom Dockerfile syntax](buildkit/dockerfile-frontend.md){: .button .outline-btn }

### Configure BuildKit

Take a deep dive into the internal BuildKit configuration to get the most out
of your builds. See also [`buildkitd.toml`](https://github.com/moby/buildkit/blob/master/docs/buildkitd.toml.md),
the configuration file for `buildkitd`.
