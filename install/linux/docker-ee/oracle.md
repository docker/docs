---
description: Instructions for installing Docker EE on Oracle Linux
keywords: requirements, installation, oracle, ol, rpm, install, uninstall, upgrade, update
redirect_from:
- /engine/installation/oracle/
- /engine/installation/linux/oracle/
- /engine/installation/linux/docker-ee/oracle/
title: Get Docker EE for Oracle Linux
---

{% assign linux-dist = "oraclelinux" %}
{% assign linux-dist-cap = "OL" %}
{% assign linux-dist-url-slug = "oraclelinux" %}
{% assign linux-dist-long = "Oracle Linux" %}
{% assign package-format = "RPM" %}
{% assign gpg-fingerprint = "77FE DA13 1A83 1D29 A418  D3E8 99E5 FF2E 7668 2BC9" %}


{% include ee-linux-install-reuse.md section="ee-install-intro" %}

## Prerequisites

This section lists what you need to consider before installing Docker EE. Items that require action are explained below.

- Use {{ linux-dist-cap }} 64-bit 7.3 or higher on RHCK 3.10.0-514 or higher.
- Use the `devicemapper` storage driver only (`direct-lvm` mode in production).
- Find the URL for your Docker EE repo at [Docker Store](https://store.docker.com/my-content){: target="_blank" class="_" }.
- Uninstall old versions of Docker.
- Remove old Docker repos from `/etc/yum.repos.d/`.
- Disable SELinux if installing or upgrading Docker EE 17.06.1.

### Architectures and storage drivers

Docker EE supports {{ linux-dist-long }} 64-bit, versions 7.3 and higher, running the Red Hat Compatible kernel (RHCK) 3.10.0-514 or higher. Older versions of {{ linux-dist-long }} are not supported.

On {{ linux-dist-long }}, Docker EE only supports the `devicemapper` storage driver. In production, you must use it in `direct-lvm` mode, which requires one or more dedicated block devices. Fast storage such as solid-state media (SSD) is recommended. Do not start Docker until properly configured per the [storage guide](/storage/storagedriver/device-mapper-driver/){: target="_blank" class="_" }.

### Find your Docker EE repo URL

{% include ee-linux-install-reuse.md section="find-ee-repo-url" %}

### Uninstall old Docker versions

The Docker EE package is called `docker-ee`. Older versions were called `docker` or `docker-engine`. Uninstall all older versions and associated dependencies. The contents of `/var/lib/docker/` are preserved, including images, containers, volumes, and networks.

```bash
$ sudo yum remove docker \
                  docker-engine \
                  docker-engine-selinux
```

## Repo install and upgrade

{% include ee-linux-install-reuse.md section="using-yum-repo" %}

{% capture selinux-warning %}
> Docker EE cannot install on {{ linux-dist-long }} with SELinux enabled
>
> If you have `selinux` enabled and you attempt to install Docker EE 17.06.1, you get an error that the `container-selinux` package cannot be found..
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
