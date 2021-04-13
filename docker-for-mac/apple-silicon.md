---
description: Docker Desktop for Apple silicon
keywords: Docker Desktop, M1, Silicon, Apple,
title: Docker Desktop for Apple silicon
toc_min: 2
toc_max: 3
redirect_from:
- /docker-for-mac/apple-m1/
---

Welcome to Docker Desktop for Apple Silicon.

## Docker Desktop for Apple silicon

2021-04-14

Click on the following link to download the GA version of Docker Desktop for Apple silicon.

> [Download](https://desktop.docker.com/mac/stable/arm64/Docker.dmg)

### System requirements

You must install **Rosetta 2** as some binaries are still Darwin/AMD64. To install Rosetta 2 manually from the command line, run the following command:

    ```
    softwareupdate --install-rosetta
    ```
    We expect to fix this in a future release.

### Known issues

- Not all images are available for ARM64 architecture. You can add `--platform linux/amd64` to run an Intel image under emulation. In particular, the [mysql](https://hub.docker.com/_/mysql?tab=tags&page=1&ordering=last_updated) image is not available for ARM64. You can work around this issue by using a [mariadb](https://hub.docker.com/_/mariadb?tab=tags&page=1&ordering=last_updated) image.

   However, attempts to run Intel-based containers on Apple Silicon machines can crash as QEMU sometimes fails to run the container. Filesystem change notification APIs (e.g. `inotify`) do not work under QEMU emulation, see [docker/for-mac#5321](https://github.com/docker/for-mac/issues/5321). Therefore, we recommend that you run ARM64 containers on Apple Silicon machines. These containers are also faster and use less memory than Intel-based containers.

   We expect this issue to become less common over time, as more and more images are rebuilt [supporting multiple architectures](https://www.docker.com/blog/multi-arch-build-and-images-the-simple-way/).

The following issues are seen when using the `virtualization.framework` experimental feature:

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

The following issues are seen when using the `qemu` backend.

- `ping` from inside a container to the Internet does not work as expected. To test the network, we recommend using `curl` or `wget`. See [docker/for-mac#5322](https://github.com/docker/for-mac/issues/5322#issuecomment-809392861).

### Fixes since Docker Desktop RC 3

ADD list of issues fixed

### Fixes since Docker Desktop RC 2

- Update to [Linux kernel 5.10.25](https://hub.docker.com/layers/docker/for-desktop-kernel/5.10.25-6594e668feec68f102a58011bb42bd5dc07a7a9b/images/sha256-80e22cd9c9e6a188a785d0e23b4cefae76595abe1e4a535449627c2794b10871?context=repo) to improve reliability.

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

Thank you for trying out the Docker Desktop for Apple silicon. Your feedback is important to us. Let us know your feedback by creating an issue in the [Docker Desktop for Mac GitHub](https://github.com/docker/for-mac/issues)repository.

We also recommend that you join the [Docker Community Slack](https://www.docker.com/docker-community) and ask questions in **#docker-desktop-mac** channel.
