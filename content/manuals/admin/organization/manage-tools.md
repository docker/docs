---
title: Manage Docker tools
weight: 45
description: Learn how to manage Docker tools for your organization
keywords: organization, tools
---

Docker subscriptions come with access to all Docker tools. In this section, learn how to manage access and view usage of the tools for your organization. For more detailed information about each tool, including how to set up and configure them, see the following manuals:

- [Docker Build Cloud](../../build-cloud/_index.md)
- [Docker Desktop](../../desktop/_index.md)
- [Docker Hub](../../docker-hub/_index.md)
- [Docker Scout](../../scout/_index.md)
- [Testcontainers Cloud](https://testcontainers.com/cloud/docs/#getting-started)

## Manage access to Docker tools

Docker Hub, Docker Build Cloud, and Docker Scout access are enabled by default
for all users with a paid subscription, excluding legacy subscriptions.
Testcontainers Cloud access by default is disabled. Use the following sections
to learn how to enable or disable access for these tools.

### Manage access to Docker Build Cloud

To manage access for Docker Build Cloud, sign in to [Docker Build
Cloud](http://app.docker.com/build) as an organization owner, select **Account
settings**, and then manage access under **Lock Docker Build Cloud**.

### Manage access to Docker Scout

To manage access for using Docker Scout on remote repositories, configure
[integrations](../../scout/explore/dashboard.md#integrations) and [repository
settings](../../scout/explore/dashboard.md#repository-settings) in the [Docker
Scout Dashboard](https://scout.docker.com/).

To manage access for using Docker Scout on local images in Docker Desktop, use
[Settings
Management](../../security/for-admins/hardened-desktop/settings-management/_index.md)
and set `sbomIndexing`.

### Manage access to Docker Hub

To manage access for Docker Hub, you can use [Registry Access
Management](../../security/for-admins/hardened-desktop/registry-access-management.md)
or [Image Access
Management](../../security/for-admins/hardened-desktop/image-access-management.md).

### Manage access to Testcontainers Cloud

To manage access for Testcontainers Cloud, sign in to [Testcontainers
Cloud](https://app.testcontainers.cloud/), select **Account**, and then select
**Users** to manage your users who have access.

## View Docker tool usage

You can view your organization's Docker tool usage and then [scale your
subscription](../../subscription/scale.md) to meet your needs. View usage for
the tools on the following pages:

- Docker Build Cloud: View the **Build minutes** page in [Docker Build Cloud
  Dashboard](http://app.docker.com/build).

- Docker Scout: View the **Billing settings** page in the [Docker Scout
  Dashboard](https://scout.docker.com/settings/billing).

- Docker Hub: View the **Usage** page in [Docker Hub](https://hub.docker.com/usage).

- Testcontainers Cloud: View the **Billing** page in the [Testcontainers Cloud
  Dashboard](https://app.testcontainers.cloud/dashboard/billing).

- Docker Desktop: View the **Insights** page in the [Docker Admin
  Console](https://app.docker.com/admin).