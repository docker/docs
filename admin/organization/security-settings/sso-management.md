---
description: Manage SSO
keywords: manage, single sign-on, SSO, sign-on
title: Manage Single Sign-On
---

{% include admin-early-access.md %}

## Manage domains

### Remove a domain from an SSO connection

1. In the **Single Sign-On Connection** table, select the **Action** icon and then **Edit connection**.
2. Select **Next** to navigate to the section where the connected domains are listed.
3. In the **Domain** drop-down, select the **Remove** icon next to the domain that you want to remove.
4. Select **Next** to confirm or change the connected organization(s).
5. Select **Next** to confirm or change the default organization and team provisioning selections.
6. Review the **Connection Summary** and select **Save**.

> **Note**
>
> If you want to re-add the domain, a new TXT record value is assigned. You must then complete the verification steps with the new TXT record value.

## Manage organizations

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

> **Important**
>
> SSO has Just-In-Time (JIT) Provisioning enabled by default. This means your users are auto-provisioned into a team called 'Company' within your organization.
>
> You can change this on a per-app basis. To prevent auto-provisioning users, you can create a security group in your IdP and configure the SSO app to authenticate and authorize only those users that are in the security group. Follow the instructions provided by your IdP:
> - [Okta](https://help.okta.com/en-us/Content/Topics/Security/policies/configure-app-signon-policies.htm)
> - [AzureAD](https://learn.microsoft.com/en-us/azure/active-directory/develop/howto-restrict-your-app-to-a-set-of-users)
{: .important}

### Add guest users when SSO is enabled

To add a guest to your organization if they aren’t verified through your IdP:


1. Sign in to [Docker Admin](https://admin.docker.com){: target="_blank" rel="noopener" class="_"}.
2. In the left navigation, select your organization in the drop-down menu.
3. Select **Members**.
4. Select **Invite**, enter the email address, and select an organization and team from the drop-down lists.
5. Select **Invite** to confirm.

### Remove users from the SSO company

To remove a user from an organization:

1. Sign in to [Docker Admin](https://admin.docker.com){: target="_blank" rel="noopener" class="_"}.
2. In the left navigation, select your oranization in the drop-down menu.
3. Select **Members**.
4. Select the action icon next to a user’s name, and then select **Remove member**.
5. Follow the on-screen instructions to remove the user.

## What's next?

- [Set up SCIM](scim.md)
- [Enable Group mapping](group-mapping.md)
