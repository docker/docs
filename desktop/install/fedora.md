---
description: Instructions for installing Docker Desktop on Fedora
keywords: fedora, rpm, update install, uninstall, upgrade, update, linux, desktop, docker desktop, docker desktop for linux, dd4l
title: Install Docker Desktop on Fedora
toc_max: 4
---

This page contains information on how to install, launch and upgrade Docker Desktop on a Fedora distribution.

## Prerequisites

To install Docker Desktop successfully, you must:
- Meet the [system requirements](../install.md#system-requirements).
- Have a 64-bit version of either Fedora 35 or Fedora 36.

Additionally, for a Gnome Desktop environment you must install AppIndicator and KStatusNotifierItem [Gnome extensions](https://extensions.gnome.org/extension/615/appindicator-support/).

For non-Gnome Desktop environments, `gnome-terminal` must be installed:

```console
$ sudo dnf install gnome-terminal
```

## Install Docker Desktop

To install Docker Desktop on Fedora:

1. Set up [Docker's package repository](../../../engine/install/fedora.md#set-up-the-repository). 

2. Download latest RPM package from the [release](../release-notes/index.md) page.

3. Install the package with dnf as follows:
    
```console
$ sudo dnf install ./docker-desktop-<version>-<arch>.rpm
```

There are a few post-install configuration steps done through the post-install script contained in the RPM package.

The post-install script:

- Sets the capability on the Docker Desktop binary to map privileged ports and set resource limits.
- Adds a DNS name for Kubernetes to `/etc/hosts`.
- Creates a link from `/usr/bin/docker` to `/usr/local/bin/com.docker.cli`.

## Launch Docker Desktop


{% include desktop-linux-launch.md %}


## Upgrade Docker Desktop

Once a new version for Docker Desktop is released, the Docker UI shows a notification. 
You need to first remove the previous version and then download the new package each time you want to upgrade Docker Desktop. Run:

```console
$ sudo dnf remove docker-desktop
$ sudo dnf install ./docker-desktop-<version>-<arch>.rpm
```


## Uninstall Docker Desktop

To remove Docker Desktop for Linux, run:

```console
$ sudo dnf remove docker-desktop
```

For a complete cleanup, remove configuration and data files at `$HOME/.docker/desktop`, the symlink at `/usr/local/bin/com.docker.cli`, and purge
the remaining systemd service files.

```console
$ rm -r $HOME/.docker/desktop
$ sudo rm /usr/local/bin/com.docker.cli
```

Remove the `credsStore` and `currentContext` properties from `$HOME/.docker/config.json`. Additionally, you must delete any edited configuration files manually. 

## Next steps

- Take a look at the [Get started](../../../get-started/index.md) training modules to learn  how to build an image and run it as a containerized application.
- Review the topics in [Develop with Docker](../../../develop/index.md) to learn how to build new applications using Docker.
