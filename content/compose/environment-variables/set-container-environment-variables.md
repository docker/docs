---
title: Set environment variables within your container's environment
description: Set, use, and manage environment variables in your container's environment
keywords: compose, orchestration, environment, environment variables, container environment variables
aliases:
- /compose/env/
- /compose/link-env-deprecated/
- /compose/environment-variables/set-environment-variables/
---

With Compose, there are two ways you can set environment variables in your containers with your Compose file. 

>**Tip**
>
> Don't use environment variables to pass sensitive information, such as passwords, in to your containers. Use [secrets](../use-secrets.md) instead.
{ .tip }

## Use the `environment` attribute

You can set environment variables directly in your container's environment with the
[`environment` attribute](../compose-file/05-services.md#environment) in your `compose.yml`. 

It supports both list and mapping syntax:

```yaml
services:
  webapp:
    environment:
      DEBUG: "true"
```
is equivalent to 
```yaml
services:
  webapp:
    environment:
      - DEBUG=true
```

See [`environment` attribute](../compose-file/05-services.md#environment) for more examples on how to use it. 

### Additional information 
- You can choose not to set a value and pass the environment variables from your shell straight through to your containers. It works in the same way as `docker run -e VARIABLE ...`:
  ```yaml
  web:
    environment:
      - DEBUG
  ```
  The value of the `DEBUG` variable in the container is taken from the value for the same variable in the shell in which Compose is run. 
  Note that in this case no warning is issued if the `DEBUG` variable in the shell environment is not set. 

- You can also take advantage of [interpolation](set-variables.md#interpolation-syntax).
  ```yaml
  web:
    environment:
      - DEBUG=${DEBUG}
  ```
  The result is similar to the one above but Compose gives you a warning if the `DEBUG` variable is not set in the shell environment.

## Use the `env_file` attribute


Container environment can also be set using [env files](env-file.md) along with the [`env_file` attribute](../compose-file/05-services.md#env_file).

```yaml
services:
  webapp:
    env_file: "webapp.env"
```

The [`env_file` attribute](../compose-file/05-services.md#env_file) lets you use multiple `.env` files in your Compose application. It also helps you keep your environment variables separate from your main configuration file, providing a more organized and secure way to manage sensitive information, as you do not need to place your `.env` file in the root of your project's directory. 

Using an env file for service environment allow to use the same file for use by a plain `docker run --env-file ...` command,
or to share same env file within multiple services without the need to duplicate a long `environment` yaml block.
In addition, as env file support [interpolation](set-variables.md#interpolation-syntax), it is possible to combine those with values set by `environment`. 

> **Important**
>
> Interpolation in env files is a Docker Compose CLI feature.
>
> It is not supported when running `docker run --env-file ...`.
{ .important }

### Additional information 
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

## Further resources
- [Understand environment variable precedence](envvars-precedence.md).
- [Set or change predefined environment variables](envvars.md)
- [Explore best practices](best-practices.md)
