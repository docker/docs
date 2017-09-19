---
description: Learn about the Docker Universal Control Plane architecture, and the requirements to install it on production.
keywords: docker, ucp, install, checklist
redirect_from:
- /ucp/plan-production-install/
- /ucp/installation/plan-production-install/
title: Plan a production installation
---

Docker Universal Control Plane can be installed on-premises, or
on a virtual private cloud. If you've never used Docker UCP before,
you should start by [installing it on a sandbox](../install-sandbox.md).

This article explains what you need to consider before deploying
Docker Universal Control Plane.

## System requirements

Before installing UCP, you should make sure all nodes of your cluster
comply with the [system requirements](system-requirements.md).


## Hostname strategy

Docker UCP requires the Docker CS Engine to run. Before installing Docker CS
Engine on the cluster nodes, you should plan for a common naming strategy.

Decide if you want to use short hostnames like `engine01` or Fully Qualified
Domain Names (FQDN) likes `engine01.docker.vm`. Independently of your choice,
ensure your naming strategy is consistent across the cluster, since UCP uses
the hostnames.

As an example, if your cluster has 4 hosts you can name them:

```bash
engine01.docker.vm
engine02.docker.vm
engine03.docker.vm
engine04.docker.vm
```

## Static IP addresses

Docker UCP requires each node on the cluster to have a static IP address.
Before installing UCP, ensure your network and nodes are configured to support
this.

## Time Synchronization

In distributed systems such as Docker UCP, time synchronization is critical
to ensure proper operation. As a best practice to ensure consistency between
then engines in a UCP cluster, all engines should regularly synchronize time
with a NTP server. If a server's clock is skewed, unexpected behavior may
cause poor performance or even failures.

## Load balancing strategy

UCP Docker UCP does not include a load-balancer. You can configure your own
load-balancer to balance user requests across all controller nodes.

If you plan on using a load balancer, you need to decide whether you are going
to add the nodes to the load balancer using their IP address, or their FQDN.
Independently of what you choose, it should be consistent across the  nodes.

After that, you should take note of all IPs or FQDNs before starting the
installation.

## Load balancing UCP and DTR

By default, both UCP and DTR use port 443. If you plan on deploying UCP and DTR,
your load balancer needs to distinguish traffic between the two by IP address
or port number.

* If you want to configure your load balancer to listen on port 443:
    * Use one load balancer for UCP, and another for DTR,
    * Use the same load balancer with multiple virtual IPs.
* Configure your load balancer to expose UCP or DTR on a port other than 443.

If you want to install UCP in a high-availability configuration that uses
a load balancer in front of your UCP controllers, include the appropriate IP
address and fully qualified domain name of the load balancer's VIP by using
one or more `--san` flags in the [install command](../reference/install.md)
or when you're asked for additional SANs in interactive mode.

## Using external CAs

You can customize UCP to use certificates signed by an external Certificate
Authority. If you decide to use your own CAs take in consideration that:

* During the installation you need to copy the ca.pem, cert.pem, and key.pem
files across all controller hosts,
* The ca.pem is the root CA public certificate
* The cert.pem is the server cert plus any intermediate CA public certificates,
* The cert.pem should have SANs for all addresses used to reach UCP,
* The key.pem is the server private key,

You can have a certificate for each controller, with a common SAN. As an
example, on a three node cluster you can have:

* engine01.docker.vm with SAN ducp.docker.vm
* engine02.docker.vm with SAN ducp.docker.vm
* engine03.docker.vm with SAN ducp.docker.vm

## File transfer across hosts

Make sure you can transfer file between the hosts on the cluster. You will
need to replicate CAs across controller nodes.

For this, you can tools like `scp` or `rsync`, or configure the hosts to use
a network file system.


## Where to go next

* [UCP system requirements](system-requirements.md)
* [Install UCP for production](install-production.md)