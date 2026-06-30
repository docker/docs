---
title: OIDC connections rulesets and subject claims
linkTitle: Rulesets and subject claims
description: Configure rulesets and subject claims to control access for OIDC connections
keywords: oidc connections, rulesets, subject claims, openid connect, token claims, access control, enterprise security, admin
tags: [admin]
weight: 20
---

{{< summary-bar feature_name="OIDC connections" >}}

Rulesets and subject claims define what actions your GitHub workflows can take with your Docker resources. This page explains how to configure rulesets and set subject claims to authorize GitHub Workflow behaviors.

## Rulesets

A ruleset is a set of conditions that Docker evaluates against an incoming GitHub ID token. When a workflow triggers an OIDC exchange, Docker checks the token against every ruleset defined in your connection. If a ruleset's conditions are satisfied, Docker grants access based on the parameters set by that ruleset.

Each ruleset contains the following fields:

- **Label**: A name for the ruleset.
- **Rules**: One or more conditions based on OIDC token claims, such as the repository name, branch, or workflow path.
  - These are expressed as subject claim strings.
  - See [Subject claims](#subject-claims).
- **Resources**: The Docker resources a workflow can access when the ruleset matches. See [Resources](#resources).
- **Scopes**: The permissions granted on those resources, such as read or write access.

You can define between 1 and 5 rulesets per connection. Use multiple rulesets to apply different access levels across different workflows or branches.

> [!TIP]
> If more than one ruleset matches an incoming token,
> Docker merges the resources from all matching rulesets
> and grants access to the combined set.

## Subject claims

A subject claim is the `sub` field in a GitHub-issued JWT ID token. It encodes details of a workflow into a single string, identifying the workflow by organization, repository, branch, environment, and so on.

Docker uses the subject claim as the primary condition when evaluating your ruleset rules. The default subject claim format is:

```text
repo:<org>/<repo>:ref:refs/heads/<branch>
```

For example:

```text
repo:octo-org/octo-repo:ref:refs/heads/main
```

The exact format varies and depends on what triggered the workflow.

- A branch push, pull request, tag, or environment deployment each produces a different `sub` value.
- Refer to [GitHub's OpenID Connect Reference](https://docs.github.com/en/actions/reference/security/oidc) for the full list of formats.

You can use wildcards to match across repositories or branches:

| Pattern                                        | Matches                                         |
| :--------------------------------------------- | :---------------------------------------------- |
| `repo:my-org/my-repo:ref:refs/heads/main`      | Only the `main` branch of a specific repository |
| `repo:my-org/*`                                | All repos in the organization                   |
| `repo:my-org/my-repo:ref:refs/heads/release-*` | All branches starting with `release-`           |

## Resources

Resources define the Docker resources a workflow can access when a ruleset matches. You specify resources per ruleset, alongside the scopes that determine the level of access granted.

Docker Hub repositories are supported resources.

## What's next

- Learn about [OIDC connections](/manuals/enterprise/security/oidc-connections/_index.md)
- [Create or manage OIDC connections](/manuals/enterprise/security/oidc-connections/create-manage.md)
