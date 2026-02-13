---
description: Learn about settings in Docker Hub
keywords: Docker Hub, Hub, repositories, settings
title: Settings
weight: 25
---

You can configure the following settings in Docker Hub:

- [Default privacy](#default-privacy): Settings for all repositories within each
  namespace
- [Notifications](#notifications): Personal settings for autobuild notifications

## Default privacy

You can configure the following default privacy settings for all repositories in
a namespace:

- [Disable creation of public repos](#disable-creation-of-public-repos): Prevent
  organization users from creating public repositories (organization namespaces
  only)
- [Configure default repository privacy](#configure-default-repository-privacy):
  Set the default repository privacy for new repositories


### Disable creation of public repos

{{< summary-bar feature_name="Disable public repositories" >}}

Organization owners and editors can prevent creating public repositories within
organization namespaces. You cannot configure this setting for personal account
namespaces.

> [!NOTE]
>
> Enabling this feature does not affect existing public repositories. Any public
> repositories that already exist will remain public. To make them private, you
> must change their visibility in the individual repository settings.

To configure the disable public repositories setting for an organization
namespace:

1. Sign in to [Docker Hub](https://hub.docker.com).
2. Select **My Hub**.
3. Select your organization from the top-left account drop-down.
4. Select **Settings** > **Default privacy**.
5. Toggle **Disable public repositories** to your desired setting.
6. Select **Save**.

### Configure default repository privacy

Use the default repository privacy setting to automatically set privacy for
repositories created via `docker push` commands when the repository doesn't
exist yet. In this case, Docker Hub automatically creates the repository with
the default repository privacy for that namespace.

> [!NOTE]
>
> You cannot configure the default repository privacy setting when **Disable
> public repositories** is enabled.

To configure the default repository privacy for a namespace:

1. Sign in to [Docker Hub](https://hub.docker.com).
2. Select **My Hub**.
3. Select your organization or account from the top-left account drop-down.
4. Select **Settings** > **Default privacy**.
5. In **Default repository privacy**, select the desired default privacy setting:

   - **Public**: All new repositories appear in Docker Hub search results and can be
     pulled by everyone.
   - **Private**: All new repositories don't appear in Docker Hub search results
     and are only accessible to you and collaborators. In addition, if the
     repository is created in an organization's namespace, then the repository
     is accessible to those with applicable roles or permissions.

6. Select **Save**.

## Notifications

You can send notifications to your email for all your repositories using
autobuilds.

### Configure autobuild notifications

1. Sign in to [Docker Hub](https://hub.docker.com).
2. Select **My Hub**.
3. Select your personal account from the top-left account drop-down.
4. Select **Settings** > **Notifications**.
5. Select the notifications to receive by email:

   - **Off**: No notifications.
   - **Only failures**: Only notifications about failed builds.
   - **Everything**: Notifications for successful and failed builds.

6. Select **Save**.
