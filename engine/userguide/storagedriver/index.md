---
description: Learn how to select the proper storage driver for your container.
keywords: container, storage, driver, aufs, btrfs, devicemapper, zfs, overlay, overlay2
title: Docker storage drivers
---

Ideally, very little data is written to a container's writable layer, and you
use Docker volumes to write data. However, some workloads require you to be able
to write to the container's writable layer. This is where storage drivers come
in.

Docker uses a series of different storage drivers to manage the filesystems
within images and running containers. These storage drivers are different from
[Docker volumes](/engine/tutorials/dockervolumes.md), which manage storage
which can be shared among multiple containers.

Docker relies on driver technology to manage the storage and interactions
associated with images and the containers that run them. This section contains
the following pages:

* [About images, containers, and storage drivers](imagesandcontainers.md)
* [Select a storage driver](selectadriver.md)
* [AUFS storage driver in practice](aufs-driver.md)
* [Btrfs storage driver in practice](btrfs-driver.md)
* [Device Mapper storage driver in practice](device-mapper-driver.md)
* [OverlayFS in practice](overlayfs-driver.md)
* [ZFS storage in practice](zfs-driver.md)

If you are new to Docker containers make sure you read
[about images, containers, and storage drivers](imagesandcontainers.md) first.
It explains key concepts and technologies that can help you when working with
storage drivers.

### Acknowledgment

The Docker storage driver material was created in large part by our guest author
Nigel Poulton with a bit of help from Docker's own Jérôme Petazzoni. In his
spare time Nigel creates
[IT training videos](http://www.pluralsight.com/author/nigel-poulton) and co-hosts
the weekly
[In Tech We Trust podcast](http://intechwetrustpodcast.com/). Follow him on
[Twitter](https://twitter.com/nigelpoulton).

