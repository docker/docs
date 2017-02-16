---
description: Persistent Data Volumes
keywords: aws persistent data volumes
title: Docker for AWS Persistent Data Volumes
---

## What is Cloudstor?
Cloudstor is a Docker managed volume plugin that comes pre-installed (and pre-configured) in swarms deployed on Docker for AWS. Using volumes created through Cloudstor, swarm tasks can mount a persistent data volume that stays attached to the swarm tasks no matter which swarm node they get scheduled or migrated to. Currently Cloudstor heavily relies on shared storage infrastructure provided by the cloud platform to allow swarm tasks to create/mount their persistent volumes on any node in the swarm. In a future release we will introduce support for direct attached storage to satisfy very low latency/high IOPs requirements.

## Using Cloudstor
After creating a swarm on Docker for AWS and SSHing into any of the managers, you should be able to find cloudstor already installed and configured for the stack/resource group:

```
$ docker plugin ls
ID                  NAME                        DESCRIPTION                       ENABLED
f416c95c0dcc        docker4x/cloudstor:aws-v1.13.1-beta18   cloud storage plugin for Docker   true
```
**Note**: Make note of the plugin tag name, it will change between versions, and yours may be different then listed here.

You can create swarm services that requires data persistence using the `--mount` flag and specifying Cloudstor as the driver. Here are some examples:

### Swarm tasks mounting the same volume:
```
docker service create --replicas 5 --name ping1 --mount type=volume,volume-driver=docker4x/cloudstor:aws-v1.13.1-beta18,source=sharedvol1,destination=/shareddata alpine ping docker.com
```
Here all replicas/tasks of the service `ping1` share the same persistent volume `sharedvol1` mounted at `/shareddata` path within the container. Docker Swarm takes care of interacting with the Cloudstor plugin to make sure the common backing store is mounted on all nodes in the swarm where tasks of the service are scheduled. Each task needs to make sure they don't write concurrently on the same file at the same time and cause corruptions since the volume is shared.

With the above example, you can make sure that the volume is indeed shared by logging into one of the containers in one swarm node, writing to a file under `/shareddata/` and reading the file under `/shareddata/` from another container (in the same node or a different node).

### Swarm tasks mounting unique volumes:
```
docker service create --replicas 5 --name ping2 --mount type=volume,volume-driver=docker4x/cloudstor:aws-v1.13.1-beta18,source={{.Service.Name}}-{{.Task.Slot}}-vol,destination=/mydata alpine ping docker.com
```
Here the templatized notation is used to indicate to Docker Swarm that a unique volume be created and mounted for each replica/task of the service `ping2`. After initial creation of the volumes corresponding to the tasks they are attached to (in the nodes the tasks are scheduled in), if a task is rescheduled on a different node, Docker Swarm will interact with the Cloudstor plugin to create and mount the volume corresponding to the task on the node the task got scheduled on. It's highly recommended that you use the `.Task.Slot` template to make sure task N always gets access to vol N no matter which node it is executing on/scheduled to.

In the above example, each task has it's own volume mounted at `/mydata/` and the files under there are unique to the task mounting the volume.

### Listing and Removal of volumes created by Cloudstor
You can use `docker volume ls` to enumerate all volumes created on a node including those backed by Cloudstor. Note that if a swarm service task starts off in a node and has a Cloudstor volume associated and later gets rescheduled to a different node, `docker volume ls` in the initial node will continue to list the Cloudstor volume that was created for the task that no longer executes on the node although the volume is mounted elsewhere. Do NOT prune/rm the volumes that gets enumerated on a node without any tasks associated since these actions will result in data loss if the same volume is mounted in another node (i.e. the volume shows up in the `docker volume ls` output on another node in the swarm). We can try to detect this and block/handle in post-Beta.

### Choose your IO performance
In AWS, in case you want a higher level of IO performance like the maxIO mode for EFS, a perfmode parameter can be specified as volume-opt:
```
docker service create --replicas 5 --name ping3 --mount type=volume,volume-driver=docker4x/cloudstor:aws-v1.13.1-beta18,source={{.Service.Name}}-{{.Task.Slot}}-vol5,destination=/mydata,volume-opt=perfmode=maxio alpine ping docker.com
```
