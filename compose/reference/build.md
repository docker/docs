---
description: Build or rebuild services.
keywords: fig, composition, compose, docker, orchestration, cli, build
title: docker-compose build
notoc: true

---

```
Usage: build [options] [--build-arg key=val...] [SERVICE...]

Options:
    --compress              Compress the build context using gzip.
    --force-rm              Always remove intermediate containers.
    --no-cache              Do not use cache when building the image.
    --pull                  Always attempt to pull a newer version of the image.
    -m, --memory MEM        Sets memory limit for the build container.
    --build-arg key=val     Set build-time variables for services.
    --parallel              Build images in parallel.
```

Services are built once and then tagged, by default as `project_service`. For
example, `composetest_db`. If the Compose file specifies an
[image](../compose-file/index.md#image) name, the image is
tagged with that name, substituting any variables beforehand. See
[variable substitution](../compose-file/index.md#variable-substitution).

If you change a service's Dockerfile or the contents of its
build directory, run `docker-compose build` to rebuild it.
