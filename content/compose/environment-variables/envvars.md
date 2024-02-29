---
description: Compose pre-defined environment variables
keywords: fig, composition, compose, docker, orchestration, cli, reference
title: Set or change pre-defined environment variables in Docker Compose
aliases:
- /compose/reference/envvars/
---

Compose already comes with pre-defined environment variables. It also inherits common Docker CLI environment variables, such as `DOCKER_HOST` and `DOCKER_CONTEXT`. See [Docker CLI environment variable reference](/engine/reference/commandline/cli/#environment-variables) for details.

This page contains information on how you can set or change the following pre-defined environment variables if you need to:

- `COMPOSE_CONVERT_WINDOWS_PATHS`
- `COMPOSE_FILE`
- `COMPOSE_PROFILES`
- `COMPOSE_PROJECT_NAME`
- `DOCKER_CERT_PATH`
- `COMPOSE_PARALLEL_LIMIT`
- `COMPOSE_IGNORE_ORPHANS`
- `COMPOSE_REMOVE_ORPHANS`
- `COMPOSE_PATH_SEPARATOR`
- `COMPOSE_ANSI`
- `COMPOSE_STATUS_STDOUT`
- `COMPOSE_ENV_FILES`

## Methods to override 

You can set or change the pre-defined environment variables:
- Within your Compose file using the [`environment` attribute](set-environment-variables.md#use-the-environment-attribute)
- With an [environment file](env-file.md) 
- From the command line
- From your [shell](set-variables.md#substitute-from-the-shell)

When changing or setting any environment variables, be aware of [Environment variable precedence](envvars-precedence.md).

## Configure

### COMPOSE\_PROJECT\_NAME

Sets the project name. This value is prepended along with the service name to
the container's name on startup.

For example, if your project name is `myapp` and it includes two services `db` and `web`, 
then Compose starts containers named `myapp-db-1` and `myapp-web-1` respectively.

Compose can set the project name in different ways. The level of precedence (from highest to lowest) for each method is as follows:

1. The `-p` command line flag 
2. `COMPOSE_PROJECT_NAME`
3. The top level `name:` variable from the config file (or the last `name:` from
  a series of config files specified using `-f`)
4. The `basename` of the project directory containing the config file (or
  containing the first config file specified using `-f`)
5. The `basename` of the current directory if no config file is specified

Project names must contain only lowercase letters, decimal digits, dashes, and
underscores, and must begin with a lowercase letter or decimal digit. If the
`basename` of the project directory or current directory violates this
constraint, you must use one of the other mechanisms.

See also the [command-line options overview](../reference/index.md#command-options-overview-and-help) and [using `-p` to specify a project name](../reference/index.md#use--p-to-specify-a-project-name).

### COMPOSE\_FILE

Specifies the path to a Compose file. Specifying multiple Compose files is supported.

- Default behavior: If not provided, Compose looks for a file named `compose.yaml` or `docker-compose.yaml` in the current directory and, if not found, then Compose searches each parent directory recursively until a file by that name is found.
- Default separator: When specifying multiple Compose files, the path separators are, by default, on:
    * Mac and Linux: `:` (colon),
    * Windows: `;` (semicolon).

The path separator can also be customized using `COMPOSE_PATH_SEPARATOR`.  

Example: `COMPOSE_FILE=docker-compose.yml:docker-compose.prod.yml`.  

See also the [command-line options overview](../reference/index.md#command-options-overview-and-help) and [using `-f` to specify name and path of one or more Compose files](../reference/index.md#use--f-to-specify-name-and-path-of-one-or-more-compose-files).

### COMPOSE\_PROFILES

Specifies one or more profiles to be enabled on `compose up` execution.
Services with matching profiles are started as well as any services for which no profile has been defined.

For example, calling `docker compose up`with `COMPOSE_PROFILES=frontend` selects services with the 
`frontend` profile as well as any services without a profile specified.

* Default separator: specify a list of profiles using a comma as separator.

Example: `COMPOSE_PROFILES=frontend,debug`  
This example enables all services matching both the `frontend` and `debug` profiles and services without a profile.

See also [Using profiles with Compose](../profiles.md) and the [`--profile` command-line option](../reference/index.md#use---profile-to-specify-one-or-more-active-profiles).

### COMPOSE\_CONVERT\_WINDOWS\_PATHS

When enabled, Compose performs path conversion from Windows-style to Unix-style in volume definitions.

* Supported values: 
    * `true` or `1`, to enable,
    * `false` or `0`, to disable.
* Defaults to: `0`.

### COMPOSE\_PATH\_SEPARATOR

Specifies a different path separator for items listed in `COMPOSE_FILE`.

* Defaults to:
    * On macOS and Linux to `:`,
    * On Windows to`;`.

### COMPOSE\_IGNORE\_ORPHANS

When enabled, Compose doesn't try to detect orphaned containers for the project.

* Supported values: 
    * `true` or `1`, to enable,
    * `false` or `0`, to disable.
* Defaults to: `0`.

### COMPOSE\_PARALLEL\_LIMIT

Specifies the maximum level of parallelism for concurrent engine calls.

### COMPOSE\_ANSI

Specifies when to print ANSI control characters. 

* Supported values:
  * `auto`, Compose detects if TTY mode can be used. Otherwise, use plain text mode.
  * `never`, use plain text mode.
  * `always` or `0`, use TTY mode.
* Defaults to: `auto`.

### COMPOSE\_STATUS\_STDOUT

When enabled, Compose writes its internal status and progress messages to `stdout` instead of `stderr`. 
The default value is false to clearly separate the output streams between Compose messages and your container's logs.

* Supported values:
  * `true` or `1`, to enable,
  * `false` or `0`, to disable.
* Defaults to: `0`.

### COMPOSE\_ENV\_FILES

Lets you specify which environment files Compose should use if `--env-file` isn't used.

When using multiple environment files, use a comma as a separator. For example, 

```console
COMPOSE_ENV_FILES=.env.envfile1, .env.envfile2
```

If `COMPOSE_ENV_FILES` is not set, and you don't provide `--env-file` in the CLI, Docker Compose uses the default behavior, which is to look for an `.env` file in the project directory.

## Unsupported in Compose V2

The environment variables listed below have no effect in Compose V2.
For more information, see [Migrate to Compose V2](../migrate.md).

- `COMPOSE_API_VERSION`
    By default the API version is negotiated with the server. Use `DOCKER_API_VERSION`.  
    See the [Docker CLI environment variable reference](../../../engine/reference/commandline/cli/#environment-variables) page.
- `COMPOSE_HTTP_TIMEOUT`
- `COMPOSE_TLS_VERSION`
- `COMPOSE_FORCE_WINDOWS_HOST`
- `COMPOSE_INTERACTIVE_NO_CLI`
- `COMPOSE_DOCKER_CLI_BUILD`
    Use `DOCKER_BUILDKIT` to select between BuildKit and the classic builder. If `DOCKER_BUILDKIT=0` then `docker compose build` uses the classic builder to build images.
