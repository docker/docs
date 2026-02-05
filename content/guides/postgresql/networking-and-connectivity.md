---
title: Networking and Connectivity
linkTitle: Networking and Connectivity
description: This module shows how to connect to PostgreSQL in Docker in two common ways; from another container (internal network) and from your host machine (external access).
keywords:
  - Networking PostgreSQL Docker
weight: 30
---

This guide covers two common ways to connect to PostgreSQL running in Docker:

- **Container-to-container**: Connect from your application container to PostgreSQL over a private Docker network. No ports need to be exposed to the host.
- **Host-to-container**: Connect from your laptop or development machine using `localhost` and a published port.

**Prerequisite**: This guide assumes you have PostgreSQL running with persistent storage. If you don't, follow the [Immediate Setup & Data Persistence](/guides/postgresql/immediate-setup-and-data-persistence/) guide first.

## Internal Network Access (container-to-container)

When your application runs in another container, connecting to PostgreSQL through a user-defined bridge network is the recommended approach. This setup provides automatic DNS resolution, so your application can connect to PostgreSQL using the container name as the hostname, without needing to track IP addresses.

> **Why not use the default bridge network?** While containers on the default bridge network can communicate, they can only do so by IP address. Since container IP addresses change when containers restart, this would require updating your PostgreSQL connection strings each time. User-defined bridge networks solve this by providing automatic DNS resolution, ensuring your PostgreSQL connection strings remain stable even if containers restart and receive new IP addresses.

**Here's a quick comparison:**

> **Note**: The examples below show the difference in approach. To actually test this, follow the steps in this guide to set up containers on the appropriate networks first.

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

With a **user-defined network**, you simply use the container name:
```bash
# Container name works directly - no IP lookup needed
docker run --rm -it \
  --network my-app-net \
  -e PGPASSWORD=mysecretpassword \
  postgres:18 \
  psql -h postgres-dev -U postgres
```

### Step 1: Create a user-defined network

```bash
docker network create my-app-net

# Example Output
ab7f984be43a0ca15534a9ee568716ddbe869a5875077fad3ef3192e3af7d288

docker network ls
# Output
ab7f984be43a   my-app-net    bridge    local


```

### Step 2: Run PostgreSQL on that network (no port publishing)

Notice there is no `-p 5432:5432` here. This keeps PostgreSQL internal to Docker and not accessible from the host machine, which is more secure for production environments.

```bash
docker run -d --name postgres-dev \
  --network my-app-net \
  -e POSTGRES_PASSWORD=mysecretpassword \
  -v postgres_data:/var/lib/postgresql \
  postgres:18

  # Outout
CONTAINER ID  IMAGE        COMMAND                 CREATED         STATUS        PORTS     NAMES
6d351ed89efc  postgres:18  "docker-entrypoint.s…"  9 seconds ago   Up 8 seconds  5432/tcp  postgres-dev

```

### Step 3: Connect from another container using the Postgres container name

You can test connectivity with a temporary psql client container:

```bash
docker run --rm -it \
  --network my-app-net \
  -e PGPASSWORD=mysecretpassword \
  postgres:18 \
  psql -h postgres-dev -U postgres
```

**Key point**: `-h postgres-dev` works because Docker DNS resolves the container name on a user-defined network. The container name acts as the hostname.

### Connection string examples

When connecting from your application container, use these PostgreSQL connection strings:

- **PostgreSQL URI format**: 
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
  

- **PostgreSQL connection parameters**: 
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

- **Connecting to a specific database**: 
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

> **PostgreSQL connection note**: The default port `5432` is used in these examples. If you're connecting to a different PostgreSQL instance or have changed the port, update the connection string accordingly. The container name (`postgres-dev`) is resolved by Docker DNS to the container's IP address on the network.


## Connecting from the Host (external access)

To connect to PostgreSQL from your host machine using tools like `psql`, pgAdmin, DBeaver, or database management scripts, you need to publish PostgreSQL's port (`5432`) to the host. This allows external tools to reach the PostgreSQL container.

### Expose Postgres to localhost only (recommended for development)

This binds to `127.0.0.1` so it's only reachable from your local machine, not from other devices on your network. This is the most secure option for development.

```bash
docker run -d --name postgres-dev \
  -e POSTGRES_PASSWORD=mysecretpassword \
  -p 127.0.0.1:5432:5432 \
  -v postgres_data:/var/lib/postgresql \
  postgres:18
```

Now connect from your host:

- **Host**: `localhost` or `127.0.0.1`
- **Port**: `5432`

If you have `psql` installed on your host:
```bash
psql -h localhost -p 5432 -U postgres
```

You'll be prompted for the password. Alternatively, you can use the `PGPASSWORD` environment variable:
```bash
PGPASSWORD=mysecretpassword psql -h localhost -p 5432 -U postgres
```

### Connecting with PostgreSQL GUI tools

Popular PostgreSQL GUI tools can connect using these common connection details: Host: `localhost`, Port: `5432`, User: `postgres`, Database: `postgres` (or your database name).

