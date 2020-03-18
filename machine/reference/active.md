---
description: Identify active machines
keywords: machine, active, subcommand
title: docker-machine active
hide_from_sitemap: true
---

See which machine is "active" (a machine is considered active if the
`DOCKER_HOST` environment variable points to it).

```bash
$ docker-machine ls

NAME      ACTIVE   DRIVER         STATE     URL
dev       -        virtualbox     Running   tcp://192.168.99.103:2376
staging   *        digitalocean   Running   tcp://203.0.113.81:2376

$ echo $DOCKER_HOST
tcp://203.0.113.81:2376

$ docker-machine active
staging
```
