---
title: Subscription and billing
linkTitle: Subscription and billing
description: Manage Docker subscriptions, plans, billing, and payments.
keywords: subscription, billing, docker plans, payments, invoices, pricing
weight: 20
params:
  sidebar:
    group: Accounts and admin
aliases:
  - /subscription/
  - /billing/
  - /docker-hub/billing/
  - /docker-hub/billing/faq/
  - /billing/docker-hub-pricing/
grid_subscriptions:
  - title: Compare Docker plans
    description: Visit the pricing page to see what's included in different Docker plans.
    link: "https://www.docker.com/pricing?ref=Docs&refAction=DocsSubscription"
    icon: magnifying-glass
  - title: Manage plans
    description: Add a new plan, upgrade an active plan, or cancel auto-renewal.
    link: /subscription-billing/manage/plans/
    icon: shopping-cart
  - title: Explore plans
    description: Browse available Docker plans and add-ons for individuals, teams, and organizations.
    link: /subscription-billing/plans/
    icon: chart-bar
  - title: Docker Desktop license agreement
    description: Review the terms of the Docker Subscription Service Agreement.
    link: /subscription-billing/desktop-license/
    icon: document-text
grid_core:
  - title: Add or update a payment method
    description: Learn how to add or update a payment method for your personal account or organization.
    link: /subscription-billing/manage/payment-method/
    icon: credit-card
  - title: Update billing information
    description: Learn how to update billing information for your personal account or organization.
    link: /subscription-billing/manage/details/
    icon: pencil-square
  - title: View billing history
    description: Learn how to view billing history and download past invoices.
    link: /subscription-billing/manage/history/
    icon: credit-card
  - title: 3D Secure authentication
    description: Learn how 3DS works and how to troubleshoot verification issues.
    link: /subscription-billing/manage/3d-secure/
    icon: wallet
  - title: Taxes
    description: Learn how to submit a US tax exemption certificate or add a VAT number.
    link: /subscription-billing/manage/tax-certificate/
    icon: document-text
---

You can subscribe to several Docker plans that range from free to paid plans. When you upgrade a plan, you expand your usage entitlements and feature sets for Docker products. You can also top up some plans, extending usage to more users without changing your plan type.

## Docker plans

You can subscribe to plans for individual or organization accounts, or plans for specific products. The following table summarizes the available plans.

| Plans                                                                  | Billing model                                             | Types                                                     |
| ---------------------------------------------------------------------- | --------------------------------------------------------- | --------------------------------------------------------- |
| [Docker](/manuals/subscription-billing/plans/docker.md)                        | Flat-rate plans for personal and organization accounts    | Docker Personal, Docker Pro, Docker Team, Docker Business |
| [Docker Agentic Platform](/manuals/subscription-billing/plans/docker-agentic-platform.md) | Pay-as-you-go (PayGo) for cloud sandbox usage             | Docker Agentic Platform                          |
| [Docker Hardened Images (DHI)](/manuals/subscription-billing/plans/dhi.md)     | Graduated security features for hardened container images | DHI Community, DHI Select, DHI Enterprise                 |
| [Gordon](/manuals/subscription-billing/plans/gordon.md)                        | Prepaid usage for the Gordon AI agent                     | Gordon Plus, Gordon Max, Gordon Ultra                     |
| [AI Governance](/manuals/subscription-billing/plans/ai-governance.md)          | Purchase set amount of licenses                           | AI Governance                                             |
| [Docker Verified Publisher (DVP)](/manuals/subscription-billing/plans/docker-verified-publisher.md)  | Annual plans based on consuming domains                   | DVP Starter, DVP Growth                                   |

Docker plans that upgrade your account (Docker Pro or Docker Team and Business) can provide a foundation for most use cases. Some product plans may require an upgraded Docker account while other product plans let you subscribe without an upgraded account. To learn more, see [Docker plans](/manuals/subscription-billing/plans/_index.md).

## Top up your plan

Plans come with usage entitlements that can be extended without upgrading to a different plan.

| Unit              | Description                                                                           | Examples                      |
| ----------------- | ------------------------------------------------------------------------------------- | ----------------------------- |
| Seats             | Each seat extends entitlements to one more member.                                    | Docker Team, Docker Business  |
| Licenses          | Access to specific products or features.                                              | AI Governance, Docker Offload |
| Minutes           | Cloud build capacity, sold in blocks and consumed within the billing period.          | Docker Build Cloud            |
| Repositories      | Additional container repositories covered by security scanning and analysis features. | DHI                           |
| Consuming domains | Additional consuming domains tracked in publisher analytics, sold in blocks of 25.    | DVP Starter, DVP Growth       |

## Manage your plans

To subscribe to a new plan or upgrade an active plan, see [Manage plans](/manuals/subscription-billing/manage/plans.md). See [Docker plans](/manuals/subscription-billing/plans/docker.md) to learn about Docker Team, Business, and Pro. You can also <a href="https://www.docker.com/pricing/contact-sales/" id="dkr_docs_index_sales" class="link" rel="noopener">contact sales</a>.

You can use the billing portal to manage your Docker subscriptions, such
as updating payment methods, reviewing billing details, and tracking
invoice history.

## Billing

You can manage your Docker plans from the billing portal:

1. Sign in to [Docker Home](https://app.docker.com/), then choose your
   account.
1. Go to **Billing** to view the **Overview** page.
1. Select the page you want to explore.

### Usage

The billing **Usage** page helps you compare usage-based charges across
billing periods. You can track usage by changing the period, product,
and how the product is metered.

### Costs

The billing **Costs** page aggregates all costs by billing period. It
breaks down charges by resource (the product accruing a charge), the
status of your billing period, and costs to date.

### Credits

The billing **Credits** page shows credits applied to your costs. If you
received a promotional credit, you can see how it applies to your bill
from this page.

## Docker plans and billing cycle

Your invoice history is a reference to the Docker plans you subscribe
to. For information about your billing cycle and renewal dates, see
[Billing cycle](/manuals/subscription-billing/manage/details.md#billing-cycle). To upgrade or add
a new plan, see [Subscription](/manuals/subscription-billing/_index.md).

## Next steps

{{< grid items="grid_subscriptions" >}}

{{< grid items="grid_core" >}}
