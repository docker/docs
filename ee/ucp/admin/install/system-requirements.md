---
title: UCP System requirements
description: Learn about the system requirements for installing Docker Universal Control Plane.
keywords: UCP, architecture, requirements, Docker EE
redirect_from:
- /enterprise/admin/install/system-requirements/
---

Docker Universal Control Plane can be installed on-premises or on the cloud.
Before installing, be sure your infrastructure has these requirements.

## Hardware and software requirements

You can install UCP on-premises or on a cloud provider. Common requirements:

* [Docker EE Engine](/ee/supported-platforms.md) version {{ site.docker_ee_version }}
* Linux kernel version 3.10 or higher
* [A static IP address for each node in the cluster](/ee/ucp/admin/install/plan-installation/#static-ip-addresses)
 
### Minimum requirements

* 8GB of RAM for manager nodes
* 4GB of RAM for worker nodes
* 2 vCPUs for manager nodes
* 5GB of free disk space for the `/var` partition for manager nodes (A minimum of 6GB is recommended.)
* 500MB of free disk space for the `/var` partition for worker nodes

**Note**: Increased storage is required for Kubernetes manager nodes in UCP 3.1. If you are upgrading to UCP 3.1, refer to [Kubelet restarting after upgrade to Universal Control Plane 3.1](https://success.docker.com/article/kublet-restarting-after-upgrade-to-universal-control-plane-31) for information on how to increase the size of the `/var/lib/kubelet` filesystem.

### Recommended production requirements

 * 16GB of RAM for manager nodes
 * 4 vCPUs for manager nodes
 * 25-100GB of free disk space

Note that Windows container images are typically larger than Linux container images. For
this reason, you should provision more local storage for Windows
nodes and for any DTR setups that store Windows container images.

Also, make sure the nodes are running an [operating system support by Docker EE](https://success.docker.com/Policies/Compatibility_Matrix).

For highly-available installations, you also need a way to transfer files
between hosts.

> Workloads on manager nodes
>
> These requirements assume that manager nodes won't run regular workloads.
> If you plan to run additional workloads on manager nodes, you may need to
> provision more powerful nodes. If manager nodes become overloaded, the
> cluster may experience issues.

## Ports used

When installing UCP on a host, a series of ports need to be opened to incoming
traffic. Each of these ports will expect incoming traffic from a set of hosts,
indicated as the "Scope" of that port. The three scopes are:
- External: Traffic arrives from outside the cluster through end-user
  interaction.
- Internal: Traffic arrives from other hosts in the same cluster.
- Self: Traffic arrives to that port only from processes on the same host.

Make sure the following ports are open for incoming traffic on the respective
host types:

|       Hosts       |          Port           |       Scope        |                                    Purpose                                    |
| :---------------- | :---------------------- | :----------------- | :---------------------------------------------------------------------------- |
| managers, workers | TCP 179                 | Internal           | Port for BGP peers, used for kubernetes networking                            |
| managers          | TCP 443  (configurable) | External, Internal | Port for the UCP web UI and API                                               |
| managers          | TCP 2376 (configurable) | Internal           | Port for the Docker Swarm manager. Used for backwards compatibility           |
| managers          | TCP 2377 (configurable) | Internal,          | Port for control communication between swarm nodes                            |
| managers, workers | UDP 4789                | Internal,          | Port for overlay networking                                                   |
| managers          | TCP 6443 (configurable) | External, Internal | Port for Kubernetes API server endpoint                                       |
| managers, workers | TCP 6444                | Self               | Port for Kubernetes API reverse proxy                                         |
| managers, workers | TCP, UDP 7946           | Internal           | Port for gossip-based clustering                                              |
| managers, workers | TCP 9099                | Self               | Port for calico health check
| managers, workers | TCP 10250               | Internal           | Port for Kubelet                                                              |
| managers, workers | TCP 12376               | Internal           | Port for a TLS authentication proxy that provides access to the Docker Engine |
| managers, workers | TCP 12378               | Self               | Port for Etcd reverse proxy                                                   |
| managers          | TCP 12379               | Internal           | Port for Etcd Control API                                                     |
| managers          | TCP 12380               | Internal           | Port for Etcd Peer API                                                        |
| managers          | TCP 12381               | Internal           | Port for the UCP cluster certificate authority                                |
| managers          | TCP 12382               | Internal           | Port for the UCP client certificate authority                                 |
| managers          | TCP 12383               | Internal           | Port for the authentication storage backend                                   |
| managers          | TCP 12384               | Internal           | Port for the authentication storage backend for replication across managers   |
| managers          | TCP 12385               | Internal           | Port for the authentication service API                                       |
| managers          | TCP 12386               | Internal           | Port for the authentication worker                                            |
| managers          | TCP 12388               | Internal           | Internal Port for the Kubernetes API Server                                   |

## Avoid firewall conflicts

For SUSE Linux Enterprise Server 12 SP2 (SLES12), the `FW_LO_NOTRACK` flag is turned on by default in the openSUSE firewall. This speeds up packet processing on the loopback interface, and breaks certain firewall setups that need to redirect outgoing packets via custom rules on the local machine.

To turn off the FW_LO_NOTRACK option, edit the `/etc/sysconfig/SuSEfirewall2` file and set `FW_LO_NOTRACK="no"`. Save the file and restart the firewall or reboot.

For For SUSE Linux Enterprise Server 12 SP3, the default value for `FW_LO_NOTRACK` was changed to `no`.

## Enable ESP traffic

For overlay networks with encryption to work, you need to ensure that
IP protocol 50 (Encapsulating Security Payload) traffic is allowed.

## Enable IP-in-IP traffic

The default networking plugin for UCP is Calico, which uses IP Protocol
Number 4 for IP-in-IP encapsulation.

If you're deploying to AWS or another cloud provider, enable IP-in-IP
traffic for your cloud provider's security group.

## Timeout settings

Make sure the networks you're using allow the UCP components enough time
to communicate before they time out.

| Component                              | Timeout (ms) | Configurable |
|:---------------------------------------|:-------------|:-------------|
| Raft consensus between manager nodes   | 3000         | no           |
| Gossip protocol for overlay networking | 5000         | no           |
| etcd                                   | 500          | yes          |
| RethinkDB                              | 10000        | no           |
| Stand-alone cluster                    | 90000        | no           |

## Time Synchronization

In distributed systems like Docker UCP, time synchronization is critical
to ensure proper operation. As a best practice to ensure consistency between
the engines in a UCP cluster, all engines should regularly synchronize time
with a Network Time Protocol (NTP) server. If a server's clock is skewed,
unexpected behavior may cause poor performance or even failures.

## Compatibility and maintenance lifecycle

Docker EE is a software subscription that includes three products:

* Docker Engine with enterprise-grade support
* Docker Trusted Registry
* Docker Universal Control Plane

Learn more about compatibility and the maintenance lifecycle for these products:

- [Compatibility Matrix](https://success.docker.com/Policies/Compatibility_Matrix)
- [Maintenance Lifecycle](https://success.docker.com/Policies/Maintenance_Lifecycle)

## Version compatibility

UCP {{ page.ucp_version }} requires minimum versions of the following Docker components:

- Docker Enterprise Engine 18.09.0-ee-1 or higher
- DTR 2.6 or higher

## Where to go next

- [Plan your installation](plan-installation.md)
- [UCP architecture](../../ucp-architecture.md)
