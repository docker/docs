---
title: Copy image between registries with GitHub Actions
linkTitle: Copy image between registries
description: Build multi-platform images and copy them between registries with GitHub Actions
keywords: ci, github actions, gha, buildkit, buildx, registry
---

[Multi-platform images](../../building/multi-platform.md) built using Buildx can
be copied from one registry to another using the [`buildx imagetools create` command](/reference/cli/docker/buildx/imagetools/create/):

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

      - name: Login to GitHub Container Registry
        uses: docker/login-action@{{% param "login_action_version" %}}
        with:
          registry: ghcr.io
          username: ${{ github.repository_owner }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Set up QEMU
        uses: docker/setup-qemu-action@{{% param "setup_qemu_action_version" %}}
      
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@{{% param "setup_buildx_action_version" %}}

      - name: Build and push
        uses: docker/build-push-action@{{% param "build_push_action_version" %}}
        with:
          platforms: linux/amd64,linux/arm64
          push: true
          tags: |
            user/app:latest
            user/app:1.0.0

      - name: Push image to GHCR
        run: |
          docker buildx imagetools create \
            --tag ghcr.io/user/app:latest \
            --tag ghcr.io/user/app:1.0.0 \
            user/app:latest
```
