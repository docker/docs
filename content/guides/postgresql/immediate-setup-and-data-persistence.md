---
title: Immediate Setup & Data Persistence
linkTitle: Immediate Setup & Data Persistence
description: Get PostgreSQL running in Docker in under five minutes. Learn how to configure named volumes and bind mounts to persist your database across container restarts.
keywords:
  - PostgreSQL Docker
  - Docker Compose PostgreSQL
  - container database
weight: 10
---

This guide gets you from zero to a running PostgreSQL container in under five minutes, then explains how to keep your data safe across container restarts and removals.

## Overview

Running PostgreSQL in Docker requires understanding one critical concept: containers are ephemeral, but your data shouldn't be. This guide covers:

- Starting PostgreSQL with a single command
- Understanding why containers lose data by default
- Configuring volumes for persistent storage
- Translating your setup to Docker Compose

## Quick Start (Minimal Viable Container)

Run PostgreSQL immediately with this single command:

```console
$ docker run --rm --name postgres-dev \
  -e POSTGRES_PASSWORD=mysecretpassword \
  -p 5432:5432 \
  -d postgres:18
```

### Understanding the flags

| Flag | Purpose |
|------|---------|
| `--name postgres-dev` | Assigns a memorable name instead of a random string |
| `-e POSTGRES_PASSWORD=...` | Sets the superuser password (required) |
| `-p 5432:5432` | Maps host port 5432 to container port 5432 |
| `-d` | Runs the container in the background (detached mode) |

Verify the container is running:

```console
$ docker ps --filter name=postgres-dev
CONTAINER ID   IMAGE         COMMAND                  STATUS         PORTS                    NAMES
a1b2c3d4e5f6   postgres:18   "docker-entrypoint.s…"   Up 2 seconds   0.0.0.0:5432->5432/tcp   postgres-dev
```

Connect using psql from inside the container:

```console
$ docker exec -it postgres-dev psql -U postgres
psql (18.0)
Type "help" for help.

postgres=#
```

You now have a working PostgreSQL instance. But there's a problem—stop this container and your data disappears.

## The Data Persistence Problem

Containers use an ephemeral filesystem. When a container is removed, everything inside it, including your database files, is deleted.

Demonstrate this yourself:

```console
$ docker exec postgres-dev psql -U postgres -c "CREATE DATABASE testdb;"
CREATE DATABASE

$ docker exec postgres-dev psql -U postgres -c "\l" | grep testdb
 testdb    | postgres | UTF8     | libc            | en_US.utf8 | en_US.utf8 |            |           |

$ docker stop postgres-dev
postgres-dev

$ docker run --rm --name postgres-dev -e POSTGRES_PASSWORD=mysecretpassword -p 5432:5432 -d postgres:18

$ docker exec postgres-dev psql -U postgres -c "\l" | grep testdb
(no output - database is gone)
```

Your `testdb` database vanished because the new container started with a fresh filesystem. This is expected behavior—and exactly why volumes exist.

## Named Volumes

Named volumes are Docker-managed storage locations that persist independently of containers. Docker handles the filesystem location, permissions, and lifecycle.

Create a container with a named volume:

```console
$ docker run --rm --name postgres-dev \
  -e POSTGRES_PASSWORD=mysecretpassword \
  -p 5432:5432 \
  -v postgres_data:/var/lib/postgresql \
  -d postgres:18
```

The `-v postgres_data:/var/lib/postgresql` flag mounts a named volume called `postgres_data` to PostgreSQL's data directory. If the volume doesn't exist, Docker creates it automatically.

> **Note:** PostgreSQL 18+ stores data in a version-specific subdirectory under `/var/lib/postgresql`. Mounting at this level (rather than `/var/lib/postgresql/data`) allows for easier upgrades using `pg_upgrade --link`.

### Verify persistence works

```console
$ docker exec postgres-dev psql -U postgres -c "CREATE DATABASE testdb;"
CREATE DATABASE

$ docker stop postgres-dev
postgres-dev

$ docker run --rm --name postgres-dev \
  -e POSTGRES_PASSWORD=mysecretpassword \
  -p 5432:5432 \
  -v postgres_data:/var/lib/postgresql \
  -d postgres:18

$ docker exec postgres-dev psql -U postgres -c "\l" | grep testdb
 testdb    | postgres | UTF8     | libc            | en_US.utf8 | en_US.utf8 |            |           |
```

The database survived because the volume preserved it.

## Bind Mounts (Alternative)

## Verifying Persistence

## Docker Compose Version