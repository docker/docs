---
title: Docker Build GitHub Actions
linkTitle: GitHub Actions
description: Build and publish images with Docker GitHub Builder or individual Docker GitHub Actions.
keywords: ci, github actions, gha,  build, introduction, tutorial
aliases:
  - /guides/gha/
  - /ci-cd/github-actions/
  - /build/ci/github-actions/examples/
---

GitHub Actions automates your build, test, and deployment pipeline. Use
[Docker GitHub Builder](github-builder/_index.md) to build and publish images
with Docker-maintained reusable workflows. Builder manages Buildx setup,
platform distribution, metadata, and signing. Start with Builder for supported
build tasks, and compose individual actions when you need custom job steps or
controls outside the reusable workflow interface.

## Get started

Build an image from your GitHub repository and push it to Docker Hub with the
following workflow.

### Prerequisites

- A GitHub repository with your application and a `Dockerfile` at its root.
- A Docker Hub account and a
  [personal access token](/manuals/security/access-tokens/personal-access-tokens.md)
  with permission to push images.

### Configure your GitHub repository

Add your Docker Hub credentials to your GitHub repository:

1. Open your repository's **Settings**.
2. Under **Security**, go to **Secrets and variables > Actions**.
3. Under **Secrets**, create a repository secret named `DOCKERHUB_TOKEN`
   containing your Docker access token.
4. Under **Variables**, create a repository variable named `DOCKERHUB_USERNAME`
   containing your Docker Hub username.

### Set up your GitHub Actions workflow

Create `.github/workflows/docker-ci.yml` in your repository with the following
contents. Replace `my-image` with the name of your Docker Hub repository, and
change `main` if your default branch has a different name.

```yaml
name: ci

on:
  push:
    branches:
      - main
  pull_request:

jobs:
  build:
    uses: docker/github-builder/.github/workflows/build.yml@{{% param "github_builder_version" %}}
    permissions:
      contents: read
      id-token: write
    with:
      output: image
      push: ${{ github.event_name != 'pull_request' }}
      set-meta-labels: true
      set-meta-annotations: true
      meta-images: ${{ vars.DOCKERHUB_USERNAME }}/my-image
    secrets:
      registry-auths: |
        - registry: docker.io
          username: ${{ vars.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}
```

Call the reusable workflow with `uses` at the job level, rather than inside
`steps`. Grant `contents: read` to fetch the repository and `id-token: write`
for signing with GitHub's OpenID Connect token. With the default `sign: auto`,
Builder signs attestation manifests when pushing an image. See the
[Build workflow](github-builder/build.md) for inputs and outputs.

Commit the workflow file and push it to your default branch to run it. The
workflow also runs on pull requests. With `push: false`, the default
`registry-login: auto` skips registry login and `sign: auto` skips signing.
Pull requests from forks can build without Docker Hub secrets or a writable
OpenID Connect token. Builds that need private dependencies still need access
to those dependencies.

Builder checks out the repository for the triggering event in its own jobs.
Files changed in a caller job aren't automatically available to the build.

Builder uses the metadata action to generate tags from the GitHub event.
For example, a push to `main` publishes `DOCKERHUB_USERNAME/my-image:main`.
The `set-meta-labels` and `set-meta-annotations` switches also apply generated
labels and annotations; both are disabled by default. See
[Manage tags and labels](manage-tags-labels.md) and [Annotations](annotations.md)
to customize this metadata.

This workflow builds images but doesn't run application tests. To add tests,
see [Test before push](test-before-push.md). For additional build options, see
[Multi-platform images](multi-platform.md) and
[SBOM and provenance attestations](attestations.md).

## Build with individual actions

Use individual actions to load images into a job's Docker image store, run
steps before or after a build on the same runner, or configure the builder
directly. The following alternative uses the same credentials, triggers, and
metadata as the Builder example:

```yaml
name: Build and push Docker image

on:
  push:
    branches:
      - main
  pull_request:

permissions:
  contents: read

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Extract Docker image metadata
        id: meta
        uses: docker/metadata-action@{{% param "metadata_action_version" %}}
        with:
          images: ${{ vars.DOCKERHUB_USERNAME }}/my-image

      - name: Sign in to Docker Hub
        if: github.event_name != 'pull_request'
        uses: docker/login-action@{{% param "login_action_version" %}}
        with:
          username: ${{ vars.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@{{% param "setup_buildx_action_version" %}}

      - name: Build and push Docker image
        uses: docker/build-push-action@{{% param "build_push_action_version" %}}
        with:
          push: ${{ github.event_name != 'pull_request' }}
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
          annotations: ${{ steps.meta.outputs.annotations }}
```

The build action uses Git context by default: BuildKit fetches the repository
at the triggering Git reference. No checkout step is needed, and changes made
by earlier steps aren't included. Use a checkout and a path context when the
build needs files changed by earlier steps.

### Available actions

The following GitHub Actions are available:

- [Build and push Docker images](https://github.com/marketplace/actions/build-and-push-docker-images):
  build and push Docker images with BuildKit.
- [Docker Buildx Bake](https://github.com/marketplace/actions/docker-buildx-bake):
  enables using high-level builds with [Bake](../../bake/_index.md).
- [Docker Login](https://github.com/marketplace/actions/docker-login):
  sign in to a Docker registry.
- [Docker Setup Buildx](https://github.com/marketplace/actions/docker-setup-buildx):
  creates and boots a BuildKit builder.
- [Docker Metadata action](https://github.com/marketplace/actions/docker-metadata-action):
  extracts metadata from Git reference and GitHub events to generate tags,
  labels, and annotations.
- [Docker Setup Compose](https://github.com/marketplace/actions/docker-setup-compose):
  installs and sets up [Compose](../../../compose).
- [Docker Setup Docker](https://github.com/marketplace/actions/docker-setup-docker):
  installs Docker Engine.
- [Docker Setup QEMU](https://github.com/marketplace/actions/docker-setup-qemu):
  installs [QEMU](https://github.com/qemu/qemu) static binaries for
  multi-platform builds.
- [Docker Scout](https://github.com/docker/scout-action):
  analyze Docker images for security vulnerabilities.

## Examples

If you're looking for examples on how to use the Docker GitHub Actions,
refer to the following sections:

{{% sectionlinks %}}
