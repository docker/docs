---
title: Buildx
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


## High-level build options

Buildx also aims to provide support for high-level build concepts that go beyond
invoking a single build command.

BuildKit efficiently handles multiple concurrent build requests and de-duplicating
work. The build commands can be combined with general-purpose command runners
(for example, `make`). However, these tools generally invoke builds in sequence
and therefore cannot leverage the full potential of BuildKit parallelization,
or combine BuildKit’s output for the user. For this use case, we have added a
command called [`docker buildx bake`](../../engine/reference/commandline/buildx_bake.md).

The `bake` command supports building images from compose files, similar to
[`docker-compose build`](../../engine/reference/commandline/compose_build.md),
but allowing all the services to be built concurrently as part of a single
request.

There is also support for custom build rules from HCL/JSON files allowing
better code reuse and different target groups. The design of bake is in very
early stages, and we are looking for feedback from users. Let us know your
feedback by creating an issue in the [Docker Buildx](https://github.com/docker/buildx/issues){:target="_blank" rel="noopener" class="_"}
GitHub repository.
