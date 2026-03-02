---
title: Advanced Configuration and Initialization
linkTitle: Advanced Configuration and Initialization
weight: 20
description: Configure PostgreSQL initialization scripts, tune performance parameters, and set timezone and locale settings for containerized deployments.
keywords:
  - PostgreSQL Docker
  - Docker Compose PostgreSQL
  - container database
  - PostgreSQL performance tuning
---

With persistent storage configured in the previous section, you're ready to customize PostgreSQL for real-world use. This guide covers advanced configuration techniques for running PostgreSQL in Docker containers, including automated database initialization, performance tuning, and timezone configuration.

## Overview

While PostgreSQL containers can be started quickly with default settings, production environments require customized configurations. This guide explains how to:

- Automate database, schema, and user creation during container startup
- Tune PostgreSQL performance parameters for containerized workloads
- Configure timezone and locale settings

## Initialization scripts

The official PostgreSQL Docker image supports running initialization scripts automatically when the container starts for the first time. Any files placed in the `/docker-entrypoint-initdb.d/` directory are executed in alphabetical order.

### How initialization works

When the container starts, it checks whether the data directory (`/var/lib/postgresql/data`) is empty. If the directory already contains data, PostgreSQL starts immediately without running any initialization. If the directory is empty, the container runs `initdb` to create a new database cluster, then executes all scripts in `/docker-entrypoint-initdb.d/` in alphabetical order before starting PostgreSQL.

### Supported file formats

| Format | Description |
|--------|-------------|
| `.sql` | SQL commands executed directly |
| `.sql.gz` | Gzip-compressed SQL files |
| `.sh` | Shell scripts executed with bash |

> **Important:** Initialization scripts only run when the PostgreSQL data directory (`/var/lib/postgresql/data`) is empty. If you mount a volume containing existing data, initialization is skipped. This behavior prevents overwriting existing databases.

## Mounting initialization scripts

Use Docker Compose to mount initialization scripts into the container. First, create a project directory:

```console
$ mkdir -p postgres-project/init-db
$ cd postgres-project
```

Create a `compose.yaml` file:

```yaml
services:
  db:
    image: postgres:18
    volumes:
      - ./init-db:/docker-entrypoint-initdb.d
      - postgres_data:/var/lib/postgresql
    environment:
      POSTGRES_PASSWORD: mysecretpassword

volumes:
  postgres_data:
```

All scripts in the `./init-db` directory execute when the container starts for the first time. This is great for bootstrapping databases.

## Initialization script example

Create a file named `init.sql` in your `init-db` directory:

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

This script runs automatically when the container starts for the first time, creating your initial database schema.

> **Note:** Ensure initialization scripts have proper read permissions. If you encounter "Permission denied" errors, run `chmod 644 init-db/*.sql` to make the files readable by the container.

## Performance tuning

Default PostgreSQL settings are conservative to work on systems with limited resources. For production workloads, you should tune these parameters based on your container's allocated resources.

### Method 1: Custom configuration file

For complete control, mount a custom `postgresql.conf` file. First, extract the default configuration:

```console
$ docker run -i --rm postgres:18 cat /usr/share/postgresql/postgresql.conf.sample > my-postgres.conf
```

Edit `my-postgres.conf` with your desired settings, then mount it in your Compose file:

```yaml
services:
  db:
    image: postgres:18
    volumes:
      - ./my-postgres.conf:/etc/postgresql/postgresql.conf
      - ./init-db:/docker-entrypoint-initdb.d
      - postgres_data:/var/lib/postgresql
    command: postgres -c config_file=/etc/postgresql/postgresql.conf
    environment:
      POSTGRES_PASSWORD: mysecretpassword

volumes:
  postgres_data:
```

## Key configuration parameters

The following tables list important `postgresql.conf` parameters for containerized PostgreSQL deployments.

### Connection settings

| Parameter | Description | Default |
|-----------|-------------|---------|
| `listen_addresses` | IP addresses to listen on | `localhost` |
| `port` | TCP port number | `5432` |
| `max_connections` | Maximum concurrent connections | `100` |

### Memory settings

| Parameter | Description | Recommended starting value |
|-----------|-------------|---------------------------|
| `shared_buffers` | Shared memory for caching | 25% of container memory |
| `work_mem` | Memory per query operation | 4MB - 64MB |
| `maintenance_work_mem` | Memory for VACUUM, CREATE INDEX | 64MB - 256MB |
| `effective_cache_size` | Planner's cache size estimate | 50-75% of container memory |

> **Docker memory limits:** When tuning memory parameters, set explicit memory limits on your container using `deploy.resources.limits.memory` in Compose or `--memory` with `docker run`. Without limits, PostgreSQL sees the host's total RAM and may allocate more than intended. For example, if your container should use 4GB maximum, set `shared_buffers` to approximately 1GB (25%).

### I/O settings

| Parameter | Description | Recommended starting value |
|-----------|-------------|---------------------------|
| `effective_io_concurrency` | Concurrent disk I/O operations | `200` for SSDs, `2` for HDDs |

### Timeout settings

| Parameter | Description | Default |
|-----------|-------------|---------|
| `statement_timeout` | Max time for any statement | `0` (disabled) |
| `lock_timeout` | Max time to wait for a lock | `0` (disabled) |
| `deadlock_timeout` | Time before checking for deadlock | `1s` |
| `transaction_timeout` | Max time for a transaction | `0` (disabled) |

> **Note:** Setting `shared_buffers` too high in a container can exceed kernel shared memory limits. Use no more than 25-30% of the container's memory limit.

## Timezone and locale configuration

Proper localization ensures timestamps and sorting behave correctly for your application's users.

```yaml
services:
  db:
    image: postgres:18
    volumes:
      - postgres_data:/var/lib/postgresql
      - /etc/localtime:/etc/localtime:ro
      - /etc/timezone:/etc/timezone:ro
    environment:
      POSTGRES_PASSWORD: mysecretpassword
      TZ: America/New_York

volumes:
  postgres_data:
```

Alternatively, set the timezone using a PostgreSQL command-line parameter:

```yaml
services:
  db:
    image: postgres:18
    command: ["postgres", "-c", "timezone=America/New_York"]
    environment:
      POSTGRES_PASSWORD: mysecretpassword
```

### Setting the locale

Specify locale settings during database initialization using the `POSTGRES_INITDB_ARGS` environment variable:

```yaml
services:
  db:
    image: postgres:18
    volumes:
      - postgres_data:/var/lib/postgresql
    environment:
      POSTGRES_PASSWORD: mysecretpassword
      POSTGRES_INITDB_ARGS: "--encoding=UTF8 --lc-collate=en_US.UTF-8 --lc-ctype=en_US.UTF-8"

volumes:
  postgres_data:
```

This affects collation (sorting) and character processing behavior. Changing this variable after database creation has no effect—it only applies during the first run when the data directory is initialized.

## Connecting to the database

You can interact with PostgreSQL running in a container even without `psql` installed on your host machine.

### Interactive shell

Open a `psql` session inside the container:

```console
$ docker exec -it postgres-container psql -U postgres
```

Connect to a specific database:

```console
$ docker exec -it postgres-container psql -U postgres -d mydb
```