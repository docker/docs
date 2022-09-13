---
title: Multi-platform images
description: Different strategies for building multi-platform images
keywords: build, buildx, buildkit, multi-platform images
redirect_from:
- /build/buildx/multiplatform-images/
- /docker-for-mac/multi-arch/
- /mackit/multi-arch/
---

Docker images can support multiple platforms, which means that a single image
may contain variants for different architectures, and sometimes for different
operating systems, such as Windows.

When running an image with multi-platform support, `docker` automatically
selects the image that matches your OS and architecture.

Most of the Docker Official Images on Docker Hub provide a [variety of architectures](https://github.com/docker-library/official-images#architectures-other-than-amd64){: target="_blank" rel="noopener" class="_" }.
For example, the `busybox` image supports `amd64`, `arm32v5`, `arm32v6`,
`arm32v7`, `arm64v8`, `i386`, `ppc64le`, and `s390x`. When running this image
on an `x86_64` / `amd64` machine, the `amd64` variant is pulled and run.

## Building multi-platform images

Docker is now making it easier than ever to develop containers on, and for Arm
servers and devices. Using the standard Docker tooling and processes, you can
start to build, push, pull, and run images seamlessly on different compute
architectures. In most cases, you don't have to make any changes to Dockerfiles
or source code to start building for Arm.

BuildKit with Buildx is designed to work well for building for multiple
platforms and not only for the architecture and operating system that the user
invoking the build happens to run.

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

## Getting started

Run the [`docker buildx ls` command](../../engine/reference/commandline/buildx_ls.md)
to list the existing builders:

```console
$ docker buildx ls
NAME/NODE  DRIVER/ENDPOINT  STATUS   BUILDKIT PLATFORMS
default *  docker
  default  default          running  20.10.17 linux/amd64, linux/arm64, linux/arm/v7, linux/arm/v6
```

This displays the default builtin driver, that uses the BuildKit server
components built directly into the docker engine, also known as the [`docker` driver](../buildx/drivers/docker.md).

Create a new builder using the [`docker-container` driver](../buildx/drivers/docker-container.md)
which gives you access to more complex features like multi-platform builds
and the more advanced cache exporters, which are currently unsupported in the
default `docker` driver:

```console
$ docker buildx create --name mybuilder --driver docker-container --bootstrap
mybuilder
```

Switch to the new builder and inspect it:

```console
$ docker buildx use mybuilder
```

> **Note**
>
> Alternatively, run `docker buildx create --name mybuilder --driver docker-container --bootstrap --use`
> to create a new builder and switch to it using a single command.

And inspect it:

```console
$ docker buildx inspect
Name:   mybuilder
Driver: docker-container

Nodes:
Name:      mybuilder0
Endpoint:  unix:///var/run/docker.sock
Status:    running
Buildkit:  v0.10.4
Platforms: linux/amd64, linux/amd64/v2, linux/amd64/v3, linux/arm64, linux/riscv64, linux/ppc64le, linux/s390x, linux/386, linux/mips64le, linux/mips64, linux/arm/v7, linux/arm/v6
```

Now listing the existing builders again, we can see our new builder is
registered:

```console
$ docker buildx ls
NAME/NODE     DRIVER/ENDPOINT              STATUS   BUILDKIT PLATFORMS
mybuilder     docker-container
  mybuilder0  unix:///var/run/docker.sock  running  v0.10.4  linux/amd64, linux/amd64/v2, linux/amd64/v3, linux/arm64, linux/riscv64, linux/ppc64le, linux/s390x, linux/386, linux/mips64le, linux/mips64, linux/arm/v7, linux/arm/v6
default *     docker
  default     default                      running  20.10.17 linux/amd64, linux/arm64, linux/arm/v7, linux/arm/v6
```

## Example

Test the workflow to ensure you can build, push, and run multi-platform images.
Create a simple example Dockerfile, build a couple of image variants, and push
them to Docker Hub.

The following example uses a single `Dockerfile` to build an Alpine image with
cURL installed for multiple architectures:

```dockerfile
# syntax=docker/dockerfile:1
FROM alpine:3.16
RUN apk add curl
```

Build the Dockerfile with buildx, passing the list of architectures to
build for:

```console
$ docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t <username>/<image>:latest --push .
...
#16 exporting to image
#16 exporting layers
#16 exporting layers 0.5s done
#16 exporting manifest sha256:71d7ecf3cd12d9a99e73ef448bf63ae12751fe3a436a007cb0969f0dc4184c8c 0.0s done
#16 exporting config sha256:a26f329a501da9e07dd9cffd9623e49229c3bb67939775f936a0eb3059a3d045 0.0s done
#16 exporting manifest sha256:5ba4ceea65579fdd1181dfa103cc437d8e19d87239683cf5040e633211387ccf 0.0s done
#16 exporting config sha256:9fcc6de03066ac1482b830d5dd7395da781bb69fe8f9873e7f9b456d29a9517c 0.0s done
#16 exporting manifest sha256:29666fb23261b1f77ca284b69f9212d69fe5b517392dbdd4870391b7defcc116 0.0s done
#16 exporting config sha256:92cbd688027227473d76e705c32f2abc18569c5cfabd00addd2071e91473b2e4 0.0s done
#16 exporting manifest list sha256:f3b552e65508d9203b46db507bb121f1b644e53a22f851185d8e53d873417c48 0.0s done
#16 ...

#17 [auth] <username>/<image>:pull,push token for registry-1.docker.io
#17 DONE 0.0s

#16 exporting to image
#16 pushing layers
#16 pushing layers 3.6s done
#16 pushing manifest for docker.io/<username>/<image>:latest@sha256:f3b552e65508d9203b46db507bb121f1b644e53a22f851185d8e53d873417c48
#16 pushing manifest for docker.io/<username>/<image>:latest@sha256:f3b552e65508d9203b46db507bb121f1b644e53a22f851185d8e53d873417c48 1.4s done
#16 DONE 5.6s
```

> **Note**
> 
> * `<username>` must be a valid Docker ID and `<image>` and valid repository on
>   Docker Hub.
> * The `--platform` flag informs buildx to create Linux images for AMD 64-bit,
>   Arm 64-bit, and Armv7 architectures.
> * The `--push` flag generates a multi-arch manifest and pushes all the images
>   to Docker Hub.

Inspect the image using [`docker buildx imagetools` command](../../engine/reference/commandline/buildx_imagetools.md):

```console
$ docker buildx imagetools inspect <username>/<image>:latest
Name:      docker.io/<username>/<image>:latest
MediaType: application/vnd.docker.distribution.manifest.list.v2+json
Digest:    sha256:f3b552e65508d9203b46db507bb121f1b644e53a22f851185d8e53d873417c48

Manifests:
  Name:      docker.io/<username>/<image>:latest@sha256:71d7ecf3cd12d9a99e73ef448bf63ae12751fe3a436a007cb0969f0dc4184c8c
  MediaType: application/vnd.docker.distribution.manifest.v2+json
  Platform:  linux/amd64

  Name:      docker.io/<username>/<image>:latest@sha256:5ba4ceea65579fdd1181dfa103cc437d8e19d87239683cf5040e633211387ccf
  MediaType: application/vnd.docker.distribution.manifest.v2+json
  Platform:  linux/arm64

  Name:      docker.io/<username>/<image>:latest@sha256:29666fb23261b1f77ca284b69f9212d69fe5b517392dbdd4870391b7defcc116
  MediaType: application/vnd.docker.distribution.manifest.v2+json
  Platform:  linux/arm/v7
```

The image is now available on Docker Hub with the tag `<username>/<image>:latest`.
You can use this image to run a container on Intel laptops, Amazon EC2 Graviton
instances, Raspberry Pis, and on other architectures. Docker pulls the correct
image for the current architecture, so Raspberry PIs run the 32-bit Arm version
and EC2 Graviton instances run 64-bit Arm.

The digest identifies a fully qualified image variant. You can also run images
targeted for a different architecture on Docker Desktop. For  example, when
you run the following on a macOS:

 ```console
$ docker run --rm docker.io/<username>/<image>:latest@sha256:2b77acdfea5dc5baa489ffab2a0b4a387666d1d526490e31845eb64e3e73ed20 uname -m
aarch64
```

```console
$ docker run --rm docker.io/<username>/<image>:latest@sha256:723c22f366ae44e419d12706453a544ae92711ae52f510e226f6467d8228d191 uname -m
armv7l
```

In the above example, `uname -m` returns `aarch64` and `armv7l` as expected,
even when running the commands on a native macOS or Windows developer machine.

## Support on Docker Desktop

[Docker Desktop](../../desktop/index.md) provides `binfmt_misc`
multi-architecture support, which means you can run containers for different
Linux architectures such as `arm`, `mips`, `ppc64le`, and even `s390x`.

This does not require any special configuration in the container itself as it
uses [qemu-static](https://wiki.qemu.org/Main_Page){: target="_blank" rel="noopener" class="_" }
from the **Docker for Mac VM**. Because of this, you can run an ARM container,
like the `arm32v7` or `ppc64le` variants of the busybox image.
