---
title: Docker plans
linkTitle: Docker
description:
  Learn about Docker plans that upgrade personal and organization accounts,
  including usage entitlements, billing behaviors, and downgrade options.
keywords:
  docker core, docker team, docker business, docker pro, docker personal,
  subscription seats, upgrade subscription, downgrade subscription, docker
  pricing, subscription changes, build cloud minutes, testcontainers minutes
weight: 10
aliases:
  - /subscription/plans/core/
  - /subscription/products/core/
  - /subscription/core/
  - /subscription/testcontainers-cloud/
  - /subscription/products/testcontainers-cloud/
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
  - /subscription/products/build-cloud/
  - /subscription/build-cloud/
---

Docker plans upgrade your account with higher usage limits, commercial
licensing, and expanded feature sets. Plans are available for personal accounts
and organization accounts. To subscribe to a plan, see
[Setup](../setup.md).

## Usage

There are four Docker plans between individual and organization account types.

| Feature                      | Personal   | Pro        | Team       | Business  |
| ---------------------------- | ---------- | ---------- | ---------- | --------- |
| Docker Desktop               | Basic  | Commercial | Commercial | Hardened  |
| Private Hub repos            | 1          | Unlimited  | Unlimited  | Unlimited |
| Hub pull rate                | 100/hr     | Unlimited  | Unlimited  | Unlimited |
| Docker Scout repos           | 1          | 2          | Unlimited  | Unlimited |
| Gordon                       | Included   | Included   | —          | —         |
| Build Cloud minutes/month    | Free trial | 200        | 500        | 1,500     |
| Testcontainers minutes/month | Free trial | 100        | 500        | 1,500     |
| SSO / SCIM                   | —          | —          | —          | ✓         |
| Registry & Image Access Mgmt | —          | —          | —          | ✓         |
| Max users                    | 1          | 1          | 100        | Unlimited |

## Billing behaviors

Docker individual and organization plans are billed at a flat rate per user per month, with monthly or
annual billing options.
Upgrading your plan immediately extends access to all features
and entitlements.

For Docker Team and Docker Business, you can purchase more seats to extend access to new members in your organization. To learn how to manage seats, see
[Manage seats](/manuals/admin/organization/manage/manage-seats.md).

> [!TIP]
> If you're upgrading from a Personal plan to a Team plan
> and want to keep your username,
> [convert your user account into an organization](/manuals/admin/organization/setup/convert-account.md).

### Docker Build Cloud minutes

Each plan includes a base allocation of Docker Build Cloud build minutes per
month.

| Plan            | Included build minutes/month |
| --------------- | ---------------------------- |
| Docker Personal | Free trial                   |
| Docker Pro      | 200                          |
| Docker Team     | 500                          |
| Docker Business | 1,500                        |

Base minutes reset on an annual or monthly cadence, and don't accumulate. Additional purchased
minutes expire at the end of your billing period.
For example:  

- On an annual plan, additional minutes last until your annual renewal. 
- On a monthly plan, they expire at month end.

For an
overview of Docker Build Cloud features, see the
[Docker Build Cloud overview](/manuals/build-cloud/_index.md).

#### Add Build Cloud minutes 

To purchase additional minutes:

1. From [Docker Home](https://app.docker.com/), choose your organization.
1. Select Build Cloud, then Build minutes.
1. From the **Minute breakdown** table, select **Add minutes**.
1. Choose your additional minute amount.
1. Verify your billing details, continue to payment, and complete checkout.

Your additional minutes appear on the Build minutes page immediately. 

### Testcontainers Cloud minutes

Each plan includes a base allocation of Testcontainers Cloud runtime minutes
per month. Base minutes reset monthly and don't accumulate.

| Plan            | Included runtime minutes/month |
| --------------- | ------------------------------ |
| Docker Personal | Free trial                     |
| Docker Pro      | 100                            |
| Docker Team     | 500                            |
| Docker Business | 1,500                          |

You can add Testcontainers Cloud runtime minutes in two ways:

- [Contact sales](https://www.docker.com/pricing/contact-sales/) to
  pre-purchase runtime minutes at $3 per 100 minutes. Pre-purchased minutes
  expire at the end of your billing period.
- Use on-demand runtime minutes at $4 per 100 minutes, billed at the end of
  each monthly cycle.

## Downgrade

> [!NOTE]
> You can't pause or delay a plan. If an 
> invoice isn't paid by the due date, there's a 
> 15-day grace period starting from the due date.

You can cancel or downgrade at any time before your renewal date. The unused portion is not refundable, but you retain access to paid features until the end of the current billing cycle.

- When you downgrade from Docker Pro, your private repository collaborators are removed and additional private repositories are locked.
- If you have Docker Team or Docker Business:
    - Remove SSO connections and verified domains if your organization uses single sign-on. Members provisioned through SCIM without a password will be locked out.
    - Reduce team members or convert private repositories to fit your new plan limits.
- If you have a sales-assisted Docker Business plan, contact your account manager to downgrade.


