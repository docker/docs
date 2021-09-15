---
description: Compose V1 to V2
keywords: compose
title: Compose V1 to V2
---

## Context of Docker Compose evolution

Introduction of the [Compose specification](https://github.com/compose-spec/compose-spec){:target="_blank" rel="noopener" class="_"} makes a clean distinction between the Compose YAML file model and the `docker-compose` implementation. Making this change has enabled a number of enhancements, including adding the `compose` command directly into the Docker CLI,  being able to “up” a Compose application on cloud platforms by simply switching the Docker context, and launching of [Amazon ECS](/cloud/ecs-integration) and [Microsoft ACI](/cloud/aci-integration). As the Compose specification evolves, new features land faster in the Docker CLI.

While docker-compose is still supported and maintained, Compose V2 implementation relies directly on the compose-go bindings which are maintained as part of the specification. This allows us to include community proposals, experimental implementations by the Docker CLI and/or Engine, and deliver features faster to users. Compose V2 also  supports some of the newer additions to the Compose specification, such as [profiles](profiles.md) and [GPU](gpu-support.md) devices.

Additionally, Compose V2 also supports [Apple silicon](../desktop/mac/apple-silicon.md).

For more information about the flags that are supported in the new compose command, see the [docker-compose compatibility list](cli-command-compatibility.md).

## Compose V2 and the new `docker compose` command

Starting with Docker Desktop 3.4.0, you can run Compose V2 commands without modifying your invocations, by enabling the drop-in replacement of the previous `docker-compose` with the new command.  See the section [Installing Compose v2](#installing-compose-v2) for detailed instructions how to enable the drop-in replacement.

We turn this option on automatically for Docker Desktop users, so that users can seamlessly move to Docker Compose V2 without the need to upgrade any of their scripts. If you run into any problems with Compose V2, you can easily switch back to Compose V1 by either by making changes in Docker Desktop Settings, or by running the command `docker-compose disable-v2`.

Your feedback is important to us. Let us know your feedback on the new 'compose' command by creating an issue in the [Compose](https://github.com/docker/compose/issues){:target="_blank" rel="noopener" class="_"} GitHub repository.
{: .important}

## Compose V1 end of life

### September 28th, 2021 
Compose V2 is the default development branch on GitHub.
New features and bug fixes will be considered in the V2 codebase. Legacy V1 codebase will only be considered
for security issues and critical bug fixes.
Users on Mac/Windows will be defaulted into Docker Compose V2, but can still opt out through the UI and the cli. 

### March 28th, 2022: 
After a 6 months transition period, V1 is marked as deprecated
No other impacts to users than a strong signal to users

### End of Life
Actual date to be defined based on community feedback and Compose V2 adoption 
End of security fixes
No new contribution will be accepted to the V1 branch, even for security fixes.
