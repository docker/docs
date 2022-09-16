---
description: Enhanced Container Isolation - benefits, why use it, how it differs to Docker rootless, who it is for
keywords: containers, rootless, security, sysbox, runtime
title: What is Enhanced Container Isolation?
---

>Note
>
>Enhanced Container Isolation is currently in [Early Access](../../release-lifecycle.md#early-access-ea) and available to Docker Business customers only. 

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


## What it is

Wait, that last point doesn’t make sense ! I thought Admin Controls would automatically lock any security controls I create, but you’re saying that Enhanced Container Isolation is also required to prevent containers from modifying these ?
Using Admin Controls, your Docker Business admin can lock in Docker Desktop settings such as HTTP Proxies and Network settings. This means that users of Docker Desktop will then have no path within Docker Desktop to change these settings (e.g. via the user interface or CLI).
However, malicious code in a container could still potentially modify these controls without the developer knowing. Enhanced Container Isolation is an extra layer of security that prevents containers from modifying any Admin Controls or security policies, so that admins have complete peace of mind that their settings are enforced.


Containers will no longer run as root inside the Docker Desktop Linux VM and will instead run using the Linux user namespace.
As a result, user containers will be unable to modify any security configurations created by your Docker admins (e.g. Registry Access Management policies and Admin Controls).


Enhanced Container Isolation is a feature that admins can enable, which prevents containers from running as root in the Docker Desktop Linux VM. 
Allows Docker Desktop admins to lock-in configurations (e.g. Registry Access Management) such that they can’t be modified by Docker Desktop users. See below for more details on this.
Enhances container isolation using the Sysbox container runtime (see below for more info). This prevents containers from running as root inside the DD Linux VM and from potentially gaining control of it.


With Enhanced Container Isolation enabled, all containers run unprivileged in the Docker Desktop Linux VM, in user namespaces. Root access to the Linux VM is removed, privileged containers cannot be run and there is no access to the host namespaces. As a result, it becomes impossible for users to alter Admin Controls via containers. 

Enabling the feature prevents containers from running as root within Docker Desktop’s Linux VM and allows admins to lock-in sensitive security configs.


Prevent container attacks and vulnerabilities via Docker Desktop’s Hardened container runtime option,

Ensure stronger isolation, without any complex setups, using Docker Desktop’s Hardened container runtime option.



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



This page contains information on how Enterprise admins can enable Enhanced Container Isolation to 



How to configure it if you are an admin

### What do users see when the settings are enforced?

## How to enable/ get ECI
(e.g. currently developers in Docker Business customers, requires authentication, etc)

requires an Apply and restart
- Admins can lock in the use of the ‘Enhanced container isolation’ mode within their org via the ‘Admin Controls’ feature <link to Admin Controls docs>

To enable Hardened Docker Desktop, Docker Business administrators simply have to toggle on the ‘Hardened Desktop’ option within the Settings panel of their Organization’s space on Docker Hub. Your developers must then authenticate to your organization in Docker Desktop for the settings to be applied. You can follow this simple guide for ensuring developers authenticate to your organization before using Docker Desktop.

How do I enable Enhanced Container Isolation for my organization ?

In the admin-settings.json specify “enhancedContainerIsolation”: true as per the below image. 



You must then place this file on your developers machines in the following locations:

Mac - <here>
Windows - <here>
Linux - <here> 

As mentioned above, the Hardened Desktop security model is designed for organizations that don't give root/admin access to their developers on their machines. By placing this file in the above protected directories, end users will be unable to modify it. We also assume that said organizations have the ability to push this settings file to the locations specified above via device management software such as Jamf.

Important - Your Docker Desktop users must then authenticate to your organization for this configuration to take effect. You can configure the registry.json file to enforce sign in.






