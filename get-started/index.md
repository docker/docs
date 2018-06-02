---
title: "Get Started, Part 1: Orientation and setup"
keywords: get started, setup, orientation, quickstart, intro, concepts, containers
description: Get oriented on some basics of Docker before diving into the walkthrough.
redirect_from:
- /getstarted/
- /get-started/part1/
- /engine/getstarted/
- /learn/
- /engine/getstarted/step_one/
- /engine/getstarted/step_two/
- /engine/getstarted/step_three/
- /engine/getstarted/step_four/
- /engine/getstarted/step_five/
- /engine/getstarted/step_six/
- /engine/getstarted/last_page/
- /engine/getstarted-voting-app/
- /engine/getstarted-voting-app/node-setup/
- /engine/getstarted-voting-app/create-swarm/
- /engine/getstarted-voting-app/deploy-app/
- /engine/getstarted-voting-app/test-drive/
- /engine/getstarted-voting-app/customize-app/
- /engine/getstarted-voting-app/cleanup/
- /engine/userguide/intro/
- /mac/started/
- /windows/started/
- /linux/started/
- /getting-started/
- /mac/step_one/
- /windows/step_one/
- /linux/step_one/
- /engine/tutorials/dockerizing/
- /mac/step_two/
- /windows/step_two/
- /linux/step_two/
- /mac/step_three/
- /windows/step_three/
- /linux/step_three/
- /engine/tutorials/usingdocker/
- /mac/step_four/
- /windows/step_four/
- /linux/step_four/
- /engine/tutorials/dockerimages/
- /userguide/dockerimages/
- /engine/userguide/dockerimages/
- /mac/last_page/
- /windows/last_page/
- /linux/last_page/
- /mac/step_six/
- /windows/step_six/
- /linux/step_six/
- /engine/tutorials/dockerrepos/
- /userguide/dockerrepos/
- /engine/userguide/containers/dockerimages/
---

{% include_relative nav.html selected="1" %}

Welcome! We are excited that you want to learn Docker. The _Docker Get Started Tutorial_
teaches you how to:

1. Set up your Docker environment (on this page)
2. [Build an image and run it as one container](part2.md)
3. [Scale your app to run multiple containers](part3.md)
4. [Distribute your app across a cluster](part4.md)
5. [Stack services by adding a backend database](part5.md)
6. [Deploy your app to production](part6.md)

## Docker concepts

Docker is a platform for developers and sysadmins to **develop, deploy, and run**
applications with containers. The use of Linux containers to deploy applications
is called _containerization_. Containers are not new, but their use for easily
deploying applications is.

Containerization is increasingly popular because containers are:

- Flexible: Even the most complex applications can be containerized.
- Lightweight: Containers leverage and share the host kernel.
- Interchangeable: You can deploy updates and upgrades on-the-fly.
- Portable: You can build locally, deploy to the cloud, and run anywhere.
- Scalable: You can increase and automatically distribute container replicas.
- Stackable: You can stack services vertically and on-the-fly.

![Containers are portable](images/laurel-docker-containers.png){:width="300px"}

### Images and containers

A container is launched by running an image. An **image** is an executable
package that includes everything needed to run an application--the code, a
runtime, libraries, environment variables, and configuration files.

A **container** is a runtime instance of an image--what the image becomes in
memory when executed (that is, an image with state, or a user process). You can
see a list of your running containers with the command, `docker ps`, just as you
would in Linux.

### Containers and virtual machines

A **container** runs _natively_ on Linux and shares the kernel of the host
machine with other containers. It runs a discrete process, taking no more memory
than any other executable, making it lightweight.

By contrast, a **virtual machine** (VM) runs a full-blown "guest" operating
system with _virtual_ access to host resources through a hypervisor. In general,
VMs provide an environment with more resources than most applications need.

