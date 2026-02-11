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

Administrators use the [Docker Admin Console](https://app.docker.com/admin) to provision user seats, manage access tokens and SSO, and deploy Docker Desktop to their orgs. Admin Console lets you oversee and manage seats, security, and identity management from a single point of visibility.

## Set up Docker with Admin Console 

Administrators get started with Docker by accessing the Admin Console to create a company and organizations. 

- If you're a Docker Business subscriber, you have access to both company and organization features in Admin Console.
- If you're Docker Team subscriber, you only have access to organization features in Admin Console.

As an administrator, you act as an owner who can invite users with their email addresses, then assign them member roles to particular teams.

## Company and organization hierarchy

Admin Console gives administrators a bird's eye overview of a company  and its downstream organizations. Company and organizations have a hierarchical relationship:

![Diagram showing Docker’s administration hierarchy with Company at the top, followed by Organizations, Teams, and Members](./images/docker-admin-structure.webp)

Administrators can occupy either company owner or organization owner roles, each with their own permissions and seat rules. 

- Company owners can view and edit downstream organizations, or change SSO and SCIM settings. When a company owner makes a change to the company, it affects all organizations beneath them.
- Organization owners have full admin permissions to manage members, roles, and teams within their organization, but not organizations they are not the owner to.

When an administrator creates the first company from Admin Console, they assume both company and organization owner roles. If you're a Docker Team subscriber, you're the owner for that organization only and don't assume company owner permissions.    

### Company

If you're a Docker Business subscriber, then a company is the highest level of visibility an administrator can have.

Companies have the following administrator role available:

- Company owner: Can view and manage all organizations within the company.
Has full access to company-wide settings and inherits the same permissions as
organization owners.

### Organization

An organization contains teams and repositories. All Docker Team and Business
subscribers must have at least one organization.

Organizations have the following administrator role available:

- Organization owner: Can manage organization settings, users, and access
controls.

## Seats

### Team

Teams are optional and let you group members to assign repository permissions
collectively. Teams simplify permission management across projects
or functions.

### Member

A member is any Docker user added to an organization. Organization and company
owners can assign roles to members to define their level of access.

> [!NOTE]
>
> Creating a company is optional, but organizations are required for Team and
Business subscriptions.

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
