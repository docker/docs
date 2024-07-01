---
title: Packet filtering and firewalls
description: How Docker works with packet filtering, iptables, and firewalls
keywords: network, iptables, firewall
aliases:
- /network/iptables/
---

On Linux, Docker creates `iptables` and `ip6tables` rules to implement network
isolation, port publishing and filtering.

Because these rules are required for the correct functioning of Docker bridge
networks, you should not modify the rules created by Docker.

But, if you are running Docker on a host exposed to the internet, you will
probably want to add iptables policies that prevent unauthorized access to
containers or other services running on your host. This page describes how
to achieve that, and the caveats you need to be aware of.

> **Note**
> 
> Docker creates `iptables` rules for bridge networks.
> 
> No `iptables` rules are created for `ipvlan`, `macvlan` or `host` networking.

## Docker and iptables chains

In the `filter` table, Docker sets the default policy to `DROP`, and creates the
following custom `iptables` chains:

* `DOCKER-USER`
  * A placeholder for user-defined rules that will be processed before rules
    in the `DOCKER` chain.
* `DOCKER`
  * Rules that determine whether a packet that is not part of an established
    connection should be accepted, based on the port forwarding configuration
    of running containers.
* `DOCKER-ISOLATION-STAGE-1` and `DOCKER-ISOLATION-STAGE-2`
  * Rules to isolate Docker networks from each other.

In the `FORWARD` chain, Docker adds rules that pass packets that are not related
to established connections to these custom chains, as well as rules to accept
packets that are part of established connections.

In the `nat` table, Docker creates chain `DOCKER` and adds rules to implement
masquerading and port-mapping.

### Add iptables policies before Docker's rules

Packets that get accepted or rejected by rules in these custom chains will not
be seen by user-defined rules appended to the `FORWARD` chain. So, to add
additional rules to filter these packets, use the `DOCKER-USER` chain.

### Match the original IP and ports for requests

When packets arrive to the `DOCKER-USER` chain, they have already passed through
a Destination Network Address Translation (DNAT) filter. That means that the
`iptables` flags you use can only match internal IP addresses and ports of
containers.

