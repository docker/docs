---
title: Configure Docker Offload v1
linktitle: Configure
weight: 20
description: Learn how to configure build settings for Docker Offload v1.
keywords: cloud, configuration, settings, offload
aliases:
 - /offload/configuration/
---

{{< summary-bar feature_name="Docker Offload v1" >}}

You can configure Docker Offload v1 settings at different levels depending on your
role. Organization owners can manage settings for all users in their
organization, while individual developers can configure their own Docker Desktop
settings when allowed by their organization.

## Manage settings for your organization

For organization owners, you can manage Docker Offload v1 settings for all users in
your organization. For more details, see [Manage Docker
products](/manuals/accounts/organization/manage/manage-products.md). To view
usage for Docker Offload v1, see [Docker Offload v1
usage](/manuals/offload/v1/usage.md).

## Configure settings in Docker Desktop

For developers, you can enable or disable Docker Offload in Docker Desktop if
allowed by your organization. To manage this setting:

1. Open the Docker Desktop Dashboard and sign in.
2. Select the settings icon in the Docker Desktop Dashboard header.
3. In **Settings**, select **Docker Offload**.
4. Toggle **Enable Docker Offload**. When enabled, you can start Offload sessions.

### Move to the current version of Docker Offload

To use `docker offload create` and `docker --cloud`, turn off **Enable Docker
Offload** and make sure your account has a Docker Agentic Platform
subscription. Turning off the setting alone doesn't grant access. For what
changes, see [Docker Offload](/manuals/offload/_index.md).

You can keep using v1 by leaving the setting on. The two versions can't run at
the same time.