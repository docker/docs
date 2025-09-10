---
title: FAQs on companies
linkTitle: Company
weight: 30
description: Company FAQs
keywords: Docker, Docker Hub, SSO FAQs, single sign-on, company, administration, company management
tags: [FAQ]
aliases:
- /docker-hub/company-faqs/
- /faq/admin/company-faqs/
---

### Some of my organizations don’t have a Docker Business subscription. Can I still use a parent company?

Yes, but you can only add organizations with a Docker Business subscription
to a company.

### What happens if one of my organizations downgrades from Docker Business, but I still need access as a company owner?

To access and manage child organizations, the organization must have a
Docker Business subscription. If the organization isn’t included in this
subscription, the owner of the organization must manage the organization
outside of the company.

### Do company owners occupy a subscription seat?

Company owners do not occupy a seat unless one of the following is true:

- They are added as a member of an organization under your company
- SSO is enabled

Although company owners have the same access as organization owners across all
organizations in the company, it's not necessary to add them to any
organization. Doing so will cause them to occupy a seat.

When you first create a company, your account is both a company owner and an
organization owner. In that case, your account will occupy a seat as long as
you remain an organization owner.

To avoid occupying a seat, [assign another user as the organization owner](/manuals/admin/organization/members.md#update-a-member-role) and remove yourself from the organization.
You'll retain full administrative access as a company owner without using a
subscription seat.

### What permissions does the company owner have in the associated/nested organizations?

Company owners can navigate to the **Organizations** page to view all their
nested organizations in a single location. They can also view or edit organization members and change single sign-on (SSO) and System for Cross-domain Identity Management (SCIM) settings. Changes to company settings impact all users in each organization under the company.

For more information, see [Roles and permissions](/manuals/enterprise/security/roles-and-permissions.md).
