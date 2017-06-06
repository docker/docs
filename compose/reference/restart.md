---
description: Restarts Docker Compose services.
keywords: fig, composition, compose, docker, orchestration, cli,  restart
title: docker-compose restart
notoc: true
---

```
Usage: restart [options] [SERVICE...]

Options:
-t, --timeout TIMEOUT      Specify a shutdown timeout in seconds. (default: 10)
```

Restarts all stopped and running services.

If you make changes to your `docker-compose.yml` configuration these changes will not be reflected after running this command.

For example, changes to environment variables (which are added after a container is built, but before the container's command is executed) will not be updated after restarting.
