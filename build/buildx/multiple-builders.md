---
title: Using multiple builders
description: How to instantiate and work with multiple builders
keywords: build, buildx, buildkit, builders, build drivers
---

By default, Buildx uses the `docker` driver if it is supported, providing a user
experience very similar to the native `docker build`. Note that you must use a
local shared daemon to build your applications.

Buildx allows you to create new instances of isolated builders. You can use this
to get a scoped environment for your CI builds that does not change the state of
the shared daemon, or for isolating builds for different projects. You can create
a new instance for a set of remote nodes, forming a build farm, and quickly
switch between them.

You can create new instances using the [`docker buildx create`](../../engine/reference/commandline/buildx_create.md)
command. This creates a new builder instance with a single node based on your
current configuration.

To use a remote node you can specify the `DOCKER_HOST` or the remote context name
while creating the new builder. After creating a new instance, you can manage its
lifecycle using the [`docker buildx inspect`](../../engine/reference/commandline/buildx_inspect.md),
[`docker buildx stop`](../../engine/reference/commandline/buildx_stop.md), and
[`docker buildx rm`](../../engine/reference/commandline/buildx_rm.md) commands.
To list all available builders, use [`docker buildx ls`](../../engine/reference/commandline/buildx_ls.md).
After creating a new builder you can also append new nodes to it.

To switch between different builders, use [`docker buildx use <name>`](../../engine/reference/commandline/buildx_use.md).
After running this command, the build commands will automatically use this
builder.

Docker also features a [`docker context`](../../engine/reference/commandline/context.md)
command that you can use to provide names for remote Docker API endpoints. Buildx
integrates with `docker context` to ensure all the contexts automatically get a
default builder instance. You can also set the context name as the target when
you create a new builder instance or when you add a node to it.
