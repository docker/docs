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

Learn how to manage your company in the Docker Admin Console, including the
organizations within it, company owners, and members.

## Manage organizations

Manage the organizations in a company using the Docker Admin Console.

### View all organizations

1. Sign in to the [Docker Home](https://app.docker.com) and choose
   your company.
1. Select **Managed organizations**.

The **Organizations** view displays all organizations under your company.

### Add seats to an organization

If you have a self-serve subscription that has no pending subscription changes,
you can add seats using Docker Home. For more information about adding seats,
see [Manage seats](/manuals/admin/organization/manage/manage-seats.md#add-seats-to-your-subscription).

If you have a sales-assisted subscription, you must contact Docker support or
sales to add seats.

### Add organizations to a company

To add an organization to a company, ensure the following:

- You are a company owner.
- You are an organization owner of the organization you want to add.
- The organization has a Docker Business subscription.
- There’s no limit to how many organizations can exist under a company.

> [!IMPORTANT]
>
> Once you add an organization to a company, you can't remove it from the
> company.

1. Sign in to [Docker Home](https://app.docker.com) and select your company from
   the top-left account drop-down.
1. Select **Managed organizations**.
1. Select **Add organization**.
1. Choose the organization you want to add from the drop-down menu.
1. Select **Add organization** to confirm.

### Manage an organization

1. Sign in to [Docker Home](https://app.docker.com) and select your company from
   the top-left account drop-down.
1. Select **Managed organizations**.
1. Select the organization you want to manage.

For more details about managing an organization, see
[Organization administration](../organization/_index.md).

## Manage company owners

A company can have multiple owners. Company owners have visibility across the
entire company and can manage settings that apply to all organizations under
that company. They also have the same access rights as organization owners but
don’t need to be members of any individual organization.

> [!IMPORTANT]
>
> Company owners do not occupy a seat unless they are added as a member of an
> organization under your company, or SSO is enabled and the company owner signs
> in via SSO (which automatically adds them as an organization member).

### Add a company owner

1. Sign in to [Docker Home](https://app.docker.com) and select your company from
   the top-left account drop-down.
1. Select **Company owners**.
1. Select **Add owner**.
1. Specify the user's Docker ID to search for the user.
1. After you find the user, select **Add company owner**.

### Remove a company owner

1. Sign in to [Docker Home](https://app.docker.com) and select your company from
   the top-left account drop-down.
1. Select **Company owners**.
1. Locate the company owner you want to remove and select the **Actions** menu.
1. Select **Remove as company owner**.

## Manage company members

You add a user to your company by inviting them to an organization within the
company. Company owners can invite members to any organization in the company
using a Docker ID, email address, or in bulk with a CSV file of email addresses.

Members and invitations belong to individual organizations, not to the company
itself. A pending invitation occupies a seat in the organization the user is
invited to.

### Invite members to an organization

Company owners invite members at the organization level. Select the target
organization first, then follow the standard invitation steps.

1. Sign in to [Docker Home](https://app.docker.com) and select your company
   from the top-left account drop-down.
1. On the **Organizations** page, select the organization you want to add
   members to.
1. Follow the steps in
   [Manage organization members](../organization/manage/members.md#invite-members)
   to invite members using a Docker ID or email address, a CSV file, or the
   Docker Hub API.

### Resend invitations

Company owners can resend invitations across all organizations in the company
from the company-level **Users** page.

#### Resend an individual invitation

1. In [Docker Home](https://app.docker.com/), select your company from the
   top-left account drop-down.
1. Select **Users**.
1. Select the **action menu** next to the invitee, then select **Resend**.
1. Select **Invite** to confirm.

#### Bulk resend invitations

1. In [Docker Home](https://app.docker.com/), select your company from the
   top-left account drop-down.
1. Select **Users**.
1. Use the **checkboxes** next to **Usernames** to select users.
1. Select **Resend invites**.
1. Select **Resend** to confirm.

### Manage members on a team

Teams exist at the organization level, not the company level. After inviting
members to an organization, you can add them to teams within that organization.
For more details, see
[Manage members on a team](../organization/manage/members.md#manage-members-on-a-team).

## More resources

- [Video: Managing a company and nested organizations](https://youtu.be/XZ5_i6qiKho?feature=shared&t=229)
- [Video: Adding nested organizations to a company](https://youtu.be/XZ5_i6qiKho?feature=shared&t=454)
