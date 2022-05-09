---
description: Instructions for installing Docker Desktop on Fedora
keywords: requirements, apt, installation, fedora, rpm, install, uninstall, upgrade, update
title: Install Docker Desktop on Fedora
toc_max: 4
---

To get started with Docker Desktop on Fedora, make sure you
[meet the prerequisites](#prerequisites), then
[install Docker Desktop](#install-docker-desktop).

## Prerequisites

Your Fedora distribution must meet the [system requirements](../install.md#system-requirements) to install and launch Docker Desktop successfully.

Additionally, for a Gnome Desktop environment you must install AppIndicator and KStatusNotifierItem [Gnome extensions](https://extensions.gnome.org/extension/615/appindicator-support/).

For non-Gnome Desktop environments, `gnome-terminal` must be installed:

```console
$ sudo dnf install gnome-terminal
```

### OS requirements

To install Docker Desktop for Linux, you need the 64-bit version of one of these Fedora versions:

- Fedora 35
- Fedora 36

## Install Docker Desktop

Recommended approach to install Docker Desktop on Ubuntu:

1. [Set up Docker's RPM repository](../../../engine/install/fedora.md#set-up-the-repository) 

2. Download latest RPM package from the [release](../release-notes/index.md) page.

3. Install the package with dnf as follows:
    
```console
$ sudo dnf install ./docker-desktop-<version>-<arch>.rpm
```

There are a few post-install configuration steps done through the maintainers' scripts (post-install script contained in the RPM package.
The post-install script:

- sets the capability on the Docker Desktop binary to map privileged ports and set resource limits
- adds a DNS name for Kubernetes to `/etc/hosts`
- creates a link from `/usr/bin/docker` to `/usr/local/bin/com.docker.cli`

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
