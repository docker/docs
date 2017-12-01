---
title: Grant access to cluster resources
description: Learn how to grant users and teams access to cluster resources with role-based access control.
keywords: rbac, ucp, grant, role, permission, authentication, authorization
redirect_from:
- /ucp/
ui_tabs:
- version: ucp-3.0
  orhigher: true
- version: ucp-2.2
  orlower: true
next_steps:
- path: /deploy/rbac/usermgmt-create-subjects/
  title: Create and configure users and teams
- path: /deploy/rbac/usermgmt-define-roles/
  title: Create roles to authorize access
- path: /deploy/rbac/resources-isolate-volumes/
  title: Isolate volumes
---

{% if include.ui %}

Docker EE administrators can create *grants* to control how users and
organizations access resources.

A grant is made up of *subject*, *role*, and *resource group*.


{% if include.version=="ucp-3.0" %}

## Kubernetes grants

With Kubernetes orchestration, a grant is made up of *subject*, *role*, and
*namespace*.




## Swarm grants

With Swarm orchestration, a grant is made up of *subject*, *role*, and
*collection*.

![](../images/ucp-grant-model-0.svg){: .with-border}

A grant defines who (subject) has how much access (role) to a set of resources
(collection). Each grant is a 1:1:1 mapping of subject, role, collection. For
example, you can grant the "Prod Team" "Restricted Control"of the "/Production"
collection.

A common workflow for creating grants has four steps:

- Add and configure **subjects** (users and teams).
- Define custom **roles** (or use defaults) by adding permitted API operations
  per resource type.
- Group cluster **resources** into Swarm collections or Kubernetes namespaces.
- Create **grants** by marrying subject + role + resource.

![](../images/ucp-grant-model.svg){: .with-border}

### Create a Swarm grant

You can create grants after creating users, collections, and roles (if using
custom roles).

To create a grant in UCP:

1. Click **Grants** under **User Management**.
2. Click **Create Grant**.
3. On the Collections tab, click **Collections** (for Swarm) or **Namespaces** (for Kubernetes).
4. Click **View Children** until you get to the desired resource group and **Select**.
5. On the Roles tab, select a role.
6. On the Subjects tab, select a user, team, or organization to authorize.
4. Click **Create**.

> By default, all new users are placed in the `docker-datacenter` organization.
> To apply permissions to all Docker EE users, create a grant with the
> `docker-datacenter` org as a subject.


{% elsif include.version=="ucp-2.2" %}

## Swarm grants

With Swarm orchestration, a grant is made up of *subject*, *role*, and
*collection*.

![](../images/ucp-grant-model-0.svg){: .with-border}

A grant defines who (subject) has how much access (role) to a set of resources
(collection). Each grant is a 1:1:1 mapping of subject, role, collection. For
example, you can grant the "Prod Team" "Restricted Control"of the "/Production"
collection.

A common workflow for creating grants has four steps:

- Add and configure **subjects** (users and teams).
- Define custom **roles** (or use defaults) by adding permitted API operations
  per resource type.
- Group cluster **resources** into Swarm collections.
- Create **grants** by marrying subject + role + resource.

![](../images/ucp-grant-model.svg){: .with-border}

### Create a Swarm grant

You can create grants after creating users, collections, and roles (if using custom roles).

To create a grant in UCP:

1. Click **Grants** under **User Management**.
2. Click **Create Grant**.
3. On the Collections tab, click **Collections**.
4. Click **View Children** until you get to the desired resource group and **Select**.
5. On the Roles tab, select a role.
6. On the Subjects tab, select a user, team, or organization to authorize.
4. Click **Create**.

> By default, all new users are placed in the `docker-datacenter` organization.
> To apply permissions to all Docker EE users, create a grant with the
> `docker-datacenter` org as a subject.

{% endif %}
{% endif %}
