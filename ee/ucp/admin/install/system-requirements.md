---
title: UCP System requirements
description: Learn about the system requirements for installing Docker Universal Control Plane.
keywords: UCP, architecture, requirements, Docker EE
ui_tabs:
- version: ucp-3.0
  orhigher: false
- version: ucp-2.2
  orlower: true
next_steps:
- path: plan-installation/
  title: Plan your installation
- path: ../../ucp-architecture/
  title: UCP architecture
redirect_from:
- /enterprise/admin/install/system-requirements/
---

{% if include.version=="ucp-3.0" %}

Docker Universal Control Plane can be installed on-premises or on the cloud.
Before installing, be sure your infrastructure has these requirements.

## Hardware and software requirements

You can install UCP on-premises or on a cloud provider. Common requirements:

* [Docker EE Engine](/engine/installation/index.md) version 17.06 or higher
* Linux kernel version 3.10 or higher
* A static IP address

### Minimum requirements

* 8GB of RAM for manager nodes or nodes running DTR
* 4GB of RAM for worker nodes
* 3GB of free disk space

### Recommended production requirements

 * 16GB of RAM for manager nodes or nodes running DTR
 * 4 vCPUs for manager nodes or nodes running DTR
 * 25-100GB of free disk space
 
Note that Windows container images are typically larger than Linux ones and for
this reason, you should consider provisioning more local storage for Windows
nodes and for DTR setups that will store Windows container images.

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

|       Hosts       |          Port           |          Scope          |                                      Purpose                                      |
| :---------------- | :---------------------- | :---------------------- | :-------------------------------------------------------------------------------- |
| managers, workers | TCP 179                 | Internal                | Port for BGP peers, used for kubernetes networking                                |
| managers          | TCP 443  (configurable) | External, Internal      | Port for the UCP web UI and API                                                   |
| managers          | TCP 2376 (configurable) | Internal                | Port for the Docker Swarm manager. Used for backwards compatibility               |
| managers          | TCP 2377 (configurable) | Internal,               | Port for control communication between swarm nodes                                |
| managers, workers | UDP 4789                | Internal,               | Port for overlay networking                                                       |
| managers          | TCP 6443 (configurable) | External, Internal      | Port for Kubernetes API server                                                    |
| managers, workers | TCP 6444                | Self                    | Port for Kubernetes API reverse proxy                                             |
| managers, workers | TCP, UDP 7946           | Internal                | Port for gossip-based clustering                                                  |
| managers, workers | TCP 10250               | Internal                | Port for Kubelet                                                                  |
| managers, workers | TCP 12376               | Internal                | Port for a TLS authentication proxy that provides access to the Docker Engine     |
| managers, workers | TCP 12378               | Self                    | Port for Etcd reverse proxy                                                       |
| managers          | TCP 12379               | Internal                | Port for Etcd Control API															|
| managers          | TCP 12380               | Internal                | Port for Etcd Peer API                                                            |
| managers          | TCP 12381               | Internal                | Port for the UCP cluster certificate authority                                    |
| managers          | TCP 12382               | Internal                | Port for the UCP client certificate authority                                     |
| managers          | TCP 12383               | Internal                | Port for the authentication storage backend                                       |
| managers          | TCP 12384               | Internal                | Port for the authentication storage backend for replication across managers       |
| managers          | TCP 12385               | Internal                | Port for the authentication service API                                           |
| managers          | TCP 12386               | Internal                | Port for the authentication worker                                                |
| managers          | TCP 12387               | Internal                | Port for the metrics service                                                      |

For overlay networks with encryption to work, you need to ensure that
IP protocol 50 (ESP) traffic is allowed.

Also, make sure the networks you're using allow the UCP components enough time
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

- Docker EE Engine 17.06 or higher
- DTR 2.5 or higher

{% elsif include.version=="ucp-2.2" %}

[UCP System requirements](/datacenter/ucp/2.2/guides/admin/install/system-requirements.md).

{% endif %}


