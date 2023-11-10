---
title: Use containers for PHP development
keywords: php, development
description: Learn how to develop your PHP application locally using containers.
---

## Prerequisites

Complete [Containerize a PHP application](containerize.md).

## Overview

In this section, you'll learn how to set up a development environment for your containerized application. This includes:
 - Adding a local database and persisting data
 - Configuring Compose to automatically update your running Compose services as you edit and save your code
 - Creating a development container that contains the dev dependencies

## Add a local database and persist data

You can use containers to set up local services, like a database.
To do this for the sample application, you'll need to do the following:
- Update the `Dockerfile` to install extensions to connect to the database
- Update the `compose.yaml` file to add a database service and volume to persist data

### Update the Dockerfile to install extensions

To install PHP extensions, you need to update the `Dockerfile`. The following is the updated `Dockerfile` that installs the `pdo` and `pdo_pgsql` extensions.

```dockerfile
# syntax=docker/dockerfile:1

FROM composer:lts as deps

WORKDIR /app

RUN --mount=type=bind,source=composer.json,target=composer.json \
    --mount=type=bind,source=composer.lock,target=composer.lock \
    --mount=type=cache,target=/tmp/cache \
    composer install --no-dev --no-interaction

FROM php:8.2-apache as final

RUN apt-get update && apt-get install -y libpq-dev && docker-php-ext-install pdo pdo_pgsql
RUN mv "$PHP_INI_DIR/php.ini-production" "$PHP_INI_DIR/php.ini"

COPY --from=deps app/vendor/ /var/www/html/vendor
COPY ./src /var/www/html

USER www-data
```

###  Update the compose.yaml file to add a db and persist data

Open the `compose.yaml` file in an IDE or text editor. You'll notice it
already contains commented-out instructions for a PostgreSQL database and volume.

Open the `src/database.php` file in an IDE or text editor. You'll notice that it reads environment variables in order to connect to the database.

In the `compose.yaml`` file, you'll need to do update the following:
1. Uncomment the database instructions.
2. Add a secret to the server service to pass in the database password.
3. Add the database connection environment variables to the server service.
4. Uncomment the volume instructions to persist data.

The following is the updated `compose.yaml` file.

```yaml
services:
  server:
    build:
      context: .
    ports:
      - 9000:80
    depends_on:
      db:
        condition: service_healthy
    secrets:
      - db-password
    environment:
      - PASSWORD_FILE_PATH=/run/secrets/db-password
      - DB_HOST=db
      - DB_NAME=example
      - DB_USER=postgres
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

In the `docker-php-sample` directory, create a new directory named `db` and
inside that directory create a file named `password.txt`. Open `password.txt` in an IDE or text editor and add the following password. The password must be on a single line, with no additional lines in the file.

```
example
```

Save and close the `password.txt` file.

You should now have the following in your `docker-php-sample` directory.

```
├── docker-php-sample/
│ ├── .git/
│ ├── db/
│ │ └── password.txt
│ ├── src/
│ ├── tests/
│ ├── .dockerignore
│ ├── .gitignore
│ ├── compose.yaml
│ ├── Dockerfile
│ ├── README.Docker.md
│ └── README.md
```

Run the following command to start your application.

```console
$ docker compose up --build
```

