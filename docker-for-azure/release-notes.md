---
description: Release notes
keywords: azure, microsoft, iaas, tutorial, edge, stable
title: Docker for Azure Release Notes
---

{% include d4a_buttons.md %}

## Enterprise Edition
[Docker Enterprise Edition Lifecycle](https://success.docker.com/Policies/Maintenance_Lifecycle){: target="_blank"}<!--_-->

[Deploy Docker Enterprise Edition (EE) for AWS](https://store.docker.com/editions/enterprise/docker-ee-aws?tab=description){: target="_blank" class="button outline-btn"}

### 17.06 EE

- Docker engine 17.06 EE
- For Std/Adv external logging has been removed, as it is now handled by [UCP](https://docs.docker.com/datacenter/ucp/2.0/guides/configuration/configure-logs/){: target="_blank"}
- UCP 2.2.3
- DTR 2.3.3

### 17.03 EE

- Docker engine 17.03 EE
- UCP 2.1.5
- DTR 2.2.7

## Stable channel

### 18.03 CE

{{azure_blue_latest}}

Release date: 3/21/2018

- Docker Engine upgraded to [Docker 18.03.0 CE](https://github.com/docker/docker-ce/releases/tag/v18.03.0-ce){: target="_blank" class="_"}

### 17.12.1 CE

Release date: 3/1/2018

- Docker Engine upgraded to [Docker 17.12.1 CE](https://github.com/docker/docker-ce/releases/tag/v17.12.0-ce){: target="_blank" class="_"}

### 17.12 CE

Release date: 1/9/2018

- Docker Engine upgraded to [Docker 17.12.0 CE](https://github.com/docker/docker-ce/releases/tag/v17.12.0-ce){: target="_blank" class="_"}
- Kernel patch to mitigates Meltdown attacks ( CVE-2017-5754) and enable KPTI

> **Note** There was an issue in LinuxKit that prevented containers from [starting after a machine reboot](https://github.com/moby/moby/issues/36189){: target="_blank" class="_"}.

### 17.09 CE

Release date: 10/6/2017

- Docker Engine upgraded to [Docker 17.09.0 CE](https://github.com/docker/docker-ce/releases/tag/v17.09.0-ce){: target="_blank" class="_"}
- Moby mounts for early reboot support
- Docker binary bundled where needed to allow easier host interchange
- Azure VHD use full hard drive space

### 17.06.2 CE

Release date: 09/08/2017

**New**

- Docker Engine upgraded to [Docker 17.06.2 CE](https://github.com/docker/docker-ce/releases/tag/v17.06.2-ce){: target="_blank"}
- VMSS APIs updated to use 2017-03-30 to allow upgrades when customdata has changed

### 17.06.1 CE

Release date: 08/17/2017

**New**

- Docker Engine upgraded to [Docker 17.06.1 CE](https://github.com/docker/docker-ce/releases/tag/v17.06.1-ce){: target="_blank"}
- Improvements to CloudStor support
- Azure agent logs are limited to 50mb with proper logrotation

### 17.06.0 CE

Release date: 06/28/2017

**New**

- Docker Engine upgraded to [Docker 17.06.0 CE](https://github.com/docker/docker/blob/master/CHANGELOG.md)
- Introduced a new way to kick off upgrades through a container. The old upgrade.sh is no longer supported.
- Introduced new SMB mount related parameters for Cloudstor volumes [persistent storage volumes](persistent-data-volumes.md).

### 17.03.1 CE

Release date: 03/30/2017

**New**

- Docker Engine upgraded to [Docker 17.03.1 CE](https://github.com/docker/docker/blob/master/CHANGELOG.md)
- Fixed bugs in the way container logs are uploaded to File Storage in the storage account for logs

### 17.03.0 CE

Release date: 02/08/2017

**New**

- Docker Engine upgraded to [Docker 17.03.0 CE](https://github.com/docker/docker/blob/master/CHANGELOG.md)

### 1.13.1-2

Release date: 02/08/2017

**New**

- Docker Engine upgraded to [Docker 1.13.1](https://github.com/docker/docker/blob/master/CHANGELOG.md)

### 1.13.0-1

Release date: 01/18/2017

**New**

- Docker Engine upgraded to [Docker 1.13.0](https://github.com/docker/docker/blob/master/CHANGELOG.md)
- Writing to home directory no longer requires `sudo`
- Added support to perform fine grained monitoring of health status of swarm nodes, destroy unhealthy nodes and create replacement nodes
- Added support to scale the number of nodes in manager and worker vm scale sets through Azure UI/CLI for managing the number of nodes in a scale set
- Improved logging and remote diagnostics mechanisms for system containers

## Edge channel

### 18.01 CE

{{aws_blue_edge}}

**New**

Release date: 1/18/2018

- Docker Engine upgraded to [Docker 18.01.0 CE](https://github.com/docker/docker-ce/releases/tag/v18.01.0-ce){: target="_blank" class="_"}


### 17.10 CE

**New**

Release date: 10/18/2017

- Docker Engine upgraded to [Docker 17.10.0 CE](https://github.com/docker/docker-ce/releases/tag/v17.10.0-ce){: target="_blank" class="_"}
- Editions container log to stdout instead of disk, preventing hdd fill-up
- Azure VHD mounts instance HDD allowing for smaller boot disks

## Template archive

If you are looking for templates from older releases, check out the [template archive](/docker-for-azure/archive.md).
