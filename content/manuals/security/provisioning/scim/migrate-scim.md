---
title: Migrate JIT to SCIM
linkTitle: Migrate
description: >-
  Migrate JIT-provisioned Docker users to SCIM for automated user lifecycle
  management and deprovisioning.
keywords: JIT to SCIM migration, SCIM provisioning, user deprovisioning,
  identity provider, Docker Home, user lifecycle management
weight: 30
aliases:
  - /platform/security/provisioning/scim/migrate-scim/
---

{{< summary-bar feature_name="SSO" >}}

Migrate users created through Just-in-Time (JIT) provisioning so System for
Cross-domain Identity Management (SCIM) can manage their full account
lifecycle. Enabling SCIM doesn't convert existing JIT-provisioned users into
SCIM-managed users.

## Why migrate

Migrating users from JIT to SCIM provides:

- Automatic user deprovisioning when users leave your organization
- Continuous synchronization of user attributes
- Centralized user management through your identity provider
- Automated access removal

> [!IMPORTANT]
>
> SCIM can't deprovision users originally created through JIT. You must remove
> these users from the Docker organization so SCIM can provision them again as
> SCIM-managed users.

## Prerequisites

Before migrating, you must have:

- SCIM configured and tested in your organization
- A maintenance window for the migration

> [!WARNING]
>
> This migration temporarily disrupts user access. Plan to perform this
> migration during a low-usage window and communicate the timeline to affected
> users.

## Prepare for migration

### Review roles and access

Removing a member revokes their access to the organization's resources and
teams. Record each user's role and team memberships so you can verify access
after SCIM provisions the user again.

1. Review the roles and team memberships of affected users.
1. Confirm that the organization has an owner who isn't part of the migration.

> [!WARNING]
>
> Don't remove the only organization owner. Assign the Owner role to another
> member before you begin the migration.

### Verify identity provider configuration

1. Confirm that all JIT-provisioned users are assigned to the Docker
   application in your identity provider.
1. Verify that identity provider group-to-Docker-team mappings are configured
   and tested.

Users not assigned to the Docker application in your identity provider are not
re-created by SCIM after removal.

### Export user records

Export a list of JIT-provisioned users from Docker Home:

1. Sign in to [Docker Home](https://app.docker.com) and select your
   organization.
1. Select **Members**.
1. Select the **Download** icon to start the export.
1. Open the email from Docker and use the link to download the CSV file.

Keep the CSV as a record of the users included in the migration.

## Complete the migration

### Disable JIT provisioning

> [!IMPORTANT]
>
> Before disabling JIT, ensure SCIM is fully configured and tested in your
> organization. Do not disable JIT until you have verified SCIM is working
> correctly.

1. Sign in to [Docker Home](https://app.docker.com) and select your
   organization.
1. Select **Identity & auth**, then **SSO and SCIM**.
1. In the **SSO connections** table, select the **Actions** menu for your
   connection.
1. Select **Disable JIT provisioning**.
1. Select **Disable** to confirm.

Disabling JIT prevents new users from being automatically added through SSO
during the migration.

### Remove JIT-origin users

> [!IMPORTANT]
>
> Removing users temporarily interrupts their access. Confirm that SCIM is
> working before you remove them.

1. Sign in to [Docker Home](https://app.docker.com) and select your
   organization.
1. Select **Members**.
1. Identify and remove JIT-provisioned users in manageable batches.
1. Monitor for errors during removal.

> [!TIP]
>
> Use the member export, IdP assignments, and provisioning logs to identify
> JIT-provisioned users. Don't remove users based only on when they joined the
> organization.

### Verify SCIM re-provisioning

After removing JIT-provisioned users, trigger a synchronization in your IdP,
then verify that SCIM provisions the users again:

1. In your identity provider's provisioning logs, confirm successful user
   creation events for Docker.
1. In Docker Home under **Members**, confirm that users reappear.
1. Verify that group mapping adds users to the correct teams.

### Validate user access

Perform post-migration validation:

1. Select a subset of migrated users to test sign-in and access.
1. Verify that team membership matches identity provider group assignments.
1. Confirm that repository access is restored.
1. Test deprovisioning by removing a test user from your
   identity provider.

Keep audit exports and logs for compliance purposes.

## Migration results

After completing the migration:

- Migrated users are SCIM-provisioned
- User deprovisioning works through your identity provider
- No new JIT users are created

## Troubleshoot migration issues

If a user fails to reappear after removal:

1. Check that the user is assigned to the Docker application in your identity
   provider.
1. Verify SCIM is enabled in both Docker and your identity provider.
1. Trigger a manual SCIM sync in your identity provider.
1. Check provisioning logs in your identity provider for errors.

For more troubleshooting guidance, see
[Troubleshoot provisioning](/manuals/security/provisioning/troubleshoot-provisioning.md).

## Next steps

- Set up [group mapping](/manuals/security/provisioning/scim/group-mapping.md).
- [Assign roles](/manuals/security/roles-and-permissions/core-roles.md) to
  organization members.
- [Enforce sign-in](/manuals/enterprise/security/enforce-sign-in.md) for your
  organization.
