---
title: Using profiles with Compose
desription: Using profiles with Compose
keywords: cli, compose, profile, profiles reference
---

Profiles allow adjusting the Compose application model for various usages and
environments by selectively enabling services.
This is achieved by assigning each service to zero or more profiles. If
unassigned, the service is _always_ started but if assigned, it is only started
if the profile is activated.

This allows one to define additional services in a single `docker-compose.yml` file
that should only be started in specific scenarios, e.g. for debugging or
development tasks.

## Assigning profiles to services

Services are associated with profiles through the
[`profiles` attribute](compose-file/compose-file-v3.md#profiles) which takes an
array of profile names:

```yaml
version: "{{ site.compose_file_v3 }}"
services:
  frontend:
    image: frontend
    profiles: ["frontend"]

  phpmyadmin:
    image: phpmyadmin
    depends_on:
      - db
    profiles:
      - debug

  backend:
    image: backend

  db:
    image: mysql
```

Here the services `frontend` and `phpmyadmin` are assigned to the profiles
`frontend` and `debug` respectively and as such are only started when their
respective profiles are enabled.

Services without a `profiles` attribute will _always_ be enabled, i.e. in this
case running `docker compose up` would only start `backend` and `db`.

Valid profile names follow the regex format of `[a-zA-Z0-9][a-zA-Z0-9_.-]+`.

> **Note**
>
> The core services of your application should not be assigned `profiles` so
> they will always be enabled and automatically started.

## Enabling profiles

To enable a profile supply the `--profile` [command-line option](reference/index.md) or
use the [`COMPOSE_PROFILES` environment variable](reference/envvars.md#compose_profiles):

```sh
$ docker-compose --profile debug up
$ COMPOSE_PROFILES=debug docker-compose up
```

The above command would both start your application with the `debug` profile enabled.
Using the `docker-compose.yml` file above, this would start the services `backend`,
`db` and `phpmyadmin`.

Multiple profiles can be specified by passing multiple `--profile` flags or
a comma-separated list for the `COMPOSE_PROFILES` environment variable:

```sh
$ docker-compose --profile frontend --profile debug up
$ COMPOSE_PROFILES=frontend,debug docker-compose up
```

## Auto-enabling profiles and dependency resolution

When a service with assigned `profiles` is explicitly targeted on the command
line its profiles will be enabled automatically so you don't need to enable them
manually. This can be used for one-off services and debugging tools.
As an example consider this configuration:

```yaml
version: "{{ site.compose_file_v3 }}"
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
# will only start backend and db
$ docker-compose up -d

# this will run db-migrations (and - if necessary - start db)
# by implicitly enabling profile `tools`
$ docker-compose run db-migrations
```

But keep in mind that `docker compose` will only automatically enable the
profiles of the services on the command line and not of any dependencies. This
means that all services the targeted service `depends_on` must have a common
profile with it, be always enabled (by omitting `profiles`) or have a matching
profile enabled explicitly:

```yaml
version: "{{ site.compose_file_v3 }}"
services:
  web:
    image: web

  mock-backend:
    image: backend
    profiles: ["dev"]
    depends_on:
      - db

  db:
    image: mysql
    profiles: ["dev"]

  phpmyadmin:
    image: phpmyadmin
    profiles: ["debug"]
    depends_on:
      - db
```

```sh
# will only start "web"
$ docker compose up -d

# this will start mock-backend (and - if necessary - db)
# by implicitly enabling profile `dev`
$ docker compose up -d mock-backend

# this will fail because profile "dev" is disabled
$ docker compose up phpmyadmin
```

Although targeting `phpmyadmin` will automatically enable its profiles - i.e.
`debug` - it will not automatically enable the profile(s) required by `db` -
i.e. `dev`. To fix this you either have to add the `debug` profile to the `db` service:

```yaml
db:
  image: mysql
  profiles: ["debug", "dev"]
```

or enable a profile of `db` explicitly:

```sh
# profile "debug" is enabled automatically by targeting phpmyadmin
$ docker compose --profile dev up phpmyadmin
$ COMPOSE_PROFILES=dev docker compose up phpmyadmin
```


## Compose documentation

- [User guide](index.md)
- [Installing Compose](install/index.md)
- [Getting Started](gettingstarted.md)
- [Command line reference](reference/index.md)
- [Compose file reference](compose-file/index.md)
- [Sample apps with Compose](samples-for-compose.md)
