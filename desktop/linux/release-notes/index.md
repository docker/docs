---
description: Change log / release notes for Docker Desktop for Linux
keywords: Docker Desktop for Linux, release notes
title: Docker Desktop for Linux release notes
toc_min: 1
toc_max: 2
---

> **Update to the Docker Desktop terms**
>
> Commercial use of Docker Desktop in larger enterprises (more than 250
> employees OR more than $10 million USD in annual revenue) now requires a paid
> subscription. The grace period for those that will require a paid subscription
> ends on January 31, 2022. [Learn more](https://www.docker.com/blog/the-grace-period-for-the-docker-subscription-service-agreement-ends-soon-heres-what-you-need-to-know/){:
 target="_blank" rel="noopener" class="_" id="dkr_docs_cta"}.
{: .important}

This page contains information about the new features, improvements, known issues, and bug fixes in Docker Desktop releases.

Take a look at the [Docker Public Roadmap](https://github.com/docker/roadmap/projects/1){: target="_blank" rel="noopener" class="_"} to see what's coming next.

## Docker Desktop 4.8.0
2022-05-05   

> Download Docker Desktop
>
> [DEB](https://desktop-stage.docker.com/linux/main/amd64/78459/docker-desktop-4.8.0-amd64.deb) |
> [RPM](https://desktop-stage.docker.com/linux/main/amd64/78459/docker-desktop-4.8.0-x86_64.rpm)
> [Arch](https://desktop-stage.docker.com/linux/main/amd64/78459/docker-desktop-4.8.0-x86_64.pkg.tar.zst)

### Known issues

- Fedora missing tray icon on a default setup - issue with all electron applications.

- only works with systemd for now

- On Arch Linux, the electron binary crashes when launched with systemctl. Starting it directly works fine.
