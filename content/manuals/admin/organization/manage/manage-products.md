---
title: Manage usage and access for Docker products
linkTitle: Product usage and access
weight: 50
description: Learn how to manage access and usage for Docker products for your organization
keywords: organization, product access, product usage, access control, docker desktop, docker hub, docker scout, docker build cloud, docker offload, testcontainers cloud
aliases:
  - /admin/organization/manage-products/
---

{{< summary-bar feature_name="Admin orgs" >}}

Use this page to learn how to control and monitor product access and usage
for your organization's members. If you're looking for setup and
configuration instructions, see each product's manual under
[What's next](#whats-next).

## Monitor product usage for your organization

You can monitor usage for Docker products across your organization. Use the
following table to learn where you can monitor organization usage:

| Product              | Monitor usage                                                                                                                                                                                    |
| -------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| Docker Desktop       | From [Docker Home](https://app.docker.com/), view the [**Insights**](../../insights.md) page.                                                                                                    |
| Docker Hub           | From Docker Hub, view the [**Usage** page](https://hub.docker.com/usage).                                                                                                                        |
| Docker Build Cloud   | From [Docker Build Cloud](http://app.docker.com/build), view the **Build minutes** page.                                                                                                         |
| Docker Scout         | From [Docker Home](https://app.docker.com/), select **Go to Scout** to view the [**Repository settings** page](https://scout.docker.com/settings/repos).                                         |
| Testcontainers Cloud | From [Docker Home](https://app.docker.com/), select **Go to Testcontainers Cloud**, then select the menu icon. Go to the [**Billing** page](https://app.testcontainers.cloud/dashboard/billing). |
| Docker Offload       | From [Docker Home](https://app.docker.com/), select **Offload**, then **Offload activity**. See [Docker Offload usage and billing](/manuals/offload/usage.md) for more details.                  |

To learn about the included usage across Docker plans, see
[Docker subscriptions and features](https://www.docker.com/pricing?ref=Docs&refAction=DocsAdminManageProducts).

## Control access for your organization

Organization members can access Docker products that your organization is
subscribed to by default. When signed in as an organization owner, you can
use the following procedures to control access for all members.

### Docker Desktop access

To manage Docker Desktop access:

1. [Enforce sign-in](/manuals/enterprise/security/enforce-sign-in/_index.md).
1. Manage members [manually](./members.md) or use
   [provisioning](/manuals/enterprise/security/provisioning/_index.md).

With sign-in enforced, only users who are a member of your organization can
use Docker Desktop after signing in.

### Docker Hub access

To manage Docker Hub access:

1. Sign in to [Docker Home](https://app.docker.com/), then select **Docker
   Desktop**.
1. Select **Registry Access** to configure
   [Registry Access Management](/manuals/enterprise/security/hardened-desktop/registry-access-management.md).
1. Select **Image Access** to control
   [Image Access Management](/manuals/enterprise/security/hardened-desktop/image-access-management.md).

### Docker Build Cloud access

To initially set up and configure Docker Build Cloud, sign in to
[Docker Build Cloud](https://app.docker.com/build) and follow the
on-screen instructions.

To manage Docker Build Cloud access:

1. Sign in to [Docker Home](https://app.docker.com/), then select
   [Docker Build Cloud](http://app.docker.com/build).
1. Select **Account settings**.
1. Select **Lock access to Docker Build Account**.

### Docker Scout access

To initially set up and configure Docker Scout, sign in to
[Docker Scout](https://scout.docker.com/) and follow the on-screen
instructions.

To manage Docker Scout access:

1. Sign in to [Docker Home](https://app.docker.com/), then select
   [Docker Scout](https://scout.docker.com/).
1. Select your organization, then **Settings**.
1. To manage what repositories are enabled for Docker Scout analysis, select
   **Repository settings**. For more information, see
   [repository settings](../../../scout/explore/dashboard.md#repository-settings).
1. To manage access to Docker Scout for use on local images with Docker
   Desktop, use
   [Settings Management](/manuals/enterprise/security/hardened-desktop/settings-management/_index.md)
   and set `sbomIndexing` to `false` to disable, or to `true` to enable.

### Testcontainers Cloud access

To initially set up and configure Testcontainers Cloud, sign in to
[Testcontainers Cloud](https://app.testcontainers.cloud/) and follow the
on-screen instructions.

To manage access to Testcontainers Cloud:

1. Sign in to the [Testcontainers Cloud](https://app.testcontainers.cloud/),
   then select the menu icon.
1. Select **Account**, then **Settings**.
1. Choose **Lock access to Testcontainers Cloud**.

### Docker Offload access

> [!NOTE]
>
> Docker Offload isn't included in the core Docker subscription plans. To
> make Docker Offload available, you must
> [contact sales](https://www.docker.com/products/docker-offload/) and
> subscribe.

To manage Docker Offload access for your organization, use [Settings
Management](/manuals/enterprise/security/hardened-desktop/settings-management/_index.md):

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
reference](/manuals/enterprise/security/hardened-desktop/settings-management/settings-reference.md#enable-docker-offload).

## What's next

For more detailed information about each product, including how to set up
and configure them, see the following manuals:

- [Docker Desktop](../../../desktop/_index.md)
- [Docker Hub](../../../docker-hub/_index.md)
- [Docker Build Cloud](../../../build-cloud/_index.md)
- [Docker Scout](../../../scout/_index.md)
- [Testcontainers Cloud](https://testcontainers.com/cloud/docs/#getting-started)
- [Docker Offload](../../../offload/_index.md)
