---
description: Learn how to back up your Docker Trusted Registry cluster, and to recover your cluster from an existing backup.
keywords: docker, registry, high-availability, backup, recovery
title: DTR backups and recovery
---

DTR replicas rely on having a majority available at any given time for writes to
succeed. Therefore if a majority of replicas are permanently lost, the only way
to restore DTR to a working state is to recover from backups. This is why it's
very important to perform periodic backups.

## DTR data persistence

Docker Trusted Registry persists:

* **Configurations**: the cluster configurations are stored on a key-value store
that is replicated through all DTR replicas.
* **Repository metadata**: the information about the repositories and
images deployed. This information is replicated through all DTR replicas.
* **Access control**: permissions for teams and repos.
* **Notary data**: notary tags and signatures.
* **Scan results**: security scanning results for images.
* **Certificates and keys**: the certificates, public keys, and private keys
that are used for mutual TLS communication.

This data persists using named volumes on the host where DTR is running.
[Learn more about DTR named volumes](../architecture.md).

DTR also persists Docker images on the filesystem of the host running DTR, or
on a cloud provider, depending on the way DTR is configured.

## Backup DTR data

To perform a backup of a DTR node, use the `backup` command. This
command creates a backup of DTR:

* Configurations,
* Repository metadata,
* Access control,
* Notary data,
* Scan results,
* Certificates and keys used by DTR.

This data is added to a tar archive, and the result is streamed to stdout. This
is done while DTR is running without shutting down any containers.

Things DTR's backup command doesn't back up:

* The vulnerability database (if using image scanning)
* Image contents
* Users, orgs, teams

There is no way to back up the vulnerability database. You can re-download it
after restoring or re-apply your offline tar update if offline.

The backup command does not create a backup of Docker images. You should
implement a separate backup policy for the Docker images, taking in
consideration whether your DTR installation is configured to store images on the
filesystem or using a cloud provider. During restore, you need to separately
restore the image contents.

The backup command also doesn't create a backup of the users and organizations.
That data is managed by UCP, so when you create a UCP backup you're creating
a backup of the users and organizations. For this reason, when restoring DTR,
you must do it on the same UCP cluster (or one created by restoring from
backups) or else all DTR resources such as repos will be owned by non-existent
users and will not be usable despite technically existing in the database.

When creating a backup, the resulting .tar file contains sensitive information
such as private keys. You should ensure the backups are stored securely.

You can check the
[reference documentation](../../reference/cli/backup.md), for the
backup command to learn about all the available flags.

As an example, to create a backup of a DTR node, you can use:

```none
$ docker run -i --rm docker/dtr backup \
  --ucp-url <ucp-url> \
  --ucp-insecure-tls \
  --existing-replica-id <replica-id> \
  --ucp-username <ucp-admin> \
  --ucp-password <ucp-password> > /tmp/backup.tar
```

Where:

* `--ucp-url` is the address of UCP,
* `--ucp-insecure-tls` is to trust the UCP TLS certificate,
* `--existing-replica-id` is the id of the replica to backup,
* `--ucp-username`, and `--ucp-password` are the credentials of a UCP administrator.

To avoid having to pass the password as a command line parameter, you may
instead use the following approach in bash:

```none
$ read -sp 'ucp password: ' PASS; UCP_PASSWORD=$PASS docker run -i --rm -e UCP_PASSWORD docker/dtr backup \
  --ucp-url <ucp-url> \
  --ucp-insecure-tls \
  --existing-replica-id <replica-id> \
  --ucp-username <ucp-admin> > /tmp/backup.tar
```

This puts the password into a shell variable which is then passed into the
docker client command with the -e flag which in turn relays the password to the
DTR bootstrapper.

## Testing backups

To validate that the backup was correctly performed, you can print the contents
of the tar file created:

```none
$ tar -tf /tmp/backup.tar
```

The structure of the archive should look something like this:

```none
dtr-backup-v2.2.1/
dtr-backup-v2.2.1/rethink/
dtr-backup-v2.2.1/rethink/properties/
dtr-backup-v2.2.1/rethink/properties/0
...
```

To really test that the backup works, you must make a copy of your UCP cluster
by backing it up and restoring it onto separate machines. Then you can restore
DTR there from your backup and verify that it has all the data you expect to
see.

## Restore DTR data

You can restore a DTR node from a backup using the `restore` command.

Note that backups are tied to specific DTR versions and are guaranteed to work
only with those DTR versions. You can backup/restore across patch versions
at your own risk, but not across minor versions as those require more complex
migrations.

Before restoring DTR, make sure that you are restoring it on the same UCP
cluster or you've also restored UCP using its restore command. DTR does not
manage users, orgs or teams so if you try to
restore DTR on a cluster other than the one it was backed up on, DTR
repositories will be associated with users that don't exist and it will appear
as if the restore operation didn't work.

Note that to restore DTR, you must first remove any left over containers from
the previous installation. To do this, see the [uninstall
documentation](../install/uninstall.md).

The restore command performs a fresh installation of DTR, and reconfigures it with
the configuration created during a backup. The command starts by installing DTR.
Then it restores the configurations from the backup and then restores the
repository metadata. Finally, it applies all of the configs specified as flags to
the restore command.

After restoring DTR, you must make sure that it's configured to use the same
storage backend where it can find the image data. If the image data was backed
up separately, you must restore it now.

Finally, if you are using security scanning, you must re-fetch the security
scanning database through the online update or by uploading the offline tar. See
the [security scanning configuration](../admin/configure/set-up-vulnerability-scans.md)
for more detail.

You can check the
[reference documentation](../../reference/cli/restore.md), for the
restore command to learn about all the available flags.

As an example, to install DTR on the host and restore its
state from an existing backup:

```none
# Install and restore configurations from an existing backup
$ docker run -i --rm \
  docker/dtr restore \
  --ucp-url <ucp-url> \
  --ucp-insecure-tls \
  --ucp-username <ucp-admin> \
  --ucp-password <ucp-password> \
  --dtr-load-balancer <dtr-domain-name> < /tmp/backup.tar
```

Where:

* `--ucp-url` is the address of UCP,
* `--ucp-insecure-tls` is to trust the UCP TLS certificate,
* `--ucp-username`, and `--ucp-password` are the credentials of a UCP administrator,
* `--dtr-load-balancer` is the domain name or ip where DTR can be reached.

Note that if you want to avoid typing your password into the terminal you must pass
it in as an environment variable using the same approach as for the backup command:

```none
$ read -sp 'ucp password: ' PASS; UCP_PASSWORD=$PASS docker run -i --rm -e UCP_PASSWORD docker/dtr restore ...
```

After you successfully restore DTR, you can join new replicas the same way you
would after a fresh installation.

## Where to go next

* [Set up high availability](configure/set-up-high-availability.md)
* [DTR architecture](../architecture.md)
