---
description: Displays the running processes.
keywords: fig, composition, compose, docker, orchestration, cli, top
title: docker-compose top
notoc: true
---

```none
Usage: top [SERVICE...]

```

Displays the running processes.

```bash
$ docker-compose top
compose_service_a_1
PID    USER   TIME   COMMAND
----------------------------
4060   root   0:00   top

compose_service_b_1
PID    USER   TIME   COMMAND
----------------------------
4115   root   0:00   top
```
