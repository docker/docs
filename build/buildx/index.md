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

You can run Buildx in different configurations that are exposed through a driver
concept. Currently, Docker supports a "docker" driver that uses the BuildKit
library bundled into the Docker daemon binary, and a "docker-container" driver
that automatically launches BuildKit inside a Docker container.

The user experience of using Buildx is very similar across drivers. However,
there are some features that are not currently supported by the "docker" driver,
because the BuildKit library which is bundled into docker daemon uses a different
storage component. In contrast, all images built with the "docker" driver are
automatically added to the "docker images" view by default, whereas when using
other drivers, the method for outputting an image needs to be selected
with `--output`.


## High-level build options with Bake

Check out our guide about [Bake](../bake/index.md) to get started with the
[`docker buildx bake` command](../../engine/reference/commandline/buildx_bake.md).
