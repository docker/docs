---
title: Company administration overview
weight: 20
description: Learn how to manage multiple organizations using companies, including managing users, owners, and security.
keywords: company, multiple organizations, manage companies, admin console, Docker Business settings
grid:
- title: Create a company
  description: Get started by learning how to create a company.
  icon: apartment
  link: /admin/company/new-company/
- title: Manage organizations
  description: Learn how to add and manage organizations as well as seats within your
    company.
  icon: store
  link: /admin/company/organizations/
- title: Manage company owners
  description: Find out more about company owners and how to manage them.
  icon: supervised_user_circle
  link: /admin/company/owners/
- title: Manage users
  description: Explore how to manage users in all organizations.
  icon: group_add
  link: /admin/company/users/
- title: Configure single sign-on
  description: Discover how to configure SSO for your entire company.
  icon: key
  link: /security/for-admins/single-sign-on/
- title: Set up SCIM
  description: Set up SCIM to automatically provision and deprovision users in your
    company.
  icon: checklist
  link: /security/for-admins/provisioning/scim/
- title: Domain management
  description: Add and verify your company's domains.
  icon: domain_verification
  link: /security/for-admins/domain-management/
- title: FAQs
  description: Explore frequently asked questions about companies.
  link: /faq/admin/company-faqs/
  icon: help
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
- Have up to ten unique users assigned to the company owner role
- Configure SSO and SCIM for all nested organizations
- Enforce SSO for all users in the company

## Create and manage your company

Learn how to create and manage a company in the following sections.

{{< grid >}}
