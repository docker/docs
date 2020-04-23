---
description: Getting Started tutorial for Docker Engine swarm mode
keywords: tutorial, cluster management, swarm mode
title: Getting started with swarm mode
toc_max: 4
---

This tutorial introduces you to the features of Docker Engine Swarm mode. You
may want to familiarize yourself with the [key concepts](../key-concepts.md)
before you begin.

The tutorial guides you through the following activities:

* initializing a cluster of Docker Engines in swarm mode
* adding nodes to the swarm
* deploying application services to the swarm
* managing the swarm once you have everything running

This tutorial uses Docker Engine CLI commands entered on the command line of a
terminal window.

If you are brand new to Docker, see [About Docker Engine](../../index.md).

## Set up

To run this tutorial, you need the following:

* [three Linux hosts which can communicate over a network, with Docker installed](#three-networked-host-machines)
* [Docker Engine 1.12 or later installed](#docker-engine-112-or-newer)
* [the IP address of the manager machine](#the-ip-address-of-the-manager-machine)
* [open ports between the hosts](#open-protocols-and-ports-between-the-hosts)

### Three networked host machines

This tutorial requires three Linux hosts which have Docker installed and can
communicate over a network. These can be physical machines, virtual machines,
Amazon EC2 instances, or hosted in some other way. You can even use Docker Machine
from a Linux, Mac, or Windows host. Check out
[Getting started - Swarms](../../../get-started/swarm-deploy.md#prerequisites)
for one possible set-up for the hosts.

One of these machines is a manager (called `manager1`) and two of them are
workers (`worker1` and `worker2`).


>**Note**: You can follow many of the tutorial steps to test single-node swarm
as well, in which case you need only one host. Multi-node commands do not
work, but you can initialize a swarm, create services, and scale them.

###  Docker Engine 1.12 or newer

This tutorial requires Docker Engine 1.12 or newer on each of the host machines.
Install Docker Engine and verify that the Docker Engine daemon is running on
each of the machines. You can get the latest version of Docker Engine as
follows:

* [install Docker Engine on Linux machines](#install-docker-engine-on-linux-machines)

* [use Docker Desktop for Mac or Docker Desktop for Windows](#use-docker-desktop-for-mac-or-docker-desktop-for-windows)

#### Install Docker Engine on Linux machines

If you are using Linux based physical computers or cloud-provided computers as
hosts, simply follow the [Linux install instructions](../../install/index.md)
for your platform. Spin up the three machines, and you are ready. You can test both
single-node and multi-node swarm scenarios on Linux machines.

#### Use Docker Desktop for Mac or Docker Desktop for Windows

Alternatively, install the latest [Docker Desktop for Mac](../../../docker-for-mac/index.md) or
[Docker Desktop for Windows](../../../docker-for-windows/index.md) application on one
computer. You can test both single-node and multi-node swarm from this computer,
but you need to use Docker Machine to test the multi-node scenarios.

* You can use Docker Desktop for Mac or Windows to test _single-node_ features of swarm
mode, including initializing a swarm with a single node, creating services,
and scaling services. Docker "Moby" on Hyperkit (Mac) or Hyper-V (Windows)
serve as the single swarm node.

<p />

* Currently, you cannot use Docker Desktop for Mac or Docker Desktop for Windows alone to test a
_multi-node_ swarm. However, you can use the included version of
[Docker Machine](../../../machine/overview.md) to create the swarm nodes (see
[Get started with Docker Machine and a local VM](../../../machine/get-started.md)), then
follow the tutorial for all multi-node features. For this scenario, you run
commands from a Docker Desktop for Mac or Docker Desktop for Windows host, but that Docker host itself is
_not_ participating in the swarm. After you create the nodes, you can run all
swarm commands as shown from the Mac terminal or Windows PowerShell with
Docker Desktop for Mac or Docker Desktop for Windows running.

### The IP address of the manager machine

The IP address must be assigned to a network interface available to the host
operating system. All nodes in the swarm need to connect to the manager at
the IP address.

Because other nodes contact the manager node on its IP address, you should use a
fixed IP address.

You can run `ifconfig` on Linux or macOS to see a list of the
available network interfaces.

If you are using Docker Machine, you can get the manager IP with either
`docker-machine ls` or `docker-machine ip <MACHINE-NAME>` &#8212; for example,
`docker-machine ip manager1`.

The tutorial uses `manager1` : `192.168.99.100`.

### Open protocols and ports between the hosts

The following ports must be available. On some systems, these ports are open by default.

* **TCP port 2377** for cluster management communications
* **TCP** and **UDP port 7946** for communication among nodes
* **UDP port 4789** for overlay network traffic

If you plan on creating an overlay network with encryption (`--opt encrypted`),
you also need to ensure **ip protocol 50** (**ESP**) traffic is allowed.

## What's next?

After you have set up your environment, you are ready to [create a swarm](create-swarm.md).
