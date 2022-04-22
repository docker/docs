---
description: Lists containers.
keywords: fig, composition, compose, docker, orchestration, cli,  ps
title: docker-compose ps
notoc: true
---

```none
Usage: docker-compose ps [options] [SERVICE...]

Options:
    -q, --quiet          Only display IDs
    --services           Display services
    --filter KEY=VAL     Filter services by a property
    -a, --all            Show all stopped containers (including those created by the run command)
```

Lists containers.

```console
$ docker-compose ps
         Name                        Command                 State             Ports
---------------------------------------------------------------------------------------------
mywordpress_db_1          docker-entrypoint.sh mysqld      Up (healthy)  3306/tcp
mywordpress_wordpress_1   /entrypoint.sh apache2-for ...   Restarting    0.0.0.0:8000->80/tcp
```

List containers with filter option.

Supported filter: status. Values supported for status: [paused | restarting | removing | running | dead | created | exited]

```console
$ docker-compose ps --filter status=running
       Name                      Command               State           Ports
-------------------------------------------------------------------------------------
composetest_redis_1   docker-entrypoint.sh redis ...   Up      6379/tcp
composetest_web_1     flask run                        Up      0.0.0.0:8000->5000/tcp
```

List container with services and filter option

```console
$ docker-compose ps --filter status=running redis
       Name                      Command               State    Ports
-----------------------------------------------------------------------
composetest_redis_1   docker-entrypoint.sh redis ...   Up      6379/tcp
```
