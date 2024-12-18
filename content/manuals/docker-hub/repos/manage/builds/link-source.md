---
description: Link to GitHub and BitBucket
keywords: Docker, docker, registry, accounts, plans, Dockerfile, Docker Hub, trusted,
  builds, trusted builds, automated builds, GitHub
title: Configure automated builds from GitHub and BitBucket
linkTitle: Link accounts
weight: 20
aliases:
- /docker-hub/github/
- /docker-hub/bitbucket/
- /docker-hub/builds/link-source/
---

> [!NOTE]
>
> Automated builds require a Docker Pro, Team, or Business subscription.


To automate building and testing of your images, you link to your hosted source
code service to Docker Hub so that it can access your source code
repositories. You can configure this link for user accounts or
organizations.

If you are linking a source code provider to create autobuilds for a team, follow the instructions to [create a service account](index.md#service-users-for-team-autobuilds) for the team before linking the account as described below.

## Link to a GitHub user account

1. Sign in to Docker Hub.

2. Select the **Settings** icon in the top-right navigation, then select **Repository Settings**.

    ![Repository settings menu](./images/docker-hub-repo-settings-menu.png)


3. From the **Linked accounts** tab, select **Link provider** for the source provider you want to link.

    If you want to unlink your current GitHub account and relink to a new GitHub account, make sure to completely sign out of [GitHub](https://github.com/) before linking via Docker Hub.


4. Review the settings for the **Docker Hub Builder** OAuth application.

    ![Granting access to GitHub account](images/authorize-builder.png)

    > [!NOTE]
    >
    > If you are the owner of any GitHub organizations, you might see
    options to grant Docker Hub access to them from this screen. You can also
    individually edit an organization's third-party access settings to grant or
    revoke Docker Hub's access. See
    [Grant access to a GitHub organization](link-source.md#grant-access-to-a-github-organization)
    to learn more.

5. Select **Authorize docker** to save the link.

### Grant access to a GitHub organization

If you are the owner of a GitHub organization, you can grant or revoke Docker
Hub's access to the organization's repositories. Depending on the GitHub
organization's settings, you may need to be an organization owner.

If the organization has not had specific access granted or revoked before, you
can often grant access at the same time as you link your user account. In this
case, a **Grant access** button appears next to the organization name in the
link accounts screen, as shown below. If this button does not appear, you must
manually grant the application's access.

To manually grant Docker Hub access to a GitHub organization:

1. Link your user account using the instructions above.

2. From your GitHub Account settings, locate the **Organization settings**
section at the lower left.

3. Select the organization you want to give Docker Hub access to.

4. Select **Third-party access**.

    The page displays a list of third party applications and their access
    status.

5. Select the pencil icon next to **Docker Hub Builder**.

6. Select **Grant access** next to the organization.

### Revoke access to a GitHub organization

To revoke Docker Hub's access to an organization's GitHub repositories:

1. From your GitHub Account settings, locate the **Organization settings** section at the lower left.

2. Select the organization you want to revoke Docker Hub's access to.

3. From the Organization Profile menu, select **Third-party access**.
    The page displays a list of third party applications and their access status.

4. Select the pencil icon next to Docker Hub Builder.

5. On the next page, select **Deny access**.

### Unlink a GitHub user account

To revoke Docker Hub's access to your GitHub account, you must unlink it both
from Docker Hub, and from your GitHub account.

1. Select the **Settings** icon in the top-right navigation, then select **Repository Settings**.

2. From the **Linked accounts** tab, select the plug icon next to the source provider you want to remove.

3. Go to your GitHub account's **Settings** page.

4. Select **Applications** in the left navigation bar.

5. Select the `...` menu to the right of the Docker Hub Builder application and select **Revoke**.

> [!NOTE]
>
> Each repository that is configured as an automated build source
contains a webhook that notifies Docker Hub of changes in the repository.
This webhook is not automatically removed when you revoke access to a source
code provider.

## Link to a Bitbucket user account

1. Sign in to Docker Hub using your Docker ID.

2. Select the **Settings** icon in the top-right navigation, then select **Repository Settings**.

3. From the **Linked accounts** tab, select **Link provider** for the source provider you want to link.

4. If necessary, sign in to Bitbucket.

5. On the page that appears, select **Grant access**.

### Unlink a Bitbucket user account

To permanently revoke Docker Hub's access to your Bitbucket account, you must
unlink it both from Docker Hub, and revoke authorization in your Bitbucket account.

1. Sign in to Docker Hub.

2. Select the **Settings** icon in the top-right navigation, then select **Repository Settings**.

3. From the **Linked accounts** tab, select the **Plug** icon next to the source provider you want to remove.

> [!IMPORTANT]
> After unlinking the account on Docker Hub, you must also revoke the authorization on the Bitbucket end.

To revoke authorization in your Bitbucket account:

1. Go to your Bitbucket account and navigate to [**Bitbucket settings**](https://bitbucket.org/account/settings/app-authorizations/).

2. On the page that appears, select **OAuth**.

3. Select **Revoke** next to the Docker Hub line.

![Bitbucket Authorization revocation page](images/bitbucket-revoke.png)

> [!NOTE]
>
> Each repository that is configured as an automated build source
contains a webhook that notifies Docker Hub of changes in the repository. This
webhook is not automatically removed when you revoke access to a source code
provider.
