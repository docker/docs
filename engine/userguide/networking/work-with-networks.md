---
description: How to work with docker networks
keywords: commands, Usage, network, docker, cluster
title: Work with network commands
---

This article provides examples of the network subcommands you can use to
interact with Docker networks and the containers in them. The commands are
available through the Docker Engine CLI. These commands are:

* `docker network create`
* `docker network connect`
* `docker network ls`
* `docker network rm`
* `docker network disconnect`
* `docker network inspect`

While not required, it is a good idea to read [Understanding Docker
network](index.md) before trying the examples in this section. The
examples use the default `bridge` network so that you can try them
immediately. To experiment with an `overlay` network, check out
the [Getting started with multi-host networks](get-started-overlay.md) guide instead.

## Create networks

Docker Engine creates a `bridge` network automatically when you install Engine.
This network corresponds to the `docker0` bridge that Engine has traditionally
relied on. In addition to this network, you can create your own `bridge` or
`overlay` network.

A `bridge` network resides on a single host running an instance of Docker
Engine. An `overlay` network can span multiple hosts running their own engines.
If you run `docker network create` and supply only a network name, it creates a
bridge network for you.

```bash
$ docker network create simple-network

69568e6336d8c96bbf57869030919f7c69524f71183b44d80948bd3927c87f6a

$ docker network inspect simple-network
[
    {
        "Name": "simple-network",
        "Id": "69568e6336d8c96bbf57869030919f7c69524f71183b44d80948bd3927c87f6a",
        "Scope": "local",
        "Driver": "bridge",
        "IPAM": {
            "Driver": "default",
            "Config": [
                {
                    "Subnet": "172.22.0.0/16",
                    "Gateway": "172.22.0.1"
                }
            ]
        },
        "Containers": {},
        "Options": {},
        "Labels": {}
    }
]
```

Unlike `bridge` networks, `overlay` networks require some pre-existing conditions
before you can create one. These conditions are:

* Access to a key-value store. Engine supports Consul, Etcd, and ZooKeeper (Distributed store) key-value stores.
* A cluster of hosts with connectivity to the key-value store.
* A properly configured Engine `daemon` on each host in the swarm.

The `dockerd` options that support the `overlay` network are:

* `--cluster-store`
* `--cluster-store-opt`
* `--cluster-advertise`

When you create a network, Engine creates a non-overlapping subnetwork for the
network by default. You can override this default and specify a subnetwork
directly using the `--subnet` option. On a `bridge` network you can only
specify a single subnet. An `overlay` network supports multiple subnets.

> **Note** : It is highly recommended to use the `--subnet` option while creating
> a network. If the `--subnet` is not specified, the docker daemon automatically
> chooses and assigns a subnet for the network and it could overlap with another subnet
> in your infrastructure that is not managed by docker. Such overlaps can cause
> connectivity issues or failures when containers are connected to that network.

In addition to the `--subnet` option, you also specify the `--gateway`,
`--ip-range`, and `--aux-address` options.

```bash
$ docker network create -d overlay \
  --subnet=192.168.0.0/16 \
  --subnet=192.170.0.0/16 \
  --gateway=192.168.0.100 \
  --gateway=192.170.0.100 \
  --ip-range=192.168.1.0/24 \
  --aux-address="my-router=192.168.1.5" --aux-address="my-switch=192.168.1.6" \
  --aux-address="my-printer=192.170.1.5" --aux-address="my-nas=192.170.1.6" \
  my-multihost-network
```

Be sure that your subnetworks do not overlap. If they do, network creation
fails and Engine returns an error.

When creating a custom network, you can pass additional options to the driver.
The `bridge` driver accepts the following options:

| Option                                           | Equivalent  | Description                                           |
|--------------------------------------------------|-------------|-------------------------------------------------------|
| `com.docker.network.bridge.name`                 | -           | bridge name to be used when creating the Linux bridge |
| `com.docker.network.bridge.enable_ip_masquerade` | `--ip-masq` | Enable IP masquerading                                |
| `com.docker.network.bridge.enable_icc`           | `--icc`     | Enable or Disable Inter Container Connectivity        |
| `com.docker.network.bridge.host_binding_ipv4`    | `--ip`      | Default IP when binding container ports               |
| `com.docker.network.driver.mtu`                  | `--mtu`     | Set the containers network MTU                        |

The `com.docker.network.driver.mtu` option is also supported by the `overlay` driver.

The following arguments can be passed to `docker network create` for any network driver.

| Argument     | Equivalent | Description                              |
|--------------|------------|------------------------------------------|
| `--internal` | -          | Restricts external access to the network |
| `--ipv6`     | `--ipv6`   | Enable IPv6 networking                   |

The following example uses `-o` to bind to a specific IP address when binding
ports, then uses `docker network inspect` to inspect the network, and finally
attaches a new container to the new network.

