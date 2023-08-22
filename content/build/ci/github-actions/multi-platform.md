---
title: Multi-platform image with GitHub Actions
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
        name: Checkout
        uses: actions/checkout@v3
      -
        name: Docker meta
        id: meta
        uses: docker/metadata-action@v4
        with:
          images: ${{ env.REGISTRY_IMAGE }}
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
        name: Build and push by digest
        id: build
        uses: docker/build-push-action@v4
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
        uses: actions/upload-artifact@v3
        with:
          name: digests
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
        uses: actions/download-artifact@v3
        with:
          name: digests
          path: /tmp/digests
      -
        name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
      -
        name: Docker meta
        id: meta
        uses: docker/metadata-action@v4
        with:
          images: ${{ env.REGISTRY_IMAGE }}
      -
        name: Login to Docker Hub
        uses: docker/login-action@v2
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