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
mkdir /tmp/mydir && sudo mount -t nfs <nfs server>:<directory>
```

## Install DTR with NFS

One way to configure DTR to use an NFS directory is at install time:

```none
docker run -it --rm {{ page.dtr_org }}/{{ page.dtr_repo }}:{{ page.dtr_version }} install \
  --nfs-storage-url <nfs-storage-url> \
  <other options>
```

The NFS storage URL should be in the format `nfs://<nfs server>/<directory>`.

When you join replicas to the DTR cluster, the replicas will pick up that
configuration, so you don't need to specify it again.

### Reconfigure DTR to use NFS

If you're upgrading from a previous version of DTR and are already using
NFS you can continue using the same configurations.

If you want to start using the new DTR built-in support for NFS you can
reconfigure DTR:

```none
docker run -it --rm {{ page.dtr_org }}/{{ page.dtr_repo }}:{{ page.dtr_version }} reconfigure \
  --nfs-storage-url <nfs-storage-url>
```

If you want to reconfigure DTR to stop using NFS storage, leave the option
in blank:

```none
docker run -it --rm {{ page.dtr_org }}/{{ page.dtr_repo }}:{{ page.dtr_version}} reconfigure \
  --nfs-storage-url ""
```

If the IP address of your NFS server changes, even if the DNS address is kept
the same, you should reconfigure DTR to stop using NFS storage, and then
add it back again.

## Where to go next

* [Configure where images are stored](index.md)