```bash
$ docker network create -o "com.docker.network.bridge.host_binding_ipv4"="172.23.0.1" my-network

b1a086897963e6a2e7fc6868962e55e746bee8ad0c97b54a5831054b5f62672a

$ docker network inspect my-network

[
    {
        "Name": "my-network",
        "Id": "b1a086897963e6a2e7fc6868962e55e746bee8ad0c97b54a5831054b5f62672a",
        "Scope": "local",
        "Driver": "bridge",
        "IPAM": {
            "Driver": "default",
            "Options": {},
            "Config": [
                {
                    "Subnet": "172.23.0.0/16",
                    "Gateway": "172.23.0.1"
                }
            ]
        },
        "Containers": {},
        "Options": {
            "com.docker.network.bridge.host_binding_ipv4": "172.23.0.1"
        },
        "Labels": {}
    }
]

$ docker run -d -P --name redis --network my-network redis

bafb0c808c53104b2c90346f284bda33a69beadcab4fc83ab8f2c5a4410cd129

$ docker ps

CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS              PORTS                        NAMES
bafb0c808c53        redis               "/entrypoint.sh redis"   4 seconds ago       Up 3 seconds        172.23.0.1:32770->6379/tcp   redis
```

## Connect containers

You can connect an existing container to one or more networks. A container can
connect to networks which use different network drivers. Once connected, the
containers can communicate using another container's IP address or name.

For `overlay` networks or custom plugins that support multi-host
connectivity, containers connected to the same multi-host network but launched
from different hosts can also communicate in this way.

This example uses six containers, and directs you to create them as they are
needed.

### Basic container networking example

1.  First, create and run two containers, `container1` and `container2`:

    ```bash
    $ docker run -itd --name=container1 busybox

    18c062ef45ac0c026ee48a83afa39d25635ee5f02b58de4abc8f467bcaa28731

    $ docker run -itd --name=container2 busybox

    498eaaaf328e1018042c04b2de04036fc04719a6e39a097a4f4866043a2c2152
    ```

2.  Create an isolated, `bridge` network to test with.

    ```bash
    $ docker network create -d bridge --subnet 172.25.0.0/16 isolated_nw

    06a62f1c73c4e3107c0f555b7a5f163309827bfbbf999840166065a8f35455a8
    ```

3.  Connect `container2` to the network and then `inspect` the network to verify
    the connection:

    ```bash
    $ docker network connect isolated_nw container2

    $ docker network inspect isolated_nw

    [
        {
            "Name": "isolated_nw",
            "Id": "06a62f1c73c4e3107c0f555b7a5f163309827bfbbf999840166065a8f35455a8",
            "Scope": "local",
            "Driver": "bridge",
            "IPAM": {
                "Driver": "default",
                "Config": [
                    {
                        "Subnet": "172.25.0.0/16",
                        "Gateway": "172.25.0.1/16"
                    }
                ]
            },
            "Containers": {
                "90e1f3ec71caf82ae776a827e0712a68a110a3f175954e5bd4222fd142ac9428": {
                    "Name": "container2",
                    "EndpointID": "11cedac1810e864d6b1589d92da12af66203879ab89f4ccd8c8fdaa9b1c48b1d",
                    "MacAddress": "02:42:ac:19:00:02",
                    "IPv4Address": "172.25.0.2/16",
                    "IPv6Address": ""
                }
            },
            "Options": {}
        }
    ]
    ```

    Notice that `container2` is assigned an IP address automatically. Because
    you specified a `--subnet` when creating the network, the IP address was
    chosen from that subnet.

    As a reminder, `container1` is only connected to the default `bridge` network.

4.  Start a third container, but this time assign it an IP address using the
    `--ip` flag and connect it to the `isolated_nw` network using the `docker run`
    command's `--network` option:

    ```bash
    $ docker run --network=isolated_nw --ip=172.25.3.3 -itd --name=container3 busybox

    467a7863c3f0277ef8e661b38427737f28099b61fa55622d6c30fb288d88c551
    ```

    As long as the IP address you specify for the container is part of the
    network's subnet, you can assign an IPv4 or IPv6 address to a container
    when connecting it to a network, by using the `--ip` or `--ip6` flag. when
    you specify an IP address in this way while using a user-defined network,
    the configuration is preserved as part of the container's configuration and
    will be applied when the container is reloaded. Assigned IP addresses are
    preserved when using non-user-defined networks, because there is no guarantee
    that a container's subnet will not change when the Docker daemon restarts unless
    you use user-defined networks.

5.  Inspect the network resources used by `container3`. The
    output below is truncated for brevity.

    ```bash
    $ docker inspect --format='{{json .NetworkSettings.Networks}}'  container3

    {"isolated_nw":
      {"IPAMConfig":
        {
          "IPv4Address":"172.25.3.3"},
          "NetworkID":"1196a4c5af43a21ae38ef34515b6af19236a3fc48122cf585e3f3054d509679b",
          "EndpointID":"dffc7ec2915af58cc827d995e6ebdc897342be0420123277103c40ae35579103",
          "Gateway":"172.25.0.1",
          "IPAddress":"172.25.3.3",
          "IPPrefixLen":16,
          "IPv6Gateway":"",
          "GlobalIPv6Address":"",
          "GlobalIPv6PrefixLen":0,
          "MacAddress":"02:42:ac:19:03:03"}
        }
      }
    }
    ```

    Because you connected `container3` to the `isolated_nw` when you started it,
    it is not connected to the default `bridge` network at all.

