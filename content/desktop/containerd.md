---
title: containerd image store (Beta)
description: How to activate the containerd integration feature in Docker Desktop
keywords: Docker, containerd, engine, image store, lazy-pull
toc_max: 3
---

This page provides information about the ongoing integration of `containerd` for
image and file system management in the Docker Engine.

> **Beta**
>
> The containerd image store feature is currently in
> [Beta](../release-lifecycle.md/#beta). We recommend that you do not use
> this feature in production environments as this feature may change or be
> removed from future releases.
{ .experimental }

## What is the containerd image store?

`containerd` is a container runtime that manages the container lifecycle, and
provides image and filesystem management. It's a low-level building block,
designed to be integrated into other systems, such as Docker and Kubernetes.

Docker Engine already uses `containerd` for container lifecycle management, which
includes creating, starting, and stopping containers. This page describes the
next step of containerd integration for Docker Engine: the image store.

The image store is the component responsible for pushing, pulling, and storing
images. Integrating the containerd image store enables many new features in the
Docker Engine, including:

- containerd snapshotters, such as [stargz][1] for lazy-pulling images on startup,
  or [nydus][2] and [dragonfly][3] for peer-to-peer image distribution.
- Natively store and build multi-platform images, and other OCI content types
  that may emerge in the future.
- Ability to run [Wasm](wasm.md) containers

[1]: https://github.com/containerd/stargz-snapshotter
[2]: https://github.com/containerd/nydus-snapshotter
[3]: https://github.com/dragonflyoss/image-service

The image store integration is still at an early stage, so not all features are
yet supported.

## Turn on the containerd image store feature

The containerd image store beta feature is turned off by default.

To start using the feature:

1. Navigate to **Settings**.
2. From the **Features in development** tab, select **Beta features**.
3. Next to **Use containerd for pulling and storing images**, select the
   checkbox.
4. Select **Apply & Restart**

To turn off this feature, clear the **Use containerd for pulling and storing 
images** checkbox.

> **Tip**
>
> After switching to the containerd image store, images and containers from the
> default image store won't be visible. All of those containers and images
> still exist. To see them again, turn off the containerd image store feature.
{ .tip }

## Building multi-platform images

The term multi-platform image refers to a bundle of images that can run on different architectures.
Out of the box, the default builder for Docker Desktop doesn't support building multi-platform images.

```console
$ docker buildx ls | grep "DRIVER\|*"
NAME/NODE       DRIVER/ENDPOINT             STATUS  BUILDKIT PLATFORMS
default *       docker
$ docker buildx build --platform=linux/amd64,linux/arm64 .
[+] Building 0.0s (0/0)
ERROR: multiple platforms feature is currently not supported for docker driver. Please switch to a different driver (eg. "docker buildx create --use")
```

Normally, building multi-platform images requires you to create a new builder,
using a driver that supports multi-platform builds.
But even then, you can't load the multi-platform images to your local image store.

```console
$ docker buildx create --bootstrap
[+] Building 2.4s (1/1) FINISHED
 => [internal] booting buildkit
 => => pulling image moby/buildkit:buildx-stable-1
 => => creating container buildx_buildkit_objective_blackburn0
objective_blackburn
$ docker buildx build --quiet \
  --platform=linux/amd64,linux/arm64 \
  --builder=objective_blackburn \
  --load .
ERROR: docker exporter does not currently support exporting manifest lists
```

Turning on the containerd image store lets you build, and load, multi-platform images
to your local image store, all while using the default builder.



```console
$ docker info --format="{{ .Driver }}"
stargz
$ docker buildx build \   
  --platform=linux/arm64,linux/amd64 \
  --tag=user/containerd-multiplatform .
[+] Building 6.2s (11/11) FINISHED                                                                                                                                                                       
 ...
 => [internal] load build definition from Dockerfile                            0.0s
 => => transferring dockerfile: 115B                                            0.0s
 => [linux/arm64 internal] load metadata for docker.io/library/alpine:latest    2.0s
 => [linux/amd64 internal] load metadata for docker.io/library/alpine:latest    2.1s
 => [linux/amd64 1/1] FROM docker.io/library/alpine@sha256:124c7d2707904e...    0.0s
 => => resolve docker.io/library/alpine@sha256:124c7d2707904eea7431fffe91...    0.0s
 => [linux/arm64 1/1] FROM docker.io/library/alpine@sha256:124c7d2707904e...    0.0s
 => => resolve docker.io/library/alpine@sha256:124c7d2707904eea7431fffe91...    0.0s
 => exporting to image                                                          0.0s
 => => exporting layers                                                         0.0s
 ...
 => => naming to docker.io/user/containerd-multiplatform:latest                 0.0s
 => => unpacking to docker.io/user/containerd-multiplatform:latest              0.0s
$ docker images
REPOSITORY                        TAG       IMAGE ID       CREATED          SIZE
user/containerd-multiplatform     latest    7401bb14c229   14 seconds ago   3.38MB
user/containerd-multiplatform     latest    7401bb14c229   14 seconds ago   3.26MB
```



You can push the multi-platform image to Docker Hub.

```console
$ docker push user/containerd-multiplatform
Using default tag: latest
699c4e744ab4: Pushed 
878d877e4f70: Pushed 
f56be85fc22e: Pushed 
a579f49700dc: Pushed 
c41833b44d91: Pushed 
ee79e74f9211: Pushed 
d28bdb47b683: Pushed
```

Inspecting the tag on Docker Hub shows that the image is available for multiple platforms.

![Multiplatform image tag on Docker Hub](images/containerd_multiplatform.png)

## Known issues

### Docker Desktop 4.13.0 release

- Listing images with `docker images` returns the error
  `content digest not found` on ARM machines after running or pulling an image
  with the `--platform` parameter.

### Docker Desktop 4.12.0 release

- The containerd image store feature requires Buildx version 0.9.0 or newer.

  - On Docker Desktop for Linux (DD4L), validate if your locally installed
    version meets this requirement.

    > **Note**
    >
    > If you're using an older version, the Docker daemon reports the following
    > error:
    > `Multiple platforms feature is currently not supported for docker driver.`
    > `Please switch to a different driver`.
    >
    > Install a newer version of Buildx following the instructions on
    > [how to manually download Buildx](../../build/architecture.md#install-buildx).

- In Docker Desktop 4.12.0, the containerd image store feature is incompatible
  with the Kubernetes cluster support. Turn off the containerd image store
  feature if you are using the Kubernetes from Docker Desktop.
- Local registry mirror configuration isn't implemented yet with the containerd
  image store. The `registry-mirrors` and `insecure-registries` aren't taken
  into account by the Docker daemon.
- The `reference` filter isn't implemented yet and will return the error
  `invalid filter 'reference'` when listing images.
- Pulling an image may fail with the error
  `pull access denied, repository does not exist or may require authorization: server message: insufficient_scope: authorization failed`,
  in the situation where the image does not contain a manifest list. To
  workaround this issue run the `docker login` command and pull the image again.

## Feedback

Thanks for trying the new features available with `containerd`. Give feedback or
report any bugs you may find through the issues tracker on the
[feedback form](https://dockr.ly/3PODIhD).
