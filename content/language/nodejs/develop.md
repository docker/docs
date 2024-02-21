---
title: Use containers for Node.js development
keywords: node, node.js, development
description: Learn how to develop your Node.js application locally using containers.
aliases:
- /get-started/nodejs/develop/
---

## Prerequisites

Complete [Containerize a Node.js application](containerize.md).

## Overview

In this section, you'll learn how to set up a development environment for your containerized application. This includes:
 - Adding a local database and persisting data
 - Configuring your container to run a development environment
 - Debugging your containerized application

## Add a local database and persist data

You can use containers to set up local services, like a database. In this section, you'll update the `compose.yaml` file to define a database service and a volume to persist data.

Open the `compose.yaml` file in an IDE or text editor. You'll notice it
already contains commented-out instructions for a Postgres database and volume.

Open `src/persistence/postgres.js` in an IDE or text editor. You'll notice that
this application uses a Postgres database and requires some environment
variables in order to connect to the database. The `compose.yaml` file doesn't
have these variables defined.

You need to update the following items in the `compose.yaml` file:
 - Uncomment all of the database instructions.
 - Add the environment variables under the server service.
 - Add `secrets` to the server service for the database password.

The following is the updated `compose.yaml` file.

```yaml {hl_lines="7-40"}
services:
  server:
    build:
      context: .
    ports:
      - 3000:3000
    environment:
      NODE_ENV: production
      POSTGRES_HOST: db
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD_FILE: /run/secrets/db-password
      POSTGRES_DB: example
    depends_on:
      db:
        condition: service_healthy
    secrets:
      - db-password
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

> **Note**
>
> To learn more about the instructions in the Compose file, see [Compose file
> reference](/compose/compose-file/).

Before you run the application using Compose, notice that this Compose file uses
`secrets` and specifies a `password.txt` file to hold the database's password.
You must create this file as it's not included in the source repository.

In the cloned repository's directory, create a new directory named `db`. Inside the `db` directory, create a file named `password.txt`. Open `password.txt` in an IDE or text editor and add a password of your choice. The password must be on a single line with no additional lines in the file.

You should now have the following contents in your `docker-nodejs-sample`
directory.

```text
├── docker-nodejs-sample/
│ ├── db/
│ │ └── password.txt
│ ├── spec/
│ ├── src/
│ ├── .dockerignore
│ ├── .gitignore
│ ├── compose.yaml
│ ├── Dockerfile
│ ├── package-lock.json
│ ├── package.json
│ ├── README.Docker.md
│ └── README.md
```

Run the following command to start your application.

```console
$ docker compose up --build
```

Open a browser and verify that the application is running at [http://localhost:3000](http://localhost:3000).

Add some items to the todo list to test data persistence.

After adding some items to the todo list, press `ctrl+c` in the terminal to stop your application.

In the terminal, run `docker compose rm` to remove your containers and then run `docker compose up` to run your application again.

```console
$ docker compose rm
$ docker compose up --build
```

Refresh [http://localhost:3000](http://localhost:3000) in your browser and verify that the todo items persisted, even after the containers were removed and ran again.

## Configure and run a development container

You can use a bind mount to mount your source code into the container. The container can then see the changes you make to the code immediately, as soon as you save a file. This means that you can run processes, like nodemon, in the container that watch for filesystem changes and respond to them. To learn more about bind mounts, see [Storage overview](../../storage/index.md).

In addition to adding a bind mount, you can configure your Dockerfile and `compose.yaml` file to install development dependencies and run development tools.

### Update your Dockerfile for development

Open the Dockerfile in an IDE or text editor. Note that the Dockerfile doesn't
install development dependencies and doesn't run nodemon. You'll
need to update your Dockerfile to install the development dependencies and run
nodemon.

Rather than creating one Dockerfile for production, and another Dockerfile for
development, you can use one multi-stage Dockerfile for both.

Update your Dockerfile to the following multi-stage Dockerfile.

```dockerfile {hl_lines="5-26"}
# syntax=docker/dockerfile:1

ARG NODE_VERSION=18.0.0

FROM node:${NODE_VERSION}-alpine as base
WORKDIR /usr/src/app
EXPOSE 3000

FROM base as dev
RUN --mount=type=bind,source=package.json,target=package.json \
    --mount=type=bind,source=package-lock.json,target=package-lock.json \
    --mount=type=cache,target=/root/.npm \
    npm ci --include=dev
USER node
COPY . .
CMD npm run dev

