---
title: R language-specific guide
linkTitle: R
description: Containerize R apps using Docker
keywords: Docker, getting started, R, language
summary: |
  This guide details how to containerize R applications using Docker.
aliases:
  - /languages/r/
  - /guides/languages/r/
  - /language/R/build-images/
  - /language/R/run-containers/
  - /language/r/containerize/
  - /language/r/develop/
  - /language/r/configure-ci-cd/
  - /language/r/deploy/
  - /guides/r/configure-ci-cd/
  - /guides/r/containerize/
  - /guides/r/deploy/
  - /guides/r/develop/
params:
  tags: [languages]
  time: 10 minutes
---

The R language-specific guide teaches you how to containerize a R application using Docker. In this guide, you’ll learn how to:

- Containerize and run a R application
- Set up a local environment to develop a R application using containers

Start by containerizing an existing R application.

## Containerize a R application

### Prerequisites

- You have a [git client](https://git-scm.com/downloads). The examples in this section use a command-line based git client, but you can use any client.

### Overview

This section walks you through containerizing and running a R application.

### Get the sample application

The sample application uses the popular [Shiny](https://shiny.posit.co/) framework.

Clone the sample application to use with this guide. Open a terminal, change directory to a directory that you want to work in, and run the following command to clone the repository:

```console
$ git clone https://github.com/mfranzon/r-docker-dev.git && cd r-docker-dev
```

You should now have the following contents in your `r-docker-dev`
directory.

```text
├── r-docker-dev/
│ ├── src/
│ │ └── app.R
│ ├── src_db/
│ │ └── app_db.R
│ ├── compose.yaml
│ ├── Dockerfile
│ └── README.md
```

To learn more about the files in the repository, see the following:

- [Dockerfile](/reference/dockerfile.md)
- [.dockerignore](/reference/dockerfile.md#dockerignore-file)
- [compose.yaml](/reference/compose-file/_index.md)

### Run the application

Inside the `r-docker-dev` directory, run the following command in a
terminal.

```console
$ docker compose up --build
```

Open a browser and view the application at [http://localhost:3838](http://localhost:3838). You should see a simple Shiny application.

In the terminal, press `ctrl`+`c` to stop the application.

#### Run the application in the background

You can run the application detached from the terminal by adding the `-d`
option. Inside the `r-docker-dev` directory, run the following command
in a terminal.

```console
$ docker compose up --build -d
```

Open a browser and view the application at [http://localhost:3838](http://localhost:3838).

You should see a simple Shiny application.

In the terminal, run the following command to stop the application.

```console
$ docker compose down
```

For more information about Compose commands, see the [Compose CLI
reference](/reference/cli/docker/compose/).

### Summary

In this section, you learned how you can containerize and run your R
application using Docker.

Related information:

- [Docker Compose overview](/manuals/compose/_index.md)

### Next steps

In the next section, you'll learn how you can develop your application using
containers.

## Use containers for R development

### Prerequisites

Complete [Containerize a R application](./).

### Overview

In this section, you'll learn how to set up a development environment for your containerized application. This includes:

- Adding a local database and persisting data
- Configuring Compose to automatically update your running Compose services as you edit and save your code

### Get the sample application

You'll need to clone a new repository to get a sample application that includes logic to connect to the database.

Change to a directory where you want to clone the repository and run the following command.

```console
$ git clone https://github.com/mfranzon/r-docker-dev.git
```

### Configure the application to use the database

To try the connection between the Shiny application and the local database you have to modify the `Dockerfile` changing the `COPY` instruction:

```diff
-COPY src/ .
+COPY src_db/ .
```

### Add a local database and persist data

You can use containers to set up local services, like a database. In this section, you'll update the `compose.yaml` file to define a database service and a volume to persist data.

In the cloned repository's directory, open the `compose.yaml` file in an IDE or text editor.

In the `compose.yaml` file, you need to un-comment the properties for configuring the database. You must also mount the database password file and set an environment variable on the `shiny-app` service pointing to the location of the file in the container.

The following is the updated `compose.yaml` file.

```yaml
services:
  shiny-app:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - 3838:3838
    environment:
      - POSTGRES_PASSWORD_FILE=/run/secrets/db-password
    depends_on:
      db:
        condition: service_healthy
    secrets:
      - db-password
  db:
    image: postgres:18
    restart: always
    user: postgres
    secrets:
      - db-password
    volumes:
      - db-data:/var/lib/postgresql
    environment:
      - POSTGRES_DB=example
      - POSTGRES_PASSWORD_FILE=/run/secrets/db-password
    expose:
      - 5432
    healthcheck:
      test: ["CMD", "pg_isready"]
      interval: 10s
      timeout: 5s
      retries: 5
volumes:
  db-data:
secrets:
  db-password:
    file: db/password.txt
```

> [!NOTE]
>
> To learn more about the instructions in the Compose file, see [Compose file
> reference](/reference/compose-file/).

Before you run the application using Compose, notice that this Compose file specifies a `password.txt` file to hold the database's password. You must create this file as it's not included in the source repository.

In the cloned repository's directory, create a new directory named `db` and inside that directory create a file named `password.txt` that contains the password for the database. Using your favorite IDE or text editor, add the following contents to the `password.txt` file.

```text
mysecretpassword
```

Save and close the `password.txt` file.

You should now have the following contents in your `r-docker-dev`
directory.

```text
├── r-docker-dev/
│ ├── db/
│ │ └── password.txt
│ ├── src/
│ │ └── app.R
│ ├── src_db/
│ │ └── app_db.R
│ ├── requirements.txt
│ ├── .dockerignore
│ ├── compose.yaml
│ ├── Dockerfile
│ └── README.md
```

Now, run the following `docker compose up` command to start your application.

```console
$ docker compose up --build
```

Now test your DB connection opening a browser at:

```console
http://localhost:3838
```

You should see a pop-up message:

```text
DB CONNECTED
```

Press `ctrl+c` in the terminal to stop your application.

### Automatically update services

Use Compose Watch to automatically update your running Compose services as you
edit and save your code. For more details about Compose Watch, see [Use Compose
Watch](/manuals/compose/how-tos/file-watch.md).

Lines 15 to 18 in the `compose.yaml` file contain properties that trigger Docker
to rebuild the image when a file in the current working directory is changed:

```yaml {hl_lines="15-18",linenos=true}
services:
  shiny-app:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - 3838:3838
    environment:
      - POSTGRES_PASSWORD_FILE=/run/secrets/db-password
    depends_on:
      db:
        condition: service_healthy
    secrets:
      - db-password
    develop:
      watch:
        - action: rebuild
          path: .
  db:
    image: postgres:18
    restart: always
    user: postgres
    secrets:
      - db-password
    volumes:
      - db-data:/var/lib/postgresql
    environment:
      - POSTGRES_DB=example
      - POSTGRES_PASSWORD_FILE=/run/secrets/db-password
    expose:
      - 5432
    healthcheck:
      test: ["CMD", "pg_isready"]
      interval: 10s
      timeout: 5s
      retries: 5
volumes:
  db-data:
secrets:
  db-password:
    file: db/password.txt
```

Run the following command to run your application with Compose Watch.

```console
$ docker compose watch
```

Now, if you modify your `app.R` you will see the changes in real time without re-building the image!

Press `ctrl+c` in the terminal to stop your application.

### Summary

In this section, you took a look at setting up your Compose file to add a local
database and persist data. You also learned how to use Compose Watch to automatically rebuild and run your container when you update your code.

Related information:

- [Compose file reference](/reference/compose-file/)
- [Compose file watch](/manuals/compose/how-tos/file-watch.md)
- [Multi-stage builds](/manuals/build/building/multi-stage.md)
