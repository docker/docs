---
title: Migrate JIT to SCIM
linkTitle: Migrate
description: Learn how to migrate from just-in-time (JIT) to SCIM.
weight: 30
---

If you already have users provisioned through Just-in-Time (JIT) and want to
enable full SCIM lifecycle management, you need to migrate them. Users
originally created by JIT cannot be automatically de-provisioned through SCIM,
even after SCIM is enabled.

## Why migrate

Organizations using JIT provisioning may encounter limitations with user
lifecycle management, particularly around de-provisioning. Migrating to SCIM
provides:

- Automatic user de-provisioning when users leave your organization. This is
  the primary benefit for large organizations that need full automation.
- Continuous synchronization of user attributes
- Centralized user management through your identity provider
- Enhanced security through automated access control

> [!IMPORTANT]
>
> Users originally created through JIT provisioning cannot be automatically
> de-provisioned by SCIM, even after SCIM is enabled. To enable full lifecycle
> management including automatic de-provisioning through your identity provider,
> you must manually remove these users so SCIM can re-create them with proper
> lifecycle management capabilities.

This migration is most critical for larger organizations that require fully
automated user de-provisioning when employees leave the company.

## Prerequisites

Before migrating, ensure you have:

- SCIM configured and tested in your organization
- A maintenance window for the migration

> [!WARNING]
>
> This migration temporarily disrupts user access. Plan to perform this
> migration during a low-usage window and communicate the timeline to affected
> users.

## Prepare for migration

### Transfer ownership

Before removing users, ensure that any repositories, teams, or organization
resources they own are transferred to another administrator or service account.
When a user is removed from the organization, any resources they own may
become inaccessible.

1. Review repositories, organization resources, and team ownership for affected
   users.
2. Transfer ownership to another administrator.

> [!WARNING]
>
> If ownership is not transferred, repositories owned by removed users may
> become inaccessible when the user is removed. Ensure all critical resources
> are transferred before proceeding.

### Verify identity provider configuration

1. Confirm all JIT-provisioned users are assigned to the Docker application in
   your identity provider.
2. Verify identity provider group to Docker Team mappings are configured and
   tested.

Users not assigned to the Docker application in your identity provider are not
re-created by SCIM after removal.

### Export user records

Export a list of JIT-provisioned users from Docker Admin Console:

1. Sign in to [Docker Home](https://app.docker.com) and select your
   organization.
2. Select **Admin Console**, then **Members**.
3. Select **Export members** to download the member list as CSV for backup and
   reference.

Keep this CSV list of JIT-provisioned users as a rollback reference if needed.

## Complete the migration

### Disable JIT provisioning

> [!IMPORTANT]
>
> Before disabling JIT, ensure SCIM is fully configured and tested in your
> organization. Do not disable JIT until you have verified SCIM is working
> correctly.

1. Sign in to [Docker Home](https://app.docker.com) and select your organization.
2. Select **Admin Console**, then **SSO and SCIM**.
3. In the SSO connections table, select the **Actions** menu for your connection.
4. Select **Disable JIT provisioning**.
5. Select **Disable** to confirm.

Disabling JIT prevents new users from being automatically added through SSO
during the migration.

### Remove JIT-origin users

> [!IMPORTANT]
>
> Users originally created through JIT provisioning cannot be automatically
> de-provisioned by SCIM, even after SCIM is enabled. To enable full lifecycle
> management including automatic de-provisioning through your identity provider,
> you must manually remove these users so SCIM can re-create them with proper
> lifecycle management capabilities.

This step is most critical for large organizations that require fully automated
user de-provisioning when employees leave the company.

1. Sign in to [Docker Home](https://app.docker.com) and select your organization.
2. Select **Admin Console**, then **Members**.
3. Identify and remove JIT-provisioned users in manageable batches.
4. Monitor for any errors during removal.

> [!TIP]
>
> To efficiently identify JIT users, compare the member list exported before
> SCIM was enabled with the current member list. Users who existed before SCIM
> was enabled were likely provisioned via JIT.

### Verify SCIM re-provisioning

After removing JIT users, SCIM automatically re-creates user accounts:

1. In your identity provider system log, confirm "create app user" events for
   Docker.
2. In Docker Admin Console, confirm users reappear with SCIM provisioning.
3. Verify users are added to the correct teams via group mapping.

### Validate user access

Perform post-migration validation:

1. Select a subset of migrated users to test sign-in and access.
2. Verify team membership matches identity provider group assignments.
3. Confirm repository access is restored.
4. Test that de-provisioning works correctly by removing a test user from your
   identity provider.

Keep audit exports and logs for compliance purposes.

## Migration results

After completing the migration:

- All users in your organization are SCIM-provisioned
- User de-provisioning works reliably through your identity provider
- No new JIT users are created
- Consistent identity lifecycle management is maintained

## Troubleshoot migration issues

If a user fails to reappear after removal:

1. Check that the user is assigned to the Docker application in your identity
   provider.
2. Verify SCIM is enabled in both Docker and your identity provider.
3. Trigger a manual SCIM sync in your identity provider.
4. Check provisioning logs in your identity provider for errors.

For more troubleshooting guidance, see
[Troubleshoot provisioning](/manuals/enterprise/security/provisioning/troubleshoot-provisioning.md).

## Next steps

- Set up [Group mapping](/manuals/enterprise/security/provisioning/scim/group-mapping.md).
- [Assign roles](/manuals/enterprise/security/roles-and-permissions/core-roles.md) to members of your org.
- [Enforce sign in](/manuals/enterprise/security/enforce-sign-in.md), if needed.
