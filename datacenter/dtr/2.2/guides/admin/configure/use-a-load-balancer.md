---
title: Use a load balancer
description: Learn how to configure a load balancer to balance user requests across multiple Docker Trusted Registry replicas.
keywords: docker, dtr, load balancer
---

Once you’ve joined multiple DTR replicas nodes for
[high-availability](set-up-high-availability.md), you can configure your own
load balancer to balance user requests across all replicas.

![](../../images/use-a-load-balancer-1.svg)


This allows users to access DTR using a centralized domain name. If a replica
goes down, the load balancer can detect that and stop forwarding requests to
it, so that the failure goes unnoticed by users.

## Load balancing DTR

DTR does not provide a load balancing service. You can use an on-premises
or cloud-based load balancer to balance requests across multiple DTR replicas.

Make sure you configure your load balancer to:

* Load balance TCP traffic on ports 80 and 443
* Not terminate HTTPS connections
* Use the `/health` endpoint (note the lack of an `/api/v0/` in the path) on each
DTR replica, to check if the replica is healthy and if it should remain on the
load balancing pool or not

## Health check endpoints

The `/health` endpoint returns a JSON object for the replica being queried with
`"Healthy"` as one of the keys. Any response other than a 200 HTTP status code
and `"Healthy":true` means the replica is unsuitable for taking requests. If
the API server is still up, the returned JSON object will have an `"Error"` key
with more details. More specifically, these issues can be in any of these
services:

* Storage container (registry)
* Authorization (garant)
* Metadata persistence (rethinkdb)
* Content trust (notary)

Note that this endpoint is for checking the health of a *single* replica. To get
the health of every replica in a cluster, querying each individual replica is
the preferred way to do it in real time.

The `/api/v0/meta/cluster_status` endpoint returns a JSON object for the entire
cluster *as observed* by the replica being queried. Health status for the
replicas is available in the `"replica_health"` key. These statuses are taken
from a cache which is updated by each replica individually.

In addition, this endpoint returns a dump of the rethink system tables
which can be rather large (~45 KB) for a status endpoint.


## Where to go next

* [Backups and disaster recovery](../backups-and-disaster-recovery.md)
* [DTR architecture](../../architecture.md)
