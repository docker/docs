---
description: Manage SSO
keywords: manage, single sign-on, SSO, sign-on
title: Manage SSO
---

## Manage domains

### Remove a domain from an SSO connection

1. In the **Single Sign-On Connection** table, select the **Action** icon and then **Edit connection**.
2. Select **Next** to navigate to the section where the connected domains are listed.
3. In the **Domain** drop-down, select the **Remove** icon next to the domain that you want to remove.
4. Select **Next** to confirm or change the connected organization(s).
5. Select **Next** to confirm or change the default organization and team provisioning selections.
6. Review the **Connection Summary** and select **Save**.

>**Note**
>
>If you want to re-add the domain, a new TXT record value is assigned. You must then complete the verification steps with the new TXT record value.

## Manage organizations

>**Note**
>
>You must have a [company](../../docker-hub/creating-companies.md) to manage more than one organization.

### Connect an organization 

1. In the **Single Sign-On Connection** table, select the **Action** icon and then **Edit connection**.
2. Select **Next** to navigate to the section where connected organizations are listed.
3. In the **Organizations** drop-down, select the organization to add to the connection.
4. Select **Next** to confirm or change the default organization and team provisioning.
5. Review the **Connection Summary** and select **Save**.

### Remove an organization

1. In the **Single Sign-On Connection** table, select the **Action** icon and then **Edit connection**.
2. Select **Next** to navigate to the section where connected organizations are listed.
3. In the **Organizations** drop-down, select **Remove** to remove the connection.
4. Select **Next** to confirm or change the default organization and team provisioning.
5. Review the **Connection Summary** and select **Save**.

## Manage SSO connections

### Edit a connection

1. In the **Single Sign-On Connection** table, select the **Action** icon.
2. Select **Edit connection** to edit you connection.
3. Continue with the on-screen instructions.

### Delete a connection

1. In the **Single Sign-On Connection** table, select the **Action** icon.
2. Select **Delete** and **Delete connection**.
3. Continue with the on-screen instructions.

### Deleting SSO

When you disable SSO, you can delete the connection to remove the configuration settings and the added domains. Once you delete this connection, it can't be undone. Users must authenticate with their Docker ID and password or create a password reset if they don't have one.

## Manage users

### Add guest users when SSO is enabled

To add a guest to your organization in Docker Hub if they aren’t verified through your IdP:

1. Go to **Organizations** in Docker Hub, and select your organization.
2. Select **Add Member**, enter the email address, and select a team from the drop-down list.
3. Select **Add** to confirm.

### Remove users from the SSO organization

To remove a user from an organization:

1. Go to **Organizations** in Docker Hub, and select your organization.
2. From the **Members** tab, select the **x** next to a member’s name to remove them from all the teams in the organization.
3. Select **Remove** to confirm. The member receives an email notification confirming the removal.

