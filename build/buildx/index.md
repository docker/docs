---
title: Working with Buildx
description: Working with Docker Buildx
keywords: build, buildx, buildkit
redirect_from:
- /buildx/working-with-buildx/
---

## Overview

Docker Buildx is a CLI plugin that extends the docker command with the full
support of the features provided by [Moby BuildKit](https://github.com/moby/buildkit){:target="_blank" rel="noopener" class="_"}
builder toolkit. It provides the same user experience as docker build with many
new features like creating scoped builder instances and building against
multiple nodes concurrently.

## Build with Buildx

To start a new build, run the command `docker buildx build .`

```console
$ docker buildx build .
[+] Building 8.4s (23/32)
 => ...
```

Buildx builds using the BuildKit engine and does not require `DOCKER_BUILDKIT=1`
environment variable to start the builds.

The [`docker buildx build` command](../../engine/reference/commandline/buildx_build.md)
supports features available for `docker build`, including features such as
outputs configuration, inline build caching, and specifying target platform.
In addition, Buildx also supports new features that are not yet available for
regular `docker build` like building manifest lists, distributed caching, and
exporting build results to OCI image tarballs.

Buildx is flexible and can be run in different configurations that are exposed
through various "drivers". Each driver defines how and where a build should
run, and have different feature sets.

We currently support the following drivers:

* The `docker` driver ([guide](drivers/docker.md), [reference](/engine/reference/commandline/buildx_create/#driver))
* The `docker-container` driver ([guide](drivers/docker-container.md), [reference](/engine/reference/commandline/buildx_create/#driver))
* The `kubernetes` driver ([guide](drivers/kubernetes.md), [reference](/engine/reference/commandline/buildx_create/#driver))
* The `remote` driver ([guide](drivers/remote.md))

For more information on drivers, see the [drivers guide](drivers/index.md).

## High-level build options with Bake

Check out our guide about [Bake](../bake/index.md) to get started with the
[`docker buildx bake` command](../../engine/reference/commandline/buildx_bake.md).
