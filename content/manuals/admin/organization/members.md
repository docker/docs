---
title: Manage organization members
weight: 30
description: Learn how to manage organization members in Docker Hub and Docker Admin Console.
keywords: members, teams, organizations, invite members, manage team members, export member list, edit roles, organization teams, user management
aliases:
- /docker-hub/members/
---

Learn how to manage members for your organization in Docker Hub and the Docker Admin Console.

## Invite members

{{< tabs >}}
{{< tab name="Admin Console" >}}

Owners can invite new members to an organization via Docker ID, email address, or with a CSV file containing email addresses. If an invitee does not have a Docker account, they must create an account and verify their email address before they can accept an invitation to join the organization. When inviting members, their pending invitation occupies a seat.

### Invite members via Docker ID or email address

Use the following steps to invite members to your organization via Docker ID or email address.

1. Sign in to [Docker Home](https://app.docker.com) and select your organization from the top-left account drop-down.
1. Select **Members**, then **Invite**.
1. Select **Emails or usernames**.
1. Follow the on-screen instructions to invite members. Invite a maximum of 1000 members and separate multiple entries by comma, semicolon, or space.

> [!NOTE]
>
> When you invite members, you assign them a role. See [Roles and permissions](/manuals/enterprise/security/roles-and-permissions.md) for
details about the access permissions for each role.

Pending invitations appear in the table. Invitees receive an email with a link to Docker Hub where they can accept or decline the invitation.

### Invite members via CSV file

To invite multiple members to an organization via a CSV file containing email addresses:

1. Sign in to [Docker Home](https://app.docker.com) and select your organization from the top-left account drop-down.
1. Select **Members**, then **Invite**.
1. Select **CSV upload**.
1. Optional. Select **Download the template CSV file** to download an example CSV file. The following is an example of the contents of a valid CSV file.

```text
email
docker.user-0@example.com
docker.user-1@example.com
```

CSV file requirements:

- The file must contain a header row with at least one heading named email. Additional columns are allowed and are ignored in the import.
- The file must contain a maximum of 1000 email addresses (rows). To invite more than 1000 users, create multiple CSV files and perform all steps in this task for each file.

1. Create a new CSV file or export a CSV file from another application.

- To export a CSV file from another application, see the application’s documentation.
- To create a new CSV file, open a new file in a text editor, type email on the first line, type the user email addresses one per line on the following lines, and then save the file with a .csv extension.

1. Select **Browse files** and then select your CSV file, or drag and drop the CSV file into the **Select a CSV file to upload** box. You can only select one CSV file at a time.

> [!NOTE]
>
> If the amount of email addresses in your CSV file exceeds the number of available seats in your organization, you cannot continue to invite members. To invite members, you can purchase more seats, or remove some email addresses from the CSV file and re-select the new file. To purchase more seats, see [Add seats](/manuals/subscription/manage-seats.md) to your subscription or [Contact sales](https://www.docker.com/pricing/contact-sales/).

1. After the CSV file has been uploaded, select **Review**.

Valid email addresses and any email addresses that have issues appear. Email addresses may have the following issues:

- Invalid email: The email address is not a valid address. The email address will be ignored if you send invites. You can correct the email address in the CSV file and re-import the file.
- Already invited: The user has already been sent an invite email and another invite email will not be sent.
- Member: The user is already a member of your organization and an invite email will not be sent.
- Duplicate: The CSV file has multiple occurrences of the same email address. The user will be sent only one invite email.

1. Follow the on-screen instructions to invite members.

> [!NOTE]
>
> When you invite members, you assign them a role. See [Roles and permissions](/manuals/enterprise/security/roles-and-permissions.md) for
details about the access permissions for each role.

Pending invitations appear in the table. The invitees receive an email with a link to Docker Hub where they can accept or decline the invitation.

### Invite members via API

You can bulk invite members using the Docker Hub API. For more information, see the [Bulk create invites](https://docs.docker.com/reference/api/hub/latest/#tag/invites/paths/~1v2~1invites~1bulk/post) API endpoint.

{{< /tab >}}
{{< tab name="Docker Hub" >}}

{{% include "hub-org-management.md" %}}

Owners can invite new members to an organization via Docker ID, email address, or with a CSV file containing email addresses. If an invitee does not have a Docker account, they must create an account and verify their email address before they can accept an invitation to join the organization. When inviting members, their pending invitation occupies a seat.

### Invite members via Docker ID or email address

Use the following steps to invite members to your organization via Docker ID or email address.

1. Sign in to [Docker Hub](https://hub.docker.com).
1. Select **My Hub**, your organization, then **Members**.
1. Select **Invite members**.
1. Select **Emails or usernames**.
1. Follow the on-screen instructions to invite members. Invite a maximum of 1000 members and separate multiple entries by comma, semicolon, or space.

> [!NOTE]
>
> When you invite members, you assign them a role. See [Roles and permissions](/manuals/enterprise/security/roles-and-permissions.md) for
details about the access permissions for each role.

Pending invitations appear in the table. The invitees receive an email with a link to Docker Hub where they can accept or decline the invitation.

### Invite members via CSV file

To invite multiple members to an organization via a CSV file containing email addresses:

1. Sign in to [Docker Hub](https://hub.docker.com).
1. Select **My Hub**, your organization, then **Members**.
1. Select **Invite members**.
1. Select **CSV upload**.
1. Optional. Select **Download the template CSV file** to download an example CSV file. The following is an example of the contents of a valid CSV file.

```text
email
docker.user-0@example.com
docker.user-1@example.com
```

CSV file requirements:

- The file must contain a header row with at least one heading named email. Additional columns are allowed and are ignored in the import.
- The file must contain a maximum of 1000 email addresses (rows). To invite more than 1000 users, create multiple CSV files and perform all steps in this task for each file.

1. Create a new CSV file or export a CSV file from another application.

- To export a CSV file from another application, see the application’s documentation.
- To create a new CSV file, open a new file in a text editor, type email on the first line, type the user email addresses one per line on the following lines, and then save the file with a .csv extension.

1. Select **Browse files** and then select your CSV file, or drag and drop the CSV file into the **Select a CSV file to upload** box. You can only select one CSV file at a time.

> [!NOTE]
>
> If the amount of email addresses in your CSV file exceeds the number of available seats in your organization, you cannot continue to invite members. To invite members, you can purchase more seats, or remove some email addresses from the CSV file and re-select the new file. To purchase more seats, see [Add seats](/manuals/subscription/manage-seats.md) to your subscription or [Contact sales](https://www.docker.com/pricing/contact-sales/).

1. After the CSV file has been uploaded, select **Review**.

Valid email addresses and any email addresses that have issues appear. Email addresses may have the following issues:

- Invalid email: The email address is not a valid address. The email address will be ignored if you send invites. You can correct the email address in the CSV file and re-import the file.
- Already invited: The user has already been sent an invite email and another invite email will not be sent.
- Member: The user is already a member of your organization and an invite email will not be sent.
- Duplicate: The CSV file has multiple occurrences of the same email address. The user will be sent only one invite email.

1. Follow the on-screen instructions to invite members.

> [!NOTE]
>
> When you invite members, you assign them a role. See [Roles and permissions](/manuals/enterprise/security/roles-and-permissions.md) for
details about the access permissions for each role.

Pending invitations appear in the table. The invitees receive an email with a link to Docker Hub where they can accept or decline the invitation.

### Invite members via API

You can bulk invite members using the Docker Hub API. For more information, see the [Bulk create invites](https://docs.docker.com/reference/api/hub/latest/#tag/invites/paths/~1v2~1invites~1bulk/post) API endpoint.

{{< /tab >}}
{{< /tabs >}}

## Accept invitation

When an invitation is to a user's email address, they receive
a link to Docker Hub where they can accept or decline the invitation.
To accept an invitation:

1. Check your email inbox and open the Docker email with an invitation to
join the Docker organization.
1. To open the link to Docker Hub, select the **click here** link.

   > [!WARNING]
   >
   > Invitation email links expire after 14 days. If your email link has expired,
   > you can sign in to [Docker Hub](https://hub.docker.com/) with the email
   > address the link was sent to and accept the invitation from the
   > **Notifications** panel.

1. The Docker create an account page will open. If you already have an account, select **Already have an account? Sign in**.
If you do not have an account yet, create an account using the same email
address you received the invitation through.
1. Optional. If you do not have an account and created one, you must navigate
back to your email inbox and verify your email address using the Docker verification
email.
1. Once you are signed in to Docker Hub, select **My Hub** from the top-level navigation menu.
1. Select **Accept** on your invitation.

After accepting an invitation, you are now a member of the organization.

## Manage invitations

After inviting members, you can resend or remove invitations as needed.

### Resend an invitation

{{< tabs >}}
{{< tab name="Admin Console" >}}

You can send individual invitations, or bulk invitations from the Admin Console.

To resend an individual invitation:

1. Sign in to [Docker Home](https://app.docker.com/) and select your
organization.
1. Select **Members**.
1. Select the **action menu** next to the invitee and select **Resend**.
1. Select **Invite** to confirm.

To bulk resend invitations:

1. Sign in to [Docker Home](https://app.docker.com/) and select your
organization.
1. Select **Members**.
1. Use the **checkboxes** next to **Usernames** to bulk select users.
1. Select **Resend invites**.
1. Select **Resend** to confirm.

{{< /tab >}}
{{< tab name="Docker Hub" >}}

{{% include "hub-org-management.md" %}}

To resend an invitation from Docker Hub:

1. Sign in to [Docker Hub](https://hub.docker.com/).
1. Select **My Hub**, your organization, and then **Members**.
1. In the table, locate the invitee, select the **Actions** icon, and then select
**Resend invitation**.
1. Select **Invite** to confirm.

You can also resend an invitation using the Docker Hub API. For more information,
see the [Resend an invite](https://docs.docker.com/reference/api/hub/latest/#tag/invites/paths/~1v2~1invites~1%7Bid%7D~1resend/patch) API endpoint.

{{< /tab >}}
{{< /tabs >}}

### Remove an invitation

{{< tabs >}}
{{< tab name="Admin Console" >}}

To remove an invitation from the Admin Console:

1. Sign in to [Docker Home](https://app.docker.com/) and select your
organization.
1. Select **Members**.
1. Select the **action menu** next to the invitee and select **Remove invitee**.
1. Select **Remove** to confirm.

{{< /tab >}}
{{< tab name="Docker Hub" >}}

{{% include "hub-org-management.md" %}}

To remove a member's invitation from Docker Hub:

1. Sign in to [Docker Hub](https://hub.docker.com/).
1. Select **My Hub**, your organization, and then **Members**.
1. In the table, select the **Action** icon, and then select **Remove member** or **Remove invitee**.
1. Follow the on-screen instructions to remove the member or invitee.

You can also remove an invitation using the Docker Hub API. For more information,
see the [Cancel an invite](https://docs.docker.com/reference/api/hub/latest/#tag/invites/paths/~1v2~1invites~1%7Bid%7D/delete) API endpoint.

{{< /tab >}}
{{< /tabs >}}

## Manage members on a team

Use Docker Hub or the Admin Console to add or remove team members. Organization owners can add a member to one or more teams within an organization.

### Add a member to a team

{{< tabs >}}
{{< tab name="Admin Console" >}}

To add a member to a team with the Admin Console:

1. Sign in to [Docker Home](https://app.docker.com/) and select your
organization.
1. Select **Teams**.
1. Select the team name.
1. Select **Add member**. You can add the member by searching for their email address or username.

   > [!NOTE]
   >
   > An invitee must first accept the invitation to join the organization before being added to the team.

{{< /tab >}}
{{< tab name="Docker Hub" >}}

{{% include "hub-org-management.md" %}}

To add a member to a team with Docker Hub:

1. Sign in to [Docker Hub](https://hub.docker.com).
1. Select **My Hub**, your organization, and then **Members**.
1. Select the **Action** icon, and then select **Add to team**.

   > [!NOTE]
   >
   > You can also navigate to **My Hub** > **Your Organization** > **Teams** > **Your Team Name** and select **Add Member**. Select a member from the drop-down list to add them to the team or search by Docker ID or email.
1. Select the team and then select **Add**.

   > [!NOTE]
   >
   > An invitee must first accept the invitation to join the organization before being added to the team.

{{< /tab >}}
{{< /tabs >}}

### Remove members from teams

> [!NOTE]
>
> If your organization uses single sign-on (SSO) with [SCIM](/manuals/enterprise/security/provisioning/scim.md) enabled, you should remove members from your identity provider (IdP). This will automatically remove members from Docker. If SCIM is disabled, you must manually manage members in Docker.

Organization owners can remove a member from a team in Docker Hub or Admin Console. Removing the member from the team will revoke their access to the permitted resources.

{{< tabs >}}
{{< tab name="Admin Console" >}}

To remove a member from a specific team with the Admin Console:

1. Sign in to [Docker Home](https://app.docker.com/) and select your
organization.
1. Select **Teams**.
1. Select the team name.
1. Select the **X** next to the user's name to remove them from the team.
1. When prompted, select **Remove** to confirm.

{{< /tab >}}
{{< tab name="Docker Hub" >}}

{{% include "hub-org-management.md" %}}

To remove a member from a specific team with Docker Hub:

1. Sign in to [Docker Hub](https://hub.docker.com).
1. Select **My Hub**, your organization, **Teams**, and then the team.
1. Select the **X** next to the user’s name to remove them from the team.
1. When prompted, select **Remove** to confirm.

{{< /tab >}}
{{< /tabs >}}

### Update a member role

Organization owners can manage [roles](/security/for-admins/roles-and-permissions/)
within an organization. If an organization is part of a company,
the company owner can also manage that organization's roles. If you have SSO enabled, you can use [SCIM for role mapping](/security/for-admins/provisioning/scim/).

> [!NOTE]
>
> If you're the only owner of an organization, you need to assign a new owner
before you can edit your role.

{{< tabs >}}
{{< tab name="Admin Console" >}}

To update a member role in the Admin Console:

1. Sign in to [Docker Home](https://app.docker.com/) and select your
organization.
1. Select **Members**.
1. Find the username of the member whose role you want to edit. Select the
**Actions** menu, then **Edit role**.

{{< /tab >}}
{{< tab name="Docker Hub" >}}

{{% include "hub-org-management.md" %}}

To update a member role in Docker Hub:

1. Sign in to [Docker Hub](https://hub.docker.com).
1. Select **My Hub**, your organization, and then **Members**.
1. Find the username of the member whose role you want to edit. In the table, select the **Actions** icon.
1. Select **Edit role**.
1. Select their organization, select the role you want to assign, and then select **Save**.

> [!NOTE]
>
> If you're the only owner of an organization,
> you need to assign a new owner before you can edit your role.

{{< /tab >}}
{{< /tabs >}}

## Export members CSV file

{{< summary-bar feature_name="Admin orgs" >}}

Owners can export a CSV file containing all members. The CSV file for a company contains the following fields:

- Name: The user's name
- Username: The user's Docker ID
- Email: The user's email address
- Member of Organizations: All organizations the user is a member of within a company
- Invited to Organizations: All organizations the user is an invitee of within a company
- Account Created: The time and date when the user account was created

{{< tabs >}}
{{< tab name="Admin Console" >}}

To export a CSV file of your members:

1. Sign in to [Docker Home](https://app.docker.com/) and select your
organization.
1. Select **Members**.
1. Select the **download** icon to export a CSV file of all members.

{{< /tab >}}
{{< tab name="Docker Hub" >}}

{{% include "hub-org-management.md" %}}

To export a CSV file of your members:

1. Sign in to [Docker Hub](https://hub.docker.com).
1. Select **My Hub**, your organization, and then **Members**.
1. Select the **Action** icon and then select **Export users as CSV**.

{{< /tab >}}
{{< /tabs >}}
