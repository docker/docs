---
title: Manage your company
linkTitle: Manage
weight: 20
description: Learn how to manage your company, including its organizations, owners, and members, using the Docker Admin Console.
keywords: company, manage company, multiple organizations, company owners, company members, Docker Admin Console, add organization, resend invites
aliases:
  - /admin/company/manage/organizations/
  - /admin/company/manage/owners/
  - /admin/company/manage/users/
  - /admin/company/organizations/
  - /admin/company/owners/
  - /admin/company/users/
  - /docker-hub/company-owner/
---

{{< summary-bar feature_name="Company" >}}

After creating a company, you can manage multiple organizations from Docker Home. Company owners can use the company portal to invite users to specific organizations, view seat availability across organizations, and add new company owners.

## Add more organizations

Company owners can add Docker organizations with a Docker Business plan to their company, so long as they're also the organization owners for that organization. There's no limit to the number of organizations you add to a company.

> [!IMPORTANT]
>
> Once you add an organization to a company, you can't remove it from the
> company.

1. Sign in to the [Docker Home](https://app.docker.com) and select
   your company.
1. Select **Managed organizations**.
1. Select **Add organization**, then choose an organization from the dropdown.

A nested organization must keep its Docker Business subscription to stay managed
by the company. If an organization downgrades from Docker Business, you can no
longer manage it through the company, and its owner must manage it separately.

## Company owners

A company can have multiple owners. Company owners have visibility across the
entire company and can manage settings that apply to all organizations under
that company. They also have the same access rights as organization owners without
needing individual organization membership.

Company owners do not occupy a seat unless:

- They are added as a member of an organization under your company.
- SSO is enabled and the company owner signs in via SSO, which automatically
  adds them as an organization member.

When you first create a company, your account is both a company owner and an
organization owner, so it occupies a seat as long as you remain an organization
owner. To keep full company-owner access without using a seat,
[assign another user as the organization owner](/manuals/admin/organization/manage/members.md#update-a-member-role),
then remove yourself from the organization.

### Add a company owner

1. Sign in to [Docker Home](https://app.docker.com) and select your company.
1. Select **Company owners**, then choose **Add owner**.
1. Specify the user's Docker ID, then finish by selecting **Add company owner**.

### Remove company owner

1. Sign in to [Docker Home](https://app.docker.com) and select your company.
1. Select **Company owners**.
1. Find the company owner you want to remove and select the **Actions** menu, then choose **Remove as company owner**.

## Company invitations

You add a user to your company by inviting them to an organization within the
company. Company owners can invite members to any organization in the company
using a Docker ID, email address, or in bulk with a CSV file of email addresses.

Members and invitations belong to individual organizations, not to the company
itself. A pending invitation occupies a seat in the organization the user is
invited to.

### Invite members to an organization

1. Sign in to [Docker Home](https://app.docker.com) and select your company.
1. Select **Users**, then choose **Invite**.
1. Choose how you want to invite members:
   - To invite individual users, select **Emails or usernames**.
   - To invite groups of users, select **CSV upload**.
1. Add user(s) to an organization by choosing **Select an organization**.

Users receive invitations in their email with instructions to accept the invitation. After accepting the invitation, new members appear on the **Users** page. The table specifies how many organizations they're members of.

### Resend invitations

Company owners can resend invitations
from the company-level **Users** page. To resend individual invitations:

1. Select your company from [Docker Home](https://app.docker.com/).
1. Select **Users**, then locate the invitee from the users table.
1. Select the **Actions** menu, then choose **Resend**.
   - Before resending, confirm you are resending the invitation to the correct invitee.
   - The resend invitation modal displays the date you originally invited the invitee.
1. Choose **Invite** to confirm.

To bulk resend invitations:

1. From the users table, use the multi-select checkboxes next to the invitees you want to invite.
1. Select **Resend invites**, then choose **Resend** to confirm.

## Add seats to an organization

If you have a self-serve subscription that has no pending subscription changes,
you can add seats using Docker Home. For more information about adding seats,
see [Manage seats](/manuals/admin/organization/manage/manage-seats.md#add-seats-to-your-subscription).

If you have a sales-assisted subscription, you must contact Docker support or
sales to add seats.

## Manage teams

Teams exist at the organization level, not the company level. After inviting
members to an organization, you can add them to teams within that organization.
For more details, see
[Manage members on a team](../organization/manage/members.md#manage-members-on-a-team).
