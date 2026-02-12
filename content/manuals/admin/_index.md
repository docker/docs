---
title: Administration
description: Overview of administration features and roles in the Docker Admin Console
keywords: admin, administration, company, organization, Admin Console, user accounts, account management
weight: 10
params:
  sidebar:
    group: Enterprise
grid:
- title: Company administration
  description: Explore how to manage a company.
  icon: apartment
  link: /admin/company/
- title: Organization administration
  description: Learn about organization administration.
  icon: store
  link: /admin/organization/
- title: Onboard your organization
  description: Learn how to onboard and secure your organization.
  icon: explore
  link: /admin/organization/onboard
- title: Company FAQ
  description: Discover common questions and answers about companies.
  icon: help
  link: /faq/admin/company-faqs/
- title: Organization FAQ
  description: Explore popular FAQ topics about organizations.
  icon: help
  link: /faq/admin/organization-faqs/
- title: Security
  description: Explore security features for administrators.
  icon: shield_locked
  link: /security/
aliases:
- /docker-hub/admin-overview
---

Administrators can manage companies and organizations using the
[Docker Admin Console](https://app.docker.com/admin). The Admin Console
provides centralized observability, access management, and security controls
across Docker environments.

## Company and organization hierarchy

The [Docker Admin Console](https://app.docker.com/admin) provides administrators with centralized observability, access management, and controls for their company and organizations. To provide these features, Docker uses the following hierarchy and roles.

![Diagram showing Docker’s administration hierarchy with Company at the top, followed by Organizations, Teams, and Members](./images/docker-admin-structure.webp)

### Company

A company groups multiple Docker organizations for centralized configuration.
Companies are only available for Docker Business subscribers.

Companies has the company owner administrator role available: 

- Can view and manage all organizations within the company
- Has full access to company-wide settings and inherits the same permissions as organization owners
- Do not occupy a seat

### Organization

An organization contains teams and repositories. All Docker Team and Business
subscribers must have at least one organization.

Organization owners has the organization owner administrator role available. They can manage organization settings, users, and access controls, but occupy a [seat](/admin/faqs/organization-faqs/#what-is-the-difference-between-user-invitee-seat-and-member). If you want to [upgrade to a Docker Business plan](https://www.docker.com/pricing/), you can access the company owner role. 

While creating a company is optional, organizations are required for Team and Business subscriptions.

### Team

Teams are optional and let you group members to assign repository permissions
collectively. Teams simplify permission management across projects
or functions.

### Member

A member is any Docker user added to an organization. Organization and company
owners can assign roles to members to define their level of access.

## Admin Console features

Docker's [Admin Console](https://app.docker.com/admin) allows you to:

- Create and manage companies and organizations
- Assign roles and permissions to members
- Group members into teams to manage access by project or role
- Set company-wide policies, including SCIM provisioning and security
enforcement

## Manage companies and organizations

Learn how to manage companies and organizations in the following sections.

{{< grid >}}
