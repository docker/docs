---
description: Learn how to backup your Docker Trusted Registry cluster, and to recover your cluster from an existing backup.
keywords: docker, registry, high-availability, backup, recovery
redirect_from:
- /docker-trusted-registry/high-availability/backups-and-disaster-recovery/
title: Backups and disaster recovery
---

When you decide to start using Docker Trusted Registry on a production
setting, you should [configure it for high availability](index.md).

The next step is creating a backup policy and disaster recovery plan.

## DTR data persistency

Docker Trusted Registry persists:

* Configurations: the cluster configurations are stored on a key-value store
that is replicated through all DTR replicas.
* Repository metadata: the information about the repositories and
images deployed. This information is replicated through all DTR replicas.
* Certificates and keys: the certificates, public keys, and private keys that
are used for mutual TLS communication.

This data is persisted on the host where DTR is running, using named volumes.
[Learn more about DTR named volumes](../architecture.md).

DTR also persists Docker images on the filesystem of the host running DTR, or
on a cloud provider, depending on the way DTR is configured.

## Backup DTR data

To perform a backup of a DTR node, use the `docker/dtr backup` command. This
command creates a backup of DTR:

* Configurations,
* Repository metadata,
* Certificates and keys used by DTR.

These files are added to a tar archive, and the result is streamed to stdout.

The backup command does not create a backup of Docker images. You should
implement a separate backup policy for the Docker images, taking in
consideration whether your DTR installation is configured to store images on the
filesystem or using a cloud provider.

The backup command also doesn't create a backup of the users and organizations.
That data is managed by UCP, so when you create a UCP backup you're creating
a backup of the users and organizations metadata.

When creating a backup, the resulting .tar file contains sensitive information
like private keys. You should ensure the backups are stored securely.

You can check the
[reference documentation](../reference/backup.md), for the
backup command to learn about all the available flags.

As an example, to create a backup of a DTR node, you can use:

```bash
# Get the certificates used by UCP
$ curl https://<ucp-url>/ca > ucp-ca.pem

# Create the backup
$ docker run -i --rm docker/dtr backup \
  --ucp-url <ucp-url> \
  --ucp-ca "$(cat ucp-ca.pem)" \
  --existing-replica-id <replica-id> \
  --ucp-username <ucp-admin> \
  --ucp-password <ucp-password> > /var/tmp/backup.tar
```

Where:

* `--ucp-url` is the address of UCP,
* `--ucp-ca` is the UCP certificate authority,
* `--existing-replica-id` is the id of the replica to backup,
* `--ucp-username`, and `--ucp-password` are the credentials of a UCP administrator.

To validate that the backup was correctly performed, you can print the contents
of the tar file created:

```bash
$ tar -tf /tmp/backup.tar
```

## Restore DTR data

This command performs a fresh installation of DTR, and reconfigures it with
the configuration created during a backup.
When restoring, make sure you use the same version of the `docker/dtr` image that you've used to create the backup.

The command starts by installing DTR, restores the configurations stored on
etcd, and then restores the repository metadata stored on RethinkDB. You
can use the `--config-only` option, to only restore the configurations stored
on etcd.

This command does not restore Docker images. You should implement a separate
restore procedure for the Docker images stored in your registry, taking in
consideration whether your DTR installation is configured to store images on
the filesystem or using a cloud provider.

You can check the
[reference documentation](../reference/restore.md), for the
restore command to learn about all the available flags.


As an example, to install DTR on the host and restore its
state from an existing backup:

```bash
# Get the certificates used by UCP
$ curl https://<ucp-url>/ca > ucp-ca.pem

# Install and restore configurations from an existing backup
$ docker run -i --rm \
  docker/dtr restore \
  --ucp-url <ucp-url> \
  --ucp-ca "$(cat ucp-ca.pem)" \
  --ucp-username <ucp-admin> \
  --ucp-password <ucp-password> \
  --dtr-external-url <dtr-domain-name> < /var/tmp/backup.tar
```

Where:

* `--ucp-url` is the address of UCP,
* `--ucp-ca` is the UCP certificate authority,
* `--ucp-username`, and `--ucp-password` are the credentials of a UCP administrator,
* `--dtr-external-url` is the domain name or ip where DTR can be reached.


## Where to go next

* [Set up high availability](index.md)
* [DTR architecture](../architecture.md)