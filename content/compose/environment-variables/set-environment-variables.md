---
title: Ways to set environment variables in Compose
description: How to set, use and manage environment variables in Compose
keywords: compose, orchestration, environment, env file
aliases:
- /compose/env/
- /compose/link-env-deprecated/
---

Environment variables are dealt with by either the Compose file or the CLI. Both have multiple ways you can substitute in or set your environment variables. This is outlined below. 

>**Tip**
>
> Don't use environment variables to pass sensitive information, such as passwords, in to your containers. Use [secrets](../use-secrets.md) instead.
{ .tip }

## Compose file

### Substitute with an `.env` file

The `.env` file is useful if you have multiple environment variables you need to store.

Below is a simple example: 

```console
$ cat .env
TAG=v1.5

$ cat compose.yml
services:
  web:
    image: "webapp:${TAG}"
```

When you run `docker compose up`, the `web` service defined in the Compose file substitutes in the
image `webapp:v1.5` which was set in the `.env` file. You can verify this with the
[config command](../../engine/reference/commandline/compose_config.md), which prints your resolved application config to the terminal:

```console
$ docker compose config

services:
  web:
    image: 'webapp:v1.5'
```

The `.env` file should be placed at the root of the project directory next to your `compose.yaml` file. You can use an alternative path with one of the following methods:
- The [`--file` option in the CLI](../reference/index.md#use--f-to-specify-name-and-path-of-one-or-more-compose-files) 
- The [`--env-file` option in the CLI](#substitute-with---env-file)
- Using the [`env_file` attribute in the Compose file](../compose-file/05-services.md#env_file)

For more information on formatting an environment file, see [Use an environment file](env-file.md).

> **Important**
>
> Substitution from `.env` files is a Docker Compose CLI feature.
>
> It is not supported by Swarm when running `docker stack deploy`.
{ .important }

### Use the `environment` attribute

You can set environment variables in a service's containers with the
[`environment` attribute](../compose-file/05-services.md#environment) in your Compose file. It works in the same way as `docker run -e VARIABLE=VALUE ...`

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

The value of the `DEBUG` variable in the container is taken from the value for the same variable in the shell in which Compose is run. 
Note that in this case no warning will be issued if the `DEBUG` variable in the shell environment is not set. 

You can also explicitly assign a variable using a Bash-like syntax `${DEBUG}`:

```yaml
web:
  environment:
    - DEBUG=${DEBUG}
```

The result is similar to the one above but Compose will give you a warning if the `DEBUG` variable is not set in the shell environment.

See [`environment` attribute](../compose-file/05-services.md#environment) and [variable interpolation](../compose-file/12-interpolation/) for more information.

### Use the `env_file` attribute

You can pass multiple environment variables from an external file through to
a service's containers with the [`env_file` option](../compose-file/05-services.md#env_file). This works in the same way as `docker run --env-file=FILE ...`:

```yaml
web:
  env_file:
    - web-variables.env
```

If multiple files are specified, they are evaluated in order and can override values set in previous files.

> **Note**
>
>With this option, environment variables declared in the file cannot then be referenced again separately in the Compose file or used to configure Compose.

See [`env_file` attribute](../compose-file/05-services.md#env_file) for more information.

### Substitute from the shell 

It's possible to use environment variables in your shell to populate values inside a Compose file. Compose uses the variable values from the shell environment in which `docker compose` is run.

For example, suppose the shell contains `POSTGRES_VERSION=9.3` and you supply the following configuration:

```yaml
db:
  image: "postgres:${POSTGRES_VERSION}"
```

When you run `docker compose up` with this configuration, Compose looks for the `POSTGRES_VERSION` environment variable in the shell and substitutes its value in. For this example, Compose resolves the image to `postgres:9.3` before running the configuration.

If an environment variable is not set, Compose substitutes with an empty string. In the example above, if `POSTGRES_VERSION` is not set, the value for the image option is `postgres:`.

> **Note**
>
> `postgres:` is not a valid image reference. Docker expects either a reference without a tag, like `postgres` which defaults to the latest image, or with a tag such as `postgres:15`.

> **Important**
>
> Values set in the shell environment override those set in the `.env` file, the `environment` attribute, and the `env_file` attribute. For more information, see [Environment variable precedence](envvars-precedence.md).
{ .important }

## CLI

### Substitute with `--env-file`

You can set default values for multiple environment variables, in an [environment file](env-file.md) and then pass the file as an argument in the CLI.

The advantage of this method is that you can store the file anywhere and name it appropriately, for example, `.env.ci`, `.env.dev`, `.env.prod`. This file path is relative to the current working directory where the Docker Compose command is executed. Passing the file path is done using the `--env-file` option:

```console
$ docker compose --env-file ./config/.env.dev up
```

In the following example, there are two environment files, `.env` and `.env.dev`. Both have different values set for `TAG`. 

```console
$ cat .env
TAG=v1.5

$ cat ./config/.env.dev
TAG=v1.6


$ cat compose.yml
services:
  web:
    image: "webapp:${TAG}"
```

If the `--env-file` is not used in the command line, the `.env` file is loaded by default:

```console
$ docker compose config
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
> Values set in the shell environment override those set when using the `--env-file` argument in the CLI. For more information, see [Environment variable precedence](envvars-precedence.md)
{ .important }

### Set environment variables with `docker compose run --env`

Similar to `docker run --env`, you can set environment variables in a one-off
container with `docker compose run --env` or its short form `docker compose run -e`:

```console
$ docker compose run -e DEBUG=1 web python console.py
```

You can also pass a variable from the shell by not giving it a value:

```console
$ docker compose run -e DEBUG web python console.py
```

The value of the `DEBUG` variable in the container is taken from the value for
the same variable in the shell in which Compose is run.