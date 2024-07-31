---
title: Deactivating an organization
description: Learn how to deactivate a Docker organization.
keywords: Docker Hub, delete, deactivate organization, account, organization management
aliases:
- /docker-hub/deactivate-account/
---

You can deactivate an account at any time. This section describes the prerequisites and steps to deactivate an organization account. For information on deactivating a user account, see [Deactivate a user account](../accounts/deactivate-user-account.md).

>**Warning**
>
> All Docker products and services that use your Docker account or organization account will be inaccessible after deactivating your account.
{ .warning }

## Prerequisites

Before deactivating an organization, complete the following:

- Download any images and tags you want to keep:
  `docker pull -a <image>:<tag>`.

- If you have an active Docker subscription, [downgrade it to a free subscription](../subscription/core-subscription/downgrade.md).

- If you have an active Docker Scout subscription, [downgrade it to a Docker Scout Free subscription](../billing/scout-billing.md#downgrade-your-subscription).

- Remove all other members within the organization.

- Unlink your [Github and Bitbucket accounts](../docker-hub/builds/link-source.md#unlink-a-github-user-account).

## Deactivate

Once you have completed all the previous steps, you can deactivate your organization.

> **Warning**
>
> This cannot be undone. Be sure you've gathered all the data you need from your organization before deactivating it.
{ .warning }

{{< tabs >}}
{{< tab name="Admin Console" >}}

1. In Admin Console, choose the organization you want to deactivate.
2. Under **Organization settings**, select **Deactivate**.
3. Enter the organization name to confirm deactivation.
4. Select **Deactivate organization**.

{{< /tab >}}
{{< tab name="Hub" >}}

1. On Docker Hub, select **Organizations**.
2. Choose the organization you want to deactivate.
3. In **Settings**, select the **Deactivate Org** tab and then **Deactivate organization**.

{{< /tab >}}
{{< /tabs >}}
