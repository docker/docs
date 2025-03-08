---
description: Learn how to set up Docker Trusted Registry for high availability.
keywords: docker, registry, high-availability, backup, recovery
redirect_from:
- /docker-trusted-registry/high-availability/high-availability/
title: Set up high availability
---

Docker Trusted Registry (DTR) is designed for high availability.

When you first install DTR, you create a cluster with a single DTR replica.
Replicas are single instances of DTR that can be joined together to form a
cluster.
When joining new replicas to the cluster, you create new DTR instances
that are running the same set of services. Any change to the state of an
instance is replicated across all other instances.

![](../images/high-availability-1.png)

Having a DTR cluster with multiple replicas, allows you to:

* Load-balance user requests across the DTR replicas,
* Keep the DTR cluster working if a replica fails.

To make a DTR installation tolerant to failures, add additional replicas to
the DTR cluster.

| DTR replicas | Failures tolerated |
|:------------:|:------------------:|
|      1       |         0          |
|      3       |         1          |
|      5       |         2          |
|      7       |         3          |


When sizing your DTR installation for high-availability,
follow these rules of thumb:

* Don't create a DTR cluster with just two replicas. Your cluster
won't tolerate any failures, and it's possible that you experience performance
degradation.
* When a replica fails, the number of failures tolerated by your cluster
decreases. Don't leave that replica offline for long.
* Adding too many replicas to the cluster might also lead to performance
degradation, as data needs to be replicated across all replicas.

## Size your cluster

When installing DTR for production, you should have dedicated nodes for running:

* Docker Universal Control Plane (UCP) controller nodes,
* DTR replicas,
* Your own applications.

Having dedicated nodes ensures all applications will remain performant since
they have dedicated resources for them.
It also makes it easier to implement backup policies and disaster recovery
plans for UCP and DTR.

To have high-availability on UCP and DTR, you need a minimum of:

* 3 dedicated nodes to install UCP with high availability,
* 3 dedicated nodes to install DTR with high availability,
* As many nodes as you want for running your containers and applications.


![](../images/high-availability-2.png)

## Load balancing

DTR does not provide a load balancing service. You can use an on-premises
or cloud-based load balancer to balance requests across multiple DTR replicas.

Make sure you configure your load balancer to:

* Load-balance TCP traffic on ports 80 and 443,
* Use a TCP load balancer that doesn't terminate HTTPS connections,
* Use the `/load_balancer_status` endpoint on each DTR replica, to check if
the replica is healthy and if it should remain on the load balancing pool or
not.

## Where to go next

* [Backups and disaster recovery](backups-and-disaster-recovery.md)
* [DTR architecture](../architecture.md)
