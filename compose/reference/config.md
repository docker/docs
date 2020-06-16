---
description: Config validates and view the compose file.
keywords: fig, composition, compose, docker, orchestration, cli, config
title: docker-compose config
notoc: true
redirect_from:
- /compose/reference/bundle/
---

```
Usage: config [options]

Options:
    --resolve-image-digests  Pin image tags to digests.
    --no-interpolate         Don't interpolate environment variables.
    -q, --quiet              Only validate the configuration, don't print anything.
    --services               Print the service names, one per line.
    --volumes                Print the volume names, one per line.
    --hash="*"               Print the service config hash, one per line.
                             Set "service1,service2" for a list of specified services
                             or use the wildcard symbol to display all services.
```

Validate and view the Compose file.
