---
description: Mount directory from machine
keywords: machine, mount, subcommand
title: docker-machine mount
hide_from_sitemap: true
---

Mount directories from a machine to your local host, using `sshfs`.

The notation is `machinename:/path/to/dir` for the argument; you can also supply an alternative mount point (default is the same dir path).

## Example

Consider the following example:

```bash
$ mkdir foo
$ docker-machine ssh dev mkdir foo
$ docker-machine mount dev:/home/docker/foo foo
$ touch foo/bar
$ docker-machine ssh dev ls foo
bar
```


Now you can use the directory on the machine, for mounting into containers.
Any changes done in the local directory, is reflected in the machine too.

```bash
$ eval $(docker-machine env dev)
$ docker run -v /home/docker/foo:/tmp/foo busybox ls /tmp/foo
bar
$ touch foo/baz
$ docker run -v /home/docker/foo:/tmp/foo busybox ls /tmp/foo
bar
baz
```

The files are actually being transferred using `sftp` (over an ssh connection),
so this program ("sftp") needs to be present on the machine - but it usually is.


To unmount the directory again, you can use the same options but the  `-u` flag.
You can also call `fuserunmount` (or `fusermount -u`) commands directly.

```bash
$ docker-machine mount -u dev:/home/docker/foo foo
$ rmdir foo
```
**Files are actually being stored on the machine, *not* on the host.**
So make sure to make a copy of any files you want to keep, before removing it!
