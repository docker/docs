---
title: Manage Docker products
weight: 45
description: Learn how to manage access and usage for Docker products for your organization
keywords: organization, tools, products, product access, organization management
---

{{< summary-bar feature_name="Admin orgs" >}}

In this section, learn how to manage access and view usage of the Docker
products for your organization. For more detailed information about each
product, including how to set up and configure them, see the following manuals:

- [Docker Desktop](../../desktop/_index.md)
- [Docker Hub](../../docker-hub/_index.md)
- [Docker Build Cloud](../../build-cloud/_index.md)
- [Docker Scout](../../scout/_index.md)
- [Testcontainers Cloud](https://testcontainers.com/cloud/docs/#getting-started)
- [Docker Offload](../../offload/_index.md)

## Manage product access for your organization

Access to the Docker products included in your subscription is turned on by
default for all users. For an overview of products included in your
subscription, see
[Docker subscriptions and features](https://www.docker.com/pricing/).

{{< tabs >}}
{{< tab name="Docker Desktop" >}}

### Manage Docker Desktop access

To manage Docker Desktop access:

1. [Enforce sign-in](/manuals/enterprise/security/enforce-sign-in/_index.md).
1. Manage members [manually](./members.md) or use
[provisioning](/manuals/enterprise/security/provisioning/_index.md).

With sign-in enforced, only users who are a member of your organization can
use Docker Desktop after signing in.

{{< /tab >}}
{{< tab name="Docker Hub" >}}

### Manage Docker Hub access

To manage Docker Hub access, sign in to
[Docker Home](https://app.docker.com/) and configure [Registry Access Management](/manuals/enterprise/security/hardened-desktop/registry-access-management.md)
or [Image Access Management](/manuals/enterprise/security/hardened-desktop/image-access-management.md).

{{< /tab >}}
{{< tab name="Docker Build Cloud" >}}

### Manage Docker Build Cloud access

To initially set up and configure Docker Build Cloud, sign in to
[Docker Build Cloud](https://app.docker.com/build) and follow the
on-screen instructions.

To manage Docker Build Cloud access:

1. Sign in to [Docker Build Cloud](http://app.docker.com/build) as an
organization owner.
1. Select **Account settings**.
1. Select **Lock access to Docker Build Account**.

{{< /tab >}}
{{< tab name="Docker Scout" >}}

### Manage Docker Scout access

To initially set up and configure Docker Scout, sign in to
[Docker Scout](https://scout.docker.com/) and follow the on-screen instructions.

To manage Docker Scout access:

1. Sign in to [Docker Scout](https://scout.docker.com/) as an organization
owner.
1. Select your organization, then **Settings**.
1. To manage what repositories are enabled for Docker Scout analysis, select
**Repository settings**. For more information on,
see [repository settings](../../scout/explore/dashboard.md#repository-settings).
1. To manage access to Docker Scout for use on local images with Docker Desktop,
use [Settings Management](/manuals/enterprise/security/hardened-desktop/settings-management/_index.md)
and set `sbomIndexing` to `false` to disable, or to `true` to enable.

{{< /tab >}}
{{< tab name="Testcontainers Cloud" >}}

### Manage Testcontainers Cloud access

To initially set up and configure Testcontainers Cloud, sign in to
[Testcontainers Cloud](https://app.testcontainers.cloud/) and follow the
on-screen instructions.

To manage access to Testcontainers Cloud:

1. Sign in to the [Testcontainers Cloud](https://app.testcontainers.cloud/) and
select **Account**.
1. Select **Settings**, then **Lock access to Testcontainers Cloud**.

{{< /tab >}}
{{< tab name="Docker Offload" >}}

### Manage Docker Offload access

> [!NOTE]
>
> Docker Offload isn't included in the core Docker subscription plans. To make Docker Offload available, you must [sign
> up](https://www.docker.com/products/docker-offload/) and subscribe.

To manage Docker Offload access for your organization, use [Settings
Management](/manuals/enterprise/security/hardened-desktop/settings-management/_index.md):

1. Sign in to [Docker Home](https://app.docker.com/) as an organization owner.
1. Select **Admin Console** > **Desktop Settings Management**.
1. Configure the **Enable Docker Offload** setting to control whether Docker Offload features are available in Docker
   Desktop. You can configure this setting in five states:
   - **Always enabled**: Docker Offload is always enabled and users cannot disable it. The Offload
     toggle is always visible in Docker Desktop header. Recommended for VDI environments where local Docker execution is
     not possible.
   - **Enabled**: Docker Offload is enabled by default but users can disable it in Docker Desktop
     settings. Suitable for hybrid environments.
   - **Disabled**: Docker Offload is disabled by default but users can enable it in Docker Desktop
     settings.
   - **Always disabled**: Docker Offload is disabled and users cannot enable it. The option is
     visible but locked. Use when Docker Offload is not approved for organizational use.
   - **User defined**: No enforced default. Users choose whether to enable or disable Docker Offload in their
     Docker Desktop settings.
1. Select **Save**.

For more details on Settings Management, see the [Settings
reference](/manuals/enterprise/security/hardened-desktop/settings-management/settings-reference.md#enable-docker-offload).

{{< /tab >}}
{{< /tabs >}}

## Monitor product usage for your organization

To view usage for Docker products:

- Docker Desktop: View the **Insights** page in [Docker Home](https://app.docker.com/). For more details, see [Insights](./insights.md).
- Docker Hub: View the [**Usage** page](https://hub.docker.com/usage) in Docker Hub.
- Docker Build Cloud: View the **Build minutes** page in [Docker Build Cloud](http://app.docker.com/build).
- Docker Scout: View the [**Repository settings** page](https://scout.docker.com/settings/repos) in Docker Scout.
- Testcontainers Cloud: View the [**Billing** page](https://app.testcontainers.cloud/dashboard/billing) in Testcontainers Cloud.
- Docker Offload: View the **Offload** > **Offload overview** page in [Docker Home](https://app.docker.com/). For more details, see
  [Docker Offload usage and billing](/offload/usage/).

If your usage or seat count exceeds your subscription amount, you can
[scale your subscription](../../subscription/scale.md) to meet your needs.
