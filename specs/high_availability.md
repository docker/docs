# UCP High Availability

This document outlines how UCP high availability works, and general
guidelines for deploying a highly available UCP in production.
When adding nodes to your cluster, you decide which nodes you want to
be replicas, and which nodes are simply additional engines for extra
capacity.  If you are planning an HA deployment, you should have a
minimum of 3 nodes (primary + two replicas)

It is **highly** recommended that you deploy your initial 3 controller
nodes (primary + at least 2 replicas) **before** you start adding
non-replica nodes or start running workloads on your cluster.  When adding
the first replica, if an error occurrs, the cluster will be come unusable.

## Architecture

* **Primary Controller** This is the first node you run the `install` against.  It runs the following containers/services:
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

Notes:
* At present, UCP does not include a load-balancer.  Users may provide one exernally and load balance between the primary and replica nodes on port 443 for web access to the system via a single IP/hostname if desired.  If no external load balancer is used, admins should note the IP/hostname of the primary and all replicas so they can access them when needed.
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

**WARNING** You should never run a cluster with only the primary
controller and a single replica.  This will result in an HA configuration
of "2-nodes" where quorum is also "2-nodes" (to prevent split-brain.)
If either the primary or single replica were to fail, the cluster will be
unusable until they are repaired.  (So you actually have a higher failure
probability than if you just ran a non-HA setup with no replica.)  You
should have a minimum of 2 replicas (aka, "3-nodes") so that you can
tolerate at least a single failure.

**TODO** In the future this document should describe best practices for layout,
target number of nodes, etc.  For now, that's an exercise for the reader
based on etcd/raft documentation.
