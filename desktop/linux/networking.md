---
description: Networking
keywords: linux, networking
title: Networking features in Docker Desktop for Linux
---
{% assign Arch = 'Linux' %}

Docker Desktop for {{Arch}} provides several networking features to make it
easier to use.

## Features

### VPN Passthrough

Docker Desktop for {{Arch}}'s networking can work when attached to a VPN. To do this,
Docker Desktop for {{Arch}} intercepts traffic from the containers and injects it into
{{Arch}} as if it originated from the Docker application.

### Port Mapping

When you run a container with the `-p` argument, for example:

```console
$ docker run -p 80:80 -d nginx
```

Docker Desktop for {{Arch}} makes whatever is running on port 80 in the container (in
this case, `nginx`) available on port 80 of `localhost`. In this example, the
host and container ports are the same. What if you need to specify a different
host port? If, for example, you already have something running on port 80 of
your host machine, you can connect the container to a different port:

```console
$ docker run -p 8000:80 -d nginx
```

Now, connections to `localhost:8000` are sent to port 80 in the container. The
syntax for `-p` is `HOST_PORT:CLIENT_PORT`.

### HTTP/HTTPS Proxy Support

See [Proxies](index.md#proxies).

### SSH agent forwarding

Docker Desktop for Linux allows you to use the host’s SSH agent inside a container. To do this:

1. Bind mount the SSH agent socket by adding the following parameter to your `docker run` command:

    ```ssh
    --mount type=bind,src=/run/host-services/ssh-auth.sock,target=/run/host-services/ssh-auth.sock
   ```


2. Add the `SSH_AUTH_SOCK` environment variable in your container:

   `-e SSH_AUTH_SOCK="/run/host-services/ssh-auth.sock"`


To enable the SSH agent in Docker Compose, add the following flags to your service:

 ```yaml
services:
  web:
    image: nginx:alpine
    volumes:
      - type: bind
        source: /run/host-services/ssh-auth.sock
        target: /run/host-services/ssh-auth.sock
    environment:
      - SSH_AUTH_SOCK=/run/host-services/ssh-auth.sock
 ```

## Known limitations, use cases, and workarounds

Below is a summary of the current limitations on the Docker Desktop for {{Arch}}
networking stack, along with some ideas for workarounds.

### Changing internal IP addresses

The internal IP addresses used by Docker can be changed from the **Settings** panel.
After changing IPs, you need to reset the Kubernetes cluster and leave any active Swarm.

### I cannot ping my containers

Docker Desktop for Linux can't route traffic to containers.

### Per-container IP addressing is not possible

The docker bridge network is not reachable from the Linux host.

### Use cases and workarounds

There are two scenarios that the above limitations affect:

#### I want to connect from a container to a service on the host

The host has a changing IP address (or none if you have no network access). We recommend that you connect to the special DNS name
`host.docker.internal` which resolves to the internal IP address used by the
host. This is for development purposes and will not work in a production environment outside of Docker Desktop for Linux.

You can also reach the gateway using `gateway.docker.internal`.

If you have installed Python on your machine, use the following instructions as an example to connect from a container to a service on the host:

1. Run the following command to start a simple HTTP server on port 8000.

   ```console
    $ python -m http.server 8000
   ```

   If you have installed Python 2.x, run:

   ```console
    $ python -m SimpleHTTPServer 8000
   ```

2. Now, run a container, install `curl`, and try to connect to the host using the following commands:

    ```console
    $ docker run --rm -it alpine sh
    # apk add curl
    # curl http://host.docker.internal:8000
    # exit
    ```

#### I want to connect to a container from the Linux host

Port forwarding works for `localhost`; `--publish`, `-p`, or `-P` all work.
Ports exposed from the Linux VM are forwarded to the host.

Our current recommendation is to publish a port, or to connect from another
container. This is what you need to do even on Linux if the container is on an
overlay network, not a bridge network, as these are not routed.

For example, to run an `nginx` webserver:

```console
$ docker run -d -p 80:80 --name webserver nginx
```

To clarify the syntax, the following two commands both expose port `80` on the
container to port `8000` on the host:

```console
$ docker run --publish 8000:80 --name webserver nginx

$ docker run -p 8000:80 --name webserver nginx
```

To expose all ports, use the `-P` flag. For example, the following command
starts a container (in detached mode) and the `-P` exposes all ports on the
container to random ports on the host.

```console
$ docker run -d -P --name webserver nginx
```

See the [run command](../../engine/reference/commandline/run.md) for more details on
publish options used with `docker run`.
