---
description: Declare default environment variables in a file
keywords: fig, composition, compose, docker, orchestration, environment, env file
title: Declare default environment variables in file
---

Compose supports declaring default environment variables in an environment file
named `.env` placed in the folder where the `docker-compose` command is executed
*(current working directory)*.

## Syntax rules

These syntax rules apply to the `.env` file:

* Compose expects each line in an `env` file to be in `VAR=VAL` format.
* Lines beginning with `#` are processed as comments and ignored.
* Blank lines are ignored.
* There is no special handling of quotation marks. This means that
  **they are part of the VAL**.

## Compose file and CLI variables

The environment variables you define here are used for [variable
substitution](compose-file/index.md#variable-substitution) in your Compose file,
and can also be used to define the following [CLI
variables](reference/envvars.md):

- `COMPOSE_API_VERSION`
- `COMPOSE_CONVERT_WINDOWS_PATHS`
- `COMPOSE_FILE`
- `COMPOSE_HTTP_TIMEOUT`
- `COMPOSE_TLS_VERSION`
- `COMPOSE_PROJECT_NAME`
- `DOCKER_CERT_PATH`
- `DOCKER_HOST`
- `DOCKER_TLS_VERIFY`

> **Notes**
>
> * Values present in the environment at runtime always override
those defined inside the `.env` file. Similarly, values passed via command-line
arguments take precedence as well.
>
> * Environment variables defined in the `.env` file are not
automatically visible inside containers. To set container-applicable
environment variables, follow the guidelines in the topic [Environment variables
in Compose](/compose/environment-variables.md), which describes how to pass
shell environment variables through to containers, define environment variables
in Compose files, and more.

## More Compose documentation

- [User guide](index.md)
- [Command line reference](./reference/index.md)
- [Compose file reference](compose-file.md)
