---
description: Run Docker Engine in swarm mode
keywords: guide, swarm mode, node
title: Run Docker Engine in swarm mode
---

When you first install and start working with Docker Engine, swarm mode is
disabled by default. When you enable swarm mode, you work with the concept of
services managed through the `docker service` command.

There are two ways to run the Engine in swarm mode:

* Create a new swarm, covered in this article.
* [Join an existing swarm](join-nodes.md).

When you run the Engine in swarm mode on your local machine, you can create and
test services based upon images you've created or other available images. In
your production environment, swarm mode provides a fault-tolerant platform with
cluster management features to keep your services running and available.

These instructions assume you have installed the Docker Engine 1.12 or later on
a machine to serve as a manager node in your swarm.

If you haven't already, read through the [swarm mode key concepts](key-concepts.md)
and try the [swarm mode tutorial](swarm-tutorial/index.md).

## Create a swarm

When you run the command to create a swarm, the Docker Engine starts running in swarm mode.

Run [`docker swarm init`](../reference/commandline/swarm_init.md)
to create a single-node swarm on the current node. The Engine sets up the swarm
as follows:

* switches the current node into swarm mode.
* creates a swarm named `default`.
* designates the current node as a leader manager node for the swarm.
* names the node with the machine hostname.
* configures the manager to listen on an active network interface on port 2377.
* sets the current node to `Active` availability, meaning it can receive tasks
from the scheduler.
* starts an internal distributed data store for Engines participating in the
swarm to maintain a consistent view of the swarm and all services running on it.
* by default, generates a self-signed root CA for the swarm.
* by default, generates tokens for worker and manager nodes to join the
swarm.
* creates an overlay network named `ingress` for publishing service ports
external to the swarm.
* creates an overlay default IP addresses and subnet mask for your networks

The output for `docker swarm init` provides the connection command to use when
you join new worker nodes to the swarm:

```bash
$ docker swarm init
Swarm initialized: current node (dxn1zf6l61qsb1josjja83ngz) is now a manager.

To add a worker to this swarm, run the following command:

    docker swarm join \
    --token SWMTKN-1-49nj1cmql0jkz5s954yi3oex3nedyz0fb0xx14ie39trti4wxv-8vxv8rssmk743ojnwacrr2e7c \
    192.168.99.100:2377

To add a manager to this swarm, run 'docker swarm join-token manager' and follow the instructions.
```
### Configuring default address pools

By default Docker Swarm uses a default address pool `10.0.0.0/8` for global scope (overlay) networks. Every 
network that does not have a subnet specified will have a subnet sequentially allocated from this pool. In 
some circumstances it may be desirable to use a different default IP address pool for networks. 

For example, if the default `10.0.0.0/8` range conflicts with already allocated address space in your network, 
then it is desirable to ensure that networks use a different range without requiring Swarm users to specify 
each subnet with the `--subnet` command. 

To configure custom default address pools, you must define pools at Swarm initialization using the 
`--default-addr-pool` command line option. This command line option uses CIDR notation for defining the subnet mask.
To create the custom address pool for Swarm, you must define at least one default address pool, and an optional default address pool subnet mask. For example, for the `10.0.0.0/27`, use the value `27`.

Docker allocates subnet addresses from the address ranges specified by the `--default-addr-pool` option. For example, a command line option `--default-addr-pool 10.10.0.0/16` indicates that Docker will allocate subnets from that `/16` address range. If `--default-addr-pool-mask-len` were unspecified or set explicitly to 24, this would result in 256 `/24` networks of the form `10.10.X.0/24`.

The subnet range comes from the `--default-addr-pool`, (such as `10.10.0.0/16`). The size of 16 there represents the number of networks one can create within that `default-addr-pool` range. The `--default-addr-pool` option may occur multiple times with each option providing additional addresses for docker to use for overlay subnets.

The format of the command is:

```bash
$ docker swarm init --default-addr-pool <IP range in CIDR> [--default-addr-pool <IP range in CIDR> --default-addr-pool-mask-length <CIDR value>]
```

To create a default IP address pool with a /16 (class B) for the 10.20.0.0 network looks like this:

```bash
$ docker swarm init --default-addr-pool 10.20.0.0/16
```

To create a default IP address pool with a `/16` (class B) for the `10.20.0.0` and `10.30.0.0` networks, and to 
create a subnet mask of `/26` for each network looks like this:

