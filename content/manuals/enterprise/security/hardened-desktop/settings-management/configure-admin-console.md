---
title: Configure Settings Management with the Admin Console
linkTitle: Use the Admin Console
description: Configure and enforce Docker Desktop settings across your organization using the Docker Admin Console
keywords: admin console, settings management, policy configuration, enterprise controls, docker desktop
weight: 20
aliases:
 - /security/for-admins/hardened-desktop/settings-management/configure-admin-console/
---

{{< summary-bar feature_name="Admin Console" >}}

Use the Docker Admin Console to create and manage settings policies for Docker Desktop across your organization. Settings policies let you standardize configurations, enforce security requirements, and maintain consistent Docker Desktop environments.

## Prerequisites

Before you begin, make sure you have:

- [Docker Desktop 4.37.1 or later](/manuals/desktop/release-notes.md) installed
- [A verified domain](/manuals/enterprise/security/single-sign-on/configure.md#step-one-add-and-verify-your-domain)
- [Enforced sign-in](/manuals/enterprise/security/enforce-sign-in/_index.md) for your organization
- A Docker Business subscription

> [!IMPORTANT]
>
> You must add users to your verified domain for settings to take effect.

## Create a settings policy

To create a new settings policy:

1. Sign in to [Docker Home](https://app.docker.com/) and select
your organization.
1. Select **Admin Console**, then **Desktop Settings Management**.
1. Select **Create a settings policy**.
1. Provide a name and optional description.

      > [!TIP]
      >
      > You can upload an existing `admin-settings.json` file to pre-fill the form.
      Admin Console policies override local `admin-settings.json` files.

1. Choose who the policy applies to:
   - All users
   - Specific users

      > [!NOTE]
      >
      > User-specific policies override global default policies. Test your policy with a small group before applying it organization-wide.

1. Configure each setting using a state:
   - **User-defined**: Users can change the setting.
   - **Always enabled**: Setting is on and locked.
   - **Enabled**: Setting is on but can be changed.
   - **Always disabled**: Setting is off and locked.
   - **Disabled**: Setting is off but can be changed.

      > [!TIP]
      >
      > For a complete list of configurable settings, supported platforms, and configuration methods, see the [Settings reference](settings-reference.md).

1. Select **Create** to save your policy.

## Apply the policy

Settings policies take effect after Docker Desktop restarts and users sign in.

For new installations:

1. Launch Docker Desktop.
1. Sign in with your Docker account.

For existing installations:

1. Quit Docker Desktop completely.
1. Relaunch Docker Desktop.

> [!IMPORTANT]
>
> Users must fully quit and reopen Docker Desktop. Restarting from the Docker Desktop menu isn't sufficient.

Docker Desktop checks for policy updates when it launches and every 60 minutes while running.

## Verify applied settings

After you apply policies:

- Docker Desktop displays most settings as greyed out
- Some settings, particularly Enhanced Container Isolation configurations,
may not appear in the GUI
- You can verify all applied settings by checking the [`settings-store.json`
file](/manuals/desktop/settings-and-maintenance/settings.md) on your system

## Manage existing policies

From the **Desktop Settings Management** page in the Admin Console, use the **Actions** menu to:

- Edit or delete an existing settings policy
- Export a settings policy as an `admin-settings.json` file
- Promote a user-specific policy to be the new global default

## Roll back policies

To roll back a settings policy:

- Complete rollback: Delete the entire policy.
- Partial rollback: Set specific settings to **User-defined**.

When you roll back settings, users regain control over those settings configurations.
