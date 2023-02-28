---
description: Manage SSO
keywords: manage, single sign-on, SSO, sign-on
title: Manage SSO
---

# Manage domains

## Delete a domain



1. In the **Single Sign-On Connection**, select the **Action** icon and **Edit**.
2. Select **Next** to navigate to the section where the connected domains are listed.
3. In the **Domain** drop-down, select the **Remove** icon next to the domain that you want to remove.
4. Select **Next** to confirm or change the connected organization(s).
5. Select **Next** to confirm or change the default organization and team provisioning selections.
6. Review the **Connection Summary** and select **Save**.

>**Note**
    >
    >If you want to add this domain again, a new TXT record value is assigned. You must complete the verification steps with the new TXT record value.

# Manage organizations

## Connect an organization

You must have a company to connect additional organizations.

1. In the **Single Sign-on Connections** section, select the **Action** icon and **Edit**.
2. Select **Next** to navigate to the section where connected organizations are listed.
3. In the **Organizations** drop-down, select the organization to add to the connection.
4. Select **Next** to confirm or change the default organization and team provisioning.
5. Review the **Connection Summary** and select **Save**.

## Remove an organization

1. In the **Single Sign-on Connection** section, select the **Action** icon and **Edit**.
2. Select **Next** to navigate to the section where connected organizations are listed.
3. In the **Organizations** drop-down, select **Remove** to remove the connection.
4. Select **Next** to confirm or change the default organization and team provisioning.
5. Review the **Connection Summary** and select **Save**.

# Manage SSO connections

## Edit a connection

1. In the **Single Sign-On Connections**, select the **Action** icon.
2. Select **Edit Connection** to edit you connection.
3. Continue with the on-screen instructions.

## Delete a connection

1. In the **Single Sign-On Connections**, select the **Action** icon.
2. Select **Delete** and **Delete Connection**.
3. Continue with the on-screen instructions.

## Deleting SSO

When you disable SSO, you can delete the connection to remove the configuration settings and the added domains. Once you delete this connection, it can't be undone. Users must authenticate with their Docker ID and password or create a password reset if they don't have one.

# Manage users

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
















