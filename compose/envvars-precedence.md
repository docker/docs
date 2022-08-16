---
title: Environment variables precedence
description: Scenario Overview illustrating how environmental variables are resolved in Compose
keywords: compose, environment, env file
---

## Order of precedence:
1. Passed from the command-line [`docker compose run --env <KEY[=VAL]>`](../../engine/reference/commandline/compose_run/#options).
2. Passed from/set in `compose.yaml` service's configuration, [from the environment key](../../compose/compose-file/#environment).
3. Passed from/set in `compose.yaml` service's configuration, [from the env_file key](../../compose/compose-file/#env_file).
4. Passed from/set in Container Image [in ENV directive](../../engine/reference/builder/#env).

Using as example `TAG`, a environmental variable defining the version for an image.

Key to table, description for each column:
* Image - `ENV` directive in the Dockerfile
* .env file - `.env` file on the project root (or, with higher precedence, the file passed via `docker compose -–env-file <FILE>`).
* Command line - environmental variable passed via `docker compose run -e `<KEY[=VAL]>`.
* Compose file - in `environment` key from the service section in the `compose.yaml`.
* Compose file - in `env_file key` from the service section in the `compose.yaml`.

<!--
### Passed or Set
explanation here
 -->

The table below aims to provide a quick overview of how interpolation works when using all environment variables on Compose.
Each row represents a scenario and each columns represents a context where you can be setting or passing an environment variable.

| # | `.env` file                | `compose.yaml`:`env_file` key | `compose.yaml`:`environment` key  | CMD          |    Image      |  OS          |    Resolved as    |
|:-:|:--------------------------:|:-----------------------------:|:---------------------------------:|:------------:|:-------------:|:------------:|:-----------------:|
| 1 | `TAG=1.3`                  |  Unset                        |    Unset                          |    -         |     `TAG=1.6`   |   `TAG=1.4`    | `TAG=1.6`           |
| 2 | `TAG=1.3`                  |  Unset                        |    Unset                          |    `TAG`      |     `TAG=1.6`   |   `TAG=1.4`    | `TAG=1.3`           |
| 3 | `TAG=$TAG:-1.2`            |  Unset                        |    Unset                          |    `TAG`      |     `TAG=1.6`   |   `TAG=1.4`    | `TAG=1.4`           |
| 4 | `TAG=$TAG:-1.2`            |  Unset                        |    Unset                          |    -         |     `TAG=1.6`   |   `TAG=1.4`    | `TAG=1.6`           |
| 5 | `TAG=$TAG:-1.2`            |  Unset                        |    Unset                          |    `TAG`      |     `TAG=1.6`   |      -       | `TAG=1.6`           |
| 6 | `TAG=$TAG:-1.2`            |  Unset                        |    Unset                          |    `TAG=1.5`   |     `TAG=1.6`   |   `TAG=1.4`    | `TAG=1.5`           |
| 7 | `TAG=$TAG:-1.2`            |  Unset                        |    `TAG`                           |    -         |     `TAG=1.6`   |   `TAG=1.4`    | `TAG=1.4`           |
| 8 | `TAG=$TAG:-1.2`            |  Unset                        |    `TAG=1.7`                        |    `TAG`      |     `TAG=1.6`   |   `TAG=1.4`    | `TAG=1.7`           |
| 9 | `TAG=$TAG:-1.2`            |  Unset                        |    `TAG=1.7`                        |    `TAG=1.5`   |     `TAG=1.6`   |   `TAG=1.4`    | `TAG=1.5`           |
| 10| `TAG=$TAG:-1.2`            |  Unset                        |    `TAG`                           |    `TAG`      |     `TAG=1.6`   |   `TAG=1.4`    | `TAG=1.4`           |
| 11| `TAG=$TAG:-1.2`            |  Unset                        |    `TAG`                           |    `TAG`      |     `TAG=1.6`   |   `TAG=1.4`    | `TAG=1.4`           |
| 12| `TAG=$TAG:-1.2`            |  `TAG=1.8`                    |    -                              |  **`TAG=1.5`** |     `TAG=1.6`   |   `TAG=1.4`    | `TAG=1.5`           |
| 13| `TAG=$TAG:-1.2`            |  `TAG=1.8`                    |  **`TAG=1.7`**                      |    `TAG=1.5`   |     `TAG=1.6`   |   `TAG=1.4`    | `TAG=1.7`           |
| 14| -                          |**`TAG=1.8`**                  |     -                             |     -        |     `TAG=1.6`   |   `TAG=1.4`    | TAG=1.8           |
| 15| -                          |  -                            |     `TAG=1.7`                       |     -        |     `TAG=1.6`   |   `TAG=1.4`    | `TAG=1.4`           |
| 16| -                          |  -                            |     `TAG=1.7`                       |     -        |     `TAG=1.6`   |   `TAG=1.4`    | `TAG=1.4`           |
 