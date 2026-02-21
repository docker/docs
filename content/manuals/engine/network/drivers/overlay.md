---
title: Overlay network driver
description: All about using overlay networks
keywords: network, overlay, user-defined, swarm, service
aliases:
  - /config/containers/overlay/
  - /engine/userguide/networking/overlay-security-model/
  - /network/overlay/
  - /network/drivers/overlay/
  - /engine/network/tutorials/overlay/
  - /engine/userguide/networking/get-started-overlay/
---

The `overlay` network driver creates a distributed network among multiple
Docker daemon hosts. This network sits on top of (overlays) the host-specific
networks, allowing containers connected to it to communicate securely when
encryption is enabled. Docker transparently handles routing of each packet to
and from the correct Docker daemon host and the correct destination container.

You can create user-defined `overlay` networks using `docker network create`,
in the same way that you can create user-defined `bridge` networks. Services
or containers can be connected to more than one network at a time. Services or
containers can only communicate across networks they're each connected to.

Overlay networks are often used to create a connection between Swarm services,
but you can also use it to connect standalone containers running on different
hosts. When using standalone containers, it's still required that you use
Swarm mode to establish a connection between the hosts.

This page describes overlay networks in general, and when used with standalone
containers. For information about overlay for Swarm services, see
[Manage Swarm service networks](/manuals/engine/swarm/networking.md).

## Requirements

Docker hosts must be part of a swarm to use overlay networks, even when
connecting standalone containers. The following ports must be open between
participating hosts:

- `2377/tcp`: Swarm control plane (configurable)
- `4789/udp`: Overlay traffic (configurable)
- `7946/tcp` and `7946/udp`: Node communication (not configurable)

## Create an overlay network

The following table lists the ports that need to be open to each host
participating in an overlay network:

