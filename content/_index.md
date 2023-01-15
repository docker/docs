---
title: Home
hero:
  image: /assets/images/laptop.svg
  heading: Docker onboarding guide
  copy: |
    Want to use Docker but don't know where to start?
    This guide contains step-by-step instructions on how to get started with Docker.

    [Start with part 1: Overview](./guides/get-started/overview.md)
skip_read_time: true
description: Home page for Docker's documentation
keywords: Docker, documentation, manual, guide, reference, api, samples
layout: wide
notoc: true
grid:
  - title: "Docker Desktop"
    description: |
      Install on [Mac](./desktop/install/mac-install.md), [Windows](./desktop/install/windows-install.md),
      and [Linux](./desktop/install/linux-install.md).

      [Hardened desktop](./desktop/hardened-desktop/_index.md)

      [Extensions (beta)](./desktop/extensions/_index.md)
    icon: "install_desktop"
    link: "/desktop/"
  - title: "Docker Engine"
    description: |
      [Install Docker Engine](./engine/install/_index.md)

      [Configure networking](./engine/network/_index.md)

      [Configure storage](./engine/storage/_index.md)
    icon: "developer_board"
    link: "/engine/"
  - title: "Docker Build"
    description: |
      [Building images](./build/building/_index.md)

      [BuildKit](./build/buildkit/_index.md)
    icon: "build"
    link: "/build/"
  - title: "Docker Compose"
    description: |
      [Install Docker Compose](./compose/install/_index.md)

      [Environment variables in Compose](./compose/environment-variables/_index.md)

      [Networking with Compose](./compose/networking.md)
    icon: "account_tree"
    link: "/compose/"
  - title: "Docker Hub"
    description: |
      [Docker Hub Quickstart](./hub/_index.md)

      [Repositories](./hub/repos/_index.md)

      [Image scanning](./hub/vulnerability-scanning.md)
    icon: "hub"
    link: "/hub/"
  - title: "Business & Teams"
    description: |
      Getting started with [Business](#) or [Teams](#).

      [Manage members](./admin/subscription/_index.md)

      [Single sign-on (SSO)](./admin/single-sign-on/_index.md)
    icon: "admin_panel_settings"
    link: "/admin/"
  - title: "References"
    description: |
      [Docker Engine CLI](./reference/cli/_index.md)

      [Dockerfile reference](./reference/builder.md)

      [Compose specification](./compose/compose-file/_index.md)

      [Docker Extension SDK API](./reference/sdk/_index.md)
    icon: "description"
    link: "/reference/"
---

{{< grid >}}
