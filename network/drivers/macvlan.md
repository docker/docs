---
title: Macvlan network driver
description: All about using Macvlan to make your containers appear like physical
  machines on the network
keywords: network, macvlan, standalone
redirect_from:
  - /config/containers/macvlan/
  - /engine/userguide/networking/get-started-macvlan/
  - /network/macvlan/
---

Some applications, especially legacy applications or applications which monitor
network traffic, expect to be directly connected to the physical network. In
this type of situation, you can use the `macvlan` network driver to assign a MAC
address to each container's virtual network interface, making it appear to be
a physical network interface directly connected to the physical network. In this
case, you need to designate a physical interface on your Docker host to use for
the Macvlan, as well as the subnet and gateway of the network. You can even
isolate your Macvlan networks using different physical network interfaces.

Keep the following things in mind:

- You may unintentionally degrade your network due to IP address
  exhaustion or to "VLAN spread", a situation that occurs when you have an
  inappropriately large number of unique MAC addresses in your network.

- Your networking equipment needs to be able to handle "promiscuous mode",
  where one physical interface can be assigned multiple MAC addresses.

- If your application can work using a bridge (on a single Docker host) or
  overlay (to communicate across multiple Docker hosts), these solutions may be
  better in the long term.

## Options

The following table describes the driver-specific options that you can pass to
`--option` when creating a network using the `macvlan` driver.

| Option         | Default  | Description                                                                   |
| -------------- | -------- | ----------------------------------------------------------------------------- |
| `macvlan_mode` | `bridge` | Sets the Macvlan mode. Can be one of: `bridge`, `vepa`, `passthru`, `private` |
| `parent`       |          | Specifies the parent interface to use.                                        |

## Create a Macvlan network

When you create a Macvlan network, it can either be in bridge mode or 802.1Q
trunk bridge mode.

- In bridge mode, Macvlan traffic goes through a physical device on the host.

- In 802.1Q trunk bridge mode, traffic goes through an 802.1Q sub-interface
  which Docker creates on the fly. This allows you to control routing and
  filtering at a more granular level.

### Bridge mode

To create a `macvlan` network which bridges with a given physical network
interface, use `--driver macvlan` with the `docker network create` command. You
also need to specify the `parent`, which is the interface the traffic will
physically go through on the Docker host.

```console
$ docker network create -d macvlan \
  --subnet=172.16.86.0/24 \
  --gateway=172.16.86.1 \
  -o parent=eth0 pub_net
```

If you need to exclude IP addresses from being used in the `macvlan` network, such
as when a given IP address is already in use, use `--aux-addresses`:

```console
$ docker network create -d macvlan \
  --subnet=192.168.32.0/24 \
  --ip-range=192.168.32.128/25 \
  --gateway=192.168.32.254 \
  --aux-address="my-router=192.168.32.129" \
  -o parent=eth0 macnet32
```

### 802.1Q trunk bridge mode

If you specify a `parent` interface name with a dot included, such as `eth0.50`,
Docker interprets that as a sub-interface of `eth0` and creates the sub-interface
automatically.

```console
$ docker network create -d macvlan \
    --subnet=192.168.50.0/24 \
    --gateway=192.168.50.1 \
    -o parent=eth0.50 macvlan50
```

### Use an IPvlan instead of Macvlan

In the above example, you are still using a L3 bridge. You can use `ipvlan`
instead, and get an L2 bridge. Specify `-o ipvlan_mode=l2`.

```console
$ docker network create -d ipvlan \
    --subnet=192.168.210.0/24 \
    --subnet=192.168.212.0/24 \
    --gateway=192.168.210.254 \
    --gateway=192.168.212.254 \
     -o ipvlan_mode=l2 -o parent=eth0 ipvlan210
```

## Use IPv6

If you have [configured the Docker daemon to allow IPv6](../../config/daemon/ipv6.md),
you can use dual-stack IPv4/IPv6 `macvlan` networks.

```console
$ docker network create -d macvlan \
    --subnet=192.168.216.0/24 --subnet=192.168.218.0/24 \
    --gateway=192.168.216.1 --gateway=192.168.218.1 \
    --subnet=2001:db8:abc8::/64 --gateway=2001:db8:abc8::10 \
     -o parent=eth0.218 \
     -o macvlan_mode=bridge macvlan216
```

## Next steps

Learn how to use the Macvlan driver in the
[Macvlan networking tutorial](../network-tutorial-macvlan.md)
