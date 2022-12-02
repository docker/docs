---
title: "Drivers overview"
keywords: build, buildx, driver, builder, docker-container, kubernetes, remote
redirect_from:
  - /build/buildx/drivers/
---

Buildx drivers are configurations for how and where the BuildKit backend runs.
Driver settings are customizable and allows fine-grained control of the builder.
Buildx supports the following drivers:

- `docker`: uses the BuildKit library bundled into the Docker daemon.
- `docker-container`: creates a dedicated BuildKit container using Docker.
- `kubernetes`: creates BuildKit pods in a Kubernetes cluster.
- `remote`: connects directly to a manually managed BuildKit daemon.

Different drivers support different use cases. The default `docker` driver
prioritizes simplicity and ease of use. It has limited support for advanced
features like caching and output formats, and isn't configurable. Other drivers
provide more flexibility and are better at handling advanced scenarios.

The following table outlines some of the differences between drivers.

| Feature                      |  `docker`   | `docker-container` | `kubernetes` |      `remote`      |
| :--------------------------- | :---------: | :----------------: | :----------: | :----------------: |
| **Automatically load image** |     ✅      |                    |              |                    |
| **Cache export**             | Inline only |         ✅         |      ✅      |         ✅         |
| **Tarball output**           |             |         ✅         |      ✅      |         ✅         |
| **Multi-arch images**        |             |         ✅         |      ✅      |         ✅         |
| **BuildKit configuration**   |             |         ✅         |      ✅      | Managed externally |

## List available builders

Use `docker buildx ls` to see builder instances available on your system, and
the drivers they're using.

```console
$ docker buildx ls
NAME/NODE       DRIVER/ENDPOINT      STATUS   BUILDKIT PLATFORMS
default         docker
  default       default              running  20.10.17 linux/amd64, linux/386
```

Depending on your setup, you may find multiple builders in your list that use
the Docker driver. For example, on a system that runs both a manually installed
version of dockerd, as well as Docker Desktop, you might see the following
output from `docker buildx ls`:

```console
NAME/NODE       DRIVER/ENDPOINT STATUS  BUILDKIT PLATFORMS
default         docker
  default       default         running 20.10.17 linux/amd64, linux/386
desktop-linux * docker
  desktop-linux desktop-linux   running 20.10.17 linux/amd64, linux/arm64, linux/riscv64, linux/ppc64le, linux/s390x, linux/386, linux/arm/v7, linux/arm/v6
```

This is because the Docker driver builders are automatically pulled from the
available [Docker Contexts](../../../engine/context/working-with-contexts/).
When you add new contexts using `docker context create`, these will appear in
your list of buildx builders.

## Create a new builder

Use the
[`docker buildx create`](../../../engine/reference/commandline/buildx_create/)
command to create a builder, and specify the driver using the `--driver` option.

```console
$ docker buildx create --name=<builder-name> --driver=<driver> --driver-opt=<driver-options>
```

## What's next

Read about each of the Buildx drivers to learn about how they work and how to
use them:

- [Docker driver](./docker.md)
- [Docker container driver](./docker-container.md)
- [Kubernetes driver](./kubernetes.md)
- [Remote driver](./remote.md)
