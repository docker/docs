---
description: Docker explained in depth
keywords: docker, introduction, documentation, about, technology, understanding
redirect_from:
- /introduction/understanding-docker/
- /engine/userguide/basics/
- /engine/introduction/understanding-docker/
- /engine/understanding-docker/
- /engine/docker-overview/
title: Docker overview
---

Docker is an open platform for developing, shipping, and running applications.
Docker enables you to separate your applications from your infrastructure so
you can deliver software quickly. With Docker, you can manage your infrastructure
in the same ways you manage your applications. By taking advantage of Docker's
methodologies for shipping, testing, and deploying code quickly, you can
significantly reduce the delay between writing code and running it in production.

## The Docker platform

Docker provides the ability to package and run an application in a loosely isolated
environment called a container. The isolation and security allow you to run many
containers simultaneously on a given host. Containers are lightweight because
they don't need the extra load of a hypervisor, but run directly within the host
machine's kernel. This means you can run more containers on a given hardware
combination than if you were using virtual machines. You can even run Docker
containers within host machines that are actually virtual machines!

Docker provides tooling and a platform to manage the lifecycle of your containers:

* Develop your application and its supporting components using containers.
* The container becomes the unit for distributing and testing your application.
* When you're ready, deploy your application into your production environment,
  as a container or an orchestrated service. This works the same whether your
  production environment is a local data center, a cloud provider, or a hybrid
  of the two.

## Docker Engine

_Docker Engine_ is a client-server application with these major components:

* A server which is a type of long-running program called a daemon process (the
  `dockerd` command).

* A REST API which specifies interfaces that programs can use to talk to the
  daemon and instruct it what to do.

* A command line interface (CLI) client (the `docker` command).

![Docker Engine Components Flow](/engine/images/engine-components-flow.png)

The CLI uses the Docker REST API to control or interact with the Docker daemon
through scripting or direct CLI commands. Many other Docker applications use the
underlying API and CLI.

The daemon creates and manages Docker _objects_, such as images, containers,
networks, and volumes.

> **Note**: Docker is licensed under the open source Apache 2.0 license.

