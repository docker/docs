---
description: Pulls service images.
keywords: fig, composition, compose, docker, orchestration, cli,  pull
title: docker-compose pull
notoc: true
---

```
Usage: pull [options] [SERVICE...]

Options:
    --ignore-pull-failures  Pull what it can and ignores images with pull failures.
    --parallel              Deprecated, pull multiple images in parallel (enabled by default).
    --no-parallel           Disable parallel pulling.
    -q, --quiet             Pull without printing progress information
    --include-deps          Also pull services declared as dependencies
```

Pulls an image associated with a service defined in a `docker-compose.yml` or `docker-stack.yml` file, but does not start containers based on those images.

For example, suppose you have this `docker-compose.yml` file from the [Quickstart: Compose and Rails](../rails.md) sample.

```yaml
version: '2'
services:
  db:
    image: postgres
  web:
    build: .
    command: bundle exec rails s -p 3000 -b '0.0.0.0'
    volumes:
      - .:/myapp
    ports:
      - "3000:3000"
    depends_on:
      - db
```

If you run `docker-compose pull ServiceName` in the same directory as the `docker-compose.yml` file that defines the service, Docker pulls the associated image. For example, to call the `postgres` image configured as the `db` service in our example, you would run `docker-compose pull db`.

```bash
$ docker-compose pull db
Pulling db (postgres:latest)...
latest: Pulling from library/postgres
cd0a524342ef: Pull complete
9c784d04dcb0: Pull complete
d99dddf7e662: Pull complete
e5bff71e3ce6: Pull complete
cb3e0a865488: Pull complete
31295d654cd5: Pull complete
fc930a4e09f5: Pull complete
8650cce8ef01: Pull complete
61949acd8e52: Pull complete
527a203588c0: Pull complete
26dec14ac775: Pull complete
0efc0ed5a9e5: Pull complete
40cd26695b38: Pull complete
Digest: sha256:fd6c0e2a9d053bebb294bb13765b3e01be7817bf77b01d58c2377ff27a4a46dc
Status: Downloaded newer image for postgres:latest
```
