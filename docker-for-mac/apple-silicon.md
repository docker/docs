---
description: Docker Desktop for Apple silicon
keywords: Docker Desktop, M1, Silicon, Apple,
title: Docker Desktop for Apple silicon
toc_min: 2
toc_max: 3
redirect_from:
- /docker-for-mac/apple-m1/
---

Welcome to Docker Desktop for Apple silicon.

## Docker Desktop 3.3.1

2021-04-15

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

- When using the `qemu` backend, `ping` from inside a container to the Internet does not work as expected. To test the network, we recommend using `curl` or `wget`. See [docker/for-mac#5322](https://github.com/docker/for-mac/issues/5322#issuecomment-809392861).

### Fixes since Docker Desktop RC 3

- Docker Desktop now ensures the permissions of `/dev/null` are correctly set to `0666` (`rw-rw-rw-`) inside `--privileged` containers. Fixes [docker/for-mac#5527](https://github.com/docker/for-mac/issues/5527).
- Fixed an issue where ICMP echo responses (`ping` responses) had incorrect sequence numbers.
- Fixed an issue where data was dropped when a TCP stream was half-closed.
- Docker Desktop now reduces the idle CPU consumption.

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

Your feedback is important to us. Let us know your feedback by creating an issue in the [Docker Desktop for Mac GitHub](https://github.com/docker/for-mac/issues)repository.

We also recommend that you join the [Docker Community Slack](https://www.docker.com/docker-community) and ask questions in **#docker-desktop-mac** channel.
