---
title: Custom roles permissions reference
linkTitle: Permissions reference
description: Reference of permissions you can assign when creating custom roles in Docker
keywords: custom roles, permissions reference, organization management, billing, docker hub, governance, access tokens, sso, scim, security
weight: 20
---

{{< summary-bar feature_name="General admin" >}}

Custom roles are built from permissions across several categories. Use the
tables on this page when you
[create or edit a custom role](manage.md).

## Organization management

| Permission                        | Description                                                       |
| :-------------------------------- | :---------------------------------------------------------------- |
| View teams                        | View teams and team members                                       |
| Manage teams                      | Create, update, and delete teams and team members                 |
| Manage registry access            | Control which registries members can access                       |
| Manage image access               | Set policies for which images members can pull and use            |
| Update organization information   | Update organization information such as name and location         |
| Member management                 | Manage organization members, invites, and roles                   |
| View custom roles                 | View existing custom roles and their permissions                  |
| Manage custom roles               | Full access to custom role management and assignment              |
| Manage organization access tokens | Create, update, and delete organization access tokens             |
| View activity logs                | Access organization audit logs and activity history               |
| View domains                      | View domains and domain audit settings                            |
| Manage domains                    | Manage verified domains and domain audit settings                 |
| View SSO and SCIM                 | View single sign-on and user provisioning configurations          |
| Manage SSO and SCIM               | Full access to SSO and SCIM management                            |
| Manage Desktop settings           | Configure Docker Desktop settings policies and view usage reports |

## Docker Hub

| Permission          | Description                                                |
| :------------------ | :--------------------------------------------------------- |
| View repositories   | View repository details and contents                       |
| Manage repositories | Create, update, and delete repositories and their contents |

## Billing

| Permission     | Description                                      |
| :------------- | :----------------------------------------------- |
| View billing   | View organization billing information            |
| Manage billing | Complete access to managing organization billing |

## Governance

| Permission      | Description                                          |
| :-------------- | :--------------------------------------------------- |
| View policies   | View existing AI Governance policies and their rules |
| Manage policies | Full access to AI Governance policy management       |
