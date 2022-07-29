---
title: Building multi-platform images
description: Different strategies for building multi-platform images
keywords: build, buildx, buildkit, multi-platform images
---

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
