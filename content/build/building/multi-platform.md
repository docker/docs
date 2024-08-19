---
title: Multi-platform images
description: Introduction to multi-platform images and how to build them
keywords: build, buildx, buildkit, multi-platform images
aliases:
- /build/buildx/multiplatform-images/
- /desktop/multi-arch/
- /docker-for-mac/multi-arch/
- /mackit/multi-arch/
---

A multi-platform image refers to a single image that includes variants for
multiple different architectures and, in some cases, different operating
systems, like Windows. This means that whether you are using an ARM-based
system or an x86 machine, Docker automatically detects and selects the
appropriate variant for your hosts's operating system and architecture.

Many of the Docker Official Images available on Docker Hub support various
architectures. For instance, the `busybox` image includes support for these
platforms:

- x86-64 (`linux/amd64`, `linux/i386`)
- ARM architectures (`linux/arm/v5`, `linux/arm/v6`, `linux/arm/v7`, `linux/arm64`)
- PowerPC and IBM Z (`linux/ppc64le`, `linux/s390x`)

On an x86 machine, Docker will automatically use the `linux/amd64` variant
when you run a container or invoke a build.

Most Docker images use the `linux/` OS prefix to indicate they are Linux-based.
While Docker Desktop on macOS or Windows typically runs Linux containers using
a Linux VM, Docker also supports Windows containers if you're operating in
Windows container mode.

## Building multi-platform images

When triggering a build, use the `--platform` flag to define the target
platforms for the build output, such as `linux/amd64` and `linux/arm64`:

```console
$ docker build --platform linux/amd64,linux/arm64 .
```

By default, Docker can build for only one platform at a time.
To build for multiple platforms concurrently, you can:

- **Enable the containerd image store**:
  The default image store in Docker Engine doesn't support multi-platform images.
  The containerd image store does, and lets you create multi-platform images using the default builder.
  Refer to the [containerd in Docker Desktop documentation](../../desktop/containerd.md).

- **Create a custom builder**:
  Initialize a [builder](../builders/_index.md) that uses the `docker-container` driver, which supports multi-platform builds.
  For more details, see the [`docker-container` driver documentation](/build/builders/drivers/docker-container.md).

## Strategies

You can build multi-platform images using three different strategies,
depending on your use case:

