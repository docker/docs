---
title: Environment variables in Compose
description: How to set, use and manage environment variables in Compose
keywords: compose, orchestration, environment, env file
---

There are multiple parts of Compose that deal with environment variables in one
sense or another. This page should help you find the information you need.


## Substitute environment variables in Compose files

It's possible to use environment variables in your shell to populate values
inside a Compose file:

```yaml
web:
  image: "webapp:${TAG}"
```

For more information, see the
[Variable substitution](compose-file/index.md#variable-substitution) section in the
Compose file reference.


## Set environment variables in containers

You can set environment variables in a service's containers with the
['environment' key](compose-file/index.md#environment), just like with
`docker run -e VARIABLE=VALUE ...`:

```yaml
web:
  environment:
    - DEBUG=1
```

## Pass environment variables to containers

You can pass environment variables from your shell straight through to a
service's containers with the ['environment' key](compose-file/index.md#environment)
by not giving them a value, just like with `docker run -e VARIABLE ...`:

```yaml
web:
  environment:
    - DEBUG
```

The value of the `DEBUG` variable in the container is taken from the value for
the same variable in the shell in which Compose is run.

## The “env_file” configuration option

You can pass multiple environment variables from an external file through to
a service's containers with the ['env_file' option](compose-file/index.md#env_file),
just like with `docker run --env-file=FILE ...`:

```yaml
web:
  env_file:
    - web-variables.env
```

## Set environment variables with 'docker-compose run'

Just like with `docker run -e`, you can set environment variables on a one-off
container with `docker-compose run -e`:

```bash
docker-compose run -e DEBUG=1 web python console.py
```

You can also pass a variable through from the shell by not giving it a value:

```bash
docker-compose run -e DEBUG web python console.py
```

The value of the `DEBUG` variable in the container is taken from the value for
the same variable in the shell in which Compose is run.


## The “.env” file

You can set default values for any environment variables referenced in the
Compose file, or used to configure Compose, in an [environment file](env-file.md)
named `.env`:

```bash
$ cat .env
TAG=v1.5

$ cat docker-compose.yml
version: '3'
services:
  web:
    image: "webapp:${TAG}"
```

When you run `docker-compose up`, the `web` service defined above uses the
image `webapp:v1.5`. You can verify this with the
[config command](reference/config.md), which prints your resolved application
config to the terminal:

```bash
$ docker-compose config

version: '3'
services:
  web:
    image: 'webapp:v1.5'
```

Values in the shell take precedence over those specified in the `.env` file.
If you set `TAG` to a different value in your shell, the substitution in `image`
uses that instead:

```bash
$ export TAG=v2.0
$ docker-compose config

version: '3'
services:
  web:
    image: 'webapp:v2.0'
```

When you set the same environment variable in multiple files, here's the
priority used by Compose to choose which value to use:

1. Compose file
2. Shell environment variables
3. Environment file
4. Dockerfile
5. Variable is not defined

In the example below, we set the same environment variable on an Environment
file, and the Compose file:

```bash
$ cat ./Docker/api/api.env
NODE_ENV=test

$ cat docker-compose.yml
version: '3'
services:
  api:
    image: 'node:6-alpine'
    env_file:
     - ./Docker/api/api.env
    environment:
     - NODE_ENV=production
```

When you run the container, the environment variable defined in the Compose
file takes precedence.

```bash
$ docker-compose exec api node

> process.env.NODE_ENV
'production'
```

Having any `ARG` or `ENV` setting in a `Dockerfile` evaluates only if there is
no Docker Compose entry for `environment` or `env_file`.

> Specifics for NodeJS containers
>
> If you have a `package.json` entry for `script:start` like
> `NODE_ENV=test node server.js`, then this overrules any setting in your
> `docker-compose.yml` file.

## Configure Compose using environment variables

Several environment variables are available for you to configure the Docker
Compose command-line behavior. They begin with `COMPOSE_` or `DOCKER_`, and are
documented in [CLI Environment Variables](reference/envvars.md).

## Environment variables created by links

When using the ['links' option](compose-file/index.md#links) in a
[v1 Compose file](compose-file/index.md#version-1), environment variables are created
for each link. They are documented in
the [Link environment variables reference](link-env-deprecated.md).

However, these variables are deprecated. Use the link alias as a hostname instead.
