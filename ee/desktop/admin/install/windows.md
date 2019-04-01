---
title: Install Docker Desktop Enterprise on Windows
description: Learn about Docker Desktop Enterprise
keywords: Docker EE, Windows, Mac, Docker Desktop, Enterprise
---

This page contains information about the system requirements and specific instructions that help you install Docker Desktop Enterprise (DDE) on Windows.

> **Warning:** If you are using the community version of Docker Desktop, you must uninstall Docker Desktop community in order to install Docker Desktop Enterprise.

# System requirements

- Windows 10 Pro or Enterprise version 15063 or later.

- Hyper-V and Containers Windows features must be enabled.

- The following hardware prerequisites are required to successfully run Client
Hyper-V on Windows 10:

  - 64 bit processor with [Second Level Address Translation (SLAT)](http://en.wikipedia.org/wiki/Second_Level_Address_Translation)

  - 4GB system RAM

  - BIOS-level hardware virtualization support must be enabled in the
    BIOS settings:

![Virtualization Technology (VTx) must be enabled in BIOS settings](.../images/windows-prereq.png "BIOS setting information for hardware virtualization support")

# Installation

Download Docker Desktop Enterprise for [**Windows**](https://download.docker.com/win/enterprise/DockerDesktop.msi).

The DDE installer includes Docker Engine, Docker CLI client, and Docker Compose.

Double-click the `.msi` file to begin the installation and follow the on-screen instructions.

# License file

You must install the Docker Desktop Enterprise license file at the following location:

`C:\Users\Docker\AppData\Roaming\Docker\docker_subscription.lic`

If the license file is missing, you will be asked to provide it when you try to run Docker Desktop Enterprise.

# Firewall exceptions

Docker Desktop Enterprise requires the following firewall exceptions. If you do not have firewall access, or are unsure about how to set firewall exceptions, contact your system administrator.

- The process `com.docker.vpnkit` proxies all outgoing container TCP and
    UDP traffic. This includes Docker image downloading but not DNS
    resolution, which is performed over a Unix domain socket connected
    to the `mDNSResponder` system service.

- The process `com.docker.vpnkit` binds external ports on behalf of
    containers. For example, `docker run -p 80:80 nginx` binds port 80 on all
    interfaces.

- If using Kubernetes, the API server is exposed with TLS on
    `127.0.0.1:6443` by `com.docker.vpnkit`.

# Version packs

Docker Desktop Enterprise is bundled with default version pack [Enterprise 2.1 (Docker
Engine 18.09 / Kubernetes 1.11.5)](https://download.docker.com/win/enterprise/enterprise-2.1.ddvp). System administrators can install versions packs using a CLI tool to use a different version of the Docker Engine and Kubernetes for development work:

- [Docker Community (18.09/Kubernetes
    1.13.0)](https://download.docker.com/win/enterprise/community.ddvp)

- [Docker Enterprise 2.0 (17.06/Kubernetes
    1.8.11)](https://download.docker.com/win/enterprise/enterprise-2.0.ddvp)

For information on using the CLI tool for version pack installation, see [Command line installation](#command-line-installation).

Available version packs are listed within the **Version Selection** option in the Docker Desktop Enterprise menu. If more than one version pack is installed, you can select the corresponding entry to work with a different version pack. After you select a different version pack, Docker Desktop Enterprise restarts and the selected Docker Engine and Kubernetes versions are used.

# Command line installation

System administrators can use a command line executable to install and uninstall Docker Desktop Enterprise and version packs.

When you install Docker Desktop Enterprise, the command line tool is installed at the following location:

`[ApplicationPath]\dockerdesktop-admin.exe`

>**Note:** Command line installation is supported for administrators only. You must have `administrator` access to run the CLI commands.

## Version-pack install

Run the following command to install or upgrade a version pack to the version contained in the specified `.ddvp` archive:

`dockerdesktop-admin.exe -InstallVersionPack=[path-to-archive]`

>**Note:** You must stop Docker Desktop Enterprise before installing a version pack.

## Version-pack uninstall

Run the following command to uninstall the specified version pack:

 `dockerdesktop-admin.exe -UninstallVersionPack=[version-pack-name]`

>**Note:** You must stop Docker Desktop Enterprise before uninstalling a version pack.

## Application uninstall

To uninstall the application, use system settings or the `.msi` file.