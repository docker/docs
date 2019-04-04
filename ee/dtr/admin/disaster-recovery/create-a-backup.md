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

To perform a backup of a DTR node, run the [docker/dtr backup](/reference/dtr/2.6/cli/backup/) command. This
command backs up the following data:

| Data                               | Backed up | Description                                                    |
|:-----------------------------------|:----------|:---------------------------------------------------------------|
| Configurations                     | yes       | DTR settings                                                   |
| Repository metadata                | yes       | Metadata such as image architecture and size                      |
| Access control to repos and images | yes       | Data about who has access to which images                      |
| Notary data                        | yes       | Signatures and digests for images that are signed              |
| Scan results                       | yes       | Information about vulnerabilities in your images               |
| Certificates and keys              | yes       | TLS certificates and keys used by DTR                          |
| Image content                      | no        | Needs to be backed up separately, depends on DTR configuration |
| Users, orgs, teams                 | no        | Create a UCP backup to back up this data                        |
| Vulnerability database             | no        | Can be redownloaded after a restore                           |


## Back up DTR data

To create a backup of DTR you need to:

1. Back up image content
2. Back up DTR metadata

You should always create backups from the same DTR replica, to ensure a smoother
restore. If you have not previously performed a backup, the web interface displays a warning for you to do so:

![](/ee/dtr/images/backup-warning.png)

### Back up image content

Since you can configure the storage backend that DTR uses to store images,
the way you backup images depends on the storage backend you're using.

If you've configured DTR to store images on the local file system or NFS mount,
you can backup the images by using ssh to log into a node where DTR is running,
and creating a tar archive of the [dtr-registry volume](../../architecture.md):

Local images:

{% raw %}
```none
sudo tar -cf dtr-image-backup-$(date +%Y%m%d-%H_%M_%S).tar \
  /var/lib/docker/volumes/dtr-registry-$(docker ps --filter name=dtr-rethinkdb \
  --format "{{ .Names }}" | sed 's/dtr-rethinkdb-//')
```
{% endraw %}

NFS mount images:

{% raw %}
```none
sudo tar -cf dtr-image-backup-$(date +%Y%m%d-%H_%M_%S).tar \
  /var/lib/docker/volumes/dtr-registry-nfs-$(docker ps --filter name=dtr-rethinkdb \
  --format "{{ .Names }}" | sed 's/dtr-rethinkdb-//')
```
{% endraw %}

If you're using a different storage backend, follow the best practices
recommended for that system.


### Back up DTR metadata

To create a DTR backup, load your UCP client bundle, and run the following
command. For your convenience, this command automatically populates your DTR version and replica ID:

```none
DTR_VERSION=$(docker container inspect $(docker container ps -f name=dtr-registry -q) | \
  grep -m1 -Po '(?<=DTR_VERSION=)\d.\d.\d'); \
REPLICA_ID=$(docker ps --filter name=dtr-rethinkdb --format "{{ .Names }}" | head -1 | \
  sed 's|.*/||' | sed 's/dtr-rethinkdb-//'); \
read -p 'ucp-url (The UCP URL including domain and port): ' UCP_URL; \
read -p 'ucp-username (The UCP administrator username): ' UCP_ADMIN; \
read -sp 'ucp password: ' UCP_PASSWORD; \
docker run --log-driver none -i --rm \
  --env UCP_PASSWORD=$UCP_PASSWORD \
  docker/dtr:$DTR_VERSION backup \
  --ucp-username $UCP_ADMIN \
  --ucp-url $UCP_URL \
  --ucp-ca "$(curl https://${UCP_URL}/ca)" \
  --existing-replica-id $REPLICA_ID > dtr-metadata-${DTR_VERSION}-backup-$(date +%Y%m%d-%H_%M_%S).tar
```

This command automatically completes the following tasks:

1. The correct DTR version is automatically set for the backup command using 
the running DTR version.
2. The Replica ID is set automatically for the backup. If you'd prefer to back-up 
a specific replica, the ID can be set manually by modifying the value of the 
`--existing-replica-id` flag. 
3. The UCP password is collected without being saved to disk or printed to the screen.
4. The UCP CA certificate is automatically retrieved and verified (best practice). If
verification is not desired, replace the `--ucp-ca` flag with 
`--ucp-insecure-tls` (not recommended).
5. The backup filename includes the backed-up DTR version and timestamp of the backup.

You can learn more about the supported flags in
the [reference documentation](/reference/dtr/2.5/cli/backup.md).

By default the backup command doesn't pause the DTR replica being backed up to 
prevent interruptions of user access to DTR. Since the replica
is not stopped changes that happen while the backup is taking place may not be saved.

You can use the `--offline-backup` option to stop the DTR replica while taking
the backup. If you do this, remove the replica from the load balancing pool to avoid
user interruption.

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
