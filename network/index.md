---
title: Overview
description: Overview of Docker networks and networking concepts
keywords: networking, bridge, routing, routing mesh, overlay, ports
redirect_from:
- /engine/userguide/networking/
- /engine/userguide/networking/dockernetworks/
- /articles/networking/
---

One of the reasons Docker containers and services are so powerful is that
you can connect them together, or connect them to non-Docker workloads. Docker
containers and services do not even need to be aware that they are deployed on
Docker, or whether their peers are also Docker workloads or not. Whether your
Docker hosts run Linux, Windows, or a mix of the two, you can use Docker to
manage them in a platform-agnostic way.

This topic defines some basic Docker networking concepts and prepares you to
design and deploy your applications to take full advantage of these
capabilities.

Most of this content applies to all Docker installations. However,
[a few advanced features](#docker-ee-networking-features) are only available to
Docker EE customers.

## Scope of this topic

This topic does **not** go into OS-specific details about how Docker networks
work, so you will not find information about how Docker manipulates `iptables`
rules on Linux or how it manipulates routing rules on Windows servers, and you
will not find detailed information about how Docker forms and encapsulates
packets or handles encryption. See [Docker and iptables](/network/iptables.md)
and
[Docker Reference Architecture: Designing Scalable, Portable Docker Container Networks](https://success.docker.com/Architecture/Docker_Reference_Architecture%3A_Designing_Scalable%2C_Portable_Docker_Container_Networks)
for a much greater depth of technical detail.

In addition, this topic does not provide any tutorials for how to create,
manage, and use Docker networks. Each section includes links to relevant
tutorials and command references.

## Network drivers

Docker's networking subsystem is pluggable, using drivers. Several drivers
exist by default, and provide core networking functionality:

- `bridge`: The default network driver. If you don't specify a driver, this is
  the type of network you are creating. **Bridge networks are usually used when
  your applications run in standalone containers that need to communicate.** See
  [bridge networks](bridge.md).

- `host`: For standalone containers, remove network isolation between the
  container and the Docker host, and use the host's networking directly. `host`
  is only available for swarm services on Docker 17.06 and higher. See
  [use the host network](host.md).

- `overlay`: Overlay networks connect multiple Docker daemons together and
  enable swarm services to communicate with each other. You can also use overlay
  networks to facilitate communication between a swarm service and a standalone
  container, or between two standalone containers on different Docker daemons.
  This strategy removes the need to do OS-level routing between these
  containers. See [overlay networks](overlay.md).

- `macvlan`: Macvlan networks allow you to assign a MAC address to a container,
  making it appear as a physical device on your network. The Docker daemon
  routes traffic to containers by their MAC addresses. Using the `macvlan`
  driver is sometimes the best choice when dealing with legacy applications that
  expect to be directly connected to the physical network, rather than routed
  through the Docker host's network stack. See
  [Macvlan networks](macvlan.md).

- `none`: For this container, disable all networking. Usually used in
  conjunction with a custom network driver. `none` is not available for swarm
  services. See
  [disable container networking](none.md).

- [Network plugins](/engine/extend/plugins_services/): You can install and use
  third-party network plugins with Docker. These plugins are available from
  [Docker Store](https://store.docker.com/search?category=network&q=&type=plugin)
  or from third-party vendors. See the vendor's documentation for installing and
  using a given network plugin.


### Network driver summary

- **User-defined bridge networks** are best when you need multiple containers to
  communicate on the same Docker host.
- **Host networks** are best when the network stack should not be isolated from
  the Docker host, but you want other aspects of the container to be isolated.
- **Overlay networks** are best when you need containers running on different
  Docker hosts to communicate, or when multiple applications work together using
  swarm services.
- **Macvlan networks** are best when you are migrating from a VM setup or
  need your containers to look like physical hosts on your network, each with a
  unique MAC address.
- **Third-party network plugins** allow you to integrate Docker with specialized
  network stacks.

## Docker EE networking features

The following two features are only possible when using Docker EE and managing
your Docker services using Universal Control Plane (UCP):

- The [HTTP routing mesh](/datacenter/ucp/2.2/guides/admin/configure/use-domain-names-to-access-services/)
  allows you to share the same network IP address and port among multiple
  services. UCP routes the traffic to the appropriate service using the
  combination of hostname and port, as requested from the client.

- [Session stickiness](/datacenter/ucp/2.2/guides/user/services/use-domain-names-to-access-services/#sticky-sessions) allows you to specify information in the HTTP header
  which UCP uses to route subsequent requests to the same service task, for
  applications which require stateful sessions.

## Networking tutorials

Now that you understand the basics about Docker networks, deepen your
understanding using the following tutorials:

- [Standalone networking tutorial](network-tutorial-standalone.md)
- [Host networking tutorial](network-tutorial-host.md)
- [Overlay networking tutorial](network-tutorial-overlay.md)
- [Macvlan networking tutorial](network-tutorial-macvlan.md)

