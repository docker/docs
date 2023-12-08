---
title: Containerize a PHP application
keywords: php, containerize, initialize, apache, composer
description: Learn how to containerize a PHP application.
---

## Prerequisites

* You have installed the latest version of [Docker
  Desktop](../../get-docker.md).
* You have a [git client](https://git-scm.com/downloads). The examples in this
  section use a command-line based git client, but you can use any client.

## Overview

This section walks you through containerizing and running a PHP
application.

## Get the sample applications

In this guide, you will use a pre-built PHP application. The application uses Composer for library dependency management. You'll serve the application via an Apache web server.

Open a terminal, change directory to a directory that you want to work in, and
run the following command to clone the repository.

```console
$ git clone https://github.com/docker/docker-php-sample
```

The sample application is a basic hello world application and an application that increments a counter in a database. In addition, the application uses PHPUnit for testing.

## Initialize Docker assets

Now that you have an application, you can use `docker init` to create the
necessary Docker assets to containerize your application. Inside the
`docker-php-sample` directory, run the `docker init` command in a terminal.
`docker init` provides some default configuration, but you'll need to answer a
few questions about your application. For example, this application uses PHP
version 8.2. Refer to the following `docker init` example and use the same
answers for your prompts.

```console
$ docker init
Welcome to the Docker Init CLI!

This utility will walk you through creating the following files with sensible defaults for your project:
  - .dockerignore
  - Dockerfile
  - compose.yaml
  - README.Docker.md

Let's get started!

? What application platform does your project use? PHP with Apache
? What version of PHP do you want to use? 8.2
? What's the relative directory (with a leading .) for your app? ./src
? What local port do you want to use to access your server? 9000
```

You should now have the following contents in your `docker-php-sample`
directory.

```text
├── docker-php-sample/
│ ├── .git/
│ ├── src/
│ ├── tests/
│ ├── .dockerignore
│ ├── .gitignore
│ ├── compose.yaml
│ ├── composer.json
│ ├── composer.lock
│ ├── Dockerfile
│ ├── README.Docker.md
│ └── README.md
```

To learn more about the files that `docker init` added, see the following:
 - [Dockerfile](../../engine/reference/builder.md)
 - [.dockerignore](../../engine/reference/builder.md#dockerignore-file)
 - [compose.yaml](../../compose/compose-file/_index.md)

## Run the application

Inside the `docker-php-sample` directory, run the following command in a
terminal.

```console
$ docker compose up --build
```

Open a browser and view the application at [http://localhost:9000/hello.php](http://localhost:9000/hello.php). You should see a simple hello world application.

In the terminal, press `ctrl`+`c` to stop the application.

### Run the application in the background

You can run the application detached from the terminal by adding the `-d`
option. Inside the `docker-php-sample` directory, run the following command
in a terminal.

```console
$ docker compose up --build -d
```

Open a browser and view the application at [http://localhost:9000/hello.php](http://localhost:9000/hello.php). You should see a simple hello world application.

In the terminal, run the following command to stop the application.

```console
$ docker compose down
```

For more information about Compose commands, see the [Compose CLI
reference](../../compose/reference/_index.md).

## Summary

In this section, you learned how you can containerize and run a simple PHP
application using Docker.

Related information:
 - [docker init reference](../../engine/reference/commandline/init.md)

## Next steps

In the next section, you'll learn how you can develop your application using
Docker containers.

{{< button text="Develop your application" url="develop.md" >}}