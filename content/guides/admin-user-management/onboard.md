---
title: Onboarding and managing roles and ermissions in Docker
description: Learn how to manage roles, invite members, and implement scalable access control in Docker for secure and efficient collaboration.
keywords: sso, scim, jit, invite members, docker hub, docker admin console, onboarding, security
weight: 20
---

Learn how to invite owners, add members, and implement advanced tools like SSO and SCIM for secure and efficient access control when onboarding your organization and developers.

## Step 1: Invite owners

When you create a Docker organization, you automatically become its sole owner. While optional, adding additional owners can significantly ease the process of onboarding and managing your organization by distributing administrative responsibilities.

### Why add owners?

 - Shared responsibilities: Distribute administrative tasks, such as managing roles and permissions.

 - Improved continuity: Ensure seamless operations in case the primary owner is unavailable.

For detailed information on owners, see [Roles and permissions](/manuals/security/for-admins/roles-and-permissions.md)

## Step 2: Invite members

Members are granted controlled access to resources and enjoy enhanced organizational benefits.

### Why invite members?

 - Enhanced visibility: Gain insights into user activity, making it easier to monitor access and enforce security policies.

 - Streamlined collaboration: Help members collaborate effectively by granting access to shared resources and repositories.

 - Improved resource management: Organize and track users within your organization, ensuring optimal allocation of resources.

 - Access to enhanced features: Members benefit from organization-wide perks, such as increased pull limits and access to premium Docker features.

 - Security control: Apply and enforce security settings at an organizational level, reducing risks associated with unmanaged accounts.

For detailed information, see [Manage organization members](/manuals/admin/organization/members.md)

## Step 3: Future-proof user management

A robust, future-proof approach to user management combines automated provisioning, centralized authentication, and dynamic access control. Implementing these practices ensures a scalable, secure, and efficient environment.

### Secure user authentication with single sign-on (SSO)

Integrating Docker with your identity provider streamlines user access and enhances security.

SSO: 

 - Simplifies sign in, as users sign in with their organizational credentials.

 - Reduces password-related vulnerabilities.

 - Simplifies onboarding as it works seamlessly with SCIM and group mapping for automated provisioning.

[SSO documentation](/manuals/security/for-admins/single-sign-on/_index.md)

### Automate onboarding with SCIM and JIT provisioning

Streamline user provisioning and role management with [SCIM](/manuals/security/for-admins/provisioning/scim.md) and [Just-in-Time (JIT) provisioning](/manuals/security/for-admins/provisioning/just-in-time.md).

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

For more information on how it works, see [Group mapping](/manuals/security/for-admins/provisioning/group-mapping.md)
