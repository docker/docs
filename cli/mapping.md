---
title: Mapping ports and volumes
description: |
  Quick overview of the syntax you use when you want to map host resources to
  the container, and container ports to the host.
keywords: 
---

The Docker CLI uses a special mapping syntax for the when you want to mount files
into a container, or publish ports from the container to the host.

The syntax for these mappings uses the pattern `HOST:CONTAINER` where:

- `HOST` is the location of the mapping on the host
- `CONTAINER` is the location in the container

## Publish ports

To publish a container port to a port on the host, you use the `--publish` flag
with the `docker run` command. The format of the flag is
`HOST_PORT:CONTAINER_PORT`.

For example, to run an Nginx server in a container, and publish it's HTTP port
`80` to the address `localhost:8080` on the host, you would use:

```console
docker run --publish 127.0.0.1:8080:80 nginx
```

## Bind mounts

Bind mounts are when you mount a part of the filesystem on the host machine
into a container. You can use the `--volume` (short form `-v`) flag with the
`docker run` command to make files from your host show up in a container.

For example, to mount the current working directory on your host's filesystem
into the `/app` directory in a container, you would use:

```console
$ docker run -v .:/app <image>
```

Alternatively, when mapping the host filesystem into containers, you can use
the `--mount` flag to create a bind mount, which is more explicit:

```console
$ docker run --mount type=bind,source=.,target=/app <image>
```

> **Note**
>
> Docker Engine 20.10 and earlier versions don't support relative paths for the
> `-v` and `--mount` flags.