| Ports                  | Description                                                                                                                                             |
| :--------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `2377/tcp`             | The default Swarm control plane port, is configurable with [`docker swarm join --listen-addr`](/reference/cli/docker/swarm/join/#--listen-addr-value) |
| `4789/udp`             | The default overlay traffic port, configurable with [`docker swarm init --data-path-addr`](/reference/cli/docker/swarm/init/#data-path-port)          |
| `7946/tcp`, `7946/udp` | Used for communication among nodes, not configurable                                                                                                    |

To create an overlay network that containers on other Docker hosts can connect to,
run the following command:

```console
$ docker network create -d overlay --attachable my-attachable-overlay
```

The `--attachable` option enables both standalone containers
and Swarm services to connect to the overlay network.
Without `--attachable`, only Swarm services can connect to the network.

You can specify the IP address range, subnet, gateway, and other options. See
`docker network create --help` for details.

## Encrypt traffic on an overlay network

Use the `--opt encrypted` flag to encrypt the application data
transmitted over the overlay network:

```console
$ docker network create \
  --opt encrypted \
  --driver overlay \
  --attachable \
  my-attachable-multi-host-network
```

This enables IPsec encryption at the level of the Virtual Extensible LAN (VXLAN).
This encryption imposes a non-negligible performance penalty,
so you should test this option before using it in production.

> [!WARNING]
>
> Don't attach Windows containers to encrypted overlay networks.
>
> Overlay network encryption isn't supported on Windows.
> Swarm doesn't report an error when a Windows host
> attempts to connect to an encrypted overlay network,
> but networking for the Windows containers is affected as follows:
>
> - Windows containers can't communicate with Linux containers on the network
> - Data traffic between Windows containers on the network isn't encrypted

## Attach a container to an overlay network

Adding containers to an overlay network gives them the ability to communicate
with other containers without having to set up routing on the individual Docker
daemon hosts. A prerequisite for doing this is that the hosts have joined the same Swarm.

To join an overlay network named `multi-host-network` with a `busybox` container:

```console
$ docker run --network multi-host-network busybox sh
```

> [!NOTE]
>
> This only works if the overlay network is attachable
> (created with the `--attachable` flag).

## Container discovery

Publishing ports of a container on an overlay network opens the ports to other
containers on the same network. Containers are discoverable by doing a DNS lookup
using the container name.

| Flag value                      | Description                                                                                                                                                 |
| ------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `-p 8080:80`                    | Map TCP port 80 in the container to port `8080` on the overlay network.                                                                                     |
| `-p 8080:80/udp`                | Map UDP port 80 in the container to port `8080` on the overlay network.                                                                                     |
| `-p 8080:80/sctp`               | Map SCTP port 80 in the container to port `8080` on the overlay network.                                                                                    |
| `-p 8080:80/tcp -p 8080:80/udp` | Map TCP port 80 in the container to TCP port `8080` on the overlay network, and map UDP port 80 in the container to UDP port `8080` on the overlay network. |

## Connection limit for overlay networks

Due to limitations set by the Linux kernel, overlay networks become unstable and
inter-container communications may break when 1000 containers are co-located on
the same host.

For more information about this limitation, see
[moby/moby#44973](https://github.com/moby/moby/issues/44973#issuecomment-1543747718).

## Usage examples

This section provides hands-on examples for working with overlay networks. These
examples cover swarm services and standalone containers on multiple Docker hosts.

### Prerequisites

All examples require at least a single-node swarm. Initialize one by running
`docker swarm init` on the host. You can run these examples on multi-node
swarms as well.

### Use the default overlay network

This example shows how the default overlay network works with swarm services.
You'll create an `nginx` service and examine the network from the service
containers' perspective.

#### Prerequisites for multi-node setup

This walkthrough requires three Docker hosts that can communicate with each
other on the same network with no firewall blocking traffic between them:

- `manager`: Functions as both manager and worker
- `worker-1`: Functions as worker only
- `worker-2`: Functions as worker only

If you don't have three hosts available, you can set up three virtual machines
on a cloud provider with Docker installed.

#### Create the swarm

1. On `manager`, initialize the swarm. If the host has one network interface,
   the `--advertise-addr` flag is optional:

   ```console
   $ docker swarm init --advertise-addr=<IP-ADDRESS-OF-MANAGER>
   ```

   Save the join token displayed for use with workers.

2. On `worker-1`, join the swarm:

   ```console
   $ docker swarm join --token <TOKEN> \
     --advertise-addr <IP-ADDRESS-OF-WORKER-1> \
     <IP-ADDRESS-OF-MANAGER>:2377
   ```

3. On `worker-2`, join the swarm:

   ```console
   $ docker swarm join --token <TOKEN> \
     --advertise-addr <IP-ADDRESS-OF-WORKER-2> \
     <IP-ADDRESS-OF-MANAGER>:2377
   ```

4. On `manager`, list all nodes:

   ```console
   $ docker node ls

   ID                            HOSTNAME            STATUS              AVAILABILITY        MANAGER STATUS
   d68ace5iraw6whp7llvgjpu48 *   ip-172-31-34-146    Ready               Active              Leader
   nvp5rwavvb8lhdggo8fcf7plg     ip-172-31-35-151    Ready               Active
   ouvx2l7qfcxisoyms8mtkgahw     ip-172-31-36-89     Ready               Active
   ```

   Filter by role if needed:

   ```console
   $ docker node ls --filter role=manager
   $ docker node ls --filter role=worker
   ```

5. List Docker networks on all hosts. Each now has an overlay network called
   `ingress` and a bridge network called `docker_gwbridge`:

   ```console
   $ docker network ls

   NETWORK ID          NAME                DRIVER              SCOPE
   495c570066be        bridge              bridge              local
   961c6cae9945        docker_gwbridge     bridge              local
   ff35ceda3643        host                host                local
   trtnl4tqnc3n        ingress             overlay             swarm
   c8357deec9cb        none                null                local
   ```

The `docker_gwbridge` connects the `ingress` network to the Docker host's
network interface. If you create services without specifying a network, they
connect to `ingress`. It's recommended to use separate overlay networks for each
application or group of related applications.

#### Create the services

1. On `manager`, create a new overlay network:

   ```console
   $ docker network create -d overlay nginx-net
   ```

   The overlay network is automatically created on worker nodes when they run
   service tasks that need it.

2. On `manager`, create a 5-replica Nginx service connected to `nginx-net`:

   > [!NOTE]
   > Services can only be created on a manager.

   ```console
   $ docker service create \
     --name my-nginx \
     --publish target=80,published=80 \
     --replicas=5 \
     --network nginx-net \
     nginx
   ```

   The default `ingress` publish mode means you can browse to port 80 on any
   node and connect to one of the 5 service tasks, even if no tasks run on that
   node.

3. Monitor service creation progress:

   ```console
   $ docker service ls
   ```

4. Inspect the `nginx-net` network on all hosts. The `Containers` section lists
   all service tasks connected to the overlay network from that host.

5. From `manager`, inspect the service:

   ```console
   $ docker service inspect my-nginx
   ```

   Notice the information about ports and endpoints.

6. Create a second network and update the service to use it:

   ```console
   $ docker network create -d overlay nginx-net-2
   $ docker service update \
     --network-add nginx-net-2 \
     --network-rm nginx-net \
     my-nginx
   ```

7. Verify the update completed:

   ```console
   $ docker service ls
   ```

   Inspect both networks to verify containers moved from `nginx-net` to
   `nginx-net-2`.

   > [!NOTE]
   > Overlay networks are automatically created on swarm worker nodes as needed,
   > but aren't automatically removed.

8. Clean up:

   ```console
   $ docker service rm my-nginx
   $ docker network rm nginx-net nginx-net-2
   ```

### Use a user-defined overlay network

This example shows the recommended approach for production services using custom
overlay networks.

#### Prerequisites

This assumes the swarm is already set up and you're on a manager node.

#### Steps

1. Create a user-defined overlay network:

   ```console
   $ docker network create -d overlay my-overlay
   ```

2. Start a service using the overlay network, publishing port 80 to port 8080:

   ```console
   $ docker service create \
     --name my-nginx \
     --network my-overlay \
     --replicas 1 \
     --publish published=8080,target=80 \
     nginx:latest
   ```

3. Verify the service task is connected to the network:

   ```console
   $ docker network inspect my-overlay
   ```

   Check the `Containers` section for the `my-nginx` service task.

4. Clean up:

   ```console
   $ docker service rm my-nginx
   $ docker network rm my-overlay
   ```

### Use an overlay network for standalone containers

This example demonstrates DNS container discovery between standalone containers
on different Docker hosts using an overlay network.

#### Prerequisites

You need two Docker hosts that can communicate with each other with the
following ports open between them:

- TCP port 2377
- TCP and UDP port 7946
- UDP port 4789

This example refers to the hosts as `host1` and `host2`.

#### Steps

1. Set up the swarm:

   On `host1`, initialize a swarm:

   ```console
   $ docker swarm init
   Swarm initialized: current node (vz1mm9am11qcmo979tlrlox42) is now a manager.

   To add a worker to this swarm, run the following command:

       docker swarm join --token SWMTKN-1-5g90q48weqrtqryq4kj6ow0e8xm9wmv9o6vgqc5j320ymybd5c-8ex8j0bc40s6hgvy5ui5gl4gy 172.31.47.252:2377
   ```

   On `host2`, join the swarm using the token from the previous output:

   ```console
   $ docker swarm join --token <your_token> <your_ip_address>:2377
   This node joined a swarm as a worker.
   ```

   If the join fails, run `docker swarm leave --force` on `host2`, verify
   network and firewall settings, and try again.

2. On `host1`, create an attachable overlay network:

   ```console
   $ docker network create --driver=overlay --attachable test-net
   uqsof8phj3ak0rq9k86zta6ht
   ```

   Note the returned network ID.

3. On `host1`, start an interactive container that connects to `test-net`:

   ```console
   $ docker run -it --name alpine1 --network test-net alpine
   / #
   ```

4. On `host2`, list available networks. Notice that `test-net` doesn't exist yet:

   ```console
   $ docker network ls
   NETWORK ID          NAME                DRIVER              SCOPE
   ec299350b504        bridge              bridge              local
   66e77d0d0e9a        docker_gwbridge     bridge              local
   9f6ae26ccb82        host                host                local
   omvdxqrda80z        ingress             overlay             swarm
   b65c952a4b2b        none                null                local
   ```

5. On `host2`, start a detached, interactive container that connects to
   `test-net`:

   ```console
   $ docker run -dit --name alpine2 --network test-net alpine
   fb635f5ece59563e7b8b99556f816d24e6949a5f6a5b1fbd92ca244db17a4342
   ```

   > [!NOTE]
   > Automatic DNS container discovery only works with unique container names.

6. On `host2`, verify that `test-net` was created with the same network ID as on
   `host1`:

   ```console
   $ docker network ls
   NETWORK ID          NAME                DRIVER              SCOPE
   ...
   uqsof8phj3ak        test-net            overlay             swarm
   ```

7. On `host1`, ping `alpine2` from within `alpine1`:

   ```console
   / # ping -c 2 alpine2
   PING alpine2 (10.0.0.5): 56 data bytes
   64 bytes from 10.0.0.5: seq=0 ttl=64 time=0.600 ms
   64 bytes from 10.0.0.5: seq=1 ttl=64 time=0.555 ms

   --- alpine2 ping statistics ---
   2 packets transmitted, 2 packets received, 0% packet loss
   round-trip min/avg/max = 0.555/0.577/0.600 ms
   ```

   The two containers communicate over the overlay network connecting the two
   hosts. You can also run another container on `host2` and ping `alpine1`:

   ```console
   $ docker run -it --rm --name alpine3 --network test-net alpine
   / # ping -c 2 alpine1
   / # exit
   ```

8. On `host1`, close the `alpine1` session (which stops the container):

   ```console
   / # exit
   ```

9. Clean up. You must stop and remove containers on each host independently:

   On `host2`:

   ```console
   $ docker container stop alpine2
   $ docker network ls
   $ docker container rm alpine2
   ```

   When you stop `alpine2`, `test-net` disappears from `host2`.

   On `host1`:

   ```console
   $ docker container rm alpine1
   $ docker network rm test-net
   ```

## Next steps

- Learn about [networking from the container's point of view](../_index.md)
- Learn about [standalone bridge networks](bridge.md)
- Learn about [Macvlan networks](macvlan.md)
