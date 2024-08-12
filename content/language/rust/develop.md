---
title: Develop your Rust application
keywords: rust, local, development, run,
description: Learn how to develop your Rust application locally.
---

## Prerequisites

* You have installed the latest version of [Docker Desktop](../../get-docker.md).
* You have completed the walkthroughs in the Docker Desktop [Learning Center](../../desktop/get-started.md) to learn about Docker concepts.
* You have a [git client](https://git-scm.com/downloads). The examples in this section use a command-line based git client, but you can use any client.

## Overview

In this section, you’ll learn how to use volumes and networking in Docker. You’ll also use Docker to build your images and Docker Compose to make everything a whole lot easier.

First, you’ll take a look at running a database in a container and how you can use volumes and networking to persist your data and let your application talk with the database. Then you’ll pull everything together into a Compose file which lets you set up and run a local development environment with one command.

## Run a database in a container

Instead of downloading PostgreSQL, installing, configuring, and then running the PostgreSQL database as a service, you can use the Docker Official Image for PostgreSQL and run it in a container.

Before you run PostgreSQL in a container, create a volume that Docker can manage to store your persistent data and configuration. Use the named volumes feature that Docker provides instead of using bind mounts.

Run the following command to create your volume.

```console
$ docker volume create db-data
```

Now create a network that your application and database will use to talk to each other. The network is called a user-defined bridge network and gives you a nice DNS lookup service which you can use when creating your connection string.

```console
$ docker network create postgresnet
```

Now you can run PostgreSQL in a container and attach to the volume and network that you created previously. Docker pulls the image from Hub and runs it for you locally.
In the following command, option `--mount` is for starting the container with a volume. For more information, see [Docker volumes](/engine/storage/volumes.md).

```console
$ docker run --rm -d --mount \
  "type=volume,src=db-data,target=/var/lib/postgresql/data" \
  -p 5432:5432 \
  --network postgresnet \
  --name db \
  -e POSTGRES_PASSWORD=mysecretpassword \
  -e POSTGRES_DB=example \
  postgres
```

Now, make sure that your PostgreSQL database is running and that you can connect to it. Connect to the running PostgreSQL database inside the container.

```console
$ docker exec -it db psql -U postgres
```

You should see output like the following.

```console
psql (15.3 (Debian 15.3-1.pgdg110+1))
Type "help" for help.

postgres=#
```

In the previous command, you logged in to the PostgreSQL database by passing the `psql` command to the `db` container. Press ctrl-d to exit the PostgreSQL interactive terminal.


## Get and run the sample application

For the sample application, you'll use a variation of the backend from the react-rust-postgres application from [Awesome Compose](https://github.com/docker/awesome-compose/tree/master/react-rust-postgres).

1. Clone the sample application repository using the following command.

   ```console
   $ git clone https://github.com/docker/docker-rust-postgres
   ```

2. In the cloned repository's directory, run `docker init` to create the necessary Docker files. Refer to the following example to answer the prompts from `docker init`.

   ```console
   $ docker init
   Welcome to the Docker Init CLI!

   This utility will walk you through creating the following files with sensible defaults for your project:
     - .dockerignore
     - Dockerfile
     - compose.yaml
     - README.Docker.md

   Let's get started!

   ? What application platform does your project use? Rust
   ? What version of Rust do you want to use? 1.70.0
   ? What port does your server listen on? 8000
   ```

3. In the cloned repository's directory, open the `Dockerfile` in an IDE or text editor to update it.

   `docker init` handled creating most of the instructions in the Dockerfile, but you'll need to update it for your unique application. In addition to a `src` directory, this application includes a `migrations` directory to initialize the database. Add a bind mount for the `migrations` directory to the build stage in the Dockerfile. The following is the updated Dockerfile.

   ```dockerfile {hl_lines="28"}
   # syntax=docker/dockerfile:1

   # Comments are provided throughout this file to help you get started.
   # If you need more help, visit the Dockerfile reference guide at
   # https://docs.docker.com/reference/dockerfile/
   
   ################################################################################
   # Create a stage for building the application.
   
   ARG RUST_VERSION=1.70.0
   ARG APP_NAME=react-rust-postgres
   FROM rust:${RUST_VERSION}-slim-bullseye AS build
   ARG APP_NAME
   WORKDIR /app
   
   # Build the application.
   # Leverage a cache mount to /usr/local/cargo/registry/
   # for downloaded dependencies and a cache mount to /app/target/ for 
   # compiled dependencies which will speed up subsequent builds.
   # Leverage a bind mount to the src directory to avoid having to copy the
   # source code into the container. Once built, copy the executable to an
   # output directory before the cache mounted /app/target is unmounted.
   RUN --mount=type=bind,source=src,target=src \
       --mount=type=bind,source=Cargo.toml,target=Cargo.toml \
       --mount=type=bind,source=Cargo.lock,target=Cargo.lock \
       --mount=type=cache,target=/app/target/ \
       --mount=type=cache,target=/usr/local/cargo/registry/ \
       --mount=type=bind,source=migrations,target=migrations \
       <<EOF
   set -e
   cargo build --locked --release
   cp ./target/release/$APP_NAME /bin/server
   EOF
   
   ################################################################################
   # Create a new stage for running the application that contains the minimal
   # runtime dependencies for the application. This often uses a different base
   # image from the build stage where the necessary files are copied from the build
   # stage.
   #
   # The example below uses the debian bullseye image as the foundation for    running the app.
   # By specifying the "bullseye-slim" tag, it will also use whatever happens to    be the
   # most recent version of that tag when you build your Dockerfile. If
   # reproducability is important, consider using a digest
   # (e.g.,    debian@sha256:ac707220fbd7b67fc19b112cee8170b41a9e97f703f588b2cdbbcdcecdd8af57).
   FROM debian:bullseye-slim AS final
   
   # Create a non-privileged user that the app will run under.
   # See https://docs.docker.com/develop/develop-images/dockerfile_best-practices/   #user
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
   
   # Copy the executable from the "build" stage.
   COPY --from=build /bin/server /bin/
   
   # Expose the port that the application listens on.
   EXPOSE 8000
   
   # What the container should run when it is started.
   CMD ["/bin/server"]
   ```

4. In the cloned repository's directory, run `docker build` to build the image.

   ```console
   $ docker build -t rust-backend-image .
   ```

5. Run `docker run` with the following options to run the image as a container on the same network as the database.

   ```console
   $ docker run \
     --rm -d \
     --network postgresnet \
     --name docker-develop-rust-container \
     -p 3001:8000 \
     -e PG_DBNAME=example \
     -e PG_HOST=db \
     -e PG_USER=postgres \
     -e PG_PASSWORD=mysecretpassword \
     -e ADDRESS=0.0.0.0:8000 \
     -e RUST_LOG=debug \
     rust-backend-image
   ```

6. Curl the application to verify that it connects to the database.

   ```console
   $ curl http://localhost:3001/users
   ```

   You should get a response like the following.

   ```json
   [{"id":1,"login":"root"}]
   ```

## Use Compose to develop locally

When you run `docker init`, in addition to a `Dockerfile`, it also creates a `compose.yaml` file.

This Compose file is super convenient as you don't have to type all the parameters to pass to the `docker run` command. You can declaratively do that using a Compose file.

In the cloned repository's directory, open the `compose.yaml` file in an IDE or text editor. `docker init` handled creating most of the instructions, but you'll need to update it for your unique application.

You need to update the following items in the `compose.yaml` file:
 - Uncomment all of the database instructions.
 - Add the environment variables under the server service.

The following is the updated `compose.yaml` file.

```yaml {hl_lines=["17-23","30-55"]}
# Comments are provided throughout this file to help you get started.
# If you need more help, visit the Docker compose reference guide at
# https://docs.docker.com/compose/compose-file/

# Here the instructions define your application as a service called "server".
# This service is built from the Dockerfile in the current directory.
# You can add other services your application may depend on here, such as a
# database or a cache. For examples, see the Awesome Compose repository:
# https://github.com/docker/awesome-compose
services:
  server:
    build:
      context: .
      target: final
    ports:
      - 8000:8000
    environment:
      - PG_DBNAME=example
      - PG_HOST=db
      - PG_USER=postgres
      - PG_PASSWORD=mysecretpassword
      - ADDRESS=0.0.0.0:8000
      - RUST_LOG=debug
# The commented out section below is an example of how to define a PostgreSQL
# database that your application can use. `depends_on` tells Docker Compose to
# start the database before your application. The `db-data` volume persists the
# database data between container restarts. The `db-password` secret is used
# to set the database password. You must create `db/password.txt` and add
# a password of your choosing to it before running `docker compose up`.
    depends_on:
      db:
        condition: service_healthy
  db:
    image: postgres
    restart: always
    user: postgres
    secrets:
      - db-password
    volumes:
      - db-data:/var/lib/postgresql/data
    environment:
      - POSTGRES_DB=example
      - POSTGRES_PASSWORD_FILE=/run/secrets/db-password
    expose:
      - 5432
    healthcheck:
      test: [ "CMD", "pg_isready" ]
      interval: 10s
      timeout: 5s
      retries: 5
volumes:
  db-data:
secrets:
  db-password:
    file: db/password.txt
```

Note that the file doesn't specify a network for those 2 services. When you use Compose, it automatically creates a network and connects the services to it. For more information see [Networking in Compose](../../compose/networking.md).

Before you run the application using Compose, notice that this Compose file specifies a `password.txt` file to hold the database's password. You must create this file as it's not included in the source repository.

In the cloned repository's directory, create a new directory named `db` and inside that directory create a file named `password.txt` that contains the password for the database. Using your favorite IDE or text editor, add the following contents to the `password.txt` file.

```text
mysecretpassword
```

If you have any other containers running from the previous sections, [stop](./run-containers.md/#stop-start-and-name-containers) them now.

Now, run the following `docker compose up` command to start your application.

```console
$ docker compose up --build
```

The command passes the `--build` flag so Docker will compile your image and then start the containers.

Now test your API endpoint. Open a new terminal then make a request to the server using the curl commands:

```console
$ curl http://localhost:8000/users
```

You should receive the following response:

```json
[{"id":1,"login":"root"}]
```

## Summary

In this section, you took a look at setting up your Compose file to run your Rust application and database with a single command.

Related information:
 - [Docker volumes](/engine/storage/volumes.md)
 - [Compose overview](../../compose/index.md)

## Next steps

In the next section, you'll take a look at how to set up a CI/CD pipeline using GitHub Actions.

{{< button text="Configure CI/CD" url="configure-ci-cd.md" >}}
