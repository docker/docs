---
title: Builders
weight: 40
keywords: build, buildx, builders, buildkit, drivers, backend
description: Learn about builders and how to manage them
---

A builder is a BuildKit daemon that you can use to run your builds. BuildKit
is the build engine that solves the build steps in a Dockerfile to produce a
container image or other artifacts.

You can create and manage builders, inspect them, and even connect to builders
running remotely. You interact with builders using the Docker CLI.

## Default builder

Docker Engine automatically creates a builder that becomes the default backend
for your builds. This builder uses the BuildKit library bundled with the
daemon. This builder requires no configuration.

The default builder is directly bound to the Docker daemon and its
[context](/manuals/engine/manage-resources/contexts.md). If you change the
Docker context, your `default` builder refers to the new Docker context.

## Build drivers

Buildx implements a concept of [build drivers](drivers/_index.md) to refer to
different builder configurations. The default builder created by the daemon
uses the [`docker` driver](drivers/docker.md).

Buildx supports the following build drivers:

- `docker`: uses the BuildKit library bundled into the Docker daemon.
- `docker-container`: creates a dedicated BuildKit container using Docker.
- `kubernetes`: creates BuildKit pods in a Kubernetes cluster.
- `remote`: connects directly to a manually managed BuildKit daemon.

## Selected builder

Selected builder refers to the builder that's used by default when you run
build commands.

When you run a build, or interact with builders in some way using the CLI,
you can use the optional `--builder` flag, or the `BUILDX_BUILDER`
[environment variable](../building/variables.md#buildx_builder),
to specify a builder by name. If you don't specify a builder,
the selected builder is used.

Use the `docker buildx ls` command to see the available builder instances.
The asterisk (`*`) next to a builder name indicates the selected builder.

```console
$ docker buildx ls
NAME/NODE       DRIVER/ENDPOINT      STATUS   BUILDKIT PLATFORMS
default *       docker
  default       default              running  v0.11.6  linux/amd64, linux/amd64/v2, linux/amd64/v3, linux/386
my_builder      docker-container
  my_builder0   default              running  v0.11.6  linux/amd64, linux/amd64/v2, linux/amd64/v3, linux/386
```

### Select a different builder

To switch between builders, use the `docker buildx use <name>` command.

After running this command, the builder you specify is automatically
selected when you invoke builds.

### Difference between `docker build` and `docker buildx build`

Even though `docker build` is an alias for `docker buildx build`, there are
subtle differences between the two commands. With Buildx, the build client and
the and daemon (BuildKit) are decoupled. This means you can use multiple
builders from a single client, even remote ones.

The `docker build` command always defaults to using the default builder that
comes bundled with the Docker Engine, for ensuring backwards compatibility with
older versions of the Docker CLI. The `docker buildx build` command, on the
other hand, checks whether you've set a different builder as the default
builder before it sends your build to BuildKit.

To use the `docker build` command with a non-default builder, you must either:

- Specify the builder explicitly, using the `--builder` flag or the `BUILDX_BUILDER` environment variable:

  ```console
  $ BUILDX_BUILDER=my_builder docker build .
  $ docker build --builder my_builder .
  ```

- Configure Buildx as the default client by running the following command:

  ```console
  $ docker buildx install
  ```

  This updates your [Docker CLI configuration file](/reference/cli/docker/_index.md#configuration-files)
  to ensure all of your build-related commands are routed via Buildx.

  > [!TIP]
  > To undo this change, run `docker buildx uninstall`.

<!-- vale Docker.We = NO -->

In general, we recommend that you use the `docker buildx build` command when
you want to use custom builders. This ensures that your [selected
builder](#selected-builder) configuration is interpreted correctly.

<!-- vale Docker.We = YES -->

## Additional information

- For information about how to interact with and manage builders,
  see [Manage builders](./manage.md)
- To learn about different types of builders,
  see [Build drivers](drivers/_index.md)
