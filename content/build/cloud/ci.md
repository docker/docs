---
title: Use Docker Build Cloud in CI
description: Speed up your continuous integration pipelines with Docker Build Cloud in CI
keywords: build, cloud build, ci, gha, gitlab, buildkite, jenkins, circle ci
aliases:
  - /hydrobuild/ci/
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
build cache or build without tagging the image. When you use Docker Build Cloud,
Buildx automatically loads the build result if you build a tagged image.
See [Loading build results](./usage/#loading-build-results) for details.

> **Note**
>
> Builds on Docker Build Cloud have a timeout limit of two hours. Builds that
> run for longer than two hours are automatically cancelled.

{{< tabs >}}
{{< tab name="GitHub Actions" >}}

> **Note**
>
> Version 4.0.0 and later of `docker/build-push-action` and
> `docker/bake-action` builds images with [provenance attestations by
> default](../ci/github-actions/attestations.md#default-provenance). Docker
> Build Cloud automatically attempts to load images to the local image store if
> you don't explicitly push them to a registry.
>
> This results in a conflicting scenario where if you build a tagged image
> without pushing it to a registry, Docker Build Cloud attempts to load images
> containing attestations. But the local image store on the GitHub runner
> doesn't support attestations, and the image load fails as a result.
>
> If you want to load images built with `docker/build-push-action` together
> with Docker Build Cloud, you must disable provenance attestations by setting
> `provenance: false` in the GitHub Action inputs (or in `docker-bake.hcl` if
> you use Bake).

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
      - name: Log in to Docker Hub
        uses: docker/login-action@v3
        with:
          username: ${{ secrets.DOCKER_USER }}
          password: ${{ secrets.DOCKER_PAT }}
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3
        with:
          version: "lab:latest"
          driver: cloud
          endpoint: "<ORG>/default"
          install: true
      - name: Build and push
        uses: docker/build-push-action@v5
        with:
          context: .
          tags: "<IMAGE>"
          # For pull requests, export results to the build cache.
          # Otherwise, push to a registry.
          outputs: ${{ github.event_name == 'pull_request' && 'type=cacheonly' || 'type=registry,push=true' }}
```

{{< /tab >}}
{{< tab name="GitLab" >}}

```yaml
default:
  image: docker:24-dind
  services:
    - docker:24-dind
  before_script:
    - docker info
    - echo "$DOCKER_PAT" | docker login --username "$DOCKER_USER" --password-stdin
    - |
      apk add curl jq
      ARCH=${CI_RUNNER_EXECUTABLE_ARCH#*/}
      BUILDX_URL=$(curl -s https://raw.githubusercontent.com/docker/actions-toolkit/main/.github/buildx-lab-releases.json | jq -r ".latest.assets[] | select(endswith(\"linux-$ARCH\"))")
      mkdir -vp ~/.docker/cli-plugins/
      curl --silent -L --output ~/.docker/cli-plugins/docker-buildx $BUILDX_URL
      chmod a+x ~/.docker/cli-plugins/docker-buildx
    - docker buildx create --use --driver cloud ${DOCKER_ORG}/default

variables:
  IMAGE_NAME: <IMAGE>
  DOCKER_ORG: <ORG>

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

{{< /tab >}}
{{< tab name="Circle CI" >}}

```yaml
version: 2.1

jobs:
  # Build multi-platform image and push to a registry
  build_push:
    machine:
      image: ubuntu-2204:current
    steps:
      - checkout

      - run: |
          mkdir -vp ~/.docker/cli-plugins/
          ARCH=amd64
          BUILDX_URL=$(curl -s https://raw.githubusercontent.com/docker/actions-toolkit/main/.github/buildx-lab-releases.json | jq -r ".latest.assets[] | select(endswith(\"linux-$ARCH\"))")
          curl --silent -L --output ~/.docker/cli-plugins/docker-buildx $BUILDX_URL
          chmod a+x ~/.docker/cli-plugins/docker-buildx

      - run: echo "$DOCKER_PAT" | docker login --username $DOCKER_USER --password-stdin
      - run: docker buildx create --use --driver cloud "<ORG>/default"

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

      - run: |
          mkdir -vp ~/.docker/cli-plugins/
          ARCH=amd64
          BUILDX_URL=$(curl -s https://raw.githubusercontent.com/docker/actions-toolkit/main/.github/buildx-lab-releases.json | jq -r ".latest.assets[] | select(endswith(\"linux-$ARCH\"))")
          curl --silent -L --output ~/.docker/cli-plugins/docker-buildx $BUILDX_URL
          chmod a+x ~/.docker/cli-plugins/docker-buildx

      - run: echo "$DOCKER_PAT" | docker login --username $DOCKER_USER --password-stdin
      - run: docker buildx create --use --driver cloud "<ORG>/default"

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

{{< /tab >}}
{{< tab name="Buildkite" >}}

The following example sets up a Buildkite pipeline using Docker Build Cloud. The
example assumes that the pipeline name is `build-push-docker` and that you
manage the Docker access token using environment hooks, but feel free to adapt
this to your needs.

Add the following `environment` hook agent's hook directory:

```bash
#!/bin/bash
set -euo pipefail

if [[ "$BUILDKITE_PIPELINE_NAME" == "build-push-docker" ]]; then
 export DOCKER_PAT="<DOCKER_PERSONAL_ACCESS_TOKEN>"
fi
```

Create a `pipeline.yml` that uses the `docker-login` plugin:

```yaml
env:
  DOCKER_ORG: <ORG>
  IMAGE_NAME: <IMAGE>

steps:
  - command: ./build.sh
    key: build-push
    plugins:
      - docker-login#v2.1.0:
          username: <DOCKER_USER>
          password-env: DOCKER_PAT # the variable name in the environment hook
```

Create the `build.sh` script:

```bash
DOCKER_DIR=/usr/libexec/docker

# Get download link for latest buildx binary.
# Set $ARCH to the CPU architecture (e.g. amd64, arm64)
UNAME_ARCH=`uname -m`
case $UNAME_ARCH in
  aarch64)
    ARCH="arm64";
    ;;
  amd64)
    ARCH="amd64";
    ;;
  *)
    ARCH="amd64";
    ;;
esac
BUILDX_URL=$(curl -s https://raw.githubusercontent.com/docker/actions-toolkit/main/.github/buildx-lab-releases.json | jq -r ".latest.assets[] | select(endswith(\"linux-$ARCH\"))")

# Download docker buildx with Hyrdobuild support
curl --silent -L --output $DOCKER_DIR/cli-plugins/docker-buildx $BUILDX_URL
chmod a+x ~/.docker/cli-plugins/docker-buildx

# Connect to your builder and set it as the default builder
docker buildx create --use --driver cloud "$DOCKER_ORG/default"

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

{{< /tab >}}
{{< tab name="Jenkins" >}}

```groovy
pipeline {
  agent any

  environment {
    ARCH = 'amd64'
    DOCKER_PAT = credentials('docker-personal-access-token')
    DOCKER_USER = credentials('docker-username')
    DOCKER_ORG = '<ORG>'
    IMAGE_NAME = '<IMAGE>'
  }

  stages {
    stage('Build') {
      environment {
        BUILDX_URL = sh (returnStdout: true, script: 'curl -s https://raw.githubusercontent.com/docker/actions-toolkit/main/.github/buildx-lab-releases.json | jq -r ".latest.assets[] | select(endswith(\\"linux-$ARCH\\"))"').trim()
      }
      steps {
        sh 'mkdir -vp ~/.docker/cli-plugins/'
        sh 'curl --silent -L --output ~/.docker/cli-plugins/docker-buildx $BUILDX_URL'
        sh 'chmod a+x ~/.docker/cli-plugins/docker-buildx'
        sh 'echo "$DOCKER_PAT" | docker login --username $DOCKER_USER --password-stdin'
        sh 'docker buildx create --use --driver cloud "$DOCKER_ORG/default"'
        // Cache-only build
        sh 'docker buildx build --platform linux/amd64,linux/arm64 --tag "$IMAGE_NAME" --output type=cacheonly .'
        // Build and push a multi-platform image
        sh 'docker buildx build --platform linux/amd64,linux/arm64 --push --tag "$IMAGE_NAME" .'
      }
    }
  }
}
```

{{< /tab >}}
{{< tab name="Shell" >}}

```bash
#!/bin/bash

# Get download link for latest buildx binary. Set $ARCH to the CPU architecture (e.g. amd64, arm64)
ARCH=amd64
BUILDX_URL=$(curl -s https://raw.githubusercontent.com/docker/actions-toolkit/main/.github/buildx-lab-releases.json | jq -r ".latest.assets[] | select(endswith(\"linux-$ARCH\"))")

# Download docker buildx with Hyrdobuild support
mkdir -vp ~/.docker/cli-plugins/
curl --silent -L --output ~/.docker/cli-plugins/docker-buildx $BUILDX_URL
chmod a+x ~/.docker/cli-plugins/docker-buildx

# Login to Docker Hub. For security reasons $DOCKER_PAT should be a Personal Access Token. See https://docs.docker.com/security/for-developers/access-tokens/
echo "$DOCKER_PAT" | docker login --username $DOCKER_USER --password-stdin

# Connect to your builder and set it as the default builder
docker buildx create --use --driver cloud "<ORG>/default"

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

{{< /tab >}}
{{< tab name="Docker Compose" >}}

Use this implementation if you want to use `docker compose build` with
Docker Build Cloud in CI.

```bash
#!/bin/bash

# Get download link for latest buildx binary. Set $ARCH to the CPU architecture (e.g. amd64, arm64)
ARCH=amd64
BUILDX_URL=$(curl -s https://raw.githubusercontent.com/docker/actions-toolkit/main/.github/buildx-lab-releases.json | jq -r ".latest.assets[] | select(endswith(\"linux-$ARCH\"))")
COMPOSE_URL=$(curl -sL \
  -H "Accept: application/vnd.github+json" \
  -H "Authorization: Bearer <GITHUB_TOKEN>" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  https://api.github.com/repos/docker/compose-desktop/releases \
  | jq "[ .[] | select(.prerelease==false and .draft==false) ] | .[0].assets.[] | select(.name | endswith(\"linux-${ARCH}\")) | .browser_download_url")

# Download docker buildx with Hyrdobuild support
mkdir -vp ~/.docker/cli-plugins/
curl --silent -L --output ~/.docker/cli-plugins/docker-buildx $BUILDX_URL
curl --silent -L --output ~/.docker/cli-plugins/docker-compose $COMPOSE_URL
chmod a+x ~/.docker/cli-plugins/docker-buildx
chmod a+x ~/.docker/cli-plugins/docker-compose

# Login to Docker Hub. For security reasons $DOCKER_PAT should be a Personal Access Token. See https://docs.docker.com/security/for-developers/access-tokens/
echo "$DOCKER_PAT" | docker login --username $DOCKER_USER --password-stdin

# Connect to your builder and set it as the default builder
docker buildx create --use --driver cloud "<ORG>/default"

# Build the image build
docker compose build
```

{{< /tab >}}
{{< /tabs >}}
