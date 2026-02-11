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

Administrators start with the [Docker Admin Console](https://app.docker.com/admin) to provision user seats, manage access tokens, SSO and SCIM, and deploy Docker Desktop to their organizations. 

## Set up Docker with Admin Console 

Administrators get started with Docker by accessing the Admin Console to create a company and organizations. 

- If you're a Docker Business subscriber, you have access to both company and organization features.
- If you're Docker Team subscriber, you have access to organization features in Admin Console.

As an administrator, you're an owner who can invite users with their email addresses, then assign them member roles to particular teams.

## Company and organization hierarchy

Admin Console gives administrators a bird's eye overview of a company  and its downstream organizations. Company and organizations have a hierarchical relationship:

![Diagram showing Docker’s administration hierarchy with Company at the top, followed by Organizations, Teams, and Members](./images/docker-admin-structure.webp)

Administrators can occupy company owner or organization owner roles (or both), where each role has its own permissions and seat rules. 

- Company owners can view and bulk edit settings and configurations for all organizations beneath them.
- Organization owners have full admin permissions to manage settings, members, roles, and teams within their organization, but not organizations they're not part of.

When an administrator creates the first company from Admin Console, they assume owner roles pursuant to their subscription type. For example:

- A Docker Business subscriber assumes owner permissions for both the first company and first organization.
- A Docker Team subscriber assumes owner permissions for the first created organization. 

### Company

The highest level of visibility an administrator can have is at the company level. A company owner views and manages all organizations within the company and has full access to company-wide settings. 

Company owners won't occupy a seat unless one of the following is true:

- They are added as a member of an organization under your company.
- SSO is enabled.

If you're a Docker team subscribe who wants access to company-level permissions, you can [upgrade to Docker Business](/subscription/change/#upgrade-your-subscription).

### Organization

An organization contains teams and repositories. All Docker Team and Business
subscribers must create one organization before inviting new members to Docker.

Organization owners manage organization settings, users, and access controls. All organizations owners occupy at least one seat, but can occupy more than one seat if they're members or owners of multiple, separate organizations. 

## Seats and user management

The number of seats an administrator can provision depends on their [subscription type](https://www.docker.com/pricing/). Once you've decided on a plan and created your first company or organization, you can send invitations to future members. 

### Seats 

A seat is a unit purchased with a subscription plan that extends access to users to an organization's repo.

  - They give administrators granular permissions around who can contribute to a repository.
  - They prevent unauthorized users from pushing to a repos they're not members of. 

For example, an organization owner takes up one seat. They can invite Docker users to an organization. Once invitees become members, organization owners can set permissions in bulk or on an individual basis to repositories affiliated with an organization.  

### Users and members 

Docker uses specific terminology to define the kind of access a Docker user has: 

- A user is someone with a Docker ID.
  - They are not necessarily affiliated with an organization.
  - They do not take up seats by default.
- An invitee is a user invited to an organization.
  - Invitees occupy one seat. 
  - This is a user state before accepting and joining an organization. 
- A member is a user who accepted an invitation to an organization.
- Teams let you group members together.
  - They are optional. 
  - They allow you to assign repository permissions in bulk. 
  - Teams can simplify permission management across projects
or functions.

## Manage companies and organizations

Learn how to manage companies and organizations in the following sections.

{{< grid >}}
