---
title: Configure layer 7 routing for production
description: Learn about Interlock, an application routing and load balancing system
  for Docker Swarm.
keywords: ucp, interlock, load balancing
ui_tabs:
- version: ucp-3.0
  orhigher: false
---

{% if include.version=="ucp-3.0" %}

The layer 7 solution that ships out of the box with UCP is highly available
and fault tolerant. It is also designed to work independently of how many
nodes you're managing with UCP.

![production deployment](../../images/interlock-deploy-production-1.svg)

For a production-grade deployment, you should tune the default deployment to
have two nodes dedicated for running the two replicas of the
`ucp-interlock-proxy` service. This makes sure:

* The proxy services have dedicated resources to handle user requests. You
can configure these nodes with higher performance network interfaces.
* No application traffic can be routed to a manager node. This makes the
deployment secure.
* The proxy service is running on two nodes. If one node fails layer 7 routing
still works.

To achieve this you need to:

1. Enable layer 7 routing. [Learn how](index.md).
2. Pick two nodes that are going to be dedicated to run the proxy service.
3. Apply labels to those nodes, so that you can constrain the proxy service to
only run on nodes with those labels.
4. Update the proxy service with the constraint.
5. Configure your load balancer to route traffic to the dedicated nodes only.

## Apply labels to nodes

In this example, we've chose node-5 and node-6 to be dedicated just for running
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

## Configure the proxy service

Now that your nodes are labelled, you can add a constraint to the
`ucp-interlock-proxy`service, so make sure it only gets scheduled on nodes
with the right label:

```bash
docker service update \
  --detach \
  --constraint-add node.labels.nodetype==loadbalancer \
  --stop-signal SIGTERM \
  --stop-grace-period 5s \
  $(docker service ls -f 'label=type=com.docker.interlock.core.proxy' -q)
```

This updates the proxy service to only be scheduled on node with the the
"loadbalancer" label. It also stops the task with a `SIGTERM` signal and gives
them five seconds to terminate, which allows the proxy service to stop accepting
new requests and finished serving existing requests from users.

Now you can check if the proxy service is running on the dedicated nodes:

```
docker service ps ucp-interlock-proxy
```

## Configure your load balancer

Once the proxy service is running on a dedicated node, configure your upstream
load balancer with the domain names or IP addresses of the nodes running
the proxy service.

This makes sure all traffic is directed to these nodes.

{% endif %}
