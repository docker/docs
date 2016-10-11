---
description: Lists containers.
keywords:
- fig, composition, compose, docker, orchestration, cli,  ps
menu:
  main:
    identifier: ps.compose
    parent: smn_compose_cli
title: ps
---

# ps

```
Usage: ps [options] [SERVICE...]

Options:
-q    Only display IDs
```

Lists containers.

```bash
$ docker-compose ps
         Name                        Command                 State             Ports         
--------------------------------------------------------------------------------------------
mywordpress_db_1          docker-entrypoint.sh mysqld      Up           3306/tcp             
mywordpress_wordpress_1   /entrypoint.sh apache2-for ...   Restarting   0.0.0.0:8000->80/tcp
```
