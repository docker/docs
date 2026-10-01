---
title: Organization accounts
linkTitle: Organization
description: Overview of administration features and roles in Docker Home
keywords: admin, administration, company, organization, Docker Home, user
  accounts, account management, organizations, manage teams, roles, members,
  permissions, organization settings, organization account, individual account,
  Docker ID, account types, owners, teams
weight: 15
grid:
  - title: Set up your organization
    description: Create, onboard, and configure your organization.
    icon: magnifying-glass-plus
    link: /accounts/organization/setup/
  - title: Manage your organization
    description: Manage members, teams, seats, and product access.
    icon: user-plus
    link: /accounts/organization/manage/
  - title: Activity logs
    description: Review member activity across your organization and repositories.
    icon: clipboard-document-list
    link: /accounts/organization/activity-logs/
  - title: Insights
    description: See how people in your organization use Docker.
    icon: chart-bar
    link: /accounts/organization/insights/
  - title: Security
    description: Explore security features for administrators.
    icon: shield-check
    link: /security/
aliases:
  - /admin/
  - /docker-hub/admin-overview
  - /admin/organization/
  - /accounts/organization/overview/
---

A Docker organization is a collection of teams and repositories that you
manage in [Docker Home](https://app.docker.com/). Organization and company
owners manage members, assign access, and enforce security.

For how individual, organization, and company accounts compare, see
[Accounts](/manuals/accounts/_index.md). For individual accounts, see
[Docker individual accounts](/manuals/accounts/individual/_index.md).
To create, convert, or onboard an organization, see
[Set up a Docker organization](/manuals/accounts/organization/setup/_index.md).

## Organization structure

The following diagram shows how organizations relate to teams and members.

![Diagram showing how teams and members relate within a Docker
organization](./images/org-structure.webp)

An organization includes owners, members, and optional teams. Organization
owners have full administrator access to manage members, roles, and teams.

### Team

Teams are optional and let you group members to assign repository permissions
collectively. Teams simplify permission management across projects
or functions.

### Member

A member is any Docker user added to an organization. Organization and company
owners can assign roles to members to define their level of access.

For details about each role and its permissions, see
[Roles and
permissions](/manuals/security/roles-and-permissions/_index.md).

For how companies relate to organizations, see
[Company structure](/manuals/accounts/company/_index.md#company-structure).

## Next steps

Learn how to manage organizations in the following sections.

{{< grid >}}
