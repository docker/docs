---
title: Interlock architecture
description: Learn about Layer 7 routing, an application routing and load balancing system
  for Docker Swarm.
keywords: ucp, interlock, load balancing
ui_tabs:
- version: ucp-3.0
  orhigher: false
---

{% if include.version=="ucp-3.0" %}

The layer 7 routing solution for swarm workloads is known as Interlock, and has
three components:

* **Interlock-proxy**: This is a proxy/load-balancing service that handles the
requests from the outside world. By default this service is a containerized
NGINX deployment.
* **Interlock-extension**: This service monitors changes in your services and
generates the configuration used by the proxy service.
* **Interlock**: This is the central piece of the layer 7 routing solution.
It uses the Docker API to monitor events, and manages the extension and
proxy services.

This is what the default configuration looks like, once you enable layer 7
routing in UCP:

![](../images/interlock-architecture-1.svg)

An Interlock service starts running on a manager node, an Interlock-extension
service starts running on a worker node, and two replicas of the
Interlock-proxy service run on worker nodes.

## Routing lifecycle

By default layer 7 routing is disabled, so an administrator first needs to
enable this service from the UCP web UI.

Once that happens, UCP creates the `ucp-interlock` overlay network. Then the
`ucp-interlock` service starts and attaches to the Docker socket and the overlay
network that was created. This allows the Interlock service to use the
Docker API. That's also the reason why this service needs to run on a manger
node.

The `ucp-interlock` service then starts the `ucp-interlock-extension` service
and attaches it to the `ucp-interlock` network. This allows both services
to communicate.

The `ucp-interlock-extension` then generates a configuration to be used by
the proxy service. By default the proxy service is NGINX, so this service
generates a standard NGING configuration.

Finally, the `ucp-interlock` service takes this configuration and uses it to
start the `ucp-interlock-proxy` service.

At this point everything is ready for you to start using this service in your
applications.

You deploy your service and apply labels to it describing how the proxy
service should route traffic to that service. Once this happens, the
`ucp-interlock-extension` service generates a new configuration based on those
labels and forwards it to the `ucp-interlock` service, which in turn uses this
to redeploy the `ucp-interlock-proxy` with the new settings.

This all happens in milliseconds and with rolling updates, so that service
is never disrupted for incoming traffic.

{% endif %}
