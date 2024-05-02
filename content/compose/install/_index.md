---
description: Learn how to install Docker Compose. Compose is available natively on
  Docker Desktop, as a Docker Engine plugin, and as a standalone tool.
keywords: install docker compose, docker compose install, install docker compose ubuntu,
  installing docker compose, docker compose download, docker compose not found, docker
  compose windows, how to install docker compose
title: Overview of installing Docker Compose
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
- [Linux](../../desktop/install/linux-install.md)
- [Mac](../../desktop/install/mac-install.md)
- [Windows](../../desktop/install/windows-install.md)

If you have already installed Docker Desktop, you can check which version of Compose you have by selecting **About Docker Desktop** from the Docker menu {{< inline-image src="../../desktop/images/whale-x.svg" alt="whale menu" >}}.

### Scenario two: Install the Compose plugin

If you already have Docker Engine and Docker CLI installed, you can install the Compose plugin from the command line, by either:
- [Using Docker's repository](linux.md#install-using-the-repository)
- [Downloading and installing manually](linux.md#install-the-plugin-manually)

> **Important**
>
>This is only available on Linux
{ .important }

### Scenario three: Install the Compose standalone 

You can [install the Compose standalone](standalone.md) on Linux or on Windows Server.

> **Warning**
>
>This install scenario is not recommended and is only supported for backward compatibility purposes.
{ .warning }
