---
title: Custom roles overview
linkTitle: Custom roles
description: Learn when to use custom roles and how they differ from Docker's built-in core roles
keywords: custom roles, permissions, access control, least privilege, docker business, organization management, docker hub, docker home, security
weight: 20
grid:
  - title: "Manage custom roles"
    description: Create, edit, assign, and delete custom roles for users and teams.
    icon: adjustments-horizontal
    link: /enterprise/security/roles-and-permissions/custom-roles/manage/
  - title: "Permissions reference"
    description: Review all permissions you can assign when building a custom role.
    icon: list-bullet
    link: /enterprise/security/roles-and-permissions/custom-roles/permissions-reference/
---

{{< summary-bar feature_name="General admin" >}}

Custom roles let you build permission sets that match your organization's
access control needs. Use them when Docker's
[core roles](/manuals/enterprise/security/roles-and-permissions/core-roles.md)
don't provide the right combination of permissions.

## What are custom roles?

With custom roles, you select permissions from categories such as user
management, team management, billing, Hub, and governance. You can assign a
custom role to individual users or to teams.

Users and teams get either a core role or a custom role, but not both.

## Prerequisites

To configure custom roles, you need:

- A Docker Business subscription
- Owner permissions in your Docker organization

## When to use custom roles

Use custom roles when:

- You need permission combinations not available in core roles
- You want specialized roles such as billing administrators, security
  auditors, or repository managers
- You need department-specific access control
- You want least-privilege access with precise permission grants

## Next steps

{{< grid >}}
