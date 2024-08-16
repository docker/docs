---
title: containerd image store
description: How to activate the containerd integration feature in Docker Desktop
keywords: Docker, containerd, engine, image store, lazy-pull
toc_max: 3
---

This page provides information about the ongoing integration of `containerd` for
image and file system management in the Docker Engine.

> [!NOTE]
> 
> After switching to the containerd image store,
> images and containers in the classic image store won't be visible.
> All of those containers and images still exist.
> To see them again, turn off the containerd image store feature.

## What is containerd?

`containerd` is an abstraction of the low-level kernel features
used to run and manage containers on a system.
It's a platform used in container software like Docker and Kubernetes.

Docker Engine already uses `containerd` for container lifecycle management,
which includes creating, starting, and stopping containers.
This page describes the next step of the containerd integration for Docker:
the containerd image store.

## Image store

The image store is the component responsible for pushing, pulling,
and storing images on the filesystem.
The classic Docker image store is limited in the types of images that it supports.
For example, it doesn't support image indices, containing manifest lists.
When you create multi-platform images, for example,
the image index resolves all the platform-specific variants of the image.
An image index is also required when building images with attestations.

The containerd image store extends range of image types
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

The containerd image store isn't enabled by default.
To enable the feature for Docker Desktop:

1. Navigate to **Settings** in Docker Desktop.
2. In the **General** tab, check **Use containerd for pulling and storing images**.
3. Select **Apply & Restart**.

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

## Feedback

Thanks for trying the new features available with `containerd`. Give feedback or
report any bugs you may find through the issues tracker on the
[feedback form](https://dockr.ly/3PODIhD).
