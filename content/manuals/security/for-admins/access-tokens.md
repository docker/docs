---
title: Organization access tokens
description: Learn how to create and manage organization access tokens
  to securely push and pull images programmatically.
keywords: docker hub, security, OAT, organization access token
linkTitle: Organization access tokens
---

{{< summary-bar feature_name="OATs" >}}

> [!WARNING]
>
> Organization access tokens (OATs) are incompatible with Docker Desktop,
> [Image Access Management (IAM)](/manuals/security/for-admins/hardened-desktop/image-access-management.md), and [Registry Access Management (RAM)](/manuals/security/for-admins/hardened-desktop/registry-access-management.md).
>
> If you use Docker Desktop, IAM, or RAM, you must use personal
> access tokens instead.

An organization access token (OAT) is like a [personal access token
(PAT)](/security/for-developers/access-tokens/), but an OAT is associated with
an organization and not a single user account. Use an OAT instead of a PAT to
let business-critical tasks access Docker Hub repositories without connecting
the token to single user. You must have a [Docker Team or Business
subscription](/subscription/core-subscription/details/) to use OATs.

OATs provide the following advantages:

- You can investigate when the OAT was last used and then disable or delete it
  if you find any suspicious activity.
- You can limit what each OAT has access to, which limits the impact if an OAT
  is compromised.
- All company or organization owners can manage OATs. If one owner leaves the
  organization, the remaining owners can still manage the OATs.
- OATs have their own Docker Hub usage limits that don't count towards your
  personal account's limits.

If you have existing [service accounts](/docker-hub/service-accounts/),
Docker recommends that you replace the service accounts with OATs. OATs offer
the following advantages over service accounts:

- Access permissions are easier to manage with OATs. You can assign access
  permissions to OATs, while service accounts require using teams for access
  permissions.
- OATs are easier to manage. OATs are centrally managed in the Admin Console.
  For service accounts, you may need to sign in to that service account to
  manage it. If using single sign-on enforcement and the service account is not
  in your IdP, you may not be able to sign in to the service account to manage
  it.
- OATs are not associated with a single user. If a user with access to the
  service account leaves your organization, you may lose access to the service
  account. OATs can be managed by any company or organization owner.

## Create an organization access token

> [!IMPORTANT]
>
> Treat access tokens like a password and keep them secret. Store your tokens
> securely in a credential manager for example.

Company or organization owners can create up to:
- 10 OATs for organizations with a Team subscription
- 100 OATs for organizations with a Business subscription

Expired tokens count towards the total amount of tokens.

To create an OAT:

1. Sign in to [Docker Home](https://app.docker.com/) and select your
organization.
1. Select **Admin Console**, then **Access tokens**.
1. Select **Generate access token**.
1. Add a label and optional description for your token. Use something that
indicates the use case or purpose of the token.
1. Select the expiration date for the token.
1. Expand the **Repository** drop-down to set access permission
scopes for your token. To set Repository access scopes:
    1. Optional. Select **Read public repositories**.
    1. Select **Add repository** and choose a repository from the drop-down.
    1. Set the scopes for your repository &mdash; **Image Push** or
    **Image Pull**.
    1. Add more repositories as needed. You can add up to 50 repositories.
1. Optional. Expand the **Organization** drop-down and select the
**Allow management access to this organization's resources** checkbox. This
setting enables organization management scopes for your token. The following
organization management scopes are available:
    - **Member Edit**: Edit members of the organization
    - **Member Read**: Read members of the organization
    - **Invite Edit**: Invite members to the organization
    - **Invite Read**: Read invites to the organization
    - **Group Edit**: Edit groups of the organization
    - **Group Read**: Read groups of the organization
1. Select **Generate token**. Copy the token that appears on the screen
   and save it. You won't be able to retrieve the token once you exit the
   screen.

## Use an organization access token

You can use an organization access token when you sign in using Docker CLI.

Sign in from your Docker CLI client with the following command, replacing
`YOUR_ORG` with your organization name:

```console
$ docker login --username <YOUR_ORG>
```

When prompted for a password, enter your organization access token instead of a
password.

## Modify existing tokens

You can rename, update the description, update the repository access,
deactivate, or delete a token as needed.

1. Sign in to [Docker Home](https://app.docker.com/) and select your
organization.
1. Select **Admin Console**, then **Access tokens**.
1. Select the actions menu in the token row, then select **Deactivate**, **Edit**, or **Delete** to modify the token. For **Inactive** tokens, you can only select **Delete**.
1. If editing a token, select **Save** after specifying your modifications.
