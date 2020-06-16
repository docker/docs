---
description: Restarts Docker Compose services.
keywords: fig, composition, compose, docker, orchestration, cli,  restart
title: docker-compose restart
notoc: true
---

```
Usage: restart [options] [SERVICE...]

Options:
  -t, --timeout TIMEOUT      Specify a shutdown timeout in seconds.
                             (default: 10)
```

Restarts all stopped and running services.

If you make changes to your `docker-compose.yml` configuration these changes are not reflected after running this command.

For example, changes to environment variables (which are added after a container is built, but before the container's command is executed) are not updated after restarting.

If you are looking to configure a service's restart policy, please refer to
[restart](../compose-file/index.md#restart) in Compose file v3 and
[restart](../compose-file/compose-file-v2.md#restart) in Compose v2. Note that if
you are [deploying a stack in swarm mode](../../engine/reference/commandline/stack_deploy.md),
you should use [restart_policy](../compose-file/index.md#restart), instead.
