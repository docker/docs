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

Unlike a personal access token, an OAT belongs to the organization. Any
organization owner can see, edit, deactivate, or delete it. Each token
has its own scopes, so it can only reach the repositories and
organization settings you choose.

> [!WARNING]
>
> Organization access tokens don't work with Docker Desktop or Image
> Access Management. For those features, use
> [personal access tokens][pat] instead.

## Best practices

- Set an expiration date and rotate tokens on a schedule
- Grant only the repositories and scopes each job needs
- Review the **Last used** column to find unused or suspicious tokens
- Store tokens in a credential manager, never in plain text or source
  code
- Deactivate or delete a token as soon as it's compromised or no longer
  needed

## Prerequisites

To create and manage organization access tokens, you need:

- A Docker Team or Docker Business subscription
- One of these roles in the organization:
  - Organization owner
  - Company owner, for organizations that belong to a company
  - A [custom role][custom-roles] that includes the
    **Manage organization access tokens** permission

Each organization can have a limited number of tokens:

- Docker Team: Up to 10 tokens
- Docker Business: Up to 100 tokens

Expired and deactivated tokens count toward the limit until you delete
them. When you reach the limit, the **Generate access token** button is
disabled until you delete a token.

## Create an organization access token

> [!IMPORTANT]
>
> Treat access tokens like passwords and keep them secure. Store tokens
> in a credential manager and never commit them to source code.

