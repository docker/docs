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

Docker Offload access is set inside the Docker Desktop Dashboard. To manage Docker Desktop settings for your
organization, you can use Setttings Management.

Follow the instructions at [Configure Settings Management with the Admin
Console](/enterprise/security/hardened-desktop/settings-management/configure-admin-console/) and set the **Enable Docker
Offload** setting to your desired value.

{{< /tab >}}
{{< /tabs >}}

## Monitor product usage for your organization

To view usage for Docker products:

- Docker Desktop: View the **Insights** page in [Docker Home](https://app.docker.com/). For more details, see [Insights](./insights.md).
- Docker Hub: View the [**Usage** page](https://hub.docker.com/usage) in Docker Hub.
- Docker Build Cloud: View the **Build minutes** page in [Docker Build Cloud](http://app.docker.com/build).
- Docker Scout: View the [**Repository settings** page](https://scout.docker.com/settings/repos) in Docker Scout.
- Testcontainers Cloud: View the [**Billing** page](https://app.testcontainers.cloud/dashboard/billing) in Testcontainers Cloud.
- Docker Offload: View the **Docker Offload** > **Usage summary** page in [Docker Billing](https://app.docker.com/billing).

If your usage or seat count exceeds your subscription amount, you can
[scale your subscription](../../subscription/scale.md) to meet your needs.
