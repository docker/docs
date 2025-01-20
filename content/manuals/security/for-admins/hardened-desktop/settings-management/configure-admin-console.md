---
description: How to configure Settings Management for Docker Desktop using the Docker Admin Console
keywords: admin, controls, rootless, enhanced container isolation
title: Configure Settings Management with the Admin Console
linkTitle: Use the Admin Console
weight: 20
params:
  sidebar:
    badge:
      color: violet
      text: EA
---

{{< summary-bar feature_name="Admin Console" >}}

This page contains information for administrators on how to configure Settings Management with the Docker Admin Console. You can specify and lock configuration parameters to create a standardized Docker Desktop environment across your Docker company or organization.

## Prerequisites

- [Download and install Docker Desktop 4.36.0 or later](/manuals/desktop/release-notes.md).
- [Verify your domain](/manuals/security/for-admins/single-sign-on/configure.md#step-one-add-and-verify-your-domain).
- [Enforce sign-in](/manuals/security/for-admins/enforce-sign-in/_index.md). The Settings Management feature requires a Docker Business
subscription, therefore your Docker Desktop users must authenticate to your
organization for configurations to take effect.

## Create a settings policy

1. Within the [Docker Admin Console](https://admin.docker.com/) navigate to the company or organization you want to define a settings policy for.
2. Under the **Security and access** section, select **Desktop Settings Management**.
3. In the top-right corner, select **Create a settings policy**.
4. Give your settings policy a name and an optional description.

   > [!TIP]
   >
   > If you have already configured Settings Management with an `admin-settings.json` file for an organization, you can upload it using the **Upload existing settings** button which then automatically populates the form for you.
   >
   > Settings policies deployed via the Docker Admin Console take precedence over manually deployed `admin-settings.json` files.

5. Assign the setting policy to all your users within the company or organization, or specific users.

   > [!NOTE]
   >
   > If a settings policy is assigned to all users, it sets the policy as the global default policy. You can only have one global settings policy at a time.
   > If a user already has a user-specific settings policy assigned, the user-specific policy takes precedence over a global policy.

   > [!TIP]
   >
   > Before setting a global settings policy, it is recommended that you first test it as a user-specific policy to make sure you're happy with the changes before proceeding.

6. Configure the settings for the policy. Go through each setting and select your chosen setting state. You can choose:
   - **User-defined**. Your developers are able to control and change this setting.
   - **Always enabled**. This means the setting is turned on and your users won't be able to edit this setting from Docker Desktop or the CLI.
   - **Enabled**. The setting is turned on and users can edit this setting from Docker Desktop or the CLI.
   - **Always disabled**. This means the setting is turned off and your users won't be able to edit this setting from Docker Desktop or the CLI.
   - **Disabled**. The setting is turned off and users can edit this setting from Docker Desktop or the CLI.
7. Select **Create**

For the settings policy to take effect:
- On a new install, users need to launch Docker Desktop and authenticate to their organization.
- On an existing install, users need to quit Docker Desktop through the Docker menu, and then re-launch Docker Desktop. If they are already signed in, they don't need to sign in again for the changes to take effect.

  > [!IMPORTANT]
  >
  > Selecting **Restart** from the Docker menu isn't enough as it only restarts some components of Docker Desktop.

To avoid disrupting your users' workflows, Docker doesn't automatically require that users re-launch once a change has been made.

> [!NOTE]
>
> Settings are synced to Docker Desktop and the CLI when a user is signed in and starts Docker Desktop, and then every 60 minutes.

If your settings policy needs to be rolled back, either delete the policy or edit the policy to set individual settings to **User-defined**.

## Settings policy actions

From the **Actions** menu on the **Desktop Settings Management** page in the Docker Admin Console, you can:
- Edit or delete an existing settings policy.
- Export a settings policy as an `admin-settings.json` file.
- Promote a policy that is applied to a select group of users, to be the new global default policy for all users.