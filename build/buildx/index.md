---
title: Buildx
description: Working with Docker Buildx
keywords: build, buildx, buildkit
redirect_from:
- /buildx/working-with-buildx/
---

## Overview

Docker Buildx is a CLI plugin that extends the docker command with the full
support of the features provided by [Moby BuildKit](https://github.com/moby/buildkit){:target="_blank" rel="noopener" class="_"}
builder toolkit. It provides the same user experience as docker build with many
new features like creating scoped builder instances and building against
multiple nodes concurrently.

## Build with Buildx

To start a new build, run the command `docker buildx build .`

```console
$ docker buildx build .
[+] Building 8.4s (23/32)
 => ...
```

Buildx builds using the BuildKit engine and does not require `DOCKER_BUILDKIT=1`
environment variable to start the builds.

The [`docker buildx build` command](../../engine/reference/commandline/buildx_build.md)
supports features available for `docker build`, including features such as
outputs configuration, inline build caching, and specifying target platform.
In addition, Buildx also supports new features that are not yet available for
regular `docker build` like building manifest lists, distributed caching, and
exporting build results to OCI image tarballs.

You can run Buildx in different configurations that are exposed through a driver
concept. Currently, Docker supports a "docker" driver that uses the BuildKit
library bundled into the Docker daemon binary, and a "docker-container" driver
that automatically launches BuildKit inside a Docker container.

The user experience of using Buildx is very similar across drivers. However,
there are some features that are not currently supported by the "docker" driver,
because the BuildKit library which is bundled into docker daemon uses a different
storage component. In contrast, all images built with the "docker" driver are
automatically added to the "docker images" view by default, whereas when using
other drivers, the method for outputting an image needs to be selected
with `--output`.

## Work with builder instances

By default, Buildx uses the `docker` driver if it is supported, providing a user
experience very similar to the native `docker build`. Note that you must use a
local shared daemon to build your applications.

Buildx allows you to create new instances of isolated builders. You can use this
to get a scoped environment for your CI builds that does not change the state of
the shared daemon, or for isolating builds for different projects. You can create
a new instance for a set of remote nodes, forming a build farm, and quickly
switch between them.

You can create new instances using the [`docker buildx create`](../../engine/reference/commandline/buildx_create.md)
command. This creates a new builder instance with a single node based on your
current configuration.

To use a remote node you can specify the `DOCKER_HOST` or the remote context name
while creating the new builder. After creating a new instance, you can manage its
lifecycle using the [`docker buildx inspect`](../../engine/reference/commandline/buildx_inspect.md),
[`docker buildx stop`](../../engine/reference/commandline/buildx_stop.md), and
[`docker buildx rm`](../../engine/reference/commandline/buildx_rm.md) commands.
To list all available builders, use [`docker buildx ls`](../../engine/reference/commandline/buildx_ls.md).
After creating a new builder you can also append new nodes to it.

To switch between different builders, use [`docker buildx use <name>`](../../engine/reference/commandline/buildx_use.md).
After running this command, the build commands will automatically use this
builder.

Docker also features a [`docker context`](../../engine/reference/commandline/context.md)
command that you can use to provide names for remote Docker API endpoints. Buildx
integrates with `docker context` to ensure all the contexts automatically get a
default builder instance. You can also set the context name as the target when
you create a new builder instance or when you add a node to it.

## Build multi-platform images

BuildKit is designed to work well for building for multiple platforms and not
only for the architecture and operating system that the user invoking the build
happens to run.

When you invoke a build, you can set the `--platform` flag to specify the target
platform for the build output, (for example, `linux/amd64`, `linux/arm64`, or
`darwin/amd64`).

When the current builder instance is backed by the `docker-container` driver,
you can specify multiple platforms together. In this case, it builds a manifest
list which contains images for all specified architectures. When you use this
image in [`docker run`](../../engine/reference/commandline/run.md) or
[`docker service`](../../engine/reference/commandline/service.md), Docker picks
the correct image based on the node's platform.

You can build multi-platform images using three different strategies that are
supported by Buildx and Dockerfiles:

1. Using the QEMU emulation support in the kernel
2. Building on multiple native nodes using the same builder instance
3. Using a stage in Dockerfile to cross-compile to different architectures

QEMU is the easiest way to get started if your node already supports it (for
example. if you are using Docker Desktop). It requires no changes to your
Dockerfile and BuildKit automatically detects the secondary architectures that
are available. When BuildKit needs to run a binary for a different architecture,
it automatically loads it through a binary registered in the `binfmt_misc`
handler.

For QEMU binaries registered with `binfmt_misc` on the host OS to work
transparently inside containers, they must be statically compiled and registered
with the `fix_binary` flag. This requires a kernel >= 4.8 and
binfmt-support >= 2.1.7. You can check for proper registration by checking if
`F` is among the flags in `/proc/sys/fs/binfmt_misc/qemu-*`. While Docker
Desktop comes preconfigured with `binfmt_misc` support for additional platforms,
for other installations it likely needs to be installed using
[`tonistiigi/binfmt`](https://github.com/tonistiigi/binfmt){:target="_blank" rel="noopener" class="_"}
image.

```console
$ docker run --privileged --rm tonistiigi/binfmt --install all
```

Using multiple native nodes provide better support for more complicated cases
that are not handled by QEMU and generally have better performance. You can
add additional nodes to the builder instance using the `--append` flag.

Assuming contexts `node-amd64` and `node-arm64` exist in `docker context ls`;

```console
$ docker buildx create --use --name mybuild node-amd64
mybuild
$ docker buildx create --append --name mybuild node-arm64
$ docker buildx build --platform linux/amd64,linux/arm64 .
```

Finally, depending on your project, the language that you use may have good
support for cross-compilation. In that case, multi-stage builds in Dockerfiles
can be effectively used to build binaries for the platform specified with
`--platform` using the native architecture of the build node. A list of build
arguments like `BUILDPLATFORM` and `TARGETPLATFORM` is available automatically
inside your Dockerfile and can be leveraged by the processes running as part
of your build.

```dockerfile
# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM golang:alpine AS build
ARG TARGETPLATFORM
ARG BUILDPLATFORM
RUN echo "I am running on $BUILDPLATFORM, building for $TARGETPLATFORM" > /log
FROM alpine
COPY --from=build /log /log
```

## High-level build options

Buildx also aims to provide support for high-level build concepts that go beyond
invoking a single build command.

BuildKit efficiently handles multiple concurrent build requests and de-duplicating
work. The build commands can be combined with general-purpose command runners
(for example, `make`). However, these tools generally invoke builds in sequence
and therefore cannot leverage the full potential of BuildKit parallelization,
or combine BuildKit’s output for the user. For this use case, we have added a
command called [`docker buildx bake`](../../engine/reference/commandline/buildx_bake.md).

The `bake` command supports building images from compose files, similar to
[`docker-compose build`](../../engine/reference/commandline/compose_build.md),
but allowing all the services to be built concurrently as part of a single
request.

There is also support for custom build rules from HCL/JSON files allowing
better code reuse and different target groups. The design of bake is in very
early stages, and we are looking for feedback from users. Let us know your
feedback by creating an issue in the [Docker Buildx](https://github.com/docker/buildx/issues){:target="_blank" rel="noopener" class="_"}
GitHub repository.
