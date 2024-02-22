---
title: Interpolation of variables in Compose files
description: How to set, use, and manage variables in a Compose file with interpolation
keywords: compose, orchestration, environment, env file, interpolation
---

A Compose file can uses variables to offer more flexibility. If you want to quickly switch 
between image tags to test multiple versions, or want to adjust a volume source to your local
environment, you don't need to edit the compose file, but can just set variables:

```yaml
services:
  web:
    image: "webapp:${TAG}"
```

When you run `docker compose up`, the `web` service defined in the Compose file get 
image defined according to variable you have set. You can verify this with the
[config command](../../reference/cli/docker/compose/config.md), which prints your resolved application config to the terminal:

```console
$ TAG=v1.5 docker compose config

services:
  web:
    image: 'webapp:v1.5'
```

This feature, known as "interpolation", allows to adjust the compose model to your needs without
having to make changes to the compose file.

### Interpolation syntax

Interpolation is inspired by bash variable substitution syntax.
Both braced (`${VAR}`) and unbraced (`$VAR`) expressions are supported.

For braced expressions, the following formats are supported:
- Direct substitution
  - `${VAR}` -> value of `VAR`
- Default value
  - `${VAR:-default}` -> value of `VAR` if set and non-empty, otherwise `default`
  - `${VAR-default}` -> value of `VAR` if set, otherwise `default`
- Required value
  - `${VAR:?error}` -> value of `VAR` if set and non-empty, otherwise exit with error
  - `${VAR?error}` -> value of `VAR` if set, otherwise exit with error
- Alternative value
  - `${VAR:+replacement}` -> `replacement` if `VAR` is set and non-empty, otherwise empty
  - `${VAR+replacement}` -> `replacement` if `VAR` is set, otherwise empty

For more information, see [Interpolation](../compose-file/12-interpolation.md) in the Compose Specification. 


### Interpolation variables

The variables used to run interpolation and resolve values in your Compose file are those set by your 
local environment. Multiple sources are involves:

- environment variables in your environment, either set globaly or explicitly by the command line
- an optional `.env` file in the project directory (parent folder of your `compose.yaml` file), or alternative 
  env file(s) explicitly declared using `--env-file` flag 

