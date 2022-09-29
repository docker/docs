---
description: admin controls for desktop
keywords: admin controls, rootless, docker desktop, hardened desktop
title: What is Admin Controls?
--- 
>Note
>
>Admin Controls is currently in [Early Access](../../../release-lifecycle.md#early-access-ea) and available to Docker Business customers only. 

Admin Controls is a feature that gives administrators the ability to enforce certain Docker Desktop settings for their organization. 

With just a few lines of JSON, administrators are able to enforce preferences for HTTP proxies, network settings, and the Docker Engine configuration. This saves significant time and cost in securing developer workflows.

### Who is it for? 

- For Organizations who wish to configure Docker Desktop to be within their organization's centralized control.
- For Organizations who want to create a standardized Docker Desktop environment at scale.
- For security conscious Docker Business customers who want to confidently manage their use of Docker Desktop within tightly regulated environments.

### What can be set?

Using the `admin-settings.json` file, administrators can:

- Switch on Enhanced Container Isolation
- Configure HTTP Proxies
- Configure network settings
- Expose daemon on tcp://localhost:2375 without TLS Resources. This is applicable to Docker Desktop on Windows only.
- Enforce the use of WSL2 based engine or Hyper-V
- Configure Docker Engine
- Turn off Docker Desktop's ability to checks for updates
- Turn off Docker Desktop's ability to send usage statistics

For more details on the syntax and options administrators can set, see [Configure Admin Controls](configure-ac.md).

### How is this different to the `settings.json` file?

Using the `settings.json` file to pre-configure Docker Desktop settings means that developers own the `settings.json` file. They can therefore adjust any settings that their admins create, for example, modifying network and proxy controls. 

The `admin-settings.json` file can only be used by an administrator with root privileges and cannot be modified by users. 

### How do I set up and enforce Admin Controls?

As an administrator, you first need to [configure a registry.json to enforce sign-in](../../../docker-hub/configure-sign-in.md). This is because your Docker Desktop users must authenticate to your organization for this configuration to take effect.

Next, you must [create and configure the admin-settings.json file](configure-ac.md).

Once this is done, Docker Desktop users receive the changed settings when they next authenticate to your organization on Docker Desktop. We do not automatically mandate that developers re-authenticate once a change has been made, so as not to disrupt your developers workflow. 

### What do users see when the settings are enforced?

Docker Desktop users will see a notification in the **Settings**, or **Preferences** using a macOS, which states **Some settings are managed by your Admin**. 

Any settings that are enforced, are grayed out in Docker Desktop and the user is unable to edit them, either via the Docker Desktop UI, CLI, or by modifying the Docker Desktop Linux VM.




