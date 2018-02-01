---
description: Learn about the system requirements for installing Docker Trusted Registry.
keywords: docker, DTR, architecture, requirements
title: Docker Trusted Registry system requirements
---

Docker Trusted Registry can be installed on-premises or on the cloud.
Before installing, be sure your infrastructure meets these requirements.

## Software requirements

You can only install DTR on a node that is being managed by Docker Universal
Control Plane.

## Ports used

When installing DTR on a node, make sure the following ports are open on that
node:

| Direction | Port    | Purpose                               |
|:---------:|:--------|:--------------------------------------|
|    in     | 80/tcp  | Web app and API client access to DTR. |
|    in     | 443/tcp | Web app and API client access to DTR. |

These ports are configurable when installing DTR.

## Compatibility and maintenance lifecycle

Docker Datacenter is a software subscription that includes 3 products:

* CS Docker Engine,
* Docker Trusted Registry,
* Docker Universal Control Plane.

[Learn more about the maintenance lifecycle for these products](https://success.docker.com/article/Compatibility_Matrix).

## Where to go next

* [DTR architecture](../../architecture.md)
* [Install DTR](index.md)
