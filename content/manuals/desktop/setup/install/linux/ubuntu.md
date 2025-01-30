---
description: Learn how to install, launch and upgrade Docker Desktop on Ubuntu. This
  quick guide will cover prerequisites, installation methods, and more.
keywords: install docker ubuntu, ubuntu install docker, install docker on ubuntu,
  docker install ubuntu, how to install docker on ubuntu, ubuntu docker install, docker
  installation on ubuntu, docker ubuntu install, docker installing ubuntu, installing
  docker on ubuntu, docker desktop for ubuntu
title: Install Docker Desktop on Ubuntu
linkTitle: Ubuntu
weight: 10
toc_max: 4
aliases:
- /desktop/linux/install/ubuntu/
- /desktop/install/ubuntu/
- /desktop/install/linux/ubuntu/
---

> **Docker Desktop terms**
>
> Commercial use of Docker Desktop in larger enterprises (more than 250
> employees OR more than $10 million USD in annual revenue) requires a [paid
> subscription](https://www.docker.com/pricing/).

This page contains information on how to install, launch and upgrade Docker Desktop on an Ubuntu distribution.

## Prerequisites

To install Docker Desktop successfully, you must:

- Meet the [general system requirements](_index.md#general-system-requirements).
- Have an x86-64 system with Ubuntu 22.04, 24.04, or the latest non-LTS version.
- For non-Gnome Desktop environments, `gnome-terminal` must be installed:
  ```console
  $ sudo apt install gnome-terminal
  ```

## Install Docker Desktop

Recommended approach to install Docker Desktop on Ubuntu:

1. Set up Docker's package repository.
   See step one of [Install using the `apt` repository](/manuals/engine/install/ubuntu.md#install-using-the-repository).

2. Download the latest [DEB package](https://desktop.docker.com/linux/main/amd64/docker-desktop-amd64.deb?utm_source=docker&utm_medium=webreferral&utm_campaign=docs-driven-download-linux-amd64). For checksums, see the [Release notes](/manuals/desktop/release-notes.md).

3. Install the package with apt as follows:

   ```console
   $ sudo apt-get update
   $ sudo apt-get install ./docker-desktop-amd64.deb
   ```

   > [!NOTE]
   >
   > At the end of the installation process, `apt` displays an error due to installing a downloaded package. You
   > can ignore this error message.
   >
   > ```text
   > N: Download is performed unsandboxed as root, as file '/home/user/Downloads/docker-desktop.deb' couldn't be accessed by user '_apt'. - pkgAcquire::Run (13: Permission denied)
   > ```

   By default, Docker Desktop is installed at `/opt/docker-desktop`.

There are a few post-install configuration steps done through the post-install script contained in the deb package.

The post-install script:

- Sets the capability on the Docker Desktop binary to map privileged ports and set resource limits.
- Adds a DNS name for Kubernetes to `/etc/hosts`.
- Creates a symlink from `/usr/local/bin/com.docker.cli` to `/usr/bin/docker`.
  This is because the classic Docker CLI is installed at `/usr/bin/docker`. The Docker Desktop installer also installs a Docker CLI binary that includes cloud-integration capabilities and is essentially a wrapper for the Compose CLI, at`/usr/local/bin/com.docker.cli`. The symlink ensures that the wrapper can access the classic Docker CLI. 

## Launch Docker Desktop

{{% include "desktop-linux-launch.md" %}}

## Upgrade Docker Desktop

Once a new version for Docker Desktop is released, the Docker UI shows a notification.
You need to download the new package each time you want to upgrade Docker Desktop and run:

```console
$ sudo apt-get install ./docker-desktop-amd64.deb
```

## Next steps

- Explore [Docker's subscriptions](https://www.docker.com/pricing/) to see what Docker can offer you.
- Take a look at the [Docker workshop](/get-started/workshop/_index.md) to learn how to build an image and run it as a containerized application.
- [Explore Docker Desktop](/manuals/desktop/use-desktop/_index.md) and all its features.
- [Troubleshooting](/manuals/desktop/troubleshoot-and-support/troubleshoot/_index.md) describes common problems, workarounds, how to run and submit diagnostics, and submit issues.
- [FAQs](/manuals/desktop/troubleshoot-and-support/faqs/general.md) provide answers to frequently asked questions.
- [Release notes](/manuals/desktop/release-notes.md) lists component updates, new features, and improvements associated with Docker Desktop releases.
- [Back up and restore data](/manuals/desktop/settings-and-maintenance/backup-and-restore.md) provides instructions
  on backing up and restoring data related to Docker.
