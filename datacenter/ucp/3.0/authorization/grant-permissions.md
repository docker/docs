---
title: Grant role-access to cluster resources
description: Learn how to grant users and teams access to cluster resources with role-based access control.
keywords: rbac, ucp, grant, role, permission, authentication, authorization, namespace, Kubernetes
redirect_from:
  - /datacenter/ucp/3.0/guides/authorization/grant-permissions/
---

Docker EE administrators can create _grants_ to control how users and
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

## Kubernetes grants

With Kubernetes orchestration, a grant is made up of *subject*, *role*, and
*namespace*.

> This section assumes that you have created objects for the grant: subject, role,
> namespace.
{: .important}

To create a Kubernetes grant in UCP:

1. Click **Grants** under **User Management**.
2. Click **Create Grant**.
3. Click **Namespaces** under **Kubernetes**.
4. Find the desired namespace and click **Select Namespace**.
5. On the **Roles** tab, select a role.
6. On the **Subjects** tab, select a user, team, organization, or service
   account to authorize.
7. Click **Create**.

## Swarm grants

With Swarm orchestration, a grant is made up of *subject*, *role*, and
*collection*.

> This section assumes that you have created objects to grant: teams/users,
> roles (built-in or custom), and a collection.

![](../images/ucp-grant-model-0.svg){: .with-border}
![](../images/ucp-grant-model.svg){: .with-border}

To create a grant in UCP:

1. Click **Grants** under **User Management**.
2. Click **Create Grant**.
3. On the Collections tab, click **Collections** (for Swarm).
4. Click **View Children** until you get to the desired collection and **Select**.
5. On the **Roles** tab, select a role.
6. On the **Subjects** tab, select a user, team, or organization to authorize.
7. Click **Create**.

> By default, all new users are placed in the `docker-datacenter` organization.
> To apply permissions to all Docker EE users, create a grant with the
> `docker-datacenter` org as a subject.
{: .important}

## Where to go next

- [Deploy a simple stateless app with RBAC](deploy-stateless-app.md)
