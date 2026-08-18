---
title: Custom roles and permissions overview
linkTitle: Custom roles
description: >-
  Use custom roles to assign tailored permissions to users and teams in
  your Docker organization
keywords: >-
  custom roles, custom permissions, permission sets, access control,
  least privilege, Docker Business, Docker Home, organization
  management, role assignment, teams, AI Governance, Docker Offload,
  security
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

{{< summary-bar feature_name="Custom roles" >}}

Custom roles are permission sets built from individual permissions, so you
can grant only the access a user or team needs. 

If Docker's predefined
permission sets meet your needs, use
[core roles](/manuals/enterprise/security/roles-and-permissions/core-roles.md)
instead.

## Prerequisites

- A Docker Business subscription
- Owner permissions in your Docker organization

## Custom roles

To create a custom role, you select permissions from organization management,
Docker Hub, billing, AI Governance, Docker Hardened Images, and Docker
Offload. You then assign custom roles you created to individual users or to teams.

Users and teams get either a core role or a custom role, but not both.

## Using custom roles

Use custom roles when you need:

- Specialized roles such as billing administrators, security auditors, or
  repository managers
- Department-specific access control
- Least-privilege access with precise permission grants

## Next steps

{{< grid >}}
