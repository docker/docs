---
description: Sets the number of containers to run for a service.
keywords: fig, composition, compose, docker, orchestration, cli,  scale
title: docker-compose scale
notoc: true
---

> **Note**: This command is deprecated. Use the [up](up.md) command with the
  `--scale` flag instead. Beware that using `up` with `--scale` flag has some
  [subtle differences](https://github.com/docker/compose/issues/5251) with the `scale` command as it incorporates the behaviour
  of `up` command.

```
Usage: scale [SERVICE=NUM...]
```

Sets the number of containers to run for a service.

Numbers are specified as arguments in the form `service=num`. For example:

    docker-compose scale web=2 worker=3

>**Tip**: Alternatively, in
[Compose file version 3.x](/compose/compose-file/index.md), you can specify
[replicas](/compose/compose-file/index.md#replicas)
under the [deploy](/compose/compose-file/index.md#deploy) key as part of a
service configuration for [Swarm mode](/engine/swarm/). The `deploy` key and its sub-options (including `replicas`) only works with the `docker stack deploy` command, not `docker-compose up` or `docker-compose run`.
