---
title: Interlock architecture
description: Learn more about the architecture of the layer 7 routing solution
  for Docker swarm services.
keywords: routing, proxy
---

The layer 7 routing solution for swarm workloads is known as Interlock, and has
three components:

* **Interlock-proxy**: This is a proxy/load-balancing service that handles the
requests from the outside world. By default this service is a containerized
NGINX deployment.
* **Interlock-extension**: This is a helper service that generates the
configuration used by the proxy service.
* **Interlock**: This is the central piece of the layer 7 routing solution.
It uses the Docker API to monitor events, and manages the extension and
proxy services.

This is what the default configuration looks like, once you enable layer 7
routing in UCP:

![](../images/interlock-architecture-1.svg)

An Interlock service starts running on a manager node, an Interlock-extension
service starts running on a worker node, and two replicas of the
Interlock-proxy service run on worker nodes.

If you don't have any worker nodes in your cluster, then all Interlock
components run on manager nodes.

## Deployment lifecycle

By default layer 7 routing is disabled, so an administrator first needs to
enable this service from the UCP web UI.

Once that happens:

1. UCP creates the `ucp-interlock` overlay network.
2. UCP deploys the `ucp-interlock` service and attaches it both to the Docker
socket and the overlay network that was created. This allows the Interlock
service to use the Docker API. That's also the reason why this service needs to
run on a manger node.
3. The `ucp-interlock` service starts the `ucp-interlock-extension` service
and attaches it to the `ucp-interlock` network. This allows both services
to communicate.
4. The `ucp-interlock-extension` generates a configuration to be used by
the proxy service. By default the proxy service is NGINX, so this service
generates a standard NGINX configuration.
5. The `ucp-interlock` service takes the proxy configuration and uses it to
start the `ucp-interlock-proxy` service.

At this point everything is ready for you to start using the layer 7 routing
service with your swarm workloads.

## Routing lifecycle

Once the layer 7 routing service is enabled, you apply specific labels to
your swarm services. The labels define the hostnames that are routed to the
service, the ports used, and other routing configurations.

Once you deploy or update a swarm service with those labels:

1. The `ucp-interlock` service is monitoring the Docker API for events and
publishes the events to the `ucp-interlock-extension` service.
2. That service in turn generates a new configuration for the proxy service,
based on the labels you've added to your services.
3. The `ucp-interlock` service takes the new configuration and reconfigures the
`ucp-interlock-proxy` to start using it.

This all happens in milliseconds and with rolling updates. Even though
services are being reconfigured, users won't notice it.

