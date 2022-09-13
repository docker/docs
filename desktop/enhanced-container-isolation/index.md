---
description: Enhanced Container Isolation - benefits, why use it, how it differs to Docker rootless, who it is for
keywords: containers, rootless, security, sysbox, runtime
title: What is Enhanced Container Isolation?
---

>Note
>
>Enhanced Container Isolation is currently in [Early Access](../../release-lifecycle.md#early-access-ea) and available to Docker Business customers only. 


## What it is
Allows Docker Desktop admins to lock-in configurations (e.g. Registry Access Management) such that they can’t be modified by Docker Desktop users. See below for more details on this.
Enhances container isolation using the Sysbox container runtime (see below for more info). This prevents containers from running as root inside the DD Linux VM and from potentially gaining control of it.


Enabling the feature prevents containers from running as root within Docker Desktop’s Linux VM and allows admins to lock-in sensitive security configs.


Prevent container attacks and vulnerabilities via Docker Desktop’s Hardened container runtime option,

Ensure stronger isolation, without any complex setups, using Docker Desktop’s Hardened container runtime option.

Enhanced Container Isolation provides an additional layer of security within Docker Desktop's Linux VM so there is strong container-to-host isolation. With Enhanced Container Isolation, Docker Desktop:

- Has a secure boot to prevent modification of Docker provided binaries pre-boot (e.g. docker engine, containerd, runc, etc)
- Prevents user containers from bypassing security controls and modifying system files.
- Prevents exposure of docker daemon on TCP without TLS

By taking advantage of Sysbox, it ensures containers run using the Linux user namespace and are not root in the VM

Developers can no longer:

Gain VM root access through privileged containers
Modify files before boot
Access the root console of the VM
Bind mount and modify system files
Escape containers

Prevent the use of privileged containers gaining root access to the Desktop VM and ensure stronger isolation (Linux user namespace, procfs & sysfs virtualization, mount locking, and more !) using Docker Desktop’s Hardened container runtime..

## What the benefits of it are

As a developer
When using Docker Desktop with the Hardened container runtime option enabled
I should be prevented from doing the following:
Running privileged containers to gain root access to the DD VM
Modifying files before boot
Accessing the root console of the VM
Bind mount and modifying system files
Escaping containers
I would add: "Modifying the config of the Docker Engine (and related components) from within DD containers".

Get more control over your local Docker Desktop instances using Docker’s Hardened container runtime.

## Who is it for: 

Problem 2 - Prevent exposure of docker daemon on TCP without TLS

As an IT admin working for a Docker Business customer, I am concerned that developers will be able to expose the docker daemon on TCP without TLS.


Problem 3 - Control mechanisms such as Registry Access Management are only designed to protect against well-intentioned developers making mistakes

As an IT admin at a Docker Business customer, I’m hesitant to adopt control features like Registry Access Management because it would be easy for a malicious actor within my org to override them by changing settings within Docker Desktop’s Linux VM.

Problem 4 - I need an easy, intuitive way to implement this control mechanism
As an IT admin at a Docker Business customer, I need an easy, intuitive way to implement the Hardened container runtime option on the machines of my developers.





## how does it work and how it differs to traditional rootless docker 

- Why this approach is advantageous as compared to traditional ‘rootless Docker' or ‘rootless mode’ in “other products”
    - workload compatibility, ease of use, etc. dive in on why Sysbox is awesome for both security and workloads

As such, we want to move to a model where the Docker Desktop user whose company has opted in to the Hardened container runtime option can still run all the containers that they expect, however they cannot gain root VM access through privileged containers, they cannot modify host system files, they are running in the user namespace and they cannot escape containers (bar kernel 0-day). These specific enhancements can be attained by integrating Sysbox, the secure container runtime created by Nestybox. 

Docker Desktop runs Docker Engine within a Linux VM, which provides strong isolation between containers and the underlying host machine (e.g. the Mac or Windows device running Docker Desktop). However, this does not prevent Docker Desktop users from launching a container that runs as root in the Docker Desktop Linux VM, or from using insecure privileged containers.
With root access to the Docker Desktop Linux VM, malicious users could potentially modify security policies of the Docker Engine and Docker Extensions as well as other control mechanisms like Registry Access Management policies and proxy configs. Moreover, whilst we have not yet seen anything of this nature, it is conceptually possible for malware in containers to read files on the users host machine, which presents an information leakage vulnerability.

Enhancing container isolation by ensuring that containers never run as root inside the Docker Desktop Linux VM, therefore preventing them from potentially gaining control of it. 
Ensuring sensitive configurations within the Docker Desktop VM cannot be mounted or modified from a container. This means that the Docker Engine, proxy settings and Registry Access configs can no longer be modified from within a container. They can only be set by the admins for your organization.


Sysbox is an alternative “runc” included in the Docker Business tier. It’s included alongside the standard OCI runc container runtime, which is the component that actually creates the containers using the Linux kernel’s namespaces, cgroups, and other features. 

What makes Sysbox different from the standard “runc” runtime is that it enhances container isolation by enabling the Linux user-namespace on all containers (i.e. root in the container maps to an unprivileged user at host level), and by vetting sensitive accesses between the container and the Linux kernel. This adds an extra layer of isolation between the container and the Linux kernel. 

This is all done under the covers, without requiring special container images and in a manner that is mostly transparent to Docker Desktop users.  



Normally, to run a container with Sysbox in Docker Desktop Business Tier, a user simply adds the --runtime=sysbox-runc flag to the docker run command. 

However, when Hardened Desktop is enabled a number of security features are activated (see above). One of these security features is that the Sysbox runtime is enforced for all user containers (e.g. the --runtime=sysbox-runc flag is implicitly set on all containers). This ensures all user containers run with the enhanced isolation offered by Sysbox. 



Currently, the Docker Engine runs inside a container on the DD Linux VM. 

Security-wise, there is no real isolation between the Docker Engine and the VM’s Linux kernel, because the Docker Engine runs as root with full capabilities inside a container that shares almost all namespaces with the VM’s root user (except the mount namespace). This gives the container access to all the VM’s kernel resources. This container is spawned by containerd + runc. 
As a result, DD users can easily gain privileged access to the DD VM (e.g., by running “docker run –privileged -it alpine”) from the host. This means DD users are one step closer to gaining privileged access to the underlying host (e.g., through the interfaces between the VM and the host).




