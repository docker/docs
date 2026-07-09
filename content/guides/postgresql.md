---
title: PostgreSQL specific guide
linkTitle: PostgreSQL
description: Containerize PostgreSQL databases using Docker
keywords: Docker, getting started, postgresql, language
summary: |
  This guide explains how to containerize PostgreSQL databases using
  Docker.
aliases:
  - /guides/postgresql/advanced-configuration-and-initialization/
  - /guides/postgresql/companions-for-postgresql/
  - /guides/postgresql/immediate-setup-and-data-persistence/
  - /guides/postgresql/networking-and-connectivity/
params:
  tags: [databases]
  time: 20 minutes
---


## Immediate setup & data persistence

This guide gets you from zero to a running PostgreSQL container in under five minutes, then explains how to keep your data safe across container restarts and removals.

### Overview

Running PostgreSQL in Docker requires understanding one critical concept: containers are ephemeral, but your data shouldn't be. This guide covers:

- Starting PostgreSQL with a single command
- Understanding why containers lose data by default
- Configuring volumes for persistent storage
- Translating your setup to Docker Compose

### Quick start (minimal viable container)

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

#### Understanding the flags

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

Connect using `psql` from inside the container:

```console
$ docker exec -it postgres-dev psql -U postgres
psql (18.0)
Type "help" for help.

postgres=#
```

You now have a working PostgreSQL instance. But there's a problem—stop this container and your data disappears.

### The data persistence problem

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

### Named volumes

Named volumes are Docker-managed storage locations that persist independently of containers. Docker handles the filesystem location, permissions, and lifecycle.

Create a container with a named volume:

{{< tabs >}}
{{< tab name="Using DHIs" >}}

You must authenticate to dhi.io before you can pull Docker Hardened Images. Run `docker login dhi.io` to authenticate.

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

> [!NOTE]
>
> PostgreSQL 18+ stores data in a version-specific subdirectory under `/var/lib/postgresql`. Mounting at this level (rather than `/var/lib/postgresql/data`) allows for easier upgrades using `pg_upgrade --link`.

