---
title: Create roles and authorize operations
description: Learn how to create roles and set permissions in Docker Universal Control Plane.
keywords: rbac, authorization, authentication, users, teams, UCP
redirect_from:
- /ucp/
ui_tabs:
- version: ucp-3.0
  orhigher: true
- version: ucp-2.2
  orlower: true
---

{% if include.ui %}
{% if include.version=="ucp-3.0" %}

Docker Universal Control Plane has two types of users: administrators and
regular users. Administrators can make changes to the UCP cluster, while regular
users have permissions that range from no access to full control over resources
such as volumes, networks, images, and containers.

Users are grouped into teams and organizations.

![Diagram showing UCP permission levels](../images/role-diagram.svg)

Administrators apply *grants* to users, teams, and organizations to give
permissions to swarm resources.

## Administrator users

In Docker UCP, only users with administrator privileges can make changes to
cluster settings. This includes:

* Managing user permissions by creating grants.
* Managing cluster configurations, like adding and removing nodes.

## Roles

A role is a set of permitted API operations on a collection that you
can assign to a specific user, team, or organization by using a grant.

UCP administrators view and manage roles by navigating to the **Roles** page.

The system provides the following default roles:

- **None**:  Users have no access to swarm resources. This role maps to the
`No Access` role in UCP 2.1.x.
- **View Only**: Users can view resources but can't create them.
- **Restricted Control**: Users can view and edit resources but can't run a
service or container in a way that affects the node where it's running. Users
_cannot_: mount a node directory, `exec` into containers, or run containers in
privileged mode or with additional kernel capabilities.
- **Scheduler**: Users can view nodes and schedule workloads on them. Worker nodes
and manager nodes are affected by `Scheduler` grants. By default, all users get
a grant with the `Scheduler` role against the `/Shared` collection. Having
`Scheduler` access doesn't allow users to view workloads on these nodes--they
need the appropriate resource permissions such as `Container View`.
- **Full Control**: Users can view and edit all granted resources. They can create
containers without any restriction, but can't see the containers of other users.

![Diagram showing UCP permission levels](../images/permissions-ucp.svg)

Administrators can create a custom role that has Docker API permissions
that specify the API actions that a subject may perform.

The **Roles** page lists the available roles, including the default roles
and any custom roles that administrators have created. In the **Roles**
list, click a role to see the API operations that it uses. For example, the
`Scheduler` role has two of the node operations, `Schedule` and `View`.

## Create a custom role

Click **Create role** to create a custom role and define the API operations
that it uses. When you create a custom role, all of the APIs that you can use
are listed on the **Create Role** page. For example, you can create a custom
role that uses the node operations, `Schedule`, `Update`, and `View`, and you
might give it a name like "Node Operator".

1. Click **Roles** under **User Management**.
2. Click **Create Role**.
3. Input the role name on the **Details** page.
4. Click **Operations**.
5. Select the permitted operations per resource type.
6. Click **Create**.

![](../images/custom-role.png){: .with-border}

You can give a role a global name, like "Remove Images", which might enable
the **Remove** and **Force Remove** operations for images. You can apply a
role with the same name to different collections.

Only an administrator can create and remove roles. Roles are always enabled.
Roles can't be edited, so to change a role's API operations, you must delete it
and create it again.

You can't delete a custom role if it's used in a grant. You must first delete
the grants that use the role.




{% elsif include.version=="ucp-2.2" %}

Docker Universal Control Plane has two types of users: administrators and
regular users. Administrators can make changes to the UCP cluster, while regular
users have permissions that range from no access to full control over resources
such as volumes, networks, images, and containers.

Users are grouped into teams and organizations.

![Diagram showing UCP permission levels](../images/role-diagram.svg)

Administrators apply *grants* to users, teams, and organizations to give
permissions to swarm resources.

## Administrator users

In Docker UCP, only users with administrator privileges can make changes to
cluster settings. This includes:

* Managing user permissions by creating grants.
* Managing cluster configurations, like adding and removing nodes.

## Roles

A role is a set of permitted API operations on a collection that you
can assign to a specific user, team, or organization by using a grant.

UCP administrators view and manage roles by navigating to the **Roles** page.

The system provides the following default roles:

- **None**:  Users have no access to swarm resources. This role maps to the
`No Access` role in UCP 2.1.x.
- **View Only**: Users can view resources but can't create them.
- **Restricted Control**: Users can view and edit resources but can't run a
service or container in a way that affects the node where it's running. Users
_cannot_: mount a node directory, `exec` into containers, or run containers in
privileged mode or with additional kernel capabilities.
- **Scheduler**: Users can view nodes and schedule workloads on them. Worker nodes
and manager nodes are affected by `Scheduler` grants. By default, all users get
a grant with the `Scheduler` role against the `/Shared` collection. Having
`Scheduler` access doesn't allow users to view workloads on these nodes--they
need the appropriate resource permissions such as `Container View`.
- **Full Control**: Users can view and edit all granted resources. They can create
containers without any restriction, but can't see the containers of other users.

![Diagram showing UCP permission levels](../images/permissions-ucp.svg)

Administrators can create a custom role that has Docker API permissions
that specify the API actions that a subject may perform.

The **Roles** page lists the available roles, including the default roles
and any custom roles that administrators have created. In the **Roles**
list, click a role to see the API operations that it uses. For example, the
`Scheduler` role has two of the node operations, `Schedule` and `View`.

## Create a custom role

Click **Create role** to create a custom role and define the API operations
that it uses. When you create a custom role, all of the APIs that you can use
are listed on the **Create Role** page. For example, you can create a custom
role that uses the node operations, `Schedule`, `Update`, and `View`, and you
might give it a name like "Node Operator".

1. Click **Roles** under **User Management**.
2. Click **Create Role**.
3. Input the role name on the **Details** page.
4. Click **Operations**.
5. Select the permitted operations per resource type.
6. Click **Create**.

![](../images/custom-role.png){: .with-border}

You can give a role a global name, like "Remove Images", which might enable
the **Remove** and **Force Remove** operations for images. You can apply a
role with the same name to different collections.

Only an administrator can create and remove roles. Roles are always enabled.
Roles can't be edited, so to change a role's API operations, you must delete it
and create it again.

You can't delete a custom role if it's used in a grant. You must first delete
the grants that use the role.
{% endif %}

## Where to go next

* [Create and manage users](create-and-manage-users.md)
* [Create and manage teams](create-and-manage-teams.md)
* [Docker Reference Architecture: Securing Docker EE and Security Best Practices](https://success.docker.com/Architecture/Docker_Reference_Architecture%3A_Securing_Docker_EE_and_Security_Best_Practices)
*
{% endif %}
