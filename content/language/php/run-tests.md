---
title: Run PHP tests in a container
keywords: php, test
description: Learn how to run your PHP tests in a container.
---

## Prerequisites

Complete all the previous sections of this guide, starting with [Containerize a PHP application](containerize.md).

## Overview

Testing is an essential part of modern software development. Testing can mean a
lot of things to different development teams. There are unit tests, integration
tests and end-to-end testing. In this guide you take a look at running your unit
tests in Docker when developing and when building.

## Run tests when developing locally

The sample application already has a PHPUnit test inside the `tests` directory. When developing locally, you can use Compose to run your tests.

Run the following command in the `docker-php-sample` directory to run the tests inside a container.

```console
$ docker compose run --build --rm server ./vendor/bin/phpunit tests/HelloWorldTest.php
```

You should see output that contains the following.

```console
Hello, Docker!PHPUnit 9.6.13 by Sebastian Bergmann and contributors.

.                                                                   1 / 1 (100%)

Time: 00:00.003, Memory: 4.00 MB

OK (1 test, 1 assertion)
```

To learn more about the command, see [docker compose run](/reference/cli/docker/compose/run/).

## Run tests when building

To run your tests when building, you need to update your Dockerfile. Create a new test stage that runs the tests.

The following is the updated Dockerfile.

```dockerfile {hl_lines="26-28"}
# syntax=docker/dockerfile:1

FROM composer:lts as prod-deps
WORKDIR /app
RUN --mount=type=bind,source=./composer.json,target=composer.json \
    --mount=type=bind,source=./composer.lock,target=composer.lock \
    --mount=type=cache,target=/tmp/cache \
    composer install --no-dev --no-interaction

FROM composer:lts as dev-deps
WORKDIR /app
RUN --mount=type=bind,source=./composer.json,target=composer.json \
    --mount=type=bind,source=./composer.lock,target=composer.lock \
    --mount=type=cache,target=/tmp/cache \
    composer install --no-interaction

FROM php:8.2-apache as base
RUN docker-php-ext-install pdo pdo_mysql
COPY ./src /var/www/html

FROM base as development
COPY ./tests /var/www/html/tests
RUN mv "$PHP_INI_DIR/php.ini-development" "$PHP_INI_DIR/php.ini"
COPY --from=dev-deps app/vendor/ /var/www/html/vendor

FROM development as test
WORKDIR /var/www/html
RUN ./vendor/bin/phpunit tests/HelloWorldTest.php

FROM base as final
RUN mv "$PHP_INI_DIR/php.ini-production" "$PHP_INI_DIR/php.ini"
COPY --from=prod-deps app/vendor/ /var/www/html/vendor
USER www-data
```

Run the following command to build an image using the test stage as the target and view the test results. Include `--progress plain` to view the build output, `--no-cache` to ensure the tests always run, and `--target test` to target the test stage.

```console
$ docker build -t php-docker-image-test --progress plain --no-cache --target test .
```

You should see output containing the following.

```console
#18 [test 2/2] RUN ./vendor/bin/phpunit tests/HelloWorldTest.php
#18 0.385 Hello, Docker!PHPUnit 9.6.13 by Sebastian Bergmann and contributors.
#18 0.392
#18 0.394 .                                                                   1 / 1 (100%)
#18 0.395
#18 0.395 Time: 00:00.003, Memory: 4.00 MB
#18 0.395
#18 0.395 OK (1 test, 1 assertion)
```

To learn more about building and running tests, see the [Build with Docker guide](../../build/guide/_index.md).

## Summary

In this section, you learned how to run tests when developing locally using Compose and how to run tests when building your image.

Related information:
 - [docker compose run](/reference/cli/docker/compose/run/)
 - [Build with Docker guide](../../build/guide/index.md)

## Next steps

Next, you’ll learn how to set up a CI/CD pipeline using GitHub Actions.

{{< button text="Configure CI/CD" url="configure-ci-cd.md" >}}