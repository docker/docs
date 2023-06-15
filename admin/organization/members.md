---
description: Manage organization members
keywords: members, teams, organizations
title: Manage members
---

{% include admin-early-access.md %}

## Invite members

Organization owners can invite new members to an organization via Docker ID, email address, or via a CSV file containing email addresses. If an invitee does not have a Docker account, they must create an account and verify their email address before they can accept the invitation to join the organization. When inviting members, their pending invitation occupies a seat.

### Invite members via Docker ID or email address

Use the following steps to invite members to your organization via Docker ID or email address. To invite a large amount of members to your organization, the recommended method is to [invite members via CSV file](#invite-members-via-csv-file).

1. Sign in to [Docker Admin](https://admin.docker.com){: target="_blank" rel="noopener" class="_"}.
2. In the left navigation, select your organization in the drop-down menu.
3. Select **Members**.
4. Select **Invite Member**.
5. Select **Emails Or Docker IDs**.
6. Enter the Docker IDs or email addresses that you want to invite, up to a maximum of 1000. Separate multiple entries by a comma, semicolon, or space.
7. Select a team or type to create a new team. Docker will invite all users to that team.
8. Select **Invite** to confirm.
   > **Note**
   >
   > You can view the pending invitations in the **Members** page. The invitees receive an email with a link to the organization in Docker Hub where they can accept or decline the invitation.

### Invite members via CSV file

To invite multiple members to your organization via a CSV file containing email addresses:

1. Sign in to [Docker Admin](https://admin.docker.com){: target="_blank" rel="noopener" class="_"}.
2. In the left navigation, select your organization in the drop-down menu.
3. Select **Members**.
4. Select **Invite Member**.
5. Select **CSV Upload**.
6. Select a team or type to create a new team. Docker will invite all users to that team.
7. Select **Download the template CSV file** to optionally download an example CSV file. The following is an example of the contents of a valid CSV file.
    ```
    email
    docker.user-0@example.com
    docker.user-1@example.com
    ```
  CSV file requirements:
   -  The file must contain a header row with at least one heading named `email`. Additional columns are allowed but are ignored in the import.
   -  The file must contain a maximum of 1000 email addresses (rows). To invite more than 1000 users, create multiple CSV files and perform all steps in this task for each file.
8. Create a new CSV file or export a CSV file from another application.
  - To export a CSV file from another application, see the application’s documentation.
  - To create a new CSV file, open a new file in a text editor, type `email` on the first line, type the user email addresses one per line on the following lines, and then save the file with a .csv extension.
9. Select **Browse files** and then select your CSV file, or drag and drop the CSV file into the **Select a CSV file to upload** box. You can only select one CSV file at a time.
  > **Note**
  >
  > If the amount of email addresses in your CSV file exceeds the number of available seats in your organization, you can't continue to invite members. To invite members, you can buy more seats, or remove some email addresses from the CSV file and re-select the new file. To buy more seats, see [Add seats to your subscription](../../subscription/add-seats.md) or [Contact sales](https://www.docker.com/pricing/contact-sales/).
10. After the CSV file has been uploaded, select **Review**.
  Valid email addresses and any email addresses that have issues appear.
  Email addresses may have the following issues:
	  - **Invalid email**: The email address is not a valid address. The email address will be ignored if you send invites. You can correct the email address in the CSV file and re-import the file.
	  - **Already invited**: The user has already been sent an invite email and another invite email will not be sent.
	  - **Member**: The user is already a member of your organization and an invite email will not be sent.
	  - **Duplicate**: The CSV file has multiple occurrences of the same email address. The user will be sent only one invite email.
11. Select **Send invites**.
   > **Note**
   >
   > You can view the pending invitations in the **Members** page. The invitees receive an email with a link to the organization in Docker Hub where they can accept or decline the invitation.

## Add a member to a team

Use Docker Hub to add a member to a team. For more details, see [Add a member to a team](../../docker-hub/members.md#add-a-member-to-a-team).

## Resend invitations

To resend an invitation if the invite is pending or declined:


1. Sign in to [Docker Admin](https://admin.docker.com){: target="_blank" rel="noopener" class="_"}.
2. In the left navigation, select your organization in the drop-down menu.
3. Select **Members**.
4. Locate the invitee, select the action icon in the invitee's row, and then select **Resend invitation**.
5. Select **Invite** to confirm.

## Remove a member or invitee from an organization

To remove a member or invitee from an organization:

1. Sign in to [Docker Admin](https://admin.docker.com){: target="_blank" rel="noopener" class="_"}.
2. In the left navigation, select your organization in the drop-down menu.
3. Select **Members**.
4. Locate the user, select the action icon in the user's row, and then select **Remove member** or **Remove invitee**.
5. Select **Remove** to confirm.

## Export members

Organization owners can export a CSV file containing the organization's members.
The CSV file contains the following fields:

 * **Name**: The user's name.
 * **Username**: The user's Docker ID.
 * **Email**: The user's email address.
 * **Type**: The type of user. For example, **Invitee** for users who have not accepted the organization's invite, or **User** for users who are members of the organization.
 * **Permissions**: The user's organization permissions. For example, **Member** or **Owner**.
 * **Teams**: The teams where the user is a member. A team is not listed for invitees.
 * **Date Joined**: The time and date when the user was invited to the organization.

To export a CSV file of the organization's members:


1. Sign in to [Docker Admin](https://admin.docker.com){: target="_blank" rel="noopener" class="_"}.
2. In the left navigation, select your organization in the drop-down menu.
3. Select **Members**.
4. Select **Export members** to download the CSV file.
