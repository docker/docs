---
title: Enable IPv6 support
description: How to enable IPv6 support in the Docker daemon
keywords: daemon, network, networking, ipv6
aliases:
- /engine/userguide/networking/default_network/ipv6/
---

Before you can use IPv6 in Docker containers, you need to
enable IPv6 support in the Docker daemon. Afterward, you can choose to use
either IPv4 or IPv6 (or both) with any container or network.

IPv6 is only supported on Docker daemons running on Linux hosts.

> **Note**
>
> When using IPv6, we recommend that you enable the
> [experimental](../../release-lifecycle.md#experimental)
> `ip6tables` parameter in the daemon configuration.

## Create an IPv6 network

The following steps show you how to create a Docker network that uses IPv6.

1. Edit the Docker daemon configuration file,
   located at `/etc/docker/daemon.json`. Configure the following parameters:

   ```json
   {
     "experimental": true,
     "ip6tables": true
   }
   ```

   `ip6tables` enables additional IPv6 packet filter rules, providing network
   isolation and port mapping. This parameter requires `experimental` to be
   set to `true`.

2. Save the configuration file.
3. Restart the Docker daemon for your changes to take effect.

   ```console
   $ sudo systemctl restart docker
   ```

4. Create a new IPv6 network.

   - Using `docker network create`:

     ```console
     $ docker network create --ipv6 --subnet 2001:0DB8::/112 ip6net
     ```

   - Using a Docker Compose file:

     ```yaml
      networks:
        ip6net:
          enable_ipv6: true
          ipam:
            config:
              - subnet: 2001:0DB8::/112
     ```

You can now run containers that attach to the `ip6net` network.

```console
$ docker run --rm --network ip6net -p 80:80 traefik/whoami
```

This publishes port 80 on both IPv6 and IPv4.
You can verify the IPv6 connection by running curl,
connecting to port 80 on the IPv6 loopback address:

```console
$ curl http://[::1]:80
Hostname: ea1cfde18196
IP: 127.0.0.1
IP: ::1
IP: 172.17.0.2
IP: fe80::42:acff:fe11:2
RemoteAddr: [fe80::42:acff:fe11:2]:54890
GET / HTTP/1.1
Host: [::1]
User-Agent: curl/8.1.2
Accept: */*
```

## Use IPv6 for the default bridge network

The following steps show you how to use IPv6 on the default bridge network.

1. Edit the Docker daemon configuration file,
   located at `/etc/docker/daemon.json`. Configure the following parameters:

   ```json
   {
     "ipv6": true,
     "fixed-cidr-v6": "2001:db8:1::/64",
     "experimental": true,
     "ip6tables": true
   }
   ```

   - `ipv6` enables IPv6 networking on the default network.
   - `fixed-cidr-v6` assigns a subnet to the default bridge network,
     enabling dynamic IPv6 address allocation.
   - `ip6tables` enables additional IPv6 packet filter rules, providing network
     isolation and port mapping. This parameter requires `experimental` to be
     set to `true`.

2. Save the configuration file.
3. Restart the Docker daemon for your changes to take effect.

   ```console
   $ sudo systemctl restart docker
   ```

You can now run containers on the default bridge network.

```console
$ docker run --rm -p 80:80 traefik/whoami
```

This publishes port 80 on both IPv6 and IPv4.
You can verify the IPv6 connection by making a request
to port 80 on the IPv6 loopback address:

```console
$ curl http://[::1]:80
Hostname: ea1cfde18196
IP: 127.0.0.1
IP: ::1
IP: 172.17.0.2
IP: fe80::42:acff:fe11:2
RemoteAddr: [fe80::42:acff:fe11:2]:54890
GET / HTTP/1.1
Host: [::1]
User-Agent: curl/8.1.2
Accept: */*
```

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
    { "base": "172.20.0.0/14", "size": 16 },
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
