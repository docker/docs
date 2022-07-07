---
title: Overview of Docker Build
description: Introduction and overview of Docker Build
keywords: build, buildx, buildkit
---

Docker Build is one of the most used features of the Docker Engine - users
ranging from developers, build teams, and release teams all use Docker Build.
It uses a [client-server architecture](../get-started/overview.md#docker-architecture)
that includes several tools. The most common method is to use the Docker CLI with
[`docker build` command](../engine/reference/commandline/build.md) that sends
requests to the Docker Engine that will execute your build.

Starting with version 18.09, Docker supports a new backend for executing your
builds that is provided by the [BuildKit](https://github.com/moby/buildkit)
project. The BuildKit backend provides many benefits compared to the old
implementation. As there is a new backend, there is also a new client called [Docker Buildx](../buildx/working-with-buildx.md),
available as a CLI plugin that extends the docker command with the full
support of the features provided by BuildKit.

BuildKit is enabled by default for all users on [Docker Desktop](../desktop/index.md).
If you have installed Docker Desktop, you don't have to manually enable BuildKit.
If you have installed Docker as a [Linux package](../engine/install/index.md),
you can enable BuildKit either by using an environment variable or by making
BuildKit the default setting.

To set the BuildKit environment variable when running the
[`docker build` command](../engine/reference/commandline/build.md), run:

```console
$ DOCKER_BUILDKIT=1 docker build .
```

To enable BuildKit backend by default, set [daemon configuration](/engine/reference/commandline/dockerd/#daemon-configuration-file)
in `/etc/docker/daemon.json` feature to `true` and restart the daemon. If the
`daemon.json` file doesn't exist, create new file called `daemon.json` and then
add the following to the file:

```json
{
  "features": {
    "buildkit": true
  }
}
```

If you're using the [`docker buildx build` command](../engine/reference/commandline/buildx_build.md),
BuildKit will always being used regardless of the environment variable or backend
configuration. See [Build with Buildx](../buildx/working-with-buildx.md#build-with-buildx) guide
for more details.
