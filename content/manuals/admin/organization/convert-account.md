---
description: Convert your Docker Hub user account into an organization
title: Convert an account into an organization
keywords: docker hub, hub, organization, convert account, migrate account
weight: 35
aliases:
- /docker-hub/convert-account/
---

{{< summary-bar feature_name="Admin orgs" >}}

Learn how to convert an existing user account into an organization. This is
useful if you need multiple users to access your account and the repositories
it’s connected to. Converting it to an organization gives you better control
over permissions for these users through
[teams](/manuals/admin/organization/manage-a-team.md) and
[roles](/manuals/security/for-admins/roles-and-permissions.md).

When you convert a user account to an organization, the account is migrated to
a Docker Team subscription by default.

## Prerequisites

Before you convert a user account to an organization, you must meet
the following requirements:

- The user account that you want to convert can not be a member of a company,
any existing teams, or another organization. You must remove the account from
all teams, organizations, or companies.
-  You must have a separate Docker ID ready to assign as the owner of the
organization during conversion.

    If you want to convert your user account into an organization account and
    you don't have any other user accounts, you need to create a new user
    account to assign it as the owner of the new organization. The assigned
    owner has full administrative access to manage the organization. You can
    assign more users the owner role after the conversion.

## What happens when you convert your account

The following happens when you convert your account into
an organization:

- This process removes the email address for the account. Notifications are
instead sent to organization owners. You'll be able to reuse the
removed email address for another account after converting.
- The current subscription will automatically cancel and your new subscription
will start.
- Repository namespaces and names won't change, but converting your account
removes any repository collaborators. Once you convert the account, you'll need
to add repository collaborators as team members.
- Existing automated builds appear as if they were set up by the first owner
added to the organization.
- The user account that you add as the first owner will have full
administrative access to configure and manage the organization.
- To transfer a user's personal access tokens (PATs) to your converted
organization, you must designate the user as an organization owner. This will
ensure any PATs associated with the user's account are transferred to the
organization owner.

> [!TIP]
>
> To avoid potentially disrupting service of personal access tokens when
converting an account or changing ownership, it is recommended to
use [organization access tokens](/manuals/security/for-admins/access-tokens.md).
Organization access tokens are associated with an organization, not a
single user account.

## Convert an account into an organization

> [!IMPORTANT]
>
> Converting an account into an organization is permanent. Back up any data
 or settings you want to retain.

1. Sign in to [Docker Home](https://app.docker.com/).
1. Select your avatar in the top-right corner to open the drop-down.
1. From **Account settings**, select **Convert**.
1. Review the warning displayed about converting a user account. This action
cannot be undone and has considerable implications for your assets and the
account.
1. Enter a **Username of new owner** to set an organization owner. The new
Docker ID you specify becomes the organization’s owner. You cannot use the
same Docker ID as the account you are trying to convert.
1. Select **Confirm**. The new owner receives a notification email. Use that
owner account to sign in and manage the new organization.
