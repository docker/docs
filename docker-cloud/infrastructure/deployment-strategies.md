---
description: Schedule a deployment
keywords: schedule, deployment, container
redirect_from:
- /docker-cloud/feature-reference/deployment-strategies/
title: Container distribution strategies
---

Docker Cloud can use different distribution strategies when deploying containers
to more than one node. You can use different strategies to change how your
service distributes new containers when scaling.

## Set a deployment distribution strategy

You can set the deployment strategy when creating a service, either through the
Docker Cloud web UI, or using the API or CLI. You can also specify a
deployment strategy in the [stack file](../apps/stack-yaml-reference.md) used to
define a [service stack](../apps/stacks.md).

For all methods, the default deployment strategy is "Emptiest node".

### Emptiest node (default)

This is the default strategy, and is commonly used to balance the total load of
all services across all nodes.

A service configured to deploy using the `EMPTIEST_NODE` strategy deploys its
containers to the nodes that match its [deploy tags](../apps/deploy-tags.md)
with the **fewest total containers** at the time of each container's deployment,
regardless of the service.

### High availability

This setting is typically used to increase the service availability.

A service using the `HIGH_AVAILABILITY` strategy deploys its containers to the
node that matches its deploy tags with the **fewest containers of that service**
at the time of each container's deployment. This means that the containers are
spread across all nodes that match the deploy tags for the service.

### Every node

A service using the `EVERY_NODE` strategy deploys one container **on each node** that matches its deploy tags.

When a service uses the `EVERY_NODE` strategy:

* A new container is deployed to every new node that matches the service's deploy tags.
* The service cannot be manually scaled.
* If the service uses volumes, each container on each node has a different volume.
* If an `EVERY_NODE` "client" service is linked to a "server" service that is also using the `EVERY_NODE` strategy, containers are linked one-to-one on each node. The "client" services are *not* automatically linked to "server" services on other nodes.

> **Note**: Because of how links are configured when using the **every node**
> strategy, you cannot currently switch from **every node** to **high
> availability** or **emptiest node** and vice versa.
