---
title: "Getting Started, Part 4: Scaling Your App on a Cluster"
---

In [Getting Started, Part 3: Stateful, Multi-container Applications](part3.md),
we figured out how to relate containers to each other. We organized an
application into two simple services -- a frontend and a backend -- and defined
how they are linked together.

In Part 4, we are going to take this application, which all ran on one **host**
(a virtual or physical machine), and deploy it onto a cluster of hosts, just
like we would in production.

## Understanding Swarm clusters

Up until now you have been using Docker in a single-host mode on your local
machine, which allows the client, which is the command-line interface (CLI), to
make assumptions about how to operate. Namely, the client assumes that the
Docker Daemon is running on the same host as the client. Single-host operations
could also be done on remote machines with your client.

But Docker also can be switched into "swarm mode." A swarm is a group of hosts
that are running Docker and have been joined into a cluster. After that has
happened, you continue to run the Docker commands you're used to, but now the
concept of a "host" changes from a single virtual or physical machine, to a
swarm. And, "a single virtual or physical machine" is not referred to as a host,
it's called a node -- or, a computing resource inside your cluster.

{% capture local-instructions %}
We now have two VMs created, named `myvm1` and `myvm2`. The first one will act
as the manager, which executes `docker` commands and authenticates workers to
join the swarm, and the second will be a worker.

You can send commands to your VMs using `docker-machine ssh`. Instruct `myvm1`
to become a swarm manager with `docker swarm init` and you'll see output like
this:

```node
$ docker-machine ssh myvm1 "docker swarm init"
Swarm initialized: current node <node ID> is now a manager.

To add a worker to this swarm, run the following command:

  docker swarm join \
  --token <token> \
  <ip>:<port>
```

> **Note**: Getting an error about needing to use `--advertise-addr`? Copy the
> IP address for `myvm1` by running `docker-machine ls`, then, run
> `docker swarm join` command again, using the IP and specifying port `2377` for
> `--advertise-addr`. For example:
>
> `docker-machine ssh myvm1 "docker swarm join --advertise-addr 192.168.99.100:2377"`
>
>

As you can see, the response to `docker swarm init` contains a pre-configured
`docker swarm join` command for you to run on any nodes you want to add. Copy
this command, and send it to `myvm2` via `docker-machine ssh` to have `myvm2`
join your new swarm as a worker:

```
$ docker-machine ssh myvm2 "docker swarm join \
--token <token> \
<ip>:<port>"

This node joined a swarm as a worker.
```

> **Note**: You can also run `docker-machine ssh myvm2` with no command attached
to open a terminal session on that VM. Type `exit` when you're ready to return
to the host shell prompt.

Congratulations, you created your first swarm!

Here are some commands you might like to run to interact with your swarm a bit:

```
docker-machine env myvm1                # View basic information about your node
docker-machine ssh myvm1 "docker node ls"         # List the nodes in your swarm
docker-machine ssh myvm1 "docker node inspect <node ID>"        # Inspect a node
docker-machine ssh myvm1 "docker swarm join-token -q worker"   # View join token
docker-machine ssh myvm1   # Open an SSH session with the VM; type "exit" to end
docker-machine ssh myvm2 "docker swarm leave"  # Make the worker leave the swarm
docker-machine stop $(docker-machine ls -q)               # Stop all running VMs
docker-machine rm $(docker-machine ls -q) # Delete all VMs and their disk images
```
{% endcapture %}

{% capture local %}
#### VMs on your local machine (Mac, Linux, Windows 7 and 8)

First, you'll need a hypervisor that can create VMs, so [install
VirtualBox](https://www.virtualbox.org/wiki/Downloads) for your machine's OS.

> **Note**: If you're on a Windows system that has Hyper-V installed, such as
Windows 10, there is no need for this step and you should use Hyper-V instead.
View the instructions for Hyper-V systems by clicking the Hyper-V tab above.

Now, create a couple of VMs using `docker-machine`, using the VirtualBox driver:

```none
$ docker-machine create --driver virtualbox myvm1
$ docker-machine create --driver virtualbox myvm2
```

{{ local-instructions | markdownify }}
{% endcapture %}

{% capture localwin %}
#### VMs on your local machine (Windows 10)

First, quickly create a virtual switch for your VMs to share, so they will be
able to connect to each other.

1. Launch Hyper-V Manager
2. Click **Virtual Switch Manager** in the right-hand menu
3. Click **Create Virtual Switch** of type **External**
4. Give it the name `myswitch`, and check the box to share your host machine's
   active network adapter

Now, create a couple of virtual machines using our node management tool,
`docker-machine`:

```none
$ docker-machine create -d hyperv --hyperv-virtual-switch "myswitch" myvm1
$ docker-machine create -d hyperv --hyperv-virtual-switch "myswitch" myvm2
```

{{ local-instructions | markdownify}}
{% endcapture %}

## Set up your swarm

A swarm is made up of multiple nodes, which can be either physical or virtual
machines. The basic concept is simple enough: run `docker swarm init` to make
your current machine a manager node, and run `docker swarm join` on other
machines to have them join the swarm as a worker. Choose a tab below to see how
this plays out in various contexts.

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" href="#local">Local VMs (Mac, Linux, Windows 7 and 8)</a></li>
  <li><a data-toggle="tab" href="#localwin">Local VMs (Windows 10/Hyper-V)</a></li>
</ul>
<div class="tab-content">
  <div id="local" class="tab-pane fade in active">{{ local | markdownify }}</div>
  <div id="localwin" class="tab-pane fade">{{ localwin | markdownify }}</div>
</div>

[On to next >>](part5.md){: class="button darkblue-btn"}
