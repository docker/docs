---
description: >-
  Provision Docker organization users with SCIM, JIT, group mapping, or
  auto-provisioning, and map SSO and SAML attributes from your identity
  provider.
keywords: provision users, user provisioning, JIT, SCIM, group mapping,
  auto-provisioning, SSO, SAML, Azure AD, OIDC, identity provider, dockerOrg,
  dockerRole, dockerTeam, dockerSessionMinutes, Docker Home, admin, security
title: User provisioning overview
linkTitle: Provision
weight: 30
aliases:
  - /security/for-admins/provisioning/
  - /enterprise/security/provisioning/
---

{{< summary-bar feature_name="SSO" >}}

After you configure single sign-on (SSO), choose how Docker creates accounts
and adds users to your organization. Organization owners configure
provisioning. For a company, the company owner configures it.
SSO requires a Docker Business subscription.

## Choose a provisioning method

Provisioning creates and updates accounts from your identity provider (IdP). Docker supports these methods:

| Provisioning method | When it runs | Lifecycle management | Default setting |
| :--- | :--- | :--- | :--- |
| [System for Cross-domain Identity Management (SCIM)](/manuals/security/provisioning/scim/_index.md) | On the IdP's synchronization schedule or through Provision on Demand | Creates and updates users, synchronizes configured groups, and deprovisions users | Disabled |
| [Just-in-Time (JIT)](/manuals/security/provisioning/just-in-time.md) | When a user signs in through SSO | Creates users and applies attributes from the SSO assertion. It doesn't deprovision users | Enabled when you configure SSO |
| [Auto-provisioning](/manuals/security/provisioning/auto-provisioning.md) | When an existing Docker user signs in or verifies their email, and that address uses a verified domain | Adds the user to the organization. It doesn't create or deprovision accounts | Disabled |

After you choose a provisioning method, you can enable [group mapping](/manuals/security/provisioning/scim/group-mapping.md) to assign
users to one or more Docker organizations and teams. Group mapping works only
with SAML SSO or SCIM.

> [!IMPORTANT]
>
> Using JIT and auto-provisioning on the same domain can cause unexpected
> results. Avoid enabling both for a domain that's part of an SSO connection.

## Manage JIT and SCIM together

You must have at least one provisioning method enabled when you configure an SSO connection. By default Docker turns on JIT provisioning. If you want to turn JIT off, then you must have SCIM enabled; conversely, if you want to turn SCIM off, you must re-enable JIT. If you
enable SCIM and JIT together, choose which one manages users and attributes. Review
[how SCIM works with JIT](/manuals/security/provisioning/scim/_index.md#choose-how-scim-works-with-jit) before choosing.

## Set attribute names in the identity provider

When a user signs in through SSO, Docker identifies the account and keeps
its email address and full name current based on what the IdP sends. IdPs
use different claim names for those two attributes. In your identity provider,
set the full name and the email address each to a claim name Docker already
accepts. Docker uses the name you set to store the value on the signed-in
user's account. See
[SSO and SAML attributes](/manuals/security/authentication/single-sign-on/connect.md#sso-and-saml-attributes)
for the accepted claim names.

Organization, team, and role are not attributes you set for SSO sign-in, and
your identity provider does not include them. For SCIM, you create custom
attributes in the identity provider and set the external names to `dockerOrg`,
`dockerTeam`, and `dockerRole`. The IdP sends those names to Docker. Docker
recognizes each name, places the user in that organization and team, and
assigns the role. See
[Set up role mapping](/manuals/security/provisioning/scim/provision-scim.md#set-up-role-mapping)
for how to create the attributes.

## Next steps

Choose the provisioning method that fits your organization:

- [Add and manage domains](/manuals/security/provisioning/domain-management.md)
  to add, verify, and manage domains for auto-provisioning.
- [SCIM provisioning](/manuals/security/provisioning/scim/_index.md) to sync
  user data between your IdP and Docker.
- [Just-in-Time (JIT) provisioning](/manuals/security/provisioning/just-in-time.md)
  to create user accounts on first SSO sign-in.
- [Auto-provisioning](/manuals/security/provisioning/auto-provisioning.md) to
  add users whose email addresses match a verified domain.
- [Troubleshoot provisioning](/manuals/security/provisioning/troubleshoot-provisioning.md)
  if users get the wrong role or team after you change methods.
