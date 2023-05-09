---
title: Multi-platform image with GitHub Actions
keywords: ci, github actions, gha, buildkit, buildx, multi-platform
---

You can build [multi-platform images](../../building/multi-platform.md) using
the `platforms` option, as shown in the following example:

> **Note**
>
> - For a list of available platforms, see the [Docker Setup Buildx](https://github.com/marketplace/actions/docker-setup-buildx){:target="blank" rel="noopener" class=""}
>   action.
> - If you want support for more platforms, you can use QEMU with the [Docker Setup QEMU](https://github.com/docker/setup-qemu-action){:target="blank" rel="noopener" class=""}
>   action.

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
          platforms: linux/amd64,linux/arm64
          push: true
          tags: user/app:latest
```
{% endraw %}

## Distribute build across multiple runners

In the previous example, each platform is built on the same runner which can
take a long time depending on the number of platforms and your Dockerfile.

To solve this issue you can use a matrix strategy to distribute the build for
each platform across multiple runners and merge the resulting manifests into a
single image using the [`buildx imagetools create` command](../../../engine/reference/commandline/buildx_imagetools_create.md).

The following workflow will build the image for each platform on a dedicated
runner using a matrix strategy and upload the resulting images as artifacts.
Then, the `push` job will download the images, load them into the local Docker
daemon, push them to a temporary local registry so Buildx can create a manifest
list and push it to Docker Hub. It will copy the images from the local
registry to Docker Hub as well.

Using a local registry is useful, so we don't create unecessary tags for each
platform on the remote registry. Only `user/app:latest` will be created.

{% raw %}
```yaml
name: ci

on:
  push:
    branches:
      - "main"

env:
  TMP_LOCAL_IMAGE: localhost:5000/user/app
  REGISTRY_IMAGE: user/app
  REGISTRY_TAG: latest

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
        name: Prepare
        run: |
          mkdir -p /tmp/images
          platform=${{ matrix.platform }}
          echo "TARFILE=${platform//\//-}.tar" >> $GITHUB_ENV
          echo "TAG=${{ env.TMP_LOCAL_IMAGE }}:${platform//\//-}" >> $GITHUB_ENV
      -
        name: Set up QEMU
        uses: docker/setup-qemu-action@v2
      -
        name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
      -
        name: Build
        uses: docker/build-push-action@v4
        with:
          context: .
          platforms: ${{ matrix.platform }}
          tags: ${{ env.TAG }}
          outputs: type=docker,dest=/tmp/images/${{ env.TARFILE }}
      -
        name: Upload image
        uses: actions/upload-artifact@v3
        with:
          name: images
          path: /tmp/images/${{ env.TARFILE }}
          if-no-files-found: error
          retention-days: 1
  
  push:
    runs-on: ubuntu-latest
    needs:
      - build
    services:
      registry:
        image: registry:2
        ports:
          - 5000:5000
    steps:
      -
        name: Download images
        uses: actions/download-artifact@v3
        with:
          name: images
          path: /tmp/images
      -
        name: Load images
        run: |
          for image in /tmp/images/*.tar; do
            docker load -i $image
          done
      -
        name: Push images to local registry
        run: |
          docker push -a ${{ env.TMP_LOCAL_IMAGE }}
      -
        name: Login to Docker Hub
        uses: docker/login-action@v2
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}
      -
        name: Create manifest list and push
        run: |
          docker buildx imagetools create -t ${{ env.REGISTRY_IMAGE }}:${{ env.REGISTRY_TAG }} \
            $(docker image ls --format '{{.Repository}}:{{.Tag}}' '${{ env.TMP_LOCAL_IMAGE }}' | tr '\n' ' ')
      -
        name: Inspect image
        run: |
          docker buildx imagetools inspect ${{ env.REGISTRY_IMAGE }}:${{ env.REGISTRY_TAG }}
```
{% endraw %}
