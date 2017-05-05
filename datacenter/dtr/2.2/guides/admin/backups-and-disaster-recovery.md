---
title: DTR backups and recovery
description: Learn how to back up your Docker Trusted Registry cluster, and to recover your cluster from an existing backup.
keywords: docker, registry, high-availability, backup, recovery
---

DTR needs that a majority (n/2 + 1) of its replicas are healthy at all times
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
command backs up the following data:

| Data                               | Backed up | Description                                                    |
|:-----------------------------------|:----------|:---------------------------------------------------------------|
| Configurations                     | yes       |                                                                |
| Repository metadata                | yes       |                                                                |
| Access control to repos and images | yes       |                                                                |
| Notary data                        | yes       |                                                                |
| Scan results                       | yes       |                                                                |
| Certificates and keys              | yes       |                                                                |
| Image content                      | no        | Needs to be backed up separately, depends on DTR configuration |
| Users, orgs, teams                 | no        | Create a UCP backup to back up this data                        |
| Vulnerability database             | no        | Can be re-downloaded after a restore                           |


## Back up DTR data

To create a backup of DTR you need to:

1. Back up image content
2. Back up DTR metadata

You should always create backups from the same DTR replica, to ensure a smoother
restore.

### Back up image content

Since you can configure the storage backend that DTR uses to store images,
the way you back up images depends on the storage backend you're using.

If you've configured DTR to store images on the local filesystem or NFS mount,
you can back up the images by using ssh to log into a node where DTR is running,
and creating a tar archive of the [dtr-registry volume](../architecture.md):

```none
tar -cf /tmp/backup-images.tar dtr-registry-<replica-id>
```

If you're using a different storage backend, follow the best practices
recommended for that system.


### Back up DTR metadata

To create a DTR backup, load your UCP client bundle, and run the following
command, replacing the placeholders for the real values:

```none
read -sp 'ucp password: ' UCP_PASSWORD; \
docker run -i --rm \
  --env UCP_PASSWORD=$UCP_PASSWORD \
  docker/dtr:<version> backup \
  --ucp-url <ucp-url> \
  --ucp-insecure-tls \
  --ucp-username <ucp-username> \
  --existing-replica-id <replica-id> > /tmp/backup-metadata.tar
```

Where:

* `<version>`, the version of DTR you're running
* `<ucp-url>` is the url you use to access UCP
* `<ucp-username>` is the username of a UCP administrator
* `<replica-id>` is the ID of the DTR replica to back up


This prompts you for the UCP password, backs up the DTR metadata and saves the
result into a tar archive. You can learn more about the supported flags in
the [reference documentation](../../reference/cli/backup.md).

The backup command doesn't stop DTR, so that you can take frequent backups
without affecting your users. Also, the backup contains sensitive information
like private keys, so you can encrypt the backup by running:

```none
gpg --symmetric /tmp/backup-metadata.tar
```

This prompts you for a password to encrypt the backup, copies the backup file
and encrypts it.

### Test your backups

To validate that the backup was correctly performed, you can print the contents
of the tar file created. The backup of the images should look like:

```none
tar -tf /tmp/backup-images.tar

dtr-backup-v2.2.3/
dtr-backup-v2.2.3/rethink/
dtr-backup-v2.2.3/rethink/layers/
```

And the backup of the DTR metadata should look like:

```none
tar -tf /tmp/backup-metadata.tar

# The archive should look like this
dtr-backup-v2.2.1/
dtr-backup-v2.2.1/rethink/
dtr-backup-v2.2.1/rethink/properties/
dtr-backup-v2.2.1/rethink/properties/0
```

If you've encrypted the metadata backup, you can use:

```none
gpg -d /tmp/backup.tar.gpg | tar -t
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
owned by users that don't exist, so you'll not be able to manage the resources,
even though they're stored in the DTR data store.

When restoring, you need to use the same version of the `docker/dtr` image
that you've used when creating the update. Other versions are not guaranteed
to work.

### Stop DTR containers

Start by removing any DTR container that is still running:

```none
docker run -it --rm \
  docker/dtr:<version> destroy \
  --ucp-insecure-tls
```

### Restore images

If you had DTR configured to store images on the local filesystem, you can
extract your backup:

```none
sudo tar -xzf /tmp/image-backup.tar -C /var/lib/docker/volumes
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
  docker/dtr:<version> restore \
  --ucp-url <ucp-url> \
  --ucp-insecure-tls \
  --ucp-username <ucp-username> \
  --ucp-node <hostname> \
  --replica-id <replica-id> \
  --dtr-external-url <dtr-external-url> < /tmp/backup-metadata.tar
```

Where:

* `<version>`, the version of DTR you're running
* `<ucp-url>` is the url you use to access UCP
* `<ucp-username>` is the username of a UCP administrator
* `<hostname>` is the hostname of the node where you've restored the images
* `<replica-id>` the ID of the replica you backed up
* `<dtr-external-url>` the url that clients use to access DTR

### Re-fetch the vulnerability database

If you're scanning images, you now need to download the vulnerability database.

After you successfully restore DTR, you can join new replicas the same way you
would after a fresh installation. [Learn more](configure/set-up-vulnerability-scans.md).

## Where to go next

* [Set up high availability](configure/set-up-high-availability.md)
* [DTR architecture](../architecture.md)
