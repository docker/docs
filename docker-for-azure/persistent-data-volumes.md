---
description: Persistent data volumes
keywords: azure persistent data volumes
title: Docker for Azure persistent data volumes
---

{% include d4a_buttons.md %}

## What is Cloudstor?


Cloudstor is a modern volume plugin built by Docker. It comes pre-installed and
pre-configured in Docker Swarms deployed on Docker for Azure. Docker swarm mode
tasks and regular Docker containers can use a volume created with Cloudstor to
mount a persistent data volume. The volume stays attached to the swarm tasks no
matter which swarm node they get scheduled on or migrated to. Cloudstor relies
on shared storage infrastructure provided by Azure (specifically File Storage
shares exposed over SMB) to allow swarm tasks to create/mount their persistent
volumes on any node in the swarm.

> **Note**: Direct attached storage, which is used to satisfy very low latency /
> high IOPS requirements, is not yet supported.

You can share the same volume among tasks running the same service, or you can
use a unique volume for each task.

## Use Cloudstor

After initializing or joining a swarm on Docker for Azure, connect to any swarm
manager using SSH. Verify that the CloudStor plugin is already installed and
configured for the stack or resource group:

```bash
$ docker plugin ls

ID                  NAME                        DESCRIPTION                       ENABLED
f416c95c0dcc        cloudstor:azure             cloud storage plugin for Docker   true
```

The following examples show how to create swarm services that require data
persistence using the `--mount` flag and specifying Cloudstor as the volume
driver.

### Share the same volume among tasks

If you specify a static value for the `source` option to the `--mount` flag, a
single volume is shared among the tasks participating in the service.

Cloudstor volumes can be created to share access to persistent data across all tasks in a swarm service running in multiple nodes. Example:

```bash
$ docker service create \
  --replicas 5 \
  --name ping1 \
  --mount type=volume,volume-driver=cloudstor:azure,source=sharedvol1,destination=/shareddata \
  alpine ping docker.com
```

In this example, all replicas/tasks of the service `ping1` share the same
persistent volume `sharedvol1` mounted at `/shareddata` path within the
container. Docker takes care of interacting with the Cloudstor plugin to ensure
that the common backing store is mounted on all nodes in the swarm where service
tasks are scheduled. Your application needs to be designed to ensure that tasks
do not write concurrently on the same file at the same time, to protect against
data corruption.

You can verify that the same volume is shared among all the tasks by logging
into one of the task containers, writing a file under `/shareddata/`, and
logging into another task container to verify that the file is available there
as well.

### Use a unique volume per task

You can use a templatized notation with the `docker service create` CLI to
create and mount a unique Cloudstor volume for each task in a swarm service.

It is possible to use the templatized notation to indicate to Docker Swarm that a unique Cloudstor volume be created and mounted for each replica/task of a service. This may be useful if the tasks write to the same file under the same path which may lead to corruption in case of shared storage. Example:

```bash
{% raw %}
$ docker service create \
  --replicas 5 \
  --name ping2 \
  --mount type=volume,volume-driver=cloudstor:azure,source={{.Service.Name}}-{{.Task.Slot}}-vol,destination=/mydata \
  alpine ping docker.com
{% endraw %}
```

A unique volume is created and mounted for each task participating in the
`ping2` service. Each task mounts its own volume at `/mydata/` and all files
under that mountpoint are unique to the task mounting the volume.

If a task is rescheduled on a different node after the volume is created and
mounted, Docker interacts with the Cloudstor plugin to create and mount the
volume corresponding to the task on the new node for the task.

It is highly recommended that you use the `.Task.Slot` template to ensure that
task `N` always gets access to volume `N`, no matter which node it is executing
on/scheduled to.

### Volume options

Cloudstor creates a new File Share in Azure File Storage for each volume and
uses SMB to mount these File Shares. SMB has limited compatibility with generic
Unix file ownership and permissions-related operations. Certain containers, such
as `jenkins` and `gitlab`, define specific users and groups which perform different
file operations. These types of workloads require the Cloudstor volume to be
mounted with the UID/GID of the user specified in the Dockerfile or setup scripts. Cloudstor allows for this scenario and
provides greater control over default file permissions by exposing the following
volume options that map to SMB parameters used for mounting the backing file
share.

| Option     | Description                                                                                             | Default                 |
|:-----------|:--------------------------------------------------------------------------------------------------------|:------------------------|
| `uid`      | User ID that will own all files on the volume.                                                          | `0` = `root`            |
| `gid`      | Group ID that will own all files on the volume.                                                         | `0` = `root`            |
| `filemode` | Permissions for all files on the volume.                                                                | `0777`                  |
| `dirmode`  | Permissions for all directories on the volume.                                                          | `0777`                  |
| `share`    | Name to associate with file share so that the share can be easily located in the Azure Storage Account. | MD5 hash of volume name |

This example sets `uid` to 1000 and `share` to `sharedvol` rather than a md5 hash:

```bash
$ docker service create \
  --replicas 5 \
  --name ping1 \
  --mount type=volume,volume-driver=cloudstor:azure,source=sharedvol1,destination=/shareddata,volume-opt=uid=1000,volume-opt=share=sharedvol \
  alpine ping docker.com
```

#### List or remove volumes created by Cloudstor

You can use `docker volume ls` on any node to enumerate all volumes created by
Cloudstor across the swarm.

You can use `docker volume rm [volume name]` to remove a Cloudstor volume from
any node. If you remove a volume from one node, make sure it is not being used
by another active node, since those tasks/containers in another node will lose
access to their data.
