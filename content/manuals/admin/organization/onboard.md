---
title: Onboard your organization
weight: 20
description: Get started onboarding your Docker Team or Business organization.
keywords: business, team, organizations, get started, onboarding
toc_min: 1
toc_max: 2
aliases:
- /docker-hub/onboard/
- /docker-hub/onboard-team/
- /docker-hub/onboard-business/
---

{{< summary-bar feature_name="Admin orgs" >}}

Learn how to onboard your organization using Docker Hub or the Docker Admin Console.

Onboarding your organization lets administrators gain visibility into user activity and enforce security settings. In addition, members of your organization receive increased pull limits and other organization wide benefits. For more details, see [Docker subscriptions and features](../../subscription/details.md).

In this guide, you'll learn how to do the following:

- Identify your users to help you efficiently allocate your subscription seats
- Invite members and owners to your organization
- Secure authentication and authorization for your organization using Single Sign-On (SSO) and System for Cross-domain Identity Management (SCIM)
- Enforce sign-on for Docker Desktop to ensure security best practices

## Prerequisites

Before you start to onboard your organization, ensure that you:

- Have a Docker Team or Business subscription. See [Docker Pricing](https://www.docker.com/pricing/) for details.

  > [!NOTE]
  >
  > When purchasing a self-serve subscription, the on-screen instructions guide you through creating an organization. If you have purchased a subscription through Docker Sales and you have not yet created an organization, see [Create an organization](/admin/organization/orgs).

- Familiarize yourself with Docker concepts and terminology in the [glossary](/glossary/) and [FAQs](/faq/admin/general-faqs/).

## Step 1: Identify your Docker users

Identifying your users will ensure that you allocate your subscription seats efficiently and that all your Docker users receive the benefits of your subscription.

1. Identify the Docker users in your organization.
   - If your organization uses device management software, like MDM or Jamf, you may use the device management software to help identify Docker users. See your device management software's documentation for details. You can identify Docker users by checking if Docker Desktop is installed at the following location on each user's machine:
      - Mac: `/Applications/Docker.app`
      - Windows: `C:\Program Files\Docker\Docker`
      - Linux: `/opt/docker-desktop`
   - If your organization doesn't use device management software or your users haven't installed Docker Desktop yet, you may survey your users.
2. Instruct all your organization's Docker users to update their existing Docker account's email address to an address that's in your organization's domain, or to create a new account using an email address in your organization's domain.
   - To update an account's email address, instruct your users to sign in to [Docker Hub](https://hub.docker.com), and update the email address to their email address in your organization's domain.
   - To create a new account, instruct your users to go [sign up](https://hub.docker.com/signup) using their email address in your organization's domain.
3. Ask your Docker sales representative or [contact sales](https://www.docker.com/pricing/contact-sales/) to get a list of Docker accounts that use an email address in your organization's domain.

## Step 2: Invite owners

When you create an organization, you are the only owner. It is optional to add additional owners. Owners can help you onboard and manage your organization.

To add an owner, invite a user and assign them the owner role. For more details, see [Invite members](/admin/organization/members/).

## Step 3: Invite members

When you add users to your organization, you gain visibility into their activity and you can enforce security settings. In addition, members of your organization receive increased pull limits and other organization wide benefits.

To add a member, invite a user and assign them the member role. For more details, see [Invite members](/admin/organization/members/).

## Step 4: Manage members with SSO and SCIM

Configuring SSO and SCIM is optional and only available to Docker Business subscribers. To upgrade a Docker Team subscription to a Docker Business subscription, see [Upgrade your subscription](/subscription/upgrade/).

You can manage your members in your identity provider and automatically provision them to your Docker organization with SSO and SCIM. See the following for more details.
   - [Configure SSO](/manuals/security/for-admins/single-sign-on/configure.md) to authenticate and add members when they sign in to Docker through your identity provider.
   - Optional. [Enforce SSO](/manuals/security/for-admins/single-sign-on/connect.md) to ensure that when users sign in to Docker, they must use SSO.

     > [!NOTE]
     >
     > Enforcing single sign-on (SSO) and [Step 5: Enforce sign-in for Docker
     > Desktop](#step-5-enforce-sign-in-for-docker-desktop) are different
     > features. For more details, see
     > [Enforcing sign-in versus enforcing single sign-on (SSO)](/security/for-admins/enforce-sign-in/#enforcing-sign-in-versus-enforcing-single-sign-on-sso).

   - [Configure SCIM](/security/for-admins/provisioning/scim/) to automatically provision, add, and de-provision members to Docker through your identity provider.

## Step 5: Enforce sign-in for Docker Desktop

By default, members of your organization can use Docker Desktop without signing
in. When users don’t sign in as a member of your organization, they don’t
receive the [benefits of your organization’s subscription](../../subscription/details.md) and they can circumvent [Docker’s security features](/security/for-admins/hardened-desktop/).

There are multiple ways you can enforce sign-in, depending on your company's setup and preferences:
- [Registry key method (Windows only)](/security/for-admins/enforce-sign-in/methods/#registry-key-method-windows-only)
- [`.plist` method (Mac only)](/security/for-admins/enforce-sign-in/methods/#plist-method-mac-only)
- [`registry.json` method (All)](/security/for-admins/enforce-sign-in/methods/#registryjson-method-all)

## What's next

- [Manage Docker products](./manage-products.md) to configure access and view usage.
- Configure [Hardened Docker Desktop](/desktop/hardened-desktop/) to improve your organization’s security posture for containerized development.
- [Audit your domains](/docker-hub/domain-audit/) to ensure that all Docker users in your domain are part of your organization.

Your Docker subscription provides many more additional features. To learn more, see [Docker subscriptions and features](/subscription/details/).
