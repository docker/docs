---
title: Restrict services to worker nodes
description: Learn how to configure Docker Universal Control Plane to only allow running services in worker nodes.
keywords: ucp, configuration, worker
---

>{% include enterprise_label_shortform.md %}

Docker Universal Control Plane (UCP) is set to run on both manager and worker nodes by default, however, it can be configured to run only on worker nodes. This ensures that all cluster management functionality stays performant, while also serving to make the cluster more secure.

> **Important**
> 
> In the event that a user deploys a malicious service capable of affecting the node on which it is running, that service will not be able to strike any other nodes in the cluster or have any impact on cluster management functionality.
{: .important} 

## Swarm workloads

To change user options for deploying workloads to manager nodes:

1. Log into the UCP UI with administrator credentials.
2. Navigate to **Admin Settings**.
3. Click **Scheduler** from the left menu.

    ![](../../images/restrict-services-to-worker-nodes-1.png){: .with-border}

4. Select the **Allow users to schedule on all nodes, including UCP managers and DTR nodes.** check box if user services are allowed to run on manager nodes. 

> **Note**
> 
> Creating a grant with the `Scheduler` role against the `/` collection takes
precedence over any other grants with `Node Schedule` on subcollections.

## Kubernetes workloads

By default, UCP clusters use [Taints and Tolerations](https://kubernetes.io/docs/concepts/configuration/taint-and-toleration/) to prevent a user's workload from being deployed on UCP manager or DTR nodes. 

To view the default taint:

```bash
$ kubectl get nodes <ucpmanager> -o json | jq -r '.spec.taints | .[]'
{
  "effect": "NoSchedule",
  "key": "com.docker.ucp.manager"
}
```

> **Note**
> 
> Workloads deployed by an administrator in the `kube-system` namespace do not follow these scheduling constraints. If an administrator deploys a workload in the `kube-system` namespace, a toleration is applied to bypass this taint, which then schedules the workload on all node types. 

### Allow administrators to schedule workloads on manager / DTR nodes

To allow administrators to deploy workloads across all nodes types, select the **Allow administrators to deploy containers on UCP managers or nodes running DTR** check box in the UCP UI. 

![](../../images/restrict-services-to-worker-nodes-2.png){: .with-border}

When this option is enabled, UCP will apply a toleration to all new workloads deployed by administers, which allows the pods to be scheduled on all node types. 

For existing workloads, administrators must use `kubectl edit <object> <workload>` or the UCP UI to add the following toleration to the pod specification:

```bash
tolerations:
- key: "com.docker.ucp.manager"
  operator: "Exists"
```

To verify that the toleration was applied successfully:

```bash
$ kubectl get <object> <workload> -o json | jq -r '.spec.template.spec.tolerations | .[]'
{
  "key": "com.docker.ucp.manager",
  "operator": "Exists"
}
```

### Allow users and service accounts to schedule workloads on manager / DTR nodes

To allow Kubernetes users and service accounts to deploy workloads across all node types in a cluster, an administrator must select the **Allow all authenticated users, including service accounts, to schedule on all nodes, including UCP managers and DTR nodes.** check box in the UCP UI. 

![](../../images/restrict-services-to-worker-nodes-3.png){: .with-border}

When this option is enabled, UCP will apply a toleration to all new workloads deployed by Kubernetes, which allows the pods to be scheduled on all node types.

For existing workloads, users will need to edit the pod specification, using `kubectl edit <object> <workload>` or the UCP UI to add the following toleration to the pod specification:

```bash
tolerations:
- key: "com.docker.ucp.manager"
  operator: "Exists"
```

To verify that the toleration was applied successfully:

```bash
$ kubectl get <object> <workload> -o json | jq -r '.spec.template.spec.tolerations | .[]'
{
  "key": "com.docker.ucp.manager",
  "operator": "Exists"
}
```

> **Note**
>
> There is a `NoSchedule` taint value available on UCP managers and DTR nodes. If the option to schedule managers and DTR nodes is disabled, a toleration for that taint will not be applied to the deployments. Note that workloads are only scheduled on the nodes if the Kubernetes workload is deployed in the `kube-system` namespace.

## Where to go next

- [Deploy an application package](/ee/ucp/deploy-application-package/)
- [Deploy a Swarm workload](/ee/ucp/swarm/)
- [Deploy a Kubernetes workload](/ee/ucp/kubernetes//)
