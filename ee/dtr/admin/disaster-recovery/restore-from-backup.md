---
title: Restore from a backup
description: Learn how to restore a DTR cluster from an existing backup
keywords: dtr, disaster recovery
---

{% assign metadata_backup_file = "dtr-metadata-backup.tar" %}
{% assign image_backup_file = "dtr-image-backup.tar" %}

## Restore DTR data

If your DTR has a majority of unhealthy replicas, the one way to restore it to
a working state is by restoring from an existing backup.

To restore DTR, you need to:

1. Stop any DTR containers that might be running
2. Restore the images from a backup
3. Restore DTR metadata from a backup
4. Re-fetch the vulnerability database

You need to restore DTR on the same UCP cluster where you've created the
backup. If you restore on a different UCP cluster, all DTR resources will be
owned by users that don't exist, so you'll not be able to manage the resources,
even though they're stored in the DTR data store.

When restoring, you need to use the same version of the `docker/dtr` image
that you've used when creating the update. Other versions are not guaranteed
to work.

### Remove DTR containers

Start by removing any DTR container that is still running:

```none
docker run -it --rm \
  {{ page.dtr_org }}/{{ page.dtr_repo }}:{{ page.dtr_version }} destroy \
  --ucp-insecure-tls
```

### Restore images

If you had DTR configured to store images on the local filesystem, you can
extract your backup:

```none
sudo tar -xf {{ image_backup_file }} -C /var/lib/docker/volumes
```

If you're using a different storage backend, follow the best practices
recommended for that system.

### Restore DTR metadata

You can restore the DTR metadata with the `docker/dtr restore` command. This
performs a fresh installation of DTR, and reconfigures it with
the configuration created during a backup.

Load your UCP client bundle, and run the following command, replacing the
placeholders for the real values:

```bash
read -sp 'ucp password: ' UCP_PASSWORD;
```

This prompts you for the UCP password. Next, run the following to restore DTR from your backup. You can learn more about the supported flags in [docker/dtr restore](/reference/dtr/2.6/cli/restore).

```bash
docker run -i --rm \
  --env UCP_PASSWORD=$UCP_PASSWORD \
  {{ page.dtr_org }}/{{ page.dtr_repo }}:{{ page.dtr_version }} restore \
  --ucp-url <ucp-url> \
  --ucp-insecure-tls \
  --ucp-username <ucp-username> \
  --ucp-node <hostname> \
  --replica-id <replica-id> \
  --dtr-external-url <dtr-external-url> < {{ metadata_backup_file }}
```

Where:

* `<ucp-url>` is the url you use to access UCP
* `<ucp-username>` is the username of a UCP administrator
* `<hostname>` is the hostname of the node where you've restored the images
* `<replica-id>` the id of the replica you backed up
* `<dtr-external-url>`the url that clients use to access DTR

#### DTR 2.5 and below

If you're using NFS as a storage backend, also include `--nfs-storage-url` as
part of your restore command, otherwise DTR is restored but starts using a
local volume to persist your Docker images. 

#### DTR 2.5 (with experimental online garbage collection) and DTR 2.6.0-2.6.3

>  When running DTR 2.5 (with experimental online garbage collection) and 2.6.0 to 2.6.3, there is an issue with [reconfiguring and restoring DTR with `--nfs-storage-url`](/ee/dtr/release-notes#version-26) which leads to erased tags. Make sure to [back up your DTR metadata](/ee/dtr/admin/disaster-recovery/create-a-backup/#back-up-dtr-metadata) before you proceed. To work around the `--nfs-storage-url`flag issue, manually create a storage volume on each DTR node. To [restore DTR](/reference/dtr/2.6/cli/restore/) from an existing backup, use `docker/dtr restore` with `--dtr-storage-volume` and the new volume. 
>
> See [Restore to a Local NFS Volume]( https://success.docker.com/article/dtr-26-lost-tags-after-reconfiguring-storage#restoretoalocalnfsvolume) for Docker's recommended recovery strategy.  
{: .warning}

### Re-fetch the vulnerability database

If you're scanning images, you now need to download the vulnerability database.

After you successfully restore DTR, you can join new replicas the same way you
would after a fresh installation. [Learn more](../configure/set-up-vulnerability-scans.md).

## Where to go next

- [docker/dtr restore](/reference/dtr/2.6/cli/restore/)