```bash
$ docker swarm init --default-addr-pool 10.20.0.0/16 --default-addr-pool 10.30.0.0/16 --default-addr-pool-mask-length 26
```

In this example, `docker network create -d overlay net1` will result in `10.20.0.0/26` as the allocated subnet for `net1`, 
and `docker network create -d overlay net2` will result in `10.20.0.64/26` as the allocated subnet for `net2`. This continues until 
all the subnets are exhausted. 

Refer to the following pages for more information:
- [Swarm networking](./networking.md) for more information about the default address pool usage
- `docker swarm init` [CLI reference](../reference/commandline/swarm_init.md) for more detail on the `--default-addr-pool` flag.

### Configure the advertise address

Manager nodes use an advertise address to allow other nodes in the swarm access
to the Swarmkit API and overlay networking. The other nodes on the swarm must be
able to access the manager node on its advertise address.

If you don't specify an advertise address, Docker checks if the system has a
single IP address. If so, Docker uses the IP address with the listening port
`2377` by default. If the system has multiple IP addresses, you must specify the
correct `--advertise-addr` to enable inter-manager communication and overlay
networking:

```bash
$ docker swarm init --advertise-addr <MANAGER-IP>
```

You must also specify the `--advertise-addr` if the address where other nodes
reach the first manager node is not the same address the manager sees as its
own. For instance, in a cloud setup that spans different regions, hosts have
both internal addresses for access within the region and external addresses that
you use for access from outside that region. In this case, specify the external
address with `--advertise-addr` so that the node can propagate that information
to other nodes that subsequently connect to it.

Refer to the `docker swarm init` [CLI reference](../reference/commandline/swarm_init.md)
for more detail on the advertise address.

### View the join command or update a swarm join token

Nodes require a secret token to join the swarm. The token for worker nodes is
different from the token for manager nodes. Nodes only use the join-token at the
moment they join the swarm. Rotating the join token after a node has already
joined a swarm does not affect the node's swarm membership. Token rotation
ensures an old token cannot be used by any new nodes attempting to join the
swarm.

To retrieve the join command including the join token for worker nodes, run:

```bash
$ docker swarm join-token worker

To add a worker to this swarm, run the following command:

    docker swarm join \
    --token SWMTKN-1-49nj1cmql0jkz5s954yi3oex3nedyz0fb0xx14ie39trti4wxv-8vxv8rssmk743ojnwacrr2e7c \
    192.168.99.100:2377

This node joined a swarm as a worker.
```

To view the join command and token for manager nodes, run:

```bash
$ docker swarm join-token manager

To add a worker to this swarm, run the following command:

    docker swarm join \
    --token SWMTKN-1-59egwe8qangbzbqb3ryawxzk3jn97ifahlsrw01yar60pmkr90-bdjfnkcflhooyafetgjod97sz \
    192.168.99.100:2377
```

Pass the `--quiet` flag to print only the token:

```bash
$ docker swarm join-token --quiet worker

SWMTKN-1-49nj1cmql0jkz5s954yi3oex3nedyz0fb0xx14ie39trti4wxv-8vxv8rssmk743ojnwacrr2e7c
```

Be careful with the join tokens because they are the secrets necessary to join
the swarm. In particular, checking a secret into version control is a bad
practice because it would allow anyone with access to the application source
code to add new nodes to the swarm. Manager tokens are especially sensitive
because they allow a new manager node to join and gain control over the whole
swarm.

We recommend that you rotate the join tokens in the following circumstances:

* If a token was checked-in by accident into a version control system, group
chat or accidentally printed to your logs.
* If you suspect a node has been compromised.
* If you wish to guarantee that no new nodes can join the swarm.

Additionally, it is a best practice to implement a regular rotation schedule for
any secret including swarm join tokens. We recommend that you rotate your tokens
at least every 6 months.

Run `swarm join-token --rotate` to invalidate the old token and generate a new
token. Specify whether you want to rotate the token for `worker` or `manager`
nodes:

```bash
$ docker swarm join-token  --rotate worker

To add a worker to this swarm, run the following command:

    docker swarm join \
    --token SWMTKN-1-2kscvs0zuymrsc9t0ocyy1rdns9dhaodvpl639j2bqx55uptag-ebmn5u927reawo27s3azntd44 \
    192.168.99.100:2377
```

## Learn more

* [Join nodes to a swarm](join-nodes.md)
* `swarm init` [command line reference](../reference/commandline/swarm_init.md)
* [Swarm mode tutorial](swarm-tutorial/index.md)
