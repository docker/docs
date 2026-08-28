---
title: Build and share a containerized application
linkTitle: Containerize an application
description: Run a container and an application stack, build an image, and share it through Docker Hub.
keywords: Docker, get started, containerize application, containers, Docker Compose, images, Docker Hub
weight: 1
aliases:
  - /get-started/run-an-app/
---

An application can depend on a particular runtime, database, and supporting
tools. Recreating that environment on every machine takes time and produces
different results.

In this 15-minute tutorial, you'll run a container, start a multi-container
application, package the application as an image, and share the image through
Docker Hub.

## Before you start

- [Install Docker Desktop](../get-docker.md) and start it
- Install [Git](https://git-scm.com/downloads)
- Create a [Docker account](https://app.docker.com/signup)

## Run a container

Start with an application that is already packaged as an image:

```console
$ docker run -d --name welcome -p 8080:80 docker/welcome-to-docker
```

Open [http://localhost:8080](http://localhost:8080). Docker pulled the image and
started a container from it. The `-p` flag made the container's web server
available on your host.

Remove the container before continuing:

```console
$ docker rm -f welcome
```

The image remains available locally, but the running application is gone.

## Run an application stack

Clone a prepared application and open its directory:

```console
$ git clone https://github.com/docker/getting-started-todo-app
$ cd getting-started-todo-app
```

Start the development stack:

```console
$ docker compose up --build -d
```

Open [http://localhost](http://localhost), then add an item to the to-do list.

The project's `compose.yaml` file describes the frontend, API, database,
database management interface, and proxy. Compose built the application
services, pulled the supporting images, and connected the containers.

Check the stack:

```console
$ docker compose ps
```

Each row represents one part of the application. The stack can move between
machines as one version-controlled definition.

## Build an image

The repository also contains a `Dockerfile` that packages the frontend and API
as one production image. Build it and replace `<YOUR_DOCKER_USERNAME>` with
your Docker username:

```console
$ docker build -t <YOUR_DOCKER_USERNAME>/getting-started-todo-app .
```

Run the image as a new container:

```console
$ docker run -d --name todo -p 8080:3000 <YOUR_DOCKER_USERNAME>/getting-started-todo-app
```

Open [http://localhost:8080](http://localhost:8080). This time, the complete
application is running from the image you built.

## Share the image

In [Docker Home](https://app.docker.com), create a public repository named
`getting-started-todo-app` in your personal namespace.

Sign in from the command line, then push the image:

```console
$ docker login
$ docker push <YOUR_DOCKER_USERNAME>/getting-started-todo-app
```

Open the repository in Docker Home. The image is ready for another machine or
deployment system to pull and run.

## Clean up

Remove the standalone container and the development stack:

```console
$ docker rm -f todo
$ docker compose down --volumes
```

## What you proved

You ran software packaged by someone else, started an application stack from a
declarative file, built your own image, and published it. Containers carry the
application environment, Compose coordinates multiple containers, and a
registry makes images available beyond one machine.

Choose a language-specific [Docker guide](/guides/) to containerize an
application of your own.
