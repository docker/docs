---
title: Join Linux nodes to your cluster
description: Learn how to scale a Docker Enterprise Edition cluster by adding manager and worker nodes.
keywords: Docker EE, UCP, cluster, scale, worker, manager
---

Docker EE is designed for scaling horizontally as your applications grow in
size and usage. You can add or remove nodes from the cluster to scale it
to your needs. You can join Windows Server 2016, IBM z System, and Linux nodes
to the cluster.

Because Docker EE leverages the clustering functionality provided by Docker
Engine, you use the [docker swarm join](/engine/swarm/swarm-tutorial/add-nodes.md)
command to add more nodes to your cluster. When you join a new node, Docker EE
services start running on the node automatically.

## Node roles

When you join a node to a cluster, you specify its role: manager or worker.

- **Manager**: Manager nodes are responsible for cluster management
  functionality and dispatching tasks to worker nodes. Having multiple
  manager nodes allows your swarm to be highly available and tolerant of
  node failures.

  Manager nodes also run all Docker EE components in a replicated way, so
  by adding additional manager nodes, you're also making the cluster highly
  available.
  [Learn more about the Docker EE architecture.](/enterprise/docker-ee-architecture.md)

- **Worker**: Worker nodes receive and execute your services and applications.
  Having multiple worker nodes allows you to scale the computing capacity of
  your cluster.

  When deploying Docker Trusted Registry in your cluster, you deploy it to a
  worker node.

## Join a node to the cluster

You can join Windows Server 2016, IBM z System, and Linux nodes to the cluster,
but only Linux nodes can be managers.

To join nodes to the cluster, go to the Docker EE web UI and navigate to the
**Nodes** page.

1.  Click **Add Node** to add a new node.
2.  Select the type of node to add, **Windows** or **Linux**.
2.  Click **Manager** if you want to add the node as a manager.
3.  Check the **Use a custom listen address** option to specify the address
    and port where new node listens for inbound cluster management traffic.
4.  Check the **Use a custom listen address** option to specify the
    IP address that's advertised to all members of the cluster for API access.

![](../../../images/join-nodes-to-cluster-2.png){: .with-border}

Copy the displayed command, use SSH to log in to the host that you want to
join to the cluster, and run the `docker swarm join` command on the host.

To add a Windows node, click **Windows** and follow the instructions in
[Join Windows worker nodes to a cluster](join-windows-nodes-to-cluster.md).

After you run the join command in the node, the node is displayed on the
**Nodes** page in the Docker EE web UI. From there, you can change the node's
cluster configuration, including its assigned orchestrator type.
[Learn how to change the orchestrator for a node](../set-orchestrator-type.md).    

## Pause or drain a node

Once a node is part of the cluster, you can configure the node's availability
so that it is:

- **Active**: the node can receive and execute tasks.
- **Paused**: the node continues running existing tasks, but doesn't receive
  new tasks.
- **Drained**: the node won't receive new tasks. Existing tasks are stopped and
  replica tasks are launched in active nodes.

Pause or drain a node from the **Edit Node** page:

1.  In the Docker EE web UI, browse to the **Nodes** page and select the node.
2.  In the details pane, click **Configure** and select **Details** to open
    the **Edit Node** page.
3.  In the **Availability** section, click **Active**, **Pause**, or **Drain**.  
4.  Click **Save** to change the availability of the node.

![](../../../images/join-nodes-to-cluster-3.png){: .with-border}

## Promote or demote a node

You can promote worker nodes to managers to make UCP fault tolerant. You can
also demote a manager node into a worker.

To promote or demote a manager node:

1.  Navigate to the **Nodes** page, and click the node that you want to demote.
2.  In the details pane, click **Configure** and select **Details** to open
    the **Edit Node** page.
3.  In the **Role** section, click **Manager** or **Worker**.
4.  Click **Save** and wait until the operation completes.
5.  Navigate to the **Nodes** page, and confirm that the node role has changed.

If you're load-balancing user requests to Docker EE across multiple manager
nodes, don't forget to remove these nodes from your load-balancing pool when
you demote them to workers.

## Remove a node from the cluster

You can remove worker nodes from the cluster at any time:

1.  Navigate to the **Nodes** page and select the node.
2.  In the details pane, click **Actions** and select **Remove**.
3.  Click **Confirm** when you're prompted.

Since manager nodes are important to the cluster overall health, you need to
be careful when removing one from the cluster.

To remove a manager node:

1. Make sure all nodes in the cluster are healthy. Don't remove manager nodes
if that's not the case.
2. Demote the manager node into a worker.
3. Now you can remove that node from the cluster.

## Use the CLI to manage your nodes

You can use the Docker CLI client to manage your nodes from the CLI. To do
this, configure your Docker CLI client with a [UCP client bundle](../../../user-access/cli.md).

Once you do that, you can start managing your UCP nodes:

```bash
docker node ls
```
