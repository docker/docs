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
* 4.00 GB of RAM
* 3.00 GB of available disk space
* A static IP address

Also, make sure the nodes are running one of these operating systems:

* CentOS 7.1 or 7.2
* Red Hat Enterprise Linux 7.0, 7.1, 7.2, or 7.3
* Ubuntu 14.04 LTS or 16.04 LTS
* SUSE Linux Enterprise 12
* Oracle Linux 7.3

For highly-available installations, you also need a way to transfer files
between hosts.

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
