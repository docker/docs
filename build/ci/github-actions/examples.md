---
title: Example workflows
description: Docker GitHub Actions workflow examples.
keywords: CI, GitHub Actions, examples
---

This page showcases different examples of how you can customize and use the
Docker GitHub Actions in your CI pipelines.

## Push to multi-registries

The following workflow will connect you to Docker Hub and [GitHub Container Registry](https://github.com/docker/login-action#github-container-registry){: target="_blank" rel="noopener" class="_" }
and push the image to both registries:

{% raw %}
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
      -
        name: Checkout
        uses: actions/checkout@v3
      -
        name: Set up QEMU
        uses: docker/setup-qemu-action@v2
      -
        name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
      -
        name: Login to Docker Hub
        uses: docker/login-action@v2
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}
      -
        name: Login to GitHub Container Registry
        uses: docker/login-action@v2
        with:
          registry: ghcr.io
          username: ${{ github.repository_owner }}
          password: ${{ secrets.GITHUB_TOKEN }}
      -
        name: Build and push
        uses: docker/build-push-action@v3
        with:
          context: .
          platforms: linux/amd64,linux/arm64
          push: true
          tags: |
            user/app:latest
            user/app:1.0.0
            ghcr.io/user/app:latest
            ghcr.io/user/app:1.0.0
```
{% endraw %}

## Manage tags and labels

If you want an "automatic" tag management and [OCI Image Format Specification](https://github.com/opencontainers/image-spec/blob/master/annotations.md){: target="_blank" rel="noopener" class="_" }
for labels, you can do it in a dedicated setup step. The following workflow
will use the [Docker Metadata Action](https://github.com/docker/metadata-action){: target="_blank" rel="noopener" class="_" }
to handle tags and labels based on GitHub Actions events and Git metadata:

{% raw %}
```yaml
name: ci

on:
  schedule:
    - cron: "0 10 * * *"
  push:
    branches:
      - "**"
    tags:
      - "v*.*.*"
  pull_request:
    branches:
      - "main"

jobs:
  docker:
    runs-on: ubuntu-latest
    steps:
      -
        name: Checkout
        uses: actions/checkout@v3
      -
        name: Docker meta
        id: meta
        uses: docker/metadata-action@v4
        with:
          # list of Docker images to use as base name for tags
          images: |
            name/app
            ghcr.io/username/app
          # generate Docker tags based on the following events/attributes
          tags: |
            type=schedule
            type=ref,event=branch
            type=ref,event=pr
            type=semver,pattern={{version}}
            type=semver,pattern={{major}}.{{minor}}
            type=semver,pattern={{major}}
            type=sha
      -
        name: Set up QEMU
        uses: docker/setup-qemu-action@v2
      -
        name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
      -
        name: Login to Docker Hub
        if: github.event_name != 'pull_request'
        uses: docker/login-action@v2
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}
      -
        name: Login to GHCR
        if: github.event_name != 'pull_request'
        uses: docker/login-action@v2
        with:
          registry: ghcr.io
          username: ${{ github.repository_owner }}
          password: ${{ secrets.GITHUB_TOKEN }}
      -
        name: Build and push
        uses: docker/build-push-action@v3
        with:
          context: .
          push: ${{ github.event_name != 'pull_request' }}
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
```
{% endraw %}

## Multi-platform images

You can build [multi-platform images](../../building/multi-platform.md) using
the `platforms` option, as described in the following example.

> **Note**
>
> - For a list of available platforms, see the [Docker Setup Buildx](https://github.com/marketplace/actions/docker-setup-buildx){: target="_blank" rel="noopener" class="_" }
>   action.
> - If you want support for more platforms, you can use QEMU with the [Docker Setup QEMU](https://github.com/docker/setup-qemu-action){: target="_blank" rel="noopener" class="_" }
>   action.

{% raw %}
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
      -
        name: Checkout
        uses: actions/checkout@v3
      -
        name: Set up QEMU
        uses: docker/setup-qemu-action@v2
      -
        name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
      -
        name: Login to Docker Hub
        uses: docker/login-action@v2
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}
      -
        name: Build and push
        uses: docker/build-push-action@v3
        with:
          context: .
          platforms: linux/amd64,linux/arm64
          push: true
          tags: user/app:latest
```

