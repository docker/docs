---
title: Control user access
description: Control user access to your verified domains, Docker resources, and more.
weight: 20
---

In high-security environments, controlling access to Docker resources is
paramount. By verifying your organization's domains and implementing
group-based access controls, you can ensure that only authorized users can
access your Docker resources.

This module guides you through the process of verifying domains and setting up
group mappings to enforce strict access controls.

## Prerequisites

Before you begin, ensure you have:

- A Docker Business subscription
- Organization owner access to your Docker organization or company
- Access to your Domain Name System (DNS) provider to add TXT records
- Access to your Identity Provider (IdP) to configure group mappings

## Step two: Enable auto-provisioning

Auto-provisioning automatically adds users to your organization when they sign
in with an email address that matches your verified domain. This simplifies
user management and ensures consistent security settings.

To enable auto-provisioning:

1. In the [Admin Console](https://app.docker.com/admin), navigate to
the **Domain management** page and locate your verified domain.
1. Select the **Actions** menu, then **Enable auto-provisioning**.
1. Confirm the action in the pop-up modal.

> [!NOTE]
>
> Auto-provisioning is optional and does not create accounts for new users, it
adds existing unassociated users to your organization. For domains that are
using SSO, Just-in-Time (JIT) provisioning overrides auto-provisioning.

## Step three: Configure group mapping

Group mapping automates permissions management by linking identity provider
groups to Docker roles and teams. This ensures consistent access control
policies and reduces manual errors in role assignments.

1. Create groups in your IdP:
    1. Use the format `organization:team` that matches the name of your Docker
    organization and teams. For example, `docker:developers`.
    1. Assign users to the appropriate groups in your IdP.
1. Configure group mapping in Docker:
    1. In the Admin Console, navigate to
    **Security and access** > **Provisioning** > **Group mapping**.
    1. Add the group names following the `organization:team` format.
    1. Docker will automatically assign users to the corresponding teams based
    on their group membership in your IdP.

> [!NOTE]
>
> When groups are synced, Docker creates a team if it doesn’t already exist.
For detailed instructions, see [Group mapping]().

## Step four: Assign roles and permissions

Assigning appropriate roles to users ensures they have the necessary
permissions without over-provisioning access.

- Member: Non-administrative role; can view other members in the same
organization.
- Editor: Partial administrative access; can create, edit, and delete
repositories, and edit existing team’s access permissions.
- Organization owner: Full administrative access; can manage repositories,
teams, members, settings, and billing.

For more information on roles and permissions, see [Roles and permissions]().

## Best practices

- Use verified domains: Ensure all users sign in with email addresses from
your verified domains to maintain control over access.
- Implement group mapping: Automate user assignments to teams and roles to
reduce manual errors and maintain consistent access policies.
- Regularly audit access: Create a schedule to review team memberships and role
assignments to ensure they align with current organizational needs.
- Limit privileged access: Assign the Organization Owner role sparingly to
minimize the risk of unauthorized changes.
