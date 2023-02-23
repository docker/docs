---
title: Example workflows
description: Docker GitHub Actions workflow examples.
keywords: ci, github actions, gha, examples
---

This page showcases different examples of how you can customize and use the
Docker GitHub Actions in your CI pipelines.

## Push to multi-registries

The following workflow will connect you to Docker Hub and [GitHub Container Registry](https://github.com/docker/login-action#github-container-registry){:target="blank" rel="noopener" class=""}
and push the image to both registries:

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
        name: Set up QEMU
        uses: docker/setup-qemu-action@v2
      -
        name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
      -
        name: Login to Docker Hub
        uses: docker/login-action@v2
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}
      -
        name: Login to GitHub Container Registry
        uses: docker/login-action@v2
        with:
          registry: ghcr.io
          username: ${{ github.repository_owner }}
          password: ${{ secrets.GITHUB_TOKEN }}
      -
        name: Build and push
        uses: docker/build-push-action@v4
        with:
          context: .
          platforms: linux/amd64,linux/arm64
          push: true
          tags: |
            user/app:latest
            user/app:1.0.0
            ghcr.io/user/app:latest
            ghcr.io/user/app:1.0.0
```
{% endraw %}

## Manage tags and labels

If you want an "automatic" tag management and [OCI Image Format Specification](https://github.com/opencontainers/image-spec/blob/master/annotations.md){:target="blank" rel="noopener" class=""}
for labels, you can do it in a dedicated setup step. The following workflow
will use the [Docker Metadata Action](https://github.com/docker/metadata-action){:target="blank" rel="noopener" class=""}
to handle tags and labels based on GitHub Actions events and Git metadata:

{% raw %}
```yaml
name: ci

on:
  schedule:
    - cron: "0 10 * * *"
  push:
    branches:
      - "**"
    tags:
      - "v*.*.*"
  pull_request:
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
        name: Docker meta
        id: meta
        uses: docker/metadata-action@v4
        with:
          # list of Docker images to use as base name for tags
          images: |
            name/app
            ghcr.io/username/app
          # generate Docker tags based on the following events/attributes
          tags: |
            type=schedule
            type=ref,event=branch
            type=ref,event=pr
            type=semver,pattern={{version}}
            type=semver,pattern={{major}}.{{minor}}
            type=semver,pattern={{major}}
            type=sha
      -
        name: Set up QEMU
        uses: docker/setup-qemu-action@v2
      -
        name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
      -
        name: Login to Docker Hub
        if: github.event_name != 'pull_request'
        uses: docker/login-action@v2
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}
      -
        name: Login to GHCR
        if: github.event_name != 'pull_request'
        uses: docker/login-action@v2
        with:
          registry: ghcr.io
          username: ${{ github.repository_owner }}
          password: ${{ secrets.GITHUB_TOKEN }}
      -
        name: Build and push
        uses: docker/build-push-action@v4
        with:
          context: .
          push: ${{ github.event_name != 'pull_request' }}
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
```
{% endraw %}

## Export image to Docker

You may want your build result to be available in the Docker client through
`docker images` to be able to use it in another step of your workflow:

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
          load: true
          tags: myimage:latest
      -
        name: Inspect
        run: |
          docker image inspect myimage:latest
```

## Test your image before pushing it

In some cases, you might want to validate that the image works as expected
before pushing it.

The following workflow implements several steps to achieve this:

- Build and export the image to Docker
- Test your image
- Multi-platform build and push the image

{% raw %}
```yaml
name: ci

on:
  push:
    branches:
      - "main"

env:
  TEST_TAG: user/app:test
  LATEST_TAG: user/app:latest

jobs:
  docker:
    runs-on: ubuntu-latest
    steps:
      -
        name: Checkout
        uses: actions/checkout@v3
      -
        name: Set up QEMU
        uses: docker/setup-qemu-action@v2
      -
        name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
      -
        name: Login to Docker Hub
        uses: docker/login-action@v2
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}
      -
        name: Build and export to Docker
        uses: docker/build-push-action@v4
        with:
          context: .
          load: true
          tags: ${{ env.TEST_TAG }}
      -
        name: Test
        run: |
          docker run --rm ${{ env.TEST_TAG }}
      -
        name: Build and push
        uses: docker/build-push-action@v4
        with:
          context: .
          platforms: linux/amd64,linux/arm64
          push: true
          tags: ${{ env.LATEST_TAG }}
```
{% endraw %}

> **Note**
>
> This workflow doesn't actually build the `linux/amd64` image twice. The image
> is built once, and the following steps uses the internal cache for from the
> first `Build and push` step. The second `Build and push` step only builds
> `linux/arm64`.

## Local registry

