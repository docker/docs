---
title: Security
linkTitle: Security
description: >
  Secure Docker accounts, manage access, and control membership for
  individuals and organizations in Docker Home.
keywords: docker, docker hub, security, 2FA, access tokens, SSO, OIDC,
  provisioning, roles, Docker Home
weight: 40
aliases:
  - /security/for-developers/
  - /platform/security/
params:
  sidebar:
    group: Accounts and admin
grid:
  - title: Authentication
    description: Two-factor authentication, single sign-on, and OIDC connections.
    icon: key
    link: /security/authentication/
  - title: Access tokens
    description: Personal and organization access tokens for the Docker CLI and automation.
    icon: lock-closed
    link: /security/access-tokens/
  - title: Provisioning
    description: Add users with SCIM, JIT, auto-provisioning, and domain management.
    icon: arrow-path
    link: /security/provisioning/
  - title: Roles and permissions
    description: Assign core or custom roles to control access in your organization.
    icon: shield-check
    link: /security/roles-and-permissions/
---

Security helps individual users and organization owners secure their
accounts, manage access, and control membership. You configure these
settings in [Docker Home](https://app.docker.com/).

## Individual accounts

You sign in with your individual account.

- [Two-factor authentication](/manuals/security/authentication/2fa/_index.md)
(2FA) adds a time-based one-time password (TOTP) from an authenticator
app to your password.
- A [personal access token](/manuals/security/access-tokens/personal-access-tokens.md)
(PAT) authenticates the Docker CLI and tools without your password, and
is required for CLI sign-in when 2FA is on or single sign-on (SSO) is
enforced.

## Organization accounts

Organization and company owners set up how members sign in, add them to
the organization, configure automation, and control what members can do.

- [Single sign-on](/manuals/security/authentication/single-sign-on/_index.md)
(SSO) federates sign-in through your identity provider, which can cover
one organization or every organization in a company.
- [Provisioning](/manuals/security/provisioning/_index.md) adds users with
System for Cross-domain Identity Management (SCIM), Just-in-Time (JIT)
provisioning, auto-provisioning, or domain matching.
- An [organization access token](/manuals/security/access-tokens/organization-access-tokens.md)
(OAT) stays with the organization when membership changes.
- [OIDC connections](/manuals/security/authentication/oidc-connections/_index.md)
use OpenID Connect to authenticate GitHub Actions with short-lived
tokens, as an alternative to a long-lived OAT.
- [Roles and permissions](/manuals/security/roles-and-permissions/_index.md)
control what members can do after they join.

## Next steps

{{< grid >}}
