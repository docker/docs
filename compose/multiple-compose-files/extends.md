---
description: How to use Docker Compose's extends keyword to share configuration between files and projects
keywords: fig, composition, compose, docker, orchestration, documentation, docs
title: Share Compose configurations between files and projects
redirect: 
 - /compose/extends/
---
{% include compose-eol.md %}

Docker Compose's [`extends` attribute](../compose-file/05-services.md#extends) lets you share common configurations
among different files, or even different projects entirely. 

Extending services
is useful if you have several services that reuse a common set of configuration
options. Using `extends` you can define a common set of service options in one
place and refer to it from anywhere.  You can refer to another compose file and select a service you want to also use in your own application, with the ability to override some attributes for your own needs.

When defining any service in `docker-compose.yml`, you can declare that you are
extending another service like this:

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

You can go further and define (or re-define) configuration locally in
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
have a common configuration.  The example below is a Compose app with
two services: a web application and a queue worker. Both services use the same
codebase and share many configuration options.

In a **common.yml** we define the common configuration:

```yaml
services:
  app:
    build: .
    environment:
      CONFIG_FILE_PATH: /code/config
      API_KEY: xxxyyy
    cpu_shares: 5
```

In a **docker-compose.yml** we define the concrete services which use the
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

Another common use case is running one off or administrative tasks against one
or more services in a Compose app. This example demonstrates running a
database backup.

Start with a **docker-compose.yml**.

```yaml
services:
  web:
    image: example/my_web_app:latest
    depends_on:
       db

  db:
    image: postgres:latest
```

In a **docker-compose.admin.yml** add a new service to run the database
export or backup.

```yaml
services:
  dbadmin:
     build: database_admin/
     depends_on:
        - db
```

To start a normal environment run `docker compose up -d`. To run a database
backup, include the `docker-compose.admin.yml` as well.

```console
$ docker compose -f docker-compose.yml -f docker-compose.admin.yml \
  run dbadmin db-backup
```

## Exceptions and limitations

Keep in mind that `volumes_from` and `depends_on` are never shared between
services using `extends`. These exceptions exist to avoid implicit
dependencies; you always define `volumes_from` locally. This ensures
dependencies between services are clearly visible when reading the current file.
Defining these locally also ensures that changes to the referenced file don't
break anything.

That’s a good solution as long as you only need a single service to be shared, and you know about its internal details so you know how to tweak configuration. But this isn’t an acceptable solution when you want to reuse someone else's configuration as a “black box” and don’t know about its own dependencies.

## Reference information

- [`extends`](../compose-file/05-services.md#extends)
