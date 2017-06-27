---
description: Persistent data volumes
keywords: azure persistent data volumes
title: Docker for Azure persistent data volumes
---

{% include d4a_buttons.md %}

## What is Cloudstor?

Cloudstor is a modern volume plugin built by Docker. It comes pre-installed and pre-configured in Docker Swarms deployed through Docker for Azure. Docker Swarm mode tasks as well as regular Docker containers can use a volume created with Cloudstor to mount a persistent data volume. The volume  stays attached to the swarm tasks no matter which swarm node they get scheduled on or migrated to. Cloudstor relies on shared storage infrastructure provided by Azure (specifically File Storage shares exposed over SMB) to allow swarm tasks to create/mount their persistent volumes on any node in the Docker Swarm. In a future release we will introduce support for direct attached/relocatable storage to satisfy very low latency/high IOPs requirements.

## Use Cloudstor

After creating a swarm on Docker for Azure and connecting to any manager using SSH, verify that Cloudstor is already installed and configured for the stack/resource group:

```bash
$ docker plugin ls
ID                  NAME                        DESCRIPTION                       ENABLED
f416c95c0dcc        cloudstor:azure             cloud storage plugin for Docker   true
```

The following examples show how to create swarm services that require data persistence using the --mount flag and specifying Cloudstor as the driver.

### Share the same volume between tasks:

Cloudstor volumes can be created to share access to persistent data across all tasks in a swarm service running in multiple nodes. Example:

```bash
docker service create --replicas 5 --name ping1 \
    --mount type=volume,volume-driver=cloudstor:azure,source=sharedvol1,destination=/shareddata \
    alpine ping docker.com
```

Here all replicas/tasks of the service `ping1` share the same persistent volume `sharedvol1` mounted at `/shareddata` path within the container. Docker Swarm takes care of interacting with the Cloudstor plugin to make sure the common backing store is mounted on all nodes in the swarm where tasks of the service are scheduled. Each task needs to make sure they don't write concurrently on the same file at the same time and cause corruptions since the volume is shared.

With the above example, you can make sure that the volume is indeed shared by logging into one of the containers in one swarm node, writing to a file under `/shareddata/` and reading the file under `/shareddata/` from another container (in the same node or a different node).

### Use a unique volume per task:

It is possible to use the templatized notation to indicate to Docker Swarm that a unique Cloudstor volume be created and mounted for each replica/task of a service. This may be useful if the tasks write to the same file under the same path which may lead to corruption in case of shared storage. Example:

```bash
{% raw %}
docker service create --replicas 5 --name ping2 \
    --mount type=volume,volume-driver=cloudstor:azure,source={{.Service.Name}}-{{.Task.Slot}}-vol,destination=/mydata \
    alpine ping docker.com
{% endraw %}
```

Here the templatized notation is used to indicate to Docker Swarm that a unique volume be created and mounted for each replica/task of the service `ping2`. After initial creation of the volumes corresponding to the tasks that will use them (in the nodes the tasks are scheduled in), if a task is rescheduled on a different node, Docker Swarm will interact with the Cloudstor plugin to create and mount the volume corresponding to the task on the node the task got scheduled on. It's highly recommended that you use the `.Task.Slot` template to make sure task N always gets access to vol N no matter which node it is executing on/scheduled to.

In the above example, each task has it's own volume mounted at `/mydata/` and the files under there are unique to the task mounting the volume.

### Volume options

Cloudstor creates a new File Share in Azure File Storage for each volume and uses SMB to mount them. SMB however is fairly limited in the area of being compatible with generic Unix file ownership and permissions related operations. Certain workloads (e.g. Jenkins, Gitlab) define specific users and groups that perform different file operations and requires the Cloudstor volume to be mounted with the corresponding UID/GID. To allow for this scenario as well as greater control over default file permissions, Cloudstor exposes the following volume options that map to SMB parameters used for mounting the backing file share.

1. `uid` : User ID that will own all files on the volume. Default: 0 = root
2. `gid` : Group ID that will own all files on the volume. Default: 0 = root
3. `filemode` : Permissions for all files on the volume. Default: 0777
4. `dirmode` : Permissions for all directories on the volume. Default: 0777
5. `share` : Name to associate with file share so that the share can be easily located in the Azure Storage Account. Default: MD5 hash of volume name

Example usage with `uid` set to 1000 and share name set to `sharedvol` rather than a md5 hash:

```bash
docker service create --replicas 5 --name ping1 \
    --mount type=volume,volume-driver=cloudstor:azure,source=sharedvol1,destination=/shareddata,volume-opt=uid=1000,volume-opt=share=s
haredvol \
    alpine ping docker.com
```

#### List or remove volumes created by Cloudstor

You can use `docker volume ls` on any node to enumerate all volumes created by Cloudstor across the swarm. You can use `docker volume rm [volume name]` to remove a cloudstor volume from any node. Please be aware that if you remove a volume from one node it may still be under active usage in another node and those tasks in the other node will lose access to their data.
