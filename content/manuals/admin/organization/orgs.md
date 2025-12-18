---
title: Create your organization
weight: 10
description: Learn how to create an organization.
keywords: docker organizations, organization, create organization, docker teams, docker admin console, organization management
aliases:
  - /docker-hub/orgs/
---

{{< summary-bar feature_name="Admin orgs" >}}

This page describes how to create an organization.

## Prerequisites

Before you begin creating an organization:

- You need a [Docker ID](/accounts/create-account/)
- Review the [Docker subscriptions and features](https://www.docker.com/pricing/)
  to determine what subscription to choose for your organization

## Create an organization

There are multiple ways to create an organization. You can either:

- Create a new organization using the **Create Organization** option in the
Admin Console or Docker Hub
- Convert an existing user account to an organization

The following section contains instructions on how to create a new organization. For prerequisites and
detailed instructions on converting an existing user account to an organization, see
[Convert an account into an organization](/manuals/admin/organization/convert-account.md).

To create an organization:

1. Sign in to [Docker Home](https://app.docker.com/) and navigate to the bottom
of the organization list.
1. Select **Create new organization**.
1. Choose a subscription for your organization, a billing cycle, and specify how many seats you need. See [Docker Pricing](https://www.docker.com/pricing/) for details on the features offered in the Team and Business subscription.
1. Select **Continue to profile**.
1. Select **Create an organization** to create a new one.
1. Enter an **Organization namespace**. This is the official, unique name for
your organization in Docker Hub. It's not possible to change the name of the
organization after you've created it.

   > [!NOTE]
   >
   > You can't use the same name for the organization and your Docker ID. If you want to use your Docker ID as the organization name, then you must first [convert your account into an organization](/manuals/admin/organization/convert-account.md).

1. Enter your **Company name**. This is the full name of your company. Docker
displays the company name on your organization page and in the details of any
public images you publish. You can update the company name anytime by navigating
to your organization's **Settings** page.
1. Select **Continue to billing** to continue.
1. Enter your organization's billing information and select **Continue to payment** to continue to the billing portal.
1. Provide your payment details and select **Purchase**.

You've now created an organization.

## View an organization

To view an organization in the Admin Console:

1. Sign in to [Docker Home](https://app.docker.com) and select your
organization.
1. From the left-hand navigation menu, select **Admin Console**.

The Admin Console contains many options that let you to
configure your organization.

## Merge organizations

> [!WARNING]
>
> If you are merging organizations, it is recommended to do so at the _end_ of
> your billing cycle. When you merge an organization and downgrade another, you
> will lose seats on your downgraded organization. Docker does not offer
> refunds for downgrades.

If you have multiple organizations that you want to merge into one, complete
the following steps:

1. Based on the number of seats from the secondary organization, [purchase additional seats](../../subscription/manage-seats.md) for the primary organization account that you want to keep.
1. Manually add users to the primary organization and remove existing users from the secondary organization.
1. Manually move over your data, including all repositories.
1. Once you're done moving all of your users and data, [downgrade](../../subscription/change.md) the secondary account to a free subscription. Note that Docker does not offer refunds for downgrading organizations mid-billing cycle.

> [!TIP]
>
> If your organization has a Docker Business subscription with a purchase
order, contact Support or your Account Manager at Docker.

## More resources

- [Video: Docker Hub Organizations](https://www.youtube.com/watch?v=WKlT1O-4Du8)
