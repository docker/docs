---
title: Multi-platform image with GitHub Actions
linkTitle: Multi-platform image
description: Build for multiple architectures with GitHub Actions using QEMU emulation or multiple native builders
keywords: ci, github actions, gha, buildkit, buildx, multi-platform
---

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
        uses: docker/login-action@v3
        with:
          username: ${{ vars.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}

      - name: Set up QEMU
        uses: docker/setup-qemu-action@v3

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Build and push
        uses: docker/build-push-action@v6
        with:
          platforms: linux/amd64,linux/arm64
          push: true
          tags: user/app:latest
```

## Build and load multi-platform images

The default Docker setup for GitHub Actions runners does not support loading
multi-platform images to the local image store of the runner after building
them. To load a multi-platform image, you need to enable the containerd image
store option for the Docker Engine.

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
        uses: docker/setup-docker-action@v4
        with:
          daemon-config: |
            {
              "debug": true,
              "features": {
                "containerd-snapshotter": true
              }
            }

      - name: Login to Docker Hub
        uses: docker/login-action@v3
        with:
          username: ${{ vars.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}

      - name: Set up QEMU
        uses: docker/setup-qemu-action@v3

      - name: Build and push
        uses: docker/build-push-action@v6
        with:
          platforms: linux/amd64,linux/arm64
          load: true
          tags: user/app:latest
```

## Distribute multi-platform build across runners

In the previous example, each platform is built on the same runner which can
take a long time depending on the number of platforms and your Dockerfile.

To solve this issue, you can use the following reusable workflows for both
[`docker/build-push-action`](https://github.com/crazy-max/.github?tab=readme-ov-file#build-distribute-mp)
and [`docker/bake-action`](https://github.com/crazy-max/.github?tab=readme-ov-file#bake-distribute-mp)
actions to distribute multi-platform builds across runners efficiently.

{{< tabs >}}
{{< tab name="build-push-action" >}}

```yaml {hl_lines=9}
name: ci

on:
  push:
  pull_request:

jobs:
  build:
    uses: crazy-max/.github/.github/workflows/build-distribute-mp.yml@main
    with:
      push: ${{ github.event_name != 'pull_request' }}
      cache: true
      meta-image: user/app
      build-platforms: linux/amd64,linux/arm64
      login-username: ${{ vars.DOCKERHUB_USERNAME }}
    secrets:
      login-password: ${{ secrets.DOCKERHUB_TOKEN }}
```

Here are the main inputs for this reusable workflow:

| Name              | Type     | Default | Description                                                                                                                                                                                                                                                |
|-------------------|----------|---------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `runner`          | String   | `auto`¹ | Runner instance (e.g., `ubuntu-latest`).                                                                                                                                                                                                                   |
| `push`            | Bool     | `false` | Push image to registry.                                                                                                                                                                                                                                    |
| `cache`           | Bool     | `false` | Enable GitHub Actions cache backend.                                                                                                                                                                                                                       |
| `cache-scope`     | String   |         | Which scope GitHub Actions cache object belongs to if `cache` enabled.                                                                                                                                                                                     |
| `cache-mode`      | String   | `min`   | Cache layers to export if `cache` enabled (one of `min` or `max`).                                                                                                                                                                                         |
| `summary`         | Bool     | `true`  | Enable [build summary](https://docs.docker.com/build/ci/github-actions/build-summary/) generation.                                                                                                                                                         |
| `meta-image`      | String   |         | Image to use as base name for tags. This input is similar to [`images` input in `docker/metadata-action`](https://github.com/docker/metadata-action?tab=readme-ov-file#images-input) used in this reusable workflow but accepts a single image name.       |
| `build-platforms` | List/CSV |         | List of target platforms for build. This input is similar to [`platforms` input in `docker/build-push-action`](https://github.com/docker/build-push-action?tab=readme-ov-file#inputs) used in this reusable workflow. At least two platforms are required. |
| `login-registry`  | String   |         | Server address of Docker registry. If not set then will default to Docker Hub. This input is similar to [`registry` input in `docker/login-action`](https://github.com/docker/login-action?tab=readme-ov-file#inputs) used in this reusable workflow.      |
| `login-username`² | String   |         | Username used to log against the Docker registry. This input is similar to [`username` input in `docker/login-action`](https://github.com/docker/login-action?tab=readme-ov-file#inputs) used in this reusable workflow.                                   |
| `login-password`  | String   |         | Specifies whether the given registry is ECR (auto, true or false). This input is similar to [`password` input in `docker/login-action`](https://github.com/docker/login-action?tab=readme-ov-file#inputs) used in this reusable workflow.                  |

> [!NOTE]
> ¹ `auto` will choose the best matching runner depending on the target
> platform being built (either `ubuntu-latest` or `ubuntu-24.04-arm`).
> 
> ² `login-username` can be used as either an input or secret.

You can find the list of available inputs directly in [the reusable workflow](https://github.com/crazy-max/.github/blob/main/.github/workflows/build-distribute-mp.yml).

{{< /tab >}}
{{< tab name="bake-action" >}}

```hcl
variable "DEFAULT_TAG" {
  default = "app:local"
}

// Special target: https://github.com/docker/metadata-action#bake-definition
target "docker-metadata-action" {
  tags = ["${DEFAULT_TAG}"]
}

// Default target if none specified
group "default" {
  targets = ["image-local"]
}

target "image" {
  inherits = ["docker-metadata-action"]
}

target "image-local" {
  inherits = ["image"]
  output = ["type=docker"]
}

target "image-all" {
  inherits = ["image"]
  platforms = [
    "linux/amd64",
    "linux/arm/v6",
    "linux/arm/v7",
    "linux/arm64"
  ]
}
```

```yaml {hl_lines=9}
name: ci

on:
  push:
  pull_request:

jobs:
  build:
    uses: crazy-max/.github/.github/workflows/build-distribute-mp.yml@main
    with:
      target: image-all
      push: ${{ github.event_name != 'pull_request' }}
      cache: true
      meta-image: user/app
      login-username: ${{ vars.DOCKERHUB_USERNAME }}
    secrets:
      login-password: ${{ secrets.DOCKERHUB_TOKEN }}
```

Here are the main inputs for this reusable workflow:

| Name              | Type   | Default | Description                                                                                                                                                                                                                                           |
|-------------------|--------|---------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `runner`          | String | `auto`¹ | Runner instance (e.g., `ubuntu-latest`).                                                                                                                                                                                                              |
| `target`          | String |         | Multi-platform target to build. This input is similar to [`targets` input in `docker/bake-action`](https://github.com/docker/build-push-action?tab=readme-ov-file#inputs) used in this reusable workflow but accepts a single target.                 |
| `push`            | Bool   | `false` | Push image to registry.                                                                                                                                                                                                                               |
| `cache`           | Bool   | `false` | Enable GitHub Actions cache backend.                                                                                                                                                                                                                  |
| `cache-scope`     | String |         | Which scope GitHub Actions cache object belongs to if `cache` enabled.                                                                                                                                                                                |
| `cache-mode`      | String | `min`   | Cache layers to export if `cache` enabled (one of `min` or `max`).                                                                                                                                                                                    |
| `summary`         | Bool   | `true`  | Enable [build summary](https://docs.docker.com/build/ci/github-actions/build-summary/) generation.                                                                                                                                                    |
| `meta-image`      | String |         | Image to use as base name for tags. This input is similar to [`images` input in `docker/metadata-action`](https://github.com/docker/metadata-action?tab=readme-ov-file#images-input) used in this reusable workflow but accepts a single image name.  |
| `login-registry`  | String |         | Server address of Docker registry. If not set then will default to Docker Hub. This input is similar to [`registry` input in `docker/login-action`](https://github.com/docker/login-action?tab=readme-ov-file#inputs) used in this reusable workflow. |
| `login-username`² | String |         | Username used to log against the Docker registry. This input is similar to [`username` input in `docker/login-action`](https://github.com/docker/login-action?tab=readme-ov-file#inputs) used in this reusable workflow.                              |
| `login-password`  | String |         | Specifies whether the given registry is ECR (auto, true or false). This input is similar to [`password` input in `docker/login-action`](https://github.com/docker/login-action?tab=readme-ov-file#inputs) used in this reusable workflow.             |

> [!NOTE]
> ¹ `auto` will choose the best matching runner depending on the target
> platform being built (either `ubuntu-latest` or `ubuntu-24.04-arm`).
> 
> ² `login-username` can be used as either an input or secret.

You can find the list of available inputs directly in [the reusable workflow](https://github.com/crazy-max/.github/blob/main/.github/workflows/bake-distribute-mp.yml).

{{< /tab >}}
{{< /tabs >}}
