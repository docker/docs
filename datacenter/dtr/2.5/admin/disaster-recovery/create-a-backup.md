---
title: Create a backup
description: Learn how to create a backup of Docker Trusted Registry, for disaster recovery.
keywords: dtr, disaster recovery
---

{% assign metadata_backup_file = "dtr-metadata-backup.tar" %}
{% assign image_backup_file = "dtr-image-backup.tar" %}


## Data managed by DTR

Docker Trusted Registry maintains data about:

| Data                               | Description                                                                                                                                        |
|:-----------------------------------|:---------------------------------------------------------------------------------------------------------------------------------------------------|
| Configurations                     | The DTR cluster configurations                                                                                                                     |
| Repository metadata                | The metadata about the repositories and images deployed                                                                                            |
| Access control to repos and images | Permissions for teams and repositories                                                                                                             |
| Notary data                        | Notary tags and signatures                                                                                                                         |
| Scan results                       | Security scanning results for images                                                                                                               |
| Certificates and keys              | The certificates, public keys, and private keys that are used for mutual TLS communication                                                         |
| Images content                     | The images you push to DTR. This can be stored on the file system of the node running DTR, or other storage system, depending on the configuration |

This data is persisted on the host running DTR, using named volumes.
[Learn more about DTR named volumes](../../architecture.md).

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

If you've configured DTR to store images on the local file system or NFS mount,
you can backup the images by using ssh to log into a node where DTR is running,
and creating a tar archive of the [dtr-registry volume](../../architecture.md):

{% raw %}
```none
sudo tar -cf {{ image_backup_file }} \
$(dirname $(docker volume inspect --format '{{.Mountpoint}}' dtr-registry-<replica-id>))
```
{% endraw %}

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

* `<ucp-url>` is the url you use to access UCP.
* `<ucp-username>` is the username of a UCP administrator.
* `<replica-id>` is the id of the DTR replica to backup.

This prompts you for the UCP password, backups up the DTR metadata and saves the
result into a tar archive. You can learn more about the supported flags in
the [reference documentation](/reference/dtr/2.5/cli/backup.md).

By default the backup command doesn't stop the DTR replica being backed up.
This allows performing backups without affecting your users. Since the replica
is not stopped, it's possible that happen while the backup is taking place, won't
be persisted.

You can use the `--offline-backup` option to stop the DTR replica while taking
the backup. If you do this, remove the replica from the load balancing pool.

Also, the backup contains sensitive information
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
tar -tf {{ metadata_backup_file }}

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
