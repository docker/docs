---
description: Persistent data volumes
keywords: aws persistent data volumes
title: Docker for AWS persistent data volumes
---

{% include d4a_buttons.md %}

## What is Cloudstor?

Cloudstor is a modern volume plugin built by Docker. It comes pre-installed and
pre-configured in Docker swarms deployed through Docker for AWS. Docker swarm
mode tasks and regular Docker containers can use a volume created with
Cloudstor to mount a persistent data volume. In Docker for AWS, Cloudstor has
two `backing` options:

- `relocatable` data volumes are backed by EBS.
- `shared` data volumes are backed by EFS.

When you use the Docker CLI to create a swarm service along with the persistent
volumes used by the service tasks, you can create three different types of
persistent volumes:

- Unique `relocatable` Cloudstor volumes mounted by each task in a swarm service.
- Global `shared` Cloudstor volumes mounted by all tasks in a swarm service.
- Unique `shared` Cloudstor volumes mounted by each task in a swarm service.

 Examples of each type of volume are described below.

## Relocatable Cloudstor volumes

 Workloads running in a Docker service that require access to low latency/high
 IOPs persistent storage, such as a database engine, can use a `relocatable`
 Cloudstor volume backed by EBS. When you create the volume, you can specify the
 type of EBS volume appropriate for the workload (such as `gp2`, `io1`, `st1`,
 `sc1`). Each `relocatable` Cloudstor volume is backed by a single EBS volume.

 If a swarm task using a `relocatable` Cloudstor volume gets rescheduled to
 another node within the same availability zone as the original node where the
 task was running, Cloudstor detaches the backing EBS volume from the original
 node and attaches it to the new target node automatically.

 If the swarm task gets rescheduled to a node in a different availability zone,
 Cloudstor transfers the contents of the backing EBS volume to the destination
 availability zone using a snapshot, and cleans up the EBS volume in the
 original availability zone. To minimize the time necessary to create the
 snapshot to transfer data across availability zones, Cloudstor periodically
 takes snapshots of EBS volumes to ensure there is never a large number of
 writes that need to be transferred as part of the final snapshot when
 transferring the EBS volume across availability zones.

 Typically the snapshot-based transfer process across availability zones takes
 between 2 and 5 minutes unless the work load is write-heavy. For extremely
 write-heavy workloads generating several GBs of fresh/new data every few
 minutes, the transfer may take longer than 5 minutes. The time required to
 snapshot and transfer increases sharply beyond 10 minutes if more than 20 GB of
 writes have been generated since the last snapshot interval. A swarm task is
 not started until the volume it mounts becomes available.

 Sharing/mounting the same Cloudstor volume backed by EBS among multiple tasks
 is not a supported scenario and leads to data loss. If you need a Cloudstor
 volume to share data between tasks, choose the appropriate EFS backed `shared`
 volume option. Using a `relocatable` Cloudstor volume backed by EBS is
 supported on all AWS regions that support EBS. The default `backing` option is
 `relocatable` if EFS support is not selected during setup/installation or if
 EFS is not supported in a region.

## Shared Cloudstor volumes

When multiple swarm service tasks need to share data in a persistent storage
volume, you can use a `shared` Cloudstor volume backed by EFS. Such a volume and
its contents can be mounted by multiple swarm service tasks without the risk of
data loss, since EFS makes the data available to all swarm nodes over NFS.

When swarm tasks using a `shared` Cloudstor volume get rescheduled from one node
to another within the same or across different availability zones, the
persistent data backed by EFS volumes is always available. `shared` Cloudstor
volumes only work in those AWS regions where EFS is supported. If EFS
Support is selected during setup/installation, the default "backing" option for
Cloudstor volumes is set to `shared` so that EFS is used by default.

