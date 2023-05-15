---
title: Network drivers overview
description: Overview of Docker network drivers and related concepts
keywords: networking, drivers, bridge, routing, routing mesh, overlay, ports
---

Docker's networking subsystem is pluggable, using drivers. Several drivers
exist by default, and provide core networking functionality:

- `bridge`: The default network driver. If you don't specify a driver, this is
  the type of network you are creating. Bridge networks are commonly used when
  your application runs in a container that needs to communicate with other
  containers on the same host. See [bridge networks](bridge.md).

- `host`: For standalone containers, remove network isolation between the
  container and the Docker host, and use the host's networking directly. See
  [use the host network](host.md).

- `overlay`: Overlay networks connect multiple Docker daemons together and
  enable Swarm services and containers to communicate across nodes. This
  strategy removes the need to do OS-level routing.
  See [overlay networks](overlay.md).

- `ipvlan`: IPvlan networks give users total control over both IPv4 and IPv6
  addressing. The VLAN driver builds on top of that in giving operators complete
  control of layer 2 VLAN tagging and even IPvlan L3 routing for users
  interested in underlay network integration. See [IPvlan networks](ipvlan.md).

- `macvlan`: Macvlan networks allow you to assign a MAC address to a container,
  making it appear as a physical device on your network. The Docker daemon
  routes traffic to containers by their MAC addresses. Using the `macvlan`
  driver is sometimes the best choice when dealing with legacy applications that
  expect to be directly connected to the physical network, rather than routed
  through the Docker host's network stack. See
  [Macvlan networks](macvlan.md).

- `none`: For this container, disable all networking. `none` is not available
  for Swarm services. See [disable container networking](none.md).

- [Network plugins](/engine/extend/plugins_services/): You can install and use
  third-party network plugins with Docker.

### Network driver summary

- The default bridge network is good for running containers that don't require
  special networking capabilities.
- User-defined bridge networks enable containers on the same Docker host to
  communicate with each other. A user-defined network typically defines an
  isolated network for multiple containers belonging to a common project or
  component.
- Host network shares the host's network with the container. When you use this
  driver, the container's network isn't isolated from the host.
- Overlay networks are best when you need containers running on different
  Docker hosts to communicate, or when multiple applications work together
  using Swarm services.
- Macvlan networks are best when you are migrating from a VM setup or need your
  containers to look like physical hosts on your network, each with a unique
  MAC address.
- IPvlan is similar to Macvlan, but doesn't assign unique MAC addresses to
  containers. Consider using IPvlan when there's a restriction on the number of
  MAC addresses that can be assigned to a network interface or port.
- Third-party network plugins allow you to integrate Docker with specialized
  network stacks.

## Networking tutorials

Now that you understand the basics about Docker networks, deepen your
understanding using the following tutorials:

- [Standalone networking tutorial](../network-tutorial-standalone.md)
- [Host networking tutorial](../network-tutorial-host.md)
- [Overlay networking tutorial](../network-tutorial-overlay.md)
- [Macvlan networking tutorial](../network-tutorial-macvlan.md)
