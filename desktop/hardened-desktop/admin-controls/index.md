---
description: admin controls for desktop
keywords: admin controls, rootless, docker desktop, hardened desktop
title: What is Admin Controls?
--- 
>Note
>
>Admin Controls is available to Docker Business customers only. 

Admin Controls is a feature that gives admins the ability to configure certain Docker Desktop settings for their organization. 

With a few lines of JSON, admins can configure Docker Desktop settings for HTTP proxies, network settings, and the Docker Engine configuration. This saves significant time and cost in securing developer workflows. If admins also enable [Enhanced Container Isolation](../enhanced-container-isolation/index.md), it ensures that any configurations set with Admin Controls cannot be modified by containers.

### Who is it for? 

- For Organizations who wish to configure Docker Desktop to be within their organization's centralized control.
- For Organizations who want to create a standardized Docker Desktop environment at scale.
- For security conscious Docker Business customers who want to confidently manage their use of Docker Desktop within tightly regulated environments.

## How does it work?

Administrators can configure several Docker Desktop settings using the `admin-settings.json` file. This file is located within the Docker Desktop host and can only be accessed by users with root or admin privileges. 

As the `admin-settings.json` overrides the `settings.json` file, which contains the settings configured by regular Docker Desktop users and is located in a directory users can write to, admins can configure default values for several Docker Desktop settings and lock them in so regular users are prevented from modifying them.

### What can be set?

Using the `admin-settings.json` file, admins can:

- Enable [Enhanced Container Isolation](../enhanced-container-isolation/index.md)
- Configure HTTP Proxies
- Configure network settings
- Enforce the use of WSL2 based engine or Hyper-V
- Configure Docker Engine
- Turn off Docker Desktop's ability to checks for updates
- Turn off Docker Desktop's ability to send usage statistics

For more details on the syntax and options admins can set, see [Configure Admin Controls](configure-ac.md).

### How do I set up and enforce Admin Controls?

As an administrator, you first need to [configure a registry.json to enforce sign-in](../../../docker-hub/configure-sign-in.md). This is because your Docker Desktop users must authenticate to your organization for this configuration to take effect.

Next, you must [create and configure the admin-settings.json file](configure-ac.md).

Once this is done, Docker Desktop users receive the changed settings when they re-launch Docker Desktop and re-authenticate to your organization. Docker doesn't automatically mandate that developers re-launch and re-authenticate once a change has been made, so as not to disrupt your developers workflow. 

### What do users see when the settings are enforced?

Docker Desktop users see a notification in the **Settings**, or **Preferences** if using a macOS, which states **Some settings are managed by your Admin**. 

Any settings that are enforced, are grayed out in Docker Desktop and the user is unable to edit them, either via the Docker Desktop UI, CLI, or the `settings.json` file. In addition, if Enhanced Container Isolation is enforced, users can't use privileged containers or similar techniques to modify enforced settings within the Docker Desktop Linux VM, for example, reconfigure proxy and networking of reconfigure Docker Engine.

![Proxy settings grayed out](/assets/images/grayed-setting.png){:width="750px"}
