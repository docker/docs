---
title: Create and manage organization access tokens
linkTitle: Organization access tokens
description: >-
  Create an organization access token, choose its repository and
  organization scopes, and use it in CI/CD and the Docker Hub API.
keywords: >-
  organization access token, OAT, Docker Hub, CI/CD authentication,
  docker login, token scopes, repository permissions, Docker Hub API,
  automation, Docker Team, Docker Business
weight: 20
aliases:
  - /security/for-admins/access-tokens/
  - /enterprise/security/access-tokens/
---

{{< summary-bar feature_name="OATs" >}}

An organization access token (OAT) lets automated systems sign in to
Docker Hub as your organization instead of as a person. Use an OAT for
CI/CD pipelines, deployment jobs, and other automation that must keep
working when people join or leave the organization.

Unlike a personal access token (PAT), an OAT belongs to the
organization. Any organization owner can see, edit, deactivate, or
delete it. Each token has its own scopes, so it can only reach the
repositories and organization settings you choose.

> [!WARNING]
>
> OATs don't work with Docker Desktop or Image Access Management. For
> those features, use [PATs][pat] instead.

## Best practices

- Set an expiration date and rotate tokens on a schedule
- Grant only the repositories and scopes each job needs
- Review the **Last used** column to find unused or suspicious tokens
- Store tokens in a credential manager, never in plain text or source
  code
- Deactivate or delete a token as soon as it's compromised or no longer
  needed

## Prerequisites

To create and manage OATs, you need:

- A Docker Team or Docker Business subscription
- One of these roles in the organization:
  - Organization owner
  - Company owner, for organizations that belong to a company
  - A [custom role][custom-roles] that includes the
    **Manage organization access tokens** permission

## Create

> [!IMPORTANT]
>
> Treat access tokens like passwords and keep them secure. Store tokens
> in a credential manager and never commit them to source code.

To create an OAT:

1. Sign in to [Docker Home](https://app.docker.com/) and select your
   organization.
1. Select **Identity & auth**, then **Access tokens**.
1. Select **Generate access token**.
1. Configure the token:
   - **Label:** A name that describes what the token is for. This field
     is required.
   - **Access token description:** Optional. Up to 200 characters.
   - **Expiration date:** Select **30 days**, **90 days**, or
     **Custom**. With **Custom**, choose a date and time up to one year
     from today. The default, **None**, creates a token that doesn't
     expire. You can't change the expiration date after you create the
     token.
1. In **Resources**, expand **Repository** to choose which repositories
   the token can reach:
   - Optional. Select **Read public repositories** to let the token pull
     from any public repository.
   - Select **Add repository**, then choose a repository or
     **All `<organization>` repositories** from the drop-down.
   - Select one or more scopes for that repository. See
     [Repository scopes][repo-scopes].
   - Repeat for up to 50 repositories. Each repository has its own
     scopes.
1. Optional. Expand **Organization** and select the organization-level
   scopes the token needs, such as reading members or creating
   repositories. See [Organization scopes][org-scopes].
1. Optional. Expand **Docker Build Cloud** or **Docker Governance** to
   grant access to those products. See
   [Product scopes][product-scopes].
1. Select **Generate token**. Copy the token and save it. Docker shows
   the token once and doesn't store it. You can't retrieve it after you
   leave the page.

## Sign in

Run `docker login` with your organization name as the username. When
the CLI asks for a password, paste the OAT.

```console
$ docker login --username <YOUR_ORGANIZATION_NAME>
Password: [paste your OAT here]
```

## Update, deactivate, or delete

You can rename a token, change its description or scopes, deactivate
it, activate it again, or delete it.

1. Sign in to [Docker Home](https://app.docker.com/) and select your
   organization.
1. Select **Identity & auth**, then **Access tokens**.

   The list shows each token's label, status, who created it, when it
   was created, when it was last used, and when it expires.

1. Select the actions menu at the end of a token row, then select
   **Deactivate**, **Activate**, **Edit**, or **Delete**.

   A deactivated token stops working until you activate it again.
   Expired tokens can only be deleted. Deleting a token is permanent.

1. If you selected **Edit**, change the label, description, or
   resources, then select **Update token**.

## Token limits

Each organization can have a limited number of tokens:

- Docker Team: Up to 10 tokens
- Docker Business: Up to 100 tokens

Expired and deactivated tokens count toward the limit until you delete
them. When you reach the limit, the **Generate access token** button is
disabled until you delete a token.

## Next steps

- [Choose a PAT or OAT][overview]
- [Look up scopes and Docker Hub API support][reference]
- [Create a PAT][pat]
- [Set up OIDC connections for GitHub Actions][oidc]
- [Review custom role permissions][custom-roles]

[overview]: /manuals/security/access-tokens/_index.md
[reference]: /manuals/security/access-tokens/reference.md
[repo-scopes]: /manuals/security/access-tokens/reference.md#repository-scopes
[org-scopes]: /manuals/security/access-tokens/reference.md#organization-scopes
[product-scopes]: /manuals/security/access-tokens/reference.md#product-scopes
[pat]: /manuals/security/access-tokens/personal-access-tokens.md
[oidc]: /manuals/security/authentication/oidc-connections/_index.md
[custom-roles]: /manuals/security/roles-and-permissions/custom-roles/permissions-reference.md
