---
description: Learn about the permission levels available in Docker Trusted Registry.
keywords: docker, registry, security, permissions
title: Permission levels
---

Docker Trusted Registry allows you to define fine-grain permissions over image
repositories.

## Administrator users

Users are shared across Docker Datacenter. When you create a new user in Docker
Universal Control Plane, that user becomes available in DTR and vice versa.
When you create an administrator user in DTR, that user is a Docker Datacenter
administrator, with permissions to:

* Manage users across Docker Datacenter,
* Manage DTR repositories and settings,
* Manage the whole UCP cluster.

## Team permission levels

Teams allow you to define the permissions a set of user has for a set of
repositories. Three permission levels are available:

| Repository operation  | read | read-write | admin |
|:----------------------|:----:|:----------:|:-----:|
| View/ browse          |  x   |     x      |   x   |
| Pull                  |  x   |     x      |   x   |
| Push                  |      |     x      |   x   |
| Delete tags           |      |     x      |   x   |
| Edit description      |      |            |   x   |
| Set public or private |      |            |   x   |
| Manage user access    |      |            |   x   |
| Delete repository     |      |            |       |

Team permissions are additive. When a user is a member of multiple teams, they
have the highest permission level defined by those teams.

## Overall permissions

Here's an overview of the permission levels available in DTR:

* Anonymous users: Can search and pull public repositories.
* Users: Can search and pull public repos, and create and manage their own
repositories.
* Team member: Everything a user can do, plus the permissions granted by the teams the user is member of.
* Team admin: Everything a team member can do, and can also add members to the team.
* Organization admin: Everything a team admin can do, can create new teams, and add members to the organization.
* DDC admin: Can manage anything across UCP and DTR.

## Where to go next

* [Authentication and authorization](index.md)
