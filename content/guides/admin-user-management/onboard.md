---
title: Onboarding and managing roles and permissions in Docker
description: Learn how to manage roles, invite members, and implement scalable access control in Docker for secure and efficient collaboration.
keywords: sso, scim, jit, invite members, docker hub, docker admin console, onboarding, security
weight: 20
---

This page guides you through onboarding owners and members, and using tools like SSO and SCIM to future-proof onboarding going forward.

## Invite owners

When you create a Docker organization, you automatically become its sole owner. While optional, adding additional owners can significantly ease the process of onboarding and managing your organization by distributing administrative responsibilities. It also ensures continuity and prevents blockers if the primary owner is unavailable.

For detailed information on owners, see [Roles and permissions](/manuals/enterprise/security/roles-and-permissions.md).

## Invite members and assign roles

Members are granted controlled access to resources and enjoy enhanced organizational benefits. When you invite members to join your Docker organization, you immediately assign them a role.

### Benefits of inviting members

- Enhanced visibility: Gain insights into user activity, making it easier to monitor access and enforce security policies.
- Streamlined collaboration: Help members collaborate effectively by granting access to shared resources and repositories.
- Improved resource management: Organize and track users within your organization, ensuring optimal allocation of resources.
- Access to enhanced features: Members benefit from organization-wide perks, such as increased pull limits and access to premium Docker features.
- Security control: Apply and enforce security settings at an organizational level, reducing risks associated with unmanaged accounts.

For detailed information, see [Manage organization members](/manuals/admin/organization/members.md).

## Future-proof user management

A robust, future-proof approach to user management combines automated provisioning, centralized authentication, and dynamic access control. Implementing these practices ensures a scalable, secure, and efficient environment.

### Secure user authentication with single sign-on (SSO)

Integrating Docker with your identity provider streamlines user access and enhances security.

SSO:

- Simplifies sign in, as users sign in with their organizational credentials.
- Reduces password-related vulnerabilities.
- Simplifies onboarding as it works seamlessly with SCIM and group mapping for automated provisioning.

For more information, see the [SSO documentation](/manuals/enterprise/security/single-sign-on/_index.md).

### Automate onboarding with SCIM and JIT provisioning

Streamline user provisioning and role management with [SCIM](/manuals/enterprise/security/provisioning/scim.md) and [Just-in-Time (JIT) provisioning](/manuals/enterprise/security/provisioning/just-in-time.md).

With SCIM you can:

- Sync users and roles automatically with your identity provider.
- Automate adding, updating, or removing users based on directory changes.

With JIT provisioning you can:

- Automatically add users upon first sign in based on [group mapping](#simplify-access-with-group-mapping).
- Reduce overhead by eliminating pre-invite steps.

### Simplify access with group mapping

Group mapping automates permissions management by linking identity provider groups to Docker roles and teams.

It also:

- Reduces manual errors in role assignments.
- Ensures consistent access control policies.
- Help you scale permissions as teams grow or change.

For more information on how it works, see [Group mapping](/manuals/enterprise/security/provisioning/group-mapping.md).
