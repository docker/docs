---
description: Compose pre-defined environment variables
keywords: fig, composition, compose, docker, orchestration, cli, reference, compose environment configuration, docker env variables
title: Configure pre-defined environment variables in Docker Compose
linkTitle: Pre-defined environment variables
weight: 30
aliases:
- /compose/reference/envvars/
- /compose/environment-variables/envvars/
---

Docker Compose includes several pre-defined environment variables. It also inherits common Docker CLI environment variables, such as `DOCKER_HOST` and `DOCKER_CONTEXT`. See [Docker CLI environment variable reference](/reference/cli/docker/#environment-variables) for details.

This page explains how to set or change the following pre-defined environment variables:

- `COMPOSE_PROJECT_NAME`
- `COMPOSE_FILE`
- `COMPOSE_PROFILES`
- `COMPOSE_CONVERT_WINDOWS_PATHS`
- `COMPOSE_PATH_SEPARATOR`
- `COMPOSE_IGNORE_ORPHANS`
- `COMPOSE_REMOVE_ORPHANS`
- `COMPOSE_PARALLEL_LIMIT`
- `COMPOSE_ANSI`
- `COMPOSE_STATUS_STDOUT`
- `COMPOSE_ENV_FILES`
- `COMPOSE_DISABLE_ENV_FILE`
- `COMPOSE_MENU`
- `COMPOSE_EXPERIMENTAL`
- `COMPOSE_PROGRESS`

## Methods to override 

| Method      | Description                                  |
| ----------- | -------------------------------------------- |
| [`.env` file](/manuals/compose/how-tos/environment-variables/variable-interpolation.md) | Located in the working directory.            |
| [Shell](variable-interpolation.md#substitute-from-the-shell)       | Defined in the host operating system shell.  |
| CLI         | Passed with `--env` or `-e` flag at runtime. |

When changing or setting any environment variables, be aware of [Environment variable precedence](envvars-precedence.md).

## Configuration details

### Project and file configuration

#### COMPOSE\_PROJECT\_NAME

Sets the project name. This value is prepended along with the service name to
the container's name on startup.

For example, if your project name is `myapp` and it includes two services `db` and `web`, 
then Compose starts containers named `myapp-db-1` and `myapp-web-1` respectively.

Compose can set the project name in different ways. The level of precedence (from highest to lowest) for each method is as follows:

1. The `-p` command line flag 
2. `COMPOSE_PROJECT_NAME`
3. The top-level `name:` variable from the config file (or the last `name:` from
  a series of config files specified using `-f`)
4. The `basename` of the project directory containing the config file (or
  containing the first config file specified using `-f`)
5. The `basename` of the current directory if no config file is specified

Project names must contain only lowercase letters, decimal digits, dashes, and
underscores, and must begin with a lowercase letter or decimal digit. If the
`basename` of the project directory or current directory violates this
constraint, you must use one of the other mechanisms.

See also the [command-line options overview](/reference/cli/docker/compose/_index.md#command-options-overview-and-help) and [using `-p` to specify a project name](/reference/cli/docker/compose/_index.md#use--p-to-specify-a-project-name).

#### COMPOSE\_FILE

Specifies the path to a Compose file. Specifying multiple Compose files is supported.

- Default behavior: If not provided, Compose looks for a file named `compose.yaml` in the current directory and, if not found, then Compose searches each parent directory recursively until a file by that name is found.
- When specifying multiple Compose files, the path separators are, by default, on:
   - Mac and Linux: `:` (colon)
   - Windows: `;` (semicolon)
   For example:

      ```console
      COMPOSE_FILE=compose.yaml:compose.prod.yaml
      ```  
   The path separator can also be customized using [`COMPOSE_PATH_SEPARATOR`](#compose_path_separator).  

See also the [command-line options overview](/reference/cli/docker/compose/_index.md#command-options-overview-and-help) and [using `-f` to specify name and path of one or more Compose files](/reference/cli/docker/compose/_index.md#use--f-to-specify-the-name-and-path-of-one-or-more-compose-files).

#### COMPOSE\_PROFILES

Specifies one or more profiles to be enabled when `docker compose up` is run.

Services with matching profiles are started as well as any services for which no profile has been defined.

For example, calling `docker compose up` with `COMPOSE_PROFILES=frontend` selects services with the 
`frontend` profile as well as any services without a profile specified.

If specifying multiple profiles, use a comma as a separator.

The following example enables all services matching both the `frontend` and `debug` profiles and services without a profile. 

```console
COMPOSE_PROFILES=frontend,debug
```

See also [Using profiles with Compose](../profiles.md) and the [`--profile` command-line option](/reference/cli/docker/compose/_index.md#use-profiles-to-enable-optional-services).

#### COMPOSE\_PATH\_SEPARATOR

Specifies a different path separator for items listed in `COMPOSE_FILE`.

- Defaults to:
    - On macOS and Linux to `:`
    - On Windows to`;`

#### COMPOSE\_ENV\_FILES

Specifies which environment files Compose should use if `--env-file` isn't used.

When using multiple environment files, use a comma as a separator. For example: 

```console
COMPOSE_ENV_FILES=.env.envfile1,.env.envfile2
```

If `COMPOSE_ENV_FILES` is not set, and you don't provide `--env-file` in the CLI, Docker Compose uses the default behavior, which is to look for an `.env` file in the project directory.

#### COMPOSE\_DISABLE\_ENV\_FILE

Lets you disable the use of the default `.env` file. 

- Supported values: 
    - `true` or `1`, Compose ignores the `.env` file
    - `false` or `0`, Compose looks for an `.env` file in the project directory
- Defaults to: `0`

### Environment handling and container lifecycle

#### COMPOSE\_CONVERT\_WINDOWS\_PATHS

When enabled, Compose performs path conversion from Windows-style to Unix-style in volume definitions.

- Supported values: 
    - `true` or `1`, to enable
    - `false` or `0`, to disable
- Defaults to: `0`

#### COMPOSE\_IGNORE\_ORPHANS

When enabled, Compose doesn't try to detect orphaned containers for the project.

- Supported values: 
   - `true` or `1`, to enable
   - `false` or `0`, to disable
- Defaults to: `0`

#### COMPOSE\_REMOVE\_ORPHANS

When enabled, Compose automatically removes orphaned containers when updating a service or stack. Orphaned containers are those that were created by a previous configuration but are no longer defined in the current `compose.yaml` file.

- Supported values:
   - `true` or `1`, to enable automatic removal of orphaned containers
   - `false` or `0`, to disable automatic removal. Compose displays a warning about orphaned containers instead.
- Defaults to: `0`

#### COMPOSE\_PARALLEL\_LIMIT

Specifies the maximum level of parallelism for concurrent engine calls.

### Output 

#### COMPOSE\_ANSI

Specifies when to print ANSI control characters. 

- Supported values:
   - `auto`, Compose detects if TTY mode can be used. Otherwise, use plain text mode
   - `never`, use plain text mode
   - `always` or `0`, use TTY mode
- Defaults to: `auto`

#### COMPOSE\_STATUS\_STDOUT

When enabled, Compose writes its internal status and progress messages to `stdout` instead of `stderr`. 
The default value is false to clearly separate the output streams between Compose messages and your container's logs.

- Supported values:
   - `true` or `1`, to enable
   - `false` or `0`, to disable
- Defaults to: `0`

#### COMPOSE\_PROGRESS

{{< summary-bar feature_name="Compose progress" >}}

Defines the type of progress output, if `--progress` isn't used. 

Supported values are `auto`, `tty`, `plain`, `json`, and `quiet`.
Default is `auto`. 

### User experience

#### COMPOSE\_MENU

{{< summary-bar feature_name="Compose menu" >}}

When enabled, Compose displays a navigation menu where you can choose to open the Compose stack in Docker Desktop, switch on [`watch` mode](../file-watch.md), or use [Docker Debug](/reference/cli/docker/debug.md).

- Supported values:
   - `true` or `1`, to enable
   - `false` or `0`, to disable
- Defaults to: `1` if you obtained Docker Compose through Docker Desktop, otherwise the default is `0`

#### COMPOSE\_EXPERIMENTAL

{{< summary-bar feature_name="Compose experimental" >}}

This is an opt-out variable. When turned off it deactivates the experimental features.

- Supported values:
   - `true` or `1`, to enable
   - `false` or `0`, to disable
- Defaults to: `1`

## Unsupported in Compose V2

The following environment variables have no effect in Compose V2.
For more information, see [Migrate to Compose V2](/manuals/compose/releases/migrate.md).

- `COMPOSE_API_VERSION`
    By default the API version is negotiated with the server. Use `DOCKER_API_VERSION`.  
    See the [Docker CLI environment variable reference](/reference/cli/docker/#environment-variables) page.
- `COMPOSE_HTTP_TIMEOUT`
- `COMPOSE_TLS_VERSION`
- `COMPOSE_FORCE_WINDOWS_HOST`
- `COMPOSE_INTERACTIVE_NO_CLI`
- `COMPOSE_DOCKER_CLI_BUILD`
    Use `DOCKER_BUILDKIT` to select between BuildKit and the classic builder. If `DOCKER_BUILDKIT=0` then `docker compose build` uses the classic builder to build images.

