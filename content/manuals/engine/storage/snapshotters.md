---
title: containerd snapshotters
description: Learn how to use containerd snapshotters in Docker Engine.
keywords: containerd, snapshotters, storage, overlayfs, fuse, docker, engine, filesystem, performance
---

containerd snapshotters are components responsible for managing the filesystem
layers of containers. They provide the mechanism for storing, mounting, and
manipulating container filesystems using different backends such as overlayfs or
fuse.

By abstracting the storage implementation, snapshotters allow Docker Engine to
efficiently manage container images and their writable layers, enabling features
like fast image pulls, efficient storage usage, and support for advanced
filesystems.

For more information, see
the [containerd repository](https://github.com/containerd/containerd/tree/main/docs/snapshotters).

| Snapshotter           | Description                                                                            |
|-----------------------|----------------------------------------------------------------------------------------|
| `overlayfs` (default) | OverlayFS. The containerd implementation of the Docker/Moby `overlay2` storage driver. |
| `native`              | Native file copying driver. Akin to Docker/Moby's "vfs" driver.                        |


