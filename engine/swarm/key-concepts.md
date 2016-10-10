---
description: Introducing key concepts for Docker Engine swarm mode
keywords: docker, container, cluster, swarm mode
title: Swarm mode key concepts
---

This topic introduces some of the concepts unique to the cluster management and
orchestration features of Docker Engine 1.12.

## What is a swarm?

The cluster management and orchestration features embedded in the Docker Engine
are built using **SwarmKit**. Docker engines participating in a cluster are
running in **swarm mode**. You enable swarm mode for an engine by either
initializing a swarm or joining an existing swarm.

A **swarm** is a cluster of Docker engines, or _nodes_, where you deploy
[services](key-concepts.md#services-and-tasks). The Docker Engine CLI and API
include commands to manage swarm nodes (e.g., add or remove nodes), and deploy
and orchestrate services across the swarm.

When you run Docker without using swarm mode, you execute container
commands. When you run the Docker in swarm mode, you orchestrate services. You can run swarm services and standalone containers on the same Docker instances.

## What is a node?

A **node** is an instance of the Docker engine participating in the swarm. You can also think of this as a Docker node. You can run one or more nodes on a single physical computer or cloud server, but production swarm deployments typically include Docker nodes distributed across multiple physical and cloud machines.

To deploy your application to a swarm, you submit a service definition to a
**manager node**. The manager node dispatches units of work called
[tasks](#Services-and-tasks) to worker nodes.

Manager nodes also perform the orchestration and cluster management functions
required to maintain the desired state of the swarm. Manager nodes elect a
single leader to conduct orchestration tasks.

**Worker nodes** receive and execute tasks dispatched from manager nodes.
By default manager nodes also run services as worker nodes, but you can
configure them to run manager tasks exclusively and be manager-only
nodes. An agent runs on each worker node and reports on the tasks assigned to
it. The worker node notifies the manager node of the current state of its
assigned tasks so that the manager can maintain the desired state of each
worker.

## Services and tasks

A **service** is the definition of the tasks to execute on the worker nodes. It
is the central structure of the swarm system and the primary root of user
interaction with the swarm.

When you create a service, you specify which container image to use and which
commands to execute inside running containers.

In the **replicated services** model, the swarm manager distributes a specific
number of replica tasks among the nodes based upon the scale you set in the
desired state.

For **global services**, the swarm runs one task for the service on every
available node in the cluster.

A **task** carries a Docker container and the commands to run inside the
container. It is the atomic scheduling unit of swarm. Manager nodes assign tasks
to worker nodes according to the number of replicas set in the service scale.
Once a task is assigned to a node, it cannot move to another node. It can only
run on the assigned node or fail.

## Load balancing

The swarm manager uses **ingress load balancing** to expose the services you
want to make available externally to the swarm. The swarm manager can
automatically assign the service a **PublishedPort** or you can configure a
PublishedPort for the service. You can specify any unused port. If you do not
specify a port, the swarm manager assigns the service a port in the 30000-32767
range.

External components, such as cloud load balancers, can access the service on the
PublishedPort of any node in the cluster whether or not the node is currently
running the task for the service.  All nodes in the swarm route ingress
connections to a running task instance.

Swarm mode has an internal DNS component that automatically assigns each service
in the swarm a DNS entry. The swarm manager uses **internal load balancing** to
distribute requests among services within the cluster based upon the DNS name of
the service.

## What's next?
* Read the [swarm mode overview](index.md).
* Get started with the [swarm mode tutorial](swarm-tutorial/index.md).
