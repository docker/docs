---
description: Apple M1
keywords: Docker Desktop, M1, Silicon, Apple, tech preview, dev preview
title: Apple M1
toc_min: 2
toc_max: 3
---

Welcome to Docker Desktop for Apple Silicon.

> **Note**
>
> Docker Desktop for Apple Silicon is still under development. Please try these builds and report issues on [docker/for-mac](https://github.com/docker/for-mac).

## Docker Desktop RC 1

2021-03-18

Click the following link to download the Apple Silicon RC 1 build:

> [Download](https://desktop.docker.com/mac/stable/arm64/62015/Docker.dmg)

### Known issues

The following issues are known and are not expected to be resolved in the final GA build for M1.

- You must install Rosetta 2 as some binaries are still Darwin/AMD64. To install Rosetta 2 manually from the command line, use this command:

    ```
    softwareupdate --install-rosetta
    ```
    We expect to fix this in a future release.

- Not all images are available for ARM64 architecture. You can add `--platform linux/amd64` to run an Intel image under emulation. In particular, the [mysql](https://hub.docker.com/_/mysql?tab=tags&page=1&ordering=last_updated){: target="blank" rel="noopener" class=“”} image is not available for ARM64. You can work around this issue by using a [mariadb](https://hub.docker.com/_/mariadb?tab=tags&page=1&ordering=last_updated){: target="blank" rel="noopener" class=“”} image.

   However, attempts to run Intel-based containers on Apple M1 machines can crash as QEMU sometimes fails to run the container. Therefore, we recommend that you run ARM64 containers on M1 machines. These containers are also faster and use less memory than Intel-based containers.

   We expect this issue to become less common over time, as more and more images are rebuilt [supporting multiple architectures](https://www.docker.com/blog/multi-arch-build-and-images-the-simple-way/){: target="blank" rel="noopener" class=“”}.

The following issues will also not be fixed for the final release, but we are working on them in collaboration with our partners. We expect to be able to make improvements in future releases.

- Some VPN clients can prevent the VM running Docker from communicating with the host, preventing Docker Desktop starting correctly. See [docker/for-mac#5208](https://github.com/docker/for-mac/issues/5208){: target="blank" rel="noopener" class=“”}.

- Docker Desktop is incompatible with macOS Internet Sharing. See [docker/for-mac#5348](https://github.com/docker/for-mac/issues/5348){: target="blank" rel="noopener" class=“”}.

- Some container disk I/O is much slower than expected. See [docker/for-mac#5389](https://github.com/docker/for-mac/issues/5389).

   We believe that disk flushes are particularly slow due to the need to guarantee data is written to stable storage on the host.

- The new Apple Virtualization framework uses port 53 (DNS) when Docker Desktop starts. Therefore you cannot bind to port 53 on all interfaces with a command like `docker run -p 53:53`. See [docker/for-mac#5335](https://github.com/docker/for-mac/issues/5335).

   This is an artifact of the new `virtualization.framework` in Big Sur. A workaround is to bind to a specific IP address e.g. `docker run -p 127.0.0.1:53:53`.

- The Linux Kernel may occasionally crash. Docker now detects this problem and pops up an error dialog offering the user the ability to quickly restart Linux.

   We are still gathering data and testing alternate kernel versions and hope to make improvements in a future release.

### Fixes since Docker Desktop preview 3.1.0

- The build should update automatically to future versions.
- HTTP proxy support is working, including support for domain name based `no_proxy` rules via TLS SNI. Fixes [docker/for-mac#2732](https://github.com/docker/for-mac/issues/2732).

## Docker Desktop preview 3.1.0

2021-02-11

Click the following link to download the Apple M1 tech preview build:

> [Download](https://desktop.docker.com/mac/stable/arm64/60984/Docker.dmg)

## Known issues

The tech preview of Docker Desktop for Apple M1 currently has the following limitations:

- The tech preview build does not update automatically. You must manually install any future versions of Docker Desktop.
- You must install Rosetta 2 as some binaries are still Darwin/AMD64. To install Rosetta 2 manually from the command line, use this command:

    ```
    softwareupdate --install-rosetta
    ```

- The HTTP proxy is not enabled.

- Not all images are available for ARM64 architecture. You can add `--platform linux/amd64` to run an Intel image under emulation. In particular, the [mysql](https://hub.docker.com/_/mysql?tab=tags&page=1&ordering=last_updated){: target="blank" rel="noopener" class=“”} image is not available for ARM64. You can work around this issue by using a [mariadb](https://hub.docker.com/_/mariadb?tab=tags&page=1&ordering=last_updated){: target="blank" rel="noopener" class=“”} image.

   However, attempts to run Intel-based containers on Apple M1 machines can crash as QEMU sometimes fails to run the container. Therefore, we recommend that you run ARM64 containers on M1 machines. These containers are also faster and use less memory than Intel-based containers.

- Some VPN clients can prevent the VM running Docker from communicating with the host, preventing Docker Desktop starting correctly. See [docker/for-mac#5208](https://github.com/docker/for-mac/issues/5208){: target="blank" rel="noopener" class=“”}.

- Docker Desktop is incompatible with macOS Internet Sharing. See [docker/for-mac#5348](https://github.com/docker/for-mac/issues/5348){: target="blank" rel="noopener" class=“”}.

- The kernel may panic. If so, look in `~/Library/Containers/com.docker.docker/Data/vms/0/console.log` for a BUG or kernel panic to report.

- The new Apple Virtualization framework uses port 53 (DNS) when Docker Desktop starts. You cannot use this port to bind a container's port to the host. See [docker/for-mac#5335](https://github.com/docker/for-mac/issues/5335).


## Fixes since the Apple Silicon preview 7

**Docker Desktop preview 3.1.0 (60984)**

2021-02-11

- Kubernetes now works (although you might need to reset the cluster in our Troubleshoot menu one time to regenerate the certificates).
- osxfs file sharing works.
- The `host.docker.internal` and `vm.docker.internal` DNS entries now resolve.
- Removed hard-coded IP addresses: Docker Desktop now dynamically discovers the IP allocated by macOS.
- The updated version includes a  change that should improve disk performance.
- The **Restart** option in the Docker menu works.

## Feedback

Thank you for trying out the Docker Desktop for Apple M1 tech preview. Your feedback is important to us. Let us know your feedback by creating an issue in the [Docker Desktop for Mac GitHub](https://github.com/docker/for-mac/issues){: target="blank" rel="noopener" class=“”} repository.

We also recommend that you join the [Docker Community Slack](https://www.docker.com/docker-community){: target="blank" rel="noopener" class=“”} and ask questions in **#docker-desktop-mac** channel.

For more information about the tech preview, see our blog post [Download and Try the Tech Preview of Docker Desktop for M1](https://www.docker.com/blog/download-and-try-the-tech-preview-of-docker-desktop-for-m1/){: target="blank" rel="noopener" class=“”}.
