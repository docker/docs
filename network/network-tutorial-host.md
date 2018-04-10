---
title: Networking using the host network
description: Tutorials for networking using the host network, disabling network isolation
keywords: networking, host, standalone
---

This series of tutorials deals with networking standalone containers which bind
directly to the Docker host's network, with no network isolation. For other
networking topics, see the [overview](index.md).

## Goal

The goal of this tutorial is to start a `nginx` container which binds directly
to port 80 on the Docker host. From a networking point of view, this is the
same level of isolation as if the `nginx` process were running directly on the
Docker host and not in a container. However, in all other ways, such as storage,
process namespace, and user namespace, the `nginx` process is isolated from the
host.

## Prerequisites

- This procedure requires port 80 to be available on the Docker host. To make
  Nginx listen on a different port, see the
  [documentation for the `nginx` image](https://hub.docker.com/_/nginx/)

- The `host` networking driver only works on Linux hosts, and is not supported
  on Docker for Mac, Docker for Windows, or Docker EE for Windows Server.

## Procedure

1.  Create and start the container as a detached process.

    ```bash
    docker run --rm -dit --network host --name my_nginx nginx
    ```

2.  Access Nginx by browsing to
    [http://localhost:80/](http://localhost:80/).

3.  Examine your network stack using the following commands:

    - Examine all network interfaces and verify that a new one was not created.

      ```bash
      ip addr show
      ```

    - Verify which process is bound to port 80, using the `netstat` command. You
      need to use `sudo` because the process is owned by the Docker daemon user
      and you otherwise won't be able to see its name or PID.

      ```bash
      sudo netstat -tulpn | grep :80
      ```

4.  Stop the container.

    ```basn
    docker container stop my_nginx
    docker container rm my_nginx
    ```

## Other networking tutorials

Now that you have completed the networking tutorials for standalone containers,
you might want to run through these other networking tutorials:

- [Standalone networking tutorial](network-tutorial-standalone.md)
- [Overlay networking tutorial](network-tutorial-overlay.md)
- [Macvlan networking tutorial](network-tutorial-macvlan.md)

