---
title: Install Docker Engine
linkTitle: Install
weight: 10
description: Learn how to choose the best method for you to install Docker Engine. This client-server
  application is available on Linux, Mac, Windows, and as a static binary.
keywords: install engine, docker engine install, install docker engine, docker engine
  installation, engine install, docker ce installation, docker ce install, engine
  installer, installing docker engine, docker server install, docker desktop vs docker engine
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

This section describes how to install Docker Engine on Linux, also known as Docker CE. Docker Engine is also available for Windows, macOS, and Linux through Docker Desktop.

> [!TIP]
> To install Docker Desktop, see [Overview of Docker Desktop](/manuals/desktop/_index.md).

## Supported platforms

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

> [!NOTE]
> While the following instructions may work, Docker doesn't test or verify installation on distribution derivatives.

- For Debian derivatives (for example, BunsenLabs Linux, Kali Linux, LMDE):
  - Follow the [Debian installation instructions](debian.md).
  - Substitute your distribution version for the corresponding Debian release.
  - Check your distribution's documentation for the matching Debian release.
- For Ubuntu derivatives (for example, Kubuntu, Lubuntu, Xubuntu):
  - Follow the [Ubuntu installation instructions](ubuntu.md).
  - Substitute your distribution version for the corresponding Ubuntu release.
  - Check your distribution's documentation for the matching Ubuntu release.
- Some Linux distributions provide a Docker Engine package in their repositories. These packages are built and maintained by the distribution's maintainers and may have configuration differences or use modified source code. Docker isn't involved in releasing these packages. Report bugs or issues with these packages to your distribution's issue tracker.

Docker provides [binaries](binaries.md) for manual installation. These binaries are statically linked and work on any Linux distribution.

## Release channels

Docker Engine has two update channels: **stable** and **test**.

- The **stable** channel gives you the latest versions released for general availability.
- The **test** channel gives you pre-release versions that are ready for testing before general availability.

Use the test channel with caution. Pre-release versions include experimental and early-access features that are subject to breaking changes.

## Support

Docker Engine is an open source project, supported by the Moby project maintainers and community members. Docker doesn't provide support for Docker Engine. Docker provides support for Docker products, including Docker Desktop, which uses Docker Engine as one of its components.

For information about the open source project, refer to the [Moby project website](https://mobyproject.org/).

### Upgrade path

Patch releases are always backward compatible with their major and minor version.

### Licensing

Docker Engine is licensed under the Apache License, Version 2.0. See [LICENSE](https://github.com/moby/moby/blob/master/LICENSE) for the full license text.

## Reporting security issues

If you discover a security issue, report it immediately.

Do not file a public issue. Instead, send your report privately to security@docker.com.

Docker appreciates security reports and will publicly thank you for your contribution.

## Get started

After setting up Docker, you can learn the basics with [Getting started with Docker](/get-started/introduction/_index.md).
