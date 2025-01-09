---
title: Reproducible builds with GitHub Actions
linkTitle: Reproducible builds
description: How to create reproducible builds in GitHub Actions using the SOURCE_EPOCH environment variable
keywords: build, buildx, github actions, ci, gha, reproducible builds, SOURCE_DATE_EPOCH
---

`SOURCE_DATE_EPOCH` is a [standardized environment variable][source_date_epoch]
for instructing build tools to produce a reproducible output.
Setting the environment variable for a build makes the timestamps in the
image index, config, and file metadata reflect the specified Unix time.

[source_date_epoch]: https://reproducible-builds.org/docs/source-date-epoch/

To set the environment variable in GitHub Actions,
use the built-in `env` property on the build step.

## Unix epoch timestamps

The following example sets the `SOURCE_DATE_EPOCH` variable to 0, Unix epoch.

{{< tabs group="action" >}}
{{< tab name="`docker/build-push-action`" >}}

```yaml
name: ci

on:
  push:

jobs:
  docker:
    runs-on: ubuntu-latest
    steps:
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Build
        uses: docker/build-push-action@v6
        with:
          tags: user/app:latest
        env:
          SOURCE_DATE_EPOCH: 0
```

{{< /tab >}}
{{< tab name="`docker/bake-action`" >}}

```yaml
name: ci

on:
  push:

jobs:
  docker:
    runs-on: ubuntu-latest
    steps:
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Build
        uses: docker/bake-action@v6
        env:
          SOURCE_DATE_EPOCH: 0
```

{{< /tab >}}
{{< /tabs >}}

## Git commit timestamps

The following example sets `SOURCE_DATE_EPOCH` to the Git commit timestamp.

{{< tabs group="action" >}}
{{< tab name="`docker/build-push-action`" >}}

```yaml
name: ci

on:
  push:

jobs:
  docker:
    runs-on: ubuntu-latest
    steps:
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Get Git commit timestamps
        run: echo "TIMESTAMP=$(git log -1 --pretty=%ct)" >> $GITHUB_ENV

      - name: Build
        uses: docker/build-push-action@v6
        with:
          tags: user/app:latest
        env:
          SOURCE_DATE_EPOCH: ${{ env.TIMESTAMP }}
```

{{< /tab >}}
{{< tab name="`docker/bake-action`" >}}

```yaml
name: ci

on:
  push:

jobs:
  docker:
    runs-on: ubuntu-latest
    steps:
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Get Git commit timestamps
        run: echo "TIMESTAMP=$(git log -1 --pretty=%ct)" >> $GITHUB_ENV

      - name: Build
        uses: docker/bake-action@v6
        env:
          SOURCE_DATE_EPOCH: ${{ env.TIMESTAMP }}
```

{{< /tab >}}
{{< /tabs >}}

## Additional information

For more information about the `SOURCE_DATE_EPOCH` support in BuildKit,
see [BuildKit documentation](https://github.com/moby/buildkit/blob/master/docs/build-repro.md#source_date_epoch).
