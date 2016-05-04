<!--[metadata]>
+++
title = "System requirements"
description = "Learn about the system requirements for installing Docker Universal Control Plane."
keywords = ["docker, ucp, architecture, requirements"]
[menu.main]
parent="mn_ucp_installation"
identifier="ucp_system_requirements"
weight=0
+++
<![end-metadata]-->

# UCP system requirements

Docker Universal Control Plane can be installed on-premises or on the cloud.
Before installing, be sure your infrastructure has these requirements.

## Hardware and software requirements

You can install UCP on-premises or on a cloud provider. To install UCP,
all nodes must have:

* 1.50 GB of RAM
* 3.00 GB of available disk space
* One of the supported operating systems installed:
  * RHEL 7.0, 7.1, or 7.2
  * Ubuntu 14.04 LTS
  * CentOS 7.1
  * SUSE Linux Enterprise 12
* Linux kernel version 3.10 or higher
* CS Docker Engine version 1.10 or higher

For highly-available installations, you also need a way to transfer files
between hosts.

## Ports used

When installing UCP on a host, make sure the following ports are open:

| Hosts              | Direction | Port                | Purpose                                                                    |
|:-------------------|:---------:|:--------------------|:---------------------------------------------------------------------------|
| controllers        |    in     | 443  (configurable) | Web app and CLI client access to UCP.                                      |
| controller         |    out    | 443                 | Send anonymous usage reports to Docker.                                    |
| controllers, nodes |    in     | 2375                | Heartbeat for nodes, to ensure they are running.                           |
| controllers        |    in     | 2376 (configurable) | Swarm manager accepts requests from UCP controller.                        |
| controllers, nodes |  in, out  | 4789                | Overlay networking.                                                        |
| controllers, nodes |  in, out  | 7946                | Overlay networking.                                                        |
| controllers, nodes |    in     | 12376               | Proxy for TLS, provides access to UCP, Swarm, and Engine.                  |
| controller         |    in     | 12379               | Internal node configuration, cluster configuration, and HA.                |
| controller         |    in     | 12380               | Internal node configuration, cluster configuration, and HA.                |
| controller         |    in     | 12381               | Proxy for TLS, provides access to UCP.                                     |
| controller         |    in     | 12382               | Manages TLS and requests from swarm manager.                               |
| controller         |    in     | 12383               | Used by the authentication storage backend.                                |
| controller         |    in     | 12384               | Used by authentication storage backend for replication across controllers. |
| controller         |    in     | 12385               | The port where the authentication API is exposed.                          |
| controller         |    in     | 12386               | Used by the authentication worker.                                         |

UCP collects anonymous usage metrics, to help us improve it. These metrics
are entirely anonymous, don’t identify your company, users, applications,
or any other sensitive information. You can disable this when installing
or on the UCP settings screen.

## Where to go next

* [UCP architecture](../architecture.md)
* [Plan a production installation](plan-production-install.md)
