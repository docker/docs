---
description: Key features and use cases of Docker Compose
keywords: documentation, docs, docker, compose, orchestration, containers, uses, features
title: Compose V2 Overview
---

## Compose V2 and the new `docker compose` command

> Important
>
> The new Compose V2, which supports the `compose` command as part of the Docker
> CLI, is now available.
>
> Compose V2 integrates compose functions into the Docker platform, continuing
> to support most of the previous `docker-compose` features and flags. You can
> run Compose V2 by replacing the hyphen (`-`) with a space, using `docker compose`,
> instead of `docker-compose`.
{: .important}

If you rely on using Docker Compose as `docker-compose` (with a hyphen), you can
set up Compose V2 to act as a drop-in replacement of the previous `docker-compose`.
Refer to the [Installing Compose](../install/index.md) section for detailed instructions.

## Context of Docker Compose evolution

Introduction of the [Compose specification](https://github.com/compose-spec/compose-spec){:target="_blank" rel="noopener" class="_"}
makes a clean distinction between the Compose YAML file model and the `docker-compose`
implementation. Making this change has enabled a number of enhancements, including
adding the `compose` command directly into the Docker CLI,  being able to "up" a
Compose application on cloud platforms by simply switching the Docker context,
and launching of [Amazon ECS](../../cloud/ecs-integration.md) and [Microsoft ACI](../../cloud/aci-integration.md).
As the Compose specification evolves, new features land faster in the Docker CLI.

Compose V2 relies directly on the compose-go bindings which are maintained as part
of the specification. This allows us to include community proposals, experimental
implementations by the Docker CLI and/or Engine, and deliver features faster to
users. Compose V2 also supports some of the newer additions to the specification,
such as [profiles](../profiles.md) and [GPU](../gpu-support.md) devices.

Compose V2 has been re-written in [Go](https://go.dev), which improves integration
with other Docker command-line features, and allows it to run natively on 
[macOS on Apple silicon](../../desktop/mac/apple-silicon.md), Windows, and Linux,
without dependencies such as Python. 

For more information about compatibility with the compose v1 command-line, see the [docker-compose compatibility list](../cli-command-compatibility.md).