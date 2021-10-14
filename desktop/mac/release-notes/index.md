---
description: Change log / release notes for Docker Desktop for Mac
keywords: Docker Desktop for Mac, release notes
title: Docker Desktop for Mac release notes
toc_min: 1
toc_max: 2
redirect_from:
- /docker-for-mac/release-notes/
- /docker-for-mac/edge-release-notes/
- /mackit/release-notes/
---

> **Update to the Docker Desktop terms**
>
> Professional use of Docker Desktop in large organizations (more than 250 employees or more than $10 million in revenue) requires users to have a paid Docker subscription. While the effective date of these terms is August 31, 2021, there is a grace period until January 31, 2022 for those that require a paid subscription. For more information, see [Docker Desktop License Agreement](../../../subscription/index.md#docker-desktop-license-agreement).
{: .important}

This page contains information about the new features, improvements, known issues, and bug fixes in Docker Desktop releases.

## Docker Desktop 4.1.1
2021-10-12

> Download Docker Desktop
>
> [Mac with Intel chip](https://desktop.docker.com/mac/main/amd64/Docker.dmg?utm_source=docker&utm_medium=webreferral&utm_campaign=docs-driven-download-mac-amd64){: .button .primary-btn }
> [Mac with Apple chip](https://desktop.docker.com/mac/main/arm64/Docker.dmg?utm_source=docker&utm_medium=webreferral&utm_campaign=docs-driven-download-mac-arm64){: .button .primary-btn }

### Bug fixes and minor changes

When upgrading from 4.1.0, **the menu does not change to "Update and restart"** so you can just wait for the download to complete (icon changes) and then select "restart". This bug is fixed (below) in 4.1.1, for future upgrades.

- Fixed a bug where a `Docker.qcow2` file would be ignored on upgrade and a fresh `Docker.raw` used instead, resulting in containers and images disappearing. If a system has both files (due to the previous bug) then the most recently modified file will be used, to avoid recent containers and images disappearing again. To force the use of the old `Docker.qcow2`, delete the newer `Docker.raw` file. Fixes [docker/for-mac#5998](https://github.com/docker/for-mac/issues/5998).
- Fixed update notification overlay sometimes getting out of sync between the Settings button and the Software update button in the Dashboard.
- Fixed menu entry to install a downloaded new Desktop update. When an update is ready to install `Restart` should change to `Update and restart`.

## Docker Desktop 4.1.0
2021-09-30

> Download Docker Desktop
>
> [Mac with Intel chip](https://desktop.docker.com/mac/main/amd64/69386/Docker.dmg) |
> [Mac with Apple chip](https://desktop.docker.com/mac/main/arm64/69386/Docker.dmg)

### New

- **Software Updates**: The Settings tab now includes a new section to help you manage Docker Desktop updates. The **Software Updates** section notifies you whenever there's a new update and allows you to download the update or view information on what's included in the newer version. For more information, see [Software Updates](../index.md#software-updates).
- **Compose V2** You can now specify whether to use [Docker Compose V2](../../../compose/cli-command.md) in the General settings.
- **Volume Management**: Volume management is now available for users on any subscription, including Docker Personal. For more information, see [Explore volumes](../../dashboard.md#explore-volumes).

### Upgrades

- [Compose V2](https://github.com/docker/compose/releases/tag/v2.0.0)
- [Buildx 0.6.3](https://github.com/docker/buildx/releases/tag/v0.6.3)
- [Kubernetes 1.21.5](https://github.com/kubernetes/kubernetes/releases/tag/v1.21.5)
- [Go 1.17.1](https://github.com/golang/go/releases/tag/go1.17.1)
- [Alpine 3.14](https://alpinelinux.org/posts/Alpine-3.14.0-released.html)
- [Qemu 6.1.0](https://wiki.qemu.org/ChangeLog/6.1)
- Base distro to debian:bullseye

## Docker Desktop 4.0.1
2021-09-13

> Download Docker Desktop
>
> [Mac with Intel chip](https://desktop.docker.com/mac/main/amd64/68347/Docker.dmg) |
> [Mac with Apple chip](https://desktop.docker.com/mac/main/arm64/68347/Docker.dmg)

### Upgrades

- [Compose V2 RC3](https://github.com/docker/compose/releases/tag/v2.0.0-rc.3)
  - Compose v2 is now hosted on github.com/docker/compose.
  - Fixed go panic on downscale using `compose up --scale`.
  - Fixed  a race condition in `compose run --rm` while capturing exit code.

### Bug fixes and minor changes

- Fixed a bug where copy-paste was not available in the Docker Dashboard.

## Docker Desktop 4.0.0
2021-08-31

> Download Docker Desktop
>
> [Mac with Intel chip](https://desktop.docker.com/mac/main/amd64/67817/Docker.dmg) |
> [Mac with Apple chip](https://desktop.docker.com/mac/main/arm64/67817/Docker.dmg)

### New

Docker has [announced](https://www.docker.com/blog/updating-product-subscriptions/){: target="*blank" rel="noopener" class="*" id="dkr_docs_relnotes_btl"} updates and extensions to the product subscriptions to increase productivity, collaboration, and added security for our developers and businesses.

The updated [Docker Subscription Service Agreement](https://www.docker.com/legal/docker-subscription-service-agreement) includes a change to the terms for **Docker Desktop**.

- Docker Desktop **remains free** for small businesses (fewer than 250 employees AND less than $10 million in annual revenue), personal use, education, and non-commercial open source projects.
- It requires a paid subscription (**Pro, Team, or Business**), for as little as $5 a month, for professional use in larger enterprises.
- The effective date of these terms is August 31, 2021. There is a grace period until January 31, 2022 for those that will require a paid subscription to use Docker Desktop.
- The Docker Pro and Docker Team subscriptions now **include commercial use** of Docker Desktop.
- The existing Docker Free subscription has been renamed **Docker Personal**.
- **No changes** to Docker Engine or any other upstream **open source** Docker or Moby project.

To understand how these changes affect you, read the [FAQs](https://www.docker.com/pricing/faq){: target="*blank" rel="noopener" class="*" id="dkr_docs_relnotes_btl"}. For more information, see [Docker subscription overview](../../../subscription/index.md).

### Upgrades

- [Compose V2 RC2](https://github.com/docker/compose-cli/releases/tag/v2.0.0-rc.2)
  - Fixed project name to be case-insensitive for `compose down`. See [docker/compose-cli#2023](https://github.com/docker/compose-cli/issues/2023)
  - Fixed non-normalized project name.
  - Fixed port merging on partial reference.
- [Kubernetes 1.21.4](https://github.com/kubernetes/kubernetes/releases/tag/v1.21.4)

### Bug fixes and minor changes

- Fixed a bug where SSH was not available for builds from git URL. Fixes [for-mac#5902](https://github.com/docker/for-mac/issues/5902)
