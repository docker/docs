---
title: Docker GitHub Builder
linkTitle: GitHub Builder
description: Use Docker-maintained reusable GitHub Actions workflows to build images and artifacts with BuildKit.
keywords: ci, github actions, gha, buildkit, buildx, bake, reusable workflows
params:
  sidebar:
    badge:
      color: green
      text: New
---

Docker GitHub Builder is a set of [reusable workflows](https://docs.github.com/en/actions/how-tos/reuse-automations/reuse-workflows)
in the [`docker/github-builder` repository](https://github.com/docker/github-builder)
for building container images and local artifacts with [BuildKit](../../../buildkit/_index.md).
This section explains what the workflows solve, how they differ from wiring
together individual GitHub Actions in each repository, and when to use
[`build.yml`](build.md) or [`bake.yml`](bake.md).

If you compose a build job from `docker/login-action`, `docker/setup-buildx-action`,
`docker/metadata-action`, and either `docker/build-push-action` or
`docker/bake-action`, your repository owns every detail of how the build runs.
That approach works, but it also means every repository has to maintain its own
runner selection, [cache setup](../cache.md), [Provenance settings](../attestations.md),
signing behavior, and [multi-platform manifest handling](../multi-platform.md).
Docker GitHub Builder moves that implementation into Docker-maintained reusable
workflows, so your workflow only decides when to build and which inputs to pass.

The difference is easiest to see in the job definition. A conventional workflow
spells out each action step:

```yaml
jobs:
  docker:
    runs-on: ubuntu-latest
    steps:
      - name: Login to Docker Hub
        uses: docker/login-action@{{% param "login_action_version" %}}
        with:
          username: ${{ vars.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}

      - name: Set up QEMU
        uses: docker/setup-qemu-action@{{% param "setup_qemu_action_version" %}}
      
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@{{% param "setup_buildx_action_version" %}}
        
      - name: Docker meta
        uses: docker/metadata-action@{{% param "metadata_action_version" %}}
        id: meta
        with:
          images: name/app

      - name: Build and push
        uses: docker/build-push-action@{{% param "build_push_action_version" %}}
        with:
          push: ${{ github.event_name != 'pull_request' }}
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
          cache-from: type=gha
          cache-to: type=gha
```

With Docker GitHub Builder, the same build is a reusable workflow call:

```yaml
jobs:
  build:
    uses: docker/github-builder/.github/workflows/build.yml@{{% param "github_builder_version" %}}
    permissions:
      contents: read # to fetch the repository content
      id-token: write # for signing attestation(s) with GitHub OIDC Token
    with:
      output: image
      push: ${{ github.event_name != 'pull_request' }}
      meta-images: name/app
    secrets:
      registry-auths: |
        - registry: docker.io
          username: ${{ vars.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}
```

This model gives you a build pipeline that is maintained in the Docker
organization, uses a pinned [BuildKit](../../../buildkit/_index.md) environment,
distributes [multi-platform builds](../../../building/multi-platform.md) across
runners when that helps, and emits signed [SLSA provenance](../../../metadata/attestations/slsa-provenance.md)
that records both the source commit and the builder identity.

That tradeoff is intentional. You keep control of when the build runs and which
inputs it uses, but the build implementation itself lives in the
Docker-maintained workflow rather than in per-repository job steps.

Use [`build.yml`](build.md) when your repository builds from a Dockerfile and
the familiar `build-push-action` inputs map cleanly to your workflow. Use
[`bake.yml`](bake.md) when your repository already describes builds in a
[Bake definition](../../../bake/_index.md), or when you want Bake targets,
overrides, and variables to stay as the source of truth.

Both workflows support image output, local output, cache export to the
[GitHub Actions cache backend](../../../cache/backends/gha.md),
[SBOM generation](../../../metadata/attestations/sbom.md), and signing. The
Bake workflow adds Bake definition validation and builds one target per workflow
call.

{{% sectionlinks %}}
