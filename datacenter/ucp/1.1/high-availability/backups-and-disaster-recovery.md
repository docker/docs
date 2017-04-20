---
description: Learn how to backup your Docker Universal Control Plane cluster, and
  to recover your cluster from an existing backup.
keywords: docker, ucp, backup, restore, recovery
redirect_from:
- /ucp/high-availability/backups-and-disaster-recovery/
title: Backups and disaster recovery
---

When you decide to start using Docker Universal Control Plane on a production
setting, you should
[configure it for high availability](set-up-high-availability.md).

The next step is creating a backup policy and disaster recovery plan.

## Backup policy

Docker UCP nodes persist data using [named volumes](../architecture.md).

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

## Backup command

The example below shows how to create a backup of a UCP controller node:

```none
# Create a backup, encrypt it, and store it on /tmp/backup.tar
$ docker run --rm -i --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock \
  docker/ucp backup --interactive \
  --passphrase "secret" > /tmp/backup.tar

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

## Restore command

The example below shows how to restore a UCP controller node from an existing
backup.
When restoring, make sure you use the same version of the `docker/dtr` image that you've used to create the backup.

First find out the Id of the UCP replica you want to restore:

```none
$ docker run --rm --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock  \
  docker/ucp id
```

Then, run the restore command:

```none
$ docker run --rm -i --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock  \
  docker/ucp restore \
    --passphrase "secret" \
    --id <replica-id> < backup.tar
```


## Restore your cluster

Configuring UCP to have multiple controller nodes allows you tolerate a certain
amount of node failures. If multiple nodes fail at the same time, causing the
cluster to go down, you can use an existing backup to recover.

As an example, if you have a cluster with three controller nodes, A, B, and C,
and your most recent backup was of node A:

1. Stop controllers B and C with the `stop` command,
2. Restore controller A,
3. Uninstall UCP from controllers B and C,
4. Join nodes B and C as replica controllers to the cluster.

You should now have your cluster up and running.


## Where to go next

* [Set up high availability](set-up-high-availability.md)
* [UCP architecture](../architecture.md)
