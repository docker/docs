---
description: >-
  Provision Docker organization users with SCIM, JIT, group mapping, or
  auto-provisioning, and map SSO and SAML attributes from your identity
  provider.
keywords: provision users, user provisioning, JIT, SCIM, group mapping,
  auto-provisioning, SSO, SAML, identity provider, dockerOrg, dockerRole,
  dockerTeam, dockerSessionMinutes, Docker Home, admin, security
title: User provisioning overview
linkTitle: Provision
weight: 30
aliases:
  - /security/for-admins/provisioning/
  - /enterprise/security/provisioning/
  - /platform/security/provisioning/
grid:
  - title: Add and manage domains
    description: Add, verify, and manage domains for auto-provisioning.
    icon: globe-alt
    link: "domain-management/"
  - title: SCIM provisioning
    description: Sync user data between your IdP and Docker with SCIM.
    icon: arrow-path
    link: "scim/"
  - title: Just-in-Time (JIT) provisioning
    description: Create user accounts automatically on first SSO sign-in.
    icon: clock
    link: "just-in-time/"
  - title: Auto-provisioning
    description: Add users whose email addresses match a verified domain.
    icon: user-group
    link: "auto-provisioning/"
---

{{< summary-bar feature_name="SSO" >}}

After you configure single sign-on (SSO), provision users so they can
access your organization through automated account management.

## Provisioning methods

Provisioning automates account creation, updates, and deactivation using
data from your identity provider (IdP). Docker supports the following
methods:

| Provisioning method | When it runs | Lifecycle management | Default setting |
| :--- | :--- | :--- | :--- |
| [System for Cross-domain Identity Management (SCIM)](/manuals/security/provisioning/scim/_index.md) | On the IdP's synchronization schedule or through Provision on Demand | Creates and updates users, synchronizes configured groups, and deprovisions users | Disabled |
| [Just-in-Time (JIT)](/manuals/security/provisioning/just-in-time.md) | When a user signs in through SSO | Creates users and applies attributes from the SSO assertion. It doesn't deprovision users | Enabled when you configure SSO |
| [Auto-provisioning](/manuals/security/provisioning/auto-provisioning.md) | When an existing Docker user signs in with an email address from a verified domain | Adds the user to the organization. It doesn't create or deprovision accounts | Disabled |

[Group mapping](/manuals/security/provisioning/scim/group-mapping.md) assigns
users to Docker organizations and teams. Use it with SAML SSO or SCIM. You can
also invite users manually when automatic provisioning isn't configured.

## Default provisioning setup

Docker turns on JIT provisioning when you configure an SSO connection. If you
also enable SCIM, Docker recommends choosing one provisioning source to manage
users and attributes. Before configuring SCIM, review
[how SCIM works with JIT](/manuals/security/provisioning/scim/_index.md#choose-how-scim-works-with-jit).

For a domain that belongs to an SSO connection, JIT adds the user instead of
auto-provisioning.

## SSO attributes

Each time a user signs in through SSO, Docker reads attributes from your
IdP to set the user's identity and permissions:

| Attribute | Required | Description |
| :--- | :--- | :--- |
| Email address | Yes | Unique identifier for the user |
| Full name | Yes | User's complete name |
| Groups | No | Group-based access control |
| Docker Org | No | Organization the user belongs to |
| Docker Team | No | Team within the organization |
| Docker Role | No | Permissions in Docker |
| Docker session minutes | No | Session duration, in minutes, before users must re-authenticate with their IdP. Must be a positive integer greater than 0. If omitted, default session timeouts apply |

> [!NOTE]
>
> Default session timeouts apply when Docker session minutes is not
> specified. Docker Desktop sessions expire after 90 days or 30 days of
> inactivity. Docker Hub and Docker Home sessions expire after 24 hours.

## SAML attribute mapping

If your organization uses SAML for SSO, Docker reads these attributes
from the SAML assertion. Identity providers may use different names for
the same attributes.

| SSO attribute | SAML assertion attributes |
| :--- | :--- |
| Email address | `"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/nameidentifier"`, `"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/upn"`, `"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/emailaddress"`, `email` |
| Full name | `"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/name"`, `name`, `"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/givenname"`, `"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/surname"` |
| Groups (optional) | `"http://schemas.xmlsoap.org/claims/Group"`, `"http://schemas.microsoft.com/ws/2008/06/identity/claims/groups"`, `Groups`, `groups` |
| Docker Org (optional) | `dockerOrg` |
| Docker Team (optional) | `dockerTeam` |
| Docker Role (optional) | `dockerRole` |
| Docker session minutes (optional) | `dockerSessionMinutes`, must be a positive integer greater than 0 |

## Next steps

Choose the provisioning method that fits your organization:

{{< grid >}}

If users get the wrong role or team after you change methods, see
[Troubleshoot provisioning](/manuals/security/provisioning/troubleshoot-provisioning.md).
