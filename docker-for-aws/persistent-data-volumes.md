---
description: Persistent data volumes
keywords: aws persistent data volumes
title: Docker for AWS persistent data volumes
---

{% include d4a_buttons.md %}

## What is Cloudstor?

Cloudstor is a modern volume plugin managed by Docker. It comes pre-installed and pre-configured in Docker Swarms deployed through Docker for AWS. Docker Swarm tasks can use a volume created through Cloudstor to mount a persistent data volume. In Docker for AWS, Cloudstor can be used with two `backing` options: `local` (that uses EBS) and `shared` (that uses EFS) described below.

## Local Cloudstor volumes

 Workloads running in a Docker Service that require access to low latency/high IOPs persistent storage (e.g. a database engine) can use a `local` Cloudstor volume backed by EBS. It is possible to configure the specific type of EBS volume supported by AWS (e.g. gp2, st1, sc1) that the workload requires during volume creation. If the swarm task gets rescheduled to another node within the same availability zone (as the node where the task was running on), Cloudstor takes care of detaching and re-attaching the EBS to the appropriate node. If the swarm task gets rescheduled to a node in a different availability zone (from the node where the task was originally running), Cloudstor will transfer the contents of the EBS volume using a snapshot to the destination availability zone as well as clean up the EBS volume in the original Availability Zone. To minimize the time necessary to create the snapshot to transfer data across availability zones, Cloudstor periodically takes snapshots of the EBS volumes to ensure there is never a large amount of diff that will need to get transferred as part of the snapshot necessary during a cross availability zone task reschedule. Typically the snapshot based transfer process across availability zones takes a few minutes and for extremely write-heavy workloads, this may take several minutes. 
 
 Sharing/mounting the same Cloudstor volume backed by EBS among multiple tasks is not a supported scenario and will lead to data loss. If you need a Cloudstor volume to share data between tasks please read below for EFS backed options. `local` Cloudstor backed by EBS is supported on all AWS regions. It is the default "backing" option if EFS support is not selected during setup/installation or if EFS is not supported in a region.

## Shared Cloudstor volumes 

Workloads running in multiple swarm tasks that need to share data in persistent storage volumes can use a Cloudstor volume backed by EFS. Such a volume and it's contents can be mounted by multiple tasks on multiple swarm nodes simultaneously since EFS makes the data available to all Swarm nodes over NFS. `shared` Cloudstor volumes will only work in those AWS regions where EFS is supported (e.g. Ohio, Oregon, N.Virginia, Sydney). If EFS Support is selected during setup/installation, the default "backing" option for Cloudstor volumes is set to Shared so that EFS is used by default. Please note that `shared` Cloudstor volumes backed by EFS (or even EFS MaxIO) may not be ideal for workloads that have very low latency and high IOPSs requirements. For details, please refer to: http://docs.aws.amazon.com/efs/latest/ug/performance.html

## Use Cloudstor

After creating a swarm on Docker for AWS and connecting to any manager using SSH, verify that Cloudstor is already installed and configured for the stack/resource group:

```bash
$ docker plugin ls
ID                  NAME                        DESCRIPTION                       ENABLED
f416c95c0dcc        cloudstor:aws               cloud storage plugin for Docker   true
```

The following examples show how to create swarm services that require data persistence using the --mount flag and specifying Cloudstor as the driver.

### Share the same volume between tasks (EFS support required and enabled):

```bash
docker service create --replicas 5 --name ping1 \
    --mount type=volume,volume-driver=cloudstor:aws,source=sharedvol1,destination=/shareddata \
    alpine ping docker.com
```

The above is supported only in those regions where EFS is supported and EFS support is enabled during deployment of the Cloud Formation template. All replicas/tasks of the service `ping1` share the same persistent volume `sharedvol1` mounted at `/shareddata` path within the container. Docker Swarm takes care of interacting with the Cloudstor plugin to make sure EFS is mounted on all nodes in the swarm where tasks of the service are scheduled. Each task needs to make sure they don't write concurrently on the same file at the same time and cause corruptions since the volume is shared.

In the above example, you can make sure that the volume is indeed shared by logging into one of the containers in one swarm node, writing to a file under `/shareddata/` and reading the file under `/shareddata/` from another container (in a different node).

If high IO throughput is desired, set the perfmode parameter as volume-opt to maxio:

```bash
{% raw %}
docker service create --replicas 5 --name ping3 \
    --mount type=volume,volume-driver=docker4x/cloudstor:aws-v{{ edition_version }},source={{.Service.Name}}-{{.Task.Slot}}-vol5,destination=/mydata,volume-opt=perfmode=maxio \
    alpine ping docker.com
{% endraw %}
```

### Use a unique volume per task (using EFS with EFS support present and enabled):

```bash
{% raw %}
docker service create --replicas 5 --name ping2 \
    --mount type=volume,volume-driver=cloudstor:aws,source={{.Service.Name}}-{{.Task.Slot}}-vol,destination=/mydata \
    alpine ping docker.com
{% endraw %}
```

Here the templatized notation is used to indicate to Docker Swarm that a unique EFS backed volume be created and mounted for each replica/task of the service `ping2`. After initial creation of the volumes corresponding to the tasks they are attached to (in the nodes the tasks are scheduled in), if a task is rescheduled on a different node, Docker Swarm will interact with the Cloudstor plugin to create and mount the volume corresponding to the task on the node the task got scheduled on. Since data on EFS is available to all Swarm nodes and can be quickly mounted and accessed, the rescheduling process for tasks using EFS backed volumes typically takes a few seconds.

It is highly recommended that you use the `.Task.Slot` template to make sure task N always gets access to vol N no matter which node it is executing on/scheduled to.

In the above example, each task has it's own volume mounted at `/mydata/` and the files under there are unique to the task mounting the volume.

### Use a unique volume per task (using EBS):

```bash
{% raw %}
docker service create --replicas 5 --name ping3 \
     --mount type=volume,volume-driver=cloudstor:aws,source={{.Service.Name}}-{{.Task.Slot}}-vol,destination=/mydata,volume-opt=backing=local,volume-opt=size=25,volume-opt=ebstype=gp2 \
     alpine ping docker.com
{% endraw %}
```

Here the templatized notation is used to indicate to Docker Swarm that a unique volume be created and mounted for each replica/task of the service `ping3`. Creation of EBS volumes typically takes a few minutes. The volume option `backing=local` indicates Cloudstor should use EBS volumes to back the data. The volume option `size` is required and indicates the size of the EBS volumes to create in GB. The volume option `ebstype` is optional and indicates the type of the EBS volumes to create for example gp2, st1, sc1. For further details about EBS volume types, please see http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/EBSVolumeTypes.html.

It is highly recommended that you use the `.Task.Slot` template to make sure task N always gets access to vol N no matter which node it is executing on/scheduled to.

In the above example, each task has it's own volume mounted at `/mydata/` and the files under there are unique to the task mounting the volume.

### List or remove volumes created by Cloudstor

You can use `docker volume ls` on any node to enumerate all volumes created by Cloudstor across the swarm. You can use `docker volume rm [volume name]` to remove a cloudstor volume from any node. Please be aware that if you remove a volume from one node it may still be under active usage in another node and those tasks in the other node will lose access to their data.
