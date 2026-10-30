---
title: Add SBOM and provenance attestations with GitHub Actions
linkTitle: Attestations
description: Add SBOM and provenance attestations to your images with GitHub Actions
keywords: ci, github actions, gha, buildkit, buildx, attestations, sbom, provenance, slsa
---

Software Bill of Material (SBOM) and provenance
[attestations](../../metadata/attestations/_index.md) add metadata about the contents of
your image, and how it was built.

## Add attestations with Builder

[Docker GitHub Builder](build-workflow.md) manages provenance and signing
in the reusable workflow. It uses max-level provenance for public repositories
and min-level provenance for private repositories. These levels aren't caller
inputs. Use individual actions below if you need to override the provenance
level for a private repository.

Set `sbom: true` to add an SBOM. With the default `sign: auto`, Builder signs
attestation manifests when pushing an image, using the job's `id-token: write`
permission:

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
    with:
      output: image
      push: true
      sbom: true
      meta-images: user/app
    secrets:
      registry-auths: |
        - registry: docker.io
          username: ${{ vars.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}
```

Use [build secrets](secrets.md) for credentials. Build argument values can
appear in provenance, so don't pass secrets as build arguments.

## Use individual actions

Attestations are supported with version 4 and later of the
`docker/build-push-action`.

## Default provenance

The `docker/build-push-action` GitHub Action automatically adds provenance
attestations to your image, with the following conditions:

- If the GitHub repository is public, provenance attestations with `mode=max`
  are automatically added to the image.
- If the GitHub repository is private, provenance attestations with `mode=min`
  are automatically added to the image.
- If you're using the [`docker` exporter](../../exporters/oci-docker.md), or
  you're loading the build results to the runner with `load: true`, no
  attestations are added to the image. These output formats don't support
  attestations.

> [!WARNING]
>
> If you're using `docker/build-push-action` to build images for code in a
> public GitHub repository, the provenance attestations attached to your image
> by default contains the values of build arguments. If you're misusing build
> arguments to pass secrets to your build, such as user credentials or
> authentication tokens, those secrets are exposed in the provenance
> attestation. Refactor your build to pass those secrets using
> [secret mounts](/reference/cli/docker/buildx/build/#secret)
> instead. Also remember to rotate any secrets you may have exposed.

## Max-level provenance

It's recommended that you build your images with max-level provenance
attestations. Private repositories only add min-level provenance by default,
but you can manually override the provenance level by setting the `provenance`
input on the `docker/build-push-action` GitHub Action to `mode=max`.

Note that adding attestations to an image means you must push the image to a
registry directly, as opposed to loading the image to the local image store of
the runner. This is because the local image store doesn't support loading
images with attestations.

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
        uses: docker/login-action@{{% param "login_action_version" %}}
        with:
          username: ${{ vars.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@{{% param "setup_buildx_action_version" %}}

      - name: Extract metadata
        id: meta
        uses: docker/metadata-action@{{% param "metadata_action_version" %}}
        with:
          images: ${{ env.IMAGE_NAME }}

      - name: Build and push image
        uses: docker/build-push-action@{{% param "build_push_action_version" %}}
        with:
          push: true
          provenance: mode=max
          tags: ${{ steps.meta.outputs.tags }}
```

## SBOM

SBOM attestations aren't automatically added to the image. To add SBOM
attestations, set the `sbom` input of the `docker/build-push-action` to true.

Note that adding attestations to an image means you must push the image to a
registry directly, as opposed to loading the image to the local image store of
the runner. This is because the local image store doesn't support loading
images with attestations.

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
        uses: docker/login-action@{{% param "login_action_version" %}}
        with:
          username: ${{ vars.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@{{% param "setup_buildx_action_version" %}}

      - name: Extract metadata
        id: meta
        uses: docker/metadata-action@{{% param "metadata_action_version" %}}
        with:
          images: ${{ env.IMAGE_NAME }}

      - name: Build and push image
        uses: docker/build-push-action@{{% param "build_push_action_version" %}}
        with:
          sbom: true
          push: true
          tags: ${{ steps.meta.outputs.tags }}
```
