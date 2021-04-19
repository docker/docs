---
description: Compose command compatibility with docker-compose
keywords: documentation, docs, docker, compose, containers
title: Compose command compatibility with docker-compose
---

The `compose` command in the Docker CLI supports most of the `docker-compose` commands and flags. It is expected to be a drop-in replacement for `docker-compose`. There are a few remaining flags that have yet to be implemented, and we are prioritizing these implementations based on usage metrics and user feedback.

You can follow progress on the implementation of the remaining commands and flags in the  [Compose-CLI](https://github.com/docker/compose-cli/issues/1283){:target="_blank" rel="noopener" class="_"} GitHub repository.

If you see some Compose functionality that is not available in the `compose` command, create an issue in the [Compose-CLI](https://github.com/docker/compose-cli/issues){:target="_blank" rel="noopener" class="_"} GitHub repository so we can prioritize it.

## Commands or flags not yet implemented

The following commands have not been implemented yet, and maybe implemented at a later time.
Let us know if these commands are a higher priority for your usecases.

* `compose build --memory`

## Flags that will not be implemented

The list below includes the flags that we are not planning to support in Compose in the Docker CLI,
either because they are already deprecated in `docker-compose`, or because they are not relevant for Compose in the Docker CLI.

* `compose build --compress` Not relevant as the 'compose' command uses buildkit by default.
* `compose build --force-rm` Not relevant as commpose command is using buildkit by default.
* `compose build --no-rm` Not relevant as commpose command is using buildkit by default.
* `compose build --parallel` Not relevant as commpose command is using buildkit by default.
* `compose ps --filter KEY-VALUE` Not relevant due to its complicated usage with the `service` command and also because it is not documented properly in `docker-compose`.
* `compose pull --parallel` Deprecated in docker-compose (Still parsed, but ignored)
* `compose pull --no-parallel` Deprecated in docker-compose (Still parsed, but ignored)
* `compose rm --all` Deprecated in docker-compose.
* `compose scale` Deprecated in docker-compose (use `compose up --scale` instead)

Global flags:

* `compose --no-ansi` Deprecated in docker-compose.
* `compose --compatibility` Deprecated in docker-compose.
