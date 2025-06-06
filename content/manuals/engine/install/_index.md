---
title: Install Docker Engine (Linux)
linkTitle: Install
weight: 10
description: Overview of methods to install Docker Engine on Linux. 
keywords: install engine, docker engine install, install docker engine, docker engine
  installation, engine install, docker ce installation, docker ce install, engine
  installer, installing docker engine, docker server install, docker desktop vs docker engine, linux install
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
---

Learn how to install Docker Engine on Linux, also known as Docker CE. 

> [!NOTE] Docker Engine is also available for Windows, macOS, and Linux,
through [Docker Desktop](https://docs.docker.com/desktop/). For instructions on how to install Docker Desktop,
see: [Overview of Docker Desktop](/manuals/desktop/_index.md).

## Supported platforms
Click the name of the platform for specific install instructions.

| Platform                                       | x86_64 / amd64 | arm64 / aarch64 | arm (32-bit) | ppc64le | s390x |
| :--------------------------------------------- | :------------: | :-------------: | :----------: | :-----: | :---: |
| [CentOS](centos.md)                            |       ✅       |       ✅        |              |   ✅    |       |
| [Debian](debian.md)                            |       ✅       |       ✅        |      ✅      |   ✅    |       |
| [Fedora](fedora.md)                            |       ✅       |       ✅        |              |   ✅    |       |
| [Raspberry Pi OS (32-bit)](raspberry-pi-os.md) |                |                 |      ✅      |         |       |
| [RHEL](rhel.md)                                |       ✅       |       ✅        |              |         |  ✅   |
| [SLES](sles.md)                                |                |                 |              |         |  ✅   |
| [Ubuntu](ubuntu.md)                            |       ✅       |       ✅        |      ✅      |   ✅    |  ✅   |
| [Binaries](binaries.md)                        |       ✅       |       ✅        |      ✅      |         |       |

### Other Linux distributions

> [!NOTE] Docker doesn't test or verify installation on distribution derivatives. These guidelines may not work.

#### Debian derivatives
If you use Debian derivatives such as "BunsenLabs Linux", "Kali Linux" or "LMDE" (Debian-based Mint):

Follow the installation instructions for [Debian](debian.md) and substitute the version of your distribution for the
corresponding Debian release. Refer to the documentation of your distribution to find which Debian release corresponds 
with your derivative version.

#### Ubuntu derivatives 
If you use Ubuntu derivatives such as "Kubuntu", "Lubuntu" or "Xubuntu":

Follow the installation instructions for [Ubuntu](ubuntu.md) and substitute the version of your distribution for the 
corresponding Ubuntu release. Refer to the documentation of your distribution to find which Ubuntu release 
corresponds with your derivative version.

#### Linux Distributions with Docker
Some Linux distributions provide a package of Docker Engine through their
package repositories. These packages are built and maintained by the Linux
distributions' package maintainers and may have differences in configuration
or are built from modified source code. Docker isn't involved in releasing these
packages. Any bugs or issues involving these packages to your Linux distribution's issue tracker.

### Docker Binary
Docker provides [binaries](binaries.md) for manual installation of Docker Engine.
These binaries are statically linked and you can use them on any Linux distribution.

## Release channels

Docker Engine has two types of update channels, **stable** and **test**:

* The **stable** channel is for the latest versions released for general availability.
* The **test** channel is for pre-release versions that are ready for testing before
  general availability.

Use the test channel with caution. Pre-release versions include experimental and
early-access features that are subject to breaking changes.

## Get started

After setting up Docker, learn the basics with [Getting started with Docker](/get-started/introduction/_index.md).

## Report security issues

If you discover a security issue, immediately send a private report to security@docker.com.
**DO NOT file a public issue.**

Security reports are greatly appreciated, and Docker will publicly thank you for it.

## Support information

Docker Engine is an open source project, supported by the Moby project maintainers
and community members. Docker doesn't provide support for Docker Engine.
Docker provides support for Docker products, including Docker Desktop, which uses
Docker Engine as one of its components.

For information about the open source project, refer to the
[Moby project website](https://mobyproject.org/).

### Upgrade path

Patch releases are always backward compatible with their major and minor versions.

### Licensing

Docker Engine is licensed under the Apache License, Version 2.0. See
[LICENSE](https://github.com/moby/moby/blob/master/LICENSE) for the full
license text.
