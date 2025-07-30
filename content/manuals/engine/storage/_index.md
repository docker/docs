---
description: Overview of persisting data in containers
title: Storage
weight: 20
keywords: storage, persistence, data persistence, volumes, mounts, bind mounts, tmpfs
aliases:
  - /engine/admin/volumes/
  - /storage/
---

By default all files created inside a container are stored on a writable
container layer that sits on top of the read-only, immutable image layers.

Data written to the container layer doesn't persist when the container is
destroyed. This means that it can be difficult to get the data out of the
container if another process needs it.

The writable layer is unique per container. You can't easily extract the data
from the writeable layer to the host, or to another container.

> [!IMPORTANT]
> Starting with Docker Engine v29, Docker uses
> [containerd](./containerd.md) for managing container storage and images.

## Image storage

Docker offers two implementation of image storage:

- Default/current implementation: [Containerd](./containerd.md)
- Legacy/deprecated implementation: [Storage drivers](./drivers/_index.md)

## Storage mount options

Docker supports the following types of storage mounts for storing data outside
of the writable layer of the container:

- [Volume mounts](#volume-mounts)
- [Bind mounts](#bind-mounts)
- [tmpfs mounts](#tmpfs-mounts)
- [Named pipes](#named-pipes)

No matter which type of mount you choose to use, the data looks the same from
within the container. It is exposed as either a directory or an individual file
in the container's filesystem.

### Volume mounts

Volumes are persistent storage mechanisms managed by the Docker daemon. They
retain data even after the containers using them are removed. Volume data is
stored on the filesystem on the host, but in order to interact with the data in
the volume, you must mount the volume to a container. Directly accessing or
interacting with the volume data is unsupported, undefined behavior, and may
result in the volume or its data breaking in unexpected ways.

Volumes are ideal for performance-critical data processing and long-term
storage needs. Since the storage location is managed on the daemon host,
volumes provide the same raw file performance as accessing the host filesystem
directly.

### Bind mounts

Bind mounts create a direct link between a host system path and a container,
allowing access to files or directories stored anywhere on the host. Since they
aren't isolated by Docker, both non-Docker processes on the host and container
processes can modify the mounted files simultaneously.

Use bind mounts when you need to be able to access files from both the
container and the host.

### tmpfs mounts

A tmpfs mount stores files directly in the host machine's memory, ensuring the
data is not written to disk. This storage is ephemeral: the data is lost when
the container is stopped or restarted, or when the host is rebooted. tmpfs
mounts do not persist data either on the Docker host or within the container's
filesystem.

These mounts are suitable for scenarios requiring temporary, in-memory storage,
such as caching intermediate data, handling sensitive information like
credentials, or reducing disk I/O. Use tmpfs mounts only when the data does not
need to persist beyond the current container session.

### Named pipes

[Named pipes](https://docs.microsoft.com/en-us/windows/desktop/ipc/named-pipes)
can be used for communication between the Docker host and a container. Common
use case is to run a third-party tool inside of a container and connect to the
Docker Engine API using a named pipe.

## Next steps

- Learn more about [volumes](./volumes.md).
- Learn more about [bind mounts](./bind-mounts.md).
- Learn more about [tmpfs mounts](./tmpfs.md).
- Learn more about [storage drivers](/engine/storage/drivers/), which
  are not related to bind mounts or volumes, but allow you to store data in a
  container's writable layer.
