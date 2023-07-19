---
title: Deactivating an account or an organization
description: Learn how to deactivate a Docker Hub account or an organization
keywords: Docker Hub, delete, deactivate, account, organization
---

You can deactivate an account at any time. 

>**Warning**
>
> If your Docker Hub account or organization is linked to other Docker products and services, deactivating your account also removes access to those products and services.
{: .warning}

## Deactivate a user account

### Prerequisites

Before deactivating your Docker Hub account, ensure that you meet the following requirements:

- You must not be a member of a company or any teams or organizations. You must remove the account from all teams, organizations, or the company.

    To do this:
    1. Navigate to **Organizations** and then select the organization(s) you need to leave.
    2. Find your username in the **Members** tab.
    3. Select the **More options** menu and then select **Leave organization**.

- If you are the sole owner of an organization, either add someone to [the **owners** team](manage-a-team.md#the-owners-team) and then remove yourself from the organization, or deactivate the organization. Similarly, if you are the sole owner of a company, either add someone else as a company owner and then remove yourself, or deactivate the company.

- If you have an active subscription, [downgrade it to a Docker Personal subscription](../subscription/downgrade.md).

- Download any images and tags you want to keep. Use `docker pull -a <image>:<tag>`.

- Unlink your [Github and Bitbucket accounts](../docker-hub/builds/link-source.md#unlink-a-github-user-account).

### Deactivate

Once you have completed all the steps above, you can deactivate your account. 

1. Select your account name in the top-right of Docker Hub and from the drop-down menu, select **Account Settings**.
2. From the **Deactivate Account** tab, select **Deactivate account**. 

> This cannot be undone. Be sure you've gathered all the data you need from your account before deactivating it.
{: .warning }


## Deactivate an organization

Before deactivating an organization, please complete the following:

- Download any images and tags you want to keep:
  `docker pull -a <image>:<tag>`.

-  If you have an active subscription, [downgrade it to a **Docker Free Team** subscription](../subscription/downgrade.md).

- Remove all other members, including those in the **Owners** team, within the organization.

- Unlink your [Github and Bitbucket accounts](../docker-hub/builds/link-source.md#unlink-a-github-user-account).

### Deactivate

Once you have completed all the steps above, you can deactivate your organization. 

1. On Docker Hub, select **Organizations**.
2. Choose the organization you want to deactivate. 
3. In **Settings**, select the **Deactivate Org** tab and then **Deactivate organization**.

> This cannot be undone. Be sure you've gathered all the data you need from your organization before deactivating it.
{: .warning }
