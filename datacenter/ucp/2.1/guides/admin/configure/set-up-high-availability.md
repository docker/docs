---
description: Docker Universal Control plane has support for high availability. Learn
  how to set up your installation to ensure it tolerates failures.
keywords: docker, ucp, high-availability, replica
title: Set up high availability
---

Docker Universal Control Plane is designed for high availability (HA). You can
join multiple manager nodes to the cluster, so that if one manager node fails,
another can automatically take its place without impact to the cluster.

Having multiple manager nodes in your cluster, allows you to:

* Handle manager node failures,
* Load-balance user requests across all manager nodes.

## Size your deployment

To make the cluster tolerant to more failures, add additional replica nodes to
your cluster.

| Manager nodes | Failures tolerated |
|:-------------:|:------------------:|
|       1       |         0          |
|       3       |         1          |
|       5       |         2          |
|       7       |         3          |


For production-grade deployments, follow these rules of thumb:

* When a manager node fails, the number of failures tolerated by your cluster
decreases. Don't leave that node offline for too long.
* You should distribute your manager nodes across different availability zones.
This way your cluster can continue working even if an entire availability zone
goes down.
* Adding many manager nodes to the cluster might lead to performance
degradation, as changes to configurations need to be replicated across all
manager nodes. The maximum advisable is having 7 manager nodes.


## Where to go next

* [Scale your cluster](scale-your-cluster.md)
* [Use a load balancer](use-a-load-balancer.md)
