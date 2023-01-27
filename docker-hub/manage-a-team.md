---
description: Docker Hub Teams & Organizations
keywords: Docker, docker, registry, teams, organizations, plans, Dockerfile, Docker Hub, docs, documentation
title: Teams and Organizations
---

## Create a team

A **Team** is a group of Docker users that belong to an organization. An
organization can have multiple teams. When you first create an organization,
you’ll see that you have a team, the **owners** team, with a single member. An
organization owner can then create new teams and add members to an existing team
using their Docker ID or email address and by selecting a team the user should be part of.

The org owner can add additional org owners to the owners team to help them
manage users, teams, and repositories in the organization. See [Owners
team](#the-owners-team) for details.

To create a team:

1. Go to **Organizations** in Docker Hub, and select your organization.
2. Open the **Teams** tab and click **Create Team**.
3. Fill out your team's information and click **Create**.

### The owners team

The **owners** team is a special team created by default during the org creation
process. The owners team has full access to all repositories in the organization.

An organization owner is an administrator who is responsible to manage
repositories and add team members to the organization. They have full access to
private repositories, all teams, billing information, and org settings. An org
owner can also specify [permissions](../docker-hub/repos/configure/index.md#permissions-reference) for each team in
the organization. Only an org owner can enable [SSO](../single-sign-on/index.md)
for
the organization. When SSO is enabled for your organization, the org owner can
also manage users. Docker can auto-provision Docker IDs for new end-users or
users who'd like to have a separate Docker ID for company use through SSO
enforcement.

The org owner can also add additional org owners to help them manage users, teams, and repositories in the organization.

## Configure repository permissions

Organization owners can configure repository permissions on a per-team basis.
For example, you can specify that all teams within an organization have Read and
Write access to repositories A and B, whereas only specific teams have Admin
access. Note that org owners have full administrative access to all repositories within the organization.

To give a team access to a repository

1. Navigate to **Organizations** in Docker Hub, and select your organization.
2. Click on the **Teams** tab and select the team that you'd like to configure  repository access to.
3. Click on the **Permissions** tab and select a repository from the
   **Repository** drop-down.
4. Choose a permission from the **Permissions** drop-down list and click
   **Add**.

    ![Team Repo Permissions](images/team-repo-permission.png){:width="700px"}

### View a team's permissions for all repositories

To view a team's permissions across all repositories:

1. Open **Organizations** > **_Your Organization_** > **Teams** > **_Team Name_**.
2. Click on the **Permissions** tab, where you can view the repositories this team can access.

## Videos

You can also check out the following videos for information about creating Teams
and Organizations in Docker Hub.

- [Overview of organizations](https://www.youtube-nocookie.com/embed/G7lvSnAqed8){: target="_blank" rel="noopener" class="_"}
- [Create an organization](https://www.youtube-nocookie.com/embed/b0TKcIqa9Po){: target="_blank" rel="noopener" class="_"}
- [Working with Teams](https://www.youtube-nocookie.com/embed/MROKmtmWCVI){: target="_blank" rel="noopener" class="_"}
- [Create Teams](https://www.youtube-nocookie.com/embed/78wbbBoasIc){: target="_blank" rel="noopener" class="_"}