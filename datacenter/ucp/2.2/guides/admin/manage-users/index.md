---
title: Authentication and authorization
description: Manage access to containers, services, volumes, and networks by using role-based access control
keywords: ucp, grant, role, permission, authentication, authorization
---

With Docker Universal Control Plane, you get to control who can create and
edit container resources in your swarm, like services, images, networks,
and volumes. You can grant and manage permissions to enforce fine-grained
access control as needed.

## Grant access to swarm resources  

If you're a UCP administrator, you can create *grants* to control how users 
and organizations access swarm resources.

A grant is made up of a *subject*, a *role*, and a *resource collection*.
A grant defines who (subject) has how much access (role) 
to a set of resources (collection). Each grant is a 1:1:1 mapping of 
subject, role, collection. For example, you can create a grant that 
specifies that the "Prod Team" gains "View Only" permission against
resources in the "/Production" collection.

The usual workflow for creating grants has three steps.

1.  Set up your users and teams. For example, you might want three teams,
    Dev, QA, and Prod.
2.  Organize swarm resources into separate collections that each team will use.
3.  Grant access against resource collections for your teams.

## Subjects

A subject represents a user, team, or organization. A subject is granted a role 
for a collection of resources.

-   **User**: A person that the authentication backend validates. You can
    assign users to one or more teams and one or more organizations.
-   **Organization**: A group of users that share a specific set of
    permissions, defined by the roles of the organization.
-   **Team**: A group of users that share a set of permissions defined in the
    team itself. A team exists only as part of an organization, and all of its
    members must be members of the organization. Team members share
    organization permissions. A team can be in one organization only.
-   **Administrator**: A person who identifies which operations can be
    performed against specific resources and who can perform these actions.
    An administrator can create and manage role assignments against any user,
    team, and organization in the system. Only administrators can manage
    grants. 

## Roles

A role is a set of permitted API operations that you can assign to a specific
subject and collection by using a grant. UCP administrators view and manage
roles by navigating to the **User Management > Roles** page. 

The system provides the following default roles: 

| Built-in role        | Description |
|----------------------|-------------|
| `View Only`          | The user can view resources like services, volumes, and networks but can't create them. |
| `Restricted Control` | The user can view and edit volumes, networks, and images but can't run a service or container in a way that might affect the node where it's running. The user can't mount a node directory and can't `exec` into containers. Also, The user can't run containers in privileged mode or with additional kernel capabilities. |
| `Full Control`       | The user can view and edit volumes, networks, and images, They can create containers without any restriction, but can't see other users' containers. |
| `Scheduler`          | The user can schedule and view workloads on worker nodes. By default, all users get a grant with the `Scheduler` role against the `/Shared` collection. |
| `Admin`              | The user has full access to all resources, like volumes, networks, images, and containers. |

Administrators can create a custom role that has Docker API permissions
that specify the API actions that a subject may perform.

The **Roles** page lists the available roles, including the default roles
and any custom roles that administrators have created. In the **Roles**
list, click a role to see the API operations that it uses. For example, the
`Scheduler` role has two of the node operations, `Schedule` and `View`.

Click **Create role** to create a custom role and define the API operations
that it uses. When you create a custom role, all of the APIs that you can use
are listed on the **Create Role** page. For example, you can create a custom
role that uses all of the node operations, `Join Token`, `Schedule`,
`Update`, and `View`, and you might name it "Node Operator".

You can give a role a global name, like "Remove Images", which might enable 
the **Remove** and **Force Remove** operations for images. You can apply a
role with the same name to different resource collections.

Only an administrator can create and remove roles. An administrator
can enable and disable roles in the system. Roles can't be edited, so
to change a role's API operations, you must delete it and recreate it.

You can't delete a custom role if it's used in a grant. You must first delete
the grants that use the role. 

## Resource collections

Docker EE enables controlling access to container resources by using
*collections*. A collection is a grouping of container resources, like
volumes, networks, secrets, and services, that you access by specifying
a directory-like path. For more info, see 
[Manage access to resources by using collections](manage-access-with-collections.md).

## Transition from UCP 2.1 access control

-   Your existing access labels and permissions are migrated automatically
    during an upgrade from UCP 2.1.x.
-   Unlabeled "user-owned" resources are migrated into the user's private
    collection, in `/Shared/Private/<username>`.
-   Old access control labels are migrated into `/Shared/Legacy/<labelname>`.
-   When deploying a resource, choose a collection instead of an access label.
-   Use grants for access control, instead of unlabeled permissions.

## Where to go next

-  [Deploy a service with view-only access across an organization](deploy-view-only-service.md)
-  [Isolate volumes between two different teams](isolate-volumes-between-teams.md)
-  [Isolate swarm nodes between two different teams](isolate-nodes-between-teams.md)

