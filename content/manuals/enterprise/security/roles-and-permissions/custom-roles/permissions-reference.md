---
title: Custom role permissions reference
linkTitle: Permissions reference
description: >-
  Permissions available for Docker custom roles across organization management,
  Docker Hub, billing, and AI Governance.
keywords: >-
  custom roles, custom role permissions, Docker, Docker Hub, organization
  management, billing, AI Governance, access tokens, SSO, SCIM, security
weight: 20
---

{{< summary-bar feature_name="General admin" >}}

Custom roles use permissions from several categories. Use the following tables
to [create or edit a custom role](manage.md).

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
| Manage custom roles               | Manage and assign custom roles                                    |
| Manage organization access tokens | Create, update, and delete organization access tokens             |
| View activity logs                | Access organization audit logs and activity history               |
| View domains                      | View domains and domain audit settings                            |
| Manage domains                    | Manage verified domains and domain audit settings                 |
| View SSO and SCIM                 | View single sign-on and user provisioning configurations          |
| Manage SSO and SCIM               | Manage single sign-on and user provisioning configurations        |
| Manage Desktop settings           | Configure Docker Desktop settings policies and view usage reports |

## Docker Hub

| Permission          | Description                                                |
| :------------------ | :--------------------------------------------------------- |
| View repositories   | View repository details and contents                       |
| Manage repositories | Create, update, and delete repositories and their contents |

## Billing

| Permission     | Description                           |
| :------------- | :------------------------------------ |
| View billing   | View organization billing information |
| Manage billing | Manage organization billing           |

## Governance

| Permission      | Description                                          |
| :-------------- | :--------------------------------------------------- |
| View policies   | View existing AI Governance policies and their rules |
| Manage policies | Manage AI Governance policies and their rules        |
