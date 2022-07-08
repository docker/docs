---
description: DD for Linux on Raspoberry Pi
keywords: docker, Docker Desktop for linux,  raspberry pi
title: Install Docker Desktop on Raspberry Pi
---
This page contains information on how to install, launch, and upgrade Docker Desktop on a Raspberyy Pi. 

> **Note**
>
> Performance of Docker Desktop may be slower on your Raspberry Pi than on a macOS, Windows or Linux machine.

## Prerequisites 
To install Docker Desktop successfully, you must:

- Meet the [system requirements](../../install/#system-requirements)
- Have a 64-bit version of Raspberyy Pi OS
- Have a fast SD card
- Have a stable power supply
- Have a minimum of 4GB RAM 


## Install Docker Desktop

To install Docker Desktop for Linux on your Raspberry Pi:


1. Set up [Docker's package repository](../../../engine/install/debian.md#set-up-the-repository).

2. Download the latest DEB package from the [release](../../release-notes.md) page.

3. Install the package with apt as follows:

```console
$ sudo apt-get update
$ sudo apt-get install ./docker-desktop-<version>-arm64.deb
```


There are a few post-install configuration steps done through the maintainers' scripts in the DEB package.

The post-install script:

- Sets the capability on the Docker Desktop binary to map privileged ports and set resource limits.
- Adds a DNS name for Kubernetes to `/etc/hosts`.
- Creates a link from `/usr/bin/docker` to `/usr/local/bin/com.docker.cli`.

## Launch Docker Desktop

> **Note**
>
> It can take a couple of minutes to load Docker Desktop when you start it for the first time. 

To start Docker Desktop you can either:
- Search for **Docker Desktop** in your **Applications** menu and open it. This launches the whale menu icon and opens **Docker Dashboard**. 
- Open a terminal and run:

```console
$ systemctl --user start docker-desktop
```
### Docker Desktop vs. Docker Engine locally

When Docker Desktop starts, it creates a dedicated context that the Docker CLI
can use as a target and sets it as the current context in use. This is to avoid
a clash with a local Docker Engine that may be running on the host and
using the default context. On shutdown, Docker Desktop resets the current
context to the previous one.

The Docker Desktop installer updates Docker Compose and the Docker CLI binaries
on the host. It installs Docker Compose V2 and gives users the choice to
link it as docker-compose in **Preferences**. Docker Desktop installs
the new Docker CLI binary that includes cloud-integration capabilities in `/usr/local/bin`
and creates a symlink to the classic Docker CLi at `/usr/local/bin/com.docker.cli`.

After you’ve successfully installed Docker Desktop, you can check the versions
of these binaries by running the following command:

```console
$ docker compose version
Docker Compose version v2.5.0

$ docker --version
Docker version 20.10.14, build a224086349

$ docker version
Client: Docker Engine - Community
Cloud integration: 1.0.24
Version:           20.10.14
API version:       1.41
...
```
## Start Docker Desktop on Log in
To enable Docker Desktop to start when you log in to your Raspberry Pi, you can either: 
- Navigate to **Settings** in Docker Desktop and then select **General** > **Start Docker Desktop when you log in**.
- Open a terminal and run:

```console
$ systemctl --user enable docker-desktop
```
## Stop Docker Desktop 
To stop Docker Desktop, you can either:
- Click on the whale menu tray icon to open the Docker menu and then select **Quit Docker Desktop**
- Open a terminal and run:

```console
$ systemctl --user stop docker-desktop
```

## Upgrade Docker Desktop

Once a new version for Docker Desktop is released, the Docker UI shows a notification. You need to first remove the previous version and then download the new package each time you want to upgrade Docker Desktop. Run:

```console
$ sudo dnf remove docker-desktop
$ sudo dnf install ./docker-desktop-<version>-<arch>.rpm
```



## Uninstall 

To remove Docker Desktop for Linux, run:

```console
$ sudo apt remove docker-desktop
```

For a complete cleanup, remove configuration and data files at `$HOME/.docker/desktop`, the symlink at `/usr/local/bin/com.docker.cli`, and purge
the remaining systemd service files.

```console
$ rm -r $HOME/.docker/desktop
$ sudo rm /usr/local/bin/com.docker.cli
$ sudo apt purge docker-desktop
```

Remove the `credsStore` and `currentContext` properties from `$HOME/.docker/config.json`. Additionally, you must delete any edited configuration files manually.

## What next?

- Explore [Docker Dashboard](../../dashboard.md)
- Start a [sample application](../../dashboard.md/#start-a-sample-application) 