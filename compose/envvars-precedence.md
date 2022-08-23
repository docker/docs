---
title: Environment variables precedence
description: Scenario Overview illustrating how environmental variables are resolved in Compose
keywords: compose, environment, env file
---

## Order of precedence
1. Passed from the command-line [`docker compose run --env <KEY[=[VAL]]>`](../../engine/reference/commandline/compose_run/#options).
2. Passed from/set in `compose.yaml` service's configuration, from the [environment key](../../compose/compose-file/#environment).
3. Passed from/set in `compose.yaml` service's configuration, from the [env_file key](../../compose/compose-file/#env_file).
4. Passed from/set in Container Image in [ENV directive](../../engine/reference/builder/#env).

### Precedence quick overview table
The table below provides a quick overview of how interpolation works when using all environment variables on Compose, using `TAG`, an environmental variable defining the version for an image, as an example.

Each row represents a scenario and each column represents a context where you can be setting or passing an environment variable.

|  # |  `.env` file      |  `OS` environment  |  Image      |  `compose.yaml`:`env_file` key  |  `compose.yaml`:`environment` key  |  `run --env`  |  Resolved as  |
|:--:|:-----------------:|:------------------:|:-----------:|:-------------------------------:|:----------------------------------:|:-------------:|:-------------:|
|  1 |  `TAG=1.3`        |  `TAG=1.4`         |   -         |   -                             |   -                                |   -           |   -           |
|  2 |  `TAG=1.3`        |  `TAG=1.4`         |**`TAG=1.5`**|   -                             |   -                                |   -           |  `TAG=1.5`    |
|  3 |  `TAG=1.3`        |**`TAG=1.4`**       |  `TAG=1.5`  |   -                             |   -                                |**`TAG`**      |  `TAG=1.4`    |
|  4 |**`TAG=1.3`**      |   -                |  `TAG=1.5`  |**`TAG`**                        |   -                                |   -           |  `TAG=1.3`    |
|  5 |**`TAG=1.3`**      |   -                |  `TAG=1.5`  |   -                             |   -                                |**`TAG`**      |  `TAG=1.3`    |
|  6 |  `TAG=1.3`        |  `TAG=1.4`         |  `TAG=1.5`  |   -                             |   -                                |**`TAG=1.8`**  |  `TAG=1.8`    |
|  7 |  `TAG=1.3`        |**`TAG=1.4`**       |  `TAG=1.5`  |   -                             |**`TAG`**                           |   -           |  `TAG=1.4`    |
|  8 |  `TAG=1.3`        |**`TAG=1.4`**       |  `TAG=1.5`  |   -                             |  `TAG=1.7`                         |**`TAG`**      |  `TAG=1.4`    |
|  9 |  `TAG=1.3`        |  `TAG=1.4`         |  `TAG=1.5`  |   -                             |  `TAG=1.7`                         |**`TAG=1.8`**  |  `TAG=1.8`    |
| 10 |  `TAG=1.3`        |**`TAG=1.4`**       |  `TAG=1.5`  |   -                             |**`TAG`**                           |   -           |  `TAG=1.4`    |
| 11 |  `TAG=1.3`        |  `TAG=1.4`         |  `TAG=1.5`  |  `TAG=1.6`                      |   -                                |**`TAG=1.8`**  |  `TAG=1.8`    |
| 12 |  `TAG=1.3`        |  `TAG=1.4`         |  `TAG=1.5`  |  `TAG=1.6`                      |  `TAG=1.7`                         |**`TAG=1.8`**  |  `TAG=1.8`    |
| 13 |   -               |  `TAG=1.4`         |  `TAG=1.5`  |**`TAG=1.6`**                    |   -                                |   -           |  `TAG=1.6`    |
| 14 |   -               |  `TAG=1.4`         |  `TAG=1.5`  |   -                             |**`TAG=1.7`**                       |   -           |  `TAG=1.7`    |

Description for each column:
* `.env` file - `.env` file on the project root (or, with higher precedence, the file passed via `docker compose -–env-file <FILE>`).
* `OS` environment - OS Environment variable
* Image - `ENV` directive in the Dockerfile
* Compose file - In `env_file` key from the service section in the `compose.yaml`.
* Compose file - In `environment` key from the service section in the `compose.yaml`.
* Command line - Environmental variable passed via `docker compose run -e <KEY[=[VAL]]>`.
* Resolved as - This column expresses the result available in the container.

**Note that the "`OS`" has precedence over "`.env` file" column for variable resolution in the other columns.**