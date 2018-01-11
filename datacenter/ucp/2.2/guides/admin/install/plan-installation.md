---
title: Plan a production UCP installation
description: Learn about the Docker Universal Control Plane architecture, and the requirements to install it on production.
keywords: UCP, install, Docker EE
---

Docker Universal Control Plane helps you manage your container swarm from a
centralized place. This article explains what you need to consider before
deploying Docker Universal Control Plane for production.

## System requirements

Before installing UCP, make sure that all nodes (physical or virtual
machines) that UCP manages:

* [Comply with the system requirements](system-requirements.md), and
* Are running the same version of Docker Engine.

## Hostname strategy

Docker UCP requires Docker Enterprise Edition. Before installing Docker EE on
your swarm nodes, you should plan for a common hostname strategy.

Decide if you want to use short hostnames, like `engine01`, or Fully Qualified
Domain Names (FQDN), like `engine01.docker.vm`. Whichever you choose,
ensure that your naming strategy is consistent across the cluster, because
Docker Engine and UCP use hostnames.

For example, if your swarm has three hosts, you can name them:

```none
node1.company.example.org
node2.company.example.org
node3.company.example.org
```

## Static IP addresses

Docker UCP requires each node on the cluster to have a static IP address.
Before installing UCP, ensure your network and nodes are configured to support
this.

## Time synchronization

In distributed systems like Docker UCP, time synchronization is critical
to ensure proper operation. As a best practice to ensure consistency between
the engines in a UCP swarm, all engines should regularly synchronize time
with a Network Time Protocol (NTP) server. If a server's clock is skewed,
unexpected behavior may cause poor performance or even failures.

## Load balancing strategy

Docker UCP doesn't include a load balancer. You can configure your own
load balancer to balance user requests across all manager nodes.

If you plan to use a load balancer, you need to decide whether to
add the nodes to the load balancer using their IP addresses or their FQDNs.
Whichever you choose, be consistent across nodes. When this is decided,
take note of all IPs or FQDNs before starting the installation.

## Load balancing UCP and DTR

By default, UCP and DTR both use port 443. If you plan on deploying UCP and
DTR, your load balancer needs to distinguish traffic between the two by IP
address or port number.

* If you want to configure your load balancer to listen on port 443:
    * Use one load balancer for UCP and another for DTR,
    * Use the same load balancer with multiple virtual IPs.
* Configure your load balancer to expose UCP or DTR on a port other than 443.

If you want to install UCP in a high-availability configuration that uses
a load balancer in front of your UCP controllers, include the appropriate IP
address and FQDN of the load balancer's VIP by using
one or more `--san` flags in the [install command](../../../reference/cli/install.md)
or when you're asked for additional SANs in interactive mode.
[Learn about high availability](../configure/set-up-high-availability.md).

## Use an external Certificate Authority

You can customize UCP to use certificates signed by an external Certificate
Authority. When using your own certificates, you need to have a certificate
bundle that has:

* A ca.pem file with the root CA public certificate,
* A cert.pem file with the server certificate and any intermediate CA public
certificates. This certificate should also have SANs for all addresses used to
reach the UCP manager,
* A key.pem file with server private key.

You can have a certificate for each manager, with a common SAN. For
example, on a three-node cluster, you can have:

* node1.company.example.org with SAN ucp.company.org
* node2.company.example.org with SAN ucp.company.org
* node3.company.example.org with SAN ucp.company.org

You can also install UCP with a single externally-signed certificate
for all managers, rather than one for each manager node. In this case,
the certificate files are copied automatically to any new
manager nodes joining the cluster or being promoted to a manager role.

## Where to go next

* [UCP system requirements](system-requirements.md)
* [Install UCP](index.md)
