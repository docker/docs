---
title: Troubleshoot provisioning
linkTitle: Troubleshoot
description: Troubleshoot common user provisioning issues with SCIM and Just-in-Time provisioning
keywords: SCIM troubleshooting, user provisioning, JIT provisioning, group mapping, attribute conflicts
tags: [Troubleshooting]
toc_max: 2
aliases:
    - /enterprise/troubleshoot/troubleshoot-provisioning/
    - /enterprise/security/provisioning/troubleshoot-provisioning/
---

This page helps troubleshoot common user provisioning issues including user roles, attributes, and unexpected account behavior with SCIM and Just-in-Time (JIT) provisioning.

## Full name or team membership changes after sign-in

### Error message

This scenario doesn't usually produce an error message in Docker or your IdP.
A user's full name changes after they sign in, or a team membership
disappears after a SCIM sync and comes back the next time they sign in.

### Causes

JIT and SCIM are both enabled:

- Each SSO sign-in writes the full name from the SSO assertion to the Docker
  account, replacing a name that SCIM set. The next SCIM sync can set it
  back.
- At sign-in, JIT adds the user to the teams the SSO assertion lists. SCIM
  group sync makes each mapped `organization:team` group's membership match
  the IdP group exactly. If the IdP group doesn't include the user, the next
  sync removes the team JIT added, and the next sign-in adds it back.

Roles aren't affected. JIT sets the organization role only when it first adds
the user. A later SCIM update can change that role, and the next sign-in
leaves the SCIM role in place.

### Affected environments

Docker organizations that use SCIM while JIT is still enabled.

### Steps to replicate

1. Enable SSO for your Docker organization. JIT is turned on by default.
1. Sign in through SSO with a user whose assertion includes a team.
1. Enable SCIM and synchronize groups that don't include that team.
1. The team membership from sign-in is removed. Signing in again adds it
   back.

### Solutions

#### Turn off JIT provisioning (recommended)

You can turn off JIT only while SCIM is enabled. Follow
[Disable JIT provisioning](/manuals/security/provisioning/just-in-time.md#disable-jit-provisioning).

With JIT turned off, SCIM is the source for user creation, profile updates,
and group membership.

#### Keep JIT enabled

If you keep JIT enabled:

- Match each user's email address exactly between the SSO assertion and SCIM.
- Send the same team memberships in the SSO assertion that SCIM group mapping
  synchronizes.
- Expect each sign-in to update the full name from the SSO assertion.

While JIT is on, you still assign users to the Docker application and
maintain group mappings in your IdP. Review
[how SCIM works with JIT](/manuals/security/provisioning/scim/_index.md#choose-how-scim-works-with-jit)
before you keep both enabled.

## User is deactivated after a SCIM sync

### Cause

The user isn't assigned to the Docker application in the IdP. On the next
synchronization, Docker deactivates the Docker account. Removing the user
from a mapped `organization:team` group removes that user from the team only.

### Solution

1. Assign the user to the Docker application in your IdP.
1. Match the user's email address exactly between the SSO assertion and SCIM.
1. Trigger a SCIM synchronization in your IdP.
1. Confirm that the account is active and that the user belongs to the
   expected teams.

To use one provisioning source, turn off JIT after SCIM is working. You can
turn off JIT only while SCIM is enabled. Review
[how SCIM works with JIT](/manuals/security/provisioning/scim/_index.md#choose-how-scim-works-with-jit)
before you change the configuration.

## SCIM updates don't apply to existing users

### Cause

SCIM can update any organization member whose email domain is verified on the
SSO connection, including users created through JIT or added manually. The
Docker account stays unchanged when that domain isn't verified on the
connection, or when the email address in the IdP differs from the account.

### Solution

1. Confirm that the user's email domain is verified on the SSO connection.
1. Match the email address in the IdP to the Docker account.
1. Assign the user to the Docker application and trigger provisioning.
1. In [Docker Home](https://app.docker.com), open **Members** and confirm the
   user.

If the account is still unlinked, remove that user and provision them again.

> [!WARNING]
>
> Removing a user removes their resource ownership, such as repositories.
> Transfer ownership before you remove the user.
