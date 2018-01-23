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

Welcome! We are excited that you want to learn Docker. Docker lets you easily
**build, ship, and run** any application on any platform. Let's get started.

The Docker _Get Started_ tutorial teaches you how to:

1. Set up your Docker environment (on this page)
2. [Build an image and run it as one container](part2.md)
3. [Scale your app to run multiple containers](part3.md)
4. [Distribute your app across a cluster](part4.md)
5. [Stack services by adding a backend database](part5.md)
6. [Deploy your app to production](part6.md)

## Images and containers

An **image** is a lightweight, stand-alone, executable package that includes
everything needed to run a piece of software--the code, a runtime, libraries,
environment variables, and configuration files. You can find thousands of images
on [Docker Hub](https://hub.docker.com/explore/){: target="_blank" class="_"}.

A **container** is a runtime instance of an image--what the image becomes in
memory when executed. By default, a container runs completely isolated from the
host environment, only accessing host files and ports if configured to do so.

The use of Linux containers to build, deploy, and run applications is called
[containerization](https://en.wikipedia.org/wiki/Operating-system-level_virtualization){: target="_blank" class="_"}.
Containers are not new but their use for easily deploying applications is.

## Containers vs. virtual machines

A **virtual machine** (VM) runs a "guest" operating system (OS) with virtual (as
opposed to native) access to host resources through a hypervisor. VMs are
resource intensive, and the resulting disk image and application state is an
entanglement of OS settings, system-installed dependencies, OS security patches,
and other ephemera. VMs provide an environment with more resources than most
applications need.

By contrast, a **container** runs applications natively and shares the kernel
of the host machine with other containers, keeping it lightweight. Each
container runs in a discrete process, taking no more memory than any other
executable. The only information needed in a container is the executable and its
package dependencies, which never need to be installed on the host system.
Because containers house their dependencies, a containerized app is portable and
“runs anywhere.”

![Container stack example](https://www.docker.com/sites/default/files/Container%402x.png){:width="300px"} | ![Virtual machine stack example](https://www.docker.com/sites/default/files/VM%402x.png){:width="300px"}

## Prepare your Docker environment

Install a [maintained version](https://docs.docker.com/engine/installation/#updates-and-patches){: target="_blank" class="_"}
of Docker Community Edition (CE) or Enterprise Edition (EE) on a
[supported platform](https://docs.docker.com/engine/installation/#supported-platforms){: target="_blank" class="_"}.

> For full
[Kubernetes integration on Docker for Mac](https://docs.docker.com/docker-for-mac/kubernetes/){: target="_blank" class="_"},
install
[17.12.0-ce Edge](https://docs.docker.com/docker-for-mac/release-notes/#docker-community-edition-17120-ce-mac45-2018-01-05-edge){: target="_blank" class="_"}
or higher.

[Install Docker](/engine/installation/index.md){: class="button outline-btn"}
<div style="clear:left"></div>

### Test Docker version

Ensure that you have a supported version of Docker:

```shell
$ docker --version
Docker version 17.12.0-ce, build c97c6d6
```

Run `docker version`(without `--`) or `docker info` to view even more details
about your docker installation:

```shell
$ docker info
Containers: 0
 Running: 0
 Paused: 0
 Stopped: 0
Images: 0
Server Version: 17.12.0-ce
Storage Driver: overlay2
...
```

> **Note**: To avoid permission errors (and the use of `sudo`), add your user to
> the `docker` group. [Read more](https://docs.docker.com/engine/installation/linux/linux-postinstall/){: target="_blank" class="_"}.

### Test Docker installation

Test that your installation works by running the simple Docker image,
[hello-world](https://hub.docker.com/_/hello-world/){: target="_blank" class="_"}:

```shell
$ docker run hello-world

Unable to find image 'hello-world:latest' locally
latest: Pulling from library/hello-world
ca4f61b1923c: Pull complete
Digest: sha256:ca0eeb6fb05351dfc8759c20733c91def84cb8007aa89a5bf606bc8b315b9fc7
Status: Downloaded newer image for hello-world:latest

Hello from Docker!
This message shows that your installation appears to be working correctly.
...
```

List the `hello-world` image that was downloaded to your machine:

```shell
$ docker image ls
```

List the `hello-world` container (spawned by the image), which exited after
displaying its message. If it were still running, you would _not_ need the
`--all` option:

```shell
$ docker container ls --all
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

## Excecute Docker image
docker run hello-world

## List Docker images
docker image ls

## List Docker containers (running, all, all in quiet mode)
docker container ls
docker container ls -all
docker container ls -a -q
```

## Conclusion of part one

Container images are portable executables that make
[CI/CD](https://www.docker.com/use-cases/cicd){: target="_blank" class="_"} seamless.
For example:

- applications have no system dependencies
- updates can be pushed to any part of a distributed application
- resource density can be optimized.

With Docker, scaling your application is a matter of spinning up new
executables, not running heavy VM hosts.

[On to Part 2 >>](part2.md){: class="button outline-btn" style="margin-bottom: 30px; margin-right: 100%"}
