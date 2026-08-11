---
title: Manage custom roles in Docker
linkTitle: Manage
description: Create, edit, assign, reassign, and delete custom roles in Docker Home for organization users and teams
keywords: custom roles, manage custom roles, role assignments, access control, Docker Home, organization roles, permissions
weight: 10
---

{{< summary-bar feature_name="General admin" >}}

Create custom roles, manage their permissions, and assign them to users and
teams. For a full list of permissions, see the
[custom roles permissions reference](permissions-reference.md).

## Create a custom role

Before you can assign a custom role, create one:

1. Sign in to [Docker Home](https://app.docker.com).
1. Select **Roles**, then **Create role**.
1. Define the role:
   - Provide a **Label**
   - Enter a unique **Name** identifier. You can't change it later.
   - Add an optional **Description**
1. Set permissions for the role by expanding permission categories and
   selecting the checkboxes for permissions. For a full list of available
   permissions, see the
   [custom roles permissions reference](permissions-reference.md).
1. Select **Review** to review the configuration and selected permissions.
1. Select **Create**.

After you create a custom role, you can
[assign it to users or teams](#assign-custom-roles).

## Edit a custom role

1. Sign in to [Docker Home](https://app.docker.com).
1. Select **Roles**.
1. Find your custom role in the list, then select the **Actions** menu.
1. Select **Edit**.
1. Edit any of the following settings:
   - Label
   - Description
   - Permissions
1. Select **Save**.

## Assign custom roles

{{< tabs >}}
{{< tab name="Individual users" >}}

1. Sign in to [Docker Home](https://app.docker.com).
1. Select **Members**.
1. Locate the member you want to assign a custom role to, then select the
   **Actions** menu.
1. Select **Change role**.
1. In the **Select a role** drop-down, select your custom role.
1. Select **Save**.

{{< /tab >}}
{{< tab name="Bulk users" >}}

1. Sign in to [Docker Home](https://app.docker.com).
1. Select **Members**.
1. Use the checkboxes in the username column to select the users you want
   to assign a custom role to.
1. Select **Change role**.
1. In the **Select a role** drop-down, select your custom role or a core
   role.
1. Select **Save**.

{{< /tab >}}
{{< tab name="Teams" >}}

1. Sign in to [Docker Home](https://app.docker.com).
1. Select **Teams**.
1. Locate the team you want to assign a custom role to, then select the
   **Actions** menu.
1. Select **Assign role**.
1. Select your custom role, then select **Assign**.

The role column updates to the newly assigned role.

{{< /tab >}}
{{< /tabs >}}

## View role assignments

To see which users and teams are assigned to roles:

1. Sign in to [Docker Home](https://app.docker.com).
1. Select **Roles**.
1. In the roles list, view the **Users** and **Teams** columns for
   assignment counts.
1. Select a role to view its permissions and assignments in detail.

## Reassign custom roles

{{< tabs >}}
{{< tab name="Individual users" >}}

1. Sign in to [Docker Home](https://app.docker.com).
1. Select **Members**.
1. Locate the member you want to reassign, then select the **Actions**
   menu.
1. Select **Change role**.
1. In the **Select a role** drop-down, select the new role.
1. Select **Save**.

{{< /tab >}}
{{< tab name="Bulk users" >}}

1. Sign in to [Docker Home](https://app.docker.com).
1. Select **Members**.
1. Use the checkboxes in the username column to select the users you want
   to reassign.
1. Select **Change role**.
1. In the **Select a role** drop-down, select the new role.
1. Select **Save**.

{{< /tab >}}
{{< tab name="Teams" >}}

1. Sign in to [Docker Home](https://app.docker.com).
1. Select **Teams**.
1. Locate the team, then select the **Actions** menu.
1. Select **Change role**.
1. In the pop-up window, select a role from the drop-down, then select
   **Save**.

{{< /tab >}}
{{< /tabs >}}

## Delete a custom role

> [!IMPORTANT]
>
> Before you delete a custom role, reassign every user and team that uses
> it to a different role.

1. Sign in to [Docker Home](https://app.docker.com).
1. Select **Roles**.
1. Find your custom role in the list, then select the **Actions** menu.
1. If the role has assigned users or teams:
   - On the **Members** page, change the role for every user assigned to
     this custom role
   - On the **Teams** page, reassign every team that has this custom role
1. When no users or teams are assigned, return to **Roles**.
1. Find your custom role and select the **Actions** menu.
1. Select **Delete**.
1. In the confirmation window, select **Delete** to confirm.

## Next steps

- [Custom roles permissions reference](permissions-reference.md): Review
  permissions you can grant to a custom role
- [Core roles and permissions](/manuals/enterprise/security/roles-and-permissions/core-roles.md):
  Compare built-in Member, Editor, and Owner permissions
- [Manage organization members](/manuals/admin/organization/manage/members.md):
  Invite and manage users in your organization
