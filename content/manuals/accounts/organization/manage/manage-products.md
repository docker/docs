---
title: Manage usage and access for Docker products
linkTitle: Products
weight: 30
description: >-
  Control which Docker products organization members can use, and find
  where to check usage
keywords: organization, product access, product usage, access control, seats, licenses, docker desktop, docker hub, docker scout, docker build cloud, docker offload, testcontainers cloud, docker sandboxes, ai governance
toc_max: 2
aliases:
  - /admin/organization/manage/manage-products/
  - /admin/organization/manage-products/
---

{{< summary-bar feature_name="Admin orgs" >}}

Use this page to learn how to control and monitor product access and usage
for your organization's members.

## How product access and usage works

Members receive access to products through seats, licenses, and the controls
you configure:

- A member uses one seat on a Docker Team or Business
  [subscription](/manuals/subscription-billing/plans/_index.md). The seat
  includes the products in the plan, such as Docker Desktop and Docker Hub.
- Some products need a separate
  [license](/manuals/accounts/organization/manage/manage-licenses.md) that
  you assign to teams or members.
- After a member has a seat or license, the controls on this page set what
  they can do in each product, such as which Hub images they can pull,
  which repositories Scout analyzes, or whether Docker Offload is available
  in Docker Desktop.

Each product measures usage in its own units, such as seats, Hub pulls,
build minutes, or runtime minutes, and your subscription sets how much is
included. The Usage section for each product shows where to check that
amount, so you can top up your plan before you run out.

This page covers product access, but does not cover access tokens.
For organization access tokens,
see [Access tokens](/manuals/security/access-tokens/_index.md).

## Docker Desktop

### Access

To manage Docker Desktop access:

1. [Enforce sign-in](/manuals/desktop/enterprise/enforce-sign-in/_index.md).
1. Manage members [manually](./members.md) or use
   [provisioning](../../../security/provisioning/_index.md).

With sign-in enforced, only members of your organization can use Docker
Desktop after they sign in.

### Usage

Desktop insights shows how members use Docker Desktop, including active
users, builds, containers, images, and extensions.