To create an organization access token:

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
     [Repository scopes](#repository-scopes).
   - Repeat for up to 50 repositories. Each repository has its own
     scopes.
1. Optional. Expand **Organization** and select the organization-level
   scopes the token needs, such as reading members or creating
   repositories. See [Organization scopes](#organization-scopes).
1. Optional. Expand **Docker Build Cloud** or **Docker Governance** to
   grant access to those products. See
   [Product scopes](#product-scopes).
1. Select **Generate token**. Copy the token and save it. Docker shows
   the token once and doesn't store it. You can't retrieve it after you
   leave the page.

## Sign in with an organization access token

Run `docker login` with your organization name as the username. When
the CLI asks for a password, paste the organization access token.

```console
$ docker login --username <YOUR_ORGANIZATION_NAME>
Password: [paste your OAT here]
```

## Update, deactivate, or delete a token {#modify-existing-tokens}

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

## Available scopes

Scopes control what a token can do. You choose them when you create or
edit the token. In the Docker Home UI, each scope shows its name and a
short description; the value in the table is what the token carries.

Where a scope includes another one, selecting the higher scope grants
the lower one too.

### Repository scopes

Repository scopes apply to each repository you add, or to all
repositories in the organization if you select
**All `<organization>` repositories**.

| Scope | Value | Grants |
|---|---|---|
| Image Pull | `scope-image-pull` | Pull images |
| Image Push | `scope-image-push` | Push images. Includes Image Pull |
| Image Delete | `scope-image-delete` | Delete images and tags through registry endpoints. Includes Image Push |
| Repository Read | `scope-repository-read` | Read repository metadata, the Dockerfile, and stars |
| Repository Edit | `scope-repository-edit` | Edit privacy, categories, Dockerfile, description, and stars. Includes Repository Read |
| Repository Admin | `scope-repository-admin` | Delete the repository. Includes Repository Edit |
| Tag Read | `scope-tag-read` | List and read tags, image lists, attestations, and compose files |
| Tag Admin | `scope-tag-admin` | Delete tags. Includes Tag Read |
| Webhook Read | `scope-webhook-read` | List webhook pipelines and delivery history |
| Webhook Edit | `scope-webhook-edit` | Create webhook pipelines. Includes Webhook Read |
| Webhook Admin | `scope-webhook-admin` | Delete webhook pipelines. Includes Webhook Edit |
| Repository Group Read | `scope-repo-group-read` | List and read repository group assignments |
| Repository Group Edit | `scope-repo-group-edit` | Create and update repository group assignments. Includes Repository Group Read |
| Repository Group Admin | `scope-repo-group-admin` | Delete repository group assignments. Includes Repository Group Edit |
| Repository Settings Admin | `scope-repository-settings-admin` | Configure immutable tag rules |

### Organization scopes

Organization scopes apply to the whole organization.

| Scope | Value | Grants |
|---|---|---|
| Member Read | `scope-member-read` | Read organization members |
| Member Edit | `scope-member-edit` | Edit organization members. Includes Member Read |
| Invite Read | `scope-invite-read` | Read invitations |
| Invite Edit | `scope-invite-edit` | Edit invitations. Includes Invite Read |
| Group Read | `scope-group-read` | Read the organization's groups (teams) |
| Group Edit | `scope-group-edit` | Edit the organization's groups (teams). Includes Group Read |
| Audit Log Read | `scope-activity-read` | Read the organization's activity logs |
| SIEM Credentials Read | `scope-siem_credentials-read` | Read SIEM destination settings, including credentials |
| Registry Access Management Read | `scope-ram-read` | Read Registry Access Management settings |
| Registry Access Management Edit | `scope-ram-write` | Edit Registry Access Management settings. Includes Registry Access Management Read |
| Report Read | `scope-report-read` | Download organization usage reports |
| Repository Create | `scope-repository-create` | Create repositories in the organization namespace |
| Repository List | `scope-repository-list` | List all repositories in the namespace, including private ones |
| Registry Usage Read | `scope-registry-usage-read` | Read namespace-level registry usage metrics |

Creating a repository requires the organization-level **Repository
Create** scope. No repository scope grants it, not even
**Repository Admin** on an existing repository.

### Product scopes

These sections appear alongside **Repository** and **Organization** in
the token's resources.

| Section | Scope | Value | Grants |
|---|---|---|---|
| Docker Build Cloud | Cloud Connect | `scope-cloud-connect` | Connect to, build with, and run on Docker Build Cloud |
| Docker Governance | Audit Events Read | `scope-audit_events-read` | Read governance audit events |
| Docker Governance | Governance Policy Read | `scope-governance-policy-read` | Read governance policies |
| Docker Governance | Governance Policy Write | `scope-governance-policy-write` | Write governance policies. Includes Governance Policy Read |

## Use the Docker Hub API {#hub-api-support}

An OAT can authenticate most Docker Hub API endpoints under
`/v2/namespaces/{namespace}/repositories/`. Exchange the token for a
bearer token with your organization name as the username, then pass it
in the `Authorization` header.

### Supported endpoints

The following endpoint groups accept OAT authentication:

- Repositories: list, create, get, update, and delete
- Tags: list, get, and delete; get tag images, attestations, and compose
  files
- Dockerfile: get and update a repository's linked Dockerfile
- Repository groups: list, get, create, update, and delete assignments
- Stars: list, count, add, and remove
- Immutable tags: update and verify policies
- Repository categories, privacy, and webhook pipeline settings
- Namespace metrics

### Listing behavior

`GET /v2/namespaces/{namespace}/repositories` filters results by the
token's scopes:

- With **Repository List** (`scope-repository-list`), the response
  includes every repository, including private ones.
- Without it, the response includes only public repositories.

The filtering is silent. The response is a normal `200` with no
indication that private repositories were left out.

### Unsupported legacy endpoints

OATs only work with the namespace-scoped routes described above. The
following legacy paths reject every OAT, regardless of its scopes, with
`403 token issued from organization access token is not allowed`. Use
the replacement endpoint instead:

- `GET /v2/repositories/{namespace}/{repository}`: use
  [Get repository](/reference/api/hub/latest/operations/GetRepository/)
- `GET /v2/repositories/{namespace}`: use
  [List repositories](/reference/api/hub/latest/operations/listNamespaceRepositories/)
- `GET /v2/users/{username}/repositories`: use
  [List repositories](/reference/api/hub/latest/operations/listNamespaceRepositories/)

## Next steps

- [Choose a personal or organization access token][overview]
- [Create a personal access token][pat]
- [Set up OIDC connections for GitHub Actions][oidc]
- [Review custom role permissions][custom-roles]

[overview]: /manuals/security/access-tokens/_index.md
[pat]: /manuals/security/access-tokens/personal-access-tokens.md
[oidc]: /manuals/security/authentication/oidc-connections/_index.md
[custom-roles]: /manuals/security/roles-and-permissions/custom-roles/permissions-reference.md
