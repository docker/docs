---
title: Using the Docker Desktop CLI
linkTitle: Docker Desktop CLI
weight: 1000
description: How to use the Docker Desktop CLI
keywords: cli, docker desktop, macos, windows, linux
params:
  sidebar:
    badge:
      color: green
      text: New
---

{{% experimental title="Beta" %}}
Docker Desktop CLI is currently in [Beta](../../release-lifecycle.md#beta).
{{% /experimental %}}

The Docker Desktop CLI lets you perform key operations such as starting, stopping, restarting, and checking the status of Docker Desktop directly from the command line. It is available with Docker Desktop version 4.37 and later.

The Docker Desktop CLI provides:

- Enhanced automation and CI/CD integration: Perform Docker Desktop operations directly in CI/CD pipelines for better workflow automation.
- An improved developer experience: Restart, quit, or reset Docker Desktop from the command line, reducing dependency on the Docker Desktop Dashboard and improving flexibility and efficiency.

## Usage

```console
docker desktop COMMAND [OPTIONS]
```

## Commands

| Command              | Description                              |
|:---------------------|:-----------------------------------------|
| `start`              | Starts Docker Desktop                    |
| `stop`               | Stops Docker Desktop                     |
| `restart`            | Restarts Docker Desktop                  |
| `status`             | Displays whether Docker Desktop is running or stopped.       |
| `engine ls`          | Lists available engines (Windows only)   |
| `engine use linux`   | Switches to Linux containers (Windows only) |
| `engine use windows` | Switches to Windows containers (Windows only)                 |

## `docker desktop start`

### Usage

`docker desktop start [OPTIONS]`

The terminal displays a spinner until Docker Desktop starts.

### Options

| Option     | Description                         |
|:-----------|:------------------------------------|
| `--detach` | Runs the command in the background and immediately returns terminal control.     |
| `--timeout <seconds>` | Specify how long to wait for Docker Desktop to start before timing out.    |

## `docker desktop stop`

### Usage

`docker desktop stop [OPTIONS]`

The terminal displays a spinner until Docker Desktop is stopped. If Docker Desktop is already stopped, the terminal notifies you.

### Options

| Option     | Description                         |
|:-----------|:------------------------------------|
| `--force` | Immediately terminate all Docker Desktop processes. |

## `docker desktop restart`

### Usage

`docker desktop restart`

Stops and starts Docker Desktop in one step.

## `docker desktop status`

### Usage

`docker desktop status`

Displays whether Docker Desktop is running or stopped.

## `docker desktop engine ls` (Windows only)

### Usage

`docker desktop engine ls`

Lists available engines.

## `docker desktop engine use linux` (Windows only)

### Usage

`docker desktop engine use linux`

Switches to Linux containers.

## `docker desktop engine use windows` (Windows only)

### Usage

`docker desktop engine use windows`

Switches to Windows containers.

## Help

All commands accept the `--help` flag, which documents each command's usage.
