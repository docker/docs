---
description: Customizing docker0
keywords: docker, bridge, docker0, network
title: Customize the docker0 bridge
---

The information in this section explains how to customize the Docker default
bridge. This is a `bridge` network named `bridge` created automatically when you
install Docker.

> **Note**: The [Docker networks feature](/engine/userguide/networking/index.md)
> allows you to create user-defined networks in addition to the default bridge network.

By default, the Docker server creates and configures the host system's `docker0`
a network interface called `docker0`, which is an ethernet bridge device. If you
don't specify a different network when starting a container, the container is
connected to the bridge and all traffic coming from and going to the container
flows over the bridge to the Docker daemon, which handles routing on behalf of
the container.

Docker configures `docker0` with an IP address, netmask, and IP allocation range.
Containers which are connected to the default bridge are allocated IP addresses
within this range. Certain default settings apply to the default bridge unless
you specify otherwise. For instance, the default maximum transmission unit (MTU),
or the largest packet length that the container will allow, defaults to 1500
bytes.

You can configure the default bridge network's settings using flags to the
`dockerd` command. However, the recommended way to configure the Docker daemon
is to use the `daemon.json` file, which is located in `/etc/docker/` on Linux.
If the file does not exist, create it. You can specify one or more of the
following settings to configure the default bridge network:

```json
{
  "bip": "192.168.1.5/24",
  "fixed-cidr": "10.20.0.0/16",
  "fixed-cidr-v6": "2001:db8::/64",
  "mtu": 1500,
  "default-gateway": "10.20.1.1",
  "default-gateway-v6": "2001:db8:abcd::89",
  "dns": ["10.20.1.2","10.20.1.3"]
}
```

Restart Docker after making changes to the `daemon.json` file.

The same options are presented as flags to `dockerd`, with an explanation for
each:

- `--bip=CIDR`: supply a specific IP address and netmask for the `docker0`
  bridge, using standard CIDR notation. For example: `192.168.1.5/24`.

- `--fixed-cidr=CIDR` and `--fixed-cidr-v6=CIDRv6`: restrict the IP range from
  the `docker0` subnet, using standard CIDR notation. For example:
  `172.16.1.0/28`. This range must be an IPv4 range for fixed IPs, such as
  `10.20.0.0/16`, and must be a subset of the bridge IP range (`docker0` or set
  using `--bridge`). For example, with `--fixed-cidr=192.168.1.0/25`, IPs for
  your containers will be chosen from the first half of addresses included in
  the `192.168.1.0/24` subnet.

- `--mtu=BYTES`: override the maximum packet length on `docker0`.

- `--default-gateway=Container default Gateway IPV4 address` and
  `--default-gateway-v6=Container default gateway IPV6 address`: designates the
  default gateway for containers connected to the `docker0` bridge, which
  controls where they route traffic by default. Applicable for addresses set
  with `--bip` and `--fixed-cidr` flags. For instance, you can configure
  `--fixed-cidr=172.17.2.0/24` and `default-gateway=172.17.1.1`.

- `--dns=[]`: The DNS servers to use. For example: `--dns=172.17.2.10`.

Once you have one or more containers up and running, you can confirm that Docker
has properly connected them to the `docker0` bridge by running the `brctl`
command on the host machine and looking at the `interfaces` column of the
output. This example shows a `docker0` bridge with two containers connected:

```bash
$ sudo brctl show

bridge name     bridge id               STP enabled     interfaces
docker0         8000.3a1d7362b4ee       no              veth65f9
                                                        vethdda6
```

If the `brctl` command is not installed on your Docker host, then on Ubuntu you
should be able to run `sudo apt-get install bridge-utils` to install it.

Finally, the `docker0` Ethernet bridge settings are used every time you create a
new container. Docker selects a free IP address from the range available on the
bridge each time you `docker run` a new container, and configures the
container's `eth0` interface with that IP address and the bridge's netmask. The
Docker host's own IP address on the bridge is used as the default gateway by
which each container reaches the rest of the Internet.

```bash
# The network, as seen from a container

$ docker run --rm -it alpine /bin/ash

root@f38c87f2a42d:/# ip addr show eth0

24: eth0: <BROADCAST,UP,LOWER_UP> mtu 1500 qdisc pfifo_fast state UP group default qlen 1000
    link/ether 32:6f:e0:35:57:91 brd ff:ff:ff:ff:ff:ff
    inet 172.17.0.3/16 scope global eth0
       valid_lft forever preferred_lft forever
    inet6 fe80::306f:e0ff:fe35:5791/64 scope link
       valid_lft forever preferred_lft forever

root@f38c87f2a42d:/# ip route

default via 172.17.42.1 dev eth0
172.17.0.0/16 dev eth0  proto kernel  scope link  src 172.17.0.3

root@f38c87f2a42d:/# exit
```

Remember that the Docker host will not be willing to forward container packets
out on to the Internet unless its `ip_forward` system setting is `1` -- see the
section on
[Communicating to the outside world](container-communication.md#communicating-to-the-outside-world)
for details.
