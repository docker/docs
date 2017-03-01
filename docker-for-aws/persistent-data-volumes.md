---
description: Persistent data volumes
keywords: aws persistent data volumes
title: Docker for AWS persistent data volumes
---

## What is Cloudstor?

Cloudstor a volume plugin managed by Docker. It comes pre-installed and pre-configured in swarms deployed on Docker for AWS. Swarm tasks use a volume created through Cloudstor to mount a persistent data volume that stays attached to the swarm tasks no matter which swarm node they get scheduled or migrated to. Cloudstor relies on shared storage infrastructure provided by AWS to allow swarm tasks to create/mount their persistent volumes on any node in the swarm. In a future release we will introduce support for direct attached storage to satisfy very low latency/high IOPs requirements.

## Use Cloudstor

After creating a swarm on Docker for AWS and connecting to any manager using SSH, verify that Cloudstor is already installed and configured for the stack/resource group:

```bash
$ docker plugin ls
ID                  NAME                        DESCRIPTION                       ENABLED
f416c95c0dcc        docker4x/cloudstor:aws-v1.13.1-beta18   cloud storage plugin for Docker   true
```

**Note**: Make note of the plugin tag name, because it will change between versions, and yours may be different then listed here.

The following examples show how to create swarm services that require data persistence using the --mount flag and specifying Cloudstor as the driver.

### Share the same volume between tasks:

```bash
docker service create --replicas 5 --name ping1 \
    --mount type=volume,volume-driver=docker4x/cloudstor:aws-v1.13.1-beta18,source=sharedvol1,destination=/shareddata \
    alpine ping docker.com
```

Here all replicas/tasks of the service `ping1` share the same persistent volume `sharedvol1` mounted at `/shareddata` path within the container. Docker Swarm takes care of interacting with the Cloudstor plugin to make sure the common backing store is mounted on all nodes in the swarm where tasks of the service are scheduled. Each task needs to make sure they don't write concurrently on the same file at the same time and cause corruptions since the volume is shared.

With the above example, you can make sure that the volume is indeed shared by logging into one of the containers in one swarm node, writing to a file under `/shareddata/` and reading the file under `/shareddata/` from another container (in the same node or a different node).

### Use a unique volume per task:

```bash
{% raw %}
docker service create --replicas 5 --name ping2 \
    --mount type=volume,volume-driver=docker4x/cloudstor:aws-v1.13.1-beta18,source={{.Service.Name}}-{{.Task.Slot}}-vol,destination=/mydata \
    alpine ping docker.com
{% endraw %}
```

Here the templatized notation is used to indicate to Docker Swarm that a unique volume be created and mounted for each replica/task of the service `ping2`. After initial creation of the volumes corresponding to the tasks they are attached to (in the nodes the tasks are scheduled in), if a task is rescheduled on a different node, Docker Swarm will interact with the Cloudstor plugin to create and mount the volume corresponding to the task on the node the task got scheduled on. It's highly recommended that you use the `.Task.Slot` template to make sure task N always gets access to vol N no matter which node it is executing on/scheduled to.

In the above example, each task has it's own volume mounted at `/mydata/` and the files under there are unique to the task mounting the volume.

### List or remove volumes created by Cloudstor

You can use `docker volume ls` to enumerate all volumes created on a node including those backed by Cloudstor. Note that if a swarm service task starts off in a node and has a Cloudstor volume associated and later gets rescheduled to a different node, `docker volume ls` in the initial node will continue to list the Cloudstor volume that was created for the task that no longer executes on the node although the volume is mounted elsewhere. Do NOT prune/rm the volumes that gets enumerated on a node without any tasks associated since these actions will result in data loss if the same volume is mounted in another node (i.e. the volume shows up in the `docker volume ls` output on another node in the swarm). We can try to detect this and block/handle in post-Beta.

### Configure IO performance

If you want a higher level of IO performance like the maxIO mode for EFS, a perfmode parameter can be specified as volume-opt:

```bash
{% raw %}
docker service create --replicas 5 --name ping3 \
    --mount type=volume,volume-driver=docker4x/cloudstor:aws-v1.13.1-beta18,source={{.Service.Name}}-{{.Task.Slot}}-vol5,destination=/mydata,volume-opt=perfmode=maxio \
    alpine ping docker.com
{% endraw %}
```
