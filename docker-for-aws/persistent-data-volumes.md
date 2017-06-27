---
description: Persistent data volumes
keywords: aws persistent data volumes
title: Docker for AWS persistent data volumes
---

{% include d4a_buttons.md %}

## What is Cloudstor?

Cloudstor is a modern volume plugin built by Docker. It comes pre-installed and pre-configured in Docker Swarms deployed through Docker for AWS. Docker Swarm mode tasks as well as regular Docker containers can use a volume created with Cloudstor to mount a persistent data volume. In Docker for AWS, Cloudstor has two `backing` options: `relocatable` (that uses EBS) and `shared` (that uses EFS) described below. Using the Docker swarm CLI (to create a service along with the persistent volumes for the tasks), it is possible to create the following:
1. Unique `relocatable` Cloudstor volumes mounted by each task in a swarm mode service.
2. Global `shared` Cloudstor volumes mounted by all tasks in a swarm mode service.
3. Unique `shared` Cloudstor volumes mounted by each task in a swarm mode service.
Example for the above are detailed down below.

## Relocatable Cloudstor volumes

 Workloads running in a Docker Service that require access to low latency/high IOPs persistent storage (e.g. a database engine) can use a `relocatable` Cloudstor volume backed by EBS. The type of EBS volume (e.g. gp2, io1, st1, sc1) that the workload requires  can be specified during volume creation. Each `relocatable` Cloudstor volume is backed by a single EBS volume. If a swarm task using a `relocatable` Cloudstor volume gets rescheduled to another node within the same availability zone (as the node where the task was running on), Cloudstor takes care of detaching and re-attaching the backing EBS volume to the target node. If the swarm task gets rescheduled to a node in a different availability zone (from the node where the task was originally running), Cloudstor will transfer the contents of the backing EBS volume using a snapshot to the destination availability zone as well as clean up the EBS volume in the original Availability Zone. To minimize the time necessary to create the snapshot to transfer data across availability zones, Cloudstor periodically takes snapshots of the EBS volumes to ensure there is never a large amount of diff that will need to get transferred as part of the snapshot necessary during a cross availability zone task reschedule. Typically the snapshot based transfer process across availability zones takes about 2 to 5 minutes unless the work load is write heavy. For extremely write-heavy workloads generating several GBs of fresh/new data within minutes, the transfer may take more than 5 minutes. The time required to snapshot and transfer increases sharply beyond 10 minutes if more than 20 GB of diff data has been generated since the last snapshot interval. Note that a swarm task is not started until the volume it mounts becomes available.
 
 Sharing/mounting the same Cloudstor volume backed by EBS among multiple tasks is not a supported scenario and will lead to data loss. If you need a Cloudstor volume to share data between tasks please read below for EFS backed `shared` volume options. `relocatable` Cloudstor backed by EBS is supported on all AWS regions that support EBS. The default `backing` option is `relocatable` if EFS support is not selected during setup/installation or if EFS is not supported in a region.

## Shared Cloudstor volumes 

Workloads running in multiple swarm tasks that need to share data in persistent storage volumes can use a `shared` Cloudstor volume backed by EFS. Such a volume and it's contents can be mounted by multiple tasks on multiple swarm nodes simultaneously since EFS makes the data available to all Swarm nodes over NFS. When swarm tasks using a `shared` Cloudstor volume gets rescheduled from one node to another within the same or across different availability zones, the persistent data backed by EFS volumes is always available. `shared` Cloudstor volumes will only work in those AWS regions where EFS is supported (e.g. Ohio, Oregon, N.Virginia, Sydney). If EFS Support is selected during setup/installation, the default "backing" option for Cloudstor volumes is set to Shared so that EFS is used by default. Please note that `shared` Cloudstor volumes backed by EFS (or even EFS MaxIO) may not be ideal for workloads that have very low latency and high IOPSs requirements. For performance details of EFS backed `shared` Cloudstor volumes, please refer to: http://docs.aws.amazon.com/efs/latest/ug/performance.html

## Use Cloudstor

After creating a swarm on Docker for AWS and connecting to any manager using SSH, verify that Cloudstor is already installed and configured for the stack/resource group:

```bash
$ docker plugin ls
ID                  NAME                        DESCRIPTION                       ENABLED
f416c95c0dcc        cloudstor:aws               cloud storage plugin for Docker   true
```

The following examples show how to create swarm services that require data persistence using the --mount flag and specifying Cloudstor as the driver.

### Share the same volume between tasks (EFS support required and enabled):

In those regions where EFS is supported and EFS support is enabled during deployment of the Cloud Formation template, `shared` Cloudstor volumes can be created to share access to persistent data across all tasks in a swarm service running in multiple nodes. Example usage:

```bash
docker service create --replicas 5 --name ping1 \
    --mount type=volume,volume-driver=cloudstor:aws,source=sharedvol1,destination=/shareddata \
    alpine ping docker.com
```

All replicas/tasks of the service `ping1` above share the same persistent volume `sharedvol1` mounted at `/shareddata` path within the container. Docker Swarm takes care of interacting with the Cloudstor plugin to make sure EFS is mounted on all nodes in the swarm where tasks of the service are scheduled. Each task needs to make sure they don't write concurrently on the same file at the same time as that may lead to data corruption since the volume is shared.

