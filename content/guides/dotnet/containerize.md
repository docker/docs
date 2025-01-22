---
title: Containerize a .NET application
linkTitle: Containerize your app
weight: 10
keywords: .net, containerize, initialize
description: Learn how to containerize an ASP.NET application.
aliases:
- /language/dotnet/build-images/
- /language/dotnet/run-containers/
- /language/dotnet/containerize/
- /guides/language/dotnet/containerize/
---

## Prerequisites

* You have installed the latest version of [Docker
  Desktop](/get-started/get-docker.md).
* You have a [git client](https://git-scm.com/downloads). The examples in this
  section use a command-line based git client, but you can use any client.

## Overview

This section walks you through containerizing and running a .NET
application.

## Get the sample applications

In this guide, you will use a pre-built .NET application. The application is
similar to the application built in the Docker Blog article, [Building a
Multi-Container .NET App Using Docker
Desktop](https://www.docker.com/blog/building-multi-container-net-app-using-docker-desktop/).

Open a terminal, change directory to a directory that you want to work in, and
run the following command to clone the repository.

```console
$ git clone https://github.com/docker/docker-dotnet-sample
```

## Initialize Docker assets

Now that you have an application, you can use `docker init` to create the
necessary Docker assets to containerize your application. Inside the
`docker-dotnet-sample` directory, run the `docker init` command in a terminal.
`docker init` provides some default configuration, but you'll need to answer a
few questions about your application. Refer to the following example to answer
the prompts from `docker init` and use the same answers for your prompts.

```console
$ docker init
Welcome to the Docker Init CLI!

This utility will walk you through creating the following files with sensible defaults for your project:
  - .dockerignore
  - Dockerfile
  - compose.yaml
  - README.Docker.md

Let's get started!

? What application platform does your project use? ASP.NET Core
? What's the name of your solution's main project? myWebApp
? What version of .NET do you want to use? 8.0
? What local port do you want to use to access your server? 8080
```

You should now have the following contents in your `docker-dotnet-sample`
directory.

```text
├── docker-dotnet-sample/
│ ├── .git/
│ ├── src/
│ ├── .dockerignore
│ ├── compose.yaml
│ ├── Dockerfile
│ ├── README.Docker.md
│ └── README.md
```

To learn more about the files that `docker init` added, see the following:
 - [Dockerfile](/reference/dockerfile.md)
 - [.dockerignore](/reference/dockerfile.md#dockerignore-file)
 - [compose.yaml](/reference/compose-file/_index.md)

## Run the application

Inside the `docker-dotnet-sample` directory, run the following command in a
terminal.

```console
$ docker compose up --build
```

Open a browser and view the application at [http://localhost:8080](http://localhost:8080). You should see a simple web application.

In the terminal, press `ctrl`+`c` to stop the application.

### Run the application in the background

You can run the application detached from the terminal by adding the `-d`
option. Inside the `docker-dotnet-sample` directory, run the following command
in a terminal.

```console
$ docker compose up --build -d
```

Open a browser and view the application at [http://localhost:8080](http://localhost:8080). You should see a simple web application.

In the terminal, run the following command to stop the application.

```console
$ docker compose down
```

For more information about Compose commands, see the [Compose CLI
reference](/reference/cli/docker/compose/_index.md).

## Summary

In this section, you learned how you can containerize and run your .NET
application using Docker.

Related information:
 - [Dockerfile reference](/reference/dockerfile.md)
 - [.dockerignore file reference](/reference/dockerfile.md#dockerignore-file)
 - [Docker Compose overview](/manuals/compose/_index.md)

## Next steps

In the next section, you'll learn how you can develop your application using
Docker containers.
