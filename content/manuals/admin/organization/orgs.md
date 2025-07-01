---
title: Create your organization
weight: 10
description: Learn how to create an organization.
keywords: docker organizations, organization, create organization, docker teams, docker admin console, organization management
aliases:
  - /docker-hub/orgs/
---

{{< summary-bar feature_name="Admin orgs" >}}

This page describes how to create an organization.

## Prerequisites

Before you begin creating an organization:

- You need a [Docker ID](/accounts/create-account/)
- Review the [Docker subscriptions and features](../../subscription/details.md)
  to determine what subscription to choose for your organization

## Create an organization

There are multiple ways to create an organization. You can either:

- Create a new organization using the **Create Organization** option in the
Admin Console or Docker Hub
- Convert an existing user account to an organization

The following section contains instructions on how to create a new organization. For prerequisites and
detailed instructions on converting an existing user account to an organization, see
[Convert an account into an organization](/manuals/admin/organization/convert-account.md).

{{< tabs >}}
{{< tab name="Admin Console" >}}

To create an organization:

1. Sign in to [Docker Home](https://app.docker.com/).
1. Under **Settings and administration**, select **Go to Admin Console**.
1. Select the organization drop-down, then **Create Organization**.
1. Choose a subscription for your organization, a billing cycle, and specify
how many seats you need. See [Docker Pricing](https://www.docker.com/pricing/)
for details on the features offered in Docker's subscriptions.
1. Select **Continue to profile**.
1. Enter an **Organization namespace**. This is the official, unique name for
   your organization in Docker Hub. It's not possible to change the name of the
   organization after you've created it.

   > [!NOTE]
   >
   > You can't use the same name for the organization and your Docker ID. If
   you want to use your Docker ID as the organization name, then you must
   first [convert your account into an organization](/manuals/admin/organization/convert-account.md).

1. Enter your **Company name**. This is the full name of your company. Docker
   displays the company name on your organization page and in the details of any
   public images you publish. You can update the company name anytime by
   navigating to your organization's **Settings** page.
1. Select **Continue to billing** to continue.
1. Enter your organization's billing information and select
**Continue to payment** to continue to the billing portal.
1. Provide your payment details and select **Purchase**.

{{< /tab >}}
{{< tab name="Docker Hub" >}}

{{% include "hub-org-management.md" %}}

1. Sign in to [Docker Hub](https://hub.docker.com/) using your Docker ID,
your email address, or your social provider.
1. Select the organization drop-down, then **Create Organization**.
1. Choose a subscription for your organization, a billing cycle, and specify
how many seats you need. See [Docker Pricing](https://www.docker.com/pricing/)
for details on the features offered in Docker's subscriptions.
1. Select **Continue to profile**.
1. Enter an **Organization namespace**. This is the official, unique name for
   your organization in Docker Hub. It's not possible to change the name of the
   organization after you've created it.

   > [!NOTE]
   >
   > You can't use the same name for the organization and your Docker ID. If
   you want to use your Docker ID as the organization name, then you must
   first [convert your account into an organization](/manuals/admin/organization/convert-account.md).

1. Enter your **Company name**. This is the full name of your company. Docker
   displays the company name on your organization page and in the details of any
   public images you publish. You can update the company name anytime by
   navigating to your organization's **Settings** page.
1. Select **Continue to billing** to continue.
1. Enter your organization's billing information and select
**Continue to payment** to continue to the billing portal.
1. Provide your payment details and select **Purchase**.

{{< /tab >}}
{{< /tabs >}}

## View an organization

{{< tabs >}}
{{< tab name="Admin Console" >}}

To view an organization in the Admin Console:

1. Sign in to [Docker Home](https://app.docker.com).
1. Under Settings and administration, select **Go to Admin Console**.
1. Select your organization from the organization drop-down in the
left-hand navigation.

The Admin Console displays various options that let you to
configure your organization.

- **Members**: Displays a list of team members. You
  can invite new members using the **Invite members** button. See [Manage members](./members.md) for details.
- **Teams**: Displays a list of existing teams and the number of
  members in each team. See [Create a team](./manage-a-team.md) for details.
- **Activity** Displays the audit logs, a chronological list of activities that
  occur at organization and repository levels. It provides the org owners a
  report of all their team member activities. See [Audit logs](./activity-logs.md) for
  details.
- **Security and access**: Manage security settings. For more information, see [Security](/manuals/security/_index.md).
- **Organization settings**: Update organization information or [deactivate your organization](/manuals/admin/organization/deactivate-account.md).

{{< /tab >}}
{{< tab name="Docker Hub" >}}

{{% include "hub-org-management.md" %}}

To view an organization:

1.  Sign in to [Docker Hub](https://hub.docker.com) with a user account that is
    a member of any team in the organization.

          > [!NOTE]
          >
          > You can't sign in directly to an organization. After converting from
          > a user to an organization, sign in with the owner or team member
          > account. If you don’t see the org, you may not be a member or owner.
1.  Select **My Hub** in the top navigation bar, then choose your
    organization from the list.

The organization landing page displays various options that let you to
configure your organization.

- **Members**: Displays a list of team members. You
  can invite new members using the **Invite members** button. See [Manage members](./members.md) for details.
- **Teams**: Displays a list of existing teams and the number of
  members in each team. See [Create a team](./manage-a-team.md) for details.
- **Repositories**: Displays a list of repositories associated with the
  organization. See [Repositories](../../docker-hub/repos/_index.md) for detailed information about
  working with repositories.
- **Activity** Displays the audit logs, a chronological list of activities that
  occur at organization and repository levels. It provides the org owners a
  report of all their team member activities. See [Audit logs](./activity-logs.md) for
  details.
- **Settings**: Displays information about your
  organization, and you to view and change your repository privacy
  settings, configure org permissions such as
  [Image Access Management](/manuals/security/for-admins/hardened-desktop/image-access-management.md), configure notification settings, and [deactivate](/manuals/admin/organization/deactivate-account.md#deactivate-an-organization) You can also update your organization name and company name that appear on your organization landing page. You must be an owner to access the
  organization's **Settings** page.
- **Billing**: Displays information about your existing
  [Docker subscription](../../subscription/_index.md), including the number of seats and next payment due date. For how to access the billing history and payment methods for your organization, see [View billing history](../../billing/history.md).

{{< /tab >}}
{{< /tabs >}}

## Merge organizations

> [!WARNING]
>
> If you are merging organizations, it is recommended to do so at the _end_ of
> your billing cycle. When you merge an organization and downgrade another, you
> will lose seats on your downgraded organization. Docker does not offer
> refunds for downgrades.

If you have multiple organizations that you want to merge into one, complete
the following steps:

1. Based on the number of seats from the secondary organization,
[purchase additional seats](../../subscription/manage-seats.md) for the
primary organization account that you want to keep.
1. Manually add users to the primary organization and remove existing users
from the secondary organization.
1. Manually move over your data, including all repositories.
1. Once you're done moving all of your users and data,
[downgrade](../../subscription/change.md) the secondary account to a free
subscription. Note that Docker does not offer refunds for downgrading
organizations mid-billing cycle.

> [!TIP]
>
> If your organization has a Docker Business subscription with a purchase
order, contact Support or your Account Manager at Docker.

## More resources

- [Video: Docker Hub Organizations](https://www.youtube.com/watch?v=WKlT1O-4Du8)
