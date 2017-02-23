---
description: How to install Docker for Windows
keywords: windows, beta, alpha, install, download
title: Install Docker for Windows
---

The Docker for Windows install package includes everything you need to run
Docker on a Windows system. This topic describes pre-install considerations, and
how to download and install Docker for Windows.

> **Already have Docker for Windows?** If you already have Docker for Windows installed, and are ready to get started, skip to [Getting started](index.md) to
work through the rest of the Docker for Windows tour and information, or jump
over to the tutorials at [Learn Docker](/learn.md).

## Download Docker for Windows

If you have not already done so, please install Docker for Windows. You can
download installers from the stable or beta channel.  For more about stable and
beta channels, see the
[FAQs](/docker-for-windows/faqs.md#questions-about-stable-and-beta-channels).

<table style="width:100%">
  <tr>
    <th style="font-size: x-large; font-family: arial">Stable channel</th>
    <th style="font-size: x-large; font-family: arial">Beta channel</th>
  </tr>
  <tr valign="top">
    <td width="50%">This installer is fully baked and tested, and comes
    with the latest GA version of Docker Engine along with
    <a href="https://github.com/docker/docker/blob/master/experimental/README.md"> experimental features in Docker Engine</a>, which are enabled
    by default and configurable on
    <a href="index#daemon-experimental-mode">Docker Daemon settings for
    experimental mode</a>. <br><br>This is the best channel to use if
    you want a reliable platform to work with. (Be sure to disable
    experimental features for apps in production.) <br><br>These releases follow a version schedule with a longer lead time than the betas,
    synched with Docker Engine releases and hotfixes.<br><br>On the
    stable channel, you can select whether to send usage statistics
    and other data.
    </td>
    <td width="50%">This installer provides the latest Beta release of
    Docker for Windows, offers cutting edge features along with <a href="https://github.com/docker/docker/blob/master/experimental/README.md"> experimental features in Docker Engine</a>, which are enabled
    by default and configurable on <a href="index#daemon-experimental-mode">
    Docker Daemon settings for experimental mode</a>. <br><br>This is
    the best channel to use if you want to experiment with features
    under development, and can weather some instability and bugs. This
    channel is a continuation of the beta program, where you can
    provide feedback as the apps evolve. Releases are typically more
    frequent than for stable, often one or more per month. <br><br>
    We collect all usage data on betas across the board.</td>
  </tr>
  <tr valign="top">
  <td width="50%">
  <a class="button darkblue-btn" href="https://download.docker.com/win/stable/InstallDocker.msi">Get Docker for Windows (stable)</a><br><br>
  <a href="https://download.docker.com/win/stable/InstallDocker.msi.sha256sum"><font color="#BDBDBD" size="-1">Download checksum: InstallDocker.msi SHA256</font></a>
  </td>
  <td width="50%">
  <a class="button darkblue-btn" href="https://download.docker.com/win/beta/InstallDocker.msi">Get Docker for Windows (beta)</a><br><br>
  <a href="https://download.docker.com/win/beta/InstallDocker.msi.sha256sum"><font color="#BDBDBD" size="-1">Download checksum: InstallDocker.msi SHA256</font></a>
  </td>
  </tr>
</table>

>**Important Notes:**
>
> - Docker for Windows requires 64bit Windows 10 Pro, Enterprise and Education
>   (1511 November update, Build 10586 or later) and Microsoft Hyper-V. Please
>   see
>   [What to know before you install](/docker-for-windows/#what-to-know-before-you-install)
>   for a full list of prerequisites.
>
> - You can switch between beta and stable versions, but you must have only one
>   app installed at a time. Also, you will need to save images and export
>   containers you want to keep before uninstalling the current version before
>   installing another. For more about this, see the
>   [FAQs about beta and stable channels](/docker-for-windows/faqs.md#questions-about-stable-and-beta-channels).

##  What to know before you install

* **README FIRST for Docker Toolbox and Docker Machine users**: Docker for Windows requires Microsoft Hyper-V to run. After Hyper-V is enabled,
VirtualBox will no longer work, but any VirtualBox VM images will remain.
VirtualBox VMs created with `docker-machine` (including the `default` one
typically created during Toolbox install) will no longer start. These VMs cannot
be used side-by-side with Docker for Windows. However, you can still use
`docker-machine` to manage remote VMs.
<p />
* The current version of Docker for Windows runs on 64bit Windows 10 Pro, Enterprise and Education (1511 November update, Build 10586 or later). In the future we will support more versions of Windows 10.
<p />
* Containers and images created with Docker for Windows are shared between all user accounts on machines where it is installed. This is because all
Windows accounts will use the same VM to build and run containers. In the
future, Docker for Windows will better isolate user content.
<p />
* The Hyper-V package must be enabled for Docker for Windows to work. The Docker for Windows installer will enable it for you, if needed. (This requires a
reboot). If your system does not satisfy these requirements, you can install
[Docker Toolbox](/toolbox/overview.md), which uses Oracle Virtual Box instead of
Hyper-V.
<p />
* Virtualization must be enabled. Typically, virtualization is enabled by default. (Note that this is different from having Hyper-V enabled.) For more
detail see [Virtualization must be
enabled](troubleshoot.md#virtualization-must-be-enabled) in Troubleshooting.
<p />
* Nested virtualization scenarios, such as running Docker for Windows
on a VMWare or Parallels instance, might work, but come with no
guarantees (i.e., not officially supported). For more information, see
[Running Docker for Windows in nested virtualization scenarios](troubleshoot.md#running-docker-for-windows-in-nested-virtualization-scenarios)
<p />
* **What the Docker for Windows install includes**: The installation provides [Docker Engine](/engine/userguide/intro.md), Docker CLI client, [Docker Compose](/compose/overview.md), and [Docker Machine](/machine/overview.md).

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

* [Troubleshooting](troubleshoot.md) describes common problems,
workarounds, how to run and submit diagnostics, and submit issues.

* [FAQs](faqs.md) provides answers to frequently asked questions.

* [Release Notes](release-notes.md) lists component updates, new features, and improvements associated with Stable and Beta releases.

* [Learn Docker](/learn.md) provides general Docker tutorials.
