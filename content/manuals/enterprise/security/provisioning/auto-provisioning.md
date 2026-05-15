---
title: Auto-provisioning
linkTitle: Auto-provision
description: Auto-provision users by associating members to your organization when email addresses match a verified domain.
keywords: user provisioning, just-in-time provisioning, JIT, autoprovision, Docker Admin, admin, security
weight: 30
---

Auto-provisioning automatically adds users to your organization when they sign in with email addresses that match your verified domains. You must verify a domain before enabling auto-provisioning.

> [!IMPORTANT]
>
> For domains that are part of an SSO connection, Just-in-Time (JIT) provisioning takes precedence over auto-provisioning when adding users to an organization.

### Overview

When auto-provisioning is enabled for a verified domain:

- Users who sign in to Docker with matching email addresses are automatically added to your organization.
- Auto-provisioning only adds existing Docker users to your organization, it doesn't create new accounts.
- Users experience no changes to their sign-in process.
- Company and organization owners receive email notifications when new users are added.
- You may need to [manage seats](/manuals/subscription/manage-seats.md) to accommodate new users.

### Enable auto-provisioning

Auto-provisioning is configured per domain. To enable it:

1. Sign in to [Docker Home](https://app.docker.com) and select
your company or organization.
1. Select **Admin Console**, then **Domain management**.
1. Select the **Actions menu** next to the domain you want to enable
auto-provisioning for.
1. Select **Enable auto-provisioning**.
1. Optional. If enabling auto-provisioning at the company level, select an
organization.
1. Select **Enable** to confirm.

The **Auto-provisioning** column will update to **Enabled** for the domain.

### Disable auto-provisioning

To disable auto-provisioning for a user:

1. Sign in to [Docker Home](https://app.docker.com) and select
your organization. If your organization is part of a company, select the company
and configure the domain for the organization at the company level.
1. Select **Admin Console**, then **Domain management**.
1. Select the **Actions menu** next to your domain.
1. Select **Disable auto-provisioning**.
1. Select **Disable** to confirm.

## Next steps

To choose a different method to provision users, you can set up:

- [SCIM provisioning](/manuals/enterprise/security/provisioning/scim/_index.md) for advanced user management.
- [Group mapping](/manuals/enterprise/security/provisioning/scim/group-mapping.md) to assign users to teams automatically.
