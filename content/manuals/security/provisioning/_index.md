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

| Provisioning method | Description | Default setting in Docker | Recommended for |
| :--- | :--- | :--- | :--- |
| System for Cross-domain Identity Management (SCIM) | Syncs user data between your IdP and Docker so attributes stay current without manual updates | Disabled by default | Large organizations or frequent changes in users or roles |
| Group mapping | Maps IdP groups to Docker roles and permissions based on group membership | Disabled by default | Organizations that assign access from IdP group membership |
| Just-in-Time (JIT) | Creates and provisions user accounts when they first sign in with SSO | Enabled by default | Organizations that need minimal setup or smaller teams |
| Auto-provision | Adds users whose email addresses match a verified domain | Disabled by default | Organizations without SSO that add existing Docker users by domain |

## Default provisioning setup

Docker turns on JIT provisioning when you configure an SSO connection.
With JIT on, Docker creates a user account the first time the user signs
in through SSO.

If you need more control over user access and attributes, configure SCIM
or group mapping.

## SSO attributes

When a user signs in through SSO, Docker reads attributes from your IdP
to set identity and permissions:

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
