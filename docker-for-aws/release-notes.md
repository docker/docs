---
description: Release notes
keywords: aws, amazon, iaas, release, edge, stable
title: Docker for AWS release notes
---

{% include d4a_buttons.md %}

## Enterprise Edition
[Docker Enterprise Edition Lifecycle](https://success.docker.com/Policies/Maintenance_Lifecycle){: target="_blank" class="_"}

[Deploy Docker Enterprise Edition (EE) for AWS](https://store.docker.com/editions/enterprise/docker-ee-aws?tab=description){: target="_blank" class="button outline-btn blank_"}

### 17.06 EE

- Docker engine 17.06 EE
- For Std/Adv external logging has been removed, as it is now handled by [UCP](https://docs.docker.com/datacenter/ucp/2.0/guides/configuration/configure-logs/){: target="_blank" class="_"}
- UCP 2.2.3
- DTR 2.3.3

### 17.03 EE

- Docker engine 17.03 EE
- UCP 2.1.5
- DTR 2.2.7


> **Note** Starting with 18.02.0-CE EFS encryption option has been removed to prevent the [recreation of the EFS volume](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html){: target="_blank" class="_"}.


## Stable channel

### 18.03 CE

{{aws_blue_latest}}

Release date: 3/21/2018

- Docker Engine upgraded to [Docker 18.03.0 CE](https://github.com/docker/docker-ce/releases/tag/v18.03.0-ce){: target="_blank" class="_"}
- [Elastic Network Interfaces](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-eni.html){: target="_blank" class="_"} enabled in the AMI kernel

### 17.12.1 CE

Release date: 3/1/2018

- Docker Engine upgraded to [Docker 17.12.1 CE](https://github.com/docker/docker-ce/releases/tag/v17.12.1-ce){: target="_blank" class="_"}
- Added baked-in rules for ECR IAM role

### 17.12 CE

Release date: 1/9/2018

- Docker Engine upgraded to [Docker 17.12.0 CE](https://github.com/docker/docker-ce/releases/tag/v17.12.0-ce){: target="_blank" class="_"}
- Kernel patch to mitigates Meltdown attacks ( CVE-2017-5754) and enable KPTI

> **Note** There was an issue in LinuxKit that prevented containers from [starting after a machine reboot](https://github.com/moby/moby/issues/36189){: target="_blank" class="_"}.

### 17.09 CE

Release date: 10/6/2017

- Docker Engine upgraded to [Docker 17.09.0 CE](https://github.com/docker/docker-ce/releases/tag/v17.09.0-ce){: target="_blank" class="_"}
- CloudStor EBS updates
- Moby mounts for early reboot support

### 17.06.1 CE

Release date: 08/17/2017

**New**

- Docker Engine upgraded to [Docker 17.06.1 CE](https://github.com/docker/docker-ce/releases/tag/v17.06.1-ce){: target="_blank" class="_"}
- Improvements to CloudStor support
- Added SSL support at the LB level

### 17.06.0 CE

Release date: 06/28/2017

**New**

- Docker Engine upgraded to [Docker 17.06.0 CE](https://github.com/docker/docker-ce/releases/tag/v17.06.0-ce){: target="_blank" class="_"}
- Fixed an issue with load balancer controller that caused the ELB health check to fail.
- Added VPCID output when a VPC is created
- Added CloudStor support (EFS (in regions that support EFS), and EBS) for [persistent storage volumes](persistent-data-volumes.md)
- Added CloudFormation parameter to enable/disable CloudStor
- Changed the AutoScaleGroup Manager max size to 6, so that it correctly upgrades with 5 managers
- Added lambda support for Mumbai
- Removed the ELB Name to allow for longer stack names
- Added i3 EC2 instance types
- [Bring your own VPC] Added a VPC CIDR Parameter

### 17.03.1 CE

Release date: 03/30/2017

**New**

- Docker Engine upgraded to [Docker 17.03.1 CE](https://github.com/docker/docker/blob/master/CHANGELOG.md){: target="_blank" class="_"}
- Updated AZ for Sao Paulo

### 17.03.0 CE

Release date: 03/01/2017

**New**

- Docker Engine upgraded to [Docker 17.03.0 CE](https://github.com/docker/docker/blob/master/CHANGELOG.md){: target="_blank" class="_"}
- Added r4 EC2 instance types
- Added `ELBDNSZoneID` output to make it easier to interact with Route53


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

## Template archive

If you are looking for templates from older releases, check out the [template archive](/docker-for-aws/archive.md).
