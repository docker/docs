---
description: Learn how to backup your Docker Universal Control Plane swarm, and
  to recover your swarm from an existing backup.
keywords: ucp, backup, restore, recovery
title: Backups and disaster recovery
---

When you decide to start using Docker Universal Control Plane on a production
setting, you should
[configure it for high availability](configure/set-up-high-availability.md).

The next step is creating a backup policy and disaster recovery plan.

## Backup policy

As part of your backup policy you should regularly create backups of UCP.

To create a UCP backup, you can run the `{{ page.ucp_org }}/{{ page.ucp_repo }}:{{ page.ucp_version }} backup` command
on a single UCP manager. This command creates a tar archive with the
contents of all the [volumes used by UCP](../architecture.md) to persist data
and streams it to stdout.

You only need to run the backup command on a single UCP manager node. Since UCP
stores the same data on all manager nodes, you only need to take periodic
backups of a single manager node.

To create a consistent backup, the backup command temporarily stops the UCP
containers running on the node where the backup is being performed. User
resources, such as services, containers, and stacks are not affected by this
operation and will continue operating as expected. Any long-lasting `exec`,
`logs`, `events` or `attach` operations on the affected manager node will
be disconnected.

Additionally, if UCP is not configured for high availability, you will be
temporarily unable to:

* Log in to the UCP Web UI
* Perform CLI operations using existing client bundles

To minimize the impact of the backup policy on your business, you should:

* Configure UCP for [high availability](configure/set-up-high-availability.md).
  This allows load-balancing user requests across multiple UCP manager nodes.
* Schedule the backup to take place outside business hours.

## Backup command

The example below shows how to create a backup of a UCP manager node and
verify its contents:

```none
# Create a backup, encrypt it, and store it on /tmp/backup.tar
$ docker container run --log-driver none --rm -i --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock \
  {{ page.ucp_org }}/{{ page.ucp_repo }}:{{ page.ucp_version }} backup --interactive > /tmp/backup.tar

# Ensure the backup is a valid tar and list its contents
# In a valid backup file, over 100 files should appear in the list
# and the `./ucp-node-certs/key.pem` file should be present
$ tar --list -f /tmp/backup.tar
```

A backup file may optionally be encrypted using a passphrase, as in the
following example:

```none
# Create a backup, encrypt it, and store it on /tmp/backup.tar
$ docker container run --log-driver none --rm -i --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock \
  {{ page.ucp_org }}/{{ page.ucp_repo }}:{{ page.ucp_version }} backup --interactive \
  --passphrase "secret" > /tmp/backup.tar

# Decrypt the backup and list its contents
$ gpg --decrypt /tmp/backup.tar | tar --list
```

### Security-Enhanced Linux (SELinux)

For Docker EE 17.06 or higher, if the Docker engine has SELinux enabled,
which is typical for RHEL hosts, you need to include `--security-opt label=disable`
in the `docker` command:

```bash
$ docker container run --security-opt label=disable --log-driver none --rm -i --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock \
  {{ page.ucp_org }}/{{ page.ucp_repo }}:{{ page.ucp_version }} backup --interactive > /tmp/backup.tar
```

To find out whether SELinux is enabled in the engine, view the host's
`/etc/docker/daemon.json` file and search for the string
`"selinux-enabled":"true"`.

## Restore your swarm

The restore command can be used to create a new UCP swarm from a backup file.
When restoring, make sure you use the same version of the `docker/ucp` image that
you've used to create the backup. After the restore operation is complete, the
following data will be recovered from the backup file:

* Users, teams, and permissions.
* All UCP configuration options available under `Admin Settings`, like the
  Docker EE subscription license, scheduling options, content trust and
  authentication backends.

There are two ways to restore a UCP swarm:

* On a manager node of an existing swarm, which is not part of a UCP
  installation. In this case, a UCP swarm will be restored from the backup.
* On a docker engine that isn't participating in a swarm. In this case, a new
  swarm is created and UCP is restored on top.

To restore an existing UCP installation from a backup, you need to
uninstall UCP from the swarm by using the `uninstall-ucp` command.
[Learn to uninstall a UCP swarm](install/uninstall.md).

The example below shows how to restore a UCP swarm from an existing backup
file, presumed to be located at `/tmp/backup.tar`:

```none
$ docker container run --rm -i --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock  \
  {{ page.ucp_org }}/{{ page.ucp_repo }}:{{ page.ucp_version }} restore < /tmp/backup.tar
```

If the backup file is encrypted with a passphrase, you will need to provide the
passphrase to the restore operation:

```none
$ docker container run --rm -i --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock  \
  {{ page.ucp_org }}/{{ page.ucp_repo }}:{{ page.ucp_version }} restore --passphrase "secret" < /tmp/backup.tar
```

The restore command may also be invoked in interactive mode, in which case the
backup file should be mounted to the container rather than streamed through
stdin:

```none
$ docker container run --rm -i --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v /tmp/backup.tar:/config/backup.tar \
  {{ page.ucp_org }}/{{ page.ucp_repo }}:{{ page.ucp_version }} restore -i
```

## Disaster recovery

In the event where half or more manager nodes are lost and cannot be recovered
to a healthy state, the system is considered to have lost quorum and can only be
restored through the following disaster recovery procedure.

It is important to note that this procedure is not guaranteed to succeed with
no loss of running services or configuration data. To properly protect against
manager failures, the system should be configured for
[high availability](configure/set-up-high-availability.md).

1. On one of the remaining manager nodes, perform `docker swarm init
   --force-new-cluster`. You may also need to specify an
   `--advertise-addr` parameter which is equivalent to the `--host-address`
   parameter of the `docker/ucp install` operation. This will instantiate a new
   single-manager swarm by recovering as much state as possible from the
   existing manager. This is a disruptive operation and existing tasks may be
   either terminated or suspended.
2. Obtain a backup of one of the remaining manager nodes if one is not already
   available.
3. If UCP is still installed on the swarm, uninstall UCP using the
   `uninstall-ucp` command.
4. Perform a restore operation on the recovered swarm manager node.
5. Log in to UCP and browse to the nodes page, or use the CLI `docker node ls`
   command.
6. If any nodes are listed as `down`, you'll have to manually [remove these
   nodes](../configure/scale-your-cluster.md) from the swarm and then re-join
   them using a `docker swarm join` operation with the swarm's new join-token.

## Where to go next

* [Set up high availability](configure/set-up-high-availability.md)
* [UCP architecture](../architecture.md)
