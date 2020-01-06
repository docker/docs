---
title: Upgrade to UCP 3.2
description: Learn how to upgrade Docker Universal Control Plane with minimal impact to your users.
keywords: UCP, upgrade, update
---

>{% include enterprise_label_shortform.md %}

This page helps you upgrade Docker Universal Control Plane (UCP) to version {{
page.ucp_version }}.

## Plan the upgrade

Before upgrading to a new version of UCP, check the
[release notes](/ee/ucp/release-notes/).
There you'll find information about new features, breaking changes, and
other relevant information for upgrading to a particular version.

As part of the upgrade process, you'll upgrade the Docker Engine - Enterprise
installed on each node of the cluster to version 19.03 or higher.
You should plan for the upgrade to take place outside of business hours,
to ensure there's minimal impact to your users.

Also, don't make changes to UCP configurations while you're upgrading it.
This can lead to misconfigurations that are difficult to troubleshoot.

### Environment checklist
Complete the following checks:

#### Systems
- Confirm time sync across all nodes (and check time daemon logs for any large time drifting)
- Check system requirements `PROD=4` `vCPU/16GB` for UCP managers and DTR replicas
- Review the full UCP/DTR/Engine port requirements
- Ensure that your cluster nodes meet the minimum requirements
- Before performing any upgrade, ensure that you meet all minimum requirements listed in [UCP System requirements](/ee/ucp/admin/install/system-requirements/), including port openings (UCP 3.x added more required ports for Kubernetes), memory, and disk space. For example, manager nodes must have at least 8GB of memory.

