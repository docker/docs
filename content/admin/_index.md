---
description: Discover manuals on administration for accounts, organizations, and companies.
keywords: admin, administration, company, organization, Docker Admin, user accounts, account management
title: Administration overview
grid:
- title: Company administration
  description: Explore how to manage a company.
  icon: apartment
  link: /admin/company/
- title: Organization administration
  description: Learn about organization administration.
  icon: store
  link: /admin/organization/
- title: Company FAQ
  description: Discover common questions and answers about companies.
  icon: help
  link: /docker-hub/company-faqs/
- title: Organization FAQ
  description: Explore popular FAQ topics about organizations.
  icon: help
  link: /docker-hub/organization-faqs/
---

Administrators can manage companies and organizations using Docker Hub or Docker Admin (Early Access). Docker Admin is available in Early Access to all company owners and organization owners that have a Docker Business or Docker Team subscription.

The [Docker Admin](https://admin.docker.com) console provides administrators with centralized observability, access management, and controls for their company and organizations. To provide these features, Docker uses the following hierarchy and roles.

![Docker hierarchy](./images/docker-admin-structure.webp)

- Company: A company simplifies the management of Docker organizations and settings. Creating a company is optional and only available to Docker Business subscribers.
  - Company owner: A company can have multiple owners. Company owners have company-wide observability and can manage company-wide settings that apply to all associated organizations. In addition, company owners have the same access as organization owners for all associated organizations.
- Organization: An organization is a collection of teams and repositories. Docker Team and Business subscribers must have at least one organization.
  - Organization owner: An organization can have multiple owners. Organization owners have observability into their organization and can manage its users and settings.
- Team: A team is a group of Docker members that belong to an organization. Organization and company owners can group members into additional teams to configure repository permissions on a per-team basis. Using teams to group members is optional.
- Member: A member is a Docker user that's a member of an organization. Organization and company owners can assign roles to members to define their permissions.

{{< grid >}}