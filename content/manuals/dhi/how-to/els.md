---
title: 'Use Extended Lifecycle Support for Docker Hardened Images <span class="not-prose bg-blue-500 dark:bg-blue-400 rounded-sm px-1 text-xs text-white whitespace-nowrap">DHI Enterprise</span>'
linktitle: Use Extended Lifecycle Support
description: Learn how to use Extended Lifecycle Support with Docker Hardened Images.
weight: 39
keywords: extended lifecycle support, docker hardened images, container security, image lifecycle, vulnerability management
---

{{< summary-bar feature_name="Docker Hardened Images" >}}

With a Docker Hardened Images subscription add-on, you can use Extended
Lifecycle Support (ELS) for Docker Hardened Images. ELS provides security
patches for end-of-life (EOL) image versions, letting you maintain secure,
compliant operations while planning upgrades on your own timeline. You can use
ELS images like any other Docker Hardened Image, but you must enable ELS for
each repository you want to use with ELS.

## Discover repositories with ELS support

To find images with ELS support:

1. Go to [Docker Hub](https://hub.docker.com) and sign in.
2. Select **My Hub**.
3. In the namespace drop-down, select your organization.
4. Select **Hardened Images** > **Catalog**.
5. In **Filter by**, select **Extended Lifecycle Support**.

## Enable ELS for a repository

To enable ELS for a repository, an organization owner must [mirror](./mirror.md)
the repository to your organization.

To enable ELS when mirroring:

1. Go to [Docker Hub](https://hub.docker.com) and sign in.
2. Select **My Hub**.
3. In the namespace drop-down, select your organization.
4. Select **Hardened Images** > **Catalog**.
5. Select a DHI repository to view its details.
6. Select **Use this image** > **Mirror repository**
7. Select **Enable support for end-of-life versions** and then follow the
   on-screen instructions.

## Disable ELS for a repository

To disable ELS for a repository, you must uncheck the ELS option in the mirrored
repository's **Settings** tab, or stop mirroring the repository. To stop mirroring, see
[Stop mirroring a repository](./mirror.md#stop-mirroring-a-repository).

To update settings:

1. Go to [Docker Hub](https://hub.docker.com) and sign in.
2. Select **My Hub**.
3. In the namespace drop-down, select your organization.
4. Select **Repositories** and then select the mirrored repository.
5. Select the **Settings** tab.
6. Uncheck the **Mirror end-of-life images** option.

## Manage ELS repositories

You can view and manage your mirrored repositories with ELS like any other
mirrored DHI repository. For more details, see [Manage images](./manage.md).