---
description: How to configure Settings Management for Docker Desktop using the Docker Admin Console
keywords: admin, controls, rootless, enhanced container isolation
title: Configure with the Docker Admin Console
---

>**Note**
>
>Settings Management is available to Docker Business customers only.

This page contains information for admins on how to configure Settings Management with the Docker Admin Console to specify and lock configuration parameters to create a standardized Docker Desktop environment across the organization.

EXPLAIN
 - can be done at org level
 - what creating a policy means, can apply to all users or select users


Note: if you have an admin-settings.json form, you can switch to admin console to make it easier (add corresponding tip to json file page)

## Prerequisites

- [Download and install Docker Desktop 4.13.0 or later](/desktop/release-notes.md).
- As an administrator, you need to [enforce
  sign-in](/security/for-admins/enforce-sign-in/_index.md). This is
because the Enhanced Container Isolation feature requires a Docker Business
subscription and therefore your Docker Desktop users must authenticate to your
organization for this configuration to take effect. 

## Setup a settings policy

1. Within the [Docker Admin Console](https://admin.docker.com/) navigate to the organization you want to define
2. 