---
title: Ways to set variables for Compose file
description: How to set, use, and manage variables in Compose files
keywords: compose, orchestration, interpolation, env file
---

Compose allows to us variables for flexibility as described in [interpolation](env-file.md#interpolation). Values are resolved based on local environment variables and env files (not to confuse with [_container_'s environment](set-environment-variables.md)) as illustrated in this example:

```yaml
services:
  web:
    image: nginx:${TAG}
```

Variable syntax is the same used in env file as documented [here](env-file.md#interpolation)

The order of precedence (highest to lowest) is as follows:
1. Variables from your local shell
2. Set using an `.env` file placed in your working directory
3. Set using an `.env` file placed at base of your project directory, if distinct from working directory

The two latter are only relevant when compose is not ran with `--env-file` flag, as this disables load of a default
`.env` file. 

Working directory (a.k.a `pwd`) and project directory are distinct when:
- `COMPOSE_FILE` variable is set, and points to a compose file in a ditinct folder
- `docker compose` is ran with `--file` flag pointing to a compose file in a ditinct folder
- Working directory does not contain a default `compose.yaml` file, and Compose finds one in a parent folder, which becomes the project directory.
