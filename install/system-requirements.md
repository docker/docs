<!--[metadata]>
+++
title = "System requirements"
description = "Learn about the system requirements for installing Docker Trusted Registry."
keywords = ["docker, DTR, architecture, requirements"]
[menu.main]
parent="workw_dtr_install"
identifier="dtr_system_requirements"
weight=0
+++
<![end-metadata]-->


# DTR system requirements

Docker Trusted Registry can be installed on-premises or on the cloud.
Before installing, be sure your infrastructure has these requirements.

## Software requirements

To install DTR on a node, that node node must be part of a Docker Universal
Control Plane 1.1 cluster.

<!-- TODO check these
* x.x GB of RAM
* x.x GB of available disk space
-->


## Ports used

When installing DTR on a node, make sure the following ports are open on that
node:

| Direction | Port | Purpose                                                                          |
|:---------:|:-----|:---------------------------------------------------------------------------------|
|    in     | 80   | Web app and API client access to DTR.                                            |
|    in     | 443  | Web app and API client access to DTR.                                            |
|    out    | 443  | Check if new versions are available, and send anonymous usage reports to Docker. |

DTR collects anonymous usage metrics, to help us improve it. These metrics
are entirely anonymous, don’t identify your company, users, applications,
or any other sensitive information. You can disable this on the DTR settings
page.

## Where to go next

* [DTR architecture](../architecture.md)
* [Install DTR](install-dtr.md)
