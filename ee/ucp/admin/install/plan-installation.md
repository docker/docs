---
title: Plan a production UCP installation
description: Learn about the Docker Universal Control Plane architecture, and the requirements to install it on production.
keywords: UCP, install, Docker EE
---

Docker Universal Control Plane helps you manage your container cluster from a
centralized place. This article explains what you need to consider before
deploying Docker Universal Control Plane for production.

## System requirements

Before installing UCP you should make sure that all nodes (physical or virtual
machines) that you'll manage with UCP:

* [Comply with the system requirements](system-requirements.md), and
* Are running the same version of Docker Engine.

## Hostname strategy

Docker UCP requires Docker Enterprise Edition. Before installing Docker EE on
your cluster nodes, you should plan for a common hostname strategy.

Decide if you want to use short hostnames, like `engine01`, or Fully Qualified
Domain Names (FQDN), like `node01.company.example.com`. Whichever you choose,
confirm your naming strategy is consistent across the cluster, because
Docker Engine and UCP use hostnames.

For example, if your cluster has three hosts, you can name them:

```none
node1.company.example.com
node2.company.example.com
node3.company.example.com
```

## Static IP addresses

Docker UCP requires each node on the cluster to have a static IP address.
Before installing UCP, ensure your network and nodes are configured to support
this.

## Avoid IP range conflicts

The `service-cluster-ip-range` Kubernetes API Server flag is currently set to `10.96.0.0/16` and cannot be changed.

Swarm uses a default address pool of `10.0.0.0/16` for its overlay networks. If this conflicts with your current network implementation, please use a custom IP address pool. To specify a custom IP address pool, use the `--default-address-pool` command line option during [Swarm initialization](../../../../engine/swarm/swarm-mode.md). 

> **Note**: Currently, the UCP installation process does not support this flag. To deploy with a custom IP pool, Swarm must first be installed using this flag and UCP must be installed on top of it.

Kubernetes uses a default cluster IP pool for pods that is `192.168.0.0/16`. If it conflicts with your current networks, please use a custom IP pool by specifying `--pod-cidr` during UCP installation.

## Avoid firewall conflicts

For SUSE Linux Enterprise Server 12 SP2 (SLES12), the `FW_LO_NOTRACK` flag is turned on by default in the openSUSE firewall. This speeds up packet processing on the loopback interface, and breaks certain firewall setups that need to redirect outgoing packets via custom rules on the local machine.

To turn off the FW_LO_NOTRACK option, edit the `/etc/sysconfig/SuSEfirewall2` file and set `FW_LO_NOTRACK="no"`. Save the file and restart the firewall or reboot.

For For SUSE Linux Enterprise Server 12 SP3, the default value for `FW_LO_NOTRACK` was changed to `no`.

## Time synchronization

In distributed systems like Docker UCP, time synchronization is critical
to ensure proper operation. As a best practice to ensure consistency between
the engines in a UCP cluster, all engines should regularly synchronize time
with a Network Time Protocol (NTP) server. If a server's clock is skewed,
unexpected behavior may cause poor performance or even failures.

## Load balancing strategy

Docker UCP doesn't include a load balancer. You can configure your own
load balancer to balance user requests across all manager nodes.

If you plan to use a load balancer, you need to decide whether you'll
add the nodes to the load balancer using their IP addresses or their FQDNs.
Whichever you choose, be consistent across nodes. When this is decided,
take note of all IPs or FQDNs before starting the installation.

[Learn how to set up your load balancer](../configure/join-nodes/use-a-load-balancer.md).

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
one or more `--san` flags in the [install command](/reference/ucp/3.0/cli/install.md)
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

- [System requirements](system-requirements.md)
- [Install UCP](index.md)

