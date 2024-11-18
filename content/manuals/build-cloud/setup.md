---
title: Docker Build Cloud setup
linkTitle: Setup
weight: 10
description: How to get started with Docker Build Cloud
keywords: build, cloud build
aliases:
  - /build/cloud/setup/
---

Before you can start using Docker Build Cloud, you must add the builder to your local
environment.

## Prerequisites

To get started with Docker Build Cloud, you need to:

- Download and install Docker Desktop version 4.26.0 or later.
- Sign up for a Docker Build Cloud subscription in the [Docker Build Cloud Dashboard](https://app.docker.com/build/).

### Use Docker Build Cloud without Docker Desktop

To use Docker Build Cloud without Docker Desktop, you must download and install
a version of Buildx with support for Docker Build Cloud (the `cloud` driver).
You can find compatible Buildx binaries on the releases page of
[this repository](https://github.com/docker/buildx-desktop).

If you plan on building with Docker Build Cloud using the `docker compose
build` command, you also need a version of Docker Compose that supports Docker
Build Cloud. You can find compatible Docker Compose binaries on the releases
page of [this repository](https://github.com/docker/compose-desktop).

## Steps

You can add a cloud builder using the CLI, with the `docker buildx create`
command, or using the Docker Desktop settings GUI.

{{< tabs >}}
{{< tab name="CLI" >}}

1. Sign in to your Docker account.

   ```console
   $ docker login
   ```

2. Add the cloud builder endpoint.

   ```console
   $ docker buildx create --driver cloud <ORG>/<BUILDER_NAME>
   ```

   Replace `ORG` with the Docker Hub namespace of your Docker organization.

This creates a builder named `cloud-ORG-BUILDER_NAME`.

{{< /tab >}}
{{< tab name="Docker Desktop" >}}

1. Sign in to your Docker account using the **Sign in** button in Docker Desktop.

2. Open the Docker Desktop settings and navigate to the **Builders** tab.

3. Under **Available builders**, select **Connect to builder**.

{{< /tab >}}
{{< /tabs >}}

The builder has native support for the `linux/amd64` and `linux/arm64`
architectures. This gives you a high-performance build cluster for building
multi-platform images natively.

## Firewall configuration

To use Docker Build Cloud behind a firewall, ensure that your firewall allows
traffic to the following addresses:

- 3.211.38.21
- https://auth.docker.io
- https://build-cloud.docker.com
- https://hub.docker.com

## What's next

- See [Building with Docker Build Cloud](usage.md) for examples on how to use Docker Build Cloud.
- See [Use Docker Build Cloud in CI](ci.md) for examples on how to use Docker Build Cloud with CI systems.
