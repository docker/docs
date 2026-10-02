---
title: Set up a Docker organization
linkTitle: Setup
weight: 10
description: >
  Create, convert, or onboard a Docker organization under one namespace and
  subscription.
keywords:
  - Docker organization setup
  - create organization
  - onboard organization
  - convert Docker account
  - organization namespace
  - Docker Home
grid:
  - title: Create your organization
    description: Choose a new namespace and subscription.
    icon: building-storefront
    link: /accounts/organization/setup/orgs/
  - title: Convert your account
    description: Keep an existing Docker ID as the organization namespace.
    icon: arrows-right-left
    link: /accounts/organization/setup/convert-account/
  - title: Onboard your organization
    description: Invite members and configure sign-in.
    icon: magnifying-glass-plus
    link: /accounts/organization/setup/onboard/
  - title: Manage your organization
    description: Add members, teams, licenses, and seats after setup.
    icon: user-group
    link: /accounts/organization/manage/
  - title: Security
    description: Configure single sign-on, provisioning, and access management.
    icon: shield-check
    link: /security/
aliases:
  - /admin/organization/setup/
---

An organization groups members and teams under one namespace and one
subscription. Anyone with a [Docker ID](/manuals/accounts/_index.md) can
create an organization or convert an individual account into one.

You start by creating a new organization or converting an individual
account. After creating or converting, you can onboard your organization.

## Names versus namespaces

When you create an organization, you set two values:

- Organization namespace is the permanent, unique identifier for your
  organization. It becomes the first part of every image name you push, as
  in `namespace/image:tag`. You can't change it after you create the
  organization.
  - Docker IDs and organization namespaces must be unique.
  - If a Docker ID is `acme`, no organization can use `acme` as its
    namespace.
- Organization name is the display name shown on your organization's
  Docker profile. You can change it at any time. See
  [Change organization information](/manuals/accounts/organization/manage/general-settings.md).

## Choose how to set up

The difference between creating and converting is what happens to your existing repositories.

- Create an organization: Choose a new namespace. Your existing repositories stay under your personal Docker ID.
- Convert your account: Your Docker ID becomes the organization’s namespace. Your repositories and image names stay the same, so anyone pulling your images can keep using their existing image references.

## Next steps

{{< grid >}}
