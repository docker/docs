---
description: How to use Docker Compose's extends keyword to share configuration between files and projects
keywords: fig, composition, compose, docker, orchestration, documentation, docs
title: Share Compose configurations between files and projects
---

Compose supports two methods of sharing common configuration:

1. Extending an entire Compose file by
   [using multiple Compose files](extends.md#multiple-compose-files)
2. Extending individual services with [the `extends` field](extends.md#extending-services) (for Compose file versions up to 2.1)


## Multiple Compose files

Using multiple Compose files enables you to customize a Compose application
for different environments or different workflows.

### Understanding multiple Compose files

By default, Compose reads two files, a `docker-compose.yml` and an optional
`docker-compose.override.yml` file. By convention, the `docker-compose.yml`
contains your base configuration. The override file, as its name implies, can
contain configuration overrides for existing services or entirely new
services.

If a service is defined in both files, Compose merges the configurations using
the rules described in
[Adding and overriding configuration](extends.md#adding-and-overriding-configuration).

To use multiple override files, or an override file with a different name, you
can use the `-f` option to specify the list of files. Compose merges files in
the order they're specified on the command line. See the
[`docker-compose` command reference](reference/index.md) for more information
about using `-f`.

When you use multiple configuration files, you must make sure all paths in the
files are relative to the base Compose file (the first Compose file specified
with `-f`). This is required because override files need not be valid
Compose files. Override files can contain small fragments of configuration.
Tracking which fragment of a service is relative to which path is difficult and
confusing, so to keep paths easier to understand, all paths must be defined
relative to the base file.

### Example use case

In this section, there are two common use cases for multiple Compose files: changing a
Compose app for different environments, and running administrative tasks
against a Compose app.

#### Different environments

A common use case for multiple files is changing a development Compose app
for a production-like environment (which may be production, staging or CI).
To support these differences, you can split your Compose configuration into
a few different files:

Start with a base file that defines the canonical configuration for the
services.

**docker-compose.yml**

```yaml
web:
  image: example/my_web_app:latest
  depends_on:
    - db
    - cache

db:
  image: postgres:latest

cache:
  image: redis:latest
```

In this example the development configuration exposes some ports to the
host, mounts our code as a volume, and builds the web image.

**docker-compose.override.yml**

```yaml
web:
  build: .
  volumes:
    - '.:/code'
  ports:
    - 8883:80
  environment:
    DEBUG: 'true'

db:
  command: '-d'
  ports:
    - 5432:5432

cache:
  ports:
    - 6379:6379
```

When you run `docker-compose up` it reads the overrides automatically.

Now, it would be nice to use this Compose app in a production environment. So,
create another override file (which might be stored in a different git
repo or managed by a different team).

**docker-compose.prod.yml**

```yaml
web:
  ports:
    - 80:80
  environment:
    PRODUCTION: 'true'

cache:
  environment:
    TTL: '500'
```

To deploy with this production Compose file you can run

```console
$ docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

This deploys all three services using the configuration in
`docker-compose.yml` and `docker-compose.prod.yml` (but not the
dev configuration in `docker-compose.override.yml`).


See [production](production.md) for more information about Compose in
production.

#### Administrative tasks

Another common use case is running adhoc or administrative tasks against one
or more services in a Compose app. This example demonstrates running a
database backup.

Start with a **docker-compose.yml**.

```yaml
web:
  image: example/my_web_app:latest
  depends_on:
    - db

db:
  image: postgres:latest
```

In a **docker-compose.admin.yml** add a new service to run the database
export or backup.

```yaml
    dbadmin:
      build: database_admin/
      depends_on:
        - db
```

To start a normal environment run `docker-compose up -d`. To run a database
backup, include the `docker-compose.admin.yml` as well.

```console
$ docker-compose -f docker-compose.yml -f docker-compose.admin.yml \
  run dbadmin db-backup
```

## Extending services

> **Note**
>
> The `extends` keyword is supported in earlier Compose file formats up to Compose
> file version 2.1 (see [extends in v2](compose-file/compose-file-v2.md#extends)), but is
> not supported in Compose version 3.x. See the [Version 3 summary](compose-file/compose-versioning.md#version-3)
> of keys added and removed, along with information on [how to upgrade](compose-file/compose-versioning.md#upgrading).
> See [moby/moby#31101](https://github.com/moby/moby/issues/31101) to follow the
> discussion thread on the possibility of adding support for `extends` in some form in
> future versions. The `extends` keyword has been included in docker-compose versions 1.27
> and higher.

Docker Compose's `extends` keyword enables the sharing of common configurations
among different files, or even different projects entirely. Extending services
is useful if you have several services that reuse a common set of configuration
options. Using `extends` you can define a common set of service options in one
place and refer to it from anywhere.

Keep in mind that `volumes_from` and `depends_on` are never shared between
services using `extends`. These exceptions exist to avoid implicit
dependencies; you always define `volumes_from` locally. This ensures
dependencies between services are clearly visible when reading the current file.
Defining these locally also ensures that changes to the referenced file don't
break anything.

### Understand the extends configuration

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

### Example use case

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

## Adding and overriding configuration

Compose copies configurations from the original service over to the local one.
If a configuration option is defined in both the original service and the local
service, the local value *replaces* or *extends* the original value.

For single-value options like `image`, `command` or `mem_limit`, the new value
replaces the old value.

original service:

```yaml
services:
  myservice:
    # ...
    command: python app.py
```

local service:

```yaml
services:
  myservice:
    # ...
    command: python otherapp.py
```

result:

```yaml
services:
  myservice:
    # ...
    command: python otherapp.py
```

For the **multi-value options** `ports`, `expose`, `external_links`, `dns`,
`dns_search`, and `tmpfs`, Compose concatenates both sets of values:

original service:

```yaml
services:
  myservice:
    # ...
    expose:
      - "3000"
```

local service:

```yaml
services:
  myservice:
    # ...
    expose:
      - "4000"
      - "5000"
```

result:

```yaml
services:
  myservice:
    # ...
    expose:
      - "3000"
      - "4000"
      - "5000"
```

In the case of `environment`, `labels`, `volumes`, and `devices`, Compose
"merges" entries together with locally-defined values taking precedence. For
`environment` and `labels`, the environment variable or label name determines
which value is used:

original service:

```yaml
services:
  myservice:
    # ...
    environment:
      - FOO=original
      - BAR=original
```

local service:

```yaml
services:
  myservice:
    # ...
    environment:
      - BAR=local
      - BAZ=local
```

result

```yaml
services:
  myservice:
    # ...
    environment:
      - FOO=original
      - BAR=local
      - BAZ=local
```

Entries for `volumes` and `devices` are merged using the mount path in the
container:

original service:

```yaml
services:
  myservice:
    # ...
    volumes:
      - ./original:/foo
      - ./original:/bar
```

local service:

```yaml
services:
  myservice:
    # ...
    volumes:
      - ./local:/bar
      - ./local:/baz
```

result:

```yaml
services:
  myservice:
    # ...
    volumes:
      - ./original:/foo
      - ./local:/bar
      - ./local:/baz
```


## Compose documentation

- [User guide](index.md)
- [Installing Compose](install.md)
- [Getting Started](gettingstarted.md)
- [Command line reference](reference/index.md)
- [Compose file reference](compose-file/index.md)
- [Sample apps with Compose](samples-for-compose.md)