> **Important**
>
> [`env-file`](set-environment-variables.md#Use_the_env_file_attribute) service attribute is not involved setting variables
> used for interpolation. It is **only** used to define container's environment
> 
{ .important }


Env file (see [dedicated page](env-file.md) for syntax) also support variable interpolation, and as such can rely
on variable set by earlier declared variables, as illustrated with this example:

```console
$ cat .env
## define COMPOSE_DEBUG based on DEV_MODE, defaults to false
COMPOSE_DEBUG=${DEV_MODE:-false}

$ cat compose.yaml 
  services:
    webapp:
      image: my-webapp-image
      environment:
        - DEBUG=${COMPOSE_DEBUG}

$ DEV_MODE=true docker compose config
services:
  webapp:
    environment:
      DEBUG: "true"
```


#### Additional information 

- As of Docker Compose version 2.24.0, you can set your `.env` file to be optional by using the `env_file` attribute. When `required` is set to `false` and the `.env` file is missing, Compose silently ignores the entry.
  ```yaml
  env_file:
    - path: ./default.env
      required: true # default
    - path: ./override.env
      required: false
  ``` 

- If you define an environment variable in your `.env` file, you can reference it directly in your `compose.yml` with the [`environment` attribute](../compose-file/05-services.md#environment). For example, if your `.env` file contains the environment variable `DEBUG=1` and your `compose.yml` file looks like this:
  ```yaml
    services:
      webapp:
        image: my-webapp-image
        environment:
          - DEBUG=${DEBUG}
  ```
  Docker Compose replaces `${DEBUG}` with the value from the `.env` file
- You can use multiple `.env` files in your `compose.yml` with the [`env_file` attribute](../compose-file/05-services.md#env_file), and Docker Compose reads them in the order specified. If the same variable is defined in multiple files, the last definition takes precedence:
  ```yaml
  services:
    webapp:
      image: my-webapp-image
      env_file:
        - .env
        - .env.override
  ```
- You can place your `.env` file in a location other than the root of your project's directory, and then use one of the following methods so Compose can navigate to it:
  - The [`--env-file` option in the CLI](#substitute-with---env-file)
  - Using the [`env_file` attribute in the Compose file](../compose-file/05-services.md#env_file)
- Values in your `.env` file can be overridden from the command line by using [`docker compose run -e`](#set-environment-variables-with-docker-compose-run---env).
- Your `.env` file can be overridden by another `.env` if it is [substituted with `--env-file`](#substitute-with---env-file).

> **Important**
>
> Substitution from `.env` files is a Docker Compose CLI feature.
>
> It is not supported by Swarm when running `docker stack deploy`.
{ .important }

### Use the `environment` attribute

You can set environment variables directly in your Compose file without using an `.env` file, with the
[`environment` attribute](../compose-file/05-services.md#environment) in your `compose.yml`. It works in the same way as `docker run -e VARIABLE=VALUE ...`

```yaml
web:
  environment:
    - DEBUG=1
```

See [`environment` attribute](../compose-file/05-services.md#environment) for more examples on how to use it. 

#### Additional information 
- You can choose not to set a value and pass the environment variables from your shell straight through to your containers. It works in the same way as `docker run -e VARIABLE ...`:
  ```yaml
  web:
    environment:
      - DEBUG
  ```
  The value of the `DEBUG` variable in the container is taken from the value for the same variable in the shell in which Compose is run. 
  Note that in this case no warning is issued if the `DEBUG` variable in the shell environment is not set. 

- You can also take advantage of [interpolation](env-file.md#interpolation).
  ```yaml
  web:
    environment:
      - DEBUG=${DEBUG}
  ```
  The result is similar to the one above but Compose gives you a warning if the `DEBUG` variable is not set in the shell environment.

### Use the `env_file` attribute

The [`env_file` attribute](../compose-file/05-services.md#env_file) lets you use multiple `.env` files in your Compose application. It also helps you keep your environment variables separate from your main configuration file, providing a more organized and secure way to manage sensitive information, as you do not need to place your `.env` file in the root of your project's directory. 

It works in the same way as `docker run --env-file=FILE ...`.

```yaml
web:
  env_file:
    - web-variables.env
```
#### Additional information 
- If multiple files are specified, they are evaluated in order and can override values set in previous files.
- Environment variables declared in the `.env` file cannot then be referenced again separately in the Compose file.
- If you use both the `env_file` and `environment` attribute, environment variables set by `environment` take precedence.
- The paths to your `.env` file, specified in the `env_file` attribute,  are relative to the location of your `compose.yml` file. 
- Values in your `.env` files can be overridden from the command line by using [`docker compose run -e`](#set-environment-variables-with-docker-compose-run---env).
- Your `.env` files can be overriden by another `.env` if it is [substituted with `--env-file`](#substitute-with---env-file).
- As of Docker Compose version 2.24.0, you can set your `.env` file to be optional by using the `required` field. When `required` is set to `false` and the `.env` file is missing,
Compose silently ignores the entry.
  ```yaml
  env_file:
    - path: ./default.env
      required: true # default
    - path: ./override.env
      required: false
  ``` 

### Substitute from the shell 

You can use existing environment variables from your host machine or from the shell environment where you execute `docker compose` commands. This allows you to dynamically inject values into your Docker Compose configuration at runtime.

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

The advantage of this method is that you can store the file anywhere and name it appropriately, for example, 

 This file path is relative to the current working directory where the Docker Compose command is executed. Passing the file path is done using the `--env-file` option:

```console
$ docker compose --env-file ./config/.env.dev up
```

#### Additional information 
- This method is useful if you want to temporarily override an `.env` file that is already referenced in your `compose.yml` file. For example you may have different `.env` files for production ( `.env.prod`) and testing (`.env.test`).
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

- You can use multiple `--env-file` options to specify multiple environment files, and Docker Compose reads them in order. Later files can override variables from earlier files.
  ```console
  $ docker compose --env-file .env --env-file .env.override up
  ```
- You can override specific environment variables from the command line when starting containers. 
  ```console
  $ docker compose --env-file .env.dev up -e DATABASE_URL=mysql://new_user:new_password@new_db:3306/new_database
  ```

### Set environment variables with `docker compose run --env`

Similar to `docker run --env`, you can set environment variables temporarily with `docker compose run --env` or its short form `docker compose run -e`:

```console
$ docker compose run -e DEBUG=1 web python console.py
```
#### Additional information 

- You can also pass a variable from the shell by not giving it a value:

  ```console
  $ docker compose run -e DEBUG web python console.py
  ```

  The value of the `DEBUG` variable in the container is taken from the value for the same variable in the shell in which Compose is run.

## Further resources
- [Understand environment variable precedence](envvars-precedence.md).
- [Set or change predefined environment variables](envvars.md)
- [Explore best practices](best-practices.md)
- [Understand the syntax and formatting guidelines for environment files](env-file.md)
