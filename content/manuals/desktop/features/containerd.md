---
title: containerd image store
weight: 80
description: Learn about the containerd image store in Docker Desktop and how it extends image management capabilities.
keywords: Docker, containerd, engine, image store, lazy-pull
toc_max: 3
aliases:
  - /desktop/containerd/
---

Docker Desktop uses containerd as its image store by default. The image store
is the component responsible for pushing, pulling, and storing images on your
filesystem. The containerd image store supports features like multi-platform
images, image attestations, and alternative snapshotters.

## What is `containerd`?

`containerd` is a container runtime that provides a lightweight, consistent
interface for container lifecycle and image management. It is used under the
hood by Docker Engine for creating, starting, and stopping containers.

## What is the `containerd` image store?

The image store is the component responsible for pushing, pulling,
and storing images on the filesystem.

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

## Classic image store

The classic image store is Docker's legacy storage backend, replaced by the
containerd image store. It doesn't support image indices or manifest lists, so
you can't load multi-platform images locally or build images with attestations.

Most users have no reason to use the classic image store. It's available for
cases where you need to match older behavior or have compatibility
requirements.

## Switch image stores

The containerd image store is enabled by default in Docker Desktop version 4.34
and later. To switch between image stores:

1. Navigate to **Settings** in Docker Desktop.
2. In the **General** tab, check or clear the **Use containerd for pulling and storing images** option.
3. Select **Apply**.

> [!NOTE]
>
> Docker Desktop maintains separate image stores for the classic and containerd image stores.
> When switching between them, images and containers from the inactive store remain on disk but are hidden until you switch back.

## Build multi-platform images

The containerd image store lets you build multi-platform images
and load them to your local image store:

<script async id="asciicast-ZSUI4Mi2foChLjbevl2dxt5GD" src="https://asciinema.org/a/ZSUI4Mi2foChLjbevl2dxt5GD.js"></script>

Building multi-platform images with the classic image store is not supported:

```console
$ docker build --platform=linux/amd64,linux/arm64 .
[+] Building 0.0s (0/0)
ERROR: Multi-platform build is not supported for the docker driver.
Switch to a different driver, or turn on the containerd image store, and try again.
Learn more at https://docs.docker.com/go/build-multi-platform/
```
