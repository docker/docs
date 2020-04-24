---
title: Use host networking
description: All about exposing containers on the Docker host's network
keywords: network, host, standalone
---

If you use the `host` network mode for a container, that container's network
stack is not isolated from the Docker host (the container shares the host's
networking namespace), and the container does not get its own IP-address allocated.
For instance, if you run a container which binds to port 80 and you use `host`
networking, the container's application is available on port 80 on the host's IP
address.

> **Note**: Given that the container does not have its own IP-address when using
> `host` mode networking, [port-mapping](overlay.md#publish-ports) does not
> take effect, and the `-p`, `--publish`, `-P`, and `--publish-all` option are
> ignored, producing a warning instead:
>
> ```
> WARNING: Published ports are discarded when using host network mode
> ```

Host mode networking can be useful to optimize performance, and in situations where
a container needs to handle a large range of ports, as it does not require network
address translation (NAT), and no "userland-proxy" is created for each port.

The host networking driver only works on Linux hosts, and is not supported on
Docker Desktop for Mac, Docker Desktop for Windows, or Docker EE for Windows Server.

You can also use a `host` network for a swarm service, by passing `--network host`
to the `docker service create` command. In this case, control traffic (traffic
related to managing the swarm and the service) is still sent across an overlay
network, but the individual swarm service containers send data using the Docker
daemon's host network and ports. This creates some extra limitations. For instance,
if a service container binds to port 80, only one service container can run on a
given swarm node.

## Next steps

- Go through the [host networking tutorial](network-tutorial-host.md)
- Learn about [networking from the container's point of view](../config/containers/container-networking.md)
- Learn about [bridge networks](bridge.md)
- Learn about [overlay networks](overlay.md)
- Learn about [Macvlan networks](macvlan.md)
