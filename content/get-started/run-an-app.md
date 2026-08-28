---
title: Run an application with Docker
linkTitle: Run an application
description: Start a multi-service application with Docker Compose, change its code, and see the update without installing the application stack.
keywords: Docker, get started, containers, Docker Compose, Compose Watch
weight: 1
---

In this 10-minute tutorial, you'll start a five-service development application
with one command and see a code change update while it runs. Docker provides the
runtimes and database, so your host only needs Git and Docker Desktop.

## Before you start

- Install [Git](https://git-scm.com/downloads)
- [Install Docker Desktop](get-docker.md) and start it

## Get the application

Clone the sample application and open its directory:

```console
$ git clone https://github.com/docker/getting-started-todo-app
$ cd getting-started-todo-app
```

The project includes a `compose.yaml` file that defines its complete development
environment.

## Start the application

Start the environment and watch the project for file changes:

```console
$ docker compose watch
```

Docker builds the application images, pulls the images for supporting services,
and starts the containers. After the services start, open
[http://localhost](http://localhost) to see the to-do application.

Add an item to confirm that the application and its database are working.

## Change the application

Open `backend/src/routes/getGreeting.js` in a text editor. Change the first line
to use your own greeting:

```javascript
const GREETING = 'Hello from Docker!';
```

Save the file, then refresh [http://localhost](http://localhost). Compose Watch
copies the change into the backend container, and the application displays the
updated greeting.

## See what Docker started

Open another terminal in the project directory and list the services:

```console
$ docker compose ps
```

The application runs a frontend, an API, a MySQL database, a database management
interface, and a proxy in separate containers. Docker started the stack without
installing Node.js, MySQL, or Traefik on your host.

## Stop the application

Press `Ctrl+C` in the terminal running Compose Watch. Then remove the containers,
network, and sample database volume:

```console
$ docker compose down --volumes
```

Your source files remain in the project directory.

## What happened

The `compose.yaml` file described the full application stack. Docker created a
repeatable environment from that description, connected its services, and kept
the running code in sync with your editor.

You can hand the same project to another developer, run it in CI, or deploy it
without recreating the stack by hand.

To containerize an application of your own, choose its language from the
[Docker guides](/guides/).
