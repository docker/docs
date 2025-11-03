---
title: Configure Docker Offload
linktitle: Configure
weight: 20
description: Learn how to configure build settings for Docker Offload.
keywords: cloud, configuration, settings, offload
---

{{< summary-bar feature_name="Docker Offload" >}}

You can configure Docker Offload settings at different levels depending on your role. Organization owners can manage
settings for all users in their organization, while individual developers can configure their own Docker Desktop
settings.

## Manage settings for your organization

For organization owners, you can manage Docker Offload settings for all users in your organization. For more details,
see [Manage Docker products](../admin/organization/manage-products.md). To view usage and configure billing for Docker
Offload, see [Docker Offload usage and billing](/offload/usage/).

## Configure settings in Docker Desktop

For developers, you can manage Docker Offload settings in Docker Desktop. To manage settings:

1. Open the Docker Desktop Dashboard and sign in.
2. Select the settings icon in the Docker Desktop Dashboard header.
3. In **Settings**, select **Docker Offload**.

   Here you can:

   - Toggle **Enable Docker Offload**. When enabled, you can start Offload sessions.
   - Select **Idle timeout**. This is the duration of time between no activity and Docker Offload entering idle mode.
     For details about idle timeout, see [Understand active and idle states](#understand-active-and-idle-states).

### Understand active and idle states

Docker Offload automatically transitions between active and idle states to help
you control costs while maintaining a seamless development experience.

#### When your session is active

Your Docker Offload environment is active when you're building images, running
containers, or actively interacting with them, such as viewing logs or
maintaining an open network connection. During active state:

- Usage is charged
- A remote Docker Engine is connected to your local machine
- All container operations execute in the cloud environment

#### When your session is idle

When there's no activity, Docker Offload transitions to idle state. During idle
state:

- You are not charged for usage
- The remote connection is suspended
- No containers are running in the cloud

The idle transition delay can be configured in Docker Desktop settings, ranging
from 10 seconds to 1 hour. This setting determines how long Docker Offload
waits after detecting inactivity before transitioning to idle state.

#### How your session is preserved

If your session has been idle for less than 5 minutes and you resume activity,
your previous containers and images are preserved and remain available. This
allows you to pick up right where you left off.

However, if the idle period exceeds 5 minutes, a new session starts with a
clean environment and any containers, images, or volumes from the previous
session are deleted.

> [!NOTE]
>
> Transitioning from active to idle and back to active within 5 minutes will be
> charged as continuous usage.