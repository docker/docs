---
description: Understand how Settings Management works, who it is for, and what the
  benefits are
keywords: Settings Management, rootless, docker desktop, hardened desktop
tags: [admin]
title: What is Settings Management?
linkTitle: Settings Management
aliases:
 - /desktop/hardened-desktop/settings-management/
weight: 10
---

{{< summary-bar feature_name="Hardened Docker Desktop" >}}

Settings Management lets administrators configure and enforce Docker Desktop
settings across end-user machines. It helps maintain consistent configurations
and enhances security within your organization.

## Who is it for?

Settings Management is designed for organizations that:

- Require centralized control over Docker Desktop configurations.
- Aim to standardize Docker Desktop environments across teams.
- Operate in regulated environments and need to enforce compliance.

This feature is available with a Docker Business subscription.

## How it works

Administrators can define settings using one of the following methods:

- [Admin Console](/manuals/security/for-admins/hardened-desktop/settings-management/configure-admin-console.md): Create and assign settings policies through the
Docker Admin Console.
- [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md): Place a configuration file on the
user's machine to enforce settings.

Enforced settings override user-defined configurations and can't be modified
by developers.

## Configurable settings

Settings Management supports a broad range of Docker Desktop features,
including proxies, network configurations, and container isolation.

For a full list of settings you can enforce, see the [Settings reference](/manuals/security/for-admins/hardened-desktop/settings-management/settings-reference.md).

## Set up Settings Management

1. [Enforce sign-in](/manuals/security/for-admins/enforce-sign-in/_index.md) to
ensure all developers authenticate with your organization.
2. Choose a configuration method:
    - Use the `--admin-settings` installer flag on [macOS](/manuals/desktop/setup/install/mac-install.md#install-from-the-command-line) or [Windows](/manuals/desktop/setup/install/windows-install.md#install-from-the-command-line) to automatically create the `admin-settings.json`.
    - Manually create and configure the [`admin-settings.json` file](/manuals/security/for-admins/hardened-desktop/settings-management/configure-json-file.md).
    - Create a settings policy in the [Docker Admin Console](configure-admin-console.md).

After configuration, developers receive the enforced setting when they:

- Quit and relaunch Docker Desktop, then sign in.
- Launch and sign in to Docker Desktop for the first time.

> [!NOTE]
>
> Docker Desktop does not automatically prompt users to restart or re-authenticate
after a settings change.

## Developer experience

When settings are enforced:

- Options appear grayed out in Docker Desktop and can't be modified via the
Dashboard, CLI, or configuration files.
- If Enhanced Container Isolation is enabled, developers can't use privileged
containers or similar methods to alter enforced settings within the Docker
Desktop Linux VM.

## What's next?

- [Configure Settings Management with the `admin-settings.json` file](configure-json-file.md)
- [Configure Settings Management with the Docker Admin Console](configure-admin-console.md)

## Learn more

To see how each Docker Desktop setting maps across the Docker Dashboard, `admin-settings.json` file, and Admin Console, see the [Settings reference](settings-reference.md).
