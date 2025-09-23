---
description: Learn how Just-in-Time provisioning works with your SSO connection.
keywords: user provisioning, just-in-time provisioning, JIT, autoprovision, Docker Admin, admin, security
title: Just-in-Time provisioning
linkTitle: Just-in-Time
weight: 10
aliases:
 - /security/for-admins/provisioning/just-in-time/
---

{{< summary-bar feature_name="SSO" >}}

Just-in-Time (JIT) provisioning streamlines user onboarding by automatically creating and updating user accounts during SSO authentication. This eliminates manual account creation and ensures users have immediate access to your organization's resources. JIT verifies that users belong to the organization and assigns them to the appropriate teams based on your identity provider (IdP) configuration. When you create your SSO connection, JIT provisioning is turned on by default.

This page explains how JIT provisioning works, SSO authentication flows, and how to disable JIT provisioning.

## Prerequisites

Before you begin, you must have:

- SSO configured for your organization
- Administrator access to Docker Home and your identity provider

## SSO authentication with JIT provisioning enabled

When a user signs in with SSO and you have JIT provisioning enabled, the following steps occur automatically:

1. The system checks if a Docker account exists for the user's email address.

    - If an account exists: The system uses the existing account and updates the user's full name if necessary.
    - If no account exists: A new Docker account is created using basic user attributes (email, name, and surname). A unique username is generated based on the user's email, name, and random numbers to ensure all usernames are unique across the platform.

2. The system checks for any pending invitations to the SSO organization.

    - Invitation found: The invitation is automatically accepted.
    - Invitation includes a specific group: The user is added to that group within the SSO organization.

3. The system verifies if the IdP has shared group mappings during authentication.

    - Group mappings provided: The user is assigned to the relevant organizations and teams.
    - No group mappings provided: The system checks if the user is already part of the organization. If not, the user is added to the default organization and team configured in the SSO connection.

The following graphic provides an overview of SSO authentication with JIT enabled:

   ![JIT provisioning enabled workflow](../images/jit-enabled-flow.svg)

## SSO authentication with JIT provisioning disabled

When JIT provisioning is disabled, the following actions occur during SSO authentication:

1. The system checks if a Docker account exists for the user's email address.

    - If an account exists: The system uses the existing account and updates the user's full name if necessary.
    - If no account exists: A new Docker account is created using basic user attributes (email, name, and surname). A unique username is generated based on the user's email, name, and random numbers to ensure all usernames are unique across the platform.

2. The system checks for any pending invitations to the SSO organization.

   - Invitation found: If the user is a member of the organization or has a pending invitation, sign-in is successful, and the invitation is automatically accepted.
   - No invitation found: If the user is not a member of the organization and has no pending invitation, the sign-in fails, and an `Access denied` error appears. The user must contact an administrator to be invited to the organization.

With JIT disabled, group mapping is only available if you have [SCIM enabled](scim/#enable-scim-in-docker). If SCIM is not enabled, users won't be auto-provisioned to groups.

The following graphic provides an overview of SSO authentication with JIT disabled:

![JIT provisioning disabled workflow](../images/jit-disabled-flow.svg)

## Disable JIT provisioning

> [!WARNING]
>
> Disabling JIT provisioning may disrupt your users' access and workflows. With JIT disabled, users will not be automatically added to your organization. Users must already be a member of the organization or have a pending invitation to successfully sign in through SSO. To auto-provision users with JIT disabled, [use SCIM](./scim.md).

You may want to disable JIT provisioning for reasons such as the following:

- You have multiple organizations, have SCIM enabled, and want SCIM to be the source of truth for provisioning
- You want to control and restrict usage based on your organization's security configuration, and want to use SCIM to provision access

Users are provisioned with JIT by default. If you enable SCIM, you can disable JIT:

1. Go to [Docker Home](https://app.docker.com/) and select your organization from the top-left account drop-down.
1. Select **Admin Console**, then **SSO and SCIM**.
1. In the **SSO connections** table, select the **Action** icon, then select **Disable JIT provisioning**.
1. Select **Disable** to confirm.

## Next steps

- Configure [SCIM provisioning](/manuals/enterprise/security/provisioning/scim.md) for advanced user management.
- Set up [group mapping](/manuals/enterprise/security/provisioning/group-mapping.md) to automatically assign users to teams.
- Review [Troubleshoot provisioning](/manuals/enterprise/troubleshoot/troubleshoot-provisioning.md).
