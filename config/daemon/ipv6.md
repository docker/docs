---
title: Enable IPv6 support
description: How to enable IPv6 support in the Docker daemon
keywords: daemon, network, networking, ipv6
redirect_from:
- /engine/userguide/networking/default_network/ipv6/
---

Before you can use IPv6 in Docker containers or swarm services, you need to
enable IPv6 support in the Docker daemon. Afterward, you can choose to use
either IPv4 or IPv6 (or both) with any container, service, or network.

> **Note**
>
> - IPv6 support is experimental, use it with caution.
> - IPv6 is only supported on Docker daemons running on Linux hosts.

## Daemon configuration

To enable IPv6, you must edit the Docker daemon configuration file located at
`/etc/docker/daemon.json`. Configure the following parameters:

```json
{
  "experimental": true,
  "ipv6": true,
  "ip6tables": true,
  "fixed-cidr-v6": "2001:db8:1::/64",
  ...
}
```

This configuration makes IPv6 networking function as you would expect it to.
The `ipv6` and `fixed-cidr-v6` parameters are optional.
They assign an IPv6 subnet to the default bridge network.

After saving the configuration file, restart the Docker daemon for your
changes to take effect:

```console
$ systemctl restart docker
```

Upon restart, the daemon assigns IPv6 addresses to containers connected to the
default bridge network, and to user-defined networks configured with an IPv6 subnet.

## Dynamic IPv6 subnet allocation

If you don't explicitly configure subnets for user-defined networks,
using `docker network create --subnet=<your-subnet>`,
those networks use the default address pools of the daemon as a fallback.
The default address pools are all IPv4 pools.
This also applies to networks created from a Docker Compose file,
with `enable_ipv6` set to `true`.

To enable dynamic subnet allocation for user-defined IPv6 networks,
you must manually configure address pools of the daemon to include:

- The default IPv4 address pools
- One or more IPv6 pools of your own

The default address pool configuration is:

```json
{
  "default-address-pools": [
    { "base": "172.17.0.0/16", "size": 16 },
    { "base": "172.18.0.0/16", "size": 16 },
    { "base": "172.19.0.0/16", "size": 16 },
    { "base": "172.20.0.0/14", "size": 16 },
    { "base": "172.24.0.0/14", "size": 16 },
    { "base": "172.28.0.0/14", "size": 16 },
    { "base": "192.168.0.0/16", "size": 20 }
  ]
}
```

The following example shows a valid configuration with the default values and
an IPv6 pool. The IPv6 pool in the example provides up to 256 IPv6 subnets of
size `/112`, from an IPv6 pool of prefix length `/104`. Each `/112`-sized
subnet supports 65 536 IPv6 addresses.

> **Note**
>
> Be aware that the following known limitations exist for IPv6 pools:
>
> - The `base` value for IPv6 needs a minimum prefix length of `/64`.
>   This is due to an integer overflow in the Docker daemon.
>   See [moby/moby#42801](https://github.com/moby/moby/issues/42801).
> - The difference between the pool length and the pool size can't be larger
>   than 24. Defining an excessive number of subnets causes the daemon to
>   consume all available memory.
>   See [moby/moby#40275](https://github.com/moby/moby/issues/40275).

```json
{
  "default-address-pools": [
    { "base": "172.17.0.0/16", "size": 16 },
    { "base": "172.18.0.0/16", "size": 16 },
    { "base": "172.19.0.0/16", "size": 16 },
    { "base": "172.20.0.0/16", "size": 16 },
    { "base": "172.24.0.0/14", "size": 16 },
    { "base": "172.28.0.0/14", "size": 16 },
    { "base": "192.168.0.0/16", "size": 20 },
    { "base": "2001:db8::/104", "size": 112 }
  ]
}
```

> **Note**
>
> The address `2001:db8` in this example is
> [reserved for use in documentation][wikipedia-ipv6-reserved].
> Replace it with a valid IPv6 network.
> The default IPv4 pools are from the private address range,
> the IPv6 equivalent would be [ULA networks][wikipedia-ipv6-ula].

[wikipedia-ipv6-reserved]: https://en.wikipedia.org/wiki/Reserved_IP_addresses#IPv6
[wikipedia-ipv6-ula]: https://en.wikipedia.org/wiki/Unique_local_address

## Next steps

- [Networking overview](../../network/index.md)
