---
title: Multi-platform image with GitHub Actions
linkTitle: Multi-platform image
description: Build multi-platform images with Docker GitHub Builder or individual actions using QEMU.
keywords: ci, github actions, gha, buildkit, buildx, multi-platform
---

Use [Docker GitHub Builder](github-builder/_index.md) to build multi-platform
images across native runners and publish a combined image manifest.

## Distribute build across multiple runners

Building multiple platforms on the same runner can significantly extend build
times, particularly when dealing with complex Dockerfiles or a high number of
target platforms. If you want to split platform builds across multiple runners
without maintaining a custom matrix and merge job, use the
[Docker GitHub Builder](github-builder/_index.md). The reusable workflows
compute the per-platform matrix, run each platform on its own runner, and
create the final manifest for you.

The following workflow uses the [`build.yml` reusable workflow](github-builder/build.md)
to distribute a multi-platform Dockerfile build:

```yaml
name: ci

on:
  push:

permissions:
  contents: read

jobs:
  build:
    uses: docker/github-builder/.github/workflows/build.yml@{{% param "github_builder_version" %}}
    permissions:
      contents: read
      id-token: write
    with:
      output: image
      push: true
      platforms: linux/amd64,linux/arm64
      meta-images: user/app
      meta-tags: |
        type=ref,event=branch
        type=ref,event=pr
        type=semver,pattern={{version}}
        type=semver,pattern={{major}}.{{minor}}
    secrets:
      registry-auths: |
        - registry: docker.io
          username: ${{ vars.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}
```

With `distribute: true`, which is the default, the workflow splits the build
into one platform per runner and assembles the final multi-platform image in
its finalize phase. The default `runner` mapping selects native GitHub-hosted
Linux runners for Arm and other platforms. To customize the mapping, see [runner selection](github-builder/architecture.md#runner-selection).
If you need to control the Docker build inputs directly, see
[Build with Docker GitHub Builder build.yml](github-builder/build.md).

### With Bake

For builds defined in a Bake file, use the
[Bake reusable workflow](github-builder/bake.md#build-and-push-a-bake-target).
It distributes the target's platforms and publishes the combined manifest.

## Build with individual actions and QEMU

Use individual actions when you need to configure the job's Docker Engine or
run other steps alongside the build. The following example uses QEMU to build
both platforms on one runner.

You can build [multi-platform images](../../building/multi-platform.md) using
the `platforms` option, as shown in the following example:

> [!NOTE]
>
> - For a list of available platforms, see the [Docker Setup Buildx](https://github.com/marketplace/actions/docker-setup-buildx)
>   action.
> - If you want support for more platforms, you can use QEMU with the [Docker Setup QEMU](https://github.com/docker/setup-qemu-action)
>   action.

```yaml
name: ci

on:
  push:

jobs:
  docker:
    runs-on: ubuntu-latest
    steps:
      - name: Login to Docker Hub
        uses: docker/login-action@{{% param "login_action_version" %}}
        with:
          username: ${{ vars.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}

      - name: Set up QEMU
        uses: docker/setup-qemu-action@{{% param "setup_qemu_action_version" %}}

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@{{% param "setup_buildx_action_version" %}}

      - name: Build and push
        uses: docker/build-push-action@{{% param "build_push_action_version" %}}
        with:
          platforms: linux/amd64,linux/arm64
          push: true
          tags: user/app:latest
```

## Build and load multi-platform images

The default Docker setup for GitHub Actions runners supports building and pushing multi-platform images to registries. However, it does not support loading multi-platform images to the local image store of the runner after building them. To load a multi-platform image locally, you need to enable the containerd image store option for the Docker Engine.

There is no way to configure the default Docker setup in the GitHub Actions
runners directly, but you can use `docker/setup-docker-action` to customize the
Docker Engine and CLI settings for a job.

The following example workflow enables the containerd image store, builds a
multi-platform image, and loads the results into the GitHub runner's local
image store.

```yaml
name: ci

on:
  push:

jobs:
  docker:
    runs-on: ubuntu-latest
    steps:
      - name: Set up Docker
        uses: docker/setup-docker-action@{{% param "setup_docker_action_version" %}}
        with:
          daemon-config: |
            {
              "debug": true,
              "features": {
                "containerd-snapshotter": true
              }
            }

      - name: Login to Docker Hub
        uses: docker/login-action@{{% param "login_action_version" %}}
        with:
          username: ${{ vars.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}

      - name: Set up QEMU
        uses: docker/setup-qemu-action@{{% param "setup_qemu_action_version" %}}

      - name: Build and push
        uses: docker/build-push-action@{{% param "build_push_action_version" %}}
        with:
          platforms: linux/amd64,linux/arm64
          load: true
          tags: user/app:latest
```
