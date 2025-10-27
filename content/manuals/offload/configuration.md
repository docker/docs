---
title: Configure Docker Offload
linktitle: Configure
weight: 20
description: Learn how to configure build settings for Docker Offload.
keywords: cloud, configuration, settings, offload, idle, timeout
---

{{< summary-bar feature_name="Docker Offload" >}}

> [!NOTE]
>
> All free trial usage granted for the Docker Offload Beta expire after 90 days from the time they are granted. To
> continue using Docker Offload Beta after your usage expires, you can enable on-demand usage at [Docker Home
> Billing](https://app.docker.com/billing).
>
> Implementation and settings for the Docker Offload Beta may differ from the General Availability (GA) release. The
> following topic describes the user experience for the GA release.

For organization owners, you can manage Docker Offload settings for all users in your organization. For more details,
see [Manage Docker products](../admin/organization/manage-products.md). To view usage and configure billing for Docker
Offload, see [Docker Offload usage and billing](/offload/usage/).


For developers, you can manage Docker Offload settings in Docker Desktop. To manage settings:

1. Open the Docker Desktop Dashboard and sign in.
2. Select the settings icon in the Docker Desktop Dashboard header.
3. In **Settings**, select **Docker Offload**.

   Here you can:

   - Enable or disable Docker Offload.
   - Select the idle timeout. This is the duration of inactivity after which Docker Offload enters idle mode and no
     longer incurs usage. The default is 5 minutes.