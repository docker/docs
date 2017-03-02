---
description: Removes stopped service containers.
keywords: fig, composition, compose, docker, orchestration, cli,  rm
title: docker-compose rm
notoc: true
---

```
Usage: rm [options] [SERVICE...]

Options:
    -f, --force   Don't ask to confirm removal
    -v            Remove any anonymous volumes attached to containers
    -a, --all     Also remove one-off containers created by
                  docker-compose run
```

Removes stopped service containers.

By default, anonymous volumes attached to containers will not be removed. You
can override this with `-v`. To list all volumes, use `docker volume ls`.

Any data which is not in a volume will be lost.
