---
description: How to install Docker Desktop for Windows
keywords: windows, beta, edge, alpha, install, download
title: Install Docker Desktop on Windows
---

Docker Desktop for Windows is the [Community](https://www.docker.com/community-edition) version of Docker for Microsoft Windows.
You can download Docker Desktop for Windows from Docker Hub.

[Download from Docker
Hub](https://hub.docker.com/?overlay=onboarding){:
.button .outline-btn}

## What to know before you install

### System Requirements

  - Windows 10 64-bit: Pro, Enterprise, or Education (Build 15063 or later).
  - Hyper-V and Containers Windows features must be enabled.
  - The following hardware prerequisites are required to successfully run Client
Hyper-V on Windows 10:

     - 64 bit processor with [Second Level Address Translation (SLAT)](http://en.wikipedia.org/wiki/Second_Level_Address_Translation)
     - 4GB system RAM
    - BIOS-level hardware virtualization support must be enabled in the
    BIOS settings.  For more information, see
    [Virtualization](troubleshoot.md#virtualization-must-be-enabled).

> **Note:** Docker supports Docker Desktop on Windows based on Microsoft’s support lifecycle for Windows 10 operating system. For more information, see the [Windows lifecycle fact sheet](https://support.microsoft.com/en-us/help/13853/windows-lifecycle-fact-sheet).

**README for Docker Toolbox and Docker Machine users**: Microsoft Hyper-V is required to run Docker Desktop. The Docker Desktop Windows installer enables Hyper-V if required, and restarts your machine. When Hyper-V is enabled, VirtualBox no longer works. However, any existing VirtualBox VM images are retained.

VirtualBox VMs created with `docker-machine` (including the `default` one
typically created during Toolbox install) no longer start. These VMs cannot be
used side-by-side with Docker Desktop. However, you can still use
`docker-machine` to manage remote VMs.

### What's included in the installer

The Docker Desktop installation includes [Docker Engine](/engine/userguide/), Docker CLI client, [Docker Compose](/compose/overview.md), [Docker Machine](/machine/overview.md), and [Kitematic](/kitematic/userguide.md).

Containers and images created with Docker Desktop are shared between all
user accounts on machines where it is installed. This is because all Windows
accounts use the same VM to build and run containers.

Nested virtualization scenarios, such as running Docker Desktop on a
VMWare or Parallels instance might work, but there are no guarantees. For
more information, see [Running Docker Desktop in nested virtualization scenarios](troubleshoot.md#running-docker-desktop-for-windows-in-nested-virtualization-scenarios).

**Note**: Refer to the [Docker compatibility matrix](https://success.docker.com/article/compatibility-matrix) for complete Docker compatibility information with Windows Server.

### About Windows containers

Looking for information on using Windows containers?

* [Switch between Windows and Linux
  containers](/docker-for-windows/index.md#switch-between-windows-and-linux-containers)
  describes how you can toggle between Linux and Windows containers in Docker Desktop and points you to the tutorial mentioned above.
* [Getting Started with Windows Containers
  (Lab)](https://github.com/docker/labs/blob/master/windows/windows-containers/README.md)
  provides a tutorial on how to set up and run Windows containers on Windows 10, Windows Server 2016 and Windows Server 2019. It shows you how to use a MusicStore application
  with Windows containers.
* Docker Container Platform for Windows [articles and blog
  posts](https://www.docker.com/microsoft/) on the Docker website.

## Install Docker Desktop on Windows

1. Double-click **Docker Desktop Installer.exe** to run the installer.

    If you haven't already downloaded the installer (`Docker Desktop Installer.exe`), you can get it from
    [**Docker Hub**](https://hub.docker.com/?overlay=onboarding).
    It typically downloads to your `Downloads` folder, or you can run it from
    the recent downloads bar at the bottom of your web browser.

2. Follow the instructions on the installation wizard to accept the license, authorize the installer, and proceed with the install.

    When prompted, authorize the Docker Desktop Installer with your system password during the
    install process. Privileged access is needed to install networking
    components, links to the Docker apps, and manage the Hyper-V VMs.

3. Click **Finish** on the setup complete dialog and launch the Docker Desktop application.

## Start Docker Desktop

Docker Desktop does not start automatically after installation. To start Docker Desktop, search for Docker, and select **Docker Desktop** in the search results.

![search for Docker app](images/docker-app-search.png){:width="400px"}

When the whale icon in the status bar stays steady, Docker Desktop is up-and-running, and is accessible from any terminal window.

![whale on taskbar](images/whale-icon-systray.png)

If the whale icon is hidden in the Notifications area, click the up arrow on the
taskbar to show it. To learn more, see [Docker Settings](/docker-for-windows/index.md#docker-settings-dialog).

After installing the Docker Desktop app, you also get a pop-up success message with
suggested next steps, and a link to this documentation.

![Startup information](images/docker-app-welcome.png){:width="400px"}

When initialization is complete, click the whale icon in the Notifications area and select **About Docker Desktop** to verify that you have the latest version.

Congratulations! You are successfully running Docker Desktop on Windows.

## Uninstall Docker Desktop

To uninstall Docker Desktop from your Windows machine:

1. From the Windows **Start** menu, select **Settings** > **Apps** > **Apps & features**.
2. Select **Docker Desktop** from the **Apps & features** list and then select **Uninstall**.
3. Click **Uninstall** to confirm your selection.

> **Note:** Uninstalling Docker Desktop will destroy Docker containers and images local to the machine and remove the files generated by the application.

## Where to go next

* [Getting started](index.md) introduces Docker Desktop for Windows.
* [Get started with Docker](/get-started/) is a tutorial that teaches you how to
  deploy a multi-service stack.
* [Troubleshooting](troubleshoot.md) describes common problems, workarounds, and
  how to get support.
* [FAQs](faqs.md) provides answers to frequently asked questions.
* [Stable Release Notes](release-notes.md) or [Edge Release
  Notes](edge-release-notes.md).
