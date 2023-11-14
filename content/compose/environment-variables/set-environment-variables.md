---
title: Ways to set environment variables with Compose
description: How to set, use and manage environment variables with Compose
keywords: compose, orchestration, environment, env file
aliases:
- /compose/env/
- /compose/link-env-deprecated/
---

With Compose, there are multiple ways you can set environment variables in your containers. You can use either your Compose file, or the CLI. 

Be aware that each method is subject to [environment variable precedence](envvars-precedence.md).

>**Tip**
>
> Don't use environment variables to pass sensitive information, such as passwords, in to your containers. Use [secrets](../use-secrets.md) instead.
{ .tip }

## Compose file

### Substitute with an `.env` file

An `.env` file in Docker Compose is a text file used to define environment variables that should be made available to Docker containers when running `docker compose up`. This file typically contains key-value pairs of environment variables, and it allows you to centralize and manage configuration in one place. The `.env` file is useful if you have multiple environment variables you need to store.

The `.env` file is the default method for setting environment variables in your containers. The `.env` file should be placed at the root of the project directory next to your `compose.yaml` file. For more information on formatting an environment file, see [Use an environment file](env-file.md).

Below is a simple example: 

```console
$ cat .env
TAG=v1.5

$ cat compose.yml
services:
  web:
    image: "webapp:${TAG}"
```

When you run `docker compose up`, the `web` service defined in the Compose file [interpolates](env-file.md#interpolation) in the
image `webapp:v1.5` which was set in the `.env` file. You can verify this with the
[config command](../../engine/reference/commandline/compose_config.md), which prints your resolved application config to the terminal:

```console
$ docker compose config

services:
  web:
    image: 'webapp:v1.5'
```

#### Additional information 

- If you define an environment variable in your `.env` file, you can reference it directly in your `compose.yml` with the `environment` attribute. For example, if your `.env` file contains the environment variable `DATABASE_URL=mysql://user:password@db:3306/mydatabase` and your `compose.yml` file looks like this:
  ```yaml
    services:
      webapp:
        image: my-webapp-image
        environment:
          - DATABASE_URL=${DATABASE_URL}
  ```
  Docker Compose replaces `${DATABASE_URL}` with the value from the `.env` file
- You can use multiple `.env` files in your `compose.yml` with the `env_file` attribute, and Docker Compose reads them in the order specified. If the same variable is defined in multiple files, the last definition takes precedence:
  ```yaml
  services:
    webapp:
      image: my-webapp-image
      env_file:
        - .env
        - .env.override
  ```
- You can place your `.env` file in a location other than the root of your project's directory, and then use one of the following methods so Compose can navigate to it:
  - The [`--file` option in the CLI](../reference/index.md#use--f-to-specify-name-and-path-of-one-or-more-compose-files) 
  - Using the [`env_file` attribute in the Compose file](../compose-file/05-services.md#env_file)
- Values in your `.env` file can be overriden from the command line by using [`docker-compose up -e`](#set-environment-variables-with-docker-compose-run---env).
- Your `.env` file can be overriden by another `.env` if it is [substituted with `--env-file`](#substitute-with---env-file).

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
- You can choose not to set a value and pass the environment variables from your shell straight through to a
service's containers. It works in the same way as `docker run -e VARIABLE ...`:
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
- Values in your `.env` files can be overriden from the command line by using [`docker-compose up -e`](#set-environment-variables-with-docker-compose-run---env).
- Your `.env` files can be overriden by another `.env` if it is [substituted with `--env-file`](#substitute-with---env-file).

### Substitute from the shell 

You can use existing environment variables from your host machine or from the shell environment where you execute docker-compose commands. This allows you to dynamically inject values into your Docker Compose configuration at runtime.

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

- You can use multiple --env-file options to specify multiple environment files, and Docker Compose will read them in order. Later files can override variables from earlier files.
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