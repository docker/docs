---
description: Why Docker EE for IBM Cloud?
keywords: ibm cloud, ibm, iaas, why
title: Why Docker EE for IBM Cloud?
---

Docker Enterprise Edition for IBM Cloud (Beta) was created and is being actively developed to ensure that Docker users can enjoy a fantastic out-of-the-box experience on Docker for enterprise-grade workloads. It is now available as a beta.

As an informed user, you might be curious to know what this project offers you for running your development, staging, or production workloads.

## Native to Docker
Docker EE for IBM Cloud provides a Docker-native solution that you can use to avoid operational complexity and using unneeded additional APIs to the Docker stack.

Docker EE for IBM Cloud allows you to interact with Docker directly (including native Docker orchestration), instead of distracting you with the need to navigate extra layers on top of Docker. You can focus instead on the thing that matters most: running your workloads. This helps you and your team to deliver more value to the business faster, to speak one common "language", and to have fewer details to keep in your head at once.

The skills that you and your team have already learned, and continue to learn, using Docker on the desktop or elsewhere automatically carry over to using Docker EE for IBM Cloud. The added consistency across clouds also helps to ensure that a migration or multi-cloud strategy is easier to accomplish in the future, if desired.

## Skip the boilerplate and maintenance work
Docker EE for IBM Cloud bootstraps all of the recommended infrastructure to start using Docker on IBM Cloud automatically. You don't need to worry about rolling your own instances, security groups, or load balancers when using Docker EE for IBM Cloud.

Likewise, setting up and using Docker swarm mode functionality for container orchestration is managed across the cluster's lifecycle when you use Docker EE for IBM Cloud. Docker has already coordinated the various bits of automation you would otherwise need to glue together on your own to bootstrap Docker swarm mode on the platforms. When the cluster is finished booting, you can jump right in and start running `docker service` commands to schedule tasks for your worker nodes.

## Self-cleaning and self-healing
Even the most conscientious admin can be caught off guard by issues such as exhaustive logging or the Linux kernel unexpectedly ending memory-hungry processes. In Docker EE for IBM Cloud, your cluster is resilient to a variety of such issues by default.

You can enable or disable logging for swarms, so chatty logs don't use up all of your disk space. Likewise, the "system prune" option allows you to ensure unused Docker resources such as old images are cleaned up automatically.

The lifecycle of nodes is managed using InfraKit, so that if a node enters an unhealthy state for unforeseen reasons, the node is removed from service and replaced automatically. Container tasks that were running on the unhealthy node are rescheduled.

You can breathe easier as these self-cleaning and self-healing properties reduce the risk of downtime.

## Logging native to the platforms
Centralized logging is a critical component of many modern infrastructure stacks. To have these logs indexed and searchable proves invaluable for debugging application and system issues as they come up. With Docker EE for IBM Cloud, you can enable seamless logging to you IBM Cloud account.

# Try it today
Ready to get started? [Try Docker for IBM Cloud today](https://www.ibm.com/us-en/marketplace/docker-for-ibm-cloud).

We'd be happy to hear your feedback via e-mail at docker-for-ibmcloud-beta@docker.com.
