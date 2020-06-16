---
description: Multi-CPU Architecture Support
keywords: mac, Multi-CPU architecture support
redirect_from:
- /mackit/multi-arch/
title: Leverage multi-CPU architecture support
notoc: true
---
Docker images can support multiple architectures, which means that a single
image may contain variants for different architectures, and sometimes for different
operating systems, such as Windows.

When running an image with multi-architecture support, `docker` will
automatically select an image variant which matches your OS and architecture.

Most of the official images on Docker Hub provide a [variety of architectures](https://github.com/docker-library/official-images#architectures-other-than-amd64).
For example, the `busybox` image supports `amd64`, `arm32v5`, `arm32v6`,
`arm32v7`, `arm64v8`, `i386`, `ppc64le`, and `s390x`. When running this image
on an `x86_64` / `amd64` machine, the `x86_64` variant will be pulled and run.

**Docker Desktop** provides `binfmt_misc` multi-architecture support,
which means you can run containers for different Linux architectures
such as `arm`, `mips`, `ppc64le`, and even `s390x`.

This does not require any special configuration in the container itself as it uses
<a href="http://wiki.qemu.org/" target="_blank">qemu-static</a> from the **Docker for
Mac VM**. Because of this, you can run an ARM container, like the `arm32v7` or `ppc64le`
variants of the busybox image.

## Buildx (Experimental)

Docker is now making it easier than ever to develop containers on, and for Arm servers and devices. Using the standard Docker tooling and processes, you can start to build, push, pull, and run images seamlessly on different compute architectures. Note that you don't have to make any changes to Dockerfiles or source code to start building for Arm.

Docker introduces a new CLI command called `buildx`. You can use the `buildx` command on Docker Desktop for Mac and Windows to build multi-arch images, link them together with a manifest file, and push them all to a registry using a single command.  With the included emulation, you can transparently build more than just native images.  Buildx accomplishes this by adding new builder instances based on BuildKit, and leveraging Docker Desktop's technology stack to run non-native binaries.

For more information about the Buildx CLI command, see [Buildx](/buildx/working-with-buildx/).

### Install

1. Download the latest version of [Docker Desktop](https://hub.docker.com/editions/community/docker-ce-desktop-mac/).

1. Follow the on-screen instructions to complete the installation process. After you have successfully installed Docker Desktop, you will see the Docker icon in your task tray.

1. Click **About Docker Desktop** from the Docker menu and ensure you have installed Docker Desktop version 2.0.4.0 (33772) or higher.

![about-docker-desktop-buildx](./images/desktop-buildx-version.png)

### Build and run multi-architecture images

Run the command `docker buildx ls` to list the existing builders. This displays the default builder, which is our old builder.

```bash
$ docker buildx ls

NAME/NODE DRIVER/ENDPOINT STATUS  PLATFORMS
default * docker
  default default         running linux/amd64, linux/arm64, linux/arm/v7, linux/arm/v6
```

Create a new builder which gives access to the new multi-architecture features.

```bash
$ docker buildx create --name mybuilder

mybuilder
```

Alternatively, run `docker buildx create --name mybuilder --use` to create a new builder and switch to it using a single command.

Switch to the new builder and inspect it.

```bash
$ docker buildx use mybuilder

$ docker buildx inspect --bootstrap

[+] Building 2.5s (1/1) FINISHED
 => [internal] booting buildkit                                                   2.5s
 => => pulling image moby/buildkit:master                                         1.3s
 => => creating container buildx_buildkit_mybuilder0                              1.2s
Name:   mybuilder
Driver: docker-container

Nodes:
Name:      mybuilder0
Endpoint:  unix:///var/run/docker.sock
Status:    running

Platforms: linux/amd64, linux/arm64, linux/arm/v7, linux/arm/v6
```

Test the workflow to ensure you can build, push, and run multi-architecture images. Create a simple example Dockerfile, build a couple of image variants, and push them to Docker Hub.

```bash
$ mkdir test && cd test && cat <<EOF > Dockerfile

FROM ubuntu
RUN apt-get update && apt-get install -y curl
WORKDIR /src
COPY . .
EOF
```

```bash
$ docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t username/demo:latest --push .

[+] Building 6.9s (19/19) FINISHED
...
 => => pushing layers                                                             2.7s
 => => pushing manifest for docker.io/username/demo:latest                       2.2
 ```

Where, `username` is a valid Docker username.

>   **Notes:**
>
>  - The `--platform` flag informs buildx to generate Linux images for AMD 64-bit, Arm 64-bit, and Armv7 architectures.
>  - The `--push` flag generates a multi-arch manifest and pushes all the images to Docker Hub.

Inspect the image using `imagetools`.

```bash
$ docker buildx imagetools inspect username/demo:latest

Name:      docker.io/username/demo:latest
MediaType: application/vnd.docker.distribution.manifest.list.v2+json
Digest:    sha256:2a2769e4a50db6ac4fa39cf7fb300fa26680aba6ae30f241bb3b6225858eab76

Manifests:
  Name:      docker.io/username/demo:latest@sha256:8f77afbf7c1268aab1ee7f6ce169bb0d96b86f585587d259583a10d5cd56edca
  MediaType: application/vnd.docker.distribution.manifest.v2+json
  Platform:  linux/amd64

  Name:      docker.io/username/demo:latest@sha256:2b77acdfea5dc5baa489ffab2a0b4a387666d1d526490e31845eb64e3e73ed20
  MediaType: application/vnd.docker.distribution.manifest.v2+json
  Platform:  linux/arm64

  Name:      docker.io/username/demo:latest@sha256:723c22f366ae44e419d12706453a544ae92711ae52f510e226f6467d8228d191
  MediaType: application/vnd.docker.distribution.manifest.v2+json
  Platform:  linux/arm/v7
  ```

  The image is now available on Docker Hub with the tag `username/demo:latest`. You can use this image to run a container on Intel laptops, Amazon EC2 A1 instances, Raspberry Pis, and on other architectures. Docker pulls the correct image for the current architecture, so Raspberry Pis run the 32-bit Arm version and EC2 A1 instances run 64-bit Arm. The SHA tags identify a fully qualified image variant. You can also run images targeted for a different architecture on Docker Desktop.

  You can run the images using the SHA tag, and verify the architecture. For example, when you run the following on a macOS:

 ```bash
 $ docker run --rm docker.io/username/demo:latest@sha256:2b77acdfea5dc5baa489ffab2a0b4a387666d1d526490e31845eb64e3e73ed20 uname -m
 aarch64
```

```bash
$ docker run --rm docker.io/username/demo:latest@sha256:723c22f366ae44e419d12706453a544ae92711ae52f510e226f6467d8228d191 uname -m
armv7l
```

In the above example, `uname -m` returns `aarch64` and `armv7l` as expected, even when running the commands on a native macOS developer machine.
