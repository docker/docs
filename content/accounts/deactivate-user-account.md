---
title: Deactivating a user account
description: Learn how to deactivate a Docker user account.
keywords: Docker Hub, delete, deactivate, account, account management
---

> **Early Access**
>
> Docker Home is in [Early Access (EA)](/release-lifecycle/#early-access-ea) for select users. If your account isn't selected for EA,
> you can manage your account settings, personal access tokens, and two-factor authentication for
> your account in [Docker Hub](https://hub.docker.com/).
{ .restricted }

You can deactivate an account at any time. This section describes the prerequisites and steps to deactivate a user account. For information on deactivating an organization, see [Deactivating an organization](../admin/deactivate-account.md).

>**Warning**
>
> All Docker products and services that use your Docker account will be inaccessible after deactivating your account.
{ .warning }

## Deactivate a user account

### Prerequisites

Before deactivating your Docker account, ensure that you meet the following requirements:

- You must not be a member of a company or any teams or organizations. You must remove the account from all teams, organizations, or the company.

    To do this:
    1. Navigate to **Organizations** and then select the organization(s) you need to leave.
    2. Find your username in the **Members** tab.
    3. Select the **More options** menu and then select **Leave organization**.

- If you are the sole owner of an organization, either assign another member of the organization the owner role and then remove yourself from the organization, or deactivate the organization. Similarly, if you are the sole owner of a company, either add someone else as a company owner and then remove yourself, or deactivate the company.

- If you have an active Docker subscription, [downgrade it to a Docker Personal subscription](../subscription/core-subscription/downgrade.md).

- If you have an active Docker Build Cloud Paid subscription, [downgrade it to a Docker Build Cloud Starter subscription](../billing/build-billing.md#downgrade-your-subscription).

- If you have an active Docker Scout subscription, [downgrade it to a Docker Scout Free subscription](../billing/scout-billing.md#downgrade-your-subscription).

- Download any images and tags you want to keep. Use `docker pull -a <image>:<tag>`.

- Unlink your [Github and Bitbucket accounts](../docker-hub/builds/link-source.md#unlink-a-github-user-account).

### Deactivate

Once you have completed all the previous steps, you can deactivate your account.

> **Warning**
>
> This cannot be undone. Be sure you've gathered all the data you need from your account before deactivating it.
{ .warning }

1. Sign in to your [Docker account](https://app.docker.com/login).
2. In Docker Home, select your avatar in the top-right corner to open the drop-down.
3. Select **My Account** to go to your account settings.
4. In the **Account management** section, select **Deactivate account** to open to deactivate account page.
5. To confirm, select **Deactivate account**.
