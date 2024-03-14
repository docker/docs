---
title: Use containers for Python development
keywords: python, local, development
description: Learn how to develop your Python application locally.
---

## Prerequisites

Complete [Containerize a Python application](containerize.md).

## Overview

In this section, you'll learn how to set up a development environment for your containerized application. This includes:

- Adding a local database and persisting data
- Configuring Compose to automatically update your running Compose services as you edit and save your code

## Get the sample application

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
     - README.Docker.md

   Let's get started!

   ? What application platform does your project use? Python
   ? What version of Python do you want to use? 3.11.4
   ? What port do you want your app to listen on? 5000
   ? What is the command to run your app? python3 -m flask run --host=0.0.0.0
   ```

## Add a local database and persist data

You can use containers to set up local services, like a database. In this section, you'll update the `compose.yaml` file to define a database service and a volume to persist data.

In the cloned repository's directory, open the `compose.yaml` file in an IDE or text editor. `docker init` handled creating most of the instructions, but you'll need to update it for your unique application.

In the `compose.yaml` file, you need to uncomment all of the database instructions. In addition, you need to add the database password file as an environment variable to the server service and specify the secret file to use .

The following is the updated `compose.yaml` file.

```yaml {hl_lines="7-36"}
services:
  server:
    build:
      context: .
    ports:
      - 5000:5000
    environment:
      - POSTGRES_PASSWORD_FILE=/run/secrets/db-password
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

Before you run the application using Compose, notice that this Compose file specifies a `password.txt` file to hold the database's password. You must create this file as it's not included in the source repository.

In the cloned repository's directory, create a new directory named `db` and inside that directory create a file named `password.txt` that contains the password for the database. Using your favorite IDE or text editor, add the following contents to the `password.txt` file.

```text
mysecretpassword
```

Save and close the `password.txt` file.

You should now have the following contents in your `python-docker-dev`
directory.

```text
├── python-docker-dev/
│ ├── db/
│ │ └── password.txt
│ ├── app.py
│ ├── requirements.txt
│ ├── .dockerignore
│ ├── compose.yaml
│ ├── Dockerfile
│ ├── README.Docker.md
│ └── README.md
```

Now, run the following `docker compose up` command to start your application.

```console
$ docker compose up --build
```

Now test your API endpoint. Open a new terminal then make a request to the server using the curl commands:

```console
$ curl http://localhost:5000/initdb
$ curl http://localhost:5000/widgets
```

You should receive the following response:

```json
[]
```

The response is empty because your database is empty.

Press `ctrl+c` in the terminal to stop your application.

## Automatically update services

Use Compose Watch to automatically update your running Compose services as you
edit and save your code. For more details about Compose Watch, see [Use Compose
Watch](../../compose/file-watch.md).

Open your `compose.yaml` file in an IDE or text editor and then add the Compose
Watch instructions. The following is the updated `compose.yaml` file.

```yaml {hl_lines="14-17"}
services:
  server:
    build:
      context: .
    ports:
      - 5000:5000
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

Run the following command to run your application with Compose Watch.

```console
$ docker compose watch
```

In a terminal, curl the application to get a response.

```console
$ curl http://localhost:5000
Hello, Docker!
```

Any changes to the application's source files on your local machine will now be immediately reflected in the running container.

Open `python-docker-dev/app.py` in an IDE or text editor and update the `Hello, Docker!` string by adding a few more exclamation marks.

```diff
-    return 'Hello, Docker!'
+    return 'Hello, Docker!!!'
```

Save the changes to `app.py` and then wait a few seconds for the application to rebuild. Curl the application again and verify that the updated text appears.

```console
$ curl http://localhost:5000
Hello, Docker!!!
```

Press `ctrl+c` in the terminal to stop your application.

## Summary

In this section, you took a look at setting up your Compose file to add a local
database and persist data. You also learned how to use Compose Watch to automatically rebuild and run your container when you update your code.

Related information:
 - [Compose file reference](/compose/compose-file/)
 - [Compose file watch](../../compose/file-watch.md)
 - [Multi-stage builds](../../build/building/multi-stage.md)

## Next steps

In the next section, you'll take a look at how to set up a CI/CD pipeline using GitHub Actions.

{{< button text="Configure CI/CD" url="configure-ci-cd.md" >}}
