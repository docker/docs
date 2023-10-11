---
title: Install Docker Engine
description: Choose the best method for you to install Docker Engine. This client-server
  application is available on Linux, Mac, Windows, and as a static binary.
keywords: install engine, docker engine install, install docker engine, docker engine
  installation, engine install, docker ce installation, docker ce install, engine
  installer, installing docker engine, docker server install
aliases:
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
'yes': '![yes](/assets/images/green-check.svg){: .inline style="height: 14px; margin:
  0 auto"}'
---

> **Docker Desktop for Linux**
>
> Docker Desktop helps you build, share, and run containers on Mac and
> Windows as you do on Linux. Docker Desktop for
> Linux is now GA. For more information, see
[Docker Desktop for Linux](../../desktop/install/linux-install.md).
{ .important }

## Supported platforms

Docker Engine is available on a variety of [Linux distros](../../desktop/install/linux-install.md),
[macOS](../../desktop/install/mac-install.md), and [Windows 10](../../desktop/install/windows-install.md)
through Docker Desktop, and as a [static binary installation](binaries.md). Find
your preferred operating system below.

### Desktop


| Platform                                                               |                    x86_64 / amd64                     |               arm64 (Apple Silicon)               |
| :--------------------------------------------------------------------- | :---------------------------------------------------: | :-----------------------------------------------: |
| [Docker Desktop for Linux](../../desktop/install/linux-install.md)     |  [ ✅  ](../../desktop/install/linux-install.md)  |                                                   |
| [Docker Desktop for Mac (macOS)](../../desktop/install/mac-install.md) |   [ ✅  ](../../desktop/install/mac-install.md)   | [ ✅  ](../../desktop/install/mac-install.md) |
| [Docker Desktop for Windows](../../desktop/install/windows-install.md) | [ ✅  ](../../desktop/install/windows-install.md) |                                                   |

### Server

Docker provides `.deb` and `.rpm` packages from the following Linux distros
and architectures:

| Platform                                       | x86_64 / amd64      | arm64 / aarch64     | arm (32-bit)               | ppc64le           | s390x             |
| :--------------------------------------------- | :------------------ | :------------------ | :------------------------- | :---------------- | :---------------- |
| [CentOS](centos.md)                            | [ ✅ ](centos.md)   | [ ✅ ](centos.md)   |                            | [ ✅ ](centos.md) |                   |
| [Debian](debian.md)                            | [ ✅ ](debian.md)   | [ ✅ ](debian.md)   | [ ✅ ](debian.md)          | [ ✅ ](debian.md) |                   |
| [Fedora](fedora.md)                            | [ ✅ ](fedora.md)   | [ ✅ ](fedora.md)   |                            | [ ✅ ](fedora.md) |                   |
| [Raspberry Pi OS (32-bit)](raspberry-pi-os.md) |                     |                     | [ ✅ ](raspberry-pi-os.md) |                   |                   |
| [RHEL](rhel.md)                                |                     |                     |                            |                   | [ ✅ ](rhel.md)   |
| [SLES](sles.md)                                |                     |                     |                            |                   | [ ✅ ](sles.md)   |
| [Ubuntu](ubuntu.md)                            | [ ✅ ](ubuntu.md)   | [ ✅ ](ubuntu.md)   | [ ✅ ](ubuntu.md)          | [ ✅ ](ubuntu.md) | [ ✅ ](ubuntu.md) |
| [Binaries](binaries.md)                        | [ ✅ ](binaries.md) | [ ✅ ](binaries.md) | [ ✅ ](binaries.md)        |                   |                   |

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

* The **Stable** channel gives you the latest versions released for general availability.
* The **Test** channel gives you pre-release versions that are ready for testing before
  general availability.

Use the test channel with caution. Pre-release versions include experimental and
early-access features that are subject to breaking changes.

## Support

Docker Engine is an open source project, supported by the Moby project maintainers
and community members. Docker doesn't provide support for Docker Engine.
Docker provides support for Docker products, including Docker Desktop, which uses
Docker Engine as one of its components.

For information about the open source project, refer to the
[Moby project website](https://mobyproject.org/).

### Upgrade path

Patch releases are always backward compatible with its major and minor version.

### Licensing

Docker Engine is licensed under the Apache License, Version 2.0. See
[LICENSE](https://github.com/moby/moby/blob/master/LICENSE) for the full
license text.

## Reporting security issues

If you discover a security issue, we request that you bring it to our attention immediately.

DO NOT file a public issue. Instead, submit your report privately to security@docker.com.

Security reports are greatly appreciated, and Docker will publicly thank you for it.

## Get started

After setting up Docker, you can learn the basics with
[Getting started with Docker](../../get-started/index.md).
