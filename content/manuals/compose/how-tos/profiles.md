---
title: Using profiles with Compose
linkTitle: Use service profiles
weight: 20
description: How to use profiles with Docker Compose
keywords: cli, compose, profile, profiles reference
aliases:
- /compose/profiles/
---

{{% include "compose/profiles.md" %}}

## Assigning profiles to services

Services are associated with profiles through the
[`profiles` attribute](/reference/compose-file/services.md#profiles) which takes an
array of profile names:

```yaml
services:
  frontend:
    image: frontend
    profiles: [frontend]

  phpmyadmin:
    image: phpmyadmin
    depends_on: [db]
    profiles: [debug]

  backend:
    image: backend

  db:
    image: mysql
```

Here the services `frontend` and `phpmyadmin` are assigned to the profiles
`frontend` and `debug` respectively and as such are only started when their
respective profiles are enabled.

Services without a `profiles` attribute are always enabled. In this
case running `docker compose up` would only start `backend` and `db`.

Valid profiles names follow the regex format of `[a-zA-Z0-9][a-zA-Z0-9_.-]+`.

> [!TIP]
>
> The core services of your application shouldn't be assigned `profiles` so
> they are always enabled and automatically started.

## Start specific profiles

To start a specific profile supply the `--profile` [command-line option](/reference/cli/docker/compose.md) or
use the [`COMPOSE_PROFILES` environment variable](environment-variables/envvars.md#compose_profiles):

```console
$ docker compose --profile debug up
```
```console
$ COMPOSE_PROFILES=debug docker compose up
```

Both commands start the services with the `debug` profile enabled.
In the previous `compose.yaml` file, this starts the services
`db`, `backend` and `phpmyadmin`.

### Start multiple profiles

You can also enable
multiple profiles, e.g. with `docker compose --profile frontend --profile debug up`
the profiles `frontend` and `debug` will be enabled.

Multiple profiles can be specified by passing multiple `--profile` flags or
a comma-separated list for the `COMPOSE_PROFILES` environment variable:

```console
$ docker compose --profile frontend --profile debug up
```

```console
$ COMPOSE_PROFILES=frontend,debug docker compose up
```

If you want to enable all profiles at the same time, you can run `docker compose --profile "*"`.

## Auto-starting profiles and dependency resolution

When you explicitly target a service on the command line that has one or more profiles assigned, you do not need to enable the profile manually as Compose runs that service regardless of whether its profile is activated. This is useful for running one-off services or debugging tools.

Only the targeted service (and any of its declared dependencies via `depends_on`) is started. Other services that share the same profile will not be started unless:
- They are also explicitly targeted, or
- The profile is explicitly enabled using `--profile` or `COMPOSE_PROFILES`.

When a service with assigned `profiles` is explicitly targeted on the command
line its profiles are started automatically so you don't need to start them
manually. This can be used for one-off services and debugging tools.
As an example consider the following configuration:

```yaml
services:
  backend:
    image: backend

  db:
    image: mysql

  db-migrations:
    image: backend
    command: myapp migrate
    depends_on:
      - db
    profiles:
      - tools
```

```sh
# Only start backend and db (no profiles involved)
$ docker compose up -d

# Run the db-migrations service without manually enabling the 'tools' profile
$ docker compose run db-migrations
```

In this example, `db-migrations` runs even though it is assigned to the tools profile, because it was explicitly targeted. The `db` service is also started automatically because it is listed in `depends_on`.

If the targeted service has dependencies that are also gated behind a profile, you must ensure those dependencies are either: 
 - In the same profile
 - Started separately
 - Not assigned to any profile so are always enabled

## Stop application and services with specific profiles

As with starting specific profiles, you can use the `--profile` [command-line option](/reference/cli/docker/compose.md#use--p-to-specify-a-project-name) or
use the [`COMPOSE_PROFILES` environment variable](environment-variables/envvars.md#compose_profiles):

```console
$ docker compose --profile debug down
```
```console
$ COMPOSE_PROFILES=debug docker compose down
```

Both commands stop and remove services with the `debug` profile and services without a profile. In the following `compose.yaml` file, this stops the services `db`, `backend` and `phpmyadmin`.

```yaml
services:
  frontend:
    image: frontend
    profiles: [frontend]

  phpmyadmin:
    image: phpmyadmin
    depends_on: [db]
    profiles: [debug]

  backend:
    image: backend

  db:
    image: mysql
```

if you only want to stop the `phpmyadmin` service, you can run 

```console 
$ docker compose down phpmyadmin
``` 
or 
```console 
$ docker compose stop phpmyadmin
```

> [!NOTE]
>
> Running `docker compose down` only stops `backend` and `db`.

## Reference information

[`profiles`](/reference/compose-file/services.md#profiles)
