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
weight: 30
---


{{< include "tax-compliance.md" >}}

## Upgrade your subscription

When you upgrade to a paid subscription, you immediately have access to all the features and entitlements available in your new chosen subscription. For detailed information on features available in each subscription, see [Docker Pricing](https://www.docker.com/pricing).

1. Sign in to [Docker Home](https://app.docker.com).

2. Optional. If you're upgrading from a free user account to a Team subscription and want to keep your account name, [convert your user account into an organization](../admin/convert-account.md).

3. Select your **avatar** to expand the drop-down menu.

4. From the drop-down menu, select **Billing** and then the account you want to upgrade.

5. On the **Billing Details** tab, select **Change plan** and then choose the plan you'd like to upgrade to.

   > [!TIP]
   >
   > If your current plan is a free plan, select **Buy now**.

6. Follow the on-screen instructions.

   If you have a coupon to use, you can enter it during this step.

## Downgrade your subscription

You can downgrade your Docker subscription at anytime before the renewal date. The unused portion of the subscription isn't refundable or creditable.

When you downgrade your subscription, access to paid features is available until the next billing cycle. The downgrade takes effect on the next billing cycle.

> [!IMPORTANT]
>
> If you downgrade your personal account from a Pro subscription to a Personal subscription, note that [Personal subscriptions](details.md#docker-personal) don't include collaborators for private repositories, and only one private repository is included. When you downgrade, all collaborators will be removed and additional private repositories are locked.

Before you downgrade, consider the following aspects before you downgrade your subscription:

- Team size: You may need to reduce the number of team members and convert any private repositories to public repositories or delete them. For information on features available in each tier, see [Docker Pricing](https://www.docker.com/pricing).

- SSO and SCIM: If you want to downgrade a Docker Business subscription and your organization uses single sign-on (SSO) for user authentication, you need to remove your SSO connection and verified domains before downgrading. After removing the SSO connection, any organization members that were auto-provisioned (for example, with SCIM) need to set up a password to sign in without SSO. To do this, users can [reset their password at sign in](/accounts/create-account/#reset-your-password-at-sign-in).

>[!IMPORTANT]
>
> If you have a [sales-assisted Docker Business subscription](details.md#sales-assisted), contact your account manager to downgrade your subscription.

To downgrade your subscription:

1. Sign in to [Docker Home](https://app.docker.com).

2. Select your **avatar** and from the drop-down menu select **Billing**.

3. Choose either your personal account or an organization to downgrade.

4. In the plan section, select **Change plan**.

5. Select the plan you'd like to downgrade to.

6. Review the downgrade warning and select **Continue**.

7. Optional. Select a reason for your downgrade from the list and select **Send**.
    The **Billing** page displays a confirmation of the downgrade with details on when the downgrade changes take effect.

If you want to cancel the downgrade, select **Cancel the downgrade** on the **Plan** tab.

## Pause a subscription

You can't pause or delay a subscription. If a subscription invoice hasn't been paid on the due date, there's a 15 day grace period, including the due date.
