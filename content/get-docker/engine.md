---
title: Quick install Docker Engine
linkTitle: Docker Engine
description: Install Docker Engine on Linux with the convenience script and run a test container.
keywords: Docker Engine, quick install, Linux, convenience script
---

Docker Engine runs containers on Linux without a desktop application. For a
Mac or Windows computer, use [Docker Desktop](/get-docker/desktop.md).

## Install

Use Docker's convenience script on a development or test machine running a
[supported Linux distribution](/manuals/engine/install/_index.md):

```console
$ curl -fsSL https://get.docker.com | sudo sh
```

The script installs Docker Engine and its dependencies. For production systems,
use the package repository instructions for your distribution instead.

## Check the installation

Run a test container:

```console
$ sudo docker run hello-world
```

The container prints a confirmation message when the installation works.

## Full installation instructions

{{< card
  title="Install Docker Engine"
  description="Choose a Linux distribution for system requirements, package installation, version selection, and upgrades."
  link="/engine/install/" >}}

For details about the script, see
[Install using the convenience script](/manuals/engine/install/ubuntu.md#install-using-the-convenience-script).
