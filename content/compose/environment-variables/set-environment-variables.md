---
title: Ways to set environment variables with Compose
description: How to set, use, and manage environment variables with Compose
keywords: compose, orchestration, environment, env file
aliases:
- /compose/env/
- /compose/link-env-deprecated/
---

With Compose, there are multiple ways you can set environment variables in your containers.

>**Tip**
>
> Don't use environment variables to pass sensitive information, such as passwords, in to your containers. Use [secrets](../use-secrets.md) instead.
{ .tip }

## Compose file

### Use the `environment` attribute

You can declare variables to be set in your container using service's `environment` attribute. This one supports both a 
list and mapping syntax:

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

As any other attribute in your compose file, `environment` support value to be set by interpolation, as illustrated here:
```yaml
services:
  webapp:
    environment:
      - DEBUG=${COMPOSE_DEBUG}
```
Docker Compose then replaces `${COMPOSE_DEBUG}` with the [value from the environment](interpolation.md#Interpolation_variables). A warning is issued if the `COMPOSE_DEBUG` variable 
is not defined. 

You also can use `environment` to directly replicate a variable from local environment into container's environment:
```yaml
services:
  webapp:
    environment:
      - DEBUG
```

The value of the `DEBUG` variable in the container is taken from the value for the same variable in the environment 
in which Compose is run. If no such variable is defined in the local environment, variable is neither set on container's
environment.

### Use the `env_file` attribute

Container environment can also be set using [env files](env-file.md):
```yaml
services:
  webapp:
    env_file: "webapp.env"
```

Using an env file for service environment allow to use the same file for use by a plain `docker run --env-file ...` command,
or to share same env file within multiple services without the need to duplicate a long `environment` yaml block.
In addition, as env file support [interpolation](interpolation.md), it is possible to combine those with values set by `environment`. Read [precedence rules](envvars-precedence.md) in such case to understand how value is when declared in
multiple places.

> **Important**
>
> Interpolation in env files is a Docker Compose CLI feature.
>
> It is not supported when running `docker run --env-file ...`.
{ .important }


#### Additional information 

- As of Docker Compose version 2.24.0, you can set your `.env` file to be optional by using the `env_file` attribute. When `required` is set to `false` and the `.env` file is missing, Compose silently ignores the entry.
  ```yaml
  env_file:
    - path: ./default.env
      required: true # default
    - path: ./override.env
      required: false
  ``` 

- You can use multiple `.env` files in your `compose.yml` with the [`env_file` attribute](../compose-file/05-services.md#env_file), and Docker Compose reads them in the order specified. If the same variable is defined in multiple files, the last definition takes precedence:
  ```yaml
  services:
    webapp:
      image: my-webapp-image
      env_file:
        - .env
        - .env.override
  ```

## Further resources

- [Understand environment variable precedence](envvars-precedence.md).
- [Set or change predefined environment variables](envvars.md)
- [Explore best practices](best-practices.md)
- [Understand the syntax and formatting guidelines for environment files](env-file.md)
