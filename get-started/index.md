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

Welcome! We are excited you want to learn how to use Docker.

In this six-part tutorial, you will:

1. Get set up and oriented, on this page.
2. [Build and run your first app](part2.md)
3. [Turn your app into a scaling service](part3.md)
4. [Span your service across multiple machines](part4.md)
5. [Add a visitor counter that persists data](part5.md)
6. [Deploy your swarm to production](part6.md)

The application itself is very simple so that you are not too distracted by
what the code is doing. After all, the value of Docker is in how it can build,
ship, and run applications; it's totally agnostic as to what your application
actually does.

## Prerequisites

While we'll define concepts along the way, it is good for you to understand
[what Docker is](https://www.docker.com/what-docker) and [why you would use
Docker](https://www.docker.com/use-cases) before we begin.

We also need to assume you are familiar with a few concepts before we continue:

- IP Addresses and Ports
- Virtual Machines
- Editing configuration files
- Basic familiarity with the ideas of code dependencies and building
- Machine resource usage terms, like CPU percentages, RAM use in bytes, etc.

## A brief explanation of containers

An **image** is a lightweight, stand-alone, executable package that includes
everything needed to run a piece of software, including the code, a runtime,
libraries, environment variables, and config files.

A **container** is a runtime instance of an image&#8212;what the image becomes
in memory when actually executed. It runs completely isolated from the host
environment by default, only accessing host files and ports if configured to do
so.

Containers run apps natively on the host machine's kernel. They have better
performance characteristics than virtual machines that only get virtual access
to host resources through a hypervisor. Containers can get native access, each
one running in a discrete process, taking no more memory than any other
executable.

## Containers vs. virtual machines

Consider this diagram comparing virtual machines to containers:

### Virtual Machine diagram

![Virtual machine stack example](https://www.docker.com/sites/default/files/VM%402x.png)

Virtual machines run guest operating systems&#8212;note the OS layer in each
box. This is resource intensive, and the resulting disk image and application
state is an entanglement of OS settings, system-installed dependencies, OS
security patches, and other easy-to-lose, hard-to-replicate ephemera.

### Container diagram

![Container stack example](https://www.docker.com/sites/default/files/Container%402x.png)

Containers can share a single kernel, and the only information that needs to be
in a container image is the executable and its package dependencies, which never
need to be installed on the host system. These processes run like native
processes, and you can manage them individually by running commands like `docker
ps`&#8212;just like you would run `ps` on Linux to see active processes.
Finally, because they contain all their dependencies, there is no configuration
entanglement; a containerized app "runs anywhere."

## Setup

Before we get started, make sure your system has the latest version of Docker
installed.

[Install Docker](/engine/installation/index.md){: class="button outline-btn"}
<div style="clear:left"></div>
> **Note**: version 1.13 or higher is required

You should be able to run `docker run hello-world` and see a response like this:
> **Note**: You may need to add your user to the `docker` group in order to call this command without sudo. [Read more](https://docs.docker.com/engine/installation/linux/linux-postinstall/)
> **Note**: If there are networking issues in your setup, `docker run hello-world` may fail to execute successfully. In case you are behind a proxy server and you suspect that it blocks the connection, check the [next part] (https://docs.docker.com/get-started/part2/) of the tutorial.

```shell
$ docker run hello-world

Hello from Docker!
This message shows that your installation appears to be working correctly.

To generate this message, Docker took the following steps:
...(snipped)...
```

Now would also be a good time to make sure you are using version 1.13 or higher. Run `docker --version` to check it out.

```shell
$ docker --version
Docker version 17.05.0-ce-rc1, build 2878a85
```

If you see messages like the ones above, you are ready to begin your journey.

## Conclusion

The unit of scale being an individual, portable executable has vast
implications. It means CI/CD can push updates to any part of a distributed
application, system dependencies are not an issue, and resource density is
increased. Orchestration of scaling behavior is a matter of spinning up new
executables, not new VM hosts.

We'll be learning about all of these things, but first let's learn to walk.

[On to Part 2 >>](part2.md){: class="button outline-btn" style="margin-bottom: 30px; margin-right: 100%"}
