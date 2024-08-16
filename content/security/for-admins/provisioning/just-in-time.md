---
description: Learn how Just-in-Time provisioning works with your SSO connection.
keywords: user provisioning, just-in-time provisioning, JIT, autoprovision, Docker Hub, Docker Admin, admin, security
title: Just-in-Time provisioning
---

Just-in-Time (JIT) provisioning runs after every successful single sign-on (SSO) sign-in. JIT verifies that the user that signs in is a member of the organization and teams that they are assigned to in the IdP. When you [create your SSO connection](../single-sign-on/_index.md), JIT provisioning is turned on by default.

## SSO authentication with JIT provisioning enabled

After every successful SSO sign-in authentication, the JIT provisioner performs the following actions:

1. Checks if there's an existing Docker account with the email address of the user that just authenticated.

   a) If no account is found with the same email address, it creates a new Docker account using basic user attributes (email, name, and surname). The JIT provisioner generates a new username for this new account by using the email, name, and random numbers to make sure that all account usernames are unique in the platform.

   b) If an account exists for this email address, it uses this account and updates the full name of the user’s profile if needed.

2. Checks for any pending invitations to the SSO organization to auto-accept the invitation. If the invitation is specific to a group, the user is added to the invited group along with group mappings in the following step.

3. Checks if the IdP shared group mappings while authenticating the user.

   a) If the IdP provided group mappings for the user, the user gets added to the organizations and teams indicated by the group mappings.

   b) If the IdP didn't provide group mappings, it checks if the user is already a member of the organization, or if the SSO connection is for multiple organizations (only at company level) and if the user is a member of any of those organizations. If the user isn't a member, it adds the user to the default team and organization configured in the SSO connection.
   
   ![JIT provisioning enabled](../../images/jit-enabled-flow.svg)

## SSO authentication with JIT provisioning disabled

When you opt to disable JIT provisioning in your SSO connection, the following actions occur:

1. Checks if there's an existing Docker account with the email address of the user that just authenticated.

   a) If no account is found with the same email address, it creates a new Docker account using basic user attributes (email, name, and surname). Authentication with SSO generates a new username for this new account by using the email, name, and random numbers to make sure that all account usernames are unique in the platform.

   b) If an account exists for this email address, it uses this account and updates the full name of the user’s profile if needed.

2. Checks if there are any pending invitations to the SSO organization (or, SSO organizations if the SSO connection is managed at the company level) in order to auto-accept the invitation.

   a) If the user isn't already a member of the organization, or doesn't have a pending invitation to join, sign in fails and the user encounters an `Access denied` error. This blocks the user from joining the organization. They need to contact an administrator to invite them to join.

   b) If the user is a member of the organization, or has a pending invitation to join, then sign in is successful.

If you disable JIT provisioning when you create or edit your SSO connection, you can still use group mapping as long as you have also [enabled SCIM](/security/for-admins/provisioning/scim/#enable-scim-in-docker). When JIT provisioning is disabled and SCIM isn't enabled, users won't be auto-provisioned to groups. For instructions on disabling JIT provisioning, see [Manage how users are provisioned](/security/for-admins/single-sign-on/manage/#manage-how-users-are-provisioned).

![JIT provisioning disabled](../../images/jit-disabled-flow.svg)

## Disable JIT provisioning

You may want to disable JIT provisioning for reasons such as the following:

- You have multiple organizations, have SCIM enabled, and want SCIM to be the source of truth for provisioning
- You want to control and restrict usage based on your organization's security configuration, and want to use SCIM to provision access

> [!WARNING]
>
> Disabling JIT provisioning could potentially disrupt your users' workflows. Users must already be a member of the organization or have an invitation to the organization when they authenticate with SSO in order to sign in successfully. To auto-provision users with JIT disabled, you can [use SCIM](./scim.md).

See [Manage how users are provisioned](../single-sign-on/manage/_index.md#manage-how-users-are-provisioned) to learn how to disable JIT provisioning.
