---
description: Apple Silicon Tech Preview
keywords: Docker Desktop, M1, Silicon, Apple, tech preview, dev preview
title: Apple Silicon Tech Preview
toc_min: 2
toc_max: 3
---

Welcome to Docker Desktop for Apple Silicon.

> **Note**
>
> We encourage you to try the release candidate and report any issues in the [Docker Desktop for Mac GitHub](https://github.com/docker/for-mac) repository.

## Docker Desktop RC 2

2021-03-26

Click on the following link to download the latest release candidate of Docker Desktop for Apple Silicon.

> [Download](https://desktop.docker.com/mac/stable/arm64/62345/Docker.dmg)

In this build, we have defaulted to a `qemu`-based virtual machine, which we believe resolves some of the issues noted as known issues in the previous release candidate. You can switch between `qemu`-based and `virtualization.framework`-based virtual machines using the **Preferences** > **Experimental** tab.

The backend that we choose for the GA release will depend on user feedback and bug reports. We are interested in your feedback on whether the known issues noted below are resolved with the `qemu` backend. We are equally interested in whether the `qemu` backend introduces other issues that we have not discovered in our in-house testing. Please help us by reporting any issues you have with this build, particularly in areas where your experience with the `qemu` backend is worse than the `virtualization.framework` backend.

### Known issues

The following issues are not expected to be resolved in the final GA build for Apple Silicon.

- You must install Rosetta 2 as some binaries are still Darwin/AMD64. To install Rosetta 2 manually from the command line, use this command:

    ```
    softwareupdate --install-rosetta
    ```
    We expect to fix this in a future release.

- Not all images are available for ARM64 architecture. You can add `--platform linux/amd64` to run an Intel image under emulation. In particular, the [mysql](https://hub.docker.com/_/mysql?tab=tags&page=1&ordering=last_updated) image is not available for ARM64. You can work around this issue by using a [mariadb](https://hub.docker.com/_/mariadb?tab=tags&page=1&ordering=last_updated) image.

   However, attempts to run Intel-based containers on Apple Silicon machines can crash as QEMU sometimes fails to run the container. Filesystem change notification APIs (e.g. `inotify`) do not work under QEMU emulation, see [docker/for-mac#5321](https://github.com/docker/for-mac/issues/5321). Therefore, we recommend that you run ARM64 containers on Apple Silicon machines. These containers are also faster and use less memory than Intel-based containers.

   We expect this issue to become less common over time, as more and more images are rebuilt [supporting multiple architectures](https://www.docker.com/blog/multi-arch-build-and-images-the-simple-way/).

The following issues are seen when using the `virtualization.framework` back end.

- Some VPN clients can prevent the VM running Docker from communicating with the host, preventing Docker Desktop starting correctly. See [docker/for-mac#5208](https://github.com/docker/for-mac/issues/5208).

   This is an interaction between `vmnet.framework` (as used by `virtualization.framework` in Big Sur) and the VPN clients.

- Docker Desktop is incompatible with macOS Internet Sharing. See [docker/for-mac#5348](https://github.com/docker/for-mac/issues/5348).

   This is an interaction between `vmnet.framework` (as used by `virtualization.framework` in Big Sur) and macOS Internet Sharing. At the moment it is not possible to use Docker Desktop and macOS Internet Sharing at the same time.

- Some container disk I/O is much slower than expected. See [docker/for-mac#5389](https://github.com/docker/for-mac/issues/5389). Disk flushes are particularly slow due to the need to guarantee data is written to stable storage on the host.

   This is an artifact of the new `virtualization.framework` in Big Sur.

- TCP and UDP port 53 (DNS) are bound on the host when Docker Desktop starts. Therefore you cannot bind to port 53 on all interfaces with a command like `docker run -p 53:53`. See [docker/for-mac#5335](https://github.com/docker/for-mac/issues/5335).

   This is an artifact of the new `virtualization.framework` in Big Sur. A workaround is to bind to a specific IP address e.g. `docker run -p 127.0.0.1:53:53`.

- The Linux Kernel may occasionally crash. Docker now detects this problem and pops up an error dialog offering the user the ability to quickly restart Linux.

   We are still gathering data and testing alternate kernel versions.

### Fixes since Docker Desktop RC 1

- Inter-container HTTP and HTTPS traffic is now routed correctly. Fixes [docker/for-mac#5476](https://github.com/docker/for-mac/issues/5476).

### Fixes since Docker Desktop preview 3.1.0

- The build should update automatically to future versions.
- HTTP proxy support is working, including support for domain name based `no_proxy` rules via TLS SNI. Fixes [docker/for-mac#2732](https://github.com/docker/for-mac/issues/2732).

### Fixes since the Apple Silicon preview 7

- Kubernetes now works (although you might need to reset the cluster in our Troubleshoot menu one time to regenerate the certificates).
- osxfs file sharing works.
- The `host.docker.internal` and `vm.docker.internal` DNS entries now resolve.
- Removed hard-coded IP addresses: Docker Desktop now dynamically discovers the IP allocated by macOS.
- The updated version includes a  change that should improve disk performance.
- The **Restart** option in the Docker menu works.

## Feedback

Thank you for trying out the Docker Desktop for Apple Silicon tech preview. Your feedback is important to us. Let us know your feedback by creating an issue in the [Docker Desktop for Mac GitHub](https://github.com/docker/for-mac/issues)repository.

We also recommend that you join the [Docker Community Slack](https://www.docker.com/docker-community) and ask questions in **#docker-desktop-mac** channel.

For more information about the tech preview, see our blog post [Download and Try the Tech Preview of Docker Desktop for M1](https://www.docker.com/blog/download-and-try-the-tech-preview-of-docker-desktop-for-m1/).
