---
title: Networking with overlay networks
description: Tutorials for networking with swarm services and standalone containers
  on multiple Docker daemons
keywords: networking, bridge, routing, ports, swarm, overlay
aliases:
- /engine/userguide/networking/get-started-overlay/
- /network/network-tutorial-overlay/
---

This series of tutorials deals with networking for swarm services.
For networking with standalone containers, see
[Networking with standalone containers](/manuals/engine/network/tutorials/standalone.md). If you need to
learn more about Docker networking in general, see the [overview](/manuals/engine/network/_index.md).

This page includes the following tutorials. You can run each of them on
Linux, Windows, or a Mac, but for the last one, you need a second Docker
host running elsewhere.

- [Use the default overlay network](#use-the-default-overlay-network) demonstrates
  how to use the default overlay network that Docker sets up for you
  automatically when you initialize or join a swarm. This network is not the
  best choice for production systems.

- [Use user-defined overlay networks](#use-a-user-defined-overlay-network) shows
  how to create and use your own custom overlay networks, to connect services.
  This is recommended for services running in production.

- [Use an overlay network for standalone containers](#use-an-overlay-network-for-standalone-containers)
  shows how to communicate between standalone containers on different Docker
  daemons using an overlay network.

## Prerequisites

These require you to have at least a single-node swarm, which means that
you have started Docker and run `docker swarm init` on the host. You can run
the examples on a multi-node swarm as well.

## Use the default overlay network

In this example, you start an `alpine` service and examine the characteristics
of the network from the point of view of the individual service containers.

This tutorial does not go into operation-system-specific details about how
overlay networks are implemented, but focuses on how the overlay functions from
the point of view of a service.

### Prerequisites

This tutorial requires three physical or virtual Docker hosts which can all
communicate with one another. This tutorial assumes that the three hosts are
running on the same network with no firewall involved.

These hosts will be referred to as `manager`, `worker-1`, and `worker-2`. The
`manager` host will function as both a manager and a worker, which means it can
both run service tasks and manage the swarm. `worker-1` and `worker-2` will
function as workers only,

If you don't have three hosts handy, an easy solution is to set up three
Ubuntu hosts on a cloud provider such as Amazon EC2, all on the same network
with all communications allowed to all hosts on that network (using a mechanism
such as EC2 security groups), and then to follow the
[installation instructions for Docker Engine - Community on Ubuntu](/manuals/engine/install/ubuntu.md).

### Walkthrough

#### Create the swarm

At the end of this procedure, all three Docker hosts will be joined to the swarm
and will be connected together using an overlay network called `ingress`.

1.  On `manager`. initialize the swarm. If the host only has one network
    interface, the `--advertise-addr` flag is optional.

    ```console
    $ docker swarm init --advertise-addr=<IP-ADDRESS-OF-MANAGER>
    ```

    Make a note of the text that is printed, as this contains the token that
    you will use to join `worker-1` and `worker-2` to the swarm. It is a good
    idea to store the token in a password manager.

2.  On `worker-1`, join the swarm. If the host only has one network interface,
    the `--advertise-addr` flag is optional.

    ```console
    $ docker swarm join --token <TOKEN> \
      --advertise-addr <IP-ADDRESS-OF-WORKER-1> \
      <IP-ADDRESS-OF-MANAGER>:2377
    ```

3.  On `worker-2`, join the swarm. If the host only has one network interface,
    the `--advertise-addr` flag is optional.

    ```console
    $ docker swarm join --token <TOKEN> \
      --advertise-addr <IP-ADDRESS-OF-WORKER-2> \
      <IP-ADDRESS-OF-MANAGER>:2377
    ```

4.  On `manager`, list all the nodes. This command can only be done from a
    manager.

    ```console
    $ docker node ls

    ID                            HOSTNAME            STATUS              AVAILABILITY        MANAGER STATUS
    d68ace5iraw6whp7llvgjpu48 *   ip-172-31-34-146    Ready               Active              Leader
    nvp5rwavvb8lhdggo8fcf7plg     ip-172-31-35-151    Ready               Active
    ouvx2l7qfcxisoyms8mtkgahw     ip-172-31-36-89     Ready               Active
    ```

    You can also use the `--filter` flag to filter by role:

    ```console
    $ docker node ls --filter role=manager

    ID                            HOSTNAME            STATUS              AVAILABILITY        MANAGER STATUS
    d68ace5iraw6whp7llvgjpu48 *   ip-172-31-34-146    Ready               Active              Leader

    $ docker node ls --filter role=worker

    ID                            HOSTNAME            STATUS              AVAILABILITY        MANAGER STATUS
    nvp5rwavvb8lhdggo8fcf7plg     ip-172-31-35-151    Ready               Active
    ouvx2l7qfcxisoyms8mtkgahw     ip-172-31-36-89     Ready               Active
    ```

5.  List the Docker networks on `manager`, `worker-1`, and `worker-2` and notice
    that each of them now has an overlay network called `ingress` and a bridge
    network called `docker_gwbridge`. Only the listing for `manager` is shown
    here:

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
network interface so that traffic can flow to and from swarm managers and
workers. If you create swarm services and do not specify a network, they are
connected to the `ingress` network. It is recommended that you use separate
overlay networks for each application or group of applications which will work
together. In the next procedure, you will create two overlay networks and
connect a service to each of them.

#### Create the services

1.  On `manager`, create a new overlay network called `nginx-net`:

    ```console
    $ docker network create -d overlay nginx-net
    ```

    You don't need to create the overlay network on the other nodes, because it
    will be automatically created when one of those nodes starts running a
    service task which requires it.

2.  On `manager`, create a 5-replica Nginx service connected to `nginx-net`. The
    service will publish port 80 to the outside world. All of the service
    task containers can communicate with each other without opening any ports.

    > [!NOTE]
    >
    > Services can only be created on a manager.

    ```console
    $ docker service create \
      --name my-nginx \
      --publish target=80,published=80 \
      --replicas=5 \
      --network nginx-net \
      nginx
      ```

      The default publish mode of `ingress`, which is used when you do not
      specify a `mode` for the `--publish` flag, means that if you browse to
      port 80 on `manager`, `worker-1`, or `worker-2`, you will be connected to
      port 80 on one of the 5 service tasks, even if no tasks are currently
      running on the node you browse to. If you want to publish the port using
      `host` mode, you can add `mode=host` to the `--publish` output. However,
      you should also use `--mode global` instead of `--replicas=5` in this case,
      since only one service task can bind a given port on a given node.

3.  Run `docker service ls` to monitor the progress of service bring-up, which
    may take a few seconds.

4.  Inspect the `nginx-net` network on `manager`, `worker-1`, and `worker-2`.
    Remember that you did not need to create it manually on `worker-1` and
    `worker-2` because Docker created it for you. The output will be long, but
    notice the `Containers` and `Peers` sections. `Containers` lists all
    service tasks (or standalone containers) connected to the overlay network
    from that host.

5.  From `manager`, inspect the service using `docker service inspect my-nginx`
    and notice the information about the ports and endpoints used by the
    service.

6.  Create a new network `nginx-net-2`, then update the service to use this
    network instead of `nginx-net`:

    ```console
    $ docker network create -d overlay nginx-net-2
    ```

    ```console
    $ docker service update \
      --network-add nginx-net-2 \
      --network-rm nginx-net \
      my-nginx
    ```

7.  Run `docker service ls` to verify that the service has been updated and all
    tasks have been redeployed. Run `docker network inspect nginx-net` to verify
    that no containers are connected to it. Run the same command for
    `nginx-net-2` and notice that all the service task containers are connected
    to it.

    > [!NOTE]
    >
    > Even though overlay networks are automatically created on swarm
    > worker nodes as needed, they are not automatically removed.

8.  Clean up the service and the networks. From `manager`, run the following
    commands. The manager will direct the workers to remove the networks
    automatically.

    ```console
    $ docker service rm my-nginx
    $ docker network rm nginx-net nginx-net-2
    ```

## Use a user-defined overlay network

### Prerequisites

This tutorial assumes the swarm is already set up and you are on a manager.

### Walkthrough

1.  Create the user-defined overlay network.

    ```console
    $ docker network create -d overlay my-overlay
    ```

2.  Start a service using the overlay network and publishing port 80 to port
    8080 on the Docker host.

    ```console
    $ docker service create \
      --name my-nginx \
      --network my-overlay \
      --replicas 1 \
      --publish published=8080,target=80 \
      nginx:latest
    ```

3.  Run `docker network inspect my-overlay` and verify that the `my-nginx`
    service task is connected to it, by looking at the `Containers` section.

4.  Remove the service and the network.

    ```console
    $ docker service rm my-nginx

    $ docker network rm my-overlay
    ```

## Use an overlay network for standalone containers

This example demonstrates DNS container discovery -- specifically, how to
communicate between standalone containers on different Docker daemons using an
overlay network. Steps are:

- On `host1`, initialize the node as a swarm (manager).
- On `host2`, join the node to the swarm (worker).
- On `host1`, create an attachable overlay network (`test-net`).
- On `host1`, run an interactive [alpine](https://hub.docker.com/_/alpine/) container (`alpine1`) on `test-net`.
- On `host2`, run an interactive, and detached, [alpine](https://hub.docker.com/_/alpine/) container (`alpine2`) on `test-net`.
- On `host1`, from within a session of `alpine1`, ping `alpine2`.

### Prerequisites

For this test, you need two different Docker hosts that can communicate with
each other. Each host must have the following ports open between the two Docker
hosts:

- TCP port 2377
- TCP and UDP port 7946
- UDP port 4789

One easy way to set this up is to have two VMs (either local or on a cloud
provider like AWS), each with Docker installed and running. If you're using AWS
or a similar cloud computing platform, the easiest configuration is to use a
security group that opens all incoming ports between the two hosts and the SSH
port from your client's IP address.

This example refers to the two nodes in our swarm as `host1` and `host2`. This
example also uses Linux hosts, but the same commands work on Windows.

### Walk-through

1.  Set up the swarm.

    a.  On `host1`, initialize a swarm (and if prompted, use `--advertise-addr`
        to specify the IP address for the interface that communicates with other
        hosts in the swarm, for instance, the private IP address on AWS):


    ```console
    $ docker swarm init
    Swarm initialized: current node (vz1mm9am11qcmo979tlrlox42) is now a manager.

    To add a worker to this swarm, run the following command:

        docker swarm join --token SWMTKN-1-5g90q48weqrtqryq4kj6ow0e8xm9wmv9o6vgqc5j320ymybd5c-8ex8j0bc40s6hgvy5ui5gl4gy 172.31.47.252:2377

    To add a manager to this swarm, run 'docker swarm join-token manager' and follow the instructions.
    ```

    b.  On `host2`, join the swarm as instructed above:

    ```console
    $ docker swarm join --token <your_token> <your_ip_address>:2377
    This node joined a swarm as a worker.
    ```

    If the node fails to join the swarm, the `docker swarm join` command times
    out. To resolve, run `docker swarm leave --force` on `host2`, verify your
    network and firewall settings, and try again.

2.  On `host1`, create an attachable overlay network called `test-net`:

    ```console
    $ docker network create --driver=overlay --attachable test-net
    uqsof8phj3ak0rq9k86zta6ht
    ```

    > Notice the returned **NETWORK ID** -- you will see it again when you connect to it from `host2`.

3.  On `host1`, start an interactive (`-it`) container (`alpine1`) that connects to `test-net`:

    ```console
    $ docker run -it --name alpine1 --network test-net alpine
    / #
    ```

4.  On `host2`, list the available networks -- notice that `test-net` does not yet exist:

    ```console
    $ docker network ls
    NETWORK ID          NAME                DRIVER              SCOPE
    ec299350b504        bridge              bridge              local
    66e77d0d0e9a        docker_gwbridge     bridge              local
    9f6ae26ccb82        host                host                local
    omvdxqrda80z        ingress             overlay             swarm
    b65c952a4b2b        none                null                local
    ```

5.  On `host2`, start a detached (`-d`) and interactive (`-it`) container (`alpine2`) that connects to `test-net`:

    ```console
    $ docker run -dit --name alpine2 --network test-net alpine
    fb635f5ece59563e7b8b99556f816d24e6949a5f6a5b1fbd92ca244db17a4342
    ```

    > [!NOTE]
    >
    > Automatic DNS container discovery only works with unique container names.

6. On `host2`, verify that `test-net` was created (and has the same NETWORK ID as `test-net` on `host1`):

    ```console
    $ docker network ls
    NETWORK ID          NAME                DRIVER              SCOPE
    ...
    uqsof8phj3ak        test-net            overlay             swarm
    ```

7.  On `host1`, ping `alpine2` within the interactive terminal of `alpine1`:

    ```console
    / # ping -c 2 alpine2
    PING alpine2 (10.0.0.5): 56 data bytes
    64 bytes from 10.0.0.5: seq=0 ttl=64 time=0.600 ms
    64 bytes from 10.0.0.5: seq=1 ttl=64 time=0.555 ms

    --- alpine2 ping statistics ---
    2 packets transmitted, 2 packets received, 0% packet loss
    round-trip min/avg/max = 0.555/0.577/0.600 ms
    ```

    The two containers communicate with the overlay network connecting the two
    hosts. If you run another alpine container on `host2` that is _not detached_,
    you can ping `alpine1` from `host2` (and here we add the
    [remove option](/reference/cli/docker/container/run/#rm) for automatic container cleanup):

    ```sh
    $ docker run -it --rm --name alpine3 --network test-net alpine
    / # ping -c 2 alpine1
    / # exit
    ```

8.  On `host1`, close the `alpine1` session (which also stops the container):

    ```console
    / # exit
    ```

9.  Clean up your containers and networks:

    You must stop and remove the containers on each host independently because
    Docker daemons operate independently and these are standalone containers.
    You only have to remove the network on `host1` because when you stop
    `alpine2` on `host2`, `test-net` disappears.

    a.  On `host2`, stop `alpine2`, check that `test-net` was removed, then remove `alpine2`:

    ```console
    $ docker container stop alpine2
    $ docker network ls
    $ docker container rm alpine2
    ```

    a.  On `host1`, remove `alpine1` and `test-net`:

    ```console
    $ docker container rm alpine1
    $ docker network rm test-net
    ```

## Other networking tutorials

- [Host networking tutorial](/manuals/engine/network/tutorials/host.md)
- [Standalone networking tutorial](/manuals/engine/network/tutorials/standalone.md)
- [Macvlan networking tutorial](/manuals/engine/network/tutorials/macvlan.md)
