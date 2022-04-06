---
title: Docker storage drivers
description: Learn how to select the proper storage driver for your container.
keywords: container, storage, driver, aufs, btrfs, devicemapper, zfs, overlay, overlay2
redirect_from:
- /engine/userguide/storagedriver/
- /engine/userguide/storagedriver/selectadriver/
- /storage/storagedriver/selectadriver/
---

Ideally, very little data is written to a container's writable layer, and you
use Docker volumes to write data. However, some workloads require you to be able
to write to the container's writable layer. This is where storage drivers come
in.

Docker supports several storage drivers, using a pluggable architecture. The
storage driver controls how images and containers are stored and managed on your
Docker host. After you have read the [storage driver overview](index.md), the
next step is to choose the best storage driver for your workloads. Use the storage
driver with the best overall performance and stability in the most usual scenarios.
  
The Docker Engine provides the following storage drivers on Linux:

| Driver             | Description                                                                                                                                                                                                                                                                                                                                                                                                     |
|--------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `overlay2`         | `overlay2` is the preferred storage driver for all currently supported Linux distributions, and requires no extra configuration.                                                                                                                                                                                                                                                                                |
| `fuse-overlayfs`   | `fuse-overlayfs`is preferred only for running Rootless Docker on a host that does not provide support for rootless `overlay2`. On Ubuntu and Debian 10, the `fuse-overlayfs` driver does not need to be used, and `overlay2` works even in rootless mode. Refer to the [rootless mode documentation](../../engine/security/rootless.md) for details.                                                            |
| `btrfs` and `zfs`  | The `btrfs` and `zfs` storage drivers allow for advanced options, such as creating "snapshots", but require more maintenance and setup. Each of these relies on the backing filesystem being configured correctly.                                                                                                                                                                                              |
| `vfs`              | The `vfs` storage driver is intended for testing purposes, and for situations where no copy-on-write filesystem can be used. Performance of this storage driver is poor, and is not generally recommended for production use.                                                                                                                                                                                   |
| `aufs`             | The `aufs` storage driver Was the preferred storage driver for Docker 18.06 and older, when running on Ubuntu 14.04 on kernel 3.13 which had no support for `overlay2`. However, current versions of Ubuntu and Debian now have support for `overlay2`, which is now the recommended driver.                                                                                                                    |
| `devicemapper`     | The `devicemapper` storage driver requires `direct-lvm` for production environments, because `loopback-lvm`, while zero-configuration, has very poor performance. `devicemapper` was the recommended storage driver for CentOS and RHEL, as their kernel version did not support `overlay2`. However, current versions of CentOS and RHEL now have support for `overlay2`, which is now the recommended driver. |
| `overlay`          | The legacy `overlay` driver was used for kernels that did not support the "multiple-lowerdir" feature required for `overlay2` All currently supported Linux distributions now provide support for this, and it is therefore deprecated.                                                                                                                                                                         |

The Docker Engine has a prioritized list of which storage driver to use if no
storage driver is explicitly configured, assuming that the storage driver meets
the prerequisites, and automatically selects a compatible storage driver. You
can see the order in the [source code for Docker Engine {{ site.docker_ce_version }}](https://github.com/moby/moby/blob/{{ site.docker_ce_version }}/daemon/graphdriver/driver_linux.go#L52-L53).
{: id="storage-driver-order" }

Some storage drivers require you to use a specific format for the backing filesystem.
If you have external  requirements to use a specific backing filesystem, this may
limit your choices. See [Supported backing filesystems](#supported-backing-filesystems).

After you have narrowed down which storage drivers you can choose from, your choice
is determined by the  characteristics of your workload and the level of stability
you need. See [Other considerations](#other-considerations) for help in making
the final decision.


## Supported storage drivers per Linux distribution

> Docker Desktop, and Docker in Rootless mode
>
> Modifying the storage-driver is not supported on Docker Desktop for Mac and
> Docker Desktop for Windows, and only the default storage driver can be used.
> The comparison table below is also not applicable for Rootless mode. For the
> drivers available in rootless mode, see the [Rootless mode documentation](../../engine/security/rootless.md).

Your operating system and kernel may not support every storage driver. For
instance, `aufs` is only supported on Ubuntu and Debian, and may require extra
packages to be installed, while `btrfs` is only supported if your system uses
`btrfs` as storage.  In general, the following configurations work on recent
versions of the Linux distribution:

| Linux distribution  | Recommended storage drivers   | Alternative drivers                                |
|:--------------------|:------------------------------|:---------------------------------------------------|
| Ubuntu              | `overlay2`                    | `overlay`¹, `devicemapper`², `aufs`³, `zfs`, `vfs` |
| Debian              | `overlay2`                    | `overlay`¹, `devicemapper`², `aufs`³, `vfs`        |
| CentOS              | `overlay2`                    | `overlay`¹, `devicemapper`², `zfs`, `vfs`          |
| Fedora              | `overlay2`                    | `overlay`¹, `devicemapper`², `zfs`, `vfs`          |
| SLES 15             | `overlay2`                    | `overlay`¹, `devicemapper`², `vfs`                 |
| RHEL                | `overlay2`                    | `overlay`¹, `devicemapper`², `vfs`                 |

¹) The `overlay` storage driver is deprecated, and will be removed in a future
release. It is recommended that users of the `overlay` storage driver migrate to `overlay2`.

²) The `devicemapper` storage driver is deprecated, and will be removed in a future
release. It is recommended that users of the `devicemapper` storage driver migrate
to `overlay2`.

³) The `aufs` storage driver is deprecated, and will be removed in a future
release. It is recommended that users of the `aufs` storage driver migrate
to `overlay2`.

