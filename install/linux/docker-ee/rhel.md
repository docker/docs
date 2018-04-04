---
description: Instructions for installing Docker EE on RHEL
keywords: requirements, installation, rhel, rpm, install, uninstall, upgrade, update
redirect_from:
- /engine/installation/rhel/
- /installation/rhel/
- /engine/installation/linux/rhel/
- /engine/installation/linux/docker-ee/rhel/
title: Get Docker EE for Red Hat Enterprise Linux
---

{% assign linux-dist = "rhel" %}
{% assign linux-dist-cap = "RHEL" %}
{% assign linux-dist-url-slug = "rhel" %}
{% assign linux-dist-long = "Red Hat Enterprise Linux" %}
{% assign package-format = "RPM" %}
{% assign gpg-fingerprint = "77FE DA13 1A83 1D29 A418  D3E8 99E5 FF2E 7668 2BC9" %}


{% include ee-linux-install-reuse.md section="ee-install-intro" %}

## Prerequisites

This section lists what you need to consider before installing Docker EE. Items that require action are explained below.

- Use {{ linux-dist-cap }} 64-bit 7.1 and higher on `x86_64`, `s390x`, or `ppc64le` (not ppc64).
- Use storage driver `overlay2` or `devicemapper` (`direct-lvm` mode in production).
- Find the URL for your Docker EE repo at [Docker Store](https://store.docker.com/my-content){: target="_blank" class="_" }.
- Uninstall old versions of Docker.
- Remove old Docker repos from `/etc/yum.repos.d/`.
- Disable SELinux on IBM Power systems before install/upgrade.

### Architectures and storage drivers

Docker EE supports {{ linux-dist-long }} 64-bit, versions 7.1 and higher (7.1, 7.2, 7.3, 7.4), running on one of the following architectures: `x86_64`, `s390x` (IBM Z), or `ppc64le` (IBM Power PC, little endian format). To ensure you have `ppc64le` (and not `ppc64`), run the command, `uname -m`.

> Little endian format only
>
> On IBM Power PC systems, Docker EE only supports little endian format, `ppc64le`, even though {{ linux-dist-cap }} 7 ships both big and little endian versions.

On {{ linux-dist-long }}, Docker EE supports storage drivers, `overlay2` and `devicemapper`. In Docker EE 17.06.2-ee-5 and higher, `overlay2` is the recommended storage driver. The following limitations apply:

- [OverlayFS](/storage/storagedriver/overlayfs-driver){: target="_blank" class="_" }: If `selinux` is enabled, the `overlay2` storage driver is supported on {{ linux-dist-cap }} 7.4 or higher. If `selinux` is disabled, `overlay2` is supported on {{ linux-dist-cap }} 7.2 or higher with kernel version 3.10.0-693 and higher.

- [Device Mapper](/storage/storagedriver/device-mapper-driver/){: target="_blank" class="_" }: On production systems using `devicemapper`, you must use `direct-lvm` mode, which requires one or more dedicated block devices. Fast storage such as solid-state media (SSD) is recommended. Do not start Docker until properly configured per the [storage guide](/storage/storagedriver/device-mapper-driver/){: target="_blank" class="_" }.

### Find your Docker EE repo URL

{% include ee-linux-install-reuse.md section="find-ee-repo-url" %}

### Uninstall old Docker versions

The Docker EE package is called `docker-ee`. Older versions were called `docker` or `docker-engine`. Uninstall all older versions and associated dependencies. The contents of `/var/lib/docker/` are preserved, including images, containers, volumes, and networks.

```bash
$ sudo yum remove docker \
                  docker-client \
                  docker-client-latest \
                  docker-common \
                  docker-latest \
                  docker-latest-logrotate \
                  docker-logrotate \
                  docker-selinux \
                  docker-engine-selinux \
                  docker-engine \
                  docker-ce
```



## Repo install and upgrade

{% include ee-linux-install-reuse.md section="using-yum-repo" %}

{% capture selinux-warning %}
> Disable SELinux before installing Docker EE on IBM Z systems
>
> There is currently no support for `selinux` on IBM Z systems. If you attempt to install or upgrade Docker EE on an IBM Z system with `selinux` enabled, an error is thrown that the `container-selinux` package is not found. Disable `selinux` before installing or upgrading Docker on IBM Z.
{:.warning}
{% endcapture %}
{{ selinux-warning }}

### Set up the repository

{% include ee-linux-install-reuse.md section="set-up-yum-repo" %}

### Install from the repository

{% include ee-linux-install-reuse.md section="install-using-yum-repo" %}

### Upgrade from the repository

{% include ee-linux-install-reuse.md section="upgrade-using-yum-repo" %}



## Package install and upgrade

{% include ee-linux-install-reuse.md section="package-installation" %}

{{ selinux-warning }}

### Install with a package

{% include ee-linux-install-reuse.md section="install-using-yum-package" %}

### Upgrade with a package

{% include ee-linux-install-reuse.md section="upgrade-using-yum-package" %}


## Uninstall Docker EE

{% include ee-linux-install-reuse.md section="yum-uninstall" %}

## Next steps

{% include ee-linux-install-reuse.md section="linux-install-nextsteps" %}
