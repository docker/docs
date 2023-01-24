---
description: Compose CLI environment variables
keywords: fig, composition, compose, docker, orchestration, cli, reference
title: Change pre-defined environment variables
redirect_from:
- /compose/reference/envvars/
---

Compose already comes with pre-defined environment variables. 

This page contains information on how you can change the following pre-defined environment variables if you need to:

- `COMPOSE_CONVERT_WINDOWS_PATHS`
- `COMPOSE_FILE`
- `COMPOSE_PROFILES`
- `COMPOSE_PROJECT_NAME`
- `DOCKER_CERT_PATH`
- `DOCKER_HOST`
- `DOCKER_TLS_VERIFY`

## Methods

You can change the pre-defined environment variables:
- Within your Compose file using the [`environment` attribute](set-environment-variables.md#use-the-environment-attribute)
- With an [environment file](env-file.md) 
- From the command line
- From your [shell](set-environment-variables.md#substitute-from-the-shell)

When changing or setting any environment variables, be aware of [Environment variable precedence](envvars-precedence.md).

## Configure

### COMPOSE\_PROJECT\_NAME

Sets the project name. This value is prepended along with the service name to
the container's name on startup.

For example, if your project name is `myapp` and it includes two services `db` and `web`, 
then Compose starts containers named `myapp-db-1` and `myapp-web-1` respectively.

It defaults to the `basename` of the project directory.

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

### DOCKER\_HOST

Sets the URL of the Docker daemon. 
Defaults to `unix:///var/run/docker.sock`(same as with the Docker client).

### DOCKER\_TLS\_VERIFY

See `DOCKER_TLS_VERIFY` on the [Use the Docker command line](../../../engine/reference/commandline/cli/#environment-variables){:target="_blank" rel="noopener" class="_"} page.

### DOCKER\_CERT\_PATH

Configures the path to the `ca.pem`, `cert.pem`, and `key.pem` files used for TLS verification.  
Defaults to `~/.docker`.

See, `DOCKER_CERT_PATH` on the [Use the Docker command line](../../../engine/reference/commandline/cli/#environment-variables){:target="_blank" rel="noopener" class="_"} page.

### COMPOSE\_CONVERT\_WINDOWS\_PATHS

When enabled, Compose performs path conversion from Windows-style to Unix-style in volume definitions.

* Supported values: 
    * `true` or `1`, to enable,
    * `false` or `0`, to disable.
* Defaults to: `0`.

### COMPOSE\_PATH\_SEPARATOR

Specifies a different path separator for items listed in `COMPOSE_FILE`.

* Defaults to:
    * On Mac and Linux to `:`,
    * On Windows to`;`.

### COMPOSE\_IGNORE\_ORPHANS

When enabled, Compose doesn't try to detect orphaned containers for the project.

* Supported values: 
    * `true` or `1`, to enable,
    * `false` or `0`, to disable.
* Defaults to: `0`.

## Deprecated in Compose v2

The pre-definded environment variables listed below are deprecated in [V2](../compose-v2/index.md).  

- COMPOSE\_API\_VERSION
    By default the API version is negotiated with the server. Use `DOCKER_API_VERSION`.  
    See `DOCKER_API_VERSION` on the [Use the Docker command line](../../../engine/reference/commandline/cli/#environment-variables){:target="_blank" rel="noopener" class="_"} page.
- COMPOSE\_HTTP\_TIMEOUT
- COMPOSE\_TLS\_VERSION
- COMPOSE\_FORCE\_WINDOWS\_HOST
- COMPOSE\_PARALLEL\_LIMIT
- COMPOSE\_INTERACTIVE\_NO\_CLI
    V2 now uses the vendored code of [Docker CLI](https://github.com/docker/cli){:target="_blank" rel="noopener" class="_"}.
- COMPOSE\_DOCKER\_CLI\_BUILDx
    Use `DOCKER_BUILDKIT` to select between BuildKit and the classic builder. If `DOCKER_BUILDKIT=0` then `docker build` uses the classic builder to build images.
