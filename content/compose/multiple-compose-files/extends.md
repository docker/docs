---
description: How to use Docker Compose's extends keyword to share configuration between
  files and projects
keywords: fig, composition, compose, docker, orchestration, documentation, docs
title: Extend your Compose file
aliases:
- /compose/extends/
---

Docker Compose's [`extends` attribute](../compose-file/05-services.md#extends) lets you share common configurations
among different files, or even different projects entirely. 

Extending services
is useful if you have several services that reuse a common set of configuration
options. With `extends` you can define a common set of service options in one
place and refer to it from anywhere. You can refer to another Compose file and select a service you want to also use in your own application, with the ability to override some attributes for your own needs.

> **Important**
>
> When you use multiple Compose files, you must make sure all paths in the
files are relative to the base Compose file. This is required because extend files need not be valid
Compose files. Extend files can contain small fragments of configuration.
Tracking which fragment of a service is relative to which path is difficult and
confusing, so to keep paths easier to understand, all paths must be defined
relative to the base file.
{ .important }

## How it works

When defining any service in your `docker-compose.yml` file, you can declare that you are
extending another service:

```yaml
services:
  web:
    extends:
      file: common-services.yml
      service: webapp
```

This instructs Compose to re-use the configuration for the `webapp` service
defined in the `common-services.yml` file. Suppose that `common-services.yml`
looks like this:

```yaml
services:
  webapp:
    build: .
    ports:
      - "8000:8000"
    volumes:
      - "/data"
```

In this case, you get exactly the same result as if you wrote
`docker-compose.yml` with the same `build`, `ports` and `volumes` configuration
values defined directly under `web`.

You can go further and define, or re-define, configuration locally in
`docker-compose.yml`:

```yaml
services:
  web:
    extends:
      file: common-services.yml
      service: webapp
    environment:
      - DEBUG=1
    cpu_shares: 5

  important_web:
    extends: web
    cpu_shares: 10
```

You can also write other services and link your `web` service to them:

```yaml
services:
  web:
    extends:
      file: common-services.yml
      service: webapp
    environment:
      - DEBUG=1
    cpu_shares: 5
    depends_on:
      - db
  db:
    image: postgres
```

## Further examples

### Example one

Extending an individual service is useful when you have multiple services that
have a common configuration. The example below is a Compose app with
two services, a web application and a queue worker. Both services use the same
codebase and share many configuration options.

The `common.yml` file defines the common configuration:

```yaml
services:
  app:
    build: .
    environment:
      CONFIG_FILE_PATH: /code/config
      API_KEY: xxxyyy
    cpu_shares: 5
```

The `docker-compose.yml` defines the concrete services which use the
common configuration:

```yaml
services:
  webapp:
    extends:
      file: common.yml
      service: app
    command: /code/run_web_app
    ports:
      - 8080:8080
    depends_on:
      - queue
      - db

  queue_worker:
    extends:
      file: common.yml
      service: app
    command: /code/run_worker
    depends_on:
      - queue
```

### Example two

Another common use case for `extends` is running one off or administrative tasks against one
or more services in a Compose app. This example demonstrates running a
database backup.

The `docker-compose.yml` defines the base configuration.

```yaml
services:
  web:
    image: example/my_web_app:latest
    depends_on:
      - db

  db:
    image: postgres:latest
```

`docker-compose.admin.yml` adds a new service to run the database
export or backup.

```yaml
services:
  dbadmin:
    build: database_admin/
    depends_on:
      - db
```

To start a normal environment, run `docker compose up -d`. To run a database
backup, include the `docker-compose.admin.yml` as well.

```console
$ docker compose -f docker-compose.yml -f docker-compose.admin.yml \
  run dbadmin db-backup
```

Compose extends files in the order they're specified on the command line.

## Exceptions and limitations

`volumes_from` and `depends_on` are never shared between
services using `extends`. These exceptions exist to avoid implicit
dependencies; you always define `volumes_from` locally. This ensures
dependencies between services are clearly visible when reading the current file.
Defining these locally also ensures that changes to the referenced file don't
break anything.

`extends` is useful if you only need a single service to be shared and you are 
familiar with the file you're extending to, so you can to tweak the configuration.
But this isn’t an acceptable solution when you want to re-use someone else's
unfamiliar configurations and you don’t know about its own dependencies.

## Reference information

- [`extends`](../compose-file/05-services.md#extends)
