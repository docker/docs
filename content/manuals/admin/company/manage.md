---
title: Manage company organizations, owners, and members
linkTitle: Manage
weight: 20
description: Manage your Docker company in Docker Home. Add organizations, invite members, manage owners, resend invitations, export a member CSV, and add subscription seats.
keywords: company, manage company, Docker Home, company owners, add organization, invite members, resend invitations, export members CSV, company members, manage seats
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

After creating a company, you can manage multiple organizations from
Docker Home. Company owners can use the company portal to invite members to
specific organizations, view seat availability across organizations, and
add new company owners.

## Add more organizations

Company owners can add Docker organizations with a Docker Business plan
to their company, so long as they're also the organization owners for
that organization. There's no limit to the number of organizations you
add to a company.

> [!IMPORTANT]
>
> Once you add an organization to a company, you can't remove it from the
> company.

1. Sign in to [Docker Home](https://app.docker.com) and select
   your company.
1. Select **Managed organizations**.
1. Select **Add organization**, then select an organization from the
   drop-down list.

A nested organization must keep its Docker Business subscription to stay
managed by the company. If an organization downgrades from Docker
Business, you can no longer manage it through the company, and its owner
must manage it separately.

## Company owners

A company can have multiple owners who manage the company and all of its
organizations. For details about the company owner role and how it affects
seats, see [Company roles](/manuals/admin/company/_index.md#company-roles).

### Add a company owner

1. Sign in to [Docker Home](https://app.docker.com) and select your company.
1. Select **Company owners**, then select **Add owner**.
1. Specify the member's Docker ID, then select **Add company owner**.

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

{{< tabs >}}
{{< tab name="Email or username" >}}

1. Select **Emails or usernames**.
1. Enter the Docker IDs or email addresses of the invitees.
1. Select an organization for each invitee.

{{< /tab >}}
{{< tab name="CSV upload" >}}

1. Select **CSV upload**.
1. Upload a CSV file of email addresses.
1. Select an organization for the invitees.

{{< /tab >}}
{{< /tabs >}}

Invitees receive an email with instructions to accept. After they accept,
new members appear on the **Users** page. The table shows how many
organizations each member belongs to.

### Resend invitations

Company owners can resend invitations from the company-level **Users**
page.

{{< tabs >}}
{{< tab name="Individual" >}}

1. Select your company from [Docker Home](https://app.docker.com/).
1. Select **Users**, then locate the invitee from the table.
1. Select the **Actions** menu, then select **Resend**.
   - Before resending, confirm you selected the correct invitee.
   - The resend invitation modal displays the date you originally invited
     the invitee.
1. Select **Invite** to confirm.

{{< /tab >}}
{{< tab name="Bulk" >}}

1. Select your company from [Docker Home](https://app.docker.com/).
1. Select **Users**.
1. From the table, use the multi-select checkboxes next to the invitees
   you want to resend to.
1. Select **Resend invites**, then select **Resend** to confirm.

{{< /tab >}}
{{< /tabs >}}

## Export a member list CSV

Company owners can export a CSV file of members across organizations in the
company.

1. Sign in to [Docker Home](https://app.docker.com/) and select your company.
1. Select **Users**.
1. Select the **Download** icon. The CSV file downloads in your browser.

   {{< accordion title="CSV fields" >}}
   - Name: The member's name
   - Username: The member's Docker ID
   - Email: The member's email address
   - Member of Organizations: Organizations the member belongs to within the
     company
   - Invited to Organizations: Organizations where the invitee has a
     pending invitation
   - Account created: The time and date when the account was created

   {{< /accordion >}}

## Add seats to an organization

If you have a self-serve subscription that has no pending subscription
changes, you can add seats using Docker Home. For more information about
adding seats, see
[Manage seats](/manuals/admin/organization/manage/manage-seats.md).

If you have a sales-assisted subscription, you must contact Docker support
or sales to add seats.

## Manage teams

Teams exist at the organization level, not the company level. After inviting
members to an organization, you can add them to teams within that
organization. For more details, see
[Manage members on a team](/manuals/admin/organization/manage/members.md).

## Next steps

- [Company overview](/manuals/admin/company/_index.md)
- [Manage organization members](/manuals/admin/organization/manage/members.md)
- [Manage seats](/manuals/admin/organization/manage/manage-seats.md)
- [Roles and permissions](/manuals/enterprise/security/roles-and-permissions/_index.md)
