---
title: Define roles with authorized API operations
description: Learn how to create roles and set permissions in Docker Universal Control Plane.
keywords: rbac, authorization, authentication, users, teams, UCP
---

A role defines a set of API operations permitted against a resource set.
You apply roles to users and teams by creating grants.

![Diagram showing UCP permission levels](../images/permissions-ucp.svg)

## Default roles

You can define custom roles or use the following built-in roles:

| Built-in role        | Description                                                                                                                                                                                                                                                           |
|:---------------------|:----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `None`               | Users have no access to Swarm or Kubernetes resources. Maps to `No Access` role in UCP 2.1.x.                                                                                                                                                                         |
| `View Only`          | Users can view resources but can't create them.                                                                                                                                                                                                                       |
| `Restricted Control` | Users can view and edit resources but can't run a service or container in a way that affects the node where it's running. Users _cannot_ mount a node directory, `exec` into containers, or run containers in privileged mode or with additional kernel capabilities. |
| `Scheduler`          | Users can view nodes (worker and manager) and schedule (not view) workloads on these nodes. By default, all users are granted the `Scheduler` role against the `/Shared` collection. (To view workloads, users need permissions such as `Container View`).            |
| `Full Control`       | Users can view and edit all granted resources. They can create containers without any restriction, but can't see the containers of other users.                                                                                                                       |


## Create a custom role

The **Roles** page lists all default and custom roles applicable in the
organization.

You can give a role a global name, such as "Remove Images", which might enable the
**Remove** and **Force Remove** operations for images. You can apply a role with
the same name to different resource sets.

1. Click **Roles** under **User Management**.
2. Click **Create Role**.
3. Input the role name on the **Details** page.
4. Click **Operations**. All available API operations are displayed.
5. Select the permitted operations per resource type.
6. Click **Create**.

![](../images/custom-role-30.png){: .with-border}

> **Some important rules regarding roles**:
> - Roles are always enabled.
> - Roles can't be edited. To edit a role, you must delete and recreate it.
> - Roles used within a grant can be deleted only after first deleting the grant.
> - Only administrators can create and delete roles.

## Where to go next

- [Create and configure users and teams](create-users-and-teams-manually.md)
- [Group and isolate cluster resources](group-resources.md)
- [Grant role-access to cluster resources](grant-permissions.md)
