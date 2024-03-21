---
title: BuildKit
description: Introduction and overview of BuildKit
keywords: build, buildkit
---

## Overview

[BuildKit](https://github.com/moby/buildkit)
is an improved backend to replace the legacy builder. BuildKit is the default builder
for users on Docker Desktop, and Docker Engine as of version 23.0.

BuildKit provides new functionality and improves your builds' performance.
It also introduces support for handling more complex scenarios:

- Detect and skip executing unused build stages
- Parallelize building independent build stages
- Incrementally transfer only the changed files in your
  [build context](../building/context.md) between builds
- Detect and skip transferring unused files in your
  [build context](../building/context.md)
- Use [Dockerfile frontend](../dockerfile/frontend.md) implementations with many
  new features
- Avoid side effects with rest of the API (intermediate images and containers)
- Prioritize your build cache for automatic pruning

Apart from many new features, the main areas BuildKit improves on the current
experience are performance, storage management, and extensibility. From the
performance side, a significant update is a new fully concurrent build graph
solver. It can run build steps in parallel when possible and optimize out
commands that don't have an impact on the final result. We have also optimized
the access to the local source files. By tracking only the updates made to these
files between repeated build invocations, there is no need to wait for local
files to be read or uploaded before the work can begin.

## LLB

At the core of BuildKit is a
[Low-Level Build (LLB)](https://github.com/moby/buildkit#exploring-llb) definition format. LLB is an intermediate binary format
that allows developers to extend BuildKit. LLB defines a content-addressable
dependency graph that can be used to put together very complex build
definitions. It also supports features not exposed in Dockerfiles, like direct
data mounting and nested invocation.

{{< figure src="../images/buildkit-dag.svg" class="invertible" >}}

Everything about execution and caching of your builds is defined in LLB. The
caching model is entirely rewritten compared to the legacy builder. Rather than
using heuristics to compare images, LLB directly tracks the checksums of build
graphs and content mounted to specific operations. This makes it much faster,
more precise, and portable. The build cache can even be exported to a registry,
where it can be pulled on-demand by subsequent invocations on any host.

LLB can be generated directly using a
[golang client package](https://pkg.go.dev/github.com/moby/buildkit/client/llb) that allows defining the relationships between your
build operations using Go language primitives. This gives you full power to run
anything you can imagine, but will probably not be how most people will define
their builds. Instead, most users would use a frontend component, or LLB nested
invocation, to run a prepared set of build steps.

## Frontend

A frontend is a component that takes a human-readable build format and converts
it to LLB so BuildKit can execute it. Frontends can be distributed as images,
and the user can target a specific version of a frontend that is guaranteed to
work for the features used by their definition.

For example, to build a [Dockerfile](../../reference/dockerfile.md) with
BuildKit, you would
[use an external Dockerfile frontend](../dockerfile/frontend.md).

## Getting started

BuildKit is the default builder for users on Docker Desktop and Docker Engine
v23.0 and later.

If you have installed Docker Desktop, you don't need to enable BuildKit. If you
are running a version of Docker Engine version earlier than 23.0, you can enable
BuildKit either by setting an environment variable, or by making BuildKit the
default setting in the daemon configuration.

To set the BuildKit environment variable when running the `docker build`
command, run:

```console
$ DOCKER_BUILDKIT=1 docker build .
```

> **Note**
>
> [Buildx](../architecture.md#buildx) always uses BuildKit.

To use Docker BuildKit by default, edit the Docker daemon configuration in
`/etc/docker/daemon.json` as follows, and restart the daemon.

```json
{
  "features": {
    "buildkit": true
  }
}
```

If the `/etc/docker/daemon.json` file doesn't exist, create new file called
`daemon.json` and then add the following to the file. And restart the Docker
daemon.

> **Warning**
>
> BuildKit only fully supports building Linux containers. Experimental Windows
> container support is tracked in
> [`moby/buildkit#616`](https://github.com/moby/buildkit/issues/616)
{ .warning }
