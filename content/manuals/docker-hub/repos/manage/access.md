---
description: Discover how to manage access to repositories on Docker Hub.
keywords: Docker Hub, Hub, repository access, repository collaborators, repository privacy
title: Access management
LinkTItle: Access
weight: 50
aliases:
- /docker-hub/repos/access/
---

In this topic learn about the features available to manage access to your
repositories. This includes visibility, collaborators, roles, teams, and
organization access tokens.

## Repository visibility

The most basic repository access is controlled via the visibility. A
repository's visibility can be public or private.

With public visibility, the repository appears in Docker Hub search results and
can be pulled by everyone. To manage push access to public personal
repositories, you can use collaborators. To manage push access to public
organization repositories, you can use roles, teams, or organization access
tokens.

With private visibility, the repository doesn't appear in Docker Hub search
results and is only accessible to those with granted permission. To manage push
and pull access to private personal repositories, you can use collaborators. To
manage push and pull access to private organization repositories, you can use
roles, teams, or organization access tokens.

### Change repository visibility

When creating a repository in Docker Hub, you can set the repository visibility.
In addition, you can set the default repository visibility when a repository is
created in your personal repository settings. The following describes how to
change the visibility after the repository has been created.

To change repository visibility:

1. Sign in to [Docker Hub](https://hub.docker.com).
2. Select **My Hub** > **Repositories**.
3. Select a repository.

   The **General** page for the repository appears.

4. Select the **Settings** tab.
5. Under **Visibility settings**, select one of the following:

   - **Make public**: The repository appears in Docker Hub search results and can be
     pulled by everyone.
   - **Make private**: The repository doesn't appear in Docker Hub search results
     and is only accessible to you and collaborators. In addition, if the
     repository is in an organization's namespace, then the repository
     is accessible to those with applicable roles or permissions.

6. Type the repository's name to verify the change.
7. Select **Make public** or **Make private**.

## Collaborators

A collaborator is someone you want to give `push` and `pull` access to a
personal repository. Collaborators aren't able to perform any administrative
tasks such as deleting the repository or changing its visibility from private to
public. In addition, collaborators can't add other collaborators.

Only personal repositories can use collaborators. You can add unlimited
collaborators to public repositories, and Docker Pro accounts can add up to 1
collaborator on private repositories.

Organization repositories can't use collaborators, but can use member roles,
teams, or organization access tokens to manage access.

### Manage collaborators

1. Sign in to [Docker Hub](https://hub.docker.com).

2. Select **My Hub** > **Repositories**.

   A list of your repositories appears.

3. Select a repository.

   The **General** page for the repository appears.

4. Select the **Collaborators** tab.

5. Add or remove collaborators based on their Docker username.

You can choose collaborators and manage their access to a private
repository from that repository's **Settings** page.

## Organization roles

Organizations can use roles for individuals, giving them different
permissions in the organization. For more details, see [Roles and
permissions](/manuals/security/roles-and-permissions.md).

## Organization teams

Organizations can use teams. A team can be assigned fine-grained repository
access.

### Configure team repository permissions

You must create a team before you are able to configure repository permissions.
For more details, see [Create and manage a
team](/manuals/accounts/organization/manage/manage-a-team.md).

To configure team repository permissions:

1. Sign in to [Docker Hub](https://hub.docker.com).

2. Select **My Hub** > **Repositories**.

   A list of your repositories appears.

3. Select a repository.

   The **General** page for the repository appears.

4. Select the **Permissions** tab.

5. Add, modify, or remove a team's repository permissions.

   - Add: Specify the **Team**, select the **Permission**, and then select **Add**.
   - Modify: Specify the new permission next to the team.
   - Remove: Select the **Remove permission** icon next to the team.

## Organization access tokens (OATs)

Organizations can use OATs. OATs let you assign fine-grained repository access
permissions to tokens. For more details, see [Organization access
tokens](/manuals/security/access-tokens/organization-access-tokens.md).

## Gated distribution

{{< summary-bar feature_name="Gated distribution" >}}

Gated distribution allows publishers to securely share private container images
with external customers or partners, without giving them full organization
access or visibility into your teams, collaborators, or other repositories.
Content stays in private repositories, and external users can pull from them
without being added to your internal organization.

This feature is ideal for commercial software publishers who want to control who
can pull specific images while preserving a clean separation between internal
users and external consumers.

If you are interested in Gated Distribution contact the <a
href="https://www.docker.com/pricing/contact-sales/"
id="dkr_docs_cs_hub_gated_distribution" class="link" rel="noopener">Docker Sales
Team</a> for more information.

### Distributor members

When you invite users to an organization entitled with gated distribution, you
assign them a role that determines their level of access. For gated
distribution, external users are invited as **distributor members** within a
specific team that you create. This role grants pull-only access to the
repositories assigned to that team, and nothing else in your organization. See
[Roles and permissions](/manuals/security/roles-and-permissions.md)
for details about the access permissions for other roles.

Distributor members can't see other members in their team or organization, and
they can't see any repositories other than the ones their team has been granted
access to. This isolation is what makes the role suitable for distributing gated
images to external users, partners, or customers without exposing internal
collaborators, teams, or repositories.

Because distributor members consume licenses allocated for that role, consider
provisioning a separate organization for gated distribution rather than adding
distributor members to your main organization. This keeps regular members from
inadvertently consuming the licenses set aside for external distribution.

### Invite distributor members via API

Distributor members can only be invited using the Docker Hub API. UI-based
invitations are not supported for this role.

To invite distributor members:

1. Use the [authentication
   API](https://docs.docker.com/reference/api/hub/latest/#tag/authentication-api/operation/AuthCreateAccessToken)
   to generate a bearer token for your Docker Hub account. This token authorizes
   the API requests you use to send invites. Replace `myusername` and
   `dckr_pat_...` with your Docker ID and a [personal access
   token](/manuals/security/access-tokens.md):

   ```console
   $ TOKEN=$(curl -s -X POST "https://hub.docker.com/v2/auth/token" \
       -H "Content-Type: application/json" \
       -d '{"identifier": "myusername", "secret": "dckr_pat_..."}' \
       | jq -r .access_token)
   ```

2. Create a team to group the distributor members and assign them a shared set
   of repository permissions.

   {{< tabs >}}
   {{< tab name="Hub" >}}

   1. Sign in to [Docker Home](https://app.docker.com) and select your
      organization.
   2. Select **Teams**.
   3. Select **Create team**.
   4. Provide the team's information, then select **Create**.

   For more details, see [Create and manage a
   team](/manuals/accounts/organization/manage/manage-a-team.md).

   {{< /tab >}}
   {{< tab name="API" >}}

   Use the [teams
   API](https://docs.docker.com/reference/api/hub/latest/#tag/groups/paths/~1v2~1orgs~1%7Borg_name%7D~1groups/post):

   ```console
   $ curl -s -X POST "https://hub.docker.com/v2/orgs/example-org/groups" \
       -H "Authorization: Bearer $TOKEN" \
       -H "Content-Type: application/json" \
       -d '{"name": "customer-team", "description": "External distributor members"}'
   ```

   The response includes the new team's `id`. Save it for the next step, for
   example `GROUP_ID=12345`.

   {{< /tab >}}
   {{< /tabs >}}

3. Grant the team read-only access to the repositories you want distributor
   members to access.

   {{< tabs >}}
   {{< tab name="Hub" >}}

   1. Sign in to [Docker Hub](https://hub.docker.com).
   2. Select **My Hub** > **Repositories**.
   3. Select the repository.
   4. Select the **Permissions** tab.
   5. Specify the **Team**, select **Read-only** as the **Permission**, and
      then select **Add**.

   {{< /tab >}}
   {{< tab name="API" >}}

   Use the [repository teams
   API](https://docs.docker.com/reference/api/hub/latest/#tag/repositories/operation/CreateRepositoryGroup),
   passing the team's `id` from the previous step as `group_id`:

   ```console
   $ curl -s -X POST "https://hub.docker.com/v2/repositories/example-org/example-repo/groups" \
       -H "Authorization: Bearer $TOKEN" \
       -H "Content-Type: application/json" \
       -d "{\"group_id\": $GROUP_ID, \"permission\": \"read\"}"
   ```

   {{< /tab >}}
   {{< /tabs >}}

4. Use the [bulk create invites
   endpoint](https://docs.docker.com/reference/api/hub/latest/#tag/invites/paths/~1v2~1invites~1bulk/post)
   to send email invites. In the request body, set `role` to
   `distributor_member`, specify the `team`, and list the invitees' email
   addresses:

   ```console
   $ curl -s -X POST "https://hub.docker.com/v2/invites/bulk" \
       -H "Authorization: Bearer $TOKEN" \
       -H "Content-Type: application/json" \
       -d '{
             "org": "example-org",
             "team": "customer-team",
             "role": "distributor_member",
             "invitees": ["user@example.com"]
           }'
   ```

   This sends email invitations to the specified users and automatically assigns
   them to the team.

5. The invited user receives an email invitation from Docker Hub. When they
   select the link in the email, they sign in with their Docker ID (or create
   one if needed) to accept the invite. Once accepted, they're added to the
   organization as a distributor member with pull-only access to the
   repositories assigned to their team.
