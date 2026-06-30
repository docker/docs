---
title: DHI plans
linkTitle: Docker Hardened Images (DHI)
description:
  Manage Docker Hardened Images Select and Enterprise repositories for
  organization accounts, including purchasing, adding repositories, and
  deactivating
keywords: dhi select, dhi enterprise, docker hardened images, hardened images,
  repositories, organization subscription, secure images
weight: 30
aliases:
  - /subscription/products/dhi-select/
  - /subscription/dhi-select/
  - /subscription/plans/dhi-select/
---

[Docker Hardened Images (DHI)](/manuals/dhi/_index.md) are secure, minimal, production-ready container images maintained by Docker. DHI is available in three plans: Community, Select, and Enterprise.

> [!TIP]
> You can self-serve a DHI Select plan. See [Manage plans](../manage.md) to add DHI Select to your organization account.

## Usage

DHI Community gives you access to hardened base images from a public registry at no cost or additional setup. Any organization can pull hardened base images directly from `dhi.io`.

When you upgrade from DHI Community to DHI Select, you purchase a set number of repositories that are mirrored into your organization's namespace. Entitlements are scoped to the organization account that you assign them to during checkout. All organization members can then pull from those mirrored repositories.

DHI Enterprise extends DHI Select with unlimited customizations, full catalog access, the Hardened System Packages repository, and an Extended Lifecycle Support add-on. You must <a href="https://www.docker.com/pricing/contact-sales/" id="dkr_docs_cs_plans_dhi" class="link" rel="noopener">contact sales</a> to upgrade to DHI Enterprise.

For details on setting up and managing repositories, see [Get started with DHI Select and Enterprise](/manuals/dhi/how-to/select-enterprise.md).

## Billing behaviors

DHI Select is an annual plan billed per repository from the date your plan starts. Repositories added mid-cycle are prorated for the remainder of the billing period. You can add more repositories to your DHI Select plan by going to **Active plans** in the billing portal. For steps, see [Manage plans](../manage.md#upgrade-plans).

## Remove repositories

If you want to revert your plan to DHI Community, you must remove repositories from your DHI plan. Repository removals are deferred to the end of the current billing cycle. You can remove repositories at any time, but you cannot stop a plan mid-cycle to receive a partial refund. Repository access remains active until the cycle ends.

To remove repositories:

1. Sign in to [Docker Home](https://app.docker.com/) and go to **Billing**.
1. From **Active plans**, select **Manage** next to **Hardened Images**.
    - Select **Remove repositories** to reduce the repository count to remove repositories.
    - To keep your current repository count after renewal, select **Cancel scheduled change**. 
    - Cancellations and repository removals take effect at the end of the current annual billing cycle.

If you're subscribed to DHI Enterprise, reach out to your sales representative to change your DHI plan.
