---
title: Docker Build GitHub Actions
linkTitle: GitHub Actions
description: Docker maintains a set of official GitHub Actions for building Docker images.
keywords: ci, github actions, gha,  build, introduction, tutorial
aliases:
  - /guides/gha/
  - /ci-cd/github-actions/
  - /build/ci/github-actions/examples/
---

GitHub Actions is a popular CI/CD platform for automating your build, test, and
deployment pipeline. Docker provides a set of official GitHub Actions for you to
use in your workflows. These official actions are reusable, easy-to-use
components for building, annotating, and pushing images.

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

Using Docker's actions provides an easy-to-use interface, while still allowing
flexibility for customizing build parameters.

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
      - name: Checkout
        uses: actions/checkout@{{% param "checkout_action_version" %}}

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
          context: .
          push: ${{ github.event_name != 'pull_request' }}
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
          annotations: ${{ steps.meta.outputs.annotations }}
```

Commit the workflow file and push it to your default branch to run it. The
workflow also runs on pull requests. Pull-request builds skip signing in to
Docker Hub and pushing the image, so they don't require registry secrets,
including for pull requests from forks.

The checkout step retrieves your repository, and `context: .` uses those files
as the build context. The metadata action generates tags, labels, and
annotations from the GitHub event and passes them to the build step. For
example, a push to `main` publishes `DOCKERHUB_USERNAME/my-image:main`.
See [Manage tags and labels](/manuals/build/ci/github-actions/manage-tags-labels.md)
and [Annotations](/manuals/build/ci/github-actions/annotations.md) to customize
this metadata.

This workflow builds images but doesn't run application tests. To add tests,
see [Test before push](/manuals/build/ci/github-actions/test-before-push.md).
For additional build options, see
[Multi-platform images](/manuals/build/ci/github-actions/multi-platform.md) and
[SBOM and provenance attestations](/manuals/build/ci/github-actions/attestations.md),
which covers default provenance and how to add SBOM attestations when pushing
to a registry.

## Examples

If you're looking for examples on how to use the Docker GitHub Actions,
refer to the following sections:

{{% sectionlinks %}}
