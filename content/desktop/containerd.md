---
title: containerd image store
description: How to activate the containerd integration feature in Docker Desktop
keywords: Docker, containerd, engine, image store, lazy-pull
toc_max: 3
---

This page provides information about the ongoing integration of `containerd` for
image and file system management in the Docker Engine.

> **Note**
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
To use it, you must manually enable it in Docker Desktop settings.

Before you switch from the classic image store to containerd,
consider pruning your existing containers and images first.
Images and containers in the classic store won't be visible to you once you switch,
but they aren't deleted automatically.
This is to prevent unused resources from taking up disk space on your machine.
Once you switch over to the containerd image store,
you can't delete images and containers in the classic store unless you switch back.

Note that pruning images and containers means they're permanently lost.
If you have images or containers that you would like to keep,
you can manually migrate them to the image store, see [Migrate data between image stores](#migrate-data-between-image-stores).

To remove containers and images:

1. Stop any running containers.

   ```console
   $ docker stop $(docker ps -q)
   ```

2. Remove all containers:

   ```console
   $ docker container prune -a
   ```

3. Remove all images:

   ```console
   $ docker image prune -a
   ```

To switch to the containerd image store:

1. Open Docker Desktop and navigate to **Settings**.
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

## Migrate data between image stores

If you have container or image data that you would like to keep when migrating to a different image store,
you can export or save them as tar archives to the filesystem before pruning them.
Note that it isn't generally recommended to keep pet images and containers around in the image store
without an easy way to recreate them. If possible, you should prefer to have your
images defined in [Dockerfiles](../build/dockerfile/frontend.md).
For long-running containers and containers with custom startup parameters,
consider creating a shell script that hold the necessary arguments, or use [Docker Compose](../compose/_index.md).

Volumes, networks, and other resources are not affected by changing the image store,
only images and containers would need to be migrated or recreated.

To migrate images and containers between image stores,
you can export the data, as an image, to the local filesystem, before reimporting it.
For containers, this means first creating a new image based on an existing container,
and then exporting the image with `docker save`.

```console
$ docker ps
CONTAINER ID   IMAGE                  COMMAND                  CREATED         STATUS         PORTS      NAMES
e81d433c0cb7   postgres:16.3-alpine   "docker-entrypoint.s…"   2 seconds ago   Up 2 seconds   5432/tcp   some-postgres
$ docker container commit some-postgres my-pet-container:local
sha256:a08967fd5cb30061b4d205b6a61bb687a54b4ec966c4b4e86bec686a4d332109
$ docker image save --output "my-pet-container.tar" my-pet-container:local
```

When you're done saving the containers and images that you would like to keep,
stop your running containers and prune the resources:

1. Stop any running containers.

   ```console
   $ docker stop $(docker ps -q)
   ```

2. Remove all containers:

   ```console
   $ docker container prune -a
   ```

3. Remove all images:

   ```console
   $ docker image prune -a
   ```

After enabling the containerd image store, import the tar archive using the
[`docker image import`](../reference/cli/docker/image/import.md) command.

```console
$ docker image import \
  -m "Import local container backup" \
  my-pet-container.tar \
  my-pet-container:local
```

## Feedback

Thanks for trying the new features available with `containerd`. Give feedback or
report any bugs you may find through the issues tracker on the
[feedback form](https://dockr.ly/3PODIhD).
