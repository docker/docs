---
title: Companions for PostgreSQL
linkTitle: Companions for PostgreSQL
description: This module explains how to customize PostgreSQL for real-world use in Docker, covering automated database initialization, performance tuning, and timezone configuration once persistent storage is in place.
keywords:
  - PostgreSQL Docker
  - Docker Compose PostgreSQL
  - container database
weight: 40
---


## PostgreSQL Ecosystem Companions: pgAdmin, PgBouncer, and Performance Testing

Running a standalone PostgreSQL container is often just the beginning. What happens when thousands of connections arrive, or when you need a visual interface to manage your database?

This is where **companion tools** come into play. These applications extend PostgreSQL with capabilities the core database engine doesn't provide natively: visual administration, connection pooling, and performance benchmarking. This guide covers how to deploy pgAdmin 4, PgBouncer, Pgpool-II, and pgbench in Docker, when to use each tool, and real-world benchmark results demonstrating their performance impact.

## pgAdmin 4: Visual Management Platform

pgAdmin 4 is the industry-standard open source management tool for PostgreSQL. When deployed in Docker, it typically runs in **Server Mode**, providing a multi-user web interface to manage one or more database instances.

While you can accomplish everything from the command line using `psql`, a visual interface significantly simplifies writing complex queries, visualizing table structures, and exploring database objects.

### Key Considerations

When running pgAdmin in Docker, keep these points in mind:

- **Image**: Use the official `dpage/pgadmin4` image
- **Networking**: In a Docker Compose environment, pgAdmin connects to the database using the internal service name (for example, `db:5432`) rather than `localhost`

### Docker Compose Configuration

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

With this configuration, access the pgAdmin interface at `http://localhost:8080`. Use the email and password specified in the environment variables for initial login.

> **Important**: In production environments, pass `PGADMIN_DEFAULT_PASSWORD` as an external environment variable or use Docker secrets. Storing passwords in plain text within `docker-compose.yml` poses a security risk.

Now that you have visual database management in place, the next challenge in production environments is handling connection load. The following section explains how to manage high-volume database traffic.

## PgBouncer: Lightweight Connection Pooling

PostgreSQL creates a new process for every client connection, which consumes significant RAM. What happens when you have 1,000 concurrent users? PgBouncer solves exactly this problem.

PgBouncer is a lightweight proxy that pools connections, allowing thousands of applications to share a small number of actual database backends. Think of it as a traffic controller: everyone wants to pass through simultaneously, but the controller regulates the flow to prevent congestion.

### Pooling Modes

PgBouncer offers three distinct pooling modes:

| Mode | Description | Use Case |
|------|-------------|----------|
| **Session** | Connection assigned for entire session duration | Long-lived connections, session variables |
| **Transaction** | Connection returned after each transaction ends | Web applications, microservices (most common) |
| **Statement** | Connection returned after every SQL statement | Simple queries, no multi-statement transactions |

### When to Use PgBouncer

PgBouncer becomes essential when you encounter:

- "too many connections" errors
- High memory consumption due to connection overhead
- Many short-lived connections (web applications, serverless functions)
- Need to serve thousands of clients with limited database connections

### Complete Docker Compose Setup

To run PostgreSQL and PgBouncer together, you need three files: `docker-compose.yml`, `pgbouncer.ini`, and `userlist.txt`.

First, create the PgBouncer configuration file (`pgbouncer.ini`):

```ini
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

```
"postgres" "postgres"
```

Finally, create the Docker Compose file (`docker-compose.yml`):

```yaml
services:
  postgres:
    image: postgres:17.7
    container_name: postgres
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: benchmark
      POSTGRES_HOST_AUTH_METHOD: trust
    volumes:
      - postgres_data:/var/lib/postgresql/data
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

- PgBouncer listens on port **6432**, avoiding confusion with the direct PostgreSQL connection on port 5432
- The `depends_on` directive with `service_healthy` condition ensures PgBouncer starts only after PostgreSQL is ready
- `pool_mode = transaction` is the optimal choice for most web applications
- The Percona PgBouncer image requires mounted configuration files (without the `:ro` flag, as the entrypoint script needs to modify them)
- This example uses `trust` authentication for simplicity. In production, configure proper SCRAM-SHA-256 authentication

> **Note**: The Percona PgBouncer entrypoint script processes the configuration files on startup. Mount them without the read-only flag to avoid permission errors.




## pgbench: Performance Benchmarking

pgbench is a benchmarking utility included with the official PostgreSQL image. It allows you to simulate heavy workloads and verify how your Docker configuration performs under pressure.

### Initialize Benchmark Tables

First, create the test tables. The `-s` (scale) parameter determines data size—scale factor 50 creates approximately 5 million rows:

```bash
docker exec postgres pgbench -i -s 50 -U postgres benchmark
```

### Run Stress Tests

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

## Understanding Benchmark Results

Does PgBouncer actually make a difference? Run the benchmarks yourself to find out. Your results will vary based on your hardware, Docker configuration, network setup, and system load.

### What to Expect

When you run these benchmarks, you'll observe patterns rather than specific numbers. Think of it like comparing two different routes to work: the "faster" route depends on traffic conditions, time of day, and your vehicle.

### Key Observations

When comparing direct connections versus PgBouncer, you'll typically notice:

**1. Connection overhead differs significantly**

Direct connections require PostgreSQL to spawn a new process for each client. PgBouncer reuses existing connections. Watch the "initial connection time" metric in your results—PgBouncer often shows dramatically faster connection setup.

**2. Behavior under pressure reveals the real difference**

Try increasing the client count (`-c` parameter) gradually: 50, 100, 150, 200. At some point, direct connections will fail with "too many clients already" while PgBouncer continues handling requests. This is PgBouncer's primary value: **it prevents connection exhaustion**.

**3. Throughput varies by environment**

On some systems, direct connections show higher transactions per second (TPS) at low concurrency. On others, PgBouncer wins even with few clients. The difference depends on:
- CPU and memory available
- Docker networking overhead
- Disk I/O speed
- Whether connections are being rapidly opened and closed

