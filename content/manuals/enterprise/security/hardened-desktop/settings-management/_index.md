---
description: Understand how Settings Management works, who it's for, and the benefits it provides
keywords: Settings Management, rootless, docker desktop, hardened desktop, admin control, enterprise
tags: [admin]
title: Settings Management
linkTitle: Settings Management
aliases:
 - /desktop/hardened-desktop/settings-management/
 - /security/for-admins/hardened-desktop/settings-management/
weight: 10
---

{{< summary-bar feature_name="Hardened Docker Desktop" >}}

Settings Management lets administrators configure and enforce Docker Desktop settings across end-user machines. It helps maintain consistent configurations and enhances security within your organization.

## Who should use Settings Management?

Settings Management is designed for organizations that:

- Need centralized control over Docker Desktop configurations
- Want to standardize Docker Desktop environments across teams
- Operate in regulated environments and must enforce compliance policies

## How Settings Management works

Administrators can define settings using one of these methods:

- [Admin Console](/manuals/enterprise/security/hardened-desktop/settings-management/configure-admin-console.md): Create and assign settings policies through the
Docker Admin Console. This provides a web-based interface for managing settings
across your organization.
- [`admin-settings.json` file](/manuals/enterprise/security/hardened-desktop/settings-management/configure-json-file.md): Place a configuration file on the
user's machine to enforce settings. This method works well for automated
deployments and scripted installations.

Enforced settings override user-defined configurations and can't be modified by developers.

## Configurable settings

Settings Management supports a wide range of Docker Desktop features, including:

- Proxy configurations
- Network settings
- Container isolation options
- Registry access controls
- Resource limits
- Security policies

For a complete list of settings you can enforce, see the [Settings reference](/manuals/enterprise/security/hardened-desktop/settings-management/settings-reference.md).

## Policy precedence

When multiple policies exist, Docker Desktop applies them in this order:

1. User-specific policies: Highest priority
1. Organization default policy: Applied when no user-specific policy exists
1. Local `admin-settings.json` file: Lowest priority, overridden by Admin Console policies

## Set up Settings Management

1. [Enforce sign-in](/manuals/enterprise/security/enforce-sign-in/_index.md) to
ensure all developers authenticate with your organization.
2. Choose a configuration method:
    - Use the `--admin-settings` installer flag on [macOS](/manuals/desktop/setup/install/mac-install.md#install-from-the-command-line) or [Windows](/manuals/desktop/setup/install/windows-install.md#install-from-the-command-line) to automatically create the `admin-settings.json`.
    - Manually create and configure the [`admin-settings.json` file](/manuals/enterprise/security/hardened-desktop/settings-management/configure-json-file.md).
    - Create a settings policy in the [Docker Admin Console](configure-admin-console.md).

After configuration, developers receive the enforced settings when they:

- Quit and relaunch Docker Desktop, then sign in
- Launch and sign in to Docker Desktop for the first time

> [!NOTE]
>
> Docker Desktop doesn't automatically prompt users to restart or re-authenticate after a settings change. You may need to communicate these requirements to your developers.

## Developer experience

When settings are enforced:

- Settings options appear grayed out in Docker Desktop and can't be modified through the Dashboard, CLI, or configuration files
- If Enhanced Container Isolation is enabled, developers can't use privileged containers or similar methods to alter enforced settings within the Docker Desktop Linux VM

This ensures consistent environments while maintaining a clear visual indication of which settings are managed by administrators.

## View applied settings

When administrators apply Settings Management policies, Docker Desktop greys out most enforced settings in the GUI.

The Docker Desktop GUI doesn't currently display all centralized settings,
particularly Enhanced Container Isolation (ECI) settings that administrators
apply via the Admin Console.

As a workaround, you can check the `settings-store.json` file to view all
applied settings:

  - Mac: `~/Library/Application Support/Docker/settings-store.json`
  - Windows: `%APPDATA%\Docker\settings-store.json`
  - Linux: `~/.docker/desktop/settings-store.json`

The `settings-store.json` file contains all settings, including those that
may not appear in the Docker Desktop GUI.

## Limitations

Settings Management has the following limitations:

- Doesn't work in air-gapped or offline environments
- Not compatible with environments that restrict authentication with Docker Hub

## Next steps

Get started with Settings Management:

- [Configure Settings Management with the `admin-settings.json` file](configure-json-file.md)
- [Configure Settings Management with the Docker Admin Console](configure-admin-console.md)

