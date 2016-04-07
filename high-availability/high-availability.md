<!--[metadata]>
+++
title = "Set up high availability"
description = "Learn how to set up Docker Trusted Registry for high availability."
keywords = ["docker, registry, high-availability, backup, recovery"]
[menu.main]
parent="dtr_menu_high_availability"
identifier="dtr_high_availability"
weight=0
+++
<![end-metadata]-->

# Set up high availability

Docker Trusted Registry (DTR) is designed for high availability.
When installing DTR you can add multiple nodes to form a cluster.

Adding more nodes to your DTR cluster allows you to:

* Load-balance user requests across the DTR nodes,
* Keep the DTR cluster working if a node fails.

To make a DTR installation tolerant to node failures, add additional nodes to
the DTR cluster.

| DTR nodes | Failures tolerated |
|:---------:|:------------------:|
|     1     |         0          |
|     3     |         1          |
|     5     |         2          |
|     7     |         3          |

When sizing your DTR installation for high-availability,
follow these rules of thumb:

* Don't create a DTR cluster with just two nodes. Your cluster
won't tolerate any failures, and it's possible that you experience performance
degradation.
* When a node fails, the number of failures tolerated by your cluster
decreases. Don't leave that node offline for long.
* Adding too many nodes to the cluster might also lead to performance
degradation, as data needs to be replicated across all nodes.

## Size your cluster

When installing DTR for production, you should have separate nodes for running
Docker Universal Control Plane (DTR), Docker Trusted Registry, and your
containers.

Having dedicated nodes for UCP, DTR, and your containers, ensures they stay
performant since all applications have dedicated resources.
It also makes it easier to implement backup policies and disaster recovery
plans.

For installing DTR for production, you'll need a minimum of:

* 3 dedicated nodes to install UCP for high-availability,
* 3 dedicated nodes to install DTR for high-availability,
* As many nodes as you want for running your containers and applications.

<!-- TODO: add diagram to illustrate this -->

![](../images/architecture-3.png)
