---
title: Authentication
linkTitle: Authentication
description: Configure single sign-on, OIDC connections, and two-factor authentication.
keywords: authentication, SSO, OIDC, two-factor authentication, 2FA, Docker security
weight: 20
aliases:
  - /platform/security/authentication/
grid:
  - title: Single sign-on
    description: Authenticate users through your identity provider.
    icon: key
    link: /security/authentication/single-sign-on/
  - title: Two-factor authentication
    description: Add a TOTP security code to an individual Docker account.
    icon: device-phone-mobile
    link: /security/authentication/2fa/
  - title: OIDC connections
    description: Authenticate GitHub Actions with short-lived tokens.
    icon: lock-closed
    link: /security/authentication/oidc-connections/
---

Authentication in Docker Home is how users and workloads prove who they are
before they access Docker products.

Two-factor authentication (2FA) protects an individual account. Single
sign-on (SSO) federates sign-in for an organization or company. OpenID
Connect (OIDC) connections authenticate CI workloads such as GitHub Actions.

## Choose an authentication method

| Method | Who it covers | Who configures it | How authentication works |
| --- | --- | --- | --- |
| Two-factor authentication (2FA) | An individual Docker account | The account holder | Password plus a time-based one-time password (TOTP) from an authenticator app |
| Single sign-on (SSO) | An organization or company | An organization or company owner | Users sign in through the organization's identity provider (IdP) |
| OIDC connections | GitHub Actions and similar workloads | An organization owner or editor | Docker exchanges short-lived tokens issued per workflow run |

SSO requires a Docker Business subscription. OIDC connections require a
Docker Team or Business subscription.

To require Docker Desktop users to sign in as organization members, see
[Enforce sign-in](/manuals/enterprise/security/enforce-sign-in/_index.md).
Enforce sign-in is configured in Enterprise, not in this section.

## Next steps

{{< grid >}}
