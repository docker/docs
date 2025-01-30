---
description: Learn how to change your Docker subscription
keywords: Docker Hub, upgrade, downgrade, subscription, Pro, Team, business, pricing plan
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
weight: 30
---

{{% include "tax-compliance.md" %}}

The following sections describe how to change plans when you have a Docker
subscription plan or legacy Docker subscription plan.

> [!NOTE]
>
> Legacy Docker plans apply to Docker subscribers who last purchased or renewed
> their subscription before December 10, 2024. These subscribers will keep
> their current plan and pricing until their next renewal date that falls on or
> after December 10, 2024. To see purchase or renewal history, view your
> [billing history](../billing/history.md). For more details about legacy
> subscriptions, see [Announcing Upgraded Docker
> Plans](https://www.docker.com/blog/november-2024-updated-plans-announcement/).

## Upgrade your subscription

When you upgrade a Docker plan, you immediately have access to all the features and entitlements available in your Docker subscription plan. For detailed information on features available in each subscription, see [Docker Pricing](https://www.docker.com/pricing).

{{< tabs >}}
{{< tab name="Docker plan" >}}

To upgrade your Docker subscription:

1. Sign in to [Docker Home](https://app.docker.com/).
2. Under Settings and administration, select **Billing**.
3. Optional. If you're upgrading from a free Personal plan to a Team plan and want to keep your username, [convert your user account into an organization](../admin/organization/convert-account.md).
4. Select the account you want to upgrade in the drop-down at the top-left of the page.
5. Select **Upgrade**.
6. Follow the on-screen instructions to complete your upgrade.

{{< /tab >}}
{{< tab name="Legacy Docker plan" >}}

You can upgrade a legacy Docker Core, Docker Build Cloud, or Docker Scout subscription plan to a Docker subscription plan that includes access to all tools.

Contact [Docker sales](https://www.docker.com/pricing/contact-sales/) to upgrade your legacy Docker plan.

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
> - SSO and SCIM: If you want to downgrade a Docker Business subscription and your organization uses single sign-on (SSO) for user authentication, you need to remove your SSO connection and verified domains before downgrading. After removing the SSO connection, any organization members that were auto-provisioned (for example, with SCIM) need to set up a password to sign in without SSO. To do this, users can [reset their password at sign in](/accounts/create-account/#reset-your-password-at-sign-in).

{{< tabs >}}
{{< tab name="Docker plan" >}}

If you have a [sales-assisted Docker Business subscription](details.md#sales-assisted), contact your account manager to downgrade your subscription.

To downgrade your Docker subscription:

1. Sign in to [Docker Home](https://app.docker.com/).
2. Under Settings and administration, select **Billing**.
3. Select the account you want to downgrade in the drop-down at the top-left of the page.
4. Select the action icon and then **Cancel subscription**.
5. Review the cancellation warnings, then select **Confirm cancellation**.
6. Optional. Fill out the feedback survey, or select **Skip**.

{{< /tab >}}
{{< tab name="Legacy Docker plan" >}}

If you have a [sales-assisted Docker Business subscription](details.md#sales-assisted), contact your account manager to downgrade your subscription.

### Downgrade Legacy Docker plan

To downgrade your legacy Docker subscription:

1. Sign in to [Docker Hub Billing](https://hub.docker.com/billing).
2. Select the account you want to downgrade in the drop-down at the top-left of the page.
3. Select the link to **Manage this account on Docker Hub**.
4. In the plan section, select **Change plan**.
5. Follow the on-screen instructions to complete your downgrade.

### Downgrade Docker Build Cloud subscription

To downgrade your Docker Build Cloud subscription:

1. Sign in to [Docker Home](https://app.docker.com) and open **Docker Build Cloud**.
2. Select **Account settings**, then **Downgrade**.
3. To confirm your downgrade, type **DOWNGRADE** in the text field and select **Yes, continue**.
4. The account settings page will update with a notification bar notifying you of your downgrade date (start of next billing cycle).

{{< /tab >}}
{{< /tabs >}}

## Pause a subscription

You can't pause or delay a subscription. If a subscription invoice hasn't been paid on the due date, there's a 15 day grace period, including the due date.