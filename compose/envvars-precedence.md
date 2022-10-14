---
title: Environment variables precedence
description: Scenario Overview illustrating how environment variables are resolved in Compose
keywords: compose, environment, env file
---

This page contains information on how how interpolation works when using environment variables in Compose.

When you set the same environment variable in multiple files, there’s a precedence rule used by Compose. It aims to resolve the value for the variable in question. 

The order of precedence is as follows:
1. Passed from the command line [`docker compose run --env <KEY[=[VAL]]>`](../../engine/reference/commandline/compose_run/#options).
2. Passed from/set in `compose.yaml` service's configuration, from the [environment key](../../compose/compose-file/#environment).
3. Passed from/set in `compose.yaml` service's configuration, from the [env_file key](../../compose/compose-file/#env_file).
4. Passed from/set in Container Image in the [ENV directive](../../engine/reference/builder/#env).

### Example scenario

The following table uses `TAG`, an environment variable defining the version for an image, as an example.

#### The columns

Each column represents a context from where you can set or pass `TAG`, or environment variables in general.

Table column's description:
* `run --env` sets the environment variable using the command line instruction `docker compose run -e <KEY[=[VAL]]>`.
* `compose.yaml`:`environment` key. In `environment` key from the service section in the `compose.yaml` file.
* `compose.yaml`:`env_file` key. In `env_file` key from the service section in the `compose.yaml` file.
* Image `ENV` - `ENV` directive in the Dockerfile for the image.
* `Host OS` environment. The environment where the Docker Engine is running.
* `.env` file. The `.env` file in the project root.
* Resolved as. Gives the final result of how `TAG` is resolved for each scenario according to the precedence rule.

#### The rows

Each row represents a combination of contexts where `TAG` is set or passed simultaneously.



|  # |  `run --env`  |  `compose.yaml`:`environment` key  |  `compose.yaml`:`env_file` key  |  Image `ENV` |  `Host OS` environment  |  `.env` file      |  Resolved as  |
|:--:|:-------------:|:----------------------------------:|:-------------------------------:|:------------:|:-----------------------:|:-----------------:|:-------------:|
|  1 |   -           |   -                                |   -                             |   -          |  `TAG=1.4`              |  `TAG=1.3`        |   -           |
|  2 |   -           |   -                                |   -                             |**`TAG=1.5`** |  `TAG=1.4`              |  `TAG=1.3`        |**`TAG=1.5`**  |
|  3 |**`TAG`**      |   -                                |   -                             |  `TAG=1.5`   |**`TAG=1.4`**            |  `TAG=1.3`        |**`TAG=1.4`**  |
|  4 |   -           |   -                                |**`TAG`**                        |  `TAG=1.5`   |   -                     |**`TAG=1.3`**      |**`TAG=1.3`**  |
|  5 |**`TAG`**      |   -                                |   -                             |  `TAG=1.5`   |   -                     |**`TAG=1.3`**      |**`TAG=1.3`**  |
|  6 |**`TAG=1.8`**  |   -                                |   -                             |  `TAG=1.5`   |  `TAG=1.4`              |  `TAG=1.3`        |**`TAG=1.8`**  |
|  7 |   -           |**`TAG`**                           |   -                             |  `TAG=1.5`   |**`TAG=1.4`**            |  `TAG=1.3`        |**`TAG=1.4`**  |
|  8 |**`TAG`**      |  `TAG=1.7`                         |   -                             |  `TAG=1.5`   |**`TAG=1.4`**            |  `TAG=1.3`        |**`TAG=1.4`**  |
|  9 |**`TAG=1.8`**  |  `TAG=1.7`                         |   -                             |  `TAG=1.5`   |  `TAG=1.4`              |  `TAG=1.3`        |**`TAG=1.8`**  |
| 10 |   -           |**`TAG`**                           |   -                             |  `TAG=1.5`   |**`TAG=1.4`**            |  `TAG=1.3`        |**`TAG=1.4`**  |
| 11 |**`TAG=1.8`**  |   -                                |  `TAG=1.6`                      |  `TAG=1.5`   |  `TAG=1.4`              |  `TAG=1.3`        |**`TAG=1.8`**  |
| 12 |**`TAG=1.8`**  |  `TAG=1.7`                         |  `TAG=1.6`                      |  `TAG=1.5`   |  `TAG=1.4`              |  `TAG=1.3`        |**`TAG=1.8`**  |
| 13 |   -           |   -                                |**`TAG=1.6`**                    |  `TAG=1.5`   |  `TAG=1.4`              |   -               |**`TAG=1.6`**  |
| 14 |   -           |**`TAG=1.7`**                       |   -                             |  `TAG=1.5`   |  `TAG=1.4`              |   -               |**`TAG=1.7`**  |


> **Note**
>
> The columns _`Host OS` environment_ and _`.env` file_ is listed only for lookup. These columns can't result in a variable in the container by itself.