---
description: Overview of the Docker Compose CLI
keywords: fig, composition, compose, docker, orchestration, cli, reference, docker-compose
title: Overview of docker compose CLI
aliases:
- /compose/reference/overview/
---

{{< include "compose-eol.md" >}}

This page provides usage information for the `docker compose` command.

## Command options overview and help

You can also see this information by running `docker compose --help` from the
command line.

```none
Usage:  docker compose [OPTIONS] COMMAND

Define and run multi-container applications with Docker.

Options:
      --ansi string                Control when to print ANSI control characters ("never"|"always"|"auto") (default "auto")
      --compatibility              Run compose in backward compatibility mode
      --env-file stringArray       Specify an alternate environment file
  -f, --file stringArray           Compose configuration files
      --parallel int               Control max parallelism, -1 for unlimited (default -1)
      --profile stringArray        Specify a profile to enable
      --project-directory string   Specify an alternate working directory
                                   (default: the path of the, first specified, Compose file)
  -p, --project-name string        Project name

Commands:
  build       Build or rebuild services
  config      Parse, resolve and render compose file in canonical format
  cp          Copy files/folders between a service container and the local filesystem
  create      Creates containers for a service.
  down        Stop and remove containers, networks
  events      Receive real time events from containers.
  exec        Execute a command in a running container.
  images      List images used by the created containers
  kill        Force stop service containers.
  logs        View output from containers
  ls          List running compose projects
  pause       Pause services
  port        Print the public port for a port binding.
  ps          List containers
  pull        Pull service images
  push        Push service images
  restart     Restart service containers
  rm          Removes stopped service containers
  run         Run a one-off command on a service.
  start       Start services
  stop        Stop services
  top         Display the running processes
  unpause     Unpause services
  up          Create and start containers
  version     Show the Docker Compose version information

Run 'docker compose COMMAND --help' for more information on a command.
```

You can use Docker Compose binary, `docker compose [-f <arg>...] [options]
[COMMAND] [ARGS...]`, to build and manage multiple services in Docker containers.




## Set up environment variables

You can set [environment variables](../environment-variables/envvars.md) for various
`docker compose` options, including the `-f` and `-p` flags.

For example, the [COMPOSE_FILE environment variable](../environment-variables/envvars.md#compose_file)
relates to the `-f` flag, and `COMPOSE_PROJECT_NAME`
[environment variable](../environment-variables/envvars.md#compose_project_name) relates to the `-p` flag.

Also, you can set some of these variables in an [environment file](../environment-variables/env-file.md).

## Where to go next

* [CLI environment variables](../environment-variables/envvars.md)
* [Declare default environment variables in file](../environment-variables/env-file.md)