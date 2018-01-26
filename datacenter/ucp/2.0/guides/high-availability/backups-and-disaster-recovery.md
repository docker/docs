---
description: Learn how to backup your Docker Universal Control Plane cluster, and
  to recover your cluster from an existing backup.
keywords: docker, ucp, backup, restore, recovery
title: Backups and disaster recovery
---

When you decide to start using Docker Universal Control Plane on a production
setting, you should
[configure it for high availability](index.md).

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
containers and services are not affected by this.

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
```

## Restore command

The example below shows how to restore a UCP controller node from an existing
backup:

```none
$ docker run --rm -i --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock  \
  docker/ucp restore --passphrase "secret" < backup.tar
```

When restoring, make sure you use the same version of the `docker/dtr` image that you've used to create the backup.

## Restore your cluster

The restore command can be used to create a new UCP cluster from a backup file.
After the restore operation is complete, the following data will be copied from
the backup file:

* Users, Teams and Permissions.
* Cluster Configuration, such as the default Controller Port or the KV store
timeout.
* DDC Subscription License.
* Options on Scheduling, Content Trust, Authentication Methods and Reporting.

The restore operation may be performed against any Docker Engine, regardless of
swarm membership, as long as the target Engine is not already managed by a UCP
installation. If the Docker Engine is already part of a swarm, that swarm and
all deployed containers and services will be managed by UCP after the restore
operation completes.

As an example, if you have a cluster with three controller nodes, A, B, and C,
and your most recent backup was of node A:

1. Uninstall UCP from the swarm using the `uninstall-ucp` operation.
2. Restore one of the swarm managers, such as node B, using the most recent
   backup from node A.
3. Wait for all nodes of the swarm to become healthy UCP nodes.

You should now have your UCP cluster up and running.

Additionally, in the event where half or more controller nodes are lost and
cannot be recovered to a healthy state, the system can only be restored through
the following disaster recovery procedure. This
procedure is not guaranteed to succeed with no loss of either swarm services or
UCP configuration data:

1. On one of the remaining manager nodes, perform `docker swarm init
   --force-new-cluster`. This will instantiate a new single-manager swarm by
   recovering as much state as possible from the existing manager. This is a
   disruptive operation and any existing tasks will be either terminated or
   suspended.
2. Obtain a backup of one of the remaining manager nodes if one is not already
   available.
3. Perform a restore operation on the recovered swarm manager node.
4. For all other nodes of the cluster, perform a `docker swarm leave --force`
   and then a `docker swarm join` operation with the cluster's new join-token.
5. Wait for all nodes of the swarm to become healthy UCP nodes.

## Where to go next

* [Set up high availability](index.md)
* [UCP architecture](../architecture.md)