{% endraw %}

## Cache

This page contains examples on using the cache storage backends with GitHub
actions.

> **Note**
>
> See [Cache storage backends](../../building/cache/backends/index.md) for more
> details about cache storage backends.

### Inline cache

In most cases you want to use the [inline cache exporter](../../building/cache/backends/inline.md).
However, note that the `inline` cache exporter only supports `min` cache mode.
To use `max` cache mode, push the image and the cache separately using the
registry cache exporter with the `cache-to` option, as shown in the [registry cache example](#registry-cache).

{% raw %}
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
      -
        name: Checkout
        uses: actions/checkout@v3
      -
        name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
      -
        name: Login to Docker Hub
        uses: docker/login-action@v2
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}
      -
        name: Build and push
        uses: docker/build-push-action@v3
        with:
          context: .
          push: true
          tags: user/app:latest
          cache-from: type=registry,ref=user/app:latest
          cache-to: type=inline
```
{% endraw %}

### Registry cache

You can import/export cache from a cache manifest or (special) image
configuration on the registry with the [registry cache exporter](../../building/cache/backends/registry.md).

{% raw %}
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
      -
        name: Checkout
        uses: actions/checkout@v3
      -
        name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
      -
        name: Login to Docker Hub
        uses: docker/login-action@v2
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}
      -
        name: Build and push
        uses: docker/build-push-action@v3
        with:
          context: .
          push: true
          tags: user/app:latest
          cache-from: type=registry,ref=user/app:buildcache
          cache-to: type=registry,ref=user/app:buildcache,mode=max
```
{% endraw %}

### GitHub cache

#### Cache backend API

> **Warning**
>
> This cache exporter is experimental. Please provide feedback on [BuildKit repository](https://github.com/moby/buildkit){: target="_blank" rel="noopener" class="_" }
> if you experience any issues.
{: .warning }

The [GitHub Actions cache exporter](../../building/cache/backends/gha.md)
backend uses the [GitHub Cache API](https://github.com/tonistiigi/go-actions-cache/blob/master/api.md)
to fetch and upload cache blobs. That's why you should only use this cache
backend in a GitHub Action workflow, as the `url` (`$ACTIONS_CACHE_URL`) and
`token` (`$ACTIONS_RUNTIME_TOKEN`) attributes only get populated in a workflow
context.

{% raw %}
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
      -
        name: Checkout
        uses: actions/checkout@v3
      -
        name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
      -
        name: Login to Docker Hub
        uses: docker/login-action@v2
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}
      -
        name: Build and push
        uses: docker/build-push-action@v3
        with:
          context: .
          push: true
          tags: user/app:latest
          cache-from: type=gha
          cache-to: type=gha,mode=max
```
{% endraw %}

#### Local cache

> **Warning**
>
> At the moment, old cache entries aren't deleted, so the cache size [keeps growing](https://github.com/docker/build-push-action/issues/252){:target="_blank" rel="noopener" class="_"}.
> The following example uses the `Move cache` step as a workaround (see [`moby/buildkit#1896`](https://github.com/moby/buildkit/issues/1896){:target="_blank" rel="noopener" class="_"}
> for more info).
{: .warning }

