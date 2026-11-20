---
title: Workload specification
description: Explore all the workload attributes shared by services and jobs.
keywords: compose, compose specification, workload, services, jobs, compose file reference
weight: 22
---

A workload specification defines the attributes of an orchestrated workload, which is a
[service](services.md) or a [job](jobs.md). These attributes go beyond the configuration of the
workload's containers, which the [container specification](container-spec.md) defines. They cover
how Compose builds the image, how the workload connects to other workloads, and how the platform
checks and exposes it.

These attributes apply to services and jobs, but not to other container-based elements such as
[`pre_start` init containers](services.md#pre_start).

## Attributes

<!-- vale off(Docker.HeadingSentenceCase.yml) -->

### `build`

`build` specifies the build configuration for creating a container image from source, as defined in the [Compose Build Specification](build.md).

### `depends_on`

{{% include "compose/services-depends-on.md" %}}

A service can only depend on other services. It can't depend on a [job](jobs.md).

#### Short syntax

The short syntax variant only specifies service names of the dependencies.
Service dependencies cause the following behaviors:

- Compose creates services in dependency order. In the following
  example, `db` and `redis` are created before `web`.

- Compose removes services in dependency order. In the following
  example, `web` is removed before `db` and `redis`.

Simple example:

```yml
services:
  web:
    build: .
    depends_on:
      - db
      - redis
  redis:
    image: redis
  db:
    image: postgres:18
```

Compose guarantees dependency services have been started before
starting a dependent service.
With short syntax, Compose does not wait for dependency services to be "healthy" before
starting a dependent service.

#### Long syntax

The long form syntax enables the configuration of additional fields that can't be
expressed in the short form.

- `restart`: When set to `true` Compose restarts this service after it updates the dependency service.
  This applies to an explicit restart controlled by a Compose operation, and excludes automated restart by the container runtime
  after the container dies. Introduced in Docker Compose version [2.17.0](https://github.com/docker/compose/releases/tag/v2.17.0).

- `condition`: Sets the condition under which dependency is considered satisfied
  - `service_started`: An equivalent of the short syntax described previously
  - `service_healthy`: Specifies that a dependency is expected to be "healthy"
    (as indicated by [`healthcheck`](#healthcheck)) before starting a dependent
    service.
  - `service_completed_successfully`: Specifies that a dependency is expected to run
    to successful completion before starting a dependent service.
- `required`: When set to `false` Compose only warns you when the dependency service isn't started or available. If it's not defined
    the default value of `required` is `true`. Introduced in Docker Compose version [2.20.0](https://github.com/docker/compose/releases/tag/v2.20.0).

Service dependencies cause the following behaviors:

- Compose creates services in dependency order. In the following
  example, `db` and `redis` are created before `web`.

- Compose waits for healthchecks to pass on dependencies
  marked with `service_healthy`. In the following example, `db` is expected to
  be "healthy" before `web` is created.

- Compose removes services in dependency order. In the following
  example, `web` is removed before `db` and `redis`.

```yml
services:
  web:
    build: .
    depends_on:
      db:
        condition: service_healthy
        restart: true
      redis:
        condition: service_started
  redis:
    image: redis
  db:
    image: postgres:18
```

Compose guarantees dependency services are started before
starting a dependent service.
Compose guarantees dependency services marked with
`service_healthy` are "healthy" before starting a dependent service.

### `expose`

`expose` defines the (incoming) port or a range of ports that Compose exposes from the container. These ports must be
accessible to linked services and should not be published to the host machine. Only the internal container
ports can be specified.

Syntax is `<portnum>/[<proto>]` or `<startport-endport>/[<proto>]` for a port range.
When not explicitly set, `tcp` protocol is used.

```yml
expose:
  - "3000"
  - "8000"
  - "8080-8085/tcp"
```

> [!NOTE]
>
> If the Dockerfile for the image already exposes ports, it is visible to other containers on the network even if `expose` is not set in your Compose file.

### `healthcheck`

{{% include "compose/services-healthcheck.md" %}}

For more information on `HEALTHCHECK`, see the [Dockerfile reference](/reference/dockerfile.md#healthcheck).

```yml
healthcheck:
  test: ["CMD", "curl", "-f", "http://localhost"]
  interval: 1m30s
  timeout: 10s
  retries: 3
  start_period: 40s
  start_interval: 5s
```

`interval`, `timeout`, `start_period`, and `start_interval` are [specified as durations](extension.md#specifying-durations). Introduced in Docker Compose version [2.20.2](https://github.com/docker/compose/releases/tag/v2.20.2)

`test` defines the command Compose runs to check container health. It can be
either a string or a list. If it's a list, the first item must be either `NONE`, `CMD` or `CMD-SHELL`.
If it's a string, it's equivalent to specifying `CMD-SHELL` followed by that string.

```yml
# Hit the local web app
test: ["CMD", "curl", "-f", "http://localhost"]
```

Using `CMD-SHELL` runs the command configured as a string using the container's default shell
(`/bin/sh` for Linux). Both of the following forms are equivalent:

```yml
test: ["CMD-SHELL", "curl -f http://localhost || exit 1"]
```

```yml
test: curl -f https://localhost || exit 1
```

`NONE` disables the healthcheck, and is mostly useful to disable the Healthcheck Dockerfile instruction set by the service's Docker image. Alternatively,
the healthcheck set by the image can be disabled by setting `disable: true`:

```yml
healthcheck:
  disable: true
```

### `ports`

{{% include "compose/services-ports.md" %}}

> [!NOTE]
>
> Port mapping must not be used with `network_mode: host`. Doing so causes a runtime error because `network_mode: host` already exposes container ports directly to the host network, so port mapping isn’t needed.

#### Short syntax

The short syntax is a colon-separated string to set the host IP, host port, and container port
in the form:

`[HOST:]CONTAINER[/PROTOCOL]` where:

- `HOST` is `[IP:](port | range)` (optional). If it is not set, it binds to all network interfaces (`0.0.0.0`).
- `CONTAINER` is `port | range`.
- `PROTOCOL` restricts ports to a specified protocol either `tcp` or `udp`(optional). Default is `tcp`.

> [!WARNING]
>
> If you do not specify a host IP (such as `127.0.0.1`), Docker binds to all interfaces (`0.0.0.0`), bypassing host firewall rules. This can expose the container directly to the internet if the host has a public IP address. For more information, see [Port publishing and mapping](/manuals/engine/network/port-publishing.md).

Ports can be either a single value or a range. `HOST` and `CONTAINER` must use equivalent ranges.

You can either specify both ports (`HOST:CONTAINER`), or just the container port. In the latter case,
the container runtime automatically allocates any unassigned port of the host.

`HOST:CONTAINER` should always be specified as a (quoted) string, to avoid conflicts
with [YAML base-60 float](https://yaml.org/type/float.html).



IPv6 addresses can be enclosed in square brackets.

Examples:

```yml
ports:
  - "3000"
  - "3000-3005"
  - "8000:8000"
  - "9090-9091:8080-8081"
  - "49100:22"
  - "8000-9000:80"
  - "127.0.0.1:8001:8001"
  - "127.0.0.1:5000-5010:5000-5010"
  - "::1:6000:6000"
  - "[::1]:6001:6001"
  - "6060:6060/udp"
```

> [!NOTE]
>
> If host IP mapping is not supported by a container engine, Compose rejects
> the Compose file and ignores the specified host IP.

#### Long syntax

The long form syntax lets you configure additional fields that can't be
expressed in the short form.

- `target`: The container port.
- `published`: The publicly exposed port. It is defined as a string and can be set as a range using syntax `start-end`. It means the actual port is assigned a remaining available port, within the set range.
- `host_ip`: The host IP mapping. If it is not set, it binds to all network interfaces (`0.0.0.0`).
- `protocol`: The port protocol (`tcp` or `udp`). Defaults to `tcp`.
- `app_protocol`: The application protocol (TCP/IP level 4 / OSI level 7) this port is used for. This is optional and can be used as a hint for Compose to offer richer behavior for protocols that it understands. Introduced in Docker Compose version [2.26.0](https://github.com/docker/compose/releases/tag/v2.26.0).
- `mode`: Specifies how the port is published in a Swarm setup. If set to `host`, it publishes the port on every node in Swarm. If set to `ingress`, it allows load balancing across the nodes in Swarm. Defaults to `ingress`.
- `name`: A human-readable name for the port, used to document its usage within the service.

```yml
ports:
  - name: web
    target: 80
    host_ip: 127.0.0.1
    published: "8080"
    protocol: tcp
    app_protocol: http
    mode: host

  - name: web-secured
    target: 443
    host_ip: 127.0.0.1
    published: "8083-9000"
    protocol: tcp
    app_protocol: https
    mode: host
```

### `stdin_open`

`stdin_open` configures a service's container to run with an allocated stdin. This is the same as running a container with the
`-i` flag. For more information, see [Keep stdin open](/reference/cli/docker/container/run/#interactive).

Supported values are `true` or `false`.

### `tty`

`tty` configures a service's container to run with a TTY. This is the same as running a container with the
`-t` or `--tty` flag. For more information, see [Allocate a pseudo-TTY](/reference/cli/docker/container/run/#tty).

Supported values are `true` or `false`.
