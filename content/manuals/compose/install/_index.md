---
description: Learn how to install Docker Compose. Compose is available natively on
  Docker Desktop, as a Docker Engine plugin, and as a standalone tool.
keywords: install docker compose, docker compose install, install docker compose ubuntu,
  installing docker compose, docker compose download, docker compose not found, docker
  compose windows, how to install docker compose
title: Overview of installing Docker Compose
linkTitle: Install
weight: 20
toc_max: 3
aliases:
- /compose/compose-desktop/
- /compose/install/other/
- /compose/install/compose-desktop/
---

This page contains summary information about the available options for installing Docker Compose.

## Installation scenarios 

### Scenario one: Install Docker Desktop

The easiest and recommended way to get Docker Compose is to install Docker Desktop. Docker Desktop
includes Docker Compose along with Docker Engine and Docker CLI which are Compose prerequisites. 

Docker Desktop is available on:
- [Linux](/manuals/desktop/setup/install/linux/_index.md)
- [Mac](/manuals/desktop/setup/install/mac-install.md)
- [Windows](/manuals/desktop/setup/install/windows-install.md)

If you have already installed Docker Desktop, you can check which version of Compose you have by selecting **About Docker Desktop** from the Docker menu {{< inline-image src="../../desktop/images/whale-x.svg" alt="whale menu" >}}.

> [!NOTE] 
>
> After Docker Compose V1 was removed in Docker Desktop version [4.23.0](/desktop/release-notes/#4230) as it had reached end-of-life,
> the `docker-compose` command now points directly to the Docker Compose V2 binary, running in standalone mode. 
> If you rely on Docker Desktop auto-update, the symlink might be broken and command unavailable, as the update doesn't ask for administrator password. 
> 
> This only affects Mac users. To fix this, either recreate the symlink:
> ```console
> $ sudo rm /usr/local/bin/docker-compose
> $ sudo ln -s /Applications/Docker.app/Contents/Resources/cli-plugins/docker-compose /usr/local/bin/docker-compose
> ```
> Or enable [Automatically check configuration](/manuals/desktop/settings-and-maintenance/settings.md) which will detect and fix it for you.

### Scenario two: Install the Docker Compose plugin

> [!IMPORTANT]
>
> This install scenario is only available on Linux.

If you already have Docker Engine and Docker CLI installed, you can install the Docker Compose plugin from the command line, by either:
- [Using Docker's repository](linux.md#install-using-the-repository)
- [Downloading and installing manually](linux.md#install-the-plugin-manually)

### Scenario three: Install the Docker Compose standalone 

> [!WARNING]
>
> This install scenario is not recommended and is only supported for backward compatibility purposes.

You can [install the Docker Compose standalone](standalone.md) on Linux or on Windows Server.

