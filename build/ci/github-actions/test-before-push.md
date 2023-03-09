---
title: Test before push with GitHub Actions
keywords: ci, github actions, gha, buildkit, buildx, test
---

In some cases, you might want to validate that the image works as expected
before pushing it.

The following workflow implements several steps to achieve this:

1. Build and export the image to Docker
2. Test your image
3. Multi-platform build and push the image

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
