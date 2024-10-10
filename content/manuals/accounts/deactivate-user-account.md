---
title: Deactivate an account
weight: 30
description: Learn how to deactivate a Docker user account.
keywords: Docker Hub, delete, deactivate, account, account management
---

You can deactivate an account at any time. This section describes the prerequisites and steps to deactivate a user account. For information on deactivating an organization, see [Deactivating an organization](../admin/deactivate-account.md).

>[!Warning]
>
> All Docker products and services that use your Docker account will be inaccessible after deactivating your account.

## Prerequisites

Before deactivating your Docker account, ensure that you meet the following requirements:

- For owners, you must leave your organization or company before deactivating your Docker account.
    To do this:
    1. Sign in to the [Docker Admin Console](https://app.docker.com/admin).
    2. Select the organization you need to leave from the top-left drop-down menu.
    3. Find your username in the **Members** tab.
    4. Select the **More options** menu and then select **Leave organization**.

- If you are the sole owner of an organization, you must assign the owner role to another member of the organization and then remove yourself from the organization, or deactivate the organization. Similarly, if you are the sole owner of a company, either add someone else as a company owner and then remove yourself, or deactivate the company.

- If you have an active Docker subscription, [downgrade it to a Docker Personal subscription](../subscription/core-subscription/downgrade.md).

- If you have an active Docker Build Cloud Paid subscription, [downgrade it to a Docker Build Cloud Starter subscription](../billing/build-billing.md#downgrade-your-subscription).

- If you have an active Docker Scout subscription, [downgrade it to a Docker Scout Free subscription](../billing/scout-billing.md#downgrade-your-subscription).

- Download any images and tags you want to keep. Use `docker pull -a <image>:<tag>`.

- Unlink your [GitHub and Bitbucket accounts](../docker-hub/builds/link-source.md#unlink-a-github-user-account).

## Deactivate

Once you have completed all the previous steps, you can deactivate your account.

> [!WARNING]
>
> This cannot be undone. Be sure you've gathered all the data you need from your account before deactivating it.

1. Sign in to [Docker Home](https://app.docker.com/login).
2. Select your avatar to open the drop-down menu.
3. Select **Account settings**.
4. In the **Account management** section, select **Deactivate account**.
5. To confirm, select **Deactivate account**.
