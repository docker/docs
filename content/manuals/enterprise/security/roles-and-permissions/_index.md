---
title: Roles and permissions overview
linkTitle: Roles and permissions
description: Choose core or custom roles to control access to repositories, teams, and organization settings
keywords: roles, permissions, core roles, custom roles, member, editor, owner, access control, organization, docker hub, docker home, security
tags: [admin]
aliases:
  - /admin/organization/roles/
  - /security/for-admins/roles-and-permissions/
  - /docker-hub/roles-and-permissions/
grid:
  - title: "Core roles"
    description: Learn about Docker's built-in Member, Editor, and Owner roles with predefined permissions.
    icon: shield-check
    link: /enterprise/security/roles-and-permissions/core-roles/
  - title: "Custom roles"
    description: Create tailored permission sets that match your organization's specific needs.
    icon: adjustments-horizontal
    link: /enterprise/security/roles-and-permissions/custom-roles/
weight: 40
---

{{< summary-bar feature_name="General admin" >}}

Roles control what users can do in your Docker organization. When you invite
users or create teams, you assign roles that determine their permissions for
repositories, teams, and organization settings.

Docker provides two role types:

- [Core roles](/manuals/enterprise/security/roles-and-permissions/core-roles.md):
  Built-in Member, Editor, and Owner roles with predefined permissions
- [Custom roles](/manuals/enterprise/security/roles-and-permissions/custom-roles/_index.md):
  Permission sets you define for your organization's needs

## Core roles

Core roles are Docker's built-in roles:

- **Member**: Basic access. Members can view other organization members and
  pull images from repositories they have access to.
- **Editor**: Partial administrative access. Editors can create, edit, and
  delete repositories, and manage team permissions for repositories.
- **Owner**: Full administrative access. Owners can manage all organization
  settings, including repositories, teams, members, billing, and security
  features.

For a full permission comparison, see
[Core roles and permissions](/manuals/enterprise/security/roles-and-permissions/core-roles.md).

## Custom roles

Custom roles let you select specific permissions from categories such as
user management, team management, billing, and Hub. Use them when core roles
don't match your access control needs.

Custom roles require a Docker Business subscription and owner permissions.
For details, see
[Custom roles](/manuals/enterprise/security/roles-and-permissions/custom-roles/_index.md).

## When to use each role type

Use core roles when:

- Docker's predefined permission sets match your structure
- You want simple role assignments
- Your access control needs are standard

Use custom roles when:

- You need permission combinations not available in core roles
- You want specialized roles such as billing administrators or security
  auditors
- You need department-specific or least-privilege access control

## How roles work with team permissions

You can assign users and teams either a core role or a custom role, but not
both. Roles also work with team permissions:

1. **Role permissions**: Apply organization-wide. Custom roles can grant
   organization settings and repository management permissions.
1. **Team permissions**: Grant additional repository-specific access when
   users join teams. This is separate from role-based permissions.

This layered model lets you grant broad organizational access through roles
and specific repository access through team memberships.

## Next steps

Choose the role type that fits your organization:

{{< grid >}}
