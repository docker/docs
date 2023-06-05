---
title: Use a proxy with GitHub Actions
keywords: ci, github actions, gha, buildkit, buildx, registry
---

The Docker GitHub Actions use the
[`docker-container` driver](../../drivers/docker-container.md) by default,
which supports setting environment variables through driver options.
If you want your build to use a proxy, you can configure the driver options to
specify HTTP proxy environment variables.

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
        with:
          driver-opts: |
            env.http_proxy=http://proxy.example.com:3128
            env.https_proxy=https://proxy.example.com:3129
            env.no_proxy=*.test.example.com,.example.org,127.0.0.0/8
      -
        name: Build and push
        uses: docker/build-push-action@v4
        with:
          context: .
          platforms: linux/amd64,linux/arm64
          push: true
          tags: |
            user/app:latest
            user/app:1.0.0
```

The following environment variables are supported:

- `HTTP_PROXY`
- `HTTPS_PROXY`
- `http_proxy`
- `https_proxy`
- `NO_PROXY`
- `no_proxy`