Open a browser and view the application at [http://localhost:9000/database.php](http://localhost:9000/database.php). You should see a simple web application with text and a counter that increments every time you refresh.

Press `ctrl+c` in the terminal to stop your application.

## Verify that data persists in the database

In the terminal, run `docker compose rm` to remove your containers and then run `docker compose up` to run your application again.

```console
$ docker compose rm
$ docker compose up --build
```

Refresh [http://localhost:9000/database.php](http://localhost:9000/database.php) in your browser and verify that the previous count still exists.

Press `ctrl+c` in the terminal to stop your application.

## Automatically update services

Use Compose Watch to automatically update your running Compose services as you edit and save your code. For more details about Compose Watch, see [Use Compose Watch](../../compose/file-watch.md).

Open your `compose.yaml` file in an IDE or text editor and then add the Compose Watch instructions. The following is the updated `compose.yaml` file.

```yaml
services:
  server:
    build:
      context: .
    ports:
      - 9000:80
    depends_on:
      db:
        condition: service_healthy
    secrets:
      - db-password
    environment:
      - PASSWORD_FILE_PATH=/run/secrets/db-password
      - DB_HOST=db
      - DB_NAME=example
      - DB_USER=postgres
    develop:
      watch:
        - action: sync
          path: ./src
          target: /var/www/html
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

Open a browser and verify that the application is running at [http://localhost:9000/hello.php](http://localhost:9000/hello.php).

Any changes to the application's source files on your local machine will now be
immediately reflected in the running container.

Open `hello.php` in an IDE or text editor and update the string `Hello, world!` to `Hello, Docker!`.

Save the changes to `hello.php` and then wait a few seconds for the application to sync. Refresh [http://localhost:9000/hello.php](http://localhost:9000/hello.php) in your browser and verify that the updated text appears.

Press `ctrl+c` in the terminal to stop Compose Watch. Run `docker compose down` in the terminal to stop the application.

## Create a development container

At this point, when you run your containerized application, Composer isn't installing the dev dependencies. While this small image is good for production, it lacks the tools and dependencies you may need when developing and it doesn't include the `tests` directory. You can use multi-stage builds to build stages for both development and production in the same Dockerfile. For more details, see [Multi-stage builds](../../build/building/multi-stage.md).

While you probably don't need a multi-stage build to optimize a development image, the example below uses it so that you can see how both the development stages and the production stages can be built while sharing some common stages.


The following is the updated Dockerfile.

```dockerfile
# syntax=docker/dockerfile:1

FROM composer:lts as deps

WORKDIR /app

FROM deps as prod-deps
RUN --mount=type=bind,source=./composer.json,target=composer.json \
    --mount=type=bind,source=./composer.lock,target=composer.lock \
    --mount=type=cache,target=/tmp/cache \
    composer install --no-dev --no-interaction

FROM deps as dev-deps
RUN --mount=type=bind,source=./composer.json,target=composer.json \
    --mount=type=bind,source=./composer.lock,target=composer.lock \
    --mount=type=cache,target=/tmp/cache \
    composer install --no-interaction

FROM php:8.2-apache as base
RUN apt-get update && apt-get install -y libpq-dev && docker-php-ext-install pdo pdo_pgsql
COPY ./src /var/www/html

FROM base as development
COPY ./tests /var/www/html/tests
RUN mv "$PHP_INI_DIR/php.ini-development" "$PHP_INI_DIR/php.ini"
COPY --from=dev-deps app/vendor/ /var/www/html/vendor
USER www-data

FROM base as final
RUN mv "$PHP_INI_DIR/php.ini-production" "$PHP_INI_DIR/php.ini"
COPY --from=prod-deps app/vendor/ /var/www/html/vendor
USER www-data
```

Update your `compose.yaml` file to target the development stage.

The following is the updated `compose.yaml` file.

```yaml
services:
  server:
    build:
      context: .
      target: development
    ports:
      - 9000:80
    depends_on:
      db:
        condition: service_healthy
    secrets:
      - db-password
    environment:
      - PASSWORD_FILE_PATH=/run/secrets/db-password
      - DB_HOST=db
      - DB_NAME=example
      - DB_USER=postgres
    develop:
      watch:
        - action: sync
          path: ./src
          target: /var/www/html
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

Your containerized application will now install the dev dependencies. Continue to the next section to learn how you can run tests.

## Summary

In this section, you took a look at setting up your Compose file to add a local
database and persist data. You also learned how to use Compose Watch to automatically sync your application when you update your code. And finally, you learned how to create a development container that contains the dependencies needed for development.

Related information:
 - [Compose file reference](/compose/compose-file/)
 - [Compose file watch](../../compose/file-watch.md)
 - [Multi-stage builds](../../build/building/multi-stage.md)

## Next steps

In the next section, you'll learn how to run unit tests using Docker.

{{< button text="Run your tests" url="run-tests.md" >}}