---
title: Networking using a macvlan network
description: Tutorials for networking using a macvlan bridge network and 802.1q trunk bridge network
keywords: networking, macvlan, 802.1q, standalone
---

This series of tutorials deals with networking standalone containers which
connect to `macvlan` networks. In this type of network, the Docker host accepts
requests for multiple MAC addresses at its IP address, and routes those requests
to the appropriate container. For other networking topics, see the
[overview](index.md).

## Goal

The goal of these tutorials is to set up a bridged `macvlan` network and attach
a container to it, then set up an 802.1q trunked `macvlan` network and attach a
container to it.

## Prerequisites

- Most cloud providers block `macvlan` networking. You may need physical access
  to your networking equipment.

- The `macvlan` networking driver only works on Linux hosts, and is not supported
  on Docker for Mac, Docker for Windows, or Docker EE for Windows Server.

- You need at least version 3.9 of the Linux kernel, and version 4.0 or higher
  is recommended.

- The examples assume your ethernet interface is `eth0`. If your device has a
  different name, use that instead.

## Bridge example

In the simple bridge example, your traffic flows through `eth0` and Docker
routes traffic to your container using its MAC address. To network devices
on your network, your container appears to be physically attached to the network.

1.  Create a `macvlan` network called `my-macvlan-net`. Modify the `subnet`, `gateway`,
    and `parent` values to values that make sense in your environment.

    ```bash
    $ docker network create -d macvlan \
      --subnet=172.16.86.0/24 \
      --gateway=172.16.86.1 \
      -o parent=eth0 \
      my-macvlan-net
    ```

    You can use `docker network ls` and `docker network inspect pub_net`
    commands to verify that the network exists and is a `macvlan` network.

2.  Start an `alpine` container and attach it to the `my-macvlan-net` network. The
    `-dit` flags start the container in the background but allow you to attach
    to it. The `--rm` flag means the container is removed when it is stopped.

    ```bash
    $ docker run --rm -itd \
      --network my-macvlan-net \
      --name my-macvlan-alpine \
      alpine:latest \
      ash
    ```

3.  Inspect the `my-macvlan-alpine` container and notice the `MacAddress` key
    within the `Networks` key:

    ```none
    $ docker container inspect my-macvlan-alpine

    ...truncated...
    "Networks": {
      "my-macvlan-net": {
          "IPAMConfig": null,
          "Links": null,
          "Aliases": [
              "bec64291cd4c"
          ],
          "NetworkID": "5e3ec79625d388dbcc03dcf4a6dc4548644eb99d58864cf8eee2252dcfc0cc9f",
          "EndpointID": "8caf93c862b22f379b60515975acf96f7b54b7cf0ba0fb4a33cf18ae9e5c1d89",
          "Gateway": "172.16.86.1",
          "IPAddress": "172.16.86.2",
          "IPPrefixLen": 24,
          "IPv6Gateway": "",
          "GlobalIPv6Address": "",
          "GlobalIPv6PrefixLen": 0,
          "MacAddress": "02:42:ac:10:56:02",
          "DriverOpts": null
      }
    }
    ...truncated
    ```

4.  Check out how the container sees its own network interfaces by running a
    couple of `docker exec` commands.

    ```bash
    $ docker exec my-macvlan-alpine ip addr show eth0

    9: eth0@tunl0: <BROADCAST,MULTICAST,UP,LOWER_UP,M-DOWN> mtu 1500 qdisc noqueue state UP
    link/ether 02:42:ac:10:56:02 brd ff:ff:ff:ff:ff:ff
    inet 172.16.86.2/24 brd 172.16.86.255 scope global eth0
       valid_lft forever preferred_lft forever
    ```

    ```bash
    $ docker exec my-macvlan-alpine ip route

    default via 172.16.86.1 dev eth0
    172.16.86.0/24 dev eth0 scope link  src 172.16.86.2
    ```

