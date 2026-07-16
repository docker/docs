---
title: Docker accounts
description: >
  Learn what a Docker account is, how Docker IDs, emails, and sign-in methods
  relate, and how individual accounts differ from organizations
keywords: accounts, docker ID, username, email, Google, GitHub, sign-in, account
  management, docker account, individual account, organization account
weight: 10
params:
  sidebar:
    group: Platform
grid:
  - title: Create an account
    description: Get started with Docker and create an account.
    icon: finger-print
    link: /accounts/create-account/
  - title: Manage account
    description: Learn how to manage the settings for your account.
    icon: cog
    link: /accounts/manage-account/
  - title: Personal access tokens
    description: Learn how to create and manage access tokens for your account.
    icon: lock-closed
    link: /security/access-tokens/
  - title: Set up two-factor authentication
    description: Add an extra layer of authentication to your Docker account.
    link: /security/2fa/
    icon: device-phone-mobile
  - title: Deactivate an account
    description: Learn how to deactivate a Docker user account.
    link: /accounts/deactivate-user-account/
    icon: no-symbol
  - title: Account FAQ
    description: Explore frequently asked questions about Docker accounts.
    icon: question-mark-circle
    link: /accounts/general-faqs/
---

A Docker account is your identity on Docker. With an account, you can access
Docker products and services such as Docker Hub and Docker Desktop, manage
security settings, and join organizations.

This section covers individual Docker accounts. For organizations, companies,
and administrator roles, see [Administration](/manuals/admin/_index.md).

## Account identity

When you create a Docker account, you choose a Docker ID and a sign-in method.
Docker also ties a verified email address to the account. These pieces work
together, but they aren't the same thing:

- Your Docker account is where Docker associates your plans, Hub
  repositories, and account settings.
- Your Docker ID identifies you with a unique username.
- You authenticate with email and password, Google, or GitHub.

You can sign in with your Docker ID, email, or connected social account.
Docker still associates a verified email from that provider with your account.
You choose one sign-in method when you create your account. Docker doesn't
support linking multiple sign-in methods to the same Docker ID.

For how to create an account, see
[Create a Docker account](/manuals/accounts/create-account.md).

## Individual accounts and organization accounts

Docker has two primary account types:

- Individual accounts that are identified by a Docker ID.
- Organization accounts that are shared workspaces for teams and
  repositories.

Every organization is created and administered by one or more individual
accounts. You always sign in with your individual account, then work in the
organizations you own or belong to. Organization owners and members are
individual accounts that hold a role in that organization.

To create and manage organizations, see
[Organization overview](/manuals/admin/organization/_index.md).

## What's next

{{< grid >}}
