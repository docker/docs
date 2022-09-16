---
description: Instructions on how to set up enhanced container isolation
title: How does it work?
keywords: set up, enhanced container isolation, rootless, security
---

>Note
>
>Enhance Container Isolation is currently in [Early Access](../../../release-lifecycle.md#early-access-ea) and available to Docker Business customers only.

Enhanced Container Isolation takes advantage of the recent integration of Sysbox, the secure container runtime created by [Nestybox](https://www.nestybox.com/). 

Sysbox is an alternative “runc” included in the Docker Business tier. It’s included alongside the standard OCI runc container runtime, which is the component that actually creates the containers using the Linux kernel’s namespaces, cgroups, and other features. What makes Sysbox different from the standard “runc” runtime is that it enhances container isolation by enabling the Linux user-namespace on all containers (i.e. root in the container maps to an unprivileged user at host level), and by vetting sensitive accesses between the container and the Linux kernel. This adds an extra layer of isolation between the container and the Linux kernel. 

Without Enhanced Container Isolation Docker Desktop runs Docker Engine within a Linux VM, which provides strong isolation between containers and the underlying host machine. However, this does not prevent Docker Desktop users from launching a container that runs as root in the Docker Desktop Linux VM, or from using insecure privileged containers.

With root access to the Docker Desktop Linux VM, malicious users could potentially modify security policies of the Docker Engine and Docker Extensions as well as other control mechanisms like Registry Access Management policies and proxy configs. Moreover, whilst we have not yet seen anything of this nature, it is conceptually possible for malware in containers to read files on the users host machine, which presents an information leakage vulnerability.

### How is this different to rootless mode in Docker Engine?
