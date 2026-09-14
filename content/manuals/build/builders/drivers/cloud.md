---
title: Cloud driver
description: The cloud driver connects Buildx to managed builders in Docker Build Cloud.
keywords: build, buildx, driver, builder, cloud, Docker Build Cloud
---

The Buildx cloud driver connects to managed BuildKit instances in Docker Build
Cloud. Docker provisions and maintains the builder infrastructure, shared build
cache, and native `linux/amd64` and `linux/arm64` nodes.

The cloud driver requires Buildx version 0.37.0 or later and access to Docker
Build Cloud.

## Synopsis

Before connecting Buildx, create a cloud builder on the
[Docker Build Cloud Dashboard](https://app.docker.com/build/). Then sign in and
connect to the builder:

```console
$ docker login
$ docker buildx create --driver cloud <ORG>/<BUILDER_NAME>
```

The command adds the cloud builder to your local Buildx configuration. It does
not create a builder in Docker Build Cloud.

For service prerequisites and authentication, see
[Docker Build Cloud setup](/manuals/build-cloud/setup.md). For build examples,
see [Building with Docker Build Cloud](/manuals/build-cloud/usage.md).
