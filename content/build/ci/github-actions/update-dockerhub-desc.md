---
title: Update Docker Hub description with GitHub Actions
description: How to update the repository README in Docker Hub using with GitHub Actions
keywords: ci, github actions, gha, buildkit, buildx, docker hub
---

You can update the Docker Hub repository description using a third party action
called [Docker Hub Description](https://github.com/peter-evans/dockerhub-description)
with this action:

```yaml
name: ci

on:
  push:

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
        uses: docker/build-push-action@v6
        with:
          context: .
          push: true
          tags: user/app:latest
      
      - name: Update repo description
        uses: peter-evans/dockerhub-description@e98e4d1628a5f3be2be7c231e50981aee98723ae # v4.0.0
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}
          repository: user/app
```
