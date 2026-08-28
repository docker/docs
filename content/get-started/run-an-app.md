---
title: Run an application with Docker
linkTitle: Run an application
description: Start a multi-service application with Docker Compose, change its code, and see the update without installing the application stack.
keywords: Docker, get started, containers, Docker Compose, Compose Watch
weight: 1
---

A development application can depend on language runtimes, a database, and
supporting tools. Installing and matching that stack on every machine delays the
first useful result.

In this 10-minute tutorial, you'll start a prepared five-service application
with one command, change its code, and see the result. Your host needs only Git
and Docker Desktop.

## Before you start

- Install [Git](https://git-scm.com/downloads)
- [Install Docker Desktop](get-docker.md) and start it

## Get the application

Clone the sample application and open its directory:

```console
$ git clone https://github.com/docker/getting-started-todo-app
$ cd getting-started-todo-app
```

The prepared project gives Docker something concrete to run. You will change one
source file. Its `compose.yaml` file already defines the rest of the development
environment.

## Start the application

Start the environment and watch the project for file changes:

```console
$ docker compose watch
```

Docker builds the application images, pulls the images for supporting services,
and starts the containers. The first run prints build and startup logs in the
terminal.

After the services start, open [http://localhost](http://localhost). You should
see a greeting and an empty to-do list.

Add an item to confirm that the application and its database are working.

## Change the application

Open `backend/src/routes/getGreeting.js` in a text editor. Change the first line
to use your own greeting:

```javascript
const GREETING = 'Hello from Docker!';
```

Save the file, then refresh [http://localhost](http://localhost). The application
displays your greeting while `docker compose watch` keeps running. Compose Watch
detected the edit and copied the file into the backend container.

## See what Docker started

Open another terminal in the project directory and list the services:

```console
$ docker compose ps
```

The `STATUS` column shows each service as running, with the MySQL service marked
as healthy. One command started the frontend, API, database, database management
interface, and proxy without installing Node.js, MySQL, or Traefik on your host.

## Stop the application

Press `Ctrl+C` in the terminal running Compose Watch. Then remove the containers,
network, and sample database volume:

```console
$ docker compose down --volumes
```

Your source files remain in the project directory.

## What you proved

You started a complete development environment from a version-controlled
description, changed the running application, and removed the environment. The
project carried its required stack instead of relying on a matching stack on
your host.

To containerize an application of your own, choose its language from the
[Docker guides](/guides/).