`shared` Cloudstor volumes backed by EFS (or even EFS MaxIO) may not be ideal
for workloads that require very low latency and high IOPSs. For performance
details of EFS backed `shared` Cloudstor volumes, see [the AWS performance
guidelines](http://docs.aws.amazon.com/efs/latest/ug/performance.html).

## Use Cloudstor

After initializing or joining a swarm on Docker for AWS, connect to any swarm
manager using SSH. Verify that the CloudStor plugin is already installed and
configured for the stack or resource group:

```bash
$ docker plugin ls

ID                  NAME                        DESCRIPTION                       ENABLED
f416c95c0dcc        cloudstor:aws               cloud storage plugin for Docker   true
```

The following examples show how to create swarm services that require data
persistence using the `--mount` flag and specifying Cloudstor as the volume
driver.

### Share the same volume among tasks using EFS

In those regions where EFS is supported and EFS support is enabled during
deployment of the Cloud Formation template, you can use `shared` Cloudstor
volumes to share access to persistent data across all tasks in a swarm service
running in multiple nodes, as in the following example:

```bash
$ docker service create \
  --replicas 5 \
  --name ping1 \
  --mount type=volume,volume-driver=cloudstor:aws,source=sharedvol1,destination=/shareddata \
  alpine ping docker.com
```

All replicas/tasks of the service `ping1` share the same persistent volume
`sharedvol1` mounted at `/shareddata` path within the container. Docker takes
care of interacting with the Cloudstor plugin to ensure that EFS is mounted on
all nodes in the swarm where service tasks are scheduled. Your application needs
to be designed to ensure that tasks do not write concurrently on the same file
at the same time, to protect against data corruption.

You can verify that the same volume is shared among all the tasks by logging
into one of the task containers, writing a file under `/shareddata/`, and
logging into another task container to verify that the file is available there
as well.

The only option available for EFS is `perfmode`. You can set `perfmode` to
`maxio` for high IO throughput:

{% raw %}
```bash
$ docker service create \
  --replicas 5 \
  --name ping3 \
  --mount type=volume,volume-driver=cloudstor:aws,source={{.Service.Name}}-{{.Task.Slot}}-vol5,destination=/mydata,volume-opt=perfmode=maxio \
  alpine ping docker.com
```
{% endraw %}

You can also create `shared` Cloudstor volumes using the
`docker volume create` CLI:

```bash
$ docker volume create -d "cloudstor:aws" --opt backing=shared mysharedvol1
```

### Use a unique volume per task using EBS

If EBS is available and enabled, you can use a templatized notation with the
`docker service create` CLI to create and mount a unique `relocatable` Cloudstor
volume backed by a specified type of EBS for each task in a swarm service. New
EBS volumes typically take a few minutes to be created. Besides
`backing=relocatable`, the following volume options are available:

| Option    | Description                                                                                                                                                                                                                                                                                                    |
|:----------|:---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `size`    | Required parameter that indicates the size of the EBS volumes to create in GB.                                                                                                                                                                                                                                 |
| `ebstype` | Optional parameter that indicates the type of the EBS volumes to create (`gp2`, `io1`, `st1`, `sc1`}. The default `ebstype` is Standard/Magnetic. For further details about EBS volume types, see the [EBS volume type documentation](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/EBSVolumeTypes.html). |
| `iops`    | Required if `ebstype` specified is `io1`, which enables provisioned IOPs. Needs to be in the appropriate range as required by EBS.                                                                                                                                                                                       |

Example usage:

{% raw %}
```bash
$ docker service create \
  --replicas 5 \
  --name ping3 \
  --mount type=volume,volume-driver=cloudstor:aws,source={{.Service.Name}}-{{.Task.Slot}}-vol,destination=/mydata,volume-opt=backing=relocatable,volume-opt=size=25,volume-opt=ebstype=gp2 \
  alpine ping docker.com
```
{% endraw %}

The above example creates and mounts a distinct Cloudstor volume backed by 25 GB EBS
volumes of type `gp2` for each task of the `ping3` service. Each task mounts its
own volume at `/mydata/` and all files under that mountpoint are unique to the
task mounting the volume.

It is highly recommended that you use the `.Task.Slot` template to ensure that
task `N` always gets access to volume `N`, no matter which node it is executing
on/scheduled to. The total number of EBS volumes in the swarm should be kept
below `12 * (minimum number of nodes that are expected to be present at any
time)` to ensure that EC2 can properly attach EBS volumes to a node when another
node fails. Use EBS volumes only for those workloads where low latency and high
IOPs is absolutely necessary.

You can also create EBS backed volumes using the `docker volume create` CLI:

```bash
$ docker volume create \
  -d "cloudstor:aws" \
  --opt ebstype=io1 \
  --opt size=25 \
  --opt iops=1000 \
  --opt backing=relocatable \
  mylocalvol1
```

Sharing the same `relocatable` Cloudstor volume across multiple tasks of a
service or across multiple independent containers is not supported when
`backing=relocatable` is specified. Attempting to do so results in IO errors.

### Use a unique volume per task using EFS

If EFS is available and enabled, you can use templatized notation to create and
mount a unique EFS-backed volume into each task of a service. This is useful if
you already have too many EBS volumes or want to reduce the amount of time it
takes to transfer volume data across availability zones.

{% raw %}
```bash
$ docker service create \
  --replicas 5 \
  --name ping2 \
  --mount type=volume,volume-driver=cloudstor:aws,source={{.Service.Name}}-{{.Task.Slot}}-vol,destination=/mydata \
  alpine ping docker.com
```
{% endraw %}

Here, each task has mounted its own volume at `/mydata/` and the files under
that mountpoint are unique to that task.

When a task with only `shared` EFS volumes mounted is rescheduled on a different
node, Docker interacts with the Cloudstor plugin to create and mount the volume
corresponding to the task on the node where the task is rescheduled. Since data
on EFS is available to all swarm nodes and can be quickly mounted and accessed,
the rescheduling process for tasks using EFS-backed volumes typically takes a
few seconds, as compared to several minutes when using EBS.

It is highly recommended that you use the `.Task.Slot` template to ensure that
task `N` always gets access to volume `N` no matter which node it is executing
on/scheduled to.

### List or remove volumes created by Cloudstor

You can use `docker volume ls` on any node to enumerate all volumes created by
Cloudstor across the swarm.

You can use `docker volume rm [volume name]` to remove a Cloudstor volume from
any node. If you remove a volume from one node, make sure it is not being used
by another active node, since those tasks/containers in another node lose
access to their data.

Before deleting a Docker4AWS stack through CloudFormation, you should remove all
`relocatable` Cloudstor volumes using `docker volume rm` from within the stack.
EBS volumes corresponding to `relocatable` Cloudstor volumes are not
automatically deleted as part of the CloudFormation stack deletion. To list any
`relocatable` Cloudstor volumes and delete them after removing the Docker4AWS
stack where the volumes were created, go to the AWS portal or CLI and set a
filter with tag key set to `StackID` and the tag value set to the md5 hash of
the CloudFormation Stack ID (typical format:
`arn:aws:cloudformation:us-west-2:ID:stack/swarmname/GUID`).
