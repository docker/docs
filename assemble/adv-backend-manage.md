---
title: Advanced backend management
description: Advanced backend management for Docker Assemble
keywords: Backend, Assemble, Docker Enterprise, plugin, Spring Boot, .NET, c#, F#
---

## Backend access to host ports

Docker Assemble requires its own buildkit instance to be running in a Docker container on the local system. You can start and manage the backend using the `backend` subcommand of `docker assemble`. For more information, see [Install Docker Assemble](/install).

As the backend runs in a container with its own network namespace, it cannot access host resources directly. This is most noticeable when trying to push to a local registry as     `localhost:5000`.

The backend supports a sidecar container which proxies ports from within the backend container to the container's gateway (which is in effect a host IP). This is sufficient to allow access to host ports which have been bound to `0.0.0.0` (or to the gateway specifically), but not ones which are bound to `127.0.0.1`.

By default, port 5000 is proxied in this way, as that is the most common port used for a local registry to allow access to a local registry on `localhost:5000` (the most common setup). You can proxy other ports using the `--allow-host-port` option to docker assemble backend start.

For example, to expose port `6000` instead of port `5000`, run:

```
$ docker assemble backend start --allow-host-port 6000
```
> **Notes:**
>
> - You can repeat the `--allow-host-port` option or give it a comma separated list of ports.
> - Passing `--allow-host-port 0` disables the default and no ports are exposed. For example:
>
>    `$ docker assemble backend start --allow-host-port 0`
> - On Docker Desktop, this functionality allows the backend to access ports on the Docker Desktop VM host, rather than the Windows or macOS host. To access the the Windows or macOS host port, you can use `host.docker.internal` as usual.

## Backend sub-commands

### Info

The info sub-command describes the backend:

```
~$ docker assemble backend info
ID: 2f03e7d288e6bea770a2acba4c8c918732aefcd1946c94c918e8a54792e4540f (running)
Image: docker/assemble-backend@sha256:«…»

Sidecar containers:
 - 0f339c0cc8d7 docker-assemble-backend-username-proxy-port-5000 (running)

Found 1 worker(s):

 - 70it95b8x171u5g9jbixkscz9
   Platforms:
    - linux/amd64
   Labels:
    - com.docker.assemble.commit: «…»
    - org.mobyproject.buildkit.worker.executor: oci
    - org.mobyproject.buildkit.worker.hostname: 2f03e7d288e6
    - org.mobyproject.buildkit.worker.snapshotter: overlayfs

Build cache contains 54 entries, total size 3.65GB (0B currently in use)
```

### Stop

The stop sub-command destroys the backend container

```
~$ docker assemble backend stop
```

### Logs

The logs sub-command displays the backend logs.

```
~$ docker assemble backend logs
```

### Cache

The build cache gets lost when the backend is stopped. To avoid this, you can create a volume named `docker-assemble-backend-cache-«username»` and it will automatically be used as the build cache.

Alternatively you can specify a named docker volume to use for the cache. For example:

```
~$ docker volume create $USER-assemble-cache
username-assemble-cache
~$ docker assemble backend start --cache-volume=username-assemble-cache
Pulling image «…»: Success
Started container "docker-assemble-backend-username" (74476d3fdea7)
```

For information regarding the current cache contents, run the command `docker assemble backend cache usage`.

To clean the cache, run `docker assemble backend cache purge`.