6. Inspect the network resources used by `container2`. If you have Python
   installed, you can pretty print the output.

   ```bash
   $ docker inspect --format='{{json .NetworkSettings.Networks}}'  container2 | python -m json.tool

   {
       "bridge": {
           "NetworkID":"7ea29fc1412292a2d7bba362f9253545fecdfa8ce9a6e37dd10ba8bee7129812",
           "EndpointID": "0099f9efb5a3727f6a554f176b1e96fca34cae773da68b3b6a26d046c12cb365",
           "Gateway": "172.17.0.1",
           "GlobalIPv6Address": "",
           "GlobalIPv6PrefixLen": 0,
           "IPAMConfig": null,
           "IPAddress": "172.17.0.3",
           "IPPrefixLen": 16,
           "IPv6Gateway": "",
           "MacAddress": "02:42:ac:11:00:03"
       },
       "isolated_nw": {
           "NetworkID":"1196a4c5af43a21ae38ef34515b6af19236a3fc48122cf585e3f3054d509679b",
           "EndpointID": "11cedac1810e864d6b1589d92da12af66203879ab89f4ccd8c8fdaa9b1c48b1d",
           "Gateway": "172.25.0.1",
           "GlobalIPv6Address": "",
           "GlobalIPv6PrefixLen": 0,
           "IPAMConfig": null,
           "IPAddress": "172.25.0.2",
           "IPPrefixLen": 16,
           "IPv6Gateway": "",
           "MacAddress": "02:42:ac:19:00:02"
       }
   }
   ```

   Notice that `container2` belongs to two networks.  It joined the default `bridge`
   network when you launched it and you connected it to the `isolated_nw` in
   step 3.

   ![](images/working.png)

    eth0      Link encap:Ethernet  HWaddr 02:42:AC:11:00:03

    eth1    Link encap:Ethernet  HWaddr 02:42:AC:15:00:02

7.  Use the `docker attach` command to connect to the running `container2` and
    examine its networking stack:

    ```bash
    $ docker attach container2
    ```

    Use the `ifconfig` command to examine the container's networking stack. you
    should see two ethernet interfaces, one for the default `bridge` network,
    and the other for the `isolated_nw` network.

    ```bash
    $ sudo ifconfig -a

    eth0      Link encap:Ethernet  HWaddr 02:42:AC:11:00:03
              inet addr:172.17.0.3  Bcast:0.0.0.0  Mask:255.255.0.0
              inet6 addr: fe80::42:acff:fe11:3/64 Scope:Link
              UP BROADCAST RUNNING MULTICAST  MTU:9001  Metric:1
              RX packets:8 errors:0 dropped:0 overruns:0 frame:0
              TX packets:8 errors:0 dropped:0 overruns:0 carrier:0
              collisions:0 txqueuelen:0
              RX bytes:648 (648.0 B)  TX bytes:648 (648.0 B)

    eth1      Link encap:Ethernet  HWaddr 02:42:AC:15:00:02
              inet addr:172.25.0.2  Bcast:0.0.0.0  Mask:255.255.0.0
              inet6 addr: fe80::42:acff:fe19:2/64 Scope:Link
              UP BROADCAST RUNNING MULTICAST  MTU:1500  Metric:1
              RX packets:8 errors:0 dropped:0 overruns:0 frame:0
              TX packets:8 errors:0 dropped:0 overruns:0 carrier:0
              collisions:0 txqueuelen:0
              RX bytes:648 (648.0 B)  TX bytes:648 (648.0 B)

    lo        Link encap:Local Loopback
              inet addr:127.0.0.1  Mask:255.0.0.0
              inet6 addr: ::1/128 Scope:Host
              UP LOOPBACK RUNNING  MTU:65536  Metric:1
              RX packets:0 errors:0 dropped:0 overruns:0 frame:0
              TX packets:0 errors:0 dropped:0 overruns:0 carrier:0
              collisions:0 txqueuelen:0
              RX bytes:0 (0.0 B)  TX bytes:0 (0.0 B)
    ```

8.  The Docker embedded DNS server enables name resolution for containers
    connected to a given network. This means that any connected container can
    ping another container on the same network by its container name. From
    within `container2`, you can ping `container3` by name.

    ```bash
    / # ping -w 4 container3
    PING container3 (172.25.3.3): 56 data bytes
    64 bytes from 172.25.3.3: seq=0 ttl=64 time=0.070 ms
    64 bytes from 172.25.3.3: seq=1 ttl=64 time=0.080 ms
    64 bytes from 172.25.3.3: seq=2 ttl=64 time=0.080 ms
    64 bytes from 172.25.3.3: seq=3 ttl=64 time=0.097 ms

    --- container3 ping statistics ---
    4 packets transmitted, 4 packets received, 0% packet loss
    round-trip min/avg/max = 0.070/0.081/0.097 ms
    ```

    This functionality is not available for the default `bridge` network. Both
    `container1` and `container2` are connected to the `bridge` network, but
    you cannot ping `container1` from `container2` using the container name.


    ```bash
    / # ping -w 4 container1
    ping: bad address 'container1'
    ```

    You can still ping the IP address directly:

    ```bash
    / # ping -w 4 172.17.0.2
    PING 172.17.0.2 (172.17.0.2): 56 data bytes
    64 bytes from 172.17.0.2: seq=0 ttl=64 time=0.095 ms
    64 bytes from 172.17.0.2: seq=1 ttl=64 time=0.075 ms
    64 bytes from 172.17.0.2: seq=2 ttl=64 time=0.072 ms
    64 bytes from 172.17.0.2: seq=3 ttl=64 time=0.101 ms

    --- 172.17.0.2 ping statistics ---
    4 packets transmitted, 4 packets received, 0% packet loss
    round-trip min/avg/max = 0.072/0.085/0.101 ms
    ```

    Detach from `container2` and leave it running using `CTRL-p CTRL-q`.

