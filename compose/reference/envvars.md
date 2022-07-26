---
description: Compose CLI environment variables
keywords: fig, composition, compose, docker, orchestration, cli, reference
title: Compose CLI environment variables
---

In this section you can find the list of pre-defined environment variables you can use to configure the Docker Compose command-line behavior.  
**See also** [Declare default environment variables in file](../env-file.md) to check how to declare default environment variables in an environment file named `.env` placed in the project directory.

## COMPOSE\_PROJECT\_NAME

Sets the project name. This value is prepended along with the service name to
the container's name on startup.

For example, if your project name is `myapp` and it includes two services `db` and `web`, 
then Compose starts containers named `myapp-db-1` and `myapp-web-1` respectively.

* **Defaults to:** the `basename` of the project directory.

**See also** the [command-line options overview](index.md#command-options-overview-and-help) and [using `-p` to specify a project name](index.md#use--p-to-specify-a-project-name).

## COMPOSE\_FILE

Specifies the path to a Compose file. Specifying multiple Compose files is supported.

* **Default behavior:** If not provided, Compose looks for a file named `compose.yaml` or `docker-compose.yaml` in the current directory and, if not found, then Compose searches each parent directory recursively until a file by that name is found.

* **Default separator:** When specifying multiple Compose files, the path separators are, by default, on:
    * Mac and Linux: `:` (colon),
    * Windows: `;` (semicolon).

The path separator can also be customized using `COMPOSE_PATH_SEPARATOR`.  
Example: `COMPOSE_FILE=docker-compose.yml:docker-compose.prod.yml`.  
**See also** the [command-line options overview](index.md#command-options-overview-and-help) and [using `-f` to specify name and path of one or more Compose files](index.md#use--f-to-specify-name-and-path-of-one-or-more-compose-files).

## COMPOSE\_PROFILES

Specifies one or more profiles to be enabled on `compose up` execution.
Services with matching profiles are started **as well as any services for which no profile has been defined**.

For example, calling `docker compose up`with `COMPOSE_PROFILES=frontend` selects services with the 
`frontend` profile as well as any services without a profile specified.


* **Default separator:**  specify a list of profiles using a comma as separator.  
Example: `COMPOSE_PROFILES=frontend,debug`  
This example would enable all services matching both the `frontend` and `debug` profiles **and services without a profile**.

**See also** [Using profiles with Compose](../profiles.md) and the [`--profile` command-line option](index.md#use---profile-to-specify-one-or-more-active-profiles).

## DOCKER\_HOST

Sets the URL of the Docker daemon. 
* **Defaults to:** `unix:///var/run/docker.sock`(same as with the Docker client).

## DOCKER\_TLS\_VERIFY

See `DOCKER_TLS_VERIFY` on the [Use the Docker command line](../../../engine/reference/commandline/cli/#environment-variables){:target="_blank" rel="noopener" class="_"} page.

## DOCKER\_CERT\_PATH

Configures the path to the `ca.pem`, `cert.pem`, and `key.pem` files used for TLS verification.  
* **Defaults to:** `~/.docker`.

See, `DOCKER_CERT_PATH` on the [Use the Docker command line](../../../engine/reference/commandline/cli/#environment-variables){:target="_blank" rel="noopener" class="_"} page.

## COMPOSE\_CONVERT\_WINDOWS\_PATHS

When enabled, Compose performs path conversion from Windows-style to Unix-style in volume definitions.

* **Supported values:** 
    * `true` or `1`, to enable,
    * `false` or `0`, to disable.
* **Defaults to:** `0`.

## COMPOSE\_PATH\_SEPARATOR

Specifies a different path separator for items listed in `COMPOSE_FILE`.

* **Defaults to:**
    * On Mac and Linux to `:`,
    * On Windows to`;`.

## COMPOSE\_IGNORE\_ORPHANS

When enabled, Compose doesn't try to detect orphaned containers for the project.

* **Supported values:** 
    * `true` or `1`, to enable,
    * `false` or `0`, to disable.
* **Defaults to:** `0`.

## Deprecated in Compose v2

>**Important**
>
> The environment variables listed below are deprecated in v2.  

### COMPOSE\_API\_VERSION

Deprecated in v2.  
By default the API version is negotiated with the server. Use `DOCKER_API_VERSION`.  
See `DOCKER_API_VERSION` on the [Use the Docker command line](../../../engine/reference/commandline/cli/#environment-variables){:target="_blank" rel="noopener" class="_"} page.

### COMPOSE\_HTTP\_TIMEOUT

Deprecated in v2.

### COMPOSE\_TLS\_VERSION

Deprecated in v2.

### COMPOSE\_FORCE\_WINDOWS\_HOST

Deprecated in v2.

### COMPOSE\_PARALLEL\_LIMIT

Deprecated in v2.

### COMPOSE\_INTERACTIVE\_NO\_CLI

Deprecated in v2.  
As v2 now uses the vendored code of [Docker CLI](https://github.com/docker/cli){:target="_blank" rel="noopener" class="_"}.

### COMPOSE\_DOCKER\_CLI\_BUILDx

Deprecated in v2.  
Use `DOCKER_BUILDKIT` to select between BuildKit and the classic builder. If `DOCKER_BUILDKIT=0` then `docker build` uses the classic builder to build images.

## Related information

- [User guide](../index.md)
- [Installing Compose](../install/index.md)
- [Compose file reference](../compose-file/index.md)
- [Environment file](../env-file.md)
