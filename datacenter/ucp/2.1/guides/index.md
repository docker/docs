---
description: Learn about Docker Universal Control Plane, the enterprise-grade cluster
  management solution from Docker.
keywords: docker, ucp, overview, orchestration, clustering
title: Universal Control Plane overview
---

Docker Universal Control Plane (UCP) is the enterprise-grade cluster management
solution from Docker. You install it on-premises or in your virtual private
cloud, and it helps you manage your Docker cluster and applications from a
single place.

![](images/overview-1.png){: .with-border}

## Centralized cluster management

With Docker you can join up to thousands of physical or virtual machines
together to create a container cluster, allowing you to deploy your applications
at scale. Docker Universal Control Plane extends the functionality provided
by Docker to make it easier to manage your cluster from a centralized place.

You can manage and monitor your container cluster using a graphical UI.

![](images/overview-2.png){: .with-border}

Since UCP exposes the standard Docker API, you can continue using the tools
you already know, including the Docker CLI client, to deploy and manage your
applications.

As an example, you can use the `docker info` command to check the
status of a Docker cluster managed by UCP:

```bash
$ docker info

Containers: 30
Images: 24
Server Version: ucp/2.0.1
Role: primary
Strategy: spread
Filters: health, port, containerslots, dependency, affinity, constraint
Nodes: 2
  ucp-node-1: 192.168.99.100:12376
    └ Status: Healthy
    └ Containers: 20
  ucp-node-2: 192.168.99.101:12376
    └ Status: Healthy
    └ Containers: 10
```

## Deploy, manage, and monitor

With Docker UCP you can manage from a centralized place all the computing
resources you have available like nodes, volumes, and networks.

You can also deploy and monitor your applications and services.

## Built-in security and access control

Docker UCP has its own built-in authentication mechanism and integrates with
LDAP services. It also has Role Based Access Control (RBAC), so that you can
control who can access and make changes to your cluster and applications.

![](images/overview-3.png){: .with-border}

Docker UCP integrates with Docker Trusted Registry so that you can keep the
Docker images you use for your applications behind your firewall, where they
are safe and can't be tampered.

You can also enforce security policies and only allow running applications
that use Docker images you know and trust.

## Where to go next

* [UCP architecture](architecture.md)
* [Install UCP](admin/install/index.md)
