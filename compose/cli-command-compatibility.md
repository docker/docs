---
description: Compose command compatibility with docker-compose
keywords: documentation, docs, docker, compose, containers
title: Compose command compatibility with docker-compose
---

The `compose` command in the Docker CLI supports most of the `docker-compose` commands and flags. It is expected to be a drop-in replacement for `docker-compose`. 

If you see any Compose functionality that is not available in the `compose` command, create an issue in the [Compose](https://github.com/docker/compose/issues){:target="_blank" rel="noopener" class="_"} GitHub repository so we can prioritize it.

## Commands or flags not yet implemented

The following commands have not been implemented yet, and maybe implemented at a later time.
Let us know if these commands are a higher priority for your use cases.

`compose build --memory`: This option is not yet supported by buildkit. The flag is currently supported, but is hidden to avoid breaking existing Compose usage. It does not have any effect.

## Flags that will not be implemented

The list below includes the flags that we are not planning to support in Compose in the Docker CLI,
either because they are already deprecated in `docker-compose`, or because they are not relevant for Compose in the Docker CLI.

* The usage of `compose ps --filter KEY=VALUE` is discouraged in v2. Although `--filter` option remains available, the functionality behind this option is only partly ported. In v2, the only supported filter `KEY` is `status`. However, we recommend that you migrate to a new option `--status`, which mirrors `--filter status=xx` functionality in v2. Support for `KEY` `source` (`--filter source=[image|build]`) has been dropped, due to its complicated usage with the `service` command and also because it is not documented properly in `docker-compose`. There is no alternative in `compose ps` v2 for filtering by `source`.
* `compose rm --all` Deprecated in docker-compose.
* `compose scale` Deprecated in docker-compose (use `compose up --scale` instead)

Global flags:

* `compose --compatibility` Deprecated in docker-compose.
