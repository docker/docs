---
title: Bake targets
linkTitle: Targets
weight: 20
description: Learn how to define and use targets in Bake
keywords: bake, target, targets, buildx, docker, buildkit, default
---

A target in a Bake file represents a build invocation. It holds all the
information you would normally pass to a `docker build` command using flags.

```hcl {title=docker-bake.hcl}
target "webapp" {
  dockerfile = "webapp.Dockerfile"
  tags = ["docker.io/username/webapp:latest"]
  context = "https://github.com/username/webapp"
}
```

To build a target with Bake, pass name of the target to the `bake` command.

```console
$ docker buildx bake webapp
```

You can build multiple targets at once by passing multiple target names to the
`bake` command.

```console
$ docker buildx bake webapp api tests
```

## Default target

If you don't specify a target when running `docker buildx bake`, Bake will
build the target named `default`.

```hcl {title=docker-bake.hcl}
target "default" {
  dockerfile = "webapp.Dockerfile"
  tags = ["docker.io/username/webapp:latest"]
  context = "https://github.com/username/webapp"
}
```

To build this target, run `docker buildx bake` without any arguments:

```console
$ docker buildx bake
```

## Target properties

The properties you can set for a target closely resemble the CLI flags for
`docker build`, with a few additional properties that are specific to Bake.

For all the properties you can set for a target, see the [Bake reference](/build/bake/reference#target).

## Grouping targets

You can group targets together using the `group` block. This is useful when you
want to build multiple targets at once.

```hcl {title=docker-bake.hcl}
group "all" {
  targets = ["webapp", "api", "tests"]
}

target "webapp" {
  dockerfile = "webapp.Dockerfile"
  tags = ["docker.io/username/webapp:latest"]
  context = "https://github.com/username/webapp"
}

target "api" {
  dockerfile = "api.Dockerfile"
  tags = ["docker.io/username/api:latest"]
  context = "https://github.com/username/api"
}

target "tests" {
  dockerfile = "tests.Dockerfile"
  contexts = {
    webapp = "target:webapp",
    api = "target:api",
  }
  output = ["type=local,dest=build/tests"]
  context = "."
}
```

To build all the targets in a group, pass the name of the group to the `bake`
command.

```console
$ docker buildx bake all
```

## Additional resources

Refer to the following pages to learn more about Bake's features:

- Learn how to use [variables](./variables.md) in Bake to make your build
  configuration more flexible.
- Learn how you can use matrices to build multiple images with different
  configurations in [Matrices](./matrices.md).
- Head to the [Bake file reference](/build/bake/reference/) to learn about all
  the properties you can set in a Bake file, and its syntax.
