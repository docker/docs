---
title: Organization overview
linkTitle: Organization
weight: 10
description: >
  Learn how Docker organization accounts relate to individual accounts, and how
  to manage teams, members, permissions, and settings.
keywords: organizations, admin, overview, manage teams, roles, members,
  permissions, organization settings, organization account, individual account,
  Docker ID, account types, owners, teams
grid:
  - title: Onboard your organization
    description: Learn how to onboard and secure your organization.
    icon: magnifying-glass-plus
    link: /admin/organization/setup/onboard
  - title: Manage members
    description: Learn how to manage members.
    icon: user-plus
    link: /admin/organization/manage/members/
  - title: Activity logs
    description: Learn how to audit the activities of your members.
    icon: document-text
    link: /admin/activity-logs/
  - title: Security
    description:
      Start here to manage security and access for your organization, including
      single sign-on, provisioning, and image and registry access management.
    icon: shield-check
    link: /enterprise/security/
---

A Docker organization is a collection of teams and repositories under
centralized management. Organization administrators group members and
assign repository access at scale.

## Organization structure

The following diagram shows how organizations relate to teams and members.

![Diagram showing how teams and members relate within a Docker
organization](/admin/images/org-structure.webp)

For how organizations fit into the broader company hierarchy, see
[administration overview](/manuals/admin/_index.md#company-and-organization-hierarchy).

## Individual and organization accounts

Docker has two primary account types:

- Individual accounts that are identified by a Docker ID.
- Organization accounts that are shared workspaces for teams and
  repositories.

Every organization is created and administered by one or more individual
accounts. You always sign in with your individual account, then work in the
organizations you own or belong to. Organization owners and members are
individual accounts that hold a role in that organization. For individual
accounts, see [Accounts](/manuals/accounts/_index.md).

## Organization roles

An organization includes owners, members, and optional teams. Organization
owners have full administrator access to manage members, roles, and teams. A
team is an optional grouping of members that share the same repository
permissions.

For details about each role and its permissions, see
[Roles and permissions](/manuals/enterprise/security/roles-and-permissions/_index.md).

## What's next

Learn how to create and manage your organization in the following sections.

{{< grid >}}
