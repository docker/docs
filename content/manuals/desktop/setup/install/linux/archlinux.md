---
description: Instructions for installing Docker Desktop Arch package. Mostly meant
  for hackers who want to try out Docker Desktop on a variety of Arch-based distributions.
keywords: Arch Linux, install, uninstall, upgrade, update, linux, desktop, docker
  desktop, docker desktop for linux, dd4l
title: Install Docker Desktop on Arch-based distributions
linkTitle: Arch
aliases:
- /desktop/linux/install/archlinux/
- /desktop/install/archlinux/
- /desktop/install/linux/archlinux/
---

{{< summary-bar feature_name="Docker Desktop Archlinux" >}}

> **Docker Desktop terms**
>
> Commercial use of Docker Desktop in larger enterprises (more than 250
> employees OR more than $10 million USD in annual revenue) requires a [paid
> subscription](https://www.docker.com/pricing/).

This page contains information on how to install, launch and upgrade Docker Desktop on an Arch-based distribution. 



## Prerequisites

To install Docker Desktop successfully, you must meet the [general system requirements](_index.md#general-system-requirements).

## Install Docker Desktop

1. [Install the Docker client binary on Linux](/manuals/engine/install/binaries.md#install-daemon-and-client-binaries-on-linux). Static binaries for the Docker client are available for Linux as `docker`. You can use:

   ```console
   $ wget https://download.docker.com/linux/static/stable/x86_64/docker-{{% param "docker_ce_version" %}}.tgz -qO- | tar xvfz - docker/docker --strip-components=1
   $ cp -rp ./docker /usr/local/bin/ && rm -r ./docker
   ```

2. Download the latest Arch package from the [Release notes](/manuals/desktop/release-notes.md).

3. Install the package:

   ```console
   $ sudo pacman -U ./docker-desktop-x86_64.pkg.tar.zst
   ```

   By default, Docker Desktop is installed at `/opt/docker-desktop`.

## Launch Docker Desktop

{{% include "desktop-linux-launch.md" %}}

## Next steps

- Explore [Docker's subscriptions](https://www.docker.com/pricing/) to see what Docker can offer you.
- Take a look at the [Docker workshop](/get-started/workshop/_index.md) to learn how to build an image and run it as a containerized application.
- [Explore Docker Desktop](/manuals/desktop/use-desktop/_index.md) and all its features.
- [Troubleshooting](/manuals/desktop/troubleshoot-and-support/troubleshoot/_index.md) describes common problems, workarounds, how to run and submit diagnostics, and submit issues.
- [FAQs](/manuals/desktop/troubleshoot-and-support/faqs/general.md) provide answers to frequently asked questions.
- [Release notes](/manuals/desktop/release-notes.md) lists component updates, new features, and improvements associated with Docker Desktop releases.
- [Back up and restore data](/manuals/desktop/settings-and-maintenance/backup-and-restore.md) provides instructions
  on backing up and restoring data related to Docker.
