---
description: Sets the number of containers to run for a service.
keywords: fig, composition, compose, docker, orchestration, cli,  scale
title: docker-compose scale
notoc: true
---

> **This command is deprecated.** Use the [up](up.md) command with the
  `--scale` flag instead. Beware that using `up` with the `--scale` flag has
  some [subtle differences](https://github.com/docker/compose/issues/5251) with
  the `scale` command, as it incorporates the behaviour of the `up` command.
  {: .warning }

```
Usage: scale [options] [SERVICE=NUM...]

Options:
  -t, --timeout TIMEOUT      Specify a shutdown timeout in seconds.
                             (default: 10)
```

Sets the number of containers to run for a service.

Numbers are specified as arguments in the form `service=num`. For example:

    docker-compose scale web=2 worker=3

>**Tip**: Alternatively, in
[Compose file version 3.x](../compose-file/index.md), you can specify
[replicas](../compose-file/index.md#replicas)
under the [deploy](../compose-file/index.md#deploy) key as part of a
service configuration for [Swarm mode](/engine/swarm/). The `deploy` key and its sub-options (including `replicas`) only works with the `docker stack deploy` command, not `docker-compose up` or `docker-compose run`.
