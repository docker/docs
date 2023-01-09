---
title: Container networking
description: How networking works from the container's point of view
keywords: networking, container, standalone
redirect_from:
  - /engine/userguide/networking/configure-dns/
  - /engine/userguide/networking/default_network/configure-dns/
  - /engine/userguide/networking/default_network/binding/
  - /engine/userguide/networking/default_network/container-communication/
---

A container has no information about what kind of network it's attached to,
whether it's a [bridge](../../network/bridge.md), an [overlay](../../network/overlay.md),
a [macvlan network](../../network/macvlan.md), or a custom network plugin.
A container only sees a network interface with an IP address,
a gateway, a routing table, DNS services, and other networking details.
That is, unless the container uses the `none` network driver.
This page describes networking from the point of view of the container.

## Published ports

By default, when you create or run a container using `docker create` or `docker run`,
the container doesn't expose any of it's ports to the outside world.
To make a port available to services outside of Docker,
or to Docker containers running on a different network,
use the `--publish` or `-p` flag.
This creates a firewall rule in the container,
mapping a container port to a port on the Docker host to the outside world.
Here are some examples:

| Flag value                      | Description                                                                                                                                           |
| ------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- |
| `-p 8080:80`                    | Map TCP port 80 in the container to port `8080` on the Docker host.                                                                                   |
| `-p 192.168.1.100:8080:80`      | Map TCP port 80 in the container to port `8080` on the Docker host for connections to host IP `192.168.1.100`.                                        |
| `-p 8080:80/udp`                | Map UDP port 80 in the container to port `8080` on the Docker host.                                                                                   |
| `-p 8080:80/tcp -p 8080:80/udp` | Map TCP port 80 in the container to TCP port `8080` on the Docker host, and map UDP port `80` in the container to UDP port `8080` on the Docker host. |

## IP address and hostname

By default, the container gets an IP address for every Docker network it attaches to.
A container receives an IP address out of the IP pool of the network it attaches to.
The Docker daemon effectively acts as a DHCP server for each container.
Each network also has a default subnet mask and gateway.

When a container starts, it can only attach to a single network, using the `--network` flag.
You can connect a running container to multiple networks using the `docker network connect` command.
When you start a container using the `--network` flag,
you can specify the IP address for the container on that network using the `--ip` or `--ip6` flags.

When you connect an existing container to a different network using `docker network connect`,
you can use the `--ip` or `--ip6` flags on that command
to specify the container's IP address on the additional network.

In the same way, a container's hostname defaults to be the container's ID in Docker.
You can override the hostname using `--hostname`.
When connecting to an existing network using `docker network connect`,
you can use the `--alias` flag to specify an additional network alias for the container on that network.

## DNS services

By default, containers inherit the DNS settings of the host, as defined in the `/etc/resolv.conf` configuration file.
Containers that attach to the default `bridge` network receive a copy of this file.
Containers that attach to a
[custom network](../../network/network-tutorial-standalone.md#use-user-defined-bridge-networks)
use Docker's embedded DNS server.
The embedded DNS server forwards external DNS lookups to the DNS servers configured on the host.

Custom hosts, defined in `/etc/hosts` on the host machine, aren't inherited by containers.
To pass additional hosts into container, refer to
[add entries to container hosts file](../../engine/reference/commandline/run.md#add-host)
in the `docker run` reference documentation.
You can override these settings on a per-container basis.

| Flag           | Description                                                                                                                                                                                                                                                         |
| -------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `--dns`        | The IP address of a DNS server. To specify multiple DNS servers, use multiple `--dns` flags. If the container can't reach any of the IP addresses you specify, it uses Google's public DNS server at `8.8.8.8`. This allows containers to resolve internet domains. |
| `--dns-search` | A DNS search domain to search non-fully-qualified hostnames. To specify multiple DNS search prefixes, use multiple `--dns-search` flags.                                                                                                                            |
| `--dns-opt`    | A key-value pair representing a DNS option and its value. See your operating system's documentation for `resolv.conf` for valid options.                                                                                                                            |
| `--hostname`   | The hostname a container uses for itself. Defaults to the container's ID if not specified.                                                                                                                                                                          |

## Proxy server

If your container needs to use a proxy server, see
[Use a proxy server](../../network/proxy.md).
