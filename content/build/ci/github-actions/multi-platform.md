---
title: Multi-platform image with GitHub Actions
description: Build for multiple architectures with GitHub Actions using QEMU emulation or multiple native builders
keywords: ci, github actions, gha, buildkit, buildx, multi-platform
---

You can build [multi-platform images](../../building/multi-platform.md) using
the `platforms` option, as shown in the following example:

> **Note**
>
> - For a list of available platforms, see the [Docker Setup Buildx](https://github.com/marketplace/actions/docker-setup-buildx)
>   action.
> - If you want support for more platforms, you can use QEMU with the [Docker Setup QEMU](https://github.com/docker/setup-qemu-action)
>   action.

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
      - name: Checkout
        uses: actions/checkout@v4
      - name: Set up QEMU
        uses: docker/setup-qemu-action@v3
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3
      - name: Login to Docker Hub
        uses: docker/login-action@v3
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}
      - name: Build and push
        uses: docker/build-push-action@v5
        with:
          context: .
          platforms: linux/amd64,linux/arm64
          push: true
          tags: user/app:latest
```

## Distribute build across multiple runners

In the previous example, each platform is built on the same runner which can
take a long time depending on the number of platforms and your Dockerfile.

To solve this issue you can use a matrix strategy to distribute the build for
each platform across multiple runners and create manifest list using the
[`buildx imagetools create` command](../../../engine/reference/commandline/buildx_imagetools_create.md).

The following workflow will build the image for each platform on a dedicated
runner using a matrix strategy and push by digest. Then, the `merge` job will
create a manifest list and push it to Docker Hub.

This example also uses the [`metadata` action](https://github.com/docker/metadata-action)
to set tags and labels.

```yaml
name: ci

on:
  push:
    branches:
      - "main"

env:
  REGISTRY_IMAGE: user/app

jobs:
  build:
    runs-on: ubuntu-latest
    strategy:
      fail-fast: false
      matrix:
        platform:
          - linux/amd64
          - linux/arm/v6
          - linux/arm/v7
          - linux/arm64
    steps:
      -
        name: Prepare
        run: |
          platform=${{ matrix.platform }}
          echo "PLATFORM_PAIR=${platform//\//-}" >> $GITHUB_ENV
      -
        name: Checkout
        uses: actions/checkout@v4
      -
        name: Docker meta
        id: meta
        uses: docker/metadata-action@v5
        with:
          images: ${{ env.REGISTRY_IMAGE }}
      -
        name: Set up QEMU
        uses: docker/setup-qemu-action@v3
      -
        name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3
      -
        name: Login to Docker Hub
        uses: docker/login-action@v3
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}
      -
        name: Build and push by digest
        id: build
        uses: docker/build-push-action@v5
        with:
          context: .
          platforms: ${{ matrix.platform }}
          labels: ${{ steps.meta.outputs.labels }}
          outputs: type=image,name=${{ env.REGISTRY_IMAGE }},push-by-digest=true,name-canonical=true,push=true
      -
        name: Export digest
        run: |
          mkdir -p /tmp/digests
          digest="${{ steps.build.outputs.digest }}"
          touch "/tmp/digests/${digest#sha256:}"
      -
        name: Upload digest
        uses: actions/upload-artifact@v4
        with:
          name: digests-${{ env.PLATFORM_PAIR }}
          path: /tmp/digests/*
          if-no-files-found: error
          retention-days: 1

  merge:
    runs-on: ubuntu-latest
    needs:
      - build
    steps:
      -
        name: Download digests
        uses: actions/download-artifact@v4
        with:
          path: /tmp/digests
          pattern: digests-*
          merge-multiple: true
      -
        name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3
      -
        name: Docker meta
        id: meta
        uses: docker/metadata-action@v5
        with:
          images: ${{ env.REGISTRY_IMAGE }}
      -
        name: Login to Docker Hub
        uses: docker/login-action@v3
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}
      -
        name: Create manifest list and push
        working-directory: /tmp/digests
        run: |
          docker buildx imagetools create $(jq -cr '.tags | map("-t " + .) | join(" ")' <<< "$DOCKER_METADATA_OUTPUT_JSON") \
            $(printf '${{ env.REGISTRY_IMAGE }}@sha256:%s ' *)
      -
        name: Inspect image
        run: |
          docker buildx imagetools inspect ${{ env.REGISTRY_IMAGE }}:${{ steps.meta.outputs.version }}
```

### With Bake

It's also possible to build on multiple runners using Bake, with the
[bake action](https://github.com/docker/bake-action).

You can find a live example [in this GitHub repository](https://github.com/crazy-max/docker-linguist).

The following example achieves the same results as described in
[the previous section](#distribute-build-across-multiple-runners).

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

```yaml
name: ci

on:
  push:
    branches:
      - "main"

env:
  REGISTRY_IMAGE: user/app

jobs:
  prepare:
    runs-on: ubuntu-latest
    outputs:
      matrix: ${{ steps.platforms.outputs.matrix }}
    steps:
      -
        name: Checkout
        uses: actions/checkout@v4
      -
        name: Create matrix
        id: platforms
        run: |
          echo "matrix=$(docker buildx bake image-all --print | jq -cr '.target."image-all".platforms')" >>${GITHUB_OUTPUT}
      -
        name: Show matrix
        run: |
          echo ${{ steps.platforms.outputs.matrix }}
      -
        name: Docker meta
        id: meta
        uses: docker/metadata-action@v5
        with:
          images: ${{ env.REGISTRY_IMAGE }}
      -
        name: Rename meta bake definition file
        run: |
          mv "${{ steps.meta.outputs.bake-file }}" "/tmp/bake-meta.json"
      -
        name: Upload meta bake definition
        uses: actions/upload-artifact@v4
        with:
          name: bake-meta
          path: /tmp/bake-meta.json
          if-no-files-found: error
          retention-days: 1

  build:
    runs-on: ubuntu-latest
    needs:
      - prepare
    strategy:
      fail-fast: false
      matrix:
        platform: ${{ fromJson(needs.prepare.outputs.matrix) }}
    steps:
      -
        name: Prepare
        run: |
          platform=${{ matrix.platform }}
          echo "PLATFORM_PAIR=${platform//\//-}" >> $GITHUB_ENV
      -
        name: Checkout
        uses: actions/checkout@v4
      -
        name: Download meta bake definition
        uses: actions/download-artifact@v4
        with:
          name: bake-meta
          path: /tmp
      -
        name: Set up QEMU
        uses: docker/setup-qemu-action@v3
      -
        name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3
      -
        name: Login to Docker Hub
        uses: docker/login-action@v3
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}
      -
        name: Build
        id: bake
        uses: docker/bake-action@v4
        with:
          files: |
            ./docker-bake.hcl
            /tmp/bake-meta.json
          targets: image
          set: |
            *.tags=
            *.platform=${{ matrix.platform }}
            *.output=type=image,"name=${{ env.REGISTRY_IMAGE }}",push-by-digest=true,name-canonical=true,push=true
      -
        name: Export digest
        run: |
          mkdir -p /tmp/digests
          digest="${{ fromJSON(steps.bake.outputs.metadata).image['containerimage.digest'] }}"
          touch "/tmp/digests/${digest#sha256:}"
      -
        name: Upload digest
        uses: actions/upload-artifact@v4
        with:
          name: digests-${{ env.PLATFORM_PAIR }}
          path: /tmp/digests/*
          if-no-files-found: error
          retention-days: 1

  merge:
    runs-on: ubuntu-latest
    needs:
      - build
    steps:
      -
        name: Download meta bake definition
        uses: actions/download-artifact@v3
        with:
          name: bake-meta
          path: /tmp
      -
        name: Download digests
        uses: actions/download-artifact@v4
        with:
          path: /tmp/digests
          pattern: digests-*
          merge-multiple: true
      -
        name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3
      -
        name: Login to DockerHub
        uses: docker/login-action@v3
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}
      -
        name: Create manifest list and push
        working-directory: /tmp/digests
        run: |
          docker buildx imagetools create $(jq -cr '.target."docker-metadata-action".tags | map(select(startswith("${{ env.REGISTRY_IMAGE }}")) | "-t " + .) | join(" ")' /tmp/bake-meta.json) \
            $(printf '${{ env.REGISTRY_IMAGE }}@sha256:%s ' *)
      -
        name: Inspect image
        run: |
          docker buildx imagetools inspect ${{ env.REGISTRY_IMAGE }}:$(jq -r '.target."docker-metadata-action".args.DOCKER_META_VERSION' /tmp/bake-meta.json)
```
