---
description: Upgrading your stack
keywords: aws, amazon, iaas, tutorial
title: Docker for AWS upgrades
---

To upgrade, apply a new version of the AWS Cloudformation template that powers
Docker for AWS. Depending on changes in the next version, an upgrade involves:

 * Changing the AMI backing manager and worker nodes (the Docker engine
   ships in the AMI)
 * Upgrading service containers
 * Changing the resource setup in the VPC that hosts Docker for AWS

## Prerequisites

 * We recommend only attempting upgrades of swarms with at least 3 managers.
 A 1-manager swarm can't maintain quorum during the upgrade.

 * You can only upgrade one version at a time. Skipping a version during
  an upgrade is not supported. Downgrades are not tested.

## Upgrading

New releases are announced on [Release Notes](release-notes.md) page.

To initiate an update, use either the AWS Console or the AWS cli to initiate a
stack update. Use the S3 template URL for the new release and complete the
update wizard. This initiates a rolling upgrade of the Docker swarm, and
service state is maintained during and after the upgrade. Appropriately
scaled services should not experience downtime during an upgrade.

![Upgrade in AWS console](img/cloudformation_update.png)

Single containers started (for example) with `docker run -d` are
**not** preserved during an upgrade. This is because they're not Docker Swarm
objects, but are known only to the individual Docker engines.

> **Note** Current Docker versions, up to 18.02.0-ce, will [recreate the EFS volume](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html){: target="_blank" class="_"} when performing a stack upgrade.

## Changing instance sizes and other template parameters

In addition to upgrading Docker for AWS from one version to the next you can
also use the AWS Update Stack feature to change template parameters such as
worker count and instance type. Changing manager count is **not** supported.
