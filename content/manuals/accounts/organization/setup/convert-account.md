---
title: Convert a Docker account to an organization
linkTitle: Convert account
weight: 20
description: >-
  Convert an individual Docker account into an organization. Keep the same
  Docker ID and repository namespace, then assign a new owner for shared
  access.
keywords:
  - convert Docker account
  - convert account to organization
  - individual account to organization
  - organization owner
  - Docker organization
  - organization namespace
aliases:
  - /admin/organization/setup/convert-account/
  - /docker-hub/convert-account/
  - /admin/organization/convert-account/
toc_max: 2
---

{{< summary-bar feature_name="Admin orgs" >}}

When you convert, your
[Docker ID](/manuals/accounts/_index.md) becomes an organization instead of
an individual account. The Docker ID stays the same, so repository
namespaces and names stay the same. Convert an existing individual account
into an organization when more than one person needs access to that
account and its repositories. To choose a new namespace and keep this
Docker ID as an individual account, see
[Create a Docker organization](/manuals/accounts/organization/setup/orgs.md).

## Prerequisites

Before you convert an individual account to an organization:

- The individual account must have a verified email address.
- You need a separate Docker ID to assign as the organization owner.
- The individual account must not belong to an organization, a team, or a
  company. If it does, you must
  [leave the organization](#leave-an-organization) before you convert.

> [!TIP]
>
> After you convert an account, personal access tokens from that account
> stop working. Sign in as the new owner and create an
> [organization access token (OAT)](/manuals/security/access-tokens/organization-access-tokens.md)
> for the converted organization.

## Conversion

Converting an account into an organization does the following:

- The account gets a
  [Docker Team](/manuals/subscription-billing/plans/docker.md)
  subscription. If the individual account was on a paid plan, it keeps
  the same billing cycle.
- Docker signs you out. Your email address is no longer on any Docker
  account. You can [sign up](https://hub.docker.com/signup) with it again
  or add it to an existing account.
- Repository collaborators are removed.
  [Invite them to your organization](/manuals/accounts/organization/manage/members.md),
  then add them to a
  [team](/manuals/accounts/organization/manage/manage-a-team.md)
  that has access to the repositories.
- Existing automated builds appear as if the new owner set them up.

### Choose an organization owner

To convert, you must name an owner. That owner has full administrative access
to configure and manage the organization. You can add more owners after
conversion. Each owner signs in to a separate individual account, then selects the
organization in Docker Home.

To name someone else as the owner, enter that person's Docker ID during
conversion. After conversion, you can sign up again with
`alex@example.com`, then the new owner can invite you to the organization
if you still need access.

If you want to remain the sole owner without assigning another owner, you must create a second individual account with a different email address.
Use your new account's Docker ID.

- For example, if the account you are converting
uses `alex@example.com`, create the owner account with
`alex.admin@example.com`.
- After conversion, `alex@example.com` is no longer on any Docker account.

You can update the `alex.admin@example.com` email back to `alex@example.com` after conversion. See
[Update email address](/manuals/accounts/individual/manage-account.md#update-email-address).

## Convert to an organization

Converting an account into an organization is permanent. Back up any data
or settings you want to keep.

1. Sign in to [Docker Home](https://app.docker.com/).
1. Select your avatar in the top-right corner.
1. Select **Account settings**, then **Convert**.
1. Review the warning. You can't undo this action.
1. Enter the **Username of new owner**.
   Enter a Docker ID, not an email address. Use 4 to 30 letters and
   digits. The Docker ID must belong to a different account from the one
   you are converting. That account must be active and have a verified
   email address.
1. Select **Confirm**. For an account on a Pro plan, the button is
   **Confirm and purchase**.
   The new owner receives a notification email. Sign in with that account
   to manage the organization.

### Leave an organization

Use this procedure when the individual account is already a member of an
organization, a team, or a company. The account can be converted after it
has left each of them. If it has no memberships, skip this section and
[convert the account](#convert-to-an-organization).

If the account is the only owner of an organization or company, assign
the owner role to another member first. Leave after that member can
administer the organization or company.

To leave an organization and its teams:

1. Sign in to [Docker Home](https://app.docker.com/) and select the
   organization.
1. Select **Members** and find your username.
1. Select the **Actions** menu, then **Leave organization**.

## Next steps

- [Onboard your organization](/manuals/accounts/organization/setup/onboard.md)
- [Create an organization access token](/manuals/security/access-tokens/organization-access-tokens.md)
