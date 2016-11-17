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

You can set up your own bridge before starting Docker and use `-b BRIDGE` or
`--bridge=BRIDGE` to configure Docker to use your bridge instead. If you already
have Docker up and running with its default `docker0` still configured,
you can directly create your bridge and restart Docker with it or want to begin by
stopping the service and removing the interface:

1.  Stop Docker.

    ```bash
    $ sudo service docker stop
    ```

2.  Stop the `docker0` bridge.

    ```bash
    $ sudo ip link set dev docker0 down
    ```
3.  Delete the `docker0` bridge.

    ```bash
    $ sudo brctl delbr docker0
    ```

4.  Flush the `POSTROUTING` table from `iptables`

    ```bash
    $ sudo iptables -t nat -F POSTROUTING
    ```

Before restarting the Docker service, create and configure your own bridge, and
configure Docker to use it. This procedure creates a simple bridge to illustrate
the technique.

1.  Create a bridge called `mybridge0`.
    ```bash
    $ sudo brctl addbr mybridge0
    ```

2.  Configure the bridge.

    ```bash
    $ sudo ip addr add 192.168.5.1/24 dev mybridge0
    ```

3.  Start the bridge.
    ```bash
    $ sudo ip link set dev mybridge0 up
    ```

4.  Confirm that the bridge is up and running.
    ```bash
    $ ip addr show bridge0

    4: bridge0: <BROADCAST,MULTICAST> mtu 1500 qdisc noop state UP group default
        link/ether 66:38:d0:0d:76:18 brd ff:ff:ff:ff:ff:ff
        inet 192.168.5.1/24 scope global bridge0
           valid_lft forever preferred_lft forever
    ```

Finally, configure the Docker daemon to use the new bridge. These instructions
work for configuring Docker on systems that use `upstart` or `systemd`. If you
have  configured Docker to use a custom configuration file using the
`--config-file` flag, use that custom file when the instructions below refer to
`/etc/docker/daemon.json`. For an overview of available options for
`/etc/docker/daemon.json`, see
[Daemon configuration file](../reference/commandline/dockerd.md#daemon-configuration-file).

1.  Create or edit the `/etc/docker/daemon.json` file on your host.

    ```bash
    $ sudo nano /etc/docker/daemon.json
    ```

2.  Add the following option to use the new bridge. If this is a brand new file,
    you need to add the curly braces at the beginning and the end. Otherwise, just
    add the `bridge:` line.

    ```json
    {
      "bridge": "mybridge0"
    }
    ```

    Save and close the file.


3.  Start the Docker daemon.

    ```bash
    $ sudo service docker start
    ```

4.  Use the `brctl show` command to verify that the `docker0` bridge does not
    exist.

5.  Verify that Docker started correctly by running a `hello-world` container.

    ```bash
    $ docker run hello-world
    ```
6.  Confirm that outgoing NAT masquerading is set up.

    ```bash
    $ sudo iptables -t nat -L -n

    ...
    Chain POSTROUTING (policy ACCEPT)
    target     prot opt source               destination
    MASQUERADE  all  --  192.168.5.0/24      0.0.0.0/0
    ```

Docker can now bind containers to the new bridge. Try creating a container. Its
IP address is in your new IP address range, which Docker auto-detects.

Use the `brctl show` command to see the interfaces Docker adds and removes from
the bridge when you start and stop containers. Run `ip addr` and `ip route`
commands from within a container to see its address in the bridge's IP address
range verify that it uses the Docker host's IP address on the bridge as its
default gateway.
