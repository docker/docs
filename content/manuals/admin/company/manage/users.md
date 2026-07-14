---
title: Manage company members
linkTitle: Users
description: Learn how to invite and manage members in the organizations within your company using the Docker Admin Console.
keywords: company, company members, members, admin, Admin Console, member management, organization management, company management, bulk invite, resend invites
aliases:
  - /admin/company/users/
---

{{< summary-bar feature_name="Company" >}}

You add a user to your company by inviting them to an organization within the
company. Company owners can invite members to any organization in the company
using a Docker ID, email address, or in bulk with a CSV file of email addresses.

Members and invitations belong to individual organizations, not to the company
itself. A pending invitation occupies a seat in the organization the user is
invited to.

## Invite members to an organization

Company owners invite members at the organization level. Select the target
organization first, then follow the standard invitation steps.

1. Sign in to [Docker Home](https://app.docker.com) and select your company
   from the top-left account drop-down.
1. On the **Organizations** page, select the organization you want to add
   members to.
1. Follow the steps in
   [Manage organization members](../../organization/manage/members.md#invite-members)
   to invite members using a Docker ID or email address, a CSV file, or the
   Docker Hub API.

## Resend invitations

Company owners can resend invitations across all organizations in the company
from the company-level **Users** page.

### Resend an individual invitation

1. In [Docker Home](https://app.docker.com/), select your company from the
   top-left account drop-down.
1. Select **Users**.
1. Select the **action menu** next to the invitee, then select **Resend**.
1. Select **Invite** to confirm.

### Bulk resend invitations

1. In [Docker Home](https://app.docker.com/), select your company from the
   top-left account drop-down.
1. Select **Users**.
1. Use the **checkboxes** next to **Usernames** to select users.
1. Select **Resend invites**.
1. Select **Resend** to confirm.

## Manage members on a team

Teams exist at the organization level, not the company level. After inviting
members to an organization, you can add them to teams within that organization.
For more details, see
[Manage members on a team](../../organization/manage/members.md#manage-members-on-a-team).
