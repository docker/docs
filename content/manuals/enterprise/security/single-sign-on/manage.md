---
title: Manage single sign-on
linkTitle: Manage
description: Learn how to manage Single Sign-On for your organization or company.
keywords: manage, single sign-on, SSO, sign-on, admin console, admin, security, domains, connections, users, provisioning
aliases:
- /admin/company/settings/sso-management/
- /single-sign-on/manage/
- /security/for-admins/single-sign-on/manage/
---

{{< summary-bar feature_name="SSO" >}}

This page covers how to manage single sign-on (SSO) after initial setup,
including managing domains, connections, users, and provisioning
settings.

## Manage domains

### Add a domain

To add a domain to an existing SSO connection:

1. Sign in to [Docker Home](https://app.docker.com) and select your company or
organization from the top-left account drop-down.
1. Select **Admin Console**, then **SSO and SCIM**.
1. In the SSO connections table, select the **Actions** menu for your
connection, then select **Edit connection**.
1. Select **Next** to navigate to the domains section.
1. In the **Domains** section, select **Add domain**.
1. Enter the domain you want to add to the connection.
1. Select **Next** to confirm or change the connected organizations.
1. Select **Next** to confirm or change the default organization and
team provisioning selections.
1. Review the connection details and select **Update connection**.

### Remove a domain from an SSO connection

> [!IMPORTANT]
>
> If you use multiple identity providers with the same domain, you must remove the domain from each SSO connection individually.

1. Sign in to [Docker Home](https://app.docker.com) and select your company or organization from the top-left account drop-down.
1. Select **Admin Console**, then **SSO and SCIM**.
1. In the **SSO connections** table, select the **Actions** menu for your connection, then
**Edit connection**.
1. Select **Next** to navigate to the domains section.
1. In the **Domain** section, select the **X** icon next to the domain
you want to remove.
1. Select **Next** to confirm or change the connected organizations.
1. Select **Next** to confirm or change the default organization and
team provisioning selections.
1. Review the connection details and select **Update connection**.

> [!NOTE]
>
> When you re-add a domain, Docker assigns a new TXT record value. You must complete domain verification again with the new TXT record.

## Manage SSO connections

### View connections

To view all configured SSO connections:

1. Sign in to [Docker Home](https://app.docker.com) and select your company or organization from the top-left account drop-down.
1. Select **Admin Console**, then **SSO and SCIM**.
1. View all configured connections in the **SSO connections** table.

### Edit a connection

To modify an existing SSO connection:

1. Sign in to [Docker Home](https://app.docker.com) and select your company or organization from the top-left account drop-down.
1. Select **Admin Console**, then **SSO and SCIM**.
1. In the **SSO connections** table, select the **Actions** menu for your connection, then
**Edit connection**.
1. Follow the on-screen instructions to modify your connection settings.

### Delete a connection

To remove an SSO connection:

1. Sign in to [Docker Home](https://app.docker.com) and select your company or organization from the top-left account drop-down.
1. Select **Admin Console**, then **SSO and SCIM**.
1. In the **SSO connections** table, select the **Actions** menu for your connection, then
**Delete connection**.
1. Follow the on-screen instructions to confirm the deletion.

> [!WARNING]
>
> Deleting an SSO connection removes access for all users who authenticate through
that connection.

## Manage users and provisioning

Docker automatically provisions users through Just-in-Time (JIT) provisioning when they sign in via SSO. You can also manually manage users and configure different provisioning methods.

### How provisioning works

Docker supports the following provisioning methods:

- JIT provisioning (default): Users are automatically added to your organization
when they sign in via SSO
- SCIM provisioning: Sync users and groups from your identity provider to Docker
- Group mapping: Sync user groups from your identity provider with teams in your Docker organization
- Manual provisioning: Turn off automatic provisioning and manually invite users

For more information on provisioning methods, see [Provision users](/manuals/enterprise/security/provisioning/_index.md).

### Add guest users

To invite users who don't authenticate through your identity provider:

1. Sign in to [Docker Home](https://app.docker.com/) and select
your organization.
1. Select **Members**.
1. Select **Invite**.
1. Follow the on-screen instructions to invite the user.

The user receives an email invitation and can create a Docker account or sign
in with their existing account.

### Remove users

To remove a user from your organization:

1. Sign in to [Docker Home](https://app.docker.com/) and select
your organization.
1. Select **Members**.
1. Find the user you want to remove and select the **Actions** menu next to their name.
1. Select **Remove** and confirm the removal.

The user loses access to your organization immediately upon removal.
