---
description: Enhanced Container Isolation - benefits, why use it, how it differs to Docker rootless, who it is for
keywords: containers, rootless, security, sysbox, runtime
title: What is Enhanced Container Isolation?
---

>Note
>
>Enhanced Container Isolation is available to Docker Business customers only. 

Enhanced Container Isolation provides an additional layer of security that uses a variety of advanced techniques to harden container isolation without impacting developer productivity. 

These techniques include:
- Running all containers unprivileged through the Linux user-namespace
- Restricting containers from modifying Docker Desktop VM settings
- Vetting some critical system calls to prevent container escapes, and partially virtualizing portions of `/proc` and `/sys` inside the container for further isolation. 

This is done automatically and with minimal performance impact. 

Enhanced Container Isolation helps ensure a strong container-to-host isolation and locks in any security configurations that have been created, for instance through [Registry Access Management policies](../registry-access-management.md) or with [Admin Controls](../admin-controls/index.md). 

>Note
>
> Enhanced Container Isolation is in addition to other security techniques used by Docker. For example, reduced Linux Capabilities, Seccomp, AppArmor.

### Who is it for?

- For organizations that want to prevent container attacks and reduce vulnerabilities.
- For organizations that want to ensure stronger container isolation that is easy and intuitive to implement on developers' machines.

### What happens when Enhanced Container Isolation is switched on?

When Enhanced Container Isolation is enabled using [Admin Controls](../admin-controls/index.md), the following features are enabled: 

- All user containers are automatically run in Linux User Namespaces which ensures stronger isolation.
- The root user in the container maps to an unprivileged user at VM level.
- Users can continue using containers as usual, including bind-mounting host directories, volumes, networking configurations, etc.
- Privileged containers work, but they are only privileged within the container's Linux User Namespace, not in the Docker Desktop VM.
- Containers can no longer share namespaces with the Docker Desktop VM. For example, `--network=host`, `--pid=host`.
- Containers can no longer modify configuration files in the Docker Desktop VM.
- Containers become harder to breach. For example, sensitive system calls are vetted and portions of `/proc` and `/sys` are emulated.

For more information on how Enhanced Container Isolation work, see [How does it work?](how-eci-works.md).

>Important
>
>Enhanced Container Isolation is currently incompatible with WSL and does not protect Kubernetes pods. For more information on known limitations and workarounds, see [FAQS and known issues](faq.md).
{: .important}

### How do I switch on Enhanced Container Isolation?

As an admin, you first need to [configure a `registry.json` file to enforce sign-in](../../../docker-hub/configure-sign-in.md). This is because your Docker Desktop users must authenticate to your organization for this configuration to take effect.

Next, you must [create and configure the `admin-settings.json` file](configure-ac.md) and specify:

```JSON
{
 "enhancedContainerIsolation": {
    "value": true,
    "locked": true
    }
}
```

Once this is done, developers need to either quit, re-launch, and sign in to Docker Desktop, or launch and sign in to Docker Desktop for the first time.

### What do users see when this setting is enforced?

When Enhanced Container Isolation is enabled, users see that containers run within a Linux user-namespace. For example:

```
$ docker run -it --rm alpine
/ # cat /proc/self/uid_map 
         0     100000      65536
```

This indicates that the container's root user (0) maps to unprivileged user (100000) in the Docker Desktop VM, and that the mapping extends for a range of 64K user-IDs.

In contrast, without Enhanced Container Isolation the Linux user-namespace is not used:

```
$ docker run -it --rm alpine             
/ # cat /proc/self/uid_map                           
         0          0 4294967295
```

This means that the root user in the container (0) is in fact the root user in the Docker Desktop VM (0), reducing container isolation.
