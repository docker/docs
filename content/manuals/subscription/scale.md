---
title: Scale your subscription
linkTitle: Scale
description: Scale Docker Build Cloud and Testcontainers Cloud consumption for your subscription
keywords: scale subscription, docker build cloud minutes, testcontainers cloud minutes, usage scaling
weight: 30
---

Docker subscriptions let you scale consumption as your needs grow. All paid Docker subscriptions include base amounts of Docker Build Cloud build minutes and Testcontainers Cloud runtime minutes that you can supplement with additional capacity.

You can scale consumption for:

- Docker Build Cloud build minutes
- Testcontainers Cloud runtime minutes
- Docker Hardened Images (DHI) repositories

## Add licenses

Licenses add a layer to your Docker subscriptions, letting you assign products to select members of your
organization. You can purchase licenses for some Docker products:

- AI Governance
- Docker Offload 

Once you've purchased licenses for your organization, you can [manage license assignment](/manuals/admin/organization/manage/manage-licenses.md) from the **Members** page in Docker Home.

> [!TIP]
> To purchase licenses for AI Governance and Docker Offload, [contact sales](https://www.docker.com/pricing/contact-sales/). 

### AI Governance licenses

AI Governance licenses let admins create and apply organization-wide AI Governance policies for license-holding members.

- AI Governance licenses are separate from Docker Team or Business. 
- You can purchase AI Governance licenses without an existing subscription.
- When you purchase AI Governance licenses, you may assign them to organization members without a Docker subscription seat. 

When a member who holds an AI governance license uses AI products like Docker Sandbox, the organization policy overrides the member’s local policy rules. Members without AI Governance licenses can still use AI products. 

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

- [Contact sales](https://www.docker.com/pricing/contact-sales/) to pre-purchase runtime minutes at discounted rates
- Use unlimited runtime minutes on-demand with billing at the end of each monthly cycle

On-demand usage is billed at higher rates than pre-purchased capacity. To avoid higher on-demand charges, pre-purchase additional minutes if you expect consistent usage over your subscription's included minutes.

## Add DHI repositories

To add more hardened repositories to your DHI Select plan:

1. Sign in to [Docker Home](https://app.docker.com/) and select your
   organization.
1. Select **Billing**.
1. On the Overview page, select **Manage** next to **Hardened Images**.
1. Select how many repositories the account can use.

Purchasing eight or more hardened repositories? [Contact Docker sales](https://www.docker.com/pricing/contact-sales/) to discuss an Enterprise plan.

## What's next

- [Manage licenses](/manuals/admin/organization/manage/manage-licenses.md)
- [Manage seats](/manuals/admin/organization/manage/manage-seats.md)
- [View your consumption](../admin/organization/manage/manage-products.md#monitor-product-usage-for-your-organization) 
- [Docker Build Cloud overview](/manuals/build-cloud/_index.md)
- [DHI Select and Enterprise quickstart](/manuals/dhi/how-to/select-enterprise.md)
- [Testcontainers overview](/manuals/testcontainers.md)



