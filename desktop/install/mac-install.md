---
description: How to install Docker Desktop on Mac
keywords: mac, install, download, run, docker, local
title: Install Docker Desktop on Mac
redirect_from:
- /desktop/mac/install/
- /docker-for-mac/install/
- /engine/installation/mac/
- /installation/mac/
---

> **Docker Desktop terms**
>
> Commercial use of Docker Desktop in larger enterprises (more than 250
> employees OR more than $10 million USD in annual revenue) requires a paid
> subscription.

This page contains information about system requirements, download URLs, and instructions on how to install Docker Desktop for Mac.

[Mac with Intel chip](https://desktop.docker.com/mac/main/amd64/Docker.dmg?utm_source=docker&utm_medium=webreferral&utm_campaign=docs-driven-download-mac-amd64){: .button .primary-btn }
[Mac with Apple chip](https://desktop.docker.com/mac/main/arm64/Docker.dmg?utm_source=docker&utm_medium=webreferral&utm_campaign=docs-driven-download-mac-arm64){: .button .primary-btn }

*For checksums, see [Release notes](../release-notes.md)*

## System requirements

Your Mac must meet the following requirements to install Docker Desktop successfully.

<ul class="nav nav-tabs">
<li class="active"><a data-toggle="tab" data-target="#mac-intel">Mac with Intel chip</a></li>
<li><a data-toggle="tab" data-target="#mac-apple-silicon">Mac with Apple silicon</a></li>
</ul>
<div class="tab-content">
<div id="mac-intel" class="tab-pane fade in active" markdown="1">

### Mac with Intel chip

- macOS must be version 10.15 or newer. That is, Catalina, Big Sur, or Monterey. We recommend upgrading to the latest version of macOS.

  > **Note**
  >
  > Docker supports Docker Desktop on the most recent versions of macOS. That is, the current release of macOS and the previous two releases. As new major versions of macOS are made generally available, Docker stops supporting the oldest version and supports the newest version of macOS (in addition to the previous two releases). Docker Desktop currently supports macOS Catalina, macOS Big Sur, and macOS Monterey.

- At least 4 GB of RAM.

- VirtualBox prior to version 4.3.30 must not be installed as it is not compatible with Docker Desktop.

</div>
<div id="mac-apple-silicon" class="tab-pane fade" markdown="1">

### Mac with Apple silicon

- Beginning with Docker Desktop 4.3.0, we have removed the hard requirement to install **Rosetta 2**. There are a few optional command line tools that still require Rosetta 2 when using Darwin/AMD64. See the [Known issues section](../mac/apple-silicon.md#known-issues). However, to get the best experience, we recommend that you install Rosetta 2. To install Rosetta 2 manually from the command line, run the following command:

  ```console
  $ softwareupdate --install-rosetta
  ```

 For more information, see [Docker Desktop for Apple silicon](../mac/apple-silicon.md).

</div>
</div>

## Install and run Docker Desktop on Mac

### Install interactively

1. Double-click `Docker.dmg` to open the installer, then drag the Docker icon to
    the Applications folder.


2. Double-click `Docker.app` in the **Applications** folder to start Docker.

3. The Docker menu (![whale menu](images/whale-x.svg){: .inline}) displays the Docker Subscription Service Agreement window.

    {% include desktop-license-update.md %}

4. Select **Accept** to continue. Docker Desktop starts after you accept the terms.

    > **Important**
    >
    > If you do not agree to the terms, the Docker Desktop application will close and  you can no longer run Docker Desktop on your machine. You can choose to accept the terms at a later date by opening Docker Desktop.
    {: .important}

    For more information, see [Docker Desktop Subscription Service Agreement](https://www.docker.com/legal/docker-subscription-service-agreement). We recommend that you also read the [FAQs](https://www.docker.com/pricing/faq){: target="_blank" rel="noopener" class="_" id="dkr_docs_desktop_install_btl"}.

### Install from the command line

After downloading `Docker.dmg`, run the following commands in a terminal to install Docker Desktop in the Applications folder:

```console
$ sudo hdiutil attach Docker.dmg
$ sudo /Volumes/Docker/Docker.app/Contents/MacOS/install
$ sudo hdiutil detach /Volumes/Docker
```

As macOS typically performs security checks the first time an application is used, the `install` command can take several minutes to run.

The `install` command accepts the following flags:
- `--accept-license`: accepts the [Docker Subscription Service Agreement](https://www.docker.com/legal/docker-subscription-service-agreement){: target="_blank" rel="noopener" class="_"} now, rather than requiring it to be accepted when the application is first run
- `--allowed-org=<org name>`: requires the user to sign in and be part of the specified Docker Hub organization when running the application
- `--user=<username>`: Runs the privileged helper service once during installation, then disables it at runtime. This removes the need for the user to grant root privileges on first run. For more information, see [Privileged helper permission requirements](../mac/permission-requirements.md#permission-requirements){: target="_blank" rel="noopener" class="_"}. To find the username, enter `ls /Users` in the CLI.
- `--admin-settings`: Automatically creates an `admin-settings.json` file which is used by admins to control certain Docker Desktop settings on client machines within their organization. For more information, see [Settings Management](../hardened-desktop/settings-management/index.md).
  - It must be used together with the `--allowed-org=<org name>` flag. 
  - For example:
    `--allowed-org=<org name> --admin-settings='{"configurationFileVersion": 2, "enhancedContainerIsolation": {"value": true, "locked": false}}'`

## Where to go next

- [Docker Desktop for Apple silicon](../mac/apple-silicon.md) for detailed information about Docker Desktop for Apple silicon.
- [Troubleshooting](../troubleshoot/overview.md) describes common problems, workarounds, how
  to run and submit diagnostics, and submit issues.
- [FAQs](../faqs/general.md) provide answers to frequently asked questions.
- [Release notes](../release-notes.md) lists component updates, new features, and improvements associated with Docker Desktop releases.
- [Get started with Docker](../../get-started/index.md) provides a general Docker tutorial.
* [Back up and restore data](../backup-and-restore.md) provides instructions
  on backing up and restoring data related to Docker.
