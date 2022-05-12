---
description: Change log / release notes for Docker Desktop for Linux
keywords: Docker Desktop for Linux, release notes
title: Docker Desktop for Linux release notes
toc_min: 1
toc_max: 2
---

This page contains information about the new features, improvements, known issues, and bug fixes in Docker Desktop releases.

Take a look at the [Docker Public Roadmap](https://github.com/docker/roadmap/projects/1){: target="_blank" rel="noopener" class="_"} to see what's coming next.

## Docker Desktop 4.8.2
2022-05-17

> Download Docker Desktop
>
> [DEB](https://desktop.docker.com/linux/main/amd64/docker-desktop-4.8.2-amd64.deb?utm_source=docker&utm_medium=webreferral&utm_campaign=docs-driven-download-linux-amd64) |
> [RPM](https://desktop.docker.com/linux/main/amd64/docker-desktop-4.8.2-x86_64.rpm?utm_source=docker&utm_medium=webreferral&utm_campaign=docs-driven-download-linux-amd64) |
> [Arch package](https://desktop.docker.com/linux/main/amd64/docker-desktop-4.8.2-x86_64.pkg.tar.zst?utm_source=docker&utm_medium=webreferral&utm_campaign=docs-driven-download-linux-amd64)

### Bug fixes and minor changes

- Fixed an issue with manual proxy settings which caused problems when pulling images. Fixes [docker/for-win#12714](https://github.com/docker/for-win/issues/12714) and [docker/for-mac#6315](https://github.com/docker/for-mac/issues/6315).
- Fixed high CPU usage when extensions are disabled. Fixes [docker/for-mac#6310](https://github.com/docker/for-mac/issues/6310).
- Docker Desktop now redacts HTTP proxy passwords in log files and diagnostics.

## Docker Desktop 4.8.1
2022-05-09

> Download Docker Desktop
>
> [DEB](https://desktop-stage.docker.com/linux/main/amd64/78998/docker-desktop-4.8.1-amd64.deb) |
> [RPM](https://desktop-stage.docker.com/linux/main/amd64/78998/docker-desktop-4.8.1-x86_64.rpm) |
> [Arch package](https://desktop-stage.docker.com/linux/main/amd64/78998/docker-desktop-4.8.1-x86_64.pkg.tar.zst)

## Bugfixes and minor changes

- Fixed a bug that caused the Kubernetes cluster to be deleted when updating Docker Desktop.

## Known issues

- Changing ownership rights for files in bind mounts fails. This is due to the way we have implemented file sharing between the host and VM within which the Docker Engine runs. We aim to resolve this issue in the next release.

## Docker Desktop 4.8.0
2022-05-06

> Download Docker Desktop
>
> [DEB](https://desktop.docker.com/linux/main/amd64/docker-desktop-4.8.0-amd64.deb?utm_source=docker&utm_medium=webreferral&utm_campaign=docs-driven-download-linux-amd64) |
> [RPM](https://desktop.docker.com/linux/main/amd64/docker-desktop-4.8.0-x86_64.rpm?utm_source=docker&utm_medium=webreferral&utm_campaign=docs-driven-download-linux-amd64) |
> [Arch package](https://desktop.docker.com/linux/main/amd64/docker-desktop-4.8.0-x86_64.pkg.tar.zst?utm_source=docker&utm_medium=webreferral&utm_campaign=docs-driven-download-linux-amd64)

## Known issues

- Changing ownership rights for files in bind mounts fails. This is due to the way we have implemented file sharing between the host and VM within which the Docker Engine runs. We aim to resolve this issue in the next release.
- Currently, if you are running a Kubernetes cluster, it will be deleted when you upgrade to Docker Desktop 4.8.0. We aim to fix this in the next release.
