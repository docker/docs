---
description: How to install Docker Compose
keywords: compose, orchestration, install, installation, docker, documentation
title: Overview
toc_max: 3
redirect_from:
- /compose/compose-desktop/
---

This page contains summary information about the available options for getting Docker Compose.

## Installation scenarios 

### Scenario one: Install Docker Desktop

The easiest and recommended way to get Docker Compose is to install Docker Desktop. Docker Desktop
includes Docker Compose along with Docker Engine and Docker CLI which are Compose prerequisites. 

Docker Desktop is available on:
- [Linux](../../desktop/install/linux-install.md)
- [Mac](../../desktop/install/mac-install.md)
- [Windows](../../desktop/install/windows-install.md)

If you have already installed Docker Desktop, you can check which version of Compose you have by selecting **About Docker Desktop** from the Docker menu ![whale menu](../../desktop/images/whale-x.svg){: .inline}

### Scenario two: Install the Compose plugin

If you already have Docker Engine and Docker CLI installed, you can install the Compose plugin from the command line, by either:
- [Using Docker's repository](linux.md#install-using-the-repository)
- [Downloading and installing manually](linux.md#install-the-plugin-manually)

>Note
>
>This is only available on Linux
{: .important}

### Scenario three: Install the Compose standalone 

You can [install the Compose standalone](other.md) on Linux or on Windows Server.

>Note
>
>This install scenario is no longer supported.
{: .important}