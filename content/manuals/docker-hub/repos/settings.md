---
description: Learn about personal repository settings in Docker Hub
keywords: Docker Hub, Hub, repositories, settings
title: Personal settings for repositories
linkTitle: Personal settings
toc_max: 3
weight: 50
---

For your account, you can set personal settings for repositories, including
default repository privacy and autobuild notifications.

## Default repository privacy

When creating a new repository in Docker Hub, you are able to specify the
repository visibility. You can also change the visibility at any time in Docker Hub.

The default setting is useful if you use the `docker push` command to push to a
repository that doesn't exist yet. In this case, Docker Hub automatically
creates the repository with your default repository privacy.

### Configure default repository privacy

1. Sign in to [Docker Hub](https://hub.docker.com).
2. Select **Repositories**.
3. Near the top-right corner, select the settings icon and then **Repository Settings**.
4. Select the **Default privacy** for any new repository created.

   - **Public**: All new repositories appear in Docker Hub search results and can be
     pulled by everyone.
   - **Private**: All new repositories don't appear in Docker Hub search results
     and are only accessible to you and collaborators. In addition, if the
     repository is created in an organization's namespace, then the repository
     is accessible to those with applicable roles or permissions.

5. Select **Save**.

## Autobuild notifications

You can send notifications to your email for all your repositories using
autobuilds.

### Configure autobuild notifications

1. Sign in to [Docker Hub](https://hub.docker.com).
2. Select **Repositories**.
3. Near the top-right corner, select the settings icon and then **Repository Settings**.
4. Select the **Notifications**
5. Select the notifications to receive by email.

   - **Off**: No notifications.
   - **Only failures**: Only notifications about failed builds.
   - **Everything**: Notifications for successful and failed builds.

6. Select **Save**.
