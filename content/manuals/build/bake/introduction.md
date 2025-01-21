---
title: Introduction to Bake
linkTitle: Introduction
weight: 10
description: Get started with using Bake to build your project
keywords: bake, quickstart, build, project, introduction, getting started
---

Bake is an abstraction for the `docker build` command that lets you more easily
manage your build configuration (CLI flags, environment variables, etc.) in a
consistent way for everyone on your team.

Bake is a command built into the Buildx CLI, so as long as you have Buildx
installed, you also have access to bake, via the `docker buildx bake` command.

## Building a project with Bake

Here's a simple example of a `docker build` command:

```console
$ docker build -f Dockerfile -t myapp:latest .
```

This command builds the Dockerfile in the current directory and tags the
resulting image as `myapp:latest`.

To express the same build configuration using Bake:

```hcl {title=docker-bake.hcl}
target "myapp" {
  context = "."
  dockerfile = "Dockerfile"
  tags = ["myapp:latest"]
}
```

Bake provides a structured way to manage your build configuration, and it saves
you from having to remember all the CLI flags for `docker build` every time.
With this file, building the image is as simple as running:

```console
$ docker buildx bake myapp
```

For simple builds, the difference between `docker build` and `docker buildx
bake` is minimal. However, as your build configuration grows more complex, Bake
provides a more structured way to manage that complexity, that would be
difficult to manage with CLI flags for the `docker build`. It also provides a
way to share build configurations across your team, so that everyone is
building images in a consistent way, with the same configuration.

## The Bake file format

You can write Bake files in HCL, YAML (Docker Compose files), or JSON. In
general, HCL is the most expressive and flexible format, which is why you'll
see it used in most of the examples in this documentation, and in projects that
use Bake.

The properties that can be set for a target closely resemble the CLI flags for
`docker build`. For instance, consider the following `docker build` command:

```console
$ docker build \
  -f Dockerfile \
  -t myapp:latest \
  --build-arg foo=bar \
  --no-cache \
  --platform linux/amd64,linux/arm64 \
  .
```

The Bake equivalent would be:

```hcl {title=docker-bake.hcl}
target "myapp" {
  context = "."
  dockerfile = "Dockerfile"
  tags = ["myapp:latest"]
  args = {
    foo = "bar"
  }
  no-cache = true
  platforms = ["linux/amd64", "linux/arm64"]
}
```

## Next steps

To learn more about using Bake, see the following topics:

- Learn how to define and use [targets](./targets.md) in Bake
- To see all the properties that can be set for a target, refer to the
  [Bake file reference](/build/bake/reference/).
