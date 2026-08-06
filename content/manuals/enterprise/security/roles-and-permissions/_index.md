---
title: Roles and permissions overview
linkTitle: Roles and permissions
description: Choose between core and custom roles to control access to repositories, teams, and organization settings
keywords: roles, permissions, core roles, custom roles, member, editor, owner, access control, organization, docker hub, docker home, security
tags: [admin]
aliases:
  - /admin/organization/roles/
  - /security/for-admins/roles-and-permissions/
  - /docker-hub/roles-and-permissions/
grid:
  - title: Core roles
    description: Compare the permissions granted by the built-in Member, Editor, and Owner roles.
    icon: shield-check
    link: /enterprise/security/roles-and-permissions/core-roles/
  - title: Custom roles
    description: Build permission sets that match your organization's access control needs.
    icon: adjustments-horizontal
    link: /enterprise/security/roles-and-permissions/custom-roles/
  - title: Custom roles permissions
    description: Reference documentation for every permission you can assign to a custom role.
    icon: list-bullet
    link: /enterprise/security/roles-and-permissions/custom-roles/permissions-reference/
weight: 40
---

{{< summary-bar feature_name="General admin" >}}

Roles determine what members can do in your Docker organization. When you
invite a user or create a team, you assign a role that grants permissions
for repositories, teams, and organization settings.

Docker provides two role types. Users and teams get either a core role or a
custom role, but not both.

## Core roles

Core roles are Docker's built-in Member, Editor, and Owner roles. Their
permissions are predefined and can't be changed. Use core roles when
Docker's permission sets match how your organization works and you want
straightforward role assignments.

## Custom roles

Custom roles are permission sets you build by selecting individual
permissions, such as billing or team management. Use custom roles when you
need a combination that core roles don't offer, for example a billing
administrator or a security auditor, or when you want to grant
least-privilege access. Custom roles require a Docker Business subscription.

## Roles and team permissions

Roles apply organization-wide, and team permissions apply to specific
repositories. The two systems work together: a user's role sets their
organization-wide access, and team membership can extend their access to
individual repositories.

## Next steps

{{< grid >}}
