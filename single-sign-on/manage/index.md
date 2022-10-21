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
