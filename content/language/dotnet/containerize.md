---
title: Containerize a .NET application
keywords: .net, containerize, initialize
description: Learn how to containerize an ASP.NET application.
aliases:
- /language/dotnet/build-images/
- /language/dotnet/run-containers/
---

## Prerequisites

* You have installed the latest version of [Docker
  Desktop](../../get-docker.md).
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

## Test the simple application without Docker (optional)

You can test the sample application locally without Docker before you continue
building and running the application with Docker. This section requires you to
have the .NET SDK version 6.0 installed on your machine. Download and install
[.NET 6.0 SDK](https://dotnet.microsoft.com/download).

Open a terminal, change directory to the `docker-dotnet-sample/src` directory,
and run the following command.

```console
$ dotnet run --urls http://localhost:5000
```

Open a browser and view the application at [http://localhost:5000](http://localhost:5000). You should see a simple web application.

In the terminal, press `ctrl`+`c` to stop the application.

## Initialize Docker assets

Now that you have an application, you can use `docker init` to create the
necessary Docker assets to containerize your application. Inside the
`docker-dotnet-sample` directory, run the `docker init` command in a terminal.
Refer to the following example to answer the prompts from `docker init`.

```console
$ docker init
Welcome to the Docker Init CLI!

This utility will walk you through creating the following files with sensible defaults for your project:
  - .dockerignore
  - Dockerfile
  - compose.yaml
  - README.Docker.md

Let's get started!

? What application platform does your project use? ASP.NET
? What's the name of your solution's main project? myWebApp
? What version of .NET do you want to use? 6.0
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
 - [Dockerfile](../../engine/reference/builder.md)
 - [.dockerignore](../../engine/reference/builder.md#dockerignore-file)
 - [compose.yaml](../../compose/compose-file/_index.md)

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
reference](../../compose/reference/_index.md).

## Summary

In this section, you learned how you can containerize and run your .NET
application using Docker.

Related information:
 - [Dockerfile reference](../../engine/reference/builder.md)
 - [Build with Docker guide](../../build/guide/index.md)
 - [.dockerignore file reference](../../engine/reference/builder.md#dockerignore-file)
 - [Docker Compose overview](../../compose/_index.md)

## Next steps

In the next section, you'll learn how you can develop your application using
Docker containers.

{{< button text="Develop your application" url="develop.md" >}}
