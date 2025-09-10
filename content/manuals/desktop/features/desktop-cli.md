---
title: Use the Docker Desktop CLI
linkTitle: Docker Desktop CLI
weight: 100
description: How to use the Docker Desktop CLI
keywords: cli, docker desktop, macos, windows, linux
---

{{< summary-bar feature_name="Docker Desktop CLI" >}}

The Docker Desktop CLI lets you perform key operations such as starting, stopping, restarting, and updating Docker Desktop directly from the command line.

The Docker Desktop CLI provides:

- Simplified automation for local development: Execute Docker Desktop operations more efficiently in scripts and tests. 
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
| `update`             | Manage Docker Desktop updates. Available for Mac only with Docker Desktop version 4.38, or all OSs with Docker Desktop version 4.39 and later. |
| `logs`               | Print log entries                        |
| `disable`            | Disable a feature                        |
| `enable`             | Enable a feature                         | 
| `version`            | Show the Docker Desktop CLI plugin version information |
| `kubernetes`         | List Kubernetes images used by Docker Desktop or restart the cluster. Available with Docker Desktop version 4.44 and later.          |

For more details on each command, see the [Docker Desktop CLI reference](/reference/cli/docker/desktop/_index.md).
