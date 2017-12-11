---
description: Instructions for installing Docker EE on Oracle Linux
keywords: requirements, installation, oracle, ol, rpm, install, uninstall, upgrade, update
redirect_from:
- /engine/installation/oracle/
- /engine/installation/linux/oracle/
title: Get Docker EE for Oracle Linux
toc_max: 4
---

{% assign linux-dist = "oraclelinux" %}
{% assign linux-dist-url-slug = "oraclelinux" %}
{% assign linux-dist-long = "Oracle Linux" %}
{% assign package-format = "RPM" %}
{% assign gpg-fingerprint = "77FE DA13 1A83 1D29 A418  D3E8 99E5 FF2E 7668 2BC9" %}

To get started with Docker EE on {{ linux-dist-long }}, make sure you
[meet the prerequisites](#prerequisites), then
[install Docker](#install-docker-ee).

## Prerequisites

Docker Community Edition (Docker CE) is not supported on {{ linux-dist-long }}.

### Docker EE repository URL

{% include ee-linux-install-reuse.md section="ee-url-intro" %}

### OS requirements

To install Docker EE, you need the 64-bit version of {{ linux-dist-long }} 7.3
or higher, running the Red Hat Compatible kernel (RHCK) 3.10.0-514 or higher.
Older versions of {{ linux-dist-long }} are not supported.

In addition, you must use the `devicemapper` storage driver if you use
Docker EE. On production systems, you must use `direct-lvm` mode, which
requires one or more dedicated block devices. Fast storage such as solid-state
media (SSD) is recommended.

> **Docker EE will not install on {{ linux-dist }} with `selinux` enabled!**
>
> If you have `selinux` enabled and you attempt to install Docker EE 17.06.1,
> you will get an error that the `container-selinux` package cannot be found.
{:.warning }

### Uninstall old versions

Older versions of Docker were called `docker` or `docker-engine`. If these are
installed, uninstall them, along with associated dependencies.

```bash
$ sudo yum remove docker \
                  docker-engine \
                  docker-engine-selinux
```

It's OK if `yum` reports that none of these packages are installed.

The contents of `/var/lib/docker/`, including images, containers, volumes, and
networks, are preserved. The Docker EE package is now called `docker-ee`.

## Install Docker EE

{% include ee-linux-install-reuse.md section="ways-to-install" %}

### Install using the repository

Before you install Docker EE for the first time on a new host machine, you need
to set up the Docker repository. Afterward, you can install and update Docker EE
from the repository.

#### Set up the repository

{% include ee-linux-install-reuse.md section="set-up-yum-repo" %}

#### Install Docker EE

{% include ee-linux-install-reuse.md section="install-using-yum-repo" %}

#### Upgrade Docker EE

{% include ee-linux-install-reuse.md section="upgrade-using-yum-repo" %}

### Install from a package

{% include ee-linux-install-reuse.md section="install-using-yum-package" %}

#### Upgrade Docker EE

{% include ee-linux-install-reuse.md section="upgrade-using-yum-package" %}

## Uninstall Docker EE

{% include ee-linux-install-reuse.md section="yum-uninstall" %}

## Next steps

{% include ee-linux-install-reuse.md section="linux-install-nextsteps" %}
