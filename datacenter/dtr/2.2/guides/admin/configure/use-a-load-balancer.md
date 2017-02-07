---
title: Use a load balancer
description: Learn how to configure a load balancer to balance user requests across multiple Docker Trusted Registry replicas.
keywords: docker, dtr, load balancer
---

Once you’ve joined multiple DTR replicas nodes for high-availability, you can
configure your own load balancer to balance user requests across all replicas.

![](../../images/use-a-load-balancer-1.svg)


This allows users to access DTR using a centralized domain name. If a replica
goes down, the load balancer can detect that and stop forwarding requests to
it, so that the failure goes unnoticed by users.

## Load-balancing DTR

DTR does not provide a load balancing service. You can use an on-premises
or cloud-based load balancer to balance requests across multiple DTR replicas.

Make sure you configure your load balancer to:

* Load-balance TCP traffic on ports 80 and 443
* Not terminate HTTPS connections
* Use the `/health` endpoint on each DTR replica, to check if
the replica is healthy and if it should remain on the load balancing pool or
not

## Where to go next

* [Backups and disaster recovery](backups-and-disaster-recovery.md)
* [DTR architecture](../architecture.md)
