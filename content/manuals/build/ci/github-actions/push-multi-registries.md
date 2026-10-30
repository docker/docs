---
title: Push to multiple registries with GitHub Actions
linkTitle: Push to multiple registries
description: Push to multiple registries with GitHub Actions
keywords: ci, github actions, gha, buildkit, buildx, registry
---

Use [Docker GitHub Builder](github-builder/build.md) to publish to multiple
registries. List the image names in `meta-images` and credentials in the
`registry-auths` secret. This example publishes `latest` and `1.0.0` tags to
Docker Hub and GitHub Container Registry:

```yaml
name: ci

on:
  push:

jobs:
  build:
    uses: docker/github-builder/.github/workflows/build.yml@{{% param "github_builder_version" %}}
    permissions:
      contents: read
      id-token: write
      packages: write
    with:
      output: image
      push: true
      platforms: linux/amd64,linux/arm64
      meta-tags: |
        type=raw,value=latest
        type=raw,value=1.0.0
      meta-images: |
        user/app
        ghcr.io/user/app
    secrets:
      registry-auths: |
        - registry: docker.io
          username: ${{ vars.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}
        - registry: ghcr.io
          username: ${{ github.repository_owner }}
          password: ${{ secrets.GITHUB_TOKEN }}
```

Replace the image names with repositories you can push to. The
`packages: write` permission authorizes publishing to GitHub Container Registry.

## Push with individual actions

The following workflow will connect you to Docker Hub and GitHub Container
Registry, and push the image to both registries:

```yaml
name: ci

on:
  push:

jobs:
  docker:
    permissions:
      contents: read
      packages: write
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
            ghcr.io/user/app:latest
            ghcr.io/user/app:1.0.0
```
