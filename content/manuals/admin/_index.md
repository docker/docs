---
title: Administration
description: Overview of administration features and roles in Docker Home
keywords: admin, administration, company, organization, Docker Home, user accounts, account management
weight: 10
params:
  sidebar:
    group: Enterprise
grid:
  - title: Company administration
    description: Explore how to manage a company.
    icon: building-office-2
    link: /admin/company/
  - title: Organization administration
    description: Learn about organization administration.
    icon: building-storefront
    link: /admin/organization/
  - title: Onboard your organization
    description: Learn how to onboard and secure your organization.
    icon: magnifying-glass-plus
    link: /admin/organization/setup/onboard
  - title: Security
    description: Explore security features for administrators.
    icon: shield-check
    link: /enterprise/security/
aliases:
  - /docker-hub/admin-overview
---

Organization and company owners can manage members, control access, and enforce
security across their Docker environments. You perform these tasks in Docker
Home, which provides centralized observability, access management, and security
controls.

As an organization or company owner, you can:

- Create and manage companies and organizations
- Assign roles and permissions to members
- Group members into teams to manage access by project or role
- Set company-wide policies, including SCIM provisioning and security
  enforcement

## Company and organization hierarchy

To provide centralized administration, Docker organizes companies and
organizations into the following hierarchy and roles.

![Diagram showing Docker’s administration hierarchy with Company at the top, followed by Organizations, Teams, and Members](./images/docker-admin-structure.webp)

### Company

A company groups multiple Docker organizations for centralized configuration. A
company owner can view and manage every organization in the company and its
company-wide settings, with the same access rights as an organization owner. For
the company owner role and how it affects seats, see
[Company roles](/manuals/admin/company/_index.md#company-roles).

Companies are only available for Docker Business subscribers.

### Organization

An organization sits below the company and is where you group teams and
members and assign access to repositories. Every Docker Team and Business
subscriber has at least one organization.

Organization owners hold the organization owner administrator role and manage
organization settings, users, and access controls. Each owner occupies a
[seat](/manuals/admin/organization/organization-faqs.md#what-is-the-difference-between-user-invitee-seat-and-member).

[Upgrading to a Docker Business plan](https://www.docker.com/pricing?ref=Docs&refAction=DocsAdmin)
grants you the company owner role so you can manage multiple organizations.

### Team

Teams are optional and let you group members to assign repository permissions
collectively. Teams simplify permission management across projects
or functions.

### Member

A member is any Docker user added to an organization. Organization and company
owners can assign roles to members to define their level of access.

## What's next

Learn how to manage companies and organizations in the following sections.

{{< grid >}}
