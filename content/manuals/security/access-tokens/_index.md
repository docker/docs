---
title: Access tokens
linkTitle: Access tokens
description: Create and manage personal and organization access tokens for Docker Hub authentication.
keywords: access tokens, personal access tokens, organization access tokens, PAT, OAT, Docker security
weight: 10
aliases:
  - /platform/security/access-tokens/
grid:
  - title: Personal access tokens
    description: Authenticate the Docker CLI and tools with a token tied to your account.
    icon: lock-closed
    link: /security/access-tokens/personal-access-tokens/
  - title: Organization access tokens
    description: Grant org-owned Hub access to CI/CD and other automation.
    icon: building-office-2
    link: /security/access-tokens/organization-access-tokens/
---

Access tokens let you authenticate to Docker Hub without using your password.
Use a token for the Docker CLI, automation, and any account that has
two-factor authentication (2FA) or enforced single sign-on (SSO), because
password sign-in to the CLI is not supported in those cases.

## Choose a token type

| Token | Ownership | Use when | Limitations |
| --- | --- | --- | --- |
| Personal access token (PAT) | Tied to an individual Docker account | CLI access, local tools, and automation that should run as you. Required for CLI sign-in when 2FA is on or SSO is enforced | Access ends if the account leaves the organization or the token is revoked |
| Organization access token (OAT) | Owned by the organization. Any organization owner can manage it | CI/CD and other automation that must keep working when membership changes | Incompatible with Docker Desktop and Image Access Management |

For GitHub Actions, [OIDC connections](/manuals/security/authentication/oidc-connections/_index.md)
are an alternative to storing a long-lived organization access token.

## Next steps

{{< grid >}}
