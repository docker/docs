---
title: Cache management with GitHub Actions
keywords: ci, github actions, gha, buildkit, buildx, cache
---

This page contains examples on using the cache storage backends with GitHub
Actions.

> **Note**
>
> See [Cache storage backends](../../cache/backends/index.md) for more
> details about cache storage backends.

## Inline cache

In most cases you want to use the [inline cache exporter](../../cache/backends/inline.md).
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
        uses: docker/build-push-action@v4
        with:
          context: .
          push: true
          tags: user/app:latest
          cache-from: type=registry,ref=user/app:latest
          cache-to: type=inline
```
{% endraw %}

## Registry cache

You can import/export cache from a cache manifest or (special) image
configuration on the registry with the [registry cache exporter](../../cache/backends/registry.md).

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
        uses: docker/build-push-action@v4
        with:
          context: .
          push: true
          tags: user/app:latest
          cache-from: type=registry,ref=user/app:buildcache
          cache-to: type=registry,ref=user/app:buildcache,mode=max
```
{% endraw %}

## GitHub cache

### Cache backend API

> Experimental
>
> This cache exporter is experimental. Please provide feedback on [BuildKit repository](https://github.com/moby/buildkit){:target="blank" rel="noopener" class=""}
> if you experience any issues.
{: .experimental }

The [GitHub Actions cache exporter](../../cache/backends/gha.md)
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
        uses: docker/build-push-action@v4
        with:
          context: .
          push: true
          tags: user/app:latest
          cache-from: type=gha
          cache-to: type=gha,mode=max
```
{% endraw %}

### Local cache

> **Warning**
>
> At the moment, old cache entries aren't deleted, so the cache size [keeps growing](https://github.com/docker/build-push-action/issues/252){:target="blank" rel="noopener" class=""}.
> The following example uses the `Move cache` step as a workaround (see [`moby/buildkit#1896`](https://github.com/moby/buildkit/issues/1896){:target="blank" rel="noopener" class=""}
> for more info).
{: .warning }

You can also leverage [GitHub cache](https://docs.github.com/en/actions/using-workflows/caching-dependencies-to-speed-up-workflows){:target="blank" rel="noopener" class=""}
using the [actions/cache](https://github.com/actions/cache) and [local cache exporter](../../cache/backends/local.md)
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
        uses: docker/build-push-action@v4
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
