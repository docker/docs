---
description: Learn how to manage organization members in Docker Hub and Docker Admin Console.
keywords: members, teams, organizations, invite members, manage team members
title: Manage organization members
aliases:
- /docker-hub/members/
---

Learn how to manage members for your organization in Docker Hub and the Docker Admin Console.

{{< tabs >}}
{{< tab name="Docker Hub" >}}

{{% admin-users product="hub" %}}

{{< /tab >}}
{{< tab name="Admin Console" >}}

{{< include "admin-early-access.md" >}}

{{% admin-users product="admin" %}}

{{< /tab >}}
{{< /tabs >}}

## Manage members on a team

Use Docker Hub to add a member to a team or remove a member from a team.

### Add a member to a team

Organization owners can add a member to one or more teams within an organization.

{{< tabs >}}
{{< tab name="Docker Hub" >}}

To add a member to a team:

1. Sign in to [Docker Hub](https://hub.docker.com).
2. Select **Organizations**, your organization, and then **Members**.
3. Select the **Action** icon, and then select **Add to team**.

   > **Note**
   >
   > You can also navigate to **Organizations** > **Your Organization** > **Teams** > **Your Team Name** and select **Add Member**. Select a member from the drop-down list to add them to the team or search by Docker ID or email.
4. Select the team and then select **Add**.

   > **Note**
   >
   > The invitee must first accept the invitation to join the organization before being added to the team.

{{< /tab >}}
{{< tab name="Admin Console" >}}

{{< include "admin-early-access.md" >}}

To add a member to a team:

1. In the Admin Console, select your organization.
2. Select the team name.
3. Select **Add member**.

{{< /tab >}}
{{< /tabs >}}

### Remove a member from a team

To remove a member from a specific team:

1. Sign in to [Docker Hub](https://hub.docker.com).
2. Select **Organizations**, your organization, **Teams**, and then the team.
3. Select the **X** next to the user’s name to remove them from the team.
4. When prompted, select **Remove** to confirm.
