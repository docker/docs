---
description: Compose command compatibility with docker-compose
keywords: documentation, docs, docker, compose, containers
title: Compose command compatibility with docker-compose
---

The `compose` command in the Docker CLI supports most of the `docker-compose` commands and flags. It is expected to be a drop-in replacement for `docker-compose`. 

If you see any Compose functionality that is not available in the `compose` command, create an issue in the [Compose](https://github.com/docker/compose/issues){:target="_blank" rel="noopener" 
class="_"} GitHub repository, so we can prioritize it.

## Commands or flags not yet implemented

The following commands have not been implemented yet, and may be implemented at a later time.
Let us know if these commands are a higher priority for your use cases.

`compose build --memory`: This option is not yet supported by buildkit. The flag is currently supported, but is hidden to avoid breaking existing Compose usage. It does not have any effect.

## Flags that will not be implemented

The list below includes the flags that we are not planning to support in Compose in the Docker CLI,
either because they are already deprecated in `docker-compose`, or because they are not relevant for Compose in the Docker CLI.

* `compose ps --filter KEY-VALUE` Not relevant due to its complicated usage with the `service` command and also because it is not documented properly in `docker-compose`.
* `compose rm --all` Deprecated in docker-compose.
* `compose scale` Deprecated in docker-compose (use `compose up --scale` instead)

Global flags:

* `--compatibility` has been resignified Docker Compose V2. This now means that in the command running V2 will behave as V1 used to do.
  * One difference is in the word separator on container names. V1 used to use `_` as separator while V2 uses `-` to keep the names more hostname friendly. So when using `--compatibility` Docker 
    Compose should use `_` again. Just make sure to stick to one of them otherwise Docker Compose will not be able to recognize the container as an instance of the service.

## Config command

The config command is intended to show the configuration used by Docker Compose to run the actual project.
As we know, at some parts of the Compose file have a short and a long format. For example, the `ports` entry.
In the example below we can see the config command expanding the `ports` section:

docker-compose.yml:
```
services:
  web:
    image: nginx
    ports:
      - 80:80
```
With `$ docker compose config` the output turns into:
```
services:
  web:
    image: nginx
    networks:
      default: null
    ports:
    - mode: ingress
      target: 80
      published: 80
      protocol: tcp
networks:
  default:
    name: workspace_default
```

The result above is a full size configuration of what will be used by Docker Compose to run the project.
