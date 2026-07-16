---
title: Build drivers
description: Build drivers are configurations for how and where the BuildKit backend runs.
keywords: build, buildx, driver, builder, docker-container, kubernetes, remote
aliases:
  - /build/buildx/drivers/
  - /build/building/drivers/
  - /build/buildx/multiple-builders/
  - /build/drivers/
---

Build drivers are configurations for how and where the BuildKit backend runs.
Driver settings are customizable and allow fine-grained control of the builder.
Buildx supports the following drivers:

- `docker`: uses the BuildKit library bundled into the Docker daemon.
- `docker-container`: creates a dedicated BuildKit container using Docker.
- `kubernetes`: creates BuildKit pods in a Kubernetes cluster.
- `remote`: connects directly to a manually managed BuildKit daemon.

Different drivers support different use cases. The default `docker` driver
prioritizes simplicity and ease of use. It has limited support for advanced
features like caching and output formats, and isn't configurable. Other drivers
provide more flexibility and are better at handling advanced scenarios.

The following table outlines some differences between drivers.

| Feature                      |  `docker`   | `docker-container` | `kubernetes` |      `remote`      |
| :--------------------------- | :---------: | :----------------: | :----------: | :----------------: |
| **Automatically load image** |     ✅      |                    |              |                    |
| **Cache export**             |     ✅\*     |         ✅         |      ✅      |         ✅         |
| **Tarball output**           |             |         ✅         |      ✅      |         ✅         |
| **Multi-arch images**        |             |         ✅         |      ✅      |         ✅         |
| **BuildKit configuration**   |             |         ✅         |      ✅      | Managed externally |

\* _The `docker` driver doesn't support all cache export options.
See [Cache storage backends](/manuals/build/cache/backends/_index.md) for more information._

## Loading to local image store

Unlike when using the default `docker` driver, images built using other drivers
aren't automatically loaded into the local image store. If you don't specify an
output, the build result is exported to the build cache only.

To build an image using a non-default driver and load it to the image store,
   use the `--load` flag with the build command:

   ```console
   $ docker buildx build --load -t <image> --builder=container .
   ...
   => exporting to oci image format                                                                                                      7.7s
   => => exporting layers                                                                                                                4.9s
   => => exporting manifest sha256:4e4ca161fa338be2c303445411900ebbc5fc086153a0b846ac12996960b479d3                                      0.0s
   => => exporting config sha256:adf3eec768a14b6e183a1010cb96d91155a82fd722a1091440c88f3747f1f53f                                        0.0s
   => => sending tarball                                                                                                                 2.8s
   => importing to docker
   ```

   With this option, the image is available in the image store after the build finishes:

   ```console
   $ docker image ls
   REPOSITORY                       TAG               IMAGE ID       CREATED             SIZE
   <image>                          latest            adf3eec768a1   2 minutes ago       197MB
   ```

### Load by default

{{< summary-bar feature_name="Load by default" >}}

You can configure the custom build drivers to behave in a similar way to the
default `docker` driver, and load images to the local image store by default.
To do so, set the `default-load` driver option when creating the builder:

```console
$ docker buildx create --driver-opt default-load=true
```

Note that, just like with the `docker` driver, if you specify a different
output format with `--output`, the result will not be loaded to the image store
unless you also explicitly specify `--output type=docker` or use the `--load`
flag.

## What's next

Read about each driver:

  - [Docker driver](./docker.md)
  - [Docker container driver](./docker-container.md)
  - [Kubernetes driver](./kubernetes.md)
- [Remote driver](./remote.md)
