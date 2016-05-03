<!--[metadata]>
+++
aliases = [ "/ucp/understand_ha/"]
title ="Set up high availability"
description="Docker Universal Control plane has support for high availability. Learn how to set up your installation to ensure it tolerates failures."
keywords= ["docker, ucp, high-availability, replica"]
[menu.main]
parent="mn_ucp_high_availability"
identifier="ucp_set_high_availability"
weight=0
+++
<![end-metadata]-->


# Set up high availability

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

* Don't create a cluster with just one replica. Your cluster won't tolerate any
failures, and it's possible that you experience performance degradation.
* When a replica fails, the number of failures tolerated by your cluster
decreases. Don't leave that replica offline for long.
* Adding too many replicas to the cluster might also lead to performance
degradation, as changes to configurations need to be replicated across all
replicas.


## Load-balancing on UCP

Docker UCP does not include a load-balancer. You can configure your own
load-balancer to balance user requests across all controller replicas.
[Learn more about the UCP reference architecture](https://www.docker.com/sites/default/files/RA_UCP%20Load%20Balancing-Feb%202016_0.pdf).

Since Docker UCP uses mutual TLS, make sure you configure your load balancer to:

* Load-balance TCP traffic on ports 80 and 443,
* Not terminate HTTPS connections,
* Use the `/_ping` endpoint on each UCP controller, to check if the controller
is healthy and if it should remain on the load balancing pool or not.


## Where to go next

* [UCP architecture](../architecture.md)
* [Install UCP for production](../installation/install-production.md)
