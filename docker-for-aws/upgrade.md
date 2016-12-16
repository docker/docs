---
description: Upgrading your stack
keywords: aws, amazon, iaas, tutorial
title: Docker for AWS Upgrades
---

Docker for AWS has support upgrading from one beta version to the next. Upgrades are done by applying a new version of the AWS Cloudformation template that powers Docker for Azure. Depending on changes in the next version, an upgrade involves:

 * Changing the AMI backing manager and worker nodes (the Docker engine ships in the AMI)
 * Upgrading service containers
 * Changing the resource setup in the VPC that hosts Docker for AWS

To be notified of updates, submit your email address at [https://beta.docker.com/](https://beta.docker.com/).

## Prerequisites

 * We recommend only attempting upgrades of swarms with at least 3 managers. A 1-manager swarm may not be able to maintain quorum during the upgrade
 * Upgrades are only supported from one version to the next version, for example beta-11 to beta-12. Skipping a version during an upgrade is not supported. For example, upgrading from beta-10 to beta-12 is not supported. Downgrades are not tested.
 
## Upgrading

If you submit your email address at [https://beta.docker.com/](beta.docker.com) Docker will notify you of new releases by email. New releases are also posted on the [Release Notes](https://beta.docker.com/docs/aws/release-notes/) page.

To initiate an update, use either the AWS Console of the AWS cli to initiate a stack update. Use the S3 template URL for the new release and complete the update wizard. This will initiate a rolling upgrade of the Docker swarm, and service state will be maintained during and after the upgrade. Appropriately scaled services should not experience downtime during an upgrade.

![Upgrade in AWS console](/img/cloudformation_update.png)

Note that single containers started (for example) with `docker run -d` are **not** preserved during an upgrade. This is because the're not Docker Swarm objects, but are known only to the individual Docker engines.

## Changing instance sizes and other template parameters

In addition to upgrading Docker for AWS from one version to the next you can also use the AWS Update Stack feature to change template parameters such as worker count and instance type. Changing manager count is **not** supported.
