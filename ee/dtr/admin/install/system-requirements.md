---
title: Docker Trusted Registry system requirements
description: Learn about the system requirements for installing Docker Trusted Registry.
keywords: DTR, architecture, requirements
---

>{% include enterprise_label_shortform.md %}

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

When image scanning feature is used, we recommend that you have at least 32 GB of RAM. As developers and teams push images into DTR, the repository grows over time so you should inspect RAM, CPU, and disk usage on DTR nodes and increase resources when resource saturation is observed on regular basis.

## Ports used

When installing DTR on a node, make sure the following ports are open on that
node:

| Direction | Port    | Purpose                               |
|:---------:|:--------|:--------------------------------------|
|    in     | 80/tcp  | Web app and API client access to DTR. |
|    in     | 443/tcp | Web app and API client access to DTR. |

These ports are configurable when installing DTR.

## UCP Configuration

When installing or backing up DTR on a UCP cluster, Administrators need to be able to deploy
containers on "UCP manager nodes or nodes running DTR". This setting can be
adjusted in the [UCP Settings
menu](/ee/ucp/admin/configure/restrict-services-to-worker-nodes/).

The DTR installation or backup will fail with the following error message if
Administrators are unable to deploy on "UCP manager nodes or nodes running
DTR".

```
Error response from daemon: {"message":"could not find any nodes on which the container could be created"}
```

## Compatibility and maintenance lifecycle

Docker Enterprise Edition is a software subscription that includes three products:

* Docker Enterprise Engine
* Docker Trusted Registry
* Docker Universal Control Plane

[Learn more about the maintenance lifecycle for these products](https://success.docker.com/article/Compatibility_Matrix).

## Where to go next

- [DTR architecture](../../architecture.md)
- [Install DTR](index.md)