You can also leverage [GitHub cache](https://docs.github.com/en/actions/using-workflows/caching-dependencies-to-speed-up-workflows)
using the [actions/cache](https://github.com/actions/cache) and [local cache exporter](../../building/cache/backends/local.md)
with this action:

{% raw %}
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
      -
        name: Checkout
        uses: actions/checkout@v3
      -
        name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
      -
        name: Cache Docker layers
        uses: actions/cache@v3
        with:
          path: /tmp/.buildx-cache
          key: ${{ runner.os }}-buildx-${{ github.sha }}
          restore-keys: |
            ${{ runner.os }}-buildx-
      -
        name: Login to Docker Hub
        uses: docker/login-action@v2
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}
      -
        name: Build and push
        uses: docker/build-push-action@v3
        with:
          context: .
          push: true
          tags: user/app:latest
          cache-from: type=local,src=/tmp/.buildx-cache
          cache-to: type=local,dest=/tmp/.buildx-cache-new,mode=max
      -
        # Temp fix
        # https://github.com/docker/build-push-action/issues/252
        # https://github.com/moby/buildkit/issues/1896
        name: Move cache
        run: |
          rm -rf /tmp/.buildx-cache
          mv /tmp/.buildx-cache-new /tmp/.buildx-cache
```
{% endraw %}

## Secrets

In the following example uses and exposes the [`GITHUB_TOKEN` secret](https://docs.github.com/en/actions/security-guides/automatic-token-authentication#about-the-github_token-secret){:target="_blank" rel="noopener" class="_"}
as provided by GitHub in your workflow.

First, create a `Dockerfile` that uses the secret:

```dockerfile
# syntax=docker/dockerfile:1
FROM alpine
RUN --mount=type=secret,id=github_token \
  cat /run/secrets/github_token
```

In this example, the secret name is `github_token`. The following workflow
exposes this secret using the `secrets` input:

{% raw %}
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
      -
        name: Checkout
        uses: actions/checkout@v3
      -
        name: Set up QEMU
        uses: docker/setup-qemu-action@v2
      -
        name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
      -
        name: Build
        uses: docker/build-push-action@v3
        with:
          context: .
          platforms: linux/amd64,linux/arm64
          tags: user/app:latest
          secrets: |
            "github_token=${{ secrets.GITHUB_TOKEN }}"
```

{% endraw %}

> **Note**
>
> You can also expose a secret file to the build with the `secret-files` input:
>
> ```yaml
> secret-files: |
>   "MY_SECRET=./secret.txt"
> ```

If you're using [GitHub secrets](https://docs.github.com/en/actions/security-guides/encrypted-secrets){:target="_blank" rel="noopener" class="_"}
and need to handle multi-line value, you will need to place the key-value pair
between quotes:

{% raw %}
```yaml
secrets: |
  "MYSECRET=${{ secrets.GPG_KEY }}"
  GIT_AUTH_TOKEN=abcdefghi,jklmno=0123456789
  "MYSECRET=aaaaaaaa
  bbbbbbb
  ccccccccc"
  FOO=bar
  "EMPTYLINE=aaaa

  bbbb
  ccc"
  "JSON_SECRET={""key1"":""value1"",""key2"":""value2""}"
```
{% endraw %}

| Key              | Value                               |
|------------------|-------------------------------------|
| `MYSECRET`       | `***********************`           |
| `GIT_AUTH_TOKEN` | `abcdefghi,jklmno=0123456789`       |
| `MYSECRET`       | `aaaaaaaa\nbbbbbbb\nccccccccc`      |
| `FOO`            | `bar`                               |
| `EMPTYLINE`      | `aaaa\n\nbbbb\nccc`                 |
| `JSON_SECRET`    | `{"key1":"value1","key2":"value2"}` |

> **Note**
>
> Double escapes are needed for quote signs.

## Export image to Docker

You may want your build result to be available in the Docker client through
`docker images` to be able to use it in another step of your workflow:

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
      -
        name: Checkout
        uses: actions/checkout@v3
      -
        name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
      -
        name: Build
        uses: docker/build-push-action@v3
        with:
          context: .
          load: true
          tags: myimage:latest
      -
        name: Inspect
        run: |
          docker image inspect myimage:latest
```

## Test your image before pushing it

In some cases, you might want to validate that the image works as expected
before pushing it.

The following workflow implements several steps to achieve this:

- Build and export the image to Docker
- Test your image
- Multi-platform build and push the image

{% raw %}
```yaml
name: ci

on:
  push:
    branches:
      - "main"

env:
  TEST_TAG: user/app:test
  LATEST_TAG: user/app:latest

jobs:
  docker:
    runs-on: ubuntu-latest
    steps:
      -
        name: Checkout
        uses: actions/checkout@v3
      -
        name: Set up QEMU
        uses: docker/setup-qemu-action@v2
      -
        name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
      -
        name: Login to Docker Hub
        uses: docker/login-action@v2
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}
      -
        name: Build and export to Docker
        uses: docker/build-push-action@v3
        with:
          context: .
          load: true
          tags: ${{ env.TEST_TAG }}
      -
        name: Test
        run: |
          docker run --rm ${{ env.TEST_TAG }}
      -
        name: Build and push
        uses: docker/build-push-action@v3
        with:
          context: .
          platforms: linux/amd64,linux/arm64
          push: true
          tags: ${{ env.LATEST_TAG }}
