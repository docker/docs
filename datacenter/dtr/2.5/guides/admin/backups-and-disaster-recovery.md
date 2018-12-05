---
title: DTR backups and recovery
description: Learn how to back up your Docker Trusted Registry cluster, and to recover your cluster from an existing backup.
keywords: registry, high-availability, backup, recovery
---

{% assign image_backup_file = "backup-images.tar" %}
{% assign metadata_backup_file = "backup-metadata.tar" %}

DTR requires that a majority (n/2 + 1) of its replicas are healthy at all times
for it to work. So if a majority of replicas is unhealthy or lost, the only
way to restore DTR to a working state, is by recovering from a backup. This
is why it's important to ensure replicas are healthy and perform frequent
backups.

## Data managed by DTR

Docker Trusted Registry maintains data about:

| Data                               | Description                                                                                                                                       |
|:-----------------------------------|:--------------------------------------------------------------------------------------------------------------------------------------------------|
| Configurations                     | The DTR cluster configurations                                                                                                                    |
| Repository metadata                | The metadata about the repositories and images deployed                                                                                           |
| Access control to repos and images | Permissions for teams and repositories                                                                                                            |
| Notary data                        | Notary tags and signatures                                                                                                                        |
| Scan results                       | Security scanning results for images                                                                                                              |
| Certificates and keys              | The certificates, public keys, and private keys that are used for mutual TLS communication                                                        |
| Images content                     | The images you push to DTR. This can be stored on the filesystem of the node running DTR, or other storage system, depending on the configuration |

This data is persisted on the host running DTR, using named volumes.
[Learn more about DTR named volumes](../architecture.md).

To perform a backup of a DTR node, run the `docker/dtr backup` command. This
command backups up the following data:

| Data                               | Backed up | Description                                                    |
|:-----------------------------------|:----------|:---------------------------------------------------------------|
| Configurations                     | yes       | DTR settings                                                   |
| Repository metadata                | yes       | Metadata like image architecture and size                      |
| Access control to repos and images | yes       | Data about who has access to which images                      |
| Notary data                        | yes       | Signatures and digests for images that are signed              |
| Scan results                       | yes       | Information about vulnerabilities in your images               |
| Certificates and keys              | yes       | TLS certificates and keys used by DTR                          |
| Image content                      | no        | Needs to be backed up separately, depends on DTR configuration |
| Users, orgs, teams                 | no        | Create a UCP backup to backup this data                        |
| Vulnerability database             | no        | Can be re-downloaded after a restore                           |


## Backup DTR data

To create a backup of DTR you need to:

1. Backup image content
2. Backup DTR metadata

You should always create backups from the same DTR replica, to ensure a smoother
restore.

### Backup image content

Since you can configure the storage backend that DTR uses to store images,
the way you backup images depends on the storage backend you're using.

If you've configured DTR to store images on the local filesystem or NFS mount,
you can backup the images by using ssh to log into a node where DTR is running,
and creating a tar archive of the [dtr-registry volume](../architecture.md):

```none
sudo tar -cf {{ image_backup_file }} \
  {% raw %}$(dirname $(docker volume inspect --format '{{.Mountpoint}}' dtr-registry-<replica-id>)){% endraw %}
```

If you're using a different storage backend, follow the best practices
recommended for that system.


### Backup DTR metadata

To create a DTR backup, load your UCP client bundle, and run the following
command, replacing the placeholders for the real values:

```none
read -sp 'ucp password: ' UCP_PASSWORD; \
docker run --log-driver none -i --rm \
  --env UCP_PASSWORD=$UCP_PASSWORD \
  {{ page.dtr_org }}/{{ page.dtr_repo }}:{{ page.dtr_version }} backup \
  --ucp-url <ucp-url> \
  --ucp-insecure-tls \
  --ucp-username <ucp-username> \
  --existing-replica-id <replica-id> > {{ metadata_backup_file }}
```

Where:

* `<ucp-url>` is the url you use to access UCP
* `<ucp-username>` is the username of a UCP administrator
* `<replica-id>` is the id of the DTR replica to backup


This prompts you for the UCP password, backups up the DTR metadata and saves the
result into a tar archive. You can learn more about the supported flags in
the [reference documentation](../../reference/cli/backup.md).

The backup command doesn't stop DTR, so that you can take frequent backups
without affecting your users. Also, the backup contains sensitive information
like private keys, so you can encrypt the backup by running:

```none
gpg --symmetric {{ metadata_backup_file }}
```

This prompts you for a password to encrypt the backup, copies the backup file
and encrypts it.

### Test your backups

To validate that the backup was correctly performed, you can print the contents
of the tar file created. The backup of the images should look like:

```none
tar -tf {{ image_backup_file }}

dtr-backup-v{{ page.dtr_version }}/
dtr-backup-v{{ page.dtr_version }}/rethink/
dtr-backup-v{{ page.dtr_version }}/rethink/layers/
```

And the backup of the DTR metadata should look like:

```none
tar -tf {{ metadata_backup_file }}

# The archive should look like this
dtr-backup-v{{ page.dtr_version }}/
dtr-backup-v{{ page.dtr_version }}/rethink/
dtr-backup-v{{ page.dtr_version }}/rethink/properties/
dtr-backup-v{{ page.dtr_version }}/rethink/properties/0
```

If you've encrypted the metadata backup, you can use:

```none
gpg -d {{ metadata_backup_file }} | tar -t
```

You can also create a backup of a UCP cluster and restore it into a new
cluster. Then restore DTR on that new cluster to confirm that everything is
working as expected.

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
owned by users that don't exist, so you can't manage the resources,
even though they're stored in the DTR data store.

When restoring, you need to use the same version of the `docker/dtr` image
that you've used when creating the update. Other versions are not guaranteed
to work.

### Stop DTR containers

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
sudo tar -xf {{ image_backup_file }} -C /
```

If you're using a different storage backend, follow the best practices
recommended for that system. When restoring the DTR metadata, DTR will be
deployed with the same configurations it had when creating the backup.


### Restore DTR metadata

You can restore the DTR metadata with the `docker/dtr restore` command. This
performs a fresh installation of DTR, and reconfigures it with
the configuration created during a backup.

Load your UCP client bundle, and run the following command, replacing the
placeholders for the real values:

```none
read -sp 'ucp password: ' UCP_PASSWORD; \
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

If you're using NFS as a storage backend, also include `--nfs-storage-url` as
part of your restore command, otherwise DTR is restored but starts using a
local volume to persist your Docker images.

### Re-fetch the vulnerability database

If you're scanning images, you now need to download the vulnerability database.

After you successfully restore DTR, you can join new replicas the same way you
would after a fresh installation. [Learn more](configure/set-up-vulnerability-scans.md).

## Where to go next

* [Set up high availability](configure/set-up-high-availability.md)
* [DTR architecture](../architecture.md)