For testing purposes you may need to create a [local registry](https://hub.docker.com/_/registry){:target="blank" rel="noopener" class=""}
to push images into:

```yaml
name: ci

on:
  push:
    branches:
      - "main"

jobs:
  docker:
    runs-on: ubuntu-latest
    services:
      registry:
        image: registry:2
        ports:
          - 5000:5000
    steps:
      -
        name: Checkout
        uses: actions/checkout@v3
      -
        name: Set up QEMU
        uses: docker/setup-qemu-action@v2
      -
        name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
        with:
          driver-opts: network=host
      -
        name: Build and push to local registry
        uses: docker/build-push-action@v4
        with:
          context: .
          push: true
          tags: localhost:5000/name/app:latest
      -
        name: Inspect
        run: |
          docker buildx imagetools inspect localhost:5000/name/app:latest
```

## Share built image between jobs

As each job is isolated in its own runner, you can't use your built image
between jobs, except if you're using [self-hosted runners](https://docs.github.com/en/actions/hosting-your-own-runners/about-self-hosted-runners){:target="blank" rel="noopener" class=""}
However, you can [pass data between jobs](https://docs.github.com/en/actions/using-workflows/storing-workflow-data-as-artifacts#passing-data-between-jobs-in-a-workflow){:target="blank" rel="noopener" class=""}
in a workflow using the [actions/upload-artifact](https://github.com/actions/upload-artifact){:target="blank" rel="noopener" class=""}
and [actions/download-artifact](https://github.com/actions/download-artifact){:target="blank" rel="noopener" class=""}
actions:

```yaml
name: ci

on:
  push:
    branches:
      - "main"

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      -
        name: Checkout
        uses: actions/checkout@v3
      -
        name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
      -
        name: Build and export
        uses: docker/build-push-action@v4
        with:
          context: .
          tags: myimage:latest
          outputs: type=docker,dest=/tmp/myimage.tar
      -
        name: Upload artifact
        uses: actions/upload-artifact@v3
        with:
          name: myimage
          path: /tmp/myimage.tar

  use:
    runs-on: ubuntu-latest
    needs: build
    steps:
      -
        name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
      -
        name: Download artifact
        uses: actions/download-artifact@v3
        with:
          name: myimage
          path: /tmp
      -
        name: Load image
        run: |
          docker load --input /tmp/myimage.tar
          docker image ls -a
```

## Named contexts

You can define [additional build contexts](../../../engine/reference/commandline/buildx_build.md#build-context),
and access them in your Dockerfile with `FROM name` or `--from=name`. When
Dockerfile defines a stage with the same name it's overwritten.

This can be useful with GitHub Actions to reuse results from other builds or pin
an image to a specific tag in your workflow.

### Pin image to a tag

Replace `alpine:latest` with a pinned one:

```dockerfile
# syntax=docker/dockerfile:1
FROM alpine
RUN echo "Hello World"
```

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

### Use image in subsequent steps

By default, the [Docker Setup Buildx](https://github.com/marketplace/actions/docker-setup-buildx){:target="blank" rel="noopener" class=""}
action uses `docker-container` as a build driver, so built Docker images aren't
loaded automatically.

With named contexts you can reuse the built image:

```dockerfile
# syntax=docker/dockerfile:1
FROM alpine
RUN echo "Hello World"
```

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

## Copy images between registries

[Multi-platform images](../../building/multi-platform.md) built using Buildx can
be copied from one registry to another using the [`buildx imagetools create` command](../../../engine/reference/commandline/buildx_imagetools_create.md):

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
        name: Set up QEMU
        uses: docker/setup-qemu-action@v2
      -
        name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
      -
        name: Login to Docker Hub
        uses: docker/login-action@v2
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}
      -
        name: Login to GitHub Container Registry
        uses: docker/login-action@v2
        with:
          registry: ghcr.io
          username: ${{ github.repository_owner }}
          password: ${{ secrets.GITHUB_TOKEN }}
      -
        name: Build and push
        uses: docker/build-push-action@v4
        with:
          context: .
          platforms: linux/amd64,linux/arm64
          push: true
          tags: |
            user/app:latest
            user/app:1.0.0
      -
        name: Push image to GHCR
        run: |
          docker buildx imagetools create \
            --tag ghcr.io/user/app:latest \
            --tag ghcr.io/user/app:1.0.0 \
            user/app:latest
```
{% endraw %}

## Update Docker Hub repository description

You can update the Docker Hub repository description using a third party action
called [Docker Hub Description](https://github.com/peter-evans/dockerhub-description){:target="blank" rel="noopener" class=""}
with this action:

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
        name: Set up QEMU
        uses: docker/setup-qemu-action@v2
      -
        name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
      -
        name: Login to Docker Hub
        uses: docker/login-action@v2
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}
      -
        name: Build and push
        uses: docker/build-push-action@v4
        with:
          context: .
          push: true
          tags: user/app:latest
      -
        name: Update repo description
        uses: peter-evans/dockerhub-description@v3
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}
          repository: user/app
```
{% endraw %}