- **pgAdmin**: A web-based PostgreSQL administration and development platform
- **DBeaver**: A universal database tool that supports PostgreSQL and many other databases. Select PostgreSQL as the connection type
- **TablePlus**: A modern, native database management tool for macOS and Windows with a clean interface

All tools will prompt for the password you set with `POSTGRES_PASSWORD`.

### Expose Postgres to all network interfaces (use with caution)

To allow connections from other devices on your network, use `-p 5432:5432` instead of `-p 127.0.0.1:5432:5432`. This binds PostgreSQL to all network interfaces on your host, making it accessible from any device that can reach your host, not just localhost.

```bash
docker run -d --name postgres-dev \
  -e POSTGRES_PASSWORD=mysecretpassword \
  -p 5432:5432 \
  -v postgres_data:/var/lib/postgresql \
  postgres:18
```

> **Warning**: Exposing PostgreSQL to all network interfaces (`0.0.0.0:5432`) makes it accessible from any device that can reach your host. Only use this in trusted network environments or behind a firewall. For production, consider using a reverse proxy or VPN instead.

### PostgreSQL security considerations for external access

When exposing PostgreSQL to external access, follow these PostgreSQL-specific security practices:

- **Avoid using the `postgres` superuser**: The default `postgres` user has full database privileges. Create dedicated users with only the permissions your application needs.
- **Use strong passwords**: PostgreSQL passwords should be complex. Consider using environment variables or secrets management instead of hardcoding passwords.
- **Limit network exposure**: Binding to `127.0.0.1` (localhost only) is safer than exposing to all interfaces (`0.0.0.0`).
- **Consider SSL/TLS**: For production, configure PostgreSQL to require SSL connections. The [Advanced Configuration and Initialization](/guides/postgresql/advanced-configuration-and-initialization/) guide shows how to configure PostgreSQL settings.
- **Create application-specific users**: Use initialization scripts to create users with limited privileges. For example, a read-only user for reporting or a user that can only access specific databases.

The [Advanced Configuration and Initialization](/guides/postgresql/advanced-configuration-and-initialization/) guide shows how to use initialization scripts to create users and roles automatically.

## Using Docker Compose for networking

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

**PostgreSQL connection details for the app service**:
- Hostname: `db` (resolved by Docker DNS)
- Port: `5432` (PostgreSQL default port)
- Database: `mydb` (as specified in the connection string)
- User: `postgres` (or a custom user you've created)

> **Note**: Docker Compose automatically creates a network for your project. Services can reach each other by service name without explicit network configuration, but defining a custom network gives you more control. For PostgreSQL, this means your application can always connect using the service name, regardless of container restarts or IP changes.

## Troubleshooting

This section covers common PostgreSQL connection issues and their solutions when working with Docker networking.

### "Could not translate host name postgres-dev"

- Both containers must be on the same Docker network (`my-app-net`).
- Verify the network exists: `docker network ls`
- Check which network a container is on: `docker inspect postgres-dev | grep NetworkMode`
- Ensure you're using a user-defined network, not the default bridge network

### "Connection refused" or "could not connect to server"

- **PostgreSQL may still be initializing**: PostgreSQL takes a few seconds to start and initialize the database cluster. Wait 5-10 seconds after container start and retry.
- **Check if the PostgreSQL container is running**:

```bash
docker ps --filter name=postgres-dev
```

- **Check PostgreSQL logs for initialization or connection errors**:

```bash
docker logs postgres-dev
```

Look for messages like "database system is ready to accept connections" to confirm PostgreSQL is fully started.

- **Verify the port mapping is correct**: 

```bash
docker port postgres-dev
```

This should show `5432/tcp -> 127.0.0.1:5432` (or `0.0.0.0:5432` if bound to all interfaces).

- **Test PostgreSQL connectivity from inside the container**:

```bash
docker exec -it postgres-dev psql -U postgres -c "SELECT version();"
```

If this works but external connections fail, the issue is with port publishing, not PostgreSQL itself.

### "Password authentication failed" or "FATAL: password authentication failed for user"

- **Confirm the password**: Verify you're using the same password set in `POSTGRES_PASSWORD` when you started the container.
- **Existing volume with old credentials**: If you reused an existing volume, the password from the original initialization is still in effect. The `POSTGRES_PASSWORD` environment variable only sets the password during the first database initialization. To reset:
  - Remove the volume: `docker volume rm postgres_data`
  - Or connect with the old password
  - Or change the password after connecting: `ALTER USER postgres WITH PASSWORD 'newpassword';`
- **Try connecting with password prompt**: `psql -h localhost -U postgres -W` (the `-W` flag forces a password prompt)
- **Use PGPASSWORD environment variable**: `PGPASSWORD=mysecretpassword psql -h localhost -U postgres`
- **Check PostgreSQL authentication configuration**: If you've customized `pg_hba.conf`, verify the authentication method allows password authentication

### "Network not found"

- Ensure the network exists before starting containers: `docker network create my-app-net`
- If using Docker Compose, the network is created automatically when you run `docker compose up`