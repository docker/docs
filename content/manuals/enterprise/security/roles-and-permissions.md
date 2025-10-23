---
title: Roles and permissions
description: Control access to content, registry, and organization management with roles in your organization.
keywords: members, teams, organization, company, roles, access, docker hub, admin console, security, permissions
aliases:
- /docker-hub/roles-and-permissions/
- /security/for-admins/roles-and-permissions/
weight: 40
---

{{< summary-bar feature_name="General admin" >}}

Roles control what users can do in your organization. When you invite users, you assign them a role that determines their permissions for repositories, teams, and organization settings.

This page provides an overview of Docker roles and permissions for each role.

## Organization roles

Docker organizations have three main roles:

- Member: Non-administrative role with basic access. Members can view other organization members and pull images from repositories they have access to.
- Editor: Partial administrative access. Editors can create, edit, and delete repositories. They can also manage team permissions for repositories.
- Owner: Full administrative access. Owners can manage all organization settings, including repositories, teams, members, billing, and security features.

## Permissions by role

> [!NOTE]
>
> An owner role assigned at the company level has the same access as an owner role assigned at the organization level. For more information, see [Company overview](/admin/company/).

### Content and registry permissions

These permissions apply organization-wide, including all repositories in your organization's namespace.

| Permission                                            | Member | Editor | Owner |
| :---------------------------------------------------- | :----- | :----- | :----------------- |
| Explore images and extensions                         | ✅     | ✅     | ✅                 |
| Star, favorite, vote, and comment on content          | ✅     | ✅     | ✅                 |
| Pull images                                           | ✅     | ✅     | ✅                 |
| Create and publish an extension                       | ✅     | ✅     | ✅                 |
| Become a Verified, Official, or Open Source publisher | ❌     | ❌     | ✅                 |
| Edit and delete publisher repository logos            | ❌     | ✅     | ✅                 |
| Observe content engagement as a publisher             | ❌     | ❌     | ✅                 |
| Create public and private repositories                | ❌     | ✅     | ✅                 |
| Edit and delete repositories                          | ❌     | ✅     | ✅                 |
| Manage tags                                           | ❌     | ✅     | ✅                 |
| View repository activity                              | ❌     | ❌     | ✅                 |
| Set up Automated builds                               | ❌     | ❌     | ✅                 |
| Edit build settings                                   | ❌     | ❌     | ✅                 |
| View teams                                            | ✅     | ✅     | ✅                 |
| Assign team permissions to repositories               | ❌     | ✅     | ✅                 |

When you add members to teams, you can grant additional repository permissions
beyond their organization role:

1. Role permissions: Applied organization-wide (member or editor)
2. Team permissions: Additional permissions for specific repositories

### Organization management permissions

| Permission                                                        | Member | Editor | Owner |
| :---------------------------------------------------------------- | :----- | :----- | :----------------- |
| Create teams                                                      | ❌     | ❌     | ✅                 |
| Manage teams (including delete)                                   | ❌     | ❌     | ✅                 |
| Configure the organization's settings (including linked services) | ❌     | ❌     | ✅                 |
| Add organizations to a company                                    | ❌     | ❌     | ✅                 |
| Invite members                                                    | ❌     | ❌     | ✅                 |
| Manage members                                                    | ❌     | ❌     | ✅                 |
| Manage member roles and permissions                               | ❌     | ❌     | ✅                 |
| View member activity                                              | ❌     | ❌     | ✅                 |
| Export and reporting                                              | ❌     | ❌     | ✅                 |
| Image Access Management                                           | ❌     | ❌     | ✅                 |
| Registry Access Management                                        | ❌     | ❌     | ✅                 |
| Set up Single Sign-On (SSO) and SCIM                              | ❌     | ❌     | ✅ \*              |
| Require Docker Desktop sign-in                                    | ❌     | ❌     | ✅ \*              |
| Manage billing information (for example, billing address)                 | ❌     | ❌     | ✅                 |
| Manage payment methods (for example, credit card or invoice)              | ❌     | ❌     | ✅                 |
| View billing history                                              | ❌     | ❌     | ✅                 |
| Manage subscriptions                                              | ❌     | ❌     | ✅                 |
| Manage seats                                                      | ❌     | ❌     | ✅                 |
| Upgrade and downgrade plans                                       | ❌     | ❌     | ✅                 |

_\* If not part of a company_

### Docker Scout permissions

| Permission                                            | Member | Editor | Owner |
| :---------------------------------------------------- | :----- | :----- | :----------------- |
| View and compare analysis results                     | ✅     | ✅     | ✅                 |
| Upload analysis records                               | ✅     | ✅     | ✅                 |
| Activate and deactivate Docker Scout for a repository | ❌     | ✅     | ✅                 |
| Create environments                                   | ❌     | ❌     | ✅                 |
| Manage registry integrations                          | ❌     | ❌     | ✅                 |

### Docker Build Cloud permissions

| Permission                   | Member | Editor | Owner |
| ---------------------------- | :----- | :----- | :----------------- |
| Use a cloud builder          | ✅     | ✅     | ✅                 |
| Create and remove builders   | ✅     | ✅     | ✅                 |
| Configure builder settings   | ✅     | ✅     | ✅                 |
| Buy minutes                  | ❌     | ❌     | ✅                 |
| Manage subscription          | ❌     | ❌     | ✅                 |
