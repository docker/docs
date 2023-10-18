---
description: How to integrate Docker Scout with GitHub Actions
keywords: supply chain, security, ci, continuous integration, github actions
title: Integrate Docker Scout with GitHub Actions
---

You can use [the Docker Scout GitHub action](https://github.com/docker/scout-action) to run Docker Scout CLI commands
as part of a workflow.

The following example works in a repository containing a Docker image's
definition and contents. Triggered by a pull request, the action builds the
image and uses Docker Scout to compare the new version to the current published
version.

First, set up the rest of the workflow. There's a lot that's not specific to
Docker Scout but needed to create the images to compare. For more details on
those actions and using GitHub Actions with Docker in general, see [the GitHub Actions documentation](../../../build/ci/github-actions/index.md).

Add the following to a GitHub Actions YAML file:


```yaml
name: Docker

on:
  push:
    tags: ["*"]
    branches:
      - "main"
  pull_request:
    branches: ["**"]

env:
  # Use docker.io for Docker Hub if empty
  REGISTRY: docker.io
  IMAGE_NAME: ${{ github.repository }}
  SHA: ${{ github.event.pull_request.head.sha || github.event.after }}

jobs:
  build:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write
```


This sets up the workflow to run on pull requests and pushes to the `main`
branch, and sets up environment variables available to all workflow steps. It
then defines a job called `build` that runs on the latest Ubuntu image and sets
the permissions available to the job.

Add the following to the YAML file:


```yaml
steps:
  - name: Checkout repository
    uses: actions/checkout@v4
    with:
      ref: ${{ env.SHA }}

  - name: Setup Docker buildx
    uses: docker/setup-buildx-action@v3

  # Login against a Docker registry except on PR
  # https://github.com/docker/login-action
  - name: Log into registry ${{ env.REGISTRY }}
    uses: docker/login-action@v3
    with:
      registry: ${{ env.REGISTRY }}
      username: ${{ secrets.DOCKER_USER }}
      password: ${{ secrets.DOCKER_PAT }}

  # Extract metadata (tags, labels) for Docker
  # https://github.com/docker/metadata-action
  - name: Extract Docker metadata
    id: meta
    uses: docker/metadata-action@v5
    with:
      images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}
      labels: |
        org.opencontainers.image.revision=${{ env.SHA }}
      tags: |
        type=edge,branch=$repo.default_branch
        type=semver,pattern=v{{version}}
        type=sha,prefix=,suffix=,format=short
```


This creates workflow steps to checkout the repository, set up Docker buildx,
log into the Docker registry, and extract metadata from Git reference and GitHub
events to use in later steps.

Add the following to the YAML file:


```yaml
# Build and push Docker image with Buildx (don't push on PR)
# https://github.com/docker/build-push-action
- name: Build and push Docker image
  id: build-and-push
  uses: docker/build-push-action@v5
  with:
    context: .
    push: true
    tags: ${{ steps.meta.outputs.tags }}
    labels: ${{ steps.meta.outputs.labels }}
    cache-from: type=gha
    cache-to: type=gha,mode=max
```


This uses the extracted metadata from the previous step to build and push the
Docker image to Docker Hub. GitHub Actions skips this step on pull requests and
only runs when a pull request is merged.

Add the following to the YAML file:


```yaml
- name: Docker Scout
  id: docker-scout
  if: ${{ github.event_name == 'pull_request' }}
  uses: docker/scout-action@v1
  with:
    command: compare
    image: ${{ steps.meta.outputs.tags }}
    to: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:edge
    ignore-unchanged: true
    only-severities: critical,high
    token: ${{ secrets.DOCKER_PAT }}
```


This final step uses the Docker Scout CLI to run [the `compare` command](../../../engine/reference/commandline/scout_compare.md), comparing the new
image to the published one. It only shows critical or high-severity
vulnerabilities and ignores vulnerabilities that haven't changed since the last
analysis.

The GitHub Action outputs the comparison results as a table and a summary in the
action output.

![A screenshot showing the results of Docker Scout output in a GitHub Action](../../images/gha-output.png)
