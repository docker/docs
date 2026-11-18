---
title: Set up a company
linkTitle: Setup
weight: 10
description: >-
  Set up a Docker company in Docker Home. Create the company, add
  owners, invite members, add organizations, and manage seats.
keywords: company, hub, organization, company owner, Docker Home, company
  management, Docker Business, set up company, create company, invite
  members
toc_max: 2
aliases:
  - /admin/company/new-company/
  - /docker-hub/new-company/
---

{{< summary-bar feature_name="Company" >}}

When you set up a company in Docker Home, you create the company, then add
owners, invite members, and add more organizations as needed. You must be the
owner of the organization you want to add to your company.


## Create a company

When you create the company, the organization you selected is automatically
migrated to the new company. The name for your company can't be the same as an existing user, organization, or company namespace.

To create a new company:

1. Sign in to [Docker Home](https://app.docker.com/) and select the
   organization you want to add to your company.
1. Select **Organization settings**.
1. Select **Create a company**.
1. Enter a unique name for your company, then select **Continue**.
1. Review the migration details and then select **Create company**.

For more information on how you can add more organizations to your company,
see [Add more organizations](#add-more-organizations).

## Company owners

A company can have multiple owners who manage the company and all of its
organizations. A company can have up to 10 owners. For details about the
company owner role and how it affects seats, see
[Company structure](/manuals/accounts/company/_index.md#company-structure).

### Add a company owner

1. Sign in to [Docker Home](https://app.docker.com) and select your company.
1. Select **Company owners**, then select **Add owner**.
1. Specify the member's Docker ID, then select **Add company owner**.

**Add owner** is
unavailable when the company already has 10 owners.

### Remove a company owner

1. Sign in to [Docker Home](https://app.docker.com) and select your company.
1. Select **Company owners**.
1. Find the company owner you want to remove and select the **Actions**
   menu, then select **Remove as company owner**.

## Company invitations

You add a member to your company by inviting them to an organization within
the company. Company owners can invite members to any organization in the
company using a Docker ID, email address, or in bulk with a CSV file of
email addresses.

Members and invitations belong to individual organizations, not to the
company itself. A pending invitation occupies a seat until the invitee
accepts.

### Invite members to an organization

1. Sign in to [Docker Home](https://app.docker.com) and select your company.
1. Select **Users**, then select **Invite**.
1. Add the invitees in one of these ways:
   - **Emails or usernames.** Select **Emails or usernames**. Enter the
     Docker IDs or email addresses of the invitees. Select an organization
     for each invitee.
   - **CSV upload.** Select **CSV upload**. Upload a CSV file of email
     addresses. Select an organization for the invitees.

Invitees receive an email with instructions to accept. After they accept,
new members appear on the **Users** page. The table shows how many
organizations each member belongs to.

### Resend invitations

Company owners can resend invitations from the company-level **Users**
page.

1. Sign in to [Docker Home](https://app.docker.com) and select your company.
1. Select **Users**.
1. Resend one invitation, or resend several:
   - **One invitee.** Locate the invitee in the table. Select the
     **Actions** menu, then select **Resend**. Before resending, confirm
     you selected the correct invitee. The resend invitation modal displays
     the date you originally invited the invitee. Select **Invite** to
     confirm.
   - **Several invitees.** Select the checkboxes next to the invitees.
     Select **Resend invites**, then select **Resend** to confirm.

### Cancel an invitation

1. Sign in to [Docker Home](https://app.docker.com) and select your company.
1. Select **Users**.
1. Find the invitee and select the **Actions** menu, then select **Remove**.
1. Complete the removal in one of these ways:
   - **One organization.** Select **Remove**.
   - **Several organizations.** Select the organizations, then select
     **Remove**.

Removing an invitee cancels the pending invitation.

## Add more organizations

Company owners can add a Docker Business organization when they are also
an owner of that organization. There's no limit to the number of
organizations you add to a company.

1. Sign in to [Docker Home](https://app.docker.com) and select
   your company.
1. Select **Managed organizations**.
1. Select **Add organization**, then select an organization from the
   drop-down list.

A nested organization must keep its Docker Business subscription to stay
managed by the company. If an organization downgrades from Docker
Business, you can no longer manage it through the company, and its owner
must manage it separately.

## Add seats to an organization

If you have a self-serve subscription that has no pending subscription
changes, you can add seats using Docker Home. For more information about
adding seats, see
[Manage seats](/manuals/accounts/organization/manage/manage-seats.md#add-seats-to-your-subscription).

If you have a sales-assisted subscription, you must contact Docker support
or sales to add seats.

## Next steps

- [Manage company members](./manage.md) to change roles, remove users, or
  export a member CSV
- [Manage organization members](/manuals/accounts/organization/manage/members.md)
  for organization-level invites and teams
- [Company overview](/manuals/accounts/company/_index.md) for company
  structure and owner seats
- [Roles and permissions](/manuals/security/roles-and-permissions/_index.md)
  for what each role can do
