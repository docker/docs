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
  - /subscription/products/build-cloud/
  - /subscription/build-cloud/
  - /subscription/core-subscription/downgrade/
---

Docker plans refer to plans that upgrade your account type from the basic free plan to a paid plan. Paid Docker plans come with higher usage limits, commercial
licensing, and expanded feature sets.

- Docker Personal is free for individual developers. Docker Pro adds unlimited private repositories, Docker Build Cloud, and commercial Docker Desktop use.
- Docker Team and Docker Business are plans for organizations, with Team adding audit logs and role-based access control, and Business adding SSO, SCIM, hardened Docker Desktop, and image access management.

To upgrade your free Docker plan in the billing portal, see [Manage plans](../manage.md).

## Usage

Docker Personal and Docker Pro are Docker plans for individual account types while Docker Team and Docker Business are Docker plans for organization account types. For a full feature and pricing breakdown, see the
<a href="https://www.docker.com/pricing/" id="dkr_docs_index_pricing_docker_plans" class="link" rel="noopener">Docker pricing page</a>. 

> [!TIP]
> If you're upgrading from a Personal plan to a Team plan
> and want to keep your username,
> [convert your user account into an organization](/manuals/admin/organization/setup/convert-account.md).

## Billing behaviors

Docker individual and organization plans are billed at a flat rate per user per month, with monthly or
annual billing options.
Upgrading your plan immediately extends access to all features
and entitlements.

### Organization seats

For Docker Team and Docker Business, you can purchase more seats for new members to extend access to your paid Docker plan. To add or remove seats from your Docker plan:

1. Sign in to [Docker Home](https://app.docker.com/), then choose your organization account.
1. Go to **Billing** to view the Overview page, then go to **Active plans**.
1. From the Docker Team or Docker Business tile, select the action menu.
   - Select **Add seats** or **Remove seats** from the drop-down menu.
   - When you add or remove seats, review your current seats against your new total seats.
   - When you remove seats, you must remove members from your organization.
1. Verify your billing details, continue to payment, and complete checkout.

To learn how to manage seats from the Admin Console, see
[Manage seats](/manuals/admin/organization/manage/manage-seats.md).

### Docker Offload licenses

[Docker Offload](/manuals/offload/_index.md) licenses are available for Docker Team and Docker Business plans. Once assigned to your account, organization owners can [manage license assignments](/manuals/admin/organization/manage/manage-licenses.md) in the Admin Console.

To add Docker Offload licenses, you must <a href="https://www.docker.com/pricing/contact-sales/" id="dkr_docs_cs_plans_docker_offload" class="link" rel="noopener">contact sales</a>.

### Docker Build Cloud minutes

Each plan includes a base allocation of [Docker Build Cloud](/manuals/build-cloud/_index.md) build minutes per
month. Base minutes reset on an annual or monthly cadence, and don't accumulate. Additional purchased
minutes expire at the end of your billing period.

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

You can add Testcontainers Cloud runtime minutes in two ways:

- <a href="https://www.docker.com/pricing/contact-sales/" id="dkr_docs_cs_plans_docker_testcontainers" class="link" rel="noopener">Contact sales</a> to
  pre-purchase runtime minutes at $3 per 100 minutes. Pre-purchased minutes
  expire at the end of your billing period.
- Use on-demand runtime minutes at $4 per 100 minutes, billed at the end of
  each monthly cycle.

## Cancel a Docker plan

> [!NOTE]
> If you have a sales-assisted Docker Business plan,
> you must contact your account manager to cancel.

You can cancel at any time before your renewal date, but you can't pause or delay a plan. If an invoice isn't paid by the due date, there's a 15-day grace period starting from the due date. While the unused portion is not refundable, you still retain access to paid features until the end of the current billing cycle. Canceling your paid plans may have implications for collaborators or organization members:

- Docker Pro private repository collaborators are removed and additional private repositories are locked.
- Docker Team or Docker Business members provisioned through SCIM without a password will be locked out. Remove SSO connections and verified domains if your organization uses single sign-on.
- For paid individual and organization plans, you must convert private repositories to fit your new plan limits.

Canceling a paid plan returns your account to Docker Personal or a basic organization account.
To cancel your plan:

1. Sign in to [Docker Home](https://app.docker.com/) and go to **Billing**.
2. From **Active plans**, select the action menu next to your Docker plan.
3. Select **Cancel plan** and complete the feedback survey.
