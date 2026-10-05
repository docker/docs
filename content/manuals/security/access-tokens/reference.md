---
title: Access token reference
linkTitle: Reference
description: >-
  Personal access token permissions, organization access token scopes,
  and the Docker Hub API endpoints an organization access token can call.
keywords: >-
  access token reference, personal access token permissions,
  organization access token scopes, Docker Hub API, OAT, PAT
weight: 30
---

Look up what an access token is allowed to do. You choose a permission
or a set of scopes when you create the token.

## Personal access token permissions

Each personal access token (PAT) has one permission level that you select when you
[create the token](/manuals/security/access-tokens/personal-access-tokens.md#create-a-personal-access-token).
The token applies to every repository your account can access. You
can't limit a token to a single repository.

- **Repo Public Read-only:** View, search, and pull images from public
  repositories. This is the default.
- **Repo Read-only:** View, search, and pull images from public
  repositories and from private repositories you have access to.
- **Repo Read & Write:** Everything in **Repo Read-only**, plus push
  images to any repository your account manages.
- **Repo Read, Write, Delete:** Everything in **Repo Read & Write**,
  plus delete images and manage your repositories.
- **Cloud Sandboxes:** Authenticate to the
  [Docker Cloud Sandboxes API](/manuals/ai/sandboxes-api/authentication.md).
  Cloud Sandboxes is a paid feature.

Choose the lowest permission that covers what the token needs to do.
For example, a CI job that only pulls a private base image needs
**Repo Read-only**. You can change PAT permissions at any time if the PAT is active or inactive, but not expired. 

## Organization access token scopes

Scopes control what an organization access token (OAT) can do. You choose
them when you
[create](/manuals/security/access-tokens/organization-access-tokens.md#create-an-organization-access-token)
or
[edit](/manuals/security/access-tokens/organization-access-tokens.md#update-deactivate-or-delete)
the token. In Docker
Home, each scope shows its name and a short description while the value
in the table is what the token carries. You can change OAT scopes at any time if the OAT is active or inactive,
but not expired.

Within a related set of operations, selecting a more capable scope
also grants the less capable ones. For example:

- **Image Delete** includes **Image Push**
- **Image Push** includes **Image Pull**

If you select **Image Delete**, you don't need to select the other two.
The Grants column notes each inclusion.

### Repository scopes

Repository scopes apply to each repository you add, or to all
repositories in the organization if you select
**All `<organization>` repositories**.

| Scope | Value | Grants |
| --- | --- | --- |
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
| --- | --- | --- |
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
| --- | --- | --- | --- |
| Docker Build Cloud | Cloud Connect | `scope-cloud-connect` | Connect to, build with, and run on Docker Build Cloud |
| Docker Governance | Audit Events Read | `scope-audit_events-read` | Read governance audit events |
| Docker Governance | Governance Policy Read | `scope-governance-policy-read` | Read governance policies |
| Docker Governance | Governance Policy Write | `scope-governance-policy-write` | Write governance policies. Includes Governance Policy Read |

## Docker Hub API

An OAT can authenticate most Docker Hub API endpoints under
`/v2/namespaces/{namespace}/repositories/`. First
[create an OAT](/manuals/security/access-tokens/organization-access-tokens.md#create-an-organization-access-token),
then exchange it for a short-lived bearer token with the
[Create access token](/reference/api/hub/latest/operations/AuthCreateAccessToken/)
API. Use your organization name as the identifier and the OAT as the
secret:

```console
$ TOKEN=$(curl -s -X POST "https://hub.docker.com/v2/auth/token" \
    -H "Content-Type: application/json" \
    -d '{"identifier": "<YOUR_ORGANIZATION_NAME>", "secret": "<YOUR_OAT>"}' \
    | jq -r .access_token)
```

Pass the bearer token in the `Authorization` header:

```console
$ curl -s -H "Authorization: Bearer $TOKEN" \
    "https://hub.docker.com/v2/namespaces/<YOUR_ORGANIZATION_NAME>/repositories"
```

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

- [Create a PAT](/manuals/security/access-tokens/personal-access-tokens.md#create-a-personal-access-token)
- [Create an OAT](/manuals/security/access-tokens/organization-access-tokens.md#create-an-organization-access-token)
- [Choose a PAT or OAT](/manuals/security/access-tokens/_index.md)
