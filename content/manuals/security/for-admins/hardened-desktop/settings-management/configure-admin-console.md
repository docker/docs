---
description: How to configure Settings Management for Docker Desktop using the Docker Admin Console
keywords: admin, controls, rootless, enhanced container isolation
title: Configure Settings Management with the Admin Console
linkTitle: Use the Admin Console
weight: 20
---

{{< summary-bar feature_name="Admin Console" >}}

This page explains how administrators can use the Docker Admin Console to create
and apply settings policies for Docker Desktop. These policies help standardize
and secure Docker Desktop environments across your organization.

## Prerequisites

- [Install Docker Desktop 4.36.0 or later](/manuals/desktop/release-notes.md).
- [Verify your domain](/manuals/security/for-admins/single-sign-on/configure.md#step-one-add-and-verify-your-domain).
- [Enforce sign-in](/manuals/security/for-admins/enforce-sign-in/_index.md) to
ensure users authenticate to your organization.
- A Docker Business subscription is required.

> [!IMPORTANT]
>
> You must add users to your verified domain for settings to take effect.

## Create a settings policy

1. Go to the [Docker Admin Console](https://app.docker.com/admin) and select
your organization.
2. Under **Docker Desktop**, select **Settings Management**.
3. Select **Create a settings policy**.
4. Provide a name and optional description.

   > [!TIP]
   >
   > You can upload an existing `admin-settings.json` file to pre-fill the form.
   Admin Console policies override local `admin-settings.json` files.

5. Choose who the policy applies to:
   - All users
   - Specific users

   > [!NOTE]
   >
   > User-specific policies override the global default. Test your policy with
   a few users before rolling it out globally.

6. Configure the state for each setting:
   - **User-defined**: Users can change the setting.
   - **Always enabled**: Setting is on and locked.
   - **Enabled**: Setting is on but can be changed.
   - **Always disabled**: Setting is off and locked.
   - **Disabled**: Setting is off but can be changed.

   > [!TIP]
   >
   > For a complete list of available settings, their supported platforms, and which configuration methods they work with, see the [Settings reference](settings-reference.md).
7. Select **Create**.

To apply the policy:

- New installs: Launch Docker Desktop and sign in.
- Existing installs: Fully quit and relaunch Docker Desktop.

> [!IMPORTANT]
>
> Restarting from the Docker Desktop menu isn't enough. Users must fully quit
and relaunch Docker Desktop.

Docker Desktop checks for policy updates at launch and every 60 minutes. To roll
back a policy, either delete it or set individual settings to **User-defined**.

## Manage policies

From the **Actions** menu on the **Settings Management** page, you can:

- Edit or delete an existing settings policy
- Export a settings policy as an `admin-settings.json` file
- Promote a user-specific policy to be the new global default

## Learn more

To see how each Docker Desktop setting maps across the Docker Dashboard, `admin-settings.json` file, and Admin Console, see the [Settings reference](settings-reference.md).