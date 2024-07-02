---
description: Instructions for installing Docker Desktop on Fedora
keywords: fedora, rpm, update install, uninstall, upgrade, update, linux, desktop,
  docker desktop, docker desktop for linux, dd4l
title: Install Docker Desktop on Fedora
toc_max: 4
aliases:
- /desktop/linux/install/fedora/
---

> **Docker Desktop terms**
>
> Commercial use of Docker Desktop in larger enterprises (more than 250
> employees OR more than $10 million USD in annual revenue) requires a [paid
> subscription](https://www.docker.com/pricing/).

This page contains information on how to install, launch and upgrade Docker Desktop on a Fedora distribution.

{{< button text="RPM package" url="https://desktop.docker.com/linux/main/amd64/156455/docker-desktop-x86_64.rpm?utm_source=docker&utm_medium=webreferral&utm_campaign=docs-driven-download-linux-amd64" >}}

## Prerequisites

To install Docker Desktop successfully, you must:

- Meet the [system requirements](linux-install.md#general-system-requirements).
- Have a 64-bit version of Fedora 39 or Fedora 40.

Additionally, for a GNOME desktop environment you must install AppIndicator and KStatusNotifierItem [GNOME extensions](https://extensions.gnome.org/extension/615/appindicator-support/).

For non-GNOME desktop environments, `gnome-terminal` must be installed:

```console
$ sudo dnf install gnome-terminal
```

## Install Docker Desktop

To install Docker Desktop on Fedora:

1. Set up [Docker's package repository](../../engine/install/fedora.md#set-up-the-repository).

2. Download latest [RPM package](https://desktop.docker.com/linux/main/amd64/156455/docker-desktop-x86_64.rpm?utm_source=docker&utm_medium=webreferral&utm_campaign=docs-driven-download-linux-amd64).

3. Install the package with dnf as follows:

   ```console
   $ sudo dnf install ./docker-desktop-<arch>.rpm
   ```

There are a few post-install configuration steps done through the post-install script contained in the RPM package.

The post-install script:

- Sets the capability on the Docker Desktop binary to map privileged ports and set resource limits.
- Adds a DNS name for Kubernetes to `/etc/hosts`.
- Creates a symlink from `/usr/local/bin/com.docker.cli` to `/usr/bin/docker`.
  This is because the classic Docker CLI is installed at `/usr/bin/docker`. The Docker Desktop installer also installs a Docker CLI binary that includes cloud-integration capabilities and is essentially a wrapper for the Compose CLI, at`/usr/local/bin/com.docker.cli`. The symlink ensures that the wrapper can access the classic Docker CLI. 

## Launch Docker Desktop

{{< include "desktop-linux-launch.md" >}}

## Upgrade Docker Desktop

Once a new version for Docker Desktop is released, the Docker UI shows a notification.
You need to first remove the previous version and then download the new package each time you want to upgrade Docker Desktop. Run:

```console
$ sudo dnf remove docker-desktop
$ sudo dnf install ./docker-desktop-<arch>.rpm
```

## Next steps

- Take a look at the [Docker workshop](../../guides/workshop/_index.md) to learn how to build an image and run it as a containerized application.
- [Explore Docker Desktop](../use-desktop/index.md) and all its features.
