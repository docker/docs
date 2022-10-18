---
description: Instructions for installing Docker Desktop Arch package. Mostly meant for hackers who want to try out Docker Desktop on a variety of Arch-based distributions.
keywords: Arch Linux, install, uninstall, upgrade, update, linux, desktop, docker desktop, docker desktop for linux, dd4l
title: Install Docker Desktop on Arch-based distributions
redirect_from:
- /desktop/linux/install/archlinux/
---

This topic discusses installation of Docker Desktop from an [Arch package](https://desktop-stage.docker.com/linux/main/amd64/78459/docker-desktop-4.8.0-x86_64.pkg.tar.zst?utm_source=docker&utm_medium=webreferral&utm_campaign=docs-driven-download-linux-amd64) that Docker provides in addition to the supported platforms. Docker has not tested or verified the installation.

[Arch package (experimental)](https://desktop.docker.com/linux/main/amd64/docker-desktop-4.12.0-x86_64.pkg.tar.zst?utm_source=docker&utm_medium=webreferral&utm_campaign=docs-driven-download-linux-amd64){: .button .primary-btn }

*For checksums, see [Release notes](../release-notes.md)*

## Prerequisites

To install Docker Desktop successfully, you must meet the [system requirements](linux-install.md#system-requirements).

Additionally, for non-Gnome Desktop environments, `gnome-terminal` must be installed:

```console
$ sudo pacman -S gnome-terminal
```


## Install Docker Desktop 

1. Install client binaries. Docker does not have an Arch package repository. Binaries not included in the package must be installed manually before installing Docker Desktop. 

2. [Install Docker client binary on Linux](../../engine/install/binaries.md#install-daemon-and-client-binaries-on-linux). On Arch-based distributions, users must install the Docker client binary.
Static binaries for the Docker client are available for Linux (as `docker`).

3. Download the Arch package from the [release](../release-notes.md) page.

4. Install the package:

```console
$ sudo pacman -U ./docker-desktop-<version>-<arch>.pkg.tar.zst
```


## Launch Docker Desktop

{% include desktop-linux-launch.md %}

## Next steps

- Take a look at the [Get started](../../get-started/index.md) training modules to learn  how to build an image and run it as a containerized application.
- Review the topics in [Develop with Docker](../../develop/index.md) to learn how to build new applications using Docker.
