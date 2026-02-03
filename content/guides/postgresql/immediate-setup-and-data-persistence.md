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

> [!NOTE]
>
> [Docker Hardened Images (DHIs)](https://docs.docker.com/dhi/) are minimal, secure, and production-ready container base and application images maintained by Docker. DHIs are recommended whenever it is possible for better security. They are designed to reduce vulnerabilities and simplify compliance, freely available to everyone with no subscription required, no usage restrictions, and no vendor lock-in.

Run PostgreSQL immediately with this single command:

{{< tabs >}}
{{< tab name="Using DHIs" >}}

You must authenticate to dhi.io before you can pull Docker Hardened Images. Run `docker login dhi.io` to authenticate.

```console
docker run --rm --name postgres-dev \
  -e POSTGRES_PASSWORD=mysecretpassword \
  -p 5432:5432 \
  -d dhi.io/postgres:18
```

{{< /tab >}}

{{< tab name="Using DOIs" >}}

```console
$ docker run --rm --name postgres-dev \
  -e POSTGRES_PASSWORD=mysecretpassword \
  -p 5432:5432 \
  -d postgres:18
```

{{< /tab >}}
{{< /tabs >}}

### Understanding the flags

| Flag | Purpose |
|------|---------|
| `--rm` | Automatically removes the container when it stops |
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

{{< tabs >}}
{{< tab name="Using DHIs" >}}

```console
$ docker exec postgres-dev psql -U postgres -c "CREATE DATABASE testdb;"
CREATE DATABASE

$ docker exec postgres-dev psql -U postgres -c "\l" | grep testdb
 testdb    | postgres | UTF8     | libc            | en_US.utf8 | en_US.utf8 |            |           |

$ docker stop postgres-dev
postgres-dev

$ docker run --rm --name postgres-dev \
  -e POSTGRES_PASSWORD=mysecretpassword \
  -p 5432:5432 \
  -d dhi.io/postgres:18

$ docker exec postgres-dev psql -U postgres -c "\l" | grep testdb
(no output - database is gone)
```

{{< /tab >}}

{{< tab name="Using DOIs" >}}

```console
$ docker exec postgres-dev psql -U postgres -c "CREATE DATABASE testdb;"
CREATE DATABASE

$ docker exec postgres-dev psql -U postgres -c "\l" | grep testdb
 testdb    | postgres | UTF8     | libc            | en_US.utf8 | en_US.utf8 |            |           |

$ docker stop postgres-dev
postgres-dev

$ docker run --rm --name postgres-dev \
  -e POSTGRES_PASSWORD=mysecretpassword \
  -p 5432:5432 \
  -d postgres:18

$ docker exec postgres-dev psql -U postgres -c "\l" | grep testdb
(no output - database is gone)
```

{{< /tab >}}
{{< /tabs >}}

Your `testdb` database vanished because the new container started with a fresh filesystem. This is expected behavior—and exactly why volumes exist.

## Named Volumes

Named volumes are Docker-managed storage locations that persist independently of containers. Docker handles the filesystem location, permissions, and lifecycle.

Create a container with a named volume:

{{< tabs >}}
{{< tab name="Using DHIs" >}}

You must authenticate to dhi.io before you can pull Docker Hardened Images. Run docker login dhi.io to authenticate.

```console
$ docker run --rm --name postgres-dev \
  -e POSTGRES_PASSWORD=mysecretpassword \
  -p 5432:5432 \
  -v postgres_data:/var/lib/postgresql \
  -d dhi.io/postgres:18
```

{{< /tab >}}

{{< tab name="Using DOIs" >}}

```console
$ docker run --rm --name postgres-dev \
  -e POSTGRES_PASSWORD=mysecretpassword \
  -p 5432:5432 \
  -v postgres_data:/var/lib/postgresql \
  -d postgres:18
```

{{< /tab >}}
{{< /tabs >}}


The `-v postgres_data:/var/lib/postgresql` flag mounts a named volume called `postgres_data` to PostgreSQL's data directory. If the volume doesn't exist, Docker creates it automatically.

> **Note:** PostgreSQL 18+ stores data in a version-specific subdirectory under `/var/lib/postgresql`. Mounting at this level (rather than `/var/lib/postgresql/data`) allows for easier upgrades using `pg_upgrade --link`.

### Verify persistence works

To verify data persistence, repeat the previous test, but this time with the named volume attached in place.

{{< tabs >}}
{{< tab name="Using DHIs" >}}

```console
$ docker exec postgres-dev psql -U postgres -c "CREATE DATABASE testdb;"
CREATE DATABASE

$ docker stop postgres-dev
postgres-dev

$ docker run --rm --name postgres-dev \
  -e POSTGRES_PASSWORD=mysecretpassword \
  -p 5432:5432 \
  -v postgres_data:/var/lib/postgresql \
  -d dhi.io/postgres:18

$ docker exec postgres-dev psql -U postgres -c "\l" | grep testdb
 testdb    | postgres | UTF8     | libc            | en_US.utf8 | en_US.utf8 |            |           |
```

{{< /tab >}}

{{< tab name="Using DOIs" >}}

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

{{< /tab >}}

{{< /tabs >}}

If you see "testdb" in the output, persistence works: The database survived because the volume preserved the data directory.

### Managing volumes

List all volumes:

```console
$ docker volume ls --filter name=postgres_data
DRIVER    VOLUME NAME
local     postgres_data
```

Inspect a volume to see its details:

```console
$ docker volume inspect postgres_data
[
    {
        "CreatedAt": "2025-01-05T10:30:00Z",
        "Driver": "local",
        "Labels": null,
        "Mountpoint": "/var/lib/docker/volumes/postgres_data/_data",
        "Name": "postgres_data",
        "Options": null,
        "Scope": "local"
    }
]
```

Remove an unused volume (warning: this deletes all data):

```console
$ docker volume rm postgres_data
```

## Bind Mounts (Alternative)

Bind mounts map a specific host directory to a container path. Unlike named volumes, you control exactly where data lives on the host filesystem.

Create a directory on your host machine to store Postgres data.

{{< tabs >}}
{{< tab name="Using DHIs" >}}

```console
mkdir -p ~/postgres-data && sudo chown -R 999:999 ~/postgres-data
```

Run Postgres using a bind mount.

```console
docker run --rm --name postgres-dev \
  -e POSTGRES_PASSWORD=mysecretpassword \
  -p 5432:5432 \
  -v ~/postgres-data:/var/lib/postgresql \
  -d dhi.io/postgres:18
```

{{< /tab >}}

{{< tab name="Using DOIs" >}}

```console
$ mkdir -p ~/postgres-data
```

Run Postgres using a bind mount.

```console
$ docker run --rm --name postgres-dev \
  -e POSTGRES_PASSWORD=mysecretpassword \
  -p 5432:5432 \
  -v ~/postgres-data:/var/lib/postgresql \
  -d postgres:18
```

{{< /tab >}}
{{< /tabs >}}

### When to use bind mounts

Bind mounts are useful when you need direct filesystem access to the data directory for backup scripts that read files directly, when integrating with host-level monitoring tools, or when specific permission requirements exist. For most development and production scenarios, named volumes are simpler and less error-prone.

### Common bind mount issues

Permission errors are the most frequent problem with bind mounts. PostgreSQL runs as user `postgres` (UID 999) inside the container. If your host directory has restrictive permissions, the container fails to start.

Check logs if the container exits immediately:

```console
$ docker logs postgres-dev
```

## Docker Compose configuration

Docker Compose captures your entire configuration in a file, making setups reproducible and easier to manage as complexity grows.

Create a `compose.yaml` file:

```yaml
services:
  db:
    image: postgres:18
    container_name: postgres-dev
    environment:
      POSTGRES_PASSWORD: mysecretpassword
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql

volumes:
  postgres_data:
```

Start the database:

```console
$ docker compose up -d
```

Stop and remove containers (volume persists):

```console
$ docker compose down
```

Alternatively, you can stop, remove containers, and delete the volume:

```console
$ docker compose down -v
```

This compose file becomes the foundation for adding initialization scripts, performance tuning, and companion services covered in subsequent guides.

### Environment variables reference

The official PostgreSQL image supports these environment variables:

| Variable | Required | Description |
|----------|----------|-------------|
| `POSTGRES_PASSWORD` | Yes | Superuser password |
| `POSTGRES_USER` | No | Superuser name (default: `postgres`) |
| `POSTGRES_DB` | No | Default database name (default: value of `POSTGRES_USER`) |

## Next steps

With persistent storage configured, you're ready to customize PostgreSQL further. The next chapter of the guide covers:

- Automated schema creation with initialization scripts
- Performance tuning for containerized workloads
- Timezone and locale configuration