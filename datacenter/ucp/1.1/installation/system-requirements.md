---
description: Learn about the system requirements for installing Docker Universal Control
  Plane.
keywords: docker, ucp, architecture, requirements
redirect_from:
- /ucp/installation/system-requirements/
title: UCP System requirements
---

Docker Universal Control Plane can be installed on-premises or on the cloud.
Before installing, be sure your infrastructure has these requirements.

## Hardware and software requirements

You can install UCP on-premises or on a cloud provider. To install UCP,
all nodes must have:

* Linux kernel version 3.10 or higher
* CS Docker Engine version 1.10 or higher. Learn about the
[operating systems supported by CS Docker Engine](/cs-engine/install/).
* 2.00 GB of RAM
* 3.00 GB of available disk space
* A static IP address

For highly-available installations, you also need a way to transfer files
between hosts.

## Ports used

When installing UCP on a host, make sure the following ports are open:

| Hosts              | Direction | Port                    | Purpose                                                                    |
|:-------------------|:---------:|:------------------------|:---------------------------------------------------------------------------|
| controllers, nodes |    in     | TCP 443  (configurable) | Web app and CLI client access to UCP.                                      |
| controllers, nodes |    in     | TCP 2375                | Heartbeat for nodes, to ensure they are running.                           |
| controllers        |    in     | TCP 2376 (configurable) | Swarm manager accepts requests from UCP controller.                        |
| controllers, nodes |  in, out  | UDP 4789                | Overlay networking.                                                        |
| controllers, nodes |  in, out  | TCP + UDP 7946          | Overlay networking.                                                        |
| controllers, nodes |    in     | TCP 12376               | Proxy for TLS, provides access to UCP, Swarm, and Engine.                  |
| controller         |    in     | TCP 12379               | Internal node configuration, cluster configuration, and HA.                |
| controller         |    in     | TCP 12380               | Internal node configuration, cluster configuration, and HA.                |
| controller         |    in     | TCP 12381               | Proxy for TLS, provides access to UCP.                                     |
| controller         |    in     | TCP 12382               | Manages TLS and requests from swarm manager.                               |
| controller         |    in     | TCP 12383               | Used by the authentication storage backend.                                |
| controller         |    in     | TCP 12384               | Used by authentication storage backend for replication across controllers. |
| controller         |    in     | TCP 12385               | The port where the authentication API is exposed.                          |
| controller         |    in     | TCP 12386               | Used by the authentication worker.                                         |

## Compatibility and maintenance lifecycle

Docker Datacenter is a software subscription that includes 3 products:

* CS Docker Engine,
* Docker Trusted Registry,
* Docker Universal Control Plane.

[Learn more about the maintenance lifecycle for these products](https://success.docker.com/article/Compatibility_Matrix).

## Where to go next

* [UCP architecture](../architecture.md)
* [Plan a production installation](plan-production-install.md)