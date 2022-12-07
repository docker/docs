---
description: Manage SSO
keywords: manage, single sign-on, SSO, sign-on
title: Manage SSO
---

## Manage users when SSO is enabled

You don’t need to add users to your organization in Docker Hub manually. You just need to make sure an account for your users exists in your IdP.

 > **Note**
 >
 > When you enable SSO for your organization, a first-time user can sign in to Docker Hub using their company's domain email address. They're then added to your organization and assigned to your company's team.

To add a guest to your organization in Docker Hub if they aren’t verified through your IdP:

1. Go to **Organizations** in Docker Hub, and select your organization.
2. Select **Add Member**, enter the email address, and select a team from the drop-down list.
3. Select **Add** to confirm.

## Audit your domain for users

When your organization has configured SSO and registered a domain, you can audit your domain. Auditing your domain will identify uncaptured users that have an email associated with your domain, but they aren't a member of your organization.

Uncaptured users can pose a security threat to your environment since your organization's security settings aren't applied to the user's sessions.

You can add uncaptured users to your organization in order to apply your organization's security settings. Additionally, you can enforce sign-in to ensure that only members of your organization can sign in to Docker Desktop in your environment. For more details about enforcing sign-in, see [Configure registry.json to enforce sign-in](../../docker-hub/configure-sign-in.md).

To audit your domain:

1. Sign in to [Docker Hub](https://hub.docker.com){: target="_blank" rel="noopener" class="_"} as an administrator of your organization.

2. Select **Organizations*** and then select your organization.

3. Select **Settings** and then select **Security**.

4. In **Domain Audit**, select **Export Users** to export a CSV file containing the following columns:
  - Name: The name of the user.
  - Username: The Docker ID of the user.
  - Email: The email address of the user.
  - Date Joined: The date the user created their Docker account.

You can invite all the uncaptured users to your organization using the exported CSV file. For more details, see [Invite members via CSV file](../../docker-hub/members.md/#invite-members-via-csv-file).

## Remove users from the SSO organization

To remove a user from an organization:

1. Sign in to [Docker Hub](https://hub.docker.com){: target="_blank" rel="noopener" class="_"} as an administrator of your organization.
2. Select the organization from the list. The organization page displays a list of user.
3. Select the **x** next to a member’s name to remove them from all the teams in the organization.
4. Select **Remove** to confirm. The member will receive an email notification confirming the removal.

    > **Note**
    >
    > When you remove a member from an SSO organization, they're unable to log
    > in using their email address.


## Deleting SSO

When you disable SSO, you can delete the connection to remove the configuration settings and the added domains. Once you delete this connection, it can't be undone. Users must authenticate with their Docker ID and password or create a password reset if they don't have one.

![Delete SSO](/single-sign-on/images/delete-sso.png){:width="500px"}

## FAQs

To learn more see [FAQs](../faqs.md).
