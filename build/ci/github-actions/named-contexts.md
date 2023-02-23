---
title: Named contexts with GitHub Actions
keywords: ci, github actions, gha, buildkit, buildx, context
---

You can define [additional build contexts](../../../engine/reference/commandline/buildx_build.md#build-context),
and access them in your Dockerfile with `FROM name` or `--from=name`. When
Dockerfile defines a stage with the same name it's overwritten.

This can be useful with GitHub Actions to reuse results from other builds or pin
an image to a specific tag in your workflow.

## Pin image to a tag

Replace `alpine:latest` with a pinned one:

```dockerfile
# syntax=docker/dockerfile:1
FROM alpine
RUN echo "Hello World"
```

{% raw %}
```yaml
name: ci

on:
  push:
    branches:
      - "main"

jobs:
  docker:
    runs-on: ubuntu-latest
    steps:
      -
        name: Checkout
        uses: actions/checkout@v3
      -
        name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
      -
        name: Build
        uses: docker/build-push-action@v4
        with:
          context: .
          build-contexts: |
            alpine=docker-image://alpine:3.16
          tags: myimage:latest
```
{% endraw %}

## Use image in subsequent steps

By default, the [Docker Setup Buildx](https://github.com/marketplace/actions/docker-setup-buildx){:target="blank" rel="noopener" class=""}
action uses `docker-container` as a build driver, so built Docker images aren't
loaded automatically.

With named contexts you can reuse the built image:

```dockerfile
# syntax=docker/dockerfile:1
FROM alpine
RUN echo "Hello World"
```

{% raw %}
```yaml
name: ci

on:
  push:
    branches:
      - "main"

jobs:
  docker:
    runs-on: ubuntu-latest
    steps:
      -
        name: Checkout
        uses: actions/checkout@v3
      -
        name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
      -
        name: Build base image
        uses: docker/build-push-action@v4
        with:
          context: base
          load: true
          tags: my-base-image:latest
      -
        name: Build
        uses: docker/build-push-action@v4
        with:
          context: .
          build-contexts: |
            alpine=docker-image://my-base-image:latest
          tags: myimage:latest
```
{% endraw %}
