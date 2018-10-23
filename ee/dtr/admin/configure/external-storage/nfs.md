---
title: Use NFS
description: Learn how to integrate Docker Trusted Registry with NFS
keywords: registry, dtr, storage, nfs
---

You can configure DTR to store Docker images in an NFS directory.

Before installing or configuring DTR to use an NFS directory, make sure that:

* The NFS server has been correctly configured
* The NFS server has a fixed IP address
* All hosts running DTR have the correct NFS libraries installed


To confirm that the hosts can connect to the NFS server, try to list the
directories exported by your NFS server:

```bash
showmount -e <nfsserver>
```

You should also try to mount one of the exported directories:

```bash
mkdir /tmp/mydir && sudo mount -t nfs <nfs server>:<directory> /tmp/mydir
```

## Install DTR with NFS

One way to configure DTR to use an NFS directory is at install time:

```bash
docker run -it --rm {{ page.dtr_org }}/{{ page.dtr_repo }}:{{ dtr_version }} install \
  --nfs-storage-url <nfs-storage-url> \
  <other options>
```

The NFS storage URL should be in the format `nfs://<nfs server>/<directory>`. With **NFS v4**, you can specify additional options. See [docker/dtr install](../../../../../reference/dtr/2.6/cli/install/) for more details.

When joining replicas to a DTR cluster, the replicas will pick up your storage
configuration, so you will not need to specify it again.

### Reconfigure DTR to use NFS

When upgrading from a previous version of DTR that is already using
NFS, you can continue using the same configurations. If you want to use **NFS v4**, see [docker/dtr reconfigure](../../../../../reference/dtr/2.6/cli/reconfigure/) for more NFS options.


To take advantage of the new DTR built-in support for NFS, you can
reconfigure DTR to use NFS:

```bash
docker run -it --rm {{ page.dtr_org }}/{{ page.dtr_repo }}:{{ dtr_version }} reconfigure \
  --nfs-storage-url <nfs-storage-url>
```

To reconfigure DTR to stop using NFS storage, leave the `--nfs-storage-url` option
blank:

```bash
docker run -it --rm {{ page.dtr_org }}/{{ page.dtr_repo }}:{{ dtr_version}} reconfigure \
  --nfs-storage-url ""
```

If the IP address of your NFS server changes with the DNS address staying the same, you should still 
reconfigure DTR to stop using NFS storage, and then add it back again.

## Where to go next

- [Configure where images are stored](index.md)
