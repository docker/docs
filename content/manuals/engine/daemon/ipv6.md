---
title: Use IPv6 networking
weight: 20
description: How to enable IPv6 support in the Docker daemon
keywords: daemon, network, networking, ipv6
aliases:
- /engine/userguide/networking/default_network/ipv6/
- /config/daemon/ipv6/
---

IPv6 is only supported on Docker daemons running on Linux hosts.

## Create an IPv6 network

- Using `docker network create`:

  ```console
  $ docker network create --ipv6 ip6net
  ```

- Using `docker network create`, specifying an IPv6 subnet:

  ```console
  $ docker network create --ipv6 --subnet 2001:db8::/64 ip6net
  ```

- Using a Docker Compose file:

  ```yaml
   networks:
     ip6net:
       enable_ipv6: true
       ipam:
         config:
           - subnet: 2001:db8::/64
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
IP: 2001:db8::2
IP: fe80::42:acff:fe11:2
RemoteAddr: [2001:db8::1]:37574
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
     "fixed-cidr-v6": "2001:db8:1::/64"
   }
   ```

   - `ipv6` enables IPv6 networking on the default network.
   - `fixed-cidr-v6` assigns a subnet to the default bridge network,
     enabling dynamic IPv6 address allocation.
   - `ip6tables` enables additional IPv6 packet filter rules, providing network
     isolation and port mapping. It is enabled by-default, but can be disabled.

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
IP: 2001:db8:1::242:ac12:2
IP: fe80::42:acff:fe12:2
RemoteAddr: [2001:db8:1::1]:35558
GET / HTTP/1.1
Host: [::1]
User-Agent: curl/8.1.2
Accept: */*
```

## Dynamic IPv6 subnet allocation

If you don't explicitly configure subnets for user-defined networks,
using `docker network create --subnet=<your-subnet>`,
those networks use the default address pools of the daemon as a fallback.
This also applies to networks created from a Docker Compose file,
with `enable_ipv6` set to `true`.

If no IPv6 pools are included in Docker Engine's `default-address-pools`,
and no `--subnet` option is given, [Unique Local Addresses (ULAs)][wikipedia-ipv6-ula]
will be used when IPv6 is enabled. These `/64` subnets include a 40-bit
Global ID based on the Docker Engine's randomly generated ID, to give a
high probability of uniqueness.

The built-in default address pool configuration is shown in [Subnet allocation](../network/_index.md#subnet-allocation).
It does not include any IPv6 pools.

To use different pools of IPv6 subnets for dynamic address allocation,
you must manually configure address pools of the daemon to include:

- The default IPv4 address pools
- One or more IPv6 pools of your own

The following example shows a valid configuration with IPv4 and IPv6 pools,
both pools provide 256 subnets. IPv4 subnets with prefix length `/24` will
be allocated from a `/16` pool. IPv6 subnets with prefix length `/64` will
be allocated from a `/56` pool.

```json
{
  "default-address-pools": [
    { "base": "172.17.0.0/16", "size": 24 },
    { "base": "2001:db8::/56", "size": 64 }
  ]
}
```

> [!NOTE]
>
> The address `2001:db8::` in this example is
> [reserved for use in documentation][wikipedia-ipv6-reserved].
> Replace it with a valid IPv6 network.
>
> The default IPv4 pools are from the private address range,
> similar to the default IPv6 [ULA][wikipedia-ipv6-ula] networks.

See [Subnet allocation](../network/_index.md#subnet-allocation) for more information about
`default-address-pools`.

[wikipedia-ipv6-reserved]: https://en.wikipedia.org/wiki/Reserved_IP_addresses#IPv6
[wikipedia-ipv6-ula]: https://en.wikipedia.org/wiki/Unique_local_address

## Docker in Docker

On a host using `xtables` (legacy `iptables`) instead of `nftables`, kernel
module `ip6_tables` must be loaded before an IPv6 Docker network can be created,
It is normally loaded automatically when Docker starts.

However, if you running Docker in Docker that is not based on a recent
version of the [official `docker` image](https://hub.docker.com/_/docker), you
may need to run `modprobe ip6_tables` on your host. Alternatively, use daemon
option `--ip6tables=false` to disable `ip6tables` for the containerized Docker
Engine.

## Next steps

- [Networking overview](/manuals/engine/network/_index.md)
