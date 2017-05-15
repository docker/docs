---
description: Learn how to build your own bridge interface
keywords: docker, bridge, docker0, network
title: Build your own bridge
---

This section explains how to build your own bridge to replace the Docker default
bridge. This is a `bridge` network named `bridge` created automatically when you
install Docker.

> **Note**: The [Docker networks feature](../index.md) allows you to
create user-defined networks in addition to the default bridge network.

You can set up your own bridge before starting Docker and configure Docker to
use your bridge instead of the default `docker0` bridge.

1.  Configure the new bridge.

    ```bash
    $ sudo brctl addbr bridge0

    $ sudo ip addr add 192.168.5.1/24 dev bridge0

    $ sudo ip link set dev bridge0 up
    ```

    Confirm the new bridge's settings.

    ```bash
    $ ip addr show bridge0

    4: bridge0: <BROADCAST,MULTICAST> mtu 1500 qdisc noop state UP group default
        link/ether 66:38:d0:0d:76:18 brd ff:ff:ff:ff:ff:ff
        inet 192.168.5.1/24 scope global bridge0
           valid_lft forever preferred_lft forever
    ```

2.  Configure Docker to use the new bridge by setting the option in the
    `daemon.json` file, which is located in `/etc/docker/` on
    Linux or `C:\ProgramData\docker\config\` on Windows Server. On Docker for
    Mac or Docker for Windows, click the Docker icon, choose **Preferences**,
    and go to **Daemon**.

    If the `daemon.json` file does not exist, create it. Assuming there
    are no other settings in the file, it should have the following contents:

    ```json
    {
      "bridge": "bridge0"
    }
    ```

    Restart Docker for the changes to take effect.

3.  Confirm that the new outgoing NAT masquerade is set up.

    ```bash
    $ sudo iptables -t nat -L -n

    Chain POSTROUTING (policy ACCEPT)
    target     prot opt source               destination
    MASQUERADE  all  --  192.168.5.0/24      0.0.0.0/0
    ```

4.  Remove the now-unused `docker0` bridge.

    ```bash
    $ sudo ip link set dev docker0 down

    $ sudo brctl delbr docker0

    $ sudo iptables -t nat -F POSTROUTING
    ```

5.  Create a new container, and verify that it is in the new IP address range.


You can use the `brctl show` command to see Docker add and remove interfaces
from the bridge as you start and stop containers, and can run `ip addr` and `ip
route` inside a container to confirm that it has an address in the bridge's IP
address range and uses the Docker host's IP address on the bridge as its default
gateway to the rest of the Internet.
