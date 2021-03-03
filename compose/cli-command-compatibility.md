---
description: Docker Compose command compatibility with docker-compose
keywords: documentation, docs, docker, compose, containers
title: Docker Compose command compatibility with docker-compose
---

In the current state, the compose command in the Docker CLI supports most of the docker-compose commands and flags. It is expected to be a drop in replacement for docker-compose. However, there are a few remaining flags that have not been implemented yet. We are prioritizing the implementation of those based on usage metrics and user feedback.

You can follow progress on the implementation of the remaining commands and flags [here](https://github.com/docker/compose-cli/issues/1283){:target="_blank" rel="noopener" class="_"}.

If you see some compose functionality that is not available in the compose command, please [let us know](https://github.com/docker/compose-cli/issues){:target="_blank" rel="noopener" class="_"} so we can prioritize it.

## Commands or flags not yet implemented

The following commands have not been implemented yet, and maybe implemented at a later time.
Please let us know if these commands are a higher priority for your usecases.

* `compose build --memory`
* `compose build --no-cache`
* `compose config --no-interpolate`
* `compose config --services`
* `compose config --volumes`
* `compose config --hash`
* `compose events`
* `compose images`
* `compose port`
* `compose pull --ignore-pull-failures`
* `compose push --ignore-push-failures`
* `compose run --service-ports`

## Flags that will not be implemented

The list below includes the flags that we are not planning to support in Compose in Docker CLI,
either because they are already deprecated in docker-compose, or because they are not relevant for the Compose in Docker CLI.

* `compose build --compress` Not relevant as commpose command is using buildkit by default.
* `compose build --force-rm` Not relevant as commpose command is using buildkit by default.
* `compose build --no-rm` Not relevant as commpose command is using buildkit by default.
* `compose build --parallel` Not relevant as commpose command is using buildkit by default.
* `compose ps --filter KEY-VALUE` Not relevant as cumbersome usage with the service command and not documented properly in docker-compose.
* `compose pull --parallel` Deprecated in docker-compose
* `compose pull --no-parallel` Deprecated in docker-compose
* `compose rm --all` Deprecated in docker-compose.
* `compose scale` Deprecated in docker-compose (use `compose up --scale` instead)

Global flags:

* `compose --no-ansi` Deprecated in docker-compose.
* `compose --compatibility` Deprecated in docker-compose.
