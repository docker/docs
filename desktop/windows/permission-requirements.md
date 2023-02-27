---
description: Permission requirements for Docker Desktop for Windows
keywords: Docker Desktop, Windows, security, install
title: Understand permission requirements for Windows
redirect_from:
- /desktop/windows/privileged-helper/
---

This page contains information about the permission requirements for running and installing Docker Desktop on Windows, the functionality of the privileged helper process `com.docker.service.exe` and the reasoning behind this approach. 

It also provides clarity on running containers as `root` as opposed to having `Administrator` access on the host and the privileges of the Windows Docker engine and Windows containers.

## Permission requirements

While Docker Desktop on Windows can be run without having `Administrator` privileges, it does require them during installation. On installation the user gets a UAC prompt which allows a privileged helper service to be installed. After that, Docker Desktop can be run by users without administrator privileges, provided they are members of the `docker-users` group. The user who performs the installation is automatically added to this group, but other users must be added manually. This allows the administrator to control who has access to Docker Desktop.

The reason for this approach is that Docker Desktop needs to perform a limited set of privileged operations which are conducted by the privileged helper process `com.docker.service.exe`. This approach allows, following the principle of least privilege, `Administrator` access to be used only for the operations for which it is absolutely necessary, while still being able to use Docker Desktop as an unprivileged user.

## Privileged Helper

The privileged helper `com.docker.service.exe` is a Windows service which runs in the background with `SYSTEM` privileges. It listens on the named pipe `//./pipe/dockerBackendV2`. The developer runs the Docker Desktop application, which connects to the named pipe and sends commands to the service. This named pipe is protected, and only users that are part of the `docker-users` group can have access to it.

The service performs the following functionalities:
- Ensuring that `kubernetes.docker.internal` is defined in the Win32 hosts file. Defining the DNS name `kubernetes.docker.internal` allows Docker to share Kubernetes contexts with containers.
- Securely caching the Registry Access Management policy which is read-only for the developer.
- Creating the Hyper-V VM `"DockerDesktopVM"` and managing its lifecycle - starting, stopping and destroying it. The VM name is hard coded in the service code so the service cannot be used for creating or manipulating any other VMs.
- Getting the VHDX disk size.
- Moving the VHDX file or folder.
- Starting and stopping the Windows Docker engine and querying whether it is running.
- Deleting all Windows containers data files.
- Checking if Hyper-V is enabled.
- Checking if the bootloader activates Hyper-V.
- Checking if required Windows features are both installed and enabled.
- Conducting healthchecks and retrieving the version of the service itself.

## Containers running as root within the Linux VM

The Linux Docker daemon and containers run in a minimal, special-purpose Linux VM managed by Docker. It is immutable so users can’t extend it or change the installed software. 
This means that although containers run by default as `root`, this does not allow altering the VM and does not grant `Administrator` access to the Windows host machine. The Linux VM serves as a security boundary and limits what resources from the host can be accessed. File sharing uses a user-space crafted file server and any directories from the host bind mounted into Docker containers still retain their original permissions. It does not give the user access to any files that it doesn’t already have access to.

## Windows Containers

Unlike the Linux Docker engine and containers which run in a VM, Windows containers are an operating system feature, and run directly on the Windows host with `Administrator` privileges. For organizations which do not want their developers to run Windows containers, a `–no-windows-containers` installer flag is available from version 4.11 to disable their use.

## Networking

For network connectivity, Docker Desktop uses a user-space process (`vpnkit`), which inherits constraints like firewall rules, VPN, HTTP proxy properties etc. from the user that launched it.



