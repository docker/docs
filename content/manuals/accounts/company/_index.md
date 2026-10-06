---
title: Company overview
linkTitle: Company
weight: 20
description: Why a Docker company groups organizations, and what company owners administer.
keywords: company, multiple organizations, manage companies, Docker Home, Docker Business settings
grid:
  - title: Create a company
    description: Get started by learning how to create a company.
    icon: building-office-2
    link: /accounts/company/new-company/
  - title: Manage your company
    description: Add organizations, manage company owners, and invite members.
    icon: building-storefront
    link: /accounts/company/manage/
  - title: Configure SSO and SCIM
    description: Set up single sign-on and SCIM provisioning for your company.
    icon: key
    link: /security/authentication/single-sign-on/
  - title: Domain management
    description: Add and verify your company's domains.
    icon: check-badge
    link: /security/provisioning/domain-management/
  - title: FAQs
    description: Explore frequently asked questions about companies.
    link: /faqs/accounts/
    icon: question-mark-circle
aliases:
  - /admin/company/
  - /docker-hub/creating-companies/
---

{{< summary-bar feature_name="Company" >}}

A company is where you configure sign-in and administration for multiple
Docker organizations.

> [!TIP]
>
> Organization owners with a Docker Business
> subscription can [create a company](./new-company.md) in
> [Docker Home](https://app.docker.com/).

## Company roles

The following diagram shows the hierarchy between companies and organizations.

```mermaid
flowchart TB
  co((Company owners)) -.-> manage@{ shape: text, label: "manage" }
  manage -.-> C[Company]
  C --- O1[Organization]
  C --- O2[Organization]
```

 A company lets you:

- Administer every organization in the company through up to 10 company
  owners who don't occupy a purchased seat. Without a company, each
  organization's owners occupy a seat in that organization.
- Configure SSO and SCIM once for every organization in the company
- Verify your domains once at the company instead of in each organization.
  When you turn on auto-provisioning for a domain, you choose which
  organization new users join
- View members and invitations from every organization in one list, and
  export that list as a CSV

## Next steps

Learn how to create and manage a company in the following sections.

{{< grid >}}
