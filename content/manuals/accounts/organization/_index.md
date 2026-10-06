---
title: Organization accounts
linkTitle: Organization
description: How Docker organizations relate to members, teams, and repositories.
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

A Docker organization is a shared workspace for members and repositories
under one namespace. You manage it in [Docker Home](https://app.docker.com/).
Organization owners administer membership, access, and security.

## Organization structure

Organization owners manage the organization, which contains
members, repositories, and teams. While teams are optional, they're another way to group your members within your organization.

The following diagram shows that hierarchy.

```mermaid
flowchart TB
  oo((Organization owners)) -.-> manage@{ shape: text, label: "manage" }
  manage -.-> org[Organization]
  org -.- t["Teams (optional)"]
  t -.- mt((Members))
  org --- r[(Repositories)]
  org --- m1((Members))
  org --- m2((Members))
  classDef optional stroke-dasharray: 5 5
  class t optional
```

### Owners

Organization owners administer the organization. They invite members, assign
roles, and manage teams and repositories.

You can have multiple organization owners per organization. You can have different and its permissions, see
[Roles and
permissions](/manuals/security/roles-and-permissions/_index.md).

### Members

A member is a Docker user added to the organization. Organization and company
owners assign a role to each member. That role sets organization-wide access.
Team permissions can grant extra access to specific repositories.

### Teams

Teams group members so you can grant repository access to many people at
once. Use a team when several members need the same repositories. Members
don't have to join a team.

For how companies relate to organizations, see
[Company structure](/manuals/accounts/company/_index.md#company-structure).

## Next steps

Learn how to manage organizations in the following sections.

{{< grid >}}
