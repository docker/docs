---
title: Use NFS
description: Learn how to integrate Docker Trusted Registry with NFS
keywords: registry, dtr, storage, nfs
---

You can configure DTR to store Docker images in an NFS directory. Starting in DTR 2.6,
changing storage backends involves initializing a new metadatastore instead of reusing an existing volume. 
This helps facilitate [online garbage collection](/ee/dtr/admin/configure/garbage-collection/#under-the-hood).
See [changes to NFS reconfiguration below](/ee/dtr/admin/configure/external-storage/nfs/#reconfigure-dtr-to-use-nfs) if you have previously configured DTR to use NFS.

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
docker run -it --rm {{ page.dtr_org }}/{{ page.dtr_repo }}:{{ page.dtr_version }} install \
  --nfs-storage-url <nfs-storage-url> \
  <other options>
```

Use the format `nfs://<nfs server>/<directory>` for the NFS storage URL. To support **NFS v4**, you can now specify additional options when running [docker/dtr install](/reference/dtr/2.6/cli/install/) with `--nfs-storage-url`.

When joining replicas to a DTR cluster, the replicas will pick up your storage
configuration, so you will not need to specify it again.

### Reconfigure DTR to use NFS

To support **NFS v4**, more NFS options have been added to the CLI. See [New Features for 2.6.0 - CLI](/ee/dtr/release-notes/#260) for updates to [docker/dtr reconfigure](/reference/dtr/2.6/cli/reconfigure/). 

> When running DTR 2.5 (with experimental online garbage collection) and 2.6.0 to 2.6.3, there is an issue with [reconfiguring and restoring DTR with `--nfs-storage-url`](/ee/dtr/release-notes#version-26) which leads to erased tags. Make sure to [back up your DTR metadata](/ee/dtr/admin/disaster-recovery/create-a-backup/#back-up-dtr-metadata) before you proceed. To work around the `--nfs-storage-url` flag issue, manually create a storage volume. If DTR is already installed in your cluster, [reconfigure DTR](/reference/dtr/2.6/cli/reconfigure/) with the `--dtr-storage-volume` flag using your newly-created volume. 
>
> See [Reconfigure Using a Local NFS Volume]( https://success.docker.com/article/dtr-26-lost-tags-after-reconfiguring-storage#reconfigureusingalocalnfsvolume) for Docker's recommended recovery strategy.  
{: .warning}

#### DTR 2.6.4 

In DTR 2.6.4, a new flag, `--storage-migrated`, [has been added to `docker/dtr reconfigure`](/reference/dtr/2.6/cli/reconfigure/) which lets you indicate the migration status of your storage data during a reconfigure. [Upgrade to 2.6.4](/reference/dtr/2.6/cli/upgrade/) and follow [Best practice for data migration in 2.6.4](/ee/dtr/admin/configure/external-storage/storage-backend-migration/#best-practice-for-data-migration) when switching storage backends. The following shows you how to reconfigure DTR using an NFSv4 volume as a storage backend:

```bash
docker run --rm -it \
  docker/dtr:{{ page.dtr_version}} reconfigure \
  --ucp-url <ucp_url> \
  --ucp-username <ucp_username> \
  --nfs-storage-url <dtr-registry-nf>
  --async-nfs
  --storage-migrated
```

To reconfigure DTR to stop using NFS storage, leave the `--nfs-storage-url` option
blank:

```bash
docker run -it --rm {{ page.dtr_org }}/{{ page.dtr_repo }}:{{ page.dtr_version}} reconfigure \
  --nfs-storage-url ""
```

## Where to go next

- [Switch storage backends](storage-backend-migration.md)
- [Create a backup](/ee/dtr/admin/disaster-recovery/create-a-backup/)
- [Restore from a backup](/ee/dtr/admin/disaster-recovery/restore-from-backup/)
- [Configure where images are stored](index.md)
- CLI reference pages
  - [docker/dtr install](/reference/dtr/2.6/cli/install/)
  - [docker/dtr reconfigure](/reference/dtr/2.6/cli/reconfigure/)
  - [docker/dtr restore](/reference/dtr/2.6/cli/restore/)
