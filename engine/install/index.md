---
title: Install Docker Engine
description: Lists the installation methods
keywords: docker, installation, install, Docker Engine, Docker Engine, docker editions, stable, edge
redirect_from:
- /engine/installation/linux/
- /engine/installation/linux/frugalware/
- /engine/installation/frugalware/
- /engine/installation/linux/other/
- /engine/installation/linux/archlinux/
- /engine/installation/linux/cruxlinux/
- /engine/installation/linux/gentoolinux/
- /engine/installation/linux/docker-ce/
- /engine/installation/linux/docker-ee/
- /engine/installation/
- /en/latest/installation/
toc_max: 2
---


## Supported platforms

Docker Engine is available on a variety of [Linux platforms](#server),
[macOS](../../docker-for-mac/install.md) and [Windows 10](../../docker-for-windows/install.md)
through Docker Desktop, and as a [static binary installation](binaries.md). Find
your preferred operating system below.

#### Desktop

{% assign yes = '![yes](/images/green-check.svg){: style="height: 14px; margin: 0 auto"}' %}

| Platform                                                          | x86_64 / amd64                                   |
|:------------------------------------------------------------------|:------------------------------------------------:|
| [Docker Desktop for Mac (macOS)](../../docker-for-mac/install.md) | [{{ yes }}](../../docker-for-mac/install.md)     |
| [Docker Desktop for Windows](../../docker-for-windows/install.md) | [{{ yes }}](../../docker-for-windows/install.md) |

#### Server

Docker provides `.deb` and `.rpm` packages from the following Linux distributions
and architectures:

| Platform              | x86_64 / amd64         | ARM                      | ARM64 / AARCH64        |
|:----------------------|:-----------------------|:-------------------------|:-----------------------|
| [CentOS](centos.md)   | [{{ yes }}](centos.md) |                          | [{{ yes }}](centos.md) |
| [Debian](debian.md)   | [{{ yes }}](debian.md) | [{{ yes }}](debian.md)   | [{{ yes }}](debian.md) |
| [Fedora](fedora.md)   | [{{ yes }}](fedora.md) |                          | [{{ yes }}](fedora.md) |
| [Raspbian](debian.md) |                        | [{{ yes }}](debian.md)   | [{{ yes }}](debian.md) |
| [Ubuntu](ubuntu.md)   | [{{ yes }}](ubuntu.md) | [{{ yes }}](ubuntu.md)   | [{{ yes }}](ubuntu.md) |

##### Other Linux distributions

> **Note**
>
> While the instructions below may work, Docker does not test or verify
> installation on derivatives.

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
- Some Linux distributions are providing a package of Docker Engine through their
  package repositories. These packages are built and maintained by the Linux
  distribution's package maintainers and may have differences in configuration
  or built from modified source code. Docker is not involved in releasing these
  packages and bugs or issues involving these packages should be reported in
  your Linux distribution's issue tracker.

Docker provides [binaries](binaries.md) for manual installation of Docker Engine.
These binaries are statically linked and can be used on any Linux distribution.

## Release channels

Docker Engine has three types of update channels, **stable**, **test**,
and **nightly**:

* The **Stable** channel gives you latest releases for general availability.
* The **Test** channel gives pre-releases that are ready for testing before
  general availability (GA).
* The **Nightly** channel gives you latest builds of work in progress for the
  next major release.

### Stable

Year-month releases are made from a release branch diverged from the master
branch. The branch is created with format `<year>.<month>`, for example
`18.09`. The year-month name indicates the earliest possible calendar
month to expect the release to be generally available. All further patch
releases are performed from that branch. For example, once `v18.09.0` is
released, all subsequent patch releases are built from the `18.09` branch.

### Test

In preparation for a new year-month release, a branch is created from
the master branch with format `YY.mm` when the milestones desired by
Docker for the release have achieved feature-complete. Pre-releases
such as betas and release candidates are conducted from their respective release
branches. Patch releases and the corresponding pre-releases are performed
from within the corresponding release branch.

### Nightly

Nightly builds give you the latest builds of work in progress for the next major
release. They are created once per day from the master branch with the version
format:

    0.0.0-YYYYmmddHHMMSS-abcdefabcdef

where the time is the commit time in UTC and the final suffix is the prefix
of the commit hash, for example `0.0.0-20180720214833-f61e0f7`.

These builds allow for testing from the latest code on the master branch.

> **Note:**
> No qualifications or guarantees are made for the nightly builds.

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
