---
title: UCP System requirements
description: Learn about the system requirements for installing Docker Universal Control Plane.
keywords: UCP, architecture, requirements, Docker EE
---

Docker Universal Control Plane can be installed on-premises or on the cloud.
Before installing, be sure your infrastructure has these requirements.

## Hardware and software requirements

You can install UCP on-premises or on a cloud provider. To install UCP,
all nodes must have:

* [Docker Enterprise Edition](/engine/installation/index.md) version 17.06 or higher
* Linux kernel version 3.10 or higher
* 8.00 GB of RAM for manager nodes or nodes running DTR
* 4.00 GB of RAM for worker nodes
* 3.00 GB of available disk space
* A static IP address

Also, make sure the nodes are running one of these operating systems:

* A maintained version of CentOS 7. Archived versions aren't supported or tested.
* Red Hat Enterprise Linux 7.0, 7.1, 7.2, 7.3, or 7.4
* Ubuntu 14.04 LTS or 16.04 LTS
* SUSE Linux Enterprise 12
* Oracle Linux 7.3

For highly-available installations, you also need a way to transfer files
between hosts.

> Workloads on manager nodes
>
> These requirements assume that manager nodes won't run regular workloads.
> If you plan to run additional workloads on manager nodes, you may need to
> provision more powerful nodes. If manager nodes become overloaded, the
> swarm may experience issues.

## Ports used

When installing UCP on a host, make sure the following ports are open:

| Hosts             | Direction | Port                    | Purpose                                                                           |
|:------------------|:---------:|:------------------------|:----------------------------------------------------------------------------------|
| managers, workers |    in     | TCP 443  (configurable) | Port for the UCP web UI and API                                                   |
| managers          |    in     | TCP 2376 (configurable) | Port for the Docker Swarm manager. Used for backwards compatibility               |
| managers, workers |    in     | TCP 2377 (configurable) | Port for communication between swarm nodes                                        |
| managers, workers |  in, out  | UDP 4789                | Port for overlay networking                                                       |
| managers, workers |  in, out  | TCP, UDP 7946           | Port for gossip-based clustering                                                  |
| managers, workers |    in     | TCP 12376               | Port for a TLS proxy that provides access to UCP, Docker Engine, and Docker Swarm |
| managers          |    in     | TCP 12379               | Port for internal node configuration, cluster configuration, and HA               |
| managers          |    in     | TCP 12380               | Port for internal node configuration, cluster configuration, and HA               |
| managers          |    in     | TCP 12381               | Port for the certificate authority                                                |
| managers          |    in     | TCP 12382               | Port for the UCP certificate authority                                            |
| managers          |    in     | TCP 12383               | Port for the authentication storage backend                                       |
| managers          |    in     | TCP 12384               | Port for the authentication storage backend for replication across managers       |
| managers          |    in     | TCP 12385               | Port for the authentication service API                                           |
| managers          |    in     | TCP 12386               | Port for the authentication worker                                                |
| managers          |    in     | TCP 12387               | Port for the metrics service                                                      |

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
| Stand-alone swarm                      | 90000        | no           |

## Compatibility and maintenance lifecycle

Docker EE is a software subscription that includes 3 products:

* Docker Engine with enterprise-grade support,
* Docker Trusted Registry,
* Docker Universal Control Plane.

[Learn more about the maintenance lifecycle for these products](http://success.docker.com/Get_Help/Compatibility_Matrix_and_Maintenance_Lifecycle).

## Version compatibility

UCP 2.2 requires minimum versions of the following Docker components:

- Docker Engine 17.06 or higher
- DTR 2.3 or higher

<!--
- Docker Remote API 1.25
- Compose 1.9
-->

## Where to go next

* [UCP architecture](../../architecture.md)
* [Plan your installation](plan-installation.md)
