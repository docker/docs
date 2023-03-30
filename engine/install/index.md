---
title: Install Docker Engine 
description: Choose the best method for you to install Docker Engine. This client-server application is available on Linux, Mac, Windows, and as a static binary.
keywords: install engine, docker engine install, install docker engine, docker engine installation, engine install, docker ce installation, docker ce install, engine installer, installing docker engine, docker server install
redirect_from:
- /cs-engine/
- /cs-engine/1.12/
- /cs-engine/1.12/upgrade/
- /cs-engine/1.13/
- /cs-engine/1.13/upgrade/
- /ee/docker-ee/oracle/
- /ee/supported-platforms/
- /en/latest/installation/
- /engine/installation/
- /engine/installation/frugalware/
- /engine/installation/linux/
- /engine/installation/linux/archlinux/
- /engine/installation/linux/cruxlinux/
- /engine/installation/linux/docker-ce/
- /engine/installation/linux/docker-ee/
- /engine/installation/linux/docker-ee/oracle/
- /engine/installation/linux/frugalware/
- /engine/installation/linux/gentoolinux/
- /engine/installation/linux/oracle/
- /engine/installation/linux/other/
- /engine/installation/oracle/
- /enterprise/supported-platforms/
- /install/linux/docker-ee/oracle/
---

> **Docker Desktop for Linux**
>
> Docker Desktop helps you build, share, and run containers on Mac and
> Windows as you do on Linux. Docker Desktop for
> Linux is now GA. For more information, see
[Docker Desktop for Linux](../../desktop/install/linux-install.md).
{: .important}

## Supported platforms

Docker Engine is available on a variety of [Linux distros](../../desktop/install/linux-install.md),
[macOS](../../desktop/install/mac-install.md), and [Windows 10](../../desktop/install/windows-install.md)
through Docker Desktop, and as a [static binary installation](binaries.md). Find
your preferred operating system below.

### Desktop

{% assign yes = '![yes](/assets/images/green-check.svg){: .inline style="height: 14px; margin: 0 auto"}' %}

| Platform                                                          | x86_64 / amd64                                   | arm64 (Apple Silicon)                            |
|:------------------------------------------------------------------|:------------------------------------------------:|:------------------------------------------------:|
| [Docker Desktop for Linux](../../desktop/install/linux-install.md)        | [{{ yes }}](../../desktop/install/linux-install.md)      |                                                  |
| [Docker Desktop for Mac (macOS)](../../desktop/install/mac-install.md)    | [{{ yes }}](../../desktop/install/mac-install.md)        | [{{ yes }}](../../desktop/install/mac-install.md)        |
| [Docker Desktop for Windows](../../desktop/install/windows-install.md)    | [{{ yes }}](../../desktop/install/windows-install.md)    |                                                  |

### Server

Docker provides `.deb` and `.rpm` packages from the following Linux distros
and architectures:

| Platform                | x86_64 / amd64         | arm64 / aarch64        | arm (32-bit)           | s390x                  |
|:------------------------|:-----------------------|:-----------------------|:-----------------------|:-----------------------|
| [CentOS](centos.md)     | [{{ yes }}](centos.md) | [{{ yes }}](centos.md) |                        |                        |
| [Debian](debian.md)     | [{{ yes }}](debian.md) | [{{ yes }}](debian.md) | [{{ yes }}](debian.md) |                        |
| [Fedora](fedora.md)     | [{{ yes }}](fedora.md) | [{{ yes }}](fedora.md) |                        |                        |
| [Raspbian](debian.md)   |                        |                        | [{{ yes }}](debian.md) |                        |
| [RHEL](rhel.md)         |                        |                        |                        | [{{ yes }}](rhel.md)   |
| [SLES](sles.md)         |                        |                        |                        | [{{ yes }}](sles.md)   |
| [Ubuntu](ubuntu.md)     | [{{ yes }}](ubuntu.md) | [{{ yes }}](ubuntu.md) | [{{ yes }}](ubuntu.md) | [{{ yes }}](ubuntu.md) |
| [Binaries](binaries.md) | [{{yes}}](binaries.md) | [{{yes}}](binaries.md) | [{{yes}}](binaries.md) |                        |

### Other Linux distros

> **Note**
>
> While the instructions below may work, Docker doesn't test or verify
> installation on distro derivatives.

- Users of Debian derivatives such as "BunsenLabs Linux", "Kali Linux" or 
  "LMDE" (Debian-based Mint) should follow the installation instructions for
  [Debian](debian.md), substituting the version of their distro for the
  corresponding Debian release. Refer to the documentation of your distro to find
  which Debian release corresponds with your derivative version.
- Likewise, users of Ubuntu derivatives such as "Kubuntu", "Lubuntu" or "Xubuntu"
  should follow the installation instructions for [Ubuntu](ubuntu.md),
  substituting the version of their distro for the corresponding Ubuntu release.
  Refer to the documentation of your distro to find which Ubuntu release
  corresponds with your derivative version.
- Some Linux distros provide a package of Docker Engine through their
  package repositories. These packages are built and maintained by the Linux
  distro's package maintainers and may have differences in configuration
  or built from modified source code. Docker isn't involved in releasing these
  packages and you should report any bugs or issues involving these packages to
  your Linux distro's issue tracker.

Docker provides [binaries](binaries.md) for manual installation of Docker Engine.
These binaries are statically linked and you can use them on any Linux distro.

## Release channels

Docker Engine has two types of update channels, **stable** and **test**:

* The **Stable** channel gives you latest releases for general availability.
* The **Test** channel gives pre-releases that are ready for testing before
  general availability (GA).

### Stable

Year-month releases are made from a release branch diverged from the master
branch. The branch is created with format `<year>.<month>`, for example
`20.10`. The year-month name indicates the earliest possible calendar
month to expect the release to be generally available. All further patch
releases are performed from that branch. For example, once `v20.10.0` is
released, all subsequent patch releases are built from the `20.10` branch.

### Test

In preparation for a new year-month release, a branch is created from
the master branch with format `YY.mm` when the milestones desired by
Docker for the release have achieved feature-complete. Pre-releases
such as betas and release candidates are conducted from their respective release
branches. Patch releases and the corresponding pre-releases are performed
from within the corresponding release branch.

## Support

Docker Engine releases of a year-month branch are supported with patches as
needed for one month after the next year-month general availability release.

This means bug reports and backports to release branches are assessed
until the end-of-life date.

After the year-month branch has reached end-of-life, the branch may be
deleted from the repository.

### Backporting

Backports to the Docker products are prioritized by the Docker company. A
Docker employee or repository maintainer will endeavour to ensure sensible
bugfixes make it into _active_ releases.

If there are important fixes that ought to be considered for backport to
active release branches, be sure to highlight this in the PR description
or by adding a comment to the PR.

### Upgrade path

Patch releases are always backward compatible with its year-month version.

### Licensing

Docker is licensed under the Apache License, Version 2.0. See
[LICENSE](https://github.com/moby/moby/blob/master/LICENSE) for the full
license text.

## Reporting security issues

The Docker maintainers take security seriously. If you discover a security
issue, please bring it to their attention right away!

Please DO NOT file a public issue; instead send your report privately
to security@docker.com.

Security reports are greatly appreciated, and Docker will publicly thank you
for it.

## Get started

After setting up Docker, you can learn the basics with
[Getting started with Docker](../../get-started/index.md).
