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

Settings Management helps you control key Docker Desktop settings, like proxies and network configurations, on your developers' machines within your organization.

For an extra layer of security, you can also use Settings Management to enable and lock in [Enhanced Container Isolation](../enhanced-container-isolation/_index.md), which prevents containers from modifying any Settings Management configurations.

## Who is it for?

- For organizations that want to configure Docker Desktop to be within their organization's centralized control.
- For organizations that want to create a standardized Docker Desktop environment at scale.
- For Docker Business customers who want to confidently manage their use of Docker Desktop within tightly regulated environments.

## How does it work?

You can configure several Docker Desktop settings using either:
 - An `admin-settings.json` file. This file is located on the Docker Desktop host and can only be accessed by developers with root or administrator privileges.
 - Creating a settings policy in the Docker Admin Console 

Settings that are defined by an administrator override any previous values set by developers and ensure that these cannot be modified. 

## What features can I configure with Settings Management?

Using the `admin-settings.json` file, you can:

- Turn on and lock in [Enhanced Container Isolation](../enhanced-container-isolation/_index.md)
- Configure HTTP proxies
- Configure network settings
- Configure Kubernetes settings
- Enforce the use of WSL 2 based engine or Hyper-V
- Enforce the use of Rosetta for x86_64/amd64 emulation on Apple Silicon
- Configure Docker Engine
- Turn off Docker Desktop's ability to checks for updates
- Turn off Docker Extensions
- Turn off Docker Scout SBOM indexing
- Turn off beta and experimental features
- Turn off Docker AI ([Ask Gordon](../../../../desktop/features/gordon/_index.md))
- Turn off Docker Desktop's onboarding survey
- Control whether developers can use the Docker terminal
- Control the file sharing implementation for your developers on macOS
- Specify which paths your developers can add file shares to
- Configure Air-gapped containers

For more details on the syntax and options, see [Configure Settings Management](configure-json-file.md).

## How do I set up and enforce Settings Management?

You first need to [enforce sign-in](/manuals/security/for-admins/enforce-sign-in/_index.md) to ensure that all Docker Desktop developers authenticate with your organization. Since the Settings Management feature requires a Docker Business subscription, enforced sign-in guarantees that only authenticated users have access and that the feature consistently takes effect across all users, even though it may still work without enforced sign-in.

Next, you must either:
 - Manually [create and configure the `admin-settings.json` file](configure-json-file.md), or use the `--admin-settings` installer flag on [macOS](/manuals/desktop/setup/install/mac-install.md#install-from-the-command-line) or [Windows](/manuals/desktop/setup/install/windows-install.md#install-from-the-command-line) to automatically create the `admin-settings.json` and save it in the correct location.
 - Fill out the **Settings policy** creation form in the [Docker Admin Console](configure-admin-console.md).

Once this is done, Docker Desktop developers receive the changed settings when they either:
- Quit, re-launch, and sign in to Docker Desktop
- Launch and sign in to Docker Desktop for the first time

To avoid disrupting your developers' workflows, Docker doesn't automatically require that developers re-launch and re-authenticate once a change has been made.

## What do developers see when the settings are enforced?

Enforced settings appear grayed out in Docker Desktop. They can't be edited via the Docker Desktop Dashboard, CLI, or `settings-store.json` (or `settings.json` for Docker Desktop 4.34 and earlier).

In addition, if Enhanced Container Isolation is enforced, developers can't use privileged containers or similar techniques to modify enforced settings within the Docker Desktop Linux VM. For example, they can't reconfigure proxy and networking, or Docker Engine.

![Proxy settings grayed out](/assets/images/grayed-setting.png)

## What's next?

- [Configure Settings Management with a `.json` file](configure-json-file.md)
- [Configure Settings Management with the Docker Admin Console](configure-admin-console.md)
