---
title: Validating build configuration with GitHub Actions
linkTitle: Build checks
description: Discover how to validate your build configuration and identify best practice violations using build checks in GitHub Actions.
keywords: github actions, gha, build, checks
---

[Build checks](/manuals/build/checks.md) let you validate your `docker build`
configuration without actually running the build.

## Run checks with `docker/build-push-action`

To run build checks in a GitHub Actions workflow with the `build-push-action`,
set the `call` input parameter to `check`. With this set, the workflow fails if
any check warnings are detected for your build's configuration.

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
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}
      
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Validate build configuration
        uses: docker/build-push-action@v6
        with:
          call: check

      - name: Build and push
        uses: docker/build-push-action@v6
        with:
          push: true
          tags: user/app:latest
```

## Run checks with `docker/bake-action`

If you're using Bake and `docker/bake-action` to run your builds, you don't
need to specify any special inputs in your GitHub Actions workflow
configuration. Instead, define a Bake target that calls the `check` method,
and invoke that target in your CI.

```hcl
target "build" {
  dockerfile = "Dockerfile"
  args = {
    FOO = "bar"
  }
}
target "validate-build" {
  inherits = ["build"]
  call = "check"
}
```

```yaml
name: ci

on:
  push:

env:
  IMAGE_NAME: user/app

jobs:
  docker:
    runs-on: ubuntu-latest
    steps:
      - name: Login to Docker Hub
        uses: docker/login-action@v3
        with:
          username: ${{ vars.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Validate build configuration
        uses: docker/bake-action@v6
        with:
          targets: validate-build

      - name: Build
        uses: docker/bake-action@v6
        with:
          targets: build
          push: true
```

### Using the `call` input directly

You can also set the build method with the `call` input which is equivalent to using the `--call` flag with `docker buildx bake`

For example, to run a check without defining `call` in your Bake file:

```yaml
name: ci

on:
  push:

jobs:
  docker:
    runs-on: ubuntu-latest
    steps:
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Validate build configuration
        uses: docker/bake-action@v6
        with:
          targets: build
          call: check
```
