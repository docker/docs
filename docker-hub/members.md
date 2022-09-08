---
description: Manage organization members
keywords: members, teams, organizations
title: Manage members
---


This section describes how to manage members in your [teams and organizations](../docker-hub/orgs.md).

## Invite members

Organization owners can invite new members to an organization via Docker ID, email address, or via a CSV file containing email addresses. If an invitee does not have a Docker account, they must create an account and verify their email address before they can accept the invitation to join the organization. When inviting members, their pending invitation occupies a seat.

### Invite members via Docker ID or email address

Use the following steps to invite members to your organization via Docker ID or email address. To invite a large amount of members to your organization, the recommended method is to [invite members via CSV file](#invite-members-via-csv-file).

1. Go to **Organizations** in [Docker Hub](https://hub.docker.com){: target="_blank" rel="noopener" class="_"}, and select your organization.
2. In the **Members** tab, select **Invite Member**.
3. Select **Emails Or Docker IDs**.
4. Enter the Docker IDs or email addresses that you want to invite, up to a maximum of 1000. Separate multiple entries by a comma, semicolon, or space.
5. Select a team from the drop-down list to add all invited users to that team.
  > **Note**
  >
  >  It is recommended that you invite non-administrative users to a team other than the owners team. Members in the owners team will have full access to your organization’s administrative settings. To create a new team, see [Create a team](../docker-hub/orgs.md/#create-a-team).
6. Click **Invite** to confirm.
   > **Note**
   >
   > You can view the pending invitations in the **Members** tab. The invitees receive an email with a link to the organization in Docker Hub where they can accept or decline the invitation.


### Invite members via CSV file

To invite multiple members to your organization via a CSV file containing email addresses:

1. Go to **Organizations** in [Docker Hub](https://hub.docker.com){: target="_blank" rel="noopener" class="_"}, and select your organization.
2. In the **Members** tab, select **Invite Member**.
3. Select **CSV Upload**.
4. Select a team from the drop-down list to add all invited users to that team.
  > **Note**
  >
  >  It is recommended that you invite non-administrative users to a team other than the owners team. Members in the owners team will have full access to your organization’s administrative settings. To create a new team, see [Create a team](../docker-hub/orgs.md/#create-a-team).
5. Select **Download the template CSV file** to optionally download an example CSV file. The following is an example of the contents of a valid CSV file.
    ```
    email
    docker.user-0@example.com
    docker.user-1@example.com
    ```
  CSV file requirements:
   -  The file must contain a header row with at least one heading named `email`. Additional columns are allowed and are ignored in the import.
   -  The file must contain a maximum of 1000 email addresses (rows). To invite more than 1000 users, create multiple CSV files and perform all steps in this task for each file.
6. Create a new CSV file or export a CSV file from another application.
  - To export a CSV file from another application, see the application’s documentation.
  - To create a new CSV file, open a new file in a text editor, type `email` on the first line, type the user email addresses one per line on the following lines, and then save the file with a .csv extension.
7. Select **Browse files** and then select your CSV file, or drag and drop the CSV file into the **Select a CSV file to upload** box. You can only select one CSV file at a time.
  > **Note**
  >
  > If the amount of email addresses in your CSV file exceeds the number of available seats in your organization, you cannot continue to invite members. To invite members, you can purchase more seats, or remove some email addresses from the CSV file and re-select the new file. To purchase more seats, see [Add seats to your subscription](../subscription/add-seats.md) or [Contact sales](https://www.docker.com/pricing/contact-sales/).
8. After the CSV file has been uploaded, select **Review**.
  Valid email addresses and any email addresses that have issues appear.
  Email addresses may have the following issues:
	  - **Invalid email**: The email address is not a valid address. The email address will be ignored if you send invites. You can correct the email address in the CSV file and re-import the file.
	  - **Already invited**: The user has already been sent an invite email and another invite email will not be sent.
	  - **Member**: The user is already a member of your organization and an invite email will not be sent.
	  - **Duplicate**: The CSV file has multiple occurrences of the same email address. The user will be sent only one invite email.
9. Click **Send invites**.
   > **Note**
   >
   > You can view the pending invitations in the **Members** tab. The invitees receive an email with a link to the organization in Docker Hub where they can accept or decline the invitation.

## Add a member to a team

Organization owners can add a member to one or more teams within an organization.

To add a member to a team:

1. Navigate to **Organizations** in [Docker Hub](https://hub.docker.com){: target="_blank" rel="noopener" class="_"}, and select your organization.
2. In the **Members** tab, click the additional options from the table menu and select **Add to team**.

    > **Note**
    >
    > You can also navigate to **Organizations** > **Your Organization** > **Teams** > **Your Team Name** and click **Add Member**. Select a member from the drop-down list to add them to the team or search by Docker ID or email.

3. Select the team and click **Add**.

    > **Note**
    >
    > The invitee must first accept the invitation to join the organization before being added to the team.

## Resend invitations

To resend an invitation if the invite is pending or declined:

1. Navigate to **Organizations** in [Docker Hub](https://hub.docker.com){: target="_blank" rel="noopener" class="_"} and select your organization.
2. In the **Members** tab, locate the invitee and select **Resend invitation** from the table menu.
3. Click **Invite** to confirm.

## Remove members

To remove a member from an organization:

1. Navigate to **Organizations** in [Docker Hub](https://hub.docker.com){: target="_blank" rel="noopener" class="_"}, and select your organization.
2. In the **Members** tab, select Remove member from the table menu.
3. When prompted, click **Remove** to confirm.

To remove an invitee from an organization:

1. Navigate to **Organizations** in [Docker Hub](https://hub.docker.com){: target="_blank" rel="noopener" class="_"}, and select your organization.
2. In the **Members** tab, locate the invitee you would like to remove and select **Remove invitee** from the table menu.
3. When prompted, click **Remove** to confirm.

To remove a member from a specific team:

1. Navigate to **Organizations** in [Docker Hub](https://hub.docker.com){: target="_blank" rel="noopener" class="_"}, and select your organization.
2. Click on the **Teams** tab and select the team from the list.
3. Click the **X** next to the user’s name to remove them from the team.
4. When prompted, click **Remove** to confirm.

## Export members

Organization owners can export a CSV file containing the organization's members.
The CSV file contains the following fields:

 * **Name**: The user's name.
 * **Username**: The user's Docker ID.
 * **Email**: The user's email address.
 * **Type**: The type of user. For example, **Invitee** for users who have not accepted the organization's invite, or **User** for users who are members of the organization.
 * **Permissions**: The user's organization permissions. For example, **Member** or **Owner**.
 * **Teams**: The teams where the user is a member.  A team is not listed for invitees.
 * **Date Joined**: The time and date when the user was invited to the organization.

To export a CSV file of an organization's members:

1. Navigate to **Organizations** in [Docker Hub](https://hub.docker.com){: target="_blank" rel="noopener" class="_"}, and select your organization.
2. In the **Members** tab, select **Export members** to download the CSV file.