---
title: Check a container's health
description: Use HEALTHCHECK or docker run --health-cmd to mark a container healthy or unhealthy
keywords: containers, healthcheck, HEALTHCHECK, docker run
weight: 15
---

A health check is a command Docker runs inside the container to decide whether
the process is actually serving. The container still has a normal status
(`created`, `running`, `exited`). When a health check is set, it also has a
health status:

- `starting` — checks have not succeeded yet
- `healthy` — the last check succeeded
- `unhealthy` — too many checks in a row failed

Compose `depends_on` with `condition: service_healthy` and Swarm wait for
`healthy` before treating a replica as ready.

## In a Dockerfile

```dockerfile
HEALTHCHECK --interval=5s --timeout=3s --start-period=20s --retries=3 \
  CMD curl -fsS http://127.0.0.1:8080/health || exit 1
```

`--start-period` is grace time after the container starts. Failures in that
window do not count toward `--retries`. A success during the start period
ends it. `--retries` is consecutive failures (the first check counts) and
must be 1 or higher; `0` is ignored and the daemon uses 3.

There can be only one `HEALTHCHECK`. A later one replaces the earlier one.
`HEALTHCHECK NONE` disables a check inherited from the base image.

See the [Dockerfile reference](/reference/dockerfile.md#healthcheck).

## On `docker run`

Flags override the image's `HEALTHCHECK`:

```console
$ docker run -d --name web \
    --health-cmd='curl -fsS http://127.0.0.1:8080/health || exit 1' \
    --health-interval=5s \
    --health-retries=3 \
    --health-start-period=20s \
    nginx:alpine
```

`--health-cmd` always runs through `CMD-SHELL`. Use
`--no-healthcheck` to turn the image check off.

See [`docker run`](/reference/cli/docker/container/run/#health).

## Read the status

```console
$ docker inspect --format='{{.State.Health.Status}}' web
healthy
```

`docker ps` shows `(healthy)` next to the status when a check is configured.
