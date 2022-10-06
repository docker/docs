---
description: How to install Docker Desktop on Linux
keywords: linux, desktop, docker desktop, docker desktop for linux, dd4l, install, system requirements
title: Install Docker Desktop on Linux
redirect_from:
- /desktop/linux/install/
---
> **Docker Desktop terms**
>
> Commercial use of Docker Desktop in larger enterprises (more than 250
> employees OR more than $10 million USD in annual revenue) requires a paid
> subscription.

 This page contains information about general system requirements, supported platforms, and instructions on how to install Docker Desktop for Linux.

## System requirements

To install Docker Desktop successfully, your Linux host must meet the following general requirements:

- 64-bit kernel and CPU support for virtualization.
- KVM virtualization support. Follow the [KVM virtualization support instructions](#kvm-virtualization-support) to check if the KVM kernel modules are enabled and how to provide access to the kvm device.
- **QEMU must be version 5.2 or newer**. We recommend upgrading to the latest version.
- systemd init system.
- Gnome or KDE Desktop environment.
  - For many Linux distros, the Gnome environment does not support tray icons. To add support for tray icons, you need to install a Gnome extension. For example, [AppIndicator](https://extensions.gnome.org/extension/615/appindicator-support/).
- At least 4 GB of RAM.
- Enable configuring ID mapping in user namespaces, see [File sharing](../faqs/linuxfaqs.md#how-do-i-enable-file-sharing).

Docker Desktop for Linux runs a Virtual Machine (VM). For more information on why, see [Why Docker Desktop for Linux runs a VM](../faqs/linuxfaqs.md#why-does-docker-desktop-for-linux-run-a-vm).

> **Note:**
>
> Docker does not provide support for running Docker Desktop in nested virtualization scenarios. We recommend that you run Docker Desktop for Linux natively on supported distributions.

## Supported platforms

Docker provides `.deb` and `.rpm` packages from the following Linux distributions
and architectures:


{% assign yes = '![yes](/assets/images/green-check.svg){: .inline style="height: 14px; margin: 0 auto"}' %}


| Platform                | x86_64 / amd64         | 
|:------------------------|:-----------------------|
| [Ubuntu](ubuntu.md)     | [{{ yes }}](ubuntu.md) |
| [Debian](debian.md)     | [{{ yes }}](debian.md) |
| [Fedora](fedora.md)     | [{{ yes }}](fedora.md) |


>  **Note:**
>
> An experimental package is available for [Arch](archlinux.md)-based distributions. Docker has not tested or verified the installation.

Docker supports Docker Desktop on the current LTS release of the aforementioned distributions and the most recent version. As new versions are made available, Docker stops supporting the oldest version and supports the newest version.


### KVM virtualization support


Docker Desktop runs a VM that requires [KVM support](https://www.linux-kvm.org).

The `kvm` module should load automatically if the host has virtualization support. To load the module manually, run:

```console
$ modprobe kvm
```

Depending on the processor of the host machine, the corresponding module must be loaded:

```console
$ modprobe kvm_intel  # Intel processors

$ modprobe kvm_amd    # AMD processors
```

If the above commands fail, you can view the diagnostics by running:

```console
$ kvm-ok
```

To check if the KVM modules are enabled, run:

```console
$ lsmod | grep kvm
kvm_amd               167936  0
ccp                   126976  1 kvm_amd
kvm                  1089536  1 kvm_amd
irqbypass              16384  1 kvm
```

#### Set up KVM device user permissions


To check ownership of `/dev/kvm`, run :

```console
$ ls -al /dev/kvm
```

Add your user to the kvm group in order to access the kvm device:

```console
$ sudo usermod -aG kvm $USER
```

Log out and log back in so that your group membership is re-evaluated.


## Generic installation steps

>Important
>
>Make sure you meet the system requirements outlined earlier and follow the distro-specific prerequisites.
{: .important} 

1. Download the correct package for your Linux distribution and install it with the corresponding package manager. 
 - [Install on Debian](debian.md)
 - [Install on Fedora](fedora.md)
 - [Install on Ubuntu](ubuntu.md)
 - [Install on Arch](archlinux.md) 

2. Open your **Applications** menu in Gnome/KDE Desktop and search for **Docker Desktop**.

    ![Docker app in Applications](images/docker-app-in-apps.png)

3. Select **Docker Desktop** to start Docker. <br> The Docker menu (![whale menu](images/whale-x.svg){: .inline}) displays the Docker Subscription Service Agreement window.

4. Select **Accept** to continue. Docker Desktop starts after you accept the terms.

    > **Important**
    >
    > If you do not agree to the terms, the Docker Desktop application will close and  you can no longer run Docker Desktop on your machine. You can choose to accept the terms at a later date by opening Docker Desktop.
    {: .important}

    For more information, see [Docker Desktop Subscription Service Agreement](https://www.docker.com/legal/docker-subscription-service-agreement). We recommend that you also read the [FAQs](https://www.docker.com/pricing/faq){: target="_blank" rel="noopener" class="_" id="dkr_docs_desktop_install_btl"}.


## Where to go next

- [Troubleshooting](../troubleshoot/overview.md) describes common problems, workarounds, how to run and submit diagnostics, and submit issues.
- [FAQs](../faqs/general.md) provide answers to frequently asked questions.
- [Release notes](../release-notes.md) lists component updates, new features, and improvements associated with Docker Desktop releases.
- [Get started with Docker](../../get-started/index.md) provides a general Docker tutorial.
* [Back up and restore data](../backup-and-restore.md) provides instructions
  on backing up and restoring data related to Docker.
