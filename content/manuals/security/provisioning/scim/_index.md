---
title: SCIM provisioning overview
linkTitle: SCIM
weight: 10
description: >-
  Provision, update, and deprovision Docker users with SCIM, and choose how
  SCIM works with Just-in-Time provisioning.
keywords: SCIM, SSO, user provisioning, deprovisioning, JIT, role mapping,
  group mapping, identity provider, Provision on Demand, Okta, Entra ID
aliases:
  - /security/for-admins/scim/
  - /security/for-admins/provisioning/scim/
  - /platform/security/provisioning/scim/
---

{{< summary-bar feature_name="SSO" >}}

System for Cross-domain Identity Management (SCIM) synchronizes users and
groups between your identity provider (IdP) and Docker. It provisions
accounts, syncs profile updates, and deprovisions users throughout the
account lifecycle.

## Prerequisites

Before you begin, you must have:

- SSO configured for your organization
- Administrator access to Docker Home and your identity provider

## How SCIM works

After you enable SCIM, any user assigned to your Docker application in the
identity provider is provisioned and added to your Docker organization. SCIM
syncs profile updates from the identity provider, such as name changes, and
reactivates users who are reassigned to the application. If group mapping is
configured, SCIM also synchronizes groups.

When a user is removed from the Docker application, SCIM deactivates and
removes them from your Docker organization.

SCIM automates:

- Creating users
- Updating user profiles
- Removing and deactivating users
- Reactivating users
- Synchronizing groups when group mapping is configured

> [!NOTE]
>
> Enabling SCIM doesn't convert manually added users into SCIM-managed users.
> SCIM only provides full lifecycle management for users it provisions.

## Choose how SCIM works with JIT

Docker turns on Just-in-Time (JIT) provisioning when you configure an SSO
connection. Before you enable SCIM, choose whether SCIM or JIT will own user
provisioning.

Docker recommends using one provisioning source. Using SCIM without JIT keeps
the IdP directory authoritative for user creation, attributes, group
membership, and deprovisioning.

### Use SCIM without JIT

With JIT turned off, SCIM provisions users on the IdP's synchronization
schedule instead of when users sign in. This configuration provides continuous
attribute updates and automatic deprovisioning.

If your IdP supports Provision on Demand, you can trigger an immediate sync
for a user who needs access before the next scheduled synchronization.
Configure and test SCIM before you
[turn off JIT](/manuals/security/provisioning/just-in-time.md#disable-jit-provisioning).

### Use SCIM with JIT

JIT and SCIM run independently:

- JIT reads the SSO assertion and applies its values when a user signs in.
- SCIM reads users, attributes, and group membership from the IdP on its
  synchronization schedule.

When both are enabled, values applied during sign-in can overwrite values that
SCIM set. A JIT-provisioned user who isn't in the SCIM-mapped IdP group can
also be removed from the Docker organization during the next SCIM
synchronization.

If you keep both enabled:

- Match each user's email address exactly between the SSO assertion and SCIM.
- Add every user who can be provisioned through JIT to the SCIM-mapped group.
- Keep roles, organizations, teams, and group membership consistent in the
  IdP.
- Monitor users and assignments for changes after sign-in and SCIM
  synchronization.

Keeping a JIT-provisioned user in the mapped group doesn't convert the account
to SCIM lifecycle management. To let SCIM manage the account, follow
[Migrate JIT to SCIM](/manuals/security/provisioning/scim/migrate-scim.md).

## Next steps

- [Set up SCIM provisioning](/manuals/security/provisioning/scim/provision-scim.md)
  to enable SCIM in Docker and your identity provider.
- [Migrate JIT to SCIM](/manuals/security/provisioning/scim/migrate-scim.md)
  if users were provisioned with Just-in-Time (JIT) before you enabled SCIM.
- [Group mapping](/manuals/security/provisioning/scim/group-mapping.md) to
  sync identity provider groups with Docker teams.
- [Troubleshoot provisioning](/manuals/security/provisioning/troubleshoot-provisioning.md)
  for SCIM, JIT, and attribute issues.
