---
title: Home
hero:
  image: /assets/images/laptop.svg
  heading: Docker documentation
  copy: |
    Want to use Docker but don't know where to start?
    This guide contains step-by-step instructions on how to get started with Docker.

    {{< button type=link url=/get-started text="Start with part 1: Overview" >}}
skip_read_time: true
description: Home page for Docker's documentation
keywords: Docker, documentation, manual, guide, reference, api, samples
layout: wide
notoc: true
grid:
  - title: Download and install
    description: Download and install Docker on your machine in a few, easy steps.
    icon: downloading
    link: /get-docker
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

      [Configure networking](./network/_index.md)

      [Configure storage](./storage/_index.md)
    icon: "developer_board"
    link: "/engine/"
  - title: "Docker Build"
    description: |
      [Building images](./build/building/packaging.md)

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
      [Docker Hub Quickstart](./docker-hub/_index.md)

      [Repositories](./docker-hub/repos/_index.md)

      [Image scanning](./docker-hub/vulnerability-scanning.md)
    icon: "hub"
    link: "/docker-hub/"
  - title: "References"
    description: |
      [Docker Engine CLI](./engine/reference/commandline/cli.md)

      [Dockerfile reference](./engine/reference/builder.md)

      [Compose specification](./compose/compose-file/_index.md)

      [Docker Extension SDK API](./desktop/extensions-sdk/dev/api/overview.md)
    icon: "description"
    link: "/reference/"
  - title: Get support
    description: Find information on how to get support, and the scope of Docker support.
    icon: contact_support
    link: /support
  - title: Release notes
    description: Find out what’s new in Docker!
    icon: auto_awesome
    link: /release-notes
---

{{< grid >}}
