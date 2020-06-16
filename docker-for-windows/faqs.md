---
description: Frequently asked questions
keywords: windows faqs
title: Frequently asked questions (FAQ)
---

## Stable and Edge releases

### How do I get the Stable or the Edge version of Docker Desktop?

You can download the Stable or the Edge version of Docker Desktop from [Docker Hub](https://hub.docker.com/editions/community/docker-ce-desktop-windows/).

For installation instructions, see [Install Docker Desktop on Windows](install.md){: target="_blank" class="_"}.

### What is the difference between the Stable and Edge versions of Docker Desktop?

Two different download channels are available in the Community version of Docker Desktop:

The **Stable channel** provides a general availability release-ready installer
for a fully baked and tested, more reliable app. The Stable version of Docker
Desktop includes the latest released version of Docker Engine. The
release schedule is synced with Docker Engine releases and patch releases. On the Stable channel, you can select whether to send usage statistics and other data.

The **Edge channel** provides an installer with new features we are working on, but is not necessarily fully tested. It includes the experimental version of
Docker Engine. Bugs, crashes, and issues can occur when using the Edge version, but you get a chance to preview new functionality, experiment, and provide feedback as Docker Desktop evolves. Edge releases are typically more frequent than for Stable, often one or more per month. Usage statistics and crash reports are sent by default. You do not have the option to disable this on the Edge channel.

### Can I switch between Stable and Edge versions of Docker Desktop?

Yes, you can switch between Stable and Edge versions. You can try out the Edge releases to see what's new, then go back to Stable for other work. However, **you can only have one version of Docker Desktop installed at a time**. For more information, see [Switch between Stable and Edge versions](install.md#switch-between-stable-and-edge-versions).

## What are the system requirements for Docker Desktop?

For information about system requirements, see [Docker Desktop Windows system requirements](install.md#system-requirements).

## What is an experimental feature?

{% include experimental.md %}

## How do I?

### How do I connect to the remote Docker Engine API?

You might need to provide the location of the Engine API for Docker clients and development tools.

On Docker Desktop, clients can connect to the Docker Engine through a
**named pipe**: `npipe:////./pipe/docker_engine`, or **TCP socket** at this URL:
`tcp://localhost:2375`.

This sets `DOCKER_HOST` and `DOCKER_CERT_PATH` environment variables to the
given values (for the named pipe or TCP socket, whichever you use).

See also [Docker Engine API](../engine/api/index.md) and the Docker Desktop for Windows forums topic [How to find the remote API](https://forums.docker.com/t/how-to-find-the-remote-api/20988){: target="_blank" class="_"}.

### How do I connect from a container to a service on the host?

Windows has a changing IP address (or none if you have no network access). We recommend that you connect to the special DNS name `host.docker.internal`, which resolves to the internal IP address used by the host. This is for development purposes and will not work in a production environment outside of Docker Desktop for Windows.

The gateway is also reachable as `gateway.docker.internal`.

For more information about the networking features in Docker Desktop for Windows, see
[Networking](networking.md).

### How do I connect to a container from Windows?

We recommend that you publish a port, or connect from another container. You can use the same method on Linux if the container is on an overlay network and not a bridge network, as these are not routed.

For more information and examples, see
[I want to connect to a container from Windows](networking.md#i-want-to-connect-to-a-container-from-windows) in the [Networking](networking.md) topic.

## Volumes

### Can I change permissions on shared volumes for container-specific deployment requirements?

No, at this point, Docker Desktop does not enable you to control (`chmod`)
the Unix-style permissions on [shared volumes](index.md#file-sharing) for
deployed containers, but rather sets permissions to a default value of
[0777](http://permissions-calculator.org/decode/0777/){: target="_blank" class="_"}
(`read`, `write`, `execute` permissions for `user` and for 
`group`) which is not configurable.

For workarounds and to learn more, see
[Permissions errors on data directories for shared volumes](troubleshoot.md#permissions-errors-on-data-directories-for-shared-volumes).

### How do symlinks work on Windows?

Docker Desktop supports 2 kinds of symlink:

1. Windows native symlinks: these are visible inside containers as symlinks.
2. Symlinks created inside a container: these are represented as [mfsymlinks](https://wiki.samba.org/index.php/UNIX_Extensions#Minshall.2BFrench_symlinks) i.e. regular Windows files with special metadata. These appear as symlinks inside containers but not as symlinks on the host.

## Certificates

### How do I add custom CA certificates?

Docker Desktop supports all trusted Certificate Authorities (CAs) (root or
intermediate). Docker recognizes certs stored under Trust Root
Certification Authorities or Intermediate Certification Authorities.

 For more information on adding server and client side certs, see [Adding TLS certificates](index.md#adding-tls-certificates) in the Getting Started topic.

### How do I add client certificates?

For information on adding client certificates, see [Adding TLS certificates](index.md#adding-tls-certificates) in the Getting Started topic.

### Can I pass through a USB device to a container?

Unfortunately, it is not possible to pass through a USB device (or a
serial port) to a container as it requires support at the hypervisor level.

### Can I run Docker Desktop in nested virtualization scenarios?

Docker Desktop can run inside a Windows 10 VM running on apps like Parallels or VMware Fusion on a Mac provided that the VM is properly configured. However, problems and intermittent failures may still occur due to the way these apps virtualize the hardware. For these reasons, **Docker Desktop is not supported in nested virtualization scenarios**. It might work in some cases, and not in others. For more information, see [Running Docker Desktop in nested virtualization scenarios](/docker-for-windows/troubleshoot/#running-docker-desktop-in-nested-virtualization-scenarios).

### Can I use VirtualBox alongside Docker Desktop?

Yes, you can run VirtualBox along with Docker Desktop if you have enabled the [ Windows Hypervisor Platform](https://docs.microsoft.com/en-us/virtualization/api/){: target="_blank" class="_"} feature on your machine.

## Windows requirements

### Can I run Docker Desktop on Windows Server?

No, running Docker Desktop on Windows Server is not supported.

### How do I run Windows containers on Windows Server?

You can install a native Windows binary which allows you to develop and run
Windows containers without Docker Desktop. For more information, see the tutorial about running Windows containers on Windows Server in
[Getting Started with Windows Containers](https://github.com/docker/labs/blob/master/windows/windows-containers/README.md){: target="_blank" class="_"}.

### Can I install Docker Desktop on Windows 10 Home?

Windows 10 Home, version 2004 users can now install [Docker Desktop Stable 2.3.0.2](https://hub.docker.com/editions/community/docker-ce-desktop-windows/) or a later release with the [WSL 2 backend](wsl.md).

Docker Desktop Stable releases require the Hyper-V feature which is not available in the Windows 10 Home edition.

### Why is Windows 10 required?

Docker Desktop uses the Windows Hyper-V features. While older Windows versions have Hyper-V, their Hyper-V implementations lack features critical for Docker Desktop to work.

### Why does Docker Desktop fail to start when anti-virus software is installed?

Some anti-virus software may be incompatible with Hyper-V and Windows 10 builds which impact Docker
Desktop. For more information, see [Docker Desktop fails to start when anti-virus software is installed](troubleshoot.md#docker-desktop-fails-to-start-when-anti-virus-software-is-installed)
in [Troubleshooting](troubleshoot.md).

## Feedback

### What kind of feedback are we looking for?

Everything is fair game. We'd like your impressions on the download and install
process, startup, functionality available, the GUI, usefulness of the app,
command line integration, and so on. Tell us about problems, what you like, or
functionality you'd like to see added.

### What if I have problems or questions?

You can find information about diagnosing and troubleshooting common issues in the [Logs and Troubleshooting](troubleshoot.md) topic.

If you do not find a solution in Troubleshooting, browse issues on
[Docker Desktop for Windows issues on GitHub](https://github.com/docker/for-win/issues){: target="_blank" class="_"}
or create a new one. You can also create new issues based on diagnostics. To learn more, see
[Diagnose problems, send feedback, and create GitHub issues](troubleshoot.md#diagnose-problems-send-feedback-and-create-github-issues).

The [Docker Desktop for Windows forum](https://forums.docker.com/c/docker-for-windows){: target="_blank" class="_"}
contains discussion threads. You can also create discussion topics there,
but we recommend using the GitHub issues over the forums for better tracking and
response.

### How can I opt out of sending my usage data?

If you do not want to send usage data, use the Stable channel. For more
information, see [What is the difference between the Stable and Edge versions of Docker Desktop](#stable-and-edge-releases).

### How is personal data handled in Docker Desktop?

When uploading diagnostics to help Docker with investigating issues, the
uploaded diagnostics bundle may contain personal data such as usernames and IP
addresses. The diagnostics bundles are only accessible to Docker, Inc. employees
who are directly involved in diagnosing Docker Desktop issues. 

By default, Docker, Inc. will delete uploaded diagnostics bundles after 30 days unless they are referenced in an open issue on the
[docker/for-mac](https://github.com/docker/for-mac/issues) or
[docker/for-win](https://github.com/docker/for-win/issues) issue trackers. If an
issue is closed, Docker, Inc. will remove the referenced diagnostics bundles
within 30 days. You may also request the removal of a diagnostics bundle by
either specifying the diagnostics ID or through your GitHub ID (if the diagnostics ID is mentioned in a GitHub issue). Docker, Inc. will only use the data in the diagnostics bundle to investigate specific user issues, but may derive high-level (non-personal) metrics such as the rate of issues from it.