---
title: Invite and manage organization members
linkTitle: Members
weight: 10
description: Invite and manage organization members in Docker Home. Assign
  roles, resend or remove invitations, add members to teams, and export a CSV
  member list.
keywords: members, invite members, organization members, Docker Home, Docker
  Hub, export member list, edit roles, manage invitations, CSV invite, bulk
  invite, resend invitation, remove member, accept invitation, teams, pending
  invitations, organization owner
aliases:
  - /docker-hub/members/
  - /admin/organization/members/
---

Learn how to invite and manage members for your organization in Docker Home.

## Invite members

Owners can invite new members using a Docker ID, email address, or a CSV file
of email addresses. If an invitee doesn't have a Docker account, they must
create one and verify their email address before they can accept. Pending
invitations occupy a seat.

When you invite members, you assign them a role. See
[Roles and permissions][roles-permissions] for details about the access
permissions for each role.

{{< tabs >}}
{{< tab name="Email or username" >}}

1. Sign in to [Docker Home](https://app.docker.com) and select your
   organization from the top-left account drop-down.
1. Select **Members**, then **Invite**.
1. Select **Emails or usernames**.
1. Follow the on-screen instructions to invite members. Invite a maximum of
   1000 members and separate multiple entries by comma, semicolon, or space.

{{< /tab >}}
{{< tab name="CSV upload" >}}

1. Sign in to [Docker Home](https://app.docker.com) and select your
   organization from the top-left account drop-down. Select **Members** >
   **Invite** > **CSV upload**.
1. Optional. Select **Download the template CSV file** to download an example
   CSV file. The following is an example of the contents of a valid CSV file:

   ```text
   email
   docker.user-0@example.com
   docker.user-1@example.com
   ```

   CSV file requirements:
   - The file must contain a header row with at least one heading named
     `email`. Additional columns are allowed and are ignored in the import.
   - The file must contain a maximum of 1000 email addresses (rows). To invite
     more than 1000 members, create multiple CSV files and complete this
     procedure for each file.

1. Create a new CSV file or export a CSV file from another application.
   - To export a CSV file from another application, see that application's
     documentation.
   - To create a new CSV file, open a new file in a text editor, type `email`
     on the first line, type one email address per line on the following
     lines, then save the file with a `.csv` extension.

1. Select **Browse files** and select your CSV file, or drag and drop the CSV
   file into the **Select a CSV file to upload** box. You can select only one
   CSV file at a time.
1. After the CSV file uploads, select **Review** to identify invalid email
   addresses, invitees with a pending invitation, members already in the
   organization, or duplicate email addresses in the same CSV file.
1. Follow the on-screen instructions to invite members.

{{< /tab >}}
{{< /tabs >}}

You can also bulk invite members with the Docker Hub API. For more
information, see the [Bulk create invites][bulk-invites] API endpoint.

Pending invitations appear in the Members table. Invitees receive an email
with a link to Docker Hub to accept or decline.

## Accept an invitation

After you receive an email invitation, open the link to Docker Hub to accept
or decline it.

1. Open the Docker invitation email and select the link to Docker Hub.
1. The Docker account creation page opens. If you already have an account,
   select **Already have an account? Sign in**. If you don't have an account,
   create one using the same email address that received the invitation.
1. Optional. If you created a new account, open your email inbox and verify
   your email address using the Docker verification email.
1. After you sign in to Docker Hub, select **My Hub** from the top-level
   navigation menu.
1. Select **Accept** on your invitation.

After you accept the invitation, you are a member of the organization.

Invitation email links expire after 14 days. If your link has expired, sign in
to [Docker Hub](https://hub.docker.com/) with the email address the invitation
was sent to, then accept the invitation from the **Notifications** panel.

## Manage invitations

After inviting members, you can resend or remove invitations. Each invitee
occupies one seat. If the number of email addresses in your CSV file exceeds
the number of available seats, you can't invite more members.

> [!TIP]
>
> Need more seats for your organization?
> [Add seats](/manuals/admin/organization/manage/manage-seats.md) to your
> subscription, or see [Docker pricing][docker-pricing] for plan options.

### Resend an invitation

You can resend individual or bulk invitations from Docker Home.

To resend an individual invitation:

1. Sign in to [Docker Home](https://app.docker.com/) and select your
   organization.
1. Select **Members**.
1. Select the **Actions** menu next to the invitee, then **Resend**.
1. Select **Invite** to confirm.

To bulk resend invitations:

1. Sign in to [Docker Home](https://app.docker.com/) and select your
   organization.
1. Select **Members**.
1. Use the checkboxes next to **Usernames** to select invitees.
1. Select **Resend invites**.
1. Select **Resend** to confirm.

### Remove an invitation

1. Sign in to [Docker Home](https://app.docker.com/) and select your
   organization.
1. Select **Members**.
1. Select the **Actions** menu next to the invitee, then **Remove invitee**.
1. Select **Remove** to confirm.

## Manage members on a team

Use Docker Hub or Docker Home to add or remove team members. Organization
owners can add a member to one or more teams within an organization.

### Add a member to a team

1. Sign in to [Docker Home](https://app.docker.com/) and select your
   organization.
1. Select **Teams**.
1. Select the team name.
1. Select **Add member**. Search for the member by email address or username.

An invitee must accept the invitation before you can add them to a team.

### Remove members from teams

If your organization uses single sign-on (SSO) with
[SCIM](/manuals/enterprise/security/provisioning/scim/_index.md) enabled,
remove members from your identity provider (IdP). That removes them from
Docker automatically. If SCIM is disabled, remove members manually in Docker
using the following steps.

Removing a member from a team revokes their access to that team's permitted
resources.

1. Sign in to [Docker Home](https://app.docker.com/) and select your
   organization.
1. Select **Teams**, then select the team.
1. Select the **X** next to the member's name to remove them from the team.
1. When prompted, select **Remove** to confirm.

### Update a member role

Organization owners can manage
[roles](/manuals/enterprise/security/roles-and-permissions/_index.md) within
an organization. If an organization is part of a company, the company owner
can also manage that organization's roles. If SSO is enabled, you can use
[SCIM for role mapping][scim-role-mapping].

1. Sign in to [Docker Home](https://app.docker.com/) and select your
   organization.
1. Select **Members**.
1. Find the username of the member whose role you want to edit. Select the
   **Actions** menu, then **Edit role**.

If you're the only owner of an organization and you want to edit your role,
assign a new owner first so you can change your own role.

### Remove a member from the organization

Organization owners can remove members from the organization. Removing a
member revokes their access to the organization's resources and teams.

If your organization uses SSO with
[SCIM](/manuals/enterprise/security/provisioning/scim/_index.md) enabled,
remove members from your IdP instead.

1. Sign in to [Docker Home](https://app.docker.com/) and select your
   organization.
1. Select **Members**.
1. Find the username of the member you want to remove. Select the
   **Actions** menu, then **Remove member**.
1. Select **Remove** to confirm.

## Export a member list CSV

Organization owners can export a CSV file of all members. Docker generates
the file asynchronously and emails it to the owner when it's ready.

1. Sign in to [Docker Home](https://app.docker.com/) and select your
   organization.
1. Select **Members**.
1. Select the **Download** icon to start the export.
1. Open the email from Docker and select the link to download the CSV file.

   {{< accordion title="CSV fields" >}}
   - Name: The member's name
   - Username: The member's Docker ID
   - Email: The member's email address
   - Type: Whether the entry is a member or an invitee
   - Role: The member's role in the organization
   - Teams: Teams the member belongs to
   - Date Joined: When the member joined the organization

   {{< /accordion >}}

## Next steps

After you invite and manage members, explore these related topics:

- [Manage subscription seats](./manage-seats.md) to add seats for pending
  invitations
- [Manage license assignment](./manage-licenses.md) to control product access
- [Create and manage a team](./manage-a-team.md) to group members and set
  repository permissions
- [Roles and permissions][roles-permissions] for role definitions
- [SCIM provisioning][scim-provisioning] to automate member and role
  management

[roles-permissions]: /manuals/enterprise/security/roles-and-permissions/_index.md
[bulk-invites]: /reference/api/hub/latest/#tag/invites/paths/~1v2~1invites~1bulk/post
[docker-pricing]: https://www.docker.com/pricing?ref=Docs&refAction=DocsAdminMembers
[scim-role-mapping]: /manuals/enterprise/security/provisioning/scim/_index.md
[scim-provisioning]: /manuals/enterprise/security/provisioning/scim/_index.md
