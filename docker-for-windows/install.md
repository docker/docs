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
href="https://github.com/docker/cli/blob/master/experimental/README.md">
experimental features in Docker Engine</a> enabled by default. Experimental mode can be toggled on and off in [preferences](/docker-for-windows/index.md#daemon-experimental-mode).

We welcome your
[feedback](/docker-for-windows/index.md#giving-feedback-and-getting-help) to help us improve Docker for Windows.

For more about Stable and Edge channels, see the
[FAQs](/docker-for-windows/faqs.md#questions-about-stable-and-edge-channels).

<table style="width:100%">
  <tr>
    <th style="font-size: x-large; font-family: arial">Stable channel</th>
    <th style="font-size: x-large; font-family: arial;">Edge channel</th>
  </tr>
  <tr valign="top">
    <td width="33%">Stable is the best channel to use if you want a reliable platform to work with. Stable releases track the Docker platform stable releases.<br><br>
   On this channel, you can select whether to send usage statistics and other data. <br><br>Stable releases happen once per quarter.
    </td>
    <td width="33%">Use the Edge channel if you want to get experimental features faster, and can weather some instability and bugs. We collect usage data on Edge releases. <br><br>Edge builds are released once per month.
    </td>
  </tr>
  <tr valign="top">
    <td width="33%">
      <a class="button outline-btn" href="https://download.docker.com/win/stable/Docker%20for%20Windows%20Installer.exe">Get Docker for Windows (Stable)</a>
    </td>
    <td width="33%">
      <a class="button outline-btn" href="https://download.docker.com/win/edge/Docker%20for%20Windows%20Installer.exe">Get Docker for Windows (Edge)</a>
    </td>
  </tr>
  <tr valign="top">
    <td width="33%"><a href="https://download.docker.com/win/stable/Docker%20for%20Windows%20Installer.exe.sha256sum"><font color="#BDBDBD" size="-1">Checksum: Stable installer SHA256</font></a>
  </td>
    <td width="33%"><a href="https://download.docker.com/win/edge/Docker%20for%20Windows%20Installer.exe.sha256sum"><font color="#BDBDBD" size="-1">Checksum: Edge installer SHA256</font></a>
    </td>
  </tr>
</table>

* Legacy (`.msi`) installers are available for [Edge](https://download.docker.com/win/edge/InstallDocker.msi) and [Stable](https://download.docker.com/win/stable/InstallDocker.msi) channels.

* The Docker for Windows is designed to configure Docker development environments on Windows 10 and on Windows Server 2016. You can develop both Docker Linux containers and Docker Windows containers with Docker for Windows. To run Docker Windows containers in production, see instructions for [installing Docker EE on Windows Server 2016](/engine/installation/windows/docker-ee.md). To run Docker Linux containers in production, see [instructions for installing Docker on Linux](/engine/installation/index.md).

* Docker for Windows requires 64bit Windows 10 Pro with Hyper-V available. Please see [What to know before you install](/docker-for-windows/install.md#what-to-know-before-you-install) for a full list
of prerequisites.

* You can switch between Edge and Stable versions, but you can only have one
   app installed at a time. Also, you will need to save images and export
   containers you want to keep before uninstalling the current version before
   installing another. For more about this, see the [FAQs about Stable and Edge
   channels](/docker-for-windows/faqs.md#questions-about-stable-and-edge-channels).

##  What to know before you install

If your system does not meet the requirements to run Docker for Windows, you can install
[Docker Toolbox](/toolbox/overview.md), which uses Oracle Virtual Box instead of
Hyper-V.

* **README FIRST for Docker Toolbox and Docker Machine users**: Docker for Windows requires Microsoft Hyper-V to run.  The Docker for Windows installer will enable Hyper-V for you, if needed, and restart your machine. After Hyper-V is
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
* The current version of Docker for Windows runs on 64bit Windows 10 Pro, Enterprise and Education (1607 Anniversary Update, Build 14393 or later). In the future we will support more versions of Windows 10.
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

1. Double-click **Docker for Windows Installer.exe** to run the installer.

    If you haven't already downloaded the installer (`Docker for Windows Installer.exe`), you can get it from
    [**download.docker.com**](https://download.docker.com/win/stable/Docker%20for%20Windows%20Installer.exe).
    It typically downloads to your `Downloads folder`, or you can run it from the recent downloads bar at the
    bottom of your web browser.

2. Follow the install wizard to accept the license, authorize the installer, and proceed with the install.

    You will be asked to authorize `Docker.app` with your system password during the install process.
    Privileged access is needed to install networking components, links to the Docker apps, and manage the
    Hyper-V VMs.

3. Click **Finish** on the setup complete dialog to launch Docker.

    ![Install complete>](/docker-for-windows/images/installer-finishes.png)

## Start Docker for Windows

Docker will not start automatically. To start it, search for Docker, select the
app in the search results, and click it (or hit Return).

![search for Docker app](/docker-for-windows/images/docker-app-search.png)

When the whale in the status bar stays steady, Docker is up-and-running, and
accessible from any terminal window.

![whale on taskbar](/docker-for-windows/images/whale-taskbar-circle.png)

If the whale is hidden in the Notifications area, click the up arrow on the
taskbar to show it. To learn more, see [Docker Settings](/docker-for-windows/index.md#docker-settings).

If you just installed the app, you also get a popup success message with
suggested next steps, and a link to this documentation.

![Startup information](/docker-for-windows/images/win-install-success-popup-cloud.png)

When initialization is complete, select **About Docker** from the notification
area icon to verify that you have the latest version.

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
