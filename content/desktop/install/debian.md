---
description: Instructions for installing Docker Desktop on Debian
keywords: debian, install, uninstall, upgrade, update, linux, desktop, docker desktop,
  docker desktop for linux, dd4l
title: Install Docker Desktop on Debian
toc_max: 4
aliases:
- /desktop/linux/install/debian/
---

This page contains information on how to install, launch, and upgrade Docker Desktop on a Debian distribution.

{{< button text="DEB package" url="https://desktop.docker.com/linux/main/amd64/docker-desktop-4.26.1-amd64.deb?utm_source=docker&utm_medium=webreferral&utm_campaign=docs-driven-download-linux-amd64" >}}

_For checksums, see [Release notes](../release-notes.md)_

## Prerequisites

To install Docker Desktop successfully, you must:

- Meet the [system requirements](linux-install.md#system-requirements).
- Have a 64-bit version of Debian 12.
- For a Gnome Desktop environment, you must also install AppIndicator and KStatusNotifierItem [Gnome extensions](https://extensions.gnome.org/extension/615/appindicator-support/).

- For non-Gnome Desktop environments, `gnome-terminal` must be installed:

  ```console
  $ sudo apt install gnome-terminal
  ```

## Install Docker Desktop

Recommended approach to install Docker Desktop on Debian:

1. Set up Docker's `apt` repository.
   See step one of [Install using the `apt` repository](../../engine/install/debian.md#install-using-the-repository).

2. Download latest [DEB package](https://desktop.docker.com/linux/main/amd64/docker-desktop-4.26.1-amd64.deb?utm_source=docker&utm_medium=webreferral&utm_campaign=docs-driven-download-linux-amd64).

3. Install the package with apt as follows:

  ```console
  $ sudo apt-get update
  $ sudo apt-get install ./docker-desktop-<version>-<arch>.deb
  ```

  > **Note**
  >
  > At the end of the installation process, `apt` displays an error due to installing a downloaded package. You
  > can ignore this error message.
  >
  > ```text
  > N: Download is performed unsandboxed as root, as file '/home/user/Downloads/docker-desktop.deb' couldn't be accessed by user '_apt'. - pkgAcquire::Run (13: Permission denied)
  > ```

There are a few post-install configuration steps done through the post-install script contained in the deb package.

The post-install script:

- Sets the capability on the Docker Desktop binary to map privileged ports and set resource limits.
- Adds a DNS name for Kubernetes to `/etc/hosts`.
- Creates a symlink from `/usr/local/bin/com.docker.cli` to `/usr/bin/docker`.
  This is because the classic Docker CLI is installed at `/usr/bin/docker`. The Docker Desktop installer also installs a Docker CLI binary that includes cloud-integration capabilities and is essentially a wrapper for the Compose CLI, at`/usr/local/bin/com.docker.cli`. The symlink ensures that the wrapper can access the classic Docker CLI. 

## Launch Docker Desktop

{{< include "desktop-linux-launch.md" >}}

## Upgrade Docker Desktop

Once a new version for Docker Desktop is released, the Docker UI shows a notification.
You need to download the new package each time you want to upgrade Docker Desktop and run:

```console
$ sudo apt-get install ./docker-desktop-<version>-<arch>.deb
```

## Next steps

- Take a look at the [Get started](../../guides/get-started/_index.md) training modules to learn how to build an image and run it as a containerized application.
- [Explore Docker Desktop](../use-desktop/index.md) and all its features.
- Review the topics in [Develop with Docker](../../develop/index.md) to learn how to build new applications using Docker.
