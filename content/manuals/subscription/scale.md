---
title: Scale your subscription
linkTitle: Scale
description: Scale Docker Build Cloud, Testcontainers Cloud, Gordon, and DHI consumption for your subscription
keywords: scale subscription, docker build cloud minutes, testcontainers cloud minutes, gordon plan, docker hardened images, dhi repositories, usage scaling
weight: 40
---

Docker subscriptions include basic entitlements that you can scale with additional add-ons as your needs grow. You can purchase different add-ons that extend usage for your personal or organization accounts:

- Docker Hardened Images (DHI) Select repositories for organization accounts
- Gordon Plus, Max, and Ultra plans for personal accounts
- Seats for Docker Core subscriptions like Docker Team and Docker Business
- Licenses for individual products, like AI Governance (available for Docker Sandboxes) and Docker Offload
- Build minutes for Docker Build Cloud
- Cloud runtime minutes for Testcontainers

## Add DHI Select to organization accounts

DHI Select lets organization admins assign DHI repositories to organization accounts. All organization members can then access those DHI repositories, including members you add to the organization after initial setup.

To purchase hardened repositories with DHI Select:

1. Sign in to [Docker Home](https://app.docker.com/) and select your organization.
1. Select **Billing**, then **Browse products**.
1. Select **Hardened Images** from the products page.
1. Choose the organization account that should receive the DHI repository, then **Continue**. You can add repositories for one account at a time.
1. Select **Add repositories**, then add the number of new repositories you want to add to the account.
1. Verify your billing details, continue to payment, and complete checkout.

To manage the repositories in your active DHI Select plan:

1. From the Billing overview page, select **Manage** next to **Hardened Images**.
1. Choose from the available actions:
   - **Add repositories** lets you add additional repositories to your active plan. You will verify your billing details, continue to payment, and complete checkout.
   - **Remove repositories** lets you remove repositories.
   - **Disable auto-renewal** cancels your DHI Select plan at the end of the current billing cycle.

To learn more about DHI Select, see [Get started with DHI Select and Enterprise](/manuals/dhi/how-to/select-enterprise.md). Purchasing eight or more hardened repositories? [Contact Docker sales](https://www.docker.com/pricing/contact-sales/) to discuss an Enterprise plan.

## Add a Gordon Plus, Max, or Ultra plan to personal accounts

> [!IMPORTANT]
> Gordon subscriptions apply to personal Docker accounts only.
> If you purchase a Gordon plan while signed in with an organization account,
> the subscription applies to your personal account automatically.

Gordon Plus, Max, and Ultra plans increase your maximum Gordon usage allowance for your personal account. Gordon plans are billed at a monthly rate at the first of the month. If you purchase a Gordon plan after the first of the month, you'll be billed on a prorated basis.

To add a new Gordon plan to your account:

1. Sign in to [Docker Home](https://app.docker.com/) and select your Docker personal account.
1. Select **Billing**, then **Browse products**.
1. Select **View plans** from the **Upgrade Gordon** section of the products page.
1. Choose the Gordon subscription plan you want to apply to your personal account.
1. Verify your billing details, continue to payment, and complete checkout.

If you have an active Gordon plan and want to upgrade to a higher tier, your new usage allowance takes effect immediately. You pay a prorated charge for the higher tier for the rest of the current billing period.

To upgrade your Gordon plan:

1. From the Billing overview page, select **Manage** from the Gordon tile under your active plans.
1. Select **Upgrade plan** and choose the new plan you want to upgrade to.
1. Verify your billing details, continue to payment, and complete checkout.

To learn about Gordon usage, see [Gordon usage limits and tiers](/manuals/ai/gordon/usage-limits.md).

## Add licenses

Licenses add a layer to your Docker subscriptions, letting you assign products to select members of your
organization without consuming a Docker Core seat. You can purchase licenses for some Docker products:

- AI Governance
- Docker Offload

Once you've purchased licenses for your organization, you can [manage license assignment](/manuals/admin/organization/manage/manage-licenses.md) from the **Members** page in Docker Home.

> [!TIP]
> To purchase licenses for AI Governance and Docker Offload, <a href="https://www.docker.com/pricing/contact-sales/" id="dkr_docs_cs_subscription_scale_licenses" class="link" rel="noopener">contact sales</a>. 

### AI Governance licenses

AI Governance licenses let admins create and apply organization-wide AI Governance policies for license-holding members.

- AI Governance licenses apply to Docker Sandbox usage.
- AI Governance licenses are separate from Docker Team or Business, so you can purchase AI Governance licenses without an existing subscription.
- When you purchase AI Governance licenses, you may assign them to organization members without a Docker Core seat.

When a member who holds an AI Governance license uses Docker Sandbox, the organization policy overrides the member’s local policy rules. Members without AI Governance licenses can still use Docker Sandbox, but organization policies will not govern their usage.

### Docker Offload licenses

Docker Offload lets developers offload building and running containers to the cloud.

- You must have a Docker Team or Docker Business subscription.
- You can only assign Docker Offload licenses to members with Docker Team or Docker Business.

## Add minutes

Minutes don't roll over. Base subscription minutes reset each billing period and don't accumulate. Additional purchased minutes expire at the end of your subscription period.

For example, with an annual Docker Team subscription (500 included minutes), if you purchase 500 additional minutes, only the additional 500 minutes roll over until your annual renewal.

### Docker Build Cloud build minutes

Purchase additional build minutes through the Docker Build Cloud Dashboard:

1. Sign in to [Docker Home](https://app.docker.com/) and choose
   your organization.
1. Select **Build Cloud**, then **Build minutes**.
1. Select **Add minutes**.
1. Select your additional minute amount, then **Continue to payment**.
1. Enter your payment details and billing address.
1. Review your order and select **Pay**.

Your additional minutes appear on the Build minutes page immediately.

### Docker Testcontainers Cloud runtime minutes

You can add Testcontainers Cloud runtime minutes in two ways:

- <a href="https://www.docker.com/pricing/contact-sales/" id="dkr_docs_cs_subscription_scale_runtime" class="link" rel="noopener">Contact sales</a> to pre-purchase runtime minutes at discounted rates
- Use unlimited runtime minutes on-demand with billing at the end of each monthly cycle

On-demand usage is billed at higher rates than pre-purchased capacity. To avoid higher on-demand charges, pre-purchase additional minutes if you expect consistent usage over your subscription's included minutes.

## What's next

- [Manage licenses](/manuals/admin/organization/manage/manage-licenses.md)
- [Manage seats](/manuals/admin/organization/manage/manage-seats.md)
- [View your consumption](../admin/organization/manage/manage-products.md#monitor-product-usage-for-your-organization)
- [Docker Build Cloud overview](/manuals/build-cloud/_index.md)
- [DHI Select and Enterprise quickstart](/manuals/dhi/how-to/select-enterprise.md)
- [Testcontainers overview](/manuals/testcontainers.md)
