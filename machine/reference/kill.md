---
description: Kill (abruptly force stop) a machine.
keywords: machine, kill, subcommand
title: docker-machine kill
hide_from_sitemap: true
---

```none
Usage: docker-machine kill [arg...]

Kill (abruptly force stop) a machine

Description:
   Argument(s) are one or more machine names.
```

For example:

```bash
$ docker-machine ls

NAME   ACTIVE   DRIVER       STATE     URL
dev    *        virtualbox   Running   tcp://192.168.99.104:2376

$ docker-machine kill dev
$ docker-machine ls

NAME   ACTIVE   DRIVER       STATE     URL
dev    *        virtualbox   Stopped
```