When in doubt, the best all-around configuration is to use a modern Linux
distribution with a kernel that supports the `overlay2` storage driver, and to
use Docker volumes for write-heavy workloads instead of relying on writing data
to the container's writable layer.

The `vfs` storage driver is usually not the best choice, and primarily intended
for debugging purposes in situations where no other storage-driver is supported.
Before using the `vfs` storage driver, be sure to read about
[its performance and storage characteristics and limitations](vfs-driver.md).

The recommendations in the table above are known to work for a large number of
users. If you use a recommended configuration and find a reproducible issue,
it is likely to be fixed very quickly. If the driver that you want to use is
not recommended according to this table, you can run it at your own risk. You
can and should still report any issues you run into. However, such issues
have a lower priority than issues encountered when using a recommended
configuration.

Depending on your Linux distribution, other storage-drivers, such as `btrfs` may
be available. These storage drivers can have advantages for specific use-cases,
but may require additional set-up or maintenance, which make them not recommended
for common scenarios. Refer to the documentation for those storage drivers for
details.

## Supported backing filesystems

With regard to Docker, the backing filesystem is the filesystem where
`/var/lib/docker/` is located. Some storage drivers only work with specific
backing filesystems.

| Storage driver        | Supported backing filesystems  |
|:----------------------|:------------------------------|
| `overlay2`, `overlay` | `xfs` with ftype=1, `ext4`    |
| `fuse-overlayfs`      | any filesystem                |
| `aufs`                | `xfs`, `ext4`                 |
| `devicemapper`        | `direct-lvm`                  |
| `btrfs`               | `btrfs`                       |
| `zfs`                 | `zfs`                         |
| `vfs`                 | any filesystem                 |

## Other considerations

### Suitability for your workload

Among other things, each storage driver has its own performance characteristics
that make it more or less suitable for different workloads. Consider the
following generalizations:

- `overlay2`, `aufs`, and `overlay` all operate at the file level rather than
  the block level. This uses memory more efficiently, but the container's
  writable layer may grow quite large in write-heavy workloads.
- Block-level storage drivers such as `devicemapper`, `btrfs`, and `zfs` perform
  better for write-heavy workloads (though not as well as Docker volumes).
- For lots of small writes or containers with many layers or deep filesystems,
  `overlay` may perform better than `overlay2`, but consumes more inodes, which
  can lead to inode exhaustion.
- `btrfs` and `zfs` require a lot of memory.
- `zfs` is a good choice for high-density workloads such as PaaS.

More information about performance, suitability, and best practices is available
in the documentation for each storage driver.

### Shared storage systems and the storage driver

If your enterprise uses SAN, NAS, hardware RAID, or other shared storage
systems, they may provide high availability, increased performance, thin
provisioning, deduplication, and compression. In many cases, Docker can work on
top of these storage systems, but Docker does not closely integrate with them.

Each Docker storage driver is based on a Linux filesystem or volume manager. Be
sure to follow existing best practices for operating your storage driver
(filesystem or volume manager) on top of your shared storage system. For
example, if using the ZFS storage driver on top of a shared storage system, be
sure to follow best practices for operating ZFS filesystems on top of that
specific shared storage system.

### Stability

For some users, stability is more important than performance. Though Docker
considers all of the storage drivers mentioned here to be stable, some are newer
and are still under active development. In general, `overlay2`, `aufs`, and
`devicemapper` are the choices with the highest stability.

### Test with your own workloads

You can test Docker's performance when running your own workloads on different
storage drivers. Make sure to use equivalent hardware and workloads to match
production conditions, so you can see which storage driver offers the best
overall performance.

## Check your current storage driver

The detailed documentation for each individual storage driver details all of the
set-up steps to use a given storage driver.

To see what storage driver Docker is currently using, use `docker info` and look
for the `Storage Driver` line:

```console
$ docker info

Containers: 0
Images: 0
Storage Driver: overlay2
 Backing Filesystem: xfs
<...>
```

To change the storage driver, see the specific instructions for the new storage
driver. Some drivers require additional configuration, including configuration
to physical or logical disks on the Docker host.

> **Important**
> 
> When you change the storage driver, any existing images and containers become
> inaccessible. This is because their layers cannot be used by the new storage
> driver. If you revert your changes, you can access the old images and containers
> again, but any that you pulled or created using the new driver are then
> inaccessible.
{:.important}

## Related information

* [About images, containers, and storage drivers](index.md)
* [`aufs` storage driver in practice](aufs-driver.md)
* [`devicemapper` storage driver in practice](device-mapper-driver.md)
* [`overlay` and `overlay2` storage drivers in practice](overlayfs-driver.md)
* [`btrfs` storage driver in practice](btrfs-driver.md)
* [`zfs` storage driver in practice](zfs-driver.md)
