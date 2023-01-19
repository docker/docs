---
title: Ways to use environment variables in Compose
description: How to set, use and manage environment variables in Compose
keywords: compose, orchestration, environment, env file
redirect_from:
- /compose/env/
- /compose/link-env-deprecated/
---

Environment variables are dealt with by either the Compose file or the CLI, and both have multiple ways you can substitute in your environment variables. This is outlined below. 

## Compose file

### Substitute with an `.env` file

The `.env` file is useful if you have multiple environment variables. 

You can set default values for any environment variables referenced in the
Compose file, or used to configure Compose, in an [environment file](env-file.md)
named `.env`. The `.env` file path is as follows:

  - The `.env` file is placed at the base of the project directory.
  - The project directory can be explicitly defined with the `--file` option or `COMPOSE_FILE`
  environment variable. Otherwise, it is the current working directory where the `docker compose` command is run.
  - For versions older than `v1.28`, Compose might have trouble resolving `.env` file with `--file` or `COMPOSE_FILE`. To work around this, use `--project-directory`. This overrides the path for the `.env` file. This inconsistency is addressed in `v1.28` by limiting the file path to the project directory.

Below is an example: 

```console
$ cat .env
TAG=v1.5

$ cat docker-compose.yml
services:
  web:
    image: "webapp:${TAG}"
```

When you run `docker compose up`, the `web` service defined above uses the
image `webapp:v1.5`. You can verify this with the
[convert command](../engine/reference/commandline/compose_convert.md), which prints your resolved application config to the terminal:

```console
$ docker compose convert

services:
  web:
    image: 'webapp:v1.5'
```

> **Important**
>
>The `.env` file feature only works when you use the `docker compose up` command and does not work with `docker stack deploy`.
{: .important}

### Substitute from the shell 

It's possible to use environment variables in your shell to populate values
inside a Compose file:

```yaml
web:
  image: "webapp:${TAG}"
```

Compose uses the variable values from the shell environment in which `docker compose` is run. For example, suppose the shell contains `POSTGRES_VERSION=9.3` and you supply the following configuration:

```console
db:
  image: "postgres:${POSTGRES_VERSION}"
```

When you run `docker compose up` with this configuration, Compose looks for the `POSTGRES_VERSION` environment variable in the shell and substitutes its value in. For this example, Compose resolves the image to `postgres:9.3` before running the configuration.

If an environment variable is not set, Compose substitutes with an empty string. In the example above, if `POSTGRES_VERSION` is not set, the value for the image option is `postgres:.`

> **Important**
>
> Values set in the shell environment override those set in the `.env` file. For more information, see [Environment variable precedence](envvars-precedence.md).

### Use the `environment` attribute

You can set environment variables in a service's containers with the
['environment' attribute](compose/compose-file.md#environment). It works in the same way as `docker run -e VARIABLE=VALUE ...`

```yaml
web:
  environment:
    - DEBUG=1
```

You can choose not to set a value and pass the environment variables from your shell straight through to a
service's containers. It works in the same way as `docker run -e VARIABLE ...`:

```yaml
web:
  environment:
    - DEBUG
```

The value of the `DEBUG` variable in the container is taken from the value for
the same variable in the shell in which Compose is run.

### Use the `env_file` attribute

You can pass multiple environment variables from an external file through to
a service's containers with the ['env_file' option](compose/compose-file.md#env_file). This works in the same way as `docker run --env-file=FILE ...`:

```yaml
web:
  env_file:
    - web-variables.env
```
> **Note**
>
> By using this option, environment variables declared in the file cannot be referenced in the Compose file or used to configure Compose.

## CLI

### Substitute with `--env-file`

You can set default values for multiple environment variables , in an [environment file](env-file.md) and then pass the file as an argument in the CLI. 

The advantage of this method is that you can store the file anywhere and name it appropriately, for example, `.env.ci`, `.env.dev`, `.env.prod`. Passing the file path is done using the `--env-file` option:

```console
$ docker compose --env-file ./config/.env.dev up
```

This file path is relative to the current working directory where the Docker Compose
command is executed.

```console
$ cat .env
TAG=v1.5

$ cat ./config/.env.dev
TAG=v1.6


$ cat docker-compose.yml
services:
  web:
    image: "webapp:${TAG}"
```

The `.env` file is loaded by default:

```console
$ docker compose convert
services:
  web:
    image: 'webapp:v1.5'
```

Passing the `--env-file` argument overrides the default file path:

```console
$ docker compose --env-file ./config/.env.dev config
services:
  web:
    image: 'webapp:v1.6'
```

When an invalid file path is being passed as an `--env-file` argument, Compose returns an error:

```console
$ docker compose --env-file ./doesnotexist/.env.dev  config
ERROR: Couldn't find env file: /home/user/./doesnotexist/.env.dev
```

> **Important**
>
> Values set in the shell environment override those set in the `.env` file. For more information, see [Environment variable precedence](envvars-precedence.md)

### Set environment variables with 'docker compose run -e'

Similar to `docker run -e`, you can set environment variables on a one-off
container with `docker compose run -e`:

```console
$ docker compose run -e DEBUG=1 web python console.py
```

You can also pass a variable from the shell by not giving it a value:

```console
$ docker compose run -e DEBUG web python console.py
```

The value of the `DEBUG` variable in the container is taken from the value for
the same variable in the shell in which Compose is run.

### Use pre-defined environment variables

Several environment variables are available for you to configure the Docker
Compose command-line behavior. They begin with `COMPOSE_` or `DOCKER_`, and are
documented in [CLI Environment Variables](reference/envvars.md).
