---
title: Company overview
linkTitle: Company
weight: 20
description: Learn how to manage multiple organizations using companies, including managing users, owners, and security.
keywords: company, multiple organizations, manage companies, Docker Home, Docker Business settings
grid:
  - title: Create a company
    description: Get started by learning how to create a company.
    icon: building-office-2
    link: /admin/company/new-company/
  - title: Manage your company
    description: Add organizations, manage company owners, and invite members.
    icon: building-storefront
    link: /admin/company/manage/
  - title: Configure SSO and SCIM
    description: Set up single sign-on and SCIM provisioning for your company.
    icon: key
    link: /enterprise/security/single-sign-on/
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

A company provides a single point of visibility across multiple
organizations, for centralized organization and settings management.
Organization owners with a Docker Business subscription can create a company
and manage it through Docker Home.

## Company structure

The following diagram shows how a company relates to its associated
organizations.

![Diagram showing how companies relate to Docker organizations](/admin/images/docker-admin-structure.webp)

For the full administration hierarchy, see the
[administration overview](/manuals/admin/_index.md#company-and-organization-hierarchy).

## Company roles

A company includes one or more company owners. The creator of a company
becomes both a company owner and an organization owner, and occupies a seat
as organization owner. After creation, a company can have multiple owners,
and each owner has visibility across the entire company. They can manage
settings for every organization under it and have the same access rights as
organization owners.

- A company can have up to ten unique company owners.
- Company owners don't occupy a seat unless one of the following applies:
  - They're added as a member of an organization under the company.
  - SSO is enabled and the company owner signs in through SSO, which
    automatically adds them as an organization member.

To add or remove company owners, see
[Manage your company](/manuals/admin/company/manage.md#company-owners).

## What's next

Learn how to create and manage a company in the following sections.

{{< grid >}}
