---
title: FAQs and Known issues
description: FAQ for Enhanced Container Isolation 
keywords: enhanced container isolation, security, faq, sysbox
---


## FAQ

### With Hardened Desktop enabled, can the user still override the --runtime flag from the CLI ?

No. With Hardened Desktop enabled, Sysbox is locked as the default (and only) runtime. If a user attempts to override the runtime by launching a container with the standard runc runtime (e.g. docker run --runtime=runc), container creation will fail. The reason runc is disallowed with Hardened Desktop is that it allows users to run as root on the Docker Desktop Linux VM, thereby providing them with implicit control of the VM and the ability to do things like modifying the administrative configurations for Docker Desktop.

### With Hardened Desktop enabled, can the user still use the “--privileged” flag from the CLI?

Yes, but by virtue of using Sysbox the container will only be privileged within its assigned Linux user-namespace. It will not be privileged within the Docker Desktop Linux VM. 

For example, the container’s init process will have all Linux capabilities enabled, have read/write access to the kernel’s /proc and /sys, run without system call or other restrictions normally imposed by Docker on regular containers (e.g., seccomp, AppArmor), and see all host devices under the container’s /dev directory. However, because Sysbox launches each container within a dedicated Linux user-namespace and vets sensitive accesses to the kernel, the container can only access resources assigned to it. For example, the container can’t access resources under /proc and /sys that are not namespaced. And though it can see all host devices under /dev, it won’t have permission to access them. Also, while the container can use system calls such as “mount” and “umount”, Sysbox will prevent the container from using them to modify the container’s chroot jail.

TODO: add table to clarify.

This makes running a privileged container with Hardened Desktop much safer than a privileged container launched with the standard runc, which offers almost no isolation.

### Why not just restrict usage of the “--privileged” flag in Hardened Desktop?

Privileged containers are typically used to run advanced workloads in containers (e.g., Docker-in-Docker), to perform kernel operations (e.g. loading modules) or to access hardware devices. We wish to allow the first within Hardened Desktop (i.e., running advanced workloads), yet deny the latter two.

By virtue of allowing the –privileged flag but restricting its impact within the container's user-namespace, it’s possible to do this.


I’ve heard that Docker Desktop’s settings can also be configured via a settings.json file ? What’s the difference between Admin Controls (which uses the admin-settings.json) and the original settings.json method ?

Some organizations currently use the settings.json file to pre-configure Docker Desktop settings for their users. The problem with this approach is that developers own the settings.json file and can therefore adjust any settings that their admins create (for example, modifying network and proxy controls). The admin-settings.json on the other hand, can only be used by an admin with root privileges and as such cannot be modified by users. This means that admins can lock in settings for their users via the admin-settings.json.


With Hardened Desktop enabled, can the user still override the --runtime flag from the CLI ?

No. With Hardened Desktop enabled, Docker’s hardened container runtime (using Sysbox) is locked as the default (and only) runtime. If a user attempts to override the runtime by launching a container with the standard runc runtime (e.g. docker run --runtime=runc), container creation will fail. The reason runc is disallowed with Hardened Desktop is that it allows users to run as root on the Docker Desktop Linux VM, thereby providing them with implicit control of the VM and the ability to do things like modifying the Admin Controls for Docker Desktop.

With Hardened Desktop enabled, can the user still use the --privileged flag from the CLI?

Yes, but by virtue of using Sysbox the container will only be privileged within its assigned Linux user-namespace. It will not be privileged within the Docker Desktop Linux VM. 

For example, the container’s init process will have all Linux capabilities enabled, have read/write access to the kernel’s /proc and /sys, run without system call or other restrictions normally imposed by Docker on regular containers (e.g. seccomp, AppArmor), and see all host devices under the container’s /dev directory. However, because Sysbox launches each container within a dedicated Linux user-namespace and vets sensitive accesses to the kernel, the container can only access resources assigned to it. For example, the container can’t access resources under /proc and /sys that are not namespaced. And though it can see all host devices under /dev, it won’t have permission to access them. Also, while the container can use system calls such as “mount” and “umount”, Sysbox will prevent the container from using them to modify the container’s chroot jail.

This makes running a privileged container with Hardened Desktop much safer than a privileged container launched with the standard runc, which offers almost no isolation.

Why not just restrict usage of the --privileged flag in Hardened Desktop ?

Privileged containers are typically used to run advanced workloads in containers (e.g. Docker-in-Docker), to perform kernel operations (e.g. loading modules) or to access hardware devices. We wish to allow the first within Hardened Desktop (e.g. running advanced workloads), yet deny the latter two. By virtue of allowing the –privileged flag but restricting its impact within the container's user-namespace, it’s possible to do this.








## Known issues

Known issues?
If in DD “secure mode” all containers are launched with Sysbox, then users may experience some differences between running a container in DD and running that same container in production, because in production the container may run on another runtime (typically the OCI runc). 


Kernel Day-0 Vulnerabilities
Sysbox can’t protect against kernel day-0 vulnerabilities (e.g., flaws in user-namespace isolation). There have been a few of these recently, but fortunately they are patched pretty quickly in the Linux kernel.
Nested virtualization
Sysbox is not a solution for the problem of running DD inside VMs (which currently requires nested virtualization). Rather Sysbox adds a layer of isolation by running Docker more securely (i.e., without root privileges on the VM).
Docker Engine Limitations
When running Docker inside a Sysbox container (e.g., for extra isolation), most Docker functionality is supported. However, there may be some advanced Docker functionality that does not currently work as the environment inside the Sysbox container does not yet fully resemble that of a bare-metal machine or VM. Fixing this requires further changes in Sysbox.

