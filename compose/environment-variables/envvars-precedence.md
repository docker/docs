---
title: Environment variables precedence
description: Scenario Overview illustrating how environment variables are resolved in Compose
keywords: compose, environment, env file
redirect_from:
- /compose/envvars-precedence/
---

When you set the same environment variable in multiple files, there’s a precedence rule used by Compose. It aims to resolve the value for the variable in question.

This page contains information on the level of precedence each method of setting environmental variables takes.

The order of precedence is as follows:
1. Set using [`docker compose run -e` in the CLI](set-environment-variables.md#set-environment-variables-with-docker-compose-run--e)
2. Substituted from your [shell](set-environment-variables.md#substitute-from-the-shell)
3. Set using the [`environment` attribute in the Compose file](set-environment-variables.md#use-the-environment-attribute)
4. Use of the [`--env-file` argument](set-environment-variables.md#substitute-with---env-file) in the CLI
5. Use of the [`env_file` attribute](set-environment-variables.md#use-the-env_file-attribute) in the Compose file
6. Set using an [`.env` file](set-environment-variables.md#substitute-with-an-env-file) placed at base of your project directory
7. Set in a container image in the [ENV directive](../../engine/reference/builder.md#env).
   Having any `ARG` or `ENV` setting in a `Dockerfile` evaluates only if there is no Docker Compose entry for `environment`, `env_file` or `run --env`.

## Simple example

In the example below, we set a different value for the same environment variable in an `.env` file and with the `environment` attribute in the Compose file:

```console
$ cat ./Docker/api/api.env
NODE_ENV=test

$ cat docker-compose.yml
services:
  api:
    image: 'node:6-alpine'
    env_file:
     - ./Docker/api/api.env
    environment:
     - NODE_ENV=production
```

The environment variable defined with the `environment` attribute takes precedence.

```console
$ docker compose exec api node

> process.env.NODE_ENV
'production'
```

> Specifics for NodeJS containers
>
> If you have a `package.json` entry for `script:start` like
> `NODE_ENV=test node server.js`, then this overrules any setting in your
> `docker-compose.yml` file.


## Advanced example 

The following table uses `TAG`, an environment variable defining the version for an image, as an example.

### How the table works

Each column represents a context from where you can set a value, or substitute in a value for `TAG`.

The columns `Host OS environment` and `.env file` is listed only as an illustration lookup. In reality, they don't result in a variable in the container by itself.

Each row represents a combination of contexts where `TAG` is set, substituted, or both.


|  # |  `docker compose run -e`  |  `environment` attribute  |  `env_file` attribute  |  Image `ENV` |  `Host OS` environment  |  `.env` file      | |  Result  |
|:--:|:-------------:|:----------------------------------:|:-------------------------------:|:------------:|:-----------------------:|:-----------------:|:---:|:-------------:|
|  1 |   -           |   -                                |   -                             |   -          |  `TAG=1.4`              |  `TAG=1.3`        || - |
|  2 |   -           |   -                                |   -                             |`TAG=1.5` |  `TAG=1.4`                  |  `TAG=1.3`        ||**`TAG=1.5`**  |
|  3 |`TAG`          |   -                                |   -                             |  `TAG=1.5`   |`TAG=1.4`                |  `TAG=1.3`        ||**`TAG=1.4`**  |
|  4 |   -           |   -                                |`TAG`                            |  `TAG=1.5`   |   -                     |`TAG=1.3`          ||**`TAG=1.3`**  |
|  5 |`TAG`          |   -                                |   -                             |  `TAG=1.5`   |   -                     |`TAG=1.3`          ||**`TAG=1.3`**  |
|  6 |`TAG=1.8`      |   -                                |   -                             |  `TAG=1.5`   |  `TAG=1.4`              |  `TAG=1.3`        ||**`TAG=1.8`**  |
|  7 |   -           |`TAG`                               |   -                             |  `TAG=1.5`   |`TAG=1.4`                |  `TAG=1.3`        ||**`TAG=1.4`**  |
|  8 |`TAG`          |  `TAG=1.7`                         |   -                             |  `TAG=1.5`   |`TAG=1.4`                |  `TAG=1.3`        ||**`TAG=1.4`**  |
|  9 |`TAG=1.8`      |  `TAG=1.7`                         |   -                             |  `TAG=1.5`   |  `TAG=1.4`              |  `TAG=1.3`        ||**`TAG=1.8`**  |
| 10 |   -           |`TAG`                               |   -                             |  `TAG=1.5`   |`TAG=1.4`                |  `TAG=1.3`        ||**`TAG=1.4`**  |
| 11 |`TAG=1.8`      |   -                                |  `TAG=1.6`                      |  `TAG=1.5`   |  `TAG=1.4`              |  `TAG=1.3`        ||**`TAG=1.8`**  |
| 12 |`TAG=1.8`      |  `TAG=1.7`                         |  `TAG=1.6`                      |  `TAG=1.5`   |  `TAG=1.4`              |  `TAG=1.3`        ||**`TAG=1.8`**  |
| 13 |   -           |   -                                |`TAG=1.6`                        |  `TAG=1.5`   |  `TAG=1.4`              |   -               ||**`TAG=1.6`**  |
| 14 |   -           |`TAG=1.7`                           |   -                             |  `TAG=1.5`   |  `TAG=1.4`              |   -               ||**`TAG=1.7`**  |

