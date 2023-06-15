---
description: Manage company users
keywords: company, company users, users, admin, docker admin
title: Manage company users
---

{% include admin-early-access.md %}

## Invite members

Company owners can invite new members to an organization in the company via Docker ID, email address, or via a CSV file containing email addresses. If an invitee does not have a Docker account, they must create an account and verify their email address before they can accept the invitation to join the organization. When inviting members, their pending invitation occupies a seat.

### Invite members via Docker ID or email address

Use the following steps to invite members to an organization in your company via Docker ID or email address. To invite a large amount of members to your company, Docker recommends that you [invite members via CSV file](#invite-members-via-csv-file).

1. Sign in to [Docker Admin](https://admin.docker.com){: target="_blank" rel="noopener" class="_"}.
2. In the left navigation, select your company in the drop-down menu.
3. Select **Users**.
4. Select **Invite**.
5. Select **Emails Or Docker IDs**.
6. Enter the Docker IDs or email addresses that you want to invite, up to a maximum of 1000. Separate multiple entries by a comma, semicolon, or space.
7. Select an organization from the drop-down list to add all invited users to that organization.
8. Select a team or type to create a new team. Docker will invite all users to that team.
9. Select **Invite** to confirm.
   > **Note**
   >
   > You can view the pending invitations in the **Users** page. The invitees receive an email with a link to the organization in Docker Hub where they can accept or decline the invitation.

### Invite members via CSV file

To invite multiple members to your organization in your company via a CSV file containing email addresses:
1. Sign in to [Docker Admin](https://admin.docker.com){: target="_blank" rel="noopener" class="_"}.
2. In the left navigation, select your company in the drop-down menu.
3. Select **Users**.
4. Select **Invite**.
5. Select **CSV Upload**.
6. Select an organization from the drop-down list to add all invited users to that organization.
7. Select a team or type to create a new team. Docker will invite all users to that team.
8. Select **Download the template CSV file** to optionally download an example CSV file. The following is an example of the contents of a valid CSV file.
    ```
    email
    docker.user-0@example.com
    docker.user-1@example.com
    ```
  CSV file requirements:
   -  The file must contain a header row with at least one heading named `email`. You can add additional columns but the import will ignore them.
   -  The file can contain a maximum of 1000 email addresses (rows). To invite more than 1000 users, create multiple CSV files and perform all steps in this task for each file.
9. Create a new CSV file or export a CSV file from another application.
  - To export a CSV file from another application, see the application’s documentation.
  - To create a new CSV file, open a new file in a text editor, type `email` on the first line, type the user email addresses one per line on the following lines, and then save the file with a .csv extension.
10. Select **Browse files** and then select your CSV file, or drag and drop the CSV file into the **Select a CSV file to upload** box. You can only select one CSV file at a time.
  > **Note**
  >
  > If the amount of email addresses in your CSV file exceeds the number of available seats in your organization, you can't continue to invite members. To invite members, you can buy more seats, or remove some email addresses from the CSV file and re-select the new file. To buy more seats, see [Add seats to your subscription](../../subscription/add-seats.md) or [Contact sales](https://www.docker.com/pricing/contact-sales/).
11. After the CSV file upload completes, select **Review**.
  Valid email addresses and any email addresses that have issues appear.
  Email addresses may have the following issues:
	  - **Invalid email**: The email address isn't a valid address. The email address will be ignored if you send invites. You can correct the email address in the CSV file and re-import the file.
	  - **Already invited**: Docker has already sent an invite email and Docker won't send another invite email.
	  - **Member**: The user is already a member of your organization and Docker won't send an invite email.
	  - **Duplicate**: The CSV file has multiple occurrences of the same email address. Docker will send the user only one invite email.
12. Select **Send invites**.
   > **Note**
   >
   > You can view the pending invitations in the **Users** page. The invitees receive an email with a link to the organization in Docker Hub where they can accept or decline the invitation.

## Resend invitations

To resend an invitation if the invite is pending or declined:

1. Go to [Docker Admin](https://admin.docker.com){: target="_blank" rel="noopener" class="_"}.
2. In the left navigation, select your company in the drop-down menu.
3. Select **Users**.
4. Locate the invitee, select the action icon in the invitee's row, and then select **Resend invitation**.
5. Select **Invite** to confirm.

## Remove a member or invitee from an organization

To remove a member or invitee from an organization:

1. Sign in to [Docker Admin](https://admin.docker.com){: target="_blank" rel="noopener" class="_"}.
2. In the left navigation, select your company in the drop-down menu.
3. Select **Users**.
4. Locate the user, select the action icon in the user's row, and then select **Remove user** or **Remove invitee**.
5. Select the organizations to remove the user from.
6. Select **Remove** to confirm.

## Export users

Company owners can export a CSV file containing all of the company's users.
The CSV file contains the following fields:

 * **Name**: The user's name.
 * **Username**: The user's Docker ID.
 * **Email**: The user's email address.
 * **Type**: The type of user. For example, **Invitee** for users who haven't accepted the organization's invite, or **User** for users who are members of the organization.
 * **Permissions**: The user's organization permissions. For example, **Member** or **Owner**.
 * **Teams**: The teams where the user is a member. A team isn't listed for invitees.
 * **Date Joined**: The time and date when the user was invited to the organization.

To export a CSV file of the company's users:

1. Sign in to [Docker Admin](https://admin.docker.com){: target="_blank" rel="noopener" class="_"}.
2. In the left navigation, select your company in the drop-down menu.
3. Select **Users**.
4. Select the action icon next to **Invite**, and then select **Export users as CSV**.