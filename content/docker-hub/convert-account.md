---
description: Convert your Docker Hub user account into an organization
title: Convert an account into an organization
keywords: docker hub, hub, organization
---

You can convert an existing user account to an organization. This is useful if you need multiple users to access your account and the repositories that it’s connected to. Converting it to an organization gives you better control over permissions for these users through [teams](manage-a-team.md).

When you convert a user account to an organization, the account
is migrated to a Team plan and charged $35 per month for 5 seats. For more information,
see [Docker Pricing](https://www.docker.com/pricing).

> **Important**
>
> Once you convert your account to an organization, you can’t revert it to a user account. 
{ .important }

## Prerequisites

Before you convert a user account to an organization, ensure that you meet the following requirements:

- The user account that you wish to convert must not be a member of a company or any teams or organizations. You must remove the account from all teams, organizations, or the company.

    To do this:
    1. Navigate to **Organizations** and then select the organization(s) you need to leave.
    2. Find your username in the **Members** tab.
    3. Select the **More options** menu and then select **Leave organization**.

    If the user account is the sole owner of any organization or company, assign another user the owner role and then remove yourself from the organization or company.

-  You must have a separate Docker ID ready to assign it as the owner of the organization during conversion.

    If you wish to convert your user account into an organization account and you do not have any other user accounts, you need to create a new user account to assign it as the owner of the new organization. With the owner role assigned, this user account has full administrative access to configure and manage the organization. You can assign more users the owner role after the conversion.

## Convert an account into an organization

1. Ensure you have removed your user account from any company or teams or organizations. Also make sure that you have a new Docker ID before you convert an account. See the [Prerequisites](#prerequisites) section for details.

2. In the top-right of Docker Hub, select your account name and then from the drop-down menu, select **Account Settings**.

3. From the **Convert Account** tab, select **Convert to Organization**.

4. Review the warning displayed about converting a user account. This action cannot be undone and has considerable implications for your assets and the account.

5. Enter a **Docker ID** to set an organization owner. This is the user account that will manage the organization, and the only way to access the organization settings after conversion. You cannot use the same Docker ID as the account you are trying to convert.

6. Select **Convert and Purchase** to confirm. The new owner receives a notification email. Use that owner account to log in to your new organization.