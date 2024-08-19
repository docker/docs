---
description: Understand how Settings Management works, who it is for, and what the
  benefits are
keywords: Settings Management, rootless, docker desktop, hardened desktop
title: What is Settings Management?
aliases:
 - /desktop/hardened-desktop/settings-management/
---

>**Note**
>
>Settings Management is available to Docker Business customers only.

Settings Management is a feature that helps admins to control certain Docker Desktop settings on client machines within their organization. It is designed specifically for organizations who don’t give developers root access to their machines.

Administrators can configure controls for Docker Desktop settings such as proxies and network settings. For an extra layer of security, admins can also use Settings Management to enable and lock in [Enhanced Container Isolation](../enhanced-container-isolation/index.md) which ensures that any configurations set with Settings Management cannot be modified by containers.

It is available with [Docker Desktop 4.13.0 and later](/desktop/release-notes.md).

### Who is it for?

- For organizations that want to configure Docker Desktop to be within their organization's centralized control.
- For organizations that want to create a standardized Docker Desktop environment at scale.
- For Docker Business customers who want to confidently manage their use of Docker Desktop within tightly regulated environments.

### How does it work?

Administrators can configure several Docker Desktop settings using either:
 - An `admin-settings.json` file. This file is located on the Docker Desktop host and can only be accessed by developers with root or admin privileges.
 - Creating a settings policy in the Docker Admin Console

Settings that defined by an administrator override any previous values set by developers and ensure that these cannot be modified. 

### What features can I configure with Settings Management?

Administrators can:

- Turn on and lock in [Enhanced Container Isolation](../enhanced-container-isolation/index.md)
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
- Turn off Docker Desktop's onboarding survey
- Control the file sharing implementation for your developers on macOS
- Specify which paths your developers can add file shares to
- Configure Air-Gapped Containers

For more details on the syntax and options admins can set, see [Configure with a .json file](json-file-configure.md) or [Configure with the Docker Admin Console](admin-console-configure.md).

### How do I set up and enforce Settings Management?

As an administrator, you first need to [enforce
sign-in](/security/for-admins/enforce-sign-in/_index.md). This is
because the Enhanced Container Isolation feature requires a Docker Business
subscription and therefore your Docker Desktop users must authenticate to your
organization for this configuration to take effect. 

Next, you must either:
 - Manually [create and configure the admin-settings.json file](configure.md), or use the `--admin-settings` installer flag on [macOS](/desktop/install/mac-install.md#install-from-the-command-line) or [Windows](/desktop/install/windows-install.md#install-from-the-command-line) to automatically create the `admin-settings.json` and save it in the correct location.
 - Fill out the **Settings policy** creation form in the Docker Admin Console

Once this is done, Docker Desktop developers receive the changed settings when they either:
- Quit, re-launch, and sign in to Docker Desktop
- Launch and sign in to Docker Desktop for the first time

Docker doesn't automatically mandate that developers re-launch and re-authenticate once a change has been made, so as not to disrupt your developers' workflow.

### What do developers see when the settings are enforced?

Any settings that are enforced, are grayed out in Docker Desktop and the user is unable to edit them, either via the Docker Desktop UI, CLI, or the `settings.json` file. In addition, if Enhanced Container Isolation is enforced, developers can't use privileged containers or similar techniques to modify enforced settings within the Docker Desktop Linux VM, for example, reconfigure proxy and networking of reconfigure Docker Engine.

![Proxy settings grayed out](/assets/images/grayed-setting.png)

## What's next?

- [Configure Settings Management with a .json file](json-file-configure.md)
- [Configure Settings Management with the Docker Admin Console](admin-console-configure.md)