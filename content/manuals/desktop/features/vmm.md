---
title: Virtual Machine Manager
linkTitle: Virtual Machine Manager
keywords: virtualization software, resource allocation, mac, windows, docker desktop, vm performance, apple silicon, wsl, hyper-v
description: Learn about Docker Desktop's Virtual Machine Manager options, including Docker VMM for Mac and Windows
weight: 110
params:
  sidebar:
    badge:
      color: blue
      text: Beta
---

Docker Desktop supports multiple Virtual Machine Managers (VMMs) to power the Linux VM that runs containers. The options available depend on your platform.

## Docker VMM

{{< summary-bar feature_name="VMM" >}}

Docker VMM is a container-optimized hypervisor. From Docker Desktop 4.86, Docker VMM uses Docker's own hypervisor, replacing `libkrun` used in version 4.79 and earlier for Mac users. Built specifically for container workloads, Docker VMM:

- Returns idle memory to the host when containers aren't active, so Docker Desktop doesn't hold RAM it's not using
- Improves file I/O between container and host, reducing latency in the edit-compile-test loop
- Reduces engine and container start-up time

Because Docker controls the virtualization layer, it can be monitored and governed in ways that aren't possible with third-party backends. On Windows, Docker VMM provides a stable alternative to WSL 2 with a real VM boundary between the container environment and the host.

### Switch to Docker VMM

Docker VMM requires a minimum of 4 GB of memory allocated to the Docker Linux VM. Increase memory in **Settings** > **Resources** before switching.

{{< tabs >}}
{{< tab name="Mac (Apple Silicon)" >}}

1. Go to **Settings** > **General** > **Virtual Machine Manager**.
2. Select **Docker VMM**.
3. Select **Apply & restart**.

If you previously had Docker VMM selected, which engine runs depends on your version:

- Docker Desktop 4.79 and earlier is backed by `libkrun`
- Docker Desktop 4.86 and later is backed by Docker's own hypervisor

If you're upgrading from 4.79 or earlier, your setting is preserved and Docker Desktop switches to the new Docker VMM automatically on restart.

{{< /tab >}}
{{< tab name="Windows" >}}

1. Go to **Settings** > **General** > **Virtual Machine Manager**.
2. Select **Docker VMM**.
3. Select **Apply & restart**.

{{< /tab >}}
{{< /tabs >}}

### Known issues

- A restart of Docker Desktop may be required after switching to Docker VMM.
- Docker VMM does not support bind mount auto-shares. If you see a `file is not shared from the host` error, go to **Settings** > **Resources** > **File sharing** and add the directory you want to share.

#### Mac only

- Docker VMM does not currently support Rosetta, so emulation of amd64 architectures is slow. Docker is exploring potential solutions.
- Certain databases, such as MongoDB and Cassandra, may fail when using virtiofs with Docker VMM. This issue is expected to be resolved in a future release.

## Alternative VMMs for Mac

### Apple Virtualization framework

The Apple Virtualization framework is a stable and well-established option for managing virtual machines on Mac. It has been a reliable choice for many Mac users over the years.

### HyperKit (Legacy) for Intel-based Macs

> [!NOTE]
>
> HyperKit is deprecated. Docker recommends switching to the Apple Virtualization framework.

HyperKit is a legacy virtualization option for Intel-based Macs. Docker recommends switching to modern alternatives for better performance and to future-proof your setup.

## Alternative VMMs for Windows

### WSL 2

WSL 2 (Windows Subsystem for Linux 2) is the default Windows backend for Docker Desktop. It runs a full Linux kernel inside a lightweight VM with tight integration into the Windows host file system and networking. WSL 2 is available in both per-user and all-users installation modes and does not require administrator privileges.

For more information, see [Docker Desktop WSL 2 backend](/manuals/desktop/features/wsl/_index.md).

### Hyper-V

Hyper-V is Windows' native hypervisor. It runs the Docker Linux VM in a fully isolated virtual machine, providing a strong boundary between the container environment and the Windows host. Hyper-V is only available in all-users installation mode and requires administrator privileges.
