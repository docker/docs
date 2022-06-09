---
description: How to install Docker Compose
keywords: compose, orchestration, install, installation, docker, documentation
title: Install Docker Compose
toc_max: 2
---

On this page you can find a summary of the available options for installing Docker Compose. 

## Compose prerequisites

* Docker Compose requires Docker Engine.
* Docker Compose plugin requires Docker CLI.

## Compose installation scenarios
You can run Compose on macOS, Windows, and 64-bit Linux. Check what installation scenario fits your needs. 

Are you looking to:

* __Get latest Docker Compose and its prerequisites__:
[Install Docker Desktop for your platform](./compose-desktop.md). This is the fastest route and you get Docker Engine and Docker CLI with the Compose plugin. Docker Desktop is available for Mac, Windows and Linux.

* __Install Compose plugin:__
  + __(Mac, Win, Linux) Docker Desktop__: If you have Desktop installed then you already have the Compose plugin installed.
  + __Linux systems__: To install the Docker CLI's Compose plugins use one of these methods of installation: 
     + Using the [convenience scripts](../../engine/install/#server){: target="_blank" rel="noopener" class="_"} offered per Linux distro from the Engine install section.
     + [Setting up Docker's repository](compose-plugin#install-using-the-repository) and using it to install the compose plugin package.     
     + Other scenarios, check the [Linux install](compose-plugin#installing-compose-on-linux-systems).
  + __Windows Server__: If you want to run the Docker daemon and client directly on Microsoft Windows Server, follow the [Windows Server install instructions](compose-plugin#install-compose-on-windows-server).

## Where to go next

- [Getting Started](../gettingstarted.md)
- [Command line reference](../../reference/index.md)
- [Compose file reference](../compose-file/index.md)
- [Sample apps with Compose](../samples-for-compose.md)