1. Sign in to [Docker Home](https://app.docker.com/) and select your
   organization.
1. Select **Docker Desktop**, then
   [**Desktop insights**](/manuals/accounts/organization/insights.md).

## Docker Hub

Image Access and Registry Access are Docker Desktop controls. They limit
what members can pull and which registries Docker Desktop can reach.

### Access

1. Sign in to [Docker Home](https://app.docker.com/) and select your
   organization, then select **Docker Desktop**.
1. Select **Registry Access** to configure
   [Registry Access Management](/manuals/desktop/enterprise/hardened-desktop/registry-access-management.md).
1. Select **Image Access** to control
   [Image Access Management](/manuals/desktop/enterprise/hardened-desktop/image-access-management.md).

### Usage

Hub usage counts how many times members pull images and how much space
your organization's repositories take up.

1. Sign in to [Docker Home](https://app.docker.com/) and select your
   organization.
1. Select **Go to Hub**.
1. Select **Usage**, then **Pulls** or **Storage**. For example,
   `https://hub.docker.com/orgs/{namespace}/usage`.

## Docker Build Cloud

### Access

To set up Docker Build Cloud, sign in to
[Docker Build Cloud](https://app.docker.com/build) and follow the on-screen
instructions.

To manage Docker Build Cloud access:

1. Sign in to [Docker Home](https://app.docker.com/), then select
   [Docker Build Cloud](https://app.docker.com/build).
1. Select **Account settings**.
1. Select **Lock access to Docker Build Account**.

> [!WARNING]
>
> **Lock access to Docker Build Account** removes the ability to use cloud
> builders or build minutes, and it purges build data.

### Usage

Build minutes are the time your organization spends on cloud builds.

1. Sign in to [Docker Home](https://app.docker.com/) and select your
   organization.
1. Select **Build Cloud**, then **Build minutes**.

## Docker Scout

### Access

To set up Docker Scout, sign in to [Docker Scout](https://scout.docker.com/)
and follow the on-screen instructions.

To manage Docker Scout access:

1. Sign in to [Docker Home](https://app.docker.com/), then select
   [Docker Scout](https://scout.docker.com/).
1. Select your organization, then **Settings**.
1. To manage which repositories Docker Scout analyzes, select
   **Repository settings**. See
   [Repository settings](../../../scout/explore/dashboard.md#repository-settings).
1. To manage Docker Scout for local images in Docker Desktop, use
   [Settings Management](/manuals/desktop/enterprise/hardened-desktop/settings-management/_index.md)
   and set `sbomIndexing` to `false` to turn it off, or to `true` to turn
   it on.

### Usage

Repository settings shows which repositories Docker Scout analyzes.

1. Sign in to [Docker Home](https://app.docker.com/) and select your
   organization.
1. Select **Go to Scout**.
1. Select **Settings**, then
   [**Repository settings**](../../../scout/explore/dashboard.md#repository-settings).

## Testcontainers Cloud

### Access

To set up Testcontainers Cloud, sign in to
[Testcontainers Cloud](https://app.testcontainers.cloud/) and follow the
on-screen instructions.

To manage access to Testcontainers Cloud:

1. Sign in to [Testcontainers Cloud](https://app.testcontainers.cloud/),
   then select the menu icon.
1. Select **Account**, then **Settings**.
1. Select **Lock access to Testcontainers Cloud**.

### Usage

Billing shows the Testcontainers Cloud runtime minutes included with your
subscription.

1. Sign in to [Docker Home](https://app.docker.com/) and select your
   organization.
1. Select **Go to Testcontainers Cloud**.
1. Select the menu icon, then
   [**Billing**](https://app.testcontainers.cloud/dashboard/billing).

## Docker Offload

### Access

To manage Docker Offload access for your organization, use [Settings
Management](/manuals/desktop/enterprise/hardened-desktop/settings-management/_index.md):

1. Sign in to [Docker Home](https://app.docker.com/), then select **Docker
   Desktop**.
1. Select **Settings Management**.
1. Configure the **Enable Docker Offload** setting to control whether
   Docker Offload features are available in Docker Desktop. You can
   configure this setting in five states:
   - **Always enabled**: Docker Offload is always enabled and users cannot
     disable it. The Offload toggle is always visible in the Docker Desktop
     header. Recommended for VDI environments where local Docker execution
     is not possible.
   - **Enabled**: Docker Offload is enabled by default but users can
     disable it in Docker Desktop settings. Suitable for hybrid
     environments.
   - **Disabled**: Docker Offload is disabled by default but users can
     enable it in Docker Desktop settings.
   - **Always disabled**: Docker Offload is disabled and users cannot
     enable it. The option is visible but locked. Use when Docker Offload
     is not approved for organizational use.
   - **User defined**: No enforced default. Users choose whether to enable
     or disable Docker Offload in their Docker Desktop settings.
1. Select **Save**.

For more details on Settings Management, see the [Settings
reference](/manuals/desktop/enterprise/hardened-desktop/settings-management/settings-reference.md#enable-docker-offload).

### Usage

Offload activity shows Docker Offload sessions, including how long they
ran and which images they used.

1. Sign in to [Docker Home](https://app.docker.com/) and select your
   organization.
1. Select **Docker Offload** or **Offload**, then
   [**Offload activity**](../../../offload/usage.md).

## Next steps

If your usage or seat count exceeds your subscription amount, you can
[add seats](./manage-seats.md) or
[view available Docker plans](../../../subscription-billing/plans/_index.md).

For more detailed information about each product, including how to set up
and configure them, see the following manuals:

- [Docker Desktop](../../../desktop/_index.md)
- [Docker Hub](../../../docker-hub/_index.md)
- [Docker Build Cloud](../../../build-cloud/_index.md)
- [Docker Scout](../../../scout/_index.md)
- [Testcontainers Cloud](https://testcontainers.com/cloud/docs/#getting-started)
- [Docker Offload](../../../offload/_index.md)
- [Docker Sandboxes](../../../ai/sandboxes/_index.md)
