---
title: Introduction
description: An introduction to the Docker Build guide
keywords: build, buildkit, buildx, guide, tutorial, introduction
---

The starting resources for this guide include a simple Go project and a
Dockerfile. From this starting point, the guide illustrates various ways that
you can improve how you build the application with Docker.

## Environment setup

To follow this guide:

1. Install [Docker Desktop or Docker Engine](../../get-docker.md)
2. Clone or create a new repository from the
   [application example on GitHub](https://github.com/dockersamples/buildme)

## The application

The example project for this guide is a client-server application for
translating messages to a fictional language.

Here’s an overview of the files included in the project:

```text
.
├── Dockerfile
├── cmd
│   ├── client
│   │   ├── main.go
│   │   ├── request.go
│   │   └── ui.go
│   └── server
│       ├── main.go
│       └── translate.go
├── go.mod
└── go.sum
```

The `cmd/` directory contains the code for the two application components:
client and server. The client is a user interface for writing, sending, and
receiving messages. The server receives messages from clients, translates them,
and sends them back to the client.

## The Dockerfile

A Dockerfile is a text document in which you define the build steps for your
application. You write the Dockerfile in a domain-specific language, called the
Dockerfile syntax.

Here's the Dockerfile used as the starting point for this guide:

```dockerfile
# syntax=docker/dockerfile:1
FROM golang:{{% param "example_go_version" %}}-alpine
WORKDIR /src
COPY . .
RUN go mod download
RUN go build -o /bin/client ./cmd/client
RUN go build -o /bin/server ./cmd/server
ENTRYPOINT [ "/bin/server" ]
```

Here’s what this Dockerfile does:

1. `# syntax=docker/dockerfile:1`

   This comment is a
   [Dockerfile parser directive](../../engine/reference/builder.md#parser-directives).
   It specifies which version of the Dockerfile syntax to use. This file uses
   the `dockerfile:1` syntax which is best practice: it ensures that you have
   access to the latest Docker build features.

2. `FROM golang:{{% param "example_go_version" %}}-alpine`

   The `FROM` instruction uses version `{{% param "example_go_version" %}}-alpine` of the `golang` official image.

3. `WORKDIR /src`

   Creates the `/src` working directory inside the container.

4. `COPY . .`

   Copies the files in the build context to the working directory in the
   container.

5. `RUN go mod download`

   Downloads the necessary Go modules to the container. Go modules is the
   dependency management tool for the Go programming language, similar to
   `npm install` for JavaScript, or `pip install` for Python.

6. `RUN go build -o /bin/client ./cmd/client`

   Builds the `client` binary, which is used to send messages to be translated, into the
   `/bin` directory.

7. `RUN go build -o /bin/server ./cmd/server`

   Builds the `server` binary, which listens for client translation requests,
   into the `/bin` directory.

8. `ENTRYPOINT [ "/bin/server" ]`

   Specifies a command to run when the container starts. Starts the server
   process.

## Build the image

To build an image using a Dockerfile, you use the `docker` command-line tool.
The command for building an image is `docker build`.

Run the following command to build the image.

```console
$ docker build --tag=buildme .
```

This creates an image with the tag `buildme`. An image tag is the name of the
image.

## Run the container

The image you just built contains two binaries, one for the server and one for
the client. To see the translation service in action, run a container that hosts
the server component, and then run another container that invokes the client.

To run a container, you use the `docker run` command.

1. Run a container from the image in detached mode.

   ```console
   $ docker run --name=buildme --rm --detach buildme
   ```

   This starts a container named `buildme`.

2. Run a new command in the `buildme` container that invokes the client binary.

   ```console
   $ docker exec -it buildme /bin/client
   ```

The `docker exec` command opens a terminal user interface where you can submit
messages for the backend (server) process to translate.

When you're done testing, you can stop the container:

```console
$ docker stop buildme
```

## Summary

This section gave you an overview of the example application used in this guide,
an introduction to Dockerfiles and building. You've successfully built a
container image and created a container from it.

Related information:

- [Dockerfile reference](../../engine/reference/builder.md)
- [`docker build` CLI reference](../../engine/reference/commandline/image_build.md)
- [`docker run` CLI reference](../../engine/reference/commandline/container_run.md)

## Next steps

The next section explores how you can use layer cache to improve build speed.

{{< button text="Layers" url="layers.md" >}}
