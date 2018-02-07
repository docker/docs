---
description: Learn how to create and manage users in your Docker Universal Control
  Plane cluster.
keywords: authorize, authentication, users, teams, UCP, Docker
redirect_from:
- /ucp/user-management/create-and-manage-users/
title: Create and manage users
---

When using the UCP built-in authentication, you need to create users and
assign them with a default permission level so that they can access the
cluster.

To create a new user, go to the **UCP web UI**, and navigate to the
**Users & Teams** page.

![](../images/create-users-1.png)

Click the **Create User** button, and fill-in the user information.

![](../images/create-users-2.png)

Check the 'Is a UCP admin' option, if you want to grant permissions for the
user to change cluster configurations. Also, assign the user with a default
permission level.

Default permissions specify the resources a user has access to create or edit
in the cluster. There are four permission levels:

| Default permission level | Description                                                                                                                                                                                                  |
|:-------------------------|:-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `No Access`              | The user can't view any resource, like volumes, networks, images, or containers.                                                                                                                             |
| `View Only`              | The user can view volumes, networks and images, but can't create any containers.                                                                                                                             |
| `Restricted Control`     | The user can view and edit volumes, networks, and images. They can create containers, but can't see other users containers, run `docker exec`, or run containers that require privileged access to the host. |
| `Full Control`           | The user can view and edit volumes, networks, and images, They can create containers without any restriction, but can't see other users containers.                                                          |

[Learn more about the UCP permission levels](permission-levels.md). Finally,
click the **Create User** button, to create the user.

## Where to go next

* [Create and manage teams](create-and-manage-teams.md)
* [UCP permission levels](permission-levels.md)