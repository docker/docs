---
title: Grant role-access to cluster resources
description: Learn how to grant users and teams access to cluster resources with role-based access control.
keywords: rbac, ucp, grant, role, permission, authentication, authorization, namespace, Kubernetes
redirect_from:
  - /datacenter/ucp/3.0/guides/authorization/grant-permissions/
---

>{% include enterprise_label_shortform.md %}

Docker Enterprise administrators can create _grants_ to control how users and
organizations access [resource sets](group-resources.md).

A grant defines _who_ has _how much_ access to _what_ resources. Each grant is a
1:1:1 mapping of _subject_, _role_, and _resource set_. For example, you can
grant the "Prod Team" "Restricted Control" over services in the "/Production"
collection. 

A common workflow for creating grants has four steps:

- Add and configure **subjects** (users, teams, and service accounts).
- Define custom **roles** (or use defaults) by adding permitted API operations
  per type of resource.
- Group cluster **resources** into Swarm collections or Kubernetes namespaces.
- Create **grants** by combining subject + role + resource set.

## Creating grants
 To create a grant:
 
 1. Log in to the UCP web UI.
 2. Click **Access Control**.
 3. Click **Grants**.
 4. In the Grants window, select **Kubernetes** or **Swarm**.

### Kubernetes grants

With Kubernetes orchestration, a grant is made up of *subject*, *role*, and
*namespace*.

> Note
> 
> This section assumes that you have created objects for the grant: subject, role,
> namespace.
{: .important}

To create a Kubernetes grant (role binding) in UCP:

1. Click **Create Role Binding**.
2. Under Subject, select **Users**, **Organizations**, or **Service Account**.
    - For Users, select the user from the pull-down menu (these should have already been created as objects).
    - For Organizations, select the Organization and Team (optional) from the pull-down menu.
    - For Service Account, select the Namespace and Service Account from the pull-down menu.
3. Click **Next** to save your selections.
4. Under Resource Set, toggle the **Apply Role Binding to all namespaces (Cluster Role Binding)** switch.
5. Click **Next**.
6. Under Role, select a cluster role.
7. Click **Create**.

### Swarm grants

With Swarm orchestration, a grant is made up of *subject*, *role*, and
*collection*.

> Note
> 
> This section assumes that you have created objects to grant: teams/users,
> roles (built-in or custom), and a collection.

![](../images/ucp-grant-model-0.svg){: .with-border}
![](../images/ucp-grant-model.svg){: .with-border}

To create a Swarm grant in UCP:

1. Click **Create Grant**.
2. Under Subject, select **Users** or **Organizations**.
    - For Users, select a user from the pull-down menu.
    - For Organizations, select the Organization and Team (optional) from the pull-down menu.
3. Click **Next**.
4. Under Resource Set, click **View Children** until you get to the desired collection.
5. Click **Select Collection**.
6. Click **Next**.
7. Under Role, select a role from the pull-down menu.
8. Click **Create**.

> Note
>
> By default, all new users are placed in the `docker-datacenter` organization.
> To apply permissions to all Docker Enterprise users, create a grant with the
> `docker-datacenter` organization as a subject.
{: .important}

## Where to go next

- [Deploy a simple stateless app with RBAC](deploy-stateless-app.md)
