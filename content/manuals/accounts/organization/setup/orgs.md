---
title: Create a Docker organization
linkTitle: Create
weight: 10
description: Create a Docker organization and choose its namespace and plan.
keywords:
  - create Docker organization
  - organization namespace
  - organization name
  - Docker Team
  - Docker Business
  - Docker Home
aliases:
  - /admin/organization/setup/orgs/
  - /docker-hub/orgs/
  - /admin/organization/orgs/
---

{{< summary-bar feature_name="Admin orgs" >}}

Create an organization to group members and teams under one namespace and
one subscription. Your Docker ID stays an individual account. To use an
existing Docker ID as the namespace, see
[Convert a Docker account to an organization](/manuals/accounts/organization/setup/convert-account.md).

## Prerequisites

You need a [Docker ID](/manuals/accounts/_index.md) before you create an
organization.

> [!TIP]
>
> Review [Docker subscriptions and features](https://www.docker.com/pricing?ref=Docs&refAction=DocsAdminOrgs)
> before you choose a plan.

## Create an organization

When you create a new organization, you must select a Docker plan,
enter organization details, and verify billing details.

1. Sign in to [Docker Home](https://app.docker.com/) and select
   **Create new organization** at the bottom of the organization list.
1. On **Plan**, choose a subscription, a billing cycle, and the number of
   seats. Select **Continue to profile**.
1. On **Organization**, enter the details for the new organization.
   - If you already belong to one or more organizations, this step opens
     as **Choose an organization**, which applies the subscription to an
     existing organization.
   - Select **Create an organization** to make a new one instead. The
     picker is replaced by the **Organization namespace** and
     **Organization name** fields.
   - For what each field means, see
     [Names versus namespaces](/manuals/accounts/organization/setup/_index.md#names-versus-namespaces).
1. Select **Continue to billing**.
1. On **Billing**, enter billing information and select
   **Continue to payment**.
1. On **Payment**, enter payment details and select **Purchase**.

You can now view your new organization.

## View an organization

1. Sign in to [Docker Home](https://app.docker.com).
1. Select your organization.

Docker Home lists the options you use to configure the organization.

## Merge organizations

> [!WARNING]
>
> Merge organizations at the end of your billing cycle. When you merge an
> organization and downgrade another, you lose seats on the downgraded
> organization. Docker doesn't offer refunds for downgrades.

If you have multiple organizations that you want to merge into one:

1. Based on the number of seats from the secondary organization,
   [purchase additional seats](../manage/manage-seats.md) for the primary
   organization you want to keep.
1. Add users to the primary organization and remove them from the
   secondary organization.
1. Move your data, including repositories.
1. After the users and data are on the primary organization,
   [downgrade](../../../subscription-billing/plans/docker.md#cancel-a-docker-plan)
   the secondary account to a free subscription. Docker doesn't offer
   refunds for a downgrade in the middle of a billing cycle.

If your organization has a Docker Business subscription with a purchase
order, contact Support or your account manager at Docker.

## Next steps

- [Onboard your organization](/manuals/accounts/organization/setup/onboard.md)
- [Manage organization members](/manuals/accounts/organization/manage/members.md)