In the above example, you can make sure that the volume is indeed shared by logging into one of the containers in one Swarm node, writing to a file under `/shareddata/` and reading the file under `/shareddata/` from another container (in a different node).

The only option available for EFS currently is `perfmode`. If high IO throughput is desired, please set the `perfmode` volume-opt to `maxio`. Example:

```bash
{% raw %}
docker service create --replicas 5 --name ping3 \
    --mount type=volume,volume-driver=docker4x/cloudstor:aws-v{{ edition_version }},source={{.Service.Name}}-{{.Task.Slot}}-vol5,destination=/mydata,volume-opt=perfmode=maxio \
    alpine ping docker.com
{% endraw %}
```

It is also possible to create `shared` Cloudstor volumes using the `docker volume create` CLI. For example:

```bash
{% raw %}
docker volume create -d "cloudstor:aws" --opt backing=shared mysharedvol1
{% endraw %}
```

### Use a unique volume per task (using EBS):

A unique `relocatable` Cloudstor volume backed by a specified type of EBS can be created and mounted for each task in a swarm service using a templatized notation with the `docker service create` CLI. Creation of new EBS volumes typically takes a few minutes.  The volume options (besides `backing=local` indicating `relocatable` EBS backed volumes) are:
1. `size` : Required parameter that indicates the size of the EBS volumes to create in GB. 
2. `ebstype` : Optional parameter that indicates the type of the EBS volumes to create, for example gp2, io1, st1, sc1. The default `ebstype` is Standard/Magnetic. For further details about EBS volume types, please see http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/EBSVolumeTypes.html. 
3. `iops` : Required if `ebstype` specified is `io1` i.e. provisioned IOPs. Needs to be in the appropriate range as required by EBS. 

Example usage:

```bash
{% raw %}
docker service create --replicas 5 --name ping3 \
     --mount type=volume,volume-driver=cloudstor:aws,source={{.Service.Name}}-{{.Task.Slot}}-vol,destination=/mydata,volume-opt=backing=local,volume-opt=size=25,volume-opt=ebstype=gp2 \
     alpine ping docker.com
{% endraw %}
```
The above indicates to Docker Swarm that a distinct Cloudstor volume backed by 25 GB EBS volumes of type gp2 be created and mounted for each replica/task of the service `ping3`. Each task has it's own volume mounted at `/mydata/` and the files under there are unique to the task mounting the volume

It is highly recommended that you use the `.Task.Slot` template to make sure task N always gets access to vol N no matter which node it is executing on/scheduled to. It is recommended that the total number of EBS volumes in the swarm be kept below 12 * (minimum number of nodes that are expected to be present at any time) to ensure EC2 can properly attach EBS volumes to a node when another node fails. Thus use EBS volumes only for those workloads where low latency and high IOPs is absolutely necessary.

It is also possible to create EBS backed volumes using the `docker volume create` CLI. For example:

```bash
{% raw %}
docker volume create -d "cloudstor:aws" --opt ebstype=io1 --opt size=25 --opt iops=1000 --opt backing=local mylocalvol1
{% endraw %}
```

Sharing the same `relocatable` Cloudstor volume across multiple tasks of a service or across multiple independent containers is not supported when `backing=local` is specified. Trying to do so in containers across multiple nodes will result in IO errors.

### Use a unique volume per task (using EFS with EFS support present and enabled):

It is possible to use the templatized notation to indicate to Docker Swarm that a unique EFS backed volume be created and mounted for each replica/task of a service. This is a useful option if you already have too many EBS volumes or want to reduce the amount of time it takes to transfer volume data across availability zones. Example:

```bash
{% raw %}
docker service create --replicas 5 --name ping2 \
    --mount type=volume,volume-driver=cloudstor:aws,source={{.Service.Name}}-{{.Task.Slot}}-vol,destination=/mydata \
    alpine ping docker.com
{% endraw %}
```

In the above example, each task has it's own volume mounted at `/mydata/` and the files under there are unique to the task mounting the volume. 

When a task with only `shared` EFS volumes mounted is rescheduled on a different node, Docker Swarm will interact with the Cloudstor plugin to create and mount the volume corresponding to the task on the node the task got scheduled on. Since data on EFS is available to all Swarm nodes and can be quickly mounted and accessed, the rescheduling process for tasks using EFS backed volumes typically takes a few seconds compared to several minutes in case of EBS

It is highly recommended that you use the `.Task.Slot` template to make sure task N always gets access to vol N no matter which node it is executing on/scheduled to.

### List or remove volumes created by Cloudstor

You can use `docker volume ls` on any node to enumerate all volumes created by Cloudstor across the swarm. You can use `docker volume rm [volume name]` to remove a cloudstor volume from any node. Please be aware that if you remove a volume from one node please make sure it is not under active usage in another node as those tasks/containers in another node will lose access to their data.

Before deleting a Docker4AWS stack through CloudFormation, it is recommended that you remove all `relocatable` Cloudstor volumes using `docker volume rm` from within the stack. EBS volumes corresponding to `relocatable` Cloudstor volumes will not be deleted as part of the CloudFormation stack deletion. To list any `relocatable` Cloudstor volumes and delete them after a Docker4AWS stack (in which the volumes where created) has been deleted, go to the AWS portal or CLI and set a filter with tag key set to `StackID` and the tag value set to the md5 hash of the CloudFormation Stack ID (typical format: arn:aws:cloudformation:us-west-2:ID:stack/swarmname/GUID)
