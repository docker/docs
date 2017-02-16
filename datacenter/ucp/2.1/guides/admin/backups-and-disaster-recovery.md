---
description: Learn how to backup your Docker Universal Control Plane cluster, and
  to recover your cluster from an existing backup.
keywords: docker, ucp, backup, restore, recovery
title: Backups and disaster recovery
---

When you decide to start using Docker Universal Control Plane on a production
setting, you should
[configure it for high availability](configure/set-up-high-availability.md).

The next step is creating a backup policy and disaster recovery plan.

## Backup policy

As part of your backup policy you should regularly create backups of UCP.

To create a UCP backup, you may use the `{{ page.docker_image }} backup` command
against a single UCP manager, according to the instructions in the next section.
This command creates a tar archive with the contents of all the [volumes used by
UCP](../architecture.md) to persist data and streams it to stdout.

You only need to run the backup command on a single UCP manager node. Since UCP
stores the same data on all manager nodes, you do not need to capture periodic
backups from more than one manager node.

To create a consistent backup, the backup command temporarily stops the UCP
containers running on the node where the backup is being performed. User
resources, such as services, containers and stacks are not affected by this
operation and will continue operating as expected. Any long-lasting `exec`,
`logs`, `events` or `attach` operations against the affected manager node will
be disconnected.

Additionally, if UCP is not configured for high availability, you will be
temporarily unable to:
* Log in to the UCP Web UI
* Perform CLI operations using existing client bundles

To minimize the impact of the backup policy on your business, you should:
* Configure UCP for high availability. This allows load-balancing user requests
across multiple UCP manager nodes.
* Schedule the backup to take place outside business hours.

## Backup command

The example below shows how to create a backup of a UCP manager node and
verify its contents:

```bash
# Create a backup, encrypt it, and store it on /tmp/backup.tar
$ docker run --rm -i --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock \
  {{ page.docker_image }} backup --interactive /tmp/backup.tar

# Ensure the backup is a valid tar and list its contents
$ tar --list /tmp/backup.tar
```

A backup file may optionally be encrypted using a passphrase, as in the
following example:

```bash
# Create a backup, encrypt it, and store it on /tmp/backup.tar
$ docker run --rm -i --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock \
  {{ page.docker_image }} backup --interactive \
  --passphrase "secret" > /tmp/backup.tar

# Decrypt the backup and list its contents
$ gpg --decrypt /tmp/backup.tar | tar --list
```

## Restore your cluster

The restore command can be used to create a new UCP cluster from a backup file.
After the restore operation is complete, the following data will be recovered
from the backup file:
* Users, Teams and Permissions.
* All UCP Configuration options available under `Admin Settings`, such as the
DDC Subscription license, scheduling options, Content Trust and authentication
backends.

There restore operation can be performed in any of three environments:
* On a manager node of an existing swarm, which is not part of a UCP
installation. In this case, a UCP cluster will be restored from the backup.
* On a docker engine that is not participating in a swarm. In this case, a new
swarm will be created and UCP will be restored on top

In order to restore an existing UCP installation from a backup, you will need to
first uninstall UCP from the cluster by using the `uninstall-ucp` command

The example below shows how to restore a UCP cluster from an existing backup
file:

```bash
$ docker run --rm -i --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock  \
  {{ page.docker_image }} restore < backup.tar
```

If the backup file is encrypted with a passphrase, you will need to provide the
passphrase to the restore operation:

```bash
$ docker run --rm -i --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock  \
  {{ page.docker_image }} restore --passphrase "secret" < /tmp/backup.tar
```

The restore command may also be invoked in interactive mode, in which case the
backup file should be mounted to the container rather than streamed through
stdin:

```bash
$ docker run --rm -i --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v /tmp/backup.tar:/config/backup.tar \
  {{ page.docker_image }} restore -i
```

## Disaster Recovery

In the event where half or more manager nodes are lost and cannot be recovered
to a healthy state, the system is considered to have lost quorum and can only be
restored through the following disaster recovery procedure. 

It is important to note that this proceedure is not guaranteed to succeed with
no loss of running services or configuration data. To properly protect against
manager failures, the system should be configured for [high availability](configure/set-up-high-availability.md).

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

* [Set up high availability](configure/set-up-high-availability.md)
* [UCP architecture](../architecture.md)
