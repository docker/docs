<!--[metadata]>
+++
title = "Backups and disaster recovery"
description = "Learn how to backup your Docker Trusted Registry cluster, and to recover your cluster from an existing backup."
keywords = ["docker, registry, high-availability, backup, recovery"]
[menu.main]
parent="dtr_menu_high_availability"
identifier="dtr_backup_disaster_recovery"
weight=10
+++
<![end-metadata]-->


# Backups and disaster recovery

When you decide to start using Docker Trusted Registry on a production
setting, you should [configure it for high availability](high-availability.md).

The next step is creating a backup policy and disaster recovery plan.

## DTR data persistency

Docker Trusted Registry persists four kinds of data:

* Configurations: the cluster configurations are stored on a key-value store
that is replicated through all DTR nodes.
* Image and repository metadata: the information about the repositories and
images deployed. This information is replicated through all DTR nodes.
* Docker images: By default images are stored on the host of the filesystem
where DTR is installed.
* Certificates and keys: the certificates, public keys, and private keys that
are used for mutual TLS communication.

This data is persisted on the host machine using named volumes.
[Learn more about DTR named volumes](../architecture.md).

## Backup DTR data

To perform a backup of a DTR node, use the `docker/dtr backup`
command. This command creates a backup of DTR:

* Configurations,
* Repository metadata,
* Certificates and keys used by DTR.

These files are added to a tar archive, and the result is streamed to stdout.

The backup command does not create a backup of Docker images. You should
implement a separate backup policy for the Docker images, taking in
consideration whether your DTR installation is configured to store images on the
filesystem or using a cloud provider.

When creating a backup, the resulting .tar file contains sensitive information
like private keys. You should ensure the backups are stored securely.

To learn about the options available on the backup command, you can
[check the reference documentation](../reference/backup.md), or run:

```bash
$ docker run --rm -it docker/dtr backup --help
```

As an example, to create a backup of a DTR node, you can use:

```bash
$ docker run -it --rm docker/dtr backup \
  --insecure-tls --replica-id 8b6174866010 \
  --username admin --password password \
  --host 192.168.10.100 > /tmp/backup.tar
```

Where:

* `--insecure-tls` allows connecting to UCP without TLS,
* `--replica-id` specifies the DTR replica to backup,
* `--username, --password` are the credentials of a UCP admin user,
* `--host` is the IP address of UCP.

## Restore DTR data

You can restore a DTR node from a backup using the `docker/dtr restore`
command.
This command performs a fresh installation of DTR, and reconfigures it with
the configuration created during a backup.

The command starts by installing DTR, restores the configurations stored on
etcd, and then restores the repository metadata stored on RethinkDB. You
can use the `--config-only` option, to only restore the configurations stored
on etcd.

This command does not restore Docker images. You should implement a separate
restore procedure for the Docker images stored in your registry, taking in
consideration whether your DTR installation is configured to store images on
the filesystem or using a cloud provider.

To learn about the options available on the restore command, you can
[check the reference documentation](../reference/restore.md), or run:

```bash
$ docker run --rm -it docker/trusted-registry restore --help
```

As an example, to install DTR on the host at 192.168.10.101, and restore its
state from an existing backup:

```bash
$ docker run -i --rm -v /var/run/docker.sock:/var/run/docker.sock \
  docker/dtr restore \
  --insecure-tls \
  --username admin --password password \
  --host 192.168.10.100 --dtr-host 192.168.10.101 < /tmp/backup.tar
```

Where:

* `--insecure-tls` allows connecting to UCP without TLS,
* `--username, --password` are the credentials of a UCP admin user,
* `--host` is the IP address of UCP,
* `--dtr-host` is the IP address of the host where DTR is going to be installed.
