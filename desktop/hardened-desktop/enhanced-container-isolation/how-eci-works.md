---
description: Instructions on how to set up enhanced container isolation
title: How does it work?
keywords: set up, enhanced container isolation, rootless, security
---

>**Note**
>
>Enhance Container Isolation is available to Docker Business customers only.

Enhanced Container Isolation takes advantage of the recent integration of Sysbox, the secure container runtime created by [Nestybox](https://www.nestybox.com/). 

Sysbox is an alternative `runc` used to create a container using the Linux kernel’s namespaces, cgroups, and other features. 

Unlike the standard `runc` runtime, Sysbox enhances container isolation by using techniques such as enabling the Linux user-namespace on all containers, emulating portions of the `proc` filesystem and `sysfs` inside the container and vetting sensitive accesses between the container and the Linux kernel. This adds an extra layer of isolation between the container and the Linux kernel. 

Without Enhanced Container Isolation, Docker Desktop has Docker Engine run as root with full capabilities inside a container that shares almost all namespaces with the Linux VM’s root user. Whilst this provides strong isolation between containers and the underlying host machine, it gives the container access to all the VM’s kernel resources and does not prevent Docker Desktop users from launching a container that runs as root in the Docker Desktop Linux VM, or from using insecure privileged containers. This brings Docker Desktop users closer to gaining privileged access to the underlying host.

### How is this different to rootless mode in Docker Engine?


