---
description: Change log / release notes for Docker Desktop for Linux
keywords: Docker Desktop for Linux, release notes
title: Docker Desktop for Linux release notes
toc_min: 1
toc_max: 2
---

This page contains information about the new features, improvements, known issues, and bug fixes in Docker Desktop releases.

Take a look at the [Docker Public Roadmap](https://github.com/docker/roadmap/projects/1){: target="_blank" rel="noopener" class="_"} to see what's coming next.

## Docker Desktop 4.8.0
2022-05-06

> Download Docker Desktop
>
> [DEB](https://desktop-stage.docker.com/linux/main/amd64/78933/docker-desktop-4.8.0-amd64.deb) |
> [RPM](https://desktop-stage.docker.com/linux/main/amd64/78933/docker-desktop-4.8.0-x86_64.rpm) |
> [Arch package](https://desktop-stage.docker.com/linux/main/amd64/78933/docker-desktop-4.8.0-x86_64.pkg.tar.zst)

## Known issues

* Currently altering ownership rights for files in bind mounts will fail. This is a limitation of how we have implemented file sharing between the host and VM within which the Docker Engine runs. We expect to resolve this issue in the next release.
