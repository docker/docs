---
description: How to install Docker Desktop on Linux
keywords: linux, install, download, run, docker, local
title: Install Docker Desktop on Linux
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

> Download Docker Desktop for Linux packages
>
> [DEB](https://desktop.docker.com/linux/main/amd64/docker-desktop-4.8.0-amd64.deb?utm_source=docker&utm_medium=webreferral&utm_campaign=docs-driven-download-linux-amd64){: .button .primary-btn }
> [RPM](https://desktop.docker.com/linux/main/amd64/docker-desktop-4.8.0-x86_64.rpm?utm_source=docker&utm_medium=webreferral&utm_campaign=docs-driven-download-linux-amd64){: .button .primary-btn }


## System requirements

Your Linux host must meet the following requirements to install Docker Desktop successfully:

- 64-bit kernel and CPU support for virtualization 

- KVM virtualization support. Follow the [KVM virtualization support instructions](#kvm-virtualization-support) to check if the KVM kernel modules are enabled and how to provide access to the kvm device.

- **QEMU must be version 5.2 or newer**. We recommend upgrading to the latest version. 

- systemd init system

- Gnome or KDE Desktop environment

- At least 4 GB of RAM.

Docker Desktop for Linux runs a Virtual Machine (VM). For more information on why, see [here](overview.md#why-docker-desktop-for-linux-runs-a-vm).

> **Note:**
>
> Docker does not provide support for running Docker Desktop in nested virtualization scenarios. We recommend that you run Docker Desktop for Linux natively on the supported distributions.


## Supported platforms

Docker provides `.deb` and `.rpm` packages from the following Linux distributions
and architectures:


{% assign yes = '![yes](/images/green-check.svg){: .inline style="height: 14px; margin: 0 auto"}' %}


| Platform                | x86_64 / amd64         | 
|:------------------------|:-----------------------|
| [Ubuntu](install/ubuntu.md)     | [{{ yes }}](install/ubuntu.md) |
| [Debian](install/debian.md)     | [{{ yes }}](install/debian.md) |
| [Fedora](install/fedora.md)     | [{{ yes }}](install/fedora.md) |


>  **Note:**
>
> An experimental package is available for [Arch](install/archlinux.md)-based distributions. Docker has not tested or verified the installation.

Docker supports Docker Desktop on the current LTS release of the aforementioned distributions and the most recent version. As new versions are made generally available, Docker stops supporting the oldest version and supports the newest version.


### KVM virtualization support


Docker Desktop runs a VM that requires [KVM support](https://www.linux-kvm.org). 

The `kvm` module should load automatically if the host has virtualization support. To load the module manually run:

```console
$ modprobe kvm
```

Depending on the processor of the host machine, the coresponding module must be loaded:

```console
$ modprobe kvm_intel  # Intel processors

$ modprobe kvm_amd    # AMD processors
```

To check the KVM modules have been enabled:

```console
$ lsmod | grep kvm
kvm_amd               167936  0
ccp                   126976  1 kvm_amd
kvm                  1089536  1 kvm_amd
irqbypass              16384  1 kvm
```

#### Set up KVM device user permissions


Check ownership of `/dev/kvm`, run :

```console
$ ls -al /dev/kvm
```

Add your user to the kvm group in order to access the kvm device.

```console
$ sudo usermod -aG kvm $USER
```

Log out and log back in so that your group membership is re-evaluated.


### Generic installation steps

1. Download the correct package for your Linux distribution and install it with the corresponding package manager.

2. Open your **Applications** menu in Gnome/KDE Desktop and search for **Docker Desktop**.

    ![Docker app in Applications](images/docker-app-in-apps.png)

3. Double-click **Docker Desktop** to start Docker.

4. The Docker menu (![whale menu](images/whale-x.png){: .inline}) displays the Docker Subscription Service Agreement window. It includes a change to the terms of use for Docker Desktop.

    {% include desktop-license-update.md %}

5. Click the checkbox to accept the updated terms and then click **Accept** to continue. Docker Desktop starts after you accept the terms.

    > **Important**
    >
    > If you do not agree to the terms, the Docker Desktop application will close and  you can no longer run Docker Desktop on your machine. You can choose to accept the terms at a later date by opening Docker Desktop.
    {: .important}

    For more information, see [Docker Desktop License Agreement](../../subscription/index.md#docker-desktop-license-agreement). We recommend that you also read the [Blog](https://www.docker.com/blog/updating-product-subscriptions/){: target="_blank" rel="noopener" class="_" id="dkr_docs_desktop_install_btl"} and [FAQs](https://www.docker.com/pricing/faq){: target="_blank" rel="noopener" class="_" id="dkr_docs_desktop_install_btl"} to learn how companies using Docker Desktop may be affected.

### Quick start guide  
  
  If you've just installed the app, Docker Desktop launches the Quick Start Guide. The tutorial includes a simple exercise to build an example Docker image, run it as a container, push and save the image to Docker Hub.

   ![Docker Quick Start tutorial](images/docker-tutorial-linux.png)

Congratulations! You are now successfully running Docker Desktop. Click the Docker menu (![whale menu](images/whale-x.png){: .inline}) to see
**Settings** and other options. To run the Quick Start Guide on demand, select the Docker menu and then choose **Quick Start Guide**.

## Updates

Once a new version for Docker Desktop is released, the Docker UI shows a notification. 
 You need to download the new package each time you want to upgrade manually.
To upgrade your installation of Docker Desktop, first stop any instance of Docker Desktop running locally,
then follow the regular installation steps to install the new version on top of the existing
version.

## Uninstall Docker Desktop

Docker Desktop can be removed from a Linux host using the package manager.

Once Docker Desktop has been removed, users must remove the `credsStore` and `currentContext` properties from the `~/.docker/config.json`.

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
