---
title: Manage usage and access for Docker products
linkTitle: Products
weight: 30
description: >-
  Control which Docker products organization members can use, and find
  where to check usage
keywords: organization, product access, product usage, access control, seats, licenses, docker desktop, docker hub, docker scout, docker build cloud, docker offload, testcontainers cloud, docker sandboxes, ai governance
aliases:
  - /admin/organization/manage/manage-products/
  - /admin/organization/manage-products/
---

{{< summary-bar feature_name="Admin orgs" >}}

Use this page to learn how to control and monitor product access and usage
for your organization's members. If you're looking for setup and
configuration instructions, see each product's manual under
[Next steps](#next-steps).

## How product access works

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

This page covers the product controls. It does not cover credentials. For
personal and organization access tokens used to authenticate to Docker Hub,
see [Access tokens](/manuals/security/access-tokens/_index.md).

## Control access for your organization

Organization members can access Docker products that your organization is
subscribed to by default. When signed in as an organization owner, you can
use the following procedures to control access for all members.

### Docker Desktop access

To manage Docker Desktop access:

1. [Enforce sign-in](/manuals/desktop/enterprise/enforce-sign-in/_index.md).
1. Manage members [manually](./members.md) or use
   [provisioning](../../../security/provisioning/_index.md).

With sign-in enforced, only members of your organization can use Docker
Desktop after they sign in.

### Docker Hub access

To manage Docker Hub access:

1. Sign in to [Docker Home](https://app.docker.com/) and select your
   organization, then select **Docker Desktop**.
1. Select **Registry Access** to configure
   [Registry Access Management](/manuals/desktop/enterprise/hardened-desktop/registry-access-management.md).
1. Select **Image Access** to control
   [Image Access Management](/manuals/desktop/enterprise/hardened-desktop/image-access-management.md).

### Docker Build Cloud access

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

### Docker Scout access

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

### Testcontainers Cloud access

To set up Testcontainers Cloud, sign in to
[Testcontainers Cloud](https://app.testcontainers.cloud/) and follow the
on-screen instructions.

To manage access to Testcontainers Cloud:

1. Sign in to [Testcontainers Cloud](https://app.testcontainers.cloud/),
   then select the menu icon.
1. Select **Account**, then **Settings**.
1. Select **Lock access to Testcontainers Cloud**.

### Docker Offload access

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

## Monitor product usage for your organization

You can monitor usage for Docker products across your organization. Use the
following table to learn where you can monitor organization usage:

| Product              | Monitor usage                                                                                                                                                                                                                                                            |
| -------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| Docker Desktop       | From [Docker Home](https://app.docker.com/), select your organization, then **Docker Desktop**, then [**Desktop insights**](/manuals/accounts/organization/insights.md).                                                                                                 |
| Docker Hub           | From [Docker Home](https://app.docker.com/), select your organization, then **Go to Hub**, then **Usage**, then **Pulls** or **Storage**. For example, `https://hub.docker.com/orgs/{namespace}/usage`.                                                                  |
| Docker Build Cloud   | From [Docker Home](https://app.docker.com/), select your organization, then **Build Cloud**, then **Build minutes**.                                                                                                                                                      |
| Docker Scout         | From [Docker Home](https://app.docker.com/), select your organization, then **Go to Scout**, then **Settings**, then [**Repository settings**](https://scout.docker.com/settings/repos).                                                                                  |
| Testcontainers Cloud | From [Docker Home](https://app.docker.com/), select your organization, then **Go to Testcontainers Cloud**, then the menu icon, then [**Billing**](https://app.testcontainers.cloud/dashboard/billing).                                                                   |
| Docker Offload       | From [Docker Home](https://app.docker.com/), select your organization, then **Docker Offload** or **Offload**, then [**Offload activity**](../../../offload/usage.md).                                                                                                    |

If your usage or seat count exceeds your subscription amount, you can
[add seats](./manage-seats.md) or [view available Docker plans](../../../subscription-billing/plans/_index.md) to meet your needs.

## Next steps

For more detailed information about each product, including how to set up
and configure them, see the following manuals:

- [Docker Desktop](../../../desktop/_index.md)
- [Docker Hub](../../../docker-hub/_index.md)
- [Docker Build Cloud](../../../build-cloud/_index.md)
- [Docker Scout](../../../scout/_index.md)
- [Testcontainers Cloud](https://testcontainers.com/cloud/docs/#getting-started)
- [Docker Offload](../../../offload/_index.md)
- [Docker Sandboxes](../../../ai/sandboxes/_index.md)
