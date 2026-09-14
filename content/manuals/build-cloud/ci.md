---
title: Use Docker Build Cloud in CI
linkTitle: Continuous integration
weight: 30
description: Speed up your continuous integration pipelines with Docker Build Cloud in CI
keywords: build, cloud build, ci, gha, gitlab, buildkite, jenkins, circle ci
aliases:
  - /build/cloud/ci/
---

Using Docker Build Cloud in CI can speed up your build pipelines, which means less time
spent waiting and context switching. You control your CI workflows as usual,
and delegate the build execution to Docker Build Cloud.

Building with Docker Build Cloud in CI involves the following steps:

1. Sign in to a Docker account.
2. Set up Buildx and connect to the builder.
3. Run the build.

When using Docker Build Cloud in CI, it's recommended that you push the result to a
registry directly, rather than loading the image and then pushing it. Pushing
directly speeds up your builds and avoids unnecessary file transfers.

If you just want to build and discard the output, export the results to the
build cache or build without tagging the image. With no explicit output, Buildx
leaves an untagged result in the cloud build cache and automatically loads
eligible tagged images. See
[Loading build results](./usage/#loading-build-results) for details.

> [!NOTE]
>
> Builds on Docker Build Cloud have a timeout limit that depends on your
> subscription plan. Builds that run for longer than your limit are
> automatically cancelled.

## Setting up credentials for CI/CD

To enable your CI/CD system to build and push images using Docker Build Cloud, provide both an access token and a username. The type of token and the username you use depend on your account type and permissions.

- If you are an organization administrator or have permission to create [organization access tokens (OAT)](/manuals/security/access-tokens/organization-access-tokens.md), use an OAT and set `DOCKER_ACCOUNT` to your Docker Hub organization name.
- If you do not have permission to create OATs or are using a personal account, use a [personal access token (PAT)](/security/access-tokens/personal-access-tokens/) and set `DOCKER_ACCOUNT` to your Docker Hub username.

### Creating access tokens

#### For organization accounts

If you are an organization administrator:

- Create an [organization access token (OAT)](/manuals/security/access-tokens/organization-access-tokens.md). The token must have these permissions:
    1. **cloud-connect** scope
    2. **Read public repositories** permission
    3. **Repository access** with **Image push** permission for the target repository:
        - Expand the **Repository** drop-down.
        - Select **Add repository** and choose your target repository.
        - Set the **Image push** permission for the repository.

If you are not an organization administrator:

- Ask your organization administrator for an access token with the permissions listed above, or use a personal access token.

#### For personal accounts

- Create a [personal access token (PAT)](/security/access-tokens/personal-access-tokens/) with the following permissions:
   1. **Read & write** access.
        - Note: Building with Docker Build Cloud only requires read access, but you need write access to push images to a Docker Hub repository.


## CI platform examples

> [!IMPORTANT]
>
> These examples require Buildx version 0.37.0 or later, which includes the
> `cloud` driver. The Docker CLI version doesn't determine the installed Buildx
> plugin version, and CI runner images may include an earlier Buildx release.
>
> Check the installed version:
>
> ```console
> $ docker buildx version
> ```
>
> If the version is earlier than 0.37.0, install Buildx as a
> [Docker CLI plugin](https://github.com/docker/buildx#manual-download).

> [!NOTE]
>
> In your CI/CD configuration, set the following variables/secrets:
> - `DOCKER_ACCESS_TOKEN` — your access token (PAT or OAT). Use a secret to store the token.
> - `DOCKER_ACCOUNT` — your Docker Hub organization name (for OAT) or username (for PAT)
> - `CLOUD_BUILDER_NAME` — the name of the cloud builder you created in the [Docker Build Cloud Dashboard](https://app.docker.com/build/)
>
> This ensures your builds authenticate correctly with Docker Build Cloud.

### GitHub Actions

<!-- TODO: Confirm whether standard Buildx requires a minimum setup-buildx-action version. -->

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
      - name: Login to Docker Hub
        uses: docker/login-action@{{% param "login_action_version" %}}
        with:
          username: ${{ vars.DOCKER_ACCOUNT }}
          password: ${{ secrets.DOCKER_ACCESS_TOKEN }}
      
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@{{% param "setup_buildx_action_version" %}}
        with:
          driver: cloud
          endpoint: "${{ vars.DOCKER_ACCOUNT }}/${{ vars.CLOUD_BUILDER_NAME }}" # for example, "acme/default"
      
      - name: Build and push
        uses: docker/build-push-action@{{% param "build_push_action_version" %}}
        with:
          tags: "<IMAGE>" # for example, "acme/my-image:latest"
          # For pull requests, export results to the build cache.
          # Otherwise, push to a registry.
          outputs: ${{ github.event_name == 'pull_request' && 'type=cacheonly' || 'type=registry' }}
```

The example above uses `docker/build-push-action`, which automatically uses the
builder set up by `setup-buildx-action`. If you need to use the `docker build`
command directly instead, you have two options:

- Use `docker buildx build` instead of `docker build`
- Set the `BUILDX_BUILDER` environment variable to use the cloud builder:

  ```yaml
  - name: Set up Docker Buildx
    id: builder
    uses: docker/setup-buildx-action@{{% param "setup_buildx_action_version" %}}
    with:
      driver: cloud
      endpoint: "${{ vars.DOCKER_ACCOUNT }}/${{ vars.CLOUD_BUILDER_NAME }}"

  - name: Build
    run: |
      docker build .
    env:
      BUILDX_BUILDER: ${{ steps.builder.outputs.name }}
  ```

For more information about the `BUILDX_BUILDER` environment variable, see
[Build variables](/manuals/build/building/variables.md#buildx_builder).

### GitLab

```yaml
default:
  image: docker:cli
  services:
    - docker:dind
  before_script:
    - docker info
    - echo "$DOCKER_ACCESS_TOKEN" | docker login --username "$DOCKER_ACCOUNT" --password-stdin
    - docker buildx create --use --driver cloud ${DOCKER_ACCOUNT}/${CLOUD_BUILDER_NAME}

variables:
  IMAGE_NAME: <IMAGE>
  DOCKER_ACCOUNT: <DOCKER_ACCOUNT> # your Docker Hub organization name (or username when using a personal account)
  CLOUD_BUILDER_NAME: <CLOUD_BUILDER_NAME> # the name of the cloud builder you created in the [Docker Build Cloud Dashboard](https://app.docker.com/build/)

# Build multi-platform image and push to a registry
build_push:
  stage: build
  script:
    - |
      docker buildx build \
        --platform linux/amd64,linux/arm64 \
        --tag "${IMAGE_NAME}:${CI_COMMIT_SHORT_SHA}" \
        --push .

# Build an image and discard the result
build_cache:
  stage: build
  script:
    - |
      docker buildx build \
        --platform linux/amd64,linux/arm64 \
        --tag "${IMAGE_NAME}:${CI_COMMIT_SHORT_SHA}" \
        --output type=cacheonly \
        .
```

### Circle CI

```yaml
version: 2.1

jobs:
  # Build multi-platform image and push to a registry
  build_push:
    machine:
      image: ubuntu-2204:current
    steps:
      - checkout

      - run: echo "$DOCKER_ACCESS_TOKEN" | docker login --username $DOCKER_ACCOUNT --password-stdin
      - run: docker buildx create --use --driver cloud "${DOCKER_ACCOUNT}/${CLOUD_BUILDER_NAME}"

      - run: |
          docker buildx build \
          --platform linux/amd64,linux/arm64 \
          --push \
          --tag "<IMAGE>" .

  # Build an image and discard the result
  build_cache:
    machine:
      image: ubuntu-2204:current
    steps:
      - checkout

      - run: echo "$DOCKER_ACCESS_TOKEN" | docker login --username $DOCKER_ACCOUNT --password-stdin
      - run: docker buildx create --use --driver cloud "${DOCKER_ACCOUNT}/${CLOUD_BUILDER_NAME}"

      - run: |
          docker buildx build \
          --tag temp \
          --output type=cacheonly \
          .

workflows:
  pull_request:
    jobs:
      - build_cache
  release:
    jobs:
      - build_push
```

### Buildkite

The following example sets up a Buildkite pipeline using Docker Build Cloud. The
example assumes that the pipeline name is `build-push-docker` and that you
manage the Docker access token using environment hooks, but feel free to adapt
this to your needs.

Add the following `environment` hook agent's hook directory:

```bash
#!/bin/bash
set -euo pipefail

if [[ "$BUILDKITE_PIPELINE_NAME" == "build-push-docker" ]]; then
 export DOCKER_ACCESS_TOKEN="<DOCKER_ACCESS_TOKEN>"
fi
```

Create a `pipeline.yml` that uses the `docker-login` plugin:

```yaml
env:
  DOCKER_ACCOUNT: <DOCKER_ACCOUNT> # your Docker Hub organization name (or username when using a personal account)
  CLOUD_BUILDER_NAME: <CLOUD_BUILDER_NAME> # the name of the cloud builder you created in the [Docker Build Cloud Dashboard](https://app.docker.com/build/)
  IMAGE_NAME: <IMAGE>

steps:
  - command: ./build.sh
    key: build-push
    plugins:
      - docker-login#v2.1.0:
          username: DOCKER_ACCOUNT
          password-env: DOCKER_ACCESS_TOKEN # the variable name in the environment hook
```

Create the `build.sh` script:

```bash
# Connect to your builder and set it as the default builder
docker buildx create --use --driver cloud "${DOCKER_ACCOUNT}/${CLOUD_BUILDER_NAME}"

# Cache-only image build
docker buildx build \
    --platform linux/amd64,linux/arm64 \
    --tag "$IMAGE_NAME:$BUILDKITE_COMMIT" \
    --output type=cacheonly \
    .

# Build, tag, and push a multi-arch docker image
docker buildx build \
    --platform linux/amd64,linux/arm64 \
    --push \
    --tag "$IMAGE_NAME:$BUILDKITE_COMMIT" \
    .
```

### Jenkins

```groovy
pipeline {
  agent any

  environment {
    DOCKER_ACCESS_TOKEN = credentials('docker-access-token')
    DOCKER_ACCOUNT = credentials('docker-account')
    CLOUD_BUILDER_NAME = '<CLOUD_BUILDER_NAME>'
    IMAGE_NAME = '<IMAGE>'
  }

  stages {
    stage('Build') {
      steps {
        sh 'echo "$DOCKER_ACCESS_TOKEN" | docker login --username $DOCKER_ACCOUNT --password-stdin'
        sh 'docker buildx create --use --driver cloud "${DOCKER_ACCOUNT}/${CLOUD_BUILDER_NAME}"'
        // Cache-only build
        sh 'docker buildx build --platform linux/amd64,linux/arm64 --tag "$IMAGE_NAME" --output type=cacheonly .'
        // Build and push a multi-platform image
        sh 'docker buildx build --platform linux/amd64,linux/arm64 --push --tag "$IMAGE_NAME" .'
      }
    }
  }
}
```

### Travis CI

```yaml
language: minimal 
dist: jammy 

services:
  - docker

env:
  global:
    - IMAGE_NAME=<IMAGE> # for example, "acme/my-image:latest"

before_install: |
  echo "$DOCKER_ACCESS_TOKEN" | docker login --username "$DOCKER_ACCOUNT" --password-stdin

before_script: |
  docker buildx create --use --driver cloud "${DOCKER_ACCOUNT}/${CLOUD_BUILDER_NAME}"

script: |
  docker buildx build \
  --platform linux/amd64,linux/arm64 \
  --push \
  --tag "$IMAGE_NAME" .
```

### BitBucket Pipelines 

```yaml
# Prerequisites: $DOCKER_ACCOUNT, $CLOUD_BUILDER_NAME, $DOCKER_ACCESS_TOKEN setup as deployment variables
# This pipeline assumes $BITBUCKET_REPO_SLUG as the image name

image: docker:cli

pipelines:
  default:
    - step:
        name: Build multi-platform image
        script:
          - echo "$DOCKER_ACCESS_TOKEN" | docker login --username $DOCKER_ACCOUNT --password-stdin
          - docker buildx create --use --driver cloud "${DOCKER_ACCOUNT}/${CLOUD_BUILDER_NAME}"
          - IMAGE_NAME=$BITBUCKET_REPO_SLUG
          - docker buildx build
            --platform linux/amd64,linux/arm64
            --push
            --tag "$IMAGE_NAME" .
        services:
          - docker
```

### Shell script

```bash
#!/bin/bash

# Login to Docker Hub with an access token. See https://docs.docker.com/build-cloud/ci/#creating-access-tokens
echo "$DOCKER_ACCESS_TOKEN" | docker login --username $DOCKER_ACCOUNT --password-stdin

# Connect to your builder and set it as the default builder
docker buildx create --use --driver cloud "${DOCKER_ACCOUNT}/${CLOUD_BUILDER_NAME}"

# Cache-only image build
docker buildx build \
    --tag temp \
    --output type=cacheonly \
    .

# Build, tag, and push a multi-arch docker image
docker buildx build \
    --platform linux/amd64,linux/arm64 \
    --push \
    --tag "<IMAGE>" \
    .
```

### Docker Compose

This example requires Buildx version 0.37.0 or later. Set the cloud builder as
the default builder before running `docker compose build` in CI:

```bash
#!/bin/bash

# Login to Docker Hub with an access token. See https://docs.docker.com/build-cloud/ci/#creating-access-tokens
echo "$DOCKER_ACCESS_TOKEN" | docker login --username $DOCKER_ACCOUNT --password-stdin

# Connect to your builder and set it as the default builder
docker buildx create --use --driver cloud "${DOCKER_ACCOUNT}/${CLOUD_BUILDER_NAME}"

# Build the images
docker compose build
```
