---
title: Access tokens
linkTitle: Access tokens
description: >-
  Choose a personal or organization access token to authenticate to
  Docker Hub.
keywords: >-
  access tokens, personal access tokens, organization access tokens,
  PAT, OAT, Docker Hub, Docker security
weight: 10
grid:
  - title: Personal access tokens
    description: Authenticate the Docker CLI and tools with a token tied to your account.
    icon: lock-closed
    link: /security/access-tokens/personal-access-tokens/
  - title: Organization access tokens
    description: Grant org-owned Hub access to CI/CD and other automation.
    icon: building-office-2
    link: /security/access-tokens/organization-access-tokens/
  - title: Reference
    description: Look up PAT permissions, OAT scopes, and Hub API support.
    icon: list-bullet
    link: /security/access-tokens/reference/
---

Access tokens authenticate to Docker Hub in place of a password. Docker
Hub offers two types: personal access tokens (PATs), tied to your
account, and organization access tokens (OATs), owned by an
organization. Both work with `docker login` and in automation, and you
can deactivate or delete a token at any time without changing a
password.

## Choose a token type

| Token | Ownership | Use when | Subscription |
| --- | --- | --- | --- |
| [PAT][pat] | Your Docker account | CLI access, Docker Desktop, local tools, and automation that should run as you | [Docker Personal or Docker Pro][plans] |
| [OAT][oat] | The organization. Any organization owner can manage it | CI/CD and other automation that must keep working when membership changes | [Docker Team or Docker Business][plans] |

## PATs

Use a PAT for Docker Desktop, for environments governed by
[Image Access Management][iam], and for other tools that run as you. A
PAT is also required to sign in to the CLI when two-factor
authentication (2FA) is turned on or single sign-on (SSO) is enforced,
because the CLI doesn't accept your password in those cases.

A PAT uses [one permission level][permissions] for every repository the
account can access.

## OATs

Use an OAT for production systems that pull images during deployment,
monitoring or backup tools that check repository status or pull images,
third-party services that integrate with your repositories, and scripts
that call the [Docker Hub API][hub-api]. OATs don't work with Docker
Desktop or Image Access Management.

An OAT can be limited to specific repositories and operations, and it
has its own Docker Hub usage limits, separate from individual accounts.
To see scopes, see [Access token reference][reference].

## OIDC connections

If your automation runs in GitHub Actions, you don't need to store a token.
An [OIDC connection][oidc] lets a workflow exchange GitHub's short-lived
identity token for Docker access on each run, so there is no long-lived
credential to rotate, scope, or leak. Organization owners set up OIDC
connections, and they require a Docker Team or Docker Business
subscription.

Use an OAT instead when the automation runs outside GitHub Actions, or
when you need a credential that works with `docker login` from any
system.

## Next steps

{{< grid >}}

[pat]: /manuals/security/access-tokens/personal-access-tokens.md
[oat]: /manuals/security/access-tokens/organization-access-tokens.md
[reference]: /manuals/security/access-tokens/reference.md
[permissions]: /manuals/security/access-tokens/reference.md#access-permissions
[plans]: /manuals/subscription-billing/plans/docker.md
[hub-api]: /manuals/security/access-tokens/reference.md#hub-api-support
[iam]: /manuals/desktop/enterprise/hardened-desktop/image-access-management.md
[oidc]: /manuals/security/authentication/oidc-connections/_index.md
