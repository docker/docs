---
title: Install Docker Engine
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

This section describes how to install Docker Engine on Linux, also known as
Docker CE. Docker Engine is also available for Windows, macOS, and Linux,
through Docker Desktop. For instructions on how to install Docker Desktop, see:

- [Docker Desktop for Linux](../../desktop/install/linux-install.md)
- [Docker Desktop for Mac (macOS)](../../desktop/install/mac-install.md)
- [Docker Desktop for Windows](../../desktop/install/windows-install.md)

## Supported platforms

| Platform                                       | x86_64 / amd64    | arm64 / aarch64   | arm (32-bit)             | ppc64le         | s390x           |
| :--------------------------------------------- | :---------------- | :---------------- | :----------------------- | :-------------- | :-------------- |
| [CentOS](centos.md)                            | [✅](centos.md)   | [✅](centos.md)   |                          | [✅](centos.md) |                 |
| [Debian](debian.md)                            | [✅](debian.md)   | [✅](debian.md)   | [✅](debian.md)          | [✅](debian.md) |                 |
| [Fedora](fedora.md)                            | [✅](fedora.md)   | [✅](fedora.md)   |                          | [✅](fedora.md) |                 |
| [Raspberry Pi OS (32-bit)](raspberry-pi-os.md) |                   |                   | [✅](raspberry-pi-os.md) |                 |                 |
| [RHEL (s390x)](rhel.md)                        |                   |                   |                          |                 | [✅](rhel.md)   |
| [SLES](sles.md)                                |                   |                   |                          |                 | [✅](sles.md)   |
| [Ubuntu](ubuntu.md)                            | [✅](ubuntu.md)   | [✅](ubuntu.md)   | [✅](ubuntu.md)          | [✅](ubuntu.md) | [✅](ubuntu.md) |
| [Binaries](binaries.md)                        | [✅](binaries.md) | [✅](binaries.md) | [✅](binaries.md)        |                 |                 |

### Other Linux distros

> **Note**
>
> While the following instructions may work, Docker doesn't test or verify
> installation on distro derivatives.

- If you use Debian derivatives such as "BunsenLabs Linux", "Kali Linux" or 
  "LMDE" (Debian-based Mint) should follow the installation instructions for
  [Debian](debian.md), substitute the version of your distro for the
  corresponding Debian release. Refer to the documentation of your distro to find
  which Debian release corresponds with your derivative version.
- Likewise, if you use Ubuntu derivatives such as "Kubuntu", "Lubuntu" or "Xubuntu"
  you should follow the installation instructions for [Ubuntu](ubuntu.md),
  substituting the version of your distro for the corresponding Ubuntu release.
  Refer to the documentation of your distro to find which Ubuntu release
  corresponds with your derivative version.
- Some Linux distros provide a package of Docker Engine through their
  package repositories. These packages are built and maintained by the Linux
  distro's package maintainers and may have differences in configuration
  or are built from modified source code. Docker isn't involved in releasing these
  packages and you should report any bugs or issues involving these packages to
  your Linux distro's issue tracker.

Docker provides [binaries](binaries.md) for manual installation of Docker Engine.
These binaries are statically linked and you can use them on any Linux distro.

## Release channels

Docker Engine has two types of update channels, **stable** and **test**:

* The **stable** channel gives you the latest versions released for general availability.
* The **test** channel gives you pre-release versions that are ready for testing before
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

## Licensing

Docker Engine is licensed under the Apache License, Version 2.0. See
[LICENSE](https://github.com/moby/moby/blob/master/LICENSE) for the full
license text.

Commercial use of Docker Desktop (which includes Docker Engine and Docker CLI)
requires a [paid subscription](https://www.docker.com/pricing/)
in enterprises exceeding 250 employees OR with annual revenue surpassing $10 million USD.

To check if you're using Docker Engine through Docker Desktop,
run the `docker version` command and inspect the server version.

```console {hl_lines=12}
$ docker version
Client:
 Cloud integration: v1.0.35+desktop.13
 Version:           26.0.0
 API version:       1.45
 Go version:        go1.21.8
 Git commit:        2ae903e
 Built:             Wed Mar 20 15:14:46 2024
 OS/Arch:           darwin/arm64
 Context:           desktop-linux

Server: Docker Desktop 4.29.0 (143713)
 Engine:
  Version:          26.0.0
  API version:      1.45 (minimum version 1.24)
  Go version:       go1.21.8
  Git commit:       8b79278
  Built:            Wed Mar 20 15:18:02 2024
  OS/Arch:          linux/arm64
  Experimental:     false
 containerd:
  Version:          1.6.28
  GitCommit:        ae07eda36dd25f8a1b98dfbf587313b99c0190bb
 runc:
  Version:          1.1.12
  GitCommit:        v1.1.12-0-g51d5e94
 docker-init:
  Version:          0.19.0
  GitCommit:        de40ad0
```

## Reporting security issues

If you discover a security issue, we request that you bring it to our attention immediately.

DO NOT file a public issue. Instead, submit your report privately to security@docker.com.

Security reports are greatly appreciated, and Docker will publicly thank you for it.

## Get started

After setting up Docker, you can learn the basics with
[Getting started with Docker](../../get-started/index.md).
