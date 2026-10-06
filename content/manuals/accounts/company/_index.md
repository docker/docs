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

A company groups multiple Docker organizations for centralized configuration
and provides a single point of visibility across those organizations.
Organization owners with a Docker Business subscription can create a company
and manage it through Docker Home.

## Company structure

A company sits at the top of the hierarchy and groups multiple Docker
organizations for centralized configuration. Companies are only available
for Docker Business subscribers.

The following diagram shows how a company groups organizations, optional
teams, and members.

```mermaid
flowchart TB
  co[Company owner] -.-> C[Company]
  C --> O1[Organization]
  C --> O2[Organization]
  oo1[Organization owner] -.-> O1
  oo2[Organization owner] -.-> O2
  O1 --> T1["Team (optional)"]
  O1 --> T2["Team (optional)"]
  O2 --> T3["Team (optional)"]
  O2 --> T4["Team (optional)"]
  T1 --> M1[Member]
  T1 --> M2[Member]
  T2 --> M3[Member]
  T2 --> M4[Member]
  T3 --> M5[Member]
  T3 --> M6[Member]
  T4 --> M7[Member]
  T4 --> M8[Member]
  classDef optional stroke-dasharray: 5 5
  class T1,T2,T3,T4 optional
```

An organization sits below the company. You group teams and members there
and assign access to repositories. Every Docker Team and Business
subscriber has at least one organization.

For organization structure, including teams and members, see
[Organization accounts](/manuals/accounts/organization/_index.md).

[Upgrading to a Docker Business plan](https://www.docker.com/pricing?ref=Docs&refAction=DocsAdmin)
grants you the company owner role so you can manage multiple organizations.

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
[Manage your company](/manuals/accounts/company/manage.md#company-owners).

## Next steps

Learn how to create and manage a company in the following sections.

{{< grid >}}
