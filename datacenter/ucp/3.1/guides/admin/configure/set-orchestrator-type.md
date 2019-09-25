---
title: Set the orchestrator type for a node
description: Learn how to specify the orchestrator for nodes in a Docker Enterprise Edition cluster.
keywords: Docker EE, UCP, cluster, orchestrator
---

When you add a node to the cluster, the node's workloads are managed by a
default orchestrator, either Docker Swarm or Kubernetes. When you install
Docker EE, new nodes are managed by Docker Swarm, but you can change the
default orchestrator to Kubernetes in the administrator settings.

Changing the default orchestrator doesn't affect existing nodes in the cluster.
You can change the orchestrator type for individual nodes in the cluster
by navigating to the node's configuration page in the Docker EE web UI.

## Change the orchestrator for a node

You can change the current orchestrator for any node that's joined to a
Docker EE cluster. The available orchestrator types are **Kubernetes**,
**Swarm**, and **Mixed**.

The **Mixed** type enables workloads to be scheduled by Kubernetes and Swarm
both on the same node. Although you can choose to mix orchestrator types on the
same node, this isn't recommended for production deployments because of the
likelihood of resource contention.

Change a node's orchestrator type on the **Edit node** page:

1.  Log in to the Docker EE web UI with an administrator account.
2.  Navigate to the **Nodes** page, and click the node that you want to assign
    to a different orchestrator.
3.  In the details pane, click **Configure** and select **Details** to open
    the **Edit node** page.
4.  In the **Orchestrator properties** section, click the orchestrator type
    for the node.
5.  Click **Save** to assign the node to the selected orchestrator.

    ![](../../images/change-orchestrator-for-node-1.png){: .with-border}

## What happens when you change a node's orchestrator

When you change the orchestrator type for a node, existing workloads are
evicted, and they're not migrated to the new orchestrator automatically.
If you want the workloads to be scheduled by the new orchestrator, you must
migrate them manually. For example, if you deploy WordPress on a Swarm
node, and you change the node's orchestrator type to Kubernetes, Docker EE
doesn't migrate the workload, and WordPress continues running on Swarm. In
this case, you must migrate your WordPress deployment to Kubernetes manually.

The following table summarizes the results of changing a node's orchestrator.

|                  Workload                   |                   On orchestrator change                   |
| ------------------------------------------- | ---------------------------------------------------------- |
| Containers                                  | Container continues running in node                        |
| Docker service                              | Node is drained, and tasks are rescheduled to another node |
| Pods and other imperative resources         | Continue running in node                                   |
| Deployments and other declarative resources | Might change, but for now, continue running in node        |

If a node is running containers, and you change the node to Kubernetes, these
containers will continue running, and Kubernetes won't be aware of them, so
you'll be in the same situation as if you were running in `Mixed` node.

> Be careful when mixing orchestrators on a node.
>
> When you change a node's orchestrator, you can choose to run the node in a
> mixed mode, with both Kubernetes and Swarm workloads. The `Mixed` type
> is not intended for production use, and it may impact existing workloads
> on the node.
>
> This is because the two orchestrator types have different views of the node's
> resources, and they don't know about each other's workloads. One orchestrator
> can schedule a workload without knowing that the node's resources are already
> committed to another workload that was scheduled by the other orchestrator.
> When this happens, the node could run out of memory or other resources.
>
> For this reason, we recommend against mixing orchestrators on a production
> node. 
{: .warning}

## Set the default orchestrator type for new nodes

You can set the default orchestrator for new nodes to **Kubernetes** or
**Swarm**.

To set the orchestrator for new nodes:

1.  Log in to the Docker EE web UI with an administrator account.
2.  Open the **Admin Settings** page, and in the left pane, click **Scheduler**.
3.  Under **Set orchestrator type for new nodes** click **Swarm**
    or **Kubernetes**.
4.  Click **Save**.
    
    ![](../../images/join-nodes-to-cluster-1.png){: .with-border}

From now on, when you join a node to the cluster, new workloads on the node
are scheduled by the specified orchestrator type. Existing nodes in the cluster
aren't affected.

Once a node is joined to the cluster, you can
[change the orchestrator](#change-the-orchestrator-for-a-node) that schedules its
workloads.

## Choosing the orchestrator type

The workloads on your cluster can be scheduled by Kubernetes or by Swarm, or
the cluster can be mixed, running both orchestrator types. If you choose to
run a mixed cluster, be aware that the different orchestrators aren't aware of
each other, and there's no coordination between them.

We recommend that you make the decision about orchestration when you set up the
cluster initially. Commit to Kubernetes or Swarm on all nodes, or assign each
node individually to a specific orchestrator. Once you start deploying workloads,
avoid changing the orchestrator setting. If you do change the orchestrator for a
node, your workloads are evicted, and you must deploy them again through the
new orchestrator.

> Node demotion and orchestrator type
>
> When you promote a worker node to be a manager, its orchestrator type
> automatically changes to `Mixed`. If you demote the same node to be a worker,
> its orchestrator type remains as `Mixed`.
{: important}

## Use the CLI to set the orchestrator type

Set the orchestrator on a node by assigning the orchestrator labels,
`com.docker.ucp.orchestrator.swarm` or `com.docker.ucp.orchestrator.kubernetes`,
to `true`.

To schedule Swarm workloads on a node:

```bash
docker node update --label-add com.docker.ucp.orchestrator.swarm=true <node-id>
```

To schedule Kubernetes workloads on a node:

```bash
docker node update --label-add com.docker.ucp.orchestrator.kubernetes=true <node-id>
```

To schedule Kubernetes and Swarm workloads on a node:

```bash
docker node update --label-add com.docker.ucp.orchestrator.swarm=true <node-id>
docker node update --label-add com.docker.ucp.orchestrator.kubernetes=true <node-id>
```

> Mixed nodes
> 
> Scheduling both Kubernetes and Swarm workloads on a node is not recommended
> for production deployments, because of the likelihood of resource contention.
{: .warning}

To change the orchestrator type for a node from Swarm to Kubernetes:

```bash
docker node update --label-add com.docker.ucp.orchestrator.kubernetes=true <node-id>
docker node update --label-rm com.docker.ucp.orchestrator.swarm <node-id>
```

UCP detects the node label change and updates the Kubernetes node accordingly.

Check the value of the orchestrator label by inspecting the node:

```bash
docker node inspect <node-id> | grep -i orchestrator
```

The `docker node inspect` command returns the node's configuration, including
the orchestrator:

```bash
"com.docker.ucp.orchestrator.kubernetes": "true"
```

> Orchestrator label
>
> The `com.docker.ucp.orchestrator` label isn't displayed in the **Labels**
> list for a node in the Docker EE web UI.
{: .important}

## Set the default orchestrator type for new nodes

The default orchestrator for new nodes is a setting in the Docker EE
configuration file:

```
default_node_orchestrator = "swarm"
```

The value can be `swarm` or `kubernetes`.

## Where to go next

- [Set up Docker EE by using a config file](ucp-configuration-file.md)