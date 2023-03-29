---
title: Local registry with GitHub Actions
keywords: ci, github actions, gha, buildkit, buildx, registry
---

For testing purposes you may need to create a [local registry](https://hub.docker.com/_/registry){:target="blank" rel="noopener" class=""}
to push images into:

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
{% endraw %}
