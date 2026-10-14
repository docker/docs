---
title: Onboard a Docker organization
linkTitle: Onboard
weight: 30
description: >
  Invite members, configure SSO and SCIM, and enforce Docker Desktop
  sign-in for a Team or Business organization.
keywords:
  - onboard Docker organization
  - guided setup
  - Docker Team
  - Docker Business
  - enforce sign-in
  - SSO
  - SCIM
  - Docker Desktop
toc_max: 2
aliases:
  - /admin/organization/setup/onboard/
  - /docker-hub/onboard/
  - /docker-hub/onboard-team/
  - /docker-hub/onboard-business/
  - /admin/organization/onboard/
---

{{< summary-bar feature_name="Admin orgs" >}}

Onboarding your organization includes:

- Identifying users so you can allocate subscription seats
- Inviting members and owners
- Securing authentication and authorization
- Enforcing sign-in for Docker Desktop

These actions give administrators visibility into user activity and a way
to enforce security settings. Organization members also receive higher
pull limits and other benefits when they are signed in.

## Prerequisites

Before you onboard your organization, you need a Docker Team or Business
subscription. For details, see
[Docker subscriptions and features](https://www.docker.com/pricing?ref=Docs&refAction=DocsAdminOnboard).
When you buy a self-serve subscription, the on-screen instructions guide
you through creating an organization.

If you bought a subscription through Docker Sales and you haven't created
an organization yet, see
[Create an organization](/manuals/accounts/organization/setup/orgs.md).

## Guided onboarding

You can use a guided setup for the first onboarding tasks:

1. Sign in to [Docker Home](https://app.docker.com).
1. Select **Guided setup** at the bottom of the left sidebar.

Guided setup walks through these steps:

- **Invite your team**: Invite owners and members.
- **Manage user access**: Claim your company's domain, manage members
  with SSO, and enforce Docker Desktop sign-in.
- **Docker Desktop security**: Configure Image Access Management,
  Registry Access Management, and Settings Management.

## Manual onboarding

### Step one: Identify your Docker users

Identifying your users helps you allocate seats and makes sure they
receive your Docker subscription benefits.

1. Identify the Docker users in your organization.
   - If your organization uses device management software, such as MDM
     or Jamf, use it to find machines with Docker Desktop installed:
     - Mac: `/Applications/Docker.app`
     - Windows: `C:\Program Files\Docker\Docker` (all-users
       installation) or `%LOCALAPPDATA%\Programs\DockerDesktop`
       (per-user installation)
     - Linux: `/opt/docker-desktop`
   - If your organization doesn't use device management software, or
     users haven't installed Docker Desktop yet, ask them who uses
     Docker Desktop.
1. Ask users to update their Docker account email address to one in your
   organization's domain, or to create an account with that email.
   - To update an email address, see
     [Update email address](/manuals/accounts/individual/manage-account.md#update-email-address).
   - To create an account, users
     [sign up](https://hub.docker.com/signup) with an email address in
     your organization's domain and verify that address.
1. Identify Docker accounts that already use your organization's domain:
   - Ask your Docker sales representative, or
     <a href="https://www.docker.com/pricing/contact-sales/" id="dkr_docs_cs_org_onboarding" class="link" rel="noopener">contact sales</a>,
     for a list of Docker accounts that use an email address in your
     domain.

### Step two: Invite owners

When you create an organization, you are the only owner. Adding more
owners is optional, but additional owners can help you onboard and manage
your organization.

To add an owner, invite a user and assign the owner role. For details,
see [Invite members](/manuals/accounts/organization/manage/members.md)
and
[Roles and permissions](/manuals/security/roles-and-permissions/_index.md).

### Step three: Invite members

When you add users to your organization, you can see their activity and
enforce security settings. Members also receive higher pull limits and
other organization-wide benefits when they are signed in.

To add a member, invite a user and assign the member role. For details,
see [Invite members](/manuals/accounts/organization/manage/members.md)
and
[Roles and permissions](/manuals/security/roles-and-permissions/_index.md).

### Step four: Manage user access with SSO and SCIM

SSO and SCIM are optional and available to Docker Business subscribers.
To upgrade a Docker Team subscription to Docker Business, see
[Upgrade a plan](/manuals/subscription-billing/manage/plans.md#upgrade-plans).

Use your identity provider (IdP) to manage members and provision them to
Docker through SSO and SCIM:

- [Configure SSO](/manuals/security/authentication/single-sign-on/connect.md)
  to authenticate members and add them when they sign in through your
  identity provider.
- Optional.
  [Enforce SSO](/manuals/security/authentication/single-sign-on/connect.md)
  so users must use SSO when they sign in to Docker.

  > [!NOTE]
  >
  > Enforcing single sign-on (SSO) and enforcing Docker Desktop sign-in
  > are different features. For details, see
  > [Enforcing sign-in versus enforcing single sign-on (SSO)](/manuals/desktop/enterprise/enforce-sign-in/_index.md#enforcing-sign-in-versus-enforcing-single-sign-on-sso).

- [Configure SCIM](/manuals/security/provisioning/scim/_index.md) to
  provision, add, and deprovision members through your identity
  provider.

### Step five: Enforce sign-in for Docker Desktop

By default, members of your organization can use Docker Desktop without
signing in. When users don't sign in as a member of your organization,
they don't receive the
[benefits of your organization's subscription](https://www.docker.com/pricing?ref=Docs&refAction=DocsAdminOnboard)
and they can bypass
[Docker's security features](/manuals/desktop/enterprise/hardened-desktop/_index.md).

You can enforce sign-in in more than one way:

- [Registry key method (Windows only)](/manuals/desktop/enterprise/enforce-sign-in/methods.md#registry-key-method-windows-only)
- [`.plist` method (Mac only)](/manuals/desktop/enterprise/enforce-sign-in/methods.md#plist-method-mac-only)
- [`registry.json` method (All)](/manuals/desktop/enterprise/enforce-sign-in/methods.md#registryjson-method-all)

### Step six: Manage Docker Desktop security

Use these features to manage your organization's security posture:

- [Image Access Management](/manuals/desktop/enterprise/hardened-desktop/image-access-management.md):
  Control which types of images developers can pull from Docker Hub.
- [Registry Access Management](/manuals/desktop/enterprise/hardened-desktop/registry-access-management.md):
  Define which registries developers can access.
- [Settings Management](/manuals/desktop/enterprise/hardened-desktop/settings-management.md):
  Set and control Docker Desktop settings for your users.

## Next steps

- [Manage Docker products](../manage/manage-products.md) to configure
  access and view usage.
- Configure
  [Hardened Docker Desktop](/manuals/desktop/enterprise/hardened-desktop/_index.md)
  to tighten security for containerized development.
- [Manage your domains](/manuals/security/provisioning/domain-management.md)
  so Docker users in your domain are part of your organization.
