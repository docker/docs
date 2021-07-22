---
description: Frequently asked questions
keywords: desktop, mac, windows, faqs
redirect_from:
- /mackit/faqs/
- /docker-for-mac/faqs/
- /docker-for-windows/faqs/
title: Frequently asked questions
toc_max: 2
---

## General

### What are the system requirements for Docker Desktop?

For information about Docker Desktop system requirements, see [Docker Desktop for Mac system requirements](../docker-for-mac/install.md#system-requirements) and [Docker Desktop for Windows system requirements](../docker-for-windows/install.md#system-requirements).

### What is an experimental feature?

{% include experimental.md %}

### Where can I find information about diagnosing and troubleshooting Docker Desktop issues?

You can find information about diagnosing and troubleshooting common issues in the Troubleshooting topic. See [Mac Logs and Troubleshooting](../docker-for-mac/troubleshoot.md) topic and Windows Logs and [Windows Logs and Troubleshooting](../docker-for-windows/troubleshoot.md).

If you do not find a solution in Troubleshooting, browse issues on
[docker/for-mac](https://github.com/docker/for-mac/issues){: target="_blank" rel="noopener" class="_"} or [docker/for-win](https://github.com/docker/for-win/issues){: target="_blank" rel="noopener" class="_"} GitHub repository, or create a new one.

### How do I connect to the remote Docker Engine API?

To connect to the remote Engine API, you might need to provide the location of the Engine API for Docker clients and development tools.

Mac and Windows WSL 2 users can connect to the Docker Engine through a Unix socket: `unix:///var/run/docker.sock`.

If you are working with applications like [Apache Maven](https://maven.apache.org/){: target="_blank" rel="noopener" class="_"}
that expect settings for `DOCKER_HOST` and `DOCKER_CERT_PATH` environment
variables, specify these to connect to Docker instances through Unix sockets.

For example:

```bash
export DOCKER_HOST=unix:///var/run/docker.sock
```

Docker Desktop Windows users can connect to the Docker Engine through a **named pipe**: `npipe:////./pipe/docker_engine`, or **TCP socket** at this URL:
`tcp://localhost:2375`.

For details, see [Docker Engine API](../engine/api/index.md).

### How do I connect from a container to a service on the host?

Both Mac and Windows have a changing IP address (or none if you have no network access). On both Mac and Windows, we recommend that you connect to the special DNS name `host.docker.internal`, which resolves to the internal IP address used by the host. This is for development purposes and does not work in a production environment outside of Docker Desktop.

For more information and examples, see how to connect from a container to a service on the host
[on Mac](../docker-for-mac/networking.md#i-want-to-connect-from-a-container-to-a-service-on-the-host) and [on Windows](../docker-for-windows/networking.md#i-want-to-connect-from-a-container-to-a-service-on-the-host).

### How do I connect to a container from Mac or Windows?

We recommend that you publish a port, or connect from another container. Port forwarding works for `localhost`; `--publish`, `-p`, or `-P` all work.

For more information and examples, see
[I want to connect to a container from Mac](../docker-for-mac/networking.md#i-want-to-connect-to-a-container-from-the-mac) and [I want to connect to a container from Windows](../docker-for-windows/networking.md#i-want-to-connect-to-a-container-from-the-mac).

### How do I add custom CA certificates?

Docker Desktop supports all trusted certificate authorities (CAs) (root or intermediate). For more information on adding server and client side certs, see
[Add TLS certificates on Mac](../docker-for-mac/index.md#add-tls-certificates) and [Add TLS certificates on Windows](../docker-for-windows/index.md#adding-tls-certificates).

### Can I pass through a USB device to a container?

Unfortunately, it is not possible to pass through a USB device (or a
serial port) to a container as it requires support at the hypervisor level.

### Can I run Docker Desktop in nested virtualization scenarios?

Docker Desktop can run inside a Windows 10 VM running on apps like Parallels or
VMware Fusion on a Mac provided that the VM is properly configured. However,
problems and intermittent failures may still occur due to the way these apps
virtualize the hardware. For these reasons, **Docker Desktop is not supported in
nested virtualization scenarios**. It might work in some cases, and not in others.

For more information, see [Running Docker Desktop in nested virtualization scenarios](../docker-for-windows/troubleshoot.md#running-docker-desktop-in-nested-virtualization-scenarios).

### Docker Desktop's UI appears green, distorted, or has visual artifacts. How do I fix this?

Docker Desktop uses hardware-accelerated graphics by default, which may cause problems for some GPUs. In such cases, 
Docker Desktop will launch successfully, but some screens may appear green, distorted, 
or have some visual artifacts.

To work around this issue, disable hardware acceleration by creating a `"disableHardwareAcceleration": true` entry in Docker Desktop's settings file. This file can be found at:

- **Mac**: `~/Library/Group Containers/group.com.docker/settings.json`
- **Windows**: `C:\Users\[USERNAME]\AppData\Roaming\Docker\settings.json`
 
Completely close and restart Docker Desktop to apply the changes.

## Releases

### When will Docker Desktop move to a cumulative release stream?

Starting with version 3.0.0, Docker Desktop will be available as a single, cumulative release stream. This is the same version for both Stable and Edge users. The next release after Docker Desktop 3.0.0 will be the first to be applied as a delta update. For more information, see [Automatic updates](../docker-for-mac/install.md#automatic-updates).

### How do new users install Docker Desktop?

Each Docker Desktop release is also delivered as a full installer for new users. The same will apply if you have skipped a version, although this doesn't normally happen as updates will be applied automatically.

### How frequent will new releases be?

New releases will be available roughly monthly, similar to Edge today, unless there are critical fixes that need to be released sooner.

### How do I ensure that all users on my team are using the same version?

Previously you had to manage this yourself: now it will happen automatically as a side effect of all users being on the latest version.

### My colleague has got a new version but I haven’t got it yet.

Sometimes we may roll out a new version gradually over a few days. Therefore, if you wait, it will turn up soon. Alternatively, you can select **Check for Updates** from the Docker menu to jump the queue and get the latest version immediately.

### Where can I find information about Stable and Edge releases?

Starting with Docker Desktop 3.0.0, Stable and Edge releases are combined into a single, cumulative release stream for all users.

## Support

### Does Docker Desktop offer support?

Yes, Docker Desktop offers support for Pro and Team users. For more information, see [Docker Desktop Support](../docker-for-mac/troubleshoot.md#support).

For information about the pricing plans and to upgrade your existing account, see [Docker pricing](https://www.docker.com/pricing){: target="_blank" rel="noopener" class="_"}.

### What kind of feedback are you looking for?

Everything is fair game. We'd like your impressions on the download-install
process, startup, functionality available, the GUI, usefulness of the app,
command line integration, and so on. Tell us about the issues you are experiencing, what you like, or request a new feature through our public [Docker Roadmap](https://github.com/docker/roadmap){: target="_blank" rel="noopener" class="_"}.

### How is personal data handled in Docker Desktop?

When uploading diagnostics to help Docker with investigating issues, the uploaded diagnostics bundle may contain personal data such as usernames and IP addresses. The diagnostics bundles are only accessible to Docker, Inc.
employees who are directly involved in diagnosing Docker Desktop issues.

By default, Docker, Inc. will delete uploaded diagnostics bundles after 30 days. You may also request the removal of a diagnostics bundle by either specifying the diagnostics ID or via your GitHub ID (if the diagnostics ID is mentioned in a GitHub issue). Docker, Inc. will only use the data in the diagnostics bundle to investigate specific user issues, but may derive high-level (non personal) metrics such as the rate of issues from it.

For more information, see [Docker Data Processing Agreement](https://www.docker.com/legal/data-processing-agreement){: target="_blank" rel="noopener" class="_"}.

## Mac FAQs

### What is Docker.app?

`Docker.app` is Docker Desktop on Mac. It bundles the Docker client and Docker Engine. `Docker.app` uses the macOS Hypervisor.framework to run containers.

### Is Docker Desktop compatible with Apple silicon processors?

Yes, you can now install Docker Desktop for Mac on Apple silicon. For more information, see [Docker Desktop for Apple silicon](../docker-for-mac/apple-silicon.md).

### What is HyperKit?

HyperKit is a hypervisor built on top of the Hypervisor.framework in macOS. It runs entirely in userspace and has no other dependencies.

We use HyperKit to eliminate the need for other VM products, such as Oracle
VirtualBox or VMWare Fusion.

### What is the benefit of HyperKit?

HyperKit is thinner than VirtualBox and VMWare fusion, and the version we include is customized for Docker workloads on Mac.

### Why is com.docker.vmnetd still running after I quit the app?

The privileged helper process `com.docker.vmnetd` is started by `launchd` and
runs in the background. The process does not consume any resources unless
Docker.app connects to it, so it's safe to ignore.

## Windows FAQs

### Can I use VirtualBox alongside Docker Desktop?

Yes, you can run VirtualBox along with Docker Desktop if you have enabled the [Windows Hypervisor Platform](https://docs.microsoft.com/en-us/virtualization/api/){: target="_blank" rel="noopener" class="_"} feature on your machine.

### Why is Windows 10 required?

Docker Desktop uses the Windows Hyper-V features. While older Windows versions have Hyper-V, their Hyper-V implementations lack features critical for Docker Desktop to work.

### Can I install Docker Desktop on Windows 10 Home?

If you are running Windows 10 Home (starting with version 1903), you can install [Docker Desktop for Windows](https://hub.docker.com/editions/community/docker-ce-desktop-windows/){: target="_blank" rel="noopener" class="_"} with the [WSL 2 backend](../docker-for-windows/wsl.md).

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
Desktop. For more information, see [Docker Desktop fails to start when anti-virus software is installed](../docker-for-windows/troubleshoot.md#docker-desktop-fails-to-start-when-anti-virus-software-is-installed).

### Can I change permissions on shared volumes for container-specific deployment requirements?

Docker Desktop does not enable you to control (`chmod`)
the Unix-style permissions on [shared volumes](../docker-for-windows/index.md#file-sharing) for
deployed containers, but rather sets permissions to a default value of
[0777](http://permissions-calculator.org/decode/0777/){: target="_blank" rel="noopener" class="_"}
(`read`, `write`, `execute` permissions for `user` and for
`group`) which is not configurable.

For workarounds and to learn more, see
[Permissions errors on data directories for shared volumes](../docker-for-windows/troubleshoot.md#permissions-errors-on-data-directories-for-shared-volumes).

### How do symlinks work on Windows?

Docker Desktop supports two types of symlinks: Windows native symlinks and symlinks created inside a container.

The Windows native symlinks are visible within the containers as symlinks, whereas symlinks created inside a container are represented as [mfsymlinks](https://wiki.samba.org/index.php/UNIX_Extensions#Minshall.2BFrench_symlinks): target="_blank" rel="noopener" class="_"}. These are regular Windows files with a special metadata. Therefore the symlinks created inside a container appear as symlinks inside the container, but not on the host.
