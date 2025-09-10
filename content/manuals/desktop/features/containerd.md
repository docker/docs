---
title: containerd image store
weight: 80
description: How to activate the containerd integration feature in Docker Desktop
keywords: Docker, containerd, engine, image store, lazy-pull
toc_max: 3
aliases:
- /desktop/containerd/
---

Docker Desktop is transitioning to use containerd for image and filesystem management. This page outlines the benefits, setup process, and new capabilities enabled by the containerd image store.

> [!NOTE]
> 
> Docker Desktop maintains separate image stores for the classic and containerd image stores.
> When switching between them, images and containers from the inactive store remain on disk but are hidden until you switch back.

## What is `containerd`?

`containerd` is a container runtime that provides a lightweight, consistent interface for container lifecycle management. It is already used under the hood by Docker Engine for creating, starting, and stopping containers.

Docker Desktop’s ongoing integration of containerd now extends to the image store, offering more flexibility and modern image support.

## What is the `containerd` image store?

The image store is the component responsible for pushing, pulling,
and storing images on the filesystem.

The classic Docker image store is limited in the types of images that it supports.
For example, it doesn't support image indices, containing manifest lists.
When you create multi-platform images, for example,
the image index resolves all the platform-specific variants of the image.
An image index is also required when building images with attestations.

The containerd image store extends the range of image types
that the Docker Engine can natively interact with.
While this is a low-level architectural change,
it's a prerequisite for unlocking a range of new use cases, including:

- [Build multi-platform images](#build-multi-platform-images) and images with attestations
- Support for using containerd snapshotters with unique characteristics,
  such as [stargz][1] for lazy-pulling images on container startup,
  or [nydus][2] and [dragonfly][3] for peer-to-peer image distribution.
- Ability to run [Wasm](wasm.md) containers

[1]: https://github.com/containerd/stargz-snapshotter
[2]: https://github.com/containerd/nydus-snapshotter
[3]: https://github.com/dragonflyoss/image-service

## Enable the containerd image store

The containerd image store is enabled by default in Docker Desktop version 4.34
and later, but only for clean installs or if you perform a factory reset. If
you upgrade from an earlier version of Docker Desktop, or if you use an older
version of Docker Desktop you must manually switch to the containerd image
store.

To manually enable this feature in Docker Desktop:

1. Navigate to **Settings** in Docker Desktop.
2. In the **General** tab, check **Use containerd for pulling and storing images**.
3. Select **Apply**.

To disable the containerd image store,
clear the **Use containerd for pulling and storing images** checkbox.

## Build multi-platform images

The term multi-platform image refers to a bundle of images for multiple different architectures.
Out of the box, the default builder for Docker Desktop doesn't support building multi-platform images.

```console
$ docker build --platform=linux/amd64,linux/arm64 .
[+] Building 0.0s (0/0)
ERROR: Multi-platform build is not supported for the docker driver.
Switch to a different driver, or turn on the containerd image store, and try again.
Learn more at https://docs.docker.com/go/build-multi-platform/
```

Enabling the containerd image store lets you build multi-platform images
and load them to your local image store:

<script async id="asciicast-ZSUI4Mi2foChLjbevl2dxt5GD" src="https://asciinema.org/a/ZSUI4Mi2foChLjbevl2dxt5GD.js"></script>


