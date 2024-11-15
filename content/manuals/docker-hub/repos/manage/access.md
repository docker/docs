---
description: Discover how to manage access to repositories on Docker Hub.
keywords: Docker Hub, Hub, repository access, repository collaborators, repository privacy
title: Access management
LinkTItle: Access
weight: 50
aliases:
- /docker-hub/repos/access/
---

In this topic learn about the features available to manage access to your
repositories. This includes visibility, collaborators, roles, teams, and
organization access tokens.

## Repository visibility

The most basic repository access is controlled via the visibility. A
repository's visibility can be public or private.

With public visibility, the repository appears in Docker Hub search results and
can be pulled by everyone. To manage push access to public personal
repositories, you can use collaborators. To manage push access to public
organization repositories, you can use roles, teams, or organization access
tokens.

With private visibility, the repository doesn't appear in Docker Hub search
results and is only accessible to those with granted permission. To manage push
and pull access to private personal repositories, you can use collaborators. To
manage push and pull access to private organization repositories, you can use
roles, teams, or organization access tokens.

### Change repository visibility

When creating a repository in Docker Hub, you can set the repository visibility.
In addition, you can set the default repository visibility when a repository is
created in your personal repository settings. The following describes how to
change the visibility after the repository has been created.

To change repository visibility:

1. Sign in to [Docker Hub](https://hub.docker.com).
2. Select **Repositories**.
3. Select a repository.

   The **General** page for the repository appears.

4. Select the **Settings** tab.
5. Under **Visibility settings**, select one of the following:

   - **Make public**: The repository appears in Docker Hub search results and can be
     pulled by everyone.
   - **Make private**: The repository doesn't appear in Docker Hub search results
     and is only accessible to you and collaborators. In addition, if the
     repository is in an organization's namespace, then the repository
     is accessible to those with applicable roles or permissions.

6. Type the repository's name to verify the change.
7. Select **Make public** or **Make private**.

## Collaborators

A collaborator is someone you want to give `push` and `pull` access to a
personal repository. Collaborators aren't able to perform any administrative
tasks such as deleting the repository or changing its visibility from private to
public. In addition, collaborators can't add other collaborators.

Only personal repositories can use collaborators. You can add unlimited
collaborators to public repositories, and Docker Pro accounts can add up to 1
collaborator on private repositories.

Organization repositories can't use collaborators, but can use member roles,
teams, or organization access tokens to manage access.

### Manage collaborators

1. Sign in to [Docker Hub](https://hub.docker.com).

2. Select **Repositories**.

   A list of your repositories appears.

3. Select a repository.

   The **General** page for the repository appears.

4. Select the **Collaborators** tab.

5. Add or remove collaborators based on their Docker username.

You can choose collaborators and manage their access to a private
repository from that repository's **Settings** page.

## Organization roles

Organizations can use roles for individuals, giving them different
permissions in the organization. For more details, see [Roles and
permissions](/manuals/security/for-admins/roles-and-permissions.md).

## Organization teams

Organizations can use teams. A team can be assigned fine-grained repository
access.

### Configure team repository permissions

You must create a team before you are able to configure repository permissions.
For more details, see [Create and manage a
team](/manuals/admin/organization/manage-a-team.md).

To configure team repository permissions:

1. Sign in to [Docker Hub](https://hub.docker.com).

2. Select **Repositories**.

   A list of your repositories appears.

3. Select a repository.

   The **General** page for the repository appears.

4. Select the **Permissions** tab.

5. Add, modify, or remove a team's repository permissions.

   - Add: Specify the **Team**, select the **Permission**, and then select **Add**.
   - Modify: Specify the new permission next to the team.
   - Remove: Select the **Remove permission** icon next to the team.

## Organization access tokens (OATs)

Organizations can use OATs. OATs let you assign fine-grained repository access
permissions to tokens. For more details, see [Organization access
tokens](/manuals/security/for-admins/access-tokens.md).
