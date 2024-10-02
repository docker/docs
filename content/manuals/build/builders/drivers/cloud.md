---
title: Cloud driver
description: |
  The cloud driver lets you connect to a remote Docker Build Cloud builder.
keywords: build, buildx, driver, builder, remote, cloud
---

The `cloud` driver establishes a connection to a cloud builder, a multi-node,
high-performance remote builder managed through the Docker Build Cloud service.

> [!NOTE]
> - The `cloud` driver requires a Docker Build Cloud subscription. For more
>   information on pricing and availability, visit the Docker Build Cloud
>   [subscriptions and features page](/manuals/subscription/build-cloud/build-details.md).
> - The `cloud` driver is not supported with the open source build of Buildx.
>   You need to use a custom version of Buildx that supports this driver.
>   See [Docker Build Cloud setup](/manuals/build-cloud/setup.md) for details.

## Synopsis

```console
$ docker buildx create \
  --driver cloud \
  <org>/<builder_name>
```

The `cloud` builder does not support any driver options.

## Build cache sharing

Cloud builders are linked to a Docker organization, which means the build cache
is shared among all members using the same builder. If two builds use identical
resources or generate the same image layers, those layers will not be rebuilt.
The cache is created once and automatically reused in all subsequent builds by
any user connected to the builder.

## Multi-player builds

When using the cloud builder with Docker Desktop, build records from the entire
team are automatically available in the [Docker Desktop Builds
view](/manuals/desktop/use-desktop/builds.md). This lets team members directly
inspect logs, cache usage, build errors, and more without needing to manually
share logs.

## Multi-node by default

Cloud builders are provisioned with two high-performance nodes, supporting both
`amd64` and `arm64` architectures. This lets you build multi-platform images at
native speeds, regardless of your local architecture or CI environment.

## Further reading

For more information about Docker Build Cloud, refer to the
[product manual](/manuals/build-cloud/_index.md).
