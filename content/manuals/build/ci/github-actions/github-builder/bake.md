---
title: Bake with Docker GitHub Builder
linkTitle: Bake workflow
description: Use the Docker GitHub Builder bake.yml reusable workflow to build images and local artifacts from a Bake definition.
keywords: ci, github actions, gha, buildkit, buildx, bake, reusable workflow
weight: 30
---

The [`bake.yml` reusable workflow](https://github.com/docker/github-builder?tab=readme-ov-file#bake-reusable-workflow)
builds from a [Bake definition](../../../bake/_index.md) instead of a Dockerfile
input set. This page shows how to call the workflow for a target, how to pass
Bake overrides and variables, and how to export local output when a Bake file
is already the source of truth for your build.

## Build and push a Bake target

The following workflow builds the `image` target from `docker-bake.hcl` and
publishes the result with tags generated from [metadata inputs](../manage-tags-labels.md):

```yaml
name: ci

on:
  push:
    branches:
      - "main"
    tags:
      - "v*"
  pull_request:

permissions:
  contents: read

jobs:
  bake:
    uses: docker/github-builder/.github/workflows/bake.yml@{{% param "github_builder_version" %}}
    permissions:
      contents: read # to fetch the repository content
      id-token: write # for signing attestation(s) with GitHub OIDC Token
    with:
      output: image
      push: ${{ github.event_name != 'pull_request' }}
      target: image
      meta-images: name/app
      meta-tags: |
        type=ref,event=branch
        type=ref,event=pr
        type=semver,pattern={{version}}
    secrets:
      registry-auths: |
        - registry: docker.io
          username: ${{ vars.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}
```

Bake workflows build one target per workflow call. Groups and multi-target
builds aren't supported because [SLSA provenance](../attestations.md), digest
handling, and manifest creation are scoped to a single target.

The workflow validates the definition before the build starts and resolves
the target from the files you pass in `files`.

## Override target values and variables

Because the workflow delegates the build to Bake, you can keep using `set` and
`vars` for target-specific overrides:

```yaml
name: ci

on:
  push:
    branches:
      - "main"

permissions:
  contents: read

jobs:
  bake:
    uses: docker/github-builder/.github/workflows/bake.yml@{{% param "github_builder_version" %}}
    permissions:
      contents: read # to fetch the repository content
      id-token: write # for signing attestation(s) with GitHub OIDC Token
    with:
      output: image
      push: true
      target: image
      vars: |
        IMAGE_TAG=${{ github.sha }}
      set: |
        *.args.BUILD_RUN_ID=${{ github.run_id }}
        *.platform=linux/amd64,linux/arm64
      cache: true
      cache-scope: image
      meta-images: name/app
      meta-tags: |
        type=sha
      set-meta-annotations: true
    secrets:
      registry-auths: |
        - registry: docker.io
          username: ${{ vars.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}
```

This form fits repositories that already use Bake groups, target inheritance,
and variable expansion. The reusable workflow takes care of Buildx setup,
[GitHub Actions cache export](../../../cache/backends/gha.md),
[Provenance defaults](../../../metadata/attestations/slsa-provenance.md),
signing behavior, and the final multi-platform manifest. Metadata labels and
annotations can be merged into the Bake definition without adding a separate
metadata step to your workflow.

## Export local output from Bake

If the target should export files instead of publishing an image, switch the
workflow output to `local` and upload the artifact:

```yaml
name: ci

on:
  pull_request:

permissions:
  contents: read

jobs:
  bake:
    uses: docker/github-builder/.github/workflows/bake.yml@{{% param "github_builder_version" %}}
    permissions:
      contents: read # to fetch the repository content
      id-token: write # for signing attestation(s) with GitHub OIDC Token
    with:
      output: local
      target: binaries
      artifact-upload: true
      artifact-name: bake-output
```

With `output: local`, the workflow injects the matching local output override
into the Bake run and merges the uploaded artifacts after the per-platform
builds finish. If you need a manual Bake pattern that stays in a normal job,
see [Multi-platform image](../multi-platform.md). If your build does not need a
Bake definition, use [build.yml](build.md) instead.
