---
description: admin controls for desktop
keywords: admin, controls, rootless, 
title: What is Admin Controls?
--- 
>Note
>
>Admin Controls is currently in [Early Access](../../../release-lifecycle.md#early-access-ea) and available to Docker Business customers only. 

Admin Controls is a a feature that gives Enterprise admins the ability to enforce certain Docker Desktop settings for their organization. 

With just a few lines of JSON, admins are able to enforce preferences like HTTP proxies, Network settings and the Docker Engine configuration. This saves signficant time and cost in securing developer workflows.

### Who is it for? 

- For Organizations who wish to configure Docker Desktop to also be within their organizations centralized control.
- For Organizaitons who want to create a standardized Docker Desktop environment at scale.
- For security conscious Docker Business customers who want to confidently manage their use of Docker Desktop within tightly regulated environments.

### What can be set?

Using the `admin-settings.json` file, admins can:

- Enable Enhanced Container Isolation
- Configure HTTP Proxies
- Configure Network settings
- Expose daemon on tcp://localhost:2375 without TLS Resources (Windows only)
- Enforce the use of WSL2 based engine or Hyper-V
- Configure Docker Engine 
- Turning off checks for updates
- Turning off sending usage statistics

For more details on the syntax and options you can set, see [Configure Admin Controls](configure-ac.md).

### What do users see when the settings are enforced?

Docker Desktop users will see a notification in the **Settings**, or **Preferences** if macOS user, which states **Some settings are managed by your Admin**. 

Any settings that are enforced, are grayed out in Docker Desktop and the user is unable to edit them, either via the Docker Desktop UI, CLI, or by modifying the Docker Desktop Linux VM.

### How does this differ to the `settings.json` file?

Using the `settings.json` file to pre-configure Docker Desktop settings menas that developers own the settings.json file and can therefore adjust any settings that their admins create, for example, modifying network and proxy controls. 

The `admin-settings.json` file can only be used by an admin with root privileges and cannot be modified by users. 

### How do I set up and enforce Admin Controls?

As an Enterprise admin, you first need to [configure a registry.json to enforce sign-in](../../../docker-hub/configure-sign-in.md). This is because your Docker Desktop users must authenticate to your organization for this configuration to take effect.

Next, you must [create and configure the admin-settings.json file](configure-ac.md).

Once this is done, Docker Desktop users receive the changed settings when they next authenticate to your organization on Docker Desktop. We do not automatically mandate that developers re-authenticate once a change has been made, so as not to disrupt your developers workflow. 






