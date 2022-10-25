---
description: Settings Management for desktop
keywords: Settings Management, rootless, docker desktop, hardened desktop
title: What is Settings Management?
--- 
>**Note**
>
>Settings Management is available to Docker Business customers only. 

Settings Management is a feature that helps admins to control certain Docker Desktop settings on client machines within their organization. 

With a few lines of JSON, admins can configure controls for Docker Desktop settings such as proxies and network settings. For an extra layer of security, admins can also use Settings Management to enable [Enhanced Container Isolation](../enhanced-container-isolation/index.md) which ensures that any configurations set with Settings Management cannot be modified by containers.

It is available with [Docker Desktop 4.13.0 or later](../../release-notes.md).

### Who is it for? 

- For Organizations who wish to configure Docker Desktop to be within their organization's centralized control.
- For Organizations who want to create a standardized Docker Desktop environment at scale.
- For Docker Business customers who want to confidently manage their use of Docker Desktop within tightly regulated environments.

### How does it work?

Administrators can configure several Docker Desktop settings using an `admin-settings.json` file. This file is located on the Docker Desktop host and can only be accessed by users with root or admin privileges. 

Values that are set to `locked: true` within the `admin-settings.json` override any previous values set by users and ensure that these cannot be modified. For more information, see [Configure Settings Management](../settings-management/configure.md#step-two-configure-the-settings-you-want-to-lock-in).

### What features can I configure with Settings Management?

Using the `admin-settings.json` file, admins can:

- Enable [Enhanced Container Isolation](../enhanced-container-isolation/index.md) (currently incompatible with WSL)
- Configure HTTP proxies
- Configure network settings
- Enforce the use of WSL2 based engine or Hyper-V
- Configure Docker Engine
- Turn off Docker Desktop's ability to checks for updates
- Turn off Docker Desktop's ability to send usage statistics

For more details on the syntax and options admins can set, see [Configure Settings Management](configure.md).

### How do I set up and enforce Settings Management?

As an administrator, you first need to [configure a registry.json to enforce sign-in](../../../docker-hub/configure-sign-in.md). This is because the Settings Management feature requires a Docker Business subscription and therefore your Docker Desktop users must authenticate to your organization for this configuration to take effect.

Next, you must either manually [create and configure the admin-settings.json file](configure.md), or use the `--admin-settings` installer flag on [macOS](../../install/mac-install.md#install-from-the-command-line) or [Windows](../../install/windows-install.md#install-from-the-command-line) to automatically create the `admin-settings.json` and save it in the correct location.

Once this is done, Docker Desktop users receive the changed settings when they either:
- Quit, re-launch, and sign in to Docker Desktop
- Launch and sign in to Docker Desktop for the first time

Docker doesn't automatically mandate that developers re-launch and re-authenticate once a change has been made, so as not to disrupt your developers' workflow. 

### What do users see when the settings are enforced?

Any settings that are enforced, are grayed out in Docker Desktop and the user is unable to edit them, either via the Docker Desktop UI, CLI, or the `settings.json` file. In addition, if Enhanced Container Isolation is enforced, users can't use privileged containers or similar techniques to modify enforced settings within the Docker Desktop Linux VM, for example, reconfigure proxy and networking of reconfigure Docker Engine.

![Proxy settings grayed out](/assets/images/grayed-setting.png){:width="750px"}
