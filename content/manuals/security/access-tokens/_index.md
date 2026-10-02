---
title: Access tokens
linkTitle: Access tokens
description: Create and manage personal and organization access tokens for Docker Hub authentication.
keywords: access tokens, personal access tokens, organization access tokens, PAT, OAT, Docker security
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
---

Access tokens let you authenticate to Docker Hub without using your password.
Use a token for the Docker CLI, automation, and any account that has
two-factor authentication (2FA) or enforced single sign-on (SSO), because
password sign-in to the CLI is not supported in those cases.

## How access tokens work

Both token types share the same basic behavior:

- A token stands in for your password. Enter it at the password prompt for
  `docker login`, or wherever a tool asks for your Docker Hub password.
- Each token has its own permissions, so it can only do what you allowed when
  you created it. A leaked token can't do more than that.
- Docker shows the token value once, when you create it. Docker doesn't store
  the value, so copy it right away. If you lose it, create a new token.
- A token can have an expiration date. An expired token stops working but
  stays in your token list until you delete it.
- You can deactivate a token to stop it from working without deleting it.
  Deleting a token is permanent.

## Choose a token type

| Token | Ownership | Use when | Limitations |
| --- | --- | --- | --- |
| Personal access token (PAT) | Tied to an individual Docker account | CLI access, local tools, and automation that should run as you. Required for CLI sign-in when 2FA is on or SSO is enforced | Access ends if the account leaves the organization or the token is revoked |
| Organization access token (OAT) | Owned by the organization. Any organization owner can manage it | CI/CD and other automation that must keep working when membership changes | Incompatible with Docker Desktop and Image Access Management |

For GitHub Actions, [OIDC connections](/manuals/security/authentication/oidc-connections/_index.md)
are an alternative to storing a long-lived organization access token.

## Next steps

{{< grid >}}
