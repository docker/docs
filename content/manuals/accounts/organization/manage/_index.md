---
title: Manage your organization
linkTitle: Manage
weight: 20
description: Learn how to manage your Docker organization, including members, teams, roles, licenses, seats, product access, activity, and settings.
keywords: manage organization, members, teams, roles, licenses, seats, product access, organization management, docker home, organization settings, activity logs, deactivate organization
grid:
  - title: User access workflow
    description: Plan roles, provisioning, and ongoing access reviews.
    icon: user-group
    link: /accounts/organization/manage/user-access/
  - title: Information
    description: Update your organization's general information and settings.
    icon: pencil-square
    link: /accounts/organization/manage/general-settings/
  - title: Members
    description: Invite, manage, and assign roles to your organization members.
    icon: user-plus
    link: /accounts/organization/manage/members/
  - title: Roles
    description: Create custom roles and review the permissions each role grants.
    icon: lock-closed
    link: /security/roles-and-permissions/
  - title: Teams
    description: Create teams and manage repository access for groups of members.
    icon: user-group
    link: /accounts/organization/manage/manage-a-team/
  - title: Licenses
    description: View your license inventory and assign licenses to teams or individual members.
    icon: key
    link: /accounts/organization/manage/manage-licenses/
  - title: Seats
    description: Add or remove seats for Docker Team and Business subscriptions.
    icon: user-circle
    link: /accounts/organization/manage/manage-seats/
  - title: Products
    description: Manage access and view usage for Docker products across your organization.
    icon: squares-2x2
    link: /accounts/organization/manage/manage-products/
  - title: Activity
    description: Review organization and repository changes in Docker Home.
    icon: clipboard-document-list
    link: /accounts/organization/activity-logs/
  - title: Docker Desktop
    description: Configure and enforce Docker Desktop settings for your organization.
    icon: cog
    link: /desktop/enterprise/hardened-desktop/settings-management/
  - title: Deactivation
    description: Deactivate an organization after completing the required steps.
    icon: minus-circle
    link: /accounts/organization/manage/deactivate-account/
  - title: Security
    description: Configure SSO and SCIM, domains, access tokens, and OIDC connections.
    icon: shield-check
    link: /security/
  - title: Billing
    description: Manage payment methods and view billing history.
    icon: credit-card
    link: /subscription-billing/
aliases:
  - /admin/organization/manage/
---

As an organization owner, you manage your organization's membership, access,
product usage, and settings. You can invite members, group them into teams,
assign roles, assign or revoke licenses and seats, change access to Docker
products, review activity, update organization information, and deactivate
the organization. Docker Home also includes Docker Desktop administration
and identity settings such as SSO and SCIM, domains, access tokens, and
OIDC connections.

## Managing your organization

You manage your organization from [Docker Home](https://app.docker.com).
Organization owners can manage the areas below. Access to a specific area
depends on your role. For example, members can open **Members** and
**Teams**.

## Seats and licenses

Seats and licenses both control access, but they apply to different kinds of
plans. The following table summarizes the difference.

| Entitlement | What it grants                                          | Applies to                                       | Managed from                 |
| ----------- | ------------------------------------------------------- | ------------------------------------------------ | ---------------------------- |
| Seat        | Membership in your Docker Team or Business subscription | Docker Team or Business subscription             | Billing                      |
| License     | Access to a specific product or add-on                  | AI Governance, Docker Offload, and other add-ons | Licenses, Teams, and Members |

For details, see [Seats](/manuals/accounts/organization/manage/manage-seats.md)
and [Manage licenses](/manuals/accounts/organization/manage/manage-licenses.md).

## Next steps

Explore the following sections to manage your organization.

{{< grid >}}
