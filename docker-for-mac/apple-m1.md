---
description: Apple M1 Tech Preview
keywords: Docker Desktop, M1, Silicon, Apple, tech preview, dev preview
title: Apple M1 Tech Preview
toc_min: 2
toc_max: 3
---

Welcome to the tech preview of Docker Desktop for Apple M1. This tech preview is aimed at early adopters of Apple M1 machines, who would like to try an experimental build of Docker Desktop.

> **Note**
>
> Docker Desktop on Apple M1 chip is still under development. We recommend that you do not use tech preview builds in production environments.

## Download

Click the following link to download the Apple M1 tech preview build:

> [Download](https://desktop.docker.com/mac/stable/arm64/60984/Docker.dmg)

## Known issues

The tech preview of Docker Desktop for Apple M1 currently has the following limitations:

- The tech preview build does not update automatically. You must manually install any future versions of Docker Desktop.
- You must install Rosetta 2 as some binaries are still Darwin/AMD64. To install Rosetta 2 manually from the command line use this command:
    ```
    softwareupdate --install-rosetta
    ```
- The HTTP proxy is not enabled.
- Not all images are available for ARM64. You can add `--platform linux/amd64` to run an Intel image under emulation.

    In particular, the [mysql](https://hub.docker.com/_/mysql?tab=tags&page=1&ordering=last_updated){: target="blank" rel="noopener" class=“”} image is not available for ARM64. You can work around this issue by using a [mariadb](https://hub.docker.com/_/mariadb?tab=tags&page=1&ordering=last_updated){: target="blank" rel="noopener" class=“”} image.
- The kernel may panic. If so, look in `~/Library/Containers/com.docker.docker/Data/vms/0/console.log` for a BUG or kernel panic to report.

## Fixes since the Apple Silicon preview 7

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