```
{% endraw %}

> **Note**
>
> This workflow doesn't actually build the `linux/amd64` image twice. The image
> is built once, and the following steps uses the internal cache for from the
> first `Build and push` step. The second `Build and push` step only builds
> `linux/arm64`.

## Local registry

For testing purposes you may need to create a [local registry](https://hub.docker.com/_/registry){: target="_blank" rel="noopener" class="_" }
to push images into:

```yaml
name: ci

on:
  push:
    branches:
      - "main"

jobs:
  docker:
    runs-on: ubuntu-latest
    services:
      registry:
        image: registry:2
        ports:
          - 5000:5000
    steps:
      -
        name: Checkout
        uses: actions/checkout@v3
      -
        name: Set up QEMU
        uses: docker/setup-qemu-action@v2
      -
        name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
        with:
          driver-opts: network=host
      -
        name: Build and push to local registry
        uses: docker/build-push-action@v3
        with:
          context: .
          push: true
          tags: localhost:5000/name/app:latest
      -
        name: Inspect
        run: |
          docker buildx imagetools inspect localhost:5000/name/app:latest
```

## Share built image between jobs

As each job is isolated in its own runner, you can't use your built image
between jobs, except if you're using [self-hosted runners](https://docs.github.com/en/actions/hosting-your-own-runners/about-self-hosted-runners){: target="_blank" rel="noopener" class="_" }.
However, you can [pass data between jobs](https://docs.github.com/en/actions/using-workflows/storing-workflow-data-as-artifacts#passing-data-between-jobs-in-a-workflow){: target="_blank" rel="noopener" class="_" }
in a workflow using the [actions/upload-artifact](https://github.com/actions/upload-artifact){: target="_blank" rel="noopener" class="_" }
and [actions/download-artifact](https://github.com/actions/download-artifact){: target="_blank" rel="noopener" class="_" }
actions:

```yaml
name: ci

on:
  push:
    branches:
      - "main"

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      -
        name: Checkout
        uses: actions/checkout@v3
      -
        name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
      -
        name: Build and export
        uses: docker/build-push-action@v3
        with:
          context: .
          tags: myimage:latest
          outputs: type=docker,dest=/tmp/myimage.tar
      -
        name: Upload artifact
        uses: actions/upload-artifact@v3
        with:
          name: myimage
          path: /tmp/myimage.tar

  use:
    runs-on: ubuntu-latest
    needs: build
    steps:
      -
        name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
      -
        name: Download artifact
        uses: actions/download-artifact@v3
        with:
          name: myimage
          path: /tmp
      -
        name: Load image
        run: |
          docker load --input /tmp/myimage.tar
          docker image ls -a
