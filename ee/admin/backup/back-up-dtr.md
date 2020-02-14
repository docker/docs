---
title: Back up DTR
description: Learn how to create a DTR backup
keywords: enterprise, backup, dtr, disaster recovery
redirect_from:
 - /ee/dtr/admin/disaster-recovery/create-a-backup/
toc_max: 3
toc_min: 1
---

>{% include enterprise_label_shortform.md %}

Backups do not cause downtime for DTR.

## DTR backup contents

All metadata and authZ information for a given DTR cluster is backed up.

| Data                               | Backed up | Description                                                    |
|:-----------------------------------|:----------|:---------------------------------------------------------------|
| Configurations                     | yes       | DTR settings and cluster configurations                                                  |
| Repository metadata                | yes       | Metadata such as image architecture, repositories, images deployed, and size                      |
| Access control to repos and images | yes       | Data about who has access to which images and repositories                      |
| Notary data                        | yes       | Signatures and digests for images that are signed              |
| Scan results                       | yes       | Information about vulnerabilities in your images               |
| Certificates and keys              | yes       | Certificates, public keys, and private keys that are used for mutual TLS communication                          |
| Image content                      | no        | The images you push to DTR. This can be stored on the file system of the node running DTR, or other storage system, depending on the configuration. Needs to be backed up separately, depends on DTR configuration |
| Users, orgs, teams                 | no        | Create a UCP backup to back up this data                        |
| Vulnerability database             | no        | Can be redownloaded after a restore                           |

This data is persisted on the host running DTR, using named volumes.
[Learn more about DTR named volumes](/ee/dtr/architecture/).

## Perform DTR backup

You should always create backups from the same DTR replica, to ensure a smoother
restore. If you have not previously performed a backup, the web interface displays a warning:

![](/ee/dtr/images/backup-warning.png)

To create a DTR backup, perform the following steps:

1. Run [DTR Backup command](#run-the-dtr-backup-command-cli)
2. [Back up DTR image content](#back-up-image-content)
3. [Back up DTR metadata](#back-up-dtr-metadata)
4. [Verify your backup](#verify-your-backup)


### Run the DTR backup command (CLI)

#### Find your replica ID

Since you need your DTR replica ID during a backup, the following covers a few ways for you to determine your replica ID:

##### UCP web interface

You can find the list of replicas by navigating to **Shared Resources > Stacks** or **Swarm > Volumes** (when using [swarm mode](/engine/swarm/)) on the UCP web interface. 

##### UCP client bundle

From a terminal [using a UCP client bundle]((/ee/ucp/user-access/cli/), run:

{% raw %}
```bash
$ docker ps --format "{{.Names}}" | grep dtr

# The list of DTR containers with <node>/<component>-<replicaID>, e.g.
# node-1/dtr-api-a1640e1c15b6
```
{% endraw %}


##### SSH access

Another way to determine the replica ID is to SSH into a DTR node and run the following:

{% raw %}
```bash
$ REPLICA_ID=$(docker inspect -f '{{.Name}}' $(docker ps -q -f name=dtr-rethink) | cut -f 3 -d '-')
&& echo $REPLICA_ID
```
{% endraw %}

### Back up image content

Since you can configure the storage backend that DTR uses to store images,
the way you back up images depends on the storage backend you're using.

If you've configured DTR to store images on the local file system or NFS mount,
you can backup the images by using SSH to log in to a DTR node,
and creating a tar archive of the [dtr-registry volume](/ee/dtr/architecture/):

{% raw %}
```none
$ sudo tar -cf {{ image_backup_file }} \
-C /var/lib/docker/volumes/dtr-registry-<replica-id>
```
{% endraw %}

If you're using a different storage backend, follow the best practices
recommended for that system.


### Back up DTR metadata

To create a DTR backup, load your UCP client bundle, and run the following
command, replacing the placeholders with real values:

```bash
$ read -sp 'ucp password: ' UCP_PASSWORD;
```

This prompts you for the UCP password. Next, run the following to back up your
DTR metadata and save the result into a tar archive. You can learn more about
the supported flags in the [reference
documentation](/reference/dtr/2.7/cli/backup/).

```bash
$ docker container run \
  --rm \
  --interactive \
  --log-driver none \
  --env UCP_PASSWORD=$UCP_PASSWORD \
  {{ page.dtr_org }}/{{ page.dtr_repo }}:{{ page.dtr_version }} backup \
  --ucp-url <ucp-url> \
  --ucp-insecure-tls \
  --ucp-username <ucp-username> \
  --existing-replica-id <replica-id> > dtr-backup-v{{ page.dtr_version }}.tar.gz
```

Where:

* `<ucp-url>` is the url you use to access UCP.
* `<ucp-username>` is the username of a UCP administrator.
* `<replica-id>` is the id of the DTR replica to backup.

By default the backup command doesn't stop the DTR replica being backed up.
This means you can take frequent backups without affecting your users.

You can use the `--offline-backup` option to stop the DTR replica while taking
the backup. If you do this, remove the replica from the load balancing pool.

Also, the backup contains sensitive information
like private keys, so you can encrypt the backup by running:

```none
gpg --symmetric dtr-backup-v{{ page.dtr_version }}.tar.gz
```

This prompts you for a password to encrypt the backup, copies the backup file
and encrypts it.

## Verify your backup

To validate that the backup was correctly performed, you can print the contents
of the tar file created. The backup of the images should look like:

```none
tar -tf dtr-backup-v{{ page.dtr_version }}.tar.gz

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

### Where to go next
- [Restoring Docker Enterprise](/ee/admin/restore/)

