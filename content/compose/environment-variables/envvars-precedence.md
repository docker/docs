---
title: Environment variables precedence in Docker Compose
description: Scenario overview illustrating how environment variables are resolved
  in Compose
keywords: compose, environment, env file
aliases:
- /compose/envvars-precedence/
---

When the same environment variable is set in multiple sources, Docker Compose follows a precedence rule to determine the value for that variable.

This page contains information on the level of precedence each method of setting environmental variables takes.

The order of precedence (highest to lowest) is as follows:
1. Set using [`docker compose run -e` in the CLI](set-environment-variables.md#set-environment-variables-with-docker-compose-run---env)
2. Substituted from your [shell](set-environment-variables.md#substitute-from-the-shell), default `.env` file
or [`--env-file` argument](set-environment-variables.md#substitute-with---env-file) in the CLI.
3. Set using just the [`environment` attribute in the Compose file](set-environment-variables.md#use-the-environment-attribute)
4. Set by values from the [`env_file` attribute in the Compose file](set-environment-variables.md#use-the-environment-attribute)
5. Set in a container image in the [ENV directive](../../reference/dockerfile.md#env).
   Having any `ARG` or `ENV` setting in a `Dockerfile` evaluates only if there is no Docker Compose entry for `environment`, `env_file` or `run --env`.

## Simple example

In the example below, we set a different value for the same environment variable in an `.env` file and with the `environment` attribute in the Compose file:

```console
$ cat ./webapp.env
NODE_ENV=test

$ cat compose.yml
services:
  webapp:
    image: 'webapp'
    env_file:
     - webapp.env
    environment:
     - NODE_ENV=production
```

The environment variable defined with the `environment` attribute takes precedence.

```console
$ docker compose run webapp env | grep NODE_ENV
NODE_ENV=production
```

## Advanced example 

The following table uses `VALUE`, an environment variable defining the version for an image, as an example:

### How the table works

Each column represents a context from where you can set a value, or substitute in a value for `VALUE`.

The columns `Host OS environment` and `.env file` is listed only as an illustration lookup. In reality, they don't result in a variable in the container by itself.

Each row represents a combination of contexts where `VALUE` is set, substituted, or both. The **Result** column indicates the final value for `VALUE` in each scenario.


|  # |  `docker compose run`  |  `environment` attribute  |  `env_file` attribute  |  Image `ENV` |  `Host OS` environment  |  `.env` file      | |  Result  |
|:--:|:----------------:|:-------------------------------:|:----------------------:|:------------:|:-----------------------:|:-----------------:|:---:|:----------:|
|  1 |   -              |   -                             |   -                    |   -          |  `VALUE=1.4`            |  `VALUE=1.3`      || -               |
|  2 |   -              |   -                             |  `VALUE=1.6`           |  `VALUE=1.5` |  `VALUE=1.4`            |   -               ||**`VALUE=1.6`**  |
|  3 |   -              |  `VALUE=1.7`                    |   -                    |  `VALUE=1.5` |  `VALUE=1.4`            |   -               ||**`VALUE=1.7`**  |
|  4 |   -              |   -                             |   -                    |  `VALUE=1.5` |  `VALUE=1.4`            |  `VALUE=1.3`      ||**`VALUE=1.5`**  |
|  5 |`--env VALUE=1.8` |   -                             |   -                    |  `VALUE=1.5` |  `VALUE=1.4`            |  `VALUE=1.3`      ||**`VALUE=1.8`**  |
|  6 |`--env VALUE`     |   -                             |   -                    |  `VALUE=1.5` |  `VALUE=1.4`            |  `VALUE=1.3`      ||**`VALUE=1.4`**  |
|  7 |`--env VALUE`     |   -                             |   -                    |  `VALUE=1.5` |   -                     |  `VALUE=1.3`      ||**`VALUE=1.3`**  |
|  8 |   -              |   -                             |   `VALUE`              |  `VALUE=1.5` |  `VALUE=1.4`            |  `VALUE=1.3`      ||**`VALUE=1.4`**  |
|  9 |   -              |   -                             |   `VALUE`              |  `VALUE=1.5` |   -                     |  `VALUE=1.3`      ||**`VALUE=1.3`**  |
| 10 |   -              |  `VALUE`                        |   -                    |  `VALUE=1.5` |  `VALUE=1.4`            |  `VALUE=1.3`      ||**`VALUE=1.4`**  |
| 11 |   -              |  `VALUE`                        |   -                    |  `VALUE=1.5` |  -                      |  `VALUE=1.3`      ||**`VALUE=1.3`**  |
| 12 |`--env VALUE`     |  `VALUE=1.7`                    |   -                    |  `VALUE=1.5` |  `VALUE=1.4`            |  `VALUE=1.3`      ||**`VALUE=1.4`**  |
| 13 |`--env VALUE=1.8` |  `VALUE=1.7`                    |   -                    |  `VALUE=1.5` |  `VALUE=1.4`            |  `VALUE=1.3`      ||**`VALUE=1.8`**  |
| 14 |`--env VALUE=1.8` |   -                             |  `VALUE=1.6`           |  `VALUE=1.5` |  `VALUE=1.4`            |  `VALUE=1.3`      ||**`VALUE=1.8`**  |
| 15 |`--env VALUE=1.8` |  `VALUE=1.7`                    |  `VALUE=1.6`           |  `VALUE=1.5` |  `VALUE=1.4`            |  `VALUE=1.3`      ||**`VALUE=1.8`**  |

Explanation:

1: local environment defines values for `VALUE` variable, but compose file is not set to replicate this inside container, so no such variable is set

2: `env_file` defines an explicit value for `VALUE`, container environment is set accordingly

3: `environment` defines an explicit value for `VALUE`, container environment is set accordingly

4: Image declares variable `VALUE`, and compose file is not set to override this value, so this variable is set as defined by image

5: `docker compose run` has `--env` flag set which an explicitly value, and overrides value set by image. 

6: `docker compose run` has `--env` flag set to replicate value from environment. Host OS value takes precedence and is replicated into container's environment

7: `docker compose run` has `--env` flag set to replicate value from environment. Value from `.env` file is the selected to define container's environment

8: Compose file `env_file` is set to replicate `VALUE` from local environment. Host OS value takes precedence and is replicated into container's environment

9: Compose file `env_file` is set to replicate `VALUE` from local environment. Value from `.env` file is the selected to define container's environment

10: Compose file `environment` is set to replicate `VALUE` from local environment. Host OS value takes precedence and is replicated into container's environment

11: Compose file `environment` is set to replicate `VALUE` from local environment. Value from `.env` file is the selected to define container's environment

12: `--env` flag has higher precedence on other ways (`environment` and `env_file`) to set environment value. Host OS value takes precedence and is replicated into container's environment

13 to 15: `--env` flag has higher precedence on other ways (`environment` and `env_file`) to set environment value. 
