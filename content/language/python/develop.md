---
title: Use containers for development
keywords: python, local, development, run,
description: Learn how to develop your application locally.
---

## Prerequisites

Work through the steps to build an image and run it as a containerized application in [Run your image as a container](run-containers.md).

## Introduction

In this section, you’ll learn how to use volumes and networking in Docker. You’ll also use Docker to build your images and Docker Compose to make everything a whole lot easier.

First, you’ll take a look at running a database in a container and how you can use volumes and networking to persist your data and let your application to talk with the database. Then you’ll pull everything together into a Compose file which lets you to set up and run a local development environment with one command.

## Run a database in a container

Instead of downloading PostgreSQL, installing, configuring, and then running the PostgreSQL database on your system directly, you can use the Docker Official Image for PostgreSQL and run it in a container.

Before you run PostgreSQL in a container, create a volume that Docker can manage to store your persistent data and configuration.

Run the following command to create your volume.

```console
$ docker volume create db-data
```

Now create a network that your application and database will use to talk to each other. The network is called a user-defined bridge network and gives you a nice DNS lookup service which you can use when creating your connection string.

```console
$ docker network create postgresnet
```

Now you can run PostgreSQL in a container and attach to the volume and network that you created above. Docker pulls the image from Hub and runs it for you locally.
In the following command, option `--mount` is for starting the container with a volume. For more information, see [Docker volumes](../../storage/volumes.md).
   
   {{< tabs >}}
   {{< tab name="Mac / Linux" >}}

   ```console
   $ docker run --rm -d \
     --mount type=volume,src=db-data,target=/var/lib/postgresql/data \
     -p 5432:5432 \
     --network postgresnet \
     --name db \
     -e POSTGRES_PASSWORD=mysecretpassword \
     -e POSTGRES_DB=example \
     postgres
   ```

   {{< /tab >}}
   {{< tab name="Windows" >}}

   ```powershell
   $ docker run --rm -d `
     --mount type=volume,src=db-data,target=/var/lib/postgresql/data `
     -p 5432:5432 `
     --network postgresnet `
     --name db `
     -e POSTGRES_PASSWORD=mysecretpassword `
     -e POSTGRES_DB=example `
     postgres
   ```

   {{< /tab >}}
   {{< /tabs >}}

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

You'll need to clone a new repository to get a sample application that includes logic to connect to the database.

1. Change to a directory where you want to clone the repository and run the following command.

   ```console
   $ git clone https://github.com/docker/python-docker-dev
   ```

2. In the cloned repository's directory, run `docker init` to create the necessary Docker files. Refer to the following example to answer the prompts from `docker init`.

   ```console
   $ docker init
   Welcome to the Docker Init CLI!

   This utility will walk you through creating the following files with sensible defaults for your project:
     - .dockerignore
     - Dockerfile
     - compose.yaml

   Let's get started!

   ? What application platform does your project use? Python
   ? What version of Python do you want to use? 3.11.4
   ? What port do you want your app to listen on? 5000
   ? What is the command to run your app? python3 -m flask run --host=0.0.0.0
   ```

3. In the cloned repository's directory, run `docker build` to build the image.

   ```console
   $ docker build -t python-docker-dev .
   ```

4. If you have any containers running from the previous sections using the name `rest-server` or port 8000, [stop and remove](./run-containers.md/#stop-start-and-name-containers) them now.

5. Run `docker run` with the following options to run the image as a container on the same network as the database.
   
   {{< tabs >}}
   {{< tab name="Mac / Linux" >}}

   ```console
   $ docker run --rm -d \
     --network postgresnet \
     --name rest-server \
     -p 8000:5000 \
     -e POSTGRES_PASSWORD=mysecretpassword \
     python-docker-dev
   ```

   {{< /tab >}}
   {{< tab name="Windows" >}}

   ```powershell
   $ docker run --rm -d `
     --network postgresnet `
     --name rest-server `
     -p 8000:5000 `
     -e POSTGRES_PASSWORD=mysecretpassword `
     python-docker-dev
   ```

   {{< /tab >}}
   {{< /tabs >}}

6. Test that your application is connected to the database and is able to list the widgets.

   ```console
   $ curl http://localhost:8000/initdb
   $ curl http://localhost:8000/widgets
   ```

   You should receive the following JSON back from your service.

   ```json
   []
   ```
   This is because your database is empty.

## Use Compose to develop locally

When you run `docker init`, in addition to a `Dockerfile`, it also creates a `compose.yaml` file.

This Compose file is super convenient as you don't have to type all the parameters to pass to the `docker run` command. You can declaratively do that using a Compose file.

In the cloned repository's directory, open the `compose.yaml` file in an IDE or text editor. `docker init` handled creating most of the instructions, but you'll need to update it for your unique application.

In the `compose.yaml` file, you need to uncomment all of the database instructions. In addition, you need to add the database password as an environment variable to the server service.

The following is the updated `compose.yaml` file.

```yaml
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
    ports:
      - 5000:5000
    environment:
      - POSTGRES_PASSWORD=mysecretpassword

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

Note that the file doesn't specify a network for those 2 services. Compose automatically creates a network and connects the services to it. For more information see [Networking in Compose](../../compose/networking.md).

Before you run the application using Compose, notice that this Compose file specifies a `password.txt` file to hold the database's password. You must create this file as it's not included in the source repository.

In the cloned repository's directory, create a new directory named `db` and inside that directory create a file named `password.txt` that contains the password for the database. Using your favorite IDE or text editor, add the following contents to the `password.txt` file.

```
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
$ curl http://localhost:5000/initdb
$ curl http://localhost:5000/widgets
```

You should receive the following response:

```json
[]
```

This is because your database is empty.

# Summary

In this section, you took a look at setting up your Compose file to run your Python application and database with a single command.

Related information:
 - [Volumes](../../storage/volumes.md)
 - [Compose overview](../../compose/index.md)

## Next steps

In the next section, you'll take a look at how to set up a CI/CD pipeline using GitHub Actions.

{{< button text="Configure CI/CD" url="configure-ci-cd.md" >}}
