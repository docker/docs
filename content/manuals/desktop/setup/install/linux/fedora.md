---
description: Instructions for installing Docker Desktop on Fedora
keywords: fedora, rpm, update install, uninstall, upgrade, update, linux, desktop,
  docker desktop, docker desktop for linux, dd4l
title: Install Docker Desktop on Fedora
linkTitle: Fedora
weight: 30
toc_max: 4
aliases:
- /desktop/linux/install/fedora/
- /desktop/install/fedora/
- /desktop/install/linux/fedora/
---

> **Docker Desktop terms**
>
> Commercial use of Docker Desktop in larger enterprises (more than 250
> employees OR more than $10 million USD in annual revenue) requires a [paid
> subscription](https://www.docker.com/pricing/).

This page contains information on how to install, launch and upgrade Docker Desktop on a Fedora distribution.

## Prerequisites

To install Docker Desktop successfully, you must:

- Meet the [general system requirements](_index.md#general-system-requirements).
- Have a 64-bit version of Fedora 40 or Fedora 41.

Additionally, for a GNOME desktop environment you must install AppIndicator and KStatusNotifierItem [GNOME extensions](https://extensions.gnome.org/extension/615/appindicator-support/).

For non-GNOME desktop environments, `gnome-terminal` must be installed:

```console
$ sudo dnf install gnome-terminal
```

## Install Docker Desktop

To install Docker Desktop on Fedora:

1. Set up [Docker's package repository](/manuals/engine/install/fedora.md#set-up-the-repository).

2. Download the latest [RPM package](https://desktop.docker.com/linux/main/amd64/docker-desktop-x86_64.rpm?utm_source=docker&utm_medium=webreferral&utm_campaign=docs-driven-download-linux-amd64). For checksums, see the [Release notes](/manuals/desktop/release-notes.md).

3. Install the package with dnf as follows:

   ```console
   $ sudo dnf install ./docker-desktop-x86_64.rpm
   ```

   By default, Docker Desktop is installed at `/opt/docker-desktop`.

There are a few post-install configuration steps done through the post-install script contained in the RPM package.

The post-install script:

- Sets the capability on the Docker Desktop binary to map privileged ports and set resource limits.
- Adds a DNS name for Kubernetes to `/etc/hosts`.
- Creates a symlink from `/usr/local/bin/com.docker.cli` to `/usr/bin/docker`.
  This is because the classic Docker CLI is installed at `/usr/bin/docker`. The Docker Desktop installer also installs a Docker CLI binary that includes cloud-integration capabilities and is essentially a wrapper for the Compose CLI, at`/usr/local/bin/com.docker.cli`. The symlink ensures that the wrapper can access the classic Docker CLI. 

## Launch Docker Desktop

{{% include "desktop-linux-launch.md" %}}

## Upgrade Docker Desktop

Once a new version for Docker Desktop is released, the Docker UI shows a notification.
You need to first remove the previous version and then download the new package each time you want to upgrade Docker Desktop. Run:

```console
$ sudo dnf remove docker-desktop
$ sudo dnf install ./docker-desktop-x86_64.rpm
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
