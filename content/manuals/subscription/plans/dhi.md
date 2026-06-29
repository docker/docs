---
title: DHI plans
linkTitle: DHI
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
> To subscribe to DHI Select, see [Set up, upgrade, or downgrade a plan](../manage.md).

## Usage

DHI Community gives you access to hardened base images from a public registry at no cost or additional setup. Any organization can pull hardened base images directly from `dhi.io`.

When you upgrade from DHI Community to DHI Select, you purchase a set number of repositories that are mirrored into your organization's namespace. Entitlements are scoped to the organization account that you assign them to during checkout. All organization members can then pull from those mirrored repositories.

DHI Enterprise extends DHI Select with unlimited customizations, full catalog access, the Hardened System Packages repository, and an Extended Lifecycle Support add-on. To upgrade to DHI Enterprise, [contact Docker sales](https://www.docker.com/pricing/contact-sales/).

For details on setting up and managing repositories, see [Get started with DHI Select and Enterprise](/manuals/dhi/how-to/select-enterprise.md).

## Billing behaviors

DHI Select is an annual plan billed per repository from the date your plan starts. Repositories added mid-cycle are prorated for the remainder of the billing period.

To add more repositories, go to **Active plans** in the billing portal. For steps, see [Upgrade or downgrade a plan](../manage.md#upgrade-or-downgrade-a-plan).

## Downgrade

Cancellations and repository removals are deferred to the end of the current billing cycle. You can cancel renewal or remove repositories at any time. You cannot stop a plan mid-cycle to receive a partial refund.

To cancel or remove repositories:

1. Sign in to [Docker Home](https://app.docker.com/) and go to **Billing**.
1. From **Active plans**, select **Manage** next to your DHI plan.
1. Select **Cancel auto renewal** to cancel, or reduce the repository count to remove repositories.

- Cancellations and repository removals take effect at the end of the current annual billing cycle.
- Repository access remains active until the cycle ends.

If you're a DHI Enterprise customer, reach out to your sales representative to downgrade.
