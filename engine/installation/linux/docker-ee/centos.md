---
description: Instructions for installing Docker EE on CentOS
keywords: requirements, apt, installation, centos, rpm, install, uninstall, upgrade, update
redirect_from:
- /engine/installation/centos/
title: Get Docker EE for CentOS
toc_max: 4
---

{% assign linux-dist = "centos" %}
{% assign linux-dist-url-slug = "centos" %}
{% assign linux-dist-long = "Centos" %}
{% assign package-format = "RPM" %}
{% assign gpg-fingerprint = "77FE DA13 1A83 1D29 A418  D3E8 99E5 FF2E 7668 2BC9" %}

{% include ee-linux-install-reuse.md section="ee-install-intro" %}

## Prerequisites

Docker CE users should go to
[Get docker CE for CentOS](https://docs.docker.com/engine/installation/linux/docker-ce/centos/)
**instead of this topic**.

### Docker EE repository URL

{% include ee-linux-install-reuse.md section="ee-url-intro" %}

### OS requirements

To install Docker EE, you need a maintained version of CentOS 7. Archived
versions aren't supported or tested.

The `centos-extras` repository must be enabled. This repository is enabled by
default, but if you have disabled it, you need to
[re-enable it](https://wiki.centos.org/AdditionalResources/Repositories){: target="_blank" class="_" }.

In addition, you must use the `overlay2` or `devicemapper` storage driver if you
use Docker EE. On production systems using `devicemapper`, you must use
`direct-lvm` mode, which requires one or more dedicated block devices. Fast
storage such as solid-state media (SSD) is recommended.

### Uninstall old versions

Older versions of Docker were called `docker` or `docker-engine`. In addition,
if you are upgrading from Docker CE to Docker EE, remove the Docker CE package.

```bash
$ sudo yum remove docker \
                  docker-common \
                  docker-selinux \
                  docker-engine-selinux \
                  docker-engine \
                  docker-ce
```

It's OK if `yum` reports that none of these packages are installed.

The contents of `/var/lib/docker/`, including images, containers, volumes, and
networks, are preserved. The Docker EE package is now called `docker-ee`.

## Install Docker EE

{% include ee-linux-install-reuse.md section="ways-to-install" %}

### Install using the repository

Before you install Docker EE for the first time on a new host machine, you need
to set up the Docker EE repository. Afterward, you can install and update Docker
EE from the repository.

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
