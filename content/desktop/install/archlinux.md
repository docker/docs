---
description: Instructions for installing Docker Desktop Arch package. Mostly meant
  for hackers who want to try out Docker Desktop on a variety of Arch-based distributions.
keywords: Arch Linux, install, uninstall, upgrade, update, linux, desktop, docker
  desktop, docker desktop for linux, dd4l
title: Install Docker Desktop on Arch-based distributions
aliases:
- /desktop/linux/install/archlinux/
---

> **Docker Desktop terms**
>
> Commercial use of Docker Desktop in larger enterprises (more than 250
> employees OR more than $10 million USD in annual revenue) requires a [paid
> subscription](https://www.docker.com/pricing/).

This page contains information on how to install, launch and upgrade Docker Desktop on an Arch-based distribution. Docker has not tested or verified the installation.

## Prerequisites

To install Docker Desktop successfully, you must meet the [general system requirements](linux-install.md#general-system-requirements).

Additionally, for non-Gnome Desktop environments, `gnome-terminal` must be installed:

```console
$ sudo pacman -S gnome-terminal
```

## Install Docker Desktop

1. Install client binaries. Docker does not have an Arch package repository. Binaries not included in the package must be installed manually before installing Docker Desktop.

2. [Install Docker client binary on Linux](../../engine/install/binaries.md#install-daemon-and-client-binaries-on-linux). On Arch-based distributions, users must install the Docker client binary.
   Static binaries for the Docker client are available for Linux (as `docker`).

3. Download the latest Arch package from the [Release notes](../release-notes.md).

4. Install the package:

   ```console
   $ sudo pacman -U ./docker-desktop-<arch>.pkg.tar.zst
   ```

   Don't forget to substitute `<arch>` with the architecture you want.

   By default, Docker Desktop is installed at `/opt/docker-desktop`.

## Launch Docker Desktop

{{< include "desktop-linux-launch.md" >}}

## Next steps

- Explore [Docker's core subscriptions](https://www.docker.com/pricing/) to see what Docker can offer you.
- Take a look at the [Docker workshop](../../guides/workshop/_index.md) to learn how to build an image and run it as a containerized application.
- [Explore Docker Desktop](../use-desktop/index.md) and all its features.
- [Troubleshooting](../troubleshoot/overview.md) describes common problems, workarounds, how to run and submit diagnostics, and submit issues.
- [FAQs](../faqs/general.md) provide answers to frequently asked questions.
- [Release notes](../release-notes.md) lists component updates, new features, and improvements associated with Docker Desktop releases.
- [Back up and restore data](../backup-and-restore.md) provides instructions
  on backing up and restoring data related to Docker.
