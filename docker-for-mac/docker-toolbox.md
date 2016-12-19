---
description: Docker for Mac and Docker Toolbox
keywords: mac, windows, alpha, beta, toolbox, docker-machine, tutorial
redirect_from:
- /mackit/docker-toolbox/
title: Docker for Mac vs. Docker Toolbox
---

If you already have an installation of Docker Toolbox, please read these topics first to learn how Docker for Mac and Docker Toolbox differ, and how they can coexist.

## The Docker Toolbox environment

Docker Toolbox installs `docker`, `docker-compose` and `docker-machine` in `/usr/local/bin` on your Mac. It also installs VirtualBox. At installation time, Toolbox uses `docker-machine` to provision a VirtualBox VM called `default`, running the `boot2docker` Linux distribution, with  [Docker Engine](/engine/) with certificates located on your Mac at `$HOME/.docker/machine/machines/default`.

Before you use `docker` or `docker-compose` on your Mac, you typically use the command `eval $(docker-machine env default)` to set environment variables so that `docker` or `docker-compose` know how to talk to Docker Engine running on VirtualBox.

This setup is shown in the following diagram.

![Docker Toolbox Install](images/toolbox-install.png)


## The Docker for Mac environment

Docker for Mac is a Mac native application, that you install in `/Applications`. At installation time, it creates symlinks in `/usr/local/bin` for `docker` and `docker-compose`, to the version of the commands inside the Mac application bundle, in `/Applications/Docker.app/Contents/Resources/bin`.

Here are some key points to know about Docker for Mac before you get started:

* Docker for Mac does not use VirtualBox, but rather <a href="https://github.com/docker/HyperKit/" target="_blank">HyperKit</a>, a lightweight macOS virtualization solution built on top of Hypervisor.framework in macOS 10.10 Yosemite and higher.

* Installing Docker for Mac does not affect machines you created with Docker Machine. The install offers to copy containers and images from your local `default` machine (if one exists) to the new Docker for Mac HyperKit VM. If chosen, content from `default` is copied to the new Docker for Mac HyperKit VM, and your original `default` machine is kept as is.

* The Docker for Mac application does not use `docker-machine` to provision that VM; but rather creates and manages it directly.

* At installation time, Docker for Mac provisions an HyperKit VM based on Alpine Linux, running Docker Engine. It exposes the docker API on a socket in `/var/run/docker.sock`. Since this is the default location where `docker` will look if no environment variables are set, you can start using `docker` and `docker-compose` without setting any environment variables.

This setup is shown in the following diagram.

![Docker for Mac Install](images/docker-for-mac-install.png)

With Docker for Mac, you get only one VM, and you don't manage it. It is managed by the Docker for Mac application, which includes autoupdate to update the client and server versions of Docker.

If you need several VMs and want to manage the version of the Docker client or server you are using, you can continue to use `docker-machine`, on the same machine, as described in [Docker Toolbox and Docker for Mac coexistence](docker-toolbox.md#docker-toolbox-and-docker-for-mac-coexistence).


## Setting up to run Docker for Mac

1. Check whether Toolbox DOCKER environment variables are set:

        $ env | grep DOCKER
        DOCKER_HOST=tcp://192.168.99.100:2376
        DOCKER_MACHINE_NAME=default
        DOCKER_TLS_VERIFY=1
        DOCKER_CERT_PATH=/Users/victoriabialas/.docker/machine/machines/default

    If this command returns no output, you are ready to use Docker for Mac.

    If it returns output (as shown in the example), you need to unset the `DOCKER` environment variables to make the client talk to the Docker for Mac Engine (next step).

2. Run the `unset` command on the following `DOCKER` environment variables to unset them in the current shell.

        unset DOCKER_TLS_VERIFY
        unset DOCKER_CERT_PATH
        unset DOCKER_MACHINE_NAME
        unset DOCKER_HOST

  Now, this command should return no output.

          $ env | grep DOCKER

  If you are using a Bash shell, you can use `unset ${!DOCKER_*}` to unset all DOCKER environment variables at once. (This will not work in other shells such as `zsh`; you will need to unset each variable individually.)

>**Note**: If you have a shell script as part of your profile that sets these `DOCKER` environment variables automatically each time you open a command window, then you will need to unset these each time you want to use Docker for Mac.

> **Warning**: If you install Docker for Mac on a machine where Docker Toolbox is installed, it will replace the `docker` and `docker-compose` command lines in `/usr/local/bin` with symlinks to its own versions.


## Docker Toolbox and Docker for Mac coexistence

You can use Docker for Mac and Docker Toolbox together on the same machine. When you want to use Docker for Mac, make sure all DOCKER environment variables are unset. You can do this in bash with `unset ${!DOCKER_*}`. When you want to use one of the VirtualBox VMs you have set with `docker-machine`, just run a `eval $(docker-machine env default)` (or the name of the machine you want to target). This will switch the current command shell to talk to the specified Toolbox machine.

This setup is represented in the following diagram.

![Docker Toolbox and Docker for Mac coexistence](images/docker-for-mac-and-toolbox.png)


## Using different versions of Docker tools

The coexistence setup works as is as long as your VirtualBox VMs provisioned with `docker-machine` run the same version of Docker Engine as Docker for Mac. If you need to use VMs running older versions of Docker Engine, you can use a tool like <a href="https://github.com/getcarina/dvm" target="_blank">Docker Version Manager</a> to manage several versions of docker client.


### Checking component versions

Ideally, the Docker CLI client and Docker Engine should be the same version. Mismatches between client and server, and among host machines you might have created with Docker Machine can cause problems (client can't talk to the server or host machines).

If you already have <a href="/toolbox/overview/" target="_blank">Docker Toolbox</a> installed, and then install Docker for Mac, you might get a newer version of the Docker client. Run `docker version` in a command shell to see client and server versions. In this example, the client installed with Docker for Mac is `Version: 1.11.1` and the server (which was installed earlier with Toolbox) is Version: 1.11.0.

    $ docker version
    Client:
    Version:      1.11.1
    ...

    Server:
    Version:      1.11.0
    ...

Also, if you created machines with Docker Machine (installed with Toolbox) then upgraded or installed Docker for Mac, you might have machines running different versions of Engine. Run `docker-machine ls` to view version information for the machines you created. In this example, the DOCKER column shows that each machine is running a different version of server.

    $ docker-machine ls
    NAME             ACTIVE   DRIVER         STATE     URL                         SWARM   DOCKER    ERRORS
    aws-sandbox      -        amazonec2      Running   tcp://52.90.113.128:2376            v1.10.0
    default          *        virtualbox     Running   tcp://192.168.99.100:2376           v1.10.1
    docker-sandbox   -        digitalocean   Running   tcp://104.131.43.236:2376           v1.10.0

You might also run into a similar situation with Docker Universal Control Plan (UCP).

There are a few ways to address this problem and keep using your older machines. One solution is to use a version manager like <a href="https://github.com/getcarina/dvm" target="_blank">DVM</a>.

## How do I uninstall Docker Toolbox?

You might decide that you do not need Toolbox now that you have Docker for Mac,
and want to uninstall it. For details on how to perform a clean uninstall of
Toolbox on the Mac, see [How to uninstall
Toolbox](/toolbox/toolbox_install_mac.md#how-to-uninstall-toolbox) in the
Toolbox Mac topics.