> Note
> 
> If you are upgrading a cluster to UCP 3.0.2 or higher on Microsoft
> Azure, please ensure that all of the Azure [prerequisites](install-on-azure.md/#azure-prerequisites)
> are met.

#### Storage
- Check `/var/` storage allocation and increase if it is over 70% usage.
- In addition, check all nodes’ local file systems for any disk storage issues (and DTR back-end storage, for example, NFS).
- If not using Overlay2 storage drivers please take this opportunity to do so, you will find stability there. Note that the transition from Device mapper to Overlay2 is a destructive rebuild.

#### Operating system
- If cluster nodes OS branch is older (Ubuntu 14.x, RHEL 7.3, etc), consider patching all relevant packages to the most recent (including kernel).
- Rolling restart of each node before upgrade (to confirm in-memory settings are the same as startup-scripts).
- Run `check-config.sh` on each cluster node (after rolling restart) for any kernel compatibility issues. Latest version of the script can be found here: https://github.com/moby/moby/blob/master/contrib/check-config.sh

#### Procedural
- Perform Swarm, UCP and DTR backups before upgrading
- Gather Compose file/service/stack files
- Generate a UCP Support dump (for point in time) before upgrading
- Preinstall Engine/UCP/DTR images. If your cluster is offline (with no connection to the internet), then Docker provides tarballs containing all of the required container images [here](/ee/ucp/admin/install/upgrade-offline/). If your cluster is
online, you can pull the required container images onto your nodes with the following command:

```bash
$ docker run --rm {{ page.ucp_org }}/{{ page.ucp_repo }}:{{ page.ucp_version }} images --list | xargs -L 1 docker pull
```
- Load troubleshooting packages (netshoot, etc)
- Best order for upgrades: Engine, UCP, and then DTR. Note that the scope of this topic is limited to upgrade instructions for UCP.

#### Upgrade strategy

For each worker node that requires an upgrade, you can upgrade that node in place or you can replace the node
with a new worker node. The type of upgrade you perform depends on what is needed for each node:

- [Automated, in-place cluster upgrade](#automated-in-place-cluster-upgrade): Performed on any manager node. Automatically upgrades the entire cluster.
- Manual cluster upgrade: Performed using the CLI. Automatically upgrades manager nodes and allows you to control the upgrade order of worker nodes. This type of upgrade is more advanced than the automated, in-place cluster upgrade.
  - [Upgrade existing nodes in place](#phased-in-place-cluster-upgrade): Performed using the CLI. Automatically upgrades manager nodes and allows you to control the order of worker node upgrades.
  - [Replace all worker nodes using blue-green deployment](#replace-existing-worker-nodes-using-blue-green-deployment): Performed using the CLI. This type of upgrade allows you to stand up a new cluster in parallel to the current code and cut over when complete. This type of upgrade allows you to join new worker nodes, schedule workloads to run on new nodes, pause, drain, and remove old worker nodes in batches of multiple nodes rather than one at a time, and shut down servers to remove worker nodes. This type of upgrade is the most advanced.

## Back up your cluster

Before starting an upgrade, make sure that your cluster is healthy. If a problem
occurs, this makes it easier to find and troubleshoot it.

[Create a backup](/ee/admin/backup/back-up-ucp/) of your cluster.
This allows you to recover if something goes wrong during the upgrade process.

> Note
> 
> The backup archive is version-specific, so you can't use it during the
> upgrade process. For example, if you create a backup archive for a UCP 2.2
> cluster, you can't use the archive file after you upgrade to UCP 3.0.

## Upgrade Docker Engine

For each node that is part of your cluster, upgrade the Docker Engine
installed on that node to Docker Engine version 19.03 or higher. Be sure
to install the Docker Enterprise Edition.

Starting with the manager nodes, and then worker nodes:

1. Log into the node using ssh.
2. Upgrade the Docker Engine to version 18.09.0 or higher. See [Upgrade Docker EE](/ee/upgrade/).
3. Make sure the node is healthy.

> Note
> 
> In your browser, navigate to **Nodes** in the UCP web interface, and check that the node is healthy and is part of the cluster.

## Upgrade UCP
When upgrading UCP to version {{ page.ucp_version }}, you can choose from
different upgrade workflows:

> Note
>  
> In all upgrade workflows, manager nodes are automatically upgraded in place. You cannot control the order
> of manager node upgrades.
{: .important}

- [Automated, in-place cluster upgrade](#automated-in-place-cluster-upgrade): Performed on any
manager node. Automatically upgrades the entire cluster.
- Manual cluster upgrade: Performed using the CLI. Automatically upgrades manager
nodes and allows you to control the upgrade order of worker nodes. This type of upgrade is more
advanced than the automated, in-place cluster upgrade.
  - [Upgrade existing nodes in place](#phased-in-place-cluster-upgrade): Performed using the CLI. Automatically upgrades manager nodes and allows you to control the order of worker node upgrades.
  - [Replace all worker nodes using blue-green deployment](#replace-existing-worker-nodes-using-blue-green-deployment): Performed using the CLI. This type of upgrade allows you to stand up a new cluster in parallel to the current code and cut over when complete. This type of upgrade allows you to join new worker nodes, schedule workloads to run on new nodes, pause, drain, and remove old worker nodes in batches of multiple nodes rather than one at a time, and shut down servers to remove worker nodes. This type of upgrade is the most advanced.

### Use the CLI to perform an upgrade

There are two different ways to upgrade a UCP Cluster via the CLI. The first is
an automated process; this approach will update all UCP components on all nodes
within the UCP Cluster. The upgrade process is done node by node, but once the
user has initiated an upgrade it will work its way through the entire cluster.

The second UCP upgrade method is a phased approach, once an upgrade has been
initiated this method will upgrade all UCP components on a single UCP worker
nodes, giving the user more control to migrate workloads and control traffic
when upgrading the cluster.

### Automated in-place cluster upgrade

This is the traditional approach to upgrading UCP and is often used when the
order in which UCP worker nodes is upgraded is NOT important.

To upgrade UCP, ensure all Docker engines have been upgraded to the
corresponding new version. Then a user should SSH to one UCP manager node and run
the following command. The upgrade command should not be run on a workstation
with a client bundle.

```
$ docker container run --rm -it \
  --name ucp \
  --volume /var/run/docker.sock:/var/run/docker.sock \
  {{ page.ucp_org }}/{{ page.ucp_repo }}:{{ page.ucp_version }} \
  upgrade \
  --interactive
```

The upgrade command will print messages regarding the progress of the upgrade as
it automatically upgrades UCP on all nodes in the cluster.

### Phased in-place cluster upgrade

The phased approach of upgrading UCP, introduced in UCP 3.2, allows granular
control of the UCP upgrade process. A user can temporarily run UCP worker nodes
with different versions of the Docker Engine and UCP. This workflow is useful
when a user wants to manually control how workloads and traffic are migrated
around a cluster during an upgrade. This process can also be used if a user
wants to add additional worker node capacity during an upgrade to handle
failover. Worker nodes can be added to a partially upgraded UCP Cluster,
workloads migrated across, and previous worker nodes then taken offline and
upgraded.

To start a phased upgrade of UCP, first all manager nodes will need to be
upgraded to the new UCP version. To tell UCP to upgrade the manager nodes but
not upgrade any worker nodes, pass `--manual-worker-upgrade` into the upgrade
command.

To upgrade UCP, ensure the Docker engine on all UCP manager nodes have been
upgraded to the corresponding new version. SSH to a UCP manager node and run
the following command. The upgrade command should not be run on a workstation
with a client bundle.

```
$ docker container run --rm -it \
  --name ucp \
  --volume /var/run/docker.sock:/var/run/docker.sock \
  {{ page.ucp_org }}/{{ page.ucp_repo }}:{{ page.ucp_version }} \
  upgrade \
  --manual-worker-upgrade \
  --interactive
```

The `--manual-worker-upgrade` flag will add an upgrade-hold label to all worker
nodes. UCP will be constantly monitor this label, and if that label is removed
UCP will then upgrade the node.

To trigger the upgrade on a worker node, you will have to remove the
label.

```
$ docker node update --label-rm com.docker.ucp.upgrade-hold <node name or id>
```

> Optional
> 
> Joining new worker nodes to the cluster. Once the manager nodes have
been upgraded to a new UCP version, new worker nodes can be added to the
cluster, assuming they are running the corresponding new docker engine
version.

The swarm join token can be found in the UCP UI, or while ssh'd on a UCP
manager node. More information on finding the swarm token can be found
[here](/ee/ucp/admin/configure/join-nodes/join-linux-nodes-to-cluster/).

```
$ docker swarm join --token SWMTKN-<YOUR TOKEN> <manager ip>:2377
```

### Replace existing worker nodes using blue-green deployment

This workflow is used to create a parallel environment for a new deployment, which can greatly reduce downtime, upgrades
worker node engines without disrupting workloads, and allows traffic to be migrated to the new environment with
worker node rollback capability. This type of upgrade creates a parallel environment for reduced downtime and workload disruption.

> Note
> 
> Steps 2 through 6 can be repeated for groups of nodes - you do not have to replace all worker
nodes in the cluster at one time.

1. Upgrade manager nodes

    - The `--manual-worker-upgrade` command automatically upgrades manager nodes first, and then allows you to control
    the upgrade of the UCP components on the worker nodes using node labels.

    ```
    $ docker container run --rm -it \
      --name ucp \
      --volume /var/run/docker.sock:/var/run/docker.sock \
      {{ page.ucp_org }}/{{ page.ucp_repo }}:{{ page.ucp_version }} \
      upgrade \
      --manual-worker-upgrade \
      --interactive
    ```

2. Join new worker nodes

    - New worker nodes have newer engines already installed and have the new UCP version running when they join the cluster.
    On the manager node, run commands similar to the following examples to get the Swarm Join token and add new worker nodes:
    ```
    docker swarm join-token worker
    ```
    - On the node to be joined:
    ```
    docker swarm join --token SWMTKN-<YOUR TOKEN> <manager ip>:2377
    ```
3. Join Enterprise Engine to the cluster
    ```
    docker swarm join --token SWMTKN-<YOUR TOKEN> <manager ip>:2377
    ```
4. Pause all existing worker nodes

    - This ensures that new workloads are not deployed on existing nodes.
      ```
      docker node update --availability pause <node name>
      ```
5. Drain paused nodes for workload migration

    - Redeploy workloads on all existing nodes to new nodes. Because all existing nodes are “paused”, workloads are
    automatically rescheduled onto new nodes.
    ```
    docker node update --availability drain <node name>
    ```
6. Remove drained nodes

    - After each node is fully drained, it can be shut down and removed from the cluster. On each worker node that is
    getting removed from the cluster, run a command similar to the following example :
    ```
    docker swarm leave <node name>
    ```
    - Run a command similar to the following example on the manager node when the old worker comes unresponsive:
    ```
    docker node rm <node name>
    ```
7. Remove old UCP agents

    - After upgrade completion, remove old UCP agents, which includes 390x and Windows agents, that were carried over
    from the previous install by running the following command on the manager node:
    ```
    docker service rm ucp-agent
    docker service rm ucp-agent-win
    docker service rm ucp-agent-s390x
    ```

### Troubleshooting

- Upgrade compatibility
  - The upgrade command automatically checks for multiple `ucp-worker-agents` before
      proceeding with the upgrade. The existence of multiple `ucp-worker-agents` might indicate
      that the cluster still in the middle of a prior manual upgrade and you must resolve the
      conflicting node labels issues before proceeding with the upgrade.
- Upgrade failures
  - For worker nodes, an upgrade failure can be rolled back by changing the node label back
      to the previous target version. Rollback of manager nodes is not supported.
- [Kubernetes errors in node state messages after upgrading UCP](https://success.docker.com/article/how-to-resolve-kubernetes-errors-after-upgrading-ucp)
- The following information applies if you have upgraded to UCP 3.0.0 or newer:
  - After performing a UCP upgrade from 2.2.x to 3.x.x, you might see unhealthy nodes in your UCP
      dashboard with any of the following errors listed:
      ```
      Awaiting healthy status in Kubernetes node inventory
      Kubelet is unhealthy: Kubelet stopped posting node status
      ```
  - Alternatively, you may see other port errors such as the one below in the ucp-controller
      container logs:
      ```
      http: proxy error: dial tcp 10.14.101.141:12388: connect: no route to host
      ```
- UCP 3.x.x requires additional opened ports for Kubernetes use. For ports that are used by the
latest UCP versions and the scope of port use, refer to
[this page](https://docs.docker.com/ee/ucp/admin/install/system-requirements/#ports-used).
  - If you have upgraded from UCP 2.2.x to 3.0.x, verify that the ports 179, 6443, 6444, and 10250 are
    open for Kubernetes traffic.
  - If you have upgraded to UCP 3.1.x, in addition to the ports listed above, also open
    ports 9099 and 12388.

### Recommended upgrade paths

- From UCP 3.0: UCP 3.0 > UCP 3.1 > UCP 3.2
- From UCP 2.2: UCP 2.2 > UCP 3.0 > UCP 3.1 > UCP 3.2

If you’re running a UCP version earlier than 2.1, first upgrade to the latest
2.1 version, then upgrade to 2.2. Use the following rules for your upgrade path
to UCP 2.2:

- From UCP 1.1: UCP 1.1 > UCP 2.1 > UCP 2.2
- From UCP 2.0: UCP 2.0 > UCP 2.1 > UCP 2.2

## Where to go next

- [Upgrade DTR](/ee/dtr/admin/upgrade/)
