---
description: How to install Docker for Mac
keywords: mac, beta, alpha, install, download
title: Install Docker for Mac
---

Docker for Mac is a [Docker Community Edition (CE)](https://www.docker.com/community-edition) app. The Docker for Mac
install package includes everything you need to run Docker on a Mac. This topic
describes pre-install considerations, and how to download and install Docker for
Mac.<br><br>

> **Already have Docker for Mac?** If you already have
Docker for Mac installed, and are ready to get started, skip to
[Get started with Docker for Mac](index.md) for a quick tour of
the command line, preferences, and tools.

>**Looking for Release Notes?** [Get release notes for all versions here](release-notes.md).

## Download Docker for Mac

If you have not already done so, please install Docker for Mac. You can download
installers from the Stable or beta channel.

Both Stable and Edge installers come with <a
href="https://github.com/moby/moby/blob/master/experimental/README.md">
experimental features in Docker Engine</a> enabled by default and configurable
on [Docker Daemon preferences](index.md#daemon-experimental-mode) for
experimental mode. We recommend that you disable experimental features for
apps in production.

On both channels, we welcome your <a
href="troubleshoot/#diagnose-problems-send-feedback-and-create-github-issues">feedback</a>
to help us as the apps evolve.

For more about Stable and Edge channels, see the
[FAQs](/docker-for-mac/faqs.md#stable-and-edge-channels).

<table style="width:100%">
  <tr>
    <th style="font-size: x-large; font-family: arial">Stable channel</th>
    <th style="font-size: x-large; font-family: arial">Edge channel</th>
  </tr>
  <tr valign="top">
    <td width="50%">This installer is fully baked and tested. This is the
    best channel to use if you want a reliable platform to work with. These releases follow the Docker Engine stable releases.<br><br>
   On this channel, you can select whether to send usage
   statistics and other data. <br><br>Stable builds are released once per quarter.
   </td>
    <td width="50%">This installer provides the latest Edge release of
    Docker for Mac and Engine, and typically offers new features in development. Use this channel if you want to get experimental features faster, and can weather some instability and bugs. We collect all usage data on Edge releases across the board. <br><br>Edge builds are released once per month.
    </td>
  </tr>
  <tr valign="top">
  <td width="50%">
  <a class="button outline-btn" href="https://download.docker.com/mac/stable/Docker.dmg">Get Docker for Mac (Stable)</a>
  </td>
  <td width="50%">
  <a class="button outline-btn" href="https://download.docker.com/mac/edge/Docker.dmg">Get Docker for Mac (Edge)</a>
  </td>
  </tr>
  <tr valign="top">
  <td width="50%">
  <a href="https://download.docker.com/mac/stable/Docker.dmg.sha256sum"><font color="#BDBDBD" size="-1">Checksum: Docker.dmg SHA256</font></a>
  </td>
  <td width="50%">
  <a href="https://download.docker.com/mac/edge/Docker.dmg.sha256sum"><font color="#BDBDBD" size="-1">Checksum: Docker.dmg SHA256</font></a>
  </td>
  </tr>
</table>

* Docker for Mac requires OS X El Capitan 10.11 or newer macOS release running on a 2010 or newer Mac, with Intel's  hardware support for MMU virtualization.
The app will run on 10.10.3 Yosemite, but with limited support. Please see [What
to know before you install](#what-to-know-before-you-install) for a full
explanation and list of prerequisites.

* You can switch between Edge and Stable versions, but you must have only one
   app installed at a time. Also, you will need to save images and export
   containers you want to keep before uninstalling the current version before
   installing another. For more about this, see the
   [FAQs about Stable and Edge channels](faqs.md#stable-and-edge-channels).

##  What to know before you install

* **README FIRST for Docker Toolbox and Docker Machine users**: If you are
  already running Docker on your machine, first read
  [Docker for Mac vs. Docker Toolbox](docker-toolbox.md) to understand the
  impact of this installation on your existing setup, how to set your environment
  for Docker for Mac, and how the two products can coexist.
<p />
* **Relationship to Docker Machine**: Installing Docker for Mac does not affect
  machines you created with Docker Machine. You'll get the option to copy
  containers and images from your local `default` machine (if one exists) to the
  new Docker for Mac [HyperKit](https://github.com/docker/HyperKit/) VM. When
  you are running Docker for Mac, you do not need Docker Machine nodes running
  at all locally (or anywhere else). With Docker for Mac, you have a new, native
  virtualization system running (HyperKit) which takes the place of the
  VirtualBox system. To learn more, see
  [Docker for Mac vs. Docker Toolbox](docker-toolbox.md).
<p />
* **System Requirements**: Docker for Mac will launch only if all of these
  requirements are met.
  <p />
  - Mac must be a 2010 or newer model, with Intel's hardware support for memory
    management unit (MMU) virtualization; i.e., Extended Page Tables (EPT) and
    Unrestricted Mode.
  <p />
  - OS X El Capitan 10.11 and newer macOS releases are supported. At a minimum,
    Docker for Mac requires macOS Yosemite 10.10.3 or newer, with the caveat
    that going forward 10.10.x is a use-at-your-own risk proposition.
  <p />
  - Starting with Docker for Mac Stable release 1.13 (upcoming), and concurrent
    Edge releases, we will no longer address issues specific to OS X Yosemite
    10.10. In future releases, Docker for Mac could stop working on OS X Yosemite
    10.10 due to the deprecated status of this OS X version. We recommend
    upgrading to the latest version of macOS.
  <p />
  - At least 4GB of RAM
  <p />
  - VirtualBox prior to version 4.3.30 must NOT be installed (it is incompatible
    with Docker for Mac). If you have a newer version of VirtualBox installed, it's fine.

  > **Note**: If your system does not satisfy these requirements, you can
  > install [Docker Toolbox](/toolbox/overview.md), which uses Oracle VirtualBox
  > instead of HyperKit.

* **What the install includes**: The installation provides
  [Docker Engine](/engine/userguide/), Docker CLI client,
  [Docker Compose](/compose/overview/), [Docker Machine](/machine/overview/), and [Kitematic](/kitematic/userguide.md).

## Install and Run Docker for Mac

1.  Double-click `Docker.dmg` to open the installer, then drag Moby the whale to
    the Applications folder.

	  ![Install Docker app](/docker-for-mac/images/docker-app-drag.png)

2.  Double-click `Docker.app` in the Applications folder to start Docker.  (In the example below, the Applications folder is in "grid" view mode.)

	  ![Docker app in Hockeyapp](/docker-for-mac/images/docker-app-in-apps.png)

	  You will be asked to authorize `Docker.app` with your system password after you launch it.
	  Privileged access is needed to install networking components and links to the Docker apps.

	  The whale in the top status bar indicates that Docker is running, and accessible from a terminal.

	  ![Whale in menu bar](/docker-for-mac/images/whale-in-menu-bar.png)

	  If you just installed the app, you also get a success message with suggested
    next steps and a link to this documentation. Click the whale (![whale](/docker-for-mac/images/whale-x.png))
    in the status bar to dismiss this popup.

	  ![Startup information](/docker-for-mac/images/mac-install-success-docker-cloud.png)

3.  Click the whale (![whale-x](images/whale-x.png)) to get Preferences and
    other options.

	  ![Docker context menu](images/menu.png)

4.  Select **About Docker** to verify that you have the latest version.

Congratulations! You are up and running with Docker for Mac.

## Where to go next

* [Getting started](index.md) provides an overview of Docker for Mac,
basic Docker command examples, how to get help or give feedback, and
links to all topics in the Docker for Mac guide.

* [Troubleshooting](troubleshoot.md) describes common problems,
workarounds, how to run and submit diagnostics, and submit issues.

* [FAQs](faqs.md) provides answers to frequently asked questions.

* [Release Notes](release-notes.md) lists component updates, new features, and improvements associated with Stable and Edge releases.

* [Get Started with Docker](/get-started/) provides a general Docker tutorial.
