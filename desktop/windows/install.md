---
description: How to install Docker Desktop for Windows
keywords: windows, install, download, run, docker, local
title: Install Docker Desktop on Windows
redirect_from:
- /docker-for-windows/install/
- /docker-ee-for-windows/install/
- /docker-for-windows/install-windows-home/
- /ee/docker-ee/windows/docker-ee/
- /engine/installation/windows/docker-ee/
- /install/windows/docker-ee/
- /install/windows/ee-preview/
---

Welcome to Docker Desktop for Windows. This page contains information about Docker Desktop for Windows system requirements, download URL, instructions to install and update Docker Desktop for Windows.

> Download Docker Desktop for Windows
>
> [Docker Desktop for Windows](https://desktop.docker.com/win/main/amd64/Docker%20Desktop%20Installer.exe){: .button .primary-btn }

## System requirements

Your Windows machine must meet the following requirements to successfully install Docker Desktop.

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#win-wsl2">WSL 2 backend</a></li>
<li><a data-toggle="tab" data-target="#win-hyper-v">Hyper-V backend and Windows containers</a></li>
</ul>
<div class="tab-content">
<div id="win-wsl2" class="tab-pane fade in active" markdown="1">

### WSL 2 backend

- Windows 10 64-bit: Home or Pro 2004 (build 19041) or higher, or Enterprise or Education 1909 (build 18363) or higher.
- Enable the WSL 2 feature on Windows. For detailed instructions, refer to the
    [Microsoft documentation](https://docs.microsoft.com/en-us/windows/wsl/install-win10){: target="_blank" rel="noopener" class="_"}.
- The following hardware prerequisites are required to successfully run
WSL 2 on Windows 10:

  - 64-bit processor with [Second Level Address Translation (SLAT)](https://en.wikipedia.org/wiki/Second_Level_Address_Translation){: target="_blank" rel="noopener" class="_"}
  - 4GB system RAM
  - BIOS-level hardware virtualization support must be enabled in the
    BIOS settings.  For more information, see
    [Virtualization](troubleshoot.md#virtualization-must-be-enabled).
- Download and install the [Linux kernel update package](https://docs.microsoft.com/windows/wsl/wsl2-kernel){: target="_blank" rel="noopener" class="_"}.

</div>
<div id="win-hyper-v" class="tab-pane fade" markdown="1">

### Hyper-V backend and Windows containers

- Windows 10 64-bit: Pro 2004 (build 19041) or higher, or Enterprise or Education 1909 (build 18363) or higher.

  For Windows 10 Home, see [System requirements for WSL 2 backend](#wsl-2-backend).
- Hyper-V and Containers Windows features must be enabled.
- The following hardware prerequisites are required to successfully run Client
Hyper-V on Windows 10:

  - 64 bit processor with [Second Level Address Translation (SLAT)](https://en.wikipedia.org/wiki/Second_Level_Address_Translation){: target="_blank" rel="noopener" class="_"}
  - 4GB system RAM
  - BIOS-level hardware virtualization support must be enabled in the
    BIOS settings.  For more information, see
    [Virtualization](troubleshoot.md#virtualization-must-be-enabled).

</div>
</div>

> **Note**
>
> Docker only supports Docker Desktop on Windows for those versions of Windows 10 that are still within [Microsoft’s servicing timeline](https://support.microsoft.com/en-us/help/13853/windows-lifecycle-fact-sheet){:target="_blank" rel="noopener" class="_"}.

### What's included in the installer

The Docker Desktop installation includes [Docker Engine](../../engine/index.md),
Docker CLI client, [Docker Compose](../../compose/index.md),
[Docker Content Trust](../../engine/security/trust/index.md),
[Kubernetes](https://github.com/kubernetes/kubernetes/),
and [Credential Helper](https://github.com/docker/docker-credential-helpers/).

Containers and images created with Docker Desktop are shared between all
user accounts on machines where it is installed. This is because all Windows
accounts use the same VM to build and run containers. Note that it is not possible to share containers and images between user accounts when using the Docker Desktop WSL 2 backend.

Nested virtualization scenarios, such as running Docker Desktop on a
VMWare or Parallels instance might work, but there are no guarantees. For
more information, see [Running Docker Desktop in nested virtualization scenarios](troubleshoot.md#running-docker-desktop-in-nested-virtualization-scenarios).

### About Windows containers

Looking for information on using Windows containers?

* [Switch between Windows and Linux containers](index.md#switch-between-windows-and-linux-containers)
  describes how you can toggle between Linux and Windows containers in Docker Desktop and points you to the tutorial mentioned above.
* [Getting Started with Windows Containers (Lab)](https://github.com/docker/labs/blob/master/windows/windows-containers/README.md)
  provides a tutorial on how to set up and run Windows containers on Windows 10, Windows Server 2016 and Windows Server 2019. It shows you how to use a MusicStore application
  with Windows containers.
* Docker Container Platform for Windows [articles and blog
  posts](https://www.docker.com/microsoft/) on the Docker website.

## Install Docker Desktop on Windows

1. Double-click **Docker Desktop Installer.exe** to run the installer.

    If you haven't already downloaded the installer (`Docker Desktop Installer.exe`), you can get it from
    [**Docker Hub**](https://hub.docker.com/editions/community/docker-ce-desktop-windows/).
    It typically downloads to your `Downloads` folder, or you can run it from
    the recent downloads bar at the bottom of your web browser.

2. When prompted, ensure the **Enable Hyper-V Windows Features** or the **Install required Windows components for WSL 2** option is selected on the Configuration page.

3. Follow the instructions on the installation wizard to authorize the installer and proceed with the install.

4. When the installation is successful, click **Close** to complete the installation process.

5. If your admin account is different to your user account, you must add the user to the **docker-users** group. Run **Computer Management** as an administrator and navigate to **Local Users and Groups** > **Groups** > **docker-users**. Right-click to add the user to the group.
Log out and log back in for the changes to take effect.

## Start Docker Desktop

Docker Desktop does not start automatically after installation. To start Docker Desktop:

1. Search for Docker, and select **Docker Desktop** in the search results.

      ![search for Docker app](images/docker-app-search.png){:width="300px"}

2. The Docker menu (![whale menu](images/whale-x.png){: .inline}) displays the Docker Subscription Service Agreement window. It includes a change to the terms of use for Docker Desktop.

    {% include desktop-license-update.md %}

3. Click the checkbox to indicate that you accept the updated terms and then click **Accept** to continue. Docker Desktop starts after you accept the terms.

    > **Important**
    >
    > If you do not agree to the updated terms, the Docker Desktop application will close and  you can no longer run Docker Desktop on your machine. You can choose to accept the terms at a later date by opening Docker Desktop.
    {: .important}

    For more information, see [Docker Desktop License Agreement](/subscription/#docker-desktop-license-agreement).

### Quick Start Guide

When the initialization is complete, Docker Desktop launches the **Quick Start Guide**. This tutorial includes a simple exercise to build an example Docker image, run it as a container, push and save the image to Docker Hub.

To run the Quick Start Guide on demand, right-click the Docker icon in the Notifications area (or System tray) to open the Docker Desktop menu and then select **Quick Start Guide**.

![Docker Quick Start tutorial](images/docker-tutorial-win.png){:width="450px"}

Congratulations! You are now successfully running Docker Desktop on Windows.

## Updates

{% include desktop-update.md %}

## Uninstall Docker Desktop

To uninstall Docker Desktop from your Windows machine:

1. From the Windows **Start** menu, select **Settings** > **Apps** > **Apps & features**.
2. Select **Docker Desktop** from the **Apps & features** list and then select **Uninstall**.
3. Click **Uninstall** to confirm your selection.

> **Important**
>
> Uninstalling Docker Desktop destroys Docker containers, images, volumes, and
> other Docker related data local to the machine, and removes the files generated
> by the application. Refer to the [back up and restore data](../backup-and-restore.md)
> section to learn how to preserve important data before uninstalling.

## Where to go next

* [Getting started](index.md) introduces Docker Desktop for Windows.
* [Get started with Docker](/get-started/) is a tutorial that teaches you how to
  deploy a multi-service stack.
* [Troubleshooting](troubleshoot.md) describes common problems, workarounds, and
  how to get support.
* [FAQs](../faqs.md) provide answers to frequently asked questions.
* [Release notes](release-notes.md) lists component updates, new features, and improvements associated with Docker Desktop releases.
* [Back up and restore data](../backup-and-restore.md) provides instructions on backing up and restoring data related to Docker.
