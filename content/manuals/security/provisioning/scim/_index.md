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
  - /enterprise/security/provisioning/scim/
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

When a user is no longer assigned to the Docker application in the IdP,
SCIM deactivates the Docker account.

SCIM automates:

- Creating users
- Updating user profiles
- Deactivating users
- Reactivating users
- Synchronizing groups when group mapping is configured

> [!NOTE]
>
> After you enable SCIM, it can manage and deprovision any organization member
> whose email domain is verified on the SSO connection. That includes users
> created through JIT or added manually. When your IdP pushes a user with a
> matching email address, SCIM links the existing Docker account. Members
> whose email domain isn't verified on the connection stay outside SCIM.

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

JIT and SCIM run independently. While JIT is on, you still assign users to
the Docker application and maintain group mappings in your IdP.

Two values can change back and forth when both are enabled:

- Full name. Each SSO sign-in writes the name from the SSO assertion to the
  Docker account. If SCIM set a different name, for example the IdP profile
  says "Jon Smith" but the assertion sends "Jonathan Smith", sign-in replaces
  the SCIM value. The next SCIM sync can set it back.
- Team membership. At sign-in, JIT reads the `groups` or `dockerTeam` value
  from the SSO assertion and adds the user to those teams. It never removes
  teams. SCIM group sync makes each mapped group's membership match the IdP
  group exactly. If JIT added a user to a team that the IdP group doesn't
  include, the next sync of that group removes the user from the team, and
  the next sign-in adds them back.

Roles don't move back and forth. JIT sets the organization role only when it
first adds the user. A later SCIM update can change that role, and the next
sign-in leaves the SCIM role in place.

When a user isn't assigned to the Docker application in the IdP, the next
synchronization deactivates the Docker account. Removing a user from a mapped
group removes that user from the team only.

If you keep both enabled:

- Match each user's email address exactly between the SSO assertion and SCIM.
- Assign every user who can sign in through SSO to the Docker application in
  the IdP.
- Keep the IdP authoritative for roles, organizations, teams, and group
  membership.
- After sign-in and after each SCIM synchronization, confirm that full names
  and team memberships still match the IdP.

SCIM links an existing account, including one created through JIT or added
manually, when the IdP pushes a user with a matching email address. To use
SCIM as the only provisioning source, see
[Migrate JIT to SCIM](/manuals/security/provisioning/scim/migrate-scim.md).

## Next steps

- [Set up SCIM provisioning](/manuals/security/provisioning/scim/provision-scim.md)
  to enable SCIM in Docker and your identity provider.
- [Migrate JIT to SCIM](/manuals/security/provisioning/scim/migrate-scim.md)
  to turn off JIT after SCIM is managing your users.
- [Group mapping](/manuals/security/provisioning/scim/group-mapping.md) to
  sync identity provider groups with Docker teams.
- [Troubleshoot provisioning](/manuals/security/provisioning/troubleshoot-provisioning.md)
  for SCIM, JIT, and attribute issues.
