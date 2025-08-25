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
- [Docker Offload](/manuals/offload/_index.md)
- [Docker Build Cloud](../../build-cloud/_index.md)
- [Docker Scout](../../scout/_index.md)
- [Testcontainers Cloud](https://testcontainers.com/cloud/docs/#getting-started)

## Manage product access for your organization

Access to the Docker products included in your subscription is turned on by
default for all users. For an overview of products included in your
subscription, see
[Docker subscriptions and features](/manuals/subscription/details.md).

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
{{< tab name="Docker Offload" >}}

> [!NOTE]
>
> Docker Build Cloud settings are now bundled under the Docker Offload product.
To learn more about Docker Offload, see [Docker Offload](/manuals/offload/_index.md).

### Enable Docker Offload

Enabling Docker Offload enables the feature for all members of your
organization. Docker Offload is enabled by default.

To enable Docker Offload:

1. Sign in to [Docker Home](http://app.docker.com/) as an
organization owner.
1. Select **Admin Console** > **Offload access**.
1. In the **Configuration** drop-down, select **Enabled**.
1. Select **Save**.

### Disable Docker Offload

Disabling Docker Offload removes the ability for your organization members or
CI/CD pipelines to use offload resources or consume offload minutes.

To disable Docker Offload:

1. Sign in to [Docker Home](http://app.docker.com/) as an
organization owner.
1. Select **Admin Console** > **Offload access**.
1. In the **Configuration** drop-down, select **Disabled**.
1. Select **Save**.

### Configure GPU access for Docker Offload

GPU access allows developers to use Run Docker Model Runner, MCP tools,
or other GPU accelerated containers when using cloud mode.

To configure GPU access:

1. Sign in to [Docker Home](http://app.docker.com/) as an
organization owner.
1. Select **Admin Console** > **Offload access**.
1. Use the **Allow GPU access** checkbox to toggle the feature on or off.
1. Select **Save**.

### Lock Docker Build Cloud

Locking access to Docker Build Cloud removes the ability for anybody in your
organization to utilize any cloud builders or consume build minutes.

To lock Docker Build Cloud:

1. Sign in to [Docker Home](http://app.docker.com/) as an
organization owner.
1. Select **Offload** > **Account settings**.
1. Select **Lock access to Docker Build Account**.

### Purge cloud data

You can permanently delete all build-related data from your
Docker Build Cloud account. This action removes build history, build details,
usage data, and the build cache.

To purge cloud data:

1. Sign in to [Docker Home](http://app.docker.com/) as an
organization owner.
1. Select **Offload** > **Account settings**.
1. Select **Purge cloud data**.

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
{{< /tabs >}}

## Monitor product usage for your organization

To view usage for Docker products:

- Docker Desktop: View the **Insights** page in [Docker Home](https://app.docker.com/). For more details, see [Insights](./insights.md).
- Docker Hub: View the [**Usage** page](https://hub.docker.com/usage) in Docker Hub.
- Docker Build Cloud: View the **Build minutes** page in [Docker Build Cloud](http://app.docker.com/build).
- Docker Scout: View the [**Repository settings** page](https://scout.docker.com/settings/repos) in Docker Scout.
- Testcontainers Cloud: View the [**Billing** page](https://app.testcontainers.cloud/dashboard/billing) in Testcontainers Cloud.

If your usage or seat count exceeds your subscription amount, you can
[scale your subscription](../../subscription/scale.md) to meet your needs.
