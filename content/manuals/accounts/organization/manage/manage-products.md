---
title: Manage usage and access for Docker products
linkTitle: Products
weight: 30
description: Where to control access and monitor usage for Docker products in your organization
keywords: organization, product access, product usage, access control, docker desktop, docker hub, docker scout, docker build cloud, docker offload, testcontainers cloud, docker sandboxes, ai governance
aliases:
  - /admin/organization/manage/manage-products/
  - /admin/organization/manage-products/
---

{{< summary-bar feature_name="Admin orgs" >}}

As an organization owner, you decide which Docker products members can
access and you monitor how much your organization consumes. Use this page
to learn where to change access in Docker Home and how to check usage. To
set up a product for your organization, see that product's documentation.

## Control access

Pick a product when you need to change who can use it or what they're
allowed to do.

### Docker Desktop

Restrict Docker Desktop to organization members who must sign in.

1. [Enforce sign-in](/manuals/desktop/enterprise/enforce-sign-in/_index.md).
1. Add members
   [manually](/manuals/accounts/organization/manage/members.md) or with
   [provisioning](/manuals/security/provisioning/_index.md).

### Docker Hub

Limit which Hub images members can pull, which registries Desktop can
reach, and whether members can push to personal namespaces.

1. Sign in to [Docker Home](https://app.docker.com/) and select your
   organization.
1. Select **Docker Desktop**.
1. Select the control you want to change:
   - **Image Access**, to limit which Docker Hub images members can pull.
     See
     [Image Access Management](/manuals/desktop/enterprise/hardened-desktop/image-access-management.md).
   - **Registry Access**, to limit which registries Docker Desktop can
     reach. See
     [Registry Access Management](/manuals/desktop/enterprise/hardened-desktop/registry-access-management.md).
   - **Namespace access**, to control whether members can push to their
     personal namespaces. See
     [Namespace access control](/manuals/desktop/enterprise/hardened-desktop/namespace-access.md).

### Docker Build Cloud

Assign Build Cloud seats so members can use cloud builders.

1. Sign in to [Docker Home](https://app.docker.com/) and select your
   organization.
1. Select **Build Cloud**, then **User management**.
1. Invite members and assign seats for cloud builders.

> [!WARNING]
>
> Under **Build Cloud**, on **Account settings**, **Lock access to Docker
> Build Account** removes cloud builders and build minutes for the
> organization and purges build data.

### Docker Scout

Choose which repositories Scout analyzes, and turn Scout image analysis on
or off in Docker Desktop.

To choose which repositories Docker Scout analyzes:

1. Sign in to [Docker Home](https://app.docker.com/) and select your
   organization.
1. Select **Scout** or **Go to Scout**.
1. Select **Settings**, then **Repository settings**. See
   [Repository settings](/manuals/scout/explore/dashboard.md#repository-settings).

To control Scout image analysis in Docker Desktop:

1. Sign in to [Docker Home](https://app.docker.com/) and select your
   organization.
1. Select **Docker Desktop**, then **Settings Management**.
1. In a settings policy, set **SBOM indexing** (`sbomIndexing`). See
   [Enable Scout image analysis](/manuals/desktop/enterprise/hardened-desktop/settings-management/settings-reference.md#enable-scout-image-analysis).

### Testcontainers Cloud

Open Testcontainers Cloud from Docker Home. Access controls for
Testcontainers Cloud are managed in Testcontainers Cloud.

1. Sign in to [Docker Home](https://app.docker.com/) and select your
   organization.
1. Select **Testcontainers Cloud** or **Go to Testcontainers Cloud**.
1. Continue in the
   [Testcontainers Cloud documentation](https://testcontainers.com/cloud/docs/).

### Docker Offload

Gate Docker Offload with licenses, then control whether it is available in
Docker Desktop.

> [!NOTE]
>
> Docker Offload licenses are available for Docker Team and Docker Business
> plans. To add licenses,
> [contact sales](https://www.docker.com/pricing/contact-sales/).
> Organization owners assign licenses in Docker Home. See
> [Manage licenses](/manuals/accounts/organization/manage/manage-licenses.md).

To control Docker Offload in Docker Desktop:

1. Sign in to [Docker Home](https://app.docker.com/) and select your
   organization.
1. Select **Docker Desktop**, then **Settings Management**.
1. Create or edit a settings policy, and set **Enable Docker Offload**. See
   [Configure Settings Management](/manuals/desktop/enterprise/hardened-desktop/settings-management/configure-admin-console.md)
   and
   [Enable Docker Offload](/manuals/desktop/enterprise/hardened-desktop/settings-management/settings-reference.md#enable-docker-offload).

### Docker Sandboxes

Gate Sandboxes with AI Governance licenses and organization policies for
network, filesystem, and MCP access.

> [!NOTE]
>
> Organization policies for Docker Sandboxes require an
> [AI Governance](/manuals/subscription-billing/plans/ai-governance.md)
> subscription. To add licenses,
> [contact sales](https://www.docker.com/pricing/contact-sales/).

To control what sandboxes can access:

1. Assign AI Governance licenses to teams or members. See
   [Manage licenses](/manuals/accounts/organization/manage/manage-licenses.md).
1. Sign in to [Docker Home](https://app.docker.com/) and select your
   organization.
1. In the left navigation, expand **AI Platform** and select
   **Network access**, **Filesystem access**, or **MCP access**.
1. Create a policy and set its scope to the organization or to teams. For
   rules, approval, and precedence, see
   [Organization policies](/manuals/ai/sandboxes/governance/access-controls/organization.md).

To require members to sign in before they use sandboxes, see
[Sign-in enforcement](/manuals/ai/sandboxes/governance/monitor-and-enforce/sign-in-enforcement.md).

## Monitor usage

Find the usage view for a product when you need to check consumption against
your plan, such as seats, Hub pulls, or build minutes.

| Product              | Where to view usage                                                                                                                          |
| -------------------- | -------------------------------------------------------------------------------------------------------------------------------------------- |
| Docker Desktop       | [**Desktop insights**](/manuals/accounts/organization/insights.md), under **Docker Desktop** in Docker Home                                   |
| Docker Hub           | [**Usage**](https://hub.docker.com/usage) for your organization in Docker Hub, including **Pulls** and **Storage**                           |
| Docker Build Cloud   | **Build minutes**, under **Build Cloud** in [Docker Home](https://app.docker.com/)                                                           |
| Docker Scout         | [**Overview**](/manuals/scout/explore/dashboard.md#overview) in Docker Scout                                                                 |
| Testcontainers Cloud | [Testcontainers Cloud](https://app.testcontainers.cloud/)                                                                                    |
| Docker Offload       | [**Offload activity**](/manuals/offload/usage.md), under **Docker Offload** or **Offload** in Docker Home                                    |
| Docker Sandboxes     | [**Audit logs**](/manuals/ai/sandboxes/governance/audit/_index.md), under **AI Platform** in Docker Home                                      |

If usage, seat count, or license assignments exceed what you have,
[add seats](/manuals/accounts/organization/manage/manage-seats.md),
[manage licenses](/manuals/accounts/organization/manage/manage-licenses.md),
or compare [Docker plans](/manuals/subscription-billing/plans/_index.md).

## Next steps

After you decide who can use each product and where to watch usage:

- [Manage licenses](/manuals/accounts/organization/manage/manage-licenses.md)
  to assign Docker Offload, AI Governance, and other add-on licenses to teams
  or members.
- [Manage seats](/manuals/accounts/organization/manage/manage-seats.md) to add
  capacity when invitations or usage need more Docker Team or Business seats.
- [Manage members](/manuals/accounts/organization/manage/members.md) and
  [roles and permissions](/manuals/security/roles-and-permissions/_index.md)
  to control who can change product settings.
- [Desktop insights](/manuals/accounts/organization/insights.md) for deeper
  Docker Desktop usage across the organization.
- [Settings Management](/manuals/desktop/enterprise/hardened-desktop/settings-management/_index.md)
  and
  [Hardened Docker Desktop](/manuals/desktop/enterprise/hardened-desktop/_index.md)
  to enforce Desktop policies beyond the product controls on this page.
- [Docker plans](/manuals/subscription-billing/plans/_index.md) if you need a
  product or add-on your organization does not have yet.
