+++
title = "Set up high availability"
description = "Docker Universal Control Plane"
[menu.main]
parent="mn_ucp"
+++

# Set up high availability

The UCP high availability (HA) feature allows you to replicate your UCP
controller within your cluster. When adding nodes to your cluster, you decide
which nodes you to use as *replicas* and which nodes are simply additional
Engines for extra capacity.

A replica is node in your cluster that can act as an additional UCP controller.
Should the primary controller fail, a replica can take over the controller role
for the cluster.

This document summarizes UCP's high availability feature and the concepts that
support it. It also explains general guidelines for deploying a highly available
UCP in production.

## Concepts and terminology

* **Primary Controller** This is the first host you run the bootstrapper `install` against.  It runs the following containers/services

    * **ucp-kv** This etcd container runs the replicated KV store
    * **ucp-swarm-manger** This Swarm Manager uses the replicated KV store for leader election and cluster membership tracking
    * **ucp-controller** This container runs the UCP server, using the replicated KV store for configuration state
    * **ucp-swarm-join** Runs the swarm join command to periodically publish this nodes existence to the KV store.  If the node goes down, this publishing stops, and the registration times out, and the node is automatically dropped from the cluster
    * **ucp-proxy** Runs a local TLS proxy for the docker socket to enable secure access of the local docker daemon
    * **ucp-swarm-ca[-proxy]** These **unreplicated** containers run the Swarm CA used for admin certificate bundles, and adding new nodes
    * **ucp-ca[-proxy]** These **unreplicated** containers run the (optional) UCP CA used for signing user bundles.
* **Replica Node**  This is a node you `join` to the primary using the `--replica` flag and it contributes to the availability of the cluster
    * **ucp-kv** This etcd container runs the replicated KV store
    * **ucp-swarm-manger** This Swarm Manager uses the replicated KV store for leader election and cluster membership tracking
    * **ucp-controller** This container runs the UCP server, using the replicated KV store for configuration state
    * **ucp-swarm-join** Runs the swarm join command to periodically publish this nodes existence to the KV store.  If the node goes down, this publishing stops, and the registration times out, and the node is automatically dropped from the cluster
    * **ucp-proxy** Runs a local TLS proxy for the docker socket to enable secure access of the local docker daemon
* **Non-Replica Node**  These nodes provide additional capacity, but do not enhance the availability of the UCP/Swarm infrastructure
    * **ucp-swarm-join** Runs the swarm join command to periodically publish this nodes existence to the KV store.  If the node goes down, this publishing stops, and the registration times out, and the node is automatically dropped from the cluster
    * **ucp-proxy** Runs a local TLS proxy for the docker socket to enable secure access of the local docker daemon

## Sizing your deployment

If you are planning an HA deployment, you should have a minimum of 3 controllers
configured, a primary and two replicas. Never run a cluster with only the
primary controller and a single replica.  This results in an HA
configuration of "2-nodes" where quorum is also "2-nodes" (to prevent
split-brain.)

If either the primary or single replica were to fail, the cluster is unusable until they are repaired. In fact, you actually have a higher failure
probability than if you just ran a non-HA setup with no replica.  

## Load balancing UCP cluster-store

At present, UCP does not include a load-balancer.  You may configure one your own. If you do, you can load balance between the primary and replica nodes on port `443` for web access to the system via a single IP/hostnamed.  

If an external load balancer is not used, system administrators should note the IP/hostname of the primary and all controller replicas. In this way, an administrator can access them when needed.

* Backups:
    * Users should always back up their volumes (see the other guides for a complete list of named volumes)
* The CAs (swarm and UCP) are not currently replicated.
    * Swarm CA:
        * Used for admin cert bundle generation
        * Used for adding hosts to the cluster
        * During an outage, no new admin cert bundles can be downloaded, but existing ones will still work.
        * During an outage, no new nodes can be added to the cluster, but existing nodes will continue to operate
    * UCP CA:
        * Used for user bundle generation
        * Used to sign certs for new replica nodes
        * During an outage, no new user cert bundles can be downloaded, but existing ones will still work
        * During an outage, no new replica nodes can be joined to the cluster
