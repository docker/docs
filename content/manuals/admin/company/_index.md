---
title: Company overview
linkTitle: Company
weight: 20
description: Learn how to manage multiple organizations using companies, including managing users, owners, and security.
keywords: company, multiple organizations, manage companies, admin console, Docker Business settings
grid:
  - title: Create a company
    description: Get started by learning how to create a company.
    icon: building-office-2
    link: /admin/company/new-company/
  - title: Manage organizations
    description:
      Learn how to add and manage organizations as well as seats within your
      company.
    icon: building-storefront
    link: /admin/company/manage/#manage-organizations
  - title: Manage company owners
    description: Find out more about company owners and how to manage them.
    icon: user-group
    link: /admin/company/manage/#manage-company-owners
  - title: Manage users
    description: Explore how to manage users in all organizations.
    icon: user-plus
    link: /admin/company/manage/#manage-company-members
  - title: Configure single sign-on
    description: Discover how to configure SSO for your entire company.
    icon: key
    link: /enterprise/security/single-sign-on/
  - title: Set up SCIM
    description:
      Set up SCIM to automatically provision and deprovision users in your
      company.
    icon: clipboard-document-check
    link: /enterprise/security/provisioning/scim/
  - title: Domain management
    description: Add and verify your company's domains.
    icon: check-badge
    link: /enterprise/security/domain-management/
  - title: FAQs
    description: Explore frequently asked questions about companies.
    link: /faq/admin/company-faqs/
    icon: question-mark-circle
aliases:
  - /docker-hub/creating-companies/
---

{{< summary-bar feature_name="Company" >}}

A company provides a single point of visibility across multiple organizations,
simplifying organization and settings management.

Organization owners with a Docker Business subscription can create a company
and manage it through the [Docker Admin Console](https://app.docker.com/admin).

The following diagram shows how a company relates to its associated
organizations.

![Diagram showing how companies relate to Docker organizations](/admin/images/docker-admin-structure.webp)

## Key features

With a company, administrators can:

- View and manage all nested organizations
- Configure company and organization settings centrally
- Control access to the company
- Configure SSO and SCIM for all nested organizations
- Enforce SSO for all users in the company

## Company owners

A company can have multiple owners. A company owner has visibility across the
entire company and can manage settings for every organization under it. Company
owners have the same access rights as organization owners, but don't need to be
a member of any individual organization. A company can have up to ten unique
company owners.

Company owners don't occupy a seat unless one of the following applies:

- They're added as a member of an organization under the company.
- SSO is enabled and the company owner signs in through SSO, which automatically
  adds them as an organization member.

When you first create a company, your account is both a company owner and an
organization owner, so it occupies a seat while you remain an organization
owner. To keep full company-owner access without using a seat,
[assign another user as the organization owner](/manuals/admin/organization/manage/members.md#update-a-member-role),
then remove yourself from the organization.

To add or remove company owners, see
[Manage your company](/manuals/admin/company/manage.md#company-owners).

## Create and manage your company

Learn how to create and manage a company in the following sections.

{{< grid >}}
