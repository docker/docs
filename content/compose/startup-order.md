---
description: How to control service startup and shutdown order in Docker Compose
keywords: documentation, docs, docker, compose, startup, shutdown, order
title: Control startup and shutdown order in Compose
notoc: true
---

You can control the order of service startup and shutdown with the
[depends_on](compose-file/05-services.md#depends_on) attribute. Compose always starts and stops
containers in dependency order, where dependencies are determined by
`depends_on`, `links`, `volumes_from`, and `network_mode: "service:..."`.

A good example of when you might use this is an application which needs to access a database. If both services are started with `docker compose up`, there is a chance this will fail since the application service might start before the database service and won't find a database able to handle its SQL statements. 

## Control startup

On startup, Compose does not wait until a container is "ready", only until it's running. This can cause issues if, for example, you have a relational database system that needs to start its own services before being able to handle incoming connections.

The solution for detecting the ready state of a service is  to use the `condition` attribute with one of the following options:

- `service_started`
- `service_healthy`. This specifies that a dependency is expected to be “healthy”, which is defined with `healthcheck`, before starting a dependent service.
- `service_completed_successfully`. This specifies that a dependency is expected to run to successful completion before starting a dependent service.

## Reference information 

- [`depends_on`](compose-file/05-services.md#depends_on)
- [`healthcheck`](compose-file/05-services.md#healthcheck)