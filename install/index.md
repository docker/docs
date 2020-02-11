---
title: Docker Engine overview
description: Lists the installation methods
keywords: docker, installation, install, Docker Engine - Community, Docker Engine - Enterprise, docker editions, stable, edge
redirect_from:
- /install/overview/
- /installation/
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
- /linux/
- /edge/
toc_max: 2
---

Docker Engine is an open source containerization technology for building and
containerizing your applications. Docker Engine acts as a client-server
application with:
* A server with a long-running daemon process [`dockerd`](/engine/reference/commandline/dockerd/).
* APIs which specify interfaces that programs can use to talk to and
  instruct the Docker daemon.
* A command line interface (CLI) client [`docker`](/engine/reference/commandline/cli/).

The CLI uses Docker APIs to control or interact with the Docker daemon
through scripting or direct CLI commands. Many other Docker applications use the
underlying API and CLI. The daemon creates and manage Docker objects, such as
images, containers, networks, and volumes.

Docker Engine has three types of update channels, **stable**, **test**, and **nightly**:

* **Stable** gives you latest releases for general availability.
* **Test** gives pre-releases that are ready for testing before general availability.
* **Nightly** gives you latest builds of work in progress for the next major release.

For more information, see [Release channels](#release-channels).

## Supported platforms

Docker Engine is available on a variety of Linux platforms, [Mac](/docker-for-mac/install/)
and [Windows](/docker-for-windows/install/) through Docker Desktop, Windows
Server, and as a static binary installation. Find your preferred operating
system below.

#### Desktop

{% assign green-check = '![yes](/install/images/green-check.svg){: style="height: 14px; margin: 0 auto"}' %}

| Platform                                                                    |      x86_64       |
|:----------------------------------------------------------------------------|:-----------------:|
| [Docker Desktop for Mac (macOS)](/docker-for-mac/install/)                        | {{ green-check }} |
| [Docker Desktop for Windows (Microsoft Windows 10)](/docker-for-windows/install/) | {{ green-check }} |

#### Server

{% assign green-check = '![yes](/install/images/green-check.svg){: style="height: 14px; margin: 0 auto"}' %}
{% assign install-prefix-ce = '/install/linux/docker-ce' %}

| Platform                                    | x86_64 / amd64                                         | ARM                                                    | ARM64 / AARCH64                                        | IBM Power (ppc64le)                                    | IBM Z (s390x)                                          |
|:--------------------------------------------|:-------------------------------------------------------|:-------------------------------------------------------|:-------------------------------------------------------|:-------------------------------------------------------|:-------------------------------------------------------|
| [CentOS]({{ install-prefix-ce }}/centos/) | [{{ green-check }}]({{ install-prefix-ce }}/centos/) |                                                        | [{{ green-check }}]({{ install-prefix-ce }}/centos/) |                                                        |                                                        |
| [Debian]({{ install-prefix-ce }}/debian/) | [{{ green-check }}]({{ install-prefix-ce }}/debian/) | [{{ green-check }}]({{ install-prefix-ce }}/debian/) | [{{ green-check }}]({{ install-prefix-ce }}/debian/) |                                                        |                                                        |
| [Fedora]({{ install-prefix-ce }}/fedora/) | [{{ green-check }}]({{ install-prefix-ce }}/fedora/) |                                                        | [{{ green-check }}]({{ install-prefix-ce }}/fedora/) |                                                        |                                                        |
| [Ubuntu]({{ install-prefix-ce }}/ubuntu/) | [{{ green-check }}]({{ install-prefix-ce }}/ubuntu/) | [{{ green-check }}]({{ install-prefix-ce }}/ubuntu/) | [{{ green-check }}]({{ install-prefix-ce }}/ubuntu/) | [{{ green-check }}]({{ install-prefix-ce }}/ubuntu/) | [{{ green-check }}]({{ install-prefix-ce }}/ubuntu/) |

## Release channels

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

> **Note:**
> While pre-releases are done to assist in the stabilization process, no
> guarantees are provided.

Binaries built for pre-releases are available in the test channel for
the targeted year-month release using the naming format `test-YY.mm`,
for example `test-18.09`.

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

The release channel for these builds is called `nightly`.

## Support

Docker Engine releases of a year-month branch are supported with patches as needed for 7 months after the first year-month general availability
release.

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
[Getting started with Docker](/get-started/).