1. Using emulation, via [QEMU](#qemu) support in the Linux kernel
2. Building on a single builder backed by
   [multiple nodes of different architectures](#multiple-native-nodes).
3. Using a stage in your Dockerfile to [cross-compile](#cross-compilation) to
   different architectures

### QEMU

Building multi-platform images under emulation with QEMU is the easiest way to
get started if your builder already supports it. Docker Desktop supports it out
of the box. It requires no changes to your Dockerfile, and BuildKit
automatically detects the secondary architectures that are available. When
BuildKit needs to run a binary for a different architecture, it automatically
loads it through a binary registered in the `binfmt_misc` handler.

> [!NOTE]
>
> Emulation with QEMU can be much slower than native builds, especially for
> compute-heavy tasks like compilation and compression or decompression.
>
> Use [multiple native nodes](#multiple-native-nodes) or
> [cross-compilation](#cross-compilation) instead, if possible.

#### Support on Docker Desktop

[Docker Desktop](../../desktop/_index.md) provides support for running and
building multi-platform images under emulation by default, which means you can
run containers for different Linux architectures such as `arm`, `mips`,
`ppc64le`, and even `s390x`.

This doesn't require any special configuration in the container itself as it
uses QEMU bundled within the Docker Desktop VM. Because of this, you can run
containers of non-native architectures like the `arm32v7` or `ppc64le`
automatically.

#### QEMU without Docker Desktop

If you're running Docker Engine on Linux, without Docker Desktop, you must
install statically compiled QEMU binaries and register them with
[`binfmt_misc`](https://en.wikipedia.org/wiki/Binfmt_misc). This enables QEMU
to execute non-native file formats for emulation. The QEMU binaries must be
statically compiled and registered with the `fix_binary` flag. This requires a
kernel version 4.8 or later, and `binfmt-support` version 2.1.7 or later.

Once QEMU is installed and the executable types are registered on the host OS,
they work transparently inside containers. You can verify your registration by
checking if `F` is among the flags in `/proc/sys/fs/binfmt_misc/qemu-*`. While
Docker Desktop comes preconfigured with `binfmt_misc` support for additional
platforms, for other installations it likely needs to be installed using
[`tonistiigi/binfmt`](https://github.com/tonistiigi/binfmt) image:

```console
$ docker run --privileged --rm tonistiigi/binfmt --install all
```

### Multiple native nodes

Using multiple native nodes provide better support for more complicated cases
that QEMU can't handle, and also provides better performance.

You can add additional nodes to a builder using the `--append` flag.

The following command creates a multi-node builder from Docker contexts named
`node-amd64` and `node-arm64`. This example assumes that you've already added
those contexts.

```console
$ docker buildx create --use --name mybuild node-amd64
mybuild
$ docker buildx create --append --name mybuild node-arm64
$ docker buildx build --platform linux/amd64,linux/arm64 .
```

While this approach has advantages over emulation, managing multi-node builders
introduces some overhead of setting up and managing builder clusters.
Alternatively, you can use Docker Build Cloud, a service that provides managed
multi-node builders on Docker's infrastructure. With Docker Build Cloud, you
get native multi-platform ARM and X86 builders without the burden of
maintaining them. Using cloud builders also provides additional benefits, such
as a shared build cache.

After signing up for Docker Build Cloud, add the builder to your local
environment and start building.

```console
$ docker buildx create --driver cloud <ORG>/<BUILDER_NAME>
cloud-<ORG>-<BUILDER_NAME>
$ docker build \
  --builder cloud-<ORG>-<BUILDER_NAME> \
  --platform linux/amd64,linux/arm64,linux/arm/v7 \
  --tag <IMAGE_NAME> \
  --push .
```

For more information, see [Docker Build Cloud](/build-cloud/_index.md).

### Cross-compilation

Depending on your project, if the programming language you use has good support
for cross-compilation, you can leverage multi-stage builds to build binaries
for target platforms from the native architecture of the builder. Special build
arguments, such as `BUILDPLATFORM` and `TARGETPLATFORM`, are automatically
available for use in your Dockerfile.

In the following example, the `FROM` instruction is pinned to the native
platform of the builder (using the `--platform=$BUILDPLATFORM` option) to
prevent emulation from kicking in. Then the pre-defined `$BUILDPLATFORM` and
`$TARGETPLATFORM` build arguments are interpolated in a `RUN` instruction. In
this case, the values are just printed to stdout with `echo`, but this
illustrates how you would pass them to the compiler for cross-compilation.

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

Run the [`docker buildx ls` command](../../reference/cli/docker/buildx/ls.md)
to list the existing builders:

```console
$ docker buildx ls
NAME/NODE  DRIVER/ENDPOINT  STATUS   BUILDKIT PLATFORMS
default *  docker
  default  default          running  v0.11.6  linux/amd64, linux/arm64, linux/arm/v7, linux/arm/v6
```

This displays the default builtin driver, that uses the BuildKit server
components built directly into the Docker Engine, also known as the [`docker` driver](/build/builders/drivers/docker.md).

Create a new builder using the [`docker-container` driver](/build/builders/drivers/docker-container.md)
which gives you access to more complex features like multi-platform builds
and the more advanced cache exporters, which are currently unsupported in the
default `docker` driver:

```console
$ docker buildx create --name mybuilder --bootstrap --use
```

Now listing the existing builders again, you can see that the new builder is
registered:

```console
$ docker buildx ls
NAME/NODE     DRIVER/ENDPOINT              STATUS   BUILDKIT PLATFORMS
mybuilder *   docker-container
  mybuilder0  unix:///var/run/docker.sock  running  v0.12.1  linux/amd64, linux/amd64/v2, linux/amd64/v3, linux/arm64, linux/riscv64, linux/ppc64le, linux/s390x, linux/386, linux/mips64le, linux/mips64, linux/arm/v7, linux/arm/v6
default       docker
  default     default                      running  v{{% param "buildkit_version" %}}  linux/amd64, linux/arm64, linux/arm/v7, linux/arm/v6
```

## Example

Test the workflow to ensure you can build, push, and run multi-platform images.
Create a simple example Dockerfile, build a couple of image variants, and push
them to Docker Hub.

The following example uses a single `Dockerfile` to build an Alpine image with
cURL installed for multiple architectures:

```dockerfile
# syntax=docker/dockerfile:1
FROM alpine:{{% param "example_alpine_version" %}}
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

> [!NOTE]
> 
> * `<username>` must be a valid Docker ID and `<image>` and valid repository on
>   Docker Hub.
> * The `--platform` flag informs buildx to create Linux images for x86 64-bit,
>   ARM 64-bit, and ARMv7 architectures.
> * The `--push` flag generates a multi-arch manifest and pushes all the images
>   to Docker Hub.

Inspect the image using [`docker buildx imagetools` command](../../reference/cli/docker/buildx/imagetools/_index.md):

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
image for the current architecture, so Raspberry PIs run the 32-bit ARM version
and EC2 Graviton instances run 64-bit ARM.

The digest identifies a fully qualified image variant. You can also run images
targeted for a different architecture on Docker Desktop. For example, when
you run the following on a macOS:

 ```console
$ docker run --rm docker.io/<username>/<image>:latest@sha256:2b77acdfea5dc5baa489ffab2a0b4a387666d1d526490e31845eb64e3e73ed20 uname -m
aarch64
```

```console
$ docker run --rm docker.io/<username>/<image>:latest@sha256:723c22f366ae44e419d12706453a544ae92711ae52f510e226f6467d8228d191 uname -m
armv7l
```

In the previous example, `uname -m` returns `aarch64` and `armv7l` as expected,
even when running the commands on a native macOS or Windows developer machine.
