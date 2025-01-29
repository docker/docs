---
description: Set up a dev Environments
keywords: Dev Environments, share, docker dev, Docker Desktop
title: Use the docker dev CLI plugin
aliases:
- /desktop/dev-environments/dev-cli/
---

{{% include "dev-envs-changing.md" %}}

Use the new `docker dev` CLI plugin to get the full Dev Environments experience from the terminal in addition to the Dashboard.

It is available with [Docker Desktop 4.13.0 and later](/manuals/desktop/release-notes.md). 

### Usage

```bash
docker dev [OPTIONS] COMMAND
```

### Commands

| Command              | Description                              |
|:---------------------|:-----------------------------------------|
| `check`              | Check Dev Environments                   |
| `create`             | Create a new dev environment             |
| `list`               | Lists all dev environments               |
| `logs`               | Traces logs from a dev environment       |
| `open`               | Open Dev Environment with the IDE        |
| `rm`                 | Removes a dev environment                |
| `start`              | Starts a dev environment                 |
| `stop`               | Stops a dev environment                  |
| `version`            | Shows the Docker Dev version information |

### `docker dev check`

#### Usage

`docker dev check [OPTIONS]`

#### Options

| Name, shorthand      | Description                         |
|:---------------------|:------------------------------------|
| `--format`,`-f`      | Format the output.                  |

### `docker dev create`

#### Usage

`docker dev create [OPTIONS] REPOSITORY_URL`

#### Options

| Name, shorthand      | Description                                               |
|:---------------------|:----------------------------------------------------------|
| `--detach`,`-d`      | Detach creates a Dev Env without attaching to it's logs.  |
| `--open`,`-o`        | Open IDE after a successful creation                      |

### `docker dev list`

#### Usage

`docker dev list [OPTIONS]`

#### Options

| Name, shorthand      | Description                   |
|:---------------------|:------------------------------|
| `--format`,`-f`      | Format the output             |
| `--quiet`,`-q`       | Only show dev environments names  |

### `docker dev logs`

#### Usage

`docker dev logs [OPTIONS] DEV_ENV_NAME`

### `docker dev open`

#### Usage

`docker dev open DEV_ENV_NAME CONTAINER_REF [OPTIONS]`

#### Options

| Name, shorthand      | Description           |
|:---------------------|:----------------------|
| `--editor`,`-e`      | Editor.               |

### `docker dev rm`

#### Usage

`docker dev rm DEV_ENV_NAME`

### `docker dev start`

#### Usage

`docker dev start DEV_ENV_NAME`

### `docker dev stop`

#### Usage

`docker dev stop DEV_ENV_NAME`

### `docker dev version`

#### Usage

`docker dev version [OPTIONS]`

#### Options

| Name, shorthand      | Description                                   |
|:---------------------|:----------------------------------------------|
| `--format`,`-f`      | Format the output.                            |
| `--short`,`-s`       | Shows only Docker Dev's version number.       |
