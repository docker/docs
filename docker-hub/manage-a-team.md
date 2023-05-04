---
description: Docker Hub Teams & Organizations
keywords: Docker, docker, registry, teams, organizations, plans, Dockerfile, Docker Hub, docs, documentation
title: Create and manage a team
---

A team is a group of Docker users that belong to an organization. An
organization can have multiple teams. When you first create an organization,
you’ll see that you have a team (called 'Company') and the owners team, with a single member. An
organization owner can then create new teams and add members to an existing team
using their Docker ID or email address and by selecting a team the user should be part of.

The org owner can add additional org owners to the owners team to help them
manage users, teams, and repositories in the organization. See [Owners
team](#the-owners-team) for details.

## Create a team

1. Go to **Organizations** in Docker Hub, and select your organization.
2. Select the **Teams** tab and then select **Create Team**.
3. Fill out your team's information and select **Create**.
4. [Add members to your team](members.md#add-a-member-to-a-team)

## The owners team

The owners team is a special team created by default during the org creation
process. The owners team has full access to all repositories in the organization.

An organization owner is an administrator who is responsible to manage
repositories and add team members to the organization. They have full access to
private repositories, all teams, billing information, and org settings. An org
owner can also specify [permissions](#permissions-reference) for each team in
the organization. Only an org owner can enable [SSO](../single-sign-on/index.md)
for
the organization. When SSO is enabled for your organization, the org owner can
also manage users. Docker can auto-provision Docker IDs for new end-users or
users who'd like to have a separate Docker ID for company use through SSO
enforcement.

The org owner can also add additional org owners to help them manage users, teams, and repositories in the organization.

## Configure repository permissions for a team

Organization owners can configure repository permissions on a per-team basis.
For example, you can specify that all teams within an organization have "Read and
Write" access to repositories A and B, whereas only specific teams have "Admin"
access. Note that org owners have full administrative access to all repositories within the organization.

To give a team access to a repository

1. Navigate to **Organizations** in Docker Hub, and select your organization.
2. Select the **Teams** tab and select the team that you'd like to configure repository access to.
3. Select the **Permissions** tab and select a repository from the
   **Repository** drop-down.
4. Choose a permission from the **Permissions** dropdown list and select
   **Add**.

    ![Team Repo Permissions](images/team-repo-permission.png){:width="700px"}

### Permissions reference

- `Read-only` access lets users view, search, and pull a private repository in the same way as they can a public repository.
- `Read & Write` access lets users pull, push, and view a repository. In addition, it lets users view, cancel, retry or trigger builds
- `Admin` access lets users pull, push, view, edit, and delete a
  repository. You can also edit build settings, and update the repositories description, collaborators rights, public/private visibility, and delete.

Permissions are cumulative. For example, if you have "Read & Write" permissions,
you automatically have "Read-only" permissions:

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

> **Note**
>
> A user who hasn't verified their email address only has
> `Read-only` access to the repository, regardless of the rights their team
> membership has given them.

## View a team's permissions for all repositories

To view a team's permissions across all repositories:

1. Open **Organizations** > **_Your Organization_** > **Teams** > **_Team Name_**.
2. Select the **Permissions** tab, where you can view the repositories this team can access.

## Videos

You can also check out the following videos for information about creating Teams
and Organizations in Docker Hub.

- [Overview of organizations](https://www.youtube-nocookie.com/embed/G7lvSnAqed8){: target="_blank" rel="noopener" class="_"}
- [Create an organization](https://www.youtube-nocookie.com/embed/b0TKcIqa9Po){: target="_blank" rel="noopener" class="_"}
- [Working with Teams](https://www.youtube-nocookie.com/embed/MROKmtmWCVI){: target="_blank" rel="noopener" class="_"}
- [Create Teams](https://www.youtube-nocookie.com/embed/78wbbBoasIc){: target="_blank" rel="noopener" class="_"}