---
title: Using Bake with additional contexts
description: |
  Additional contexts are useful when you want to pin image versions,
  or reference the output of other targets
keywords: build, buildx, bake, buildkit, hcl
aliases:
  - /build/customize/bake/build-contexts/
  - /build/bake/build-contexts/
---

In addition to the main `context` key that defines the build context, each
target can also define additional named contexts with a map defined with key
`contexts`. These values map to the `--build-context` flag in the [build
command](../../reference/cli/docker/buildx/build.md#build-context).

Inside the Dockerfile these contexts can be used with the `FROM` instruction or
`--from` flag.

Supported context values are:

- Local filesystem directories
- Container images
- Git URLs
- HTTP URLs
- Name of another target in the Bake file

## Pinning alpine image

```dockerfile
# syntax=docker/dockerfile:1
FROM alpine
RUN echo "Hello world"
```

```hcl
# docker-bake.hcl
target "app" {
  contexts = {
    alpine = "docker-image://alpine:3.13"
  }
}
```

## Using a secondary source directory

```dockerfile
# syntax=docker/dockerfile:1
FROM scratch AS src

FROM golang
COPY --from=src . .
```

```hcl
# docker-bake.hcl
target "app" {
  contexts = {
    src = "../path/to/source"
  }
}
```

## Using a target as a build context

To use a result of one target as a build context of another, specify the target
name with `target:` prefix.

```dockerfile
# syntax=docker/dockerfile:1
FROM baseapp
RUN echo "Hello world"
```

```hcl
# docker-bake.hcl
target "base" {
  dockerfile = "baseapp.Dockerfile"
}

target "app" {
  contexts = {
    baseapp = "target:base"
  }
}
```

In most cases you should just use a single multi-stage Dockerfile with multiple
targets for similar behavior. This case is only recommended when you have
multiple Dockerfiles that can't be easily merged into one.
