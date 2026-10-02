---
title: Configure OIDC rulesets and subject claims
linkTitle: Rulesets and subject claims
description: >-
  Rulesets and subject claims control which GitHub Actions workflows can
  push images, pull images, read public repositories, and connect to Docker
  Build Cloud.
keywords: >-
  oidc connections, rulesets, subject claims, github actions, openid connect,
  image push, image pull, docker build cloud, wildcards, access control
tags: [admin]
weight: 20
aliases:
  - /enterprise/security/oidc-connections/rulesets-claims/
---

{{< summary-bar feature_name="OIDC connections" >}}

Rulesets and subject claims define what actions your GitHub workflows can
take with your Docker resources. Use them to authorize GitHub workflow
behaviors for an OIDC connection.

## Rulesets

A ruleset tells Docker which GitHub workflow can access which Docker
resources. When a workflow triggers an OIDC exchange, Docker checks the
token's subject claim against every ruleset on the connection. If the
subject claim matches, Docker grants the resources and scopes on that
ruleset.

Each ruleset has these fields:

- **Ruleset name**: A name for the ruleset.
- **Subject claim**: One `sub` string from the GitHub ID token. The
  repository, branch, tag, pull request, or environment is part of that
  string, not a separate field. See [Subject claims](#subject-claims).
- **Resources**: What the workflow can access when the subject claim
  matches. See [Resources](#resources).
- **Scopes**: The access granted on those resources.
  - Repository: **Image Pull**, **Image Push**, and **Read public
    repositories**
  - Docker Build Cloud: **Cloud Connect**

You can add one to five rulesets on a connection. If more than one
ruleset matches a token, Docker combines their resources and grants
access to that combined set. For more information, see
[Ruleset examples](#ruleset-examples).

## Subject claims

A subject claim is the `sub` field in a GitHub-issued JWT ID token. It
encodes details of a workflow into a single string, identifying the
workflow by organization, repository, branch, environment, and so on.

On each ruleset, enter that string in **Subject claim**. The default
format is:

```text
repo:<org>/<repo>:ref:refs/heads/<branch>
```

You can use wildcards to match across repositories or branches:

| Pattern                                        | Matches                                              |
| :--------------------------------------------- | :--------------------------------------------------- |
| `repo:my-org/my-repo:ref:refs/heads/main`      | Only the `main` branch of a specific repository      |
| `repo:my-org/*`                                | All repositories in the organization                 |
| `repo:my-org/my-repo:ref:refs/heads/release-*` | All branches starting with `release-`                |

> [!NOTE]
> GitHub repositories created after July 15, 2026 use immutable
> identifiers for default subject claims. For example:
> `repo:octocat@123456/my-repo@456789:ref:refs/heads/main`. See the
> [GitHub changelog](https://github.blog/changelog/2026-04-23-immutable-subject-claims-for-github-actions-oidc-tokens/)
> for more details.

For the full list of formats, see
[GitHub's OpenID Connect reference](https://docs.github.com/en/actions/reference/security/oidc).

## Resources

A resource is a Docker product the workflow can use after the subject
claim matches. The subject claim decides whether a ruleset applies. The
resources on that ruleset decide what the workflow can reach, and the
scopes decide what it can do there.

You set resources on each ruleset, next to the scopes for those
resources. A workflow receives only the resources from rulesets that
match its token. When more than one ruleset matches, Docker combines
those resources.

Docker Hub repositories and Docker Build Cloud are the supported
resources.

- A repository resource is a Docker Hub repository in your organization.
  **Image Pull** and **Image Push** apply to that repository. **Read
  public repositories** applies to public repositories.
- A Docker Build Cloud resource uses the **Cloud Connect** scope, so the
  workflow can connect and build in Docker Build Cloud.

## Ruleset examples

These are different ways to set up rulesets on a connection. Each
example is its own setup.

### One ruleset for several resources

If you have workflows on `main` that need to push the
`octo-org/octo-repo` image and build it in Docker Build Cloud, you can
set up one ruleset:

| Ruleset name | Subject claim                                 | Resources and scopes                                                                               |
| :----------- | :-------------------------------------------- | :------------------------------------------------------------------------------------------------- |
| `main`       | `repo:octo-org/octo-repo:ref:refs/heads/main` | Repository `octo-org/octo-repo` with **Image Push**, and Docker Build Cloud with **Cloud Connect** |

A push to `main` makes GitHub set the token's subject claim to
`repo:octo-org/octo-repo:ref:refs/heads/main`, which Docker compares
against the ruleset. When the string matches the ruleset named `main`,
the run can push the image and connect to Docker Build Cloud.

Choose this setup when you want to update both permissions
together.

### Two rulesets for different branches

You may want to create different rulesets for different branches. For
example, release branches need to pull `octo-org/octo-repo` while any
workflows on `main` need to push that image. In this case it makes sense
to create two rulesets:

| Ruleset name   | Subject claim                                      | Resources and scopes                                |
| :------------- | :------------------------------------------------- | :-------------------------------------------------- |
| `release-pull` | `repo:octo-org/octo-repo:ref:refs/heads/release-*` | Repository `octo-org/octo-repo` with **Image Pull** |
| `main-push`    | `repo:octo-org/octo-repo:ref:refs/heads/main`      | Repository `octo-org/octo-repo` with **Image Push** |

The subject claim uses `release-*` to extend the permission to every branch whose
name starts with `release-`, letting the workflow pull the image for
all release branches. On the other hand, a push to `main`
matches the ruleset named `main-push`, so that workflow can push the
image.

Docker returns `access denied` on any other branch you leave off the
rulesets. For example, a branch named `feature-x` can still send its
token to Docker, but the subject claim matches neither ruleset.

Choose this setup when branches need different access. You can change
what release branches are allowed to do without editing `main-push`.

### Rulesets that match the same run

A single run can match multiple rulesets. For example:

| Ruleset name | Subject claim                                 | Resources and scopes                                   |
| :----------- | :-------------------------------------------- | :----------------------------------------------------- |
| `main-push`  | `repo:octo-org/octo-repo:ref:refs/heads/main` | Repository `octo-org/octo-repo` with **Image Push**    |
| `org-read`   | `repo:octo-org/*`                             | Public repositories with **Read public repositories** |

When you push to `main` in this instance, the subject claim matches the
ruleset named `main-push` and the ruleset named `org-read`. It matches
`org-read` because `octo-repo` is a repository in `octo-org`. Docker
combines their permissions, so the workflow that pushed to `main` can
push the image and read public repositories.

Alternatively, a push to any other branch or repository in `octo-org`
matches `org-read` only. That run can read public repositories.

Choose this setup when the organization-wide read permission and the
`main` push permission should be added and removed separately. You could
then delete `org-read` to remove read access to public repositories and
leave the `main-push` ruleset as it is.

## Next steps

- [OIDC connections overview](/manuals/security/authentication/oidc-connections/_index.md)
- [Create or manage OIDC connections](/manuals/security/authentication/oidc-connections/create-manage.md)
- [Troubleshoot OIDC connections](/manuals/security/authentication/oidc-connections/troubleshoot.md)
