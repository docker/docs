---
title: Docker Trusted Registry system requirements
description: Learn about the system requirements for installing Docker Trusted Registry.
keywords: DTR, architecture, requirements
---

Docker Trusted Registry can be installed on-premises or on the cloud.
Before installing, be sure your infrastructure has these requirements.

## Hardware and Software requirements

You can install DTR on-premises or on a cloud provider. To install DTR,
all nodes must:
* Be a worker node managed by UCP (Universal Control Plane). See [Compatibility Matrix](https://success.docker.com/article/compatibility-matrix) for version compatibility.
* Have a fixed hostname.

### Minimum requirements

* 16GB of RAM for nodes running DTR
* 2 vCPUs for nodes running DTR
* 10GB of free disk space

### Recommended production requirements

 * 16GB of RAM for nodes running DTR
 * 4 vCPUs for nodes running DTR
 * 25-100GB of free disk space
 
Note that Windows container images are typically larger than Linux ones and for
this reason, you should consider provisioning more local storage for Windows
nodes and for DTR setups that will store Windows container images.

## Ports used

When installing DTR on a node, make sure the following ports are open on that
node:

| Direction | Port    | Purpose                               |
|:---------:|:--------|:--------------------------------------|
|    in     | 80/tcp  | Web app and API client access to DTR. |
|    in     | 443/tcp | Web app and API client access to DTR. |

These ports are configurable when installing DTR.

## Compatibility and maintenance lifecycle

Docker Enterprise Edition is a software subscription that includes three products:

* Docker Enterprise Engine
* Docker Trusted Registry
* Docker Universal Control Plane

[Learn more about the maintenance lifecycle for these products](https://success.docker.com/article/Compatibility_Matrix).

## Where to go next

- [DTR architecture](../../architecture.md)
- [Install DTR](index.md)