FROM base as prod
RUN --mount=type=bind,source=package.json,target=package.json \
    --mount=type=bind,source=package-lock.json,target=package-lock.json \
    --mount=type=cache,target=/root/.npm \
    npm ci --omit=dev
USER node
COPY . .
CMD node src/index.js
```

In the Dockerfile, you first add a label `as base` to the `FROM
node:${NODE_VERSION}-alpine` statement. This lets you refer to this build stage
in other build stages. Next, you add a new build stage labeled `dev` to install
your development dependencies and start the container using `npm run dev`.
Finally, you add a stage labeled `prod` that omits the dev dependencies and runs
your application using `node src/index.js`. To learn more about multi-stage
builds, see [Multi-stage builds](../../build/building/multi-stage.md).

Next, you'll need to update your Compose file to use the new stage.

### Update your Compose file for development

To run the `dev` stage with Compose, you need to update your `compose.yaml`
file. Open your `compose.yaml` file in an IDE or text editor, and then add the
`target: dev` instruction to target the `dev` stage from your multi-stage
Dockerfile.

Also, add a new volume to the server service for the bind mount. For this application, you'll mount `./src` from your local machine to `/usr/src/app/src` in the container.

Lastly, publish port `9229` for debugging.

The following is the updated Compose file.

```yaml {hl_lines=[5,8,20,21]}
services:
  server:
    build:
      context: .
      target: dev
    ports:
      - 3000:3000
      - 9229:9229
    environment:
      NODE_ENV: production
      POSTGRES_HOST: db
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD_FILE: /run/secrets/db-password
      POSTGRES_DB: example
    depends_on:
      db:
        condition: service_healthy
    secrets:
      - db-password
    volumes:
      - ./src:/usr/src/app/src
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

### Run your development container and debug your application

Run the following command to run your application with the new changes to the `Dockerfile` and `compose.yaml` file.

```console
$ docker compose up --build
```

