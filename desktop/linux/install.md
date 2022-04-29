---
description: How to install Docker Desktop on Linux
keywords: linux, install, download, run, docker, local
title: Install Docker Desktop on Linux
redirect_from:
- /docker-for-linux/install/
---

> **Update to the Docker Desktop terms**
>
> Commercial use of Docker Desktop in larger enterprises (more than 250
> employees OR more than $10 million USD in annual revenue) now requires a paid
> subscription. The grace period for those that will require a paid subscription
> ends on January 31, 2022. [Learn
> more](https://www.docker.com/blog/the-grace-period-for-the-docker-subscription-service-agreement-ends-soon-heres-what-you-need-to-know/){:
 target="_blank" rel="noopener" class="_" id="dkr_docs_cta"}.
{: .important}

Welcome to Docker Desktop for Linux. This page contains information about Docker Desktop for Linux system requirements, download URLs, instructions to install and update Docker Desktop for Linux.

> Download Docker Desktop for Linux
>
> [DEB](https://desktop-stage.docker.com/linux/main/amd64/78459/docker-desktop-4.8.0-amd64.deb?utm_source=docker&utm_medium=webreferral&utm_campaign=docs-driven-download-linux-amd64){: .button .primary-btn }
> [RPM](https://desktop-stage.docker.com/linux/main/amd64/78459/docker-desktop-4.8.0-x86_64.rpm?utm_source=docker&utm_medium=webreferral&utm_campaign=docs-driven-download-linux-amd64){: .button .primary-btn }
> [Arch](https://desktop-stage.docker.com/linux/main/amd64/78459/docker-desktop-4.8.0-x86_64.pkg.tar.zst?utm_source=docker&utm_medium=webreferral&utm_campaign=docs-driven-download-linux-amd64){: .button .primary-btn }


## System requirements

Your Linux host must meet the following requirements to install Docker Desktop successfully.

- **QEMU must be version 5.2 or newer**. We recommend upgrading to the latest version.

  > **Note**
  >
  > Docker supports Docker Desktop on the most recent versions of Ubuntu. That is, the current and the previous release of Ubuntu. As new major versions of are made generally available, Docker stops supporting the oldest version and supports the newest version of Ubuntu (in addition to the previous two releases). Docker Desktop currently supports .

- At least 4 GB of RAM.


## Install and run Docker Desktop on ...
TODO: installation guides for several linux distros

### Install interactively

1. Download the package coresponding to your Linux distro and install it using the distro's package manager.

      ![Install Docker app](images/docker-app-drag.png)

2. Open Applications Menu and double-click `Docker Desktop` to start Docker.

    ![Docker app in Applications](images/docker-app-in-apps.png)

3. The Docker menu (![whale menu](images/whale-x.png){: .inline}) displays the Docker Subscription Service Agreement window. It includes a change to the terms of use for Docker Desktop.

    {% include desktop-license-update.md %}

4. Click the checkbox to indicate that you accept the updated terms and then click **Accept** to continue. Docker Desktop starts after you accept the terms.

    > **Important**
    >
    > If you do not agree to the terms, the Docker Desktop application will close and  you can no longer run Docker Desktop on your machine. You can choose to accept the terms at a later date by opening Docker Desktop.
    {: .important}

    For more information, see [Docker Desktop License Agreement](../../subscription/index.md#docker-desktop-license-agreement). We recommend that you also read the [Blog](https://www.docker.com/blog/updating-product-subscriptions/){: target="_blank" rel="noopener" class="_" id="dkr_docs_desktop_install_btl"} and [FAQs](https://www.docker.com/pricing/faq){: target="_blank" rel="noopener" class="_" id="dkr_docs_desktop_install_btl"} to learn how companies using Docker Desktop may be affected.

### Quick start guide  
  
  If you've just installed the app, Docker Desktop launches the Quick Start Guide. The tutorial includes a simple exercise to build an example Docker image, run it as a container, push and save the image to Docker Hub.

   ![Docker Quick Start tutorial](images/docker-tutorial-linux.png)

Congratulations! You are now successfully running Docker Desktop. Click the Docker menu (![whale menu](images/whale-x.png){: .inline}) to see
**Preferences** and other options. To run the Quick Start Guide on demand, select the Docker menu and then choose **Quick Start Guide**.

## Updates

{% include desktop-update.md %}

## Uninstall Docker Desktop

To uninstall Docker Desktop from your Linux host:

TODO: exemplify for several distros and app managers

> **Note**
>
> Uninstalling Docker Desktop destroys Docker containers, images, volumes, and
> other Docker related data local to the machine, and removes the files generated
> by the application. Refer to the [back up and restore data](../backup-and-restore.md)
> section to learn how to preserve important data before uninstalling.

## Where to go next

- [Getting started](index.md) provides an overview of Docker Desktop on Linux, basic Docker command examples, how to get help or give feedback, and links to other topics about Docker Desktop on Linux.
- [Troubleshooting](troubleshoot.md) describes common problems, workarounds, how
  to run and submit diagnostics, and submit issues.
- [FAQs](../faqs.md) provide answers to frequently asked questions.
- [Release notes](release-notes/index.md) lists component updates, new features, and improvements associated with Docker Desktop releases.
- [Get started with Docker](../../get-started/index.md) provides a general Docker tutorial.
* [Back up and restore data](../backup-and-restore.md) provides instructions
  on backing up and restoring data related to Docker.
