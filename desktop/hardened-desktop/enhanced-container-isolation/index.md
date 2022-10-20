---
description: Enhanced Container Isolation - benefits, why use it, how it differs to Docker rootless, who it is for
keywords: containers, rootless, security, sysbox, runtime
title: What is Enhanced Container Isolation?
---

>**Note**
>
>Enhanced Container Isolation is available to Docker Business customers only. 

Enhanced Container Isolation provides an additional layer of security that uses a variety of advanced techniques to harden container isolation without impacting developer productivity. It is available with [Docker Desktop 4.13.0 or later](../../release-notes.md).

These techniques include:
- Running all containers unprivileged through the Linux user-namespace.
- Restricting containers from modifying Docker Desktop VM settings.
- Vetting some critical system calls to prevent container escapes, and partially virtualizing portions of `/proc` and `/sys` inside the container for further isolation. 
- Preventing console access to the Docker Desktop VM.

This is done automatically and with minimal functional or performance impact. 

Enhanced Container Isolation helps ensure strong container isolation and also locks in any security configurations that have been created, for instance through [Registry Access Management policies](../registry-access-management.md) or with [Settings Management](../settings-management/index.md). 

>**Note**
>
> Enhanced Container Isolation is in addition to other container security techniques used by Docker. For example, reduced Linux Capabilities, Seccomp, AppArmor.

### Who is it for?

- For organizations that want to prevent container attacks and reduce vulnerabilities.
- For organizations that want to ensure stronger container isolation that is easy and intuitive to implement on developers' machines.

### What happens when Enhanced Container Isolation is enabled?

When Enhanced Container Isolation is enabled using [Settings Management](../settings-management/index.md), the following features are enabled: 

- All user containers are automatically run in Linux User Namespaces which ensures stronger isolation.
- The root user in the container maps to an unprivileged user at VM level.
- Users can continue using containers as usual, including bind mounting host directories, volumes, networking configurations, etc.
- Privileged containers work, but they are only privileged within the container's Linux User Namespace, not in the Docker Desktop VM.
- Containers can no longer share namespaces with the Docker Desktop VM. For example, `--network=host`, `--pid=host`.
- Containers can no longer modify configuration files in the Docker Desktop VM.
- Console access to the Desktop VM is forbidden for all users.
- Containers become harder to breach. For example, sensitive system calls are vetted and portions of `/proc` and `/sys` are emulated.

For more information on how Enhanced Container Isolation work, see [How does it work](how-eci-works.md).

>**Important**
>
>Enhanced Container Isolation is currently incompatible with WSL and does not protect Kubernetes pods. For more information on known limitations and workarounds, see [FAQs and known issues](faq.md).
{: .important}

### How do I enable Enhanced Container Isolation?

As an admin, you first need to [configure a `registry.json` file to enforce sign-in](../../../docker-hub/configure-sign-in.md). This is because the Enhanced Container Isolation feature requires a Docker Business subscription and therefore your Docker Desktop users must authenticate to your organization for this configuration to take effect.

Next, you must [create and configure the `admin-settings.json` file](../settings-management/configure.md) and specify:

```JSON
{
 "configurationFileVersion": 2,
 "enhancedContainerIsolation": {
    "value": true,
    "locked": true
    }
}
```

For this to take effect:

- On a new install, developers need to launch Docker Desktop and authenticate to their organization.
- On an existing install, developers need to quit Docker Desktop through the Docker menu, and then relaunch Docker Desktop. If they are already signed in, they don’t need to sign in again for the changes to take effect.

>Important
  >
  >Selecting **Restart** from the Docker menu isn't enough as it only restarts some components of Docker Desktop.
  {: .important}

### What do users see when this setting is enforced?

When Enhanced Container Isolation is enabled, users see that containers run within a Linux user namespace. 

To check, run:

```
$ docker run --rm alpine cat /proc/self/uid_map 
```

The following output displays:

```
         0     100000      65536
```

This indicates that the container's root user (0) maps to unprivileged user (100000) in the Docker Desktop VM, and that the mapping extends for a range of 64K user-IDs.

In contrast, without Enhanced Container Isolation the Linux user namespace is not used, the following displays:

```
         0          0 4294967295
```

This means that the root user in the container (0) is in fact the root user in the Docker Desktop VM (0) which reduces container isolation. The user-ID mapping varies with each new container, as each container gets an exclusive range of host User-IDs for isolation. User-ID mapping is automatically managed by Docker Desktop.

With Enhanced Container Isolation, if a process were to escape the container, it would find itself without privileges at the VM level. For further details, see [How Enhanced Container Isolation works](how-eci-works.md).

Since Enhanced Container Isolation [uses the Sysbox container runtime](how-eci-works.md) embedded in the Docker Desktop Linux VM, another way to determine if a container is running with Enhanced Container Isolation is by using `docker inspect`:

{% highlight liquid %}
docker inspect --format={% raw %}'{{.HostConfig.Runtime}}'{% endraw %} my_container
{% endhighlight %}

It outputs:

```
sysbox-runc
```

Without Enhanced Container Isolation, `docker inspect` outputs `runc`, which is the standard OCI runtime.
