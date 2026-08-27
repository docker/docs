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
  - title: Company administration
    description: Explore how to manage a company.
    icon: building-office-2
    link: /accounts/organization/company/
  - title: Security
    description: Explore security features for administrators.
    icon: shield-check
    link: /enterprise/security/
aliases:
  - /admin/
  - /docker-hub/admin-overview
  - /admin/organization/
  - /accounts/organization/overview/
---

Organization and company owners can manage members, control access, and enforce
security across their Docker environments. You perform these tasks in Docker
Home, which provides centralized observability, access management, and security
controls.

A Docker organization is a collection of teams and repositories under
centralized management. Organization administrators group members and
assign repository access at scale.

As an organization or company owner, you can:

- Create and manage companies and organizations
- Assign roles and permissions to members
- Group members into teams to manage access by project or role
- Set company-wide policies, including SCIM provisioning and security
  enforcement

## Individual and organization accounts

Docker has two primary account types:

- Individual accounts that are identified by a Docker ID.
- Organization accounts that are shared workspaces for teams and
  repositories.

Every organization is created and administered by one or more individual
accounts. You always sign in with your individual account, then work in the
organizations you own or belong to. Organization owners and members are
individual accounts that hold a role in that organization. For individual
accounts, see [Docker individual accounts](/manuals/accounts/individual/_index.md).

## Organization structure

The following diagram shows how organizations relate to teams and members.

![Diagram showing how teams and members relate within a Docker
organization](./images/org-structure.webp)

An organization includes owners, members, and optional teams. Organization
owners have full administrator access to manage members, roles, and teams. A
team is an optional grouping of members that share the same repository
permissions.

For details about each role and its permissions, see
[Roles and
permissions](/manuals/enterprise/security/roles-and-permissions/_index.md).

## Company and organization hierarchy

To provide centralized administration, Docker organizes companies and
organizations into the following hierarchy and roles.

![Diagram showing Docker’s administration hierarchy with Company at the top, followed by Organizations, Teams, and Members](./images/docker-admin-structure.webp)

### Company

A company groups multiple Docker organizations for centralized configuration.
Companies are only available for Docker Business subscribers. For company
structure, owners, and seats, see
[Company overview](/manuals/accounts/organization/company/_index.md).

### Organization

An organization sits below the company and is where you group teams and
members and assign access to repositories. Every Docker Team and Business
subscriber has at least one organization.

Organization owners hold the organization owner administrator role and manage
organization settings, users, and access controls. Each owner occupies a
[seat](/manuals/faqs/organization-faqs.md#what-is-the-difference-between-user-invitee-seat-and-member).

[Upgrading to a Docker Business plan](https://www.docker.com/pricing?ref=Docs&refAction=DocsAdmin)
grants you the company owner role so you can manage multiple organizations.

### Team

Teams are optional and let you group members to assign repository permissions
collectively. Teams simplify permission management across projects
or functions.

### Member

A member is any Docker user added to an organization. Organization and company
owners can assign roles to members to define their level of access.

## Next steps

Learn how to manage companies and organizations in the following sections.

{{< grid >}}
