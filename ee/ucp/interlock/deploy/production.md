---
title: Configure layer 7 routing for production
description: Learn how to configure the layer 7 routing solution for a production
  environment.
keywords: routing, proxy
---

The layer 7 solution that ships out of the box with UCP is highly available
and fault tolerant. It is also designed to work independently of how many
nodes you're managing with UCP.

![production deployment](../../images/interlock-deploy-production-1.svg)

For a production-grade deployment, you should tune the default deployment to
have two nodes dedicated for running the two replicas of the
`ucp-interlock-proxy` service. This ensures:

* The proxy services have dedicated resources to handle user requests. You
can configure these nodes with higher performance network interfaces.
* No application traffic can be routed to a manager node. This makes your
deployment secure.
* The proxy service is running on two nodes. If one node fails, layer 7 routing
continues working.

To achieve this you need to:

1. Enable layer 7 routing. [Learn how](index.md).
2. Pick two nodes that are going to be dedicated to run the proxy service.
3. Apply labels to those nodes, so that you can constrain the proxy service to
only run on nodes with those labels.
4. Update the `ucp-interlock` service to deploy proxies using that constraint.
5. Configure your load balancer to route traffic to the dedicated nodes only.

## Apply labels to nodes

In this example, we chose node-5 and node-6 to be dedicated just for running
the proxy service. To apply labels to those nodes run:

```bash
docker node update --label-add nodetype=loadbalancer <node>
```

To make sure the label was successfully applied, run:

{% raw %}
```bash
docker node inspect --format '{{ index .Spec.Labels "nodetype" }}' <node>
```
{% endraw %}

The command should print "loadbalancer".

## Configure the ucp-interlock service

Now that your nodes are labelled, you need to update the `ucp-interlock-proxy`
service configuration to deploy the proxy service with the correct constraints.

Add a constraint to the `ucp-interlock-proxy` service to update the running service:

```bash
docker service update \
    --constraint-add node.labels.nodetype==loadbalancer \
    ucp-interlock-proxy
```

Then add the constraint to the `ProxyConstraints` array in the `interlock-proxy` service
configuration so it takes effect if Interlock is restored from backup:

```toml
[Extensions]
  [Extensions.default]
    ProxyConstraints = ["node.labels.com.docker.ucp.orchestrator.swarm==true", "node.platform.os==linux", "node.labels.nodetype==loadbalancer"]
```

[Learn how to configure ucp-interlock](configure.md).

Once reconfigured you can check if the proxy service is running on the dedicated nodes:

```bash
docker service ps ucp-interlock-proxy
```

## Configure your load balancer

Once the proxy service is running on dedicated nodes, configure your upstream
load balancer with the domain names or IP addresses of those nodes.

This makes sure all traffic is directed to these nodes.

