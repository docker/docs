---
description: Enhanced Container Isolation - benefits, why use it, how it differs to Docker rootless, who it is for
keywords: containers, rootless, security, sysbox, runtime
title: What is Enhanced Container Isolation?
---

>Note
>
>Enhanced Container Isolation is currently in [Early Access](../../release-lifecycle.md#early-access-ea) and available to Docker Business customers only. 

Enhanced Container Isolation provides an additional layer of security was it which prevents containers from running as root in the Docker Desktop Linux VM. This ensures a strong container-to-host isolation and locks in any security configurations that have been created, for instance Registry Access Management policies and Admin Controls. 

With Enhanced Container Isolation, Docker Desktop:

- Has a secure boot to prevent modification of Docker provided binaries pre-boot (e.g. docker engine, containerd, runc, etc)
- Prevents user containers from bypassing security controls and modifying system files.
- Prevents exposure of docker daemon on TCP without TLS

It also means all containers run unprivileged in the Docker Desktop Linux VM, in user namespaces. Root access to the Linux VM is removed, privileged containers cannot be run and there is no access to the host namespaces. As a result, it becomes impossible for users to alter Admin Controls via containers. 

When Enhanced Container Isolation is enabled using [Admin Controls](../admin-controls/index.md), developers can no longer:

- Gain VM root access through privileged containers
- Modify files before boot
- Modify the config of the Docker Engine (and related components) from within Docker Desktop containers
- Access the root console of the VM
- Bind mount and modify system files
- Escape containers

For more information on how Enhanced Container Isolation work, see [How does it work?](how-eci-works.md).

### Who is it for?

- For organizations that want to prevent container attacks and reduce vulnerabilities.
- For organizations that want ensure stronger container isolation that is easy and intuitive to implement on developers' machines.

### How to enable Enahnced Container Isolation

As an Enterprise admin, you first need to [configure a registry.json to enforce sign-in](../../../docker-hub/configure-sign-in.md). This is because your Docker Desktop users must authenticate to your organization for this configuration to take effect.

Next, you must [create and configure the admin-settings.json file](configure-ac.md) and specify:

```JSON
{
 "enhancedContainerIsolation": {
    "value": true,
    "locked": true
    }
}
```

Once this is done, Docker Desktop users receive the changed settings when they next authenticate to your organization on Docker Desktop. We do not automatically mandate that developers re-authenticate once a change has been made, so as not to disrupt your developers workflow. 

### What do users see when the settings are enforced?

Docker Desktop users will see a notification in the **Settings**, or **Preferences** if macOS user, which states **Some settings are managed by your Admin**. 

Any settings that are enforced, are grayed out in Docker Desktop and the user is unable to edit them, either via the Docker Desktop UI, CLI, or by modifying the Docker Desktop Linux VM.
