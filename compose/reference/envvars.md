---
description: Compose CLI environment variables
keywords: fig, composition, compose, docker, orchestration, cli, reference
title: Compose CLI environment variables
---

In this section you can find the list of pre-defined environment variables you can use to configure the Docker Compose command-line behavior.

**See also** [Declare default environment variables in file](../env-file.md) to check how to declare default environment variables in an environment file named `.env` placed in the project directory.

## COMPOSE\_PROJECT\_NAME

(Optional) Sets the project name. This value is prepended along with the service name to
the container on start up. 

For example, if your project name is `myapp` and it includes two services `db` and `web`, 
then Compose starts containers named `myapp-db-1` and `myapp-web-1` respectively.

* **Defaults to:** the `basename` of the project directory. 

**See also** the [`-p` command-line option](index.md#command-options-overview-and-help).

## COMPOSE\_FILE

Specifies the path to a Compose file. Specifying multiple Compose files is supported. 

* **Default behavior:** If not provided, Compose looks for a file named `docker-compose.yml` in the current directory and, if not found there, Compose then searches each parent directory recursively until a file by that name is found. 

* **Default separator:**  path separators are by default on:
    * Linux and macOS: `:` (colon), 
    * Windows: `;` (semicolon). 
The path separator can also be customized using `COMPOSE_PATH_SEPARATOR`.

Example: `COMPOSE_FILE=docker-compose.yml:docker-compose.prod.yml`. 

**See also** the [`-f` command-line option](index.md#command-options-overview-and-help).

## COMPOSE\_PROFILES

Specifies one or more profiles to be enabled on `compose-up` execution. 
Services with matching profiles are started *as well as any services for which no profile has been defined*.

For example, calling `docker-compose up`with `COMPOSE_PROFILES=frontend` selects services with the 
`frontend` profile as well as any services without a profile specified.


* **Default separator:**  specify a list of profiles using a comma as separator.
Example: `COMPOSE_PROFILES=frontend,debug`
This example would enable all services matching both the `frontend` and `debug` profiles *and* services without a profile.

**See also** [Using profiles with Compose](../profiles.md) and the [`--profile` command-line option](index.md#use---profile-to-specify-one-or-more-active-profiles).

## COMPOSE\_API\_VERSION

Specifies the API version of current environment. 
The Docker API only supports requests from clients reporting a specific version. 

If you receive a **client and server don't have same version** error while using `docker-compose`, 
you can work around this error by setting this environment variable. 
Set this version value in order to match the server version.



**Important**

> **This setting is only intended as a workaround and it is not officially supported.**
>
> Set and use this variable only as a last resort workaround when you need
> to *temporarily* run Docker with a mismatch between the client and server version. 
> For example, when you can upgrade the client but need to wait to upgrade the server.
>
> This mismatch between server and client prevents some features from working properly. 
> What features fail when using this configuration depends on the
> Docker client and server versions in question. 
>
> It is highly recommended that you upgrade client and server and remove this setting as 
> as soon as possible. Also, perform these actions to check if any problems you might 
> be having are resolved before notifying support.

## DOCKER\_HOST

Sets the URL of the Docker daemon. 
* **Defaults to:** `unix:///var/run/docker.sock`(same as with the Docker client).

## DOCKER\_TLS\_VERIFY

When set to anything other than an empty string, enables TLS communication with
the Docker daemon.

## DOCKER\_CERT\_PATH

Configures the path to the `ca.pem`, `cert.pem`, and `key.pem` files used for TLS verification.

* **Defaults to:** `~/.docker`.

## COMPOSE\_HTTP\_TIMEOUT

Defines the timeout period (in seconds) for requests to the Docker daemon after which Compose considers the requests have failed.

* **Defaults to:** 60 seconds.

## COMPOSE\_TLS\_VERSION

Defines which TLS version is used for TLS communication with the Docker daemon. 

* **Supported values:** `TLSv1`, `TLSv1_1`, `TLSv1_2`.
* **Defaults to:** `TLSv1`.

## COMPOSE\_CONVERT\_WINDOWS\_PATHS

When enabled, Compose performs path conversion from Windows-style to Unix-style in volume definitions.

* **Supported values:** 
    * `true` or `1`, to enable, 
    * `false` or `0`, to disable.

* **Defaults to:** `0`

## COMPOSE\_PATH\_SEPARATOR

Specifies a different path separator for items listed in `COMPOSE_FILE`. 

Default path separator in Linux and macOS the path separator is `:`, on Windows it is `;`.

## COMPOSE\_FORCE\_WINDOWS\_HOST

When enabled, volume declarations using the [short syntax](../compose-file/compose-file-v3.md#short-syntax-3) are parsed assuming the host path is a Windows path, even if Compose is
running on a UNIX-based system.

* **Supported values:** 
    * `true` or `1`, to enable, 
    * `false` or `0`, to disable.

## COMPOSE\_IGNORE\_ORPHANS

When enabled Compose doesn't try to detect orphaned containers for the project.

* **Supported values:** 
    * `true` or `1`, to enable, 
    * `false` or `0`, to disable.

## COMPOSE\_PARALLEL\_LIMIT

Sets a limit for the number of operations Compose can execute in parallel.

* **Supported values:** minimum value is `2`.
* **Defaults to:** `64`.

## COMPOSE\_INTERACTIVE\_NO\_CLI

When enabled Compose doesn't attempt to use the Docker CLI for interactive `run` and `exec` operations.

**Note:** on Windows the CLI is required for `run`and `exec` operations, so this option doesn't change the behavior of Compose on Windows.

* **Supported values:** 
    * `true` or `1`, to enable,
    * `false` or `0`, to disable.

## COMPOSE\_DOCKER\_CLI\_BUILD

Configure whether to use the Compose python client for building images or the native docker cli. 
By default, Compose uses the Docker CLI to perform builds,
which allows you to use [BuildKit](../../develop/develop-images/build_enhancements.md#to-enable-buildkit-builds) to perform builds.

Set `COMPOSE_DOCKER_CLI_BUILD=0` to disable native builds, and to use the built-in python client.

## Related information

- [User guide](../index.md)
- [Installing Compose](../install/index.md)
- [Compose file reference](../compose-file/index.md)
- [Environment file](../env-file.md)
