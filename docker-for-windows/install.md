---
description: How to install Docker for Windows
keywords: windows, beta, edge, alpha, install, download
title: Install Docker for Windows
---

Docker for Windows is a [Docker Community Edition
(CE)](https://www.docker.com/community-edition) app. The Docker for Windows
install package includes everything you need to run Docker on a Windows system.
This topic describes pre-install considerations, and how to download and install
Docker for Windows.<br><br>

> **Already have Docker for Windows?**
> If you already have Docker for
Windows installed, and are ready to get started, skip to
[Get started with Docker for Windows](index.md) for a quick tour of
the command line, settings, and tools.

>**Looking for Release Notes?** [Get release notes for all
versions here](release-notes.md).

## Download Docker for Windows

If you have not already done so, please install Docker for Windows. You can
download installers from the **Stable** or **Edge** channel.

Both Stable and Edge installers come with <a
href="https://github.com/moby/moby/blob/master/experimental/README.md">
experimental features in Docker Engine</a> enabled by default and configurable
on [Docker Daemon preferences](index.md#daemon-experimental-mode) for
experimental mode. We recommend that you disable experimental features for
apps in production.

On both channels, we welcome your
[feedback](index.md#giving-feedback-and-getting-help) to help us as the apps
evolve.

For more about Stable and Edge channels, see the
[FAQs](/docker-for-windows/faqs.md#questions-about-stable-and-edge-channels).

<table style="width:100%">
  <tr>
    <th style="font-size: x-large; font-family: arial">Stable channel</th>
    <th style="font-size: x-large; font-family: arial;">Edge channel</th>
    <th style="font-size: x-large; font-family: arial;">Legacy Edge installer</th>
  </tr>
  <tr valign="top">
    <td width="33%">This installer is fully baked and tested. This is the
    best channel to use if you want a reliable platform to work with. These releases follow the Docker Engine stable releases.<br><br>
   On this channel, you can select whether to send usage
   statistics and other data. <br><br>Stable builds are released once per quarter.
    </td>
    <td width="33%">This new installer includes experimental support for Windows Server 2016 as a part of the latest Edge release of
    Docker for Windows and Engine. <br><br>Use this channel if you want to get experimental features faster, and can weather some instability and bugs. We collect all usage data on Edge releases across the board. <br><br>Edge builds are released once per month.
    </td>
    <td width="33%">We recommend that all Edge users try the new installer. <br><br>However, if you have problems with the new installer experience, you can use the legacy Edge installer and still get all other Edge features.
    </td>
  </tr>
  <tr valign="top">
  <td width="33%">
  <a class="button outline-btn" href="https://download.docker.com/win/stable/InstallDocker.msi">Get Docker for Windows (Stable)</a>
  </td>
  <td width="33%">
  <a class="button outline-btn" href="https://download.docker.com/win/edge/Docker%20for%20Windows%20Installer.exe">Get Docker for Windows (Edge)</a>
  </td>
  <td width="33%">
  <a class="button outline-btn" href="https://download.docker.com/win/edge/InstallDocker.msi">Get Docker for Windows (Edge) (legacy installer)</a>
  </td>
  </tr>
  <tr valign="top">
  <td width="33%"><a href="https://download.docker.com/win/stable/InstallDocker.msi.sha256sum"><font color="#BDBDBD" size="-1">Checksum: Stable InstallDocker.msi SHA256</font></a>
  </td>
  <td width="33%"><a href="https://download.docker.com/win/edge/Docker%20for%20Windows%20Installer.exe.sha256sum"><font color="#BDBDBD" size="-1">Checksum: New Edge InstallDocker.exe SHA256</font></a>
  </td>
  <td width="33%"><a href="https://download.docker.com/win/edge/InstallDocker.msi.sha256sum"><font color="#BDBDBD" size="-1">Checksum: Legacy Edge InstallDocker.msi SHA256</font></a>
  </td>
  </tr>
</table>

* Docker for Windows requires 64bit Windows 10 Pro and Microsoft Hyper-V. Please see [What to know before you install](/docker-for-windows/#what-to-know-before-you-install) for a full list
of prerequisites.

* You can switch between Edge and Stable versions, but you must have only one
   app installed at a time. Also, you will need to save images and export
   containers you want to keep before uninstalling the current version before
   installing another. For more about this, see the [FAQs about Stable and Edge
   channels](/docker-for-windows/faqs.md#questions-about-stable-and-edge-channels).

##  What to know before you install

If your system does not satisfy these requirements, you can install
[Docker Toolbox](/toolbox/overview.md), which uses Oracle Virtual Box instead of
Hyper-V.

* **README FIRST for Docker Toolbox and Docker Machine users**: Docker for Windows requires Microsoft Hyper-V to run.  The Docker for Windows installer
will enable it for you, if needed, which requires a reboot. After Hyper-V is
enabled, VirtualBox will no longer work, but any VirtualBox VM images will
remain. VirtualBox VMs created with `docker-machine` (including the `default`
one typically created during Toolbox install) will no longer start. These VMs
cannot be used side-by-side with Docker for Windows. However, you can still use
`docker-machine` to manage remote VMs.
<p />
* Virtualization must be enabled. Typically, virtualization is enabled by default. (Note that this is different from having Hyper-V enabled.) For more
detail see [Virtualization must be
enabled](troubleshoot.md#virtualization-must-be-enabled) in Troubleshooting.
<p />
* The current version of Docker for Windows runs on 64bit Windows 10 Pro, Enterprise and Education (1511 November update, Build 10586 or later). In the future we will support more versions of Windows 10.
<p />
* Containers and images created with Docker for Windows are shared between all user accounts on machines where it is installed. This is because all
Windows accounts will use the same VM to build and run containers. In the
future, Docker for Windows will better isolate user content.
<p />
* Nested virtualization scenarios, such as running Docker for Windows
on a VMWare or Parallels instance, might work, but come with no
guarantees (i.e., not officially supported). For more information, see
[Running Docker for Windows in nested virtualization scenarios](troubleshoot.md#running-docker-for-windows-in-nested-virtualization-scenarios)
<p />
* **What the Docker for Windows install includes**: The installation provides [Docker Engine](/engine/userguide/), Docker CLI client, [Docker Compose](/compose/overview.md), [Docker Machine](/machine/overview.md), and [Kitematic](/kitematic/userguide.md).

### About Windows containers and Windows Server 2016

Looking for information on using Windows containers?

* [Switch between Windows and Linux containers](/docker-for-windows/index.md#switch-between-windows-and-linux-containers) describes the Linux / Windows containers toggle in Docker for Windows and points you to the tutorial mentioned above.
<p />
* [Getting Started with Windows Containers (Lab)](https://github.com/docker/labs/blob/master/windows/windows-containers/README.md)
provides a tutorial on how to set up and run Windows containers on Windows 10 or
with Windows Server 2016. It shows you how to use a MusicStore application with
Windows containers.
<p />
* [Setup - Windows Server 2016 (Lab)](https://github.com/docker/labs/blob/master/windows/windows-containers/Setup-Server2016.md) specifically describes environment setup.
<p />
* Docker Container Platform for Windows Server 2016 [articles and blog posts](https://www.docker.com/microsoft/) on the Docker website

## Install Docker for Windows

1. Double-click `InstallDocker.msi` to run the installer.

    If you haven't already downloaded the installer (`InstallDocker.msi`), you can get it [**here**](https://download.docker.com/win/stable/InstallDocker.msi). It typically downloads to your `Downloads folder`, or you can run it from the recent downloads bar at the bottom of your web browser.

2. Follow the install wizard to accept the license, authorize the installer, and proceed with the install.

    You will be asked to authorize `Docker.app` with your system password during the install process. Privileged access is needed to install networking components, links to the Docker apps, and manage the Hyper-V VMs.

3. Click **Finish** on the setup complete dialog to launch Docker.

    ![Install complete>](/docker-for-windows/images/installer-finishes.png)

## Start Docker for Windows

When the installation finishes, Docker starts automatically.

The whale in the status bar indicates that Docker is running, and accessible from a terminal.

If you just installed the app, you also get a popup success message with suggested next steps, and a link to this documentation.

![Startup information](/docker-for-windows/images/win-install-success-popup-cloud.png)

When initialization is complete, select **About Docker** from the notification area icon to verify that you have the latest version.

Congratulations! You are up and running with Docker for Windows.

## Where to go next

* [Getting started](index.md) provides an overview of Docker for Windows,
basic Docker command examples, how to get help or give feedback, and
links to all topics in the Docker for Windows guide.

* [Get started with Docker](/get-started/) teaches you how to define and deploy
applications with Docker.

* [Troubleshooting](troubleshoot.md) describes common problems,
workarounds, how to run and submit diagnostics, and submit issues.

* [FAQs](faqs.md) provides answers to frequently asked questions.

* [Release Notes](release-notes.md) lists component updates, new features, and improvements associated with Stable and Edge releases.
