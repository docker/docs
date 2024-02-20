---
title: Containerize a Java application
keywords: java, containerize, initialize, maven, build
description: Learn how to containerize a Java application.
aliases:
 - /language/java/build-images/
 - /language/java/run-containers/
---

## Prerequisites

- You have installed the latest version of [Docker Desktop](../../../get-docker.md).
  Docker adds new features regularly and some parts of this guide may
  work only with the latest version of Docker Desktop.
* You have a [Git client](https://git-scm.com/downloads). The examples in this
  section use a command-line based Git client, but you can use any client.

## Overview

This section walks you through containerizing and running a Java
application.

## Get the sample applications

Clone the sample application that you'll be using to your local development machine. Run the following command in a terminal to clone the repository.

```console
$ git clone https://github.com/spring-projects/spring-petclinic.git
```

The sample application is a Spring Boot application built using Maven. For more details, see `readme.md` in the repository.

## Initialize Docker assets

Now that you have an application, you can use `docker init` to create the
necessary Docker assets to containerize your application. Inside the
`spring-petclinic` directory, run the `docker init` command in a terminal.
`docker init` provides some default configuration, but you'll need to answer a
few questions about your application. Use the answers in the following example in order to follow along with this guide.

The sample application already contains Docker assets. You'll be prompted to overwrite the existing Docker assets. To continue with this guide, select `y` to overwrite them.

```console
$ docker init
Welcome to the Docker Init CLI!

This utility will walk you through creating the following files with sensible defaults for your project:
  - .dockerignore
  - Dockerfile
  - compose.yaml
  - README.Docker.md

Let's get started!

WARNING: The following Docker files already exist in this directory:
  - docker-compose.yml
? Do you want to overwrite them? Yes
? What application platform does your project use? Java
? What's the relative directory (with a leading .) for your app? ./src
? What version of Java do you want to use? 17
? What port does your server listen on? 8080
```

In the previous example, notice the `WARNING`. `docker-compose.yaml` already
exists, so `docker init` overwrites that file rather than creating a new
`compose.yaml` file. This prevents having multiple Compose files in the
directory. Both names are supported, but Compose prefers the canonical
`compose.yaml`.

You should now have the following three new files in your `spring-petclinic`
directory.

- [Dockerfile](../../engine/reference/builder.md)
- [.dockerignore](../../engine/reference/builder.md#dockerignore-file)
- [docker-compose.yaml](../../compose/compose-file/_index.md)


You can open the files in a code or text editor, then read the comments to learn
more about the instructions, or visit the links in the previous list.

## Run the application

Inside the `spring-petclinic` directory, run the following command in a
terminal.

```console
$ docker compose up --build
```

The first time you build and run the app, Docker downloads dependencies and builds the app. It may take several minutes depending on your network connection.

Open a browser and view the application at [http://localhost:8080](http://localhost:8080). You should see a simple app for a pet clinic.

In the terminal, press `ctrl`+`c` to stop the application.

### Run the application in the background

You can run the application detached from the terminal by adding the `-d`
option. Inside the `docker-php-sample` directory, run the following command
in a terminal.

```console
$ docker compose up --build -d
```

Open a browser and view the application at [http://localhost:8080](http://localhost:8080). You should see a simple app for a pet clinic.

In the terminal, run the following command to stop the application.

```console
$ docker compose down
```

For more information about Compose commands, see the
[Compose CLI reference](../../compose/reference/_index.md).

## Summary

In this section, you learned how you can containerize and run a Java
application using Docker.

Related information:
 - [docker init reference](../../engine/reference/commandline/init.md)

## Next steps

In the next section, you'll learn how you can develop your application using
Docker containers.

{{< button text="Develop your application" url="develop.md" >}}