---
title: "Get Started, Part 4: Swarms"
keywords: swarm, scale, cluster, machine, vm, manager, worker, deploy, ssh, orchestration
description: Learn how to create clusters of Dockerized machines.
---
{% include_relative nav.html selected="4" %}

## Prerequisites

- [Install Docker version 1.13 or higher](/engine/installation/index.md).

- Get [Docker Compose](/compose/overview.md) as described in [Part 3 prerequisites](/get-started/part3.md#prerequisites).

- Get [Docker Machine](/machine/overview.md), which is pre-installed with
[Docker for Mac](/docker-for-mac/index.md) and [Docker for
Windows](/docker-for-windows/index.md), but on Linux systems you need to
[install it directly](/machine/install-machine/#installing-machine-directly). On pre Windows 10 systems _without Hyper-V_, use [Docker
Toolbox](https://docs.docker.com/toolbox/overview.md).

- Read the orientation in [Part 1](index.md).
- Learn how to create containers in [Part 2](part2.md).
- Make sure you have published the `friendlyhello` image you created by
[pushing it to a registry](/get-started/part2.md#share-your-image). We will be using that shared image here.
- Be sure your image works as a deployed container by running this command, and visting `http://localhost/` (slotting in your info for `username`,
`repo`, and `tag`):

  ```shell
  docker run -p 80:80 username/repo:tag
  ```
- Have a copy of your `docker-compose.yml` from [Part 3](part3.md) handy.

## Introduction

In [part 3](part3.md), you took an app you wrote in [part 2](part2.md), and
defined how it should run in production by turning it into a service, scaling it
up 5x in the process.

Here in part 4, you deploy this application onto a cluster, running it on
multiple machines. Multi-container, multi-machine applications are made possible
by joining multiple machines into a "Dockerized" cluster called a **swarm**.

## Understanding Swarm clusters

A swarm is a group of machines that are running Docker and joined into
a cluster. After that has happened, you continue to run the Docker commands
you're used to, but now they are executed on a cluster by a **swarm manager**.
The machines in a swarm can be physical or virtual. After joining a swarm, they
are referred to as **nodes**.

Swarm managers can use several strategies to run containers, such as "emptiest
node" -- which fills the least utilized machines with containers. Or "global",
which ensures that each machine gets exactly one instance of the specified
container. You instruct the swarm manager to use these strategies in the Compose
file, just like the one you have already been using.

Swarm managers are the only machines in a swarm that can execute your commands,
or authorize other machines to join the swarm as **workers**. Workers are just
there to provide capacity and do not have the authority to tell any other
machine what it can and cannot do.

Up until now, you have been using Docker in a single-host mode on your local
machine. But Docker also can be switched into **swarm mode**, and that's what
enables the use of swarms. Enabling swarm mode instantly makes the current
machine a swarm manager. From then on, Docker will run the commands you execute
on the swarm you're managing, rather than just on the current machine.

{% capture local-instructions %}
You now have two VMs created, named `myvm1` and `myvm2` (as `docker-machine ls`
shows). The first one will act as the manager, which executes `docker` commands
and authenticates workers to join the swarm, and the second will be a worker.

You can send commands to your VMs using `docker-machine ssh`. Instruct `myvm1`
to become a swarm manager with `docker swarm init` and you'll see output like
this:

```shell
$ docker-machine ssh myvm1 "docker swarm init"
Swarm initialized: current node <node ID> is now a manager.

To add a worker to this swarm, run the following command:

  docker swarm join \
  --token <token> \
  <ip>:<port>
```

> Got an error about needing to use `--advertise-addr`?
>
> Copy the
> IP address for `myvm1` by running `docker-machine ls`, then run the
> `docker swarm init` command again, using that IP and specifying port `2377`
> (the port for swarm joins) with `--advertise-addr`. For example:
>
> ```
> docker-machine ssh myvm1 "docker swarm init --advertise-addr 192.168.99.100:2377"
> ```
{: .note-vanilla}

As you can see, the response to `docker swarm init` contains a pre-configured
`docker swarm join` command for you to run on any nodes you want to add. Copy
this command, and send it to `myvm2` via `docker-machine ssh` to have `myvm2`
join your new swarm as a worker:

```shell
$ docker-machine ssh myvm2 "docker swarm join \
--token <token> \
<ip>:<port>"

This node joined a swarm as a worker.
```

Congratulations, you have created your first swarm.

> **Note**: You can also run `docker-machine ssh myvm2` with no command attached
to open a terminal session on that VM. Type `exit` when you're ready to return
to the host shell prompt. It may be easier to paste the join command in that
way.

Use `ssh` to connect to the (`docker-machine ssh myvm1`), and run `docker node ls` to view the nodes in this swarm:

```shell
docker@myvm1:~$ docker node ls
ID                            HOSTNAME            STATUS              AVAILABILITY        MANAGER STATUS
brtu9urxwfd5j0zrmkubhpkbd     myvm2               Ready               Active              
rihwohkh3ph38fhillhhb84sk *   myvm1               Ready               Active              Leader
```

Type `exit` to get back out of that machine.

Alternatively, wrap commands in `docker-machine ssh` to keep from having to directly log in and out. For example:

```shell
docker-machine ssh myvm1 "docker node ls"
```


{% endcapture %}

{% capture local %}
#### VMs on your local machine (Mac, Linux, Windows 7 and 8)

First, you'll need a hypervisor that can create VMs, so [install
VirtualBox](https://www.virtualbox.org/wiki/Downloads) for your machine's OS.

> **Note**: If you're on a Windows system that has Hyper-V installed, such as
Windows 10, there is no need to install VirtualBox and you should use Hyper-V
instead. View the instructions for Hyper-V systems by clicking the Hyper-V tab
above.

Now, create a couple of VMs using `docker-machine`, using the VirtualBox driver:

```none
$ docker-machine create --driver virtualbox myvm1
$ docker-machine create --driver virtualbox myvm2
```

{{ local-instructions }}
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

```shell
$ docker-machine create -d hyperv --hyperv-virtual-switch "myswitch" myvm1
$ docker-machine create -d hyperv --hyperv-virtual-switch "myswitch" myvm2
```

{{ local-instructions }}
{% endcapture %}

## Set up your swarm

A swarm is made up of multiple nodes, which can be either physical or virtual
machines. The basic concept is simple enough: run `docker swarm init` to enable
swarm mode and make your current machine a swarm manager, then run
`docker swarm join` on other machines to have them join the swarm as workers.
Choose a tab below to see how this plays out in various contexts. We'll use VMs
to quickly create a two-machine cluster and turn it into a swarm.

### Create a cluster

<ul class="nav nav-tabs">
  <li class="active"><a data-toggle="tab" href="#local">Local VMs (Mac, Linux, Windows 7 and 8)</a></li>
  <li><a data-toggle="tab" href="#localwin">Local VMs (Windows 10/Hyper-V)</a></li>
</ul>
<div class="tab-content">
  <div id="local" class="tab-pane fade in active" markdown="1">{{ local }}</div>
  <div id="localwin" class="tab-pane fade" markdown="1">{{ localwin }}</div>
</div>

## Deploy your app on a cluster

The hard part is over. Now you just repeat the process you used in [part
3](part3.md) to deploy on your new swarm. Just remember that only swarm managers
like `myvm1` execute Docker commands; workers are just for capacity.

Copy the file `docker-compose.yml` you created in part 3 to the swarm manager
`myvm1`'s home directory (alias: `~`) by using the `docker-machine scp` command:

```
docker-machine scp docker-compose.yml myvm1:~
```

Now have `myvm1` use its powers as a swarm manager to deploy your app, by sending
the same `docker stack deploy` command you used in part 3 to `myvm1` using
`docker-machine ssh`:

```
docker-machine ssh myvm1 "docker stack deploy -c docker-compose.yml getstartedlab"
```

And that's it, the app is deployed on a cluster.

Wrap all the commands you used in part 3 in a call to `docker-machine ssh`, and
they'll all work as you'd expect. Only this time, you'll see that the containers
have been distributed between both `myvm1` and `myvm2`.

```
$ docker-machine ssh myvm1 "docker stack ps getstartedlab"

ID            NAME        IMAGE              NODE   DESIRED STATE
jq2g3qp8nzwx  test_web.1  username/repo:tag  myvm1  Running
88wgshobzoxl  test_web.2  username/repo:tag  myvm2  Running
vbb1qbkb0o2z  test_web.3  username/repo:tag  myvm2  Running
ghii74p9budx  test_web.4  username/repo:tag  myvm1  Running
0prmarhavs87  test_web.5  username/repo:tag  myvm2  Running
```

### Accessing your cluster

You can access your app from the IP address of **either** `myvm1` or `myvm2`.
The network you created is shared between them and load-balancing. Run
`docker-machine ls` to get your VMs' IP addresses and visit either of them on a
browser, hitting refresh (or just `curl` them). You'll see five possible
container IDs all cycling by randomly, demonstrating the load-balancing.

The reason both IP addresses work is that nodes in a swarm participate in an
ingress **routing mesh**. This ensures that a service deployed at a certain port
within your swarm always has that port reserved to itself, no matter what node
is actually running the container. Here's a diagram of how a routing mesh for a
service called `my-web` published at port `8080` on a three-node swarm would
look:

![routing mesh diagram](/engine/swarm/images/ingress-routing-mesh.png)

> Having connectivity trouble?
>
> Keep in mind that in order to use the ingress network in the swarm,
> you need to have the following ports open between the swarm nodes
> before you enable swarm mode:
>
> - Port 7946 TCP/UDP for container network discovery.
> - Port 4789 UDP for the container ingress network.
{: .note-vanilla}

## Iterating and scaling your app

From here you can do everything you learned about in part 3.

Scale the app by changing the `docker-compose.yml` file.

Change the app behavior by editing code.

In either case, simply run `docker stack deploy` again to deploy these
changes.

You can join any machine, physical or virtual, to this swarm, using the
same `docker swarm join` command you used on `myvm2`, and capacity will be added
to your cluster. Just run `docker stack deploy` afterwards, and your app will
take advantage of the new resources.

## Cleanup

You can tear down the stack with `docker stack rm`. For example:

```
docker-machine ssh myvm1 "docker stack rm getstartedlab"
```

> Keep the swarm or remove it?
>
> At some point later, you can remove this swarm if you want to with
> `docker-machine ssh myvm2 "docker swarm leave"` on the worker
> and `docker-machine ssh myvm1 "docker swarm leave --force"` on the
> manager, but _you'll need this swarm for part 5, so please keep it
> around for now_.
{: .note-vanilla}

[On to Part 5 >>](part5.md){: class="button outline-btn" style="margin-bottom: 30px"}

## Recap and cheat sheet (optional)

Here's [a terminal recording of what was covered on this page](https://asciinema.org/a/113837):

<script type="text/javascript" src="https://asciinema.org/a/113837.js" id="asciicast-113837" speed="2" async></script>

In part 4 you learned what a swarm is, how nodes in swarms can be managers or
workers, created a swarm, and deployed an application on it. You saw that the
core Docker commands didn't change from part 3, they just had to be targeted to
run on a swarm master. You also saw the power of Docker's networking in action,
which kept load-balancing requests across containers, even though they were
running on different machines. Finally, you learned how to iterate and scale
your app on a cluster.

Here are some commands you might like to run to interact with your swarm a bit:

```shell
docker-machine create --driver virtualbox myvm1 # Create a VM (Mac, Win7, Linux)
docker-machine create -d hyperv --hyperv-virtual-switch "myswitch" myvm1 # Win10
docker-machine env myvm1                # View basic information about your node
docker-machine ssh myvm1 "docker node ls"         # List the nodes in your swarm
docker-machine ssh myvm1 "docker node inspect <node ID>"        # Inspect a node
docker-machine ssh myvm1 "docker swarm join-token -q worker"   # View join token
docker-machine ssh myvm1   # Open an SSH session with the VM; type "exit" to end
docker-machine ssh myvm2 "docker swarm leave"  # Make the worker leave the swarm
docker-machine ssh myvm1 "docker swarm leave -f" # Make master leave, kill swarm
docker-machine start myvm1            # Start a VM that is currently not running
docker-machine stop $(docker-machine ls -q)               # Stop all running VMs
docker-machine rm $(docker-machine ls -q) # Delete all VMs and their disk images
docker-machine scp docker-compose.yml myvm1:~     # Copy file to node's home dir
docker-machine ssh myvm1 "docker stack deploy -c <file> <app>"   # Deploy an app
```
