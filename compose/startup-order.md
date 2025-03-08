---
description: How to control service startup order in Docker Compose
keywords: documentation, docs, docker, compose, startup, order
title: Control startup order in Compose
notoc: true
---

You can control the order of service startup with the
[depends_on](compose-file.md#depends-on) option. Compose always starts
containers in dependency order, where dependencies are determined by
`depends_on`, `links`, `volumes_from`, and `network_mode: "service:..."`.

However, Compose does not wait until a container is "ready" (whatever that means
for your particular application) - only until it's running. There's a good
reason for this.

The problem of waiting for a database (for example) to be ready is really just
a subset of a much larger problem of distributed systems. In production, your
database could become unavailable or move hosts at any time. Your application
needs to be resilient to these types of failures.

To handle this, design your application to attempt to re-establish a connection to
the database after a failure. If the application retries the connection,
it can eventually connect to the database.

The best solution is to perform this check in your application code, both at
startup and whenever a connection is lost for any reason. However, if you don't
need this level of resilience, you can work around the problem with a wrapper
script:

-   Use a tool such as [wait-for-it](https://github.com/vishnubob/wait-for-it),
    [dockerize](https://github.com/jwilder/dockerize), or sh-compatible
    [wait-for](https://github.com/Eficode/wait-for). These are small
    wrapper scripts which you can include in your application's image to
    poll a given host and port until it's accepting TCP connections.

    For example, to use `wait-for-it.sh` or `wait-for` to wrap your service's command:

        version: "2"
        services:
          web:
            build: .
            ports:
              - "80:8000"
            depends_on:
              - "db"
            command: ["./wait-for-it.sh", "db:5432", "--", "python", "app.py"]
          db:
            image: postgres

    >**Tip**: There are limitations to this first solution. For example, it doesn't verify when a specific service is really ready. If you add more arguments to the command, use the `bash shift` command with a loop, as shown in the next example.

-   Alternatively, write your own wrapper script to perform a more application-specific health
    check. For example, you might want to wait until Postgres is definitely
    ready to accept commands:

        #!/bin/bash
        # wait-for-postgres.sh

        set -e

        host="$1"
        shift
        cmd="$@"

        until PGPASSWORD=$POSTGRES_PASSWORD psql -h "$host" -U "postgres" -c '\q'; do
          >&2 echo "Postgres is unavailable - sleeping"
          sleep 1
        done

        >&2 echo "Postgres is up - executing command"
        exec $cmd

    You can use this as a wrapper script as in the previous example, by setting:

    ```none
    command: ["./wait-for-postgres.sh", "db", "python", "app.py"]
    ```


## Compose documentation

- [Installing Compose](install.md)
- [Get started with Django](django.md)
- [Get started with Rails](rails.md)
- [Get started with WordPress](wordpress.md)
- [Command line reference](./reference/index.md)
- [Compose file reference](compose-file.md)