5.  Stop the container (Docker removes it because of the `--rm` flag), and remove
    the network.

    ```bash
    $ docker container stop my-macvlan-alpine

    $ docker network rm my-macvlan-net
    ```

## 802.1q trunked bridge example

In the 802.1q trunked bridge example, your traffic flows through a sub-interface
of `eth0` (called `eth0.10`) and Docker routes traffic to your container using
its MAC address. To network devices on your network, your container appears to
be physically attached to the network.

1.  Create a `macvlan` network called `my-8021q-macvlan-net`. Modify the
    `subnet`, `gateway`, and `parent` values to values that make sense in your
    environment.

    ```bash
    $ docker network create -d macvlan \
      --subnet=172.16.86.0/24 \
      --gateway=172.16.86.1 \
      -o parent=eth0.10 \
      my-8021q-macvlan-net
    ```

    You can use `docker network ls` and `docker network inspect pub_net`
    commands to verify that the network exists, is a `macvlan` network, and
    has parent `eth0.10`. You can use `ip addr show` on the Docker host to
    verify that the interface `eth0.10` exists and has a separate IP address

2.  Start an `alpine` container and attach it to the `my-8021q-macvlan-net`
    network. The `-dit` flags start the container in the background but allow
    you to attach to it. The `--rm` flag means the container is removed when it
    is stopped.

    ```bash
    $ docker run --rm -itd \
      --network my-8021q-macvlan-net \
      --name my-second-macvlan-alpine \
      alpine:latest \
      ash
    ```

3.  Inspect the `my-second-macvlan-alpine` container and notice the `MacAddress`
    key within the `Networks` key:

    ```none
    $ docker container inspect my-second-macvlan-alpine

    ...truncated...
    "Networks": {
      "my-8021q-macvlan-net": {
          "IPAMConfig": null,
          "Links": null,
          "Aliases": [
              "12f5c3c9ba5c"
          ],
          "NetworkID": "c6203997842e654dd5086abb1133b7e6df627784fec063afcbee5893b2bb64db",
          "EndpointID": "aa08d9aa2353c68e8d2ae0bf0e11ed426ea31ed0dd71c868d22ed0dcf9fc8ae6",
          "Gateway": "172.16.86.1",
          "IPAddress": "172.16.86.2",
          "IPPrefixLen": 24,
          "IPv6Gateway": "",
          "GlobalIPv6Address": "",
          "GlobalIPv6PrefixLen": 0,
          "MacAddress": "02:42:ac:10:56:02",
          "DriverOpts": null
      }
    }
    ...truncated
    ```

4.  Check out how the container sees its own network interfaces by running a
    couple of `docker exec` commands.

    ```bash
    $ docker exec my-second-macvlan-alpine ip addr show eth0

    11: eth0@if10: <BROADCAST,MULTICAST,UP,LOWER_UP,M-DOWN> mtu 1500 qdisc noqueue state UP
    link/ether 02:42:ac:10:56:02 brd ff:ff:ff:ff:ff:ff
    inet 172.16.86.2/24 brd 172.16.86.255 scope global eth0
       valid_lft forever preferred_lft forever
    ```

    ```bash
    $ docker exec my-second-macvlan-alpine ip route

    default via 172.16.86.1 dev eth0
    172.16.86.0/24 dev eth0 scope link  src 172.16.86.2
    ```

5.  Stop the container (Docker removes it because of the `--rm` flag), and remove
    the network.

    ```bash
    $ docker container stop my-second-macvlan-alpin

    $ docker network rm my-8021q-macvlan-net
    ```

## Other networking tutorials

Now that you have completed the networking tutorial for `macvlan` networks,
you might want to run through these other networking tutorials:

- [Standalone networking tutorial](network-tutorial-standalone.md)
- [Overlay networking tutorial](network-tutorial-overlay.md)
- [Host networking tutorial](network-tutorial-host.md)

