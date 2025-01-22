---
description: Using tmpfs mounts
title: tmpfs mounts
weight: 30
keywords: storage, persistence, data persistence, tmpfs
aliases:
  - /engine/admin/volumes/tmpfs/
  - /storage/tmpfs/
---

[Volumes](volumes.md) and [bind mounts](bind-mounts.md) let you share files
between the host machine and container so that you can persist data even after
the container is stopped.

If you're running Docker on Linux, you have a third option: tmpfs mounts.
When you create a container with a tmpfs mount, the container can create
files outside the container's writable layer.

As opposed to volumes and bind mounts, a tmpfs mount is temporary, and only
persisted in the host memory. When the container stops, the tmpfs mount is
removed, and files written there won't be persisted.

tmpfs mounts are best used for cases when you do not want the data to persist
either on the host machine or within the container. This may be for security
reasons or to protect the performance of the container when your application
needs to write a large volume of non-persistent state data.

> [!IMPORTANT]
> tmpfs mounts in Docker map directly to
> [tmpfs](https://en.wikipedia.org/wiki/Tmpfs) in the Linux kernel. As such,
> the temporary data may be written to a swap file, and thereby persisted to
> the filesystem.

## Mounting over existing data

If you create a tmpfs mount into a directory in the container in which files or
directories exist, the pre-existing files are obscured by the mount. This is
similar to if you were to save files into `/mnt` on a Linux host, and then
mounted a USB drive into `/mnt`. The contents of `/mnt` would be obscured by
the contents of the USB drive until the USB drive was unmounted.

With containers, there's no straightforward way of removing a mount to reveal
the obscured files again. Your best option is to recreate the container without
the mount.

## Limitations of tmpfs mounts

- Unlike volumes and bind mounts, you can't share tmpfs mounts between containers.
- This functionality is only available if you're running Docker on Linux.
- Setting permissions on tmpfs may cause them to [reset after container restart](https://github.com/docker/for-linux/issues/138). In some cases [setting the uid/gid](https://github.com/docker/compose/issues/3425#issuecomment-423091370) can serve as a workaround.

## Syntax

To mount a tmpfs with the `docker run` command, you can use either the
`--mount` or `--tmpfs` flag.

```console
$ docker run --mount type=tmpfs,dst=<mount-path>
$ docker run --tmpfs <mount-path>
```

In general, `--mount` is preferred. The main difference is that the `--mount`
flag is more explicit. On the other hand, `--tmpfs` is less verbose and gives
you more flexibility as it lets you set more mount options.

The `--tmpfs` flag cannot be used with swarm services. You must use `--mount`.

### Options for --tmpfs

The `--tmpfs` flag consists of two fields, separated by a colon character
(`:`).

```console
$ docker run --tmpfs <mount-path>[:opts]
```

The first field is the container path to mount into a tmpfs. The second field
is optional and lets you set mount options. Valid mount options for `--tmpfs`
include:

| Option       | Description                                                                                 |
| ------------ | ------------------------------------------------------------------------------------------- |
| `ro`         | Creates a read-only tmpfs mount.                                                            |
| `rw`         | Creates a read-write tmpfs mount (default behavior).                                        |
| `nosuid`     | Prevents `setuid` and `setgid` bits from being honored during execution.                    |
| `suid`       | Allows `setuid` and `setgid` bits to be honored during execution (default behavior).        |
| `nodev`      | Device files can be created but are not functional (access results in an error).            |
| `dev`        | Device files can be created and are fully functional.                                       |
| `exec`       | Allows the execution of executable binaries in the mounted file system.                     |
| `noexec`     | Does not allow the execution of executable binaries in the mounted file system.             |
| `sync`       | All I/O to the file system is done synchronously.                                           |
| `async`      | All I/O to the file system is done asynchronously (default behavior).                       |
| `dirsync`    | Directory updates within the file system are done synchronously.                            |
| `atime`      | Updates file access time each time the file is accessed.                                    |
| `noatime`    | Does not update file access times when the file is accessed.                                |
| `diratime`   | Updates directory access times each time the directory is accessed.                         |
| `nodiratime` | Does not update directory access times when the directory is accessed.                      |
| `size`       | Specifies the size of the tmpfs mount, for example, `size=64m`.                             |
| `mode`       | Specifies the file mode (permissions) for the tmpfs mount (for example, `mode=1777`).       |
| `uid`        | Specifies the user ID for the owner of the tmpfs mount (for example, `uid=1000`).           |
| `gid`        | Specifies the group ID for the owner of the tmpfs mount (for example, `gid=1000`).          |
| `nr_inodes`  | Specifies the maximum number of inodes for the tmpfs mount (for example, `nr_inodes=400k`). |
| `nr_blocks`  | Specifies the maximum number of blocks for the tmpfs mount (for example, `nr_blocks=1024`). |

```console {title="Example"}
$ docker run --tmpfs /data:noexec,size=1024,mode=1777
```

Not all tmpfs mount features available in the Linux mount command are supported
with the `--tmpfs` flag. If you require advanced tmpfs options or features, you
may need to use a privileged container or configure the mount outside of
Docker.

> [!CAUTION]
> Running containers with `--privileged` grants elevated permissions and can
> expose the host system to security risks. Use this option only when
> absolutely necessary and in trusted environments.

```console
$ docker run --privileged -it debian sh
/# mount -t tmpfs -o <options> tmpfs /data
```

### Options for --mount

The `--mount` flag consists of multiple key-value pairs, separated by commas
and each consisting of a `<key>=<value>` tuple. The order of the keys isn't
significant.

```console
$ docker run --mount type=tmpfs,dst=<mount-path>[,<key>=<value>...]
```

Valid options for `--mount type=tmpfs` include:

| Option                         | Description                                                                                                            |
| :----------------------------- | :--------------------------------------------------------------------------------------------------------------------- |
| `destination`, `dst`, `target` | Size of the tmpfs mount in bytes. If unset, the default maximum size of a tmpfs volume is 50% of the host's total RAM. |
| `tmpfs-size`                   | Size of the tmpfs mount in bytes. If unset, the default maximum size of a tmpfs volume is 50% of the host's total RAM. |
| `tmpfs-mode`                   | File mode of the tmpfs in octal. For instance, `700` or `0770`. Defaults to `1777` or world-writable.                  |

```console {title="Example"}
$ docker run --mount type=tmpfs,dst=/app,tmpfs-size=21474836480,tmpfs-mode=1770
```

## Use a tmpfs mount in a container

To use a `tmpfs` mount in a container, use the `--tmpfs` flag, or use the
`--mount` flag with `type=tmpfs` and `destination` options. There is no
`source` for `tmpfs` mounts. The following example creates a `tmpfs` mount at
`/app` in a Nginx container. The first example uses the `--mount` flag and the
second uses the `--tmpfs` flag.

{{< tabs >}}
{{< tab name="`--mount`" >}}

```console
$ docker run -d \
  -it \
  --name tmptest \
  --mount type=tmpfs,destination=/app \
  nginx:latest
```

Verify that the mount is a `tmpfs` mount by looking in the `Mounts` section of
the `docker inspect` output:

```console
$ docker inspect tmptest --format '{{ json .Mounts }}'
[{"Type":"tmpfs","Source":"","Destination":"/app","Mode":"","RW":true,"Propagation":""}]
```

{{< /tab >}}
{{< tab name="`--tmpfs`" >}}

```console
$ docker run -d \
  -it \
  --name tmptest \
  --tmpfs /app \
  nginx:latest
```

Verify that the mount is a `tmpfs` mount by looking in the `Mounts` section of
the `docker inspect` output:

```console
$ docker inspect tmptest --format '{{ json .Mounts }}'
{"/app":""}
```

{{< /tab >}}
{{< /tabs >}}

Stop and remove the container:

```console
$ docker stop tmptest
$ docker rm tmptest
```

## Next steps

- Learn about [volumes](volumes.md)
- Learn about [bind mounts](bind-mounts.md)
- Learn about [storage drivers](/engine/storage/drivers/)