Open a browser and verify that the application is running at [http://localhost:3000](http://localhost:3000).

Any changes to the application's source files on your local machine will now be
immediately reflected in the running container.

Open `docker-nodejs-sample/src/static/js/app.js` in an IDE or text editor and update the button text on line 109 from `Add Item` to `Add`.

```diff
+                         {submitting ? 'Adding...' : 'Add'}
-                         {submitting ? 'Adding...' : 'Add Item'}
```

Refresh [http://localhost:3000](http://localhost:3000) in your browser and verify that the updated text appears.

You can now connect an inspector client to your application for debugging. For
more details about inspector clients, see the [Node.js
documentation](https://nodejs.org/en/docs/guides/debugging-getting-started).

## Debug your containers

While running your application in a container, you may come across issues that aren't related to the application's code. Issues may be related to the container environment that your application is running in. In this section, you'll walk through the steps you may take to debug a scenario where the application container gives an error about no database connection.

To debug a database connection issue:

1. Verify that the database container is running.
   
   You can use the `docker ps` command to verify which containers are running and get the container IDs.

   ```console
   $ docker ps

   CONTAINER ID   IMAGE                         COMMAND                  CREATED          STATUS                    PORTS                                            NAMES
   78b36184f881   docker-nodejs-sample-server   "docker-entrypoint.s…"   54 seconds ago   Up 42 seconds             0.0.0.0:3000->3000/tcp, 0.0.0.0:9229->9229/tcp   docker-nodejs-sample-server-1
   5f6271e301a8   postgres                      "docker-entrypoint.s…"   54 seconds ago   Up 52 seconds (healthy)   5432/tcp                                         docker-nodejs-sample-db-1
   ```

   In the previous example, you can see the `postgres` container is running and its ID is `5f6271e301a8`.

2. Inspect the logs of the containers.

   You can view the logs of containers by running the `docker logs` command with
   the container ID. Replace the container ID with your container ID.

   ```console
   $ docker logs 5f6271e301a8

   ...
   2024-02-21 20:34:40.416 UTC [49] LOG:  checkpoint starting: time
   2024-02-21 20:34:46.882 UTC [49] LOG:  checkpoint complete: wrote 67 buffers (0.4%); 0 WAL file(s) added, 0 removed, 0 recycled; write=6.431 s, sync=0.023 s, total=6.466 s; sync files=31, longest=0.015 s, average=0.001 s; distance=354 kB, estimate=354 kB; lsn=0/196BB18, redo lsn=0/196BAE0
   ```
   
   You can look through the container logs for any error messages. For example,
   you can look for credential or network configuration errors.

3. Get a shell into the container to interact with the environment and execute
   commands.

   You can use the `docker exec` command to get a shell into some containers.
   Many containers may use slim images that have no shell. The easiest way to
   get a shell into any container, even those running slim images, is to use the
   `docker debug` command. In debugging scenarios, the `docker debug` command is
   a replacement for `docker exec` and has several advantages, such as not
   modifying the environment of containers and providing access to a
   customizable toolbox.

   Use one of the following commands to get a shell into your container.

   {{< tabs >}}
   {{< tab name="docker exec" >}}

   1. Get a shell into the database container with the following command.
      Replace the container ID with your own container ID.

      ```console
      $ docker exec -u 0 -it 5f6271e301a8 /bin/sh
      ```
      
      The following is a breakdown of the `docker exec` command options.
      - `-u 0`: This option specifies the user that the command should run as. 0
        is the user ID for the root user, so this command will run with root
        privileges inside the container. Specifying the user ensures that the
        command has the appropriate permissions for its operations.
      - `-i` (or `--interactive`): This flag keeps the STDIN (standard input)
        open even if not attached. It allows you to interact with the command
        line of the container through your terminal. This is useful for
        interactive commands that require input from the user.
      - `-t` (or `--tty`): This allocates a pseudo-TTY, which is essentially a
        virtual terminal inside the container. This makes it possible to
        interact with the command line interface of the container in a more
        user-friendly way, as if you were using a terminal session directly on
        the host. It supports input, output, and error flow, along with terminal
        resizing and signals.
      - `5f6271e301a8`: This is the container ID or name. It specifies which
        container you want to execute the command in. Docker uses this to
        identify the specific container instance you're targeting for your
        command.
      - `/bin/sh`: This is the command you want to run inside the container. In
        this case, it's launching a shell (sh), which allows you to interact
        with the container's file system and execute further commands inside it.
        If you're using Git Bash in Windows, use `//bin/sh`. When a container
        has no shell, you can use `docker debug`.

   2. Inside the shell, install tools to debug the issue. The following command
      installs ping.

      ```console
      # apt-get update && apt-get iputils-ping
      ```

      Note that installing tools in the container will increase its size and
      attack surface. In addition, any tools installed aren't available in
      future containers. To avoid these issues, Docker recommends using `docker
      debug`.

  3. Inside the shell, use the installed tools. The following command pings the
     `server` service to verify connectivity.

      ```console
      # ping -c 1 server
      PING server (172.26.0.3) 56(84) bytes of data.
      64 bytes from docker-nodejs-sample-server-1.docker-nodejs-sample_default (172.26.0.3): icmp_seq=1 ttl=64 time=0.131 ms
      ```
   {{< /tab >}}
   {{< tab name="docker debug" >}}

   > **Note**
   >
   > Docker Debug requires a [Pro, Team, or Business subscription](../../subscription/details.md). You must [sign in](../../desktop/get-started.md) to use this command.

   1. Get a shell into the database container with the following command.
      Replace the container ID with your own container ID.

      ```console
      $ docker debug 5f6271e301a8
      ```

   2. Inside the shell, use tools from Docker Debug's toolbox. The following
     command pings the `server` service to verify connectivity.

      ```console
      docker > ping -c 1 server
      PING server (172.26.0.3) 56(84) bytes of data.
      64 bytes from docker-nodejs-sample-server-1.docker-nodejs-sample_default (172.26.0.3): icmp_seq=1 ttl=64 time=0.066 ms
      ```

      Docker Debug's toolbox comes with many standard Linux tools pre-installed,
      such as `ping`. To learn how to install more tools in your toolbox, see
      [docker debug](../../reference/cli/docker/debug.md).

   {{< /tab >}}
   {{< /tabs >}}

## Summary

In this section, you took a look at setting up your Compose file to add a local
database and persist data. You also learned how to create a multi-stage
Dockerfile and set up a bind mount for development. Finally, you learned how to debug container issues.

Related information:
 - [Volumes top-level element](/compose/compose-file/07-volumes/)
 - [Services top-level element](/compose/compose-file/05-services/)
 - [Multi-stage builds](../../build/building/multi-stage.md)
 - [docker CLI reference](../../reference/cli/docker/_index.md)

## Next steps

In the next section, you'll learn how to run unit tests using Docker.

{{< button text="Run your tests" url="run-tests.md" >}}
