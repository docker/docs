---
description: Create creates containers for a service.
keywords: fig, composition, compose, docker, orchestration, cli, create
title: docker-compose create
notoc: true
---

> **This command is deprecated.** Use the [up](up.md) command with `--no-start`
instead.
{: .warning }

```none
Usage: docker-compose create [options] [SERVICE...]

Options:
    --force-recreate       Recreate containers even if their configuration and
                           image haven't changed. Incompatible with --no-recreate.
    --no-recreate          If containers already exist, don't recreate them.
                           Incompatible with --force-recreate.
    --no-build             Don't build an image, even if it's missing.
    --build                Build images before creating containers.
```
