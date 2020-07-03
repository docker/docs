---
description: Convert your Docker Hub user account into an organization
title: Convert an account into an organization
keywords: docker hub, hub, organization
---

You can convert an existing user account to an organization. You might want to do this if you need multiple users to access your account and the repositories that it’s connected to. Converting it to an organization gives you better control over permissions for these users through teams.

> **Note:**
>
> Once you convert your account to an organization, you can’t revert it to a user account, so make sure you really want to convert the user account.

## Prerequisites

Before you convert a user account to an organization, ensure that you have completed the following steps:

1. The user account that you wish to convert must not be a member of any teams or organizations. You must remove the account from all teams and organizations.

    Go to **Organizations**, select an organization from the list, and then click the Leave organization arrow next to your username in the members list.

    If the user account is the sole owner of any organization, add someone to the "owners" team and then remove yourself from the organization.

2. You must have a separate Docker ID ready to assign it as the owner of the organization during conversion.

    If you wish to convert your user account into an organization account and you do not have any other user accounts, you need to create a new user account to assign it as the owner of the new organization. This user account then becomes the first member of the "owners" team and has full administrative access to configure and manage the organization. You can add more users into the "owners" team after the conversion.

## Convert a Community account into an organization

1. Ensure you have removed your user account from all teams and organizations and that you have a new Docker ID before you convert an account. See the [Prerequisites](#prerequisites) section for details.

2. Click on your account name in the top navigation, then go to your **Account Settings**.

3. Under the **Convert Account** tab, click **Convert to Organization**.

4. Carefully review the warning displayed about converting a user account. This cannot be undone and will have considerable implications for your assets and the account.

5. As part of the conversion, you must enter a **Docker ID** to set an organization owner. This is the user account that will manage the organization, and the only way to access the organization settings after conversion. You cannot use the same Docker ID as the account you are trying to convert.

6. Click **Convert** to confirm. The new owner will receive a notification email. Use that owner account to log into your new organization.

    Your Community account has now been converted to an organization.

## Convert a Pro account into an organization

>**Note:**
>
> When you convert a Pro or a legacy individual repository plan to an organization, the account
will be migrated to a Team plan and will be charged $35 per month for 5 seats. For more information,
see [Docker Hub Pricing](https://hub.docker.com/pricing).

1. Ensure you have removed your user account from all teams and organizations and that you have a new Docker ID before you convert an account. See the [Prerequisites](#prerequisites) section for details.

2. Click on your account name in the top navigation bar, then go to your **Account Settings**.

3. Under the **Convert Account** tab, click **Convert to Organization**.

4. Carefully review the warning displayed about converting a user account. This cannot be undone and will have considerable implications for your assets and the account.

5. As part of the conversion, you must enter a **Docker ID** to set an organization owner. This is the user account that will manage the organization, and the only way to access the organization settings after conversion. You cannot use the same Docker ID as the account you are trying to convert.

6. Click **Convert** to confirm. The new owner will receive a notification email. Use that owner account to log into your new organization.

    Your Pro user account has now been converted to an organization.
