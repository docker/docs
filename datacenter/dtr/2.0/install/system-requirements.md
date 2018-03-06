---
description: Learn about the system requirements for installing Docker Trusted Registry.
keywords: docker, DTR, architecture, requirements
redirect_from:
- /docker-trusted-registry/install/system-requirements/
title: Docker Trusted Registry system requirements
---

Docker Trusted Registry can be installed on-premises or on the cloud.
Before installing, be sure your infrastructure has these requirements.

## Software requirements

To install DTR on a node, that node must be part of a Docker Universal
Control Plane 1.1 cluster.

## Ports used

When installing DTR on a node, make sure the following ports are open on that
node:

| Direction | Port | Purpose                                                                          |
|:---------:|:-----|:---------------------------------------------------------------------------------|
|    in     | 80   | Web app and API client access to DTR.                                            |
|    in     | 443  | Web app and API client access to DTR.                                            |
|    out    | 443  | Check if new versions are available, and send anonymous usage reports to Docker. |

The inbound ports are configurable.

DTR collects anonymous usage metrics, to help us improve it. These metrics
are entirely anonymous, don’t identify your company, users, applications,
or any other sensitive information. You can disable this on the DTR settings
page.

## Compatibility and maintenance lifecycle

Docker Datacenter is a software subscription that includes 3 products:

* CS Docker Engine,
* Docker Trusted Registry,
* Docker Universal Control Plane.

[Learn more about the maintenance lifecycle for these products](https://success.docker.com/article/Compatibility_Matrix).

## Where to go next

* [DTR architecture](../architecture.md)
* [Install DTR](index.md)
