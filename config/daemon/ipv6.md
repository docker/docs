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
  ...
}
```

You can optionally also configure the `fixed-cidr-v6` key, if you want to
assign an IPv6 subnet to the default bridge network:

```diff
{
  "experimental": true,
  "ipv6": true,
  "ip6tables": true,
+ "fixed-cidr-v6": "2001:db8:1::/64",
  ...
}
```

After saving the configuration file, reload the Docker daemon for your
changes to take effect:

```console
$ systemctl reload docker
```

You can now create networks with the `--ipv6` flag and assign containers IPv6
addresses using the `--ip6` flag.

## Dynamic IPv6 subnet allocation

If you want dynamic IPv6 subnet allocation, you must explicitly configure the
`default-address-pools` parameter to include:

- The default address pool values
- One or more custom IPv6 supernets

> **Note**
>
> Be aware that the following known limitations exist:
>
> - Supernets can't have a size larger than 80. This is due to an integer
>   overflow in the Docker daemon. See
>   [moby/moby#42801](https://github.com/moby/moby/issues/42801)
> - The difference between the supernet length and the pool size can't be
>   larger than 24. Otherwise, the daemon consumes all available memory. See
>   [moby/moby#40275](https://github.com/moby/moby/issues/40275)

The default address pool configuration is:

```json
{
  "default-address-pools": [
    { "base": "172.17.0.0/16", "size": 16 },
    { "base": "172.18.0.0/16", "size": 16 },
    { "base": "172.19.0.0/16", "size": 16 },
    { "base": "172.20.0.0/16", "size": 16 },
    { "base": "172.24.0.0/14", "size": 16 },
    { "base": "172.28.0.0/14", "size": 16 },
    { "base": "172.28.0.0/16", "size": 20 }
  ]
}
```

The following example shows a valid configuration with the default values and
an IPv6 supernet, with a prefix length of 64 and a size of 80.

```json
{
  "default-address-pools": [
    { "base": "172.17.0.0/16", "size": 16 },
    { "base": "172.18.0.0/16", "size": 16 },
    { "base": "172.19.0.0/16", "size": 16 },
    { "base": "172.20.0.0/16", "size": 16 },
    { "base": "172.24.0.0/14", "size": 16 },
    { "base": "172.28.0.0/14", "size": 16 },
    { "base": "172.28.0.0/16", "size": 20 },
    { "base": "2001:db8::/64", "size": 80 }
  ]
}
```


## Next steps

- [Networking overview](../../network/index.md)
