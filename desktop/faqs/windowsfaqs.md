---
description: Frequently asked questions
keywords: desktop, windows, faqs
title: Frequently asked questions for Windows
---


## Windows FAQs

### Can I use VirtualBox alongside Docker Desktop?

Yes, you can run VirtualBox along with Docker Desktop if you have enabled the [Windows Hypervisor Platform](https://docs.microsoft.com/en-us/virtualization/api/){: target="_blank" rel="noopener" class="_"} feature on your machine.

### Why is Windows 10 or Windows 11 required?

Docker Desktop uses the Windows Hyper-V features. While older Windows versions have Hyper-V, their Hyper-V implementations lack features critical for Docker Desktop to work.

### Can I install Docker Desktop on Windows 10 Home?

If you are running Windows 10 Home (starting with version 1903), you can install [Docker Desktop for Windows](https://hub.docker.com/editions/community/docker-ce-desktop-windows/){: target="_blank" rel="noopener" class="_"} with the [WSL 2 backend](../windows/wsl.md).

### Can I run Docker Desktop on Windows Server?

No, running Docker Desktop on Windows Server is not supported.

### How do I run Windows containers on Windows Server?

You can install a native Windows binary which allows you to develop and run
Windows containers without Docker Desktop. For more information, see the tutorial about running Windows containers on Windows Server in
[Getting Started with Windows Containers](https://github.com/docker/labs/blob/master/windows/windows-containers/README.md){: target="_blank" rel="noopener" class="_"}.

### Why do I see the `Docker Desktop Access Denied` error message when I try to start Docker Desktop?

Docker Desktop displays the **Docker Desktop - Access Denied** error if a Windows user is not part of the **docker-users** group.

If your admin account is different to your user account, add the **docker-users** group. Run **Computer Management** as an administrator and navigate to **Local Users* and Groups** > **Groups** > **docker-users**.

Right-click to add the user to the group. Log out and log back in for the changes to take effect.

### Why does Docker Desktop fail to start when anti-virus software is installed?

Some anti-virus software may be incompatible with Hyper-V and Windows 10 builds which impact Docker
Desktop. For more information, see [Docker Desktop fails to start when anti-virus software is installed](../windows/troubleshoot.md#docker-desktop-fails-to-start-when-anti-virus-software-is-installed).

### Can I change permissions on shared volumes for container-specific deployment requirements?

Docker Desktop does not enable you to control (`chmod`)
the Unix-style permissions on [shared volumes](../windows/index.md#file-sharing) for
deployed containers, but rather sets permissions to a default value of
[0777](http://permissions-calculator.org/decode/0777/){: target="_blank" rel="noopener" class="_"}
(`read`, `write`, `execute` permissions for `user` and for
`group`) which is not configurable.

For workarounds and to learn more, see
[Permissions errors on data directories for shared volumes](../windows/troubleshoot.md#permissions-errors-on-data-directories-for-shared-volumes).

### How do symlinks work on Windows?

Docker Desktop supports two types of symlinks: Windows native symlinks and symlinks created inside a container.

The Windows native symlinks are visible within the containers as symlinks, whereas symlinks created inside a container are represented as [mfsymlinks](https://wiki.samba.org/index.php/UNIX_Extensions#Minshall.2BFrench_symlinks){:target="_blank" rel="noopener" class="_"}. These are regular Windows files with a special metadata. Therefore the symlinks created inside a container appear as symlinks inside the container, but not on the host.

### File sharing with Kubernetes and WSL 2

Docker Desktop mounts the Windows host filesystem under `/run/desktop` inside the container running Kubernetes.
See the [Stack Overflow post](https://stackoverflow.com/questions/67746843/clear-persistent-volume-from-a-kubernetes-cluster-running-on-docker-desktop/69273405#69273){:target="_blank" rel="noopener" class="_"} for an example of how to configure a Kubernetes Persistent Volume to represent directories on the host.
