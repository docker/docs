---
description: Learn how to change your Docker subscription
keywords: Docker Hub, upgrade, downgrade, subscription, Pro, Team, business, pricing
title: Change your subscription
aliases:
- /docker-hub/upgrade/
- /docker-hub/billing/upgrade/
- /subscription/upgrade/
- /subscription/downgrade/
- /subscription/core-subscription/upgrade/
- /subscription/core-subscription/downgrade/
- /docker-hub/cancel-downgrade/
- /docker-hub/billing/downgrade/
- /billing/scout-billing/
- /billing/subscription-management/
weight: 30
---

{{% include "tax-compliance.md" %}}

The following sections describe how to change plans when you have a Docker
subscription or legacy Docker subscription.

> [!NOTE]
>
> Legacy Docker plans apply to Docker subscribers who last purchased or renewed
> their subscription before December 10, 2024. These subscribers will keep
> their current subscription and pricing until their next renewal date that falls on or
> after December 10, 2024. To see purchase or renewal history, view your
> [billing history](../billing/history.md). For more details about legacy
> subscriptions, see [Announcing Upgraded Docker
> Plans](https://www.docker.com/blog/november-2024-updated-plans-announcement/).

## Upgrade your subscription

When you upgrade a Docker subscription, you immediately have access to all the features and entitlements available in your Docker subscription. For detailed information on features available in each subscription, see [Docker Pricing](https://www.docker.com/pricing).

{{< tabs >}}
{{< tab name="Docker subscription" >}}

To upgrade your Docker subscription:

1. Sign in to [Docker Home](https://app.docker.com/) and select the organization
you want to upgrade.
1. Select **Billing**.
1. Optional. If you're upgrading from a free Personal subscription to a Team subscription and want to keep your username, [convert your user account into an organization](../admin/organization/convert-account.md).
1. Select **Upgrade**.
1. Follow the on-screen instructions to complete your upgrade.

> [!NOTE]
>
> If you choose to pay using a US bank account, you must verify the account. For
> more information, see [Verify a bank account](manuals/billing/payment-method.md#verify-a-bank-account).

{{< /tab >}}
{{< tab name="Legacy Docker subscription" >}}

You can upgrade a legacy Docker Core, Docker Build Cloud, or Docker Scout subscription to a Docker subscription that includes access to all tools.

Contact [Docker sales](https://www.docker.com/pricing/contact-sales/) to upgrade your legacy Docker subscription.

{{< /tab >}}
{{< /tabs >}}

## Downgrade your subscription

You can downgrade your Docker subscription at anytime before the renewal date. The unused portion of the subscription isn't refundable or creditable.

When you downgrade your subscription, access to paid features is available until the next billing cycle. The downgrade takes effect on the next billing cycle.

> [!IMPORTANT]
>
> If you downgrade your personal account from a Pro subscription to a Personal subscription, note that [Personal subscriptions](details.md#docker-personal) don't include collaborators for private repositories. Only one private repository is included with a Personal subscription. When you downgrade, all collaborators will be removed and additional private repositories are locked.
> Before you downgrade, consider the following:
> - Team size: You may need to reduce the number of team members and convert any private repositories to public repositories or delete them. For information on features available in each tier, see [Docker Pricing](https://www.docker.com/pricing).
> - SSO and SCIM: If you want to downgrade a Docker Business subscription and your organization uses single sign-on (SSO) for user authentication, you need to remove your SSO connection and verified domains before downgrading. After removing the SSO connection, any organization members that were auto-provisioned (for example, with SCIM) need to set up a password to sign in without SSO. To do this, users can [reset their password at sign in](/accounts/create-account/#reset-your-password).

{{< tabs >}}
{{< tab name="Docker subscription" >}}

If you have a [sales-assisted Docker Business subscription](details.md#sales-assisted), contact your account manager to downgrade your subscription.

To downgrade your Docker subscription:

1. Sign in to [Docker Home](https://app.docker.com/) and select
the organization you want to downgrade.
1. Select **Billing**.
1. Select the action icon and then **Cancel subscription**.
1. Fill out the feedback survey to continue with cancellation.

{{< /tab >}}
{{< tab name="Legacy Docker subscription" >}}

If you have a [sales-assisted Docker Business subscription](details.md#sales-assisted), contact your account manager to downgrade your subscription.

### Downgrade Legacy Docker subscription

To downgrade your legacy Docker subscription:

1. Sign in to [Docker Hub](https://hub.docker.com/billing).
1. Select the organization you want to downgrade, then select **Billing**.
1. To downgrade, you must navigate to the upgrade plan page. Select **Upgrade**.
1. On the upgrade page, select **Downgrade** in the **Free Team** plan card.
1. Follow the on-screen instructions to complete your downgrade.

### Downgrade Docker Build Cloud subscription

To downgrade your Docker Build Cloud subscription:

1. Sign in to [Docker Home](https://app.docker.com) and select **Build Cloud**.
1. Select **Account settings**, then **Downgrade**.
1. To confirm your downgrade, type **DOWNGRADE** in the text field and select **Yes, continue**.
1. The account settings page will update with a notification bar notifying you of your downgrade date (start of next billing cycle).

{{< /tab >}}
{{< /tabs >}}

## Pause a subscription

You can't pause or delay a subscription. If a subscription invoice hasn't been paid on the due date, there's a 15 day grace period, including the due date.