---
title: Core roles and permissions
linkTitle: Core roles
description: Compare Member, Editor, and Owner permissions for content, registry, organization, Scout, and Build Cloud
keywords: core roles, member, editor, owner, permissions, organization, company, docker hub, docker scout, docker build cloud, oidc, teams, security
weight: 10
---

{{< summary-bar feature_name="General admin" >}}

Core roles are Docker's built-in roles with predefined permission sets. This
page compares what Member, Editor, and Owner can do across content, registry,
organization management, Docker Scout, and Docker Build Cloud.

## Core role summary

Docker organizations have three core roles:

- **Member**: Non-administrative role with basic access. Members can view
  other organization members and pull images from repositories they have
  access to.
- **Editor**: Partial administrative access. Editors can create, edit, and
  delete repositories. They can also manage team permissions for
  repositories.
- **Owner**: Full administrative access. Owners can manage all organization
  settings, including repositories, teams, members, billing, and security
  features.

A company owner has the same organization management permissions as an
organization owner, but some content and registry permissions don't apply
to company owners (for example, repository pull and push). For more
information, see
[Company overview](/manuals/admin/company/_index.md).

## Content and registry permissions

These permissions apply organization-wide.

| Permission                                            | Member | Editor | Owner |
| :---------------------------------------------------- | :----- | :----- | :---- |
| Explore images and extensions                         | ✅     | ✅     | ✅    |
| Star, favorite, vote, and comment on content          | ✅     | ✅     | ✅    |
| Pull images                                           | ✅     | ✅     | ✅    |
| Create and publish an extension                       | ✅     | ✅     | ✅    |
| Become a Verified, Official, or Open Source publisher | ❌     | ❌     | ✅    |
| Edit and delete publisher repository logos            | ❌     | ✅     | ✅    |
| Configure DVP analytics settings                      | ❌     | ✅     | ✅    |
| Observe content engagement as a publisher             | ❌     | ❌     | ✅    |
| Create public and private repositories                | ❌     | ✅     | ✅    |
| Disable public repositories                           | ❌     | ✅     | ✅    |
| Edit and delete repositories                          | ❌     | ✅     | ✅    |
| Manage tags                                           | ❌     | ✅     | ✅    |
| View repository activity                              | ❌     | ❌     | ✅    |
| Set up Automated builds                               | ❌     | ❌     | ✅    |
| Edit build settings                                   | ❌     | ❌     | ✅    |
| View teams                                            | ✅     | ✅     | ✅    |
| Assign team permissions to repositories               | ❌     | ✅     | ✅    |
| Manage OIDC connections                               | ❌     | ✅     | ✅    |

You can grant repository permissions beyond a member's organization role:

- Role permissions: Applied organization-wide (member or editor)
- Team permissions: Additional permissions for specific repositories

To extend access to private repositories, configure team permissions.
Custom roles can grant organization-wide permissions to manage repositories
(create, edit, delete) but don't grant pull access to private repositories.
Use team permissions for that.

## Organization management permissions

| Permission                                                        | Member | Editor | Owner |
| :---------------------------------------------------------------- | :----- | :----- | :---- |
| Create teams                                                      | ❌     | ❌     | ✅    |
| Manage teams (including delete)                                   | ❌     | ❌     | ✅    |
| Configure the organization's settings (including linked services) | ❌     | ❌     | ✅    |
| Add organizations to a company                                    | ❌     | ❌     | ✅    |
| Invite members                                                    | ❌     | ❌     | ✅    |
| Manage members                                                    | ❌     | ❌     | ✅    |
| Manage member roles and permissions                               | ❌     | ❌     | ✅    |
| View member activity                                              | ❌     | ❌     | ✅    |
| Export and reporting                                              | ❌     | ❌     | ✅    |
| Image Access Management                                           | ❌     | ❌     | ✅    |
| Registry Access Management                                        | ❌     | ❌     | ✅    |
| Namespace access control                                          | ❌     | ❌     | ✅    |
| Set up Single Sign-On (SSO) and SCIM                              | ❌     | ❌     | ✅ \* |
| Require Docker Desktop sign-in                                    | ❌     | ❌     | ✅ \* |
| Manage billing information (for example, billing address)         | ❌     | ❌     | ✅    |
| Manage payment methods (for example, credit card or invoice)      | ❌     | ❌     | ✅    |
| View billing history                                              | ❌     | ❌     | ✅    |
| Manage subscriptions                                              | ❌     | ❌     | ✅    |
| Manage seats                                                      | ❌     | ❌     | ✅    |
| Upgrade and downgrade plans                                       | ❌     | ❌     | ✅    |

> [!TIP]
>
> For more granular access control,
> [upgrade to a Docker Business plan](https://www.docker.com/pricing?ref=Docs&refAction=DocsEnterpriseCoreRoles)
> to use [custom roles](/manuals/enterprise/security/roles-and-permissions/custom-roles/_index.md).

_\* If not part of a company_

## Docker Scout permissions

| Permission                                            | Member | Editor | Owner |
| :---------------------------------------------------- | :----- | :----- | :---- |
| View and compare analysis results                     | ✅     | ✅     | ✅    |
| Upload analysis records                               | ✅     | ✅     | ✅    |
| Activate and deactivate Docker Scout for a repository | ❌     | ✅     | ✅    |
| Create environments                                   | ❌     | ❌     | ✅    |
| Manage registry integrations                          | ❌     | ❌     | ✅    |

## Docker Build Cloud permissions

| Permission                 | Member | Editor | Owner |
| -------------------------- | :----- | :----- | :---- |
| Use a cloud builder        | ✅     | ✅     | ✅    |
| Create and remove builders | ✅     | ✅     | ✅    |
| Configure builder settings | ✅     | ✅     | ✅    |
| Buy minutes                | ❌     | ❌     | ✅    |
| Manage subscription        | ❌     | ❌     | ✅    |

## Next steps

- [Custom roles](/manuals/enterprise/security/roles-and-permissions/custom-roles/_index.md):
  Create tailored permission sets on a Docker Business plan
- [Manage organization members](/manuals/admin/organization/manage/members.md):
  Invite users and assign roles
- [Company overview](/manuals/admin/company/_index.md): Understand company
  owner permissions versus organization owner permissions
