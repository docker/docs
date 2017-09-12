---
description: Learn how to select the proper storage driver for your container.
keywords: container, storage, driver, AUFS, btfs, devicemapper,zvfs
title: Select a storage driver
---

Ideally, very little data is written to a container's writable layer, and you
use Docker volumes to write data. However, some workloads require you to be able
to write to the container's writable layer. This is where storage drivers come
in.

Docker supports several different storage drivers, using a pluggable
architecture. The storage driver controls how images and containers are stored
and managed on your Docker host.

After you have read the [storage driver overview](imagesandcontainers.md), the
next step is to choose the best storage driver for your workloads. In making
this decision, there are three high-level factors to consider:

- If multiple storage drivers are supported in your kernel, Docker has a
  prioritized list of which storage driver to use if no storage driver is
  explicitly configured, assuming that the prerequisites for that storage driver
  are met:

  - If `aufs` is available, default to it, because it is the oldest storage
    driver. However, it is not universally available.

  - If possible, the storage driver with the least amount of configuration is
    used, such as `btrfs` or `zfs`. Each of these relies on the backing
    filesystem being configured correctly.

  - Otherwise, try to use the storage driver with the best overall performance
    and stability in the most usual scenarios.

    - `overlay2` is preferred, followed by `overlay`. Neither of these requires
      extra configuration.

    - `devicemapper` is next, but requires `direct-lvm` for production
      environments, because `loopback-lvm`, while zero-configuration, has very
      poor performance.

  The selection order is defined in Docker's source code. You can see the order
  for Docker 17.03 by looking at
  [the source code](https://github.com/moby/moby/blob/v17.03.1-ce/daemon/graphdriver/driver_linux.go#L54-L63).
  For a different Docker version, change the URL to that version.
  {: id="storage-driver-order" }

- Your choice may be limited by your Docker edition, operating system, and
  distribution. For instance, `aufs` is only supported on Ubuntu and Debian,
  while `btrfs` is only supported on SLES, which is only supported with Docker
  EE. See
  [Support storage drivers per Linux distribution](#supported-storage-drivers-per-linux-distribution).

- Some storage drivers require you to use a specific format for the backing
  filesystem. If you have external requirements to use a specific backing
  filesystem, this may limit your choices. See
  [Supported backing filesystems](#supported-backing-filesystems).

- After you have narrowed down which storage drivers you can choose from, your
  choice will be determined by the characteristics of your workload and the
  level of stability you need. See [Other considerations](#other-considerations)
  for help making the final decision.

## Supported storage drivers per Linux distribution

At a high level, the storage drivers you can use is partially determined by
the Docker edition you use.

In addition, Docker does not recommend any configuration that requires you to
disable security features of your operating system, such as the need to disable
`selinux` if you use the `overlay` or `overlay2` driver on CentOS.

### Docker EE and CS-Engine

For Docker EE and CS-Engine, the definitive resource for which storage drivers
are supported is the
[Product compatibility matrix](https://success.docker.com/Policies/Compatibility_Matrix).
In order to get commercial support from Docker, you must use a supported
configuration.

### Docker CE

For Docker CE, only some configurations are tested, and your operating system's
kernel may not support every storage driver. In general, the following
configurations work on recent versions of the Linux distribution:

| Linux distribution  | Recommended storage drivers                                                                           |
|:--------------------|:------------------------------------------------------------------------------------------------------|
| Docker CE on Ubuntu | `aufs`, `devicemapper`, `overlay2` (Ubuntu 14.04.4 or later, 16.04 or later), `overlay`, `zfs`, `vfs` |
| Docker CE on Debian | `aufs`, `devicemapper`, `overlay2` (Debian Stretch), `overlay`, `vfs`                                 |
| Docker CE on CentOS | `devicemapper`, `vfs`                                                                                 |
| Docker CE on Fedora | `devicemapper`, `overlay2` (Fedora 26 or later, experimental), `overlay` (experimental), `vfs`        |

When in doubt, the best all-around configuration is to use a modern Linux
distribution with a kernel that supports the `overlay2` storage driver, and to
use Docker volumes for write-heavy workloads instead of relying on writing data
to the container's writable layer.

The `vfs` storage driver is usually not the best choice. Before using the `vfs`
storage driver, be sure to read about
[its performance and storage characteristics and limitations](vfs-driver.md).

> **Expectations for non-recommended storage drivers**: Commercial support is
> not available for Docker CE, and you can technically use any storage driver
> that is available for your platform. For instance, you can use `btrfs` with
> Docker CE, even though it is not recommended on any platform for Docker CE,
> and you do so at your own risk.
>
> The recommendations in the table above are based on automated regression
> testing and the configurations that are known to work for a large number of
> users. If you use a recommended configuration and find a reproducible issue,
> it is likely to be fixed very quickly. If the driver that you want to use is
> not recommended according to this table, you can run it at your own risk. You
> can and should still report any issues you run into. However, such issues will
> have a lower priority than issues encountered when using a recommended
> configuration.

### Docker for Mac and Docker for Windows

Docker for Mac and Docker for Windows are intended for development, rather
than production. Modifying the storage driver on these platforms is not
supported.

## Supported backing filesystems

With regard to Docker, the backing filesystem is the filesystem where
`/var/lib/docker/` is located. Some storage drivers only work with specific
backing filesystems.

| Storage driver        | Supported backing filesystems |
|:----------------------|:------------------------------|
| `overlay`, `overlay2` | `ext4`, `xfs`                 |
| `aufs`                | `ext4`, `xfs`                 |
| `devicemapper`        | `direct-lvm`                  |
| `btrfs`               | `btrfs`                       |
| `zfs`                 | `zfs`                         |


## Other considerations

### Suitability for your workload

Among other things, each storage driver has its own performance characteristics
that make it more or less suitable for different workloads. Consider the
following generalizations:

- `aufs`, `overlay`, and `overlay2` all operate at the file level rather than
  the block level. This uses memory more efficiently, but the container's
  writable layer may grow quite large in write-heavy workloads.
- Block-level storage drivers such as `devicemapper`, `btrfs`, and `zfs` perform
  better for write-heavy workloads (though not as well as Docker volumes).
- For lots of small writes or containers with many layers or deep filesystems,
  `overlay` may perform better than `overlay2`.
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
and are still under active development. In general, `aufs`, `overlay`, and
`devicemapper` are the choices with the highest stability.

### Experience and expertise

Choose a storage driver that your organization is comfortable maintaining. For
example, if you use RHEL or one of its downstream forks, you may already have
experience with LVM and Device Mapper. If so, the `devicemapper` driver might
be the best choice.

### Test with your own workloads

You can test Docker's performance when running your own workloads on different
storage drivers. Make sure to use equivalent hardware and workloads to match
production conditions, so you can see which storage driver offers the best
overall performance.

## Check and set your current storage driver

The detailed documentation for each individual storage driver details all of the
set-up steps to use a given storage driver. This is a very high-level summary of
how to change the storage driver.

> **Important**: Some storage driver types, such as `devicemapper`, `btrfs`, and
> `zfs`, require additional set-up at the operating system level before you can
> use them with Docker.

To see what storage driver Docker is currently using, use `docker info` and look
for the `Storage Driver` line:

```bash
$ docker info

Containers: 0
Images: 0
Storage Driver: overlay
 Backing Filesystem: extfs
<output truncated>
```



To set the storage driver, you can use the `--storage-driver` flag when starting
the Docker daemon manually, or (recommended) set the option in the `daemon.json`
file, which is located in `/etc/docker/` on Linux and
`C:\ProgramData\docker\config\` on Windows Server.

> **Note**: Using the `--storage-driver` flag when running `dockerd` manually or
> using an init script is not recommended. We recommend setting the option in
> the `daemon.json` file instead, because this mechanism is cross-platform,
> and will not create configuration conflicts with the default init scripts for
> your operating system.

Changing the storage driver on Docker for Mac or Docker for Windows is not
supported.

If the `daemon.json` file does not exist, create it. Assuming there are no other
settings in the file, it should have the following contents:

```json
{
  "storage-driver": "devicemapper"
}
```

You can specify any valid storage driver in place of `devicemapper`.

Restart Docker for the changes to take effect. After restarting, run
`docker info` again to verify that the new storage driver is being used.

## Related information

* [About images, containers, and storage drivers](imagesandcontainers.md)
* [`aufs` storage driver in practice](aufs-driver.md)
* [`devicemapper` storage driver in practice](device-mapper-driver.md)
* [`overlay` and `overlay2` storage drivers in practice](overlayfs-driver.md)
* [`btrfs` storage driver in practice](btrfs-driver.md)
* [`zfs` storage driver in practice](zfs-driver.md)