9.  Currently, `container2` is attached to both `bridge` and `isolated_nw`
    networks, so it can communicate with both `container1` and `container3`.
    However, `container3` and `container1` do not have any networks in common,
    so they cannot communicate. To verify this, attach to `container3` and try
    to ping `container1` by IP address.

    ```bash
    $ docker attach container3

    $ ping 172.17.0.2
    PING 172.17.0.2 (172.17.0.2): 56 data bytes
    ^C

    --- 172.17.0.2 ping statistics ---
    10 packets transmitted, 0 packets received, 100% packet loss

    ```

    Detach from `container3` and leave it running using `CTRL-p CTRL-q`.

>You can connect a container to a network even if the container is not running.
However, `docker network inspect` only displays information on running containers.

### Linking containers without using user-defined networks

After you complete the steps in
[Basic container networking examples](#basic-container-networking-examples),
`container2` can resolve `container3`'s name automatically because both containers
are connected to the `isolated_nw` network. However, containers connected to the
default `bridge` network cannot resolve each other's container name. If you need
containers to be able to communicate with each other over the `bridge` network,
you need to use the legacy [link](default_network/dockerlinks.md) feature.
This is the only use case where using `--link` is recommended. You should
strongly consider using user-defined networks instead.

Using the legacy `link` flag adds the following features for communication
between communication on the default `bridge` network:

* the ability to resolve container names to IP addresses
* the ability to define a network alias as an alternate way to refer to the linked container, using `--link=CONTAINER-NAME:ALIAS`
* secured container connectivity (in isolation via `--icc=false`)
* environment variable injection

To reiterate, all of these features are provided by default when you use a
user-defined network, with no additional configuration required. **Additionally,
you get the ability to dynamically attach to and detach from multiple networks.**

* automatic name resolution using DNS
* supports the `--link` option to provide name alias for the linked container
* automatic secured isolated environment for the containers in a network
* environment variable injection

The following example briefly describes how to use `--link`.

1.  Continuing with the above example, create a new container, `container4`, and
    connect it to the network `isolated_nw`. In addition, link it to container
    `container5` (which does not exist yet!) using the `--link` flag.

    ```bash
    $ docker run --network=isolated_nw -itd --name=container4 --link container5:c5 busybox

    01b5df970834b77a9eadbaff39051f237957bd35c4c56f11193e0594cfd5117c
    ```

    This is a little tricky, because `container5` does not exist yet. When
    `container5` is created, `container4` will be able to resolve the name `c5` to
    `container5`'s IP address.

    >**Note**: Any link between containers created with *legacy link* is static in
    nature and hard-binds the container with the alias. It does not tolerate
    linked container restarts. The new *link* functionality in user defined
    networks supports dynamic links between containers, and tolerates restarts and
    IP address changes in the linked container.

    Since you have not yet created container `container5` trying to ping it will result
    in an error. Attach to `container4` and try to ping either `container5` or `c5`:

    ```bash
    $ docker attach container4

    $ ping container5

    ping: bad address 'container5'

    $ ping c5

    ping: bad address 'c5'

    ```
    Detach from `container4` and leave it running using `CTRL-p CTRL-q`.

2.  Create another container named `container5`, and link it to `container4`
    using the alias `c4`.

    ```bash
    $ docker run --network=isolated_nw -itd --name=container5 --link container4:c4 busybox

    72eccf2208336f31e9e33ba327734125af00d1e1d2657878e2ee8154fbb23c7a
    ```

    Now attach to `container4` and try to ping `c5` and `container5`.

    ```bash
    $ docker attach container4

    / # ping -w 4 c5
    PING c5 (172.25.0.5): 56 data bytes
    64 bytes from 172.25.0.5: seq=0 ttl=64 time=0.070 ms
    64 bytes from 172.25.0.5: seq=1 ttl=64 time=0.080 ms
    64 bytes from 172.25.0.5: seq=2 ttl=64 time=0.080 ms
    64 bytes from 172.25.0.5: seq=3 ttl=64 time=0.097 ms

    --- c5 ping statistics ---
    4 packets transmitted, 4 packets received, 0% packet loss
    round-trip min/avg/max = 0.070/0.081/0.097 ms

    / # ping -w 4 container5
    PING container5 (172.25.0.5): 56 data bytes
    64 bytes from 172.25.0.5: seq=0 ttl=64 time=0.070 ms
    64 bytes from 172.25.0.5: seq=1 ttl=64 time=0.080 ms
    64 bytes from 172.25.0.5: seq=2 ttl=64 time=0.080 ms
    64 bytes from 172.25.0.5: seq=3 ttl=64 time=0.097 ms

    --- container5 ping statistics ---
    4 packets transmitted, 4 packets received, 0% packet loss
    round-trip min/avg/max = 0.070/0.081/0.097 ms
    ```
    Detach from `container4` and leave it running using `CTRL-p CTRL-q`.

3.  Finally, attach to `container5` and verify that you can ping `container4`.

    ```bash
    $ docker attach container5

    / # ping -w 4 c4
    PING c4 (172.25.0.4): 56 data bytes
    64 bytes from 172.25.0.4: seq=0 ttl=64 time=0.065 ms
    64 bytes from 172.25.0.4: seq=1 ttl=64 time=0.070 ms
    64 bytes from 172.25.0.4: seq=2 ttl=64 time=0.067 ms
    64 bytes from 172.25.0.4: seq=3 ttl=64 time=0.082 ms

    --- c4 ping statistics ---
    4 packets transmitted, 4 packets received, 0% packet loss
    round-trip min/avg/max = 0.065/0.070/0.082 ms

    / # ping -w 4 container4
    PING container4 (172.25.0.4): 56 data bytes
    64 bytes from 172.25.0.4: seq=0 ttl=64 time=0.065 ms
    64 bytes from 172.25.0.4: seq=1 ttl=64 time=0.070 ms
    64 bytes from 172.25.0.4: seq=2 ttl=64 time=0.067 ms
    64 bytes from 172.25.0.4: seq=3 ttl=64 time=0.082 ms

    --- container4 ping statistics ---
    4 packets transmitted, 4 packets received, 0% packet loss
    round-trip min/avg/max = 0.065/0.070/0.082 ms
    ```
    Detach from `container5` and leave it running using `CTRL-p CTRL-q`.

### Network alias scoping example

When you link containers, whether using the legacy `link` method or using
user-defined networks, any aliases you specify only have meaning to the
container where they are specified, and won't work on other containers on the
default `bridge` network.

In addition, if a container belongs to multiple networks, a given linked alias
is scoped within a given network. Thus, a container can be linked to different
aliases in different networks, and the aliases will not work for containers which
are not on the same network.

The following example illustrates these points.

1.  Create another network named `local_alias`:

    ```bash
    $ docker network create -d bridge --subnet 172.26.0.0/24 local_alias
    76b7dc932e037589e6553f59f76008e5b76fa069638cd39776b890607f567aaa
    ```

2.  Next, connect `container4` and `container5` to the new network `local_alias`
    with the aliases `foo` and `bar`:

    ```bash
    $ docker network connect --link container5:foo local_alias container4
    $ docker network connect --link container4:bar local_alias container5
    ```

3. Attach to `container4` and try to ping `container4` (yes, the same one) using alias `foo`, then
   try pinging container `container5` using alias `c5`:

    ```bash
    $ docker attach container4

    / # ping -w 4 foo
    PING foo (172.26.0.3): 56 data bytes
    64 bytes from 172.26.0.3: seq=0 ttl=64 time=0.070 ms
    64 bytes from 172.26.0.3: seq=1 ttl=64 time=0.080 ms
    64 bytes from 172.26.0.3: seq=2 ttl=64 time=0.080 ms
    64 bytes from 172.26.0.3: seq=3 ttl=64 time=0.097 ms

    --- foo ping statistics ---
    4 packets transmitted, 4 packets received, 0% packet loss
    round-trip min/avg/max = 0.070/0.081/0.097 ms

    / # ping -w 4 c5
    PING c5 (172.25.0.5): 56 data bytes
    64 bytes from 172.25.0.5: seq=0 ttl=64 time=0.070 ms
    64 bytes from 172.25.0.5: seq=1 ttl=64 time=0.080 ms
    64 bytes from 172.25.0.5: seq=2 ttl=64 time=0.080 ms
    64 bytes from 172.25.0.5: seq=3 ttl=64 time=0.097 ms

    --- c5 ping statistics ---
    4 packets transmitted, 4 packets received, 0% packet loss
    round-trip min/avg/max = 0.070/0.081/0.097 ms
    ```

    Both pings succeed, but the subnets are different, which means that the
    networks are different.

    Detach from `container4` and leave it running using `CTRL-p CTRL-q`.

4.  Disconnect `container5` from the `isolated_nw` network. Attach to `container4`
    and try pinging `c5` and `foo`.

    ```bash
    $ docker network disconnect isolated_nw container5

    $ docker attach container4

    / # ping -w 4 c5
    ping: bad address 'c5'

    / # ping -w 4 foo
    PING foo (172.26.0.3): 56 data bytes
    64 bytes from 172.26.0.3: seq=0 ttl=64 time=0.070 ms
    64 bytes from 172.26.0.3: seq=1 ttl=64 time=0.080 ms
    64 bytes from 172.26.0.3: seq=2 ttl=64 time=0.080 ms
    64 bytes from 172.26.0.3: seq=3 ttl=64 time=0.097 ms

    --- foo ping statistics ---
    4 packets transmitted, 4 packets received, 0% packet loss
    round-trip min/avg/max = 0.070/0.081/0.097 ms

    ```

    You can no longer reach containers on the `isolated_nw` network from `container5`.
    However, you can still reach `container4` (from `container4`) using the alias
    `foo`.

    Detach from `container4` and leave it running using `CTRL-p CTRL-q`.

### Limitations of `docker network`

Although `docker network` is the recommended way to control the networks your
containers use, it does have some limitations.

#### Environment variable injection

Environment variable injection is static in nature and environment variables
cannot be changed after a container is started. The legacy `--link` flag shares
all environment variables to the linked container, but the `docker network` command
has no equivalent. When you connect a container to a network using `docker network`, no
environment variables can be dynamically shared among containers.

#### Use network-scoped aliases

Legacy links provide outgoing name resolution that is isolated within the
container in which the alias is configured. Network-scoped aliases do not allow
for this one-way isolation, but provide the alias to all members of the network.

The following example illustrates this limitation.

1.  Create another container called `container6` in the network `isolated_nw`
    and give it the network alias `app`.

    ```bash
    $ docker run --network=isolated_nw -itd --name=container6 --network-alias app busybox

    8ebe6767c1e0361f27433090060b33200aac054a68476c3be87ef4005eb1df17
    ```

2.  Attach to `container4`. Try pinging the container by name (`container6`) and by
    network alias (`app`). Notice that the IP address is the same.

    ```bash
    $ docker attach container4

    / # ping -w 4 app
    PING app (172.25.0.6): 56 data bytes
    64 bytes from 172.25.0.6: seq=0 ttl=64 time=0.070 ms
    64 bytes from 172.25.0.6: seq=1 ttl=64 time=0.080 ms
    64 bytes from 172.25.0.6: seq=2 ttl=64 time=0.080 ms
    64 bytes from 172.25.0.6: seq=3 ttl=64 time=0.097 ms

    --- app ping statistics ---
    4 packets transmitted, 4 packets received, 0% packet loss
    round-trip min/avg/max = 0.070/0.081/0.097 ms

    / # ping -w 4 container6
    PING container5 (172.25.0.6): 56 data bytes
    64 bytes from 172.25.0.6: seq=0 ttl=64 time=0.070 ms
    64 bytes from 172.25.0.6: seq=1 ttl=64 time=0.080 ms
    64 bytes from 172.25.0.6: seq=2 ttl=64 time=0.080 ms
    64 bytes from 172.25.0.6: seq=3 ttl=64 time=0.097 ms

    --- container6 ping statistics ---
    4 packets transmitted, 4 packets received, 0% packet loss
    round-trip min/avg/max = 0.070/0.081/0.097 ms
    ```

    Detach from `container4` and leave it running using `CTRL-p CTRL-q`.


3.  Connect `container6` to the `local_alias` network with the network-scoped
    alias `scoped-app`.

    ```bash
    $ docker network connect --alias scoped-app local_alias container6
    ```

    Now `container6` is aliased as `app` in network `isolated_nw`
    and as `scoped-app` in network `local_alias`.


4.  Try to reach these aliases from `container4` (which is connected to both
    these networks) and `container5` (which is connected only to `isolated_nw`).

    ```bash
    $ docker attach container4

    / # ping -w 4 scoped-app
    PING foo (172.26.0.5): 56 data bytes
    64 bytes from 172.26.0.5: seq=0 ttl=64 time=0.070 ms
    64 bytes from 172.26.0.5: seq=1 ttl=64 time=0.080 ms
    64 bytes from 172.26.0.5: seq=2 ttl=64 time=0.080 ms
    64 bytes from 172.26.0.5: seq=3 ttl=64 time=0.097 ms

    --- foo ping statistics ---
    4 packets transmitted, 4 packets received, 0% packet loss
    round-trip min/avg/max = 0.070/0.081/0.097 ms
    ```

    Detach from `container4` and leave it running using `CTRL-p CTRL-q`.

    ```bash
    $ docker attach container5

    / # ping -w 4 scoped-app
    ping: bad address 'scoped-app'

    ```

    Detach from `container5` and leave it running using `CTRL-p CTRL-q`.

This shows that an alias is scoped to the network where it is defined, and only
containers connected to that network can access the alias.


#### Resolving multiple containers to a single alias

Multiple containers can share the same network-scoped alias within the same
network. This provides a sort of DNS round-robin high availability. This may not
be reliable when using software such as Nginx, which caches clients by IP
address.

The following example illustrates how to set up and use network aliases.

> **Note**: Those using network aliases for DNS round-robin high availability
> should consider using swarm services instead. Swarm services
> provide a similar load-balancing feature out of the box. If you connect to any
> node, even a node that isn't participating in the service. Docker sends
> the request to a random node which is participating in the service and
> manages all the communication.

1.  Launch `container7` in `isolated_nw` with the same alias as `container6`,
    which is `app`.

    ```bash
    $ docker run --network=isolated_nw -itd --name=container7 --network-alias app busybox

    3138c678c123b8799f4c7cc6a0cecc595acbdfa8bf81f621834103cd4f504554
    ```

    When multiple containers share the same alias, one of those containers
    will resolve to the alias. If that container is unavailable, another
    container with the alias will be resolved. This provides a sort of high
    availability within the cluster.

    > **Note**: When the IP address is resolved, the container chosen to resolve
    > it is not completely predictable. For that reason, in the exercises below,
    > you may get different results in some steps. If the step assumes the result
    > returned is `container6` but you get `container7`, this is why.

2.  Start a continuous ping from `container4` to the `app` alias.

    ```bash
    $ docker attach container4

    $ ping app
    PING app (172.25.0.6): 56 data bytes
    64 bytes from 172.25.0.6: seq=0 ttl=64 time=0.070 ms
    64 bytes from 172.25.0.6: seq=1 ttl=64 time=0.080 ms
    64 bytes from 172.25.0.6: seq=2 ttl=64 time=0.080 ms
    64 bytes from 172.25.0.6: seq=3 ttl=64 time=0.097 ms
    ...
    ```

    The IP address that is returned belongs to `container6`.

3.  In another terminal, stop `container6`.
    ```bash
    $ docker stop container6
    ```

    In the terminal attached to `container4`, observe the `ping` output.
    It will pause when `container6` goes down, because the `ping` command
    looks up the IP when it is first invoked, and that IP is no longer reachable.
    However, the `ping` command has a very long timeout by default, so no error
    occurs.

4.  Exit the `ping` command using `CTRL+C` and run it again.

    ```bash
    $ ping app

    PING app (172.25.0.7): 56 data bytes
    64 bytes from 172.25.0.7: seq=0 ttl=64 time=0.095 ms
    64 bytes from 172.25.0.7: seq=1 ttl=64 time=0.075 ms
    64 bytes from 172.25.0.7: seq=2 ttl=64 time=0.072 ms
    64 bytes from 172.25.0.7: seq=3 ttl=64 time=0.101 ms
    ...
    ```

    The `app` alias now resolves to the IP address of `container7`.

5.  For one last test, restart `container6`.

    ```bash
    $ docker start container6
    ```

    In the terminal attached to `container4`, run the `ping` command again. It
    might now resolve to `container6` again. If you start and stop the `ping`
    several times, you will see responses from each of the containers.

    ```bash
    $ docker attach container4

    $ ping app
    PING app (172.25.0.6): 56 data bytes
    64 bytes from 172.25.0.6: seq=0 ttl=64 time=0.070 ms
    64 bytes from 172.25.0.6: seq=1 ttl=64 time=0.080 ms
    64 bytes from 172.25.0.6: seq=2 ttl=64 time=0.080 ms
    64 bytes from 172.25.0.6: seq=3 ttl=64 time=0.097 ms
    ...
    ```

    Stop the ping with `CTRL+C`. Detach from `container4` and leave it running
    using `CTRL-p CTRL-q`.

## Disconnecting containers

You can disconnect a container from a network at any time using the `docker network
disconnect` command.

1.  Disconnect `container2` from the `isolated_nw` network, then inspect `container2`
    and the `isolated_nw` network.

    ```bash
    $ docker network disconnect isolated_nw container2

    $ docker inspect --format='{{json .NetworkSettings.Networks}}'  container2 | python -m json.tool

    {
        "bridge": {
            "NetworkID":"7ea29fc1412292a2d7bba362f9253545fecdfa8ce9a6e37dd10ba8bee7129812",
            "EndpointID": "9e4575f7f61c0f9d69317b7a4b92eefc133347836dd83ef65deffa16b9985dc0",
            "Gateway": "172.17.0.1",
            "GlobalIPv6Address": "",
            "GlobalIPv6PrefixLen": 0,
            "IPAddress": "172.17.0.3",
            "IPPrefixLen": 16,
            "IPv6Gateway": "",
            "MacAddress": "02:42:ac:11:00:03"
        }
    }


    $ docker network inspect isolated_nw

    [
        {
            "Name": "isolated_nw",
            "Id": "06a62f1c73c4e3107c0f555b7a5f163309827bfbbf999840166065a8f35455a8",
            "Scope": "local",
            "Driver": "bridge",
            "IPAM": {
                "Driver": "default",
                "Config": [
                    {
                        "Subnet": "172.21.0.0/16",
                        "Gateway": "172.21.0.1/16"
                    }
                ]
            },
            "Containers": {
                "467a7863c3f0277ef8e661b38427737f28099b61fa55622d6c30fb288d88c551": {
                    "Name": "container3",
                    "EndpointID": "dffc7ec2915af58cc827d995e6ebdc897342be0420123277103c40ae35579103",
                    "MacAddress": "02:42:ac:19:03:03",
                    "IPv4Address": "172.25.3.3/16",
                    "IPv6Address": ""
                }
            },
            "Options": {}
        }
    ]
    ```


2.  When a container is disconnected from a network, it can no longer communicate
    with other containers connected to that network, unless it has other networks
    in common with them. Verify that `container2` can no longer reach `container3`,
    which is on the `isolated_nw` network.

    ```bash
    $ docker attach container2

    / # ifconfig
    eth0      Link encap:Ethernet  HWaddr 02:42:AC:11:00:03  
              inet addr:172.17.0.3  Bcast:0.0.0.0  Mask:255.255.0.0
              inet6 addr: fe80::42:acff:fe11:3/64 Scope:Link
              UP BROADCAST RUNNING MULTICAST  MTU:9001  Metric:1
              RX packets:8 errors:0 dropped:0 overruns:0 frame:0
              TX packets:8 errors:0 dropped:0 overruns:0 carrier:0
              collisions:0 txqueuelen:0
              RX bytes:648 (648.0 B)  TX bytes:648 (648.0 B)

    lo        Link encap:Local Loopback  
              inet addr:127.0.0.1  Mask:255.0.0.0
              inet6 addr: ::1/128 Scope:Host
              UP LOOPBACK RUNNING  MTU:65536  Metric:1
              RX packets:0 errors:0 dropped:0 overruns:0 frame:0
              TX packets:0 errors:0 dropped:0 overruns:0 carrier:0
              collisions:0 txqueuelen:0
              RX bytes:0 (0.0 B)  TX bytes:0 (0.0 B)

    / # ping container3
    PING container3 (172.25.3.3): 56 data bytes
    ^C
    --- container3 ping statistics ---
    2 packets transmitted, 0 packets received, 100% packet loss
    ```

3.  Verify that `container2` still has full connectivity to the default `bridge`
    network.

    ```bash
    / # ping container1
    PING container1 (172.17.0.2): 56 data bytes
    64 bytes from 172.17.0.2: seq=0 ttl=64 time=0.119 ms
    64 bytes from 172.17.0.2: seq=1 ttl=64 time=0.174 ms
    ^C
    --- container1 ping statistics ---
    2 packets transmitted, 2 packets received, 0% packet loss
    round-trip min/avg/max = 0.119/0.146/0.174 ms
    / #
    ```

4.  Remove `container4`, `container5`, `container6`, and `container7`.

    ```bash
    $ docker stop container4 container5 container6 container7

    $ docker rm container4 container5 container6 container7
    ```

### Handling stale network endpoints

In some scenarios, such as ungraceful docker daemon restarts in a
multi-host network, the daemon cannot clean up stale connectivity endpoints.
Such stale endpoints may cause an error if a new container is connected
to that network with the same name as the stale endpoint:

```none
ERROR: Cannot start container bc0b19c089978f7845633027aa3435624ca3d12dd4f4f764b61eac4c0610f32e: container already connected to network multihost
```

To clean up these stale endpoints, remove the container and disconnect it
from the network forcibly (`docker network disconnect -f`). Now you can
successfully connect the container to the network.

```bash
$ docker run -d --name redis_db --network multihost redis

ERROR: Cannot start container bc0b19c089978f7845633027aa3435624ca3d12dd4f4f764b61eac4c0610f32e: container already connected to network multihost

$ docker rm -f redis_db

$ docker network disconnect -f multihost redis_db

$ docker run -d --name redis_db --network multihost redis

7d986da974aeea5e9f7aca7e510bdb216d58682faa83a9040c2f2adc0544795a
```

## Remove a network

When all the containers in a network are stopped or disconnected, you can
remove a network. If a network has connected endpoints, an error occurs.

1.  Disconnect `container3` from `isolated_nw`.

    ```bash
    $ docker network disconnect isolated_nw container3
      ```

2.  Inspect `isolated_nw` to verify that no other endpoints are connected to it.

    ```bash
    $ docker network inspect isolated_nw

    [
        {
            "Name": "isolated_nw",
            "Id": "06a62f1c73c4e3107c0f555b7a5f163309827bfbbf999840166065a8f35455a8",
            "Scope": "local",
            "Driver": "bridge",
            "IPAM": {
                "Driver": "default",
                "Config": [
                    {
                        "Subnet": "172.21.0.0/16",
                        "Gateway": "172.21.0.1/16"
                    }
                ]
            },
            "Containers": {},
            "Options": {}
        }
    ]
    ```

3.  Remove the `isolated_nw` network.

    ```bash
    $ docker network rm isolated_nw
    ```

4.  List all your networks to verify that `isolated_nw` no longer exists:

    ```bash
    $ docker network ls

    NETWORK ID          NAME                DRIVER              SCOPE
    4bb8c9bf4292        bridge              bridge              local
    43575911a2bd        host                host                local
    76b7dc932e03        local_alias         bridge              local
    b1a086897963        my-network          bridge              local
    3eb020e70bfd        none                null                local
    69568e6336d8        simple-network      bridge              local
    ```

## Related information

* [network create](../../reference/commandline/network_create.md)
* [network inspect](../../reference/commandline/network_inspect.md)
* [network connect](../../reference/commandline/network_connect.md)
* [network disconnect](../../reference/commandline/network_disconnect.md)
* [network ls](../../reference/commandline/network_ls.md)
* [network rm](../../reference/commandline/network_rm.md)
