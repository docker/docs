---
description: Restart a machine
keywords: machine, restart, subcommand
title: docker-machine restart
hide_from_sitemap: true
---

```none
Usage: docker-machine restart [arg...]

Restart a machine

Description:
   Argument(s) are one or more machine names.
```

Restart a machine. Oftentimes this is equivalent to
`docker-machine stop; docker-machine start`. But some cloud driver try to implement a clever restart which keeps the same
IP address.

```bash
$ docker-machine restart dev

Waiting for VM to start...
```
