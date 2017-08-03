---
description: Learn about the system requirements for installing Docker Universal Control
  Plane.
keywords: docker, ucp, architecture, requirements
title: UCP System requirements
---

Docker Universal Control Plane can be installed on-premises or on the cloud.
Before installing, be sure your infrastructure has these requirements.

## Hardware and software requirements

You can install UCP on-premises or on a cloud provider. To install UCP,
all nodes must have:

* Linux kernel version 3.10 or higher
* CS Docker Engine version 1.13.0 or higher
* 4.00 GB of RAM
* 3.00 GB of available disk space
* A static IP address

For highly-available installations, you also need a way to transfer files
between hosts.

## Network requirements

When installing UCP on a host, make sure the following ports are open:

| Hosts             | Direction | Port                    | Purpose                                                                           |
|:------------------|:---------:|:------------------------|:----------------------------------------------------------------------------------|
| managers, workers |    in     | TCP 443  (configurable) | Port for the UCP web UI and API                                                   |
| managers          |    in     | TCP 2376 (configurable) | Port for the Docker Swarm manager. Used for backwards compatibility               |
| managers, workers |    in     | TCP 2377 (configurable) | Port for communication between swarm nodes                                        |
| managers, workers |  in, out  | TCP, UDP 4789           | Port for overlay networking                                                       |
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

Also, make sure the networks you're using allow the UCP components enough time to
communicate before they time out.

| Component                              | Timeout (ms) | Configurable |
|:---------------------------------------|:-------------|:-------------|
| Raft consensus between manager nodes   | 3000         | no           |
| Gossip protocol for overlay networking | 5000         | no           |
| etcd                                   | 500          | yes          |
| rethinkDB                              | 10000        | no           |
| Stand-alone swarm                      | 90000        | no           |


## Compatibility and maintenance lifecycle

Docker Datacenter is a software subscription that includes 3 products:

* CS Docker Engine,
* Docker Trusted Registry,
* Docker Universal Control Plane.

[Learn more about the maintenance lifecycle for these products](http://success.docker.com/Get_Help/Compatibility_Matrix_and_Maintenance_Lifecycle).

## Version compatibility

UCP 2.1 requires minimum versions of the following Docker components:

- Docker Engine 1.13.0
- Docker Remote API 1.25
- Compose 1.9

## Where to go next

* [UCP architecture](../../architecture.md)
* [Plan your installation](plan-installation.md)
