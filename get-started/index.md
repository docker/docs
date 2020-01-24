---
title: "Orientation and setup"
keywords: get started, setup, orientation, quickstart, intro, concepts, containers, docker desktop
description: Get oriented on some basics of Docker and install Docker Desktop.
redirect_from:
- /getstarted/
- /get-started/part1/
- /get-started/part6/
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

Welcome! We are excited that you want to learn Docker. The _Docker Community QuickStart_
teaches you how to:

1. Set up your Docker environment (on this page)
2. [Build an image and run it as one container](part2.md)
3. [Share your containerized applications on Docker Hub](part3.md)

## Docker concepts

Docker is a platform for developers and sysadmins to **build, share, and run**
applications with containers. The use of containers to deploy applications
is called _containerization_. Containers are not new, but their use for easily
deploying applications is.

Containerization is increasingly popular because containers are:

- Flexible: Even the most complex applications can be containerized.
- Lightweight: Containers leverage and share the host kernel,
  making them much more efficient in terms of system resources than virtual machines.
- Portable: You can build locally, deploy to the cloud, and run anywhere.
- Loosely coupled: Containers are highly self sufficient and encapsulated,
  allowing you to replace or upgrade one without disrupting others.
- Scalable: You can increase and automatically distribute container replicas across a datacenter.
- Secure: Containers apply aggressive constraints and isolations to processes without
  any configuration required on the part of the user.

### Images and containers

Fundamentally, a container is nothing but a running process,
with some added encapsulation features applied to it in order to keep it isolated from the host
and from other containers.
One of the most important aspects of container isolation is that each container interacts
with its own, private filesystem; this filesystem is provided by a Docker **image**.
An image includes everything needed to run an application -- the code or binary,
runtimes, dependencies, and any other filesystem objects required.

### Containers and virtual machines

A container runs _natively_ on Linux and shares the kernel of the host
machine with other containers. It runs a discrete process, taking no more memory
than any other executable, making it lightweight.

By contrast, a **virtual machine** (VM) runs a full-blown "guest" operating
system with _virtual_ access to host resources through a hypervisor. In general,
VMs incur a lot of overhead beyond what is being consumed by your application logic.

![Container stack example](/images/Container%402x.png){:width="300px"} | ![Virtual machine stack example](/images/VM%402x.png){:width="300px"}

### Orchestration

The portability and reproducibility of a containerized process mean we have 
an opportunity to move and scale our containerized applications across 
clouds and datacenters; containers effectively guarantee that those 
applications will run the same way anywhere, allowing us to quickly and 
easily take advantage of all these environments. 
Furthermore, as we scale our applications up, we'll 
want some tooling to help automate the maintenance of those applications, 
able to replace failed containers automatically and manage the rollout of updates 
and reconfigurations of those containers during their lifecycle. 

Tools to manage, scale, and maintain containerized applications are called 
_orchestrators_, and the most common examples of these are _Kubernetes_ and 
_Docker Swarm_. Development environment deployments of both of these 
orchestrators are provided by Docker Desktop, which we'll use throughout 
this guide to create our first orchestrated, containerized application.

## Install Docker Desktop

The best way to get started developing containerized applications is with Docker Desktop, for OSX or Windows. Docker Desktop will allow you to easily set up Kubernetes or Swarm on your local development machine, so you can use all the features of the orchestrator you're developing applications for right away, no cluster required. Follow the installation instructions appropriate for your operating system:

 - [OSX](/docker-for-mac/install/){: target="_blank" class="_"}
 - [Windows](/docker-for-windows/install/){: target="_blank" class="_"}


## Conclusion

At this point, you've installed Docker Desktop on your development machine, and confirmed that you can run simple containerized workloads in Kubernetes and Swarm. In the next section, we'll start developing our first containerized application.

[On to Part 2 >>](part2.md){: class="button outline-btn" style="margin-bottom: 30px; margin-right: 100%"}

## CLI References

Further documentation for all CLI commands used in this article are available here:

- [`kubectl apply`](https://kubernetes.io/docs/reference/generated/kubectl/kubectl-commands#apply)
- [`kubectl get`](https://kubernetes.io/docs/reference/generated/kubectl/kubectl-commands#get)
- [`kubectl logs`](https://kubernetes.io/docs/reference/generated/kubectl/kubectl-commands#logs)
- [`kubectl delete`](https://kubernetes.io/docs/reference/generated/kubectl/kubectl-commands#delete)
- [`docker swarm init`](https://docs.docker.com/engine/reference/commandline/swarm_init/)
- [`docker service *`](https://docs.docker.com/engine/reference/commandline/service/)
