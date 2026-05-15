---
title: Change your subscription
description: Upgrade or downgrade your Docker subscription and understand billing changes
keywords: upgrade subscription, downgrade subscription, docker pricing, subscription changes
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

You can upgrade or downgrade your Docker subscription at any time to match your changing needs. This page explains how to make subscription changes and what to expect with billing and feature access.

## Upgrade your subscription

When you upgrade your Docker subscription, you immediately get access to all features and entitlements in your new subscription tier.

For detailed feature information, see [Docker Pricing](https://www.docker.com/pricing?ref=Docs&refAction=DocsSubscriptionChange).

To upgrade your subscription:

1. Sign in to [Docker Home](https://app.docker.com/) and select the organization
   you want to upgrade.
1. Select **Billing**.
1. Optional. If you're upgrading from a free Personal subscription to a Team subscription and want to keep your username, [convert your user account into an organization](../admin/organization/setup/convert-account.md).
1. Select **Upgrade**.
1. Follow the on-screen instructions to complete your upgrade. If you choose to pay using a US bank account, you must verify the account. For more information, see [Verify a bank account](/manuals/billing/payment-method.md#verify-a-bank-account).

## Downgrade your subscription

You can downgrade your Docker subscription at any time before the renewal date. The unused portion isn't refundable, but you retain access to paid features until the next billing cycle.

### Downgrade considerations

Consider the following before downgrading:

- Team size and repositories: You may need to reduce team members and convert private repositories to public or delete them based on your new subscription limits.
- SSO and SCIM: If downgrading from Docker Business and your organization uses single sign-on, remove your SSO connection and verified domains first. Organization members who were auto-provisioned through SCIM need to reset their passwords to sign in without SSO.
- Private repository collaborators: Personal subscriptions don't include collaborators for private repositories. When downgrading from Pro to Personal, all collaborators are removed and additional private repositories are locked.

For feature limits in each tier, see [Docker Pricing](https://www.docker.com/pricing?ref=Docs&refAction=DocsSubscriptionChange).

To downgrade your subscription:

1. Sign in to [Docker Home](https://app.docker.com/) and select
   the organization you want to downgrade.
1. Select **Billing**.
1. Select the action icon and then **Cancel subscription**.
1. Fill out the feedback survey to continue with cancellation.

> [!IMPORTANT]
>
> If you have a sales-assisted Docker Business subscription, contact your account manager to downgrade your subscription.

## Subscription pause policy

You can't pause or delay a subscription. If a subscription invoice isn't paid by the due date, there's a 15-day grace period starting from the due date.
