---
description: Docker Compose in the Docker CLI
keywords: documentation, docs, docker, compose, containers
title: Docker Compose in the Docker CLI
---

The Docker CLI now supports the compose command to include (most of) the docker-compose features and flags in the Docker CLI, without the need for a separate tool.

You can just replace the dash with a space when you use docker-compose to just adopt “docker compose”. You can also use them interchangeably, so that you are not locked-in with the new compose command and can still use docker-compose if needed.

With the introduction of the [Compose specification](https://github.com/compose-spec/compose-spec){:target="_blank" rel="noopener" class="_"}, a clean distinction has been made between the Compose YAML file model and the docker-compose implementation. This made it possible to introduce [Amazon ECS](/cloud/ecs-integration) and [Microsoft ACI](/cloud/aci-integration) support in the Docker CLI so one can “up” a Compose application on cloud platforms just by switching Docker context. With this solid backbone in place, we have added compose in the Docker CLI as a first class command.
As the Compose specification evolves, new features will land in the Docker CLI faster. While docker-compose is still alive and will be maintained, Compose in the Docker CLI Go implementation relies directly on the compose-go bindings which  are maintained as part of the specification. This will allow it to more quickly include community proposals, experimental implementations by the Docker CLI and/or Engine, and deliver features to end-users. Compose in the Docker CLI already supports some of the newer additions to the Compose specification: profiles and GPU devices.
