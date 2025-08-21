---
title: Create and manage a team
weight: 40
description: Learn how to create and manage teams for your organization
keywords: docker, registry, teams, organizations, plans, Dockerfile, Docker
  Hub, docs, documentation, repository permissions, configure repository access, team management
aliases:
- /docker-hub/manage-a-team/
---

{{< summary-bar feature_name="Admin orgs" >}}

You can create teams for your organization in the Admin Console or Docker Hub,
and configure team repository access in Docker Hub.

A team is a group of Docker users that belong to an organization. An
organization can have multiple teams. An organization owner can create new
teams and add members to an existing team using their Docker ID or email
address. Members aren't required to be part of a team to be associated with an
organization.

The organization owner can add additional organization owners to help them
manage users, teams, and repositories in the organization by assigning them
the owner role.

## What is an organization owner?

An organization owner is an administrator who has the following permissions:

- Manage repositories and add team members to the organization
- Access private repositories, all teams, billing information, and
organization settings
- Specify [permissions](#permissions-reference) for each team in the
organization
- Enable [SSO](/manuals/enterprise/security/single-sign-on/_index.md) for the
organization

When SSO is enabled for your organization, the organization owner can
also manage users. Docker can auto-provision Docker IDs for new end-users or
users who'd like to have a separate Docker ID for company use through SSO
enforcement.

Organization owners can add others with the owner role to help them
manage users, teams, and repositories in the organization.

For more information on roles, see
[Roles and permissions](/manuals/enterprise/security/roles-and-permissions.md).

## Create a team

{{< tabs >}}
{{< tab name="Admin Console" >}}

1. Sign in to [Docker Home](https://app.docker.com) and select your
organization.
1. Select **Teams**.

{{< /tab >}}
{{< tab name="Docker Hub" >}}

{{% include "hub-org-management.md" %}}

1. Sign in to [Docker Hub](https://hub.docker.com).
1. Select **My Hub** and choose your organization.
1. Select the **Teams** and then select **Create Team**.
1. Fill out your team's information and select **Create**.
1. [Add members to your team](members.md#add-a-member-to-a-team).

{{< /tab >}}
{{< /tabs >}}

## Set team repository permissions

You must create a team before you are able to configure repository permissions.
For more details, see [Create and manage a
team](/manuals/admin/organization/manage-a-team.md).

To set team repository permissions:

1. Sign in to [Docker Hub](https://hub.docker.com).
1. Select **My Hub** > **Repositories**.

   A list of your repositories appears.

1. Select a repository.

   The **General** page for the repository appears.

1. Select the **Permissions** tab.
1. Add, modify, or remove a team's repository permissions.

   - Add: Specify the **Team**, select the **Permission**, and then select **Add**.
   - Modify: Specify the new permission next to the team.
   - Remove: Select the **Remove permission** icon next to the team.

### Permissions reference

- `Read-only` access lets users view, search, and pull a private repository
in the same way as they can a public repository.
- `Read & Write` access lets users pull, push, and view a repository. In
addition, it lets users view, cancel, retry or trigger builds.
- `Admin` access lets users pull, push, view, edit, and delete a
  repository. You can also edit build settings and update the repository’s
  description, collaborator permissions, public/private visibility, and delete.

Permissions are cumulative. For example, if you have "Read & Write" permissions,
you automatically have "Read-only" permissions.

The following table shows what each permission level allows users to do:

| Action | Read-only | Read & Write | Admin |
|:------------------:|:---------:|:------------:|:-----:|
| Pull a Repository | ✅ | ✅ | ✅ |
| View a Repository | ✅ | ✅ | ✅ |
| Push a Repository | ❌ | ✅ | ✅ |
| Edit a Repository | ❌ | ❌ | ✅ |
| Delete a Repository | ❌ | ❌ | ✅ |
| Update a Repository Description | ❌ | ❌ | ✅ |
| View Builds | ✅ | ✅ | ✅ |
| Cancel Builds | ❌ | ✅ | ✅ |
| Retry Builds | ❌ | ✅ | ✅ |
| Trigger Builds | ❌ | ✅ | ✅ |
| Edit Build Settings | ❌ | ❌ | ✅ |

> [!NOTE]
>
> A user who hasn't verified their email address only has `Read-only` access to
the repository, regardless of the rights their team membership has given them.

## View team permissions for all repositories

To view a team's permissions across all repositories:

1. Sign in to [Docker Hub](https://hub.docker.com).
1. Select **My Hub** and choose your organization.
1. Select **Teams** and choose your team name.
1. Select the **Permissions** tab, where you can view the repositories this
team can access.

## Delete a team

Organization owners can delete a team. When you remove a team from your
organization, this action revokes member access to the team's permitted
resources. It won't remove users from other teams that they belong to, and it
won't delete any resources.

{{< tabs >}}
{{< tab name="Admin Console" >}}

1. Sign in to [Docker Home](https://app.docker.com/) and select your
organization.
1. Select **Teams**.
1. Select the **Actions** icon next to the name of the team you want to delete.
1. Select **Delete team**.
1. Review the confirmation message, then select **Delete**.

{{< /tab >}}
{{< tab name="Docker Hub" >}}

{{% include "hub-org-management.md" %}}

1. Sign in to [Docker Hub](https://hub.docker.com).
1. Select **My Hub** and choose your organization.
1. Select **Teams**.
1. Select the name of the team that you want to delete.
1. Select **Settings**.
1. Select **Delete Team**.
1. Review the confirmation message, then select **Delete**.

{{< /tab >}}
{{< /tabs >}}

## More resources

- [Video: Docker teams](https://youtu.be/WKlT1O-4Du8?feature=shared&t=348)
- [Video: Roles, teams, and repositories](https://youtu.be/WKlT1O-4Du8?feature=shared&t=435)
