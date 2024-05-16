---
title: Deactivating an account or an organization
description: Learn how to deactivate a Docker Hub account or an organization
keywords: Docker Hub, delete, deactivate, account, organization
aliases:
- /docker-hub/deactivate-account/
---

You can deactivate an account at any time. 

>**Warning**
>
> All Docker products and services that use your Docker account or organization account will be inaccessible after deactivating your account.
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

- If you have an active Docker Build Cloud Team subscription, [downgrade it to a Docker Build Cloud Starter subscription](../billing/build-billing.md#downgrade-your-subscription).

- If you have an active Docker Scout subscription, [downgrade it to a Docker Scout Free subscription](../billing/scout-billing.md#downgrade-your-subscription).

- Download any images and tags you want to keep. Use `docker pull -a <image>:<tag>`.

- Unlink your [Github and Bitbucket accounts](../docker-hub/builds/link-source.md#unlink-a-github-user-account).

### Deactivate

Once you have completed all the steps above, you can deactivate your account. 

1. Select your account name in the top-right of Docker Hub and from the drop-down menu, select **My Account**.
2. From the **Deactivate Account** tab, select **Deactivate account**. 

> This cannot be undone. Be sure you've gathered all the data you need from your account before deactivating it.
{ .warning }


## Deactivate an organization

Before deactivating an organization, complete the following:

- Download any images and tags you want to keep:
  `docker pull -a <image>:<tag>`.

- If you have an active Docker subscription, [downgrade it to a **Docker Free Team** subscription](../subscription/core-subscription/downgrade.md).

- If you have an active Docker Scout subscription, [downgrade it to a Docker Scout Free subscription](../billing/scout-billing.md#downgrade-your-subscription).

- Remove all other members within the organization.

- Unlink your [Github and Bitbucket accounts](../docker-hub/builds/link-source.md#unlink-a-github-user-account).

### Deactivate

Once you have completed all previous the steps, you can deactivate your organization.

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
