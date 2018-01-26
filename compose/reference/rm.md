---
description: Removes stopped service containers.
keywords: fig, composition, compose, docker, orchestration, cli,  rm
title: docker-compose rm
notoc: true
---

```none
Usage: rm [options] [SERVICE...]

Options:
    -f, --force   Don't ask to confirm removal
    -s, --stop    Stop the containers, if required, before removing
    -v            Remove any anonymous volumes attached to containers
```

Removes stopped service containers.

By default, anonymous volumes attached to containers are not removed. You
can override this with `-v`. To list all volumes,  use `docker volume ls`.

Any data which is not in a volume is lost.

Running the command with no options also removes one-off containers created
by `docker-compose up` or `docker-compose run`:

```none
$ docker-compose rm
Going to remove djangoquickstart_web_run_1
Are you sure? [yN] y
Removing djangoquickstart_web_run_1 ... done
```