---
title: Migrate JIT to SCIM
linkTitle: Migrate
description: >-
  Move from Just-in-Time provisioning to SCIM so your identity provider
  manages Docker user lifecycle.
keywords: JIT to SCIM migration, SCIM provisioning, user deprovisioning,
  identity provider, Docker Home, user lifecycle management
weight: 30
aliases:
  - /enterprise/security/provisioning/scim/migrate-scim/
---

{{< summary-bar feature_name="SSO" >}}

Move from Just-in-Time (JIT) provisioning to System for Cross-domain Identity
Management (SCIM) as the only source of user provisioning. After SCIM is
enabled, it can manage organization members whose email domain is verified on
the SSO connection, including users created through JIT. When your identity
provider (IdP) pushes a user with a matching email address, SCIM links the
existing Docker account.

## Why migrate

With JIT turned off, your identity provider stays authoritative for who has
access:

- Users are deprovisioned when they leave your organization
- User attributes and group membership stay synchronized with the IdP

Docker recommends SCIM with JIT turned off. If your IdP supports Provision
on Demand, use it when a user needs access before the next scheduled
synchronization.

## Prerequisites

Before you migrate:

- [Set up and test SCIM](provision-scim.md) in Docker and your IdP.
- Confirm each user's email address matches exactly between the IdP and
  Docker.
- In the IdP, set the group memberships and any `dockerRole` values you want
  to keep.

## Assign users in your IdP

1. Assign every user who should belong to the Docker organization to the
   Docker application in your IdP.
1. Confirm that group-to-team mappings are configured and tested. See
   [Group mapping](group-mapping.md).

When a user isn't assigned to the Docker application, the next
synchronization deactivates the Docker account.

## Sync and verify

Trigger a synchronization, or use Provision on Demand, so SCIM links the
existing accounts.

1. In your IdP's provisioning logs, confirm that provisioning succeeded for
   those users.
1. In [Docker Home](https://app.docker.com), select your organization, then
   **Members**.
1. Confirm that the users are still members and that their roles and teams
   match the IdP.

To compare the Docker member list with the IdP, export it:

1. On the **Members** page, select **Export members**.
1. Docker emails you a link to download the CSV file.

## Disable JIT provisioning

Turn off JIT after you have verified the linked accounts. You can turn off
JIT only while SCIM is enabled.

1. Go to [Docker Home](https://app.docker.com/) and select your organization
   from the top-left account drop-down.
1. Select **Identity & auth**, then **SSO and SCIM**.
1. In the **SSO connections** table, select the **Action** icon, then select
   **Disable JIT provisioning**.
1. Select **Disable** to confirm.

With JIT turned off, users must already be members, have a pending
invitation, or be provisioned through SCIM.

## Resolve an unlinked account

Members whose email domain isn't verified on the SSO connection stay outside
SCIM. A different email address in the IdP also leaves the existing Docker
account unlinked.

1. Compare the email address in the IdP with the Docker account.
1. Confirm that the user's email domain is verified on the SSO connection.
1. Assign the user to the Docker application and run provisioning again.
1. Confirm the user under **Members**.

If the account is still unlinked, remove that user so SCIM can provision
them again.

> [!WARNING]
>
> Removing a user removes their resource ownership, such as repositories.
> Transfer ownership before you remove the user. Don't remove the only
> organization owner. Assign the Owner role to another member first.

1. In Docker Home, select **Members** and remove the user.
1. Trigger provisioning from your IdP.
1. Confirm that the user reappears with the expected role and teams.

For more troubleshooting guidance, see
[Troubleshoot provisioning](/manuals/security/provisioning/troubleshoot-provisioning.md).

## Migration results

After you turn off JIT:

- SCIM manages linked users, including users originally created through JIT
- Your IdP deprovisions users by deactivating their Docker accounts
- Sign-in no longer adds users through JIT

## Next steps

- Set up [group mapping](/manuals/security/provisioning/scim/group-mapping.md).
- [Assign roles](/manuals/security/roles-and-permissions/core-roles.md) to
  organization members.
- [Enforce sign-in](/manuals/desktop/enterprise/enforce-sign-in/_index.md) for
  your organization.
