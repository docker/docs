---
description: Introduction and Overview of Compose
keywords: documentation, docs, docker, development, environments, containers
title: Overview
redirect_from:
- /devenvironments/
---

Docker Development Environments are a set of tools that help you build, test, and share your applications. They are designed to work with any language, framework, or runtime. Docker Development Environments are built on top of Docker Desktop and the Docker Engine.

It enables you to:
- Build your application in a container
- Test your application in a container
- Debug your application in a container
- Run your application to a container

This simplifies the development process by providing a consistent environment for your application.

## Getting started

To get started, you need to install Docker Desktop. Docker Desktop is available for Windows, Mac, and Linux. You can download [Docker Desktop](https://www.docker.com/products/docker-desktop/).

Once you have installed Docker Desktop, you can start using Docker Development Environments.

## Clone a sample application

To get started, you can choose one of the [Awesome Compose](https://github.com/docker/awesome-compose) projects.
For example, you can use the [Flask sample application](https://github.com/docker/awesome-compose/tree/master/flask)

```shell
$ docker dev create https://github.com/docker/awesome-compose/tree/master/flask
```

This command will clone the sample application in your development environment container.
