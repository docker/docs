---
description: Instructions for installing Docker EE on RHEL
keywords: requirements, installation, rhel, rpm, install, uninstall, upgrade, update
redirect_from:
- /engine/installation/rhel/
- /installation/rhel/
- /engine/installation/linux/rhel/
title: Get Docker EE for Red Hat Enterprise Linux
toc_max: 4
---

{% assign linux-dist = "rhel" %}
{% assign linux-dist-url-slug = "rhel" %}
{% assign linux-dist-long = "Red Hat Enterprise Linux" %}
{% assign package-format = "RPM" %}
{% assign gpg-fingerprint = "77FE DA13 1A83 1D29 A418  D3E8 99E5 FF2E 7668 2BC9" %}

{% include ee-linux-install-reuse.md section="ee-install-intro" %}

## Prerequisites

Docker Community Edition (Docker CE) is not supported on {{ linux-dist-long }}.

### Docker EE repository URL

{% include ee-linux-install-reuse.md section="ee-url-intro" %}

### OS requirements

To install Docker EE, you need the 64-bit version of {{ linux-dist-long }}
running on `x86_64`, `s390x` (IBM Z), or `ppc64le` (IBM Power) architectures.

In addition, you must use the `overlay2` or `devicemapper` storage driver.
Beginning with Docker EE 17.06.2-ee-5 the `overlay2` storage driver is the
recommended storage driver.

The following limitations apply:

**OverlayFS**:

- The `overlay2` storage driver is only supported on RHEL 7.2 or higher.
- If `selinux` is enabled, the `overlay2` storage driver is only supported on
  RHEL 7.4 or higher.

**Devicemapper**:

- On production systems using `devicemapper`, you must use `direct-lvm` mode,
  which requires one or more dedicated block devices. Fast storage such as
  solid-state media (SSD) is recommended.

{% capture selinux-warning %}
> **Warning**: There is currently no support for `selinux` on IBM Z systems. If
> you try to install Docker EE on an IBM Z system with `selinux` enabled, you get
> an error about the `container-selinux` package, which is missing from Red Hat's
> repository for IBM Z. The only current workaround is to disable `selinux`
> before installing or upgrading Docker on IBM Z.
{:.warning}
{% endcapture %}
{{ selinux-warning }}

### Uninstall old versions

Older versions of Docker were called `docker` or `docker-engine`. If these are
installed, uninstall them, along with associated dependencies.

```bash
$ sudo yum remove docker \
                  docker-common \
                  docker-selinux \
                  docker-engine-selinux \
                  docker-engine
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

{{ selinux-warning }}

{% include ee-linux-install-reuse.md section="install-using-yum-repo" %}

#### Upgrade Docker EE

{% include ee-linux-install-reuse.md section="upgrade-using-yum-repo" %}

### Install from a package

{{ selinux-warning }}

{% include ee-linux-install-reuse.md section="install-using-yum-package" %}

#### Upgrade Docker EE

{% include ee-linux-install-reuse.md section="upgrade-using-yum-package" %}

## Uninstall Docker EE

{% include ee-linux-install-reuse.md section="yum-uninstall" %}

## Next steps

{% include ee-linux-install-reuse.md section="linux-install-nextsteps" %}
