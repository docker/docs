---
description: Networking on Docker Desktop
keywords: netowrking, docker desktop
title: Explore networking features
redirect_from:
- /desktop/linux/networking/
- /docker-for-mac/networking/
- /mackit/networking/
- /desktop/mac/networking/
- /docker-for-win/networking/
- /docker-for-windows/networking/
- /desktop/windows/networking/
---

Docker Desktop provides several networking features to make it easier to
use.

## Features for all platforms

### VPN Passthrough

Docker Desktop networking can work when attached to a VPN. To do this,
Docker Desktop intercepts traffic from the containers and injects it into
the host as if it originated from the Docker application.

### Port Mapping

When you run a container with the `-p` argument, for example:

```console
$ docker run -p 80:80 -d nginx
```

Docker Desktop makes whatever is running on port 80 in the container (in
this case, `nginx`) available on port 80 of `localhost`. In this example, the
host and container ports are the same. If, for example, you already have something running on port 80 of
your host machine, you can connect the container to a different port:

```console
$ docker run -p 8000:80 -d nginx
```

Now, connections to `localhost:8000` are sent to port 80 in the container. The
syntax for `-p` is `HOST_PORT:CLIENT_PORT`.

### HTTP/HTTPS Proxy Support

See:
- [Proxies on Linux](linux/index.md#proxies)
- [Proxies on Mac](mac/index.md#proxies)
- [Proxies on Windows](windows/index.md#proxies)

## Features for Mac and Linux

### SSH agent forwarding

Docker Desktop on Mac and Linux allows you to use the host’s SSH agent inside a container. To do this:

1. Bind mount the SSH agent socket by adding the following parameter to your `docker run` command:

    `--mount type=bind,src=/run/host-services/ssh-auth.sock,target=/run/host-services/ssh-auth.sock`

1. Add the `SSH_AUTH_SOCK` environment variable in your container:

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

## Known limitations for all platforms

### Changing internal IP addresses

The internal IP addresses used by Docker can be changed from the Settings, if you're a Windows user, or Preferences, if you use Mac or Linux. After changing IPs, it is necessary to reset the Kubernetes cluster and to leave any active Swarm.

### There is no docker0 bridge on the host

Because of the way networking is implemented in Docker Desktop, you cannot
see a `docker0` interface on the host.  This interface is actually within the
virtual machine.

### I cannot ping my containers

Docker Desktop can't route traffic to Linux containers.  However if you're a Windows user, you can
ping the Windows containers.

### Per-container IP addressing is not possible

The docker bridge network is not reachable from the host.
However if you are a Windows user, it works with Windows containers.

## Use cases and workarounds for all platforms

### I want to connect from a container to a service on the host

The host has a changing IP address (or none if you have no network access). We recommend that you connect to the special DNS name
`host.docker.internal` which resolves to the internal IP address used by the
host. This is for development purpose and does not work in a production environment outside of Docker Desktop.

You can also reach the gateway using `gateway.docker.internal`.

If you have installed Python on your machine, use the following instructions as an example to connect from a container to a service on the host:

1. Run the following command to start a simple HTTP server on port 8000.

    `python -m http.server 8000`

    If you have installed Python 2.x, run `python -m SimpleHTTPServer 8000`.

2. Now, run a container, install `curl`, and try to connect to the host using the following commands:

    ```console
    $ docker run --rm -it alpine sh
    # apk add curl
    # curl http://host.docker.internal:8000
    # exit
    ```

#### I want to connect to a container from the host

Port forwarding works for `localhost`; `--publish`, `-p`, or `-P` all work.
Ports exposed from Linux are forwarded to the host.

Our current recommendation is to publish a port, or to connect from another
container. This is what you need to do even on Linux if the container is on an
overlay network, not a bridge network, as these are not routed.

For example, to run an `nginx` webserver:

```console
$ docker run -d -p 80:80 --name webserver nginx
```

To clarify the syntax, the following two commands both publish container's port `80` to host's port `8000`:

```console
$ docker run --publish 8000:80 --name webserver nginx

$ docker run -p 8000:80 --name webserver nginx
```

To publish all ports, use the `-P` flag. For example, the following command
starts a container (in detached mode) and the `-P` flag publishes all exposed ports of the
container to random ports on the host.

```console
$ docker run -d -P --name webserver nginx
```

See the [run command](../engine/reference/commandline/run.md) for more details on
publish options used with `docker run`.

