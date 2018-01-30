---
description: Docker Universal Control plane has support for high availability. Learn how to set up your installation to ensure it tolerates failures.
keywords: docker, ucp, high-availability, replica
redirect_from:
- /ucp/understand_ha/
- /ucp/high-availability/set-up-high-availability/
title: Set up high availability
---

Docker Universal Control Plane is designed for high availability (HA).

When setting up a UCP cluster, you can add additional nodes to serve as
replicas of the controller. In that case, you’ll have multiple nodes, each
running the same set of containers.
[Learn more about the UCP architecture](../architecture.md).

Adding replica nodes to your cluster allows you to:

* Handle controller node failures,
* Load-balance user requests across the controller and replica nodes.


## Size your deployment

To make the cluster tolerant to more failures, add additional replica nodes to
your cluster:

| Controller and replicas | Failures tolerated |
|:-----------------------:|:------------------:|
|            1            |         0          |
|            3            |         1          |
|            5            |         2          |
|            7            |         3          |


When sizing your cluster, follow these rules of thumb:

* Don't create a cluster with just one replica. Your cluster can't tolerate any
failures, and it's possible that you experience performance degradation.
* When a replica fails, the number of failures tolerated by your cluster
decreases. Don't leave that replica offline for long.
* Adding too many replicas to the cluster might also lead to performance
degradation, as changes to configurations need to be replicated across all
replicas.

## Replicating CAs

When configuring UCP for high-availability, you need to ensure the CAs running
on each UCP controller node are interchangeable. This is done by transferring
root certificates and keys for the CAs to each controller node on the cluster.
[Learn how to replicate CAs for high availability](replicate-cas.md)

## Load-balancing on UCP

Docker UCP does not include a load-balancer. You can configure your own
load-balancer to balance user requests across all controller replicas.
[Learn more about the UCP reference architecture](https://www.docker.com/sites/default/files/RA_UCP%20Load%20Balancing-Feb%202016_0.pdf).

Since Docker UCP uses mutual TLS, make sure you configure your load balancer to:

* Load-balance TCP traffic on ports 80 and 443,
* Use a TCP load balancer that doesn't terminate HTTPS connections,
* Use the `/_ping` endpoint on each UCP controller, to check if the controller
is healthy and if it should remain on the load balancing pool or not.


## Where to go next

* [UCP architecture](../architecture.md)
* [Install UCP for production](../installation/install-production.md)