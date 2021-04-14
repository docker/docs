---
description: How to install Docker Desktop on Mac
keywords: mac, install, download, run, docker, local
title: Install Docker Desktop on Mac
---

Welcome to Docker Desktop for Mac. This page contains information about Docker Desktop for Mac system requirements, download URLs, installation instructions, and automatic updates.

Click one of the following buttons to download Docker Desktop for Mac:

[Mac with Intel chip](https://desktop.docker.com/mac/stable/amd64/Docker.dmg){: .button .primary-btn}

[Mac with Apple chip](https://desktop.docker.com/mac/stable/arm64/Docker.dmg){: .button .primary-btn}

By downloading Docker Desktop, you agree to the terms of the [Docker Software End User License Agreement](https://www.docker.com/legal/docker-software-end-user-license-agreement){: target="_blank" rel="noopener" class="_"} and the [Docker Data Processing Agreement](https://www.docker.com/legal/data-processing-agreement){: target="_blank" rel="noopener" class="_"}.

## System requirements

Your Mac must meet the following requirements to successfully install Docker Desktop.

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#mac-intel">Mac with Intel chip</a></li>
<li><a data-toggle="tab" data-target="#mac-apple-silicon">Mac with Apple silicon</a></li>
</ul>
<div class="tab-content">
<div id="mac-intel" class="tab-pane fade in active" markdown="1">

### Mac with Intel chip

- **macOS must be version 10.14 or newer**. That is, Mojave, Catalina, or Big Sur. We recommend upgrading to the latest version of macOS.

  If you experience any issues after upgrading your macOS to version 10.15, you must install the latest version of Docker Desktop to be compatible with this version of macOS.

  > **Note**
  >
  > Docker supports Docker Desktop on the most recent versions of macOS. That is, the current release of macOS and the previous two releases. As new major versions of macOS are made generally available, Docker stops supporting the oldest version and supports the newest version of macOS (in addition to the previous two releases). Docker Desktop currently supports macOS Mojave, macOS Catalina, and macOS Big Sur.

- At least 4 GB of RAM.

- VirtualBox prior to version 4.3.30 must not be installed as it is not compatible with Docker Desktop.

</div>
<div id="mac-apple-silicon" class="tab-pane fade" markdown="1">

### Mac with Apple silicon

- You must install **Rosetta 2** as some binaries are still Darwin/AMD64. To install Rosetta 2 manually from the command line, run the following command:

  ```bash
    softwareupdate --install-rosetta
  ```

 For more information, see [Docker Desktop for Apple silicon](apple-silicon.md).

</div>
</div>

## What's included in the installer

The Docker Desktop installation includes
  [Docker Engine](../engine/index.md), Docker CLI client,
  [Docker Compose](../compose/index.md), [Notary](../notary/getting_started.md), [Kubernetes](https://github.com/kubernetes/kubernetes/), and [Credential Helper](https://github.com/docker/docker-credential-helpers/).

## Install and run Docker Desktop on Mac

1. Double-click `Docker.dmg` to open the installer, then drag the Docker icon to
    the Applications folder.

      ![Install Docker app](images/docker-app-drag.png)

2. Double-click `Docker.app` in the Applications folder to start Docker. (In the example below, the Applications folder is in "grid" view mode.)

    ![Docker app in Hockeyapp](images/docker-app-in-apps.png)

    The Docker menu in the top status bar indicates that Docker Desktop is running, and accessible from a terminal.

      ![Whale in menu bar](images/whale-in-menu-bar.png)

    If you've just installed the app, Docker Desktop launches the onboarding tutorial. The tutorial includes a simple exercise to build an example Docker image, run it as a container, push and save the image to Docker Hub.

    ![Docker Quick Start tutorial](images/docker-tutorial-mac.png)

3. Click the Docker menu (![whale menu](images/whale-x.png){: .inline}) to see
**Preferences** and other options.

4. Select **About Docker** to verify that you have the latest version.

Congratulations! You are now successfully running Docker Desktop.

If you would like to rerun the tutorial, go to the Docker Desktop menu 
and select **Learn**.

## Automatic updates

Starting with Docker Desktop 3.0.0, updates to Docker Desktop will be available automatically as delta updates from the previous version.

When an update is available, Docker Desktop displays an icon to indicate the availability of a newer version. You can start downloading the update in the background whenever it is convenient for you. Alternatively, you can click **Snooze** to pause the update.

After downloading the update, all you need to do is to click **Update and restart** from the Docker menu. This installs the latest update and restarts Docker Desktop for the changes to take effect.

## Uninstall Docker Desktop

To uninstall Docker Desktop from your Mac:

1. From the Docker menu, select **Troubleshoot** and then select **Uninstall**.
2. Click **Uninstall** to confirm your selection.

> **Important**
> 
> Uninstalling Docker Desktop destroys Docker containers, images, volumes, and
> other Docker related data local to the machine, and removes the files generated
> by the application. Refer to the [back up and restore data](../desktop/backup-and-restore.md)
> section to learn how to preserve important data before uninstalling.

## Where to go next

- [Getting started](index.md) provides an overview of Docker Desktop on Mac, basic Docker command examples, how to get help or give feedback, and links to other topics about Docker Desktop on Mac.
- [Docker Desktop for Apple silicon](apple-silicon.md) for detailed information about Docker Desktop for Apple silicon.
- [Troubleshooting](troubleshoot.md) describes common problems, workarounds, how
  to run and submit diagnostics, and submit issues.
- [FAQs](../desktop/faqs.md) provide answers to frequently asked questions.
- [Release notes](release-notes.md) lists component updates, new features, and improvements associated with Docker Desktop releases.
- [Get started with Docker](../get-started/index.md) provides a general Docker tutorial.
* [Back up and restore data](../desktop/backup-and-restore.md) provides instructions
  on backing up and restoring data related to Docker.
