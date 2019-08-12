---
title: Disaster recovery
description: Learn disaster recovery procedures for Docker Enterprise
keywords: enterprise, recovery, disaster recovery, dtr, ucp, swarm
redirect_from:
  - /ee/dtr/admin/disaster-recovery/
  - /ee/dtr/admin/disaster-recovery/repair-a-single-replica/
  - /ee/dtr/admin/disaster-recovery/repair-a-cluster/
---

Disaster recovery procedures should be performed in the following order:

1. [Docker Swarm](#swarm-disaster-recovery).
2. [Universal Control Plane (UCP)](#ucp-disaster-recovery).
3. [Docker Trusted Registry (DTR)](#dtr-disaster-recovery).

## Swarm disaster recovery

### Recover from losing the quorum

Swarm is resilient to failures and the swarm can recover from any number
of temporary node failures (machine reboots or crash with restart) or other
transient errors. However, a swarm cannot automatically recover if it loses a
quorum. Tasks on existing worker nodes continue to run, but administrative
tasks are not possible, including scaling or updating services and joining or
removing nodes from the swarm. The best way to recover is to bring the missing
manager nodes back online. If that is not possible, continue reading for some
options for recovering your swarm.

In a swarm of `N` managers, a quorum (a majority) of manager nodes must always
be available. For example, in a swarm with 5 managers, a minimum of 3 must be
operational and in communication with each other. In other words, the swarm can
tolerate up to `(N-1)/2` permanent failures beyond which requests involving
swarm management cannot be processed. These types of failures include data
corruption or hardware failures.

If you lose the quorum of managers, you cannot administer the swarm. If you have
lost the quorum and you attempt to perform any management operation on the swarm,
an error occurs:

```none
Error response from daemon: rpc error: code = 4 desc = context deadline exceeded
```

The best way to recover from losing the quorum is to bring the failed nodes back
online. If you can't do that, the only way to recover from this state is to use
the `--force-new-cluster` action from a manager node. This removes all managers
except the manager the command was run from. The quorum is achieved because
there is now only one manager. Promote nodes to be managers until you have the
desired number of managers.

```bash
# From the node to recover
$ docker swarm init --force-new-cluster --advertise-addr node01:2377
```

When you run the `docker swarm init` command with the `--force-new-cluster`
flag, the Docker Engine where you run the command becomes the manager node of a
single-node swarm which is capable of managing and running services. The manager
has all the previous information about services and tasks, worker nodes are
still part of the swarm, and services are still running. You need to add or
re-add  manager nodes to achieve your previous task distribution and ensure that
you have enough managers to maintain high availability and prevent losing the
quorum.

### Force the swarm to rebalance

Generally, you do not need to force the swarm to rebalance its tasks. When you
add a new node to a swarm, or a node reconnects to the swarm after a
period of unavailability, the swarm does not automatically give a workload to
the idle node. This is a design decision. If the swarm periodically shifted tasks
to different nodes for the sake of balance, the clients using those tasks would
be disrupted. The goal is to avoid disrupting running services for the sake of
balance across the swarm. When new tasks start, or when a node with running
tasks becomes unavailable, those tasks are given to less busy nodes. The goal
is eventual balance, with minimal disruption to the end user.

In Docker 1.13 and higher, you can use the `--force` or `-f` flag with the
`docker service update` command to force the service to redistribute its tasks
across the available worker nodes. This causes the service tasks to restart.
Client applications may be disrupted. If you have configured it, your service
uses a [rolling update](/engine/swarm/swarm-tutorial/rolling-update/).

If you use an earlier version and you want to achieve an even balance of load
across workers and don't mind disrupting running tasks, you can force your swarm
to re-balance by temporarily scaling the service upward. Use
`docker service inspect --pretty <servicename>` to see the configured scale
of a service. When you use `docker service scale`, the nodes with the lowest
number of tasks are targeted to receive the new workloads. There may be multiple
under-loaded nodes in your swarm. You may need to scale the service up by modest
increments a few times to achieve the balance you want across all the nodes.

When the load is balanced to your satisfaction, you can scale the service back
down to the original scale. You can use `docker service ps` to assess the current
balance of your service across nodes.

See also
[`docker service scale`](/engine/reference/commandline/service_scale/) and
[`docker service ps`](/engine/reference/commandline/service_ps/).

## UCP disaster recovery

In the event half or more manager nodes are lost and cannot be recovered
to a healthy state, the system is considered to have lost quorum and can only be
restored through the following disaster recovery procedure. 

### Recover a UCP cluster from an existing backup

1. If UCP is still installed on the swarm, uninstall UCP using the `uninstall-ucp` command. 
   > **Note**: If the restore is happening on new machines, skip this step.
2. Perform a [restore from an existing backup](/ee/admin/restore/) on any node. If there is an 
existing swarm, the restore operation must be performed on a manager node. If no swarm exists, 
the restore operation will create one. 

### Recover a UCP cluster without an existing backup (not recommended)
If your cluster has lost quorum, you can still perform a backup of one of the remaining nodes.

> **Important**: Performing a backup after losing quorum is not guaranteed to succeed with
no loss of running services or configuration data. To properly protect against
manager failures, the system should be configured for
[high availability](/ee/ucp/admin/configure/join-nodes/), and backups should be performed regularly
in order to have complete backup data.

1. On one of the remaining manager nodes, perform `docker swarm init --force-new-cluster`. You might also need to specify an
   `--advertise-addr` parameter, which is equivalent to the `--host-address`
   parameter of the `docker/ucp install` operation. This instantiates a new
   single-manager swarm by recovering as much state as possible from the
   existing manager. This is a disruptive operation and existing tasks might be
   either terminated or suspended.
2. [Create a backup](/ee/admin/backup/) of the remaining manager node.
3. If UCP is still installed on the swarm, uninstall UCP using the
   `uninstall-ucp` command.
4. Perform a [restore](/ee/admin/restore/) on the recovered swarm manager node.
5. Log in to UCP and browse to the nodes page, or use the CLI `docker node ls`
   command.
6. If any nodes are listed as `down`, you'll have to manually [remove these
   nodes](/ee/ucp/admin/configure/scale-your-cluster/) from the swarm and then re-join
   them using a `docker swarm join` operation with the swarm's new join-token.  
7. [Create a backup](/ee/admin/backup/) of the restored cluster.

###  Recreate objects within Orchestrators that Docker Enterprise supports

Kubernetes currently backs up the declarative state of Kube objects in etcd. However, for Swarm, there is no way to take the state and export it to a declarative format, since the objects that are embedded within the Swarm raft logs are not easily transferable to other nodes or clusters. 

For disaster recovery, to recreate swarm related workloads requires having the original scripts used for deployment. Alternatively, you can recreate workloads by manually recreating output from `docker inspect` commands.

## DTR disaster recovery

Docker Trusted Registry is a clustered application. You can join multiple
replicas for high availability. For a DTR cluster to be healthy, a majority of its replicas (n/2 + 1) need to
be healthy and be able to communicate with the other replicas. This is also
known as maintaining quorum.

This means that there are three failure scenarios possible.

### Replica is unhealthy but cluster maintains quorum

One or more replicas are unhealthy, but the overall majority (n/2 + 1) is still
healthy and able to communicate with one another.

![Failure scenario 1](/ee/dtr/images/dr-overview-1.svg)

In this example the DTR cluster has five replicas but one of the nodes stopped
working, and the other has problems with the DTR overlay network.

Even though these two replicas are unhealthy the DTR cluster has a majority
of replicas still working, which means that the cluster is healthy.

In this case you should repair the unhealthy replicas, or remove them from
the cluster and join new ones.

#### Repair a single replica

When one or more DTR replicas are unhealthy but the overall majority
(n/2 + 1) is healthy and able to communicate with one another, your DTR
cluster is still functional and healthy.

![Cluster with two nodes unhealthy](/ee/dtr/images/repair-replica-1.svg)

Given that the DTR cluster is healthy, there's no need to execute any disaster
recovery procedures like restoring from a backup.

Instead, you should:

1. Remove the unhealthy replicas from the DTR cluster.
2. Join new replicas to make DTR highly available.

Since a DTR cluster requires a majority of replicas to be healthy at all times,
the order of these operations is important. If you join more replicas before
removing the ones that are unhealthy, your DTR cluster might become unhealthy.

##### Split-brain scenario

To understand why you should remove unhealthy replicas before joining new ones,
imagine you have a five-replica DTR deployment, and something goes wrong with
the overlay network connection the replicas, causing them to be separated in
two groups.

![Cluster with network problem](/ee/dtr/images/repair-replica-2.svg)

Because the cluster originally had five replicas, it can work as long as
three replicas are still healthy and able to communicate (5 / 2 + 1 = 3).
Even though the network separated the replicas in two groups, DTR is still
healthy.

If at this point you join a new replica instead of fixing the network problem
or removing the two replicas that got isolated from the rest, it's possible
that the new replica ends up in the side of the network partition that has
less replicas.

![cluster with split brain](/ee/dtr/images/repair-replica-3.svg)

When this happens, both groups now have the minimum amount of replicas needed
to establish a cluster. This is also known as a split-brain scenario, because
both groups can now accept writes and their histories start diverging, making
the two groups effectively two different clusters.

##### Remove replicas

To remove unhealthy replicas, you'll first have to find the replica ID
of one of the replicas you want to keep, and the replica IDs of the unhealthy
replicas you want to remove.

You can find the list of replicas by navigating to **Shared Resources > Stacks** or **Swarm > Volumes** (when using [swarm mode](/engine/swarm/)) on the UCP web interface, or by using the UCP
client bundle to run:

{% raw %}
```bash
$ docker ps --format "{{.Names}}" | grep dtr

# The list of DTR containers with <node>/<component>-<replicaID>, e.g.
# node-1/dtr-api-a1640e1c15b6
```
{% endraw %}

Another way to determine the replica ID is to SSH into a DTR node and run the following:

{% raw %}
```bash
$ REPLICA_ID=$(docker inspect -f '{{.Name}}' $(docker ps -q -f name=dtr-rethink) | cut -f 3 -d '-')
&& echo $REPLICA_ID
```
{% endraw %}

Then use the UCP client bundle to remove the unhealthy replicas:

```bash
$ docker container run \
  --rm \
  --interactive \
  {{ page.dtr_org }}/{{ page.dtr_repo }}:{{ page.dtr_version }} remove \
  --existing-replica-id <healthy-replica-id> \
  --replica-ids <unhealthy-replica-id> \
  --ucp-insecure-tls \
  --ucp-url <ucp-url> \
  --ucp-username <user> \
  --ucp-password <password>
```

You can remove more than one replica at the same time, by specifying multiple
IDs with a comma.

![Healthy cluster](/ee/dtr/images/repair-replica-4.svg)

##### Join replicas

Once you've removed the unhealthy nodes from the cluster, you should join new
ones to make sure your cluster is highly available.

Use your UCP client bundle to run the following command which prompts you for
the necessary parameters:

```bash
$ docker container run \
  --rm \
  --interactive \
  {{ page.dtr_org }}/{{ page.dtr_repo }}:{{ page.dtr_version }} join \
  --ucp-node <ucp-node-name> \
  --ucp-insecure-tls
```

[Learn more about high availability](/ee/dtr/admin/configure/set-up-high-availability/).

### The majority of replicas are unhealthy

If a majority of replicas are unhealthy, making the cluster lose quorum, but at
least one replica is still healthy, or at least the data volumes for DTR are
accessible from that replica, you can repair the cluster without having to restore from
a backup. This minimizes the amount of data loss. The following image provides an example of this scenario.

![Failure scenario 2](/ee/dtr/images/dr-overview-2.svg)
 
#### Repair a cluster

For a DTR cluster to be healthy, a majority of its replicas (n/2 + 1) need to
be healthy and be able to communicate with the other replicas. This is known
as maintaining quorum.

In a scenario where quorum is lost, but at least one replica is still
accessible, you can use that replica to repair the cluster. That replica doesn't
need to be completely healthy. The cluster can still be repaired as the DTR
data volumes are persisted and accessible.

![Unhealthy cluster](/ee/dtr/images/repair-cluster-1.svg)

Repairing the cluster from an existing replica minimizes the amount of data lost.
If this procedure doesn't work, you'll have to
[restore from an existing backup](/ee/admin/restore/).

##### Diagnose an unhealthy cluster

When a majority of replicas are unhealthy, causing the overall DTR cluster to
become unhealthy, operations like `docker login`, `docker pull`, and `docker push`
present `internal server error`.

Accessing the `/_ping` endpoint of any replica also returns the same error.
It's also possible that the DTR web UI is partially or fully unresponsive.

##### Perform an emergency repair

Use the `docker/dtr emergency-repair` command to try to repair an unhealthy
DTR cluster, from an existing replica.

This command checks the data volumes for the DTR

This command checks the data volumes for the DTR replica are uncorrupted,
redeploys all internal DTR components and reconfigured them to use the existing
volumes.

It also reconfigures DTR removing all other nodes from the cluster, leaving DTR
as a single-replica cluster with the replica you chose.

Start by finding the ID of the DTR replica that you want to repair from.
You can find the list of replicas by navigating to **Shared Resources > Stacks** or **Swarm > Volumes** (when using [swarm mode](/engine/swarm/)) on the UCP web interface, or by using
a UCP client bundle to run:

{% raw %}
```bash
$ docker ps --format "{{.Names}}" | grep dtr

# The list of DTR containers with <node>/<component>-<replicaID>, e.g.
# node-1/dtr-api-a1640e1c15b6
```
{% endraw %}

Another way to determine the replica ID is to SSH into a DTR node and run the following:

{% raw %}
```bash
$ REPLICA_ID=$(docker inspect -f '{{.Name}}' $(docker ps -q -f name=dtr-rethink) | cut -f 3 -d '-')
&& echo $REPLICA_ID
```
{% endraw %}

Then, use your UCP client bundle to run the emergency repair command:

```bash
$ docker container run \
  --rm \
  --interactive \
  {{ page.dtr_org }}/{{ page.dtr_repo }}:{{ page.dtr_version }} emergency-repair \
  --ucp-insecure-tls \
  --existing-replica-id <replica-id>
```

If the emergency repair procedure is successful, your DTR cluster now has a
single replica. You should now
[join more replicas for high availability](/ee/dtr/admin/configure/set-up-high-availability/).

![Healthy cluster](/ee/dtr/images/repair-cluster-2.svg)

If the emergency repair command fails, try running it again using a different
replica ID. As a last resort, you can restore your cluster from an existing
backup.

### All replicas are unhealthy

This is a total disaster scenario where all DTR replicas are lost, causing
the data volumes for all DTR replicas to get corrupted or be lost.

![Failure scenario 3](/ee/dtr/images/dr-overview-3.svg)

In a disaster scenario like this, you'll have to restore DTR from an existing
backup. Restoring from a backup should be only used as a last resort, since
doing an emergency repair might prevent some data loss.

[Create a backup](/ee/admin/backup/).

## Where to go next

- [Create a backup](/ee/admin/backup/)
- [Set up high availability](/ee/ucp/admin/configure/join-nodes/)
