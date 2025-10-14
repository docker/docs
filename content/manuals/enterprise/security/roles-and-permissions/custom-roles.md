---
title: Custom roles
description:
keywords:
---

{{< summary-bar feature_name="General admin" >}}

Custom roles allow you to create tailored permission sets that match your
organization's specific needs. This page covers custom roles, and steps
to create and manage them.

## What are custom roles?

Custom roles let you create tailored permission sets for your organization. You
can assign custom roles to individual users or teams.
Users get either a core role or custom role, but not both.

Use custom roles when Docker's default roles don't fit your needs.

## Prerequisites

To configure custom roles, you need owner permissions in your Docker
organization.

## Create a custom role

Before you can assign a custom role to users, you must create one in the
Admin Console:

1. Sign in to [Docker Home](https://app.docker.com).
1. Select **Admin Console**, then **User management**.
1. Select **Roles**, then **Create role**.
1. Create a name and describe what the role is for:
    - Provide a **Display name**
    - Enter a unique **Name** identifier (can't be changed later)
    - Add an optional **Description**
1. Set permissions for the role by expanding permission categories and selecting
the checkboxes for permissions. For a full list of available permissions, see
the [custom roles permissions reference](#custom-roles-permissions-reference).
1. Select **Review** to review your custom roles configruation and see a summary
of selected permissions.
1. Select **Create**.

With a custom role created, you can now [assign custom roles to users](#assign-custom-roles-to-users).

## Edit a custom role

1. Sign in to [Docker Home](https://app.docker.com).
1. Select **Admin Console**, then **User management**.
1. Select **Roles**.
1. Find your custom role from the list, and select the **Actions menu**.
1. Select **Edit**.
1. You can edit the following custom role settings:
    - Display name
    - Description
    - Permissions
1. After you have finished editing, select **Save**.

## Assign custom roles

{{< tabs >}}
{{< tab name="Individual users" >}}

1. Sign in to [Docker Home](https://app.docker.com).
1. Select **Members**.
1. Locate the member you want to assign a custom role to, then select the
**Actions menu**.
1. In the drop-down, select **Change role**.
1. In the **Select a role** drop-down, select your custom role.
1. Select **Save**.

{{< /tab >}}
{{< tab name="Bulk users" >}}

1. Sign in to [Docker Home](https://app.docker.com).
1. Select **Members**.
1. Use the checkboxes in the username column to select all users you want
to assign a custom role to.
1. Select **Change role**.
1. In the **Select a role** drop-down, select your custom role.
1. Select **Save**.

{{< /tab >}}
{{< tab name="Teams" >}}

1. Sign in to [Docker Home](https://app.docker.com).
1. Select **Teams**.
1. Locate the team you want to assign a custom role to, then select
the **Actions menu**.
1. Select **Assign role**.
1. Select your custom role, then select **Assign**.

The role column will update to the newly assigned role.

{{< /tab >}}
{{< /tabs >}}

## View role assignments

To see which users and teams are assigned to roles:

1. Sign in to [Docker Home](https://app.docker.com).
1. Select **Admin Console**, then **User management**.
1. Select **Roles**.
1. In the roles list, view the **Users** and **Teams** columns to see
assignment counts.
1. Select a specific role to view its permissions adn assignments in detail.

## Reassign custom roles

{{< tabs >}}
{{< tab name="Individual users" >}}

1. Sign in to [Docker Home](https://app.docker.com).
1. Select **Admin Console**, then **User management**.
1. Select **Roles**.
1. Find your custom role from the list, and select the **Actions menu**.
1. Select **Reassign**.
1. On the reassignment page, **Select a role** to reassign, then select **Save**.

{{< /tab >}}
{{< tab name="Bulk users" >}}


{{< /tab >}}
{{< tab name="Teams" >}}

1. Sign in to [Docker Home](https://app.docker.com).
1. Select **Teams**.
1. Locate the team, then select the **Actions menu**.
1. Select **Change role**.
1. In the pop-up window, select a role from the drop-down menu, then
select **Save**.

{{< /tab >}}
{{< /tabs >}}

## Duplicate a custom role

1. Sign in to [Docker Home](https://app.docker.com).
1. Select **Admin Console**, then **User management**.
1. Select **Roles**.
1. Find your custom role from the list, and select the **Actions menu**.
1. Select **Duplicate**.
1. Modify the duplicated role's name, description, and permissions as needed.
1. Select **Create** to save the new role.

## Delete a custom role

If you have users or teams assigned to a role, you must reassign them to new roles before deleting.

1. Sign in to [Docker Home](https://app.docker.com).
1. Select **Admin Console**, then **User management**.
1. Select **Roles**.
1. Find your custom role from the list, and select the **Actions menu**.
1. If the role has assigned users or teams, select **Reassign** first to move
them to different roles.
1. Once no users or teams are assigned, select the **Actions menu** again.
1. Select **Delete**.
1. In the confirmation window, select **Delete** to confirm.

## Custom roles permissions reference

Custom roles can included any combination of the following permissions.

### User and role management permissions

- **Invite members**: Send organization invitations
- **Manage members**: Remove users from the organizatino
- **Manage member roles**: Assign roles to users
- **Create custom roles**: Create, edit, and delete custom roles
- **View member activity**: View activity logs in the organization
- **Export and reporting**: Export users and activity logs

### Team management permissions

- **Create teams**:
- **Manage teams**:

### Organization configuration permissions

### Billing permissions

### Hub permissions