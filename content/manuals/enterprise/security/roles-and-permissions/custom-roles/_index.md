---
title: Custom roles and permissions overview
linkTitle: Custom roles
description: Use custom roles to assign tailored permissions to users and teams in your Docker organization
keywords: custom roles, custom permissions, permission sets, access control, least privilege, docker business, docker home, organization management, role assignment, teams, security
weight: 20
grid:
  - title: Manage custom roles
    description: Create, edit, assign, and delete custom roles for users and teams.
    icon: adjustments-horizontal
    link: /enterprise/security/roles-and-permissions/custom-roles/manage/
  - title: Permissions reference
    description: Review every permission you can assign when building a custom role.
    icon: list-bullet
    link: /enterprise/security/roles-and-permissions/custom-roles/permissions-reference/
---

{{< summary-bar feature_name="General admin" >}}

Custom roles are permission sets built from individual permissions. This page
defines custom roles and explains when to use them.
If Docker's predefined permission sets meet your needs, use
[core roles](/manuals/enterprise/security/roles-and-permissions/core-roles.md)
instead.

## Custom roles

With custom roles, you select permissions from categories such as
organization management, Docker Hub, billing, and governance. You can assign
a custom role to individual users or to teams. Users and teams get either a
core role or a custom role, but not both.

Before configuring custom roles, you need:

- A Docker Business subscription
- Owner permissions in your Docker organization

## Using custom roles

Use custom roles when you need:

- Specialized roles such as billing administrators, security auditors, or
  repository managers
- Department-specific access control
- Least-privilege access with precise permission grants

## Next steps

{{< grid >}}
