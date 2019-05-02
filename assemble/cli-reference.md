---
title: Docker Assemble CLI reference
description: Docker Assemble CLI reference
keywords: Docker, assemble, Spring Boot, ASP .NET, backend
---

This page provides information about the `docker assemble` command.

## Overview

Docker Assemble (`docker assemble`) is a CLI plugin which provides a language and framework-aware tool that enables users to build an application into an optimized Docker container.

For more information about Docker Assemble, see [Docker Assemble](/assemble/install/).

## `docker assemble` commands

To view the commands and sub-commands available in `docker assemble`, run:

`docker assemble --help`

```
Usage:  docker assemble [OPTIONS] COMMAND

assemble is a high-level build tool

Options:
      --addr string   backend address (default
                      "docker-container://docker-assemble-backend-Usha-Mandya")

Management Commands:
  backend     Manage build backend service

Commands:
  build       Build a project into a container
  version     Print the version number of docker assemble

Run 'docker assemble COMMAND --help' for more information on a command.
```

### backend

The `docker assemble backend` command allows you to manage and build backend services. Docker Assemble requires its own buildkit instance to be running in a Docker container on the local system.

```
Usage:  docker assemble backend [OPTIONS] COMMAND

Manage build backend service

Options:
      --addr string   backend address (default
                      "docker-container://docker-assemble-backend-username")

Management Commands:
  cache       Manage build cache

Commands:
  info        Print information about build backend service
  logs        Show logs for build backend service
  start       Start build backend service
  stop        Stop build backend service

Run 'docker assemble backend COMMAND --help' for more information on a command.
```

For example:

```
docker assemble backend start
Pulling image «…»: Success
Started backend container "docker-assemble-backend-username" (3e627bb365a4)
```

For more information about `backend`, see  [Advanced backend management](/assemble/adv-backend-manage).

### build

The `docker assemble build` command enables you to build a project into a container.

```
Usage:  docker assemble build [PATH]

Build a project into a container

Options:
      --addr string           backend address (default
                              "docker-container://docker-assemble-backend-username")
      --label KEY=VALUE       label to write into the image as KEY=VALUE
      --name NAME             build image with repository NAME (default
                              taken from project metadata)
      --namespace NAMESPACE   build image within repository NAMESPACE
                              (default no namespace)
  -o, --option OPTION=VALUE   set an option as OPTION=VALUE
      --port stringArray      port to expose from container
      --progress string       set type of progress (auto, plain, tty).
                              Use plain to show container output (default
                              "auto")
      --push                  push result to registry, not local image store
      --push-insecure         push result to insecure (http) registry,
                              not local image store
      --tag TAG               tag image with TAG (default taken from
                              project metadata or "latest")
```

For example:

```
~$ docker assemble build docker-springframework
«…»
Successfully built: docker.io/library/hello-boot:1
```

## version

The `docker assemble version` command displays the version number of Docker Assemble.

```
Usage:  docker assemble version

Print the version number of docker assemble

Options:
      --addr string   backend address (default
                      "docker-container://docker-assemble-backend-username")
```

For example:

```
> docker assemble version
docker assemble v0.31.0
commit: d089e2be00b0f7d7f565aeba11cb8bc6dd56a40b
buildkit: 2bd8e6cb2b42
os/arch: windows/amd64
```