For more details, see [Docker Architecture](#docker-architecture) below.

## What can I use Docker for?

**Fast, consistent delivery of your applications**

Docker streamlines the development lifecycle by allowing developers to work in
standardized environments using local containers which provide your applications
and services. Containers are great for continuous integration and continuous
delivery (CI/CD) workflows.

Consider the following example scenario:

- Your developers write code locally and share their work with their colleagues
  using Docker containers.
- They use Docker to push their applications into a test environment and execute
  automated and manual tests.
- When developers find bugs, they can fix them in the development environment
  and redeploy them to the test environment for testing and validation.
- When testing is complete, getting the fix to the customer is as simple as
  pushing the updated image to the production environment.

**Responsive deployment and scaling**

Docker's container-based platform allows for highly portable workloads. Docker
containers can run on a developer's local laptop, on physical or virtual
machines in a data center, on cloud providers, or in a mixture of environments.

Docker's portability and lightweight nature also make it easy to dynamically
manage workloads, scaling up or tearing down applications and services as
business needs dictate, in near real time.

**Running more workloads on the same hardware**

Docker is lightweight and fast. It provides a viable, cost-effective alternative
to hypervisor-based virtual machines, so you can use more of your compute
capacity to achieve your business goals. Docker is perfect for high density
environments and for small and medium deployments where you need to do more with
fewer resources.

## Docker architecture

Docker uses a client-server architecture. The Docker *client* talks to the
Docker *daemon*, which does the heavy lifting of building, running, and
distributing your Docker containers. The Docker client and daemon *can*
run on the same system, or you can connect a Docker client to a remote Docker
daemon. The Docker client and daemon communicate using a REST API, over UNIX
sockets or a network interface.

![Docker Architecture Diagram](/engine/images/architecture.svg)

### The Docker daemon

The Docker daemon (`dockerd`) listens for Docker API requests and manages Docker
objects such as images, containers, networks, and volumes. A daemon can also
communicate with other daemons to manage Docker services.

### The Docker client

The Docker client (`docker`) is the primary way that many Docker users interact
with Docker. When you use commands such as `docker run`, the client sends these
commands to `dockerd`, which carries them out. The `docker` command uses the
Docker API. The Docker client can communicate with more than one daemon.

### Docker registries

A Docker _registry_ stores Docker images. Docker Hub is a public
registry that anyone can use, and Docker is configured to look for images on
Docker Hub by default. You can even run your own private registry.

When you use the `docker pull` or `docker run` commands, the required images are
pulled from your configured registry. When you use the `docker push` command,
your image is pushed to your configured registry.

### Docker objects

When you use Docker, you are creating and using images, containers, networks,
volumes, plugins, and other objects. This section is a brief overview of some
of those objects.

#### Images

An _image_ is a read-only template with instructions for creating a Docker
container. Often, an image is _based on_ another image, with some additional
customization. For example, you may build an image which is based on the `ubuntu`
image, but installs the Apache web server and your application, as well as the
configuration details needed to make your application run.

You might create your own images or you might only use those created by others
and published in a registry. To build your own image, you create a _Dockerfile_
with a simple syntax for defining the steps needed to create the image and run
it. Each instruction in a Dockerfile creates a layer in the image. When you
change the Dockerfile and rebuild the image, only those layers which have
changed are rebuilt. This is part of what makes images so lightweight, small,
and fast, when compared to other virtualization technologies.

#### Containers

A container is a runnable instance of an image. You can create, start, stop,
move, or delete a container using the Docker API or CLI. You can connect a
container to one or more networks, attach storage to it, or even create a new
image based on its current state.

By default, a container is relatively well isolated from other containers and
its host machine. You can control how isolated a container's network, storage,
or other underlying subsystems are from other containers or from the host
machine.

A container is defined by its image as well as any configuration options you
provide to it when you create or start it. When a container is removed, any changes to
its state that are not stored in persistent storage disappear.

##### Example `docker run` command

The following command runs an `ubuntu` container, attaches interactively to your
local command-line session, and runs `/bin/bash`.

```bash
$ docker run -i -t ubuntu /bin/bash
```

When you run this command, the following happens (assuming you are using
the default registry configuration):

1.  If you do not have the `ubuntu` image locally, Docker pulls it from your
    configured registry, as though you had run `docker pull ubuntu` manually.

2.  Docker creates a new container, as though you had run a `docker container create`
    command manually.

3.  Docker allocates a read-write filesystem to the container, as its final
    layer. This allows a running container to create or modify files and
    directories in its local filesystem.

4.  Docker creates a network interface to connect the container to the default
    network, since you did not specify any networking options. This includes
    assigning an IP address to the container. By default, containers can
    connect to external networks using the host machine's network connection.

5.  Docker starts the container and executes `/bin/bash`. Because the container
    is running interactively and attached to your terminal (due to the `-i` and `-t`
    flags), you can provide input using your keyboard while the output is logged to
    your terminal.

6.  When you type `exit` to terminate the `/bin/bash` command, the container
    stops but is not removed. You can start it again or remove it.

#### Services

Services allow you to scale containers across multiple Docker daemons, which
all work together as a _swarm_ with multiple _managers_ and _workers_. Each
member of a swarm is a Docker daemon, and all the daemons communicate using
the Docker API. A service allows you to define the desired state, such as the
number of replicas of the service that must be available at any given time.
By default, the service is load-balanced across all worker nodes. To
the consumer, the Docker service appears to be a single application. Docker
Engine supports swarm mode in Docker 1.12 and higher.

## The underlying technology
Docker is written in [Go](https://golang.org/) and takes advantage of several
features of the Linux kernel to deliver its functionality.

### Namespaces
Docker uses a technology called `namespaces` to provide the isolated workspace
called the *container*. When you run a container, Docker creates a set of
*namespaces* for that container.

These namespaces provide a layer of isolation. Each aspect of a container runs
in a separate namespace and its access is limited to that namespace.

Docker Engine uses namespaces such as the following on Linux:

 - **The `pid` namespace:** Process isolation (PID: Process ID).
 - **The `net` namespace:** Managing network interfaces (NET:
 Networking).
 - **The `ipc` namespace:** Managing access to IPC
 resources (IPC: InterProcess Communication).
 - **The `mnt` namespace:** Managing filesystem mount points (MNT: Mount).
 - **The `uts` namespace:** Isolating kernel and version identifiers. (UTS: Unix
Timesharing System).

### Control groups
Docker Engine on Linux also relies on another technology called _control groups_
(`cgroups`). A cgroup limits an application to a specific set of resources.
Control groups allow Docker Engine to share available hardware resources to
containers and optionally enforce limits and constraints. For example,
you can limit the memory available to a specific container.

### Union file systems
Union file systems, or UnionFS, are file systems that operate by creating layers,
making them very lightweight and fast. Docker Engine uses UnionFS to provide
the building blocks for containers. Docker Engine can use multiple UnionFS variants,
including AUFS, btrfs, vfs, and DeviceMapper.

### Container format
Docker Engine combines the namespaces, control groups, and UnionFS into a wrapper
called a container format. The default container format is `libcontainer`. In
the future, Docker may support other container formats by integrating with
technologies such as BSD Jails or Solaris Zones.

## Next steps
- Read about [installing Docker](../get-docker.md).
- Get hands-on experience with the [Getting started with Docker](index.md)
    tutorial.
