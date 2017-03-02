---
description: Sets the number of containers to run for a service.
keywords: fig, composition, compose, docker, orchestration, cli,  scale
title: docker-compose scale
notoc: true
---

```
Usage: scale [SERVICE=NUM...]
```

Sets the number of containers to run for a service.

Numbers are specified as arguments in the form `service=num`. For example:

    docker-compose scale web=2 worker=3

>**Tip:** Alternatively, in
[Compose file version 3.x](/compose/compose-file/index.md), you can specify
[`replicas`](/compose/compose-file/index.md#replicas)
under [`deploy`](/compose/compose-file/index.md#deploy) as part of the
service configuration for [Swarm mode](/engine/swarm/).