#### Verify persistence works

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
```

{{< /tab >}}

{{< /tabs >}}

If you see `testdb` in the output, persistence works: The database survived because the volume preserved the data directory.

#### Managing volumes

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

### Bind mounts (alternative)

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

#### When to use bind mounts

Bind mounts are useful when you need direct filesystem access to the data directory for backup scripts that read files directly, when integrating with host-level monitoring tools, or when specific permission requirements exist. For most development and production scenarios, named volumes are simpler and less error-prone.

#### Common bind mount issues

Permission errors are the most frequent problem with bind mounts. PostgreSQL runs as user `postgres` (UID 999) inside the container. If your host directory has restrictive permissions, the container fails to start.

Check logs if the container exits immediately:

```console
$ docker logs postgres-dev
```

### Docker Compose configuration

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

#### Environment variables reference

The official PostgreSQL image supports these environment variables:

| Variable | Required | Description |
|----------|----------|-------------|
| `POSTGRES_PASSWORD` | Yes | Superuser password |
| `POSTGRES_USER` | No | Superuser name (default: `postgres`) |
| `POSTGRES_DB` | No | Default database name (default: value of `POSTGRES_USER`) |

### Next steps

With persistent storage configured, you're ready to customize PostgreSQL further. The next chapter of the guide covers:

- Automated schema creation with initialization scripts
- Performance tuning for containerized workloads
- Timezone and locale configuration

## Advanced Configuration and Initialization

With persistent storage configured in the previous section, you're ready to customize PostgreSQL for real-world use. This guide covers advanced configuration techniques for running PostgreSQL in Docker containers, including automated database initialization, performance tuning, and timezone configuration.

### Overview

While PostgreSQL containers can be started quickly with default settings, production environments require customized configurations. This guide explains how to:

- Automate database, schema, and user creation during container startup
- Tune PostgreSQL performance parameters for containerized workloads
- Configure timezone and locale settings

### Initialization scripts

The official PostgreSQL Docker image supports running initialization scripts automatically when the container starts for the first time. Any files placed in the `/docker-entrypoint-initdb.d/` directory are executed in alphabetical order.

#### How initialization works

When the container starts, it checks whether the PostgreSQL data directory is empty. If the directory already contains data, PostgreSQL starts immediately without running any initialization. If the directory is empty, the container runs `initdb` to create a new database cluster, then executes all scripts in `/docker-entrypoint-initdb.d/` in alphabetical order before starting PostgreSQL.

#### Supported file formats

| Format | Description |
|--------|-------------|
| `.sql` | SQL commands executed directly |
| `.sql.gz` | Gzip-compressed SQL files |
| `.sh` | Shell scripts executed with bash |

> [!IMPORTANT]
>
> Initialization scripts only run when the PostgreSQL data directory (`/var/lib/postgresql/data`) is empty. If you mount a volume containing existing data, initialization is skipped. This behavior prevents overwriting existing databases.

### Mounting initialization scripts

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

### Initialization script example

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

> [!NOTE]
>
> Ensure initialization scripts have proper read permissions. If you encounter "Permission denied" errors, run `chmod 644 init-db/*.sql` to make the files readable by the container.

### Performance tuning

Default PostgreSQL settings are conservative to work on systems with limited resources. For production workloads, you should tune these parameters based on your container's allocated resources.

#### Method 1: Custom configuration file

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

### Key configuration parameters

The following tables list important `postgresql.conf` parameters for containerized PostgreSQL deployments.

#### Connection settings

| Parameter | Description | Default |
|-----------|-------------|---------|
| `listen_addresses` | IP addresses to listen on | `localhost` |
| `port` | TCP port number | `5432` |
| `max_connections` | Maximum concurrent connections | `100` |

#### Memory settings

| Parameter | Description | Recommended starting value |
|-----------|-------------|---------------------------|
| `shared_buffers` | Shared memory for caching | 25% of container memory |
| `work_mem` | Memory per query operation | 4MB - 64MB |
| `maintenance_work_mem` | Memory for VACUUM, CREATE INDEX | 64MB - 256MB |
| `effective_cache_size` | Planner's cache size estimate | 50-75% of container memory |

##### Docker memory limits

When tuning memory parameters, set explicit memory limits on your container using `deploy.resources.limits.memory` in Compose or `--memory` with `docker run`. Without limits, PostgreSQL sees the host's total RAM and may allocate more than intended. For example, if your container should use 4GB maximum, set `shared_buffers` to approximately 1GB (25%).

#### I/O settings

| Parameter | Description | Recommended starting value |
|-----------|-------------|---------------------------|
| `effective_io_concurrency` | Concurrent disk I/O operations | `200` for SSDs, `2` for HDDs |

#### Timeout settings

| Parameter | Description | Default |
|-----------|-------------|---------|
| `statement_timeout` | Max time for any statement | `0` (disabled) |
| `lock_timeout` | Max time to wait for a lock | `0` (disabled) |
| `deadlock_timeout` | Time before checking for deadlock | `1s` |
| `transaction_timeout` | Max time for a transaction | `0` (disabled) |

> [!NOTE]
>
> Setting `shared_buffers` too high in a container can exceed kernel shared memory limits. Use no more than 25-30% of the container's memory limit.

### Timezone and locale configuration

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

#### Setting the locale

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

### Connecting to the database

You can interact with PostgreSQL running in a container even without `psql` installed on your host machine.

#### Interactive shell

Open a `psql` session inside the container:

```console
$ docker exec -it postgres-container psql -U postgres
```

Connect to a specific database:

```console
$ docker exec -it postgres-container psql -U postgres -d mydb
```

## Networking and connectivity

This guide covers two common ways to connect to PostgreSQL running in Docker:

- Container-to-container: Connect from your application container to PostgreSQL over a private Docker network. No ports need to be exposed to the host.
- Host-to-container: Connect from your laptop or development machine using `localhost` and a published port.

Prerequisite: This guide assumes you have PostgreSQL running with persistent storage. If you don't, follow the [Immediate Setup & Data Persistence](/guides/postgresql/immediate-setup-and-data-persistence/) guide first.

### Internal network access (container-to-container)

When your application runs in another container, connecting to PostgreSQL through a user-defined bridge network is the recommended approach. This setup provides automatic DNS resolution, so your application can connect to PostgreSQL using the container name as the hostname, without needing to track IP addresses.

> [!NOTE]
> Why not use the default bridge network? While containers on the default bridge network can communicate, they can only do so by IP address. Since container IP addresses change when containers restart, this would require updating your PostgreSQL connection strings each time. User-defined bridge networks solve this by providing automatic DNS resolution, ensuring your PostgreSQL connection strings remain stable even if containers restart and receive new IP addresses.

Here's a quick comparison:

> [!NOTE]
>
> The following examples show the difference in approach. To actually test this, follow the steps in this guide to set up containers on the appropriate networks first.

With the default bridge network, you'd need to find the IP address first:
```bash
# Get the container's IP address (changes on restart)
docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' postgres-dev
# Output: 172.17.0.2

# Then connect using that IP address from another container
# (No --network flag needed - containers default to bridge network)
docker run --rm -it \
  -e PGPASSWORD=mysecretpassword \
  postgres:18 \
  psql -h 172.17.0.2 -U postgres
```

With a user-defined network, you simply use the container name:
```bash
# Container name works directly - no IP lookup needed
docker run --rm -it \
  --network my-app-net \
  -e PGPASSWORD=mysecretpassword \
  postgres:18 \
  psql -h postgres-dev -U postgres
```

#### Step 1: Create a user-defined network

```bash
docker network create my-app-net

# Example Output
ab7f984be43a0ca15534a9ee568716ddbe869a5875077fad3ef3192e3af7d288

docker network ls
# Output
ab7f984be43a   my-app-net    bridge    local


```

#### Step 2: Run PostgreSQL on that network (no port publishing)

Notice there is no `-p 5432:5432` here. This keeps PostgreSQL internal to Docker and not accessible from the host machine, which is more secure for production environments.

```bash
docker run -d --name postgres-dev \
  --network my-app-net \
  -e POSTGRES_PASSWORD=mysecretpassword \
  -v postgres_data:/var/lib/postgresql \
  postgres:18

  # Output
CONTAINER ID  IMAGE        COMMAND                 CREATED         STATUS        PORTS     NAMES
6d351ed89efc  postgres:18  "docker-entrypoint.s…"  9 seconds ago   Up 8 seconds  5432/tcp  postgres-dev

```

#### Step 3: Connect from another container using the Postgres container name

You can test connectivity with a temporary `psql` client container:

```bash
docker run --rm -it \
  --network my-app-net \
  -e PGPASSWORD=mysecretpassword \
  postgres:18 \
  psql -h postgres-dev -U postgres
```

Key point: `-h postgres-dev` works because Docker DNS resolves the container name on a user-defined network. The container name acts as the hostname.

#### Connection string examples

When connecting from your application container, use these PostgreSQL connection strings:

- PostgreSQL URI format:
  This is the standard PostgreSQL connection URI format that combines all connection parameters into a single string, widely supported by PostgreSQL clients and libraries.

  ```bash
  postgresql://postgres:mysecretpassword@postgres-dev:5432/postgres
  ```

  This command demonstrates passing a PostgreSQL URI connection string as an environment variable to a container, which your application can then read to connect to the database.

  Example usage in a Docker run command:
  ```bash
  docker run --rm -it \
    --network my-app-net \
    -e DATABASE_URL="postgresql://postgres:mysecretpassword@postgres-dev:5432/postgres" \
    alpine:latest \
    sh -c 'echo "DATABASE_URL is set to: $DATABASE_URL"'
  ```


- PostgreSQL connection parameters:
  This format uses key-value pairs separated by spaces, which many PostgreSQL client libraries accept as an alternative to URI format.
  ```bash
  host=postgres-dev
  port=5432
  user=postgres
  password=mysecretpassword
  dbname=postgres
  ```

  Example usage in application code (Python with psycopg2):
  ```python
  conn = psycopg2.connect(
      host="postgres-dev",
      port=5432,
      user="postgres",
      password="mysecretpassword",
      dbname="postgres"
  )
  ```

- Connecting to a specific database:
  Replace the database name in the connection string to connect to a specific database instead of the default `postgres` database.
  If you created a custom database (e.g., `testdb`), use:
  ```bash
  postgresql://postgres:mysecretpassword@postgres-dev:5432/testdb
  ```

  Example with SSL disabled (common in Docker networks):
  Add `?sslmode=disable` to the connection string when connecting within a private Docker network where SSL encryption isn't required.
  ```bash
  postgresql://postgres:mysecretpassword@postgres-dev:5432/testdb?sslmode=disable
  ```

> [!NOTE]
>
> The default port `5432` is used in these examples. If you're connecting to a different PostgreSQL instance or have changed the port, update the connection string accordingly. The container name (`postgres-dev`) is resolved by Docker DNS to the container's IP address on the network.


### Connecting from the host (external access)

To connect to PostgreSQL from your host machine using tools like `psql`, `pgAdmin`, `DBeaver`, or database management scripts, you need to publish PostgreSQL's port (`5432`) to the host. This allows external tools to reach the PostgreSQL container.

#### Expose Postgres to localhost only (recommended for development)

This binds to `127.0.0.1` so it's only reachable from your local machine, not from other devices on your network. This is the most secure option for development.

```bash
docker run -d --name postgres-dev \
  -e POSTGRES_PASSWORD=mysecretpassword \
  -p 127.0.0.1:5432:5432 \
  -v postgres_data:/var/lib/postgresql \
  postgres:18
```

Now connect from your host:

- Host: `localhost` or `127.0.0.1`
- Port: `5432`

If you have `psql` installed on your host:
```bash
psql -h localhost -p 5432 -U postgres
```

You'll be prompted for the password. Alternatively, you can use the `PGPASSWORD` environment variable:
```bash
PGPASSWORD=mysecretpassword psql -h localhost -p 5432 -U postgres
```

#### Connecting with PostgreSQL GUI tools

Popular PostgreSQL GUI tools can connect using these common connection details: Host: `localhost`, Port: `5432`, User: `postgres`, Database: `postgres` (or your database name).

- pgAdmin: A web-based PostgreSQL administration and development platform
- DBeaver: A universal database tool that supports PostgreSQL and many other databases. Select PostgreSQL as the connection type
- TablePlus: A modern, native database management tool for macOS and Windows with a clean interface

All tools will prompt for the password you set with `POSTGRES_PASSWORD`.

#### Expose Postgres to all network interfaces (use with caution)

To allow connections from other devices on your network, use `-p 5432:5432` instead of `-p 127.0.0.1:5432:5432`. This binds PostgreSQL to all network interfaces on your host, making it accessible from any device that can reach your host, not just localhost.

```bash
docker run -d --name postgres-dev \
  -e POSTGRES_PASSWORD=mysecretpassword \
  -p 5432:5432 \
  -v postgres_data:/var/lib/postgresql \
  postgres:18
```

> [!WARNING]
>
> Exposing PostgreSQL to all network interfaces (`0.0.0.0:5432`) makes it accessible from any device that can reach your host. Only use this in trusted network environments or behind a firewall. For production, consider using a reverse proxy or VPN instead.

#### PostgreSQL security considerations for external access

When exposing PostgreSQL to external access, follow these PostgreSQL-specific security practices:

- Avoid using the `postgres` superuser: The default `postgres` user has full database privileges. Create dedicated users with only the permissions your application needs.
- Use strong passwords: PostgreSQL passwords should be complex. Consider using environment variables or secrets management instead of `hardcoding` passwords.
- Limit network exposure: Binding to `127.0.0.1` (localhost only) is safer than exposing to all interfaces (`0.0.0.0`).
- Consider SSL/TLS: For production, configure PostgreSQL to require SSL connections. The [Advanced Configuration and Initialization](/guides/postgresql/advanced-configuration-and-initialization/) guide shows how to configure PostgreSQL settings.
- Create application-specific users: Use initialization scripts to create users with limited privileges. For example, a read-only user for reporting or a user that can only access specific databases.

The [Advanced configuration and initialization](/guides/postgresql/advanced-configuration-and-initialization/) guide shows how to use initialization scripts to create users and roles automatically.

### Using Docker Compose for networking

Docker Compose automatically creates a network for your services, making networking configuration simpler. Here's an example that combines both internal and external access:

```yaml
services:
  db:
    image: postgres:18
    container_name: postgres-dev
    environment:
      POSTGRES_PASSWORD: mysecretpassword
    volumes:
      - postgres_data:/var/lib/postgresql
    ports:
      - "127.0.0.1:5432:5432"  # Expose to localhost only
    networks:
      - app-network

  app:
    build: ./my-app
    environment:
      DATABASE_URL: postgresql://postgres:mysecretpassword@db:5432/mydb
    networks:
      - app-network
    depends_on:
      - db

volumes:
  postgres_data:

networks:
  app-network:
    driver: bridge
```

In this PostgreSQL-focused setup:
- The `app` service connects to PostgreSQL using the service name (`db`) as the hostname in the connection string
- PostgreSQL is accessible from your host at `localhost:5432` for external tools
- Both services are isolated on a custom network, providing network-level security
- The `depends_on` directive ensures PostgreSQL starts before your application

PostgreSQL connection details for the app service:
- Hostname: `db` (resolved by Docker DNS)
- Port: `5432` (PostgreSQL default port)
- Database: `mydb` (as specified in the connection string)
- User: `postgres` (or a custom user you've created)

> [!NOTE]
>
> Docker Compose automatically creates a network for your project. Services can reach each other by service name without explicit network configuration, but defining a custom network gives you more control. For PostgreSQL, this means your application can always connect using the service name, regardless of container restarts or IP changes.

### Troubleshooting

This section covers common PostgreSQL connection issues and their solutions when working with Docker networking.

#### "Could not translate host name postgres-dev"

- Both containers must be on the same Docker network (`my-app-net`).
- Verify the network exists: `docker network ls`
- Check which network a container is on: `docker inspect postgres-dev | grep NetworkMode`
- Ensure you're using a user-defined network, not the default bridge network

#### "Connection refused" or "could not connect to server"

- PostgreSQL may still be initializing: PostgreSQL takes a few seconds to start and initialize the database cluster. Wait 5-10 seconds after container start and retry.
- Check if the PostgreSQL container is running:

  ```bash
  docker ps --filter name=postgres-dev
  ```

- Check PostgreSQL logs for initialization or connection errors:

  ```bash
  docker logs postgres-dev
  ```

  Look for messages like "database system is ready to accept connections" to confirm PostgreSQL is fully started.

- Verify the port mapping is correct:

  ```bash
  docker port postgres-dev
  ```

  This should show `5432/tcp -> 127.0.0.1:5432` (or `0.0.0.0:5432` if bound to all interfaces).

- Test PostgreSQL connectivity from inside the container:

  ```bash
  docker exec -it postgres-dev psql -U postgres -c "SELECT version();"
  ```

  If this works but external connections fail, the issue is with port publishing, not PostgreSQL itself.

#### "Password authentication failed" or "FATAL: password authentication failed for user"

- Confirm the password: Verify you're using the same password set in `POSTGRES_PASSWORD` when you started the container.
- Existing volume with old credentials: If you reused an existing volume, the password from the original initialization is still in effect. The `POSTGRES_PASSWORD` environment variable only sets the password during the first database initialization. To reset:
  - Remove the volume: `docker volume rm postgres_data`
  - Or connect with the old password
  - Or change the password after connecting: `ALTER USER postgres WITH PASSWORD 'newpassword';`
- Try connecting with password prompt: `psql -h localhost -U postgres -W` (the `-W` flag forces a password prompt)
- Use PGPASSWORD environment variable: `PGPASSWORD=mysecretpassword psql -h localhost -U postgres`
- Check PostgreSQL authentication configuration: If you've customized `pg_hba.conf`, verify the authentication method allows password authentication

#### "Network not found"

- Ensure the network exists before starting containers: `docker network create my-app-net`
- If using Docker Compose, the network is created automatically when you run `docker compose up`

## Companions for PostgreSQL

### PostgreSQL ecosystem companions: pgAdmin, PgBouncer, and performance testing

Running a standalone PostgreSQL container is often just the beginning. What happens when thousands of connections arrive, or when you need a visual interface to manage your database?

This is where **companion tools** come into play. These applications extend PostgreSQL with capabilities the core database engine doesn't provide natively: visual administration, connection pooling, and performance benchmarking. This guide covers how to deploy pgAdmin 4, PgBouncer, Pgpool-II, and `pgbench` in Docker, when to use each tool, and real-world benchmark results demonstrating their performance impact.

### pgAdmin 4: Visual management platform

pgAdmin 4 is the industry-standard open source management tool for PostgreSQL. When deployed in Docker, it typically runs in **Server Mode**, providing a multi-user web interface to manage one or more database instances.

While you can accomplish everything from the command line using `psql`, a visual interface significantly simplifies writing complex queries, visualizing table structures, and exploring database objects.

#### Key considerations

When running pgAdmin in Docker, keep these points in mind:

- **Image**: Use the official `dpage/pgadmin4` image
- **Networking**: In a Docker Compose environment, pgAdmin connects to the database using the internal service name (for example, `db:5432`) rather than `localhost`

#### Docker Compose configuration

To quickly deploy pgAdmin:

```yaml
pgadmin:
  image: dpage/pgadmin4:8.14
  environment:
    PGADMIN_DEFAULT_EMAIL: admin@example.com
    PGADMIN_DEFAULT_PASSWORD: secure_password
  volumes:
    - pgadmin_data:/var/lib/pgadmin
  ports:
    - "8080:80"
```

With this configuration, access the pgAdmin interface at `http://localhost:8080`. Use the email and password specified in the environment variables for initial sign in.

> [!IMPORTANT]
>
> In production environments, pass `PGADMIN_DEFAULT_PASSWORD` as an external environment variable or use Docker secrets. Storing passwords in plain text within `docker-compose.yml` poses a security risk.

Now that you have visual database management in place, the next challenge in production environments is handling connection load. The following section explains how to manage high-volume database traffic.

### PgBouncer: Lightweight connection pooling

PostgreSQL creates a new process for every client connection, which consumes significant RAM. What happens when you have 1,000 concurrent users? PgBouncer solves exactly this problem.

PgBouncer is a lightweight proxy that pools connections, allowing thousands of applications to share a small number of actual database backends. Think of it as a traffic controller: everyone wants to pass through simultaneously, but the controller regulates the flow to prevent congestion.

#### Pooling modes

PgBouncer offers three distinct pooling modes:

| Mode | Description | Use case |
|------|-------------|----------|
| **Session** | Connection assigned for entire session duration | Long-lived connections, session variables |
| **Transaction** | Connection returned after each transaction ends | Web applications, microservices (most common) |
| **Statement** | Connection returned after every SQL statement | Simple queries, no multi-statement transactions |

#### When to use PgBouncer

PgBouncer becomes essential when you encounter:

- "too many connections" errors
- High memory consumption due to connection overhead
- Many short-lived connections (web applications, serverless functions)
- Need to serve thousands of clients with limited database connections

#### Complete Docker Compose setup

To run PostgreSQL and PgBouncer together, you need three files: `docker-compose.yml`, `pgbouncer.ini`, and `userlist.txt`.

First, create the PgBouncer configuration file (`pgbouncer.ini`):

```bash
[databases]
benchmark = host=postgres port=5432 dbname=benchmark user=postgres

[pgbouncer]
listen_addr = 0.0.0.0
listen_port = 6432
auth_type = trust
auth_file = /etc/pgbouncer/userlist.txt
admin_users = postgres
pool_mode = transaction
max_client_conn = 1000
default_pool_size = 50
min_pool_size = 10
reserve_pool_size = 10
max_db_connections = 100
```

Next, create the user authentication file (`userlist.txt`):

```bash
"postgres" "postgres"
```

Finally, create the Docker Compose file (`docker-compose.yml`):

```yaml
services:
  postgres:
    image: postgres:18
    container_name: postgres
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: benchmark
      POSTGRES_HOST_AUTH_METHOD: trust
    volumes:
      - postgres_data:/var/lib/postgresql
    ports:
      - "5432:5432"
    networks:
      - pgnet
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 5s
      retries: 5

  pgbouncer:
    image: percona/percona-pgbouncer:1.25.0
    container_name: pgbouncer
    volumes:
      - ./pgbouncer.ini:/etc/pgbouncer/pgbouncer.ini
      - ./userlist.txt:/etc/pgbouncer/userlist.txt
    ports:
      - "6432:6432"
    networks:
      - pgnet
    depends_on:
      postgres:
        condition: service_healthy

volumes:
  postgres_data:

networks:
  pgnet:
    driver: bridge
```

Key configuration notes:

- `PgBouncer` listens on port **6432**, avoiding confusion with the direct PostgreSQL connection on port 5432
- The `depends_on` directive with `service_healthy` condition ensures PgBouncer starts only after PostgreSQL is ready
- `pool_mode = transaction` is the optimal choice for most web applications
- The [Percona PgBouncer image](https://hub.docker.com/r/percona/percona-pgbouncer) requires mounted configuration files (without the `:ro` flag, as the entrypoint script needs to modify them)
- This example uses `trust` authentication for simplicity. In production, configure proper SCRAM-SHA-256 authentication

> [!NOTE]
>
> The `Percona PgBouncer` entrypoint script processes the configuration files on startup. Mount them without the read-only flag to avoid permission errors.




### `pgbench`: Performance benchmarking

`pgbench` is a benchmarking utility included with the official PostgreSQL image. It allows you to simulate heavy workloads and verify how your Docker configuration performs under pressure.

#### Initialize benchmark tables

First, create the test tables. The `-s` (scale) parameter determines data size—scale factor 50 creates approximately 5 million rows:

```bash
docker exec postgres pgbench -i -s 50 -U postgres benchmark
```

#### Run stress tests

Key parameters:

- `-c`: Number of simulated clients
- `-j`: Number of threads
- `-T`: Duration in seconds

Test with direct PostgreSQL connection:

```bash
docker exec postgres pgbench -h localhost -U postgres -c 50 -j 4 -T 60 benchmark
```

Test through PgBouncer:

```bash
docker exec postgres pgbench -h pgbouncer -p 6432 -U postgres -c 50 -j 4 -T 60 benchmark
```

### Understanding benchmark results

Does PgBouncer actually make a difference? Run the benchmarks yourself to find out. Your results will vary based on your hardware, Docker configuration, network setup, and system load.

#### What to expect

When you run these benchmarks, you'll observe patterns rather than specific numbers. Think of it like comparing two different routes to work: the "faster" route depends on traffic conditions, time of day, and your vehicle.

#### Key observations

When comparing direct connections versus PgBouncer, you'll typically notice:

##### 1. Connection overhead differs significantly

Direct connections require PostgreSQL to spawn a new process for each client. PgBouncer reuses existing connections. Watch the "initial connection time" metric in your results—PgBouncer often shows dramatically faster connection setup.

##### 2. Behavior under pressure reveals the real difference

Try increasing the client count (`-c` parameter) gradually: 50, 100, 150, 200. At some point, direct connections will fail with "too many clients already" while PgBouncer continues handling requests. This is PgBouncer's primary value: **it prevents connection exhaustion**.

##### 3. Throughput varies by environment

On some systems, direct connections show higher transactions per second (TPS) at low concurrency. On others, PgBouncer wins even with few clients. The difference depends on:
- CPU and memory available
- Docker networking overhead
- Disk I/O speed
- Whether connections are being rapidly opened and closed

