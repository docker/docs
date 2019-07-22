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

- Use {{ linux-dist-cap }} 64-bit 7.4 and higher on `x86_64`, or `s390x`.
- Use storage driver `overlay2` or `devicemapper` (`direct-lvm` mode in production).
- Find the URL for your Docker EE repo at [Docker Hub](https://hub.docker.com/my-content){: target="_blank" class="_" }.
- Uninstall old versions of Docker.
- Remove old Docker repos from `/etc/yum.repos.d/`.
- Disable SELinux on `s390x` (IBM Z) systems before install/upgrade.

### Architectures and storage drivers

Docker EE supports {{ linux-dist-long }} 64-bit, versions 7.4 and higher running on one of the following architectures: `x86_64`, or `s390x` (IBM Z). See [Compatability Matrix](https://success.docker.com/article/compatibility-matrix){: target="_blank" class="_" }) for specific details.

> Little-endian format only
>
> On IBM Power systems, Docker EE only supports little-endian format, `ppc64le`, even though {{ linux-dist-cap }} 7 ships both big and little-endian versions.

On {{ linux-dist-long }}, Docker EE supports storage drivers, `overlay2` and `devicemapper`. In Docker EE 17.06.2-ee-5 and higher, `overlay2` is the recommended storage driver. The following limitations apply:

- [OverlayFS](/storage/storagedriver/overlayfs-driver){: target="_blank" class="_" }: If `selinux` is enabled, the `overlay2` storage driver is supported on {{ linux-dist-cap }} 7.4 or higher. If `selinux` is disabled, `overlay2` is supported on {{ linux-dist-cap }} 7.2 or higher with kernel version 3.10.0-693 and higher.

- [Device Mapper](/storage/storagedriver/device-mapper-driver/){: target="_blank" class="_" }: On production systems using `devicemapper`, you must use `direct-lvm` mode, which requires one or more dedicated block devices. Fast storage such as solid-state media (SSD) is recommended. Do not start Docker until properly configured per the [storage guide](/storage/storagedriver/device-mapper-driver/){: target="_blank" class="_" }.

### FIPS 140-2 cryptographic module support

[Federal Information Processing Standards (FIPS) Publication 140-2](https://csrc.nist.gov/csrc/media/publications/fips/140/2/final/documents/fips1402.pdf) is a United States Federal security requirement for cryptographic modules. 

With Docker EE Basic license for versions 18.03 and later, Docker provides FIPS 140-2 support in RHEL 7.3, 7.4 and 7.5. This includes a FIPS supported cryptographic module. If the RHEL implementation already has FIPS support enabled, FIPS is automatically enabled in the Docker engine.

To verify the FIPS-140-2 module is enabled in the Linux kernel, confirm the file `/proc/sys/crypto/fips_enabled` contains `1`.

```
$ cat /proc/sys/crypto/fips_enabled
1
```

> **Note**: FIPS is only supported in the Docker Engine EE. UCP and DTR currently do not have support for FIPS-140-2.

To enable FIPS 140-2 compliance on a system that is not in FIPS 140-2 mode, do the following:

Create a file called `/etc/systemd/system/docker.service.d/fips-module.conf`. It needs to contain the following:

```
[Service]
Environment="DOCKER_FIPS=1"
```

Reload the Docker configuration to systemd.

`$ sudo systemctl daemon-reload`

Restart the Docker service as root.

`$ sudo systemctl restart docker`

To confirm Docker is running with FIPS-140-2 enabled, run the `docker info` command:

{% raw %}
```
docker info --format {{.SecurityOptions}}
[name=selinux name=fips]
```
{% endraw %}

### Disabling FIPS-140-2 

If the system has the FIPS 140-2 cryptographic module installed on the operating system, 
it is possible to disable FIPS-140-2 compliance. 

To disable FIPS 140-2 in Docker but not the operating system, set the value `DOCKER_FIPS=0` 
in the `/etc/systemd/system/docker.service.d/fips-module.conf`.

Reload the Docker configuration to systemd.

`$ sudo systemctl daemon-reload`

Restart the Docker service as root.

`$ sudo systemctl restart docker`

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
                  docker-engine
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