![Container stack example](https://www.docker.com/sites/default/files/Container%402x.png){:width="300px"} | ![Virtual machine stack example](https://www.docker.com/sites/default/files/VM%402x.png){:width="300px"}

## Prepare your Docker environment

Install a [maintained version](https://docs.docker.com/engine/installation/#updates-and-patches){: target="_blank" class="_"}
of Docker Community Edition (CE) or Enterprise Edition (EE) on a
[supported platform](https://docs.docker.com/engine/installation/#supported-platforms){: target="_blank" class="_"}.

> For full Kubernetes Integration
>
> - [Kubernetes on Docker for Mac](/docker-for-mac/kubernetes/){: target="_blank" class="_"}
is available in [17.12 Edge (mac45)](/docker-for-mac/edge-release-notes/#docker-community-edition-17120-ce-mac45-2018-01-05){: target="_blank" class="_"} or
[17.12 Stable (mac46)](/docker-for-mac/release-notes/#docker-community-edition-17120-ce-mac46-2018-01-09){: target="_blank" class="_"} and higher.
> - [Kubernetes on Docker for Windows](/docker-for-windows/kubernetes/){: target="_blank" class="_"}
is available in
[18.02 Edge (win50)](/docker-for-windows/edge-release-notes/#docker-community-edition-18020-ce-rc1-win50-2018-01-26){: target="_blank" class="_"} and higher edge channels only. 

[Install Docker](/engine/installation/index.md){: class="button outline-btn"}
<div style="clear:left"></div>

### Test Docker version

1.  Run `docker --version` and ensure that you have a supported version of Docker:

    ```shell
    docker --version

    Docker version 17.12.0-ce, build c97c6d6
    ```

2.  Run `docker info` or (`docker version` without `--`) to view even more details about your docker installation:

    ```shell
    docker info

    Containers: 0
     Running: 0
     Paused: 0
     Stopped: 0
    Images: 0
    Server Version: 17.12.0-ce
    Storage Driver: overlay2
    ...
    ```

> To avoid permission errors (and the use of `sudo`), add your user to the `docker` group. [Read more](https://docs.docker.com/engine/installation/linux/linux-postinstall/){: target="_blank" class="_"}.

### Test Docker installation

1.  Test that your installation works by running the simple Docker image,
[hello-world](https://hub.docker.com/_/hello-world/){: target="_blank" class="_"}:

    ```shell
    docker run hello-world

    Unable to find image 'hello-world:latest' locally
    latest: Pulling from library/hello-world
    ca4f61b1923c: Pull complete
    Digest: sha256:ca0eeb6fb05351dfc8759c20733c91def84cb8007aa89a5bf606bc8b315b9fc7
    Status: Downloaded newer image for hello-world:latest

    Hello from Docker!
    This message shows that your installation appears to be working correctly.
    ...
    ```

2.  List the `hello-world` image that was downloaded to your machine:

    ```shell
    docker image ls
    ```

3.  List the `hello-world` container (spawned by the image) which exits after
    displaying its message. If it were still running, you would not need the `--all` option:

    ```shell
    docker container ls --all

    CONTAINER ID     IMAGE           COMMAND      CREATED            STATUS
    54f4984ed6a8     hello-world     "/hello"     20 seconds ago     Exited (0) 19 seconds ago
    ```

## Recap and cheat sheet

```shell
## List Docker CLI commands
docker
docker container --help

## Display Docker version and info
docker --version
docker version
docker info

## Execute Docker image
docker run hello-world

## List Docker images
docker image ls

## List Docker containers (running, all, all in quiet mode)
docker container ls
docker container ls --all
docker container ls -aq
```

## Conclusion of part one

Containerization makes [CI/CD](https://www.docker.com/use-cases/cicd){: target="_blank" class="_"} seamless. For example:

- applications have no system dependencies
- updates can be pushed to any part of a distributed application
- resource density can be optimized.

With Docker, scaling your application is a matter of spinning up new
executables, not running heavy VM hosts.

[On to Part 2 >>](part2.md){: class="button outline-btn" style="margin-bottom: 30px; margin-right: 100%"}
