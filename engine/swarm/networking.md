---
description: Use swarm mode overlay networking features
keywords: swarm, networking, ingress, overlay, service discovery
title: Manage swarm service networks
---

A Docker swarm generates two different kinds of traffic:

- **Control and management plane traffic**: This includes swarm management
  messages, such as requests to join or leave the swarm. This traffic is
  always encrypted.

- **Application data plane traffic**: This includes container traffic and
  traffic to and from external clients.

This topic discusses how to manage the application data for your swarm services.
For more details about swarm networking in general, see the
[Docker networking reference architecture](https://success.docker.com/Architecture/Docker_Reference_Architecture%3A_Designing_Scalable%2C_Portable_Docker_Container_Networks){: target="_blank" class="_" }.

The following three network concepts are important to swarm services:

- **Overlay networks** manage communications among the Docker daemons
  participating in the swarm. You can create overlay networks, in the same way
  as user-defined networks for standalone containers. You can attach a service
  to one or more existing overlay networks as well, to enable service-to-service
  communication. Overlay networks are Docker networks that use the `overlay`
  network driver.

- The **ingress network** is a special overlay network that facilitates
  load balancing among a service's nodes. When any swarm node receives a
  request on a published port, it hands that request off to a module called
  `IPVS`. `IPVS` keeps track of all the IP addresses participating in that
  service, selects one of them, and routes the request to it, over the
  `ingress` network.

  The `ingress` network is created automatically when you initialize or join a
  swarm. Most users do not need to customize its configuration, but Docker 17.05
  and higher allows you to do so.

- The **docker_gwbridge** is a bridge network that connects the overlay
  networks (including the `ingress` network) to an individual Docker daemon's
  physical network. By default, each container a service is running is connected
  to its local Docker daemon host's `docker_gwbridge` network.

  The `docker_gwbridge` network is created automatically when you initialize or
  join a swarm. Most users do not need to customize its configuration, but
  Docker allows you to do so.


## Firewall considerations

Docker daemons participating in a swarm need the ability to communicate with
each other over the following ports:

* Port `7946` TCP/UDP for container network discovery.
* Port `4789` UDP for the container overlay network.

## Create an overlay network

To create an overlay network, specify the `overlay` driver when using the
`docker network create` command:

```bash
$ docker network create \
  --driver overlay \
  my-network
```

The above command doesn't specify any custom options, so Docker assigns a
subnet and uses default options. You can see information about the network using
`docker network inspect`.

When no containers are connected to the overlay network, its configuration is
not very exciting:

```bash
$ docker network inspect my-network
[
    {
        "Name": "my-network",
        "Id": "fsf1dmx3i9q75an49z36jycxd",
        "Created": "0001-01-01T00:00:00Z",
        "Scope": "swarm",
        "Driver": "overlay",
        "EnableIPv6": false,
        "IPAM": {
            "Driver": "default",
            "Options": null,
            "Config": []
        },
        "Internal": false,
        "Attachable": false,
        "Ingress": false,
        "Containers": null,
        "Options": {
            "com.docker.network.driver.overlay.vxlanid_list": "4097"
        },
        "Labels": null
    }
]
```

In the above output, notice that the driver is `overlay` and that the scope is
`swarm`, rather than `local`, `host`, or `global` scopes you might see in
other types of Docker networks. This scope indicates that only hosts which are
participating in the swarm can access this network.

The network's subnet and gateway are dynamically configured when a service
connects to the network for the first time. The following example shows
the same network as above, but with three containers of a `redis` service
connected to it.

```bash
$ docker network inspect my-network
[
    {
        "Name": "my-network",
        "Id": "fsf1dmx3i9q75an49z36jycxd",
        "Created": "2017-05-31T18:35:58.877628262Z",
        "Scope": "swarm",
        "Driver": "overlay",
        "EnableIPv6": false,
        "IPAM": {
            "Driver": "default",
            "Options": null,
            "Config": [
                {
                    "Subnet": "10.0.0.0/24",
                    "Gateway": "10.0.0.1"
                }
            ]
        },
        "Internal": false,
        "Attachable": false,
        "Ingress": false,
        "Containers": {
            "0e08442918814c2275c31321f877a47569ba3447498db10e25d234e47773756d": {
                "Name": "my-redis.1.ka6oo5cfmxbe6mq8qat2djgyj",
                "EndpointID": "950ce63a3ace13fe7ef40724afbdb297a50642b6d47f83a5ca8636d44039e1dd",
                "MacAddress": "02:42:0a:00:00:03",
                "IPv4Address": "10.0.0.3/24",
                "IPv6Address": ""
            },
            "88d55505c2a02632c1e0e42930bcde7e2fa6e3cce074507908dc4b827016b833": {
                "Name": "my-redis.2.s7vlybipal9xlmjfqnt6qwz5e",
                "EndpointID": "dd822cb68bcd4ae172e29c321ced70b731b9994eee5a4ad1d807d9ae80ecc365",
                "MacAddress": "02:42:0a:00:00:05",
                "IPv4Address": "10.0.0.5/24",
                "IPv6Address": ""
            },
            "9ed165407384f1276e5cfb0e065e7914adbf2658794fd861cfb9b991eddca754": {
                "Name": "my-redis.3.hbz3uk3hi5gb61xhxol27hl7d",
                "EndpointID": "f62c686a34c9f4d70a47b869576c37dffe5200732e1dd6609b488581634cf5d2",
                "MacAddress": "02:42:0a:00:00:04",
                "IPv4Address": "10.0.0.4/24",
                "IPv6Address": ""
            }
        },
        "Options": {
            "com.docker.network.driver.overlay.vxlanid_list": "4097"
        },
        "Labels": {},
        "Peers": [
            {
                "Name": "moby-e57c567e25e2",
                "IP": "192.168.65.2"
            }
        ]
    }
]
```

### Customize an overlay network

There may be situations where you don't want to use the default configuration
for an overlay network. For a full list of configurable options, run the
command `docker network create --help`. The following are some of the most
common options to change.

#### Configure the subnet and gateway

By default, the network's subnet and gateway are configured automatically when
the first service is connected to the network. You can configure these when
creating a network using the `--subnet` and `--gateway` flags. The following
example extends the previous one by configuring the subnet and gateway.

```bash
$ docker network create \
  --driver overlay \
  --subnet 10.0.9.0/24 \
  --gateway 10.0.9.99 \
  my-network
```

##### Using custom default address pools

To customize subnet allocation for your Swarm networks, you can [optionally configure them](swarm-mode.md) during `swarm init`.

For example, the following command is used when initializing Swarm:

```bash
$ docker swarm init --default-addr-pool 10.20.0.0/16 --default-addr-pool-mask-length 26`
```

Whenever a user creates a network, but does not use the `--subnet` command line option, the subnet for this network will be allocated sequentially from the next available subnet from the pool. If the specified network is already allocated, that network will not be used for Swarm. 

Multiple pools can be configured if discontiguous address space is required. However, allocation from specific pools is not supported. Network subnets will be allocated sequentially from the IP pool space and subnets will be reused as they are deallocated from networks that are deleted.

The default mask length can be configured and is the same for all networks. It is set to `/24` by default. To change the default subnet mask length, use the `--default-addr-pool-mask-length` command line option.

> **Note**: Default address pools can only be configured on `swarm init` and cannot be altered after cluster creation.

##### Overlay network size limitations

Docker recommends creating overlay networks with `/24` blocks. The `/24` overlay network blocks, which limits the network to 256 IP addresses. 

This recommendation addresses [limitations with swarm mode](https://github.com/moby/moby/issues/30820). 
If you need more than 256 IP addresses, do not increase the IP block size. You can either use `dnsrr` 
endpoint mode with an external load balancer, or use multiple smaller overlay networks. See 
[Configure service discovery](#configure-service-discovery) or more information about different endpoint modes.

#### Configure encryption of application data

Management and control plane data related to a swarm is always encrypted.
For more details about the encryption mechanisms, see the
[Docker swarm mode overlay network security model](../../network/overlay.md).

Application data among swarm nodes is not encrypted by default. To encrypt this
traffic on a given overlay network, use the `--opt encrypted` flag on `docker
network create`. This enables IPSEC encryption at the level of the vxlan. This
encryption imposes a non-negligible performance penalty, so you should test this
option before using it in production.


## Attach a service to an overlay network

To attach a service to an existing overlay network, pass the `--network` flag to
`docker service create`, or the `--network-add` flag to `docker service update`.

```bash
$ docker service create \
  --replicas 3 \
  --name my-web \
  --network my-network \
  nginx
```

Service containers connected to an overlay network can communicate with
each other across it.

To see which networks a service is connected to, use `docker service ls` to find
the name of the service, then `docker service ps <service-name>` to list the
networks. Alternately, to see which services' containers are connected to a
network, use `docker network inspect <network-name>`. You can run these commands
from any swarm node which is joined to the swarm and is in a `running` state.

### Configure service discovery

**Service discovery** is the mechanism Docker uses to route a request from your
service's external clients to an individual swarm node, without the client
needing to know how many nodes are participating in the service or their
IP addresses or ports. You don't need to publish ports which are used between
services on the same network. For instance, if you have a
[WordPress service that stores its data in a MySQL service](http://training.play-with-docker.com/swarm-service-discovery/),
and they are connected to the same overlay network, you do not need to publish
the MySQL port to the client, only the WordPress HTTP port.

Service discovery can work in two different ways: internal connection-based
load-balancing at Layers 3 and 4 using the embedded DNS and a virtual IP (VIP),
or external and customized request-based load-balancing at Layer 7 using DNS
round robin (DNSRR). You can configure this per service.

- By default, when you attach a service to a network and that service publishes
  one or more ports, Docker assigns the service a virtual IP (VIP), which is the
  "front end" for clients to reach the service. Docker keeps a list of all
  worker nodes in the service, and routes requests between the client and one of
  the nodes. Each request from the client might be routed to a different node.

- If you configure a service to use DNS round-robin (DNSRR) service discovery,
  there is not a single virtual IP. Instead, Docker sets up DNS entries for the
  service such that a DNS query for the service name returns a list of IP
  addresses, and the client connects directly to one of these.

  DNS round-robin is useful in cases where you want to use your own load
  balancer, such as HAProxy. To configure a service to use DNSRR, use the flag
  `--endpoint-mode dnsrr` when creating a new service or updating an existing
  one.

## Customize the ingress network

Most users never need to configure the `ingress` network, but Docker 17.05 and
higher allow you to do so. This can be useful if the automatically-chosen subnet
conflicts with one that already exists on your network, or you need to customize
other low-level network settings such as the MTU.

Customizing the `ingress` network involves removing and recreating it. This is
usually done before you create any services in the swarm. If you have existing
services which publish ports, those services need to be removed before you can
remove the `ingress` network.

During the time that no `ingress` network exists, existing services which do not
publish ports continue to function but are not load-balanced. This affects
services which publish ports, such as a WordPress service which publishes port
80.

1.  Inspect the `ingress` network using `docker network inspect ingress`, and
    remove any services whose containers are connected to it. These are services
    that publish ports, such as a WordPress service which publishes port 80. If
    all such services are not stopped, the next step fails.

2.  Remove the existing `ingress` network:

    ```bash
    $ docker network rm ingress

    WARNING! Before removing the routing-mesh network, make sure all the nodes
    in your swarm run the same docker engine version. Otherwise, removal may not
    be effective and functionality of newly created ingress networks will be
    impaired.
    Are you sure you want to continue? [y/N]
    ```

3.  Create a new overlay network using the `--ingress` flag, along  with the
    custom options you want to set. This example sets the MTU to 1200, sets
    the subnet to `10.11.0.0/16`, and sets the gateway to `10.11.0.2`.

    ```bash
    $ docker network create \
      --driver overlay \
      --ingress \
      --subnet=10.11.0.0/16 \
      --gateway=10.11.0.2 \
      --opt com.docker.network.mtu=1200 \
      my-ingress
    ```

    > **Note**: You can name your `ingress` network something other than
    > `ingress`, but you can only have one. An attempt to create a second one
    > fails.

4.  Restart the services that you stopped in the first step.

## Customize the docker_gwbridge

The `docker_gwbridge` is a virtual bridge that connects the overlay networks
(including the `ingress` network) to an individual Docker daemon's physical
network. Docker creates it automatically when you initialize a swarm or join a
Docker host to a swarm, but it is not a Docker device. It exists in the kernel
of the Docker host. If you need to customize its settings, you must do so before
joining the Docker host to the swarm, or after temporarily removing the host
from the swarm.

You need to have the `brctl` application installed on your operating system in
order to delete an existing bridge. The package name is `bridge-utils`.

1.  Stop Docker.

2.  Use the `brctl show docker_gwbridge` command to check whether a bridge
    device exists called `docker_gwbridge`. If so, remove it using
    `brctl delbr docker_gwbridge`.

3.  Start Docker. Do not join or initialize the swarm.

4.  Create or re-create the `docker_gwbridge` bridge with your custom settings.
    This example uses the subnet `10.11.0.0/16`. For a full list of customizable
    options, see [Bridge driver options](../reference/commandline/network_create.md#bridge-driver-options).

    ```bash
    $ docker network create \
    --subnet 10.11.0.0/16 \
    --opt com.docker.network.bridge.name=docker_gwbridge \
    --opt com.docker.network.bridge.enable_icc=false \
    --opt com.docker.network.bridge.enable_ip_masquerade=true \
    docker_gwbridge
    ```

5.  Initialize or join the swarm.

## Use a separate interface for control and data traffic

By default, all swarm traffic is sent over the same interface, including control
and management traffic for maintaining the swarm itself and data traffic to and
from the service containers.

In Docker 17.06 and higher, it is possible to separate this traffic by passing
the `--data-path-addr` flag when initializing or joining the swarm. If there are
multiple interfaces, `--advertise-addr` must be specified explicitly, and
`--data-path-addr` defaults to `--advertise-addr` if not specified. Traffic about
joining, leaving, and managing the swarm is sent over the
`--advertise-addr` interface, and traffic among a service's containers is sent
sent over the `--data-path-addr` interface. These flags can take an IP address or
a network device name, such as `eth0`.

This example initializes a swarm with a separate `--data-path-addr`. It assumes
that your Docker host has two different network interfaces: 10.0.0.1 should be
used for control and management traffic and 192.168.0.1 should be used for
traffic relating to services.

```bash
$ docker swarm init --advertise-addr 10.0.0.1 --data-path-addr 192.168.0.1
```

This example joins the swarm managed by host `192.168.99.100:2377` and sets the
`--advertise-addr` flag to `eth0` and the `--data-path-addr` flag to `eth1`.

```bash
$ docker swarm join \
  --token SWMTKN-1-49nj1cmql0jkz5s954yi3oex3nedyz0fb0xx14ie39trti4wxv-8vxv8rssmk743ojnwacrr2d7c \
  --advertise-addr eth0 \
  --data-path-addr eth1 \
  192.168.99.100:2377
```

## Learn more

* [Deploy services to a swarm](services.md)
* [Swarm administration guide](admin_guide.md)
* [Docker CLI reference](../reference/commandline/docker.md)
* [Swarm mode tutorial](swarm-tutorial/index.md)
* [Docker networking reference architecture](https://success.docker.com/Architecture/Docker_Reference_Architecture%3A_Designing_Scalable%2C_Portable_Docker_Container_Networks){: target="_blank" class="_" }
