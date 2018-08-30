---
description: Manage network in Docker Cloud
keywords: nodes, hosts, infrastructure, Cloud
title: Docker Cloud infrastructure overview (Standard Mode)
---

Docker Cloud uses an agent and system containers to deploy and manage nodes (hosts) on your behalf. All nodes accessible to your account are connected by an overlay or mesh network, regardless of host or cloud service provider.

## Deploy nodes from Docker Cloud

When you use Docker Cloud to deploy nodes on a hosted provider, the service stores your cloud provider credentials and then deploys nodes for you using the services' API to perform actions on your behalf.

## Bring your own host

If you are using [Bring Your Own Host](byoh.md), Docker Cloud provides a script that:

- installs the Docker Cloud Agent on the host
- downloads and installs the latest Docker CS Engine version and the AUFS storage driver
- sets up TLS certificates and the Docker security configuration
- registers the host with Docker Cloud under your user account

Once this connection is established, the Docker Cloud Agent manages the node and performs updates when the user requests them, and can also create and maintain a reverse tunnel to Docker Cloud if firewall restrictions prevent a direct connection.

## Internal networking

Docker Cloud communicates with the Docker daemon running in the node using the following IPs, on port **2375/tcp**.

- 52.204.126.235/32
- 52.6.30.174/32
- 52.205.192.142/32
- 52.205.2.114/32

If the port is not accessible, Docker Cloud creates a secure reverse tunnel from the nodes to Docker Cloud.

When you add a node on Docker Cloud, the node joins the Weave private overlay network for containers in other nodes by connecting on ports **6783/tcp** and **6783/udp**. (You should make sure these ports are open.)

## Node management

Nodes managed by Docker Cloud are connected to any other nodes owned by the user or organization, regardless of the host or service provider.

Docker Cloud uses system containers to do the following:

- Set up a secure overlay network between all nodes using Weave
- Create a stream of Docker events from nodes to Docker Cloud
- Synchronize node clocks
- Rotate container logs when they exceed 10 MB
- Remove `Terminated` images (images not used by a container for 30 minutes)

  > **Note**: If this is not sufficient for your needs, you can add a logging container to your services.

## Internal overlay network

Docker Cloud creates a per-user overlay network which connects all containers across all of the user's hosts. This network connects all of your containers on the `10.7.0.0/16` subnet, and gives every container a local IP. This IP persists on each container even if the container is redeployed and ends up on a different host. Every container can reach any other container on any port within the subnet.

## External access

The easiest way to access nodes is to ensure that your public ssh key is available to them. You can quickly copy your public key to all of the nodes in your Docker Cloud account by running the **authorizedkeys** container. See [SSHing into a node](ssh-into-a-node.md) for more information.

## What's in this section?
The pages in this section explain how to link Docker Cloud to your infrastructure providers or your own hosts, and how to manage your nodes from within Docker Cloud.

* [SSH into a Docker Cloud-managed node](ssh-into-a-node.md)
* Read more about [Deployment strategies](deployment-strategies.md)
* Learn how to [Upgrade Docker Engine on a node](docker-upgrade.md)
* [Use the Docker Cloud Agent to Bring your Own Host](byoh.md)
* [Link to Amazon Web Services hosts](link-aws.md)
    * [Using Docker Cloud on AWS FAQ](cloud-on-aws-faq.md)
* [Link to DigitalOcean hosts](link-do.md)
* [Link to Microsoft Azure hosts](link-azure.md)
* [Link to Packet hosts](link-packet.md)
    * [Using Docker Cloud and Packet FAQ](cloud-on-packet.net-faq.md)
* [Link to SoftLayer hosts](link-softlayer.md)
