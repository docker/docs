---
description: Instructions on how to set up enhanced container isolation
title: How does it work?
keywords: set up, enhanced container isolation, rootless, security
---

>Note
>
>Enhance Container Isolation is currently in [Early Access](../../../release-lifecycle.md#early-access-ea) and available to Docker Business customers only.

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







