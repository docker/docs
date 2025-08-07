---
title: Networking overview
linkTitle: Networking
weight: 30
description: Learn how networking works from the container's point of view
keywords: networking, container, standalone, IP address, DNS resolution
aliases:
- /articles/networking/
- /config/containers/container-networking/
- /engine/tutorials/networkingcontainers/
- /engine/userguide/networking/
- /engine/userguide/networking/configure-dns/
- /engine/userguide/networking/default_network/binding/
- /engine/userguide/networking/default_network/configure-dns/
- /engine/userguide/networking/default_network/container-communication/
- /engine/userguide/networking/dockernetworks/
- /network/
---

Container networking refers to the ability for containers to connect to and
communicate with each other, or to non-Docker workloads.

Containers have networking enabled by default, and they can make outgoing
connections. A container has no information about what kind of network it's
attached to, or whether their peers are also Docker workloads or not. A
container only sees a network interface with an IP address, a gateway, a
routing table, DNS services, and other networking details. That is, unless the
container uses the `none` network driver.

This page describes networking from the point of view of the container,
and the concepts around container networking.
This page doesn't describe OS-specific details about how Docker networks work.
For information about how Docker manipulates `iptables` rules on Linux,
see [Packet filtering and firewalls](packet-filtering-firewalls.md).

## User-defined networks

You can create custom, user-defined networks, and connect multiple containers
to the same network. Once connected to a user-defined network, containers can
communicate with each other using container IP addresses or container names.

The following example creates a network using the `bridge` network driver and
running a container in the created network:

```console
$ docker network create -d bridge my-net
$ docker run --network=my-net -itd --name=container3 busybox
```

### Drivers

The following network drivers are available by default, and provide core
networking functionality:

| Driver    | Description                                                              |
| :-------- | :----------------------------------------------------------------------- |
| `bridge`  | The default network driver.                                              |
| `host`    | Remove network isolation between the container and the Docker host.      |
| `none`    | Completely isolate a container from the host and other containers.       |
| `overlay` | Overlay networks connect multiple Docker daemons together.               |
| `ipvlan`  | IPvlan networks provide full control over both IPv4 and IPv6 addressing. |
| `macvlan` | Assign a MAC address to a container.                                     |

For more information about the different drivers, see [Network drivers
overview](./drivers/_index.md).

### Connecting to multiple networks

A container can be connected to multiple networks.

For example, a frontend container may be connected to a bridge network
with external access, and a
[`--internal`](/reference/cli/docker/network/create/#internal) network
to communicate with containers running backend services that do not need
external network access.

A container may also be connected to different types of network. For example,
an `ipvlan` network to provide internet access, and a `bridge` network for
access to local services.

When sending packets, if the destination is an address in a directly connected
network, packets are sent to that network. Otherwise, packets are sent to
a default gateway for routing to their destination. In the example above,
the `ipvlan` network's gateway must be the default gateway.

The default gateway is selected by Docker, and may change whenever a
container's network connections change.
To make Docker choose a specific default gateway when creating the container
or connecting a new network, set a gateway priority. See option `gw-priority`
for the [`docker run`](/reference/cli/docker/container/run.md) and
[`docker network connect`](/reference/cli/docker/network/connect.md) commands.

The default `gw-priority` is `0` and the gateway in the network with the
highest priority is the default gateway. So, when a network should always
be the default gateway, it is enough to set its `gw-priority` to `1`.

```console
$ docker run --network name=gwnet,gw-priority=1 --network anet1 --name myctr myimage
$ docker network connect anet2 myctr
```

## Container networks

In addition to user-defined networks, you can attach a container to another
container's networking stack directly, using the `--network
container:<name|id>` flag format.

The following flags aren't supported for containers using the `container:`
networking mode:

- `--add-host`
- `--hostname`
- `--dns`
- `--dns-search`
- `--dns-option`
- `--mac-address`
- `--publish`
- `--publish-all`
- `--expose`

The following example runs a Redis container, with Redis binding to
`localhost`, then running the `redis-cli` command and connecting to the Redis
server over the `localhost` interface.

```console
$ docker run -d --name redis example/redis --bind 127.0.0.1
$ docker run --rm -it --network container:redis example/redis-cli -h 127.0.0.1
```

## Published ports

When you create or run a container using `docker create` or `docker run`, all
ports of containers on bridge networks are accessible from the Docker host and
other containers connected to the same network. Ports are not accessible from
outside the host or, with the default configuration, from containers in other
networks.

Use the `--publish` or `-p` flag to make a port available outside the host,
and to containers in other bridge networks.

For more information about port mapping, including how to disable it and use
direct routing to containers, see
[port publishing](./port-publishing.md).

## IP address and hostname

When creating a network, IPv4 address allocation is enabled by default, it
can be disabled using `--ipv4=false`. IPv6 address allocation can be enabled
using `--ipv6`.

```console
$ docker network create --ipv6 --ipv4=false v6net
```

By default, the container gets an IP address for every Docker network it attaches to.
A container receives an IP address out of the IP subnet of the network.
The Docker daemon performs dynamic subnetting and IP address allocation for containers.
Each network also has a default subnet mask and gateway.

You can connect a running container to multiple networks,
either by passing the `--network` flag multiple times when creating the container,
or using the `docker network connect` command for already running containers.
In both cases, you can use the `--ip` or `--ip6` flags to specify the container's IP address on that particular network.

In the same way, a container's hostname defaults to be the container's ID in Docker.
You can override the hostname using `--hostname`.
When connecting to an existing network using `docker network connect`,
you can use the `--alias` flag to specify an additional network alias for the container on that network.

## DNS services

Containers use the same DNS servers as the host by default, but you can
override this with `--dns`.

By default, containers inherit the DNS settings as defined in the
`/etc/resolv.conf` configuration file.
Containers that attach to the default `bridge` network receive a copy of this file.
Containers that attach to a
[custom network](tutorials/standalone.md#use-user-defined-bridge-networks)
use Docker's embedded DNS server.
The embedded DNS server forwards external DNS lookups to the DNS servers configured on the host.

You can configure DNS resolution on a per-container basis, using flags for the
`docker run` or `docker create` command used to start the container.
The following table describes the available `docker run` flags related to DNS
configuration.

| Flag           | Description                                                                                                                                                                                                                                           |
| -------------- |-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `--dns`        | The IP address of a DNS server. To specify multiple DNS servers, use multiple `--dns` flags. DNS requests will be forwarded from the container's network namespace so, for example, `--dns=127.0.0.1` refers to the container's own loopback address. |
| `--dns-search` | A DNS search domain to search non-fully qualified hostnames. To specify multiple DNS search prefixes, use multiple `--dns-search` flags.                                                                                                              |
| `--dns-opt`    | A key-value pair representing a DNS option and its value. See your operating system's documentation for `resolv.conf` for valid options.                                                                                                              |
| `--hostname`   | The hostname a container uses for itself. Defaults to the container's ID if not specified.                                                                                                                                                            |

### Custom hosts

Your container will have lines in `/etc/hosts` which define the hostname of the
container itself, as well as `localhost` and a few other common things. Custom
hosts, defined in `/etc/hosts` on the host machine, aren't inherited by
containers. To pass additional hosts into a container, refer to [add entries to
container hosts file](/reference/cli/docker/container/run.md#add-host) in the
`docker run` reference documentation.

## Proxy server

If your container needs to use a proxy server, see
[Use a proxy server](/manuals/engine/daemon/proxy.md).
