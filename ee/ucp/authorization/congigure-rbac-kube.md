---
title: Enable Kuberenetes RBAC
description: Learn how configure role-based access control for Kubernetes
keywords: Kuberenetes, ucp, RBAC
---

UCP 3.0 used its own role-based asccess control (RBAC) for Kubernetes clusters. New in UCP 3.1 is the ability to use Kube RBAC. The benefits of doing this are:

- Many ecosystem applications and integrations expect Kube RBAC as a part of their YAML files to provide access to service accounts. - Organizations planning to run UCP both on-premesis as well as in hosted cloud services want to run Kubernetes applications on both sets of environments, without manually changing RBAC for their YAML file.

Kubernetes RBAC is turned on by default when customers upgrade to UCP 3.1. See [RBAC authorization in Kubernetes](https://v1-8.docs.kubernetes.io/docs/admin/authorization/rbac/) for more information about Kubernetes.

Starting with UCP 3.1, Kubernetes & Swarm roles have separate views. You can view all the roles for a particular cluster under **Access Control** then **Users**. Select Kubernetes or Swarm to view the specific roles for each.

## Creating roles

Kubernetes provides 2 types of roles:

- `ClusterRoleBinding` which applies to all namespaces
- `RoleBinding1` which applies to a specific namespace

You create Kubernetes roles either through the CLI using `kubectl` or through the UCP web interface.

To create a Kuberenetes role in the UCP web interface:

1. 1 Go to the UCP web UI.
2. Navigate to the **Access Control**.
3. In the lefthand menu, select  **Grants**.

![Roles in UCP](../images/kube-rbac-grants.png)

4. Select the **Kuberneters** tab at the top of the window.
5.
