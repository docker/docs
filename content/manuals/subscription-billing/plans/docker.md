---
title: Docker plans
linkTitle: Docker
description:
  Learn about Docker plans that upgrade individual and organization accounts,
  including usage entitlements, billing cycle, and downgrade options.
keywords:
  docker team, docker business, docker pro, docker personal,
  subscription seats, upgrade subscription, downgrade subscription, docker
  pricing, subscription changes, build cloud minutes, testcontainers minutes,
  docker offload licenses, cancel subscription
weight: 10
aliases:
  - /subscription/plans/docker/
  - /subscription/plans/core/
  - /subscription/products/core/
  - /subscription/core/
  - /subscription/testcontainers-cloud/
  - /subscription/products/testcontainers-cloud/
  - /subscription/products/build-cloud/
  - /subscription/build-cloud/
  - /subscription/core-subscription/downgrade/
---

Docker plans upgrade your account from the free Docker Personal plan to a
paid plan with higher usage limits, commercial licensing, and expanded
features.

- Docker Personal is free for individual developers. Docker Pro adds
  unlimited private repositories, Docker Build Cloud, and commercial Docker
  Desktop use.
- Docker Team and Docker Business are plans for organizations. Team adds
  audit logs and role-based access control. Business adds single sign-on
  (SSO), SCIM, Hardened Docker Desktop, and Image Access Management.

To upgrade to Docker Pro, go to the
<a href="https://www.docker.com/pricing/" id="dkr_docs_pricing_docker_pro" class="link" rel="noopener">Docker pricing page</a>
and select **Buy now**. To upgrade to Docker Team or Docker Business, see
[Manage plans](/manuals/subscription-billing/manage/plans.md).

## Usage

Docker Personal and Docker Pro apply to individual accounts. Docker Team and
Docker Business apply to organization accounts. For a full feature and
pricing comparison, see the
<a href="https://www.docker.com/pricing/" id="dkr_docs_index_pricing_docker_plans" class="link" rel="noopener">Docker pricing page</a>.

> [!TIP]
> If you're upgrading from a Personal plan to a Team plan and want to keep
> your username,
> [convert your user account into an organization](/manuals/accounts/organization/setup/convert-account.md).

## Billing cycle

Docker plans are billed at a flat rate per user per month, with monthly or
annual billing options. Upgrading your plan immediately extends access to
all features and entitlements.

### Switch to annual billing

When you manage the subscription yourself and the plan is billed monthly,
the action menu includes **Switch to annual billing**. The review shows the
annual total and a **Prorated credit** for the unused part of the current
monthly period. The amount due today is labeled **Prorated for this year**.

1. Sign in to [Docker Home](https://app.docker.com/) and go to **Billing**.
1. From **Active subscriptions**, select the action menu next to your Docker
   plan.
1. Select **Switch to annual billing**.
1. Review the annual total and prorated credit, then complete checkout.

### Organization seats

For Docker Team and Docker Business, each seat extends your paid plan to one
more member. You can remove only unused seats. To remove more, first
[remove members](/manuals/accounts/organization/manage/members.md#remove-member-from-organization)
from your organization.

To add or remove seats:

1. Sign in to [Docker Home](https://app.docker.com/), then choose your
   organization account.
1. Go to **Billing** to view the Overview page, then go to
   **Active subscriptions**.
1. From the Docker Team or Docker Business tile, select the action menu, then
   select **Add seats** or **Remove seats**.
1. Choose the number of seats and review your new total.
1. For added seats, verify your billing details, continue to payment, and
   complete checkout.

To learn how to manage seats from Docker Home, see
[Manage seats](/manuals/accounts/organization/manage/manage-seats.md).

### Docker Offload licenses

[Docker Offload](/manuals/offload/_index.md) licenses are available for
Docker Team and Docker Business plans. Once assigned to your account,
organization owners can
[manage licenses](/manuals/accounts/organization/manage/manage-licenses.md)
in Docker Home.

To add Docker Offload licenses, you must
<a href="https://www.docker.com/pricing/contact-sales/" id="dkr_docs_cs_plans_docker_offload" class="link" rel="noopener">contact sales</a>.

### Docker Build Cloud minutes

Each plan includes a base allocation of
[Docker Build Cloud](/manuals/build-cloud/_index.md) build minutes per
month. Base minutes reset on an annual or monthly cadence and don't
accumulate. Additional purchased minutes expire at the end of your billing
period.

To purchase additional minutes:

1. From [Docker Home](https://app.docker.com/), choose your organization.
1. Select **Build Cloud**, then **Build minutes**.
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
> If you have a sales-assisted Docker Business plan, contact your account
> manager to cancel.

You can cancel at any time before your renewal date, but you can't pause or
delay a plan. The unused portion of your plan isn't refundable, but you keep
access to paid features until the end of the current billing cycle. If an
invoice isn't paid by the due date, a 15-day grace period starts from the
due date.

> [!NOTE]
>
> Billing periods are calculated in UTC. If your time zone is behind UTC,
> for example US Pacific Time at UTC−7 or UTC−8, canceling or making changes
> late in your local day can fall into the next billing day.

Canceling a paid plan can affect collaborators and organization members:

- Docker Pro private repository collaborators are removed, and additional
  private repositories are locked.
- Docker Team or Docker Business members provisioned through SCIM without a
  password are locked out. Remove SSO connections and verified domains if
  your organization uses single sign-on.
- For paid individual and organization plans, you must convert private
  repositories to fit your new plan limits.

Canceling a paid plan returns your account to Docker Personal or a basic
organization account. To cancel your plan:

1. Sign in to [Docker Home](https://app.docker.com/) and go to **Billing**.
1. From **Active subscriptions**, select the action menu next to your Docker
   plan.
1. Select **Cancel plan** and complete the feedback survey.

If you cancel and the plan is still in its current billing cycle, the action
menu includes **Reactivate**.

## Next steps

- To add or upgrade a plan, see
  [Manage plans](/manuals/subscription-billing/manage/plans.md)
- To manage seats for your organization, see
  [Manage seats](/manuals/accounts/organization/manage/manage-seats.md)
- To track usage across plans, see
  [Usage](/manuals/subscription-billing/manage/_index.md#usage)
