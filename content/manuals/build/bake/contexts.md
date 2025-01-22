---
title: Using Bake with additional contexts
linkTitle: Contexts
weight: 80
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
command](/reference/cli/docker/buildx/build.md#build-context).

Inside the Dockerfile these contexts can be used with the `FROM` instruction or
`--from` flag.

Supported context values are:

- Local filesystem directories
- Container images
- Git URLs
- HTTP URLs
- Name of another target in the Bake file

## Pinning alpine image

```dockerfile {title=Dockerfile}
# syntax=docker/dockerfile:1
FROM alpine
RUN echo "Hello world"
```

```hcl {title=docker-bake.hcl}
target "app" {
  contexts = {
    alpine = "docker-image://alpine:3.13"
  }
}
```

## Using a secondary source directory

```dockerfile {title=Dockerfile}
FROM golang
COPY --from=src . .
```

```hcl {title=docker-bake.hcl}
# Running `docker buildx bake app` will result in `src` not pointing
# to some previous build stage but to the client filesystem, not part of the context.
target "app" {
  contexts = {
    src = "../path/to/source"
  }
}
```

## Using a target as a build context

To use a result of one target as a build context of another, specify the target
name with `target:` prefix.

```dockerfile {title=baseapp.Dockerfile}
FROM scratch
```
```dockerfile {title=Dockerfile}
# syntax=docker/dockerfile:1
FROM baseapp
RUN echo "Hello world"
```

```hcl {title=docker-bake.hcl}
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

## Deduplicate context transfer

> [!NOTE]
>
> As of Buildx version 0.17.0 and later, Bake automatically de-duplicates
> context transfer for targets that share the same context. In addition to
> Buildx version 0.17.0, the builder must be running BuildKit version 0.16.0 or
> later, and the Dockerfile syntax must be `docker/dockerfile:1.10` or later.
>
> If you meet these requirements, you don't need to manually de-duplicate
> context transfer as described in this section.
>
> - To check your Buildx version, run `docker buildx version`.
> - To check your BuildKit version, run `docker buildx inspect --bootstrap` and
>   look for the `BuildKit version` field.
> - To check your Dockerfile syntax version, check the `syntax`
>   [parser directive](/reference/dockerfile.md#syntax) in your Dockerfile. If
>   it's not present, the default version whatever comes bundled with your
>   current version of BuildKit. To set the version explicitly, add
>   `#syntax=docker/dockerfile:1.10` at the top of your Dockerfile.

When you build targets concurrently, using groups, build contexts are loaded
independently for each target. If the same context is used by multiple targets
in a group, that context is transferred once for each time it's used. This can
result in significant impact on build time, depending on your build
configuration. For example, say you have a Bake file that defines the following
group of targets:

```hcl {title=docker-bake.hcl}
group "default" {
  targets = ["target1", "target2"]
}

target "target1" {
  target = "target1"
  context = "."
}

target "target2" {
  target = "target2"
  context = "."
}
```

In this case, the context `.` is transferred twice when you build the default
group: once for `target1` and once for `target2`.

If your context is small, and if you are using a local builder, duplicate
context transfers may not be a big deal. But if your build context is big, or
you have a large number of targets, or you're transferring the context over a
network to a remote builder, context transfer becomes a performance bottleneck.

To avoid transferring the same context multiple times, you can define a named
context that only loads the context files, and have each target that needs
those files reference that named context. For example, the following Bake file
defines a named target `ctx`, which is used by both `target1` and `target2`:

```hcl {title=docker-bake.hcl}
group "default" {
  targets = ["target1", "target2"]
}

target "ctx" {
  context = "."
  target = "ctx"
}

target "target1" {
  target = "target1"
  contexts = {
    ctx = "target:ctx"
  }
}

target "target2" {
  target = "target2"
  contexts = {
    ctx = "target:ctx"
  }
}
```

The named context `ctx` represents a Dockerfile stage, which copies the files
from its context (`.`). Other stages in the Dockerfile can now reference the
`ctx` named context and, for example, mount its files with `--mount=from=ctx`.

```dockerfile {title=Dockerfile}
FROM scratch AS ctx
COPY --link . .

FROM golang:alpine AS target1
WORKDIR /work
RUN --mount=from=ctx \
    go build -o /out/client ./cmd/client \

FROM golang:alpine AS target2
WORKDIR /work
RUN --mount=from=ctx \
    go build -o /out/server ./cmd/server
```
