---
title: Use host networking
description: All about exposing containers on the Docker host's network
keywords: network, host, standalone
---

If you use the `host` network driver for a container, that container's network
stack is not isolated from the Docker host. For instance, if you run a container
which binds to port 80 and you use `host` networking, the container's
application will be available on port 80 on the host's IP address.

In Docker 17.06 and higher, you can also use a `host` network for a swarm
service, by passing `--network host` to the `docker container create` command.
In this case, control traffic (traffic related to managing the swarm and the
service) is still sent across an overlay network, but the individual swarm
service containers send data using the Docker daemon's host network and ports.
This creates some extra limitations. For instance, if a service container binds
to port 80, only one service container can run on a given swarm node.

If your container or service publishes no ports, host networking has no effect.

## Next steps

-  Go through the [host networking tutorial](/network/network-tutorial-host.md)
- Learn about [networking from the container's point of view](/config/containers/container-networking.md)
- Learn about [bridge networks](/network/bridge.md)
- Learn about [overlay networks](/network/overlay.md)
- Learn about [Macvlan networks](/network/macvlan.md)
