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

Now that you have an application, you can create the necessary Docker assets to containerize it. You can choose between using the official .NET images or Docker Hardened Images (DHI).

> [Docker Hardened Images (DHIs)](https://docs.docker.com/dhi/) are minimal, secure, and production-ready container base and application images maintained by Docker. DHI images are recommended for better security—they are designed to reduce vulnerabilities and simplify compliance.

{{< tabs >}}
{{< tab name="Using Docker Hardened Images" >}}

Docker Hardened Images (DHIs) for .NET are available in the [Docker Hardened Images catalog](https://hub.docker.com/hardened-images/catalog/dhi/aspnetcore). Docker Hardened Images are freely available to everyone with no subscription required. You can pull and use them like any other Docker image after signing in to the DHI registry. For more information, see the [DHI quickstart](/dhi/get-started/) guide.

1. Sign in to the DHI registry:
   ```console
   $ docker login dhi.io
   ```

2. Pull the .NET SDK DHI (check the catalog for available versions):
   ```console
   $ docker pull dhi.io/dotnet:10-sdk
   ```

3. Pull the ASP.NET Core runtime DHI (check the catalog for available versions):
   ```console
   $ docker pull dhi.io/aspnetcore:10
   ```

You can use `docker init` to generate Docker assets, then modify the Dockerfile to use DHI images:

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
? What version of .NET do you want to use? 10.0
? What local port do you want to use to access your server? 8080
```

In the following Dockerfile, the `FROM` instructions use `dhi.io/dotnet:10-sdk` and `dhi.io/aspnetcore:10` as the base images.

```dockerfile {title=Dockerfile}
# syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM dhi.io/dotnet:10-sdk AS build
ARG TARGETARCH
COPY . /source
WORKDIR /source/src
RUN --mount=type=cache,id=nuget,target=/root/.nuget/packages \
    dotnet publish -a ${TARGETARCH/amd64/x64} --use-current-runtime --self-contained false -o /app

FROM dhi.io/aspnetcore:10
WORKDIR /app
COPY --from=build /app .
ENTRYPOINT ["dotnet", "myWebApp.dll"]
```

> [!NOTE]
>
> DHI runtime images already run as a non-root user (`nonroot`, UID 65532), so there's no need to create a user or specify `USER` in your Dockerfile. This reduces the attack surface and simplifies your configuration.

{{< /tab >}}
{{< tab name="Using the official .NET 10 image" >}}

You can use `docker init` to create the necessary Docker assets. Inside the `docker-dotnet-sample` directory, run the `docker init` command in a terminal. `docker init` provides some default configuration, but you'll need to answer a few questions about your application. Refer to the following example to answer the prompts from `docker init` and use the same answers for your prompts.

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
? What version of .NET do you want to use? 10.0
? What local port do you want to use to access your server? 8080
```

This generates a Dockerfile using the official .NET 10 images from Microsoft Container Registry:

```dockerfile {title=Dockerfile}
# syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM mcr.microsoft.com/dotnet/sdk:10.0-alpine AS build
ARG TARGETARCH
COPY . /source
WORKDIR /source/src
RUN --mount=type=cache,id=nuget,target=/root/.nuget/packages \
    dotnet publish -a ${TARGETARCH/amd64/x64} --use-current-runtime --self-contained false -o /app

FROM mcr.microsoft.com/dotnet/aspnet:10.0-alpine AS final
WORKDIR /app
COPY --from=build /app .
ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser
USER appuser
ENTRYPOINT ["dotnet", "myWebApp.dll"]
```

{{< /tab >}}
{{< /tabs >}}

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

To learn more about the files, see the following:
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
reference](/reference/cli/docker/compose/).

## Summary

In this section, you learned how you can containerize and run your .NET
application using Docker.

Related information:
 - [Dockerfile reference](/reference/dockerfile.md)
 - [.dockerignore file reference](/reference/dockerfile.md#dockerignore-file)
 - [Docker Compose overview](/manuals/compose/_index.md)
 - [Docker Hardened Images](/dhi/)

## Next steps

In the next section, you'll learn how you can develop your application using
Docker containers.
