---
title: Upgrade to UCP 3.2
description: Learn how to upgrade Docker Universal Control Plane with minimal impact to your users.
keywords: UCP, upgrade, update
---

With UCP version {{ page.ucp_version }}, upgrades provide improved progress information for install and upgrade as well as 
no downtime for user workloads.

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

- Systems:
    
    - Confirm time sync across all nodes (and check time daemon logs for any large time drifting)
    - Check system requirements `PROD=4` `vCPU/16GB` for UCP managers and DTR replicas
    - Review the full UCP/DTR/Engine port requirements 
    - Ensure that your cluster nodes meet the minimum requirements   
    - Before performing any upgrade, ensure that you meet all minimum requirements listed 
    in [UCP System requirements](/ee/ucp/admin/install/system-requirements/), including port openings (UCP 3.x added more 
    required ports for Kubernetes), memory, and disk space. For example, manager nodes must have at least 8GB of memory.
    > **Note**: If you are upgrading a cluster to UCP 3.0.2 or higher on Microsoft
    > Azure then please ensure all of the Azure [prerequisites](install-on-azure.md/#azure-prerequisites) 
    > are met.

- Storage:

    - Check `/var/` storage allocation and increase if it is over 70% usage.
    - In addition, check all nodes’ local filesystems for any disk storage issues (and DTR backend storage, for example, NFS).
    - If not using Overlay2 storage drivers please take this opportunity to do so, you will find stability there. (Note: 
    The transition from Device mapper to Overlay2 is a destructive rebuild.)

- Operating system:

    - If cluster nodes OS branch is older (Ubuntu 14.x, RHEL 7.3, etc), consider patching all relevant packages to the 
    most recent (including kernel).
    - Rolling restart of each node before upgrade (to confirm in-memory settings are the same as startup-scripts).
    - Run `check-config.sh` on each cluster node (after rolling restart) for any kernel compatibility issues.

- Procedural:

    - Perform a Swarm, UCP and DTR backups pre-upgrade
    - Gather compose file/service/stack files 
    - Generate a UCP Support dump (for point in time) pre-upgrade
    - Preload Engine/UCP/DTR images (in case of offline upgrade method)
        - To prepare existing cluster nodes with new UCP images, pull the needed images onto new workers that you are 
        planning to join using the new upgrade flow.
        ```
        docker login -u <username>
        docker run --rm docker/ucp:3.2.0 images --list: | xargs -L 1 docker pull
        ```
    - Load troubleshooting packages (netshoot, etc)
    - Best order for upgrades: Engine, UCP, and then DTR. Note: The scope of this topic is limited to upgrade instructions for UCP. 

- Upgrade strategy:
For each worker node that requires an upgrade, you can upgrade that node in place or you can replace the node 
with a new worker node. The type of upgrade you perform depends on what is needed for each node:

    - [Automated, in-place cluster upgrade](#automated-in-place-cluster-upgrade): Performed on any 
    manager node. Automatically upgrades the entire cluster. 
    - Manual cluster upgrade: Performed using the CLI or the UCP UI. Automatically upgrades manager 
    nodes and allows you to control the upgrade order of worker nodes. This type of upgrade is more 
    advanced than the automated, in-place cluster upgrade.
        - [Upgrade existing nodes in place](#upgrade-existing-nodes-in-place): Performed using the CLI. 
        Automatically upgrades manager nodes and allows you to control the order of worker node upgrades.
        - [Replace all worker nodes using blue-green deployment](#replace-existing-worker-nodes-using-blue-green-deployment):
        Performed using the CLI. This type of upgrade allows you to 
        stand up a new cluster in parallel to the current code 
        and cut over when complete. This type of upgrade allows you to join new worker nodes, 
        schedule workloads to run on new nodes, pause, drain, and remove old worker nodes 
        in batches of multiple nodes rather than one at a time, and shut down servers to 
        remove worker nodes. This type of upgrade is the most advanced.   

[Upgrade UCP offline](https://docs.docker.com/ee/ucp/admin/install/upgrade-offline/).

## Back up your cluster

Before starting an upgrade, make sure that your cluster is healthy. If a problem
occurs, this makes it easier to find and troubleshoot it.

[Create a backup](/ee/ucp/admin/backup/) of your cluster.
This allows you to recover if something goes wrong during the upgrade process.

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

    In your browser, navigate to **Nodes** in the UCP web interface,
    and check that the node is healthy and is part of the cluster.

## Upgrade UCP
When upgrading Docker Universal Control Plane (UCP) to version {{ page.ucp_version }}, you can choose from 
different upgrade workflows:
> **Important**: In all upgrade workflows, manager nodes are automatically upgraded in place. You cannot control the order 
of manager node upgrades.

- [Automated, in-place cluster upgrade](#automated-in-place-cluster-upgrade): Performed on any 
manager node. Automatically upgrades the entire cluster. 
- Manual cluster upgrade: Performed using the CLI or the UCP UI. Automatically upgrades manager 
nodes and allows you to control the upgrade order of worker nodes. This type of upgrade is more 
advanced than the automated, in-place cluster upgrade.
    - [Upgrade existing nodes in place](#upgrade-existing-nodes-in-place): Performed using the CLI. 
        Automatically upgrades manager nodes and allows you to control the order of worker node upgrades.
    - [Replace all worker nodes using blue-green deployment](#replace-existing-worker-nodes-using-blue-green-deployment): 
    Performed using the CLI. This type of upgrade allows you to 
        stand up a new cluster in parallel to the current code 
        and cut over when complete. This type of upgrade allows you to join new worker nodes, 
        schedule workloads to run on new nodes, pause, drain, and remove old worker nodes 
        in batches of multiple nodes rather than one at a time, and shut down servers to 
        remove worker nodes. This type of upgrade is the most advanced.

### Use the web interface to perform an upgrade

> **Note**: If you plan to add nodes to the UCP cluster, use the [CLI](#use-the-cli-to-perform-an-upgrade) for the upgrade. 

When an upgrade is available for a UCP installation, a banner appears.

![](../../images/upgrade-ucp-1.png){: .with-border}

Clicking this message takes an admin user directly to the upgrade process.
It can be found under the **Upgrade** tab of the **Admin Settings** section.

![](../../images/upgrade-ucp-2.png){: .with-border}

In the **Available Versions** dropdown, select the version you want to update. Copy and paste 
the CLI command provided into a terminal on a manager node to perform the upgrade.

During the upgrade, the web interface will be unavailable, and you should wait
until completion before continuing to interact with it. When the upgrade
completes, you'll see a notification that a newer version of the web interface
is available and a browser refresh is required to see it.

### Use the CLI to perform an upgrade

To upgrade using the CLI, log into a UCP manager node using SSH, and run:

```
# Get the latest version of UCP
docker image pull {{ page.ucp_org }}/{{ page.ucp_repo }}:{{ page.ucp_version }}
```

### Automated in-place cluster upgrade
This workflow is performed when:

- You are not adding new nodes, AND
- the order of worker node upgrades is NOT important.

1. On a manager node, to upgrade the entire cluster (both manager and worker nodes), run commands similar to the following examples: 
    ```
    export ucp_version=3.2.0
    docker image pull docker/ucp:$ucp_version
    docker container run --rm -it \
      --name ucp \
      -v /var/run/docker.sock:/var/run/docker.sock \
      docker/ucp:$ucp_version \
      upgrade --interactive
    ```
        
    **Note**: The `docker image pull` command in the previous example is not needed if the images are pulled as described 
    in the **Procedural** section of [prerequisites](#procedural).

### Upgrade existing nodes in place
This workflow is performed when:

- You are upgrading existing worker nodes in place. This workflow also includes adding additional worker nodes if needed.

1. Upgrade manager nodes
        
    - Run the upgrade on a manager node. The `--manual-worker-upgrade` option automatically upgrades manager nodes first and then 
    allows you to control the upgrade of the UCP components on worker nodes using node labels, as shown in the following example.
    ```
    export ucp_version=3.2.0
    docker image pull docker/ucp:$ucp_version
    docker container run --rm -it \
      --name ucp \
      -v /var/run/docker.sock:/var/run/docker.sock \
      docker/ucp:$ucp_version \
      upgrade --manual-worker-upgrade \
      --interactive
    ```
    Note: The `docker image pull` command in this example is not needed if the images are pulled as described 
    in the **Procedural** section of [prerequisites](#procedural).  
2. Upgrade worker nodes in place
        
    - On the manager node, after the `--manual-worker-upgrade` command completes, remove the `upgrade-hold` label on one or 
    more nodes to upgrade the UCP components on those worker nodes in place: 
    ```
    docker node update --label-rm com.docker.ucp.upgrade-hold <node name or id>
    ```
3. (Optional) Join new worker nodes
       
    - New worker nodes have newer engines already installed and have the new UCP version running when they join the cluster. 
    On the manager node, run commands similar to the following examples to get the Swarm Join token and add 
    new worker nodes:
    ```
    # Get Swarm Join token
    docker swarm join-token worker
    ```
    - On the node to be joined:
    ```
    docker swarm join --token SWMTKN-<YOUR TOKEN> <manager ip>:2377
    ```

### Replace existing worker nodes using blue-green deployment
This workflow is used to create a parallel environment for a new deployment, which can greatly reduce downtime, upgrades 
worker node engines without disrupting workloads, and allows traffic to be migrated to the new environment with 
worker node rollback capability. This type of upgrade creates a parallel environment for reduced downtime and workload disruption.

> **Note**: Steps 2 through 6 can be repeated for groups of nodes - you do not have to replace all worker 
nodes in the cluster at one time.

1. Upgrade manager nodes
        
    - The `--manual-worker-upgrade` command automatically upgrades manager nodes first, and then allows you to control 
    the upgrade of the UCP components on the worker nodes using node labels.
    ```
    export ucp_version=3.2.0
    docker run --rm -it -v /var/run/docker.sock:/var/run/docker.sock docker/ucp:$ucp_version upgrade 
      -i --manual-worker-upgrade
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
### Verify the upgrade
A successful upgrade exits without errors. If you perform a  manual upgrade, you can use the CLI or UI to verify that all nodes are marked as `Ready`.

### Optional: Configure L7 routing
[Configure Interlock or the Ingress controller to point L7 traffic to the new services](/ee/ucp/interlock/).

### Troubleshooting

    - Upgrade compatibility
      
        - The upgrade command automatically checks for multiple `ucp-worker-agents` before 
          proceeding with the upgrade. The existence of multiple `ucp-worker-agents` might indicate 
          that the cluster still in the middle of a prior manual upgrade and you must resolve the 
          conflicting node labels issues before proceeding with the upgrade.

    - Upgrade failures
        - For worker nodes, an upgrade failure can be rolled back by changing the node label back 
          to the previous target version. Rollback of manager nodes is not supported. 

    - Kubernetes errors in node state messages after upgrading UCP 
    (from https://github.com/docker/kbase/how-to-resolve-kubernetes-errors-after-upgrading-ucp/readme.md)

    - The following information applies If you have upgraded to UCP 3.0.0 or newer:
      
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

        - If you have upgraded from UCP 2.2.x to 3.0.x, verify that the ports 179, 6443, 6444 and 10250 are 
        open for Kubernetes traffic.

        - If you have upgraded to UCP 3.1.x, in addition to the ports listed above, do also open 
        ports 9099 and 12388.
  
### Recommended upgrade paths

From UCP 3.0: UCP 3.0 -> UCP 3.1 -> UCP 3.2
From UCP 2.2: UCP 2.2 -> UCP 3.0 -> UCP 3.1 -> UCP 3.2

If you’re running a UCP version earlier than 2.1, first upgrade to the latest 2.1 version, then upgrade to 2.2. Use the following rules for your upgrade path to UCP 2.2:

From UCP 1.1: UCP 1.1 -> UCP 2.1 -> UCP 2.2
From UCP 2.0: UCP 2.0 -> UCP 2.1 -> UCP 2.2

## Where to go next

- [Upgrade DTR](/e/dtr/admin/upgrade/)
