---
description: How to set, use and manage environment variables in Compose
keywords: fig, composition, compose, docker, orchestration, environment, variables, env file
title: Environment variables in Compose
---

There are multiple parts of Compose that deal with environment variables in one sense or another. This page should help you find the information you need.


## Substituting environment variables in Compose files

It's possible to use environment variables in your shell to populate values inside a Compose file:

    web:
      image: "webapp:${TAG}"

For more information, see the [Variable substitution](compose-file.md#variable-substitution) section in the Compose file reference.


## Setting environment variables in containers

You can set environment variables in a service's containers with the ['environment' key](compose-file.md#environment), just like with `docker run -e VARIABLE=VALUE ...`:

    web:
      environment:
        - DEBUG=1


## Passing environment variables through to containers

You can pass environment variables from your shell straight through to a service's containers with the ['environment' key](compose-file.md#environment) by not giving them a value, just like with `docker run -e VARIABLE ...`:

    web:
      environment:
        - DEBUG

The value of the `DEBUG` variable in the container will be taken from the value for the same variable in the shell in which Compose is run.


## The “env_file” configuration option

You can pass multiple environment variables from an external file through to a service's containers with the ['env_file' option](compose-file.md#envfile), just like with `docker run --env-file=FILE ...`:

    web:
      env_file:
        - web-variables.env


## Setting environment variables with 'docker-compose run'

Just like with `docker run -e`, you can set environment variables on a one-off container with `docker-compose run -e`:

    docker-compose run -e DEBUG=1 web python console.py

You can also pass a variable through from the shell by not giving it a value:

    docker-compose run -e DEBUG web python console.py

The value of the `DEBUG` variable in the container will be taken from the value for the same variable in the shell in which Compose is run.


## The “.env” file

You can set default values for any environment variables referenced in the Compose file, or used to configure Compose, in an [environment file](env-file.md) named `.env`:

    $ cat .env
    TAG=v1.5

    $ cat docker-compose.yml
    version: '3'
    services:
      web:
        image: "webapp:${TAG}"

When you run `docker-compose up`, the `web` service defined above uses the image `webapp:v1.5`. You can verify this with the [config command](reference/config.md), which prints your resolved application config to the terminal:

    $ docker-compose config
    version: '3'
    services:
      web:
        image: 'webapp:v1.5'

Values in the shell take precedence over those specified in the `.env` file. If you set `TAG` to a different value in your shell, the substitution in `image` uses that instead:

    $ export TAG=v2.0
    $ docker-compose config
    version: '3'
    services:
      web:
        image: 'webapp:v2.0'
   
When values are provided with both with shell `environment` variable and with an `env_file` configuration file, values of environment variables will be taken **from environment key first and then from environment file, then from a `Dockerfile` `ENV`–entry**:

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

You can test this with for e.g. a _NodeJS_ container in the CLI:

    $ docker-compose exec api node
    > process.env.NODE_ENV
    'production'

Having any `ARG` or `ENV` setting in a `Dockerfile` will evaluate only if there is _no_ Docker _Compose_ entry for `environment` or `env_file`.

_Spcecifics for NodeJS containers:_ If you have a `package.json` entry for `script:start` like `NODE_ENV=test node server.js`, then this will overrule _any_ setting in your `docker-compose.yml` file.

## Configuring Compose using environment variables

Several environment variables are available for you to configure the Docker Compose command-line behaviour. They begin with `COMPOSE_` or `DOCKER_`, and are documented in [CLI Environment Variables](reference/envvars.md).

## Environment variables created by links

When using the ['links' option](compose-file.md#links) in a [v1 Compose file](compose-file.md#version-1), environment variables will be created for each link. They are documented in the [Link environment variables reference](link-env-deprecated.md). Please note, however, that these variables are deprecated - you should just use the link alias as a hostname instead.
