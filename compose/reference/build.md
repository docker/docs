---
description: Build or rebuild services.
keywords: fig, composition, compose, docker, orchestration, cli, build
title: docker-compose build
notoc: true

---

```none
Usage: build [options] [--build-arg key=val...] [SERVICE...]

Options:
    --build-arg key=val     Set build-time variables for services.
    --compress              Compress the build context using gzip.
    --force-rm              Always remove intermediate containers.
    -m, --memory MEM        Set memory limit for the build container.
    --no-cache              Do not use cache when building the image.
    --no-rm                 Do not remove intermediate containers after a successful build.
    --parallel              Build images in parallel.
    --progress string       Set type of progress output (`auto`, `plain`, `tty`).
    --pull                  Always attempt to pull a newer version of the image.
    -q, --quiet             Don't print anything to `STDOUT`.
```

Services are built once and then tagged, by default as `project_service`. For
example, `composetest_db`. If the Compose file specifies an
[image](../compose-file/compose-file-v3.md#image) name, the image is
tagged with that name, substituting any variables beforehand. See
[variable substitution](../compose-file/compose-file-v3.md#variable-substitution).

If you change a service's Dockerfile or the contents of its
build directory, run `docker-compose build` to rebuild it.

## Native build using the docker CLI

Compose by default uses the `docker` CLI to perform builds (also known as "native
build"). By using the `docker` CLI, Compose can take advantage of features such
as [BuildKit](../../develop/develop-images/build_enhancements.md), which are not
supported by Compose itself. BuildKit is enabled by default on Docker Desktop,
but requires the `DOCKER_BUILDKIT=1` environment variable to be set on other
platforms.

Refer to the [Compose CLI environment variables](envvars.md#COMPOSE_DOCKER_CLI_BUILD)
section to learn how to switch between "native build" and "compose build".
