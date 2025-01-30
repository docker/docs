---
title: Using the Docker Desktop CLI
linkTitle: Docker Desktop CLI
weight: 120
description: How to use the Docker Desktop CLI
keywords: cli, docker desktop, macos, windows, linux
params:
  sidebar:
    badge:
      color: green
      text: New
---

{{< summary-bar feature_name="Docker Desktop CLI" >}}

The Docker Desktop CLI lets you perform key operations such as starting, stopping, restarting, and checking the status of Docker Desktop directly from the command line.

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
| `engine use`         | Switch between Linux and Windows containers (Windows only) |
| `update`             | Manage Docker Desktop updates. Available for Mac only and with Docker Desktop version 4.38 and later. |


For more details on each command, see the [Docker Desktop CLI reference](/reference/cli/docker/desktop/_index.md).