```

## Named contexts

You can define [additional build contexts](../../../engine/reference/commandline/buildx_build.md#build-context),
and access them in your Dockerfile with `FROM name` or `--from=name`. When
Dockerfile defines a stage with the same name it's overwritten.

This can be useful with GitHub Actions to reuse results from other builds or pin
an image to a specific tag in your workflow.

### Pin image to a tag

Replace `alpine:latest` with a pinned one:

```dockerfile
# syntax=docker/dockerfile:1
FROM alpine
RUN echo "Hello World"
```

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
      -
        name: Checkout
        uses: actions/checkout@v3
      -
        name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
      -
        name: Build
        uses: docker/build-push-action@v3
        with:
          context: .
          build-contexts: |
            alpine=docker-image://alpine:3.16
          tags: myimage:latest
```

### Use image in subsequent steps

By default, the [Docker Setup Buildx](https://github.com/marketplace/actions/docker-setup-buildx){: target="_blank" rel="noopener" class="_" }
action uses `docker-container` as a build driver, so built Docker images aren't
loaded automatically.

With named contexts you can reuse the built image:

```dockerfile
# syntax=docker/dockerfile:1
FROM alpine
RUN echo "Hello World"
```

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
      -
        name: Checkout
        uses: actions/checkout@v3
      -
        name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
      -
        name: Build base image
        uses: docker/build-push-action@v3
        with:
          context: base
          load: true
          tags: my-base-image:latest
      -
        name: Build
        uses: docker/build-push-action@v3
        with:
          context: .
          build-contexts: |
            alpine=docker-image://my-base-image:latest
          tags: myimage:latest
```

## Builder configuration

This section contains instructions on configuring your BuildKit build instances
when using GitHub Actions.

### Append additional nodes to the builder

Buildx supports running builds on multiple machines. This is useful for building
[multi-platform images](../../building/multi-platform.md) on native nodes for
more complicated cases that aren't handled by QEMU. Building on native nodes
generally has better performance, and allows you to distribute the build across
multiple machines.

You can append nodes to the builder you're creating using the `append` option.
It takes input in the form of a YAML string document to remove limitations
intrinsically linked to GitHub Actions: you can only use strings in the input
fields:

| Name              | Type   | Description                                                                                                                                                                                                                                                             |
|-------------------|--------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `name`            | String | [Name of the node](../../../engine/reference/commandline/buildx_create.md#node). If empty, it's the name of the builder it belongs to, with an index number suffix. This is useful to set it if you want to modify/remove a node in an underlying step of you workflow. |
| `endpoint`        | String | [Docker context or endpoint](../../../engine/reference/commandline/buildx_create.md#description) of the node to add to the builder                                                                                                                                      |
| `driver-opts`     | List   | List of additional [driver-specific options](../../../engine/reference/commandline/buildx_create.md#driver-opt)                                                                                                                                                         |
| `buildkitd-flags` | String | [Flags for buildkitd](../../../engine/reference/commandline/buildx_create.md#buildkitd-flags) daemon                                                                                                                                                                    |
| `platforms`       | String | Fixed [platforms](../../../engine/reference/commandline/buildx_create.md#platform) for the node. If not empty, values take priority over the detected ones.                                                                                                             |

Here is an example using remote nodes with the [`remote` driver](../../building/drivers/remote.md)
and [TLS authentication](#tls-authentication):

{% raw %}
```yaml
name: ci

on:
  push:

jobs:
  buildx:
    runs-on: ubuntu-latest
    steps:
      -
        name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
        with:
          driver: remote
          endpoint: tcp://oneprovider:1234
          append: |
            - endpoint: tcp://graviton2:1234
              platforms: linux/arm64
            - endpoint: tcp://linuxone:1234
              platforms: linux/s390x
        env:
          BUILDER_NODE_0_AUTH_TLS_CACERT: ${{ secrets.ONEPROVIDER_CA }}
          BUILDER_NODE_0_AUTH_TLS_CERT: ${{ secrets.ONEPROVIDER_CERT }}
          BUILDER_NODE_0_AUTH_TLS_KEY: ${{ secrets.ONEPROVIDER_KEY }}
          BUILDER_NODE_1_AUTH_TLS_CACERT: ${{ secrets.GRAVITON2_CA }}
          BUILDER_NODE_1_AUTH_TLS_CERT: ${{ secrets.GRAVITON2_CERT }}
          BUILDER_NODE_1_AUTH_TLS_KEY: ${{ secrets.GRAVITON2_KEY }}
          BUILDER_NODE_2_AUTH_TLS_CACERT: ${{ secrets.LINUXONE_CA }}
          BUILDER_NODE_2_AUTH_TLS_CERT: ${{ secrets.LINUXONE_CERT }}
          BUILDER_NODE_2_AUTH_TLS_KEY: ${{ secrets.LINUXONE_KEY }}
```
{% endraw %}

### Authentication for remote builders

The following examples show how to handle authentication for remote builders,
using SSH or TLS.

#### SSH authentication

To be able to connect to an SSH endpoint using the [`docker-container` driver](../../building/drivers/docker-container.md),
you have to set up the SSH private key and configuration on the GitHub Runner:

{% raw %}
```yaml
name: ci

on:
  push:

jobs:
  buildx:
    runs-on: ubuntu-latest
    steps:
      -
        name: Set up SSH
        uses: MrSquaare/ssh-setup-action@523473d91581ccbf89565e12b40faba93f2708bd # v1.1.0
        with:
          host: graviton2
          private-key: ${{ secrets.SSH_PRIVATE_KEY }}
          private-key-name: aws_graviton2
      -
        name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
        with:
          endpoint: ssh://me@graviton2
```
{% endraw %}

#### TLS authentication

You can also [set up a remote BuildKit instance](../../building/drivers/remote.md#example-remote-buildkit-in-docker-container)
using the remote driver. To ease the integration in your workflow, you can use
an environment variables that sets up authentication using the BuildKit client
certificates for the `tcp://`:

- `BUILDER_NODE_<idx>_AUTH_TLS_CACERT`
- `BUILDER_NODE_<idx>_AUTH_TLS_CERT`
- `BUILDER_NODE_<idx>_AUTH_TLS_KEY`

The `<idx>` placeholder is the position of the node in the list of nodes.

{% raw %}
```yaml
name: ci

on:
  push:

jobs:
  buildx:
    runs-on: ubuntu-latest
    steps:
      -
        name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
        with:
          driver: remote
          endpoint: tcp://graviton2:1234
        env:
          BUILDER_NODE_0_AUTH_TLS_CACERT: ${{ secrets.GRAVITON2_CA }}
          BUILDER_NODE_0_AUTH_TLS_CERT: ${{ secrets.GRAVITON2_CERT }}
          BUILDER_NODE_0_AUTH_TLS_KEY: ${{ secrets.GRAVITON2_KEY }}
```
{% endraw %}

### Daemon configuration

You can provide a [BuildKit configuration](../../buildkit/toml-configuration.md)
to your builder if you're using the [`docker-container` driver](../../building/drivers/docker-container.md)
(default) with the `config` or `config-inline` inputs:

### Registry mirror

You can configure a registry mirror using an inline block directly in your
workflow with the `config-inline` input:

```yaml
name: ci

on:
  push:

jobs:
  buildx:
    runs-on: ubuntu-latest
    steps:
      -
        name: Checkout
        uses: actions/checkout@v3
      -
        name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
        with:
          config-inline: |
            [registry."docker.io"]
              mirrors = ["mirror.gcr.io"]
```

For more information about using a registry mirror, see [Registry mirror](../../buildkit/configure.md#registry-mirror).

#### Max parallelism

You can limit the parallelism of the BuildKit solver which is particularly
useful for low-powered machines.

You can use the `config-inline` input like the previous example, or you can use
a dedicated BuildKit config file from your repository if you want with the
`config` input:

```toml
# .github/buildkitd.toml
[worker.oci]
  max-parallelism = 4
```

```yaml
name: ci

on:
  push:

jobs:
  buildx:
    runs-on: ubuntu-latest
    steps:
      -
        name: Checkout
        uses: actions/checkout@v3
      -
        name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
        with:
          config: .github/buildkitd.toml
```

### Standalone mode

If you don't have the Docker CLI installed on the GitHub Runner, the Buildx
binary gets invoked directly, instead of calling it as a Docker CLI plugin. This
can be useful if you want to use the `kubernetes` driver in your self-hosted
runner:

```yaml
name: ci

on:
  push:

jobs:
  buildx:
    runs-on: ubuntu-latest
    steps:
      -
        name: Checkout
        uses: actions/checkout@v3
      -
        name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
        with:
          driver: kubernetes
      -
        name: Build
        run: |
          buildx build .
```

## Isolated builders

The following example shows how you can select different builders for different
jobs.

An example scenario where this might be useful is when you are using a monorepo,
and you want to pinpoint different packages to specific builders. For example,
some packages may be particularly resource-intensive to build and require more
compute. Or they require a builder equipped with a particular capability or
hardware.

For more information about remote builder, see [`remote` driver](../../building/drivers/remote.md)
and the [append builder nodes example](#append-additional-nodes-to-the-builder).

{% raw %}
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
      -
        name: Checkout
        uses: actions/checkout@v3
      -
        uses: docker/setup-buildx-action@v2
        id: builder1
      -
        uses: docker/setup-buildx-action@v2
        id: builder2
      -
        name: Builder 1 name
        run: echo ${{ steps.builder1.outputs.name }}
      -
        name: Builder 2 name
        run: echo ${{ steps.builder2.outputs.name }}
      -
        name: Build against builder1
        uses: docker/build-push-action@v3
        with:
          builder: ${{ steps.builder1.outputs.name }}
          context: .
          target: mytarget1
      -
        name: Build against builder2
        uses: docker/build-push-action@v3
        with:
          builder: ${{ steps.builder2.outputs.name }}
          context: .
          target: mytarget2
```
{% endraw %}

## Copy images between registries

[Multi-platform images](../../building/multi-platform.md) built using Buildx can
be copied from one registry to another using the [`buildx imagetools create` command](../../../engine/reference/commandline/buildx_imagetools_create.md):

{% raw %}
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
      -
        name: Checkout
        uses: actions/checkout@v3
      -
        name: Set up QEMU
        uses: docker/setup-qemu-action@v2
      -
        name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
      -
        name: Login to Docker Hub
        uses: docker/login-action@v2
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}
      -
        name: Login to GitHub Container Registry
        uses: docker/login-action@v2
        with:
          registry: ghcr.io
          username: ${{ github.repository_owner }}
          password: ${{ secrets.GITHUB_TOKEN }}
      -
        name: Build and push
        uses: docker/build-push-action@v3
        with:
          context: .
          platforms: linux/amd64,linux/arm64
          push: true
          tags: |
            user/app:latest
            user/app:1.0.0
      -
        name: Push image to GHCR
        run: |
          docker buildx imagetools create \
            --tag ghcr.io/user/app:latest \
            --tag ghcr.io/user/app:1.0.0 \
            user/app:latest
```
{% endraw %}

## Update Docker Hub repository description

You can update the Docker Hub repository description using a third party action
called [Docker Hub Description](https://github.com/peter-evans/dockerhub-description){: target="_blank" rel="noopener" class="_" }
with this action:

{% raw %}
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
      -
        name: Checkout
        uses: actions/checkout@v3
      -
        name: Set up QEMU
        uses: docker/setup-qemu-action@v2
      -
        name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
      -
        name: Login to Docker Hub
        uses: docker/login-action@v2
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}
      -
        name: Build and push
        uses: docker/build-push-action@v3
        with:
          context: .
          push: true
          tags: user/app:latest
      -
        name: Update repo description
        uses: peter-evans/dockerhub-description@v2
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_PASSWORD }}
          repository: user/app
```
{% endraw %}