If you want to match traffic based on the original IP and port in the network
request, you must use the
[`conntrack` iptables extension](https://ipset.netfilter.org/iptables-extensions.man.html#lbAO).
For example:

```console
$ sudo iptables -I DOCKER-USER -p tcp -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT
$ sudo iptables -I DOCKER-USER -p tcp -m conntrack --ctorigdst 198.51.100.2 --ctorigdstport 80 -j ACCEPT
```

> **Important**
>
> Using the `conntrack` extension may result in degraded performance.
{ .important }

## Port publishing and mapping

By default, for both IPv4 and IPv6, the daemon blocks access to ports that have not
been published. Published container ports are mapped to host IP addresses.
To do this, it uses iptables to perform Network Address Translation (NAT),
Port Address Translation (PAT), and masquerading.

For example, `docker run -p 8080:80 [...]` creates a mapping
between port 8080 on any address on the Docker host, and the container's
port 80. Outgoing connections from the container will masquerade, using
the Docker host's IP address.

### Restrict external connections to containers

By default, all external source IPs are allowed to connect to ports that have
been published to the Docker host's addresses.

To allow only a specific IP or network to access the containers, insert a
negated rule at the top of the `DOCKER-USER` filter chain. For example, the
following rule drops packets from all IP addresses except `192.0.2.2`:

```console
$ iptables -I DOCKER-USER -i ext_if ! -s 192.0.2.2 -j DROP
```

You will need to change `ext_if` to correspond with your
host's actual external interface. You could instead allow connections from a
source subnet. The following rule only allows access from the subnet `192.0.2.0/24`:

```console
$ iptables -I DOCKER-USER -i ext_if ! -s 192.0.2.0/24 -j DROP
```

Finally, you can specify a range of IP addresses to accept using `--src-range`
(Remember to also add `-m iprange` when using `--src-range` or `--dst-range`):

```console
$ iptables -I DOCKER-USER -m iprange -i ext_if ! --src-range 192.0.2.1-192.0.2.3 -j DROP
```

You can combine `-s` or `--src-range` with `-d` or `--dst-range` to control both
the source and destination. For instance, if the Docker host has addresses
`2001:db8:1111::2` and `2001:db8:2222::2`, you can make rules specific to
`2001:db8:1111::2` and leave `2001:db8:2222::2` open.

`iptables` is complicated. There is a lot more information at [Netfilter.org HOWTO](https://www.netfilter.org/documentation/HOWTO/NAT-HOWTO.html).

### Direct routing

Port mapping ensures that published ports are accessible on the host's
network addresses, which are likely to be routable for any external
clients. No routes are normally set up in the host's network for container
addresses that exist within a host.

But, particularly with IPv6 you may prefer to avoid using NAT and instead
arrange for external routing to container addresses.

To access containers on a bridge network from outside the Docker host,
you must set up routing to the bridge network via an address on the Docker
host. This can be achieved using static routes, Border Gateway Protocol
(BGP), or any other means appropriate for your network.

The bridge network driver has options
`com.docker.network.bridge.gateway_mode_ipv6=<nat|routed>` and
`com.docker.network.bridge.gateway_mode_ipv4=<nat|routed>`.

The default is `nat`, NAT and masquerading rules are set up for each
published container port. With mode `routed`, no NAT or masquerading rules
are set up, but `iptables` are still set up so that only published container
ports are accessible.

In `routed` mode, a host port in a `-p` or `--publish` port mapping is
not used, and the host address is only used to decide whether to apply
the mapping to IPv4 or IPv6. So, when a mapping only applies to `routed`
mode, only addresses `0.0.0.0` or `::1` are allowed, and a host port
must not be given.

Mapped container ports, in `nat` or `routed` mode, are accessible from
any remote address, if routing is set up in the network, unless the
Docker host's firewall has additional restrictions.

#### Example

Create a network suitable for direct routing for IPv6, with NAT enabled
for IPv4:
```console
$ docker network create --ipv6 --subnet 2001:db8::/64 -o com.docker.network.bridge.gateway_mode_ipv6=routed mynet
```

Create a container with a published port:
```console
$ docker run --network=mynet -p 8080:80 myimage
```

Then:
- Only container port 80 will be open, for IPv4 and IPv6. It is accessible
  from anywhere, if there is routing to the container's address, and access
  is not blocked by the host's firewall.
- For IPv6, using `routed` mode, port 80 will be open on the container's IP
  address. Port 8080 will not be opened on the host's IP addresses, and
  outgoing packets will use the container's IP address.
- For IPv4, using the default `nat` mode, the container's port 80 will be
  accessible via port 8080 on the host's IP addresses, as well as directly.
  Connections originating from the container will masquerade, using the 
  host's IP address.

In `docker inspect`, this port mapping will be shown as follows. Note that
there is no `HostPort` for IPv6, because it is using `routed` mode:
```console
$ docker container inspect <id> --format "{{json .NetworkSettings.Ports}}"
{"80/tcp":[{"HostIp":"0.0.0.0","HostPort":"8080"},{"HostIp":"::","HostPort":""}]}
```

Alternatively, to make the mapping IPv6-only, disabling IPv4 access to the
container's port 80, use the unspecified IPv6 address `[::]` and do not
include a host port number:
```console
$ docker run --network mynet -p '[::]::80'
```

### Setting the default bind address for containers

By default, when a container's ports are mapped without any specific host
address, the Docker daemon binds published container ports to all host
addresses (`0.0.0.0` and `[::]`).

For example, the following command publishes port 8080 to all network
interfaces on the host, on both IPv4 and IPv6 addresses, potentially
making them available to the outside world.

```console
docker run -p 8080:80 nginx
```

You can change the default binding address for published container ports so that
they're only accessible to the Docker host by default. To do that, you can
configure the daemon to use the loopback address (`127.0.0.1`) instead.

> **Warning**
>
> Hosts within the same L2 segment (for example, hosts connected to the same
> network switch) can reach ports published to localhost.
> For more information, see
> [moby/moby#45610](https://github.com/moby/moby/issues/45610)
{ .warning }

To configure this setting for user-defined bridge networks, use
the `com.docker.network.bridge.host_binding_ipv4`
[driver option](./drivers/bridge.md#options) when you create the network.

```console
$ docker network create mybridge \
  -o "com.docker.network.bridge.host_binding_ipv4=127.0.0.1"
```

> **Note**
>
> - Setting the default binding address to `::` means port bindings with no host
>   address specified will work for any IPv6 address on the host. But, `0.0.0.0`
>   means any IPv4 or IPv6 address.
> - Changing the default bind address doesn't have any effect on Swarm services.
>   Swarm services are always exposed on the `0.0.0.0` network interface.

#### Default bridge

To set the default binding for the default bridge network, configure the `"ip"`
key in the `daemon.json` configuration file:

```json
{
  "ip": "127.0.0.1"
}
```

This changes the default binding address to `127.0.0.1` for published container
ports on the default bridge network.
Restart the daemon for this change to take effect.
Alternatively, you can use the `dockerd --ip` flag when starting the daemon.

## Docker on a router

Docker sets the policy for the `FORWARD` chain to `DROP`. This will prevent
your Docker host from acting as a router.

If you want your system to function as a router, you must add explicit
`ACCEPT` rules to the `DOCKER-USER` chain. For example:

```console
$ iptables -I DOCKER-USER -i src_if -o dst_if -j ACCEPT
```

## Prevent Docker from manipulating iptables

It is possible to set the `iptables` or `ip6tables` keys to `false` in
[daemon configuration](https://docs.docker.com/reference/cli/dockerd/), but
this option is not appropriate for most users. It is likely to break
container networking for the Docker Engine.

All ports of all containers will be accessible from the network, and none
will be mapped from Docker host IP addresses.

It is not possible to completely prevent Docker from creating `iptables`
rules, and creating rules after-the-fact is extremely involved and beyond
the scope of these instructions.

## Integration with firewalld

If you are running Docker with the `iptables` option set to `true`, and
[firewalld](https://firewalld.org) is enabled on your system, Docker
automatically creates a `firewalld` zone called `docker`, with target `ACCEPT`.

All network interfaces created by Docker (for example, `docker0`) are inserted
into the `docker` zone.

Docker also creates a forwarding policy called `docker-forwarding` that allows
forwarding from `ANY` zone to the `docker` zone.

## Docker and ufw

[Uncomplicated Firewall](https://launchpad.net/ufw)
(ufw) is a frontend that ships with Debian and Ubuntu,
and it lets you manage firewall rules. Docker and ufw use iptables in ways
that make them incompatible with each other.

When you publish a container's ports using Docker, traffic to and from that
container gets diverted before it goes through the ufw firewall settings.
Docker routes container traffic in the `nat` table, which means that packets
are diverted before it reaches the `INPUT` and `OUTPUT` chains that ufw uses.
Packets are routed before the firewall rules can be applied,
effectively ignoring your firewall configuration.
