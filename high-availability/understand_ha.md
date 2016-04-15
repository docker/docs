<!--[metadata]>
+++
aliases = [ "/ucp/understand_ha/"]
title ="Set up high availability"
description="Docker Universal Control plane has support for high availability. Learn how to set up your installation to ensure it tolerates failures."
keywords= ["replica, controller, availability, high, ucp"]
[menu.main]
parent="mn_ucp_high_availability"
+++
<![end-metadata]-->


# Set up high availability

Docker Universal Control Plane is designed for high availability (HA).
When setting up a UCP cluster, you can add additional nodes to serve as
replicas of the controller.

Adding replica nodes to your cluster allows you to:

* Load-balance user requests across the controller and replica nodes,
* Maintain the cluster state in case of failure.

This page explains some of the components of UCP that provide support for
high availability. It also provides some guidelines on how to set up UCP to
ensure it can handle failures.

## Understand high availability terms and containers

A Docker UCP installation is made of several nodes:

* Controller node: the node that handles user requests,
* Replica nodes: replicas of the controller node, for high-availability,
* Nodes: the nodes that run your containers.

The **controller** is the first node added to cluster by running the
 `ucp install`. This node runs the following containers:

| Name                | Description                                                                                                                                                                                                                       |
|:--------------------|:----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| ucp-kv              | This container runs the key-value store used by UCP. Don't use this key-value store in your applications, since it's for internal use only.                                                                                       |
| ucp-swarm-manager   | This Swarm manager uses the replicated KV store for leader election and cluster membership tracking.                                                                                                                              |
| ucp-controller      | This container runs the UCP server, using the replicated KV store for configuration state.                                                                                                                                        |
| ucp-swarm-join      | Runs the `swarm join` command to periodically publish this nodes existence to the KV store. If the node goes down, this publishing stops, and the registration times out, and the node is automatically dropped from the cluster. |
| ucp-proxy           | Runs a local TLS proxy for the docker socket to enable secure access of the local docker daemon.                                                                                                                                  |
| ucp-cluster-root-ca | These containers run the Swarm CA used for admin certificate bundles, and adding new nodes.                                                                                                                                       |
| ucp-client-root-ca  | These containers run the (optional) UCP CA used for signing user bundles.                                                                                                                                                         |

A **replica node** is a node you add to the cluster by running the
`ucp join --replica` command. These nodes run the following containers:

| Name              | Description                                                                                                                                                                                                                       |
|:------------------|:----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| ucp-kv            | This etcd container runs the replicated KV store.                                                                                                                                                                                 |
| ucp-swarm-manager | This Swarm manager uses the replicated KV store for leader election and cluster membership tracking.                                                                                                                              |
| ucp-controller    | This container runs the UCP server, using the replicated KV store for configuration state.                                                                                                                                        |
| ucp-swarm-join    | Runs the `swarm join` command to periodically publish this nodes existence to the KV store. If the node goes down, this publishing stops, and the registration times out, and the node is automatically dropped from the cluster. |
| ucp-proxy         | Runs a local TLS proxy for the docker socket to enable secure access of the local docker daemon.                                                                                                                                  |

The remaining **non-replica nodes** provide additional capacity to the cluster,
to run your own containers and applications. They don't contribute to the
high-availability of the UCP cluster. These nodes run the following containers:

| Name           | Description                                                                                                                                                                                                                       |
|:---------------|:----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| ucp-swarm-join | Runs the `swarm join` command to periodically publish this nodes existence to the KV store. If the node goes down, this publishing stops, and the registration times out, and the node is automatically dropped from the cluster. |
| ucp-proxy      | Runs a local TLS proxy for the docker socket to enable secure access of the local docker daemon.                                                                                                                                  |

## Size your deployment

To make the cluster tolerant to more failures, add additional replica nodes to
your cluster. For an high-availability deployment of UCP, you should have at
least one controller and two replicas in your cluster.

| Controller and replicas | Failures tolerated |
|:-----------------------:|:------------------:|
|            1            |         0          |
|            3            |         1          |
|            5            |         2          |
|            7            |         3          |


When sizing your cluster, follow these rules of thumb:

* Don't create a cluster with just one controller and one replica. Your cluster
won't tolerate any failures, and it's possible that you experience performance
degradation.
* When a replica fails, the number of failures tolerated by your cluster
decreases. Don't leave that replica offline for long.
* Adding too many replicas to the cluster might also lead to performance
degradation, as each value stored in the key-value store needs to be
replicated across all replicas.


## Load-balancing on UCP

At present, UCP does not include a load-balancer. You may configure one your
own. If you do, you can load balance between the primary and replica nodes on
port `443` for web access to the system via a single IP/hostname.  

If an external load balancer is not used, system administrators should note the
IP/hostname of the primary and all controller replicas. In this way, an
administrator can access them when needed.



## Where to go next

* [UCP architecture](../architecture.md)
* [Install UCP for production](../installation/install-production.md)
