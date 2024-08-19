---
description: How to configure Settings Management for Docker Desktop using the Docker Admin Console
keywords: admin, controls, rootless, enhanced container isolation
title: Configure with the Docker Admin Console
---

>**Note**
>
>Settings Management is available to Docker Business customers only.

This page contains information for administrators on how to configure Settings Management with the Docker Admin Console to specify and lock configuration parameters to create a standardized Docker Desktop environment across your Docker company or organization.

## Prerequisites

- [Download and install Docker Desktop 4.34.0 or later](/desktop/release-notes.md).
- As an administrator, you need to [enforce
  sign-in](/security/for-admins/enforce-sign-in/_index.md). This is
because the Settings Management feature requires a Docker Business
subscription and therefore your Docker Desktop users must authenticate to your
organization for configurations to take effect. 

## Setup a settings policy

1. Within the [Docker Admin Console](https://admin.docker.com/) navigate to the company or organization you want to define a settings policy for. 
2. Under the **Security and access** section, select **Desktop Settings Management** policy. 
3. Select **Create settings policy**.
4. Give your settings policy a name and an optional description.
5. Assign the setting policy to all your users within the company or organization, or specific users. 

   > [!NOTE]
   >
   > If a policy is assigned to all users, it sets the policy as the global default policy.

6. Configure the settings for the policy. Go through each setting and select your chosen setting state. You can choose:
   - **User-defined**. 
   - **Always enabled**. This means the setting is turned on and your users won't be able to edit this setting from Docker Desktop or the CLI.
   - **Enabled**. The setting is turned on and users can edit this setting from Docker Desktop or the CLI.
   - **Always disabled**. This means the setting is turned off and your users won't be able to edit this setting from Docker Desktop or the CLI.
   - **Disabled**. The setting is turned off and users can edit this setting from Docker Desktop or the CLI.
   
   > [!TIP]
   >
   > If you have already configured Settings Management with an `admin-settings.json` file for an organization, you can upload it using the **Upload existing settings** button which then automatically populates the form for you. 
   
7. Select **Create**

For the settings policy to take effect:
- On a new install, users need to launch Docker Desktop and authenticate to their organization.
- On an existing install, users need to quit Docker Desktop through the Docker menu, and then relaunch Docker Desktop. If they are already signed in, they don't need to sign in again for the changes to take effect.
  > [!IMPORTANT]
  >
  > Selecting **Restart** from the Docker menu isn't enough as it only restarts some components of Docker Desktop.

Docker doesn't automatically mandate that users re-launch and sign in once a change has been made so as not to disrupt your users' workflow.

> [!NOTE]
>
> Settings are synced to Docker Desktop and the CLI when a user is signed in and starts Docker Desktop, and then every 60 minutes. 

## Settings policy actions

From the **Desktop Settings Management** view in the Docker Admin Console, you can:
- Edit or delete an existing settings policy. 
- Export a settings policy as an `admin-settings.json` file.
- Set a policy that is applied to a select group of users, to be the new global default policy for all users. 