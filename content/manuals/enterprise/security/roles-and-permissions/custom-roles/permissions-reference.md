---
title: Custom role permissions reference
linkTitle: Permissions reference
description: >-
  Permissions available for Docker custom roles across organization management,
  Docker Hub, billing, AI Governance, Docker Hardened Images, and Docker
  Offload.
keywords: >-
  custom roles, custom role permissions, Docker, Docker Hub, organization
  management, billing, AI Governance, access tokens, SSO, SCIM, OIDC, DHI,
  Docker Offload, security
weight: 20
---

{{< summary-bar feature_name="Custom roles" >}}

Custom roles use permissions from organization management, Docker Hub,
billing, AI Governance, Docker Hardened Images, and Docker Offload. Use
the following tables to [create or edit a custom role](manage.md).

## Organization management

| Permission                        | Description                                                                                     |
| :-------------------------------- | :---------------------------------------------------------------------------------------------- |
| View teams                        | View teams and team members                                                                     |
| Manage teams                      | Create, update, and delete teams and team members                                               |
| Manage registry access            | Control which registries members can access                                                     |
| Manage image access               | Set policies for which images members can pull and use                                          |
| Update organization information   | Update organization information such as name and location                                       |
| Member management                 | Manage organization members, invites, and roles                                                 |
| View custom roles                 | View existing custom roles and their permissions                                                |
| Manage custom roles               | Full access to custom role management and assignment                                            |
| Manage organization access tokens | Create, update, and delete repositories in this org. Push/pull or registry actions not included |
| View activity logs                | Access organization audit logs and activity history                                             |
| View usage reports                | Download organization usage reports (pulls, storage)                                            |
| View domains                      | View domains and domain audit settings                                                          |
| Manage domains                    | Manage verified domains and domain audit settings                                               |
| View SSO and SCIM                 | View single sign-on and user provisioning configurations                                        |
| Manage SSO and SCIM               | Full access to SSO and SCIM management                                                          |
| Manage Desktop settings           | Configure Docker Desktop settings policies and view usage reports                               |
| View OIDC connections             | View OIDC connections and their configuration                                                   |
| Manage OIDC connections           | View, create, edit, and delete OIDC connections                                                 |

## Docker Hub

| Permission                         | Description                                                                                  |
| :--------------------------------- | :------------------------------------------------------------------------------------------- |
| View repositories                  | View repository details and contents                                                         |
| Manage repositories                | Full repository management including settings, webhooks, privacy, Dockerfile, immutable tags |
| Manage repository team permissions | Add and remove teams from a repository, and manage team access level                         |

## Billing

| Permission     | Description                                      |
| :------------- | :----------------------------------------------- |
| View billing   | View organization billing information            |
| Manage billing | Complete access to managing organization billing |

## AI Governance

| Permission            | Description                                          |
| :-------------------- | :--------------------------------------------------- |
| View policies         | View existing AI Governance policies and their rules |
| Manage policies       | Full access to AI Governance policy management       |
| View audit logs       | View audit events for the organization               |
| Manage audit settings | Update audit configuration for the organization      |

## DHI (Docker Hardened Images)

| Permission         | Description                                      |
| :----------------- | :----------------------------------------------- |
| Create DHI mirrors | Create Docker Hardened Image mirror repositories |

## Docker Offload

| Permission        | Description                                    |
| :---------------- | :--------------------------------------------- |
| Offload Read-Only | View Offload account status, leases, and zones |

## Next steps

- [Manage custom roles](manage.md): Create, assign, and delete custom
  roles
- [Core roles and permissions](/manuals/enterprise/security/roles-and-permissions/core-roles.md):
  Compare built-in Member, Editor, and Owner permissions
