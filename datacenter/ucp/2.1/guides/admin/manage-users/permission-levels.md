---
description: Learn about the permission levels available in Docker Universal
  Control Plane.
keywords: authorization, authentication, users, teams, UCP
title: Permission levels
---

Docker Universal Control Plane has two types of users: administrators and
regular users. Administrators can make changes to the UCP cluster, while
regular users have permissions that range from no access to full control over
volumes, networks, images, and containers.

## Administrator users

In Docker UCP, only users with administrator privileges can make changes to
cluster settings. This includes:

* Managing user and team permissions,
* Managing cluster configurations like adding and removing nodes to the cluster.

## Default permission levels

Regular users can't change cluster settings, and they are assigned with a
default permission level.

The default permission level specify the resources a user access or edit. You
can choose from four permission levels that range from no
access to full control over the resources.

| Default permission level | Description                                                                                                                                                                                                  |
|:-------------------------|:-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `No Access`              | The user can't view any resource, like volumes, networks, images, or containers.                                                                                                                             |
| `View Only`              | The user can view volumes, networks, and images, but can't create any containers.                                                                                                                             |
| `Restricted Control`     | The user can view and edit volumes, networks, and images. They can create containers, but can't see other users' containers, run `docker exec`, or run containers that require privileged access to the host. |
| `Full Control`           | The user can view and edit volumes, networks, and images, They can create containers without any restriction, but can't see other users' containers.                                                          |

If a user has Restricted Control or Full Control default permissions, they can create resources without labels, and only the user and Admins can see and access the resources. Default permissions also affect ability for a user to access things that can't have labels, images and nodes.  

## Team permission levels

Teams and labels give the administrator fine-grained control over permissions. Each team can have multiple labels. Each label has a key of `com.docker.ucp.access.label`. The label is then applied to the containers, services, networks, secrets, and volumes. Labels are not currently available for nodes and images. DTR has its own permissions.

There are four permission levels:

| Team permission level | Description                                                                                                                                          |
|:----------------------|:-----------------------------------------------------------------------------------------------------------------------------------------------------|
| `No Access`           | The user can't view containers with this label.                                                                                                      |
| `View Only`           | The user can view but can't create containers with this label.                                                                                       |
| `Restricted Control`  | The user can view and create containers with this label. The user can't run `docker exec`, or containers that require privileged access to the host. |
| `Full Control`        | The user can view and create containers with this label, without any restriction.                                                                    |

## Where to go next

* [Create and manage users](create-and-manage-users.md)
* [Create and manage teams](create-and-manage-teams.md)
* [Docker Reference Architecture: Securing Docker Datacenter and Security Best Practices](https://success.docker.com/KBase/Docker_Reference_Architecture%3A_Securing_Docker_Datacenter_and_Security_Best_Practices)
