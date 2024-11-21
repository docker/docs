---
title: Base images
weight: 70
description: Learn about base images and how they're created
keywords: images, base image, examples
aliases:
- /articles/baseimages/
- /engine/articles/baseimages/
- /engine/userguide/eng-image/baseimages/
- /develop/develop-images/baseimages/
---

All Dockerfiles start from a base image.
A base is the image that your image extends.
It refers to the contents of the `FROM` instruction in the Dockerfile.

```dockerfile
FROM debian
```

For most cases, you don't need to create your own base image. Docker Hub
contains a vast library of Docker images that are suitable for use as a base
image in your build. [Docker Official
Images](../../docker-hub/image-library/trusted-content.md#docker-official-images)
are specifically designed as a set of hardened, battle-tested images that
support a wide variety of platforms, languages, and frameworks. There are also
[Docker Verified
Publisher](../../docker-hub/image-library/trusted-content.md#verified-publisher-images)
images, created by trusted publishing partners, verified by Docker.

## Create a base image

If you need to completely control the contents of your image, you can create
your own base image from a Linux distribution of your choosing, or use the
special `FROM scratch` base:

```dockerfile
FROM scratch
```

The `scratch` image is typically used to create minimal images containing only
just what an application needs. See [Create a minimal base image using scratch](#create-a-minimal-base-image-using-scratch).

To create a distribution base image, you can use a root filesystem, packaged as
a `tar` file, and import it to Docker with `docker import`. The process for
creating your own base image depends on the Linux distribution you want to
package. See [Create a full image using tar](#create-a-full-image-using-tar).

## Create a minimal base image using scratch

The reserved, minimal `scratch` image serves as a starting point for
building containers. Using the `scratch` image signals to the build process
that you want the next command in the `Dockerfile` to be the first filesystem
layer in your image.

While `scratch` appears in Docker's [repository on Docker Hub](https://hub.docker.com/_/scratch),
you can't pull it, run it, or tag any image with the name `scratch`.
Instead, you can refer to it in your `Dockerfile`.
For example, to create a minimal container using `scratch`:

```dockerfile
# syntax=docker/dockerfile:1
FROM scratch
ADD hello /
CMD ["/hello"]
```

Assuming an executable binary named `hello` exists at the root of the [build context](/manuals/build/concepts/context.md).
You can build this Docker image using the following `docker build` command:

```console
$ docker build --tag hello .
```

To run your new image, use the `docker run` command:

```console
$ docker run --rm hello
```

This example image can only successfully execute as long as the `hello` binary
doesn't have any runtime dependencies. Computer programs tend to depend on
certain other programs or resources to exist in the runtime environment. For
example:

- Programming language runtimes
- Dynamically linked C libraries
- CA certificates

When building a base image, or any image, this is an important aspect to
consider. And this is why creating a base image using `FROM scratch` can be
difficult, for anything other than small, simple programs. On the other hand,
it's also important to include only the things you need in your image, to
reduce the image size and attack surface.

## Create a full image using tar

In general, start with a working machine that is running
the distribution you'd like to package as a base image, though that is
not required for some tools like Debian's [Debootstrap](https://wiki.debian.org/Debootstrap),
which you can also use to build Ubuntu images.

For example, to create an Ubuntu base image:

```dockerfile
$ sudo debootstrap focal focal > /dev/null
$ sudo tar -C focal -c . | docker import - focal

sha256:81ec9a55a92a5618161f68ae691d092bf14d700129093158297b3d01593f4ee3

$ docker run focal cat /etc/lsb-release

DISTRIB_ID=Ubuntu
DISTRIB_RELEASE=20.04
DISTRIB_CODENAME=focal
DISTRIB_DESCRIPTION="Ubuntu 20.04 LTS"
```

There are more example scripts for creating base images in
[the Moby GitHub repository](https://github.com/moby/moby/blob/master/contrib).

## More resources

For more information about building images and writing Dockerfiles, see:

* [Dockerfile reference](/reference/dockerfile.md)
* [Dockerfile best practices](/manuals/build/building/best-practices.md)
* [Docker Official Images](../../docker-hub/image-library/trusted-content.md#docker-official-images)
