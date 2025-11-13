---
title: Onboard your organization
weight: 20
description: Get started onboarding your Docker Team or Business organization.
keywords: business, team, organizations, get started, onboarding, Admin Console, organization management,
toc_min: 1
toc_max: 3
aliases:
- /docker-hub/onboard/
- /docker-hub/onboard-team/
- /docker-hub/onboard-business/
---

{{< summary-bar feature_name="Admin orgs" >}}

Learn how to onboard your organization using the Admin Console or Docker Hub.

Onboarding your organization includes:

- Identifying users to help you allocate your subscription seats
- Invite members and owners to your organization
- Secure authentication and authorization for your organization
- Enforce sign-in for Docker Desktop to ensure security best practices

These actions help administrators gain visibility into user activity and
enforce security settings. Organization members also receive increased pull
limits and other benefits when they are signed in.

## Prerequisites

Before you start onboarding your organization, ensure you:

- Have a Docker Team or Business subscription. For more details, see
[Docker subscriptions and features](/manuals/subscription/details.md).

  > [!NOTE]
  >
  > When purchasing a self-serve subscription, the on-screen instructions
  guide you through creating an organization. If you have purchased a
  subscription through Docker Sales and you have not yet created an
  organization, see [Create an organization](/manuals/admin/organization/orgs.md).

- Familiarize yourself with Docker concepts and terminology in
the [administration overview](../_index.md).

## Onboard with guided setup

The Admin Console has a guided setup to help you
onboard your organization. The guided setup's steps consist of basic onboarding
tasks. If you want to onboard outside of the guided setup,
see [Recommended onboarding steps](/manuals/admin/organization/onboard.md#recommended-onboarding-steps).

To onboard using the guided setup,
navigate to the [Admin Console](https://app.docker.com) and
select **Guided setup** in the left-hand navigation.

The guided setup walks you through the following onboarding steps:

- **Invite your team**: Invite owners and members.
- **Manage user access**: Add and verify a domain, manage users with SSO, and
enforce Docker Desktop sign-in.
- **Docker Desktop security**: Configure image access management, registry
access management, and settings management.

## Recommended onboarding steps

### Step one: Identify your Docker users

Identifying your users helps you allocate seats efficiently and ensures they
receive your Docker subscription benefits.

1. Identify the Docker users in your organization.
   - If your organization uses device management software, like MDM or Jamf,
   you can use the device management software to help identify Docker users.
   See your device management software's documentation for details. You can
   identify Docker users by checking if Docker Desktop is installed at the
   following location on each user's machine:
      - Mac: `/Applications/Docker.app`
      - Windows: `C:\Program Files\Docker\Docker`
      - Linux: `/opt/docker-desktop`
   - If your organization doesn't use device management software or your
   users haven't installed Docker Desktop yet, you can survey your users to
   identify who is using Docker Desktop.
1. Ask users to update their Docker account's email address to one associated
with your organization's domain, or create a new account with that email.
   - To update an account's email address, instruct your users to sign in
   to [Docker Hub](https://hub.docker.com), and update the email address to
   their email address in your organization's domain.
   - To create a new account, instruct your users to
   [sign up](https://hub.docker.com/signup) using their email address associated
   with your organization's domain.
1. Identify Docker accounts associated with your organization's domain:
   - Ask your Docker sales representative or
   [contact sales](https://www.docker.com/pricing/contact-sales/) to get a list
   of Docker accounts that use an email address in your organization's domain.

### Step two: Invite owners

Owners can help you onboard and manage your organization.

When you create an organization, you are the only owner. It is optional to
add additional owners.

To add an owner, invite a user and assign them the owner role. For more
details, see [Invite members](/manuals/admin/organization/members.md) and
[Roles and permissions](/manuals/enterprise/security/roles-and-permissions.md).

### Step three: Invite members

When you add users to your organization, you gain visibility into their
activity and you can enforce security settings. Your members also
receive increased pull limits and other organization wide benefits when
they are signed in.

To add a member, invite a user and assign them the member role.
For more details, see [Invite members](/manuals/admin/organization/members.md) and
[Roles and permissions](/manuals/enterprise/security/roles-and-permissions.md).

### Step four: Manage user access with SSO and SCIM

Configuring SSO and SCIM is optional and only available to Docker Business
subscribers. To upgrade a Docker Team subscription to a Docker Business
subscription, see [Change your subscription](/manuals/subscription/change.md).

Use your identity provider (IdP) to manage members and provision them to Docker
automatically via SSO and SCIM. See the following for more details:

   - [Configure SSO](/manuals/enterprise/security/single-sign-on/configure.md)
   to authenticate and add members when they sign in to Docker through your
   identity provider.
   - Optional.
   [Enforce SSO](/manuals/enterprise/security/single-sign-on/connect.md) to
   ensure that when users sign in to Docker, they must use SSO.

     > [!NOTE]
     >
     > Enforcing single sign-on (SSO) and enforcing Docker Desktop sign in
     are different features. For more details, see
     > [Enforcing sign-in versus enforcing single sign-on (SSO)](/manuals/enterprise/security/enforce-sign-in/_index.md#enforcing-sign-in-versus-enforcing-single-sign-on-sso).

   - [Configure SCIM](/manuals/enterprise/security/provisioning/scim.md) to
   automatically provision, add, and de-provision members to Docker through
   your identity provider.

### Step five: Enforce sign-in for Docker Desktop

By default, members of your organization can use Docker Desktop without signing
in. When users don’t sign in as a member of your organization, they don’t
receive the
[benefits of your organization’s subscription](../../subscription/details.md)
and they can circumvent [Docker’s security features](/manuals/enterprise/security/hardened-desktop/_index.md).

There are multiple ways you can enforce sign-in, depending on your organization's
Docker configuration:
- [Registry key method (Windows only)](/manuals/enterprise/security/enforce-sign-in/methods.md#registry-key-method-windows-only)
- [`.plist` method (Mac only)](/manuals/enterprise/security/enforce-sign-in/methods.md#plist-method-mac-only)
- [`registry.json` method (All)](/manuals/enterprise/security/enforce-sign-in/methods.md#registryjson-method-all)

### Step six: Manage Docker Desktop security

Docker offers the following security features to manage your organization's
security posture:

- [Image Access Management](/manuals/enterprise/security/hardened-desktop/image-access-management.md): Control which types of images your developers can pull from Docker Hub.
- [Registry Access Management](/manuals/enterprise/security/hardened-desktop/registry-access-management.md): Define which registries your developers can access.
- [Settings management](/manuals/enterprise/security/hardened-desktop/settings-management.md): Set and control Docker Desktop settings for your users.

## What's next

- [Manage Docker products](./manage-products.md) to configure access and view usage.
- Configure [Hardened Docker Desktop](/desktop/hardened-desktop/) to improve your organization’s security posture for containerized development.
- [Manage your domains](/manuals/enterprise/security/domain-management.md) to ensure that all Docker users in your domain are part of your organization.

Your Docker subscription provides many more additional features. To learn more,
see [Docker subscriptions and features](/subscription/details/).
