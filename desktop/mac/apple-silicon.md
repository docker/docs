---
description: Docker Desktop for Apple silicon
keywords: Docker Desktop, M1, Silicon, Apple,
title: Docker Desktop for Apple silicon
redirect_from:
- /docker-for-mac/apple-m1/
- /docker-for-mac/apple-silicon/
---

Docker Desktop for Mac on Apple silicon is now available as a GA release. This enables you to develop applications with your choice of local development environments, and extends development pipelines for ARM-based applications.

Docker Desktop for Apple silicon also supports multi-platform images, which allows you to build and run images for both x86 and ARM architectures without having to set up a complex cross-compilation development environment. Additionally, you can use [docker buildx](../../engine/reference/commandline/buildx.md){:target="_blank" rel="noopener" class="_"} to seamlessly integrate multi-platform builds into your build pipeline, and use [Docker Hub](https://hub.docker.com/){:target="_blank" rel="noopener" class="_"} to identify and share repositories that provide multi-platform images.

Download Docker Desktop for Mac on Apple silicon:

> Download Docker Desktop
>
>
> [Mac with Apple chip](https://desktop.docker.com/mac/main/arm64/Docker.dmg?utm_source=docker&utm_medium=webreferral&utm_campaign=docs-driven-download-mac-arm64){: .button .primary-btn }

<br>

### System requirements

Beginning with Docker Desktop 4.3.0, we have removed the hard requirement to install **Rosetta 2**. There are a few optional command line tools that still require Rosetta 2 when using Darwin/AMD64. See the Known issues section below. However, to get the best experience, we recommend that you install Rosetta 2. To install Rosetta 2 manually from the command line, run the following command:

```console
$ softwareupdate --install-rosetta
```

### Known issues

- Some command line tools do not work when Rosetta 2 is not installed.
  - The old version 1.x of `docker-compose`. We recommend that you use Compose V2 instead. Either type `docker compose` or enable the **Use Docker Compose V2** option in the [General preferences tab](../settings/mac-settings.md#general).
  - The `docker scan` command and the underlying `snyk` binary.
  - The `docker-credential-ecr-login` credential helper.
- Not all images are available for ARM64 architecture. You can add `--platform linux/amd64` to run an Intel image under emulation. In particular, the [mysql](https://hub.docker.com/_/mysql?tab=tags&page=1&ordering=last_updated) image is not available for ARM64. You can work around this issue by using a [mariadb](https://hub.docker.com/_/mariadb?tab=tags&page=1&ordering=last_updated) image.

   However, attempts to run Intel-based containers on Apple silicon machines under emulation can crash as qemu sometimes fails to run the container. In addition, filesystem change notification APIs (`inotify`) do not work under qemu emulation. Even when the containers do run correctly under emulation, they will be slower and use more memory than the native equivalent.

   In summary, running Intel-based containers on Arm-based machines should be regarded as "best effort" only. We recommend running arm64 containers on Apple silicon machines whenever possible, and encouraging container authors to produce arm64, or multi-arch, versions of their containers. We expect this issue to become less common over time, as more and more images are rebuilt [supporting multiple architectures](https://www.docker.com/blog/multi-arch-build-and-images-the-simple-way/).
- `ping` from inside a container to the Internet does not work as expected.  To test the network, we recommend using `curl` or `wget`. See [docker/for-mac#5322](https://github.com/docker/for-mac/issues/5322#issuecomment-809392861).
- Users may occasionally experience data drop when a TCP stream is half-closed.

### Fixes since Docker Desktop RC 3

- Docker Desktop now ensures the permissions of `/dev/null` and other devices are correctly set to `0666` (`rw-rw-rw-`) inside `--privileged` containers. Fixes [docker/for-mac#5527](https://github.com/docker/for-mac/issues/5527).
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

Your feedback is important to us. Let us know your feedback by creating an issue in the [Docker Desktop for Mac GitHub](https://github.com/docker/for-mac/issues) repository.

We also recommend that you join the [Docker Community Slack](https://www.docker.com/docker-community) and ask questions in **#docker-desktop-mac** channel.
