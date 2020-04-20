---
title: Install Docker Desktop Enterprise on Mac
description: Learn about Docker Desktop Enterprise
keywords: Docker EE, Mac, Docker Desktop, Enterprise
redirect_from:
- /ee/desktop/admin/install/mac/
---

This page contains information about the system requirements and specific instructions that help you install Docker Desktop Enterprise (DDE) on Mac. If you are using the Community version of Docker Desktop, you must uninstall Docker Desktop Community in order to install DDE.

[Download Docker Desktop Enterprise for Mac](https://download.docker.com/mac/enterprise/Docker.pkg){: .button .outline-btn}

> **Note:** By downloading DDE, you agree to the terms of the [Docker Software End User License Agreement](https://www.docker.com/legal/docker-software-end-user-license-agreement){: target="_blank" class="_"} and the [Docker Data Processing Agreement (DPA)](https://www.docker.com/legal/data-processing-agreement){: target="_blank" class="_"}.

## System requirements

- Mac hardware must be a 2010 or newer model, with Intel’s hardware support for memory management unit (MMU) virtualization, including Extended Page Tables (EPT) and Unrestricted Mode. You can check to see if your machine has this support by running the following command in a terminal: `sysctl kern.hv_support`

- macOS must be version 10.13 or newer. We recommend upgrading to the latest version of macOS.

    If you experience any issues after upgrading your macOS to version 10.15, you must install the latest version of Docker Desktop to be compatible with this version of macOS.

- At least 4GB of RAM

- VirtualBox prior to version 4.3.30 must NOT be installed (it is incompatible with Docker for Mac). If you have a newer version of VirtualBox installed, it’s fine.

> **Note:** Docker supports Docker Desktop Enterprise on the most recent versions of macOS. That is, the current release of macOS and the previous two releases. As new major versions of macOS are made generally available, Docker will stop supporting the oldest version and support the newest version of macOS (in addition to the previous two releases).

## Installation

The DDE installer includes Docker Engine, Docker CLI client, and Docker Compose.

Double-click the `.pkg` file to begin the installation and follow the on-screen instructions. When the installation is complete, click the Launchpad icon in the Dock and then **Docker** to start Docker Desktop.

Mac administrators can use the command line option `\$ sudo installer -pkg Docker.pkg -target /` for fine tuning and mass installation. After running this command, you can start Docker Desktop from the Applications folder on each machine.

Administrators can configure additional settings by modifying the administrator configuration file. For more information, see [Configure Desktop Enterprise for Mac](/desktop/enterprise/admin/configure/mac-admin).

## License file

Install the Docker Desktop Enterprise license file at the following location:

`/Library/Group Containers/group.com.docker/docker_subscription.lic`

You must create the path if it doesn't already exist. If the license file is missing, you will be asked to provide it when you try to run Docker Desktop Enterprise. Contact your system administrator to obtain the license file.

## Firewall exceptions

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

## Version packs

Docker Desktop Enterprise is bundled with default version pack [Enterprise 3.0 (Docker Engine 19.03 / Kubernetes 1.14)](https://download.docker.com/mac/enterprise/enterprise-3.0.ddvp). System administrators can install version packs using a command line tool to use a different version of the Docker Engine and Kubernetes for development work:

- [Docker Enterprise 2.0 (17.06/Kubernetes 1.8.11)](https://download.docker.com/mac/enterprise/enterprise-2.0.ddvp)

- [Docker Enterprise 2.1 (18.09/Kubernetes 1.11.5)](https://download.docker.com/mac/enterprise/enterprise-2.1.ddvp)

For information on using the CLI tool for version pack installation, see [Command line installation](#command-line-installation).

> **Note:** It is not possible to install the version packs using the Docker Desktop user interface or by double-clicking the `.ddvp` file.

Available version packs are listed within the **Version Selection** option in the Docker Desktop menu. If more than one version pack is installed, you can select the corresponding entry to work with a different version pack. After you select a different version pack, Docker Desktop restarts and the selected Docker Engine and Kubernetes versions are used.

If more than one version pack is installed, you can select the corresponding entry to work with a different version pack. After you select a different version pack, Docker Desktop restarts and the selected Docker Engine and Kubernetes versions are used.

## Command line installation

System administrators can use a command line executable to install and uninstall Docker Desktop Enterprise and version packs.

When you install Docker Desktop Enterprise, the command line tool is installed at the following location:

[ApplicationPath]/Contents/Resources/bin/dockerdesktop-admin

>**Note:** Command line installation is supported for administrators only. You must have `sudo` access privilege to run the CLI commands.

### Version-pack install

Run the following command to install or upgrade a version pack to the version contained in the specified `.ddvp` archive:

    sudo /Applications/Docker.app/Contents/Resources/bin/dockerdesktop-admin version-pack install [path-to-archive]

 >**Note:** You must stop Docker Desktop before installing a version pack.

### Version-pack uninstall

 Run the following command to uninstall the specified version pack:

    sudo /Applications/Docker.app/Contents/Resources/bin/dockerdesktop-admin version-pack uninstall [version-pack-name]

>**Note:** You must stop Docker Desktop before uninstalling a version pack.

### Application uninstall

Run the following command to uninstall the application:

    sudo /Applications/Docker.app/Contents/Resources/bin/dockerdesktop-admin app uninstall

The `sudo` command uninstalls files such as version packs that are installed by an administrator, but are not accessible by users.
