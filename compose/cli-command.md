---
description: Docker Compose
keywords: compose, V2
title: Compose V2
---

## Compose V2 and the new `docker compose` command

> Important
>
> The new Compose V2, which supports the `compose` command as part of the Docker CLI, is now available.
>
> Compose V2 integrates compose functions into the Docker platform, continuing to support most of the previous `docker-compose` features and flags. You can test the Compose V2 by simply replacing the dash (`-`) with a space, and by running `docker compose`, instead of `docker-compose`.
{: .important}

Starting with Docker Desktop 3.4.0, you can run Compose V2 commands without modifying your invocations, by enabling the drop-in replacement of the previous `docker-compose` with the new command.  See the section [Installing Compose](install.md) for detailed instructions.

## Context of Docker Compose evolution

Introduction of the [Compose specification](https://github.com/compose-spec/compose-spec){:target="_blank" rel="noopener" class="_"} makes a clean distinction between the Compose YAML file model and the `docker-compose` implementation. Making this change has enabled a number of enhancements, including adding the `compose` command directly into the Docker CLI,  being able to “up” a Compose application on cloud platforms by simply switching the Docker context, and launching of [Amazon ECS](/cloud/ecs-integration) and [Microsoft ACI](/cloud/aci-integration). As the Compose specification evolves, new features land faster in the Docker CLI.

Compose V2 implementation relies directly on the compose-go bindings which are maintained as part of the specification. This allows us to include community proposals, experimental implementations by the Docker CLI and/or Engine, and deliver features faster to users. Compose V2 also  supports some of the newer additions to the Compose specification, such as [profiles](profiles.md) and [GPU](gpu-support.md) devices.

Additionally, Compose V2 also supports [Apple silicon](../desktop/mac/apple-silicon.md).

For more information about the flags that are supported in the new compose command, see the [docker-compose compatibility list](cli-command-compatibility.md).

## Where to go next

- [User guide](index.md)
- [Installing Compose](install.md)