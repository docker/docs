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
                            `EXPERIMENTAL` flag for native builder.
                            To enable, run with `COMPOSE_DOCKER_CLI_BUILD=1`)
    --pull                  Always attempt to pull a newer version of the image.
    -q, --quiet             Don't print anything to `STDOUT`.
```

Services are built once and then tagged, by default as `project_service`. For
example, `composetest_db`. If the Compose file specifies an
[image](../compose-file/index.md#image) name, the image is
tagged with that name, substituting any variables beforehand. See
[variable substitution](../compose-file/index.md#variable-substitution).

If you change a service's Dockerfile or the contents of its
build directory, run `docker-compose build` to rebuild it.
