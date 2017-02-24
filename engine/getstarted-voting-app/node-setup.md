---
description: Setup for voting app example
keywords: docker-machine, multi-container, services, swarm mode, cluster, voting app
title: Set up Dockerized machines
---

The first step in getting the voting app deployed is to set up Docker machines
for the swarm nodes. You could create these Docker hosts on different physical
machines, virtual machines, or cloud providers.

For this example, we use [Docker Machine](/machine/get-started.md) to create two
virtual machines on a single system. We'll also verify the setup, and
run some basic commmands to interact with the machines.

## Prerequisites

* **Docker Machine** - These steps rely on use of
[Docker Machine](/machine/get-started.md) (`docker-machine`), which
comes auto-installed with both Docker for Mac and Docker for Windows.

* **VirtualBox driver on Docker for Mac** - On Docker for Mac, you'll use `docker-machine` with
the `virtualbox` driver to create machines. If you had a legacy installation
of Docker Toolbox, you already have Oracle VirtualBox installed as part of
that. If you started fresh with Docker for Mac, then you need to install
VirtualBox independently. We recommend doing this rather than using the Toolbox
installer because it can [conflict](/docker-for-mac/docker-toolbox.md) with
Docker for Mac. You can [download VirtualBox
here](https://www.virtualbox.org/wiki/Downloads). Click the link for `OS X
hosts`, click the `.dmg` intaller, and follow the instructions to install. You
do not need to start it, as we are simply using the driver.

* **Hyper-V driver on Docker for Windows** - On Docker for Windows, you
will use `docker-machine` with the [`Hyper-V`](/machine/drivers/hyper-v/) driver
to create machines. You will need to follow the instructions in the [Hyper-V
example](/machine/drivers/hyper-v#example) reference topic to set up a new
external network switch (a one-time task), reboot, and then
[create the machines (nodes)](/machine/drivers/hyper-v.md#create-the-nodes-with-docker-machine-and-the-microsoft-hyper-v-driver)
in an elevated PowerShell per those instructions.

### Commands to create machines

The Docker Machine commands to create local virtual machines on Mac and Windows are are as follows. The rest of the `docker-machine` commands are the same on both.

#### Mac

```
docker-machine create --driver virtualbox HOST-NAME
```

#### Windows

```
docker-machine create -d hyperv --hyperv-virtual-switch "NETWORK-SWITCH"
MACHINE-NAME`
```

This must be done in an elevated PowerShell, using a custom-created external network switch. See [Hyper-V example](/machine/drivers/hyper-v#example).

## Create manager and worker machines

Create two machines and name them to anticipate what their roles will be in the swarm:

* manager

* worker

Here is an example of creating the `manager` on Docker for Mac. Create this one, then do the same for `worker`.

```
$  docker-machine create --driver virtualbox manager
Running pre-create checks...
Creating machine...
(manager) Copying /Users/victoria/.docker/machine/cache/boot2docker.iso to /Users/victoria/.docker/machine/machines/manager/boot2docker.iso...
(manager) Creating VirtualBox VM...
(manager) Creating SSH key...
(manager) Starting the VM...
(manager) Check network to re-create if needed...
(manager) Waiting for an IP...
Waiting for machine to be running, this may take a few minutes...
Detecting operating system of created instance...
Waiting for SSH to be available...
Detecting the provisioner...
Provisioning with boot2docker...
Copying certs to the local machine directory...
Copying certs to the remote machine...
Setting Docker configuration on the remote daemon...
Checking connection to Docker...
Docker is up and running!
To see how to connect your Docker Client to the Docker Engine running on this virtual machine, run: docker-machine env manager
```

## Verify machines are running and get IP addresses

Use `docker-machine ls` to verify that the machines are
running and to get their IP addresses.

```
$ docker-machine ls
NAME       ACTIVE   DRIVER       STATE     URL                         SWARM   DOCKER        ERRORS
manager   *        virtualbox   Running   tcp://192.168.99.100:2376           v1.13.0-rc6
worker    -        virtualbox   Running   tcp://192.168.99.101:2376           v1.13.0-rc6
```

You now have two "Dockerized" machines, each running
Docker Engine, accessible through the
[Docker CLI](/engine/reference/commandline/docker.md), and available
to become swarm nodes.

You can also get the IP address of a particular machine:

```
$ docker-machine ip manager
192.168.99.100
```

You will need the IP address of the manager for a later step.

## Interacting with the machines

There are several ways to interact with these machines directly on the command line or programatically. We'll cover two methods for managing the machines directly from the command line:

#### Manage the machines from a pre-configured shell

You can use `docker-machine` to set up environment variables in a shell that connect to the Docker client on a virtual machine. With this setup, the Docker commands you type in your local shell will run on the given machine. We'll set up a shell to talk to our manager machine.

1.  Run `docker-machine env manager` to get environment variables for the manager.

    ```
    $ docker-machine env manager
    export DOCKER_TLS_VERIFY="1"
    export DOCKER_HOST="tcp://192.168.99.100:2376"
    export DOCKER_CERT_PATH="/Users/victoriabialas/.docker/machine/machines/manager"
    export DOCKER_MACHINE_NAME="manager"
    export DOCKER_API_VERSION="1.25"
    # Run this command to configure your shell:
    # eval $(docker-machine env manager)
    ```

2.  Connect your shell to the manager.
    ```
    $ eval $(docker-machine env manager)
    ```

    This sets environment variables for the current shell that the Docker client will read which specify the TLS settings.

    **Note**: If you are using `fish`, or a Windows shell such as
    Powershell/`cmd.exe` the above method will not work as described.
    Instead, see the [env command reference documentation](/machine/reference/env.md) to learn how to set the environment variables for your shell.

3.  Run `docker-machine ls` again.

    ```
    $ docker-machine ls
    NAME      ACTIVE   DRIVER       STATE     URL                         SWARM   DOCKER        ERRORS
    manager   *        virtualbox   Running   tcp://192.168.99.100:2376           v1.13.0-rc6   
    worker    -        virtualbox   Running   tcp://192.168.99.101:2376           v1.13.0-rc6   
    ```

    The asterisk next `manager` indicates that the current shell is connected to that machine. Docker commands run in this shell will execute on the `manager.` (Note that you could change this by re-running the above commands to connect to the `worker`, or open multiple terminals to talk to multiple machines.)

#### ssh into a machine

You can use the command `docker-machine ssh <MACHINE-NAME>` to log into a machine:

```
$ docker-machine ssh manager
                        ##         .
                  ## ## ##        ==
               ## ## ## ## ##    ===
           /"""""""""""""""""\___/ ===
      ~~~ {~~ ~~~~ ~~~ ~~~~ ~~~ ~ /  ===- ~~~
           \______ o           __/
             \    \         __/
              \____\_______/
 _                 _   ____     _            _
| |__   ___   ___ | |_|___ \ __| | ___   ___| | _____ _ __
| '_ \ / _ \ / _ \| __| __) / _` |/ _ \ / __| |/ / _ \ '__|
| |_) | (_) | (_) | |_ / __/ (_| | (_) | (__|   <  __/ |
|_.__/ \___/ \___/ \__|_____\__,_|\___/ \___|_|\_\___|_|

  WARNING: this is a build from test.docker.com, not a stable release.

Boot2Docker version 1.13.0-rc6, build HEAD : 5ab2289 - Wed Jan 11 23:37:52 UTC 2017
Docker version 1.13.0-rc6, build 2f2d055
```

We'll use this method later in the example.

## What's next?

In the next step, we'll [create a swarm](create-swarm.md) across these two
Docker machines.
