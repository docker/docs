<!--[metadata]>
+++
title ="Backups and disaster recovery"
description="Learn how to backup your Docker Universal Control Plane cluster, and to recover your cluster from an existing backup."
keywords= ["docker, ucp, backup, restore, recovery"]
[menu.main]
parent="mn_ucp_high_availability"
weight=20
+++
<![end-metadata]-->

# Backups and disaster recovery

When you decide to start using Docker Universal Control Plane on a production
setting, you should
[configure it for high availability](set-up-high-availability.md).

The next step is creating a backup policy and disaster recovery plan.

## Backup policy

Docker UCP nodes persist data using [named volumes](../architecture.md):

* Controller nodes persist cluster configurations, certificates, and keys
used to issue certificates and user bundles. This data is replicated on every
controller node in the cluster.
* Nodes are stateless. They only store certificates for mutual TLS, that
can be regenerated.

As part of your backup policy you should regularly create backups of the
controller nodes. Since the nodes used for running user containers don't
persist data, you can decide not to create any backups for them.

To perform a backup of a UCP controller node, use the `docker/ucp backup`
command. This creates a tar archive with the contents of the volumes used by
UCP on that node, and streams it to stdout.

To create a consistent backup, the backup command temporarily stops the UCP
containers running on the node where the backup is being performed. User
containers are not affected by this.

To have minimal impact on your business, you should:

* Schedule the backup to take place outside business hours.
* Configure UCP for high availability. This allows load-balancing user requests
across multiple UCP controller nodes.

## Backup UCP data

To learn about the options available on the `docker/ucp backup` command, you can
check the reference documentation, or run:

```bash
$ docker run --rm docker/ucp backup --help
```

When creating a backup, the resulting tar archive contains sensitive information
like private keys. To ensure this information is kept private you should run
the backup command with the `--passphrase` option. This encrypts
the backup with a passphrase of your choice.

The example below shows how to create a backup of a UCP controller node:

```bash
# Create a backup, encrypt it, and store it on /tmp/backup.tar
$ docker run --rm -i --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock \
  docker/ucp --interactive --passphrase "secret" > /tmp/backup.tar

Do you want proceed with the backup? (y/n):
$ y

INFO[0000] Temporarily Stopping local UCP containers to ensure a consistent backup
INFO[0000] Beginning backup
INFO[0001] Backup completed successfully
INFO[0002] Resuming stopped UCP containers

# Decrypt the backup and list its contents
$ gpg --decrypt /tmp/backup.tar | tar --list

Enter passphrase: secret

/ucp-client-root-ca/
./ucp-client-root-ca/cert.pem
./ucp-client-root-ca/config.json
./ucp-client-root-ca/key.pem
./ucp-cluster-root-ca/
# output snipped
```

## Where to go next

* [Set up high availability](set-up-high-availability.md)
* [UCP architecture](../architecture.md)
