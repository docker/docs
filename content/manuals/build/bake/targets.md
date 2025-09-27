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
    webapp = "target:webapp"
    api = "target:api"
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

## Pattern matching for targets and groups

Bake supports shell-style wildcard patterns when specifying target or grouped targets.
This makes it easier to build multiple targets without listing each one explicitly.

Supported patterns:

- `*` matches any sequence of characters
- `?` matches any single character
- `[abc]` matches any character in brackets

> [!NOTE]
>
> Always wrap wildcard patterns in quotes. Without quotes, your shell will expand the
> wildcard to match files in the current directory, which usually causes errors.

Examples: 

```console
# Match all targets starting with 'foo-'
$ docker buildx bake "foo-*"

# Match all targets
$ docker buildx bake "*"

# Matches: foo-baz, foo-caz, foo-daz, etc.
$ docker buildx bake "foo-?az"

# Matches: foo-bar, boo-bar
$ docker buildx bake "[fb]oo-bar"

# Matches: mtx-a-b-d, mtx-a-b-e, mtx-a-b-f
$ docker buildx bake "mtx-a-b-*"
``` 

You can also combine multiple patterns:

```console
$ docker buildx bake "foo*" "tests"
```

## Additional resources

Refer to the following pages to learn more about Bake's features:

- Learn how to use [variables](./variables.md) in Bake to make your build
  configuration more flexible.
- Learn how you can use matrices to build multiple images with different
  configurations in [Matrices](./matrices.md).
- Head to the [Bake file reference](/build/bake/reference/) to learn about all
  the properties you can set in a Bake file, and its syntax.
