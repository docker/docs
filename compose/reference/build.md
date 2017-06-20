---
description: docker-compose build
keywords: fig, composition, compose, docker, orchestration, cli, build
title: docker-compose build
notoc: true

---

```
Usage: build [options] [--build-arg key=val...] [SERVICE...]

Options:
    --force-rm              Always remove intermediate containers.
    --no-cache              Do not use cache when building the image.
    --pull                  Always attempt to pull a newer version of the image.
    --build-arg key=val     Set build-time variables for one service.
```

Services are built once and then tagged, by default as `project_service`, e.g.,
`composetest_db`. If the Compose file specifies an
[image](/compose/compose-file/index.md#image) name, the image will be
tagged with that name, substituting any variables beforehand. See [variable
substitution](#variable-substitution)

If you change a service's Dockerfile or the contents of its
build directory, run `docker-compose build` to rebuild it.
