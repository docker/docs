---
description: Release notes
keywords: aws, amazon, iaas, release
title: Docker for AWS Release Notes
---

{% include d4a_buttons.md %}

## Stable Channel

### 1.13.1-2 (stable)
Release date: 02/08/2017

{{aws_button_latest}}

**New**

- Docker Engine upgraded to [Docker 1.13.1](https://github.com/docker/docker/blob/master/CHANGELOG.md)

### 1.13.0-1 (stable)
Release date: 01/18/2017

**New**

- Docker Engine upgraded to [Docker 1.13.0](https://github.com/docker/docker/blob/master/CHANGELOG.md)
- Change ELB health check from TCP to HTTP

## Beta Channel

### 1.13.1-beta18
Release date: 02/16/2017

**New**

- Docker Engine upgraded to [Docker 1.13.1](https://github.com/docker/docker/blob/master/CHANGELOG.md)
- Added a second CloudFormation template that allows you to [install Docker for AWS into a pre-existing VPC](index.md#install-into-an-existing-vpc).
- Added Swarm wide support for [persistent storage volumes](persistent-data-volumes.md)
- Added the following engine labels
    - **os** (linux)
    - **region** (us-east-1, etc)
    - **availability_zone** (us-east-1a, etc)
    - **instance_type** (t2.micro, etc)
    - **node_type** (worker, manager)

### 1.13.1-rc2-beta17
Release date: 02/07/2017

**New**

- Docker Engine upgraded to [Docker 1.13.1-rc2](https://github.com/docker/docker/blob/master/CHANGELOG.md)

### 1.13.1-rc1-beta16
Release date: 02/01/2017

**New**

- Docker Engine upgraded to [Docker 1.13.1-rc1](https://github.com/docker/docker/blob/master/CHANGELOG.md)

### 1.13.0-rc5-beta15
Release date: 01/10/2017

**New**

- Docker Engine upgraded to [Docker 1.13.0-rc5](https://github.com/docker/docker/blob/master/CHANGELOG.md)

### 1.13.0-rc4-beta14
Release date: 12/21/2016

**New**

- Docker Engine upgraded to [Docker 1.13.0-rc4](https://github.com/docker/docker/blob/master/CHANGELOG.md)
- Previously we always only used 2 Availability Zones per region, even if the region had more. We now dynamically pick the best number of Availability Zones to use based on the region. If a region only has two AZs it will only use 2. If it has three or more, it will use 3
- Changed the AutoScaleGroup HealthCheck from an EC2 check to an ELB check
- Removed password prompt when ssh key is invalid
- Added new Canada Central region `ca-central-1`
- Added new London region `eu-west-2`
- Made recovery improvements when primary swarm node crashes

### 1.13.0-rc3-beta13
Release date: 12/06/2016

**New**

- Docker Engine upgraded to [Docker 1.13.0-rc3](https://github.com/docker/docker/blob/master/CHANGELOG.md)
- New option to decide if you want to send container logs to CloudWatch. (previously it was always on)
- SSH access has been added to the worker nodes
- The Docker daemon no longer listens on port 2375
- Added a `swarm-exec` to execute a docker command across all of the swarm nodes. See [Executing Docker commands in all swarm nodes](/docker-for-aws/deploy.md#execute-docker-commands-in-all-swarm-nodes) for more details.

### 1.13.0-rc2-beta12
Release date: 11/23/2016

**New**

- Docker Engine upgraded to [Docker 1.13.0-rc2](https://github.com/docker/docker/blob/master/CHANGELOG.md)
- New option to cleanup unused resources on your Swarm using new Docker prune command available in 1.13
- New option to pick the size of the ephemeral storage volume size on workers and managers
- New option to pick the disk type for the ephemeral storage on workers and managers
- Changed the Cloud Watch container log name from container "ID" to "Container Name-ID"


### 1.13.0-rc1-beta11

Release date: 11/17/2016

**New**

- Docker Engine upgraded to [Docker 1.13.0-rc1](https://github.com/docker/docker/blob/master/CHANGELOG.md)
- Changes to port 2375 access. For security reasons we locked down access to port 2375 in the following ways.
    - You can't connect to port 2375 on managers from workers (changed)
    - You can't connect to port 2375 on workers from other workers (changed)
    - You can't connect to port 2375 on managers and workers from the public internet (no change)
    - You can connect to port 2375 on workers from managers (no change)
    - You can connect to port 2375 on managers from other managers (no change)
- Added changes to the way we manage swarm tokens to make it more secure.

**Important**

- Due to some changes with the IP ranges in the subnets in Beta10, it will not be possible to upgrade from beta 10 to beta 11. You will need to start from scratch using beta11. We are sorry for any issues this might cause. We needed to make the change, and it was decided it was best to do it now, while still in private beta to limit the impact.


### 1.12.3-beta10

Release date: 10/27/2016

**New**

- Docker Engine upgraded to Docker 1.12.3
- Fixed the shell container that runs on the managers, to remove a ssh host key that was accidentally added to the image.
This could have led to a potential man in the middle (MITM) attack. The ssh host key is now generated on host startup, so that each host has its own key.
- The SSH ELB for connecting to the managers by SSH has been removed because it is no longer possible to SSH into the managers without getting a security warning
- You can connect to each manager using SSH by following our deploy [guide](/docker-for-aws/deploy.md)
- Added new region us-east-2 (Ohio)
- Fixed some bugs related to upgrading the swarm
- SSH keypair is now a required field in CloudFormation
- Fixed networking dependency issues in CloudFormation template that could result in a stack failure.

### 1.12.2-beta9

Release date: 10/12/2016

**New**

- Docker Engine upgraded to Docker 1.12.2
- Can better handle scaling swarm nodes down and back up again
- Container logs are now sent to CloudWatch
- Added a diagnostic command (docker-diagnose), to more easily send us diagnostic information in case of errors for troubleshooting
- Added sudo support to the shell container on manager nodes
- Change SQS default message timeout to 12 hours from 4 days
- Added support for region 'ap-south-1': Asia Pacific (Mumbai)

**Deprecated:**

- Port 2375 will be closed in next release. If you relay on this being open, please plan accordingly.

### 1.12.2-RC3-beta8

Release date: 10/06/2016

**New**

 * Docker Engine upgraded to 1.12.2-RC3

### 1.12.2-RC2-beta7

Release date: 10/04/2016

**New**

 * Docker Engine upgraded to 1.12.2-RC2

### 1.12.2-RC1-beta6

Release date: 9/29/2016

**New**

 * Docker Engine upgraded to 1.12.2-RC1


### 1.12.1-beta5

Release date: 8/18/2016

**New**

 * Docker Engine upgraded to 1.12.1

**Errata**

 * Upgrading from previous Docker for AWS versions to 1.12.0-beta4 is not possible because of RC-incompatibilities between Docker 1.12.0 release candidate 5 and previous release candidates.

### 1.12.0-beta4

Release date: 7/28/2016

**New**

 * Docker Engine upgraded to 1.12.0

**Errata**

 * Upgrading from previous Docker for AWS versions to 1.12.0-beta4 is not possible because of RC-incompatibilities between Docker 1.12.0 release candidate 5 and previous release candidates.

### 1.12.0-rc5-beta3

(internal release)

### 1.12.0-rc4-beta2

Release date: 7/13/2016

**New**

 * Docker Engine upgraded to 1.12.0-rc4
 * EC2 instance tags
 * Beta Docker for AWS sends anonymous analytics

**Errata**

 * When upgrading, old Docker nodes may not be removed from the swarm and show up when running `docker node ls`. Marooned nodes can be removed with `docker node rm`

### 1.12.0-rc3-beta1

**New**

 * First release of Docker for AWS!
 * CloudFormation based installer
 * ELB integration for running public-facing services
 * Swarm access with SSH
 * Worker scaling using AWS ASG

**Errata**

 * To assist with debugging, the Docker Engine API is available internally in the AWS VPC on TCP port 2375. These ports cannot be accessed from outside the cluster, but could be used from within the cluster to obtain privileged access on other cluster nodes. In future releases, direct remote access to the Docker API will not be available.
 * Likewise, swarm-mode is configured to auto-accept both manager and worker nodes inside the VPC. This policy will be changed to be more restrictive by default in the future.
