---
description: Instructions for installing Docker Desktop on RHEL
keywords: red hat, red hat enterprise linux, rhel, rpm,
  update install, uninstall, upgrade, update, linux,
  desktop, docker desktop, docker desktop for linux, dd4l
title: Install Docker Desktop on RHEL
toc_max: 4
download-url-base: https://download.docker.com/linux/rhel
---

> **Early Access**
>
> Docker Desktop on RHEL is in [Early Access](../../release-lifecycle.md).
{ .restricted }

> **Docker Desktop terms**
>
> Commercial use of Docker Desktop in larger enterprises (more than 250
> employees OR more than $10 million USD in annual revenue) requires a [paid
> subscription](https://www.docker.com/pricing/).

This page contains information on how to install, launch and upgrade Docker Desktop on a Red Hat Enterprise Linux (RHEL) distribution.

## Prerequisites

To install Docker Desktop successfully, you must:

- Meet the [system requirements](linux-install.md#general-system-requirements).
- Have a 64-bit version of either RHEL 8 or RHEL 9.

If you don't have `pass` installed, or it can't be installed, you must enable
[CodeReady Linux Builder (CRB) repository](https://access.redhat.com/articles/4348511)
and
[Extra Packages for Enterprise Linux (EPEL)](https://docs.fedoraproject.org/en-US/epel/).

   {{< tabs group="os_version" >}}
   {{< tab name="RHEL 9" >}}
```console
$ sudo subscription-manager repos --enable codeready-builder-for-rhel-9-$(arch)-rpms
$ sudo dnf install https://dl.fedoraproject.org/pub/epel/epel-release-latest-9.noarch.rpm
$ sudo dnf install pass
```

   {{< /tab >}}
   {{< tab name="RHEL 8" >}}
```console
$ sudo subscription-manager repos --enable codeready-builder-for-rhel-8-$(arch)-rpms
$ sudo dnf install https://dl.fedoraproject.org/pub/epel/epel-release-latest-8.noarch.rpm
$ sudo dnf install pass
```

   {{< /tab >}}
   {{< /tabs >}}

Additionally, for a GNOME desktop environment you must install AppIndicator and KStatusNotifierItem [GNOME extensions](https://extensions.gnome.org/extension/615/appindicator-support/). You must also enable EPEL.

   {{< tabs group="os_version" >}}
   {{< tab name="RHEL 9" >}}
```console
$ # enable EPEL as described above
$ sudo dnf install gnome-shell-extension-appindicator
$ sudo gnome-extensions enable appindicatorsupport@rgcjonas.gmail.com
```

   {{< /tab >}}
   {{< tab name="RHEL 8" >}}
```console
$ # enable EPEL as described above
$ sudo dnf install gnome-shell-extension-appindicator
$ sudo dnf install gnome-shell-extension-desktop-icons
$ sudo gnome-shell-extension-tool -e appindicatorsupport@rgcjonas.gmail.com
```

   {{< /tab >}}
   {{< /tabs >}}

For non-GNOME desktop environments, `gnome-terminal` must be installed:

```console
$ sudo dnf install gnome-terminal
```

## Install Docker Desktop

To install Docker Desktop on RHEL:

1. Set up Docker's package repository as follows:

   ```console
   $ sudo dnf config-manager --add-repo {{% param "download-url-base" %}}/docker-ce.repo
   ```

2. Download the latest RPM package.

3. Install the package with dnf as follows:

   ```console
   $ sudo dnf install ./docker-desktop-<version>-<arch>-rhel.rpm
   ```

There are a few post-install configuration steps done through the post-install script contained in the RPM package.

The post-install script:

- Sets the capability on the Docker Desktop binary to map privileged ports and set resource limits.
- Adds a DNS name for Kubernetes to `/etc/hosts`.
- Creates a symlink from `/usr/local/bin/com.docker.cli` to `/usr/bin/docker`.
  This is because the classic Docker CLI is installed at `/usr/bin/docker`. The Docker Desktop installer also installs a Docker CLI binary that includes cloud-integration capabilities and is essentially a wrapper for the Compose CLI, at`/usr/local/bin/com.docker.cli`. The symlink ensures that the wrapper can access the classic Docker CLI. 
- Creates a symlink from `/usr/libexec/qemu-kvm` to `/usr/local/bin/qemu-system-x86_64`.

## Launch Docker Desktop

{{< include "desktop-linux-launch.md" >}}

>**Tip**
>
> To attach Red Hat subscription data to containers, see [Red Hat verified solution](https://access.redhat.com/solutions/5870841).
>
> For example:
> ```console
> $ docker run --rm -it -v "/etc/pki/entitlement:/etc/pki/entitlement" -v "/etc/rhsm:/etc/rhsm-host" -v "/etc/yum.repos.d/redhat.repo:/etc/yum.repos.d/redhat.repo" registry.access.redhat.com/ubi9
> ```
{ .tip }

## Upgrade Docker Desktop

Once a new version for Docker Desktop is released, the Docker UI shows a notification.
You need to first remove the previous version and then download the new package each time you want to upgrade Docker Desktop. Run:

```console
$ sudo dnf remove docker-desktop
$ sudo dnf install ./docker-desktop-<version>-<arch>-rhel.rpm
```

## Next steps

- Take a look at the [Get started](../../guides/get-started/_index.md) training modules to learn how to build an image and run it as a containerized application.
- [Explore Docker Desktop](../use-desktop/index.md) and all its features.
- Review the topics in [Develop with Docker](../../develop/index.md) to learn how to build new applications using Docker.
