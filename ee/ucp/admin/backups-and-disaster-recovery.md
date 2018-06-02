---
title: Backups and disaster recovery
description: Learn how to backup your Docker Universal Control Plane swarm, and
  to recover your swarm from an existing backup.
keywords: ucp, backup, restore, recovery
---

When you decide to start using Docker Universal Control Plane on a production
setting, you should
[configure it for high availability](configure/join-nodes/index.md).

The next step is creating a backup policy and disaster recovery plan.

## Data managed by UCP

UCP maintains data about:

| Data                  | Description                                                                                                          |
|:----------------------|:---------------------------------------------------------------------------------------------------------------------|
| Configurations        | The UCP cluster configurations, as shown by `docker config ls`, including Docker EE license and swarm and client CAs |
| Access control        | Permissions for teams to cluster resources, including resource sets, grants, and roles                               |
| Certificates and keys | The certificates, public keys, and private keys that are used for authentication and mutual TLS communication        |
| Metrics data          | Monitoring data gathered by UCP                                                                                      |
| Organizations         | Your users, teams, and orgs                                                                                          |
| Volumes               | All [UCP named volumes](../architecture/#volumes-used-by-ucp), which include all UCP component certs and data        |

This data is persisted on the host running UCP, using named volumes.
[Learn more about UCP named volumes](../ucp-architecture.md).

UCP won't backup your routing mesh settings. After restoring you need to
[re-enable the routing mesh](../interlock/deploy/index.md). If you've customized
your layer 7 routing deployment, you'll need to re-apply those customizations too.

## Backup steps

Back up your Docker EE components in the following order:

1. [Back up your swarm](/engine/swarm/admin_guide/#back-up-the-swarm)
2. Back up UCP
3. [Back up DTR](../../../../dtr/2.5/guides/admin/backups-and-disaster-recovery.md)

## Backup policy

As part of your backup policy you should regularly create backups of UCP.
DTR is backed up independently.
[Learn about DTR backups and recovery](../../../../dtr/2.5/guides/admin/backups-and-disaster-recovery.md).

To create a UCP backup, run the `{{ page.ucp_org }}/{{ page.ucp_repo }}:{{ page.ucp_version }} backup` command
on a single UCP manager. This command creates a tar archive with the
contents of all the [volumes used by UCP](../ucp-architecture.md) to persist data
and streams it to `stdout`. The backup doesn't include the swarm-mode state,
like service definitions and overlay network definitions.

You only need to run the backup command on a single UCP manager node. Since UCP
stores the same data on all manager nodes, you only need to take periodic
backups of a single manager node.

To create a consistent backup, the backup command temporarily stops the UCP
containers running on the node where the backup is being performed. User
resources, such as services, containers, and stacks are not affected by this
operation and will continue operating as expected. Any long-lasting `exec`,
`logs`, `events`, or `attach` operations on the affected manager node will
be disconnected.

Additionally, if UCP is not configured for high availability, you will be
temporarily unable to:

* Log in to the UCP Web UI
* Perform CLI operations using existing client bundles

To minimize the impact of the backup policy on your business, you should:

* Configure UCP for [high availability](configure/join-nodes/index.md).
  This allows load-balancing user requests across multiple UCP manager nodes.
* Schedule the backup to take place outside business hours.

## Backup command

The example below shows how to create a backup of a UCP manager node and
verify its contents:

```bash
# Create a backup, encrypt it, and store it on /tmp/backup.tar
docker container run \
  --log-driver none --rm \
  --interactive \
  --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock \
  {{ page.ucp_org }}/{{ page.ucp_repo }}:{{ page.ucp_version }} backup \
  --id <ucp-instance-id> \
  --passphrase "secret" > /tmp/backup.tar

# Decrypt the backup and list its contents
$ gpg --decrypt /tmp/backup.tar | tar --list
```

### Security-Enhanced Linux (SELinux)

For Docker EE 17.06 or higher, if the Docker engine has SELinux enabled,
which is typical for RHEL hosts, you need to include `--security-opt label=disable`
in the `docker` command:

```bash
docker container run --security-opt label=disable --log-driver none --rm -i --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock \
  {{ page.ucp_org }}/{{ page.ucp_repo }}:{{ page.ucp_version }} backup --interactive > /tmp/backup.tar
```

To find out whether SELinux is enabled in the engine, view the host's
`/etc/docker/daemon.json` file and search for the string
`"selinux-enabled":"true"`.

## Restore UCP

To restore an existing UCP installation from a backup, you need to
uninstall UCP from the swarm by using the `uninstall-ucp` command.
[Learn to uninstall UCP](install/uninstall.md).

If you restore UCP using a different Docker swarm than the one where UCP was
previously deployed on, UCP will start using new TLS certificates. Existing
client bundles won't work anymore, so users need to download new ones. 

When restoring, make sure you use the same version of the `docker/ucp` image
that you've used to create the backup. The example below shows how to restore
UCP from an existing backup file, presumed to be located at
`/tmp/backup.tar`:

```none
docker container run --rm -i --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock  \
  {{ page.ucp_org }}/{{ page.ucp_repo }}:{{ page.ucp_version }} restore < /tmp/backup.tar
```

If the backup file is encrypted with a passphrase, you will need to provide the
passphrase to the restore operation:

```none
docker container run --rm -i --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock  \
  {{ page.ucp_org }}/{{ page.ucp_repo }}:{{ page.ucp_version }} restore --passphrase "secret" < /tmp/backup.tar
```

The restore command may also be invoked in interactive mode, in which case the
backup file should be mounted to the container rather than streamed through
`stdin`:

```none
docker container run --rm -i --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v /tmp/backup.tar:/config/backup.tar \
  {{ page.ucp_org }}/{{ page.ucp_repo }}:{{ page.ucp_version }} restore -i
```

### UCP and Swarm

UCP restore recovers the following assets from the backup file:

* Users, teams, and permissions.
* All UCP configuration options available under `Admin Settings`, like the
  Docker EE subscription license, scheduling options, content trust, and
  authentication backends.

UCP restore does not include swarm assets such as cluster membership, services, networks,
secrets, etc.  [Learn to backup a swarm](/engine/swarm/admin_guide/#back-up-the-swarm).

There are two ways to restore UCP:

* On a manager node of an existing swarm which does not have UCP installed.
  In this case, UCP restore will use the existing swarm.
* On a docker engine that isn't participating in a swarm. In this case, a new
  swarm is created and UCP is restored on top.

## Disaster recovery

In the event where half or more manager nodes are lost and cannot be recovered
to a healthy state, the system is considered to have lost quorum and can only be
restored through the following disaster recovery procedure. If your cluster has
lost quorum, you can still take a backup of one of the remaining nodes, but we
recommend making backups regularly.

It is important to note that this procedure is not guaranteed to succeed with
no loss of running services or configuration data. To properly protect against
manager failures, the system should be configured for
[high availability](configure/join-nodes/index.md).

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
   nodes](configure/scale-your-cluster.md) from the swarm and then re-join
   them using a `docker swarm join` operation with the swarm's new join-token.

## Where to go next

- [UCP architecture](../ucp-architecture.md)
- [Set up high availability](configure/join-nodes/index.md)
