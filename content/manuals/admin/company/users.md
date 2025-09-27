---
title: Manage company members
description: Learn how to manage company users in the Docker Admin Console.
keywords: company, company users, users, admin, Admin Console, member management, organization management, company management, bulk invite, resend invites
---

{{< summary-bar feature_name="Company" >}}

Company owners can invite new members to an organization via Docker ID,
email address, or in bulk with a CSV file containing email
addresses.

If an invitee does not have a Docker account, they must create an account and
verify their email address before they can accept an invitation to join the
organization. Pending invitations occupy seats for the organization
the user is invited to.

## Invite members via Docker ID or email address

Use the following steps to invite members to your organization via Docker ID or
email address.

1. Sign in to [Docker Home](https://app.docker.com) and select
your company.
1. On the **Organizations** page, select the organization you want
to invite members to.
1. Select **Members**, then **Invite**.
1. Select **Emails or usernames**.
1. Follow the on-screen instructions to invite members.
   Invite a maximum of 1000 members and separate multiple entries by comma,
   semicolon, or space.

   > [!NOTE]
   >
   > When you invite members, you assign them a role.
   > See [Roles and permissions](/security/for-admins/roles-and-permissions/)
   > for details about the access permissions for each role.

   Pending invitations appear on the Members page. The invitees receive an
   email with a link to Docker Hub where they can accept or decline the
   invitation.

## Invite members via CSV file

To invite multiple members to an organization via a CSV file containing email
addresses:

1. Sign in to [Docker Home](https://app.docker.com) and select
your company.
1. On the **Organizations** page, select the organization you want
to invite members to.
1. Select **Members**, then **Invite**.
1. Select **CSV upload**.
1. Select **Download the template CSV file** to optionally download an example
CSV file. The following is an example of the contents of a valid CSV file.

   ```text
   email
   docker.user-0@example.com
   docker.user-1@example.com
   ```

   CSV file requirements:

   - The file must contain a header row with at least one heading named `email`.
   Additional columns are allowed and are ignored in the import.
   - The file must contain a maximum of 1000 email addresses (rows). To invite
   more than 1000 users, create multiple CSV files and perform all steps in
   this task for each file.

1. Create a new CSV file or export a CSV file from another application.

   - To export a CSV file from another application, see the application’s
   documentation.
   - To create a new CSV file, open a new file in a text editor, type `email`
   on the first line, type the user email addresses one per line on the
   following lines, and then save the file with a .csv extension.

1. Select **Browse files** and then select your CSV file, or drag and drop the
CSV file into the **Select a CSV file to upload** box. You can only select
one CSV file at a time.

   > [!NOTE]
   >
   > If the amount of email addresses in your CSV file exceeds the number of
   available seats in your organization, you cannot continue to invite members.
   To invite members, you can purchase more seats, or remove some email
   addresses from the CSV file and re-select the new file. To purchase more
   seats, see [Add seats to your subscription](/subscription/add-seats/) or
   [Contact sales](https://www.docker.com/pricing/contact-sales/).

1. After the CSV file has been uploaded, select **Review**.

   Valid email addresses and any email addresses that have issues will appear.
   Email addresses may have the following issues:

   - Invalid email: The email address is not a valid address. The email address
   will be ignored if you send invites. You can correct the email address in
   the CSV file and re-import the file.
   - Already invited: The user has already been sent an invite email and another
   invite email will not be sent.
   - Member: The user is already a member of your organization and an invite
   email will not be sent.
   - Duplicate: The CSV file has multiple occurrences of the same email address.
   The user will be sent only one invite email.

1. Follow the on-screen instructions to invite members.

   > [!NOTE]
   >
   > When you invite members, you assign them a role.
   > See [Roles and permissions](/security/for-admins/roles-and-permissions/)
   > for details about the access permissions for each role.

Pending invitations appear on the Members page. The invitees receive an email
with a link to Docker Hub where they can accept or decline the invitation.

## Resend invitations to users

You can resend individual invitations, or bulk invitations from the Admin Console.

### Resend individual invitations

1. In [Docker Home](https://app.docker.com/), select your company from
the top-left account drop-down.
2. Select **Admin Console**, then **Users**.
3. Select the **action menu** next to the invitee and select **Resend**.
4. Select **Invite** to confirm.

### Bulk resend invitation

1. In [Docker Home](https://app.docker.com/), select your company from
the top-left account drop-down.
2. Select **Admin Console**, then **Users**.
3. Use the **checkboxes** next to **Usernames** to bulk select users.
4. Select **Resend invites**.
5. Select **Resend** to confirm.

## Invite members via API

You can bulk invite members using the Docker Hub API. For more information,
see the [Bulk create invites](https://docs.docker.com/reference/api/hub/latest/#tag/invites/paths/~1v2~1invites~1bulk/post) API endpoint.

## Manage members on a team

Use Docker Hub to add a member to a team or remove a member from a team. For
more details, see [Manage members](../organization/members.md#manage-members-on-a-team).
